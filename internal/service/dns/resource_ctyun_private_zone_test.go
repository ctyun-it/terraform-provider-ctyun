package dns_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

// 测试用例1: 基础创建测试
func TestAccCtyunPrivateZone_Basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_private_zone." + rnd
	resourceFile := "resource_ctyun_private_zone.tf"

	datasourceFile := "datasource_ctyun_private_zones.tf"
	datasourceName := "data.ctyun_private_zones." + dnd

	// 测试数据
	zoneName := "test-zone." + rnd + ".internal.com"
	description := "Test private zone"
	proxyPattern := "zone"
	vpcID := fmt.Sprintf(`"%s"`, dependence.vpcID) // 假设在依赖中定义了vpcID
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建基础Private Zone
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneName, description, proxyPattern,
					300, // ttl
					vpcID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "proxy_pattern", proxyPattern),
					resource.TestCheckResourceAttr(resourceName, "ttl", "300"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id_list.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "vpc_id_list.*", dependence.vpcID),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "update_time"),
				),
			},
			// 2. 更新测试 - 修改TTL和描述
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneName, "Updated description",
					"record", // 修改proxy_pattern
					600,      // 修改ttl
					vpcID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "proxy_pattern", "record"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "600"),
				),
			},
			// datasource验证
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneName, description, proxyPattern,
					300, // ttl
					vpcID) + utils.LoadTestCase(
					datasourceFile, dnd, fmt.Sprintf("%s.id", resourceName), zoneName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "zones.#", "1"),
				),
			},
			// 3. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneName, "Updated description",
					"record", // 修改proxy_pattern
					600,      // 修改ttl
					vpcID,
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例2: 测试带标签的Private Zone
func TestAccCtyunPrivateZone_WithTags(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_private_zone." + rnd
	resourceFile := "resource_ctyun_private_zone_tag.tf"

	zoneName := "test-zone." + rnd + ".internal.com"
	vpcID := dependence.vpcID
	vpcID1 := dependence.vpcID1
	vpcIds1 := fmt.Sprintf(`"%s"`, vpcID)
	vpcIds := fmt.Sprintf(`"%s","%s"`, vpcID, vpcID1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建带标签的Private Zone
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneName, "Zone with tags", "zone",
					300, vpcIds1,
					`[{
						key   = "environment"
						value = "test"
					},{
						key   = "project"
						value = "terraform"
					}]`,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
				),
			},
			// 2. 更新标签
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneName, "Zone with updated tags", "zone", 1000, vpcIds,
					`[{
						key   = "environment"
						value = "production"
					},{
						key   = "team"
						value = "devops"
					}]`,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "description", "Zone with updated tags"),
				),
			},
			// 3. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneName, "Zone with updated tags", "zone", 1000, vpcIds,
					`[{
						key   = "environment"
						value = "production"
					},{
						key   = "team"
						value = "devops"
					}]`,
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例3: 测试多VPC关联
func TestAccCtyunPrivateZone_MultipleVPCs(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceFile := "resource_ctyun_private_zone.tf"
	resourceName := "ctyun_private_zone." + rnd
	//resourceFile := "resource_ctyun_private_zone.tf"

	zoneName := "test-zone-" + rnd + ".internal.com"
	// 假设有多个VPC用于测试
	vpcID1 := dependence.vpcID
	vpcID2 := dependence.vpcID1 // 另一个VPC
	vpcID3 := dependence.vpcID2 // 另一个VPC
	vpcID4 := dependence.vpcID3 // 另一个VPC
	vpcID5 := dependence.vpcID4 // 另一个VPC
	vpcID6 := dependence.vpcID5 // 另一个VPC

	vpcIds := fmt.Sprintf(`"%s","%s","%s","%s"`, vpcID1, vpcID2, vpcID3, vpcID4)
	updatedVpcIds := fmt.Sprintf(`"%s","%s","%s","%s"`, vpcID1, vpcID2, vpcID3, vpcID6)
	updatedVpcIds1 := fmt.Sprintf(`"%s","%s"`, vpcID4, vpcID5)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建关联多个VPC的Private Zone
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, zoneName, "multiple vpc test", "zone", 300,
					vpcIds),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id_list.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceName, "vpc_id_list.*", vpcID1),
					resource.TestCheckTypeSetElemAttr(resourceName, "vpc_id_list.*", vpcID2),
				),
			},
			// 2. 更新VPC列表
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, zoneName, "multiple vpc test", "zone", 2147483647,
					updatedVpcIds),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vpc_id_list.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceName, "vpc_id_list.*", vpcID1),
					resource.TestCheckTypeSetElemAttr(resourceName, "vpc_id_list.*", vpcID2),
					resource.TestCheckTypeSetElemAttr(resourceName, "vpc_id_list.*", vpcID3),
					resource.TestCheckTypeSetElemAttr(resourceName, "vpc_id_list.*", vpcID6),
				),
			},
			// 2.1 更新VPC列表
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, zoneName, "multiple vpc test", "zone", 2147483647,
					updatedVpcIds1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vpc_id_list.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "vpc_id_list.*", vpcID4),
					resource.TestCheckTypeSetElemAttr(resourceName, "vpc_id_list.*", vpcID5),
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
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{}, // 可选忽略
			},
			// 3.1 只导入ID测试
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
				ImportStateVerifyIgnore: []string{"region_id"}, // 可选忽略
			},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, zoneName, "multiple vpc test", "zone", 2147483647,
					updatedVpcIds1),
				Destroy: true,
			},
		},
	})
}
