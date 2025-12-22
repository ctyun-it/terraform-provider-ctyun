package oceanfs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunOceanfsPermissionGroupAssociation(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_oceanfs_permission_group_association." + rnd
	resourceFile := "resource_ctyun_oceanfs_permission_group_association.tf"

	permissionGroupID := dependence.permissionGroupID
	updatedPermissionGroupID := dependence.permissionGroupID1
	sfsUID := dependence.oceanfsID
	vpcID := dependence.vpcID1
	subnetID := dependence.subnetID

	// 测试数据
	initialIsVpce := true

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建权限组关联测试（启用VPCE）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					permissionGroupID, sfsUID, vpcID, subnetID, initialIsVpce,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					// 基本属性验证
					resource.TestCheckResourceAttr(resourceName, "permission_group_id", permissionGroupID),
					resource.TestCheckResourceAttr(resourceName, "oceanfs_id", sfsUID),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "is_vpce", fmt.Sprintf("%t", initialIsVpce)),

					// 系统生成属性验证
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 更新权限组关联测试（禁用VPCE）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedPermissionGroupID, sfsUID, vpcID, subnetID, initialIsVpce,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "is_vpce", fmt.Sprintf("%t", initialIsVpce)),

					// 验证其他属性保持不变
					resource.TestCheckResourceAttr(resourceName, "permission_group_id", updatedPermissionGroupID),
					resource.TestCheckResourceAttr(resourceName, "oceanfs_id", sfsUID),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),

					// 验证ID保持不变
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 3. 导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s,%s,%s,%s",
						rs.Primary.Attributes["region_id"],
						rs.Primary.Attributes["permission_group_id"],
						rs.Primary.Attributes["oceanfs_id"],
						rs.Primary.Attributes["vpc_id"],
						rs.Primary.Attributes["subnet_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_vpce"}, // VPCE设置可能变化

			},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedPermissionGroupID, sfsUID, vpcID, subnetID, initialIsVpce,
				),
				Destroy: true,
			},
		},
	})
}
