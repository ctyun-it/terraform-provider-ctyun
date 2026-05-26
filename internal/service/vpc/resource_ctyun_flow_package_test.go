package vpc_test

import (
	"fmt"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccCtyunFlowPackage 共享流量包资源的 acceptance 测试
func TestAccCtyunFlowPackage(t *testing.T) {
	rnd := utils.GenerateRandomString()

	// 资源名称：对应Schema中定义的 ctyun_flow_packages
	resourceName := "ctyun_flow_package." + rnd
	// 测试配置文件名称
	resourceFile := "resource_ctyun_flow_package.tf"

	resource.Test(t, resource.TestCase{
		// 资源销毁检查函数
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("共享流量包资源销毁失败")
			}
			return nil
		},
		// 提供者工厂配置
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 步骤1：创建共享流量包资源
			{
				Config: utils.LoadTestCase(resourceFile, rnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cycle_type", "month"),
					resource.TestCheckResourceAttr(resourceName, "spec", "10"),
				),
			},
			// 步骤2：销毁资源
			{
				Config:  utils.LoadTestCase(resourceFile, rnd),
				Destroy: true,
			},
		},
	})
}
