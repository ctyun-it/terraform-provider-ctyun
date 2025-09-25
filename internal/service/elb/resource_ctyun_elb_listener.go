package elb

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
)

var (
	_ resource.Resource                = &CtyunElbListener{}
	_ resource.ResourceWithConfigure   = &CtyunElbListener{}
	_ resource.ResourceWithImportState = &CtyunElbListener{}
)

type CtyunElbListener struct {
	meta *common.CtyunMetadata
}

func NewCtyunElbListener() resource.Resource {
	return &CtyunElbListener{}
}

func (c *CtyunElbListener) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunElbListener) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_listener"
}

func (c *CtyunElbListener) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026756/10140276**`,
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
			"loadbalancer_id": schema.StringAttribute{
				Required:    true,
				Description: "负载均衡实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 32),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z\\x{4e00}-\\x{9fa5}][a-zA-Z0-9_\\-\\x{4e00}-\\x{9fa5}]*$"), "必须以拉丁字母或中文开头，只能包含拉丁字母、中文、数字、下划线和连字符"),
					stringvalidator.RegexMatches(regexp.MustCompile(`^([^h]|h[^t]|ht[^t]|htt[^p]|http[^s]|https.).*$`), "不能以http:或https:开头"),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·！@#￥%……&*（） —— -+={}\\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 128),
					validator2.Desc(),
					stringvalidator.RegexMatches(regexp.MustCompile(`^([^h]|h[^t]|ht[^t]|htt[^p]|http[^s]|https.).*$`), "不能以http:或https:开头"),
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "监听协议。取值范围：TCP、UDP、HTTP、HTTPS",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ListenerProtocols...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol_port": schema.Int32Attribute{
				Required:    true,
				Description: "负载均衡实例监听端口。取值：1-65535，protocol_port。不支持更新",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"certificate_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "证书ID。当protocol为HTTPS时，此参数必填，支持更新",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.ListenerProtocolHTTPS),
					),
				},
			},
			"ca_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否开启双向认证。true（开启），false（不开启），支持更新",
			},
			"client_certificate_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "双向认证的证书ID，当ca_enabled=ture，必填。支持更新",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("ca_enabled"),
						types.BoolValue(true),
					),
				},
			},
			"access_control_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "访问控制ID，当access_control_type=white或者black，必填。支持更新",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("access_control_type"),
						types.StringValue(business.ListenerAccessControlTypeWhite),
						types.StringValue(business.ListenerAccessControlTypeBlack),
					),
				},
			},
			"access_control_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "访问控制类型。取值范围：Close（未启用）、White（白名单）、Black（黑名单），支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ListenerAccessControlTypes...),
				},
			},
			"forwarded_for_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "x-forward-for功能。false（未开启）、true（开启），支持更新",
			},
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "监听器ID",
			},
			"default_action_type": schema.StringAttribute{
				Required:    true,
				Description: "默认规则动作类型。取值范围：forward、redirect，支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ListenerDefaultActionTypes...),
				},
			},
			"redirect_listener_id": schema.StringAttribute{
				Optional:    true,
				Description: "重定向监听器ID，当default_action_type为redirect时，此字段必填。支持更新",
				Validators: []validator.String{
					validator2.ConflictsWithEqualString(
						path.MatchRoot("default_action_type"),
						types.StringValue(business.ListenerDefaultActionTypeForward),
					),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("default_action_type"),
						types.StringValue(business.ListenerDefaultActionTypeRedirect),
					),
				},
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称",
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
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "监听器状态: DOWN/ACTIVE，可以控制监听器开关。支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ElbRuleStatus...),
				},
				Default: stringdefault.StaticString(business.ElbRuleStatusACTIVE),
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
			},
			"updated_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间，为UTC格式",
			},
			"enable_nat_64": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否开启nat64，elb需要支持ipv6能力，支持更新",
			},
			"listener_qps": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "qps 大小，仅支持协议为 HTTP / HTTPS，的监听器，支持更新",
				Validators: []validator.Int32{
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("protocol"),
						types.StringValue(business.ListenerProtocolTCP),
						types.StringValue(business.ListenerProtocolUDP),
					),
				},
			},
			"establish_timeout": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "建立连接超时时间，单位秒，取值范围：1 - 1800。不支持协议为 UDP / HTTP / HTTPS 的监听器，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(1, 1800),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("protocol"),
						types.StringValue(business.ListenerProtocolHTTP),
						types.StringValue(business.ListenerProtocolHTTPS),
						types.StringValue(business.ListenerProtocolUDP),
					),
				},
			},
			"idle_timeout": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "链接空闲断开超时时间，单位秒，取值范围：1 - 300,不支持协议为 TCP / UDP 的监听器，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(1, 300),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("protocol"),
						types.StringValue(business.ListenerProtocolTCP),
						types.StringValue(business.ListenerProtocolUDP),
					),
				},
			},
			"response_timeout": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "响应超时，单位秒，取值范围：1 - 300。不支持协议为 TCP / UDP 的监听器，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(1, 300),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("protocol"),
						types.StringValue(business.ListenerProtocolTCP),
						types.StringValue(business.ListenerProtocolUDP),
					),
				},
			},
			"listener_cps": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "cps大小，仅支持协议为 TCP / UDP 的监听器。支持更新",
				Validators: []validator.Int32{
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("protocol"),
						types.StringValue(business.ListenerProtocolHTTP),
						types.StringValue(business.ListenerProtocolHTTPS),
					),
				},
			},
			"target_groups": schema.ListNestedAttribute{
				Optional:    true,
				Description: "后端服务组，最多只支持添加一个后端服务组。当default_action_type=forward时，target_groups不能为空。支持更新",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"target_group_id": schema.StringAttribute{
							Required:    true,
							Description: "后端服务组ID，支持更新",
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
						"weight": schema.Int32Attribute{
							Optional:    true,
							Computed:    true,
							Default:     int32default.StaticInt32(100),
							Description: "后端主机权重，取值范围：1-256。默认为100，支持更新",
							Validators: []validator.Int32{
								int32validator.Between(1, 256),
							},
						},
					},
				},
				Validators: []validator.List{
					validator2.AlsoRequiresEqualList(
						path.MatchRoot("default_action_type"),
						types.StringValue(business.ListenerDefaultActionTypeForward),
					),
					listvalidator.SizeAtMost(1),
				},
			},
		},
	}
}

func (c CtyunElbListener) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunElbListenerConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	//err = c.CheckBeforeCreateElbListener(ctx, plan)
	//if err != nil {
	//	return
	//}

	// 创建
	err = c.CreateElbListener(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	err = c.setNat64(ctx, &plan, &plan)
	if err != nil {
		return
	}
	//设置QPS,仅支持协议为 HTTP / HTTPS 的监听器
	err = c.setQPS(ctx, &plan, &plan, false)
	if err != nil {
		return
	}
	// 设置CPS，仅支持协议为 TCP / UDP 的监听器。
	err = c.setCPS(ctx, &plan, &plan, false)
	if err != nil {
		return
	}
	//设置establish_timeout,不支持协议为 UDP / HTTP / HTTPS 的监听器
	err = c.setEstablishTimeout(ctx, &plan, &plan, false)
	if err != nil {
		return
	}
	//设置idle_timeout,不支持协议为 TCP / UDP 的监听器
	err = c.setIdleTimeout(ctx, &plan, &plan, false)
	if err != nil {
		return
	}
	//设置response_timeout,不支持协议为 TCP / UDP 的监听器
	err = c.setResponseTimeout(ctx, &plan, &plan, false)
	if err != nil {
		return
	}
	if plan.Status.ValueString() == business.ElbRuleStatusDOWN {
		err = c.stopListener(ctx, plan)
		if err != nil {
			return
		}
	}
	// 创建后反查创建后的nat信息
	err = c.getAndMergeListener(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunElbListener) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunElbListenerConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeListener(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunElbListener) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置
	var plan CtyunElbListenerConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunElbListenerConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 更新负载均衡监听器信息
	err = c.updateListenerInfo(ctx, &state, &plan)
	if err != nil {
		return
	}
	//停止监听器,若plan.status=DOWN，并且state.status=ACTIVE则触发停止监听器
	if plan.Status.ValueString() == business.ElbRuleStatusDOWN && state.Status.ValueString() == business.ElbRuleStatusACTIVE {
		err = c.stopListener(ctx, plan)
		if err != nil {
			return
		}
	}
	//启动监听器
	if plan.Status.ValueString() == business.ElbRuleStatusACTIVE && state.Status.ValueString() == business.ElbRuleStatusDOWN {
		err = c.startListener(ctx, state)
		if err != nil {
			return
		}
	}
	//设置NAT64,仅支持开启 IPv6 的负载均衡。
	err = c.setNat64(ctx, &state, &plan)
	if err != nil {
		return
	}
	//设置QPS,仅支持协议为 HTTP / HTTPS 的监听器
	err = c.setQPS(ctx, &state, &plan, true)
	if err != nil {
		return
	}
	// 设置CPS，仅支持协议为 TCP / UDP 的监听器。
	err = c.setCPS(ctx, &state, &plan, true)
	if err != nil {
		return
	}
	//设置establish_timeout,不支持协议为 UDP / HTTP / HTTPS 的监听器
	err = c.setEstablishTimeout(ctx, &state, &plan, true)
	if err != nil {
		return
	}
	//设置idle_timeout,不支持协议为 TCP / UDP 的监听器
	err = c.setIdleTimeout(ctx, &state, &plan, true)
	if err != nil {
		return
	}
	//设置response_timeout,不支持协议为 TCP / UDP 的监听器
	err = c.setResponseTimeout(ctx, &state, &plan, true)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeListener(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunElbListener) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunElbListenerConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &ctelb.CtelbDeleteListenerRequest{
		ClientToken: uuid.NewString(),
		RegionID:    state.RegionID.ValueString(),
		ListenerID:  state.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbDeleteListenerApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

}

func (c *CtyunElbListener) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunElbListenerConfig
	var id string
	err = terraform_extend.Split(request.ID, &id)
	if err != nil {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(cfg.RegionID.ValueString(), common.ExtraRegionId)
	cfg.RegionID = types.StringValue(regionId)
	azName := c.meta.GetExtraIfEmpty(cfg.AzName.ValueString(), common.ExtraAzName)
	cfg.AzName = types.StringValue(azName)

	cfg.ID = types.StringValue(id)
	err = c.getAndMergeListener(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &cfg)...)
}

func (c *CtyunElbListener) CreateElbListener(ctx context.Context, plan *CtyunElbListenerConfig) (err error) {

	params := &ctelb.CtelbCreateListenerRequest{
		ClientToken:    uuid.NewString(),
		RegionID:       plan.RegionID.ValueString(),
		LoadBalancerID: plan.LoadBalancerID.ValueString(),
		Name:           plan.Name.ValueString(),
		Protocol:       plan.Protocol.ValueString(),
		ProtocolPort:   plan.ProtocolPort.ValueInt32(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		params.Description = plan.Description.ValueString()
	}
	if !plan.CertificateID.IsNull() && !plan.CertificateID.IsUnknown() {
		params.CertificateID = plan.CertificateID.ValueString()
	}
	if !plan.CaEnabled.IsNull() && !plan.CaEnabled.IsUnknown() {
		params.CaEnabled = plan.CaEnabled.ValueBoolPointer()
	}
	if !plan.ClientCertificateID.IsNull() && !plan.ClientCertificateID.IsUnknown() {
		params.ClientCertificateID = plan.ClientCertificateID.ValueString()
	}
	if !plan.AccessControlID.IsNull() && !plan.AccessControlID.IsUnknown() {
		params.AccessControlID = plan.AccessControlID.ValueString()
	}
	if !plan.AccessControlType.IsNull() && !plan.AccessControlType.IsUnknown() {
		params.AccessControlType = plan.AccessControlType.ValueString()
	}
	if !plan.ForwardedForEnabled.IsNull() && !plan.ForwardedForEnabled.IsUnknown() {
		params.ForwardedForEnabled = plan.ForwardedForEnabled.ValueBoolPointer()
	}

	// 处理defaultAction
	var defaultAction ctelb.CtelbCreateListenerDefaultActionRequest
	defaultAction.RawType = plan.DefaultActionType.ValueString()
	if plan.DefaultActionType.ValueString() == business.ListenerDefaultActionTypeRedirect && plan.RedirectListenerID.IsNull() {
		err = fmt.Errorf("创建负载均衡监听器时，若默认规则类型=redirect时，重定向监听器ID不能为空")
	}
	defaultAction.RedirectListenerID = plan.RedirectListenerID.ValueString()
	defaultAction.ForwardConfig = &ctelb.CtelbCreateListenerDefaultActionForwardConfigRequest{}
	var targetGroupList []TargetGroupModel
	// 判断default_action_type为不同值的时候，各个参数的传参情况
	if plan.DefaultActionType.ValueString() == business.ListenerDefaultActionTypeForward {
		if plan.TargetGroups.IsNull() || plan.TargetGroups.IsUnknown() {
			err = fmt.Errorf("当default_action_type=forward时，target_groups不能为空")
			return err
		}
	} else if plan.DefaultActionType.ValueString() == business.ListenerDefaultActionTypeRedirect {
		if plan.RedirectListenerID.IsNull() {
			err = fmt.Errorf("当default_action_type=redirect时，redirect_listener_id不能为空")
			return err
		}
	}
	var targetGroups []*ctelb.CtelbCreateListenerDefaultActionForwardConfigTargetGroupsRequest
	diags := plan.TargetGroups.ElementsAs(ctx, &targetGroupList, false)
	if diags.HasError() {
		return
	}
	for _, targetGroupItem := range targetGroupList {
		var targetGroup ctelb.CtelbCreateListenerDefaultActionForwardConfigTargetGroupsRequest
		if targetGroupItem.TargetGroupID.IsNull() {
			err = errors.New("创建转发规则时，targetGroupID不能为空")
			return
		}
		targetGroup.TargetGroupID = targetGroupItem.TargetGroupID.ValueString()
		if !targetGroupItem.Weight.IsNull() {
			targetGroup.Weight = targetGroupItem.Weight.ValueInt32()
		}
		targetGroups = append(targetGroups, &targetGroup)
	}
	defaultAction.ForwardConfig.TargetGroups = targetGroups
	params.DefaultAction = &defaultAction
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbCreateListenerApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 保存listener id
	if len(resp.ReturnObj) != 1 {
		err = fmt.Errorf("返回的监听器ID数量有误！")
		return
	}
	plan.ID = types.StringValue(resp.ReturnObj[0].ID)
	return nil
}

func (c *CtyunElbListener) updateListenerInfo(ctx context.Context, state *CtyunElbListenerConfig, plan *CtyunElbListenerConfig) (err error) {
	params := &ctelb.CtelbUpdateListenerRequest{
		ClientToken: uuid.NewString(),
		RegionID:    state.RegionID.ValueString(),
		ListenerID:  state.ID.ValueString(),
	}
	if !plan.Name.Equal(state.Name) {
		params.Name = plan.Name.ValueString()
	}
	if !plan.Description.IsNull() && !plan.Description.Equal(state.Description) {
		params.Description = plan.Description.ValueString()
	}
	if !plan.CertificateID.IsNull() && !plan.CertificateID.Equal(state.CertificateID) {
		params.CertificateID = plan.CertificateID.ValueString()
	}
	if !plan.CaEnabled.IsNull() && !plan.CaEnabled.Equal(state.CaEnabled) {
		params.CaEnabled = plan.CaEnabled.ValueBoolPointer()
	}
	if !plan.ClientCertificateID.IsNull() && !plan.ClientCertificateID.Equal(state.ClientCertificateID) {
		params.ClientCertificateID = plan.ClientCertificateID.ValueString()
	}
	if !plan.AccessControlID.IsNull() && !plan.AccessControlID.Equal(state.AccessControlID) {
		params.AccessControlID = plan.AccessControlID.ValueString()
	}
	if !plan.AccessControlType.IsNull() && !plan.AccessControlType.Equal(state.AccessControlType) {
		params.AccessControlType = plan.AccessControlType.ValueString()
	}
	if !plan.ForwardedForEnabled.IsNull() && !plan.ForwardedForEnabled.Equal(state.ForwardedForEnabled) {
		params.ForwardedForEnabled = plan.ForwardedForEnabled.ValueBoolPointer()
	}

	if plan.DefaultActionType.ValueString() != "" {
		// 处理defaultAction
		defaultAction := &ctelb.CtelbUpdateListenerDefaultActionRequest{}
		defaultAction.RawType = plan.DefaultActionType.ValueString()
		if plan.DefaultActionType.ValueString() == business.ListenerDefaultActionTypeRedirect {
			if plan.RedirectListenerID.ValueString() == "" {
				err = fmt.Errorf("当DefaultActionType=redirect时，redirectListenerID不能为空")
				return
			}
			defaultAction.RedirectListenerID = plan.RedirectListenerID.ValueString()
		} else if plan.DefaultActionType.ValueString() == business.ListenerDefaultActionTypeForward {
			defaultAction.ForwardConfig = &ctelb.CtelbUpdateListenerDefaultActionForwardConfigRequest{}
			if plan.TargetGroups.IsNull() || plan.TargetGroups.IsUnknown() {
				err = fmt.Errorf("当DefaultActionType=forward时，targetGroups不能为空")
				return
			}
			// 后端服务组
			var targetGroupList []TargetGroupModel
			var targetGroups []*ctelb.CtelbUpdateListenerDefaultActionForwardConfigTargetGroupsRequest
			diags := plan.TargetGroups.ElementsAs(ctx, &targetGroupList, false)
			if diags.HasError() {
				return
			}
			for _, targetGroupItem := range targetGroupList {
				var targetGroup ctelb.CtelbUpdateListenerDefaultActionForwardConfigTargetGroupsRequest
				if targetGroupItem.TargetGroupID.IsNull() {
					err = errors.New("targetGroupID不能为空")
					return
				}
				targetGroup.TargetGroupID = targetGroupItem.TargetGroupID.ValueString()
				if targetGroupItem.Weight.ValueInt32() != 0 {
					targetGroup.Weight = targetGroupItem.Weight.ValueInt32()
				}
				targetGroups = append(targetGroups, &targetGroup)
			}
			defaultAction.ForwardConfig.TargetGroups = targetGroups
		}
		params.DefaultAction = defaultAction
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateListenerApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

func (c *CtyunElbListener) getAndMergeListener(ctx context.Context, plan *CtyunElbListenerConfig) (err error) {
	params := &ctelb.CtelbShowListenerRequest{
		RegionID:   plan.RegionID.ValueString(),
		ListenerID: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbShowListenerApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回listener详情
	respObj := resp.ReturnObj[0]
	plan.Description = types.StringValue(respObj.Description)
	plan.Status = types.StringValue(respObj.Status)
	plan.CreatedTime = types.StringValue(respObj.CreatedTime)
	plan.UpdatedTime = types.StringValue(respObj.UpdatedTime)
	plan.ListenerQps = types.Int32Value(respObj.Qps)
	plan.ResponseTimeout = types.Int32Value(respObj.ResponseTimeout)
	plan.EstablishTimeout = types.Int32Value(respObj.EstablishTimeout)
	plan.IdleTimeout = types.Int32Value(respObj.IdleTimeout)
	plan.ListenerCps = types.Int32Value(respObj.Cps)
	plan.Name = types.StringValue(respObj.Name)
	plan.CertificateID = types.StringValue(respObj.CertificateID)
	plan.ClientCertificateID = types.StringValue(respObj.ClientCertificateID)
	plan.AccessControlID = types.StringValue(respObj.AccessControlID)
	plan.CaEnabled = utils.SecBoolValue(respObj.CaEnabled)
	plan.ForwardedForEnabled = utils.SecBoolValue(respObj.ForwardedForEnabled)
	plan.AccessControlType = types.StringValue(respObj.AccessControlType)
	if respObj.Nat64 == 0 {
		plan.EnableNat64 = types.BoolValue(false)
	} else if respObj.Nat64 == 1 {
		plan.EnableNat64 = types.BoolValue(true)
	}

	// 更新defaultAction
	plan.DefaultActionType = types.StringValue(respObj.DefaultAction.RawType)
	targetGroupList := respObj.DefaultAction.ForwardConfig.TargetGroups
	var targetGroups []TargetGroupsModel
	for _, targetGroupItem := range targetGroupList {
		var targetGroup TargetGroupsModel
		targetGroup.TargetGroupID = types.StringValue(targetGroupItem.TargetGroupID)
		targetGroup.Weight = types.Int32Value(targetGroupItem.Weight)
		targetGroups = append(targetGroups, targetGroup)
	}
	var diags diag.Diagnostics
	plan.TargetGroups, diags = types.ListValueFrom(ctx, utils.StructToTFObjectTypes(TargetGroupsModel{}), targetGroups)
	if diags.HasError() {
		return
	}
	return
}

// 启动监听器
func (c *CtyunElbListener) startListener(ctx context.Context, state CtyunElbListenerConfig) (err error) {
	params := &ctelb.CtelbStartListenerRequest{
		ClientToken: uuid.NewString(),
		RegionID:    state.RegionID.ValueString(),
		ListenerID:  state.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbStartListenerApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

func (c *CtyunElbListener) stopListener(ctx context.Context, state CtyunElbListenerConfig) (err error) {
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbStopListenerApi.Do(ctx, c.meta.SdkCredential, &ctelb.CtelbStopListenerRequest{
		ClientToken: uuid.NewString(),
		RegionID:    state.RegionID.ValueString(),
		ListenerID:  state.ID.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// 设置监听器 Nat64
func (c *CtyunElbListener) setNat64(ctx context.Context, state *CtyunElbListenerConfig, plan *CtyunElbListenerConfig) (err error) {
	// 验证负载均衡是否开启ipv6
	if !c.enableElbIpv6(ctx, state) {
		if plan.EnableNat64.ValueBool() {
			err = errors.New("elb未开启ipv6，不支持nat64")
		}
		return
	}
	params := &ctelb.CtelbUpdateListenerNat64Request{
		RegionID:    state.RegionID.ValueString(),
		ListenerID:  state.ID.ValueString(),
		EnableNat64: plan.EnableNat64.ValueBool(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateListenerNat64Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

func (c *CtyunElbListener) setQPS(ctx context.Context, state *CtyunElbListenerConfig, plan *CtyunElbListenerConfig, isUpdate bool) (err error) {
	// 仅支持协议为 HTTP / HTTPS 的监听器
	if state.Protocol.ValueString() == business.ListenerProtocolUDP || state.Protocol.ValueString() == business.ListenerProtocolTCP {
		return
	}
	if plan.ListenerQps.IsNull() || plan.ListenerQps.IsUnknown() {
		return
	}
	// 如何plan qps 和state qps相同，直接返回
	if isUpdate && plan.ListenerQps.ValueInt32() == state.ListenerQps.ValueInt32() {
		return
	}
	params := &ctelb.CtelbUpdateListenerQpsRequest{
		RegionID:    state.RegionID.ValueString(),
		ListenerID:  state.ID.ValueString(),
		ListenerQps: plan.ListenerQps.ValueInt32(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateListenerQpsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

func (c *CtyunElbListener) setCPS(ctx context.Context, state *CtyunElbListenerConfig, plan *CtyunElbListenerConfig, isUpdate bool) (err error) {
	// 仅支持协议为 TCP / UDP 的监听器
	if state.Protocol.ValueString() != business.ListenerProtocolUDP && state.Protocol.ValueString() != business.ListenerProtocolTCP {
		return
	}
	if plan.ListenerCps.IsNull() || plan.ListenerCps.IsUnknown() {
		return
	}
	// 如果plan cps与state cps相同，直接返回
	if isUpdate && plan.ListenerCps.ValueInt32() == state.ListenerCps.ValueInt32() {
		return
	}
	params := &ctelb.CtelbUpdateListenerCpsRequest{
		RegionID:    state.RegionID.ValueString(),
		ListenerID:  state.ID.ValueString(),
		ListenerCps: plan.ListenerCps.ValueInt32(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateListenerCpsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

func (c *CtyunElbListener) setEstablishTimeout(ctx context.Context, state *CtyunElbListenerConfig, plan *CtyunElbListenerConfig, isUpdate bool) (err error) {
	//不支持协议为 UDP / HTTP / HTTPS 的监听器
	if state.Protocol.ValueString() != business.ListenerProtocolTCP {
		return
	}

	if plan.EstablishTimeout.IsNull() || plan.EstablishTimeout.IsUnknown() {
		return
	}
	if isUpdate && plan.EstablishTimeout.ValueInt32() == state.EstablishTimeout.ValueInt32() {
		return
	}
	params := &ctelb.CtelbUpdateListenerEstabTimeoutRequest{
		RegionID:         state.RegionID.ValueString(),
		ListenerID:       state.ID.ValueString(),
		EstablishTimeout: plan.EstablishTimeout.ValueInt32(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateListenerEstabTimeoutApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

func (c *CtyunElbListener) setIdleTimeout(ctx context.Context, state *CtyunElbListenerConfig, plan *CtyunElbListenerConfig, isUpdate bool) (err error) {
	// 不支持协议为 TCP / UDP 的监听器
	if state.Protocol.ValueString() == business.ListenerProtocolTCP || state.Protocol.ValueString() == business.ListenerProtocolUDP {
		return
	}
	if plan.IdleTimeout.IsNull() || plan.IdleTimeout.IsUnknown() {
		return
	}
	if isUpdate && plan.IdleTimeout.ValueInt32() == state.IdleTimeout.ValueInt32() {
		return
	}

	params := &ctelb.CtelbUpdateListenerIdleTimeoutRequest{
		RegionID:    state.RegionID.ValueString(),
		ListenerID:  state.ID.ValueString(),
		IdleTimeout: plan.IdleTimeout.ValueInt32(),
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateListenerIdleTimeoutApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

func (c *CtyunElbListener) setResponseTimeout(ctx context.Context, state *CtyunElbListenerConfig, plan *CtyunElbListenerConfig, isUpdate bool) (err error) {
	// 不支持协议为 TCP / UDP 的监听器
	if state.Protocol.ValueString() == business.ListenerProtocolTCP || state.Protocol.ValueString() == business.ListenerProtocolUDP {
		return
	}
	if plan.ResponseTimeout.IsNull() || plan.ResponseTimeout.IsUnknown() {
		return
	}
	if isUpdate && plan.ResponseTimeout.ValueInt32() == state.ResponseTimeout.ValueInt32() {
		return
	}

	params := &ctelb.CtelbUpdateListenerResponseTimeoutRequest{
		RegionID:        state.RegionID.ValueString(),
		ListenerID:      state.ID.ValueString(),
		ResponseTimeout: plan.ResponseTimeout.ValueInt32(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateListenerResponseTimeoutApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

func (c *CtyunElbListener) enableElbIpv6(ctx context.Context, state *CtyunElbListenerConfig) bool {
	params := &ctelb.CtelbShowLoadBalancerRequest{
		RegionID: state.RegionID.ValueString(),
		ElbID:    state.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbShowLoadBalancerApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return false
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return false
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return false
	}
	if resp.ReturnObj[0].Ipv6Address != "" {
		return true
	} else {
		return false
	}
}

type CtyunElbListenerConfig struct {
	RegionID            types.String `tfsdk:"region_id"`             //区域ID
	LoadBalancerID      types.String `tfsdk:"loadbalancer_id"`       //负载均衡实例ID
	Name                types.String `tfsdk:"name"`                  //唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32
	Description         types.String `tfsdk:"description"`           //支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128
	Protocol            types.String `tfsdk:"protocol"`              //监听协议。取值范围：TCP、UDP、HTTP、HTTPS
	ProtocolPort        types.Int32  `tfsdk:"protocol_port"`         //负载均衡实例监听端口。取值：1-65535
	CertificateID       types.String `tfsdk:"certificate_id"`        //证书ID。当protocol为HTTPS时,此参数必选
	CaEnabled           types.Bool   `tfsdk:"ca_enabled"`            //是否开启双向认证。false（不开启）、true（开启）
	ClientCertificateID types.String `tfsdk:"client_certificate_id"` //双向认证的证书ID
	DefaultActionType   types.String `tfsdk:"default_action_type"`   //默认规则动作
	RedirectListenerID  types.String `tfsdk:"redirect_listener_id"`  //重定向监听器ID，当type为redirect时，此字段必填
	TargetGroups        types.List   `tfsdk:"target_groups"`         //后端服务组
	AccessControlID     types.String `tfsdk:"access_control_id"`     //访问控制ID
	AccessControlType   types.String `tfsdk:"access_control_type"`   //访问控制类型。取值范围：Close（未启用）、White（白名单）、Black（黑名单）
	ForwardedForEnabled types.Bool   `tfsdk:"forwarded_for_enabled"` //x forward for功能。false（未开启）、true（开启）
	ID                  types.String `tfsdk:"id"`                    //监听器 ID
	AzName              types.String `tfsdk:"az_name"`               //可用区名称
	ProjectID           types.String `tfsdk:"project_id"`            //项目ID
	Status              types.String `tfsdk:"status"`                //监听器状态: DOWN / ACTIVE
	CreatedTime         types.String `tfsdk:"created_time"`          //创建时间，为UTC格式
	UpdatedTime         types.String `tfsdk:"updated_time"`          //更新时间，为UTC格式
	EnableNat64         types.Bool   `tfsdk:"enable_nat_64"`         //是否开启 nat64
	ListenerQps         types.Int32  `tfsdk:"listener_qps"`          //qps 大小
	EstablishTimeout    types.Int32  `tfsdk:"establish_timeout"`     //建立连接超时时间，单位秒，取值范围： 1 - 1800
	IdleTimeout         types.Int32  `tfsdk:"idle_timeout"`          //链接空闲断开超时时间，单位秒，取值范围：1 - 300
	ResponseTimeout     types.Int32  `tfsdk:"response_timeout"`      //响应超时，单位秒，取值范围：1 - 300
	ListenerCps         types.Int32  `tfsdk:"listener_cps"`          //cps 大小
}

type TargetGroupsModel struct {
	TargetGroupID types.String `tfsdk:"target_group_id"`
	Weight        types.Int32  `tfsdk:"weight"`
}
