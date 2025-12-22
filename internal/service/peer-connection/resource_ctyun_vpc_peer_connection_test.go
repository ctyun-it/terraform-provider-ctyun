package peer_connection_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunVpcPeerConnection_Basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_vpc_peer_connection." + rnd
	resourceFile := "resource_ctyun_vpc_peer_connection.tf"

	datasourceName := "data.ctyun_vpc_peer_connections." + dnd
	datasourceFile := "datasource_vpc_peer_connections.tf"

	// 配置测试环境需要的动态值
	requestVpcID := dependence.vpcID1
	acceptVpcID := dependence.vpcID2
	name := "peer-conn-" + utils.GenerateRandomString()
	description := "test basic peer connection"
	projectID := "0"

	updatedName := "peer-conn-updated-" + utils.GenerateRandomString()
	updatedDescription := "updated test peer connection"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// 确保两个VPC存在且不在同一个CIDR段
			fmt.Printf("Using VPCs: request=%s, accept=%s\n", requestVpcID, acceptVpcID)
		},
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 基础创建测试（同一个租户）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, projectID, name, description, requestVpcID, acceptVpcID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "request_vpc_id", requestVpcID),
					resource.TestCheckResourceAttr(resourceName, "accept_vpc_id", acceptVpcID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(resourceName, "request_vpc_name"),
					resource.TestCheckResourceAttrSet(resourceName, "request_vpc_cidr"),
					resource.TestCheckResourceAttrSet(resourceName, "accept_vpc_name"),
					resource.TestCheckResourceAttrSet(resourceName, "accept_vpc_cidr"),
				),
			},
			// 2. 资源更新测试（更新名称和描述）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, projectID, updatedName, updatedDescription, requestVpcID, acceptVpcID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "request_vpc_id", requestVpcID),
					resource.TestCheckResourceAttr(resourceName, "accept_vpc_id", acceptVpcID),
				),
			},
			// 3. datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, projectID, updatedName, updatedDescription, requestVpcID, acceptVpcID) +
					utils.LoadTestCase(datasourceFile, dnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "peer_connections.#")),
			},
			// 4. import state验证
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
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accept_email", "description"},
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
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["project_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accept_email", "description"},
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
						rs.Primary.Attributes["instance_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accept_email", "description", "project_id"},
			},
			// 5. 销毁资源
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, projectID, updatedName, updatedDescription, requestVpcID, acceptVpcID),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunVpcPeerConnection_WithTags(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_vpc_peer_connection." + rnd
	resourceFile := "resource_ctyun_vpc_peer_connection_tags.tf"

	// 配置测试环境需要的动态值
	requestVpcID := dependence.vpcID1
	acceptVpcID := dependence.vpcID2
	name := "peer-conn-tags-" + utils.GenerateRandomString()
	description := "test peer connection with tags"
	projectID := "0"

	// 标签配置
	tags1 := `[
		{"key": "environment", "value": "test"},
		{"key": "project", "value": "terraform"}
	]`
	tags2 := `[
		{"key": "environment", "value": "production"},
		{"key": "owner", "value": "devops"},
		{"key": "version", "value": "v1.0"}
	]`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			fmt.Printf("Testing VPC peer connection with tags\n")
		},
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建带标签的对等连接
			{
				Config: utils.LoadTestCase(resourceFile, rnd, projectID, name, description, requestVpcID, acceptVpcID, tags1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "environment"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.1.key", "project"),
					resource.TestCheckResourceAttr(resourceName, "tags.1.value", "terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.1.id"),
				),
			},
			// 2. 更新标签（增删改）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, projectID, name, description, requestVpcID, acceptVpcID, tags2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", map[string]string{
						"key":   "environment",
						"value": "production",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", map[string]string{
						"key":   "owner",
						"value": "devops",
					}),

					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", map[string]string{
						"key":   "version",
						"value": "v1.0",
					}),
				),
			},
			// 3. 移除所有标签
			{
				Config: utils.LoadTestCase(resourceFile, rnd, projectID, name, description, requestVpcID, acceptVpcID, "[]"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
				),
			},
			// 4. 销毁资源
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, projectID, name, description, requestVpcID, acceptVpcID, "[]"),
				Destroy: true,
			},
		},
	})
}

//func TestAccCtyunVpcPeerConnection_CrossAccount(t *testing.T) {
//	rnd := utils.GenerateRandomString()
//	resourceName := "ctyun_vpc_peer_connection." + rnd
//	resourceFile := "resource_ctyun_vpc_peer_connection_cross_account.tf"
//
//	attchResourceName := "ctyun_vpc_peer_connection_attach." + rnd
//	attchResourceFile := "resource_ctyun_vpc_peer_connection_attach.tf"
//
//	// 配置测试环境需要的动态值
//	requestVpcID := dependence.vpcID1
//	acceptVpcID := dependence.crossAccountVpcID // 跨账号VPC
//	acceptEmail := dependence.crossAccountEmail // 对端账号邮箱
//	name := "peer-conn-cross-" + utils.GenerateRandomString()
//	description := "test cross account peer connection"
//	projectID := "0"
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			fmt.Printf("Testing cross-account VPC peer connection\n")
//		},
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			// 1. 创建跨账号对等连接
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, projectID, name, description, requestVpcID, acceptVpcID, acceptEmail),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttr(resourceName, "name", name),
//					resource.TestCheckResourceAttr(resourceName, "accept_email", acceptEmail),
//					resource.TestCheckResourceAttr(resourceName, "request_vpc_id", requestVpcID),
//					resource.TestCheckResourceAttr(resourceName, "accept_vpc_id", acceptVpcID),
//				),
//			},
//			// 2. 更新名称和描述
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, projectID, name+"-new", description+" updated", requestVpcID, acceptVpcID, acceptEmail),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "name", name+"-new"),
//					resource.TestCheckResourceAttr(resourceName, "description", description+" updated"),
//				),
//			},
//			// 3. 验证对等连接状态为已接受
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, projectID, name+"-new", description+" updated", requestVpcID, acceptVpcID, acceptEmail) +
//					utils.LoadTestCase(attchResourceFile, rnd, fmt.Sprintf("%s.id", resourceName), "enable"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					//resource.TestCheckResourceAttr(resourceName, "status", "agree"),
//					resource.TestCheckResourceAttrSet(attchResourceName, "id")),
//			},
//			// 4. 销毁资源
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, projectID, name+"-new", description+" updated", requestVpcID, acceptVpcID, acceptEmail) +
//					utils.LoadTestCase(attchResourceFile, rnd, fmt.Sprintf("%s.id", resourceName), "enable"),
//				Destroy: true,
//			},
//		},
//	})
//}
