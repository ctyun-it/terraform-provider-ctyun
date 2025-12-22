package ec_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunExpressConnectVpcInstance(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_ec_vpc_instance." + rnd
	resourceFile := "resource_ctyun_ec_vpc_instance.tf"

	datasourceFile := "datasource_ctyun_ec_vpc_instances.tf"
	datasourceName := "data.ctyun_ec_vpc_instances." + dnd
	// 从环境变量获取测试依赖资源
	ecID := dependence.expressConnectID
	cgwID := dependence.cloudGatewayId
	rtbID := dependence.rtbID
	vpcID := dependence.vpcID
	//exclusiveID := os.Getenv("CTYUN_EXCLUSIVE_ID")

	// 测试数据 - 子网配置
	initialSubnets := fmt.Sprintf(`"%s"`, dependence.subnetID)
	updatedSubnets := fmt.Sprintf(`"%s", "%s"`, dependence.subnetID, dependence.subnetID2)
	routeLearn := 1
	routeSync := 1

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建VPC实例测试（默认路由学习和同步）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, rtbID, vpcID, routeLearn, routeSync,
					initialSubnets,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ec_id", ecID),
					resource.TestCheckResourceAttr(resourceName, "cgw_id", cgwID),
					resource.TestCheckResourceAttr(resourceName, "rtb_id", rtbID),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "route_learn", "1"),
					resource.TestCheckResourceAttr(resourceName, "route_sync", "1"),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "subnets.0", dependence.subnetID),
					//resource.TestCheckResourceAttr(resourceName, "subnets.0.ip_version", "ipv4"),
					//resource.TestCheckResourceAttr(resourceName, "subnets.0.cidr", "192.168.1.0/24"),
					//resource.TestCheckResourceAttr(resourceName, "subnets.0.subnet_name", "test-subnet-1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 更新子网列表测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, rtbID, vpcID, routeLearn, routeSync,
					updatedSubnets,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "2"),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, rtbID, vpcID, routeLearn, routeSync,
					updatedSubnets,
				) + utils.LoadTestCase(datasourceFile, dnd, ecID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "vpc_instances.#")),
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
					return fmt.Sprintf("%s,%s,%s",
						rs.Primary.ID,
						rs.Primary.Attributes["ec_id"],
						rs.Primary.Attributes["cgw_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"exclusive_id", "project_id"}, // 子网列表可能变化

			},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, rtbID, vpcID, routeLearn, routeSync,
					updatedSubnets,
				),
				Destroy: true,
			},
		},
	})
}
