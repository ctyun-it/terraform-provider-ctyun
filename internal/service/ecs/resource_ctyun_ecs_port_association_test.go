package ecs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunEcsPortAssociation_all(t *testing.T) {
	rnd := utils.GenerateRandomString()
	name := "ctyun_ecs_port_association." + rnd
	configFile := "resource_ctyun_ecs_port_association.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				// 测试基本创建场景
				Config: utils.LoadTestCase(configFile, rnd, dependence.instanceID, dependence.ecsPortForAssociationId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "region_id"),
					resource.TestCheckResourceAttr(name, "instance_id", dependence.instanceID),
				),
			},
			{
				// 测试更新场景
				Config: utils.LoadTestCase(configFile, rnd, dependence.instanceID, dependence.ecsPortForAssociationId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "region_id"),
					resource.TestCheckResourceAttr(name, "instance_id", dependence.instanceID),
					resource.TestCheckResourceAttrSet(name, "port_id"),
				),
			},
			{
				// 测试导入功能 (始终使用完整的三参数格式)
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[name]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", name)
					}
					return fmt.Sprintf("%s,%s,%s",
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["port_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerifyIgnore: []string{
					"az_name",
					"project_id",
				},
			},
			{
				// 测试销毁解绑场景
				Config:  utils.LoadTestCase(configFile, rnd, dependence.instanceID, dependence.ecsPortForAssociationId),
				Destroy: true,
			},
		},
	})
}
