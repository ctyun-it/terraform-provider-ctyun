package crs_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunCrsOpensourceImages(t *testing.T) {
	dnd := utils.GenerateRandomString()
	datasourceName := "data.ctyun_crs_opensource_images." + dnd
	datasourceFile := "datasource_ctyun_crs_opensource_images.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{

				Config: utils.LoadTestCase(datasourceFile, dnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(datasourceName, "images.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrSet(datasourceName, "images.0.image_url"),
					resource.TestCheckResourceAttrSet(datasourceName, "images.0.image_url_internal"),
				),
			},
		},
	})
}
