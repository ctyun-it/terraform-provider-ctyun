package vpc_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunNetworkInterfaces_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dataSourceName := "data.ctyun_ports." + rnd + "_filtered"
	dataSourceFile := "data_ctyun_network_interfaces.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// 测试基本查询
				Config: utils.LoadTestCase(
					dataSourceFile, rnd, dependence.vpcID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "region_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_interfaces.#"),
				),
			},
		},
	})
}
