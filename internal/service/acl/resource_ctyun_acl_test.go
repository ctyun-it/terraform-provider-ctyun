package acl_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunAcl(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_acl." + rnd
	resourceFile := "resource_ctyun_acl.tf"
	resourceFile1 := "resource_ctyun_acl_all.tf"
	projectID := ""
	vpcID := dependence.vpcID

	// 测试数据
	aclName := "test-acl-" + rnd
	initialDescription := "Initial ACL description"
	updatedDescription := "Updated ACL description"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建ACL测试（默认启用）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					vpcID, aclName,
					initialDescription,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "name", aclName),
					resource.TestCheckResourceAttr(resourceName, "description", initialDescription),
					resource.TestCheckResourceAttr(resourceName, "apply_to_public_lb", "false"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 更新ACL测试（修改名称、描述和启用状态）
			{
				Config: utils.LoadTestCase(
					resourceFile1, rnd, projectID,
					vpcID, aclName+"-updated",
					updatedDescription, false, false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", aclName+"-updated"),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "apply_to_public_lb", "false"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
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
					return fmt.Sprintf("%s,%s,%s",
						rs.Primary.Attributes["id"],
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"apply_to_public_lb", "enabled"}, // 可选忽略
			},
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
				ImportStateVerifyIgnore: []string{"apply_to_public_lb", "enabled", "project_id"}, // 可选忽略
			},
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
						rs.Primary.Attributes["project_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"apply_to_public_lb", "enabled"}, // 可选忽略
			},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile1, rnd, projectID,
					vpcID, aclName+"-updated",
					updatedDescription, true, false,
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例2：创建时禁用ACL
func TestAccCtyunAclDisabled(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_acl." + rnd
	resourceFile := "resource_ctyun_acl_all.tf"
	projectID := "0"
	vpcID := dependence.vpcID

	dataSourceName := "data.ctyun_acls." + dnd
	datasourceFile := "datasource_ctyun_acls.tf"

	// 测试数据
	aclName := "test-acl-disabled-" + rnd
	description := "Disabled ACL"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建禁用状态的ACL
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					vpcID, aclName,
					description, false, false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 启用acl
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					vpcID, aclName,
					description, false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},

			// 3. datasource验证
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					vpcID, aclName,
					description, false, true) +
					utils.LoadTestCase(
						datasourceFile, dnd, fmt.Sprintf("%s.id", resourceName), projectID, aclName, 1, 50),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "acls.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "acls.0.name", aclName),
					resource.TestCheckResourceAttr(dataSourceName, "acls.0.enabled", "true")),
			},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					vpcID, aclName,
					description, false, true),
				Destroy: true,
			},
		},
	})
}
