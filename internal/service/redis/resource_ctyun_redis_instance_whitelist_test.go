package redis_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunRedisInstanceWhitelists(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_redis_instance_whitelist." + rnd
	datasourceName := "data.ctyun_redis_instance_whitelists." + dnd
	resourceFile := "resource_ctyun_redis_instance_whitelist.tf"
	datasourceFile := "datasource_ctyun_redis_instance_whitelists.tf"

	initName := "init_redis_instance_whitelist_ip-" + rnd
	instanceId := dependence.instanceId
	ip := "10.0.0.1"
	updateIp := "10.0.0.2"

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
				Config: utils.LoadTestCase(resourceFile, rnd, instanceId, initName, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instanceId, initName, updateIp),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "ip", updateIp),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instanceId, initName, updateIp) +
					utils.LoadTestCase(datasourceFile, dnd, instanceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "rows.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "rows.0.name", initName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					instanceId := ds.Attributes["instance_id"]

					name := ds.Attributes["name"]
					regionId := ds.Attributes["region_id"]

					return fmt.Sprintf("%s,%s,%s", name, instanceId, regionId), nil
				},
				ImportStateVerifyIgnore: []string{},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					instanceId := ds.Attributes["instance_id"]

					name := ds.Attributes["name"]

					return fmt.Sprintf("%s,%s", name, instanceId), nil
				},
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instanceId, initName, updateIp) +
					utils.LoadTestCase(datasourceFile, dnd, instanceId),
				Destroy: true,
			},
		},
	})
}
