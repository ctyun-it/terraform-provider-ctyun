package redis_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunRedisInstance(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	and := utils.GenerateRandomString()

	resourceName := "ctyun_redis_instance." + rnd
	datasourceName := "data.ctyun_redis_instances." + dnd
	associationName := "ctyun_redis_association_eip." + and
	resourceFile := "resource_ctyun_redis_instance.tf"
	datasourceFile := "datasource_ctyun_redis_instances.tf"
	associationFile := "resource_ctyun_redis_association_eip.tf"

	initName := "tf-redis-" + utils.GenerateRandomString()
	initPassword := "P@ss" + utils.GenerateRandomString()
	initEngineVersion := "6.0"
	initMaintenanceTime := "00:00-02:00"
	initProtectionStatus := "true"

	updatedPassword := "P@ss" + utils.GenerateRandomString()
	updatedEngineVersion := "7.0"
	updatedMaintenanceTime := "02:00-04:00" // 必须是整点
	updatedProtectionStatus := "false"

	var shardCount, copiesCount string
	if dependence.redisEngineEdition == business.RedisEditionDirectClusterSingle ||
		dependence.redisEngineEdition == business.RedisEditionDirectCluster ||
		dependence.redisEngineEdition == business.RedisEditionClusterOriginalProxy {
		shardCount = "shard_count = 3"
	}
	if dependence.redisEngineEdition == business.RedisEditionOriginalMultipleReadLvs ||
		dependence.redisEngineEdition == business.RedisEditionStandardDual ||
		dependence.redisEngineEdition == business.RedisEditionDirectCluster ||
		dependence.redisEngineEdition == business.RedisEditionClusterOriginalProxy {
		copiesCount = "copies_count = 2"
	}

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
					dependence.redisVersion,
					dependence.redisEngineEdition,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					shardCount, copiesCount,
					initPassword, initEngineVersion, initMaintenanceTime, initProtectionStatus,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", initName),
					resource.TestCheckResourceAttr(resourceName, "version", dependence.redisVersion),
					resource.TestCheckResourceAttr(resourceName, "edition", dependence.redisEngineEdition),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", dependence.securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "password", initPassword),
					resource.TestCheckResourceAttr(resourceName, "engine_version", initEngineVersion),
					resource.TestCheckResourceAttr(resourceName, "maintenance_time", initMaintenanceTime),
					resource.TestCheckResourceAttr(resourceName, "protection_status", initProtectionStatus),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			// 更新属性，同时绑定eip
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					dependence.redisVersion,
					dependence.redisEngineEdition,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					shardCount, copiesCount,
					updatedPassword, updatedEngineVersion, updatedMaintenanceTime, updatedProtectionStatus,
				) + utils.LoadTestCase(
					associationFile, and,
					dependence.eipAddress,
					resourceName+".id",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "password", updatedPassword),
					resource.TestCheckResourceAttr(resourceName, "engine_version", updatedEngineVersion),
					resource.TestCheckResourceAttr(resourceName, "maintenance_time", updatedMaintenanceTime),
					resource.TestCheckResourceAttr(resourceName, "protection_status", updatedProtectionStatus),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			// 验证绑定关系导入
			{
				ResourceName:            associationName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			// 通过查询进行检查
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					dependence.redisVersion,
					dependence.redisEngineEdition,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					shardCount, copiesCount,
					updatedPassword, updatedEngineVersion, updatedMaintenanceTime, updatedProtectionStatus,
				) + utils.LoadTestCase(
					associationFile, and,
					dependence.eipAddress,
					resourceName+".id",
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".instance_name",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.instance_name", initName),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.eip_address", dependence.eipAddress),
				),
			},
			// 解绑
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					dependence.redisVersion,
					dependence.redisEngineEdition,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					shardCount, copiesCount,
					updatedPassword, updatedEngineVersion, updatedMaintenanceTime, updatedProtectionStatus,
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".instance_name",
				),
			},
			// 检查解绑
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					dependence.redisVersion,
					dependence.redisEngineEdition,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					shardCount, copiesCount,
					updatedPassword, updatedEngineVersion, updatedMaintenanceTime, updatedProtectionStatus,
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".instance_name",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.instance_name", initName),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.eip_address", ""),
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
					"cycle_count",
					"cycle_type",
					"edition",
					"host_type",
					"password",
					"project_id",
					"version",
					"master_order_id",
				},
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					dependence.redisVersion,
					dependence.redisEngineEdition,
					dependence.vpcID,
					dependence.subnetID,
					dependence.securityGroupID,
					shardCount, copiesCount,
					updatedPassword, updatedEngineVersion, updatedMaintenanceTime, updatedProtectionStatus,
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					resourceName+".instance_name",
				),
				Destroy: true,
			},
		},
	})
}
