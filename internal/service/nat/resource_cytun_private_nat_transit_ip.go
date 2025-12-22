package nat

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctnat"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &ctyunPrivateNatTransitIpResource{}
	_ resource.ResourceWithConfigure   = &ctyunPrivateNatTransitIpResource{}
	_ resource.ResourceWithImportState = &ctyunPrivateNatTransitIpResource{}
)

type ctyunPrivateNatTransitIpResource struct {
	meta *common.CtyunMetadata
}

func NewCtyunPrivateNatTransitIpResource() resource.Resource {
	return &ctyunPrivateNatTransitIpResource{}
}

func (c *ctyunPrivateNatTransitIpResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_private_nat_transit_ip"
}

func (c *ctyunPrivateNatTransitIpResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026759/10378390`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "中转IP的ID，格式为regionID:natGatewayID:address",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，默认使用provider ctyun总region_id 或者环境变量",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"nat_gateway_id": schema.StringAttribute{
				Required:    true,
				Description: "私网NAT网关ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"address": schema.StringAttribute{
				Required:    true,
				Description: "中转IP地址，必须在中转网段指定的网络范围内",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.Ip(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "中转IP状态: running代表运行中, freeze代表已冻结, expired代表已到期",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_default": schema.BoolAttribute{
				Computed:    true,
				Description: "是否为默认中转地址",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"snat_count": schema.Int32Attribute{
				Computed:    true,
				Description: "在使用此中转IP的snat数量",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"dnat_count": schema.Int32Attribute{
				Computed:    true,
				Description: "在使用此中转IP的dnat数量",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *ctyunPrivateNatTransitIpResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunPrivateNatTransitIpConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建中转IP
	createResp, err := c.meta.Apis.SdkCtNatApis.CtnatCreatePrivatenatIPApi.Do(ctx, c.meta.SdkCredential, &ctnat.CtnatCreatePrivatenatIPRequest{
		RegionID:     plan.RegionID.ValueString(),
		NatGatewayID: plan.NatGatewayID.ValueString(),
		Address:      plan.Address.ValueString(),
	})
	if err != nil {
		return
	} else if createResp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", createResp.Message, createResp.Description)
		return
	}

	// 设置ID
	id := fmt.Sprintf("%s,%s,%s", plan.RegionID.ValueString(), plan.NatGatewayID.ValueString(), plan.Address.ValueString())
	plan.ID = types.StringValue(id)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建后反查创建后的中转IP信息
	err = c.getAndMergePrivateNatTransitIp(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunPrivateNatTransitIpResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunPrivateNatTransitIpConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergePrivateNatTransitIp(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(request.State.Set(ctx, &state)...)
}

func (c *ctyunPrivateNatTransitIpResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// 中转IP不支持更新操作，直接将计划状态设置为当前状态
	// 由于所有属性都需要替换资源，实际上不会执行更新操作
	response.Diagnostics.AddError(
		"不支持更新操作",
		"中转IP不支持更新操作，如需修改请先删除再重新创建",
	)
}

func (c *ctyunPrivateNatTransitIpResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunPrivateNatTransitIpConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 删除中转IP
	resp, err := c.meta.Apis.SdkCtNatApis.CtnatDeletePrivatenatIPApi.Do(ctx, c.meta.SdkCredential, &ctnat.CtnatDeletePrivatenatIPRequest{
		RegionID:     state.RegionID.ValueString(),
		NatGatewayID: state.NatGatewayID.ValueString(),
		Addresses:    []string{state.Address.ValueString()},
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
}

func (c *ctyunPrivateNatTransitIpResource) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)

	c.meta = meta
}

func (c *ctyunPrivateNatTransitIpResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [address],[natGateWayID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunPrivateNatTransitIpConfig
	var ID, regionID, natGateWayID, address string
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &address, &natGateWayID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &address, &natGateWayID, &regionID)
		if err != nil {
			return
		}
	}
	ID = fmt.Sprintf("%s,%s,%s", regionID, natGateWayID, address)

	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	if natGateWayID == "" {
		err = fmt.Errorf("natGateWayID不能为空")
	}
	if address == "" {
		err = fmt.Errorf("address 不能为空")
	}
	config.RegionID = types.StringValue(regionID)
	config.NatGatewayID = types.StringValue(natGateWayID)
	config.ID = types.StringValue(ID)
	config.Address = types.StringValue(address)
	err = c.getAndMergePrivateNatTransitIp(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)

}

func (c *ctyunPrivateNatTransitIpResource) getAndMergePrivateNatTransitIp(ctx context.Context, cfg *CtyunPrivateNatTransitIpConfig) (err error) {
	resp, err := c.meta.Apis.SdkCtNatApis.CtnatQueryPrivatenatIPApi.Do(ctx, c.meta.SdkCredential, &ctnat.CtnatQueryPrivatenatIPRequest{
		RegionID:     cfg.RegionID.ValueString(),
		NatGatewayID: cfg.NatGatewayID.ValueString(),
		Address:      cfg.Address.ValueString(),
	})
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}

	// 查找匹配的中转IP
	var targetIp *ctnat.CtnatQueryPrivatenatIPReturnObjResponse
	for _, ip := range resp.ReturnObj {
		if ip.Address == cfg.Address.ValueString() {
			targetIp = ip
			break
		}
	}

	if targetIp == nil {
		err = fmt.Errorf("transit ip not found")
		return
	}

	cfg.Status = types.StringValue(targetIp.Status)
	if targetIp.IsDefault != nil {
		cfg.IsDefault = types.BoolValue(*targetIp.IsDefault)
	}
	cfg.SnatCount = types.Int32Value(targetIp.SnatCnt)
	cfg.DnatCount = types.Int32Value(targetIp.DnarCnt)

	return nil
}

type CtyunPrivateNatTransitIpConfig struct {
	ID           types.String `tfsdk:"id"`
	RegionID     types.String `tfsdk:"region_id"`      // 区域id
	NatGatewayID types.String `tfsdk:"nat_gateway_id"` // 私网NAT网关ID
	Address      types.String `tfsdk:"address"`        // 中转IP地址
	Status       types.String `tfsdk:"status"`         // 中转IP状态
	IsDefault    types.Bool   `tfsdk:"is_default"`     // 是否为默认中转地址
	SnatCount    types.Int32  `tfsdk:"snat_count"`     // 在使用此中转IP的snat数量
	DnatCount    types.Int32  `tfsdk:"dnat_count"`     // 在使用此中转IP的dnat数量
}
