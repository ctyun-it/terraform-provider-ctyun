package elb

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

var (
	_ resource.Resource                = &ctyunElbTarget{}
	_ resource.ResourceWithConfigure   = &ctyunElbTarget{}
	_ resource.ResourceWithImportState = &ctyunElbTarget{}
)

type ctyunElbTarget struct {
	meta *common.CtyunMetadata
}

func NewCtyunElbTarget() resource.Resource {
	return &ctyunElbTarget{}
}

func (c *ctyunElbTarget) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_target"
}
func (c *ctyunElbTarget) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunElbTarget) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *ctyunElbTarget) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026756/10196689**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池Id，默认使用provider ctyun总region_id 或者环境变量",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"target_group_id": schema.StringAttribute{
				Required:    true,
				Description: "后端服务组Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "描述，支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:'{},./;'[,]·~！@#￥%……&*（） —— -+={},，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1)},
			},
			"instance_type": schema.StringAttribute{
				Required:    true,
				Description: "实例类型。取值范围：VM-虚拟云主机、BM-物理机、ECI-弹性容器",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ElbTargetInstanceType...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "云主机或物理机，或弹性容器实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"instance_ip": schema.StringAttribute{
				Optional:    true,
				Description: "后端实例ip",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.Ip(),
				},
			},
			"protocol_port": schema.Int32Attribute{
				Required:    true,
				Description: "协议端口。取值范围：1-65535，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"weight": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "后端实例权重。取值范围：1-256，默认为100，支持更新",
				Default:     int32default.StaticInt32(100),
				Validators: []validator.Int32{
					int32validator.Between(1, 256),
				},
			},
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "后端主机服务(elb_target)ID",
			},
			"health_check_status": schema.StringAttribute{
				Computed:    true,
				Description: "IPv4的健康检查状态: offline / online / unknown",
			},
			"health_check_status_ipv6": schema.StringAttribute{
				Computed:    true,
				Description: "IPv6的健康检查状态: offline / online / unknown",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "状态: DOWN / ACTIVE",
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
			},
			"updated_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间，为UTC格式",
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				// az时候有必要设定默认值
				Default: defaults.AcquireFromGlobalString(common.ExtraAzName, false),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
		},
	}
}

func (c *ctyunElbTarget) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunElbTargetConfig

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.CrateElbTarget(ctx, &plan)
	if err != nil {
		return
	}
	err = c.getAndMergeElbTarget(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunElbTarget) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunElbTargetConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 确认该rule是否异常
	if !c.acquireAndSetIdIfOrderNotFinished(ctx, &state, response) {
		return
	}
	// 查询远端
	err = c.getAndMergeElbTarget(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunElbTarget) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置
	var plan CtyunElbTargetConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunElbTargetConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
	}

	// 更新后端主机信息
	err = c.updateElbTarget(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeElbTarget(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (c *ctyunElbTarget) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunElbTargetConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	params := &ctelb.CtelbDeleteTargetRequest{
		ClientToken: uuid.NewString(),
		RegionID:    state.RegionID.ValueString(),
		TargetID:    state.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbDeleteTargetApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}
	return
}

func (c *ctyunElbTarget) CrateElbTarget(ctx context.Context, plan *CtyunElbTargetConfig) (err error) {
	if plan.RegionID.IsNull() {
		err = errors.New("创建ELB后端主机时，regionID不能为空")
		return
	}

	params := &ctelb.CtelbCreateTargetRequest{
		ClientToken:   uuid.NewString(),
		RegionID:      plan.RegionID.ValueString(),
		TargetGroupID: plan.TargetGroupID.ValueString(),
		InstanceType:  plan.InstanceType.ValueString(),
		InstanceID:    plan.InstanceID.ValueString(),
		ProtocolPort:  plan.ProtocolPort.ValueInt32(),
	}
	if !plan.InstanceIP.IsNull() {
		params.InstanceIP = plan.InstanceIP.ValueString()
	}
	if !plan.Weight.IsNull() {
		params.Weight = plan.Weight.ValueInt32()
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbCreateTargetApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}

	// 获取规则id
	if len(resp.ReturnObj) != 1 {
		err = fmt.Errorf("创建后端主机时，返回id数量有误，当前id数量为：" + strconv.Itoa(len(resp.ReturnObj)))
		return
	}
	plan.ID = types.StringValue(resp.ReturnObj[0].ID)
	return
}

func (c *ctyunElbTarget) getAndMergeElbTarget(ctx context.Context, plan *CtyunElbTargetConfig) (err error) {
	// 获取后端主机详情
	params := &ctelb.CtelbShowTargetRequest{
		RegionID: plan.RegionID.ValueString(),
		TargetID: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbShowTargetApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}
	//解析返回值
	returnObj := resp.ReturnObj

	plan.Description = types.StringValue(returnObj.Description)
	plan.ProtocolPort = types.Int32Value(returnObj.ProtocolPort)
	plan.HealthCheckStatus = types.StringValue(returnObj.HealthCheckStatus)
	plan.HealthCheckStatusIpv6 = types.StringValue(returnObj.HealthCheckStatusIpv6)
	plan.Status = types.StringValue(returnObj.Status)
	plan.CreatedTime = types.StringValue(returnObj.CreatedTime)
	plan.UpdatedTime = types.StringValue(returnObj.UpdatedTime)
	plan.Weight = types.Int32Value(returnObj.Weight)

	return
}

func (c *ctyunElbTarget) updateElbTarget(ctx context.Context, state *CtyunElbTargetConfig, plan *CtyunElbTargetConfig) (err error) {
	if state.ProtocolPort.ValueInt32() == plan.ProtocolPort.ValueInt32() && state.Weight.ValueInt32() == plan.Weight.ValueInt32() {
		return
	}
	if plan.ProtocolPort.IsNull() && plan.Weight.IsNull() {
		return
	}

	params := &ctelb.CtelbUpdateTargetRequest{
		RegionID: state.RegionID.ValueString(),
		TargetID: state.ID.ValueString(),
		Weight:   100,
	}
	if !plan.ProtocolPort.IsNull() {
		params.ProtocolPort = plan.ProtocolPort.ValueInt32()
	}
	if !plan.Weight.IsNull() {
		params.Weight = plan.Weight.ValueInt32()
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateTargetApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}
	return
}

func (c *ctyunElbTarget) acquireAndSetIdIfOrderNotFinished(ctx context.Context, state *CtyunElbTargetConfig, response *resource.ReadResponse) bool {
	if state.ID.IsNull() {
		// 该rule没有id，为非法id。移除当前状态并返回
		response.State.RemoveResource(ctx)
		return false
	}
	return true
}

type CtyunElbTargetConfig struct {
	RegionID              types.String `tfsdk:"region_id"`       //区域ID
	TargetGroupID         types.String `tfsdk:"target_group_id"` //后端服务组ID
	Description           types.String `tfsdk:"description"`     //支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:'{},./;'[,]·~！@#￥%……&*（） —— -+={},
	InstanceType          types.String `tfsdk:"instance_type"`   //实例类型。取值范围：VM、BM、ECI、IP
	InstanceID            types.String `tfsdk:"instance_id"`     //实例ID
	InstanceIP            types.String `tfsdk:"instance_ip"`     //后端服务 ip
	ProtocolPort          types.Int32  `tfsdk:"protocol_port"`   //协议端口。取值范围：1-65535
	Weight                types.Int32  `tfsdk:"weight"`          //权重。取值范围：1-256，默认为100
	ID                    types.String `tfsdk:"id"`              //后端服务组ID
	AzName                types.String `tfsdk:"az_name"`
	ProjectID             types.String `tfsdk:"project_id"`
	HealthCheckStatus     types.String `tfsdk:"health_check_status"`
	HealthCheckStatusIpv6 types.String `tfsdk:"health_check_status_ipv6"`
	Status                types.String `tfsdk:"status"`
	CreatedTime           types.String `tfsdk:"created_time"` //创建时间，为UTC格式
	UpdatedTime           types.String `tfsdk:"updated_time"` //更新时间，为UTC格式
}
