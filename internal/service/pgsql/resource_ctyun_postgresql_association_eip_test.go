package pgsql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunPgsqlAssociationEip(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_association_eip." + rnd
	resourceFile := "resource_ctyun_postgresql_association_eip.tf"

	specsDatasourceName := "data.ctyun_postgresql_specs." + dnd
	specsDatasourceFile := "datasource_ctyun_postgresql_specs.tf"

	eipId := dependence.eipID
	instId := dependence.pgsqlID
	//instId := "1bb2c455994c419ca0acadbc436c44c8"

	//prodType := "1"
	//prodCode := "POSTGRESQL"
	instanceType := "S"

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
				Config: utils.LoadTestCase(specsDatasourceFile, dnd, instanceType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(specsDatasourceName, "specs.#", utils.AtLeastOne),
				),
			},
			// 绑定eip
			{
				Config: utils.LoadTestCase(resourceFile, rnd, eipId, instId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "eip_id", eipId),
					resource.TestCheckResourceAttr(resourceName, "eip_status", "1"),
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
			// 解绑
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, eipId, instId),
				Destroy: true,
			},
		},
	})
}
