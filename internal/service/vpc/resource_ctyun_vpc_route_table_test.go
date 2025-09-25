package vpc_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunVpcRouteTable(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_vpc_route_table." + rnd
	datasourceName := "data.ctyun_vpc_route_tables." + dnd
	resourceFile := "resource_ctyun_vpc_route_table.tf"
	datasourceFile := "datasource_ctyun_vpc_route_tables.tf"

	initName := "terraform-unit"
	updatedName := "terraform-route-table"
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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, dependence.vpcID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, dependence.vpcID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, dependence.vpcID) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "route_tables.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "route_tables.0.name", updatedName),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, dependence.vpcID) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Destroy: true,
			},
		},
	})
}
