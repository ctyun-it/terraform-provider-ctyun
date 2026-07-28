package vpc_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunVipsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccVipsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ctyun_vips.test", "vips.#"),
				),
			},
		},
	})
}

func testAccVipsDataSourceConfig() string {
	return `
data "ctyun_vips" "test" {
  page_no   = 1
  page_size = 10
}
`
}
