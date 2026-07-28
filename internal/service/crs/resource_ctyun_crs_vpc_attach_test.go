package crs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunCrsVpcAttach(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_crs_vpc_attach." + rnd
	resourceFile := "resource_ctyun_crs_vpc_attach.tf"

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
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.vpcID,
					dependence.subnetID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"subnet_id",
				},
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
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["vpc_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"subnet_id"},
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s",
						rs.Primary.Attributes["vpc_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"subnet_id", "region_id"},
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.vpcID,
					dependence.subnetID,
				),
				Destroy: true,
			},
		},
	})
}
