package ec

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

var (
	_ resource.Resource                = &CtyunEcCloudGateway{}
	_ resource.ResourceWithConfigure   = &CtyunEcCloudGateway{}
	_ resource.ResourceWithImportState = &CtyunEcCloudGateway{}
)

func NewCtyunEcCloudGateway() resource.Resource {
	return &CtyunEcCloudGateway{}
}

type CtyunEcCloudGateway struct {
	meta *common.CtyunMetadata
}

type CtyunEcCloudGatewayConfig struct {
	ID          types.String `tfsdk:"id"`
	EcID        types.String `tfsdk:"ec_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Region      types.Int64  `tfsdk:"region"`
	DcName      types.String `tfsdk:"region_name"`
	RegionID    types.String `tfsdk:"region_id"`
	CreateTime  types.String `tfsdk:"create_time"`
	RtbID       types.String `tfsdk:"rtb_id"`
}

func (c *CtyunEcCloudGateway) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ec_cloud_gateway"
}

func (c *CtyunEcCloudGateway) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026763/10038220`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "云网关实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rtb_id": schema.StringAttribute{
				Computed:    true,
				Description: "云网关实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ec_id": schema.StringAttribute{
				Required:    true,
				Description: "云间高速实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "云网关名称 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,

				Description: "云网关描述  支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"region": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "地域信息，取值如下: 1：中国大陆（默认） 2:亚太",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.OneOf(1, 2),
				},
			},
			"region_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraAzName, true),
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
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *CtyunEcCloudGateway) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunEcCloudGateway) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEcCloudGatewayConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeCreate(ctx, &plan)
	if err != nil {
		return
	}
	err = c.createCgwBill(ctx, &plan)
	if err != nil {
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

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunEcCloudGateway) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcCloudGatewayConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunEcCloudGateway) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEcCloudGatewayConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &plan)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunEcCloudGateway) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcCloudGatewayConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *CtyunEcCloudGateway) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ecId],[cgwID]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunEcCloudGatewayConfig

	var ecId, cgwID string

	err = terraform_extend.Split(request.ID, &ecId, &cgwID)
	if err != nil {
		return
	}
	if ecId == "" {
		err = fmt.Errorf("ecId不能为空")
		return
	}
	if cgwID == "" {
		err = fmt.Errorf("cgwID不能为空")
		return
	}

	config.EcID = types.StringValue(ecId)
	config.ID = types.StringValue(cgwID)

	// 查询远端
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunEcCloudGateway) checkBeforeCreate(ctx context.Context, c2 *CtyunEcCloudGatewayConfig) (err error) {
	return nil
}
func (c *CtyunEcCloudGateway) create(ctx context.Context, plan *CtyunEcCloudGatewayConfig) (err error) {
	// 创建云网关实例
	createReq := &ec.EcEcCreateGatewayRequest{
		CgwName: plan.Name.ValueString(),
		DcName:  plan.DcName.ValueString(),
		DcID:    plan.RegionID.ValueString(),
		EcID:    plan.EcID.ValueString(),
		DcType:  "CNP",
	}

	if !plan.Description.IsNull() {
		createReq.Description = plan.Description.ValueStringPointer()
	}

	if !plan.Region.IsNull() {
		region := int32(plan.Region.ValueInt64())
		createReq.Region = &region
	}

	tflog.Info(ctx, "创建云网关实例", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcCreateGatewayApi.Do(ctx, c.meta.SdkCredential, createReq)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}
	if resp.ReturnObj.CgwID == nil {
		return fmt.Errorf("API return error. CgwID is nil")
	}

	plan.ID = types.StringValue(*resp.ReturnObj.CgwID)

	return
}
func (c *CtyunEcCloudGateway) getAndMerge(ctx context.Context, plan *CtyunEcCloudGatewayConfig) (err error) {
	// 查询云网关实例
	listReq := &ec.EcEcListGatewayRequest{
		EcID:  plan.EcID.ValueString(),
		CgwID: plan.ID.ValueStringPointer(),
	}

	resp, err := c.meta.Apis.SdkEcApis.EcEcListGatewayApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if resp.ReturnObj == nil || len(resp.ReturnObj.Results) == 0 {
		return common.InvalidReturnObjError
	}
	result := resp.ReturnObj.Results[0]
	if result.CgwID == nil {
		return fmt.Errorf("API return error. CgwID is nil")
	}

	if result.CgwName == nil {
		return fmt.Errorf("API return error. CgwName is nil")
	}

	if result.EcID == nil {
		return fmt.Errorf("API return error. EcID is nil")
	}

	if result.DcID == nil {
		return fmt.Errorf("API return error. RegionID is nil")
	}

	if result.DcName == nil {
		return fmt.Errorf("API return error. DcName is nil")
	}

	plan.ID = types.StringValue(*result.CgwID)
	plan.Name = types.StringValue(*result.CgwName)
	plan.EcID = types.StringValue(*result.EcID)
	plan.RegionID = types.StringValue(*result.DcID)
	plan.DcName = types.StringValue(*result.DcName)
	plan.Region = types.Int64Value(*result.Region)
	plan.Description = types.StringValue(*resp.Description)
	//

	if result.Region != nil {
		plan.Region = types.Int64Value(int64(*result.Region))
	}

	if result.CgwDescription != nil {
		plan.Description = types.StringValue(*result.CgwDescription)
	}

	if result.CreateDate != nil {
		plan.CreateTime = types.StringValue(utils.FromBJTimeToUTCZ(*result.CreateDate))
	}
	if result.DefaultRtbID != nil {
		plan.RtbID = types.StringValue(*result.DefaultRtbID)
	}
	return
}
func (c *CtyunEcCloudGateway) update(ctx context.Context, plan *CtyunEcCloudGatewayConfig) (err error) {
	// 更新云网关实例
	updateReq := &ec.EcEcUpdateGatewayRequest{
		CgwID: plan.ID.ValueString(),
	}

	if !plan.Name.IsNull() {
		updateReq.CgwName = plan.Name.ValueStringPointer()
	}

	if !plan.Description.IsNull() {
		updateReq.CgwDescription = plan.Description.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkEcApis.EcEcUpdateGatewayApi.Do(ctx, c.meta.SdkCredential, updateReq)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}
	return
}
func (c *CtyunEcCloudGateway) delete(ctx context.Context, state CtyunEcCloudGatewayConfig) (err error) {
	// 删除云网关实例
	deleteReq := &ec.EcEcDeleteGatewayRequest{
		CgwID: state.ID.ValueString(),
	}

	tflog.Info(ctx, "删除云网关实例", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcDeleteGatewayApi.Do(ctx, c.meta.SdkCredential, deleteReq)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}
	return
}

// createCgwBill 创建云网关计费
func (c *CtyunEcCloudGateway) createCgwBill(ctx context.Context, plan *CtyunEcCloudGatewayConfig) (err error) {
	// 实现创建云网关计费的逻辑
	// 使用 c.meta.Apis.SdkEcApis.EcEcCgwBillNewApi 调用新创建的订购API
	oderType := "1"
	oderState := "1"
	// 构造查询参数（这里需要根据实际业务需求进行调整）
	queryReq := &ec.EcEcTgwOrderQueryRequest{
		EcID:       plan.EcID.ValueString(),
		OrderType:  &oderType,
		OrderState: &oderState,
	}
	queryResp, err := c.meta.Apis.SdkEcApis.EcEcTgwOrderQueryApi.Do(ctx, c.meta.SdkCredential, queryReq)
	if err != nil {
		return
	} else if *queryResp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *queryResp.Message)
	} else if queryResp.ReturnObj == nil || len(queryResp.ReturnObj.Results) == 0 {
		// 构造请求参数（这里需要根据实际业务需求进行调整）
		req := &ec.EcEcCgwBillNewRequest{
			EcID: plan.EcID.ValueString(),
			// RegionID, ClientToken, PayVoucherPrice 等参数根据实际需求添加
		}

		tflog.Info(ctx, "创建云网关计费", map[string]interface{}{
			"ec_id": plan.EcID.ValueString(),
		})
		var resp *ec.EcEcCgwBillNewResponse
		resp, err = c.meta.Apis.SdkEcApis.EcEcCgwBillNewApi.Do(ctx, c.meta.SdkCredential, req)
		if err != nil {
			return
		} else if *resp.StatusCode != common.NormalStatusCode {
			if resp.Message != nil && strings.Contains(*resp.Message, "该用户账号存在云企业路由器计费订单") {
				err = nil
				return
			}
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return
		} else if resp.ReturnObj == nil {
			err = common.InvalidReturnObjError
			return
		}
	}

	return
}
