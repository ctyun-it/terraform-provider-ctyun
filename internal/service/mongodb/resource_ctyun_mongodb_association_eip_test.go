package mongodb_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunMongodbAssociationEip(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_mongodb_association_eip." + rnd
	resourceFile := "resource_ctyun_mongodb_association_eip.tf"

	eipId := dependence.eipID
	instId := dependence.mongodbID

	instanceType := "1"

	specDatasourceName := "data.ctyun_mongodb_specs." + dnd
	specDatasourceFile := "datasource_ctyun_mongodb_specs.tf"
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
			// 绑定IP验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, eipId, instId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "eip_id", eipId),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instId),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, eipId, instId) +
					utils.LoadTestCase(specDatasourceFile, dnd, instanceType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(specDatasourceName, "specs.#"),
				),
			},
			//import验证
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s,%s,%s",
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["eip_id"],
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerifyIgnore: []string{
					"master_order_id", "project_id",
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s,%s",
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["eip_id"],
						rs.Primary.Attributes["project_id"],
					), nil
				},
				ImportStateVerifyIgnore: []string{
					"master_order_id", "project_id",
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["eip_id"],
					), nil
				},
				ImportStateVerifyIgnore: []string{
					"master_order_id", "project_id", "region_id",
				},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, eipId, instId),
				Destroy: true,
			},
		},
	})
}
