package sfs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunSfsPermissionGroupRule(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_sfs_permission_rule." + rnd
	resourceFile := "resource_sfs_permission_rule.tf"

	datasourceName := "data.ctyun_sfs_permission_rules." + dnd
	datasourceFile := "datasource_ctyun_permission_rules1.tf"

	// 从环境变量获取测试依赖资源
	permissionGroupFuid := dependence.sfsPermissionGroupID

	// 基础配置参数
	authAddr := "192.168.0.0/24"
	rwPermission := "rw"
	priority := 100

	// 更新配置参数
	updatedAuthAddr := "192.168.1.0/24"
	updatedRwPermission := "ro"
	updatedPriority := 200

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
				Config: utils.LoadTestCase(resourceFile, rnd, permissionGroupFuid, authAddr, rwPermission, priority),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "permission_group_fuid", permissionGroupFuid),
					resource.TestCheckResourceAttr(resourceName, "auth_addr", authAddr),
					resource.TestCheckResourceAttr(resourceName, "rw_permission", rwPermission),
					resource.TestCheckResourceAttr(resourceName, "permission_rule_priority", fmt.Sprintf("%d", priority)),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
				),
			},
			// 2. 资源更新测试（授权地址、读写权限、优先级）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, permissionGroupFuid, updatedAuthAddr, updatedRwPermission, updatedPriority),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "permission_group_fuid", permissionGroupFuid),
					resource.TestCheckResourceAttr(resourceName, "auth_addr", updatedAuthAddr),
					resource.TestCheckResourceAttr(resourceName, "rw_permission", updatedRwPermission),
					resource.TestCheckResourceAttr(resourceName, "permission_rule_priority", fmt.Sprintf("%d", updatedPriority)),
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
					return fmt.Sprintf("%s,%s,%s",
						rs.Primary.ID,
						rs.Primary.Attributes["region_id"],
						rs.Primary.Attributes["permission_group_fuid"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{}, // 不需要忽略任何字段
			},
			// 4. 验证datasource
			{
				Config: utils.LoadTestCase(resourceFile, rnd, permissionGroupFuid, updatedAuthAddr, updatedRwPermission, updatedPriority) +
					utils.LoadTestCase(datasourceFile, dnd, permissionGroupFuid),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "permission_group_fuid", permissionGroupFuid)),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, permissionGroupFuid, updatedAuthAddr, updatedRwPermission, updatedPriority),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunSfsPermissionGroupRuleIPv6ToIPv4(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_sfs_permission_rule." + rnd
	resourceFile := "resource_sfs_permission_rule.tf"

	datasourceName := "data.ctyun_sfs_permission_rules." + dnd
	datasourceFile := "datasource_ctyun_permission_rules.tf"

	// 从环境变量获取测试依赖资源
	permissionGroupFuid := dependence.sfsPermissionGroupID

	// IPv6格式的初始配置参数
	ipv6AuthAddr := "2001:db8::/32"
	rwPermission := "rw"
	priority := 2

	// IPv4格式的更新配置参数
	ipv4AuthAddr := "192.168.1.0/24"
	updatedRwPermission := "ro"
	updatedPriority := 1

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
			// 1. 基础创建测试（IPv6地址）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, permissionGroupFuid, ipv6AuthAddr, rwPermission, priority),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "permission_group_fuid", permissionGroupFuid),
					resource.TestCheckResourceAttr(resourceName, "auth_addr", ipv6AuthAddr),
					resource.TestCheckResourceAttr(resourceName, "rw_permission", rwPermission),
					resource.TestCheckResourceAttr(resourceName, "permission_rule_priority", fmt.Sprintf("%d", priority)),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
				),
			},
			// 2. 资源更新测试（IPv4地址）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, permissionGroupFuid, ipv4AuthAddr, updatedRwPermission, updatedPriority),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "permission_group_fuid", permissionGroupFuid),
					resource.TestCheckResourceAttr(resourceName, "auth_addr", ipv4AuthAddr),
					resource.TestCheckResourceAttr(resourceName, "rw_permission", updatedRwPermission),
					resource.TestCheckResourceAttr(resourceName, "permission_rule_priority", fmt.Sprintf("%d", updatedPriority)),
				),
			},
			// 3. 验证datasource
			{
				Config: utils.LoadTestCase(resourceFile, rnd, permissionGroupFuid, ipv4AuthAddr, updatedRwPermission, updatedPriority) +
					utils.LoadTestCase(datasourceFile, dnd, permissionGroupFuid, fmt.Sprintf("%s.id", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "permission_rules.0.auth_addr", ipv4AuthAddr),
					resource.TestCheckResourceAttr(datasourceName, "permission_rules.0.permission_group_fuid", permissionGroupFuid),
					resource.TestCheckResourceAttr(datasourceName, "permission_rules.0.rw_permission", updatedRwPermission),
					resource.TestCheckResourceAttr(datasourceName, "permission_rules.0.permission_rule_priority", fmt.Sprintf("%d", updatedPriority)),
					resource.TestCheckResourceAttr(datasourceName, "permission_rules.0.user_permission", "no_root_squash"),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, permissionGroupFuid, ipv4AuthAddr, updatedRwPermission, updatedPriority),
				Destroy: true,
			},
		},
	})
}
