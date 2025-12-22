package ec_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunEcSdwanInstance_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_ec_sdwan_instance." + rnd
	resourceFile := "resource_ctyun_ec_sdwan_instance.tf"
	datasourceFile := "datasource_ctyun_ec_sdwan_instances.tf"
	datasourceName := "data.ctyun_ec_sdwan_instances." + dnd

	// 从环境变量获取测试依赖资源
	ecID := dependence.expressConnectID
	cgwID := dependence.cloudGatewayId
	sdwanID := dependence.sdwanId
	rtbID := dependence.rtbID

	weights := 60
	routeLearn := 1
	routeSync := 1

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建SDWAN网络实例测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, sdwanID, rtbID, weights, routeLearn, routeSync,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ec_id", ecID),
					resource.TestCheckResourceAttr(resourceName, "cgw_id", cgwID),
					resource.TestCheckResourceAttr(resourceName, "sdwan_id", sdwanID),
					resource.TestCheckResourceAttr(resourceName, "rtb_id", rtbID),
					resource.TestCheckResourceAttr(resourceName, "weights", "60"),
					resource.TestCheckResourceAttr(resourceName, "route_learn", "1"),
					resource.TestCheckResourceAttr(resourceName, "route_sync", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 更新权重测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, sdwanID, rtbID, 80, routeLearn, routeSync,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "weights", "80"),
				),
			},
			// 3. 数据源测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, sdwanID, rtbID, 80, routeLearn, routeSync,
				) + utils.LoadTestCase(datasourceFile, dnd, ecID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "instances.#")),
			},
			// 4. 导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["id"],
						rs.Primary.Attributes["ec_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cgw_id", "rtb_id", "sdwan_id"}, // 子网列表可能变化

			},
			// 5. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					ecID, cgwID, sdwanID, rtbID, 80, routeLearn, routeSync,
				),
				Destroy: true,
			},
		},
	})
}
