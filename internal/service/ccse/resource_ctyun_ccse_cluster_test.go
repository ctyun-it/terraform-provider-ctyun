package ccse_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunClusterStandard(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_ccse_cluster." + rnd
	datasourceName := "data.ctyun_ccse_clusters." + dnd
	resourceFile := "resource_ctyun_ccse_cluster_standard.tf"
	datasourceFile := "datasource_ctyun_ccse_clusters.tf"

	clusterName := "tf-" + utils.GenerateRandomString()
	clusterSeries := "cce.standard"

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
				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_info.cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "base_info.cluster_series", clusterSeries),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".base_info.cluster_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "records.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "records.0.cluster_name", clusterName),
					resource.TestCheckResourceAttr(datasourceName, "records.0.cluster_series", clusterSeries),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"base_info.nat_gateway_spec",
					"base_info.nginx_ingress_lb_spec",
					"base_info.nginx_ingress_network",
					"base_info.auto_renew",
					"base_info.cluster_domain",
					"base_info.cluster_name",
					"base_info.container_runtime",
					"base_info.cycle_type",
					"base_info.deploy_type",
					"base_info.elb_prod_code",
					"base_info.enable_api_server_eip",
					"base_info.enable_snat",
					"base_info.install_als",
					"base_info.install_als_cube_event",
					"base_info.install_ccse_monitor",
					"base_info.install_nginx_ingress",
					"base_info.ip_vlan",
					"base_info.network_policy",
					"base_info.pod_subnet_id_list.#",
					"base_info.pod_subnet_id_list.0",
					"base_info.project_id",
					"master_host.az_infos",
					"master_host.data_disks",
					"master_host.item_def_name",
					"master_host.sys_disk",
					"master_host.sys_disk.size",
					"master_host.sys_disk.type",
					"master_order_id",
					"master_host",
					"slave_host.az_infos",
					"slave_host.data_disks",
					"slave_host.instance_type",
					"slave_host.item_def_name",
					"slave_host.mirror_id",
					"slave_host.mirror_type",
					"slave_host.sys_disk",
					"slave_host.sys_disk.size",
					"slave_host.sys_disk.type",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".base_info.cluster_name"),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunClusterManaged(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_ccse_cluster." + rnd
	datasourceName := "data.ctyun_ccse_clusters." + dnd
	resourceFile := "resource_ctyun_ccse_cluster_managed.tf"
	datasourceFile := "datasource_ctyun_ccse_clusters.tf"

	clusterName := "tf-" + rnd
	clusterSeries := "cce.managed"

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
				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_info.cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "base_info.cluster_series", clusterSeries),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".base_info.cluster_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "records.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "records.0.cluster_name", clusterName),
					resource.TestCheckResourceAttr(datasourceName, "records.0.cluster_series", clusterSeries),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".base_info.cluster_name"),
				Destroy: true,
			},
		},
	})
}

//func TestAccCtyunClusterStandardEbm(t *testing.T) {
//	t.Parallel()
//	rnd := utils.GenerateRandomString()
//	dnd := utils.GenerateRandomString()
//
//	resourceName := "ctyun_ccse_cluster." + rnd
//	datasourceName := "data.ctyun_ccse_clusters." + dnd
//	resourceFile := "resource_ctyun_ccse_cluster_standard_ebm.tf"
//	datasourceFile := "datasource_ctyun_ccse_clusters.tf"
//
//	clusterName := "tf-ebm-" + utils.GenerateRandomString()
//	clusterSeries := "cce.standard"
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
//				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName, dependence.deviceType, dependence.ebmMirrorName),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "base_info.cluster_name", clusterName),
//					resource.TestCheckResourceAttr(resourceName, "base_info.cluster_series", clusterSeries),
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
//				),
//			},
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName, dependence.deviceType, dependence.ebmMirrorName) +
//					utils.LoadTestCase(datasourceFile, dnd, resourceName+".base_info.cluster_name"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(datasourceName, "records.#", "1"),
//					resource.TestCheckResourceAttr(datasourceName, "records.0.cluster_name", clusterName),
//					resource.TestCheckResourceAttr(datasourceName, "records.0.cluster_series", clusterSeries),
//					resource.ComposeAggregateTestCheckFunc(
//						func(s *terraform.State) error {
//							time.Sleep(30 * time.Second)
//							return nil
//						},
//					),
//				),
//			},
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, clusterName, clusterSeries, dependence.vpcID, dependence.subnetID, dependence.flavorName, dependence.deviceType, dependence.ebmMirrorName) +
//					utils.LoadTestCase(datasourceFile, dnd, resourceName+".base_info.cluster_name"),
//				Destroy: true,
//			},
//		},
//	})
//}
