package ec

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

var (
	_ resource.Resource                = &CtyunEcPacket{}
	_ resource.ResourceWithConfigure   = &CtyunEcPacket{}
	_ resource.ResourceWithImportState = &CtyunEcPacket{}
)

func NewCtyunEcPacket() resource.Resource {
	return &CtyunEcPacket{}
}

type CtyunEcPacket struct {
	meta *common.CtyunMetadata
}

type CtyunEcPacketConfig struct {
	ID                   types.String `tfsdk:"id"`
	EcID                 types.String `tfsdk:"ec_id"`
	RegionID             types.String `tfsdk:"region_id"`
	Name                 types.String `tfsdk:"name"`
	Bandwidth            types.Int64  `tfsdk:"bandwidth"`
	CycleType            types.String `tfsdk:"cycle_type"`
	CycleCount           types.Int64  `tfsdk:"cycle_count"`
	AreaA                types.String `tfsdk:"area_a"`
	AreaB                types.String `tfsdk:"area_b"`
	PayVoucherPrice      types.String `tfsdk:"pay_voucher_price"`
	MasterOrderID        types.String `tfsdk:"master_order_id"`
	MasterOrderNO        types.String `tfsdk:"master_order_no"`
	MasterResourceID     types.String `tfsdk:"master_resource_id"`
	MasterResourceStatus types.String `tfsdk:"master_resource_status"`
	ResourceID           types.String `tfsdk:"resource_id"` // 添加ResourceID字段用于后续操作
}

func (c *CtyunEcPacket) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ec_packet"
}

func (c *CtyunEcPacket) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026763/10038220`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ec_id": schema.StringAttribute{
				Required:    true,
				Description: "云间高速ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
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
					stringvalidator.LengthAtLeast(1),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "带宽包名字",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"bandwidth": schema.Int64Attribute{
				Required:    true,
				Description: "带宽，单位MB",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围：month：按月，year：按年",
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_count": schema.Int64Attribute{
				Required:    true,
				Description: "订购时长，当cycle_type=month，支持订购1-11个月；当cycle_type=year，支持订购1-3年",
				Validators: []validator.Int64{
					validator2.AlsoRequiresEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeMonth),
						types.StringValue(business.OrderCycleTypeYear),
					),
					validator2.ConflictsWithEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
					validator2.CycleCount(1, 11, 1, 3),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"area_a": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "区域A类型，取值如下：china: 中国大陆, APAC:亚太，默认china",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("china", "APAC"),
				},
				Default: stringdefault.StaticString("china"),
			},
			"area_b": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "区域B类型，取值如下：china: 中国大陆, APAC:亚太，默认china",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("china", "APAC"),
				},
				Default: stringdefault.StaticString("china"),
			},
			"pay_voucher_price": schema.StringAttribute{
				Optional:    true,
				Description: "代金券金额，只适用于预付费客户自动支付，若代金券支付金额传0或者控制符，则不适用代金券支付（小数会只保留2位，非负）",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"master_order_no": schema.StringAttribute{
				Computed:    true,
				Description: "订单号",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"master_resource_id": schema.StringAttribute{
				Computed:    true,
				Description: "主资源ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"master_resource_status": schema.StringAttribute{
				Computed:    true,
				Description: "主资源状态，只有主订单资源会返回",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "云间高速带宽包资源ID，用于升配、续订、退订等操作",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (c *CtyunEcPacket) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunEcPacket) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEcPacketConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *CtyunEcPacket) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcPacketConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	// 读取操作成功，保持当前状态
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func (c *CtyunEcPacket) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan, state CtyunEcPacketConfig

	// 获取计划状态和当前状态
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// 检查是否是带宽变更（升配）
	if !plan.Bandwidth.Equal(state.Bandwidth) {
		err := c.upgrade(ctx, &plan, &state)
		if err != nil {
			resp.Diagnostics.AddError("升配带宽包失败", err.Error())
			return
		}
	}

	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *CtyunEcPacket) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcPacketConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// 执行退订操作
	err = c.refund(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunEcPacket) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ecId],[resourceId]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunEcPacketConfig

	var ecId, resourceId string

	err = terraform_extend.Split(request.ID, &ecId, &resourceId)
	if err != nil {
		return
	}

	if ecId == "" {
		err = fmt.Errorf("ecId不能为空")
		return
	}
	if resourceId == "" {
		err = fmt.Errorf("resourceId不能为空")
		return
	}

	config.EcID = types.StringValue(ecId)
	config.ResourceID = types.StringValue(resourceId)

	// 查询远端
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}

	// 设置ID字段，确保导入时有正确的ID（修复：保持与Create中一致的格式）
	config.ID = types.StringValue(fmt.Sprintf("%s,%s", ecId, resourceId))

	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunEcPacket) getAndMerge(ctx context.Context, state *CtyunEcPacketConfig) (err error) {
	// 使用API查询带宽包信息
	listReq := &ec.EcEcPacketListPacketRequest{
		EcID: state.EcID.ValueString(),
	}

	// 如果ResourceID存在，使用它来查询特定资源
	if !state.ResourceID.IsNull() && !state.ResourceID.IsUnknown() {
		listReq.ResourceID = state.ResourceID.ValueStringPointer()
	}

	tflog.Info(ctx, "查询云间高速带宽包信息", map[string]interface{}{
		"ec_id":       state.EcID.ValueString(),
		"resource_id": state.ResourceID.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcPacketListPacketApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		return
	} else if *resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 检查返回结果
	if resp.ReturnObj != nil && resp.ReturnObj.Results != nil && len(resp.ReturnObj.Results) > 0 {
		result := resp.ReturnObj.Results[0]

		// 更新状态信息
		if result.PacketName != nil {
			state.Name = types.StringValue(*result.PacketName)
		}

		if result.Rate != nil {
			state.Bandwidth = types.Int64Value(int64(*result.Rate))
		}

		if result.ResourceID != nil {
			state.ResourceID = types.StringValue(*result.ResourceID)
		}

		if result.PacketID != nil {
			state.ID = types.StringValue(*result.PacketID)
		}
	}

	return
}

// 升配带宽包
func (c *CtyunEcPacket) upgrade(ctx context.Context, plan, state *CtyunEcPacketConfig) (err error) {
	// 只有当resource_id存在时才能执行升配操作
	if state.ResourceID.IsNull() || state.ResourceID.IsUnknown() {
		return fmt.Errorf("无法执行升配操作：ResourceID为空")
	}
	clientToken := uuid.NewString()
	upgradeReq := &ec.EcEcOrderPacketUpgradeRequest{
		EcID:        plan.EcID.ValueString(),
		RegionID:    plan.RegionID.ValueString(),
		Bandwidth:   int32(plan.Bandwidth.ValueInt64()),
		ResourceID:  state.ResourceID.ValueString(),
		ClientToken: &clientToken,
	}

	tflog.Info(ctx, "升配云间高速带宽包", map[string]interface{}{
		"ec_id":       plan.EcID.ValueString(),
		"resource_id": state.ResourceID.ValueString(),
		"bandwidth":   plan.Bandwidth.ValueInt64(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcOrderPacketUpgradeApi.Do(ctx, c.meta.SdkCredential, upgradeReq)
	if err != nil {
		return
	} else if *resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 更新状态信息
	if resp.ReturnObj != nil {
		if resp.ReturnObj.MasterOrderID != nil {
			plan.MasterOrderID = types.StringValue(*resp.ReturnObj.MasterOrderID)
		}
		if resp.ReturnObj.MasterOrderNO != nil {
			plan.MasterOrderNO = types.StringValue(*resp.ReturnObj.MasterOrderNO)
		}
	}

	return nil
}

// 续订带宽包
func (c *CtyunEcPacket) renew(ctx context.Context, plan, state *CtyunEcPacketConfig) (err error) {
	// 只有当resource_id存在时才能执行续订操作
	if state.ResourceID.IsNull() || state.ResourceID.IsUnknown() {
		return fmt.Errorf("无法执行续订操作：ResourceID为空")
	}
	clientToken := uuid.NewString()
	renewReq := &ec.EcEcOrderPacketRenewRequest{
		EcID:        plan.EcID.ValueString(),
		RegionID:    plan.RegionID.ValueString(),
		ResourceID:  state.ResourceID.ValueString(),
		CycleType:   plan.CycleType.ValueString(),
		CycleCount:  int32(plan.CycleCount.ValueInt64()),
		ClientToken: &clientToken,
	}

	tflog.Info(ctx, "续订云间高速带宽包", map[string]interface{}{
		"ec_id":       plan.EcID.ValueString(),
		"resource_id": state.ResourceID.ValueString(),
		"cycle_type":  plan.CycleType.ValueString(),
		"cycle_count": plan.CycleCount.ValueInt64(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcOrderPacketRenewApi.Do(ctx, c.meta.SdkCredential, renewReq)
	if err != nil {
		return
	} else if *resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 更新状态信息
	if resp.ReturnObj != nil {
		if resp.ReturnObj.MasterOrderID != nil {
			plan.MasterOrderID = types.StringValue(*resp.ReturnObj.MasterOrderID)
		}
		if resp.ReturnObj.MasterOrderNO != nil {
			plan.MasterOrderNO = types.StringValue(*resp.ReturnObj.MasterOrderNO)
		}
	}

	return nil
}

// 退订带宽包
func (c *CtyunEcPacket) refund(ctx context.Context, state *CtyunEcPacketConfig) (err error) {
	// 只有当resource_id存在时才能执行退订操作
	if state.ResourceID.IsNull() || state.ResourceID.IsUnknown() {
		return fmt.Errorf("无法执行退订操作：ResourceID为空")
	}
	clientToken := uuid.NewString()
	refundReq := &ec.EcEcOrderPacketRefundRequest{
		EcID:        state.EcID.ValueString(),
		RegionID:    state.RegionID.ValueString(),
		ResourceID:  state.ResourceID.ValueString(),
		ClientToken: &clientToken,
	}
	tflog.Info(ctx, "退订云间高速带宽包", map[string]interface{}{
		"ec_id":       state.EcID.ValueString(),
		"resource_id": state.ResourceID.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcOrderPacketRefundApi.Do(ctx, c.meta.SdkCredential, refundReq)
	if err != nil {
		return
	} else if *resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 更新状态信息
	if resp.ReturnObj != nil {
		if resp.ReturnObj.MasterOrderID != nil {
			state.MasterOrderID = types.StringValue(*resp.ReturnObj.MasterOrderID)
		}
		if resp.ReturnObj.MasterOrderNO != nil {
			state.MasterOrderNO = types.StringValue(*resp.ReturnObj.MasterOrderNO)
		}
	}

	return nil
}

func (c *CtyunEcPacket) create(ctx context.Context, plan *CtyunEcPacketConfig) (err error) {
	// 创建云间高速带宽包订购订单
	clientToken := uuid.NewString()
	newReq := &ec.EcEcOrderPacketNewRequest{
		EcID:        plan.EcID.ValueString(),
		RegionID:    plan.RegionID.ValueString(),
		PacketName:  plan.Name.ValueString(),
		Bandwidth:   int32(plan.Bandwidth.ValueInt64()),
		AreaA:       plan.AreaA.ValueStringPointer(),
		CycleType:   strings.ToUpper(plan.CycleType.ValueString()),
		CycleCount:  int32(plan.CycleCount.ValueInt64()),
		AreaB:       plan.AreaB.ValueStringPointer(),
		ClientToken: &clientToken,
	}

	if !plan.AreaA.IsNull() {
		areaA := plan.AreaA.ValueString()
		newReq.AreaA = &areaA
	}

	if !plan.AreaB.IsNull() {
		areaB := plan.AreaB.ValueString()
		newReq.AreaB = &areaB
	}

	if !plan.PayVoucherPrice.IsNull() {
		price := plan.PayVoucherPrice.ValueString()
		newReq.PayVoucherPrice = &price
	}

	tflog.Info(ctx, "创建云间高速带宽包订购订单", map[string]interface{}{
		"ec_id": plan.EcID.ValueString(),
		"name":  plan.Name.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcOrderPacketNewApi.Do(ctx, c.meta.SdkCredential, newReq)
	if err != nil {
		return
	} else if *resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	if resp.ReturnObj != nil {
		if resp.ReturnObj.MasterOrderID != nil {
			plan.MasterOrderID = types.StringValue(*resp.ReturnObj.MasterOrderID)
		}
		if resp.ReturnObj.MasterOrderNO != nil {
			plan.MasterOrderNO = types.StringValue(*resp.ReturnObj.MasterOrderNO)
		}
		if resp.ReturnObj.MasterResourceID != nil {
			plan.MasterResourceID = types.StringValue(*resp.ReturnObj.MasterResourceID)
		}
		if resp.ReturnObj.MasterResourceStatus != nil {
			plan.MasterResourceStatus = types.StringValue(*resp.ReturnObj.MasterResourceStatus)
		}

		// 设置资源明细中的ResourceID，用于后续操作（升配、续订、退订）
		if len(resp.ReturnObj.Resources) > 0 && resp.ReturnObj.Resources[0].ResourceID != nil {
			plan.ResourceID = types.StringValue(*resp.ReturnObj.Resources[0].ResourceID)
		} else if resp.ReturnObj.MasterResourceID != nil {
			// 如果Resources为空，则使用MasterResourceID作为备选
			plan.ResourceID = types.StringValue(*resp.ReturnObj.MasterResourceID)
		}
	}
	// 设置返回值
	return nil
}
