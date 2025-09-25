package elb_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunElbAcl(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_elb_acl." + rnd
	resourceFile := "resource_ctyun_elb_acl.tf"

	datasourceName := "data.ctyun_elb_acls." + dnd
	datasourceFile := "datasource_ctyun_elb_acls.tf"

	name := "acl_" + utils.GenerateRandomString()
	sourceIps := `"127.0.0.1/32","192.168.0.0/16"`

	updatedName := "acl_new_" + utils.GenerateRandomString()
	updatedSourceIps := `"127.0.0.1/32","192.168.0.0/16","192.168.10.0"`
	updatedSourceIps2 := `"192.168.0.0/16","192.168.10.0"`

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
			// 1.1 Create验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, sourceIps),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "source_ips.#", "2"),
				),
			},
			// 1.2 update验证, name更新，sourceIps增加(增加的是ip地址)
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedSourceIps),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "source_ips.#", "3"),
				),
			},
			// 1.3 update验证, name更新，sourceIps减少
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedSourceIps2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "source_ips.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "source_ips.0", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "source_ips.1", "192.168.10.0"),
				),
			},
			// 1.4 datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedSourceIps2) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("ids=%s.id", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "acls.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "acls.0.name", updatedName),
				),
			},
			// destroy
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, updatedName, updatedSourceIps2),
				Destroy: true,
			},
		},
	})
}
