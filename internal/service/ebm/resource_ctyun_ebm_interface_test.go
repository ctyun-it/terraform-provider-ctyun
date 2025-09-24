package ebm_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunEbmInterface(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_ebm_interface." + rnd
	resourceFile := "resource_ctyun_ebm_interface.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.ebmID,
					dependence.securityGroupID,
					dependence.subnetID,
					dependence.az2,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.0", dependence.securityGroupID),
					resource.TestCheckResourceAttrSet(resourceName, "interface_id"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.ebmID,
					dependence.securityGroupID2,
					dependence.subnetID,
					dependence.az2,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.0", dependence.securityGroupID2),
					resource.TestCheckResourceAttrSet(resourceName, "interface_id"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{

				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionID := ds.Attributes["region_id"]
					instanceID := ds.Attributes["instance_id"]
					azName := ds.Attributes["az_name"]
					return fmt.Sprintf("%s,%s,%s,%s", instanceID, id, regionID, azName), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.ebmID,
					dependence.securityGroupID2,
					dependence.subnetID,
					dependence.az2,
				),
				Destroy: true,
			},
		},
	})
}
