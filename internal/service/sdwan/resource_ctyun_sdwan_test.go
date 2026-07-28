package sdwan_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunSdwan_basic(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_sdwan." + rnd
	resourceFile := "resource_ctyun_sdwan.tf"

	name := utils.GenerateRandomString()
	desc := "provider测试创建专用"
	descUpdate := "provider测试更新专用"

	datasourceName := "data.ctyun_sdwans." + dnd
	datasourceFile := "datasource_ctyun_sdwans.tf"
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
			// 创建SD-WAN测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "description", desc),
				),
			},
			// 更新SD-WAN测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, descUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", descUpdate),
				),
			},
			//3. 导入测试
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"}, // 项目ID可能变化

			},
			//datasource 测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, descUpdate) + "\n" + utils.LoadTestCase(datasourceFile, dnd, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttrSet(datasourceName, "sdwans.#"),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, name, descUpdate),
				Destroy: true,
			},
		},
	})
}
