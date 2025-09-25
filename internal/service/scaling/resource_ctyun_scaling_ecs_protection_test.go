package scaling_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunScalingEcsProtection(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_scaling_ecs_protection." + rnd
	resourceFile := "resource_ctyun_scaling_ecs_protection.tf"

	// 从环境变量获取测试依赖资源
	scalingGroupID := dependence.scalingGroupID
	instanceIDList := fmt.Sprintf(`["%s"]`, dependence.instanceUUID3)

	protectStatus := true

	updatedProtectStatus := false

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 基础创建测试（开启保护）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, scalingGroupID, instanceIDList, protectStatus),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "true"),
				),
			},
			// 2. 资源更新测试（关闭保护并添加实例）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, scalingGroupID, instanceIDList, updatedProtectStatus),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "protect_status", "false"),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, scalingGroupID, instanceIDList, updatedProtectStatus),
				Destroy: true,
			},
		},
	})
}
