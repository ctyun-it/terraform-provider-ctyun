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

	datasourceName := "data.ctyun_mongodb_association_eips." + dnd
	datasourceFile := "datasource_ctyun_mysql_association_eips.tf"
	eipId := dependence.eipID
	//eipId := "eip-140rfs2and"
	//eipAddress := "150.223.193.123"
	instId := dependence.mongodbID
	hostIp := dependence.hostIP
	//instId := "c1ef217509294dc6a4d6b6ec24a46586"
	//hostIp := "192.168.128.3"

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
				Config: utils.LoadTestCase(resourceFile, rnd, eipId, instId, hostIp),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "eip_id", eipId),
					resource.TestCheckResourceAttr(resourceName, "inst_id", instId),
				),
			},
			//datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, eipId, instId, hostIp) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf(`eip_id="%s"`, eipId)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "eips.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "eips.0.bind_status", "1"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, eipId, instId, hostIp) +
					utils.LoadTestCase(specDatasourceFile, dnd, instanceType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(specDatasourceName, "specs.#"),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, eipId, instId, hostIp),
				Destroy: true,
			},
		},
	})
}
