package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMongodbWhiteList_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_mongodb_white_list." + rnd
	resourceFile := "resource_ctyun_mongodb_white_list.tf"

	instance_id := dependence.mongodbID
	groupName := utils.GenerateRandomString()
	ipType := "ipv4"
	whiteListType := "2"
	ipList := "[\"10.138.16.8\",\"10.138.16.10\"]"
	ipListUpdate := "[\"10.138.16.8\",\"10.138.16.119\",\"10.138.16.118\"]"

	datasourceName := "data.ctyun_mongodb_white_lists." + rnd
	datasourceFile := "data_source_ctyun_mongodb_white_lists.tf"

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
			// 基本功能验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id, groupName, ipType, whiteListType, ipList),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "group_name", groupName),
					resource.TestCheckResourceAttr(resourceName, "ip_type", ipType),
					resource.TestCheckResourceAttr(resourceName, "white_list_type", whiteListType),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "2"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id, groupName, ipType, whiteListType, ipListUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "group_name", groupName),
					resource.TestCheckResourceAttr(resourceName, "ip_type", ipType),
					resource.TestCheckResourceAttr(resourceName, "white_list_type", whiteListType),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "3"),
				),
			},
			// 3. 导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s,%s",
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["group_name"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"}, // 项目ID可能变化

			},
			//datasource 测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id, groupName, ipType, whiteListType, ipListUpdate) + "\n" + utils.LoadTestCase(datasourceFile, rnd, instance_id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttrSet(datasourceName, "white_lists.#"),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, instance_id, groupName, ipType, whiteListType, ipListUpdate),
				Destroy: true,
			},
		},
	})
}
