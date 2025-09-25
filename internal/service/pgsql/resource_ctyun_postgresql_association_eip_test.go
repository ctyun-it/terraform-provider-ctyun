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

	//
	resourceFile := "resource_ctyun_postgresql_association_eip.tf"

	datasourceName := "data.ctyun_mysql_association_eips." + dnd
	datasourceFile := "datasource_ctyun_pgsql_association_eips.tf"

	specsDatasourceName := "data.ctyun_postgresql_specs." + dnd
	specsDatasourceFile := "datasource_ctyun_postgresql_specs.tf"

	eipId := dependence.eipID
	instId := dependence.PgsqlID
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
			// 绑定eip
			{
				Config: utils.LoadTestCase(resourceFile, rnd, eipId, instId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "eip_id", eipId),
					resource.TestCheckResourceAttr(resourceName, "eip_status", "1"),
				),
			},
			// resource验证
			//datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, eipId, instId) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf(`eip_id="%s"`, eipId)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "eips.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "eips.0.bind_status", "1"),
				),
			},
			// spec datasource验证
			{
				Config: utils.LoadTestCase(specsDatasourceFile, dnd, instanceType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(specsDatasourceName, "specs.#", utils.AtLeastOne),
				),
			},
			// 解绑
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, eipId, instId),
				Destroy: true,
			},
		},
	})
}
