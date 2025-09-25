package ecs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunAffinityGroup(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_ecs_affinity_group." + rnd
	datasourceName := "data.ctyun_ecs_affinity_groups." + dnd
	resourceFile := "resource_ctyun_ecs_affinity_group.tf"
	datasourceFile := "datasource_ctyun_ecs_affinity_groups.tf"

	initName := "init-affinity-group"
	policy := "anti-affinity"
	updatedName := "updated-affinity-group"

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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, policy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "affinity_group_name", initName),
					resource.TestCheckResourceAttr(resourceName, "affinity_group_policy", policy),
					resource.TestCheckResourceAttrSet(resourceName, "affinity_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, policy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "affinity_group_name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "affinity_group_policy", policy),
					resource.TestCheckResourceAttrSet(resourceName, "affinity_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, policy) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "groups.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "groups.0.affinity_group_name", updatedName),
					resource.TestCheckResourceAttr(datasourceName, "groups.0.affinity_group_policy", policy),
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
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, policy) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Destroy: true,
			},
		},
	})
}
