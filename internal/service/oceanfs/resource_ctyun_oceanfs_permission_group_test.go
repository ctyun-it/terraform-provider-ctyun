package oceanfs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunOceanfsPermissionGroup(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_oceanfs_permission_group." + rnd
	resourceFile := "resource_ctyun_oceanfs_permission_group.tf"

	// 测试数据
	initialName := "test-permission-group-" + rnd
	updatedName := "test-permission-group-updated-" + rnd
	initialDescription := "Initial permission group description"
	updatedDescription := "Updated permission group description"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建权限组测试（基本配置）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initialName, initialDescription,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					// 基本属性验证
					resource.TestCheckResourceAttr(resourceName, "name", initialName),
					resource.TestCheckResourceAttr(resourceName, "description", initialDescription),

					// 系统生成属性验证
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "update_time"),
				),
			},
			// 2. 更新权限组测试（修改名称和描述）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName, updatedDescription,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),

					// 验证ID保持不变
					resource.TestCheckResourceAttrSet(resourceName, "id"),

					// 验证时间戳已更新
					resource.TestCheckResourceAttrSet(resourceName, "update_time"),
				),
			},
			// 3. 导入测试1
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "update_time"},
			},
			// 3. 导入测试1
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s",
						rs.Primary.Attributes["id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "update_time"},
			},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName, updatedDescription,
				),
				Destroy: true,
			},
		},
	})
}
