package nat_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunSNat1(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	// 创建nat信息
	natGatewayID := dependence.natID
	// 创建snat信息
	resourceName := "ctyun_nat_snat." + rnd
	datasourceName := "data.ctyun_nat_snats." + dnd
	resourceFile := "resource_ctyun_nat_snat.tf"
	datasourceFile := "datasource_ctyun_snat.tf"

	initSourceCidr := `source_cidr="192.168.0.0/24"`
	updatedSourceCidr := fmt.Sprintf(`source_cidr="%s"`, "192.168.128.0/24")

	//natGateWayId := "natgw-asdsmh8scy"
	//var natGatewayId string
	snatIps := fmt.Sprintf(`["%s"]`, dependence.eipID)
	updatedSnatIps := fmt.Sprintf(`["%s","%s"]`, dependence.eipID, dependence.eipID1)

	//updateDescription := utils.GenerateRandomString()
	//var id string

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
				// 1. 创建nat,snat
				// 1.1 resource create验证:
				// subnetType = 0(自定义情况),sourceCIDR必传
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayID, initSourceCidr, snatIps, "我是一条description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_cidr", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "snat_ips.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "snat_id"),
					resource.TestCheckResourceAttr(resourceName, "description", "我是一条description"),
				),
			},
			{
				// 1.2 resource update source_cidr验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayID, updatedSourceCidr, updatedSnatIps, "我是一条description plus"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_cidr", "192.168.128.0/24"),
					resource.TestCheckResourceAttr(resourceName, "snat_ips.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "description", "我是一条description plus"),
				),
			},
			{
				// 1.3. datasource验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayID, updatedSourceCidr, updatedSnatIps, "我是一条description plus") +
					utils.LoadTestCase(datasourceFile, dnd, natGatewayID, fmt.Sprintf(`snat_id=%s.snat_id`, resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "snats.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "snats.0.nat_gateway_id", natGatewayID),
					resource.TestCheckResourceAttr(datasourceName, "snats.0.subnet_cidr", "192.168.128.0/24"),
					//resource.TestCheckResourceAttr(datasourceName, "snats.0.subnet_id", updatedSubnetId),
				),
			},
			// 1.4 资源销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, natGatewayID, updatedSourceCidr, updatedSnatIps, "我是一条description plus"),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunSNat2(t *testing.T) {
	rnd := utils.GenerateRandomString()

	// 创建nat信息
	natGatewayID := dependence.natID
	// 创建snat信息
	resourceName := "ctyun_nat_snat." + rnd
	resourceFile := "resource_ctyun_nat_snat.tf"
	sourceSubnetId := dependence.subnetID1
	updatedSubnetId := dependence.subnetID2
	tfSourceSubnetID := fmt.Sprintf(`source_subnet_id="%s"`, sourceSubnetId)
	updatedTfSourceSubnetID := fmt.Sprintf(`source_subnet_id="%s"`, updatedSubnetId)

	//natGateWayId := "natgw-asdsmh8scy"
	//var natGatewayId string
	snatIps := fmt.Sprintf(`["%s"]`, dependence.eipID)
	updatedSnatIps := fmt.Sprintf(`["%s","%s"]`, dependence.eipID, dependence.eipID1)

	//updateDescription := utils.GenerateRandomString()
	//var id string

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
			// 2 subnetType = 1(有vpcId 的子网情况),sourceSubnetId必传
			{
				// 2.1 resource create验证1:
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayID, tfSourceSubnetID, updatedSnatIps, "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", natGatewayID),
					resource.TestCheckResourceAttr(resourceName, "source_subnet_id", sourceSubnetId),
					resource.TestCheckResourceAttr(resourceName, "snat_ips.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "snat_id"),
				),
			},
			{
				// 2.2 resource update source_subnet_id验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayID, updatedTfSourceSubnetID, snatIps, "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_subnet_id", updatedSubnetId),
					resource.TestCheckResourceAttr(resourceName, "snat_ips.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerifyIgnore: []string{
					"nat_gateway_id",
				},
			},
			// 2.3destroy
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, natGatewayID, updatedTfSourceSubnetID, snatIps, "test"),
				Destroy: true,
			},
		},
	})
}
