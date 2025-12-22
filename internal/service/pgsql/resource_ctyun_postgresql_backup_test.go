package pgsql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunPostgresqlBackup(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_backup." + rnd
	resourceFile := "resource_ctyun_postgresql_backup.tf"

	// 从环境变量获取测试依赖资源

	projectID := "0"
	instanceID := dependence.pgsqlID

	// 测试数据
	backupName := "test_backup_" + rnd
	description := "Test backup created by Terraform"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建备份测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					projectID,
					instanceID, backupName,
					description,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceID),
					resource.TestCheckResourceAttr(resourceName, "name", backupName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "backup_type"),
					resource.TestCheckResourceAttrSet(resourceName, "backup_result"),
				),
			},
			// 2. 导入测试
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
				ImportStateVerifyIgnore: []string{"description", "start_time", "end_time", "backup_result", "backup_type"},
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
				ImportStateVerifyIgnore: []string{"description", "start_time", "end_time", "backup_result", "backup_type", "project_id"},
			},
			// 3. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					instanceID, backupName,
					description,
				),
				Destroy: true,
			},
		},
	})
}
