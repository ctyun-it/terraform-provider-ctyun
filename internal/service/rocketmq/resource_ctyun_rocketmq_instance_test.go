package rocketmq_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunRocketmqInstanceCluster(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_rocketmq_instance." + rnd
	datasourceName := "data.ctyun_rocketmq_instances." + dnd
	resourceFile := "resource_ctyun_rocketmq_instance.tf"
	datasourceFile := "datasource_ctyun_rocketmq_instances.tf"

	cycleResourceFile := "resource_ctyun_rocketmq_instance_on_demand.tf"

	zone := dependence.zone
	nodeNum := 4
	diskSize := 100
	diskType := "SAS"

	initName := "tf-cluster-init-" + utils.GenerateRandomString()

	updatedName := "tf-cluster-updated-" + utils.GenerateRandomString()
	updatedNum := 6
	updatedDiskSize := 1000

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
					dependence.rocketmqClusterSpecName,
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
					resource.TestCheckResourceAttr(resourceName, "spec_name", dependence.rocketmqClusterSpecName),
					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(nodeNum)),
					resource.TestCheckTypeSetElemAttr(resourceName, "zone_list.*", zone),
					resource.TestCheckResourceAttr(resourceName, "disk_type", diskType),
					resource.TestCheckResourceAttr(resourceName, "disk_size", strconv.Itoa(diskSize)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
					func(s *terraform.State) error {
						ds := s.RootModule().Resources[resourceName].Primary
						createTime, expireTime := ds.Attributes["create_time"], ds.Attributes["expire_time"]
						if utils.IsEmptyOrRfc3339(createTime) && utils.IsEmptyOrRfc3339(expireTime) {
							return nil
						}
						return fmt.Errorf("time format doesn't match")
					},
				),
			},
			{
				Config: utils.LoadTestCase(
					cycleResourceFile, rnd,
					initName,
					dependence.rocketmqClusterSpecName,
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
					resource.TestCheckResourceAttrWith(datasourceName, "instances.#", utils.AtLeastOne),
				),
			},
			// 更新属性
			{
				Config: utils.LoadTestCase(
					cycleResourceFile, rnd,
					updatedName,
					dependence.rocketmqClusterSpecName2,
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
					resource.TestCheckResourceAttr(resourceName, "spec_name", dependence.rocketmqClusterSpecName2),
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
					cycleResourceFile, rnd,
					updatedName,
					dependence.rocketmqClusterSpecName2,
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
					resource.TestCheckResourceAttrWith(datasourceName, "instances.#", utils.AtLeastOne),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.Attributes["id"]

					return fmt.Sprintf("%s", id), nil
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
					"expire_time",
					"auto_renew",
				},
			},
			{
				Config: utils.LoadTestCase(
					cycleResourceFile, rnd,
					updatedName,
					dependence.rocketmqClusterSpecName2,
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

//func TestAccCtyunRocketmqInstanceClusterOnDemand(t *testing.T) {
//	t.Parallel()
//	rnd := utils.GenerateRandomString()
//	dnd := utils.GenerateRandomString()
//
//	resourceName := "ctyun_rocketmq_instance." + rnd
//	datasourceName := "data.ctyun_rocketmq_instances." + dnd
//	resourceFile := "resource_ctyun_rocketmq_instance_on_demand.tf"
//	datasourceFile := "datasource_ctyun_rocketmq_instances.tf"
//
//	cycleResourceFile := "resource_ctyun_rocketmq_instance.tf"
//
//	zone := os.Getenv("CTYUN_AZ_NAME")
//	nodeNum := 1
//	diskSize := 100
//	diskType := "SAS"
//
//	initName := "tf-single-init-" + utils.GenerateRandomString()
//
//	updatedName := "tf-single-updated-" + utils.GenerateRandomString()
//
//	resource.Test(t, resource.TestCase{
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			{
//				// 创建
//				Config: utils.LoadTestCase(
//					resourceFile, rnd,
//					initName,
//					dependence.rocketmqSingleSpecName,
//					nodeNum,
//					zone,
//					diskSize,
//					diskType,
//					dependence.vpcID,
//					dependence.subnetID,
//					dependence.securityGroupID,
//				),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "instance_name", initName),
//					resource.TestCheckResourceAttr(resourceName, "spec_name", dependence.rocketmqSingleSpecName),
//					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(nodeNum)),
//					resource.TestCheckTypeSetElemAttr(resourceName, "zone_list.*", zone),
//					resource.TestCheckResourceAttr(resourceName, "disk_type", diskType),
//					resource.TestCheckResourceAttr(resourceName, "disk_size", strconv.Itoa(diskSize)),
//					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
//					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
//					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
//					resource.TestCheckResourceAttrSet(resourceName, "actual_cycle_type"),
//					resource.TestCheckResourceAttrSet(resourceName, "endpoint"),
//					resource.TestCheckResourceAttrSet(resourceName, "ssl_endpoint"),
//				),
//			},
//			// 更新属性
//			{
//				Config: utils.LoadTestCase(
//					cycleResourceFile, rnd,
//					updatedName,
//					dependence.rocketmqSingleSpecName2,
//					nodeNum,
//					zone,
//					diskSize,
//					diskType,
//					dependence.vpcID,
//					dependence.subnetID,
//					dependence.securityGroupID,
//				),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "instance_name", updatedName),
//					resource.TestCheckResourceAttr(resourceName, "spec_name", dependence.rocketmqSingleSpecName2),
//					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(nodeNum)),
//					resource.TestCheckTypeSetElemAttr(resourceName, "zone_list.*", zone),
//					resource.TestCheckResourceAttr(resourceName, "disk_type", diskType),
//					resource.TestCheckResourceAttr(resourceName, "disk_size", strconv.Itoa(diskSize)),
//					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
//					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
//					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
//				),
//			},
//
//			{
//				Config: utils.LoadTestCase(
//					cycleResourceFile, rnd,
//					updatedName,
//					dependence.rocketmqSingleSpecName2,
//					nodeNum,
//					zone,
//					diskSize,
//					diskType,
//					dependence.vpcID,
//					dependence.subnetID,
//					dependence.securityGroupID,
//				) + utils.LoadTestCase(
//					datasourceFile, dnd,
//					resourceName+".id",
//				),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(datasourceName, "instances.#", "1"),
//					resource.TestCheckResourceAttr(datasourceName, "instances.0.instance_name", updatedName),
//				),
//			},
//
//			{
//				ResourceName: resourceName,
//				ImportState:  true,
//				ImportStateIdFunc: func(s *terraform.State) (string, error) {
//					ds := s.RootModule().Resources[resourceName].Primary
//					id := ds.Attributes["id"]
//					regionId := ds.Attributes["region_id"]
//					if id == "" || regionId == "" {
//						return "", fmt.Errorf("id or region_id is required")
//					}
//					return fmt.Sprintf("%s,%s", id, regionId), nil
//				},
//				ImportStateVerify: true,
//				ImportStateVerifyIgnore: []string{
//					"cycle_count",
//					"cycle_type",
//					"disk_size",
//					"disk_type",
//					"master_order_id",
//					"node_num",
//					"project_id",
//					"security_group_id",
//					"subnet_id",
//					"vpc_id",
//					"zone_list",
//					"expire_time",
//				},
//			},
//			{
//				Config: utils.LoadTestCase(
//					cycleResourceFile, rnd,
//					updatedName,
//					dependence.rocketmqSingleSpecName2,
//					nodeNum,
//					zone,
//					diskSize,
//					diskType,
//					dependence.vpcID,
//					dependence.subnetID,
//					dependence.securityGroupID,
//				) + utils.LoadTestCase(
//					datasourceFile, dnd,
//					resourceName+".id",
//				),
//				Destroy: true,
//			},
//		},
//	})
//}
