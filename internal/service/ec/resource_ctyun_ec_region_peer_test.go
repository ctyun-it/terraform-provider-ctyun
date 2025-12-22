package ec_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunExpressConnectRegionPeer(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_ec_region_peer." + rnd
	resourceFile := "resource_ctyun_ec_region_peer.tf"

	datasourceFile := "datasource_ctyun_ec_regions_peers.tf"
	datasourceName := "data.ctyun_ec_region_peers." + dnd

	// 从环境变量获取测试依赖资源
	ecID := dependence.expressConnectID
	srcCgwID := dependence.cgwID1
	dstCgwID := dependence.cgwID2
	packetID := dependence.packetID

	// 测试数据
	peerName := "test-region-peer-" + rnd
	initialRate := int32(1) // 10 Mbps
	updatedRate := int32(2) // 2 Mbps

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建跨域连接测试（默认路由学习）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					peerName, ecID, srcCgwID, dstCgwID,
					packetID, initialRate, 1, // route_learn = 1 (默认值)
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", peerName),
					resource.TestCheckResourceAttr(resourceName, "ec_id", ecID),
					resource.TestCheckResourceAttr(resourceName, "src_cgw_id", srcCgwID),
					resource.TestCheckResourceAttr(resourceName, "dst_cgw_id", dstCgwID),
					resource.TestCheckResourceAttr(resourceName, "packet_id", packetID),
					resource.TestCheckResourceAttr(resourceName, "rate", fmt.Sprintf("%d", initialRate)),
					resource.TestCheckResourceAttr(resourceName, "route_learn", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "src_region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "dst_region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "peer_type"),
					resource.TestCheckResourceAttrSet(resourceName, "update_time"),
				),
			},
			// 2. 更新带宽测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					peerName, ecID, srcCgwID, dstCgwID,
					packetID, updatedRate, 1,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rate", fmt.Sprintf("%d", updatedRate)),
					resource.TestCheckResourceAttr(resourceName, "route_learn", "1"),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					peerName, ecID, srcCgwID, dstCgwID,
					packetID, updatedRate, 1,
				) + utils.LoadTestCase(datasourceFile, dnd, ecID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "region_peers.#")),
			},
			//// 3. 更新路由学习设置测试
			//{
			//	Config: utils.LoadTestCase(
			//		resourceFile, rnd,
			//		peerName, ecID, srcCgwID, dstCgwID,
			//		packetID, updatedRate, 0, // route_learn = 0 (关闭路由学习)
			//	),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr(resourceName, "route_learn", "0"),
			//	),
			//},
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
						rs.Primary.Attributes["id"],
						rs.Primary.Attributes["ec_id"],
						rs.Primary.Attributes["packet_id"],
						rs.Primary.Attributes["src_cgw_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_time", "route_learn"}, // 更新时间可能变化
			},
			// 5. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					peerName, ecID, srcCgwID, dstCgwID,
					packetID, updatedRate, 1,
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例2: 不同带宽值测试
func TestAccCtyunExpressConnectRegionPeerDifferentRates(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_ec_region_peer." + rnd
	resourceFile := "resource_ctyun_ec_region_peer.tf"

	// 从环境变量获取测试依赖资源
	ecID := dependence.expressConnectID
	srcCgwID := dependence.cgwID1
	dstCgwID := dependence.cgwID2
	packetID := dependence.packetID

	// 测试不同的带宽值
	rates := []int32{1, 2, 5}

	for _, rate := range rates {
		t.Run(fmt.Sprintf("Rate_%dMbps", rate), func(t *testing.T) {
			peerName := fmt.Sprintf("test-peer-%s-%d", rnd, rate)

			// 等待函数

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					// 1. 创建连接测试
					{
						Config: utils.LoadTestCase(
							resourceFile, rnd,
							peerName, ecID, srcCgwID, dstCgwID,
							packetID, rate, 0,
						),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "rate", fmt.Sprintf("%d", rate)),
							resource.TestCheckResourceAttrSet(resourceName, "id"),
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
								rs.Primary.Attributes["id"],
								rs.Primary.Attributes["ec_id"],
								rs.Primary.Attributes["packet_id"],
								rs.Primary.Attributes["src_cgw_id"],
							), nil
						},
						ImportStateVerify: true,
						ImportStateVerifyIgnore: []string{
							"exclusive_id",
							"project_id",
							"route_learn",
						},
					},
					// 2. 清理资源
					{
						Config: utils.LoadTestCase(
							resourceFile, rnd,
							peerName, ecID, srcCgwID, dstCgwID,
							packetID, rate, 0,
						),
						Destroy: true,
					},
				},
			})
		})
	}
}
