package acl_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunAclRule(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_acl_rule." + rnd
	resourceFile := "resource_ctyun_acl_rule.tf"

	projectID := ""
	aclID := dependence.aclID

	// 测试数据
	initialSourceIP := "192.168.1.0/24"
	initialDestIP := "10.0.0.0/16"
	updatedSourceIP := "172.16.0.0/16"
	updatedDestIP := "192.168.0.0/24"
	initialPortRange := "8080:8085"
	updatedPortRange := "9000:9010"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建入站规则测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					aclID, "ingress",
					"tcp", "ipv4",
					initialPortRange, initialPortRange,
					initialSourceIP, initialDestIP,
					"accept", true, "Ingress rule description",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "acl_id", aclID),
					resource.TestCheckResourceAttr(resourceName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "source_port", initialPortRange),
					resource.TestCheckResourceAttr(resourceName, "destination_port", initialPortRange),
					resource.TestCheckResourceAttr(resourceName, "source_ip_address", initialSourceIP),
					resource.TestCheckResourceAttr(resourceName, "destination_ip_address", initialDestIP),
					resource.TestCheckResourceAttr(resourceName, "action", "accept"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "Ingress rule description"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 更新入站规则测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					aclID, "ingress",
					"tcp", "ipv4",
					updatedPortRange, updatedPortRange,
					updatedSourceIP, updatedDestIP,
					"drop", false, "Updated ingress rule",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_port", updatedPortRange),
					resource.TestCheckResourceAttr(resourceName, "destination_port", updatedPortRange),
					resource.TestCheckResourceAttr(resourceName, "source_ip_address", updatedSourceIP),
					resource.TestCheckResourceAttr(resourceName, "destination_ip_address", updatedDestIP),
					resource.TestCheckResourceAttr(resourceName, "action", "drop"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated ingress rule"),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					aclID, "ingress",
					"tcp", "ipv4",
					updatedPortRange, updatedPortRange,
					updatedSourceIP, updatedDestIP,
					"drop", false, "Updated ingress rule",
				),
				Destroy: true,
			},
			// 3. 创建出站规则测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					aclID, "egress",
					"udp", "ipv4",
					initialPortRange, initialPortRange,
					initialSourceIP, initialDestIP,
					"accept", true, "Egress rule description",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "direction", "egress"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "udp"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 4. 更新出站规则测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					aclID, "egress",
					"udp", "ipv4",
					updatedPortRange, updatedPortRange,
					updatedSourceIP, updatedDestIP,
					"drop", false, "Updated egress rule",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_port", updatedPortRange),
					resource.TestCheckResourceAttr(resourceName, "destination_port", updatedPortRange),
					resource.TestCheckResourceAttr(resourceName, "source_ip_address", updatedSourceIP),
					resource.TestCheckResourceAttr(resourceName, "destination_ip_address", updatedDestIP),
					resource.TestCheckResourceAttr(resourceName, "action", "drop"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated egress rule"),
				),
			},
			// 5. 导入测试
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
						rs.Primary.Attributes["direction"],
						rs.Primary.Attributes["acl_id"],
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"}, // 可选忽略
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
						rs.Primary.Attributes["direction"],
						rs.Primary.Attributes["acl_id"],
						rs.Primary.Attributes["project_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"}, // 可选忽略
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
						rs.Primary.Attributes["direction"],
						rs.Primary.Attributes["acl_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"}, // 可选忽略
			},
			// 6. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					aclID, "egress",
					"udp", "ipv4",
					updatedPortRange, updatedPortRange,
					updatedSourceIP, updatedDestIP,
					"drop", false, "Updated egress rule",
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例2: 测试IPv6规则
func TestAccCtyunAclRuleIPv6(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_acl_rule." + rnd
	resourceFile := "resource_ctyun_acl_rule.tf"

	datasourceName := "data.ctyun_acl_rules." + dnd
	datasourceFile := "datasource_acl_rules.tf"

	projectID := "0"
	aclID := dependence.aclID

	// 测试数据
	sourceIP := "2001:db8::/32"
	destIP := "2001:db8:1::/48"
	portRange := "8080:8085"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建IPv6入站规则
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					aclID, "ingress",
					"tcp", "ipv6",
					portRange, portRange,
					sourceIP, destIP,
					"accept", true, "IPv6 ingress rule",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip_version", "ipv6"),
					resource.TestCheckResourceAttr(resourceName, "source_ip_address", sourceIP),
					resource.TestCheckResourceAttr(resourceName, "destination_ip_address", destIP),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. datasource验证
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					aclID, "ingress",
					"tcp", "ipv6",
					portRange, portRange,
					sourceIP, destIP,
					"accept", true, "IPv6 ingress rule",
				) + utils.LoadTestCase(datasourceFile, dnd, aclID, projectID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "in_policy_id.#", "3"),
				),
			},
			// 3. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					aclID, "ingress",
					"tcp", "ipv6",
					portRange, portRange,
					sourceIP, destIP,
					"accept", true, "IPv6 ingress rule",
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例3: 测试不同协议规则
func TestAccCtyunAclRuleProtocols(t *testing.T) {

	rnd := utils.GenerateRandomString()
	rnd1 := utils.GenerateRandomString()
	resourceName := "ctyun_acl_rule." + rnd
	resourceFile := "resource_ctyun_acl_rule_priority.tf"
	resourceFile1 := "resource_ctyun_acl_rule_none_port.tf"

	projectID := "0"
	aclID := dependence.aclID

	// 测试数据
	sourceIP := "192.168.1.0/24"
	destIP := "10.0.0.0/16"
	portRange := "8080:8085"
	//"all", "icmp",
	// 测试所有支持的协议
	protocols := []string{"tcp", "udp"}
	protocols1 := []string{"all", "icmp"}
	for _, protocol := range protocols {
		t.Run(protocol, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					// 1. 创建规则
					{
						Config: utils.LoadTestCase(
							resourceFile, rnd, projectID,
							aclID, "ingress",
							protocol, "ipv4",
							portRange, portRange,
							sourceIP, destIP,
							"accept", true, "accept Protocol test: "+protocol, 1,
						) + utils.LoadTestCase(resourceFile, rnd1, projectID,
							aclID, "ingress",
							protocol, "ipv4",
							portRange, portRange,
							sourceIP, destIP,
							"drop", true, "drop Protocol test: "+protocol, 2,
						),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "protocol", protocol),
							resource.TestCheckResourceAttrSet(resourceName, "id"),
						),
					},
					// 2. 清理资源
					{
						Config: utils.LoadTestCase(
							resourceFile, rnd, projectID,
							aclID, "ingress",
							protocol, "ipv4",
							portRange, portRange,
							sourceIP, destIP,
							"accept", true, "accept Protocol test: "+protocol, 1,
						) + utils.LoadTestCase(resourceFile, rnd1, projectID,
							aclID, "ingress",
							protocol, "ipv4",
							portRange, portRange,
							sourceIP, destIP,
							"drop", true, "drop Protocol test: "+protocol, 2,
						),
						Destroy: true,
					},
				},
			})
		})
	}

	for _, protocol := range protocols1 {
		t.Run(protocol, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					// 1. 创建规则
					{
						Config: utils.LoadTestCase(
							resourceFile1, rnd, projectID,
							aclID, "ingress",
							protocol, "ipv4",
							sourceIP, destIP,
							"accept", true, "accept Protocol test: "+protocol, 1,
						) + utils.LoadTestCase(resourceFile1, rnd1, projectID,
							aclID, "ingress",
							protocol, "ipv4",
							sourceIP, destIP,
							"drop", true, "drop Protocol test: "+protocol, 2,
						),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "protocol", protocol),
							resource.TestCheckResourceAttrSet(resourceName, "id"),
						),
					},
					// 2. 清理资源
					{
						Config: utils.LoadTestCase(
							resourceFile1, rnd, projectID,
							aclID, "ingress",
							protocol, "ipv4",
							sourceIP, destIP,
							"accept", true, "accept Protocol test: "+protocol, 1,
						) + utils.LoadTestCase(resourceFile1, rnd1, projectID,
							aclID, "ingress",
							protocol, "ipv4",
							sourceIP, destIP,
							"drop", true, "drop Protocol test: "+protocol, 2,
						),
						Destroy: true,
					},
				},
			})
		})
	}
}
