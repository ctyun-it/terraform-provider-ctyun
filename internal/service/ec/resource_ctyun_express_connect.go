package ec

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"time"
)

var (
	_ resource.Resource                = &CtyunExpressConnect{}
	_ resource.ResourceWithConfigure   = &CtyunExpressConnect{}
	_ resource.ResourceWithImportState = &CtyunExpressConnect{}
)

func NewCtyunExpressConnect() resource.Resource {
	return &CtyunExpressConnect{}
}

type CtyunExpressConnect struct {
	meta *common.CtyunMetadata
}

type CtyunExpressConnectConfig struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Status      types.Int64  `tfsdk:"status"`
	CreateTime  types.String `tfsdk:"create_time"`
	RegionId    types.String `tfsdk:"region_id"`
}

func (c *CtyunExpressConnect) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_express_connect"
}

func (c *CtyunExpressConnect) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026763/10038220`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "云间高速实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "名称 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "描述信息 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"status": schema.Int64Attribute{
				Computed:    true,
				Description: "运行状态，取值范围: 0:创建中 2:运行中 18:删除中 21:设置中 22:更新带宽中 24:更新中",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
		},
	}
}

func (c *CtyunExpressConnect) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunExpressConnect) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunExpressConnectConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeCreate(ctx, &plan)
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *CtyunExpressConnect) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunExpressConnectConfig
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

func (c *CtyunExpressConnect) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunExpressConnectConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, plan)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *CtyunExpressConnect) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunExpressConnectConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err = c.deleteCgwBill(ctx, &state)
	if err != nil {
		return
	}
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *CtyunExpressConnect) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunExpressConnectConfig
	config.ID = types.StringValue(request.ID)

	// 查询远端
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunExpressConnect) checkBeforeCreate(ctx context.Context, c2 *CtyunExpressConnectConfig) (err error) {
	return nil
}
func (c *CtyunExpressConnect) create(ctx context.Context, plan *CtyunExpressConnectConfig) (err error) {
	// 创建云间高速实例
	createReq := &ec.EcEcCreateRequest{
		EcName: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() {
		createReq.EcDescription = plan.Description.ValueStringPointer()
	}

	tflog.Info(ctx, "创建云间高速实例", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcCreateApi.Do(ctx, c.meta.SdkCredential, createReq)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}
	plan.ID = types.StringValue(*resp.ReturnObj.EcID)
	plan.Status = types.Int64Value(int64(*resp.ReturnObj.Status))
	return
}
func (c *CtyunExpressConnect) getAndMerge(ctx context.Context, plan *CtyunExpressConnectConfig) (err error) {
	// 查询云间高速实例
	listReq := &ec.EcEcListRequest{
		EcID:     plan.ID.ValueStringPointer(),
		PageNo:   func() *int32 { i := int32(1); return &i }(),
		PageSize: func() *int32 { i := int32(1); return &i }(),
	}

	resp, err := c.meta.Apis.SdkEcApis.EcEcListApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if len(resp.ReturnObj.Results) == 0 {
		return fmt.Errorf("no express connect instance found")
	}
	result := resp.ReturnObj.Results[0]
	plan.Name = types.StringValue(*result.EcName)
	plan.Status = types.Int64Value(int64(*result.Status))
	plan.CreateTime = types.StringValue(utils.FromBJTimeToUTCZ(utils.SecString(result.CreateDate)))

	if result.EcDescription != nil {
		plan.Description = types.StringValue(*result.EcDescription)
	}

	return
}
func (c *CtyunExpressConnect) update(ctx context.Context, plan CtyunExpressConnectConfig) (err error) {
	// 更新云间高速实例
	updateReq := &ec.EcEcUpdateRequest{
		EcID:   plan.ID.ValueString(),
		EcName: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() {
		updateReq.EcDescription = plan.Description.ValueStringPointer()
	}

	tflog.Info(ctx, "更新云间高速实例", map[string]interface{}{
		"id":   plan.ID.ValueString(),
		"name": plan.Name.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcUpdateApi.Do(ctx, c.meta.SdkCredential, updateReq)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}
	return
}
func (c *CtyunExpressConnect) delete(ctx context.Context, state CtyunExpressConnectConfig) (err error) {
	// 删除云间高速实例
	deleteReq := &ec.EcEcDeleteRequest{
		EcID: state.ID.ValueString(),
	}

	tflog.Info(ctx, "删除云间高速实例", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcDeleteApi.Do(ctx, c.meta.SdkCredential, deleteReq)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}

	return
}

// deleteCgwBill 删除云网关计费
func (c *CtyunExpressConnect) deleteCgwBill(ctx context.Context, state *CtyunExpressConnectConfig) (err error) {
	// 实现删除云网关计费的逻辑
	// 使用 c.meta.Apis.SdkEcApis.EcEcCgwBillRefundApi 调用新创建的退订API
	oderType := "1"
	oderState := "1"
	// 构造查询参数（这里需要根据实际业务需求进行调整）
	queryReq := &ec.EcEcTgwOrderQueryRequest{
		EcID:       state.ID.ValueString(),
		OrderType:  &oderType,
		OrderState: &oderState,
	}
	queryResp, err := c.meta.Apis.SdkEcApis.EcEcTgwOrderQueryApi.Do(ctx, c.meta.SdkCredential, queryReq)
	if err != nil {
		return
	} else if *queryResp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *queryResp.Message)
	} else if queryResp.ReturnObj == nil || len(queryResp.ReturnObj.Results) == 0 {
		return

	} else {
		resourceID := queryResp.ReturnObj.Results[0].ResourceID
		// 构造请求参数（这里需要根据实际业务需求进行调整）
		req := &ec.EcEcCgwBillRefundRequest{
			EcID:       state.ID.ValueString(),
			RegionID:   state.RegionId.ValueString(), // 使用实际的RegionID
			ResourceID: *resourceID,                  // 使用实际的ResourceID
			// ClientToken 参数根据实际需求添加
		}

		tflog.Info(ctx, "删除云网关计费", map[string]interface{}{
			"ec_id": state.ID.ValueString(),
		})
		var resp *ec.EcEcCgwBillRefundResponse
		resp, err = c.meta.Apis.SdkEcApis.EcEcCgwBillRefundApi.Do(ctx, c.meta.SdkCredential, req)
		if err != nil {
			return
		} else if *resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return
		} else if resp.ReturnObj == nil {
			err = common.InvalidReturnObjError
			return
		}

		// 轮询查询订单状态，确保删除完成
		// 最多轮询30次，每次间隔5秒
		for i := 0; i < 30; i++ {
			time.Sleep(5 * time.Second)
			queryResp, err := c.meta.Apis.SdkEcApis.EcEcTgwOrderQueryApi.Do(ctx, c.meta.SdkCredential, queryReq)
			if err != nil {
				// 查询失败，记录日志但继续轮询
				tflog.Warn(ctx, "查询订单状态失败", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			// 检查API响应状态
			if queryResp.StatusCode == nil {
				tflog.Warn(ctx, "查询订单状态失败，StatusCode为空")
				continue
			} else if *queryResp.StatusCode != common.NormalStatusCode {
				// API返回错误，记录日志但继续轮询
				tflog.Warn(ctx, "查询订单状态失败", map[string]interface{}{
					"message": func() string {
						if queryResp.Message != nil {
							return *queryResp.Message
						}
						return "unknown error"
					}(),
					"description": func() string {
						if queryResp.Description != nil {
							return *queryResp.Description
						}
						return "unknown description"
					}(),
				})
				continue
			}

			// 检查返回结果
			if queryResp.ReturnObj != nil && queryResp.ReturnObj.Results != nil {
				// 如果没有查询到订单，说明删除已完成
				if len(queryResp.ReturnObj.Results) == 0 {
					tflog.Info(ctx, "确认云间高速实例已成功删除")
					break
				}

				// 如果查询到订单，继续轮询
				tflog.Info(ctx, "云间高速实例仍在删除中", map[string]interface{}{
					"result_count": len(queryResp.ReturnObj.Results),
				})
			}

			// 如果是最后一次轮询，仍然查询到订单，则记录警告
			if i == 29 && queryResp.ReturnObj != nil && queryResp.ReturnObj.Results != nil && len(queryResp.ReturnObj.Results) > 0 {
				tflog.Warn(ctx, "轮询结束但仍未确认删除完成", map[string]interface{}{
					"result_count": len(queryResp.ReturnObj.Results),
				})
			}
		}
	}

	return nil
}
