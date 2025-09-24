package ebs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunSnapshotPolicy(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_ebs_snapshot_policy." + rnd
	datasourceName := "data.ctyun_ebs_snapshot_policies." + dnd
	resourceFile := "resource_ctyun_ebs_snapshot_policy.tf"
	datasourceFile := "datasource_ctyun_ebs_snapshot_policies.tf"
	bindFile := "resource_ctyun_ebs_snapshot_policy_association.tf"

	initName := "init-snapshot-policy-" + rnd
	updatedName := "updated-snapshot-policy-" + rnd

	diskId := dependence.ebsID

	associationResourceName := "ctyun_ebs_snapshot_policy_association." + dnd

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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2.更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 3.停用
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 4.启用
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 5.查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, true) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "snapshot_policies.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "snapshot_policies.0.name", updatedName),
				),
			},
			// 6.绑定云硬盘
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, true) +
					utils.LoadTestCase(bindFile, dnd, resourceName+".id", diskId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:            associationResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			// 5.解绑云硬盘
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
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
					"project_id", "is_enabled", "bound_disk_num",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, true) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Destroy: true,
			},
		},
	})
}
