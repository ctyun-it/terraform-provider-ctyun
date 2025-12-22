package ebm_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunEbmInterface(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_ebm_interface." + rnd
	resourceFile := "resource_ctyun_ebm_interface.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.ebmID,
					dependence.securityGroupID,
					dependence.subnetID,
					dependence.az2,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.0", dependence.securityGroupID),
					resource.TestCheckResourceAttrSet(resourceName, "interface_id"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.ebmID,
					dependence.securityGroupID2,
					dependence.subnetID,
					dependence.az2,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.0", dependence.securityGroupID2),
					resource.TestCheckResourceAttrSet(resourceName, "interface_id"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 多种导入方式测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s,%s,%s",
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["interface_id"],
						rs.Primary.Attributes["az_name"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s,%s",
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["interface_id"],
						rs.Primary.Attributes["az_name"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.ebmID,
					dependence.securityGroupID2,
					dependence.subnetID,
					dependence.az2,
				),
				Destroy: true,
			},
		},
	})
}
