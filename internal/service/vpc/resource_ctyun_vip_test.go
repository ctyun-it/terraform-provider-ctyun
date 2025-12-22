package vpc_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunVip_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_vip." + rnd
	resourceFile := "resource_ctyun_vip.tf"

	// 测试参数
	subnetId := dependence.subnetID
	vpcId := dependence.vpcID
	ipAddress := "192.168.2.101"
	vipType := "v4"

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
			{
				// 测试创建带所有参数的HaVip
				Config: utils.LoadTestCase(resourceFile, rnd, subnetId, vpcId, ipAddress, vipType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetId),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
					resource.TestCheckResourceAttr(resourceName, "vip_type", vipType),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv4_address"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					if id == "" || regionId == "" {
						return "", fmt.Errorf("id or region_id is required")
					}
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ip_address",
					"ipv6_address",
					"vip_type",
					"project_id",
				},
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					return fmt.Sprintf("%s", id), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ip_address",
					"ipv6_address",
					"vip_type",
					"project_id",
				},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, subnetId, vpcId, ipAddress, vipType),
				Destroy: true,
			},
		},
	})
}
