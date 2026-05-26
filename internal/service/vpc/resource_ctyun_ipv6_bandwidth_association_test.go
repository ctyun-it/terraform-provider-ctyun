package vpc_test

import (
	"fmt"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunIPv6BandwidthAssociationTest(t *testing.T) {
	rnd := utils.GenerateRandomString()

	// 资源名称
	resourceName := "ctyun_ipv6_bandwidth_association." + rnd
	// 测试配置文件
	resourceFile := "resource_ipv6_bandwidth_association.tf"

	ipv6BandwidthID := dependence.ipv6BandwidthID
	vipIpv6Address := dependence.vipIpv6Address

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
				Config: utils.LoadTestCase(resourceFile, rnd, ipv6BandwidthID, vipIpv6Address),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ipv6_bandwidth_id", ipv6BandwidthID),
					resource.TestCheckResourceAttr(resourceName, "ipv6", vipIpv6Address),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
				),
			},
			// Step2: 导入测试（格式：id,region_id）
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					regionId := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s,%s", ipv6BandwidthID, vipIpv6Address, regionId), nil
				},
				ImportStateVerify: true,
			},
			// Step3: 导入测试（格式：仅id）
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return fmt.Sprintf("%s,%s", ipv6BandwidthID, vipIpv6Address), nil
				},
				ImportStateVerify: true,
			},
			// Step5: 销毁资源
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, ipv6BandwidthID, vipIpv6Address),
				Destroy: true,
			},
		},
	})
}
