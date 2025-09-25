package vpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctvpc2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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
	"time"
)

type ctyunBandwidth struct {
	meta *common.CtyunMetadata
}

type CtyunBandwidthConfig struct {
	CycleType  types.String `tfsdk:"cycle_type"`
	Bandwidth  types.Int32  `tfsdk:"bandwidth"`
	CycleCount types.Int64  `tfsdk:"cycle_count"`
	Name       types.String `tfsdk:"name"`
	Id         types.String `tfsdk:"id"`
	Status     types.String `tfsdk:"status"`
	ProjectId  types.String `tfsdk:"project_id"`
	RegionId   types.String `tfsdk:"region_id"`
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
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "共享带宽id",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "共享带宽命名，单账户单资源池下，命名需唯一，长度为2-63个字符，只能由数字、字母、-组成，不能以数字、-开头，且不能以-结尾，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 63),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z\\x{4e00}-\\x{9fa5}][0-9a-zA-Z_\\x{4e00}-\\x{9fa5}-]+$"), "共享带宽名称不符合规则"),
				},
			},
			"bandwidth": schema.Int32Attribute{
				Required:    true,
				Description: "共享带宽的带宽峰值（Mbit/s），取值范围5-1000，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(5, 1000),
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
				Description: "订购时长, 该参数在cycle_type为month或year时才生效，当cycle_type=month，支持订购1-11个月；当cycle_type=year，支持订购1-3年",
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
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
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
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
		},
	}
}

func (c *ctyunBandwidth) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunBandwidthConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	id, err := c.loopCreate(ctx, plan)
	if err != nil {
		return
	}
	plan.Id = types.StringValue(id)

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
			Bandwidth:   int(plan.Bandwidth.ValueInt32()),
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

// 导入命令：terraform import [配置标识].[导入配置名称] [bandwidthId],[regionId],[projectId]
func (c *ctyunBandwidth) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunBandwidthConfig
	var bandwidthId, regionId, projectId string
	err := terraform_extend.Split(request.ID, &bandwidthId, &regionId, &projectId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.Id = types.StringValue(bandwidthId)
	cfg.RegionId = types.StringValue(regionId)
	cfg.ProjectId = types.StringValue(projectId)

	instance, err := c.getAndMergeBandwidth(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
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
	cfg.Bandwidth = types.Int32Value(int32(resp.Bandwidth))
	cfg.Name = types.StringValue(resp.Name)
	return &cfg, nil
}

// loopCreate 循环执行create
func (c *ctyunBandwidth) loopCreate(ctx context.Context, plan CtyunBandwidthConfig) (id string, err error) {
	clientToken := uuid.NewString()
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			id, err = c.create(ctx, clientToken, plan)
			if err != nil {
				return false
			}
			if id != "" {
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = errors.New("创建时未获取到共享带宽id")
	}
	return
}

// create 创建共享带宽
func (c *ctyunBandwidth) create(ctx context.Context, clientToken string, plan CtyunBandwidthConfig) (id string, err error) {
	params := &ctvpc2.CtvpcCreateBandwidthRequest{
		ClientToken: clientToken,
		RegionID:    plan.RegionId.ValueString(),
		CycleType:   plan.CycleType.ValueString(),
		CycleCount:  int32(plan.CycleCount.ValueInt64()),
		Name:        plan.Name.ValueString(),
		Bandwidth:   plan.Bandwidth.ValueInt32(),
		ProjectID:   plan.ProjectId.ValueStringPointer(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateBandwidthApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	id = utils.SecString(resp.ReturnObj.BandwidthID)
	return
}
