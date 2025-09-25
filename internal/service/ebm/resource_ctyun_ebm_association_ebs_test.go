package ebm_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunEbmAssociationEbs(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_ebm_association_ebs." + rnd
	resourceFile := "resource_ctyun_ebm_association_ebs.tf"

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
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.ebsID,
					dependence.ebmID,
					dependence.az2,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.ebsID,
					dependence.ebmID,
					dependence.az2,
				),
				Destroy: true,
			},
		},
	})
}
