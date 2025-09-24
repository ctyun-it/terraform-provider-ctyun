package sfs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunSfsPermissionGroup(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_sfs_permission_group." + rnd
	resourceFile := "resource_ctyun_sfs_permission_group.tf"

	// 基础配置参数
	name := "permission-group-" + rnd
	description := "Test permission group created by Terraform"

	// 更新配置参数
	updatedName := "updated-" + name
	updatedDescription := "Updated description for permission group"

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
			// 1. 基础创建测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "sfs_count"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_rule_count"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_is_default"),
				),
			},
			// 2. 资源更新测试（名称和描述）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					// 验证其他属性保持不变
					resource.TestCheckResourceAttrSet(resourceName, "sfs_count"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_rule_count"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_is_default"),
				),
			},
			// 3. 资源导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					// 构造导入ID: "id,region_id"
					return fmt.Sprintf("%s,%s",
						rs.Primary.ID,
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{}, // 不需要忽略任何字段
			},
			// 4. 清理资源
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunSfsPermissionGroupNoneDesc(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_sfs_permission_group." + rnd
	resourceFile := "resource_ctyun_sfs_permission_group_none_desc.tf"

	// 基础配置参数
	name := "permission-group-" + rnd
	//description := "Test permission group created by Terraform"

	// 更新配置参数
	updatedName := "updated-" + name
	//updatedDescription := "Updated description for permission group"

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
			// 1. 基础创建测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "sfs_count"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_rule_count"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_is_default"),
				),
			},
			// 2. 资源更新测试（名称和描述）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					// 验证其他属性保持不变
					resource.TestCheckResourceAttrSet(resourceName, "sfs_count"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_rule_count"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_is_default"),
				),
			},
			// 3. 资源导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					// 构造导入ID: "id,region_id"
					return fmt.Sprintf("%s,%s",
						rs.Primary.ID,
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"}, // 不需要忽略任何字段
			},
			// 4. 清理资源
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, updatedName),
				Destroy: true,
			},
		},
	})
}
