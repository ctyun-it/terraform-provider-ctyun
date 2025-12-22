package kafka_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunKafkaConsumerGroups(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_kafka_consumer_group." + rnd
	datasourceName := "data.ctyun_kafka_consumer_groups." + dnd
	resourceFile := "resource_ctyun_kafka_consumer_group.tf"
	datasourceFile := "datasource_ctyun_kafka_consumer_groups.tf"

	initName := "init-kafka_consumer_group-" + rnd
	instanceID := dependence.instanceID
	topicName := dependence.topicName

	resetConfig := fmt.Sprintf(`reset_config = {
      topic_name = "%s"
      type       = 1
      timestamp       = 1571299747516
    }`, topicName)
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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, instanceID, "desc", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, instanceID, "desc-update", resetConfig),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "description", "desc-update"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, instanceID, "desc-update", "") +
					utils.LoadTestCase(datasourceFile, dnd, initName, instanceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "consumer_groups.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "consumer_groups.0.name", initName),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					regionId := ds.Attributes["region_id"]
					instanceId := ds.Attributes["instance_id"]
					groupName := ds.Attributes["name"]
					return fmt.Sprintf("%s,%s,%s", instanceId, groupName, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project_id",
				},
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					instanceId := ds.Attributes["instance_id"]
					groupName := ds.Attributes["name"]
					return fmt.Sprintf("%s,%s", instanceId, groupName), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project_id",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, instanceID, "desc-update", "") +
					utils.LoadTestCase(datasourceFile, dnd, initName, instanceID),
				Destroy: true,
			},
		},
	})
}
