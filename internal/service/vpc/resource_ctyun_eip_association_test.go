package vpc_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunEipAssociation(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_eip_association." + rnd
	resourceFile := "resource_ctyun_eip_association.tf"

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
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.ecsID, dependence.eipID),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project_id",
				},
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					regionId := ds.Attributes["region_id"]
					eipId := ds.Attributes["eip_id"]
					return fmt.Sprintf("%s,%s", eipId, regionId), nil // eipId is not used
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project_id",
					"cycle_type",
					"master_order_id",
					"demand_billing_type",
				},
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					eipId := ds.Attributes["eip_id"]
					return fmt.Sprintf("%s", eipId), nil // eipId is not used
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project_id",
					"cycle_type",
					"master_order_id",
					"demand_billing_type",
				},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, dependence.ecsID, dependence.eipID),
				Destroy: true,
			},
		},
	})
}
