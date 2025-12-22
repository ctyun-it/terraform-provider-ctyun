package ec_test

import (
	"fmt"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccExpressConnect_update(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_express_connect." + rnd
	resourceFile := "resource_ctyun_express_connect.tf"
	dnd := utils.GenerateRandomString()

	datasourceName := "data.ctyun_express_connects." + dnd
	datasourceFile := "datasource_ctyun_express_connects.tf"
	name := utils.GenerateRandomString()
	description := "Initial description"
	updatedDescription := "Updated description"

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
				Config: utils.LoadTestCase(resourceFile, rnd, name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, updatedDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
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
					if ID == "" {
						return "", fmt.Errorf("id is not set")
					}

					return fmt.Sprintf("%s", ID), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"resource_id",
					"region_id",
				},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, name, updatedDescription),
				Destroy: true,
			},
			{
				Config: utils.LoadTestCase(datasourceFile, dnd, dependence.expressConnectID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "id", dependence.expressConnectID),
				),
			},
		},
	})
}
