package ec_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunExpressConnectRoute(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_ec_route." + rnd
	resourceFile := "resource_ctyun_ec_route.tf"
	resourceFile1 := "resource_ctyun_ec_route_blackhole.tf"

	// 从环境变量获取测试依赖资源
	ecID := dependence.expressConnectID
	cgwID := dependence.cloudGatewayId
	rtbID := dependence.rtbID
	cidr := "192.168.1.3/32"
	ipVersion := "ipv4"
	nextHopID := dependence.vpcInstanceVpcID
	nextHopType := "vpc"

	initDescription := "init description for route test"
	initDescription1 := "init black hole description for route test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建普通路由测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, rtbID, cidr, ipVersion, initDescription,
					false, nextHopType, nextHopID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ec_id", ecID),
					resource.TestCheckResourceAttr(resourceName, "cgw_id", cgwID),
					resource.TestCheckResourceAttr(resourceName, "rtb_id", rtbID),
					resource.TestCheckResourceAttr(resourceName, "cidr", cidr),
					resource.TestCheckResourceAttr(resourceName, "next_hop_type", nextHopType),
					resource.TestCheckResourceAttr(resourceName, "next_hop_id", nextHopID),
					resource.TestCheckResourceAttr(resourceName, "ip_version", ipVersion),
					resource.TestCheckResourceAttr(resourceName, "description", initDescription),
					resource.TestCheckResourceAttr(resourceName, "is_black_hole_route", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, rtbID, cidr, ipVersion, initDescription,
					false, nextHopType, nextHopID),
				Destroy: true,
			},

			// 2. 创建黑洞路由测试
			{
				Config: utils.LoadTestCase(
					resourceFile1, rnd,
					ecID, cgwID, rtbID, cidr, ipVersion, initDescription1,
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cidr", cidr),
					resource.TestCheckResourceAttr(resourceName, "ip_version", ipVersion),
					resource.TestCheckResourceAttr(resourceName, "description", initDescription1),
					resource.TestCheckResourceAttr(resourceName, "is_black_hole_route", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 3. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile1, rnd,
					ecID, cgwID, rtbID, cidr, ipVersion, initDescription1,
					true,
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例2: IPv6路由测试
//func TestAccCtyunExpressConnectRouteIPv6(t *testing.T) {
//	rnd := utils.GenerateRandomString()
//	resourceName := "ctyun_ec_route." + rnd
//	resourceFile := "resource_ctyun_ec_route.tf"
//
//	ecID := dependence.expressConnectID
//	cgwID := dependence.cloudGatewayId
//	rtbID := dependence.rtbID
//	cidr := "240e:982:da91:::/88"
//	ipVersion := "ipv6"
//	description := "IPv6 route description"
//	isBlackHole := false
//	nextHopType := "vpc"
//	nextHopID := dependence.vpcInstanceVpcID
//
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			// 1. 创建IPv6路由测试
//			{
//				Config: utils.LoadTestCase(
//					resourceFile, rnd,
//					ecID, cgwID, rtbID,
//					cidr, ipVersion, description, isBlackHole, nextHopType, nextHopID),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "ip_version", "ipv6"),
//					resource.TestCheckResourceAttr(resourceName, "cidr", "2001:db8::/32"),
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//				),
//			},
//			// 2. 清理资源
//			{
//				Config: utils.LoadTestCase(
//					resourceFile, rnd,
//					ecID, cgwID, rtbID,
//					cidr, ipVersion, description, isBlackHole, nextHopType, nextHopID),
//				Destroy: true,
//			},
//		},
//	})
//}

// 测试用例4: 不同下一跳类型测试
func TestAccCtyunExpressConnectRouteNextHopTypes(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceFile := "resource_ctyun_ec_route.tf"

	datasourceFile := "datasource_ctyun_ec_routes.tf"
	datasourceName := "data.ctyun_ec_routes." + dnd
	// 从环境变量获取测试依赖资源
	ecID := dependence.expressConnectID
	cgwID := dependence.cloudGatewayId
	rtbID := dependence.rtbID

	// 测试不同的下一跳类型
	nextHopTypesAndIDMap := map[string]string{
		"vpc": dependence.vpcInstanceVpcID,
		//"cda":   "",
		//"vpn":   "",
		"cross": dependence.regionPeerID,
	}
	i := 1
	for nextHopType, nextHopID := range nextHopTypesAndIDMap {
		t.Run(nextHopType, func(t *testing.T) {
			resourceName := "ctyun_ec_route." + rnd + "_" + nextHopType
			cidr := fmt.Sprintf("192.168.1.%d/32", i+10)
			ipVersion := "ipv4"
			description := fmt.Sprintf("Route with %s next hop", nextHopType)

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					// 1. 创建路由测试
					{
						Config: utils.LoadTestCase(
							resourceFile, rnd+"_"+nextHopType,
							ecID, cgwID, rtbID,
							cidr, ipVersion, description, false, nextHopType, nextHopID,
						),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "next_hop_type", nextHopType),
							resource.TestCheckResourceAttrSet(resourceName, "id"),
						),
					},
					{
						Config: utils.LoadTestCase(
							resourceFile, rnd+"_"+nextHopType,
							ecID, cgwID, rtbID,
							cidr, ipVersion, description, false, nextHopType, nextHopID,
						) + utils.LoadTestCase(datasourceFile, dnd,
							ecID, cgwID, rtbID),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttrSet(datasourceName, "routes.#")),
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
							return fmt.Sprintf("%s,%s,%s,%s,%s",
								rs.Primary.ID,
								rs.Primary.Attributes["ec_id"],
								rs.Primary.Attributes["cgw_id"],
								rs.Primary.Attributes["rtb_id"],
								rs.Primary.Attributes["next_hop_id"],
							), nil
						},
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"exclusive_id", "is_black_hole_route"}, // 子网列表可能变化

					},
					// 2. 清理资源
					{
						Config: utils.LoadTestCase(
							resourceFile, rnd+"_"+nextHopType,
							ecID, cgwID, rtbID,
							cidr, ipVersion, description, false, nextHopType, nextHopID,
						),
						Destroy: true,
					},
				},
			})
		})
		i = i + 1
	}
}
