package rabbitmq_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunRabbitmqExchange(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_rabbitmq_exchange." + rnd
	datasourceName := "data.ctyun_rabbitmq_exchanges." + dnd
	resourceFile := "resource_ctyun_rabbitmq_exchange.tf"
	datasourceFile := "datasource_ctyun_rabbitmq_exchanges.tf"

	name := utils.GenerateRandomString()
	eType := "topic"

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
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.instanceID, name, eType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", dependence.instanceID),
					resource.TestCheckResourceAttr(resourceName, "auto_delete", "false"),
					resource.TestCheckResourceAttr(resourceName, "durable", "false"),
					resource.TestCheckResourceAttr(resourceName, "internal", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.instanceID, name, eType) +
					utils.LoadTestCase(datasourceFile, dnd, dependence.instanceID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "exchanges.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "exchanges.0.name", name),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					name1 := ds.Attributes["name"]
					vhost := ds.Attributes["vhost"]
					instanceID := ds.Attributes["instance_id"]
					regionId := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s,%s,%s", name1, vhost, instanceID, regionId), nil
				},
				ImportStateVerifyIgnore: []string{
					"internal",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.instanceID, name, eType) +
					utils.LoadTestCase(datasourceFile, dnd, dependence.instanceID, name),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunRabbitmqExchangeAll(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_rabbitmq_exchange." + rnd
	resourceFile := "resource_ctyun_rabbitmq_exchange_all.tf"
	aName := utils.GenerateRandomString()
	aType := "x-delayed-message"
	xDelayedType := "direct"
	alternate := dependence.exchangeName

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
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.instanceID, aName, aType, xDelayedType, alternate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", dependence.instanceID),
					resource.TestCheckResourceAttr(resourceName, "auto_delete", "true"),
					resource.TestCheckResourceAttr(resourceName, "durable", "true"),
					resource.TestCheckResourceAttr(resourceName, "internal", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", aType),
					resource.TestCheckResourceAttr(resourceName, "x_delayed_type", xDelayedType),
					resource.TestCheckResourceAttr(resourceName, "alternate_exchange", alternate),
					resource.TestCheckResourceAttr(resourceName, "name", aName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
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
					"alternate_exchange",
					"internal",
					"x_delayed_type",
				},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, dependence.instanceID, aName, aType, xDelayedType, alternate),
				Destroy: true,
			},
		},
	})
}
