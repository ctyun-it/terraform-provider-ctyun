package sdwan_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunSdwanAclRule_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_sdwan_acl_rule." + rnd
	resourceFile := "resource_ctyun_sdwan_acl_rule.tf"
	datasourceName := "data.ctyun_sdwan_acl_rules." + rnd
	datasourceFile := "datasource_ctyun_sdwan_acl_rules.tf"

	// 先创建一个ACL用于测试
	aclRnd := utils.GenerateRandomString()
	aclName := utils.GenerateRandomString()
	aclId := dependence.AclId

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
			// 创建SD-WAN ACL Rule测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, aclId, "in", "tcp", "IPv4", "192.168.0.0/24", "80-443", "50", "allow", "192.168.0.0/24", "80-443"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "direction", "in"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "dst_cidr", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "dst_port_range", "80-443"),
					resource.TestCheckResourceAttr(resourceName, "priority", "50"),
					resource.TestCheckResourceAttr(resourceName, "action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "src_cidr", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "src_port_range", "80-443"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 更新SD-WAN ACL Rule测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, aclId, "out", "udp", "IPv4", "172.16.0.0/16", "50-443", "60", "deny", "172.16.0.0/16", "50-443"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "direction", "out"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "udp"),
					resource.TestCheckResourceAttr(resourceName, "dst_cidr", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "dst_port_range", "50-443"),
					resource.TestCheckResourceAttr(resourceName, "priority", "60"),
					resource.TestCheckResourceAttr(resourceName, "action", "deny"),
					resource.TestCheckResourceAttr(resourceName, "src_cidr", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "src_port_range", "50-443"),
				),
			},
			// 导入测试
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["acl_id"],
						rs.Primary.ID,
					), nil
				},
				ImportStateVerifyIgnore: []string{"ip_version"}, // 项目ID可能变化

			},
			// datasource 测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, aclId, "out", "tcp", "IPv4", "192.168.0.0/24", "80-443", "50", "deny", "10.0.0.0/24", "80-443") + "\n" +
					utils.LoadTestCase(datasourceFile, rnd, aclId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
				),
			},
			{
				Config:  utils.LoadTestCase("resource_ctyun_sdwan_acl.tf", aclRnd, aclName, "in", "udp", "IPv4", "10.0.0.0/16", "-1/-1", 100, "allow", "10.0.0.0/16", "-1/-1"),
				Destroy: true,
			},
		},
	})
}
