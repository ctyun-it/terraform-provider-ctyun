package ec_test

import (
	"fmt"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccEcCloudGateway_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_ec_cloud_gateway." + rnd
	resourceFile := "resource_ctyun_ec_cloud_gateway.tf"
	datasourceName := "data.ctyun_ec_cloud_gateways." + dnd
	datasourceFile := "datasource_ctyun_ec_cloud_gateways.tf"

	name := utils.GenerateRandomString()
	description := "terrform 测试专用"
	updatedDescription := "Updated description"
	regionType := "CNP"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.expressConnectID, name, description, regionType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.expressConnectID, name, updatedDescription, regionType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},

			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.expressConnectID, name, updatedDescription, regionType) +
					utils.LoadTestCase(datasourceFile, dnd, dependence.expressConnectID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "ec_id", dependence.expressConnectID),
				),
			},
			{
				// 测试导入
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceName)
					}
					ID := rs.Primary.ID
					EcId := rs.Primary.Attributes["ec_id"]
					return fmt.Sprintf("%s,%s", EcId, ID), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"region_id",
				},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, name, updatedDescription),
				Destroy: true,
			},
		},
	})
}
