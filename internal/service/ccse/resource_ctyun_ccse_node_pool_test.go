package ccse_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunCcseNodePool(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_ccse_node_pool." + rnd
	datasourceName := "data.ctyun_ccse_node_pools." + dnd
	resourceFile := "resource_ctyun_ccse_node_pool.tf"
	datasourceFile := "datasource_ctyun_ccse_node_pools.tf"

	initName := "init-pool"
	initNodeNum := 0
	initVisibilityPostHostScript := "YWJj"
	initVisibilityHostScript := "MTIz"
	initSysDiskType := "SATA"
	initSysDiskSize := 100
	initDataDiskType := "SATA"
	initDataDiskSize := 200
	initCycleType := "on_demand"
	initExtra := ""
	password := "P@ss" + utils.GenerateRandomString()
	updatedName := "updated-pool"
	updateNodeNum := 2
	updatedVisibilityPostHostScript := "MTIz"
	updatedVisibilityHostScript := "YWJj"
	updatedSysDiskType := "SAS"
	updatedSysDiskSize := 200
	updatedDataDiskType := "SAS"
	updatedDataDiskSize := 400
	updatedCycleType := "month"
	updatedExtra := `cycle_count             = 1
auto_renew = true
`

	var id string
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
				Config: utils.LoadTestCase(resourceFile, rnd,
					initName,
					initVisibilityPostHostScript,
					initVisibilityHostScript,
					initSysDiskType,
					initSysDiskSize,
					initDataDiskType,
					initDataDiskSize,
					initCycleType,
					initExtra,
					dependence.flavorName,
					dependence.clusterID,
					initNodeNum,
					password,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", map[bool]string{true: "true", false: "false"}[false]),
					resource.TestCheckResourceAttr(resourceName, "visibility_post_host_script", initVisibilityPostHostScript),
					resource.TestCheckResourceAttr(resourceName, "visibility_host_script", initVisibilityHostScript),
					resource.TestCheckResourceAttr(resourceName, "sys_disk.type", initSysDiskType),
					resource.TestCheckResourceAttr(resourceName, "sys_disk.size", strconv.Itoa(initSysDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "data_disks.0.type", initDataDiskType),
					resource.TestCheckResourceAttr(resourceName, "data_disks.0.size", strconv.Itoa(initDataDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", initCycleType),
					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(initNodeNum)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					updatedName,
					updatedVisibilityPostHostScript,
					updatedVisibilityHostScript,
					updatedSysDiskType,
					updatedSysDiskSize,
					updatedDataDiskType,
					updatedDataDiskSize,
					updatedCycleType,
					updatedExtra,
					dependence.flavorName,
					dependence.clusterID,
					updateNodeNum,
					password,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", map[bool]string{true: "true", false: "false"}[true]),
					resource.TestCheckResourceAttr(resourceName, "visibility_post_host_script", updatedVisibilityPostHostScript),
					resource.TestCheckResourceAttr(resourceName, "visibility_host_script", updatedVisibilityHostScript),
					resource.TestCheckResourceAttr(resourceName, "sys_disk.type", updatedSysDiskType),
					resource.TestCheckResourceAttr(resourceName, "sys_disk.size", strconv.Itoa(updatedSysDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "data_disks.0.type", updatedDataDiskType),
					resource.TestCheckResourceAttr(resourceName, "data_disks.0.size", strconv.Itoa(updatedDataDiskSize)),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", updatedCycleType),
					resource.TestCheckResourceAttr(resourceName, "cycle_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_num", strconv.Itoa(updateNodeNum)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					updatedName,
					updatedVisibilityPostHostScript,
					updatedVisibilityHostScript,
					updatedSysDiskType,
					updatedSysDiskSize,
					updatedDataDiskType,
					updatedDataDiskSize,
					updatedCycleType,
					updatedExtra,
					dependence.flavorName,
					dependence.clusterID,
					updateNodeNum,
					password,
				) + utils.LoadTestCase(datasourceFile, dnd, dependence.clusterID, resourceName+".name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "records.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "records.0.name", updatedName),
					resource.TestCheckResourceAttr(datasourceName, "records.0.auto_renew", map[bool]string{true: "true", false: "false"}[true]),
					resource.TestCheckResourceAttr(datasourceName, "records.0.visibility_post_host_script", updatedVisibilityPostHostScript),
					resource.TestCheckResourceAttr(datasourceName, "records.0.visibility_host_script", updatedVisibilityHostScript),
					resource.TestCheckResourceAttr(datasourceName, "records.0.sys_disk.type", updatedSysDiskType),
					resource.TestCheckResourceAttr(datasourceName, "records.0.sys_disk.size", strconv.Itoa(updatedSysDiskSize)),
					resource.TestCheckResourceAttr(datasourceName, "records.0.data_disks.0.type", updatedDataDiskType),
					resource.TestCheckResourceAttr(datasourceName, "records.0.data_disks.0.size", strconv.Itoa(updatedDataDiskSize)),
					resource.TestCheckResourceAttr(datasourceName, "records.0.cycle_type", updatedCycleType),
					resource.TestCheckResourceAttr(datasourceName, "records.0.cycle_count", "1"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id = ds.ID
					regionId := ds.Attributes["region_id"]
					clusterId := ds.Attributes["cluster_id"]
					return fmt.Sprintf("%s,%s,%s", id, clusterId, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"az_infos",
					"use_affinity_group",
					"sys_disk",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					updatedName,
					updatedVisibilityPostHostScript,
					updatedVisibilityHostScript,
					updatedSysDiskType,
					updatedSysDiskSize,
					updatedDataDiskType,
					updatedDataDiskSize,
					updatedCycleType,
					updatedExtra,
					dependence.flavorName,
					dependence.clusterID,
					updateNodeNum,
					password,
				) + utils.LoadTestCase(datasourceFile, dnd, dependence.clusterID, resourceName+".name"),
				Destroy: true,
			},
		},
	})
}
