package mysql_test

import (
	"encoding/json"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

// 测试MySQL账户资源
func TestAccCtyunMysqlAccount(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_account." + rnd
	resourceFile := "resource_ctyun_mysql_account.tf"

	dataSourceName := "data.ctyun_mysql_accounts." + dnd
	datasourceFile := "datasource_ctyun_mysql_accounts.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	mysqlInstanceID := dependence.mysqlID
	accountPassword := utils.GenerateRandomString() + "&R3?=@"
	newPassword := utils.GenerateRandomString() + "New2@"
	accountName := "test_account_" + rnd

	// 测试数据库名称（确保在MySQL实例中存在）
	testDB1 := "test_db1"
	testDB2 := "test_db2"
	testDB3 := "test_db3"

	// 初始权限配置
	initialPrivileges := []map[string]string{
		{"grant_schema": testDB1, "privilege": "read_only"},
		{"grant_schema": testDB3, "privilege": "dml"},
	}
	initialPrivilegesStr, err := json.Marshal(initialPrivileges)
	if err != nil {
		t.Fatal(err)
	}

	// 更新后的权限配置
	updatedPrivileges := []map[string]string{
		{"grant_schema": testDB1, "privilege": "rw"},
		{"grant_schema": testDB2, "privilege": "ddl"},
	}

	updatedPrivilegesStr, err := json.Marshal(updatedPrivileges)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建测试（带初始权限）
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd, mysqlInstanceID, projectID,
					accountName, accountPassword,
					initialPrivilegesStr, "Initial description",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", mysqlInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "name", accountName),
					resource.TestCheckResourceAttr(resourceName, "description", "Initial description"),
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.0.grant_schema", testDB1),
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.0.privilege", "read_only"),
				),
			},
			// 2. 更新测试（修改密码、权限和描述）
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd, mysqlInstanceID, projectID,
					accountName, newPassword,
					updatedPrivilegesStr, "Updated description",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", mysqlInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "name", accountName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.0.grant_schema", testDB1),
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.0.privilege", "rw"),
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.1.grant_schema", testDB2),
					resource.TestCheckResourceAttr(resourceName, "schema_privilege_list.1.privilege", "ddl"),
				),
			},
			// 3. 资源导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					// 构造导入ID: "id,region_id"
					return fmt.Sprintf("%s,%s,%s,%s",
						rs.Primary.Attributes["name"],
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "description", "project_id"},
			},
			// 3. 资源导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					// 构造导入ID: "id,region_id"
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["name"],
						rs.Primary.Attributes["instance_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "description", "project_id"},
			},
			{

				Config: utils.LoadTestCase(resourceFile,
					rnd, mysqlInstanceID, projectID,
					accountName, newPassword,
					updatedPrivilegesStr, "Updated description",
				) + utils.LoadTestCase(datasourceFile, dnd, mysqlInstanceID, accountName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "accounts.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "accounts.0.name", accountName),
					resource.TestCheckResourceAttrSet(dataSourceName, "accounts.0.schema_privilege_list.#"),
				),
			},
			{

				Config: utils.LoadTestCase(resourceFile,
					rnd, mysqlInstanceID, projectID,
					accountName, newPassword,
					updatedPrivilegesStr, "Updated description",
				),
				Destroy: true,
			},
		},
	})
}
