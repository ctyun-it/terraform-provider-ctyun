package vpc_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunDhcpOptionSet_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_dhcpoptionset." + rnd
	resourceFile := "resource_ctyun_dhcpoptionset.tf"

	description := "Example DHCP option set for demonstration"
	domainName := "example.com"
	dnsList := `"8.8.8.8", "8.8.4.4"`

	updatedDescription := "Updated DHCP option set for demonstration"
	updatedDomainName := "updated.example.com"
	updatedDnsList := `"1.1.1.1", "8.8.8.8", "8.8.4.4"`
	dnd := utils.GenerateRandomString()

	datasourceName := "data.ctyun_dhcpoptionsets." + dnd
	datasourceFile := "datasource_ctyun_dhcpoptionsets.tf"
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
				// 测试创建带所有参数的DhcpOptionSet
				Config: utils.LoadTestCase(resourceFile, rnd, description, domainName, dnsList),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rnd),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "domain_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "dns_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_list.0", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "dns_list.1", "8.8.4.4"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				// 测试更新DhcpOptionSet
				Config: utils.LoadTestCase(resourceFile, rnd, updatedDescription, updatedDomainName, updatedDnsList),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rnd),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "domain_name", updatedDomainName),
					resource.TestCheckResourceAttr(resourceName, "dns_list.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "dns_list.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "dns_list.1", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "dns_list.2", "8.8.4.4"),
				),
			},

			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedDescription, updatedDomainName, updatedDnsList) +
					utils.LoadTestCase(datasourceFile, dnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "dhcpoptionsets.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "total_count"),
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
						return "", fmt.Errorf("id or region_id is required")
					}
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"region_id",
				},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, updatedDescription, updatedDomainName, updatedDnsList),
				Destroy: true,
			},
		},
	})
}
