package vpc_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunBandwidthAssociationEip(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_bandwidth_association_eip." + rnd
	resourceFile := "resource_ctyun_bandwidth_association_eip.tf"

	datasourceFile := "datasource_ctyun_bandwidths.tf"
	listDatasourceFile := "datasource_ctyun_bandwidths_list.tf"

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
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.bandwidthID, dependence.eipID),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.bandwidthID, dependence.eipID) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf(`"`+dependence.bandwidthID+`"`)),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.bandwidthID, dependence.eipID) +
					utils.LoadTestCase(listDatasourceFile, dnd),
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
				Config:  utils.LoadTestCase(resourceFile, rnd, dependence.bandwidthID, dependence.eipID),
				Destroy: true,
			},
		},
	})
}
