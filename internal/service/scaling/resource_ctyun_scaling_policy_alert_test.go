package scaling_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
	"time"
)

func TestAccCtyunScalingPolicyAlert(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_scaling_policy." + rnd
	resourceFile := "resource_ctyun_scaling_policy_alert.tf"
	updatedResourceFile := "resource_ctyun_scaling_policy_alert_update.tf"

	datasourceFile := "datasource_ctyun_scaling_policies.tf"
	datasourceName := "data.ctyun_scaling_policies." + dnd

	groupID, err := strconv.ParseInt(dependence.scalingGroupID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	// 创建参数
	name := "test-alert-policy-" + rnd
	policyType := "alert"
	operateUnit := "percent"
	operateCount := 10
	action := "increase"
	cooldown := 300
	triggerName := "cpu-high-alert-" + rnd
	triggerMetricName := "cpu_util"
	triggerStatistics := "avg"
	triggerComparisonOperator := "ge"
	triggerThreshold := 80
	triggerPeriod := "5m"
	triggerEvaluationCount := 3

	// 更新参数
	updatedName := "updated-alert-policy-" + rnd
	updatedOperateUnit := "count"
	updatedAction := "set"
	updatedTriggerName := "memory-high-alert-" + rnd
	updatedTriggerMetricName := "mem_util"
	updatedTriggerStatistics := "max"
	updatedTriggerComparisonOperator := "gt"
	updatedOperateCount := 15
	updatedCooldown := 600
	updatedTriggerPeriod := "1h"
	updatedTriggerThreshold := 85
	updatedTriggerEvaluationCount := 5
	resource.Test(t, resource.TestCase{
		// 检查资源是否被销毁
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource %s still exists", resourceName)
			}
			return nil
		},
		// 使用ProtoV6ProviderFactories
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 创建伸缩策略（告警类型）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, groupID, name, policyType, operateUnit, operateCount, action, cooldown, triggerName,
					triggerMetricName, triggerStatistics, triggerComparisonOperator, triggerThreshold, triggerPeriod, triggerEvaluationCount),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "policy_type", policyType),
					resource.TestCheckResourceAttr(resourceName, "operate_unit", operateUnit),
					resource.TestCheckResourceAttr(resourceName, "operate_count", fmt.Sprintf("%d", operateCount)),
					resource.TestCheckResourceAttr(resourceName, "action", action),
					resource.TestCheckResourceAttr(resourceName, "cooldown", fmt.Sprintf("%d", cooldown)),
					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "trigger_metric_name", triggerMetricName),
					resource.TestCheckResourceAttr(resourceName, "trigger_statistics", triggerStatistics),
					resource.TestCheckResourceAttr(resourceName, "trigger_comparison_operator", triggerComparisonOperator),
					resource.TestCheckResourceAttr(resourceName, "trigger_threshold", fmt.Sprintf("%d", triggerThreshold)),
					resource.TestCheckResourceAttr(resourceName, "trigger_period", triggerPeriod),
					resource.TestCheckResourceAttr(resourceName, "trigger_evaluation_count", fmt.Sprintf("%d", triggerEvaluationCount)),
					resource.TestCheckResourceAttr(resourceName, "status", "enable"),
				),
			},
			{
				Config: utils.LoadTestCase(updatedResourceFile, rnd, groupID, updatedName, policyType, updatedOperateUnit, updatedOperateCount, updatedAction, updatedCooldown, updatedTriggerName,
					updatedTriggerMetricName, updatedTriggerStatistics, updatedTriggerComparisonOperator, updatedTriggerThreshold, updatedTriggerPeriod, updatedTriggerEvaluationCount, "disable"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "operate_unit", updatedOperateUnit),
					resource.TestCheckResourceAttr(resourceName, "operate_count", fmt.Sprintf("%d", updatedOperateCount)),
					resource.TestCheckResourceAttr(resourceName, "action", updatedAction),
					resource.TestCheckResourceAttr(resourceName, "cooldown", fmt.Sprintf("%d", updatedCooldown)),
					resource.TestCheckResourceAttr(resourceName, "trigger_name", updatedTriggerName),
					resource.TestCheckResourceAttr(resourceName, "trigger_metric_name", updatedTriggerMetricName),
					resource.TestCheckResourceAttr(resourceName, "trigger_statistics", updatedTriggerStatistics),
					resource.TestCheckResourceAttr(resourceName, "trigger_comparison_operator", updatedTriggerComparisonOperator),
					resource.TestCheckResourceAttr(resourceName, "trigger_threshold", fmt.Sprintf("%d", updatedTriggerThreshold)),
					resource.TestCheckResourceAttr(resourceName, "trigger_period", updatedTriggerPeriod),
					resource.TestCheckResourceAttr(resourceName, "trigger_evaluation_count", fmt.Sprintf("%d", updatedTriggerEvaluationCount)),
					resource.TestCheckResourceAttr(resourceName, "status", "disable"),
				),
			},
			// datasource 验证
			{
				Config: utils.LoadTestCase(updatedResourceFile, rnd, groupID, updatedName, policyType, updatedOperateUnit, updatedOperateCount, updatedAction, updatedCooldown, updatedTriggerName,
					updatedTriggerMetricName, updatedTriggerStatistics, updatedTriggerComparisonOperator, updatedTriggerThreshold, updatedTriggerPeriod, updatedTriggerEvaluationCount, "disable") +
					utils.LoadTestCase(datasourceFile, dnd, groupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "scaling_policies.#"),
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
					groupId := ds.Attributes["group_id"]
					if id == "" || regionId == "" {
						return "", fmt.Errorf("id or region_id is required")
					}
					return fmt.Sprintf("%s,%s,%s,%s", id, regionId, groupId, policyType), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_execute", "target_disable_scale_in"},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, groupID, updatedName, policyType, operateUnit, operateCount, action, updatedCooldown, triggerName,
					triggerMetricName, triggerStatistics, triggerComparisonOperator, updatedTriggerThreshold, triggerPeriod, updatedTriggerEvaluationCount),
				Destroy: true,
			},
		}})
}

func TestAccCtyunScalingPolicyRegular(t *testing.T) {
	
	// 生成随机名称避免冲突
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_scaling_policy." + rnd

	resourceFile := "resource_ctyun_scaling_policy_regular.tf"
	updatedResourceFile := "resource_ctyun_scaling_policy_regular_update.tf"

	datasourceFile := "datasource_ctyun_scaling_policies.tf"
	datasourceName := "data.ctyun_scaling_policies." + dnd

	groupID, err := strconv.ParseInt(dependence.scalingGroupID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	// 测试依赖项（需要提前创建好）

	// 创建参数
	name := "test-regular-policy-" + rnd
	policyType := "regular"
	operateUnit := "count"
	operateCount := 2
	action := "increase"
	executionTime := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05") // 明天执行

	// 更新参数
	updatedName := "updated-regular-policy-" + rnd
	updatedOperateUnit := "percent"
	updatedOperateCount := 10
	updatedAction := "decrease"
	updatedExecutionTime := time.Now().Add(48 * time.Hour).Format("2006-01-02 15:04:05") // 后天执行

	resource.Test(t, resource.TestCase{
		// 检查资源是否被销毁
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource %s still exists", resourceName)
			}
			return nil
		},
		// 使用ProtoV6ProviderFactories
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 创建定时策略
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					groupID, name, policyType, operateUnit, operateCount, action, executionTime,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "policy_type", policyType),
					resource.TestCheckResourceAttr(resourceName, "operate_unit", operateUnit),
					resource.TestCheckResourceAttr(resourceName, "operate_count", fmt.Sprintf("%d", operateCount)),
					resource.TestCheckResourceAttr(resourceName, "action", action),
					resource.TestCheckResourceAttr(resourceName, "execution_time", executionTime),
				),
			},
			// 更新所有字段
			{
				Config: utils.LoadTestCase(updatedResourceFile, rnd,
					groupID, updatedName, policyType, updatedOperateUnit, updatedOperateCount, updatedAction, updatedExecutionTime, "disable",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "operate_unit", updatedOperateUnit),
					resource.TestCheckResourceAttr(resourceName, "operate_count", fmt.Sprintf("%d", updatedOperateCount)),
					resource.TestCheckResourceAttr(resourceName, "action", updatedAction),
					resource.TestCheckResourceAttr(resourceName, "execution_time", updatedExecutionTime),
				),
			},
			{
				// datasource 验证
				Config: utils.LoadTestCase(updatedResourceFile, rnd,
					groupID, updatedName, policyType, updatedOperateUnit, updatedOperateCount, updatedAction, updatedExecutionTime, "disable") +
					utils.LoadTestCase(datasourceFile, dnd, groupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "scaling_policies.#"),
				),
			},
			{

				Config: utils.LoadTestCase(updatedResourceFile, rnd,
					groupID, updatedName, policyType, updatedOperateUnit, updatedOperateCount, updatedAction, updatedExecutionTime, "disable",
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunScalingPolicyPeriod(t *testing.T) {

	// 生成随机名称避免冲突
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_scaling_policy." + rnd

	resourceFile := "resource_ctyun_scaling_policy_period.tf"
	updatedResourceFile := "resource_ctyun_scaling_policy_period_update.tf"

	datasourceFile := "datasource_ctyun_scaling_policies.tf"
	datasourceName := "data.ctyun_scaling_policies." + dnd

	// 测试依赖项（需要提前创建好）
	groupID, err := strconv.ParseInt(dependence.scalingGroupID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	// 创建参数
	name := "test-period-policy-" + rnd
	policyType := "period"
	operateUnit := "count"
	operateCount := 2
	action := "decrease"
	cycle := "monthly"
	day := "[1, 31]" // 每月1号和15号执行
	effectiveFrom := time.Now().AddDate(0, 0, 1).Format("2006-01-02 15:04:05")
	effectiveTill := time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05") // 一年后
	executionTime := time.Now().AddDate(0, 0, 2).Format("2006-01-02 15:04:05") // 明天执行

	// 更新参数
	updatedName := "updated-period-policy-" + rnd
	updatedOperateUnit := "percent"
	updatedOperateCount := 10
	updatedAction := "decrease"
	updatedCycle := "weekly"
	updatedDay := "[1, 4, 6]"                                                         // 每周二、四、六执行（周一为1，周日为7）
	updatedEffectiveFrom := time.Now().AddDate(0, 1, 1).Format("2006-01-02 15:04:05") // 一个月后
	updatedEffectiveTill := time.Now().AddDate(1, 1, 0).Format("2006-01-02 15:04:05") // 一年一个月后
	updatedExecutionTime := time.Now().AddDate(0, 1, 2).Format("2006-01-02 15:04:05") // 后天执行

	resource.Test(t, resource.TestCase{
		// 检查资源是否被销毁
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource %s still exists", resourceName)
			}
			return nil
		},
		// 使用ProtoV6ProviderFactories
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 创建周期策略
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					groupID, name, policyType, operateUnit, operateCount, action, cycle, day, effectiveFrom, effectiveTill, executionTime,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "policy_type", policyType),
					resource.TestCheckResourceAttr(resourceName, "operate_unit", operateUnit),
					resource.TestCheckResourceAttr(resourceName, "operate_count", fmt.Sprintf("%d", operateCount)),
					resource.TestCheckResourceAttr(resourceName, "action", action),
					resource.TestCheckResourceAttr(resourceName, "cycle", cycle),
					resource.TestCheckResourceAttr(resourceName, "day.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "effective_from", effectiveFrom),
					resource.TestCheckResourceAttr(resourceName, "effective_till", effectiveTill),
					resource.TestCheckResourceAttr(resourceName, "execution_time", executionTime),
				),
			},
			// 更新所有字段
			{
				Config: utils.LoadTestCase(updatedResourceFile, rnd,
					groupID, updatedName, policyType, updatedOperateUnit, updatedOperateCount, updatedAction, updatedCycle, updatedDay, updatedEffectiveFrom, updatedEffectiveTill, updatedExecutionTime, "disable",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "operate_unit", updatedOperateUnit),
					resource.TestCheckResourceAttr(resourceName, "operate_count", fmt.Sprintf("%d", updatedOperateCount)),
					resource.TestCheckResourceAttr(resourceName, "action", updatedAction),
					resource.TestCheckResourceAttr(resourceName, "cycle", updatedCycle),
					resource.TestCheckResourceAttr(resourceName, "day.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "effective_from", updatedEffectiveFrom),
					resource.TestCheckResourceAttr(resourceName, "effective_till", updatedEffectiveTill),
					resource.TestCheckResourceAttr(resourceName, "execution_time", updatedExecutionTime),
				),
			},
			{
				// datasource 验证
				Config: utils.LoadTestCase(updatedResourceFile, rnd,
					groupID, updatedName, policyType, updatedOperateUnit, updatedOperateCount, updatedAction, updatedCycle, updatedDay, updatedEffectiveFrom, updatedEffectiveTill, updatedExecutionTime, "disable",
				) + utils.LoadTestCase(datasourceFile, dnd, groupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "scaling_policies.#"),
				),
			},
			// 销毁资源（通过空配置）
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					groupID, updatedName, policyType, updatedOperateUnit, updatedOperateCount, updatedAction, updatedCycle, updatedDay, updatedEffectiveFrom, updatedEffectiveTill,
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunScalingPolicyPeriodDay(t *testing.T) {

	// 生成随机名称避免冲突
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_scaling_policy." + rnd

	resourceFile := "resource_ctyun_scaling_policy_period_daily.tf"
	updatedResourceFile := "resource_ctyun_scaling_policy_period_daily_update.tf"
	// 测试依赖项（需要提前创建好）
	groupID, err := strconv.ParseInt(dependence.scalingGroupID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	// 创建参数
	name := "test-period-policy-" + rnd
	policyType := "period"
	operateUnit := "count"
	operateCount := 2
	action := "decrease"
	cycle := "daily"
	//day := "[1, 15]" // 每月1号和15号执行
	effectiveFrom := time.Now().Add(1 * time.Hour).Format("2006-01-02 15:04:05")
	effectiveTill := time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05")   // 一年后
	executionTime := time.Now().Add(2 * time.Hour).Format("2006-01-02 15:04:05") // 2小时后执行

	// 更新参数
	updatedName := "updated-period-policy-" + rnd
	updatedOperateUnit := "percent"
	updatedOperateCount := 10
	updatedAction := "decrease"
	//updatedCycle := "weekly"
	//updatedDay := "[2, 4, 6]"                                                         // 每周二、四、六执行（周一为1，周日为7）
	updatedEffectiveFrom := time.Now().Add(2 * time.Hour).Format("2006-01-02 15:04:05")
	updatedEffectiveTill := time.Now().AddDate(1, 1, 0).Format("2006-01-02 15:04:05")   // 一年一个月后
	updatedExecutionTime := time.Now().Add(3 * time.Hour).Format("2006-01-02 15:04:05") // 3小时后执行

	resource.Test(t, resource.TestCase{
		// 检查资源是否被销毁
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource %s still exists", resourceName)
			}
			return nil
		},
		// 使用ProtoV6ProviderFactories
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 创建周期策略
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					groupID, name, policyType, operateUnit, operateCount, action, cycle, effectiveFrom, effectiveTill, executionTime,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "policy_type", policyType),
					resource.TestCheckResourceAttr(resourceName, "operate_unit", operateUnit),
					resource.TestCheckResourceAttr(resourceName, "operate_count", fmt.Sprintf("%d", operateCount)),
					resource.TestCheckResourceAttr(resourceName, "action", action),
					resource.TestCheckResourceAttr(resourceName, "cycle", cycle),
					resource.TestCheckResourceAttr(resourceName, "effective_from", effectiveFrom),
					resource.TestCheckResourceAttr(resourceName, "effective_till", effectiveTill),
					resource.TestCheckResourceAttr(resourceName, "execution_time", executionTime),
				),
			},
			// 更新所有字段
			{
				Config: utils.LoadTestCase(updatedResourceFile, rnd,
					groupID, updatedName, policyType, updatedOperateUnit, updatedOperateCount, updatedAction, cycle, updatedEffectiveFrom, updatedEffectiveTill, updatedExecutionTime, "disable",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "operate_unit", updatedOperateUnit),
					resource.TestCheckResourceAttr(resourceName, "operate_count", fmt.Sprintf("%d", updatedOperateCount)),
					resource.TestCheckResourceAttr(resourceName, "action", updatedAction),
					resource.TestCheckResourceAttr(resourceName, "cycle", cycle),
					resource.TestCheckResourceAttr(resourceName, "effective_from", updatedEffectiveFrom),
					resource.TestCheckResourceAttr(resourceName, "effective_till", updatedEffectiveTill),
					resource.TestCheckResourceAttr(resourceName, "execution_time", updatedExecutionTime),
				),
			},
			{
				Config: utils.LoadTestCase(updatedResourceFile, rnd,
					groupID, updatedName, policyType, updatedOperateUnit, updatedOperateCount, updatedAction, cycle, updatedEffectiveFrom, updatedEffectiveTill, updatedExecutionTime, "enable",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "operate_unit", updatedOperateUnit),
					resource.TestCheckResourceAttr(resourceName, "operate_count", fmt.Sprintf("%d", updatedOperateCount)),
					resource.TestCheckResourceAttr(resourceName, "action", updatedAction),
					resource.TestCheckResourceAttr(resourceName, "cycle", cycle),
					resource.TestCheckResourceAttr(resourceName, "effective_from", updatedEffectiveFrom),
					resource.TestCheckResourceAttr(resourceName, "effective_till", updatedEffectiveTill),
					resource.TestCheckResourceAttr(resourceName, "execution_time", updatedExecutionTime),
				),
			},
			// 销毁资源（通过空配置）
			{
				Config: utils.LoadTestCase(updatedResourceFile, rnd,
					groupID, updatedName, policyType, updatedOperateUnit, updatedOperateCount, updatedAction, cycle, updatedEffectiveFrom, updatedEffectiveTill, updatedExecutionTime, "enable",
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunScalingPolicyTarget(t *testing.T) {

	// 生成随机名称避免冲突
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_scaling_policy." + rnd

	resourceFile := "resource_ctyun_scaling_policy_target.tf"
	updatedResourceFile := "resource_ctyun_scaling_policy_target_update.tf"

	datasourceFile := "datasource_ctyun_scaling_policies.tf"
	datasourceName := "data.ctyun_scaling_policies." + dnd

	// 测试依赖项（需要提前创建好）
	groupID, err := strconv.ParseInt(dependence.scalingGroupID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	// 创建参数
	name := "test-target-policy-" + rnd
	policyType := "target"
	cooldown := 300
	targetMetricName := "cpu_util"
	targetValue := 50
	scaleOutEvaluationCount := 3
	scaleInEvaluationCount := 10
	operateRange := 10
	disableScaleIn := false

	// 更新参数
	updatedName := "updated-target-policy-" + rnd
	updatedCooldown := 600
	updatedTargetMetricName := "network_incoming_bytes_rate_inband"
	updatedTargetValue := 100
	updatedScaleOutEvaluationCount := 5
	updatedScaleInEvaluationCount := 15
	updatedOperateRange := 15
	updatedDisableScaleIn := true

	resource.Test(t, resource.TestCase{
		// 检查资源是否被销毁
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource %s still exists", resourceName)
			}
			return nil
		},
		// 使用ProtoV6ProviderFactories
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 创建目标追踪策略
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					groupID, name, policyType, cooldown, targetMetricName, targetValue,
					scaleOutEvaluationCount, scaleInEvaluationCount, operateRange, disableScaleIn,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "policy_type", policyType),
					resource.TestCheckResourceAttr(resourceName, "cooldown", fmt.Sprintf("%d", cooldown)),
					resource.TestCheckResourceAttr(resourceName, "target_metric_name", targetMetricName),
					resource.TestCheckResourceAttr(resourceName, "target_value", fmt.Sprintf("%d", targetValue)),
					resource.TestCheckResourceAttr(resourceName, "target_scale_out_evaluation_count", fmt.Sprintf("%d", scaleOutEvaluationCount)),
					resource.TestCheckResourceAttr(resourceName, "target_scale_in_evaluation_count", fmt.Sprintf("%d", scaleInEvaluationCount)),
					resource.TestCheckResourceAttr(resourceName, "target_operate_range", fmt.Sprintf("%d", operateRange)),
					resource.TestCheckResourceAttr(resourceName, "target_disable_scale_in", fmt.Sprintf("%t", disableScaleIn)),
				),
			},
			// 更新所有字段
			{
				Config: utils.LoadTestCase(updatedResourceFile, rnd, groupID, updatedName, policyType, updatedCooldown, updatedTargetMetricName, updatedTargetValue,
					updatedScaleOutEvaluationCount, updatedScaleInEvaluationCount, updatedOperateRange, updatedDisableScaleIn, "disable",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cooldown", fmt.Sprintf("%d", updatedCooldown)),
					resource.TestCheckResourceAttr(resourceName, "target_metric_name", updatedTargetMetricName),
					resource.TestCheckResourceAttr(resourceName, "target_value", fmt.Sprintf("%d", updatedTargetValue)),
					resource.TestCheckResourceAttr(resourceName, "target_scale_out_evaluation_count", fmt.Sprintf("%d", updatedScaleOutEvaluationCount)),
					resource.TestCheckResourceAttr(resourceName, "target_scale_in_evaluation_count", fmt.Sprintf("%d", updatedScaleInEvaluationCount)),
					resource.TestCheckResourceAttr(resourceName, "target_operate_range", fmt.Sprintf("%d", updatedOperateRange)),
					resource.TestCheckResourceAttr(resourceName, "target_disable_scale_in", fmt.Sprintf("%t", updatedDisableScaleIn)),
				),
			},
			// datasource验证
			{
				Config: utils.LoadTestCase(updatedResourceFile, rnd, groupID, updatedName, policyType, updatedCooldown, updatedTargetMetricName, updatedTargetValue,
					updatedScaleOutEvaluationCount, updatedScaleInEvaluationCount, updatedOperateRange, updatedDisableScaleIn, "disable",
				) + utils.LoadTestCase(datasourceFile, dnd, groupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "scaling_policies.#"),
					//resource.TestCheckResourceAttr(datasourceName, "scaling_policies.#", "1"),
					//resource.TestCheckResourceAttr(datasourceName, "scaling_policies.0.name", updatedName),
					//resource.TestCheckResourceAttr(datasourceName, "scaling_policies.0.policy_type", policyType),
					//resource.TestCheckResourceAttr(datasourceName, "scaling_policies.0.cooldown", fmt.Sprintf("%d", updatedCooldown)),
					//resource.TestCheckResourceAttr(datasourceName, "scaling_policies.0.target_metric_name", updatedTargetMetricName),
					//resource.TestCheckResourceAttr(datasourceName, "scaling_policies.0.target_value", fmt.Sprintf("%d", updatedTargetValue)),
					//resource.TestCheckResourceAttr(datasourceName, "scaling_policies.0.target_scale_out_evaluation_count", fmt.Sprintf("%d", updatedScaleOutEvaluationCount)),
					//resource.TestCheckResourceAttr(datasourceName, "scaling_policies.0.target_scale_in_evaluation_count", fmt.Sprintf("%d", updatedScaleInEvaluationCount)),
					//resource.TestCheckResourceAttr(datasourceName, "scaling_policies.0.target_operate_range", fmt.Sprintf("%d", updatedOperateRange)),
					//resource.TestCheckResourceAttr(datasourceName, "scaling_policies.0.target_disable_scale_in", fmt.Sprintf("%t", updatedDisableScaleIn)),
				),
			},
			// 销毁资源（通过空配置）
			{
				Config: utils.LoadTestCase(updatedResourceFile, rnd, groupID, updatedName, policyType, updatedCooldown, updatedTargetMetricName, updatedTargetValue,
					updatedScaleOutEvaluationCount, updatedScaleInEvaluationCount, updatedOperateRange, updatedDisableScaleIn, "disable",
				),
				Destroy: true,
			},
		},
	})
}
