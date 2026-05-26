package rocketmq_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunRocketmqGroup(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_rocketmq_group." + rnd
	datasourceName := "data.ctyun_rocketmq_groups." + dnd
	resourceFile := "resource_ctyun_rocketmq_group.tf"
	datasourceFile := "datasource_ctyun_rocketmq_groups.tf"
	name := utils.GenerateRandomString()
	consumeEnable := true
	firstConsumeMechanism := int32(1)
	pullMechanism := int32(1)
	retryMaxTimes := int32(16)
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
					consumeEnable,
					firstConsumeMechanism,
					pullMechanism,
					retryMaxTimes,
					remark,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", dependence.instanceID),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "consume_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "first_consume_mechanism", "1"),
					resource.TestCheckResourceAttr(resourceName, "pull_mechanism", "1"),
					resource.TestCheckResourceAttr(resourceName, "retry_max_times", "16"),
					resource.TestCheckResourceAttr(resourceName, "remark", remark),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				// 更新 consume_enable 字段为 false
				Config: utils.LoadTestCase(resourceFile,
					rnd, dependence.instanceID, name,
					false,
					firstConsumeMechanism,
					pullMechanism,
					retryMaxTimes,
					remark,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "consume_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "remark", remark),
				),
			},
			{
				// 恢复 consume_enable 字段为 true
				Config: utils.LoadTestCase(resourceFile,
					rnd, dependence.instanceID, name,
					consumeEnable,
					firstConsumeMechanism,
					pullMechanism,
					retryMaxTimes,
					remarkUpdate,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "consume_enable", "true"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd, dependence.instanceID, name,
					consumeEnable,
					firstConsumeMechanism,
					pullMechanism,
					retryMaxTimes,
					remark,
				) +
					utils.LoadTestCase(datasourceFile, dnd, dependence.instanceID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "groups.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "groups.0.name", name),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					name := ds.Attributes["name"]
					instanceId := ds.Attributes["instance_id"]
					regionId := ds.Attributes["region_id"]
					if name == "" || instanceId == "" || regionId == "" {
						return "", fmt.Errorf("name, instance_id and region_id are required")
					}
					return fmt.Sprintf("%s,%s,%s", name, instanceId, regionId), nil
				},
				ImportStateVerify: true,
			},
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd, dependence.instanceID, name,
					consumeEnable,
					firstConsumeMechanism,
					pullMechanism,
					retryMaxTimes,
					remark,
				) +
					utils.LoadTestCase(datasourceFile, dnd, dependence.instanceID, name),
				Destroy: true,
			},
		},
	})
}
