package nat_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccCtyunDNat(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_nat_dnat." + rnd
	datasourceName := "data.ctyun_nat_dnats." + dnd

	resourceFile := "resource_ctyun_nat_dnat.tf"
	datasourceFile := "datasource_ctyun_nat_dnat.tf"

	natGatewayId := dependence.natID
	dnatType := "custom"
	internalPort := utils.GenerateRandomPort(0, 65535)
	updatedInternalPort := utils.GenerateRandomPort(0, 65535)
	externalPort := utils.GenerateRandomPort(0, 1024)
	updatedExternalPort := utils.GenerateRandomPort(0, 1024)

	internalIp := "127.0.0.1"
	updatedInternalIp := "127.0.0.2"

	protocol := "tcp"
	updatedProtocol := "udp"

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
			// 1resource create/ delete 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, dependence.eipID, protocol, externalPort, internalPort, dnatType, internalIp),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_port", strconv.Itoa(internalPort)),
					resource.TestCheckResourceAttr(resourceName, "external_port", strconv.Itoa(externalPort)),
					resource.TestCheckResourceAttr(resourceName, "internal_ip", internalIp),
					resource.TestCheckResourceAttr(resourceName, "external_id", dependence.eipID),
					resource.TestCheckResourceAttr(resourceName, "protocol", protocol),
					resource.TestCheckResourceAttr(resourceName, "dnat_type", dnatType),
					resource.TestCheckResourceAttrSet(resourceName, "dnat_id"),
				),
			},
			{
				//	2 resource update验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, dependence.eipID1, updatedProtocol, updatedExternalPort, updatedInternalPort, dnatType, updatedInternalIp),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_port", strconv.Itoa(updatedInternalPort)),
					resource.TestCheckResourceAttr(resourceName, "external_port", strconv.Itoa(updatedExternalPort)),
					resource.TestCheckResourceAttr(resourceName, "internal_ip", updatedInternalIp),
					resource.TestCheckResourceAttr(resourceName, "external_id", dependence.eipID1),
					resource.TestCheckResourceAttr(resourceName, "protocol", updatedProtocol),
					resource.TestCheckResourceAttr(resourceName, "dnat_type", dnatType),
					resource.TestCheckResourceAttrSet(resourceName, "dnat_id"),
				),
			},
			{
				// 3 datasource 验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, dependence.eipID1, updatedProtocol, updatedExternalPort, updatedInternalPort, dnatType, updatedInternalIp) +
					utils.LoadTestCase(datasourceFile, dnd, natGatewayId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "dnats.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.internal_port", strconv.Itoa(updatedInternalPort)),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.external_port", strconv.Itoa(updatedExternalPort)),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.protocol", updatedProtocol),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.internal_ip", updatedInternalIp),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					ngID := ds.Attributes["nat_gateway_id"]
					regionId := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s,%s", id, ngID, regionId), nil
				},
				ImportStateVerifyIgnore: []string{
					"dnat_type",
					"internal_ip",
				},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, natGatewayId, dependence.eipID1, updatedProtocol, updatedExternalPort, updatedInternalPort, dnatType, updatedInternalIp),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunDNat2(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_nat_dnat." + rnd
	datasourceName := "data.ctyun_nat_dnats." + dnd

	resourceFile := "resource_ctyun_nat_dnat2.tf"
	datasourceFile := "datasource_ctyun_nat_dnat.tf"

	natGatewayId := dependence.natID
	dnatType := "instance"
	serverType := "VM"
	internalPort := utils.GenerateRandomPort(0, 65535)
	updatedInternalPort := utils.GenerateRandomPort(0, 65535)
	externalPort := utils.GenerateRandomPort(0, 1024)
	updatedExternalPort := utils.GenerateRandomPort(0, 1024)

	protocol := "tcp"
	updatedProtocol := "udp"

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
			// 1resource create/ delete 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, dependence.eipID, protocol, externalPort, internalPort, dnatType, serverType, dependence.ecsID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_port", strconv.Itoa(internalPort)),
					resource.TestCheckResourceAttr(resourceName, "external_port", strconv.Itoa(externalPort)),
					resource.TestCheckResourceAttr(resourceName, "instance_id", dependence.ecsID),
					resource.TestCheckResourceAttr(resourceName, "external_id", dependence.eipID),
					resource.TestCheckResourceAttr(resourceName, "protocol", protocol),
					resource.TestCheckResourceAttr(resourceName, "dnat_type", dnatType),
					resource.TestCheckResourceAttrSet(resourceName, "dnat_id"),
				),
			},
			// 1resource create/ delete 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, dependence.eipID1, updatedProtocol, updatedExternalPort, updatedInternalPort, dnatType, serverType, dependence.ecsID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_port", strconv.Itoa(updatedInternalPort)),
					resource.TestCheckResourceAttr(resourceName, "external_port", strconv.Itoa(updatedExternalPort)),
					resource.TestCheckResourceAttr(resourceName, "instance_id", dependence.ecsID),
					resource.TestCheckResourceAttr(resourceName, "external_id", dependence.eipID1),
					resource.TestCheckResourceAttr(resourceName, "protocol", updatedProtocol),
					resource.TestCheckResourceAttr(resourceName, "dnat_type", dnatType),
					resource.TestCheckResourceAttrSet(resourceName, "dnat_id"),
				),
			},
			{
				// 3 datasource 验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, dependence.eipID1, updatedProtocol, updatedExternalPort, updatedInternalPort, dnatType, serverType, dependence.ecsID) +
					utils.LoadTestCase(datasourceFile, dnd, natGatewayId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "dnats.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.internal_port", strconv.Itoa(updatedInternalPort)),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.external_port", strconv.Itoa(updatedExternalPort)),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.protocol", updatedProtocol),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.instance_id", dependence.ecsID),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, natGatewayId, dependence.eipID1, updatedProtocol, updatedExternalPort, updatedInternalPort, dnatType, serverType, dependence.ecsID),
				Destroy: true,
			},
		},
	})
}
