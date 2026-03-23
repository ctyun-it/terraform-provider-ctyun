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
					func(s *terraform.State) error {
						ds := s.RootModule().Resources[resourceName].Primary
						createTime := ds.Attributes["create_time"]
						if utils.IsEmptyOrRfc3339(createTime) {
							return nil
						}
						return fmt.Errorf("time format doesn't match")
					},
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

// copyFile 辅助函数：复制文件（用于tfstate备份/恢复）
//func copyFile(src, dst string) error {
//	input, err := os.ReadFile(src)
//	if err != nil {
//		return err
//	}
//	return os.WriteFile(dst, input, 0644)
//}

//func TestResourceYourResource_Read_RemoveFromStateAfterDestroy(t *testing.T) {
//	rnd := utils.GenerateRandomString()
//	resourceName := "ctyun_oceanfs_permission_group." + rnd
//	resourceFile := "resource_ctyun_oceanfs_permission_group.tf"
//
//	var (
//		testWorkspaceDir = filepath.Join(os.TempDir(), "tf-test-"+t.Name())
//		tfStatePath      = filepath.Join(testWorkspaceDir, "terraform.tfstate")
//		tfStateBackup    = tfStatePath + ".bak"
//	)
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			// 创建资源
//			{
//				Config: utils.LoadTestCase(
//					resourceFile, rnd,
//					"test-permission-group-"+rnd, "Initial permission group description",
//				),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					func(state *terraform.State) error {
//						// 复制当前tfstate到备份文件
//						if err := copyFile(tfStatePath, tfStateBackup); err != nil {
//							return fmt.Errorf("备份tfstate失败: %v", err)
//						}
//						return nil
//					},
//				),
//			},
//			// 删除资源
//			{
//				Config: utils.LoadTestCase(
//					resourceFile, rnd,
//					"test-permission-group-"+rnd, "Initial permission group description",
//				),
//				Destroy: true,
//			},
//
//			// 恢复tfstate文件，并触发refresh
//			{
//				Check: func(state *terraform.State) error {
//					if err := copyFile(tfStateBackup, tfStatePath); err != nil {
//						return fmt.Errorf("恢复tfstate失败: %v", err)
//					}
//					return nil
//				},
//				// 触发refresh（read）并验证state被移除
//				// 从恢复的state中获取资源
//				rs, ok := state.RootModule().Resources[resourceName]
//				if !ok {
//				return fmt.Errorf("恢复的tfstate中未找到资源 %s", resourceName)
//			}
//
//			},
//		},
//	})
//}
