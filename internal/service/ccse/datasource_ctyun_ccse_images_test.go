package ccse_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunCcseImages(t *testing.T) {
	ecs := utils.GenerateRandomString() + "ecs"
	ebm := utils.GenerateRandomString() + "ebm"
	datasourceFile := "datasource_ctyun_ccse_images.tf"

	ecsData := "data.ctyun_ccse_images." + ecs
	ebmData := "data.ctyun_ccse_images." + ebm

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(datasourceFile, ecs, dependence.flavorName, ebm, dependence.deviceType),
				Check: resource.ComposeAggregateTestCheckFunc(
					// ECS 数据源检查
					resource.TestCheckResourceAttrWith(ecsData, "images.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrSet(ecsData, "images.0.id"),
					resource.TestCheckResourceAttrSet(ecsData, "images.0.name"),

					// EBM 数据源检查
					resource.TestCheckResourceAttrWith(ebmData, "images.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrSet(ebmData, "images.0.id"),
					resource.TestCheckResourceAttrSet(ebmData, "images.0.name"),
				),
			},
		},
	})
}
