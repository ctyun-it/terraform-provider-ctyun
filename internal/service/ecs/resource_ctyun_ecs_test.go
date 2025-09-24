package ecs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccCtyunEcs(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	and := utils.GenerateRandomString()

	resourceName := "ctyun_ecs." + rnd
	datasourceName := "data.ctyun_ecs_instances." + dnd
	resourceFile := "resource_ctyun_ecs.tf"
	datasourceFile := "datasource_ctyun_ecs_instances.tf"
	associationFile := "resource_ctyun_ecs_affinity_group_association.tf"

	instanceName := "tf-test-ecs"

	initDisplayName := "tf-test-init-ecs"
	initSysDiskSize := 60
	initExtra := `cycle_type         = "on_demand"`
	nextShelveExtra := `cycle_type         = "on_demand" 
   status = "shelve"`
	nextStartExtra := `cycle_type         = "on_demand" 
   status = "running"`

	updatedDisplayName := "tf-test-updated-ecs"
	updatedSysDiskSize := 100
	updatedExtra := fmt.Sprintf(
		`status = "running"
  cycle_type         = "month"
  cycle_count         = 1
  security_group_ids     = ["%s"]`, dependence.securityGroupID)

	nextExtra := fmt.Sprintf(
		`status = "stopped"
   cycle_type         = "month"
   cycle_count         = 1
  security_group_ids     = ["%s"]
  is_destroy_instance  = true`, dependence.securityGroupID)

	affinityGroupAssociationResourceName := "ctyun_ecs_affinity_group_association." + and

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
			// 1.创建
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceName,
					initDisplayName,
					dependence.flavorID,
					dependence.imageID,
					initSysDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.keyPairName,
					initExtra,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", initDisplayName),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", dependence.flavorID),
					resource.TestCheckResourceAttr(resourceName, "image_id", dependence.imageID),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", strconv.Itoa(initSysDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "key_pair_name", dependence.keyPairName),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			// 2.节省关机
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceName,
					initDisplayName,
					dependence.flavorID,
					dependence.imageID,
					initSysDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.keyPairName,
					nextShelveExtra,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "shelve"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			//3.开机
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceName,
					initDisplayName,
					dependence.flavorID,
					dependence.imageID,
					initSysDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.keyPairName,
					nextStartExtra,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 4.更新名称、密钥、安全组属性，并进行按需转包周期
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceName,
					updatedDisplayName,
					dependence.flavorID,
					dependence.imageID,
					initSysDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.keyPairName2,
					updatedExtra,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", updatedDisplayName),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", dependence.flavorID),
					resource.TestCheckResourceAttr(resourceName, "image_id", dependence.imageID),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", strconv.Itoa(initSysDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "key_pair_name", dependence.keyPairName2),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckTypeSetElemAttr(resourceName, "security_group_ids.*", dependence.securityGroupID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},

			// 5.关机，然后更新规格、系统盘大小，并进行包周期转按需(不生效 去掉)，同时关联主机组
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceName,
					updatedDisplayName,
					dependence.flavorID2,
					dependence.imageID,
					updatedSysDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.keyPairName2,
					nextExtra,
				) + utils.LoadTestCase(associationFile, and, resourceName+".id", dependence.affinityGroupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", updatedDisplayName),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", dependence.flavorID2),
					resource.TestCheckResourceAttr(resourceName, "image_id", dependence.imageID),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", strconv.Itoa(updatedSysDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", dependence.vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", dependence.subnetID),
					resource.TestCheckResourceAttr(resourceName, "key_pair_name", dependence.keyPairName2),
					resource.TestCheckResourceAttr(resourceName, "status", "stopped"),
					resource.TestCheckTypeSetElemAttr(resourceName, "security_group_ids.*", dependence.securityGroupID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			// 6.检查是否关联成功
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceName,
					updatedDisplayName,
					dependence.flavorID2,
					dependence.imageID,
					updatedSysDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.keyPairName2,
					nextExtra,
				) + utils.LoadTestCase(associationFile, and, resourceName+".id", dependence.affinityGroupID) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "instances.0.affinity_group.affinity_group_id", dependence.affinityGroupID),
				),
			},
			{
				ResourceName:            affinityGroupAssociationResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			// 7.解绑主机组
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceName,
					updatedDisplayName,
					dependence.flavorID2,
					dependence.imageID,
					updatedSysDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.keyPairName2,
					nextExtra,
				),
			},
			// 8.检查是否解绑成功
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceName,
					updatedDisplayName,
					dependence.flavorID2,
					dependence.imageID,
					updatedSysDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.keyPairName2,
					nextExtra,
				) + utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(datasourceName, "instances.0.affinity_group.affinity_group_id"),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceName,
					updatedDisplayName,
					dependence.flavorID2,
					dependence.imageID,
					updatedSysDiskSize,
					dependence.vpcID,
					dependence.subnetID,
					dependence.keyPairName2,
					nextExtra,
				) + utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Destroy: true,
			},
		},
	},
	)
}
