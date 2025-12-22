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

func TestAccCtyunPostgresqlDatabase(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_database." + rnd
	resourceFile := "resource_ctyun_postgresql_database_none_charset.tf"

	charsetDatasourceName := "data.ctyun_postgresql_character_set." + dnd
	charsetDatasourceFile := "datasource_ctyun_postgresql_character_set.tf"

	collationDatasourceName := "data.ctyun_postgresql_collation_time_zone." + dnd
	collationDatasourceFile := "datasource_ctyun_postgresql_collation_time_zone.tf"

	databasesDatasourceName := "data.ctyun_postgresql_databases." + dnd
	databasesDatasourceFile := "datasource_ctyun_postgresql_databases.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	instanceID := dependence.pgsqlID
	ownerAccount := "kqjwyk"

	// 测试数据
	dbName := "test_db_" + rnd
	charset := "UTF8"
	initialDescription := "Initial_database_description"
	updatedDescription := "Updated_database_description"

	// 等待函数
	wait20Seconds := func() {
		t.Logf("等待20秒...")
		time.Sleep(20 * time.Second)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// collation datasource 验证
			{
				Config: utils.LoadTestCase(collationDatasourceFile, dnd, instanceID, projectID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(collationDatasourceName, "standard_time_offset"),
					resource.TestCheckResourceAttrSet(collationDatasourceName, "time_zone"),
					resource.TestCheckResourceAttrSet(collationDatasourceName, "collations.#")),
			},
			// char set datasource验证
			{
				Config: utils.LoadTestCase(charsetDatasourceFile, dnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(charsetDatasourceName, "character_set.#"),
				),
			},

			// 1. 创建数据库测试（UTF8字符集）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					projectID,
					instanceID, dbName,
					charset,
					ownerAccount, initialDescription,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceID),
					resource.TestCheckResourceAttr(resourceName, "name", dbName),
					resource.TestCheckResourceAttr(resourceName, "charset_name", charset),
					resource.TestCheckResourceAttr(resourceName, "description", initialDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "owner", ownerAccount),
				),
			},
			// 2. 更新数据库描述测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					instanceID, dbName,
					charset,
					ownerAccount, updatedDescription,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
			// datasource 验证
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					instanceID, dbName,
					charset,
					ownerAccount, updatedDescription,
				) + utils.LoadTestCase(
					databasesDatasourceFile, dnd, instanceID, dbName),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(databasesDatasourceName, "databases.#")),
			},
			// 4. 导入测试
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
				ImportStateVerifyIgnore: []string{"charset_collate", "charset_type", "owner", "charset_name", "description"}, // 可选忽略
				PreConfig: func() {
					wait20Seconds()
				},
			},
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
				ImportStateVerifyIgnore: []string{"charset_collate", "charset_type", "owner", "charset_name", "description", "project_id"}, // 可选忽略
				PreConfig: func() {
					wait20Seconds()
				},
			},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					instanceID, dbName,
					charset,
					ownerAccount, updatedDescription,
				),
				Destroy: true,
				PreConfig: func() {
					wait20Seconds()
				},
			},
		},
	})
}

// 测试用例2：使用其他字符集（需要提供排序规则和类型）
func TestAccCtyunPostgresqlDatabaseWithOtherCharset(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_database." + rnd
	resourceFile := "resource_ctyun_postgresql_database.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	instanceID := dependence.pgsqlID
	ownerAccount := dependence.accountName

	// 测试数据
	dbName := "test_db_" + rnd
	charset := dependence.charsetName
	collate := dependence.collateName
	charType := dependence.collateType
	description := "Database_with_specific_charset"

	// 等待函数
	wait20Seconds := func() {
		t.Logf("等待20秒...")
		time.Sleep(20 * time.Second)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建数据库测试（指定字符集、排序规则和类型）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					projectID,
					instanceID, dbName,
					charset, collate, charType,
					ownerAccount, description,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceID),
					resource.TestCheckResourceAttr(resourceName, "name", dbName),
					resource.TestCheckResourceAttr(resourceName, "charset_name", charset),
					resource.TestCheckResourceAttr(resourceName, "charset_collate", collate),
					resource.TestCheckResourceAttr(resourceName, "charset_type", charType),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "owner", ownerAccount),
				),
				PreConfig: func() {
					wait20Seconds()
				},
			},
			// 2. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					projectID,
					instanceID, dbName,
					charset, collate, charType,
					ownerAccount, description,
				),
				Destroy: true,
			},
		},
	})
}
