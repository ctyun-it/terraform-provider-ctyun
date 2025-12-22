package pgsql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
	"time"
)

func TestAccCtyunPostgresqlAccount(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_account." + rnd
	resourceFile := "resource_ctyun_postgresql_account.tf"
	lockResourceFile := "resource_ctyun_postgresql_account_lock.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	instanceID := dependence.pgsqlID
	testDB1 := "test1"
	testDB2 := "test"

	// 测试数据
	accountName := "test_account_" + rnd
	initialPassword := "Te123!" + utils.GenerateRandomString()
	updatedPassword := "updatED123!" + utils.GenerateRandomString()
	initialDescription := "Initial_account_description"
	updatedDescription := "Updated_account_description"

	// 等待函数
	wait10Seconds := func() {
		t.Logf("等待10秒...")
		time.Sleep(10 * time.Second)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建账户测试（普通账号）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					projectID,
					instanceID, accountName,
					initialPassword, "normal", "[]",
					initialDescription,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceID),
					resource.TestCheckResourceAttr(resourceName, "name", accountName),
					resource.TestCheckResourceAttr(resourceName, "user_type", "normal"),
					resource.TestCheckResourceAttr(resourceName, "description", initialDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "rol_can_login"),
				),
			},
			// 2. 更新账户测试（添加数据库授权）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					instanceID, accountName,
					initialPassword, "normal",
					fmt.Sprintf(`[
						{ grant_schema = "%s", privilege = "readwrite" },
						{ grant_schema = "%s", privilege = "readonly" }
					]`, testDB1, testDB2),
					initialDescription,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.#", "2"),
				),
			},
			// 3. 更新账户测试（修改密码、描述和权限）
			{
				SkipFunc: func() (bool, error) {
					return testDB1 == "", nil
				},
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					projectID,
					instanceID, accountName,
					updatedPassword, "normal",
					fmt.Sprintf(`[
						{ grant_schema = "%s", privilege = "readonly" }
					]`, testDB1),
					updatedDescription,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.0.grant_schema", testDB1),
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.0.privilege", "readonly"),
				),
			},
			// 4. 更新账户测试（锁定账户）
			{
				Config: utils.LoadTestCase(
					lockResourceFile, rnd,
					projectID,
					instanceID, accountName,
					updatedPassword, "normal", "[]",
					updatedDescription, true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "is_lock", "true"),
				),
			},
			// 5. 导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s,%s,%s",
						rs.Primary.Attributes["name"],
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "is_lock", "user_type", "description"}, // 密码敏感，导入时不验证
				PreConfig: func() {
					wait10Seconds()
				},
			},
			// 5. 导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["name"],
						rs.Primary.Attributes["instance_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "is_lock", "user_type", "description", "project_id"}, // 密码敏感，导入时不验证
				PreConfig: func() {
					wait10Seconds()
				},
			},
			// 6. 清理资源
			{
				Config: utils.LoadTestCase(
					lockResourceFile, rnd,
					projectID,
					instanceID, accountName,
					updatedPassword, "normal", "[]",
					updatedDescription, true,
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunPostgresqlAdvancedAccount(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_account." + rnd
	resourceFile := "resource_ctyun_postgresql_account_lock.tf"

	datasourceName := "data.ctyun_postgresql_accounts." + dnd
	datasourceFile := "datasource_ctyun_postgresql_accounts.tf"

	projectID := "0"
	instanceID := dependence.pgsqlID
	testDB := "test"

	// 测试数据
	accountName := "admin_" + rnd
	password := "Ad123!" + utils.GenerateRandomString()
	description := "Advanced_account_for_administration"

	// 等待函数
	wait10Seconds := func() {
		t.Logf("等待10秒...")
		time.Sleep(10 * time.Second)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建高权限账户测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					projectID,
					instanceID, accountName,
					password, "advanced", "[]", // 无数据库授权
					description, false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceID),
					resource.TestCheckResourceAttr(resourceName, "name", accountName),
					resource.TestCheckResourceAttr(resourceName, "user_type", "advanced"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "is_lock", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					// 验证高权限账户特有属性
					resource.TestCheckResourceAttr(resourceName, "rol_create_role", "true"),
					resource.TestCheckResourceAttr(resourceName, "rol_create_db", "true"),
				),
				PreConfig: func() {
					wait10Seconds()
				},
			},
			// datasource 验证
			{
				Config: utils.LoadTestCase(
					datasourceFile, dnd, instanceID, accountName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "accounts.0.name", accountName),
				),
			},
			// 2. 验证高权限账户的数据库授权能力
			{
				SkipFunc: func() (bool, error) {
					return testDB == "", nil
				},
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					projectID,
					instanceID, accountName,
					password, "advanced", fmt.Sprintf(`[
						{ grant_schema = "%s", privilege = "readwrite" }
					]`, testDB),
					description, false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.0.grant_schema", testDB),
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.0.privilege", "readwrite"),
				),
				PreConfig: func() {
					wait10Seconds()
				},
			},
			// 3. 验证高权限账户的锁定功能
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					instanceID, accountName,
					password, "advanced", "[]", // 无数据库授权
					description, true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "is_lock", "true"),
					resource.TestCheckResourceAttr(resourceName, "rol_can_login", "false"),
				),
				PreConfig: func() {
					wait10Seconds()
				},
			},

			// 5. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					instanceID, accountName,
					password, "advanced", "[]", // 无数据库授权
					description, true,
				),
				Destroy: true,
				PreConfig: func() {
					wait10Seconds()
				},
			},
		},
	})
}
