package resource

import (
	"context"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
	defaults2 "terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
	"time"
)

type ctyunBandwidth struct {
	meta *common.CtyunMetadata
}

func NewCtyunBandwidth() resource.Resource {
	return &ctyunBandwidth{}
}

func (c *ctyunBandwidth) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_bandwidth"
}

func (c *ctyunBandwidth) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026761**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "共享带宽id",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "共享带宽命名，单账户单资源池下，命名需唯一，长度为2-63个字符，只能由数字、字母、-组成，不能以数字、-开头，且不能以-结尾",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 63),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z\u4e00-\u9fa5][0-9a-zA-Z_\u4e00-\u9fa5-]+$"), "共享带宽名称不符合规则"),
				},
			},
			"bandwidth": schema.Int64Attribute{
				Required:    true,
				Description: "共享带宽的带宽峰值（Mbit/s），必须大于等于5",
				Validators: []validator.Int64{
					int64validator.AtLeast(5),
				},
			},
			"cycle_type": schema.StringAttribute{
				Optional:    true,
				Description: "订购周期类型，取值范围：month：按月，year：按年、on_demand：按需。当此值为month或者year时，cycle_count为必填",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypes...),
				},
			},
			"cycle_count": schema.Int64Attribute{
				Optional:    true,
				Description: "订购时长, 该参数在cycle_type为month或year时才生效，当cycleType=month，支持续订1-11个月；当cycleType=year，支持续订1-3年",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
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
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "共享带宽状态: active：有效，expired：已过期，freezing：冻结",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目id，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
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

func (c *ctyunBandwidth) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunBandwidthConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := plan.RegionId.ValueString()
	projectId := plan.ProjectId.ValueString()
	createRequest := ctvpc.BandwidthCreateRequest{
		RegionId:    regionId,
		ProjectId:   projectId,
		ClientToken: uuid.NewString(),
		CycleType:   plan.CycleType.ValueString(),
		Bandwidth:   plan.Bandwidth.ValueInt64(),
		Name:        plan.Name.ValueString(),
		CycleCount:  plan.CycleCount.ValueInt64(),
	}

	var resp *ctvpc.BandwidthCreateResponse
	var err error
	retryer, _ := business.NewRetryer(time.Second*5, 20)
	retryerResult := retryer.Start(
		func(currentTime int) bool {
			do, requestError := c.meta.Apis.CtVpcApis.BandwidthCreateApi.Do(ctx, c.meta.Credential, &createRequest)
			if requestError != nil {
				// 接口如果出现异常，且不等于配置的重复轮询的异常码时 则退出轮询
				if requestError.ErrorCode() != common.OpenapiOrderInprogress {
					err = requestError
					return false
				}
				return true
			}
			resp = do
			return false
		},
	)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	if retryerResult.ReturnReason == business.ReachMaxLoopTime {
		msg := "创建带宽下单已达到最大的轮询次数"
		response.Diagnostics.Append(diag.NewErrorDiagnostic(msg, msg))
		return
	}

	plan.Id = types.StringValue(resp.BandwidthId)
	plan.RegionId = types.StringValue(regionId)
	plan.ProjectId = types.StringValue(projectId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, ctyunRequestError := c.getAndMergeBandwidth(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunBandwidth) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunBandwidthConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeBandwidth(ctx, state)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunBandwidth) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan CtyunBandwidthConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var state CtyunBandwidthConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	// 判断名字是否相同
	if !plan.Name.Equal(state.Name) {
		_, err := c.meta.Apis.CtVpcApis.BandwidthChangeNameApi.Do(ctx, c.meta.Credential, &ctvpc.BandwidthChangeNameRequest{
			BandwidthId: state.Id.ValueString(),
			RegionId:    state.RegionId.ValueString(),
			ClientToken: uuid.NewString(),
			Name:        plan.Name.ValueString(),
		})
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
	}
	// 判断带宽大小是否相同，不同要走修改带宽接口
	if !plan.Bandwidth.Equal(state.Bandwidth) {
		resp, err := c.meta.Apis.CtVpcApis.BandwidthChangeSpecApi.Do(ctx, c.meta.Credential, &ctvpc.BandwidthChangeSpecRequest{
			BandwidthId: state.Id.ValueString(),
			RegionId:    state.RegionId.ValueString(),
			ClientToken: uuid.NewString(),
			Bandwidth:   int(plan.Bandwidth.ValueInt64()),
		})

		var masterOrderId string
		if err == nil {
			masterOrderId = resp.MasterOrderId
		} else {
			if err.ErrorCode() != common.OpenapiOrderInprogress {
				response.Diagnostics.AddError(err.Error(), err.Error())
				return
			}
			id, err := c.getMasterOrderIdIfOrderInProgress(err)
			if err != nil {
				response.Diagnostics.AddError(err.Error(), err.Error())
				return
			}
			masterOrderId = id
		}

		helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
		_, err2 := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId)
		if err2 != nil {
			response.Diagnostics.AddError(err2.Error(), err2.Error())
			return
		}
	}

	instance, ctyunRequestError := c.getAndMergeBandwidth(ctx, state)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunBandwidth) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunBandwidthConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	resp, err := c.meta.Apis.CtVpcApis.BandwidthDeleteApi.Do(ctx, c.meta.Credential, &ctvpc.BandwidthDeleteRequest{
		BandwidthId: state.Id.ValueString(),
		RegionId:    state.RegionId.ValueString(),
		ProjectId:   state.ProjectId.ValueString(),
		ClientToken: uuid.NewString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	err2 := helper.RefundLoop(ctx, c.meta.Credential, resp.MasterOrderId)
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}
}

func (c *ctyunBandwidth) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// getMasterOrderIdIfOrderInProgress 获取masterOrderId
func (c *ctyunBandwidth) getMasterOrderIdIfOrderInProgress(err ctyunsdk.CtyunRequestError) (string, error) {
	resp := struct {
		MasterOrderId string `json:"masterOrderID"`
		MasterOrderNo string `json:"masterOrderNO"`
	}{}
	if err.CtyunResponse() == nil {
		return "", err
	}
	_, err = err.CtyunResponse().ParseByStandardModel(&resp)
	if err != nil {
		return "", err
	}
	return resp.MasterOrderId, err
}

// getAndMergeBandwidth 查询带宽
func (c *ctyunBandwidth) getAndMergeBandwidth(ctx context.Context, cfg CtyunBandwidthConfig) (*CtyunBandwidthConfig, error) {
	resp, err := c.meta.Apis.CtVpcApis.BandwidthDescribeApi.Do(ctx, c.meta.Credential, &ctvpc.BandwidthDescribeRequest{
		RegionId:    cfg.RegionId.ValueString(),
		BandwidthId: cfg.Id.ValueString(),
		ProjectId:   cfg.ProjectId.ValueString(),
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiSharedbandwidthNotFound {
			return nil, nil
		}
		return nil, err
	}

	statusResp, err2 := business.BandwidthStatusMap.ToOriginalScene(resp.Status, business.BandwidthStatusMapScene1)
	if err2 != nil {
		return nil, err2
	}
	cfg.Status = types.StringValue(statusResp.(string))
	cfg.Bandwidth = types.Int64Value(int64(resp.Bandwidth))
	cfg.Name = types.StringValue(resp.Name)
	return &cfg, nil
}

type CtyunBandwidthConfig struct {
	CycleType  types.String `tfsdk:"cycle_type"`
	Bandwidth  types.Int64  `tfsdk:"bandwidth"`
	CycleCount types.Int64  `tfsdk:"cycle_count"`
	Name       types.String `tfsdk:"name"`
	Id         types.String `tfsdk:"id"`
	Status     types.String `tfsdk:"status"`
	ProjectId  types.String `tfsdk:"project_id"`
	RegionId   types.String `tfsdk:"region_id"`
}
