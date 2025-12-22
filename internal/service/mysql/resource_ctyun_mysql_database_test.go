package mysql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunMysqlDatabase(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_database." + rnd
	resourceFile := "resource_ctyun_mysql_database.tf"
	datasourceName := "data.ctyun_mysql_databases." + dnd
	datasourceFile := "datasource_ctyun_mysql_databases.tf"

	charsetDatasourceName := "data.ctyun_mysql_character_set." + dnd
	charsetDatasourceFile := "datasource_ctyun_mysql_character_set.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	mysqlInstanceID := dependence.mysqlID

	// 数据库配置
	dbName := "test_db_" + rnd

	initialDescription := "Initial database description"
	updatedDescription := "Updated database description"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// char set datasource验证
			{
				Config: utils.LoadTestCase(charsetDatasourceFile, dnd, mysqlInstanceID, projectID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(charsetDatasourceName, "character_set.#"),
				),
			},
			// 1. 创建数据库测试
			{
				Config: utils.LoadTestCase(charsetDatasourceFile, dnd, mysqlInstanceID, projectID) + utils.LoadTestCase(
					resourceFile, rnd, mysqlInstanceID, projectID,
					dbName, fmt.Sprintf("%s.character_set.0.charset", charsetDatasourceName), initialDescription,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", mysqlInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "name", dbName),
					resource.TestCheckResourceAttr(resourceName, "description", initialDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 更新数据库描述测试
			{
				Config: utils.LoadTestCase(charsetDatasourceFile, dnd, mysqlInstanceID, projectID) + utils.LoadTestCase(
					resourceFile, rnd, mysqlInstanceID, projectID,
					dbName, fmt.Sprintf("%s.character_set.0.charset", charsetDatasourceName), updatedDescription,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", mysqlInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "name", dbName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// database datasource验证
			{
				Config: utils.LoadTestCase(charsetDatasourceFile, dnd, mysqlInstanceID, projectID) + utils.LoadTestCase(
					resourceFile, rnd, mysqlInstanceID, projectID, dbName, fmt.Sprintf("%s.character_set.0.charset", charsetDatasourceName), updatedDescription) + utils.LoadTestCase(datasourceFile, dnd, mysqlInstanceID, projectID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "databases.#"),
				),
			},
			// 3. 导入测试
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
				ImportStateVerifyIgnore: []string{"charset_name", "description"}, // 不需要忽略任何字段
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
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["name"],
						rs.Primary.Attributes["instance_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"charset_name", "description", "project_id"}, // 不需要忽略任何字段
			},
			// 5. 清理资源
			{
				Config: utils.LoadTestCase(charsetDatasourceFile, dnd, mysqlInstanceID, projectID) + utils.LoadTestCase(
					rnd, resourceFile,
					mysqlInstanceID, projectID,
					dbName, fmt.Sprintf("%s.character_set.0.charset", charsetDatasourceName), updatedDescription,
				),
				Destroy: true,
			},
		},
	})
}
