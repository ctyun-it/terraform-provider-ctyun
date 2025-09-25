package ebm_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunEbm(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_ebm." + rnd
	datasourceName := "data.ctyun_ebms." + dnd
	resourceFile := "resource_ctyun_ebm.tf"
	datasourceFile := "datasource_ctyun_ebms.tf"

	initName := "init"
	initHostname := "init-hostname"
	initPassword := "P@ss-" + utils.GenerateRandomString()
	initStatus := "running"

	updatedName := "updated"
	updatedHostname := "updated-hostname"
	updatedPassword := "P@sstf-" + utils.GenerateRandomString()
	updatedStatus := "stopped"

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
			// 创建
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName, initHostname, initPassword, initStatus,
					dependence.deviceType,
					dependence.imageUUID,
					dependence.securityGroupID,
					dependence.vpcID,
					dependence.systemRaid,
					dependence.dataRaid,
					dependence.subnetID,
					dependence.az2,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", initName),
					resource.TestCheckResourceAttr(resourceName, "hostname", initHostname),
					resource.TestCheckResourceAttr(resourceName, "status", initStatus),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName, updatedHostname, updatedPassword, updatedStatus,
					dependence.deviceType,
					dependence.imageUUID,
					dependence.securityGroupID,
					dependence.vpcID,
					dependence.systemRaid,
					dependence.dataRaid,
					dependence.subnetID,
					dependence.az2,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "hostname", updatedHostname),
					resource.TestCheckResourceAttr(resourceName, "status", updatedStatus),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName, updatedHostname, updatedPassword, updatedStatus,
					dependence.deviceType,
					dependence.imageUUID,
					dependence.securityGroupID,
					dependence.vpcID,
					dependence.systemRaid,
					dependence.dataRaid,
					dependence.subnetID,
					dependence.az2,
				) + utils.LoadTestCase(datasourceFile, dnd, resourceName+".id", dependence.az2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.instance_name", updatedName),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.hostname", updatedHostname),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.status", updatedStatus),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionID := ds.Attributes["region_id"]
					azName := ds.Attributes["az_name"]
					return fmt.Sprintf("%s,%s,%s", id, regionID, azName), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"auto_renew",
					"subnet_id",
					"data_volume_raid_uuid",
					"system_volume_raid_uuid",
					"master_order_id",
					"password",
					"project_id",
					"user_data",
					"cycle_type",
					"image_uuid",
					"cycle_count",
				},
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName, updatedHostname, updatedPassword, initStatus,
					dependence.deviceType,
					dependence.imageUUID,
					dependence.securityGroupID,
					dependence.vpcID,
					dependence.systemRaid,
					dependence.dataRaid,
					dependence.subnetID,
					dependence.az2,
				) + utils.LoadTestCase(datasourceFile, dnd, resourceName+".id", dependence.az2),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName, updatedHostname, updatedPassword, initStatus,
					dependence.deviceType,
					dependence.imageUUID,
					dependence.securityGroupID,
					dependence.vpcID,
					dependence.systemRaid,
					dependence.dataRaid,
					dependence.subnetID,
					dependence.az2,
				) + utils.LoadTestCase(datasourceFile, dnd, resourceName+".id", dependence.az2),
				Destroy: true,
			},
		},
	})
}
