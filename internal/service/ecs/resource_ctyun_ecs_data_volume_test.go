package ecs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunEcsDataVolume(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_ecs_data_volume." + rnd
	resourceFile := "resource_ctyun_ecs_data_volume.tf"
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
					dependence.ecsID,
					dependence.ebsID3,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", dependence.ecsID),
					resource.TestCheckResourceAttr(resourceName, "ebs_ids.0", dependence.ebsID3),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.ecsID,
					fmt.Sprintf(`%s","%s","%s`, dependence.ebsID, dependence.ebsID2, dependence.ebsID3),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", dependence.ecsID),
					resource.TestCheckResourceAttr(resourceName, "ebs_ids.0", dependence.ebsID),
					resource.TestCheckResourceAttr(resourceName, "ebs_ids.1", dependence.ebsID2),
					resource.TestCheckResourceAttr(resourceName, "ebs_ids.2", dependence.ebsID3),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
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
						rs.Primary.Attributes["instance_id"],
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
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					dependence.ecsID,
					fmt.Sprintf(`%s","%s","%s`, dependence.ebsID2, dependence.ebsID, dependence.ebsID3),
				),
				Destroy: true,
			},
		},
	},
	)
}
