package vpc

import (
	"context"
	"errors"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCtyunSecurityGroupRule() resource.Resource {
	return &ctyunSecurityGroupRule{}
}

type ctyunSecurityGroupRule struct {
	meta                 *common.CtyunMetadata
	securityGroupService *business.SecurityGroupService
}

func (c *ctyunSecurityGroupRule) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_security_group_rule"
}

func (c *ctyunSecurityGroupRule) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026730/10225510`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "id",
			},
			"security_group_id": schema.StringAttribute{
				Required:    true,
				Description: "安全组id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SecurityGroupValidate(),
				},
			},
			"direction": schema.StringAttribute{
				Required:    true,
				Description: "规则方向，egress：出方向，ingress：入方向",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.SecurityGroupRuleDirections...),
				},
			},
			"action": schema.StringAttribute{
				Required:    true,
				Description: "拒绝策略，accept：允许，drop：拒绝",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.SecurityGroupRuleActions...),
				},
			},
			"priority": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "优先级：1~100，取值越小优先级越大，默认优先级为50",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Between(1, 100),
				},
				Default: int64default.StaticInt64(1),
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "协议类型: tcp、udp、icmp、any，当此值填写any时，range的值不能设置",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.SecurityGroupRuleProtocols...),
				},
			},
			"ether_type": schema.StringAttribute{
				Required:    true,
				Description: "IP类型：ipv4、ipv6",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.SecurityGroupRuleEtherTypes...),
				},
			},
			"dest_cidr_ip": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "远端地址，为cidr地址格式，如果不填默认为0.0.0.0/0",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.Cidr(),
				},
				Default: stringdefault.StaticString("0.0.0.0/0"),
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "描述，长度1-128，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(128),
				},
			},
			"range": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "安全组开放的传输层协议相关的源端端口范围，格式如：8000-9000，如果仅开放单一端口则直接填写，如：22，中间不能有空格以及其他特殊字符；如果protocol的值为any，请保证此值留空，如果protocol的值为tcp或udp，此值必填",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.SecurityGroupRuleProtocolTcp),
						types.StringValue(business.SecurityGroupRuleProtocolUdp),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.SecurityGroupRuleProtocolAny),
					),
					validator2.Range("-", 1, 65535),
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

func (c *ctyunSecurityGroupRule) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunSecurityGroupRuleConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.checkCreate(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	var id string
	regionId := plan.RegionId.ValueString()
	direction := plan.Direction.ValueString()
	securityGroupId := plan.SecurityGroupId.ValueString()
	clientToken := uuid.NewString()

	requestDirection, err := business.SecurityGroupRuleDirectionMap.FromOriginalScene(direction, business.SecurityGroupRuleDirectionMapScene1)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	requestAction, err := business.SecurityGroupRuleActionMap.FromOriginalScene(plan.Action.ValueString(), business.SecurityGroupRuleActionMapScene1)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	requestProtocol, err := business.SecurityGroupRuleProtocolMap.FromOriginalScene(plan.Protocol.ValueString(), business.SecurityGroupRuleProtocolMapScene1)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	requestEthertype, err := business.SecurityGroupRuleEtherTypeMap.FromOriginalScene(plan.EtherType.ValueString(), business.SecurityGroupRuleEtherTypeMapScene1)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	requestPriority := int(plan.Priority.ValueInt64())
	requestDestCidrIp := plan.DestCidrIp.ValueString()
	requestDescription := plan.Description.ValueString()
	requestRange := plan.Range.ValueString()

	if business.IsEgress(direction) {
		resp, err := c.meta.Apis.CtEcsApis.SecurityGroupRuleEgressCreateApi.Do(ctx, c.meta.Credential, &ctecs.SecurityGroupRuleEgressCreateRequest{
			RegionId:        regionId,
			SecurityGroupId: securityGroupId,
			ClientToken:     clientToken,
			SecurityGroupRules: []ctecs.SecurityGroupRuleEgressCreateSecurityGroupRulesRequest{
				{
					Direction:   requestDirection.(string),
					Action:      requestAction.(string),
					Protocol:    requestProtocol.(string),
					Ethertype:   requestEthertype.(string),
					Priority:    requestPriority,
					DestCidrIp:  requestDestCidrIp,
					Description: requestDescription,
					Range:       requestRange,
				},
			},
		})
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
		targetId, err2 := c.checkAndGetId(resp.SgRuleIds)
		if err2 != nil {
			response.Diagnostics.AddError(err2.Error(), err2.Error())
			return
		}
		id = targetId
	} else {
		resp, err := c.meta.Apis.CtEcsApis.SecurityGroupRuleIngressCreateApi.Do(ctx, c.meta.Credential, &ctecs.SecurityGroupRuleIngressCreateRequest{
			RegionId:        regionId,
			SecurityGroupId: securityGroupId,
			ClientToken:     clientToken,
			SecurityGroupRules: []ctecs.SecurityGroupRuleIngressCreateSecurityGroupRulesRequest{
				{
					Direction:   requestDirection.(string),
					Action:      requestAction.(string),
					Protocol:    requestProtocol.(string),
					Ethertype:   requestEthertype.(string),
					Priority:    requestPriority,
					DestCidrIp:  requestDestCidrIp,
					Description: requestDescription,
					Range:       requestRange,
				},
			},
		})
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
		targetId, err2 := c.checkAndGetId(resp.SgRuleIds)
		if err2 != nil {
			response.Diagnostics.AddError(err2.Error(), err2.Error())
			return
		}
		id = targetId
	}

	plan.Id = types.StringValue(id)
	plan.RegionId = types.StringValue(regionId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeSecurityGroupRule(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunSecurityGroupRule) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunSecurityGroupRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	instance, err := c.getAndMergeSecurityGroupRule(ctx, state)
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

func (c *ctyunSecurityGroupRule) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state CtyunSecurityGroupRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	var plan CtyunSecurityGroupRuleConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)

	requestId := state.Id.ValueString()
	requestRegionId := state.RegionId.ValueString()
	requestSecurityGroupId := state.SecurityGroupId.ValueString()
	requestClientToken := uuid.NewString()
	requestDescription := plan.Description.ValueString()

	// 判断描述是否相同
	if !state.Description.Equal(plan.Description) {
		direction := plan.Direction.ValueString()
		var err error
		if business.IsEgress(direction) {
			_, err = c.meta.Apis.CtVpcApis.SecurityGroupRuleEgressModifyApi.Do(ctx, c.meta.Credential, &ctvpc.SecurityGroupRuleEgressModifyRequest{
				RegionId:            requestRegionId,
				SecurityGroupId:     requestSecurityGroupId,
				SecurityGroupRuleId: requestId,
				ClientToken:         requestClientToken,
				Description:         requestDescription,
			})

		} else {
			_, err = c.meta.Apis.CtVpcApis.SecurityGroupRuleIngressModifyApi.Do(ctx, c.meta.Credential, &ctvpc.SecurityGroupRuleIngressModifyRequest{
				RegionId:            requestRegionId,
				SecurityGroupId:     requestSecurityGroupId,
				SecurityGroupRuleId: requestId,
				ClientToken:         requestClientToken,
				Description:         requestDescription,
			})
		}
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
	}

	instance, err := c.getAndMergeSecurityGroupRule(ctx, state)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunSecurityGroupRule) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunSecurityGroupRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	var err error

	requestRegionId := state.RegionId.ValueString()
	requestSecurityGroupId := state.SecurityGroupId.ValueString()
	requestSecurityGroupRuleId := state.Id.ValueString()
	requestClientToken := uuid.NewString()

	direction := state.Direction.ValueString()
	if business.IsEgress(direction) {
		_, err = c.meta.Apis.CtVpcApis.SecurityGroupRuleEgressRevokeApi.Do(ctx, c.meta.Credential, &ctvpc.SecurityGroupRuleEgressRevokeRequest{
			RegionId:            requestRegionId,
			SecurityGroupId:     requestSecurityGroupId,
			SecurityGroupRuleId: requestSecurityGroupRuleId,
			ClientToken:         requestClientToken,
		})
	} else {
		_, err = c.meta.Apis.CtVpcApis.SecurityGroupRuleIngressRevokeApi.Do(ctx, c.meta.Credential, &ctvpc.SecurityGroupRuleIngressRevokeRequest{
			RegionId:            requestRegionId,
			SecurityGroupId:     requestSecurityGroupId,
			SecurityGroupRuleId: requestSecurityGroupRuleId,
			ClientToken:         requestClientToken,
		})
	}
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [securityGroupRuleId],[securityGroupId],[regionId]
func (c *ctyunSecurityGroupRule) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunSecurityGroupRuleConfig
	var securityGroupRuleId, securityGroupId, regionId string
	err := terraform_extend.Split(request.ID, &securityGroupRuleId, &securityGroupId, &regionId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.Id = types.StringValue(securityGroupRuleId)
	cfg.SecurityGroupId = types.StringValue(securityGroupId)
	cfg.RegionId = types.StringValue(regionId)

	instance, err := c.getAndMergeSecurityGroupRule(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunSecurityGroupRule) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.securityGroupService = business.NewSecurityGroupService(meta)
}

// getAndMergeSecurityGroupRule 查询安全组规则
func (c *ctyunSecurityGroupRule) getAndMergeSecurityGroupRule(ctx context.Context, cfg CtyunSecurityGroupRuleConfig) (*CtyunSecurityGroupRuleConfig, error) {
	request := &ctvpc.SecurityGroupRuleDescribeRequest{
		SecurityGroupRuleId: cfg.Id.ValueString(),
		SecurityGroupId:     cfg.SecurityGroupId.ValueString(),
		RegionId:            cfg.RegionId.ValueString(),
	}
	response, err := c.meta.Apis.CtVpcApis.SecurityGroupRuleDescribeApi.Do(ctx, c.meta.Credential, request)
	if err != nil {
		// 如果查询不到信息会报异常，此时直接返回空
		if err.ErrorCode() == common.OpenapiSecurityGroupRuleNotFound {
			return nil, nil
		}
		return nil, err
	}

	protocol, err2 := business.SecurityGroupRuleProtocolMap.ToOriginalScene(response.Protocol, business.SecurityGroupRuleProtocolMapScene1)
	if err2 != nil {
		return nil, err2
	}
	ethertype, err2 := business.SecurityGroupRuleEtherTypeMap.ToOriginalScene(response.Ethertype, business.SecurityGroupRuleEtherTypeMapScene1)
	if err2 != nil {
		return nil, err2
	}
	cfg.Id = types.StringValue(response.Id)
	cfg.SecurityGroupId = types.StringValue(response.SecurityGroupId)
	cfg.Direction = types.StringValue(response.Direction)
	cfg.Action = types.StringValue(response.Action)
	cfg.Priority = types.Int64Value(int64(response.Priority))
	cfg.Protocol = types.StringValue(protocol.(string))
	cfg.Range = types.StringValue(response.Range)
	cfg.EtherType = types.StringValue(ethertype.(string))
	cfg.DestCidrIp = types.StringValue(response.DestCidrIp)
	cfg.Description = types.StringValue(response.Description)
	return &cfg, nil
}

// checkAndGetId 校验创建是否成功，如果后台没有返回id则说明规则已经存在了
func (c *ctyunSecurityGroupRule) checkAndGetId(ids []string) (string, error) {
	if len(ids) == 0 {
		return "", errors.New("对应的规则已经存在")
	}
	return ids[0], nil
}

// checkCreate 校验创建动作是否能执行
func (c *ctyunSecurityGroupRule) checkCreate(ctx context.Context, plan CtyunSecurityGroupRuleConfig) error {
	return c.securityGroupService.MustExist(ctx, plan.SecurityGroupId.ValueString(), plan.RegionId.ValueString())
}

type CtyunSecurityGroupRuleConfig struct {
	Id              types.String `tfsdk:"id"`
	SecurityGroupId types.String `tfsdk:"security_group_id"`
	Direction       types.String `tfsdk:"direction"`
	Action          types.String `tfsdk:"action"`
	Priority        types.Int64  `tfsdk:"priority"`
	Protocol        types.String `tfsdk:"protocol"`
	EtherType       types.String `tfsdk:"ether_type"`
	DestCidrIp      types.String `tfsdk:"dest_cidr_ip"`
	Description     types.String `tfsdk:"description"`
	Range           types.String `tfsdk:"range"`
	RegionId        types.String `tfsdk:"region_id"`
}
