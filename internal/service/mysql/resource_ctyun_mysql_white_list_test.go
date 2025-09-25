package mysql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunMysqlWhiteListTest(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_white_list." + rnd
	datasourceName := "data.ctyun_mysql_white_lists." + dnd

	resourceFile := "resource_ctyun_mysql_white_list.tf"
	datasourceFile := "datasource_ctyun_mysql_white_lists.tf"
	prodInstId := dependence.mysqlID
	groupName := "tf_test"
	groupWhiteList := `"192.168.1.1", "30.8.7.*"`
	updatedGroupWhiteList := `"192.168.1.1", "30.8.7.*", "192.168.1.2", "192.168.1.3"`

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
			// 创建白名单
			{
				Config: utils.LoadTestCase(resourceFile, rnd, prodInstId, groupName, groupWhiteList),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", groupName),
					resource.TestCheckResourceAttr(resourceName, "group_white_list_count", "2"),
				),
			},
			// 更新白名单
			{
				Config: utils.LoadTestCase(resourceFile, rnd, prodInstId, groupName, updatedGroupWhiteList),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", groupName),
					resource.TestCheckResourceAttr(resourceName, "group_white_list_count", "4"),
				),
			},
			// datasource 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, prodInstId, groupName, updatedGroupWhiteList) +
					utils.LoadTestCase(datasourceFile, dnd, prodInstId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "white_lists.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "white_lists.0.group_name", groupName),
					resource.TestCheckResourceAttr(datasourceName, "white_lists.0.group_white_list_count", "4"),
				),
			},
			// 销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, prodInstId, groupName, updatedGroupWhiteList),
				Destroy: true,
			},
		},
	})
}
