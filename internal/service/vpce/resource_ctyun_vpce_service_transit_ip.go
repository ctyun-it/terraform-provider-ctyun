package vpce

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
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
	_ resource.Resource                = &ctyunVpceServiceTransitIP{}
	_ resource.ResourceWithConfigure   = &ctyunVpceServiceTransitIP{}
	_ resource.ResourceWithImportState = &ctyunVpceServiceTransitIP{}
)

type ctyunVpceServiceTransitIP struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpceServiceTransitIP() resource.Resource {
	return &ctyunVpceServiceTransitIP{}
}

func (c *ctyunVpceServiceTransitIP) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpce_service_transit_ip"
}

type CtyunVpceServiceTransitIPConfig struct {
	ID                types.String `tfsdk:"id"`
	EndpointServiceID types.String `tfsdk:"endpoint_service_id"`
	RegionID          types.String `tfsdk:"region_id"`
	SubnetID          types.String `tfsdk:"subnet_id"`
	TransitIP         types.String `tfsdk:"transit_ip"`
	CreateTime        types.String `tfsdk:"create_time"`
	UpdateTime        types.String `tfsdk:"update_time"`
}

func (c *ctyunVpceServiceTransitIP) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10042658/10048507`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID，使用中转IP地址，和transit_ip相等",
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
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"transit_ip": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "中转IP地址",
				Validators: []validator.String{
					validator2.Ip(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"create_time": schema.StringAttribute{
				Description: "创建时间，为UTC格式",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Description: "更新时间，为UTC格式",
				Computed:    true,
			},
		},
	}
}

func (c *ctyunVpceServiceTransitIP) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunVpceServiceTransitIPConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建
	ip, err := c.create(ctx, plan)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(ip)
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunVpceServiceTransitIP) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpceServiceTransitIPConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "endpointServiceID ensure") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunVpceServiceTransitIP) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (c *ctyunVpceServiceTransitIP) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpceServiceTransitIPConfig
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

func (c *ctyunVpceServiceTransitIP) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunVpceServiceTransitIP) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[endpointServiceID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunVpceServiceTransitIPConfig
	var ip, endpointServiceID, regionID string
	// 根据分隔符数量判断是否输入了 regionID
	if strings.Count(request.ID, common.ImportSeparator) == 1 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &ip, &endpointServiceID)

		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &ip, &endpointServiceID, &regionID)
		if err != nil {
			return
		}
	}

	if ip == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if endpointServiceID == "" {
		err = fmt.Errorf("endpointServiceID不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.EndpointServiceID = types.StringValue(endpointServiceID)
	cfg.ID = types.StringValue(ip)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// create 创建
func (c *ctyunVpceServiceTransitIP) create(ctx context.Context, plan CtyunVpceServiceTransitIPConfig) (ip string, err error) {
	params := &ctvpc.CtvpcCreateEndpointServiceTransitIPRequest{
		ClientToken:       uuid.NewString(),
		RegionID:          plan.RegionID.ValueString(),
		SubnetID:          plan.SubnetID.ValueString(),
		EndpointServiceID: plan.EndpointServiceID.ValueString(),
	}
	transitIP := plan.TransitIP.ValueString()
	if transitIP != "" {
		params.TransitIP = &transitIP
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateEndpointServiceTransitIPApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	ip = resp.ReturnObj.TransitIP
	return
}

// getAndMerge 从远端查询
func (c *ctyunVpceServiceTransitIP) getAndMerge(ctx context.Context, plan *CtyunVpceServiceTransitIPConfig) (err error) {
	params := &ctvpc.CtvpcListEndpointServiceTransitIPRequest{
		RegionID:          plan.RegionID.ValueString(),
		EndpointServiceID: plan.EndpointServiceID.ValueString(),
		PageSize:          50,
		PageNo:            1,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListEndpointServiceTransitIPApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var exist bool
	for _, ip := range resp.ReturnObj.TransitIPs {
		if utils.SecStringValue(ip.TransitIP) == plan.ID {
			plan.SubnetID = utils.SecStringValue(ip.SubnetID)
			plan.TransitIP = utils.SecStringValue(ip.TransitIP)
			plan.CreateTime = utils.SecStringValue(ip.CreatedAt)
			plan.UpdateTime = utils.SecStringValue(ip.UpdatedAt)
			exist = true
		}
	}
	if !exist {
		err = common.InvalidReturnObjResultsError
		return
	}

	return
}

// delete 删除
func (c *ctyunVpceServiceTransitIP) delete(ctx context.Context, plan CtyunVpceServiceTransitIPConfig) (err error) {
	params := &ctvpc.CtvpcDeleteEndpointServiceTransitIPRequest{
		RegionID:          plan.RegionID.ValueString(),
		EndpointServiceID: plan.EndpointServiceID.ValueString(),
		TransitIP:         plan.TransitIP.ValueString(),
		ClientToken:       uuid.NewString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteEndpointServiceTransitIPApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}
