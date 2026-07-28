package vpc_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunIPv6Bandwidth(t *testing.T) {
	rnd := utils.GenerateRandomString()

	// 资源名称
	resourceName := "ctyun_ipv6_bandwidth." + rnd
	// 测试配置文件
	resourceFile := "resource_ctyun_ipv6_bandwidth.tf"

	// 初始配置
	initName := "ipv6_init"
	initBandwidth := "1"
	// 更新配置
	updatedName := "ipv6_updated"
	updatedBandwidth := "5"

	resource.Test(t, resource.TestCase{
		// 资源销毁校验：测试完成后确保资源已删除
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("ipv6 bandwidth resource destroy failed")
			}
			return nil
		},
		// 加载测试Provider
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Step1: 创建IPv6带宽 + 校验初始属性
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, initBandwidth),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth", initBandwidth),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
				),
			},
			// Step2: 导入测试（格式：id,region_id）
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerify: true,
			},
			// Step3: 导入测试（格式：仅id）
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					return fmt.Sprintf("%s", id), nil
				},
				ImportStateVerify: true,
			},
			// Step4: 更新名称和带宽 + 校验更新结果
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedBandwidth),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth", updatedBandwidth),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Step5: 销毁资源
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, updatedName, updatedBandwidth),
				Destroy: true,
			},
		},
	})
}
