package rabbitmq_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunRabbitmqInstanceCluster(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_rabbitmq_instance." + rnd
	datasourceName := "data.ctyun_rabbitmq_instances." + dnd
	resourceFile := "resource_ctyun_rabbitmq_instance.tf"
	datasourceFile := "datasource_ctyun_rabbitmq_instances.tf"

	zone := os.Getenv("CTYUN_AZ_NAME")
	nodeNum := 3
	diskSize := 100
	diskType := "SAS"

	initName := "tf-cluster-init-" + utils.GenerateRandomString()

	updatedName := "tf-cluster-updated-" + utils.GenerateRandomString()
	updatedNum := 5
	updatedDiskSize := 200

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
				// 创建
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					dependence.rabbitmqClusterSpecName,
					nodeNum,
					zone,
					diskSize,
					diskType,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", initName),
					resource.TestCheckResourceAttr(resourceName, "spec_name", dependence.rabbitmqClusterSpecName),
					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(nodeNum)),
					resource.TestCheckTypeSetElemAttr(resourceName, "zone_list.*", zone),
					resource.TestCheckResourceAttr(resourceName, "disk_type", diskType),
					resource.TestCheckResourceAttr(resourceName, "disk_size", strconv.Itoa(diskSize)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			// 更新属性
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					dependence.rabbitmqClusterSpecName2,
					updatedNum,
					zone,
					updatedDiskSize,
					diskType,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "spec_name", dependence.rabbitmqClusterSpecName2),
					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(updatedNum)),
					resource.TestCheckTypeSetElemAttr(resourceName, "zone_list.*", zone),
					resource.TestCheckResourceAttr(resourceName, "disk_type", diskType),
					resource.TestCheckResourceAttr(resourceName, "disk_size", strconv.Itoa(updatedDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},

			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					dependence.rabbitmqClusterSpecName2,
					updatedNum,
					zone,
					updatedDiskSize,
					diskType,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".id",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.instance_name", updatedName),
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
					"cycle_count",
					"cycle_type",
					"disk_size",
					"disk_type",
					"master_order_id",
					"node_num",
					"project_id",
					"security_group_id",
					"subnet_id",
					"vpc_id",
					"zone_list",
				},
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					dependence.rabbitmqClusterSpecName2,
					updatedNum,
					zone,
					updatedDiskSize,
					diskType,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".id",
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunRabbitmqInstanceClusterOnDemand(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_rabbitmq_instance." + rnd
	datasourceName := "data.ctyun_rabbitmq_instances." + dnd
	resourceFile := "resource_ctyun_rabbitmq_instance_on_demand.tf"
	datasourceFile := "datasource_ctyun_rabbitmq_instances.tf"

	zone := os.Getenv("CTYUN_AZ_NAME")
	nodeNum := 1
	diskSize := 100
	diskType := "SAS"

	initName := "tf-single-init-" + utils.GenerateRandomString()

	updatedName := "tf-single-updated-" + utils.GenerateRandomString()

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
				// 创建
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					dependence.rabbitmqSingleSpecName,
					nodeNum,
					zone,
					diskSize,
					diskType,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", initName),
					resource.TestCheckResourceAttr(resourceName, "spec_name", dependence.rabbitmqSingleSpecName),
					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(nodeNum)),
					resource.TestCheckTypeSetElemAttr(resourceName, "zone_list.*", zone),
					resource.TestCheckResourceAttr(resourceName, "disk_type", diskType),
					resource.TestCheckResourceAttr(resourceName, "disk_size", strconv.Itoa(diskSize)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			// 更新属性
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					dependence.rabbitmqSingleSpecName2,
					nodeNum,
					zone,
					diskSize,
					diskType,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "spec_name", dependence.rabbitmqSingleSpecName2),
					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(nodeNum)),
					resource.TestCheckTypeSetElemAttr(resourceName, "zone_list.*", zone),
					resource.TestCheckResourceAttr(resourceName, "disk_type", diskType),
					resource.TestCheckResourceAttr(resourceName, "disk_size", strconv.Itoa(diskSize)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},

			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					dependence.rabbitmqSingleSpecName2,
					nodeNum,
					zone,
					diskSize,
					diskType,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".id",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.instance_name", updatedName),
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
					"cycle_count",
					"cycle_type",
					"disk_size",
					"disk_type",
					"master_order_id",
					"node_num",
					"project_id",
					"security_group_id",
					"subnet_id",
					"vpc_id",
					"zone_list",
				},
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					dependence.rabbitmqSingleSpecName2,
					nodeNum,
					zone,
					diskSize,
					diskType,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".id",
				),
				Destroy: true,
			},
		},
	})
}
