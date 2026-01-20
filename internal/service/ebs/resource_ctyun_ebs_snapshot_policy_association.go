package ebs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctebs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctebs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &ctyunEbsSnapshotPolicyAssociation{}
	_ resource.ResourceWithConfigure   = &ctyunEbsSnapshotPolicyAssociation{}
	_ resource.ResourceWithImportState = &ctyunEbsSnapshotPolicyAssociation{}
)

/*
将快照策略和云硬盘绑定
*/

func NewCtyunEbsSnapshotPolicyAssociation() resource.Resource {
	return &ctyunEbsSnapshotPolicyAssociation{}
}

type ctyunEbsSnapshotPolicyAssociation struct {
	meta *common.CtyunMetadata
	name string
}

func (c *ctyunEbsSnapshotPolicyAssociation) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs_snapshot_policy_association"
	c.name = response.TypeName
}

type CtyunEbsSnapshotPolicyAssociationConfig struct {
	ID               types.String `tfsdk:"id"`
	SnapshotPolicyID types.String `tfsdk:"snapshot_policy_id"`
	RegionID         types.String `tfsdk:"region_id"`
	DiskID           types.String `tfsdk:"disk_id"`
}

func (c *ctyunEbsSnapshotPolicyAssociation) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027696/10118856`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"snapshot_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "云硬盘自动快照策略id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
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
		},
	}
}

func (c *ctyunEbsSnapshotPolicyAssociation) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEbsSnapshotPolicyAssociationConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeBindDisks(ctx, plan)
	if err != nil {
		return
	}
	err = c.createSnapshotOrder(ctx, plan)
	if err != nil {
		return
	}
	// 实际创建
	err = c.create(ctx, plan)
	if err != nil {
		return
	}
	err = c.checkAfterBindDisks(ctx, plan)
	if err != nil {
		return
	}
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunEbsSnapshotPolicyAssociation) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (c *ctyunEbsSnapshotPolicyAssociation) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbsSnapshotPolicyAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "未关联") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEbsSnapshotPolicyAssociation) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbsSnapshotPolicyAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkBeforeDissociate(ctx, state)
	if err != nil {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDissociation(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunEbsSnapshotPolicyAssociation) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// create 创建
func (c *ctyunEbsSnapshotPolicyAssociation) create(ctx context.Context, plan CtyunEbsSnapshotPolicyAssociationConfig) (err error) {
	params := &ctebs2.EbsApplyPolicyEbsSnapRequest{
		RegionID:         plan.RegionID.ValueString(),
		SnapshotPolicyID: plan.SnapshotPolicyID.ValueString(),
		TargetDiskIDs:    plan.DiskID.ValueString(),
	}
	// 创建实例
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsApplyPolicyEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.StatusCode == common.NormalStatusCode {
		return
	}

	return
}

func (c *ctyunEbsSnapshotPolicyAssociation) createSnapshotOrder(ctx context.Context, plan CtyunEbsSnapshotPolicyAssociationConfig) (err error) {
	clientToken := uuid.NewString()
	params := &ctebs2.EbsCreateOrderEbsSnapRequest{
		RegionID:    plan.RegionID.ValueString(),
		ClientToken: &clientToken,
		DiskID:      plan.DiskID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsCreateOrderEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		if *resp.Description == "云硬盘已经开通过快照服务" {
			return
		}
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

func (c *ctyunEbsSnapshotPolicyAssociation) checkBeforeBindDisks(ctx context.Context, cfg CtyunEbsSnapshotPolicyAssociationConfig) (err error) {
	// 获取实例详情
	params := &ctebs2.EbsQueryPolicyEbsSnapRequest{
		RegionID:         cfg.RegionID.ValueString(),
		SnapshotPolicyID: cfg.SnapshotPolicyID.ValueStringPointer(),
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
	} else if resp.ReturnObj.SnapshotPolicyTotalCount != 1 {
		return fmt.Errorf("自动快照策略必须存在")
	}
	diskId := cfg.DiskID.ValueString()
	//查询云硬盘
	_, err2 := c.meta.Apis.CtEbsApis.EbsShowApi.Do(ctx, c.meta.Credential, &ctebs.EbsShowRequest{
		RegionId: cfg.RegionID.ValueString(),
		DiskId:   diskId,
	})
	if err2 != nil {
		// 实例已经被退订的情况
		if err2.ErrorCode() == common.EcsInstanceNotFound {
			return nil
		}
		return err2
	}
	return
}

// checkAfterBindDisks 绑定后检查
func (c *ctyunEbsSnapshotPolicyAssociation) checkAfterBindDisks(ctx context.Context, plan CtyunEbsSnapshotPolicyAssociationConfig) (err error) {
	var executeSuccessFlag bool
	var snapshotPolicyID string
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			snapshotPolicyID, err = c.getBindingDisks(ctx, plan)
			if err != nil {
				return false
			}
			if snapshotPolicyID == plan.SnapshotPolicyID.ValueString() {
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("云硬盘自动快照策略 %s 和云硬盘 %s 未关联", plan.SnapshotPolicyID.String(), plan.DiskID.ValueString())
	}
	return nil
}

// checkBeforeDissociate 解绑前检查
func (c *ctyunEbsSnapshotPolicyAssociation) checkBeforeDissociate(ctx context.Context, plan CtyunEbsSnapshotPolicyAssociationConfig) (err error) {

	snapshotPolicyID, err := c.getBindingDisks(ctx, plan)
	if err != nil {
		return
	}
	if snapshotPolicyID != plan.SnapshotPolicyID.ValueString() {
		err = fmt.Errorf("云硬盘自动快照策略 %s 和云硬盘 %s 未关联", plan.SnapshotPolicyID.String(), plan.DiskID.ValueString())
		return
	}
	return
}

// checkAfterDissociation 解绑后检查
func (c *ctyunEbsSnapshotPolicyAssociation) checkAfterDissociation(ctx context.Context, plan CtyunEbsSnapshotPolicyAssociationConfig) (err error) {
	var executeSuccessFlag bool
	var snapshotPolicyID string
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			snapshotPolicyID, err = c.getBindingDisks(ctx, plan)
			if err != nil {
				return false
			}
			if snapshotPolicyID != plan.SnapshotPolicyID.ValueString() {
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		return fmt.Errorf("云硬盘自动快照策略 %s 和云硬盘 %s 解绑失败", plan.SnapshotPolicyID.ValueString(), plan.DiskID.ValueString())
	}
	return nil
}

// dissociate 解绑
func (c *ctyunEbsSnapshotPolicyAssociation) delete(ctx context.Context, plan CtyunEbsSnapshotPolicyAssociationConfig) (err error) {
	params := &ctebs2.EbsCancelPolicyEbsSnapRequest{
		RegionID:      plan.RegionID.ValueString(),
		TargetDiskIDs: plan.DiskID.ValueString(),
	}

	// 创建实例
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsCancelPolicyEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.StatusCode == common.NormalStatusCode {
		return
	}

	return
}

func (c *ctyunEbsSnapshotPolicyAssociation) getBindingDisks(ctx context.Context, plan CtyunEbsSnapshotPolicyAssociationConfig) (snapshotPolicyID string, err error) {
	diskId := plan.DiskID.ValueString()
	// 组装请求体
	params := &ctebs2.EbsQueryEbsByIDRequest{
		RegionID: plan.RegionID.ValueString(),
		DiskID:   diskId,
	}

	// 调用API
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsQueryEbsByIDApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return "", err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return "", err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return "", err
	}
	return resp.ReturnObj.SnapshotPolicyID, nil
}

// getAndMerge 查询绑定关系
func (c *ctyunEbsSnapshotPolicyAssociation) getAndMerge(ctx context.Context, plan *CtyunEbsSnapshotPolicyAssociationConfig) (err error) {
	policyId, diskID, regionID := plan.SnapshotPolicyID.ValueString(), plan.DiskID.ValueString(), plan.RegionID.ValueString()
	snapshotPolicyID, err := c.getBindingDisks(ctx, *plan)
	if err != nil {
		return
	}
	if snapshotPolicyID != plan.SnapshotPolicyID.ValueString() {
		err = fmt.Errorf("云硬盘自动快照策略 %s 和云硬盘 %s 未关联  regionID： %s", policyId, diskID, regionID)
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", policyId, diskID, regionID))
	return
}

func (c *ctyunEbsSnapshotPolicyAssociation) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := c.name + "导入失败：" + err.Error()
			detail := fmt.Sprintf("导入命令：terraform import %s.[导入配置名称] [snapshot_policy_id],[disk_id],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunEbsSnapshotPolicyAssociationConfig

	var policyID, diskID, regionId string
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &policyID, &diskID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &policyID, &diskID, &regionId)
		if err != nil {
			return
		}
	}
	if policyID == "" {
		err = fmt.Errorf("snapshot_policy_id不能为空")
		return
	}
	if diskID == "" {
		err = fmt.Errorf("disk_id不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}
	cfg.SnapshotPolicyID = types.StringValue(policyID)
	cfg.RegionID = types.StringValue(regionId)
	cfg.DiskID = types.StringValue(diskID)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}
