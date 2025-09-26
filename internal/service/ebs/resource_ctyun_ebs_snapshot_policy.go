package ebs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctebs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"time"
)

func NewCtyunEbsSnapshotPolicy() resource.Resource {
	return &ctyunEbsSnapshotPolicy{}
}

type ctyunEbsSnapshotPolicy struct {
	meta       *common.CtyunMetadata
	ebsService *business.EbsService
}

func (c *ctyunEbsSnapshotPolicy) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs_snapshot_policy"
}

func (c *ctyunEbsSnapshotPolicy) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027696/10118840`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "云硬盘快照策略id",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "云硬盘快照策略名称，长度为2-63字符，头尾不支持输入空格。支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 63),
				},
			},
			"repeat_weekdays": schema.StringAttribute{
				Required:    true,
				Description: "创建快照的重复日期，0-6分别代表周日-周六，多个日期用英文逗号隔开。支持更新",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[0-6](,[0-6])*$`),
						"必须是由逗号分隔的0-6之间的整数组成",
					),
				},
			},
			"repeat_times": schema.StringAttribute{
				Required:    true,
				Description: "创建快照的重复时间，0-23分别代表零点-23点，多个时间用英文逗号隔开。支持更新",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^([0-9]|1[0-9]|2[0-3])(,([0-9]|1[0-9]|2[0-3]))*$`),
						"必须是由逗号分隔的0-23之间的整数组成",
					),
				},
			},
			"retention_time": schema.Int64Attribute{
				Required:    true,
				Description: "创建快照的保留时间，输入范围为[-1，1-65535]，-1代表永久保留。单位为天。支持更新",
				Validators: []validator.Int64{
					int64validator.Between(-1, 65535),
				},
			},
			"is_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否启用策略，取值范围：true：启用，false：不启用，默认为true。支持更新",
				Default:     booldefault.StaticBool(true),
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},

			"snapshot_policy_status": schema.StringAttribute{
				Computed:    true,
				Description: "自动快照策略状态，取值范围：activated:启用，nonactivated：停用",
			},
			"bound_disk_num": schema.Int64Attribute{
				Computed:    true,
				Description: "关联云硬盘的数量",
			},
			"snapshot_policy_create_time": schema.StringAttribute{
				Computed:    true,
				Description: "策略创建时间",
			},
		},
	}
}

func (c *ctyunEbsSnapshotPolicy) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEbsSnapshotPolicyConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 校验创建动作的前置条件
	err = c.checkCreate(ctx, plan)
	if err != nil {
		return
	}

	// 实际创建
	id, err := c.create(ctx, &plan)
	if err != nil {
		return
	}
	plan.Id = types.StringValue(id)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)

	//轮询快照状态为可用状态
	err = c.StartedLoop(ctx, &plan)
	if err != nil {
		return
	}
	// 查询信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunEbsSnapshotPolicy) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunEbsSnapshotPolicyConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunEbsSnapshotPolicyConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新
	if plan.IsEnabled.ValueBool() != state.IsEnabled.ValueBool() {
		err = c.updateStatus(ctx, plan, state)
		if err != nil {
			return
		}
	} else {
		if plan.Name.ValueString() != state.Name.ValueString() {
			//修改自动快照策略名称时不允许与已有的策略名称重复。
			err := c.checkName(ctx, plan)
			if err != nil {
				return
			}
		}
		err = c.update(ctx, plan, state)
	}

	if err != nil {
		return
	}
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

// checkCreate 校验创建动作是否能执行
func (c *ctyunEbsSnapshotPolicy) checkCreate(ctx context.Context, plan CtyunEbsSnapshotPolicyConfig) error {

	// 单个用户的单个资源池下自动快照策略名称不能重复。
	err := c.checkName(ctx, plan)
	if err != nil {
		return err
	}

	return nil
}

// update 修改云硬盘自动快照策略
func (c *ctyunEbsSnapshotPolicy) update(ctx context.Context, plan, state CtyunEbsSnapshotPolicyConfig) (err error) {
	params := &ctebs2.EbsModifyPolicyEbsSnapRequest{
		RegionID:           state.RegionId.ValueString(),
		SnapshotPolicyID:   state.Id.ValueString(),
		SnapshotPolicyName: plan.Name.ValueStringPointer(),
		RepeatWeekdays:     plan.RepeatWeekdays.ValueStringPointer(),
		RetentionTime:      int32(plan.RetentionTime.ValueInt64()),
		RepeatTimes:        plan.RepeatTimes.ValueStringPointer(),
	}
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsModifyPolicyEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return
}

const (
	StatusActivated    = "activated"
	StatusNonActivated = "nonactivated"
)

// updateStatus 启用或关闭云硬盘自动快照策略
func (c *ctyunEbsSnapshotPolicy) updateStatus(ctx context.Context, plan, state CtyunEbsSnapshotPolicyConfig) (err error) {
	targetStatus := StatusNonActivated
	if plan.IsEnabled.ValueBool() {
		targetStatus = StatusActivated
	}
	params := &ctebs2.EbsModifyPolicyStatusEbsSnapRequest{
		RegionID:         state.RegionId.ValueString(),
		SnapshotPolicyID: state.Id.ValueString(),
		TargetStatus:     targetStatus,
	}
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsModifyPolicyStatusEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if *resp.ReturnObj.SnapshotPolicyJobResult != "任务执行成功" {
		err = fmt.Errorf("启用或关闭云硬盘自动快照策略 任务执行失败. Message: %s SnapshotPolicyID: %s", *resp.Message, plan.Id.ValueString())
		return
	}

	return
}

// getAndMerge 查询
func (c *ctyunEbsSnapshotPolicy) getAndMerge(ctx context.Context, cfg *CtyunEbsSnapshotPolicyConfig) (err error) {
	// 获取实例详情
	params := &ctebs2.EbsQueryPolicyEbsSnapRequest{
		RegionID:         cfg.RegionId.ValueString(),
		SnapshotPolicyID: cfg.Id.ValueStringPointer(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsQueryPolicyEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.ReturnObj.SnapshotPolicyTotalCount == 0 {
		err = fmt.Errorf("no snapshot details found for snapshot ID: %s", cfg.Id.ValueString())
		return
	}
	// 处理返回的数据
	policy := resp.ReturnObj.SnapshotPolicyList[0]

	// 设置字段值
	cfg.SnapshotPolicyStatus = types.StringValue(*policy.SnapshotPolicyStatus)
	cfg.RetentionTime = types.Int64Value(int64(policy.RetentionTime))
	cfg.Name = types.StringValue(*policy.SnapshotPolicyName)
	cfg.BoundDiskNum = types.Int64Value(int64(policy.BoundDiskNum))
	cfg.SnapshotPolicyCreateTime = types.StringValue(*policy.SnapshotPolicyCreateTime)
	cfg.Id = types.StringValue(*policy.SnapshotPolicyID)
	cfg.RepeatWeekdays = types.StringValue(*policy.RepeatWeekdays)
	cfg.RepeatTimes = types.StringValue(*policy.RepeatTimes)
	// 根据 SnapshotPolicyStatus 设置 is_enabled
	isEnabled := false
	if *policy.SnapshotPolicyStatus == StatusActivated {
		isEnabled = true
	}
	cfg.IsEnabled = types.BoolValue(isEnabled)

	return
}

func (c *ctyunEbsSnapshotPolicy) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbsSnapshotPolicyConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEbsSnapshotPolicy) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbsSnapshotPolicyConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunEbsSnapshotPolicy) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ebsService = business.NewEbsService(meta)
}

// create 创建
func (c *ctyunEbsSnapshotPolicy) create(ctx context.Context, plan *CtyunEbsSnapshotPolicyConfig) (id string, err error) {

	params := &ctebs2.EbsCreatePolicyEbsSnapRequest{
		RegionID:           plan.RegionId.ValueString(),
		ProjectID:          plan.ProjectId.ValueString(),
		SnapshotPolicyName: plan.Name.ValueString(),
		RepeatWeekdays:     plan.RepeatWeekdays.ValueString(),
		RetentionTime:      int32(plan.RetentionTime.ValueInt64()),
		RepeatTimes:        plan.RepeatTimes.ValueString(),
		IsEnabled:          plan.IsEnabled.ValueBoolPointer(),
	}

	// 创建实例
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsCreatePolicyEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	if resp.ReturnObj != nil && resp.ReturnObj.SnapshotPolicyID != nil {
		id = *resp.ReturnObj.SnapshotPolicyID
	} else {
		err = fmt.Errorf("failed to get snapshot policy ID from response")
	}

	return
}

// checkName 重名判断
func (c *ctyunEbsSnapshotPolicy) checkName(ctx context.Context, plan CtyunEbsSnapshotPolicyConfig) (err error) {
	params := &ctebs2.EbsQueryPolicyEbsSnapRequest{
		RegionID:           plan.RegionId.ValueString(),
		SnapshotPolicyName: plan.Name.ValueStringPointer(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsQueryPolicyEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.ReturnObj.SnapshotPolicyTotalCount > 0 {
		// 如果查询结果存在，说明名称已存在，返回错误
		err = fmt.Errorf("name '%s' already exists", plan.Name.ValueString())
		return
	}
	return
}

// delete 删除
func (c *ctyunEbsSnapshotPolicy) delete(ctx context.Context, plan CtyunEbsSnapshotPolicyConfig) (err error) {

	params := &ctebs2.EbsDeletePolicyEbsSnapRequest{
		RegionID:         plan.RegionId.ValueString(),
		SnapshotPolicyID: plan.Id.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsDeletePolicyEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}

func (c *ctyunEbsSnapshotPolicy) StartedLoop(ctx context.Context, state *CtyunEbsSnapshotPolicyConfig, loopCount ...int) (err error) {
	count := 30
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			// 获取实例详情
			params := &ctebs2.EbsQueryPolicyEbsSnapRequest{
				RegionID:         state.RegionId.ValueString(),
				SnapshotPolicyID: state.Id.ValueStringPointer(),
			}
			// 调用API
			resp, err := c.meta.Apis.SdkCtEbsApis.EbsQueryPolicyEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
			if err != nil {
				return false
			} else if resp.StatusCode == common.ErrorStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			} else if resp.ReturnObj.SnapshotPolicyTotalCount == 0 {
				err = fmt.Errorf("no snapshot details found for snapshot ID: %s", state.Id.ValueString())
				return false
			}

			runningStatus := resp.ReturnObj.SnapshotPolicyList[0].SnapshotPolicyStatus
			if *runningStatus == "activated" {
				return false
			}
			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未到达启动状态！")
	}
	return
}

func (c *ctyunEbsSnapshotPolicy) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunEbsSnapshotPolicyConfig
	var id, regionID string
	err = terraform_extend.Split(request.ID, &id, &regionID)
	if err != nil {
		return
	}
	cfg.RegionId = types.StringValue(regionID)
	cfg.Id = types.StringValue(id)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

type CtyunEbsSnapshotPolicyConfig struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	RepeatWeekdays types.String `tfsdk:"repeat_weekdays"`
	RepeatTimes    types.String `tfsdk:"repeat_times"`
	RetentionTime  types.Int64  `tfsdk:"retention_time"`
	IsEnabled      types.Bool   `tfsdk:"is_enabled"`
	ProjectId      types.String `tfsdk:"project_id"`
	RegionId       types.String `tfsdk:"region_id"`

	SnapshotPolicyStatus     types.String `tfsdk:"snapshot_policy_status"`
	BoundDiskNum             types.Int64  `tfsdk:"bound_disk_num"`
	SnapshotPolicyCreateTime types.String `tfsdk:"snapshot_policy_create_time"`
}
