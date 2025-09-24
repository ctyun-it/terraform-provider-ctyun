package scaling_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

// 多az，单安全组->多安全组
func TestAccCtyunScaling(t *testing.T) {

	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()

	resourceName := "ctyun_scaling_group." + rnd
	resourceFile := "resource_ctyun_scaling.tf"

	//datasourceName := "data.ctyun_scalings." + dnd
	//datasourceFile := "datasource_ctyun_scalings.tf"

	securityGroupIDList := fmt.Sprintf(`["%s"]`, dependence.securityGroupID)
	name := "scaling-group-" + utils.GenerateRandomString()
	healthMode := "server"
	subnetIDList := fmt.Sprintf(`["%s"]`, dependence.subnetID)
	moveOutStrategy := "earlier_config"
	vpcId := dependence.vpcID
	minCount := 1
	maxCount := 1

	expectedCount := 1
	healthPeriod := 300
	useLb := 1
	lbList := fmt.Sprintf(`[{"port": 12306, "lb_id": "%s", "weight": 1, "host_group_id": "%s"}]`, dependence.loadbalancerID, dependence.targetGroupID)
	scalingConfigID, err := strconv.ParseInt(dependence.scalingConfigID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	scalingConfigID1, err := strconv.ParseInt(dependence.scalingConfigID1, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	configList := fmt.Sprintf(`[%d]`, scalingConfigID)
	azStrategy := "uniform_distribution"

	updateSecurityGroupIDList := fmt.Sprintf(`["%s", "%s"]`, dependence.securityGroupID, dependence.securityGroupID1)
	updateName := "scaling-group-" + utils.GenerateRandomString()
	healthModeLb := "lb"
	//updatedSubnetIDList := fmt.Sprintf(`["%s", "%s"]`, dependence.subnetID, dependence.subnetID1)
	updatemoveOutStrategy := "earlier_vm"
	updateMinCount := 2
	updatedMaxCount := 2
	updatedExpectedCount := 2
	updateHealthPeriod := 10080
	updatedConfigList := fmt.Sprintf(`[%d, %d]`, scalingConfigID, scalingConfigID1)
	azStrategyPriorityDistribution := "priority_distribution"
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
			// 创建弹性组
			{
				Config: utils.LoadTestCase(resourceFile, rnd, securityGroupIDList, name, healthMode, subnetIDList, moveOutStrategy, vpcId, minCount, maxCount, expectedCount,
					healthPeriod, useLb, lbList, configList, azStrategy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "health_mode", healthMode),
					resource.TestCheckResourceAttr(resourceName, "move_out_strategy", moveOutStrategy),
					resource.TestCheckResourceAttr(resourceName, "min_count", fmt.Sprintf("%d", minCount)),
					resource.TestCheckResourceAttr(resourceName, "max_count", fmt.Sprintf("%d", maxCount)),
					resource.TestCheckResourceAttr(resourceName, "expected_count", fmt.Sprintf("%d", expectedCount)),
					resource.TestCheckResourceAttr(resourceName, "health_period", fmt.Sprintf("%d", healthPeriod)),
					resource.TestCheckResourceAttr(resourceName, "use_lb", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "az_strategy", azStrategy),
				),
			},
			// 更新name-> new; health_mode:云主机健康检查->云主机健康检查, subnet不修改，移除策略 earlier_config->earlier_vm，min_count: 1->2, max_count: 1->2， expected_count: 2->2, health_period:300->10080,
			// configList 增加一个config配置，az_strategy : uniform_distribution-> priority_distribution
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updateSecurityGroupIDList, updateName, healthModeLb, subnetIDList, updatemoveOutStrategy, vpcId, updateMinCount, updatedMaxCount, updatedExpectedCount,
					updateHealthPeriod, useLb, lbList, updatedConfigList, azStrategyPriorityDistribution),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "health_mode", healthModeLb),
					resource.TestCheckResourceAttr(resourceName, "move_out_strategy", updatemoveOutStrategy),
					resource.TestCheckResourceAttr(resourceName, "min_count", fmt.Sprintf("%d", updateMinCount)),
					resource.TestCheckResourceAttr(resourceName, "max_count", fmt.Sprintf("%d", updatedMaxCount)),
					resource.TestCheckResourceAttr(resourceName, "expected_count", fmt.Sprintf("%d", updatedExpectedCount)),
					resource.TestCheckResourceAttr(resourceName, "health_period", fmt.Sprintf("%d", updateHealthPeriod)),
					resource.TestCheckResourceAttr(resourceName, "use_lb", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "az_strategy", azStrategyPriorityDistribution),
				),
			},

			// 销毁
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updateSecurityGroupIDList, updateName, healthModeLb, subnetIDList, updatemoveOutStrategy, vpcId, updateMinCount, updatedMaxCount, updatedExpectedCount,
					updateHealthPeriod, useLb, lbList, updatedConfigList, azStrategyPriorityDistribution),
				Destroy: true,
			},
		}},
	)
}

// 多az，使用LB，停用弹性伸缩组，更新子网和lb。 校验datasource。创建时不绑定config，更新时添加
func TestAccCtyunScaling1(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_scaling_group." + rnd
	resourceName1 := "ctyun_scaling_group." + rnd
	resourceFile := "resource_ctyun_scaling.tf"
	resourceFile1 := "resource_ctyun_scaling1.tf"

	datasourceName := "data.ctyun_scalings." + dnd
	datasourceFile := "datasource_ctyun_scalings.tf"

	securityGroupIDList := fmt.Sprintf(`["%s"]`, dependence.securityGroupID)
	name := "scaling-group-" + utils.GenerateRandomString()
	healthMode := "lb"
	subnetIDList := fmt.Sprintf(`["%s"]`, dependence.subnetID)
	moveOutStrategy := "later_config"
	vpcId := dependence.vpcID
	minCount := 1
	maxCount := 1

	expectedCount := 1
	healthPeriod := 300
	useLb := 1
	lbList := fmt.Sprintf(`[{"port": 12306, "lb_id": "%s", "weight": 1, "host_group_id": "%s"}]`, dependence.loadbalancerID, dependence.targetGroupID)
	scalingConfigID, err := strconv.ParseInt(dependence.scalingConfigID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	scalingConfigID1, err := strconv.ParseInt(dependence.scalingConfigID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	configList := fmt.Sprintf(`[%d, %d]`, scalingConfigID, scalingConfigID1)
	azStrategy := "uniform_distribution"
	status := "disable"
	deleteProtection := "disable"

	updatedSubnetIDList := fmt.Sprintf(`["%s", "%s","%s", "%s", "%s"]`, dependence.subnetID, dependence.subnetID1, dependence.subnetID1, dependence.subnetID1, dependence.subnetID1)
	updatedLbList := fmt.Sprintf(`[{"port": 12306, "lb_id": "%s", "weight": 1, "host_group_id": "%s"},{"port": 12305, "lb_id": "%s", "weight": 1, "host_group_id": "%s"} ]`,
		dependence.loadbalancerID, dependence.targetGroupID, dependence.loadbalancerID1, dependence.targetGroupID1)

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
			// 创建弹性组
			{
				Config: utils.LoadTestCase(resourceFile, rnd, securityGroupIDList, name, healthMode, subnetIDList, moveOutStrategy, vpcId, minCount, maxCount, expectedCount,
					healthPeriod, useLb, lbList, configList, azStrategy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "health_mode", healthMode),
					resource.TestCheckResourceAttr(resourceName, "move_out_strategy", moveOutStrategy),
					resource.TestCheckResourceAttr(resourceName, "min_count", fmt.Sprintf("%d", minCount)),
					resource.TestCheckResourceAttr(resourceName, "max_count", fmt.Sprintf("%d", maxCount)),
					resource.TestCheckResourceAttr(resourceName, "expected_count", fmt.Sprintf("%d", expectedCount)),
					resource.TestCheckResourceAttr(resourceName, "health_period", fmt.Sprintf("%d", healthPeriod)),
					resource.TestCheckResourceAttr(resourceName, "use_lb", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "az_strategy", azStrategy),
				),
			},
			// 验证datasource
			{
				Config: utils.LoadTestCase(resourceFile, rnd, securityGroupIDList, name, healthMode, subnetIDList, moveOutStrategy, vpcId, minCount, maxCount, expectedCount,
					healthPeriod, useLb, lbList, configList, azStrategy) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("%s.id", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "scaling_list.0.name", name),
					resource.TestCheckResourceAttr(datasourceName, "scaling_list.0.health_period", fmt.Sprintf("%d", healthPeriod)),
					resource.TestCheckResourceAttr(datasourceName, "scaling_list.0.health_mode", healthMode),
					resource.TestCheckResourceAttr(datasourceName, "scaling_list.0.max_count", fmt.Sprintf("%d", maxCount)),
					resource.TestCheckResourceAttr(datasourceName, "scaling_list.0.min_count", fmt.Sprintf("%d", minCount)),
					resource.TestCheckResourceAttr(datasourceName, "scaling_list.0.expected_count", fmt.Sprintf("%d", expectedCount)),
					resource.TestCheckResourceAttr(datasourceName, "scaling_list.0.move_out_strategy", moveOutStrategy),
					resource.TestCheckResourceAttr(datasourceName, "scaling_list.0.use_lb", "1"),
				),
			},
			// 关机并更新subnet和lb
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, securityGroupIDList, name, healthMode, updatedSubnetIDList, moveOutStrategy, vpcId, minCount, maxCount,
					healthPeriod, useLb, updatedLbList, configList, azStrategy, status, deleteProtection),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName1, "id"),
					resource.TestCheckResourceAttr(resourceName1, "name", name),
					resource.TestCheckResourceAttr(resourceName1, "health_mode", healthMode),
					resource.TestCheckResourceAttr(resourceName1, "move_out_strategy", moveOutStrategy),
					resource.TestCheckResourceAttr(resourceName1, "min_count", fmt.Sprintf("%d", minCount)),
					resource.TestCheckResourceAttr(resourceName1, "max_count", fmt.Sprintf("%d", maxCount)),
					resource.TestCheckResourceAttr(resourceName1, "expected_count", fmt.Sprintf("%d", expectedCount)),
					resource.TestCheckResourceAttr(resourceName1, "health_period", fmt.Sprintf("%d", healthPeriod)),
					resource.TestCheckResourceAttr(resourceName1, "use_lb", "1"),
					resource.TestCheckResourceAttr(resourceName1, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName1, "az_strategy", azStrategy),
					resource.TestCheckResourceAttr(resourceName1, "subnet_id_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName1, "lb_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName1, "status", status),
					resource.TestCheckResourceAttr(resourceName1, "delete_protection", deleteProtection),
				),
			},
			// 开启弹性伸缩服务
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, securityGroupIDList, name, healthMode, updatedSubnetIDList, moveOutStrategy, vpcId, minCount, maxCount,
					healthPeriod, useLb, updatedLbList, configList, azStrategy, "enable", deleteProtection),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName1, "id"),
					resource.TestCheckResourceAttr(resourceName1, "name", name),
					resource.TestCheckResourceAttr(resourceName1, "health_mode", healthMode),
					resource.TestCheckResourceAttr(resourceName1, "move_out_strategy", moveOutStrategy),
					resource.TestCheckResourceAttr(resourceName1, "min_count", fmt.Sprintf("%d", minCount)),
					resource.TestCheckResourceAttr(resourceName1, "max_count", fmt.Sprintf("%d", maxCount)),
					resource.TestCheckResourceAttr(resourceName1, "expected_count", fmt.Sprintf("%d", expectedCount)),
					resource.TestCheckResourceAttr(resourceName1, "health_period", fmt.Sprintf("%d", healthPeriod)),
					resource.TestCheckResourceAttr(resourceName1, "use_lb", "1"),
					resource.TestCheckResourceAttr(resourceName1, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName1, "az_strategy", azStrategy),
					resource.TestCheckResourceAttr(resourceName1, "subnet_id_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName1, "lb_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName1, "status", "enable"),
					resource.TestCheckResourceAttr(resourceName1, "delete_protection", deleteProtection),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, securityGroupIDList, name, healthMode, updatedSubnetIDList, moveOutStrategy, vpcId, minCount, maxCount,
					healthPeriod, useLb, updatedLbList, configList, azStrategy, status, deleteProtection),
				Destroy: true,
			},
		}})
}

// 多az，创建时，手动增加云主机；更新时，手动增减云主机
func TestAccCtyunScalingEcs(t *testing.T) {

	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()

	resourceName := "ctyun_scaling_group." + rnd
	resourceFile := "resource_ctyun_scaling_add_ecs_list.tf"

	resourceFile1 := "resource_ctyun_scaling_ecs_list_destroy.tf"
	//datasourceName := "data.ctyun_scalings." + dnd
	//datasourceFile := "datasource_ctyun_scalings.tf"

	securityGroupIDList := fmt.Sprintf(`["%s"]`, dependence.securityGroupID)
	name := "scaling-group-" + utils.GenerateRandomString()
	healthMode := "server"
	subnetIDList := fmt.Sprintf(`["%s"]`, dependence.subnetID)
	moveOutStrategy := "earlier_config"
	vpcId := dependence.vpcID
	minCount := 1
	maxCount := 50

	expectedCount := 1
	healthPeriod := 300
	useLb := 2
	//lbList := fmt.Sprintf(`[{"port": 12306, "lb_id": "%s", "weight": 1, "host_group_id": "%s"}]`, dependence.loadbalancerID, dependence.targetGroupID)
	scalingConfigID, err := strconv.ParseInt(dependence.scalingConfigID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	scalingConfigID1, err := strconv.ParseInt(dependence.scalingConfigID1, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	configList := fmt.Sprintf(`[%d]`, scalingConfigID)
	azStrategy := "uniform_distribution"

	updateSecurityGroupIDList := fmt.Sprintf(`["%s", "%s"]`, dependence.securityGroupID, dependence.securityGroupID1)
	updateName := "scaling-group-" + utils.GenerateRandomString()
	healthModeLb := "lb"
	addInstanceUUIDList := fmt.Sprintf(`["%s","%s"]`, dependence.instanceUUID, dependence.instanceUUID2)

	//updatedSubnetIDList := fmt.Sprintf(`["%s", "%s"]`, dependence.subnetID, dependence.subnetID1)
	updateMoveOutStrategy := "earlier_vm"
	updateMinCount := 1
	updatedMaxCount := 50
	updatedExpectedCount := 4
	updateHealthPeriod := 10080
	updatedConfigList := fmt.Sprintf(`[%d, %d]`, scalingConfigID, scalingConfigID1)
	updatedAddInstanceUUIDList := fmt.Sprintf(`["%s"]`, dependence.instanceUUID1)
	updatedRemoveInstanceUUIDList := fmt.Sprintf(`["%s"]`, dependence.instanceUUID2)
	azStrategyPriorityDistribution := "priority_distribution"
	//isDestroy := true
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
			// 创建弹性组
			{
				Config: utils.LoadTestCase(resourceFile, rnd, securityGroupIDList, name, healthMode, subnetIDList, moveOutStrategy, vpcId, minCount, maxCount, expectedCount,
					healthPeriod, useLb, configList, addInstanceUUIDList, azStrategy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "health_mode", healthMode),
					resource.TestCheckResourceAttr(resourceName, "move_out_strategy", moveOutStrategy),
					resource.TestCheckResourceAttr(resourceName, "min_count", fmt.Sprintf("%d", minCount)),
					resource.TestCheckResourceAttr(resourceName, "max_count", fmt.Sprintf("%d", maxCount)),
					resource.TestCheckResourceAttr(resourceName, "expected_count", fmt.Sprintf("%d", expectedCount)),
					resource.TestCheckResourceAttr(resourceName, "health_period", fmt.Sprintf("%d", healthPeriod)),
					resource.TestCheckResourceAttr(resourceName, "use_lb", "2"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
				),
			},
			// 更新name-> new; health_mode:云主机健康检查->云主机健康检查, subnet不修改，移除策略 earlier_config->earlier_vm，min_count: 1->2, max_count: 1->2， expected_count: 2->2, health_period:300->10080,
			// configList 增加一个config配置，az_strategy : uniform_distribution-> priority_distribution
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, updateSecurityGroupIDList, updateName, healthModeLb, subnetIDList, updateMoveOutStrategy, vpcId, updateMinCount, updatedMaxCount, updatedExpectedCount,
					updateHealthPeriod, useLb, updatedConfigList, updatedAddInstanceUUIDList, updatedRemoveInstanceUUIDList, azStrategyPriorityDistribution),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "health_mode", healthModeLb),
					resource.TestCheckResourceAttr(resourceName, "move_out_strategy", updateMoveOutStrategy),
					resource.TestCheckResourceAttr(resourceName, "min_count", fmt.Sprintf("%d", updateMinCount)),
					resource.TestCheckResourceAttr(resourceName, "max_count", fmt.Sprintf("%d", updatedMaxCount)),
					resource.TestCheckResourceAttr(resourceName, "expected_count", fmt.Sprintf("%d", updatedExpectedCount)),
					resource.TestCheckResourceAttr(resourceName, "health_period", fmt.Sprintf("%d", updateHealthPeriod)),
					resource.TestCheckResourceAttr(resourceName, "use_lb", "2"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
				),
			},
			// 销毁
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, updateSecurityGroupIDList, updateName, healthModeLb, subnetIDList, updateMoveOutStrategy, vpcId, updateMinCount, updatedMaxCount, updatedExpectedCount,
					updateHealthPeriod, useLb, updatedConfigList, updatedAddInstanceUUIDList, updatedRemoveInstanceUUIDList, azStrategyPriorityDistribution),
				Destroy: true,
			},
		}},
	)
}

// 多az，创建时，不添加云主机；更新时，手动增减云主机
func TestAccCtyunScalingEcs1(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_scaling_group." + rnd
	resourceFile1 := "resource_ctyun_scaling.tf"
	resourceFile := "resource_ctyun_scaling_ecs_list1.tf"

	//datasourceName := "data.ctyun_scaling_ecs_list." + dnd
	datasourceFile := "datasource_ecs_list.tf"

	//datasourceName1 := "data.ctyun_scaling_activities." + dnd
	datasourceFile1 := "datasource_scaling_activities.tf"

	ecsResourceName := "ctyun_scaling_ecs_protection." + rnd
	ecsResourceFile := "resource_ctyun_scaling_ecs_protection.tf"

	securityGroupIDList := fmt.Sprintf(`["%s"]`, dependence.securityGroupID)
	name := "scaling-group-" + utils.GenerateRandomString()
	healthMode := "server"
	subnetIDList := fmt.Sprintf(`["%s"]`, dependence.subnetID)
	moveOutStrategy := "earlier_config"
	vpcId := dependence.vpcID
	minCount := 1
	maxCount := 1

	expectedCount := 1
	healthPeriod := 300
	useLb := 1
	lbList := fmt.Sprintf(`[{"port": 12306, "lb_id": "%s", "weight": 1, "host_group_id": "%s"}]`, dependence.loadbalancerID, dependence.targetGroupID)
	scalingConfigID, err := strconv.ParseInt(dependence.scalingConfigID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	scalingConfigID1, err := strconv.ParseInt(dependence.scalingConfigID1, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	configList := fmt.Sprintf(`[%d]`, scalingConfigID)
	azStrategy := "uniform_distribution"

	updateSecurityGroupIDList := fmt.Sprintf(`["%s", "%s"]`, dependence.securityGroupID, dependence.securityGroupID1)
	updateName := "scaling-group-" + utils.GenerateRandomString()
	healthModeLb := "lb"

	//updatedSubnetIDList := fmt.Sprintf(`["%s", "%s"]`, dependence.subnetID, dependence.subnetID1)
	updateMoveOutStrategy := "earlier_vm"
	updateMinCount := 1
	updatedMaxCount := 50
	updatedExpectedCount := 1
	updateHealthPeriod := 10080
	updatedConfigList := fmt.Sprintf(`[%d, %d]`, scalingConfigID, scalingConfigID1)
	updatedInstanceUUIDList := fmt.Sprintf(`["%s", "%s"]`, dependence.instanceUUID1, dependence.instanceUUID)
	//updatedProtectStatus := "enable"
	protectStatus := true
	azStrategyPriorityDistribution := "priority_distribution"

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
			// 创建弹性组
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, securityGroupIDList, name, healthMode, subnetIDList, moveOutStrategy, vpcId, minCount, maxCount, expectedCount,
					healthPeriod, useLb, lbList, configList, azStrategy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "health_mode", healthMode),
					resource.TestCheckResourceAttr(resourceName, "move_out_strategy", moveOutStrategy),
					resource.TestCheckResourceAttr(resourceName, "min_count", fmt.Sprintf("%d", minCount)),
					resource.TestCheckResourceAttr(resourceName, "max_count", fmt.Sprintf("%d", maxCount)),
					resource.TestCheckResourceAttr(resourceName, "expected_count", fmt.Sprintf("%d", expectedCount)),
					resource.TestCheckResourceAttr(resourceName, "health_period", fmt.Sprintf("%d", healthPeriod)),
					resource.TestCheckResourceAttr(resourceName, "use_lb", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "az_strategy", azStrategy),
				),
			},
			// 启动弹性伸缩，其他都不做
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updateSecurityGroupIDList, updateName, healthModeLb, subnetIDList, updateMoveOutStrategy, vpcId, updateMinCount, updatedMaxCount,
					updateHealthPeriod, useLb, lbList, updatedConfigList, updatedInstanceUUIDList, azStrategyPriorityDistribution),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "health_mode", healthModeLb),
					resource.TestCheckResourceAttr(resourceName, "move_out_strategy", updateMoveOutStrategy),
					resource.TestCheckResourceAttr(resourceName, "min_count", fmt.Sprintf("%d", updateMinCount)),
					resource.TestCheckResourceAttr(resourceName, "max_count", fmt.Sprintf("%d", updatedMaxCount)),
					resource.TestCheckResourceAttr(resourceName, "expected_count", fmt.Sprintf("%d", updatedExpectedCount)),
					resource.TestCheckResourceAttr(resourceName, "health_period", fmt.Sprintf("%d", updateHealthPeriod)),
					resource.TestCheckResourceAttr(resourceName, "use_lb", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
				),
			},
			// 对机器开启保护
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updateSecurityGroupIDList, updateName, healthModeLb, subnetIDList, updateMoveOutStrategy, vpcId, updateMinCount, updatedMaxCount,
					updateHealthPeriod, useLb, lbList, updatedConfigList, updatedInstanceUUIDList, azStrategyPriorityDistribution) +
					utils.LoadTestCase(ecsResourceFile, rnd, fmt.Sprintf(`%s.id`, resourceName), updatedInstanceUUIDList, protectStatus),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(ecsResourceName, "region_id"),
					resource.TestCheckResourceAttr(ecsResourceName, "protect_status", "true"),
				),
			},
			// 对机器关闭保护
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updateSecurityGroupIDList, updateName, healthModeLb, subnetIDList, updateMoveOutStrategy, vpcId, updateMinCount, updatedMaxCount,
					updateHealthPeriod, useLb, lbList, updatedConfigList, updatedInstanceUUIDList, azStrategyPriorityDistribution) +
					utils.LoadTestCase(ecsResourceFile, rnd, fmt.Sprintf(`%s.id`, resourceName), updatedInstanceUUIDList, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(ecsResourceName, "region_id"),
					resource.TestCheckResourceAttr(ecsResourceName, "protect_status", "false"),
				),
			},
			// 验证datasource, group下esc数量
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updateSecurityGroupIDList, updateName, healthModeLb, subnetIDList, updateMoveOutStrategy, vpcId, updateMinCount, updatedMaxCount,
					updateHealthPeriod, useLb, lbList, updatedConfigList, updatedInstanceUUIDList, azStrategyPriorityDistribution) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("%s.id", resourceName)),
				//Check: resource.ComposeAggregateTestCheckFunc(
				//	resource.TestCheckResourceAttrSet(datasourceName, "ecs_list")),
			},
			// 验证伸缩活动的datasource
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updateSecurityGroupIDList, updateName, healthModeLb, subnetIDList, updateMoveOutStrategy, vpcId, updateMinCount, updatedMaxCount,
					updateHealthPeriod, useLb, lbList, updatedConfigList, updatedInstanceUUIDList, azStrategyPriorityDistribution) +
					utils.LoadTestCase(datasourceFile1, dnd, fmt.Sprintf("%s.id", resourceName)),
				//Check: resource.ComposeAggregateTestCheckFunc(
				//	resource.TestCheckResourceAttrSet(datasourceName1, "scaling_activities")),
			},
			// 销毁
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updateSecurityGroupIDList, updateName, healthModeLb, subnetIDList, updateMoveOutStrategy, vpcId, updateMinCount, updatedMaxCount,
					updateHealthPeriod, useLb, lbList, updatedConfigList, updatedInstanceUUIDList, azStrategyPriorityDistribution),
				Destroy: true,
			},
		}},
	)
}

// 不填写expected count
func TestAccCtyunScalingNoneExpectedCount(t *testing.T) {

	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_scaling_group." + rnd
	resourceFile := "resource_ctyun_scaling2.tf"

	securityGroupIDList := fmt.Sprintf(`["%s"]`, dependence.securityGroupID)
	name := "scaling-group-" + utils.GenerateRandomString()
	healthMode := "lb"
	subnetIDList := fmt.Sprintf(`["%s"]`, dependence.subnetID)
	moveOutStrategy := "later_config"
	vpcId := dependence.vpcID
	minCount := 4
	maxCount := 50

	healthPeriod := 300
	useLb := 1
	lbList := fmt.Sprintf(`[{"port": 12306, "lb_id": "%s", "weight": 1, "host_group_id": "%s"}]`, dependence.loadbalancerID, dependence.targetGroupID)
	scalingConfigID, err := strconv.ParseInt(dependence.scalingConfigID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	scalingConfigID1, err := strconv.ParseInt(dependence.scalingConfigID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	configList := fmt.Sprintf(`[%d, %d]`, scalingConfigID, scalingConfigID1)
	azStrategy := "uniform_distribution"
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
			// 创建弹性组
			{
				Config: utils.LoadTestCase(resourceFile, rnd, securityGroupIDList, name, healthMode, subnetIDList, moveOutStrategy, vpcId, minCount, maxCount,
					healthPeriod, useLb, lbList, configList, azStrategy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "health_mode", healthMode),
					resource.TestCheckResourceAttr(resourceName, "move_out_strategy", moveOutStrategy),
					resource.TestCheckResourceAttr(resourceName, "min_count", fmt.Sprintf("%d", minCount)),
					resource.TestCheckResourceAttr(resourceName, "max_count", fmt.Sprintf("%d", maxCount)),
					resource.TestCheckResourceAttr(resourceName, "health_period", fmt.Sprintf("%d", healthPeriod)),
					resource.TestCheckResourceAttr(resourceName, "use_lb", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "az_strategy", azStrategy),
				),
			},
			//  资源导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					projectId := ds.Attributes["project_id"]
					if id == "" || regionId == "" {
						return "", fmt.Errorf("id or region_id is required")
					}
					return fmt.Sprintf("%s,%s,%s,%s", id, regionId, projectId, vpcId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"add_instance_uuid_list", "remove_instance_uuid_list", "is_destroy", "expected_count"},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, securityGroupIDList, name, healthMode, subnetIDList, moveOutStrategy, vpcId, minCount, maxCount,
					healthPeriod, useLb, lbList, configList, azStrategy),
				Destroy: true,
			},
		}})
}
