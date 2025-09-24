package ecs_test

import (
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunEcsFlavors_basic(t *testing.T) {
	rName := utils.GenerateRandomString()
	dataSourceName := "data.ctyun_ecs_flavors." + rName
	resourceFile := "datasource_ctyun_ecs_flavors.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCtyunEcsFlavorsConfig_basic(resourceFile, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.cpu"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.ram"),
				),
			},
		},
	})
}

func testAccCtyunEcsFlavorsConfig_basic(file, name string) string {
	return utils.LoadTestCase(file, name)
}

func TestAccCtyunEcsFlavors_byType(t *testing.T) {
	dnd := utils.GenerateRandomString()
	datasourceFile := "datasource_ctyun_ecs_flavors_type.tf"
	cpu := 2
	ram := 4
	arch := "x86"
	series := "C"
	cpuType := "CPU_C7"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(datasourceFile, dnd, cpu, ram, arch, series, cpuType),
			},
		},
	})
}
