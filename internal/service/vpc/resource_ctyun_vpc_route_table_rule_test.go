package vpc_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunVpcRouteTableRule(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_vpc_route_table_rule." + rnd
	datasourceName := "data.ctyun_vpc_route_table_rules." + dnd
	resourceFile := "resource_ctyun_vpc_route_table_rule.tf"
	datasourceFile := "datasource_ctyun_vpc_route_table_rules.tf"

	initDestination := "188.188.0.0/16"
	initDescription := "test"
	updatedDescription := "updated"

	var id string
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
				Config: utils.LoadTestCase(resourceFile, rnd, initDestination, initDescription, dependence.vpcID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", initDescription),
					resource.TestCheckResourceAttr(resourceName, "destination", initDestination),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "rule_id"),
					func(state *terraform.State) error {
						rs, ok := state.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("resource not found")
						}
						id = rs.Primary.ID
						return nil
					},
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initDestination, updatedDescription, dependence.vpcID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "destination", initDestination),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initDestination, updatedDescription, dependence.vpcID) +
					utils.LoadTestCase(datasourceFile, dnd, "ctyun_vpc_route_table.route.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						ds := s.RootModule().Resources[datasourceName].Primary

						count, err := strconv.Atoi(ds.Attributes["rules.#"])
						if err != nil || count == 0 {
							return fmt.Errorf("rules 无效: %v", ds.Attributes)
						}

						for i := 0; i < count; i++ {
							if ds.Attributes[fmt.Sprintf("rules.%d.rule_id", i)] == id {
								return nil
							}
						}
						return fmt.Errorf("未找到目标元素")
					},
				)},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionID := ds.Attributes["region_id"]
					tableID := ds.Attributes["route_table_id"]
					return fmt.Sprintf("%s,%s,%s", id, tableID, regionID), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initDestination, updatedDescription, dependence.vpcID) +
					utils.LoadTestCase(datasourceFile, dnd, "ctyun_vpc_route_table.route.id"),
				Destroy: true,
			},
		},
	})
}
