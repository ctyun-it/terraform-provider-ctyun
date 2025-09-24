package vpc_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunSecurityGroupRule(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_security_group_rule." + rnd
	resourceFile := "resource_ctyun_security_group_rule.tf"

	ingressDirection := "ingress"
	ingressAction := "drop"
	ingressProtocol := "any"
	ingressEtherType := "ipv4"

	egressDirection := "egress"
	egressAction := "accept"
	egressProtocol := "icmp"
	egressEtherType := "ipv4"

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
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.securityGroupID, ingressDirection, ingressAction, ingressProtocol, ingressEtherType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "direction", ingressDirection),
					resource.TestCheckResourceAttr(resourceName, "action", ingressAction),
					resource.TestCheckResourceAttr(resourceName, "protocol", ingressProtocol),
					resource.TestCheckResourceAttr(resourceName, "ether_type", ingressEtherType),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.securityGroupID, egressDirection, egressAction, egressProtocol, egressEtherType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "direction", egressDirection),
					resource.TestCheckResourceAttr(resourceName, "action", egressAction),
					resource.TestCheckResourceAttr(resourceName, "protocol", egressProtocol),
					resource.TestCheckResourceAttr(resourceName, "ether_type", egressEtherType),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					sgID := ds.Attributes["security_group_id"]
					return fmt.Sprintf("%s,%s,%s", id, sgID, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project_id",
				},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, dependence.securityGroupID, egressDirection, egressAction, egressProtocol, egressEtherType),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunSecurityGroupRuleAllField(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_security_group_rule." + rnd
	resourceFile := "resource_ctyun_security_group_rule_all_field.tf"

	ingressDirection := "ingress"
	ingressAction := "drop"
	ingressProtocol := "udp"
	ingressEtherType := "ipv4"
	ingressPriority := 10
	ingressRange := "200"
	ingressDestCidrIp := "192.168.0.0/24"
	ingressDescription := "first"
	ingressDescription2 := "first1"

	egressDirection := "egress"
	egressAction := "accept"
	egressProtocol := "tcp"
	egressEtherType := "ipv6"
	egressPriority := 100
	egressRange := "200-300"
	egressDestCidrIp := "::/0"
	egressDescription := "first2"
	egressDescription2 := "first3"

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
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.securityGroupID,
					ingressDirection, ingressAction, ingressProtocol, ingressEtherType,
					ingressPriority, ingressRange, ingressDestCidrIp, ingressDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "direction", ingressDirection),
					resource.TestCheckResourceAttr(resourceName, "action", ingressAction),
					resource.TestCheckResourceAttr(resourceName, "protocol", ingressProtocol),
					resource.TestCheckResourceAttr(resourceName, "ether_type", ingressEtherType),
					resource.TestCheckResourceAttr(resourceName, "protocol", ingressProtocol),
					resource.TestCheckResourceAttr(resourceName, "range", ingressRange),
					resource.TestCheckResourceAttr(resourceName, "dest_cidr_ip", ingressDestCidrIp),
					resource.TestCheckResourceAttr(resourceName, "description", ingressDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.securityGroupID,
					ingressDirection, ingressAction, ingressProtocol, ingressEtherType,
					ingressPriority, ingressRange, ingressDestCidrIp, ingressDescription2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "direction", ingressDirection),
					resource.TestCheckResourceAttr(resourceName, "action", ingressAction),
					resource.TestCheckResourceAttr(resourceName, "protocol", ingressProtocol),
					resource.TestCheckResourceAttr(resourceName, "ether_type", ingressEtherType),
					resource.TestCheckResourceAttr(resourceName, "protocol", ingressProtocol),
					resource.TestCheckResourceAttr(resourceName, "range", ingressRange),
					resource.TestCheckResourceAttr(resourceName, "dest_cidr_ip", ingressDestCidrIp),
					resource.TestCheckResourceAttr(resourceName, "description", ingressDescription2),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.securityGroupID,
					egressDirection, egressAction, egressProtocol, egressEtherType,
					egressPriority, egressRange, egressDestCidrIp, egressDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "direction", egressDirection),
					resource.TestCheckResourceAttr(resourceName, "action", egressAction),
					resource.TestCheckResourceAttr(resourceName, "protocol", egressProtocol),
					resource.TestCheckResourceAttr(resourceName, "ether_type", egressEtherType),
					resource.TestCheckResourceAttr(resourceName, "protocol", egressProtocol),
					resource.TestCheckResourceAttr(resourceName, "range", egressRange),
					resource.TestCheckResourceAttr(resourceName, "dest_cidr_ip", egressDestCidrIp),
					resource.TestCheckResourceAttr(resourceName, "description", egressDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.securityGroupID,
					egressDirection, egressAction, egressProtocol, egressEtherType,
					egressPriority, egressRange, egressDestCidrIp, egressDescription2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "direction", egressDirection),
					resource.TestCheckResourceAttr(resourceName, "action", egressAction),
					resource.TestCheckResourceAttr(resourceName, "protocol", egressProtocol),
					resource.TestCheckResourceAttr(resourceName, "ether_type", egressEtherType),
					resource.TestCheckResourceAttr(resourceName, "protocol", egressProtocol),
					resource.TestCheckResourceAttr(resourceName, "range", egressRange),
					resource.TestCheckResourceAttr(resourceName, "dest_cidr_ip", egressDestCidrIp),
					resource.TestCheckResourceAttr(resourceName, "description", egressDescription2),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					sgID := ds.Attributes["security_group_id"]
					return fmt.Sprintf("%s,%s,%s", id, sgID, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project_id",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.securityGroupID,
					egressDirection, egressAction, egressProtocol, egressEtherType,
					egressPriority, egressRange, egressDestCidrIp, egressDescription),
				Destroy: true,
			},
		},
	})
}
