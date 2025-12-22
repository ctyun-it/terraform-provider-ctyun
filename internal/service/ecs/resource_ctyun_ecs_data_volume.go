package ecs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

var (
	_ resource.Resource                = &CtyunEcsDataVolume{}
	_ resource.ResourceWithConfigure   = &CtyunEcsDataVolume{}
	_ resource.ResourceWithImportState = &CtyunEcsDataVolume{}
)

func NewCtyunEcsDataVolume() resource.Resource {
	return &CtyunEcsDataVolume{}
}

type CtyunEcsDataVolume struct {
	meta       *common.CtyunMetadata
	ecsService *business.EcsService
	ebsService *business.EbsService
	jobLooper  *business.GeneralJobHelper
}

func (c *CtyunEcsDataVolume) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ecs_data_volume"
}

func (c *CtyunEcsDataVolume) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027696/10169293`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "云主机ID，多可用区资源池下，云硬盘和云主机必须在同个az才能支持挂载",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"ebs_ids": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "磁盘ID列表",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(validator2.UUID()),
					listvalidator.SizeBetween(1, 10),
					listvalidator.UniqueValues(),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
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

func (c *CtyunEcsDataVolume) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEcsDataVolumeConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkBeforeAttach(ctx, plan)
	if err != nil {
		return
	}
	tflog.Info(ctx, "尝试绑定的磁盘列表：", map[string]interface{}{"ebsIDs": plan.EbsIDs})
	for _, ebsID := range plan.EbsIDs {
		err = c.attach(ctx, ebsID, plan.InstanceID.ValueString(), plan.RegionID.ValueString())
		if err != nil {
			return
		}
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *CtyunEcsDataVolume) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsDataVolumeConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "不存在") {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (c *CtyunEcsDataVolume) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (c *CtyunEcsDataVolume) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsDataVolumeConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	for _, ebsID := range state.EbsIDs {
		err = c.detach(ctx, ebsID, state.InstanceID.ValueString(), state.RegionID.ValueString())
		if err != nil {
			return
		}
	}
}

func (c *CtyunEcsDataVolume) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ecsId],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunEcsDataVolumeConfig
	var ecsId, regionId string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		ecsId = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &ecsId, &regionId)
		if err != nil {
			return
		}
	}

	if ecsId == "" {
		err = fmt.Errorf("ecsId不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}

	cfg.InstanceID = types.StringValue(ecsId)
	cfg.RegionID = types.StringValue(regionId)

	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *CtyunEcsDataVolume) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ebsService = business.NewEbsService(meta)
	c.ecsService = business.NewEcsService(meta)
	c.jobLooper = business.NewGeneralJobHelper(c.meta.Apis.CtEcsApis.JobShowApi)
}

// checkBeforeAttach 挂载前检查
func (c *CtyunEcsDataVolume) checkBeforeAttach(ctx context.Context, plan CtyunEcsDataVolumeConfig) (err error) {
	err = c.ecsService.CheckEcsStatus(ctx, plan.InstanceID.ValueString(), plan.RegionID.ValueString())
	if err != nil {
		return
	}
	for _, ebsID := range plan.EbsIDs {
		err = c.ebsService.MustExist(ctx, ebsID, plan.RegionID.ValueString())
		if err != nil {
			return
		}
	}
	return
}

// attach 云主机挂载数据盘
func (c *CtyunEcsDataVolume) attach(ctx context.Context, ebsID, instanceID, regionID string) (err error) {
	params := ctecs.CtecsAttachVolumeV41Request{
		DiskID:     ebsID,
		RegionID:   regionID,
		InstanceID: instanceID,
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsAttachVolumeV41Api.Do(ctx, c.meta.SdkCredential, &params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	_, err = c.jobLooper.JobLoop(ctx, c.meta.Credential, regionID, resp.ReturnObj.DiskJobID)
	if err != nil {
		return
	}
	return
}

// attach 云主机解绑数据盘
func (c *CtyunEcsDataVolume) detach(ctx context.Context, ebsID, instanceID, regionID string) (err error) {
	params := ctecs.CtecsDetachVolumeV41Request{
		DiskID:     ebsID,
		RegionID:   regionID,
		InstanceID: instanceID,
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsDetachVolumeV41Api.Do(ctx, c.meta.SdkCredential, &params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	_, err = c.jobLooper.JobLoop(ctx, c.meta.Credential, regionID, resp.ReturnObj.DiskJobID)
	if err != nil {
		return
	}
	return
}

// getAndMerge 查询绑定关系
func (c *CtyunEcsDataVolume) getAndMerge(ctx context.Context, cfg *CtyunEcsDataVolumeConfig) (err error) {
	volumes, err := c.ecsService.GetEcsAttachedVolume(ctx, cfg.InstanceID.ValueString(), cfg.RegionID.ValueString())
	if err != nil {
		return
	}
	if len(volumes) == 0 {
		err = fmt.Errorf("can't find any attached volumes")
	}
	cfg.EbsIDs = volumes[1:]
	cfg.ID = types.StringValue(fmt.Sprintf("%s,%s", cfg.InstanceID.ValueString(), cfg.RegionID.ValueString()))
	return
}

type CtyunEcsDataVolumeConfig struct {
	ID         types.String `tfsdk:"id"`
	InstanceID types.String `tfsdk:"instance_id"`
	RegionID   types.String `tfsdk:"region_id"`
	EbsIDs     []string     `tfsdk:"ebs_ids"`
}
