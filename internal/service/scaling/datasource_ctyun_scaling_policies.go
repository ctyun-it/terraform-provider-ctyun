package scaling

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/scaling"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunScalingPolicies{}
	_ datasource.DataSourceWithConfigure = &CtyunScalingPolicies{}
)

type CtyunScalingPolicies struct {
	meta *common.CtyunMetadata
}

func NewCtyunScalingPolicies() datasource.DataSource {
	return &CtyunScalingPolicies{}
}

func (c *CtyunScalingPolicies) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_scaling_policies"
}

func (c *CtyunScalingPolicies) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunScalingPolicies) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027725/10241454`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id",
			},
			"group_id": schema.Int64Attribute{
				Required:    true,
				Description: "伸缩组ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "页码",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页行数，取值范围:[1~100]，默认值为10",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
				},
			},
			"scaling_policies": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_id": schema.Int64Attribute{
							Computed:    true,
							Description: "伸缩策略ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "伸缩策略名称",
						},
						"policy_type": schema.StringAttribute{
							Computed:    true,
							Description: "策略类型: alert-告警, regular-定时, period-周期, target-目标追踪",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "启用状态: enable-启用, disable-停用",
						},
						"action": schema.StringAttribute{
							Computed:    true,
							Description: "执行动作: increase-增加, decrease-减少, set-设置为",
						},
						"operate_count": schema.Int32Attribute{
							Computed:    true,
							Description: "调整值",
						},
						"operate_unit": schema.StringAttribute{
							Computed:    true,
							Description: "操作单位: count-个数, percent-百分比",
						},
						"cooldown": schema.Int32Attribute{
							Computed:    true,
							Description: "冷却时间或预热时间 (秒)",
						},
						"execution_time": schema.StringAttribute{
							Computed:    true,
							Description: "触发时间",
						},
						"cycle": schema.StringAttribute{
							Computed:    true,
							Description: "循环方式: monthly-按月循环, weekly-按周循环, daliy-按天循环",
						},
						"effective_from": schema.StringAttribute{
							Computed:    true,
							Description: "周期策略生效开始时间",
						},
						"effective_till": schema.StringAttribute{
							Computed:    true,
							Description: "周期策略生效截止时间",
						},
						"day": schema.SetAttribute{
							ElementType: types.Int64Type,
							Computed:    true,
							Description: "执行日期",
						},
						"group_id": schema.Int32Attribute{
							Computed:    true,
							Description: "伸缩组ID",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "企业项目ID",
						},
						"create_date": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"update_date": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间",
						},
						"trigger_id": schema.StringAttribute{
							Computed:    true,
							Description: "告警规则ID",
						},
						"trigger_name": schema.StringAttribute{
							Computed:    true,
							Description: "告警规则名称",
						},
						"trigger_metric_name": schema.StringAttribute{
							Computed:    true,
							Description: "监控指标名称",
						},
						"trigger_statistics": schema.StringAttribute{
							Computed:    true,
							Description: "聚合方法",
						},
						"trigger_comparison_operator": schema.StringAttribute{
							Computed:    true,
							Description: "比较符",
						},
						"trigger_threshold": schema.Int32Attribute{
							Computed:    true,
							Description: "阈值",
						},
						"trigger_period": schema.StringAttribute{
							Computed:    true,
							Description: "监控周期",
						},
						"trigger_evaluation_count": schema.Int32Attribute{
							Computed:    true,
							Description: "连续出现次数",
						},
						"trigger_cooldown": schema.Int32Attribute{
							Computed:    true,
							Description: "冷却时间 (秒)",
						},
						"trigger_status": schema.Int32Attribute{
							Computed:    true,
							Description: "告警规则状态: 0-启用, 1-停用",
						},
						"target_metric_name": schema.StringAttribute{
							Computed:    true,
							Description: "监控指标名称",
						},
						"target_value": schema.Int32Attribute{
							Computed:    true,
							Description: "追踪目标值",
						},
						"target_scale_out_evaluation_count": schema.Int32Attribute{
							Computed:    true,
							Description: "扩容连续告警次数",
						},
						"target_scale_in_evaluation_count": schema.Int32Attribute{
							Computed:    true,
							Description: "缩容连续告警次数",
						},
						"target_operate_range": schema.Int32Attribute{
							Computed:    true,
							Description: "缩容波动范围",
						},
						"target_disable_scale_in": schema.BoolAttribute{
							Computed:    true,
							Description: "是否禁用缩容",
						},
					},
				},
				Description: "弹性伸缩策略列表",
			},
		},
	}
}

func (c *CtyunScalingPolicies) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunScalingPoliciesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	if regionId == "" {
		err = errors.New("region id 为空")
		return
	}
	params := &scaling.ScalingRuleListRequest{
		RegionID: regionId,
		GroupID:  config.GroupID.ValueInt64(),
		PageNo:   1,
		PageSize: 10,
	}
	if !config.PageNo.IsNull() && !config.PageNo.IsUnknown() {
		params.PageNo = config.PageNo.ValueInt32()
	}
	if !config.PageSize.IsNull() && !config.PageSize.IsUnknown() {
		params.PageSize = config.PageSize.ValueInt32()
	}

	resp, err := c.meta.Apis.SdkScalingApis.ScalingRuleListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询伸缩策略列表失败，接口返回nil。对应的策略组id为：{%d}, 请联系研发确认！", config.GroupID.ValueInt64())
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var scalingPolicies []CtyunScalingPolicyList
	ruleList := resp.ReturnObj.RuleList
	for _, rule := range ruleList {
		scalingPolicy := CtyunScalingPolicyList{}
		scalingPolicy.RuleID = types.Int64Value(rule.RuleID)
		scalingPolicy.Name = types.StringValue(rule.Name)
		scalingPolicy.PolicyType = types.StringValue(business.ScalingPolicyTypeDictRev[rule.RuleType])
		scalingPolicy.Status = types.StringValue(business.ScalingPolicyStatusDictRev[rule.Status])
		scalingPolicy.Action = types.StringValue(business.ActionDictRev[rule.Action])
		scalingPolicy.OperateUnit = types.StringValue(business.OperateUnitDictRev[rule.OperateUnit])
		scalingPolicy.Cooldown = types.Int32Value(rule.Cooldown)
		scalingPolicy.ExecutionTime = types.StringValue(rule.ExecutionTime)
		scalingPolicy.Cycle = types.StringValue(business.CycleDictRev[rule.Cycle])
		scalingPolicy.EffectiveFrom = types.StringValue(rule.EffectiveFrom)
		scalingPolicy.EffectiveTill = types.StringValue(rule.EffectiveTill)
		day, diagnostics := types.SetValueFrom(ctx, types.Int32Type, rule.Day)
		if diagnostics != nil && diagnostics.HasError() {
			err = errors.New(diagnostics[0].Detail())
			return
		}
		scalingPolicy.Day = day
		scalingPolicy.GroupID = types.Int32Value(rule.ScalingGroupID)
		scalingPolicy.ProjectID = types.StringValue(rule.ProjectIDEcs)
		scalingPolicy.CreateDate = types.StringValue(rule.CreateDate)
		scalingPolicy.UpdateDate = types.StringValue(rule.UpdateDate)
		// trigger info
		if rule.TriggerObj != nil {
			scalingPolicy.TriggerID = types.StringValue(rule.TriggerObj.TriggerID)
			scalingPolicy.TriggerName = types.StringValue(rule.TriggerObj.Name)
			scalingPolicy.TriggerMetricName = types.StringValue(rule.TriggerObj.MetricName)
			scalingPolicy.TriggerStatistics = types.StringValue(rule.TriggerObj.Statistics)
			scalingPolicy.TriggerComparisonOperator = types.StringValue(rule.TriggerObj.ComparisonOperator)
			scalingPolicy.TriggerThreshold = types.Int32Value(rule.TriggerObj.Threshold)
			scalingPolicy.TriggerPeriod = types.StringValue(rule.TriggerObj.Period)
			scalingPolicy.TriggerEvaluationCount = types.Int32Value(rule.TriggerObj.EvaluationCount)
			scalingPolicy.TriggerCooldown = types.Int32Value(rule.TriggerObj.Cooldown)
			scalingPolicy.TriggerStatus = types.Int32Value(rule.TriggerObj.Status)
		} else {
			scalingPolicy.TriggerID = types.StringValue("")
			scalingPolicy.TriggerName = types.StringValue("")
			scalingPolicy.TriggerMetricName = types.StringValue("")
			scalingPolicy.TriggerStatistics = types.StringValue("")
			scalingPolicy.TriggerComparisonOperator = types.StringValue("")
			scalingPolicy.TriggerThreshold = types.Int32Value(-1)
			scalingPolicy.TriggerPeriod = types.StringValue("")
			scalingPolicy.TriggerEvaluationCount = types.Int32Value(-1)
			scalingPolicy.TriggerCooldown = types.Int32Value(-1)
			scalingPolicy.TriggerStatus = types.Int32Value(-1)
		}

		// target info
		if rule.TargetObj != nil {
			scalingPolicy.TargetMetricName = types.StringValue(rule.TargetObj.MetricName)
			scalingPolicy.TargetValue = types.Int32Value(rule.TargetObj.TargetValue)
			scalingPolicy.TargetScaleOutEvaluationCount = types.Int32Value(rule.TargetObj.ScaleOutEvaluationCount)
			scalingPolicy.TargetScaleInEvaluationCount = types.Int32Value(rule.TargetObj.ScaleInEvaluationCount)
			scalingPolicy.TargetOperateRange = types.Int32Value(rule.TargetObj.OperateRange)
			scalingPolicy.TargetDisableScaleIn = types.BoolValue(*rule.TargetObj.DisableScaleIn)
		} else {
			scalingPolicy.TargetMetricName = types.StringValue("")
			scalingPolicy.TargetValue = types.Int32Value(-1)
			scalingPolicy.TargetScaleOutEvaluationCount = types.Int32Value(-1)
			scalingPolicy.TargetScaleInEvaluationCount = types.Int32Value(-1)
			scalingPolicy.TargetOperateRange = types.Int32Value(-1)
			scalingPolicy.TargetDisableScaleIn = types.BoolValue(false)
		}

		scalingPolicies = append(scalingPolicies, scalingPolicy)
	}

	config.ScalingPolicies = scalingPolicies
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	return
}

type CtyunScalingPolicyList struct {
	RuleID                        types.Int64  `tfsdk:"rule_id"`                           // 伸缩策略ID
	Name                          types.String `tfsdk:"name"`                              // 伸缩策略名称
	PolicyType                    types.String `tfsdk:"policy_type"`                       // 策略类型: 1-告警, 2-定时, 3-周期, 4-目标追踪
	Status                        types.String `tfsdk:"status"`                            // 启用状态: 1-启用, 2-停用
	Action                        types.String `tfsdk:"action"`                            // 执行动作: 1-增加, 2-减少, 3-设置为
	OperateCount                  types.Int32  `tfsdk:"operate_count"`                     // 调整值
	OperateUnit                   types.String `tfsdk:"operate_unit"`                      // 操作单位: 1-个数, 2-百分比
	Cooldown                      types.Int32  `tfsdk:"cooldown"`                          // 冷却时间或预热时间 (秒)
	ExecutionTime                 types.String `tfsdk:"execution_time"`                    // 触发时间
	Cycle                         types.String `tfsdk:"cycle"`                             // 循环方式: 1-按月循环, 2-按周循环, 3-按天循环
	EffectiveFrom                 types.String `tfsdk:"effective_from"`                    // 周期策略生效开始时间
	EffectiveTill                 types.String `tfsdk:"effective_till"`                    // 周期策略生效截止时间
	Day                           types.Set    `tfsdk:"day"`                               // 执行日期 (列表)
	GroupID                       types.Int32  `tfsdk:"group_id"`                          // 伸缩组ID
	ProjectID                     types.String `tfsdk:"project_id"`                        // 企业项目ID
	CreateDate                    types.String `tfsdk:"create_date"`                       // 创建时间
	UpdateDate                    types.String `tfsdk:"update_date"`                       // 更新时间
	TriggerID                     types.String `tfsdk:"trigger_id"`                        // 告警规则ID
	TriggerName                   types.String `tfsdk:"trigger_name"`                      // 告警规则名称
	TriggerMetricName             types.String `tfsdk:"trigger_metric_name"`               // 监控指标名称
	TriggerStatistics             types.String `tfsdk:"trigger_statistics"`                // 聚合方法
	TriggerComparisonOperator     types.String `tfsdk:"trigger_comparison_operator"`       // 比较符
	TriggerThreshold              types.Int32  `tfsdk:"trigger_threshold"`                 // 阈值
	TriggerPeriod                 types.String `tfsdk:"trigger_period"`                    // 监控周期
	TriggerEvaluationCount        types.Int32  `tfsdk:"trigger_evaluation_count"`          // 连续出现次数
	TriggerCooldown               types.Int32  `tfsdk:"trigger_cooldown"`                  // 冷却时间 (秒)
	TriggerStatus                 types.Int32  `tfsdk:"trigger_status"`                    // 告警规则状态: 0-启用, 1-停用
	TargetMetricName              types.String `tfsdk:"target_metric_name"`                // 监控指标名称
	TargetValue                   types.Int32  `tfsdk:"target_value"`                      // 追踪目标值
	TargetScaleOutEvaluationCount types.Int32  `tfsdk:"target_scale_out_evaluation_count"` // 扩容连续告警次数
	TargetScaleInEvaluationCount  types.Int32  `tfsdk:"target_scale_in_evaluation_count"`  // 缩容连续告警次数
	TargetOperateRange            types.Int32  `tfsdk:"target_operate_range"`              // 缩容波动范围
	TargetDisableScaleIn          types.Bool   `tfsdk:"target_disable_scale_in"`           // 是否禁用缩容
}

type CtyunScalingPoliciesConfig struct {
	RegionID        types.String             `tfsdk:"region_id"`        // 资源池ID
	GroupID         types.Int64              `tfsdk:"group_id"`         // 伸缩组ID
	PageNo          types.Int32              `tfsdk:"page_no"`          // 页码
	PageSize        types.Int32              `tfsdk:"page_size"`        // 分页查询时设置的每页行数，取值范围:[1~100]，默认值为10
	ScalingPolicies []CtyunScalingPolicyList `tfsdk:"scaling_policies"` // 弹性伸缩策略列表
}
