package vpc_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCtyunNetResourcesByTag_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCtyunNetResourcesByTagConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "region_id"),
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "total_count"),
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "current_count"),
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "total_page"),
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "resources.#"),
				),
			},
		},
	})
}

func TestAccCtyunNetResourcesByTag_withLabelKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCtyunNetResourcesByTagConfigWithLabelKey(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "region_id"),
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "total_count"),
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "current_count"),
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "total_page"),
				),
			},
			{
				Config: testAccCtyunNetResourcesByTagConfigWithLabelValue(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "region_id"),
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "total_count"),
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "current_count"),
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "total_page"),
				),
			},
		},
	})
}

func TestAccCtyunNetResourcesByTag_withPagination(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCtyunNetResourcesByTagConfigWithPagination(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "region_id"),
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "total_count"),
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "current_count"),
					resource.TestCheckResourceAttrSet("data.ctyun_net_resources_by_tag.test", "total_page"),
					resource.TestCheckResourceAttr("data.ctyun_net_resources_by_tag.test", "page_number", "1"),
					resource.TestCheckResourceAttr("data.ctyun_net_resources_by_tag.test", "page_size", "10"),
				),
			},
		},
	})
}

func testAccCtyunNetResourcesByTagConfig() string {
	return `
data "ctyun_net_resources_by_tag" "test" {
  label_id= "1c24ddb1ff534de9a4bcd13c3b680b59"
}

output "resources_count" {
  value = length(data.ctyun_net_resources_by_tag.test.resources)
}

output "total_count" {
  value = data.ctyun_net_resources_by_tag.test.total_count
}

output "resources" {
  value = data.ctyun_net_resources_by_tag.test.resources
}
`
}

func testAccCtyunNetResourcesByTagConfigWithLabelKey() string {
	return `
data "ctyun_net_resources_by_tag" "test" {
  label_key = "key"
}
`
}
func testAccCtyunNetResourcesByTagConfigWithLabelValue() string {
	return `
data "ctyun_net_resources_by_tag" "test" {
  label_value = "value"
}
`
}

func testAccCtyunNetResourcesByTagConfigWithPagination() string {
	return `
data "ctyun_net_resources_by_tag" "test" {
  label_value = "value"
  page_number = 1
  page_size   = 10
}
`
}
