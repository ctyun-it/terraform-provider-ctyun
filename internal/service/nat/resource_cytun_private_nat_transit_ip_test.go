package nat_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunPrivateNatTransitIp(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_private_nat_transit_ip." + rnd
	datasourceName := "data.ctyun_private_nat_transit_ips." + dnd

	resourceFile := "resource_ctyun_private_nat_transit_ip.tf"
	datasourceFile := "datasource_ctyun_private_nat_transit_ip.tf"

	natGatewayId := dependence.privateNatID
	address := "192.168.128.55"
	updatedAddress := "192.168.128.56"

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. resource create/delete 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, address),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "address", address),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", natGatewayId),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				// 2. resource update验证 (通过重新创建)
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, updatedAddress),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "address", updatedAddress),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", natGatewayId),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					region_id := ds.Attributes["region_id"]
					nat_gateway_id := ds.Attributes["nat_gateway_id"]
					address1 := ds.Attributes["address"]
					return fmt.Sprintf("%s,%s,%s", address1, nat_gateway_id, region_id), nil
				},
				ImportStateVerifyIgnore: []string{},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					nat_gateway_id := ds.Attributes["nat_gateway_id"]
					address1 := ds.Attributes["address"]
					return fmt.Sprintf("%s,%s", address1, nat_gateway_id), nil
				},
				ImportStateVerifyIgnore: []string{},
			},
			{
				// 3. datasource 验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, updatedAddress) +
					utils.LoadTestCase(datasourceFile, dnd, natGatewayId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "transit_ips.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "transit_ips.0.address"),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, natGatewayId, updatedAddress),
				Destroy: true,
			},
		},
	})
}
