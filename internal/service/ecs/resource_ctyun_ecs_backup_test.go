package ecs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunBackup(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_ecs_backup." + rnd
	datasourceName := "data.ctyun_ecs_backups." + dnd
	resourceFile := "resource_ctyun_ecs_backup.tf"
	datasourceFile := "datasource_ctyun_ecs_backups.tf"

	initName := "init-backup"
	updatedName := "updated-backup-" + rnd
	instanceId := dependence.ecsID
	repositoryID := "0cd13a89-5ada-42a7-95e8-60fb9705eecc"

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
				Config: utils.LoadTestCase(resourceFile, rnd, repositoryID, instanceId, initName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, repositoryID, instanceId, updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, repositoryID, instanceId, updatedName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "backups.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "backups.0.name", updatedName),
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
					"full_backup",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, repositoryID, instanceId, updatedName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Destroy: true,
			},
		},
	})
}
