package vpc_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunDhcpOptionSetAssociationVpc_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_dhcpoptionset_association_vpc." + rnd
	resourceFile := "resource_ctyun_dhcpoptionset_association_vpc.tf"
	datasourceName := "data.ctyun_dhcpoptionset_association_vpcs." + dnd
	datasourceFile := "datasource_ctyun_dhcpoptionset_association_vpcs.tf"

	// 测试参数
	dhcpOptionSetsId := dependence.dhcpID
	vpcIds := fmt.Sprintf(`"%s"`, dependence.vpcID)
	updatedVpcIds := fmt.Sprintf(`"%s","%s"`, dependence.vpcID, dependence.vpcID)

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
				// 测试创建DHCP选项集与VPC绑定关系
				Config: utils.LoadTestCase(resourceFile, rnd, dhcpOptionSetsId, vpcIds),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "dhcp_option_sets_id", dhcpOptionSetsId),
					resource.TestCheckResourceAttr(resourceName, "vpc_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dhcpOptionSetsId, vpcIds) +
					utils.LoadTestCase(datasourceFile, dnd, dhcpOptionSetsId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "dhcp_option_sets_id", dhcpOptionSetsId),
					resource.TestCheckResourceAttr(datasourceName, "vpcs.#", "1"),
				),
			},
			{
				// 测试更新DHCP选项集与VPC绑定关系
				Config: utils.LoadTestCase(resourceFile, rnd, dhcpOptionSetsId, updatedVpcIds),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "dhcp_option_sets_id", dhcpOptionSetsId),
				),
			},

			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.Attributes["dhcp_option_sets_id"]
					regionId := ds.Attributes["region_id"]
					if id == "" || regionId == "" {
						return "", fmt.Errorf("id or region_id is required")
					}
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},

			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.Attributes["dhcp_option_sets_id"]
					return fmt.Sprintf("%s", id), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"region_id",
				},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, dhcpOptionSetsId, vpcIds),
				Destroy: true,
			},
		},
	})
}
