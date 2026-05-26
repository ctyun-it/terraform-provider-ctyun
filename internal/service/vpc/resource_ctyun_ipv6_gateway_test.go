package vpc_test

import (
	"fmt"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunIPv6GatewayTest(t *testing.T) {
	rnd := utils.GenerateRandomString()

	// 资源名称
	resourceName := "ctyun_ipv6_gateway." + rnd
	// 测试配置文件
	resourceFile := "resource_ipv6_gateway.tf"

	vpcID := dependence.vpcID2

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
				Config: utils.LoadTestCase(resourceFile, rnd, vpcID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
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
			// Step5: 销毁资源
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, vpcID),
				Destroy: true,
			},
		},
	})
}
