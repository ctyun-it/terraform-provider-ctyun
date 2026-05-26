package rocketmq_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunRocketmqTopic(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_rocketmq_topic." + rnd
	datasourceName := "data.ctyun_rocketmq_topics." + dnd
	resourceFile := "resource_ctyun_rocketmq_topic.tf"
	datasourceFile := "datasource_ctyun_rocketmq_topics.tf"
	name := utils.GenerateRandomString()
	order := false
	perm := 6
	remark := "test remark"
	remarkUpdate := "test remark update"

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
				// 创建
				Config: utils.LoadTestCase(resourceFile,
					rnd, dependence.instanceID, name,
					order,
					perm,
					remark,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", dependence.instanceID),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "order", "false"),
					resource.TestCheckResourceAttr(resourceName, "perm", "6"),
					resource.TestCheckResourceAttr(resourceName, "remark", remark),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				// 更新 perm 字段为 2（禁读/只写）
				Config: utils.LoadTestCase(resourceFile,
					rnd, dependence.instanceID, name,
					order,
					2,
					remarkUpdate,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "perm", "2"),
					resource.TestCheckResourceAttr(resourceName, "remark", remarkUpdate),
				),
			},
			{
				// 更新 perm 字段为 4（禁写/只读）
				Config: utils.LoadTestCase(resourceFile,
					rnd, dependence.instanceID, name,
					order,
					4,
					remark,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "perm", "4"),
				),
			},
			{
				// 恢复 perm 字段为 6（读写）
				Config: utils.LoadTestCase(resourceFile,
					rnd, dependence.instanceID, name,
					order,
					6,
					remark,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "perm", "6"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd, dependence.instanceID, name,
					order,
					perm,
					remark,
				) +
					utils.LoadTestCase(datasourceFile, dnd, dependence.instanceID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "topics.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "topics.0.name", name),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					topicname := ds.Attributes["name"]
					instanceId := ds.Attributes["instance_id"]
					regionId := ds.Attributes["region_id"]
					if topicname == "" || instanceId == "" || regionId == "" {
						return "", fmt.Errorf("name, instance_id and region_id are required")
					}
					return fmt.Sprintf("%s,%s,%s", instanceId, topicname, regionId), nil
				},
				ImportStateVerify: true,
			},
			//{
			//	ResourceName: resourceName,
			//	ImportState:  true,
			//	ImportStateIdFunc: func(s *terraform.State) (string, error) {
			//		ds := s.RootModule().Resources[resourceName].Primary
			//		name := ds.Attributes["name"]
			//		instanceId := ds.Attributes["instance_id"]
			//		return fmt.Sprintf("%s,%s", name, instanceId), nil
			//	},
			//	ImportStateVerify: true,
			//},
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd, dependence.instanceID, name,
					order,
					perm,
					remark,
				) +
					utils.LoadTestCase(datasourceFile, dnd, dependence.instanceID, name),
				Destroy: true,
			},
		},
	})
}
