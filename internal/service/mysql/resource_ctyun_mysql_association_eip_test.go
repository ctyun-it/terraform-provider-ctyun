package mysql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunMysqlAssociationEip(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_mysql_association_eip." + rnd
	resourceFile := "resource_ctyun_mysql_association_eip.tf"
	eipId := dependence.eipID
	eipAddress := dependence.eipAddress
	instId := dependence.mysqlID

	instance_series := "S"

	specDatasourceName := "data.ctyun_mysql_specs." + dnd
	specDatasourceFile := "datasource_ctyun_mysql_specs.tf"
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
			// datasource
			{
				Config: utils.LoadTestCase(specDatasourceFile, dnd, instance_series),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(specDatasourceName, "specs.#", "8"),
				),
			},
			// eip 绑定
			{
				Config: utils.LoadTestCase(resourceFile, rnd, eipId, eipAddress, instId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "eip_id", eipId),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instId),
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
					"master_order_id",
					"project_id",
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
				Config:  utils.LoadTestCase(resourceFile, rnd, eipId, eipAddress, instId),
				Destroy: true,
			},
		},
	})
}
