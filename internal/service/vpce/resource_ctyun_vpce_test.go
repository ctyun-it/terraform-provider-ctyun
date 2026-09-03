package vpce_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunVpce(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_vpce." + rnd
	datasourceName := "data.ctyun_vpces." + dnd
	resourceFile := "resource_ctyun_vpce.tf"
	datasourceFile := "datasource_ctyun_vpces.tf"

	initName := "init"
	initWhitelistFlag := "true"
	initWhitelistCidr := `whitelist_cidr = ["192.168.1.0/24"]`
	updatedName := "updated"
	updatedWhitelistFlag := "false"
	updatedWhitelistCidr := ``

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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, initWhitelistFlag, initWhitelistCidr, dependence.vpcID, dependence.subnetID, dependence.vpceServiceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "whitelist_flag", "true"),
					resource.TestCheckTypeSetElemAttr(resourceName, "whitelist_cidr.*", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
					func(s *terraform.State) error {
						ds := s.RootModule().Resources[resourceName].Primary
						createTime := ds.Attributes["create_time"]
						if utils.IsEmptyOrRfc3339(createTime) {
							return nil
						}
						return fmt.Errorf("time format doesn't match")
					},
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedWhitelistFlag, updatedWhitelistCidr, dependence.vpcID, dependence.subnetID, dependence.vpceServiceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "whitelist_flag", "false"),
					resource.TestCheckResourceAttr(resourceName, "whitelist_cidr.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, initWhitelistFlag, initWhitelistCidr, dependence.vpcID, dependence.subnetID, dependence.vpceServiceID) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "vpces.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "vpces.0.name", updatedName),
					resource.TestCheckTypeSetElemAttr(datasourceName, "vpces.0.whitelist_cidr.*", "192.168.1.0/24"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id = ds.ID
					regionId := ds.Attributes["region_id"]
					if id == "" || regionId == "" {
						return "", fmt.Errorf("id or region_id is required")
					}
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"master_order_id",
					"project_id",
				},
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id = ds.ID
					return fmt.Sprintf("%s", id), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"master_order_id",
					"project_id",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, initWhitelistFlag, initWhitelistCidr, dependence.vpcID, dependence.subnetID, dependence.vpceServiceID) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunVpceFields(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_vpce." + rnd
	resourceFile := "resource_ctyun_vpce_fields.tf"

	initName := "init"
	initWhitelistFlag := "true"
	initWhitelistCidr := `whitelist_cidr = ["192.168.1.0/24"]`
	initDescription := "init description"
	initDeleteProtection := "true"
	initEnableDns := "false"

	updatedName := "updated"
	updatedWhitelistFlag := "true"
	updatedWhitelistCidr := `whitelist_cidr = ["192.168.1.0/24"]`
	updatedDescription := "updated description"
	updatedDeleteProtection := "false"
	updatedEnableDns := "true"

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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, initWhitelistFlag, initWhitelistCidr, dependence.vpcID, dependence.subnetID, dependence.vpceServiceID, initDescription, initDeleteProtection, initEnableDns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "description", initDescription),
					resource.TestCheckResourceAttr(resourceName, "delete_protection", initDeleteProtection),
					resource.TestCheckResourceAttr(resourceName, "enable_dns", initEnableDns),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedWhitelistFlag, updatedWhitelistCidr, dependence.vpcID, dependence.subnetID, dependence.vpceServiceID, updatedDescription, updatedDeleteProtection, updatedEnableDns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "delete_protection", updatedDeleteProtection),
					resource.TestCheckResourceAttr(resourceName, "enable_dns", updatedEnableDns),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					return fmt.Sprintf("%s", id), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"master_order_id",
					"project_id",
				},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, updatedName, updatedWhitelistFlag, updatedWhitelistCidr, dependence.vpcID, dependence.subnetID, dependence.vpceServiceID, updatedDescription, updatedDeleteProtection, updatedEnableDns),
				Destroy: true,
			},
		},
	})
}
