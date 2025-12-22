package acl_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunSubnetAssociationAcl(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_subnet_association_acl." + rnd
	resourceFile := "resource_ctyun_subnet_association_acl.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	subnetID := dependence.subnetID
	aclID1 := dependence.aclID
	aclID2 := dependence.aclID2

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建关联测试（绑定到ACL1）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					aclID1, subnetID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "acl_id", aclID1),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 更新关联测试（绑定到ACL2）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					aclID2, subnetID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "acl_id", aclID2),
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
						rs.Primary.Attributes["id"],
						rs.Primary.Attributes["subnet_id"],
						rs.Primary.Attributes["acl_id"],
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{}, // 不需要忽略任何字段
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s,%s,%s",
						rs.Primary.Attributes["id"],
						rs.Primary.Attributes["subnet_id"], // 假设资源代码需要name
						rs.Primary.Attributes["acl_id"],    // 假设资源代码需要vpc_id
						rs.Primary.Attributes["project_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{}, // 不需要忽略任何字段
			},
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
						rs.Primary.Attributes["subnet_id"],
						rs.Primary.Attributes["acl_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
			// 4. 清理资源（解绑）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					aclID2, subnetID,
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例2: 多个子网关联测试
func TestAccCtyunSubnetAssociationAclMultiple(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName1 := "ctyun_subnet_association_acl." + rnd + "_1"
	resourceName2 := "ctyun_subnet_association_acl." + rnd + "_2"
	resourceFile := "resource_ctyun_subnet_association_acl.tf"

	projectID := "0"
	subnetID1 := dependence.subnetID
	subnetID2 := dependence.subnetID2
	aclID := dependence.aclID

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建多个关联测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd+"_1", projectID,
					aclID, subnetID1,
				) + utils.LoadTestCase(resourceFile, rnd+"_2", projectID,
					aclID, subnetID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "subnet_id", subnetID1),
					resource.TestCheckResourceAttr(resourceName1, "acl_id", aclID),
					resource.TestCheckResourceAttrSet(resourceName1, "id"),
					resource.TestCheckResourceAttr(resourceName2, "subnet_id", subnetID2),
					resource.TestCheckResourceAttr(resourceName2, "acl_id", aclID),
					resource.TestCheckResourceAttrSet(resourceName2, "id"),
				),
			},
			// 2. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd+"_1", projectID,
					aclID, subnetID1,
				) + utils.LoadTestCase(resourceFile, rnd+"_2", projectID,
					aclID, subnetID2),
				Destroy: true,
			},
		},
	})
}
