package pgsql_test

//
//import (
//	"fmt"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
//	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
//	"github.com/hashicorp/terraform-plugin-testing/terraform"
//	"testing"
//)
//
//func TestAccCtyunPgsqlSqlCollectorPolicy(t *testing.T) {
//	t.Parallel()
//	rnd := utils.GenerateRandomString()
//	resourceName := "ctyun_postgresql_sql_collector_policy." + rnd
//	resourceFile := "resource_ctyun_postgresql_sql_collector_policy.tf"
//
//	instanceId := dependence.pgsqlID
//	sql_collector_status := "enable"
//	log_interval := 10
//	log_interval_update := 5
//
//	resource.Test(t, resource.TestCase{
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, instanceId, sql_collector_status, log_interval),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//				),
//			},
//
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, instanceId, sql_collector_status, log_interval_update),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//				),
//			},
//			// 2. 导入测试
//			{
//				ResourceName: resourceName,
//				ImportState:  true,
//				ImportStateIdFunc: func(s *terraform.State) (string, error) {
//					rs, ok := s.RootModule().Resources[resourceName]
//					if !ok {
//						return "", fmt.Errorf("resource not found: %s", resourceName)
//					}
//					return fmt.Sprintf("%s,%s,%s",
//						rs.Primary.Attributes["instance_id"],
//						rs.Primary.Attributes["project_id"],
//						rs.Primary.Attributes["region_id"],
//					), nil
//				},
//				ImportStateVerify:       true,
//				ImportStateVerifyIgnore: []string{}, // 不需要忽略任何字段
//			},
//			// 3. 导入测试
//			{
//				ResourceName: resourceName,
//				ImportState:  true,
//				ImportStateIdFunc: func(s *terraform.State) (string, error) {
//					rs, ok := s.RootModule().Resources[resourceName]
//					if !ok {
//						return "", fmt.Errorf("resource not found: %s", resourceName)
//					}
//					return fmt.Sprintf("%s",
//						rs.Primary.Attributes["instance_id"],
//					), nil
//				},
//				ImportStateVerify:       true,
//				ImportStateVerifyIgnore: []string{}, // 不需要忽略任何字段
//			},
//
//			// destroy
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, instanceId, sql_collector_status, log_interval_update),
//				Destroy: true,
//			},
//		},
//	})
//}
