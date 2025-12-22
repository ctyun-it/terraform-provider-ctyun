package hpfs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccCtyunHpfs(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_hpfs." + rnd
	resourceFile := "resource_ctyun_hpfs.tf"

	datasourceName := "data.ctyun_hpfs_instances." + dnd
	datasourceFile := "datasource_ctyun_hfps_instances.tf"
	sfsProtocol := "hpfs"
	cycleType := "on_demand"
	sfsName := "hpfs-" + utils.GenerateRandomString()
	updatedSfsName := "hpfs-" + utils.GenerateRandomString() + "-new"
	sfsSize := 512
	updatedSfsSize := 1024
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
			// 开通hpfs，
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, cycleType, sfsName, sfsSize),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "protocol", sfsProtocol),
					resource.TestCheckResourceAttr(resourceName, "name", sfsName),
					resource.TestCheckResourceAttr(resourceName, "size", strconv.Itoa(sfsSize)),
				),
			},
			// 变配sfs name 和 SIZE规格 512->1024
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, cycleType, updatedSfsName, updatedSfsSize),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "protocol", sfsProtocol),
					resource.TestCheckResourceAttr(resourceName, "name", updatedSfsName),
					resource.TestCheckResourceAttr(resourceName, "size", strconv.Itoa(updatedSfsSize)),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, cycleType, updatedSfsName, updatedSfsSize) +
					utils.LoadTestCase(datasourceFile, dnd, "available", sfsProtocol),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "instances.0.protocol", sfsProtocol),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.status", "available"),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, sfsProtocol, cycleType, updatedSfsName, updatedSfsSize),
				Destroy: true,
			},
		},
	})
}

// 指定集群和baseline
func TestAccCtyunHpfs1(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_hpfs." + rnd
	resourceFile := "resource_ctyun_hpfs1.tf"
	sfsProtocol := "hpfs"
	cluster := dependence.clusterName
	baseline := "200"
	sfsName := "hpfs-" + utils.GenerateRandomString()
	updatedSfsName := "hpfs-" + utils.GenerateRandomString() + "-new"
	sfsSize := 512
	updatedSfsSize := 512
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
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, sfsName, sfsSize, cluster, baseline),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "protocol", sfsProtocol),
					resource.TestCheckResourceAttr(resourceName, "name", sfsName),
					resource.TestCheckResourceAttr(resourceName, "size", strconv.Itoa(sfsSize)),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", cluster),
					resource.TestCheckResourceAttr(resourceName, "baseline", baseline),
				),
			},
			// 变配sfs name 和 SIZE规格 512->1024
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, updatedSfsName, updatedSfsSize, cluster, baseline),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "protocol", sfsProtocol),
					resource.TestCheckResourceAttr(resourceName, "name", updatedSfsName),
					resource.TestCheckResourceAttr(resourceName, "size", strconv.Itoa(updatedSfsSize)),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", cluster),
					resource.TestCheckResourceAttr(resourceName, "baseline", baseline),
				),
			},
			{
				Config:       utils.LoadTestCase(resourceFile, rnd, sfsProtocol, updatedSfsName, updatedSfsSize, cluster, baseline),
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					return fmt.Sprintf("%s", id), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"az_name", "cycle_type", "master_order_id", "update_time"},
			},
			// importState 2
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					projectID := ds.Attributes["project_id"]
					regionID := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s,%s", id, projectID, regionID), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"az_name", "cycle_type", "master_order_id", "update_time"},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, sfsProtocol, updatedSfsName, updatedSfsSize, cluster, baseline),
				Destroy: true,
			},
		},
	})
}
