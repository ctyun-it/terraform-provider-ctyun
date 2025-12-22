package kafka_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunKafkaAcl(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_kafka_acl." + rnd
	datasourceName := "data.ctyun_kafka_acls." + dnd
	resourceFile := "resource_ctyun_kafka_acl.tf"
	datasourceFile := "datasource_ctyun_kafka_acls.tf"

	initName := "init-kafka-acl-" + rnd
	instanceId := dependence.instanceID

	initUseNewTopic := false
	updateUseNewTopic := true
	userName := dependence.userName

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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, instanceId, initUseNewTopic, userName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "use_new_topic", fmt.Sprint(initUseNewTopic)),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, instanceId, updateUseNewTopic, userName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "use_new_topic", fmt.Sprint(updateUseNewTopic)),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, instanceId, updateUseNewTopic, userName) +
					utils.LoadTestCase(datasourceFile, dnd, initName, instanceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "acls.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "acls.0.name", initName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					regionId := ds.Attributes["region_id"]
					instanceId := ds.Attributes["instance_id"]
					name := ds.Attributes["name"]
					return fmt.Sprintf("%s,%s,%s", instanceId, name, regionId), nil
				},
				ImportStateVerifyIgnore: []string{"use_new_topic"},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					instanceId := ds.Attributes["instance_id"]
					name := ds.Attributes["name"]
					return fmt.Sprintf("%s,%s", instanceId, name), nil
				},
				ImportStateVerifyIgnore: []string{"use_new_topic"},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, initName, instanceId, updateUseNewTopic, userName) + utils.LoadTestCase(datasourceFile, dnd, initName, instanceId),
				Destroy: true,
			},
		},
	})
}
