package peer_connection_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunVpcPeerConnectionRoute_basic(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_vpc_route_table_rule." + rnd
	resourceFile := "resource_ctyun_vpc_peer_connection_route.tf"

	nextHopID := dependence.peerConnectionID // 对等连接ID作为下一跳
	rtbID := dependence.rtbID
	description := "对等连接路由添加测试"

	// 测试数据
	ipVersions := []int32{4, 6}
	destinationCIDRs := []string{"172.168.1.0/24", "3ffe::/16"}
	for idx, ipVersion := range ipVersions {
		destinationCIDR := destinationCIDRs[idx]
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				// 1. 创建VPC对等连接路由测试
				{
					Config: utils.LoadTestCase(
						resourceFile, rnd,
						ipVersion, nextHopID, destinationCIDR, rtbID, description,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						// 基本属性验证
						resource.TestCheckResourceAttr(resourceName, "ip_version", fmt.Sprintf("%d", ipVersion)),
						resource.TestCheckResourceAttr(resourceName, "next_hop_id", nextHopID),
						resource.TestCheckResourceAttr(resourceName, "destination", destinationCIDR),
						// 系统生成属性验证
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
						return fmt.Sprintf("%s,%s,%s",
							rs.Primary.Attributes["rule_id"],
							rs.Primary.Attributes["route_table_id"],
							rs.Primary.Attributes["region_id"],
						), nil
					},
					ImportStateVerify: true,
				},
				// 4. 清理资源
				{
					Config: utils.LoadTestCase(
						resourceFile, rnd,
						ipVersion, nextHopID, destinationCIDR, rtbID, description,
					),
					Destroy: true,
				},
			},
		})
	}

}
