package nat_test

import (
	"fmt"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNewCtyunNatResource(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_nat." + rnd
	datasourceName := "data.ctyun_nats." + dnd
	initDescription := "terraform provider 开发测试"
	resourceFile := "resource_ctyun_nat.tf"
	datasourceFile := "datasource_ctyun_nat.tf"

	vpcId := dependence.vpcID
	spec := "1"
	updatedSpec := "2" // 各类规格已经都试过包括2（中型），3（大型），4（超大型）
	onDemandCycleType := "on_demand"
	monthCycleType := "month"
	cycleCount := fmt.Sprintf(`cycle_count=%d`, 1)

	yearCycleType := "year"

	initName := utils.GenerateRandomString()

	updatedName := utils.GenerateRandomString()
	updatedDescription := utils.GenerateRandomString()

	tcpExpireTime := 10800
	udpExpireTime := 10800
	icmpExpireTime := 800
	tcpDelayCloseTime := 5

	tcpExpireTimeUpdate := 1080
	udpExpireTimeUpdate := 1080
	icmpExpireTimeUpdate := 80
	tcpDelayCloseTimeUpdate := 50

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
			// 1.resource create验证, cycle_type=按需
			// 1.1 Create验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, vpcId, spec, initName, initDescription, onDemandCycleType, "", tcpExpireTime, udpExpireTime, icmpExpireTime, tcpDelayCloseTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", initDescription),
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttrSet(resourceName, "nat_gateway_id"),
					func(s *terraform.State) error {
						ds := s.RootModule().Resources[resourceName].Primary
						createTime, updateTime, expireTime := ds.Attributes["create_time"], ds.Attributes["update_time"], ds.Attributes["expire_time"]
						if utils.IsEmptyOrRfc3339(createTime) && utils.IsEmptyOrRfc3339(updateTime) && utils.IsEmptyOrRfc3339(expireTime) {
							return nil
						}
						return fmt.Errorf("time format doesn't match")
					},
				),
			},
			// 1.2 resource update验证，更新nat name和description
			{
				Config: utils.LoadTestCase(resourceFile, rnd, vpcId, spec, updatedName, updatedDescription, onDemandCycleType, "", tcpExpireTimeUpdate, udpExpireTimeUpdate, icmpExpireTimeUpdate, tcpDelayCloseTimeUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttrSet(resourceName, "nat_gateway_id"),
				),
			},
			// 1.3 resource nat 变配验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, vpcId, updatedSpec, updatedName, updatedDescription, onDemandCycleType, "", tcpExpireTimeUpdate, udpExpireTimeUpdate, icmpExpireTimeUpdate, tcpDelayCloseTimeUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "spec", updatedSpec),
					resource.TestCheckResourceAttrSet(resourceName, "nat_gateway_id"),
				),
			},
			// 1.4 datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, vpcId, updatedSpec, updatedName, updatedDescription, onDemandCycleType, "", tcpExpireTimeUpdate, udpExpireTimeUpdate, icmpExpireTimeUpdate, tcpDelayCloseTimeUpdate) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf(`nat_gateway_id=%s.nat_gateway_id`, resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr(datasourceName, "nats.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "nats.0.name", updatedName),
					resource.TestCheckResourceAttr(datasourceName, "nats.0.description", updatedDescription),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerifyIgnore: []string{
					"master_order_id",
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					return fmt.Sprintf("%s", id), nil
				},
				ImportStateVerifyIgnore: []string{
					"master_order_id",
				},
			},
			// 1.5  销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, vpcId, spec, updatedName, updatedDescription, onDemandCycleType, "", tcpExpireTimeUpdate, udpExpireTimeUpdate, icmpExpireTimeUpdate, tcpDelayCloseTimeUpdate),
				Destroy: true,
			},
			// 2 cycle_type = month类型
			// 2.1 Create验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, vpcId, spec, initName, initDescription, monthCycleType, cycleCount, tcpExpireTimeUpdate, udpExpireTimeUpdate, icmpExpireTimeUpdate, tcpDelayCloseTimeUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", initDescription),
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttrSet(resourceName, "nat_gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", "month"),
					resource.TestCheckResourceAttr(resourceName, "cycle_count", "1"),
					func(s *terraform.State) error {
						ds := s.RootModule().Resources[resourceName].Primary
						createTime, updateTime, expireTime := ds.Attributes["create_time"], ds.Attributes["update_time"], ds.Attributes["expire_time"]
						if utils.IsEmptyOrRfc3339(createTime) && utils.IsEmptyOrRfc3339(updateTime) && utils.IsEmptyOrRfc3339(expireTime) {
							return nil
						}
						return fmt.Errorf("time format doesn't match")
					},
				),
			},
			// 2.2 续费验证
			//{
			//	Config: utils.LoadTestCase(resourceFile, rnd, vpcId, spec, initName, initDescription, monthCycleType, updatedCycleCount),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr(resourceName, "description", initDescription),
			//		resource.TestCheckResourceAttr(resourceName, "name", initName),
			//		resource.TestCheckResourceAttrSet(resourceName, "nat_gateway_id"),
			//		resource.TestCheckResourceAttr(resourceName, "cycle_type", "month"),
			//		resource.TestCheckResourceAttr(resourceName, "cycle_count", "2"),
			//	),
			//},
			// 销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, vpcId, spec, initName, initDescription, monthCycleType, cycleCount, tcpExpireTimeUpdate, udpExpireTimeUpdate, icmpExpireTimeUpdate, tcpDelayCloseTimeUpdate),
				Destroy: true,
			},
			// 3 cycle_type = year类型
			// 3.1 Create验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, vpcId, spec, initName, initDescription, yearCycleType, cycleCount, tcpExpireTimeUpdate, udpExpireTimeUpdate, icmpExpireTimeUpdate, tcpDelayCloseTimeUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", initDescription),
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttrSet(resourceName, "nat_gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", "year"),
					resource.TestCheckResourceAttr(resourceName, "cycle_count", "1"),
				),
			},
			// 销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, vpcId, spec, initName, initDescription, yearCycleType, cycleCount, tcpExpireTime, udpExpireTime, icmpExpireTime, tcpDelayCloseTime),
				Destroy: true,
			},
		},
	})
}
