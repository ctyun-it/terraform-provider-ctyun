package ecs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunSnapshot(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_ecs_snapshot." + rnd
	datasourceName := "data.ctyun_ecs_snapshots." + dnd
	resourceFile := "resource_ctyun_ecs_snapshot.tf"
	datasourceFile := "datasource_ctyun_ecs_snapshots.tf"

	initName := "init-snapshot-" + rnd
	updatedName := "updated-snapshot-" + rnd
	instanceId := dependence.ecsID

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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, instanceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceId),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, instanceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, instanceId) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "snapshots.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "snapshots.0.name", updatedName),
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
					"project_id",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, instanceId) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Destroy: true,
			},
		},
	})
}
