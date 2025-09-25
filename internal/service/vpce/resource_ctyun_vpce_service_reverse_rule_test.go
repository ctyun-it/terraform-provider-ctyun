package vpce_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunVpceServiceReverseRule(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_vpce_service_reverse_rule." + rnd
	datasourceName := "data.ctyun_vpce_service_reverse_rules." + dnd
	resourceFile := "resource_ctyun_vpce_service_reverse_rule.tf"
	datasourceFile := "datasource_ctyun_vpce_service_reverse_rules.tf"

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
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.reverseVpceServiceID, dependence.vpceID, dependence.transitIP, dependence.targetIP),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "endpoint_service_id", dependence.reverseVpceServiceID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", dependence.vpceID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
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
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.reverseVpceServiceID, dependence.vpceID, dependence.transitIP, dependence.targetIP) +
					utils.LoadTestCase(datasourceFile, dnd, dependence.reverseVpceServiceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							ds := s.RootModule().Resources[datasourceName].Primary

							count, err := strconv.Atoi(ds.Attributes["rules.#"])
							if err != nil || count == 0 {
								return fmt.Errorf("rules 无效: %v", ds.Attributes)
							}

							for i := 0; i < count; i++ {
								if ds.Attributes[fmt.Sprintf("rules.%d.id", i)] == id {
									return nil
								}
							}
							return fmt.Errorf("未找到目标元素")
						}),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id = ds.ID
					regionId := ds.Attributes["region_id"]
					endpointServiceID := ds.Attributes["endpoint_service_id"]
					if id == "" || regionId == "" {
						return "", fmt.Errorf("id or region_id or endpoint_service_id is required")
					}
					return fmt.Sprintf("%s,%s,%s", id, endpointServiceID, regionId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.reverseVpceServiceID, dependence.vpceID, dependence.transitIP, dependence.targetIP) +
					utils.LoadTestCase(datasourceFile, dnd, dependence.reverseVpceServiceID),
				Destroy: true,
			},
		},
	})
}
