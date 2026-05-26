package vpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"time"
)

var (
	_ resource.Resource              = &ctyunFlowPackageService{}
	_ resource.ResourceWithConfigure = &ctyunFlowPackageService{}
)

type ctyunFlowPackageService struct {
	meta        *common.CtyunMetadata
	name        string
	orderLooper *business.OrderLooper
}

func NewCtyunFlowPackageService() resource.Resource {
	return &ctyunFlowPackageService{}
}

func (c *ctyunFlowPackageService) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_flow_package"
	c.name = response.TypeName
}

type CtyunFlowPackageConfig struct {
	ID               types.String  `tfsdk:"id"`
	PackageName      types.String  `tfsdk:"package_name"`
	RegionID         types.String  `tfsdk:"region_id"`
	CycleType        types.String  `tfsdk:"cycle_type"`
	Spec             types.Int32   `tfsdk:"spec"`
	EffectiveTime    types.String  `tfsdk:"effective_time"`
	ExpireTime       types.String  `tfsdk:"expire_time"`
	TotalVolume      types.Float64 `tfsdk:"total_volume"`
	LeftVolume       types.Float64 `tfsdk:"left_volume"`
	MasterResourceID types.String  `tfsdk:"master_resource_id"`
}

func (c *ctyunFlowPackageService) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理共享流量包", "共享流量包（Flow Package）", "https://www.ctyun.cn/document/10026753/10032118"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "ID",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"package_name": schema.StringAttribute{
				Computed:      true,
				Description:   "套餐名",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
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
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围：month（1月），year（1年）",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypesMY...),
				},
			},
			"spec": schema.Int32Attribute{
				Required:    true,
				Description: "规格说明：当 cycleType = month 时，10-10GB,50-50GB,100-100GB,500-500GB,1024-1TB,5120-5TB,10240-10TB,51200-50TB;   **当 cycleType = year 时，120-120GB,512-512GB,8192-8TB,36864-36TB,122880-120TB,614400-600TB,1048576-1PB,2097152-2PB",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					validator2.CrossFieldInt32(
						path.MatchRoot("cycle_type"),
						[]any{"month"},
						[]int32{10, 50, 100, 500, 1024, 5120, 10240, 51200},
					),
					validator2.CrossFieldInt32(
						path.MatchRoot("cycle_type"),
						[]any{"year"},
						[]int32{
							120, 512, 8192, 36864, 122880, 614400, 1048576, 2097152,
						},
					),
				},
			},
			"effective_time": schema.StringAttribute{
				Computed:      true,
				Description:   "生效时间",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"expire_time": schema.StringAttribute{
				Computed:      true,
				Description:   "过期时间",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"total_volume": schema.Float64Attribute{
				Computed:      true,
				Description:   "总流量",
				PlanModifiers: []planmodifier.Float64{float64planmodifier.UseStateForUnknown()},
			},
			"left_volume": schema.Float64Attribute{
				Computed:    true,
				Description: "剩余流量",
			},
			"master_resource_id": schema.StringAttribute{
				Computed:      true,
				Description:   "主资源ID",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (c *ctyunFlowPackageService) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunFlowPackageConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	masterResourceID, err := c.loopCreate(ctx, &plan)
	tflog.Info(ctx, fmt.Sprintf("[DEBUG] %s: create flow masterResourceID: %s", c.name, *masterResourceID))
	if err != nil {
		return
	}
	plan.MasterResourceID = types.StringValue(*masterResourceID)

	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunFlowPackageService) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunFlowPackageConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunFlowPackageService) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (c *ctyunFlowPackageService) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunFlowPackageConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 删除
	masterOrderID, err := c.delete(ctx, state)
	if err != nil {
		return
	}

	err = c.orderLooper.WaitOrderFinish(ctx, c.meta.Credential, masterOrderID)
	return
}

func (c *ctyunFlowPackageService) delete(ctx context.Context, plan CtyunFlowPackageConfig) (masterOrderID string, err error) {
	clientToken := uuid.NewString()
	resourceID, regionID := plan.MasterResourceID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcRefundFlowPackageRequest{
		RegionID:    regionID,
		ResourceID:  resourceID,
		ClientToken: clientToken,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcRefundFlowPackageApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	masterOrderID = utils.SecString(resp.ReturnObj.MasterOrderID)
	return
}

func (c *ctyunFlowPackageService) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.orderLooper = business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
}

// loopCreate 循环执行create
func (c *ctyunFlowPackageService) loopCreate(ctx context.Context, plan *CtyunFlowPackageConfig) (masterResourceID *string, err error) {
	clientToken := uuid.NewString()
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			masterResourceID, err = c.create(ctx, clientToken, plan)
			if err != nil {
				return false
			}
			if masterResourceID == nil {
				return true
			}
			if *masterResourceID != "" {
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = errors.New("创建时未获取到主资源ID")
	}
	return
}

// create 创建
func (c *ctyunFlowPackageService) create(ctx context.Context, clientToken string, plan *CtyunFlowPackageConfig) (masterResourceID *string, err error) {
	params := &ctvpc.CtvpcBuyFlowPackageRequest{
		ClientToken: clientToken,
		RegionID:    plan.RegionID.ValueString(),
		CycleType:   plan.CycleType.ValueString(),
		CycleCount:  1,
		Spec:        plan.Spec.ValueInt32(),
		Count:       1,
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcBuyFlowPackageApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", utils.SecStringValue(resp.Message), utils.SecStringValue(resp.Description))
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	masterResourceID = resp.ReturnObj.MasterResourceID
	return
}

// getAndMerge 从远端查询
func (c *ctyunFlowPackageService) getAndMerge(ctx context.Context, plan *CtyunFlowPackageConfig) (err error) {
	flowPackageResp, err := c.loopGet(ctx, plan)

	if err != nil {
		return
	}

	plan.ID = types.StringValue(*flowPackageResp.Id)
	plan.PackageName = types.StringValue(*flowPackageResp.PackageName)
	plan.TotalVolume = types.Float64Value(flowPackageResp.TotalVolumn)
	plan.LeftVolume = types.Float64Value(flowPackageResp.LeftVolumn)
	plan.EffectiveTime = types.StringValue(*flowPackageResp.EffectiveTime)
	plan.ExpireTime = types.StringValue(*flowPackageResp.ExpireTime)

	return
}

// loopCreate 循环执行create
func (c *ctyunFlowPackageService) loopGet(ctx context.Context, plan *CtyunFlowPackageConfig) (flowPackageResp *ctvpc.CtvpcListFlowPackageReturnObjResponse, err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			flowPackageResp, err = c.get(ctx, plan)
			if err != nil {
				return false
			}
			if flowPackageResp == nil {
				err = common.ResourceNotExistError
				return true
			}
			if flowPackageResp != nil {
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = errors.New("未查询到资源信息")
	}
	return
}

func (c *ctyunFlowPackageService) get(ctx context.Context, plan *CtyunFlowPackageConfig) (flowPackageResp *ctvpc.CtvpcListFlowPackageReturnObjResponse, err error) {
	regionID := plan.RegionID.ValueString()
	params := &ctvpc.CtvpcListFlowPackageRequest{
		RegionID: regionID,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListFlowPackageApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if utils.SecString(resp.ErrorCode) == common.OpenapiFlowPackageNotFound || resp.ReturnObj == nil {
		return
	}

	for _, flowPackage := range resp.ReturnObj {
		if *flowPackage.MasterResourceBundleId == plan.MasterResourceID.ValueString() {
			flowPackageResp = flowPackage
			return
		}
	}
	return
}
