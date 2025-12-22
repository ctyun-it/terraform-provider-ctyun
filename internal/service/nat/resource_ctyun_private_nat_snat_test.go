package nat_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunPrivateSNat(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	// 创建nat信息
	natGatewayID := dependence.privateNatID
	// 创建snat信息
	resourceName := "ctyun_private_nat_snat." + rnd
	datasourceName := "data.ctyun_private_nat_snats." + dnd
	resourceFile := "resource_ctyun_private_nat_snat.tf"
	datasourceFile := "datasource_ctyun_private_nat_snats.tf"

	sourceSubnetId := dependence.subnetID2
	updatedSubnetId := dependence.subnetID1

	snatIps := fmt.Sprintf(`["%s"]`, "192.168.128.100")
	updatedSnatIps := fmt.Sprintf(`["%s","%s"]`, "192.168.128.100", "192.168.128.101")

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
				// 1. 创建nat,snat
				// 1.1 resource create验证:
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayID, sourceSubnetId, snatIps, "我是一条description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_subnet_id", sourceSubnetId),
					resource.TestCheckResourceAttr(resourceName, "addresses.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "snat_id"),
					resource.TestCheckResourceAttr(resourceName, "description", "我是一条description"),
				),
			},
			{
				// 1.2 resource update source_subnet_id验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayID, updatedSubnetId, updatedSnatIps, "我是一条description plus"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_subnet_id", updatedSubnetId),
					resource.TestCheckResourceAttr(resourceName, "addresses.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "description", "我是一条description plus"),
				),
			},
			{
				// 1.3. datasource验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayID, updatedSubnetId, updatedSnatIps, "我是一条description add") +
					utils.LoadTestCase(datasourceFile, dnd, natGatewayID, fmt.Sprintf(`snat_id=%s.snat_id`, resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "snats.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "snats.0.addresses.#", "2"),
					resource.TestCheckResourceAttr(datasourceName, "snats.0.source_subnet_id", updatedSubnetId),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					natGatewayId := ds.Attributes["nat_gateway_id"]
					return fmt.Sprintf("%s,%s", id, natGatewayId), nil
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					natGatewayId := ds.Attributes["nat_gateway_id"]
					regionId := ds.Attributes["region_id"]
					projectId := ds.Attributes["project_id"]
					return fmt.Sprintf("%s,%s,%s,%s", id, natGatewayId, projectId, regionId), nil
				},
			},
			// 1.4 资源销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, natGatewayID, updatedSubnetId, updatedSnatIps),
				Destroy: true,
			},
		},
	})
}
