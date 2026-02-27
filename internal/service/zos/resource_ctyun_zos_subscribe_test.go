package zos_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCtyunZosSubscribeTest(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_zos_subscribe." + rnd
	resourceFile := "resource_ctyun_zos_subscribe.tf"

	regionID := "200000004421"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建按需计费只读实例测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, regionID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
				),
			},

			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, regionID,
				),
				Destroy: true,
			},
		},
	})
}
