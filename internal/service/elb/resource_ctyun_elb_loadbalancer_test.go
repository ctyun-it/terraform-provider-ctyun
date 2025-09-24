package elb_test

//
//import (
//	"fmt"
//	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
//	"github.com/hashicorp/terraform-plugin-testing/terraform"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
//	"testing"
//)
//
//// 测试经典型负载均衡创建、修改基本信息和销毁
//func TestAccCtyunElbLoadBalancer(t *testing.T) {
//
//	rnd := utils.GenerateRandomString()
//	dnd := utils.GenerateRandomString()
//
//	resourceName := "ctyun_elb_loadbalancer." + rnd
//	datasourceName := "data.ctyun_elb_loadbalancers." + dnd
//	resourceFile := "resource_ctyun_elb_loadbalancer.tf"
//	datasourceFile := "datasource_ctyun_elb_loadbalancers.tf"
//	name := "elb_" + utils.GenerateRandomString()
//	slaName := "elb.s1.small"
//	resourceType := "internal"
//
//	updateName := "elb_" + utils.GenerateRandomString()
//
//	// 等代码互通了，需要更新的字段
//	vpcID := dependence.vpcID
//	subnetID := dependence.subnetID
//
//	resource.Test(t, resource.TestCase{
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			{
//				// create 验证
//				Config: utils.LoadTestCase(resourceFile, rnd, subnetID, name, slaName, resourceType, vpcID, "", "", "", ""),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
//					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
//					resource.TestCheckResourceAttr(resourceName, "name", name),
//					resource.TestCheckResourceAttr(resourceName, "sla_name", slaName),
//					resource.TestCheckResourceAttr(resourceName, "resource_type", resourceType),
//				),
//			},
//			{
//				// update 验证
//				Config: utils.LoadTestCase(resourceFile, rnd, subnetID, updateName, slaName, resourceType, vpcID, "", "", "", ""),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
//					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
//					resource.TestCheckResourceAttr(resourceName, "name", updateName),
//					resource.TestCheckResourceAttr(resourceName, "sla_name", slaName),
//					resource.TestCheckResourceAttr(resourceName, "resource_type", resourceType),
//				),
//			},
//			{
//				// datasource验证
//				Config: utils.LoadTestCase(datasourceFile, dnd),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(datasourceName, "elbs.0.vpc_id", vpcID),
//					resource.TestCheckResourceAttr(datasourceName, "elbs.0.subnet_id", subnetID),
//					resource.TestCheckResourceAttr(datasourceName, "elbs.0.name", updateName),
//					resource.TestCheckResourceAttr(datasourceName, "elbs.0.sla_name", slaName),
//					resource.TestCheckResourceAttr(datasourceName, "elbs.0.resource_type", resourceType),
//				),
//			},
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, subnetID, updateName, slaName, resourceType, vpcID, "", "", "", ""),
//				Destroy: true,
//			},
//		},
//	})
//}
