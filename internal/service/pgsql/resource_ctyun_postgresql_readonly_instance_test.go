package pgsql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunPostgresqlReadOnlyInstance(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_readonly_instance." + rnd
	resourceFile := "resource_ctyun_postgresql_readonly_instance_on_demand.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"                 // 默认项目ID
	instanceID := dependence.pgsqlID // 主实例ID
	cycleType := "on_demand"

	// 测试数据
	instanceName := "test-pg-ro-" + rnd
	flavorName := "c7.xlarge.2" // PostgreSQL 规格

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建按需计费只读实例测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceID, cycleType, flavorName,
					projectID, instanceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cycle_type", "on_demand"),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
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
					return fmt.Sprintf("%s,%s,%s",
						rs.Primary.ID,
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"cycle_type", "cycle_count", "auto_renew", "flavor_name",
					"availability_zone_name", "instance_id"}, // 不需要忽略任何字段
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
					return fmt.Sprintf("%s",
						rs.Primary.ID,
					), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"cycle_type", "cycle_count", "auto_renew", "flavor_name",
					"availability_zone_name", "instance_id", "project_id"}, // 不需要忽略任何字段
			},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceID, cycleType, flavorName,
					projectID, instanceName,
				),
				Destroy: true,
			},
		},
	})
}
