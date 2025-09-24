package ebs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunBackupPolicy(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_ebs_backup_policy." + rnd
	datasourceName := "data.ctyun_ebs_backup_policies." + dnd
	resourceFile := "resource_ctyun_ebs_backup_policy.tf"
	datasourceFile := "datasource_ctyun_ebs_backup_policies.tf"
	bindDisksFile := "resource_ctyun_ebs_backup_policy_bind_disks.tf"
	bindRepoFile := "resource_ctyun_ebs_backup_policy_bind_repo.tf"

	initName := "init-backup-policy-" + rnd
	updatedName := "updated-backup-policy-" + rnd

	diskId := dependence.ebsID
	//TODO 获取存储库ID替换
	repositoryID := "671f67c4-6131-4154-8c1d-7c5b82edd1eb"

	bindDisksResourceName := "ctyun_ebs_backup_policy_bind_disks." + dnd
	bindRepositoryResourceName := "ctyun_ebs_backup_policy_bind_repo." + dnd
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
				Config: utils.LoadTestCase(resourceFile, rnd, initName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2.更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 3.查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "backup_policies.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "backup_policies.0.name", updatedName),
				),
			},
			// 4.绑定云硬盘
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName) +
					utils.LoadTestCase(bindDisksFile, dnd, resourceName+".id", diskId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 5.检查是否关联成功
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName) +
					utils.LoadTestCase(bindDisksFile, dnd, resourceName+".id", diskId) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "backup_policies.0.resource_ids", diskId),
				),
			},
			{
				ResourceName:            bindDisksResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			// 6.解绑云硬盘
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 7.检查是否解绑成功
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "backup_policies.0.resource_ids", ""),
				),
			},
			// 8.云硬盘备份策略绑定存储库
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName) +
					utils.LoadTestCase(bindRepoFile, dnd, resourceName+".id", repositoryID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 9.检查是否关联成功
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName) +
					utils.LoadTestCase(bindRepoFile, dnd, resourceName+".id", repositoryID) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// 先检查列表不为空
					resource.TestCheckResourceAttr(datasourceName, "backup_policies.0.repository_list.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "backup_policies.0.repository_list.0.repository_id", repositoryID),
				),
			},
			{
				ResourceName:            bindRepositoryResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			// 10.云硬盘备份策略解绑存储库
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 11.检查是否解绑成功
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "backup_policies.0.repository_list.#", "0"),
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
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Destroy: true,
			},
		},
	})
}
