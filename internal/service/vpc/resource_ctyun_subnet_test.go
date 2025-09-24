package vpc_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunSubnet(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_subnet." + rnd
	datasourceName := "data.ctyun_subnets." + dnd
	resourceFile := "resource_ctyun_subnet.tf"
	datasourceFile := "datasource_ctyun_subnets.tf"

	initName := "init"
	initDescription := "description"
	initDns := "114.114.114.114"
	updatedName := "updated"
	updatedDescription := "updated-description"
	updatedDns := "8.8.8.8"

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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, initDescription, initDns, dependence.vpcID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "description", initDescription),
					resource.TestCheckTypeSetElemAttr(resourceName, "dns.*", initDns),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription, updatedDns, dependence.vpcID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckTypeSetElemAttr(resourceName, "dns.*", updatedDns),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription, updatedDns, dependence.vpcID) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "subnets.0.name", updatedName),
					resource.TestCheckResourceAttr(datasourceName, "subnets.0.description", updatedDescription),
					resource.TestCheckTypeSetElemAttr(datasourceName, "subnets.0.dns_list.*", updatedDns),
				),
			},
			{

				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					vpcID := ds.Attributes["vpc_id"]
					if id == "" || regionId == "" || vpcID == "" {
						return "", fmt.Errorf("id or region_id os vpcID is required")
					}
					return fmt.Sprintf("%s,%s,%s", id, vpcID, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project_id",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription, updatedDns, dependence.vpcID) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Destroy: true,
			},
		},
	})
}
