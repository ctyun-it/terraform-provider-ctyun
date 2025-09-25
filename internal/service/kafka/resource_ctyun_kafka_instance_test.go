package kafka_test

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

func TestAccCtyunKafkaInstanceCluster(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_kafka_instance." + rnd
	datasourceName := "data.ctyun_kafka_instances." + dnd
	resourceFile := "resource_ctyun_kafka_instance.tf"
	datasourceFile := "datasource_ctyun_kafka_instances.tf"

	engineVersion := "3.6"
	zone := os.Getenv("CTYUN_AZ_NAME")

	initName := "tf-kafka-init-" + utils.GenerateRandomString()
	initNodeNum := 3
	initDiskSize := 100
	initRetentionHours := 80

	updatedName := "tf-kafka-updated-" + utils.GenerateRandomString()
	updatedNodeNum := 5
	updatedDiskSize := 200
	updatedRetentionHours := 60

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
					engineVersion,
					dependence.kafkaClusterSpecName,
					initNodeNum,
					zone,
					dependence.kafkaClusterDiskType,
					initDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					initRetentionHours,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", initName),
					resource.TestCheckResourceAttr(resourceName, "engine_version", engineVersion),
					resource.TestCheckResourceAttr(resourceName, "spec_name", dependence.kafkaClusterSpecName),
					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(initNodeNum)),
					resource.TestCheckTypeSetElemAttr(resourceName, "zone_list.*", zone),
					resource.TestCheckResourceAttr(resourceName, "disk_type", dependence.kafkaClusterDiskType),
					resource.TestCheckResourceAttr(resourceName, "disk_size", strconv.Itoa(initDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "retention_hours", strconv.Itoa(initRetentionHours)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			// 更新属性
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					engineVersion,
					dependence.kafkaClusterSpecName2,
					updatedNodeNum,
					zone,
					dependence.kafkaClusterDiskType,
					updatedDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					updatedRetentionHours,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "engine_version", engineVersion),
					resource.TestCheckResourceAttr(resourceName, "spec_name", dependence.kafkaClusterSpecName2),
					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(updatedNodeNum)),
					resource.TestCheckTypeSetElemAttr(resourceName, "zone_list.*", zone),
					resource.TestCheckResourceAttr(resourceName, "disk_type", dependence.kafkaClusterDiskType),
					resource.TestCheckResourceAttr(resourceName, "disk_size", strconv.Itoa(updatedDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "retention_hours", strconv.Itoa(updatedRetentionHours)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},

			// 缩容
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					engineVersion,
					dependence.kafkaClusterSpecName,
					updatedNodeNum,
					zone,
					dependence.kafkaClusterDiskType,
					updatedDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					updatedRetentionHours,
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".id",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.instance_name", updatedName),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.engine_version", engineVersion),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.spec_name", dependence.kafkaClusterSpecName),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.node_num", strconv.Itoa(updatedNodeNum)),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.disk_type", dependence.kafkaClusterDiskType),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.disk_size", strconv.Itoa(updatedDiskSize)),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.subnet_id", dependence.subnetID),
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
					"auto_renew_cycle_count",
					"auto_renew",
					"cycle_count",
					"cycle_type",
					"project_id",
					"security_group_id",
					"master_order_id",
					"zone_list",
					"plain_port",
					"http_port",
					"ssl_port",
					"sasl_port",
				},
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					engineVersion,
					dependence.kafkaClusterSpecName,
					updatedNodeNum,
					zone,
					dependence.kafkaClusterDiskType,
					updatedDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					updatedRetentionHours,
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".id",
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunKafkaInstanceSingle(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_kafka_instance." + rnd
	resourceFile := "resource_ctyun_kafka_instance_on_demand.tf"

	engineVersion := "3.6"
	zone := os.Getenv("CTYUN_AZ_NAME")

	initName := "tf-kafka-init-" + utils.GenerateRandomString()
	initNodeNum := 1
	initDiskSize := 100
	initRetentionHours := 80

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
					engineVersion,
					dependence.kafkaSingleSpecName,
					initNodeNum,
					zone,
					dependence.kafkaSingleDiskType,
					initDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					initRetentionHours,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", initName),
					resource.TestCheckResourceAttr(resourceName, "engine_version", engineVersion),
					resource.TestCheckResourceAttr(resourceName, "spec_name", dependence.kafkaSingleSpecName),
					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(initNodeNum)),
					resource.TestCheckTypeSetElemAttr(resourceName, "zone_list.*", zone),
					resource.TestCheckResourceAttr(resourceName, "disk_type", dependence.kafkaClusterDiskType),
					resource.TestCheckResourceAttr(resourceName, "disk_size", strconv.Itoa(initDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "retention_hours", strconv.Itoa(initRetentionHours)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			// 更新属性
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					engineVersion,
					dependence.kafkaSingleSpecName2,
					initNodeNum,
					zone,
					dependence.kafkaSingleDiskType,
					initDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					initRetentionHours,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", initName),
					resource.TestCheckResourceAttr(resourceName, "engine_version", engineVersion),
					resource.TestCheckResourceAttr(resourceName, "spec_name", dependence.kafkaSingleSpecName2),
					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(initNodeNum)),
					resource.TestCheckTypeSetElemAttr(resourceName, "zone_list.*", zone),
					resource.TestCheckResourceAttr(resourceName, "disk_type", dependence.kafkaClusterDiskType),
					resource.TestCheckResourceAttr(resourceName, "disk_size", strconv.Itoa(initDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "retention_hours", strconv.Itoa(initRetentionHours)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					engineVersion,
					dependence.kafkaSingleSpecName2,
					initNodeNum,
					zone,
					dependence.kafkaSingleDiskType,
					initDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					initRetentionHours,
				),
				Destroy: true,
			},
		},
	})
}
