package ec

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
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
	"time"
)

var (
	_ resource.Resource                = &CtyunEcSdwanInstance{}
	_ resource.ResourceWithConfigure   = &CtyunEcSdwanInstance{}
	_ resource.ResourceWithImportState = &CtyunEcSdwanInstance{}
)

func NewCtyunEcSdwanInstance() resource.Resource {
	return &CtyunEcSdwanInstance{}
}

type CtyunEcSdwanInstance struct {
	meta *common.CtyunMetadata
}

type CtyunEcSdwanInstanceConfig struct {
	ID         types.String `tfsdk:"id"`
	EcID       types.String `tfsdk:"ec_id"`
	CgwID      types.String `tfsdk:"cgw_id"`
	SdwanID    types.String `tfsdk:"sdwan_id"`
	RtbID      types.String `tfsdk:"rtb_id"`
	Weights    types.Int64  `tfsdk:"weights"`
	RouteLearn types.Int64  `tfsdk:"route_learn"`
	RouteSync  types.Int64  `tfsdk:"route_sync"`
}

func (c *CtyunEcSdwanInstance) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ec_sdwan_instance"
}

func (c *CtyunEcSdwanInstance) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026763/10038220`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "网络实例ID",
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
			"cgw_id": schema.StringAttribute{
				Required:    true,
				Description: "云网关ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"sdwan_id": schema.StringAttribute{
				Required:    true,
				Description: "sdwan ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"rtb_id": schema.StringAttribute{
				Required:    true,
				Description: "云网关默认路由表ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"weights": schema.Int64Attribute{
				Optional:    true,
				Description: "权重，sdwan默认60，无冗余实例则不传 支持更新",
				Validators: []validator.Int64{
					int64validator.Between(0, 100),
				},
			},
			"route_learn": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "路由学习开关，开启后云网关自动学习网络实例路由，取值范围: 1:学习 0:不学习，默认学习",
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"route_sync": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "路由同步开关，开启后云网关路由自动同步到网络实例，取值范围: 1:同步 0:不同步，默认同步",
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *CtyunEcSdwanInstance) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunEcSdwanInstance) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunEcSdwanInstanceConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.create(ctx, &plan)
	if err != nil {
		return
	}

	// 通过查询获取真实的InstanceID
	tflog.Info(ctx, "Starting to get instance ID after creation", map[string]interface{}{
		"ec_id":    plan.EcID.ValueString(),
		"sdwan_id": plan.SdwanID.ValueString(),
		"cgw_id":   plan.CgwID.ValueString(),
	})

	instanceID, err := c.getInstanceID(ctx, plan.EcID.ValueString(), plan.SdwanID.ValueString(), plan.CgwID.ValueString())
	if err != nil {
		tflog.Error(ctx, "Failed to get instance ID after creation", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	plan.ID = types.StringValue(instanceID)
	tflog.Info(ctx, "Successfully got instance ID", map[string]interface{}{
		"instance_id": instanceID,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunEcSdwanInstance) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunEcSdwanInstanceConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if err == common.ResourceNotExistError {
			resp.State.RemoveResource(ctx)
			return
		}
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunEcSdwanInstance) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunEcSdwanInstanceConfig
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

func (c *CtyunEcSdwanInstance) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunEcSdwanInstanceConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *CtyunEcSdwanInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	// 支持通过InstanceID导入
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunEcSdwanInstanceConfig
	var ID, ecID string
	err = terraform_extend.Split(request.ID, &ID, &ecID)
	if err != nil {
		return
	}
	config.ID = types.StringValue(ID)
	config.EcID = types.StringValue(ecID)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunEcSdwanInstance) create(ctx context.Context, plan *CtyunEcSdwanInstanceConfig) (err error) {
	req := &ec.EcEcCreateSDWANInstanceRequest{
		EcID:    plan.EcID.ValueString(),
		CgwID:   plan.CgwID.ValueString(),
		SdwanID: plan.SdwanID.ValueString(),
		RtbID:   plan.RtbID.ValueString(),
	}

	if !plan.Weights.IsNull() {
		weights := int32(plan.Weights.ValueInt64())
		req.Weights = &weights
	}

	if !plan.RouteLearn.IsNull() {
		routeLearn := int32(plan.RouteLearn.ValueInt64())
		req.RouteLearn = &routeLearn
	}

	if !plan.RouteSync.IsNull() {
		routeSync := int32(plan.RouteSync.ValueInt64())
		req.RouteSync = &routeSync
	}

	tflog.Info(ctx, "Creating SDWAN network instance", map[string]interface{}{
		"ec_id":    plan.EcID.ValueString(),
		"cgw_id":   plan.CgwID.ValueString(),
		"sdwan_id": plan.SdwanID.ValueString(),
		"rtb_id":   plan.RtbID.ValueString(),
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcCreateSDWANInstanceApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp == nil {
		return fmt.Errorf("API return error. StatusCode is nil")
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	tflog.Info(ctx, "Successfully created SDWAN network instance")
	return nil
}

func (c *CtyunEcSdwanInstance) getAndMerge(ctx context.Context, state *CtyunEcSdwanInstanceConfig) (err error) {
	// 通过查询接口获取实例信息
	instanceID := state.ID.ValueString()

	// 根据InstanceID查询实例
	listReq := &ec.EcEcListSDWANInstanceRequest{
		EcID: state.EcID.ValueString(),
	}
	if !state.SdwanID.IsNull() {
		listReq.SdwanID = state.SdwanID.ValueStringPointer()
	}
	if !state.CgwID.IsNull() {
		listReq.CgwID = state.CgwID.ValueStringPointer()
	}

	listResp, err := c.meta.Apis.SdkEcApis.EcEcListSDWANInstanceApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		return
	} else if listResp == nil {
		return fmt.Errorf("API return error. StatusCode is nil")
	} else if *listResp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *listResp.Message)
		return
	} else if listResp.ReturnObj == nil || listResp.ReturnObj.Results == nil || len(listResp.ReturnObj.Results) == 0 {
		err = common.ResourceNotExistError
		return
	}
	// 检查返回的实例是否匹配
	found := false
	for _, result := range listResp.ReturnObj.Results {
		if result.InstanceID != nil && *result.InstanceID == instanceID {
			found = true
			// 更新状态值
			if result.RouteLearn != nil {
				state.RouteLearn = types.Int64Value(int64(*result.RouteLearn))
			}
			if result.RouteSync != nil {
				state.RouteSync = types.Int64Value(int64(*result.RouteSync))
			}
			if result.Weights != nil {
				state.Weights = types.Int64Value(int64(*result.Weights))
			}
			break
		}
	}

	if !found {
		return common.ResourceNotExistError
	}

	return nil
}

func (c *CtyunEcSdwanInstance) update(ctx context.Context, plan *CtyunEcSdwanInstanceConfig) (err error) {
	// 更新操作实际上是更新SDWAN实例的权重
	req := &ec.EcEcUpdateWeightsRequest{
		InstanceID: plan.ID.ValueString(),
		Weights:    int32(plan.Weights.ValueInt64()),
	}

	resp, err := c.meta.Apis.SdkEcApis.EcEcUpdateWeightsApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp == nil {
		return fmt.Errorf("API return error. StatusCode is nil")
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return nil
}

func (c *CtyunEcSdwanInstance) delete(ctx context.Context, state CtyunEcSdwanInstanceConfig) (err error) {
	instanceID := state.ID.ValueString()

	// 删除SDWAN网络实例
	req := &ec.EcEcDeleteSDWANInstanceRequest{
		InstanceID: instanceID,
	}

	tflog.Info(ctx, "删除SDWAN网络实例", map[string]interface{}{
		"instance_id": instanceID,
	})

	resp, err := c.meta.Apis.SdkEcApis.EcEcDeleteSDWANInstanceApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp == nil {
		return fmt.Errorf("API return error. StatusCode is nil")
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return nil
}

// getInstanceID 通过查询获取实例ID
func (c *CtyunEcSdwanInstance) getInstanceID(ctx context.Context, ecID, sdwanID, cgwID string) (string, error) {
	// 等待实例创建完成并获取实例ID
	// 使用更灵活的重试机制
	req := &ec.EcEcListSDWANInstanceRequest{
		EcID:    ecID,
		SdwanID: &sdwanID,
		CgwID:   &cgwID,
	}

	// 初始化重试器，最大轮询61次，间隔10秒
	retryer, err := business.NewRetryer(time.Second*10, 61)
	if err != nil {
		return "", fmt.Errorf("failed to create retryer: %w", err)
	}
	time.Sleep(time.Second * 10)
	var instanceID string
	result := retryer.Start(func(currentTime int) bool {
		resp, err := c.meta.Apis.SdkEcApis.EcEcListSDWANInstanceApi.Do(ctx, c.meta.SdkCredential, req)

		if err != nil {
			return false
		} else if resp == nil {
			return true
		} else if *resp.StatusCode != common.NormalStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return false
		} else if resp.ReturnObj == nil || resp.ReturnObj.Results == nil || len(resp.ReturnObj.Results) == 0 {
			err = fmt.Errorf("no SDWAN instance details found for SDWAN ID: %s", sdwanID)
			return true
		}

		// 查找实例ID
		for _, result := range resp.ReturnObj.Results {
			if result.InstanceID != nil && *result.SdwanID == sdwanID && *result.CgwID == cgwID {
				instanceID = *result.InstanceID
				return false // 成功，停止重试
			}
		}

		return true // 继续重试
	})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return "", errors.New("轮询已达最大次数，资源仍未绑定成功！")
	}
	return instanceID, nil
}
