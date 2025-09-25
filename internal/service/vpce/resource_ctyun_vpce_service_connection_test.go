package vpce_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunVpceServiceConnection(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_vpce_service_connection." + rnd
	resourceFile := "resource_ctyun_vpce_service_connection.tf"
	statusUp, statusDown := "up", "down"

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
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.reverseVpceServiceID, dependence.vpceID, statusDown),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "endpoint_service_id", dependence.reverseVpceServiceID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", dependence.vpceID),
					resource.TestCheckResourceAttr(resourceName, "status", statusDown),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.reverseVpceServiceID, dependence.vpceID, statusUp),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "endpoint_service_id", dependence.reverseVpceServiceID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", dependence.vpceID),
					resource.TestCheckResourceAttr(resourceName, "status", statusUp),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					if id == "" || regionId == "" {
						return "", fmt.Errorf("id or region_id or endpoint_service_id is required")
					}
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, dependence.reverseVpceServiceID, dependence.vpceID, statusUp),
				Destroy: true,
			},
		},
	})
}
