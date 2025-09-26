package vpce

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &ctyunVpceServiceConnection{}
	_ resource.ResourceWithConfigure   = &ctyunVpceServiceConnection{}
	_ resource.ResourceWithImportState = &ctyunVpceServiceConnection{}
)

type ctyunVpceServiceConnection struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpceServiceConnection() resource.Resource {
	return &ctyunVpceServiceConnection{}
}

func (c *ctyunVpceServiceConnection) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpce_service_connection"
}

type CtyunVpceServiceConnectionConfig struct {
	ID                types.String `tfsdk:"id"`
	RegionID          types.String `tfsdk:"region_id"`
	EndpointServiceID types.String `tfsdk:"endpoint_service_id"`
	EndpointID        types.String `tfsdk:"endpoint_id"`
	Status            types.String `tfsdk:"status"`
}

func (c *ctyunVpceServiceConnection) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10042658/10043026`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
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
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"endpoint_service_id": schema.StringAttribute{
				Required:    true,
				Description: "终端节点服务ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"endpoint_id": schema.StringAttribute{
				Required:    true,
				Description: "终端节点ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"status": schema.StringAttribute{
				Required:    true,
				Description: "连接状态，up表示申请连接，down表示断开连接，支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.VpceServiceConnectionUp, business.VpceServiceConnectionDown),
				},
			},
		},
	}
}

func (c *ctyunVpceServiceConnection) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunVpceServiceConnectionConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建
	switch plan.Status.ValueString() {
	case business.VpceServiceConnectionUp:
		err = c.acceptApply(ctx, plan)
	case business.VpceServiceConnectionDown:
		err = c.rejectApply(ctx, plan)
	}
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

func (c *ctyunVpceServiceConnection) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpceServiceConnectionConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not exist apply") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunVpceServiceConnection) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunVpceServiceConnectionConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunVpceServiceConnectionConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新
	switch plan.Status.ValueString() {
	case business.VpceServiceConnectionUp:
		err = c.acceptApply(ctx, plan)
	case business.VpceServiceConnectionDown:
		err = c.rejectApply(ctx, plan)
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

func (c *ctyunVpceServiceConnection) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	// connection没有删除操作
}

func (c *ctyunVpceServiceConnection) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [enpointServiceID],[endpointID],[regionID]
func (c *ctyunVpceServiceConnection) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunVpceServiceConnectionConfig
	var enpointServiceID, endpointID, regionID string
	err = terraform_extend.Split(request.ID, &enpointServiceID, &endpointID, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.EndpointID = types.StringValue(endpointID)
	cfg.EndpointServiceID = types.StringValue(enpointServiceID)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// getAndMerge 从远端查询
func (c *ctyunVpceServiceConnection) getAndMerge(ctx context.Context, plan *CtyunVpceServiceConnectionConfig) (err error) {
	params := &ctvpc.CtvpcShowEndpointServiceConnectionsRequest{
		RegionID:          plan.RegionID.ValueString(),
		EndpointServiceID: plan.EndpointServiceID.ValueString(),
		PageNo:            1,
		PageSize:          50,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowEndpointServiceConnectionsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	for _, con := range resp.ReturnObj {
		if utils.SecStringValue(con.EndpointID) == plan.EndpointID {
			plan.ID = types.StringValue(fmt.Sprintf("%s,%s", plan.EndpointServiceID.ValueString(), plan.EndpointID.ValueString()))
			if utils.SecString(con.ConnectionStatus) == "已连接" {
				plan.Status = types.StringValue(business.VpceServiceConnectionUp)
			} else {
				plan.Status = types.StringValue(business.VpceServiceConnectionDown)
			}
			return
		}
	}
	err = fmt.Errorf("apply not exist")
	return
}

// acceptApply允许连接
func (c *ctyunVpceServiceConnection) acceptApply(ctx context.Context, plan CtyunVpceServiceConnectionConfig) (err error) {
	params := &ctvpc.CtvpcAcceptEndpointApplyRequest{
		RegionID:          plan.RegionID.ValueString(),
		EndpointID:        plan.EndpointID.ValueString(),
		EndpointServiceID: plan.EndpointServiceID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcAcceptEndpointApplyApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}

	return
}

// rejectApply 拒绝连接
func (c *ctyunVpceServiceConnection) rejectApply(ctx context.Context, plan CtyunVpceServiceConnectionConfig) (err error) {
	params := &ctvpc.CtvpcRefuseEndpointApplyRequest{
		RegionID:          plan.RegionID.ValueString(),
		EndpointID:        plan.EndpointID.ValueString(),
		EndpointServiceID: plan.EndpointServiceID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcRefuseEndpointApplyApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}

	return
}
