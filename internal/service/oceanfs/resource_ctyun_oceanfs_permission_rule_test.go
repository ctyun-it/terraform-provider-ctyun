package oceanfs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunOceanfsPermissionRule(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_oceanfs_permission_rule." + rnd
	resourceFile := "resource_ctyun_oceanfs_permission_rule.tf"

	// 测试数据
	permissionGroupFuid := dependence.permissionGroupID
	initialAuthAddr := "192.168.1.0/24"
	updatedAuthAddr := "10.0.0.0/16"
	//ipv6AuthAddr := "2001:db8::/32"
	initialRwPermission := "ro"
	updatedRwPermission := "rw"
	initialPriority := int32(1)
	updatedPriority := int32(100)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建权限规则测试（只读权限，默认优先级）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					permissionGroupFuid, initialAuthAddr,
					initialRwPermission, initialPriority,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					// 基本属性验证
					resource.TestCheckResourceAttr(resourceName, "permission_group_id", permissionGroupFuid),
					resource.TestCheckResourceAttr(resourceName, "auth_addr", initialAuthAddr),
					resource.TestCheckResourceAttr(resourceName, "rw_permission", initialRwPermission),
					resource.TestCheckResourceAttr(resourceName, "priority", fmt.Sprintf("%d", initialPriority)),

					// 系统生成属性验证
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 更新权限规则测试（读写权限，更高优先级）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, permissionGroupFuid, updatedAuthAddr,
					updatedRwPermission, updatedPriority,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "auth_addr", updatedAuthAddr),
					resource.TestCheckResourceAttr(resourceName, "rw_permission", updatedRwPermission),
					resource.TestCheckResourceAttr(resourceName, "priority", fmt.Sprintf("%d", updatedPriority)),

					// 验证ID保持不变
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 3. 导入测试 1
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s,%s",
						rs.Primary.Attributes["id"],
						rs.Primary.Attributes["permission_group_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			// 3. 导入测试 1
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
						rs.Primary.Attributes["permission_group_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, permissionGroupFuid, updatedAuthAddr,
					updatedRwPermission, updatedPriority,
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例2: IPv6地址测试
func TestAccCtyunOceanfsPermissionRuleIPv6(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_oceanfs_permission_rule." + rnd
	resourceFile := "resource_ctyun_oceanfs_permission_rule.tf"

	permissionGroupFuid := dependence.permissionGroupID

	// 测试数据
	ipv6AuthAddr := "2001:0db8:0000:0000:0000:ff00:0042:8329"
	rwPermission := "rw"
	priority := int32(50)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建IPv6权限规则测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					permissionGroupFuid, ipv6AuthAddr,
					rwPermission, priority,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "auth_addr", ipv6AuthAddr),
					resource.TestCheckResourceAttr(resourceName, "rw_permission", rwPermission),
					resource.TestCheckResourceAttr(resourceName, "priority", fmt.Sprintf("%d", priority)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					permissionGroupFuid, ipv6AuthAddr,
					rwPermission, priority,
				),
				Destroy: true,
			},
		},
	})
}
