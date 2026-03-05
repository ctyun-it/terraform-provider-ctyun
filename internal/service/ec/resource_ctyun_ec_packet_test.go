package ec_test

import (
	"fmt"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccEcPacket_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_ec_packet." + rnd
	resourceFile := "resource_ctyun_ec_packet.tf"

	ecID := dependence.expressConnectID
	packetName := utils.GenerateRandomString()
	bandwidth := "5"
	cycleType := "month"
	cycleCount := "1"
	bandwidthUpdate := "2"
	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			// 带宽包订购是一次性操作，无需特殊清理
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(resourceFile, rnd, ecID, packetName, bandwidth, cycleType, cycleCount),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ec_id", ecID),
					resource.TestCheckResourceAttr(resourceName, "name", packetName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth", bandwidth),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "cycle_count", cycleCount),
				),
			},
			{
				// 测试导入
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s", rs.Primary.Attributes["id"], rs.Primary.Attributes["ec_id"]), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"master_order_id",
					"master_order_no",
					"master_resource_id",
					"master_resource_status",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, ecID, packetName, bandwidthUpdate, cycleType, cycleCount),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ec_id", ecID),
					resource.TestCheckResourceAttr(resourceName, "name", packetName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth", bandwidthUpdate),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "cycle_count", cycleCount),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, ecID, packetName, bandwidth, cycleType, cycleCount),
				Destroy: true,
			},
		},
	})
}
