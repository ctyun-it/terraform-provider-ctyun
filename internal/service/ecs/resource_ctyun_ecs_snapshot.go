package ecs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctecs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

func NewCtyunEcsSnapshot() resource.Resource {
	return &ctyunEcsSnapshot{}
}

type ctyunEcsSnapshot struct {
	meta       *common.CtyunMetadata
	ecsService *business.EcsService
}

func (c *ctyunEcsSnapshot) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ecs_snapshot"
}

func (c *ctyunEcsSnapshot) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026730/10335345**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "云主机快照id",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "云主机ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "云主机快照名称，长度为2-63字符，头尾不支持输入空格。支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 63),
				},
			},
			"snapshot_status": schema.StringAttribute{
				Computed:    true,
				Description: "云主机快照状态： pending：创建中, available：可用， restoring：恢复中， error：错误",
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
		},
	}
}

func (c *ctyunEcsSnapshot) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEcsSnapshotConfig
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
	id, err := c.create(ctx, plan)
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

// checkCreate 校验创建动作是否能执行
func (c *ctyunEcsSnapshot) checkCreate(ctx context.Context, plan CtyunEcsSnapshotConfig) error {
	// 1.云主机必须存在
	err := c.ecsService.MustExist(ctx, plan.InstanceId.ValueString(), plan.RegionId.ValueString())
	if err != nil {
		return err
	}
	var status string
	// 2.云主机必须处于运行中（running）或关机（stopped）状态
	status, err = c.ecsService.GetEcsStatus(ctx, plan.InstanceId.ValueString(), plan.RegionId.ValueString())
	if err != nil {
		return err
	}
	allowedStatuses := map[string]bool{
		business.EcsStatusRunning: true,
		business.EcsStatusStopped: true,
	}

	if !allowedStatuses[status] {
		return fmt.Errorf("云主机状态无效(当前:%s)，仅允许在%s或%s状态下创建快照",
			status, business.EcsStatusRunning, business.EcsStatusStopped)
	}

	// 3.云主机快照名称不能和已有云主机快照重名
	err = c.checkName(ctx, plan)
	if err != nil {
		return err
	}
	// 4.云主机所挂全部盘的状态需要为"已挂载"

	// 5. 挂载了本地盘，共享盘，XSSD盘，ISCSI磁盘模式盘的云主机不支持快照功能

	return nil
}

// getAndMerge 查询
func (c *ctyunEcsSnapshot) getAndMerge(ctx context.Context, cfg *CtyunEcsSnapshotConfig) (err error) {
	params := &ctecs2.CtecsQuerySnapshotDetailsV41Request{
		RegionID:   cfg.RegionId.ValueString(),
		SnapshotID: cfg.Id.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsQuerySnapshotDetailsV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if len(resp.ReturnObj.Results) == 0 {
		err = fmt.Errorf("no snapshot details found for snapshot ID: %s", cfg.Id.ValueString())
		return
	}

	//快照名称更新
	cfg.SnapshotName = types.StringValue(resp.ReturnObj.Results[0].SnapshotName)
	cfg.SnapshotStatus = types.StringValue(resp.ReturnObj.Results[0].SnapshotStatus)
	cfg.InstanceId = types.StringValue(resp.ReturnObj.Results[0].InstanceID)
	return
}

func (c *ctyunEcsSnapshot) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsSnapshotConfig
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

func (c *ctyunEcsSnapshot) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunEcsSnapshotConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunEcsSnapshotConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新
	err = c.updateName(ctx, plan, state)
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

func (c *ctyunEcsSnapshot) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsSnapshotConfig
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

func (c *ctyunEcsSnapshot) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ecsService = business.NewEcsService(meta)
}

// create 创建
func (c *ctyunEcsSnapshot) create(ctx context.Context, plan CtyunEcsSnapshotConfig) (Id string, err error) {

	params := &ctecs2.CtecsCreateSnapshotV41Request{
		RegionID:     plan.RegionId.ValueString(),
		SnapshotName: plan.SnapshotName.ValueString(),
		InstanceID:   plan.InstanceId.ValueString(),
	}

	// 创建实例
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsCreateSnapshotV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	Id = resp.ReturnObj.SnapshotID
	return
}

// updateName 修改名称
func (c *ctyunEcsSnapshot) updateName(ctx context.Context, plan, state CtyunEcsSnapshotConfig) (err error) {
	if plan.SnapshotName.Equal(state.SnapshotName) {
		return
	}
	params := &ctecs2.CtecsUpdateSnapshotV41Request{
		RegionID:     state.RegionId.ValueString(),
		SnapshotID:   state.Id.ValueString(),
		SnapshotName: plan.SnapshotName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsUpdateSnapshotV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

// checkName 校验名称是否重复
func (c *ctyunEcsSnapshot) checkName(ctx context.Context, plan CtyunEcsSnapshotConfig) (err error) {
	params := &ctecs2.CtecsQuerySnapshotListV41Request{
		RegionID:     plan.RegionId.ValueString(),
		SnapshotName: plan.SnapshotName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsQuerySnapshotListV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj != nil && len(resp.ReturnObj.Results) > 0 {
		// 如果查询结果存在，说明快照名称已存在，返回错误
		err = fmt.Errorf("snapshot name '%s' already exists", plan.SnapshotName.ValueString())
		return
	}
	return
}

// delete 删除
func (c *ctyunEcsSnapshot) delete(ctx context.Context, plan CtyunEcsSnapshotConfig) (err error) {
	params := &ctecs2.CtecsDeleteSnapshotV41Request{
		RegionID:   plan.RegionId.ValueString(),
		SnapshotID: plan.Id.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsDeleteSnapshotV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

func (c *ctyunEcsSnapshot) StartedLoop(ctx context.Context, state *CtyunEcsSnapshotConfig, loopCount ...int) (err error) {
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
			params := &ctecs2.CtecsQuerySnapshotDetailsV41Request{
				RegionID:   state.RegionId.ValueString(),
				SnapshotID: state.Id.ValueString(),
			}
			// 调用API
			resp, err := c.meta.Apis.SdkCtEcsApis.CtecsQuerySnapshotDetailsV41Api.Do(ctx, c.meta.SdkCredential, params)
			if err != nil {
				return false
			} else if resp.StatusCode == common.ErrorStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			} else if len(resp.ReturnObj.Results) == 0 {
				err = fmt.Errorf("no snapshot details found for snapshot ID: %s", state.Id.ValueString())
				return false
			}

			runningStatus := resp.ReturnObj.Results[0].SnapshotStatus
			if runningStatus == business.EcsSnapshotStatusAvailable {
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

func (c *ctyunEcsSnapshot) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunEcsSnapshotConfig
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

type CtyunEcsSnapshotConfig struct {
	Id             types.String `tfsdk:"id"`
	InstanceId     types.String `tfsdk:"instance_id"`
	SnapshotName   types.String `tfsdk:"name"`
	SnapshotStatus types.String `tfsdk:"snapshot_status"`
	ProjectId      types.String `tfsdk:"project_id"`
	RegionId       types.String `tfsdk:"region_id"`
}
