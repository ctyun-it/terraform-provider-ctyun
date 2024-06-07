package resource

import (
	"context"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	"terraform-provider-ctyun/internal/utils"
)

func NewCtyunEip() resource.Resource {
	return &ctyunEip{}
}

type ctyunEip struct {
	meta *common.CtyunMetadata
}

func (c *ctyunEip) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_eip"
}

func (c *ctyunEip) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026753**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "id",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "弹性ip名称。长度2-32，字母、数字，下划线，连字符，中文/英文字母开头，不能以http:或https:开头",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 32),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z\u4e00-\u9fa5][0-9a-zA-Z_\u4e00-\u9fa5-]+$"), "子网名称不符合规则"),
				},
			},
			"bandwidth": schema.Int64Attribute{
				Required:    true,
				Description: "原始弹性ip的带宽峰值，1-1024Mbps",
				Validators: []validator.Int64{
					int64validator.Between(1, 1024),
				},
			},
			"current_bandwidth": schema.Int64Attribute{
				Computed:    true,
				Description: "当前实际的带宽大小，如果绑定了共享带宽，此值显示为共享带宽的值，否则此值与bandwidth一致",
			},
			"bandwidth_type": schema.StringAttribute{
				Computed:    true,
				Description: "弹性ip的类型的带宽类型，standalone：独享，share：共享",
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，与cycle_count配合使用，month：按月，year：按年，on_demand：按需计费类型，当为按需时，demand_billing_type为必填。当此值为month或者year时，cycle_count为必填",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypes...),
				},
			},
			"cycle_count": schema.Int64Attribute{
				Optional:    true,
				Description: "订购时长，该参数在cycle_type为month或year时才生效，当cycleType=month，支持续订1-11个月；当cycleType=year，支持续订1-3年",
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
			"demand_billing_type": schema.StringAttribute{
				Optional:    true,
				Description: "按需计费类型，当cycle_type为on_demand时生效，bandwidth：按带宽，upflowc：按流量",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
					stringvalidator.OneOf(business.EipDemandBillingTypes...),
				},
			},
			"address": schema.StringAttribute{
				Computed:    true,
				Description: "ip地址",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "弹性ip状态，取值范围：active：有效，down：未绑定，error：出错，updating：更新中，banding_or_unbangding：绑定解绑中，deleting：删除中，deleted：已删除，expired：已过期",
			},
			"expire_time": schema.StringAttribute{
				Computed:    true,
				Description: "到期时间",
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订购的受理单id",
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

func (c *ctyunEip) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunEipConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := plan.RegionId.ValueString()
	projectId := plan.ProjectId.ValueString()
	resp, err := c.meta.Apis.CtVpcApis.EipCreateApi.Do(ctx, c.meta.Credential, &ctvpc.EipCreateRequest{
		ClientToken:       uuid.NewString(),
		RegionId:          regionId,
		CycleType:         plan.CycleType.ValueString(),
		CycleCount:        int(plan.CycleCount.ValueInt64()),
		Name:              plan.Name.ValueString(),
		Bandwidth:         int(plan.Bandwidth.ValueInt64()),
		DemandBillingType: plan.DemandBillingType.ValueString(),
		ProjectId:         projectId,
	})

	var id, masterOrderId string
	if err == nil {
		id = resp.EipId
		masterOrderId = resp.MasterOrderId
	} else {
		// 判断返回信息是否需要轮询
		if err.ErrorCode() != common.OpenapiOrderInprogress {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
		// 获取主订单
		moi, err := c.getMasterOrderIdIfOrderInProgress(err)
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
		response.Diagnostics.Append(response.State.Set(ctx, plan)...)
		// 轮询结果
		helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
		loop, err := helper.OrderLoop(ctx, c.meta.Credential, moi)
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
		id = loop.Uuid[0]
		masterOrderId = moi
	}

	plan.Id = types.StringValue(id)
	plan.RegionId = types.StringValue(regionId)
	plan.ProjectId = types.StringValue(projectId)
	plan.MasterOrderId = types.StringValue(masterOrderId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, ctyunRequestError := c.getAndMergeEip(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEip) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunEipConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	if !c.acquireAndSetIdIfOrderNotFinished(ctx, &state, response) {
		return
	}
	instance, err := c.getAndMergeEip(ctx, state)
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

func (c *ctyunEip) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state CtyunEipConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	var plan CtyunEipConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 判断名字是否相同
	if !state.Name.Equal(plan.Name) {
		_, err := c.meta.Apis.CtVpcApis.EipChangeNameApi.Do(ctx, c.meta.Credential, &ctvpc.EipChangeNameRequest{
			EipId:       state.Id.ValueString(),
			RegionId:    state.RegionId.ValueString(),
			ClientToken: uuid.NewString(),
			Name:        plan.Name.ValueString(),
		})
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
	}
	// 判断带宽是否相同，不同要走修改带宽接口
	if !state.Bandwidth.Equal(plan.Bandwidth) {
		resp, err := c.meta.Apis.CtVpcApis.EipModifySpecApi.Do(ctx, c.meta.Credential, &ctvpc.EipModifySpecRequest{
			EipId:       state.Id.ValueString(),
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

	instance, ctyunRequestError := c.getAndMergeEip(ctx, state)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEip) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunEipConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	resp, err := c.meta.Apis.CtVpcApis.EipDeleteApi.Do(ctx, c.meta.Credential, &ctvpc.EipDeleteRequest{
		EipId:       state.Id.ValueString(),
		RegionId:    state.RegionId.ValueString(),
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

func (c *ctyunEip) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// isBindWithShareBandwidth eip是否绑定共享带宽
func (c *ctyunEip) isBindWithShareBandwidth(bandwidthId string) bool {
	return bandwidthId != ""
}

// getAndMergeEip 查询eip
func (c *ctyunEip) getAndMergeEip(ctx context.Context, cfg CtyunEipConfig) (*CtyunEipConfig, error) {
	resp, err := c.meta.Apis.CtVpcApis.EipShowApi.Do(ctx, c.meta.Credential, &ctvpc.EipShowRequest{
		RegionId: cfg.RegionId.ValueString(),
		EipId:    cfg.Id.ValueString(),
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiEipNotFound {
			return nil, nil
		}
		return nil, err
	}

	// 带宽类型
	bandwidthType := business.EipBandwidthTypeStandalone
	if c.isBindWithShareBandwidth(resp.BandwidthId) {
		// 如果未共享带宽，那么忽略bandwidth的设值
		bandwidthType = business.EipBandwidthTypeShare
	}

	// 带宽大小
	bandwidth := types.Int64Value(int64(resp.Bandwidth))
	if bandwidthType == business.EipBandwidthTypeStandalone {
		// 独享带宽才更新bandwidth的值
		cfg.Bandwidth = bandwidth
	}

	statusResp, err2 := business.EipStatusMap.ToOriginalScene(resp.Status, business.EipStatusMapScene1)
	if err2 != nil {
		return nil, err2
	}

	cfg.Id = types.StringValue(resp.Id)
	cfg.Name = types.StringValue(resp.Name)
	cfg.CurrentBandwidth = bandwidth
	cfg.BandwidthType = types.StringValue(bandwidthType)
	cfg.Address = types.StringValue(resp.EipAddress)
	cfg.Status = types.StringValue(statusResp.(string))
	cfg.ExpireTime = types.StringValue(utils.FromRFC3339ToLocal(resp.ExpiredAt))
	return &cfg, nil
}

// getMasterOrderIdIfOrderInProgress 获取masterOrderId
func (c *ctyunEip) getMasterOrderIdIfOrderInProgress(err ctyunsdk.CtyunRequestError) (string, error) {
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

// acquireIdIfOrderNotFinished 重新获取id，如果前订单状态有问题需要重新轮询
// 返回值：数据是否有效
func (c *ctyunEip) acquireAndSetIdIfOrderNotFinished(ctx context.Context, state *CtyunEipConfig, response *resource.ReadResponse) bool {
	id := state.Id.ValueString()
	masterOrderId := state.MasterOrderId.ValueString()
	if id != "" {
		// 数据是完整的，无需处理
		return true
	}
	if state.MasterOrderId.ValueString() == "" {
		// 没有受理的订购单id，数据是不可恢复的，直接把当前状态移除并且返回
		response.State.RemoveResource(ctx)
		return false
	}
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	resp, err := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId)
	if err != nil || len(resp.Uuid) == 0 {
		// 报错了，或者受理没有返回数据的情况，那么意思是这个单子并没有开通出来，此时数据无法恢复
		response.State.RemoveResource(ctx)
		return false
	}

	// 成功把id恢复出来
	state.Id = types.StringValue(resp.Uuid[0])
	response.State.Set(ctx, state)
	return true
}

type CtyunEipConfig struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	CycleType         types.String `tfsdk:"cycle_type"`
	CycleCount        types.Int64  `tfsdk:"cycle_count"`
	Bandwidth         types.Int64  `tfsdk:"bandwidth"`
	CurrentBandwidth  types.Int64  `tfsdk:"current_bandwidth"`
	BandwidthType     types.String `tfsdk:"bandwidth_type"`
	DemandBillingType types.String `tfsdk:"demand_billing_type"`
	Address           types.String `tfsdk:"address"`
	Status            types.String `tfsdk:"status"`
	ExpireTime        types.String `tfsdk:"expire_time"`
	MasterOrderId     types.String `tfsdk:"master_order_id"`
	ProjectId         types.String `tfsdk:"project_id"`
	RegionId          types.String `tfsdk:"region_id"`
}
