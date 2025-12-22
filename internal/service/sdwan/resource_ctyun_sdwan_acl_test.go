package sdwan_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunSdwanAcl_basic(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_sdwan_acl." + rnd
	resourceFile := "resource_ctyun_sdwan_acl.tf"

	name := utils.GenerateRandomString()
	nameUpdate := utils.GenerateRandomString()

	datasourceName := "data.ctyun_sdwan_acls." + dnd
	datasourceFile := "datasource_ctyun_sdwan_acls.tf"
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
			// 创建SD-WAN ACL测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, "in", "udp", "IPv4", "10.0.0.0/16", "-1/-1", 100, "allow", "10.0.0.0/16", "-1/-1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 更新SD-WAN ACL测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, nameUpdate, "out", "tcp", "IPv4", "192.168.0.0/24", "80-443", 50, "deny", "10.0.0.0/24", "80-443"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdate),
				),
			},
			// 导入测试
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id", "rules"}, // 项目ID可能变化

			},
			// datasource 测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, nameUpdate, "out", "tcp", "IPv4", "192.168.0.0/24", "80-443", 50, "deny", "10.0.0.0/24", "80-443") + "\n" +
					utils.LoadTestCase(datasourceFile, dnd, nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttrSet(datasourceName, "acls.#"),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, nameUpdate, "out", "tcp", "IPv4", "192.168.0.0/24", "80-443", 50, "deny", "10.0.0.0/24", "80-443"),
				Destroy: true,
			},
		},
	})
}
