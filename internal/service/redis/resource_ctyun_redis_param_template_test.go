package redis_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunRedisParamTemplate(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_redis_param_template." + rnd
	datasourceName := "data.ctyun_redis_param_templates." + dnd
	resourceFile := "resource_ctyun_redis_param_template.tf"
	datasourceFile := "datasource_ctyun_redis_param_templates.tf"

	initName := "init_redis_param_template-" + rnd
	cacheMode := "ORIGINAL_67"

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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, "Initial Redis template", cacheMode, false, `[
					{
						param_name    = "maxmemory-policy"
						current_value = "allkeys-lru"
					}
				]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "cache_mode", cacheMode),
					resource.TestCheckResourceAttr(resourceName, "sys_template", "false"),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, "Updated Redis template", cacheMode, false, `[
					{
						param_name    = "maxmemory-policy"
						current_value = "allkeys-lru"
					},
					{
						param_name    = "timeout"
						current_value = "300"
					}
				]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated Redis template"),
					resource.TestCheckResourceAttr(resourceName, "params.#", "2"),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, "Updated Redis template", cacheMode, false, `[
					{
						param_name    = "maxmemory-policy"
						current_value = "allkeys-lru"
					},
					{
						param_name    = "timeout"
						current_value = "300"
					}
				]`) +
					utils.LoadTestCase(datasourceFile, dnd, "custom", "1", "10"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(datasourceName, "list.#", utils.AtLeastOne),
					resource.TestCheckResourceAttr(datasourceName, "list.0.name", initName),
				),
			},
			// 导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					regionId := ds.Attributes["region_id"]
					templateId := ds.Attributes["id"]
					return fmt.Sprintf("%s,%s", templateId, regionId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"params", "params_return"},
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary

					templateId := ds.Attributes["id"]
					return fmt.Sprintf("%s", templateId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"params", "params_return"},
			},
			// 清理
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, "Updated Redis template", cacheMode, false, `[
					{
						param_name    = "maxmemory-policy"
						current_value = "allkeys-lru"
					},
					{
						param_name    = "timeout"
						current_value = "300"
					}
				]`) +
					utils.LoadTestCase(datasourceFile, dnd, "custom", "1", "10"),
				Destroy: true,
			},
		},
	})
}
