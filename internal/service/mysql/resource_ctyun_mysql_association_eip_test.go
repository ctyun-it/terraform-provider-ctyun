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

	datasourceName := "data.ctyun_mysql_association_eips." + dnd
	datasourceFile := "datasource_ctyun_mysql_association_eips.tf"
	eipId := dependence.eipID
	//eipId := "eip-140rfs2and"
	eipAddress := dependence.eipAddress
	//eipAddress := "150.223.193.123"
	instId := dependence.mysqlID
	//instId := dependence.subnetID

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
			// 绑定IP验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, eipId, eipAddress, instId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "eip_id", eipId),
					resource.TestCheckResourceAttr(resourceName, "inst_id", instId),
				),
			},
			//datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, eipId, eipAddress, instId) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf(`eip_id="%s"`, eipId)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "eips.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "eips.0.bind_status", "1"),
					//resource.TestCheckResourceAttr(datasourceName, "eips.0.eip_id", eipId),
					//resource.TestCheckResourceAttr(datasourceName, "eips.0.eip", eipAddress),
				),
			},
			{
				Config: utils.LoadTestCase(specDatasourceFile, dnd, instance_series),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(specDatasourceName, "specs.#", "8"),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, eipId, eipAddress, instId),
				Destroy: true,
			},
		},
	})
}
