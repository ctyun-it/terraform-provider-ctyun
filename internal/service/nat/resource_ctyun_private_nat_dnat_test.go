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

func TestAccCtyunPrivateDNat(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_private_nat_dnat." + rnd
	datasourceName := "data.ctyun_private_nat_dnats." + dnd

	resourceFile := "resource_ctyun_private_nat_dnat.tf"
	datasourceFile := "datasource_ctyun_private_nat_dnats.tf"

	natGatewayId := dependence.privateNatID
	internalPort := utils.GenerateRandomPort(0, 65535)
	updatedInternalPort := utils.GenerateRandomPort(0, 65535)
	externalPort := utils.GenerateRandomPort(0, 1024)
	updatedExternalPort := utils.GenerateRandomPort(0, 1024)

	internalIp := "192.168.128.3"
	updatedInternalIp := "192.168.128.6"
	externalIp := dependence.privateNatIP1
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
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, protocol, externalPort, internalPort, internalIp),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_port", strconv.Itoa(internalPort)),
					resource.TestCheckResourceAttr(resourceName, "external_port", strconv.Itoa(externalPort)),
					resource.TestCheckResourceAttr(resourceName, "internal_ip", internalIp),
					resource.TestCheckResourceAttr(resourceName, "external_ip", externalIp),
					resource.TestCheckResourceAttr(resourceName, "protocol", protocol),
					resource.TestCheckResourceAttrSet(resourceName, "dnat_id"),
				),
			},
			{
				//	2 resource update验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, updatedProtocol, updatedExternalPort, updatedInternalPort, updatedInternalIp),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_port", strconv.Itoa(updatedInternalPort)),
					resource.TestCheckResourceAttr(resourceName, "external_port", strconv.Itoa(updatedExternalPort)),
					resource.TestCheckResourceAttr(resourceName, "internal_ip", updatedInternalIp),
					resource.TestCheckResourceAttr(resourceName, "external_ip", externalIp),
					resource.TestCheckResourceAttr(resourceName, "protocol", updatedProtocol),
					resource.TestCheckResourceAttrSet(resourceName, "dnat_id"),
				),
			},
			{
				// 3 datasource 验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, updatedProtocol, updatedExternalPort, updatedInternalPort, updatedInternalIp) +
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
				Config:  utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, updatedProtocol, updatedExternalPort, updatedInternalPort, updatedInternalIp),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunPrivateDNat2(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_private_nat_dnat." + rnd
	datasourceName := "data.ctyun_private_nat_dnats." + dnd

	resourceFile := "resource_ctyun_private_nat_dnat2.tf"
	datasourceFile := "datasource_ctyun_private_nat_dnats.tf"
	externalIp := dependence.privateNatIP1
	natGatewayId := dependence.privateNatID
	internalPort := utils.GenerateRandomPort(0, 65535)
	updatedInternalPort := utils.GenerateRandomPort(0, 65535)
	externalPort := utils.GenerateRandomPort(0, 1024)
	updatedExternalPort := utils.GenerateRandomPort(0, 1024)

	protocol := "tcp"
	updatedProtocol := "udp"
	description := "测试"

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
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, protocol, externalPort, internalPort, dependence.portId, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_port", strconv.Itoa(internalPort)),
					resource.TestCheckResourceAttr(resourceName, "external_port", strconv.Itoa(externalPort)),
					resource.TestCheckResourceAttr(resourceName, "port_id", dependence.portId),
					resource.TestCheckResourceAttr(resourceName, "external_ip", externalIp),
					resource.TestCheckResourceAttr(resourceName, "protocol", protocol),
					resource.TestCheckResourceAttrSet(resourceName, "dnat_id"),
				),
			},
			// 1resource create/ delete 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, updatedProtocol, updatedExternalPort, updatedInternalPort, dependence.portId, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_port", strconv.Itoa(updatedInternalPort)),
					resource.TestCheckResourceAttr(resourceName, "external_port", strconv.Itoa(updatedExternalPort)),
					resource.TestCheckResourceAttr(resourceName, "port_id", dependence.portId),
					resource.TestCheckResourceAttr(resourceName, "external_ip", externalIp),
					resource.TestCheckResourceAttr(resourceName, "protocol", updatedProtocol),
					resource.TestCheckResourceAttrSet(resourceName, "dnat_id"),
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
					natGatewayID := ds.Attributes["nat_gateway_id"]
					projectId := ds.Attributes["project_id"]
					return fmt.Sprintf("%s,%s,%s,%s", id, natGatewayID, projectId, regionId), nil
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					natGatewayID := ds.Attributes["nat_gateway_id"]
					return fmt.Sprintf("%s,%s", id, natGatewayID), nil
				},
			},
			{
				// 3 datasource 验证
				Config: utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, updatedProtocol, updatedExternalPort, updatedInternalPort, dependence.portId, description) +
					utils.LoadTestCase(datasourceFile, dnd, natGatewayId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "dnats.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.internal_port", strconv.Itoa(updatedInternalPort)),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.external_port", strconv.Itoa(updatedExternalPort)),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.protocol", updatedProtocol),
					resource.TestCheckResourceAttr(datasourceName, "dnats.0.port_id", dependence.portId),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, natGatewayId, externalIp, updatedProtocol, updatedExternalPort, updatedInternalPort, dependence.portId, description),
				Destroy: true,
			},
		},
	})
}
