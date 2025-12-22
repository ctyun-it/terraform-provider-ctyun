package pgsql_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCtyunPgsqlWhiteList(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_white_list." + rnd
	datasourceName := "data.ctyun_postgresql_white_lists." + dnd
	datasourceFile := "datasource_ctyun_postgresql_white_lists.tf"
	resourceFile := "resource_ctyun_postgresql_white_list.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	instanceID := dependence.pgsqlID

	// 测试数据
	initialIPs := `["192.168.1.0/24", "10.0.0.1/32"]`
	updatedIPs := `["192.168.1.0/24", "10.0.0.1/32", "172.16.0.0/16"]`
	removedIPs := `["10.0.0.1/32"]`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建白名单测试（覆盖模式）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					instanceID, "cover",
					initialIPs,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceID),
					resource.TestCheckResourceAttr(resourceName, "mode", "cover"),
					resource.TestCheckResourceAttr(resourceName, "ip_list_result.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "ip_list_result.*", "192.168.1.0/24"),
					resource.TestCheckTypeSetElemAttr(resourceName, "ip_list_result.*", "10.0.0.1/32"),
				),
			},
			//  datasource验证
			{
				Config: utils.LoadTestCase(
					datasourceFile, dnd, instanceID, "192.168.1.0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemAttr(datasourceName, "white_list.*", "192.168.1.0/24"),
				),
			},
			// 2. 更新白名单测试（追加模式）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					instanceID, "append",
					updatedIPs,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "append"),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "ip_list_result.*", "192.168.1.0/24"),
					resource.TestCheckTypeSetElemAttr(resourceName, "ip_list_result.*", "10.0.0.1/32"),
					resource.TestCheckTypeSetElemAttr(resourceName, "ip_list_result.*", "172.16.0.0/16"),
				),
			},
			// 3. 更新白名单测试（删除模式）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					instanceID, "delete",
					removedIPs,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "delete"),
					resource.TestCheckResourceAttr(resourceName, "ip_list_result.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "ip_list_result.*", "192.168.1.0/24"),
				),
			},
			// 4. 清理资源（恢复初始状态）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					instanceID, "cover",
					initialIPs,
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例2：仅使用覆盖模式
func TestAccCtyunPgsqlWhiteListCoverOnly(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_white_list." + rnd
	resourceFile := "resource_ctyun_postgresql_white_list.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	instanceID := dependence.pgsqlID

	// 测试数据
	initialIPs := `["192.168.2.0/24", "10.0.0.2/32"]`
	updatedIPs := `["192.168.3.0/24", "10.0.0.3/32"]`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建白名单测试（覆盖模式）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					instanceID, "cover",
					initialIPs,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "cover"),
					resource.TestCheckResourceAttr(resourceName, "ip_list_result.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "ip_list_result.*", "192.168.2.0/24"),
					resource.TestCheckTypeSetElemAttr(resourceName, "ip_list_result.*", "10.0.0.2/32"),
				),
			},
			// 2. 更新白名单测试（覆盖模式）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					instanceID, "cover",
					updatedIPs,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "cover"),
					resource.TestCheckResourceAttr(resourceName, "ip_list_result.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "ip_list_result.*", "192.168.3.0/24"),
					resource.TestCheckTypeSetElemAttr(resourceName, "ip_list_result.*", "10.0.0.3/32"),
				),
			},
			// 3. 清理资源（恢复初始状态）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, projectID,
					instanceID, "cover",
					initialIPs,
				),
				Destroy: true,
			},
		},
	})
}
