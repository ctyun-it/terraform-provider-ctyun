package rabbitmq_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunRabbitmqQueue(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_rabbitmq_queue." + rnd
	datasourceName := "data.ctyun_rabbitmq_queues." + dnd
	resourceFile := "resource_ctyun_rabbitmq_queue.tf"
	datasourceFile := "datasource_ctyun_rabbitmq_queues.tf"
	name := utils.GenerateRandomString()

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
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.instanceID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", dependence.instanceID),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.instanceID, name) +
					utils.LoadTestCase(datasourceFile, dnd, dependence.instanceID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "queues.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "queues.0.name", name),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"x_max_priority",
					"x_overflow",
					"x_queue_mode",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.instanceID, name) +
					utils.LoadTestCase(datasourceFile, dnd, dependence.instanceID, name),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunRabbitmqQueueAll(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_rabbitmq_queue." + rnd
	datasourceName := "data.ctyun_rabbitmq_queues." + dnd
	resourceFile := "resource_ctyun_rabbitmq_queue_all.tf"
	datasourceFile := "datasource_ctyun_rabbitmq_queues.tf"
	name := utils.GenerateRandomString()
	x_queue_mode := "lazy"
	x_overflow := "reject-publish"
	x_dead_letter_exchange := dependence.exchangeName
	x_dead_letter_routing_key := "amqkey"
	x_message_ttl := 3600000
	x_max_length := 1000
	x_expires := 1000000
	x_max_priority := 101

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
					x_queue_mode,
					x_overflow,
					x_dead_letter_exchange,
					x_dead_letter_routing_key,
					x_message_ttl,
					x_max_length,
					x_expires,
					x_max_priority,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", dependence.instanceID),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "x_queue_mode", x_queue_mode),
					resource.TestCheckResourceAttr(resourceName, "x_overflow", x_overflow),
					resource.TestCheckResourceAttr(resourceName, "x_dead_letter_exchange", x_dead_letter_exchange),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd, dependence.instanceID, name,
					x_queue_mode,
					x_overflow,
					x_dead_letter_exchange,
					x_dead_letter_routing_key,
					x_message_ttl,
					x_max_length,
					x_expires,
					x_max_priority,
				) +
					utils.LoadTestCase(datasourceFile, dnd, dependence.instanceID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "queues.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "queues.0.name", name),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					name := ds.Attributes["name"]
					vhost := ds.Attributes["vhost"]
					instanceId := ds.Attributes["instance_id"]
					regionId := ds.Attributes["region_id"]
					if name == "" || vhost == "" || instanceId == "" || regionId == "" {
						return "", fmt.Errorf("name, vhost, instance_id and region_id are required")
					}
					return fmt.Sprintf("%s,%s,%s,%s", name, vhost, instanceId, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"x_dead_letter_routing_key",
					"x_dead_letter_exchange",
					"x_expires",
					"x_max_length",
					"x_message_ttl",
					"x_max_priority",
					"x_overflow",
					"x_queue_mode",
				},
			}, {
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					name := ds.Attributes["name"]
					vhost := ds.Attributes["vhost"]
					instanceId := ds.Attributes["instance_id"]
					return fmt.Sprintf("%s,%s,%s", name, vhost, instanceId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"x_dead_letter_routing_key",
					"x_dead_letter_exchange",
					"x_expires",
					"x_max_length",
					"x_message_ttl",
					"x_max_priority",
					"x_overflow",
					"x_queue_mode",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd, dependence.instanceID, name,
					x_queue_mode,
					x_overflow,
					x_dead_letter_exchange,
					x_dead_letter_routing_key,
					x_message_ttl,
					x_max_length,
					x_expires,
					x_max_priority,
				) +
					utils.LoadTestCase(datasourceFile, dnd, dependence.instanceID, name),
				Destroy: true,
			},
		},
	})
}
