package scaling

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/scaling"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunScalingActivities{}
	_ datasource.DataSourceWithConfigure = &CtyunScalingActivities{}
)

type CtyunScalingActivities struct {
	meta *common.CtyunMetadata
}

func NewCtyunScalingActivities() datasource.DataSource {
	return &CtyunScalingActivities{}
}

func (c *CtyunScalingActivities) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_scaling_activities"
}

func (c *CtyunScalingActivities) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunScalingActivities) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027725/10216432**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"group_id": schema.Int64Attribute{
				Required:    true,
				Description: "伸缩组ID",
			},
			"active_ids": schema.SetAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				Description: "待查询伸缩活动ID列表",
			},
			"start_time": schema.Int64Attribute{
				Optional:    true,
				Description: "开始时间 (Unix时间戳，秒级)",
			},
			"end_time": schema.Int64Attribute{
				Optional:    true,
				Description: "结束时间 (Unix时间戳，秒级)",
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
			"scaling_activities": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_fail_reason": schema.StringAttribute{
							Computed:    true,
							Description: "策略失败原因",
						},
						"after_count": schema.Int64Attribute{
							Computed:    true,
							Description: "活动后计数",
						},
						"end_time": schema.StringAttribute{
							Computed:    true,
							Description: "结束时间 (格式: 2006-01-02 15:04:05)",
						},
						"before_count": schema.Int64Attribute{
							Computed:    true,
							Description: "活动前计数",
						},
						"rule_id": schema.StringAttribute{
							Computed:    true,
							Description: "伸缩策略ID",
						},
						"start_time": schema.StringAttribute{
							Computed:    true,
							Description: "开始时间 (格式: 2006-01-02 15:04:05)",
						},
						"fail_reason": schema.StringAttribute{
							Computed:    true,
							Description: "失败原因",
						},
						"instance_list": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"instance_id": schema.StringAttribute{
										Computed:    true,
										Description: "云主机ID",
									},
									"instance_name": schema.StringAttribute{
										Computed:    true,
										Description: "云主机名称",
									},
								},
							},
							Description: "涉及实例列表",
						},
						"execution_mode": schema.StringAttribute{
							Computed:    true,
							Description: "执行方式",
						},
						"group_id": schema.Int64Attribute{
							Computed:    true,
							Description: "伸缩组ID",
						},
						"rule_expect_delta": schema.Int64Attribute{
							Computed:    true,
							Description: "策略预期可变化数量",
						},
						"execution_result": schema.Int32Attribute{
							Computed:    true,
							Description: "执行结果: 0-执行中, 1-成功, 2-失败",
						},
						"execution_date": schema.StringAttribute{
							Computed:    true,
							Description: "执行时间 (格式: 2006-01-02 15:04:05)",
						},
						"rule_execution_result": schema.StringAttribute{
							Computed:    true,
							Description: "策略执行结果",
						},
						"active_id": schema.Int64Attribute{
							Computed:    true,
							Description: "伸缩活动ID",
						},
						"rule_real_delta": schema.Int32Attribute{
							Computed:    true,
							Description: "策略实际可变化数量",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "说明",
						},
					},
				},
				Description: "伸缩活动列表",
			},
		},
	}
}

func (c *CtyunScalingActivities) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunScalingActivitiesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	if regionId == "" {
		err = errors.New("region id 为空")
		return
	}

	params := &scaling.ScalingQueryActivitiesListRequest{
		RegionID: regionId,
		GroupID:  config.GroupID.ValueInt64(),
		PageNo:   1,
		PageSize: 10,
	}
	if !config.PageNo.IsNull() {
		params.PageNo = config.PageNo.ValueInt32()
	}
	if !config.PageSize.IsNull() {
		params.PageSize = config.PageSize.ValueInt32()
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingQueryActivitiesListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("获取group_id 为：%d 下的伸缩活动列表失败，返回为nil。", config.GroupID.ValueInt64())
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	var scalingActivities []scalingActivityModel
	for _, item := range resp.ReturnObj.ActiveList {
		var activity scalingActivityModel
		activity.RuleFailReason = types.StringValue(item.FailReason)
		activity.AfterCount = types.Int32Value(item.AfterCount)
		activity.EndTime = types.StringValue(item.EndTime)
		activity.BeforeCount = types.Int32Value(item.BeforeCount)
		activity.RuleID = types.StringValue(item.RuleID)
		activity.StartTime = types.StringValue(item.StartTime)
		activity.FailReason = types.StringValue(item.FailReason)
		activity.ExecutionMode = types.StringValue(business.ExecutionModeToString[item.ExecutionMode])
		activity.GroupID = types.Int64Value(item.GroupID)
		activity.RuleExpectDelta = types.Int32Value(item.RuleExpectDelta)
		activity.ExecutionResult = types.Int32Value(item.ExecutionResult)
		activity.ExecutionDate = types.StringValue(item.ExecutionDate)
		activity.RuleExecutionResult = types.StringValue(item.RuleExecutionResult)
		activity.ActiveID = types.Int64Value(item.ActiveID)
		activity.RuleRealDelta = types.Int32Value(item.RuleRealDelta)
		activity.Description = types.StringValue(item.Description)

		var activityInstanceList []CtyunScalingActivityInstance
		if item.InstanceList != nil {
			for _, instanceItem := range item.InstanceList {
				var instance CtyunScalingActivityInstance
				instance.InstanceID = types.StringValue(instanceItem.InstanceID)
				instance.InstanceName = types.StringValue(instanceItem.InstanceName)
				activityInstanceList = append(activityInstanceList, instance)
			}
		}
		instanceList, diagnostics := types.ListValueFrom(ctx, utils.StructToTFObjectTypes(CtyunScalingActivityInstance{}), activityInstanceList)
		if diagnostics != nil && diagnostics.HasError() {
			err = errors.New(diagnostics[0].Detail())
			return
		}
		activity.InstanceList = instanceList

		scalingActivities = append(scalingActivities, activity)

	}
	config.ScalingActivities = scalingActivities
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type scalingActivityModel struct {
	RuleFailReason      types.String `tfsdk:"rule_fail_reason"`      // 策略失败原因
	AfterCount          types.Int32  `tfsdk:"after_count"`           // 活动后计数
	EndTime             types.String `tfsdk:"end_time"`              // 结束时间
	BeforeCount         types.Int32  `tfsdk:"before_count"`          // 活动前计数
	RuleID              types.String `tfsdk:"rule_id"`               // 伸缩策略ID
	StartTime           types.String `tfsdk:"start_time"`            // 开始时间
	FailReason          types.String `tfsdk:"fail_reason"`           // 失败原因
	InstanceList        types.List   `tfsdk:"instance_list"`         // 虚机列表
	ExecutionMode       types.String `tfsdk:"execution_mode"`        // 执行方式
	GroupID             types.Int64  `tfsdk:"group_id"`              // 伸缩组ID
	RuleExpectDelta     types.Int32  `tfsdk:"rule_expect_delta"`     // 策略预期可变化数量
	ExecutionResult     types.Int32  `tfsdk:"execution_result"`      // 执行结果
	ExecutionDate       types.String `tfsdk:"execution_date"`        // 执行时间
	RuleExecutionResult types.String `tfsdk:"rule_execution_result"` // 策略执行结果
	ActiveID            types.Int64  `tfsdk:"active_id"`             // 伸缩活动ID
	RuleRealDelta       types.Int32  `tfsdk:"rule_real_delta"`       // 策略实际可变化数量
	Description         types.String `tfsdk:"description"`           // 说明
}

// CtyunScalingActivityInstance 伸缩活动涉及的实例
type CtyunScalingActivityInstance struct {
	InstanceID   types.String `tfsdk:"instance_id"`   // 云主机ID
	InstanceName types.String `tfsdk:"instance_name"` // 云主机名称
}
type CtyunScalingActivitiesConfig struct {
	RegionID          types.String           `tfsdk:"region_id"`          // 资源池ID
	GroupID           types.Int64            `tfsdk:"group_id"`           // 伸缩组ID
	ActiveIDs         types.Set              `tfsdk:"active_ids"`         // 待查询伸缩活动ID列表
	StartTime         types.Int64            `tfsdk:"start_time"`         // 开始时间 (Unix时间戳，秒级)
	EndTime           types.Int64            `tfsdk:"end_time"`           // 结束时间 (Unix时间戳，秒级)
	PageNo            types.Int32            `tfsdk:"page_no"`            // 页码
	PageSize          types.Int32            `tfsdk:"page_size"`          // 每页行数
	ScalingActivities []scalingActivityModel `tfsdk:"scaling_activities"` // 伸缩活动列表
}
