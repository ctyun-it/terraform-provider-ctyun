package ebm_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunEbmDeviceTypes(t *testing.T) {
	dnd := utils.GenerateRandomString()
	datasourceFile := "datasource_ctyun_ebm_device_types.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{

				Config: utils.LoadTestCase(datasourceFile, dnd, dependence.deviceType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith("data.ctyun_ebm_device_types."+dnd, "device_types.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith("data.ctyun_ebm_device_images."+dnd, "images.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith("data.ctyun_ebm_device_raids."+dnd, "raids.#", utils.AtLeastOne),
				),
			},
		},
	})
}
