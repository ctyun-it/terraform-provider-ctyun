package vpc_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunNetworkInterface_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_port." + rnd
	resourceFile := "resource_ctyun_network_interface.tf"

	name := "tf-port-test-" + utils.GenerateRandomString()
	updatedName := "tf-port-test-updated-" + utils.GenerateRandomString()

	description := "测试网络接口描述"
	updatedDescription := description + "-updated"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				// 测试创建
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					name,
					description,
					dependence.subnetID,
					dependence.securityGroupID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "port_id"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
				),
			},
			{
				// 测试更新
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					updatedDescription,
					dependence.subnetID,
					dependence.securityGroupID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
			{
				// 测试导入
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceName)
					}

					regionId := rs.Primary.Attributes["region_id"]
					portId := rs.Primary.Attributes["port_id"]
					if regionId == "" {
						return "", fmt.Errorf("region_id is not set")
					}

					return fmt.Sprintf("%s,%s", portId, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"subnet_id",
					"security_group_ids",
					"secondary_private_ip_count",
				},
			},
			{
				// 测试导入
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceName)
					}
					return fmt.Sprintf("%s", rs.Primary.Attributes["port_id"]), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"subnet_id",
					"security_group_ids",
					"secondary_private_ip_count",
				},
			},
			{
				// 测试销毁
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					updatedDescription,
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunNetworkInterface_case1(t *testing.T) {

	rnd := utils.GenerateRandomString()
	name := "ctyun_port." + rnd
	configFile := "resource_ctyun_network_interface_case1.tf"

	// 初始测试参数
	initialPortName := "test-port-" + rnd
	initialDescription := "test port description"
	subnetId := dependence.subnetID
	securityGroupId := dependence.securityGroupID
	secondaryIpCount := 1
	ipv6AddressCount := 0

	// 更新后的测试参数
	updatedPortName := "updated-port-" + rnd
	updatedDescription := "updated port description"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			// 创建测试
			{
				Config: utils.LoadTestCase(configFile, rnd, initialPortName, initialDescription, subnetId, securityGroupId, secondaryIpCount, ipv6AddressCount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "port_id"),
					resource.TestCheckResourceAttrSet(name, "mac_address"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttr(name, "name", initialPortName),
					resource.TestCheckResourceAttr(name, "description", initialDescription),
					resource.TestCheckResourceAttr(name, "subnet_id", subnetId),
					resource.TestCheckResourceAttr(name, "secondary_private_ip_count", "1"),
				),
			},
			// 更新测试
			{
				Config: utils.LoadTestCase(configFile, rnd, updatedPortName, updatedDescription, subnetId, securityGroupId, secondaryIpCount, ipv6AddressCount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "port_id"),
					resource.TestCheckResourceAttrSet(name, "mac_address"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttr(name, "name", updatedPortName),
					resource.TestCheckResourceAttr(name, "description", updatedDescription),
					resource.TestCheckResourceAttr(name, "subnet_id", subnetId),
				),
			},
			// 删除测试（通过Destroy步骤）
			{
				Config:  utils.LoadTestCase(configFile, rnd, updatedPortName, updatedDescription, subnetId, securityGroupId, secondaryIpCount, ipv6AddressCount),
				Destroy: true,
			},
		},
	})

}

func TestAccCtyunNetworkInterface_case2(t *testing.T) {

	rnd := utils.GenerateRandomString()
	name := "ctyun_port." + rnd
	configFile := "resource_ctyun_network_interface_case2.tf"

	// 初始测试参数
	initialPortName := "test-port-" + rnd
	initialDescription := "test port description"

	subnetId := dependence.subnetID
	primaryIp := "192.168.2.33" // 使用手动分配
	securityGroupId := dependence.securityGroupID
	secondaryIpCount := 1
	ipv6AddressCount := 0

	// 更新后的测试参数（只更新name和description字段）
	updatedPortName := "updated-port-" + rnd
	updatedDescription := "updated port description"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			// 创建测试 - 包含所有可选字段
			{
				Config: utils.LoadTestCase(configFile, rnd, initialPortName, initialDescription, subnetId, primaryIp, securityGroupId, secondaryIpCount, ipv6AddressCount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "port_id"),
					resource.TestCheckResourceAttrSet(name, "mac_address"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttr(name, "name", initialPortName),
					resource.TestCheckResourceAttr(name, "description", initialDescription),
					resource.TestCheckResourceAttr(name, "subnet_id", subnetId),
					resource.TestCheckResourceAttr(name, "secondary_private_ip_count", "1"),
					resource.TestCheckResourceAttr(name, "ipv6_address_count", "0"),
				),
			},
			// 更新测试 - 只更新name和description字段
			{
				Config: utils.LoadTestCase(configFile, rnd, updatedPortName, updatedDescription, subnetId, primaryIp, securityGroupId, secondaryIpCount, ipv6AddressCount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "port_id"),
					resource.TestCheckResourceAttrSet(name, "mac_address"),
					resource.TestCheckResourceAttrSet(name, "status"),
					resource.TestCheckResourceAttr(name, "name", updatedPortName),
					resource.TestCheckResourceAttr(name, "description", updatedDescription),
					resource.TestCheckResourceAttr(name, "subnet_id", subnetId),
				),
			},
			// 删除测试（通过Destroy步骤）
			{
				Config:  utils.LoadTestCase(configFile, rnd, updatedPortName, updatedDescription, subnetId, primaryIp, securityGroupId, secondaryIpCount, ipv6AddressCount),
				Destroy: true,
			},
		},
	})
}
