package rabbitmq_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunRabbitmqVhost(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_rabbitmq_vhost." + rnd
	datasourceName := "data.ctyun_rabbitmq_vhosts." + dnd
	resourceFile := "resource_ctyun_rabbitmq_vhost.tf"
	datasourceFile := "datasource_ctyun_rabbitmq_vhosts.tf"
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
					resource.TestCheckResourceAttr(datasourceName, "vhosts.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "vhosts.0.name", name),
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
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					name := ds.Attributes["name"]
					instanceId := ds.Attributes["instance_id"]
					return fmt.Sprintf("%s,%s", name, instanceId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependence.instanceID, name) +
					utils.LoadTestCase(datasourceFile, dnd, dependence.instanceID, name),
				Destroy: true,
			},
		},
	})
}
