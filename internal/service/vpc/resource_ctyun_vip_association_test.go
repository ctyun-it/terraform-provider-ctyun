package vpc_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunVipAssociation_vm(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_vip_association." + rnd
	resourceFile := "resource_ctyun_vip_association.tf"

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
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.vipId, "VM", dependence.networkInterfaceID, dependence.ecsID, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vip_id", dependence.vipId),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "VM"),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", dependence.networkInterfaceID),
					resource.TestCheckResourceAttr(resourceName, "instance_id", dependence.ecsID),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, dependence.vipId, "VM", dependence.networkInterfaceID, dependence.ecsID, ""),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunVipAssociation_network(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_vip_association." + rnd
	resourceFile := "resource_ctyun_vip_association.tf"
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
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.vipId, "NETWORK", "", "", dependence.eipID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vip_id", dependence.vipId),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "NETWORK"),
					resource.TestCheckResourceAttr(resourceName, "floating_id", dependence.eipID),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, dependence.vipId, "NETWORK", "", "", dependence.eipID),
				Destroy: true,
			},
		},
	})
}
