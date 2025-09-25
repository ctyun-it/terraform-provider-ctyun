package sfs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunSfsPermissionGroupAssociation(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_sfs_permission_group_association." + rnd
	resourceFile := "resource_ctyun_sfs_permission_group_association.tf"

	sfsUID := dependence.SfsUID
	vpcID1 := dependence.vpcID1
	permissionGroupFuid1 := dependence.sfsPermissionGroupID
	permissionGroupFuid2 := dependence.sfsPermissionGroupID1

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
			// 1. 基础创建测试（绑定第一个权限组）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, permissionGroupFuid1, sfsUID, vpcID1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "permission_group_fuid", permissionGroupFuid1),
					resource.TestCheckResourceAttr(resourceName, "sfs_uid", sfsUID),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID1),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_cidr"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_description"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_is_default"),
				),
			},
			// 2. 资源更新测试（更换为第二个权限组）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, permissionGroupFuid2, sfsUID, vpcID1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttr(resourceName, "permission_group_fuid", permissionGroupFuid2),
					resource.TestCheckResourceAttr(resourceName, "sfs_uid", sfsUID),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID1),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_cidr"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_description"),
					resource.TestCheckResourceAttrSet(resourceName, "permission_group_is_default"),
				),
			},
			// 3. 资源导入测试
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			// 4. 清理资源
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, permissionGroupFuid2, sfsUID, vpcID1),
				Destroy: true,
			},
		},
	})
}
