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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

func NewCtyunEbsSnapshot() resource.Resource {
	return &ctyunEbsSnapshot{}
}

type ctyunEbsSnapshot struct {
	meta       *common.CtyunMetadata
	ebsService *business.EbsService
}

func (c *ctyunEbsSnapshot) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (c *ctyunEbsSnapshot) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs_snapshot"
}

func (c *ctyunEbsSnapshot) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027696/10043223`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "云硬盘快照id",
			},
			"disk_id": schema.StringAttribute{
				Required:    true,
				Description: "云硬盘ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "云硬盘快照名称，长度为2-63字符，头尾不支持输入空格",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 63),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"snapshot_status": schema.StringAttribute{
				Computed:    true,
				Description: "云硬盘快照状态： pending：创建中, available：可用， restoring：恢复中， error：错误",
			},
			"retention_policy": schema.StringAttribute{
				Required:    true,
				Description: "快照保留策略，取值范围：custom：自定义保留天数，forever：永久保留",
				Validators: []validator.String{
					stringvalidator.OneOf("custom", "forever"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"retention_time": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "自定义快照保留天数。取值范围：1-65535。当快照保留策略为custom时该参数为必填，当快照保留策略为forever时，自动设置为65535",
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
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

func (c *ctyunEbsSnapshot) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEbsSnapshotConfig
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
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
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
func (c *ctyunEbsSnapshot) checkCreate(ctx context.Context, plan CtyunEbsSnapshotConfig) error {
	// 1.云硬盘必须存在
	err := c.ebsService.MustExist(ctx, plan.DiskId.ValueString(), plan.RegionId.ValueString())
	if err != nil {
		return err
	}
	// 2.云硬盘必须处于“未挂载”或“已挂载”状态 ——资料描述错误 应该是available状态  暂不校验
	resp, err := c.ebsService.GetEbsInfo(ctx, plan.DiskId.ValueString(), plan.RegionId.ValueString())
	if err != nil {
		return err
	}

	// 3.同一个租户在同一个资源池下创建的快照名称不能重复。
	err = c.checkName(ctx, plan)
	if err != nil {
		return err
	}
	// 4.磁盘模式为FCSAN或ISCSI的云硬盘不支持创建快照。
	if business.EbsDiskModeVbd != strings.ToLower(resp.DiskMode) {
		return fmt.Errorf("磁盘模式为FCSAN或ISCSI的云硬盘不支持创建快照")
	}

	return nil
}

// getAndMerge 查询
func (c *ctyunEbsSnapshot) getAndMerge(ctx context.Context, cfg *CtyunEbsSnapshotConfig) (err error) {
	// 获取实例详情
	snapshotId := cfg.Id.ValueString()
	params := &ctebs2.EbsListEbsSnapRequest{
		RegionID:   cfg.RegionId.ValueString(),
		SnapshotID: &snapshotId,
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsListEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)

	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.ReturnObj.SnapshotTotal == 0 {
		err = fmt.Errorf("no snapshot details found for snapshot ID: %s", cfg.Id.ValueString())
		return
	}
	status := resp.ReturnObj.SnapshotList[0].SnapshotStatus
	cfg.SnapshotStatus = types.StringValue(*status)
	cfg.DiskId = types.StringValue(*resp.ReturnObj.SnapshotList[0].DiskID)
	cfg.RetentionTime = types.Int64Value(resp.ReturnObj.SnapshotList[0].RetentionTime)
	cfg.SnapshotName = types.StringValue(*resp.ReturnObj.SnapshotList[0].SnapshotName)
	cfg.RetentionPolicy = types.StringValue(*resp.ReturnObj.SnapshotList[0].RetentionPolicy)
	return
}

func (c *ctyunEbsSnapshot) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbsSnapshotConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "no snapshot details found") {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEbsSnapshot) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbsSnapshotConfig
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

func (c *ctyunEbsSnapshot) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ebsService = business.NewEbsService(meta)
}

// create 创建
func (c *ctyunEbsSnapshot) create(ctx context.Context, plan *CtyunEbsSnapshotConfig) (err error) {

	params := &ctebs2.EbsCreateEbsSnapRequest{
		RegionID:        plan.RegionId.ValueString(),
		SnapshotName:    plan.SnapshotName.ValueString(),
		DiskID:          plan.DiskId.ValueString(),
		RetentionPolicy: plan.RetentionPolicy.ValueString(),
		RetentionTime:   plan.RetentionTime.ValueInt64(),
	}

	// 创建实例
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsCreateEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	}
	var masterOrderId, snapshotJobID string
	if resp != nil && resp.ErrorCode == common.EbsOrderInProgress {
		if resp.ReturnObj.MasterOrderID != "" {
			masterOrderId = resp.ReturnObj.MasterOrderID
		} else if resp.ReturnObj.Resources != nil && len(resp.ReturnObj.Resources) > 0 && resp.ReturnObj.Resources[0].OrderID != "" {
			masterOrderId = resp.ReturnObj.Resources[0].OrderID
		}

		// 轮询订单状态
		helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
		loop, err2 := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId)
		if err2 != nil {
			return err2
		}
		// 最后设置id
		id := loop.Uuid[0]
		plan.Id = types.StringValue(id)
	}
	if resp != nil && resp.ReturnObj.SnapshotJobID != "" {
		snapshotJobID = resp.ReturnObj.SnapshotJobID
		err := c.queryJob(ctx, plan, snapshotJobID)
		if err != nil {
			return err
		}
	}

	return
}

func (c *ctyunEbsSnapshot) queryJob(ctx context.Context, plan *CtyunEbsSnapshotConfig, jobID string) (err error) {
	helper := business.NewGeneralJobHelper(c.meta.Apis.CtEcsApis.JobShowApi)
	reps, err := helper.JobLoop(ctx, c.meta.Credential, plan.RegionId.ValueString(), jobID)
	if err != nil {
		return err
	} else {
		plan.Id = types.StringValue(reps.ID)
	}
	return
}

// checkName 重名判断
func (c *ctyunEbsSnapshot) checkName(ctx context.Context, plan CtyunEbsSnapshotConfig) (err error) {
	snapshotName := plan.SnapshotName.ValueString()
	params := &ctebs2.EbsListEbsSnapRequest{
		RegionID:     plan.RegionId.ValueString(),
		SnapshotName: &snapshotName,
	}
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsListEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj != nil && resp.ReturnObj.SnapshotTotal > 0 {
		// 如果查询结果存在，说明快照名称已存在，返回错误
		err = fmt.Errorf("snapshot name '%s' already exists", plan.SnapshotName.ValueString())
		return
	}
	return
}

// delete 删除
func (c *ctyunEbsSnapshot) delete(ctx context.Context, plan CtyunEbsSnapshotConfig) (err error) {
	snapshotID := plan.Id.ValueString()

	params := &ctebs2.EbsDeleteEbsSnapRequest{
		RegionID:    plan.RegionId.ValueString(),
		DiskID:      plan.DiskId.ValueString(),
		SnapshotIDs: []string{snapshotID},
		RefundOrder: true,
	}
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsDeleteEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

func (c *ctyunEbsSnapshot) StartedLoop(ctx context.Context, state *CtyunEbsSnapshotConfig, loopCount ...int) (err error) {
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
			snapshotId := state.Id.ValueString()
			params := &ctebs2.EbsListEbsSnapRequest{
				RegionID:   state.RegionId.ValueString(),
				SnapshotID: &snapshotId,
			}
			// 调用API
			resp, err := c.meta.Apis.SdkCtEbsApis.EbsListEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
			if err != nil {
				return false
			} else if resp.StatusCode == common.ErrorStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			} else if resp.ReturnObj.SnapshotTotal == 0 {
				err = fmt.Errorf("no snapshot details found for snapshot ID: %s", state.Id.ValueString())
				return false
			}

			runningStatus := resp.ReturnObj.SnapshotList[0].SnapshotStatus
			if *runningStatus == business.EbsSnapshotStatusAvailable {
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

func (c *ctyunEbsSnapshot) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunEbsSnapshotConfig
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

type CtyunEbsSnapshotConfig struct {
	Id              types.String `tfsdk:"id"`
	DiskId          types.String `tfsdk:"disk_id"`
	SnapshotName    types.String `tfsdk:"name"`
	SnapshotStatus  types.String `tfsdk:"snapshot_status"`
	RetentionPolicy types.String `tfsdk:"retention_policy"`
	RetentionTime   types.Int64  `tfsdk:"retention_time"`
	ProjectId       types.String `tfsdk:"project_id"`
	RegionId        types.String `tfsdk:"region_id"`
}
