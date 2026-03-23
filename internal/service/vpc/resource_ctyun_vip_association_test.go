package vpc_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunVipAssociation_port(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_vip_association." + rnd
	resourceFile := "resource_ctyun_vip_association_port.tf"

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
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.vipId, dependence.networkInterfaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vip_id", dependence.vipId),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", dependence.networkInterfaceID),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					vipID := ds.Attributes["vip_id"]
					portID := ds.Attributes["network_interface_id"]
					regionID := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s,%s", vipID, portID, regionID), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, dependence.vipId, dependence.networkInterfaceID),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunVipAssociation_network(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_vip_association." + rnd
	resourceFile := "resource_ctyun_vip_association_eip.tf"
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
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.vipId, dependence.eipID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vip_id", dependence.vipId),
					resource.TestCheckResourceAttr(resourceName, "floating_id", dependence.eipID),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					vipID := ds.Attributes["vip_id"]
					eipID := ds.Attributes["floating_id"]
					regionID := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s,%s", vipID, eipID, regionID), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, dependence.vipId, dependence.eipID),
				Destroy: true,
			},
		},
	})
}
