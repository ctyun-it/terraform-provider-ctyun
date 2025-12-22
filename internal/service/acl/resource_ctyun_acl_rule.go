package acl

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
)

var (
	_ resource.Resource                = &CtyunAclRule{}
	_ resource.ResourceWithConfigure   = &CtyunAclRule{}
	_ resource.ResourceWithImportState = &CtyunAclRule{}
)

type CtyunAclRule struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunAclRule) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_acl_rule"
}

func (c *CtyunAclRule) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunAclRule() resource.Resource {
	return &CtyunAclRule{}
}

func (c *CtyunAclRule) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID], [aclId], [direction],[projectId],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunAclRuleConfig
	var ID, aclId, direction, projectId, regionId string
	// 根据分隔符数量判断是否输入了regionID,projectId
	if strings.Count(request.ID, common.ImportSeparator) == 2 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		projectId = c.meta.GetExtraIfEmpty(projectId, common.ExtraProjectId)
		err = terraform_extend.Split(request.ID, &ID, &direction, &aclId)
		if err != nil {
			return
		}
	} else if strings.Count(request.ID, common.ImportSeparator) == 3 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &ID, &direction, &aclId, &projectId)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &ID, &direction, &aclId, &projectId, &regionId)
		if err != nil {
			return
		}
	}

	if ID == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}

	if aclId == "" {
		err = fmt.Errorf("aclId不能为空")
		return
	}
	if direction == "" {
		err = fmt.Errorf("direction不能为空")
		return
	}
	config.ID = types.StringValue(ID)
	config.RegionID = types.StringValue(regionId)
	config.AclID = types.StringValue(aclId)
	if projectId != "" {
		config.ProjectID = types.StringValue(projectId)
	}
	config.Direction = types.StringValue(direction)

	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunAclRule) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10028588",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
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
			"acl_id": schema.StringAttribute{
				Required:    true,
				Description: "acl_id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.AclID(),
				},
			},
			"direction": schema.StringAttribute{
				Required:    true,
				Description: "acl类型，支持更新。取值范围：ingress, egress",
				Validators: []validator.String{
					stringvalidator.OneOf(business.AclDirectionIngress, business.AclDirectionEgress),
				},
			},
			"priority": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "优先级，支持更新。取值范围： 1 - 32766，不填默认100",
				Default:     int32default.StaticInt32(100),
				Validators: []validator.Int32{
					int32validator.Between(1, 32766),
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "协议，支持更新。取值范围：all, icmp, tcp, udp",
				Validators: []validator.String{
					stringvalidator.OneOf(business.AclRuleProtocols...),
				},
			},
			"ip_version": schema.StringAttribute{
				Required:    true,
				Description: "IP版本，支持更新。取值范围：ipv4, ipv6",
				Validators: []validator.String{
					stringvalidator.OneOf("ipv4", "ipv6"),
				},
			},
			"destination_port": schema.StringAttribute{
				Optional:    true,
				Description: "目的地址端口范围，支持更新。示例 8080:8085",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^(?:[1-9]\d{0,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5]):(?:[1-9]\d{0,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$`), "输入的端口范围有误!"),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.AclRuleProtocolTCP)),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.AclRuleProtocolUDP)),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.AclRuleProtocolALL)),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.AclRuleProtocolICMP)),
				},
			},
			"source_port": schema.StringAttribute{
				Optional:    true,
				Description: "源地址端口范围，支持更新。示例： 8080:8085",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^(?:[1-9]\d{0,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5]):(?:[1-9]\d{0,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$`), "输入的端口范围有误!"),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.AclRuleProtocolTCP)),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.AclRuleProtocolUDP)),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.AclRuleProtocolALL)),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.AclRuleProtocolICMP)),
				},
			},
			"source_ip_address": schema.StringAttribute{
				Required:    true,
				Description: "源地址，支持更新。支持cidr格式",
				Validators: []validator.String{
					validator2.Cidr(),
				},
			},
			"destination_ip_address": schema.StringAttribute{
				Required:    true,
				Description: "目的地址，支持更新。支持cidr格式",
				Validators: []validator.String{
					validator2.Cidr(),
				},
			},
			"action": schema.StringAttribute{
				Required:    true,
				Description: "策略，支持更新。取值范围：accept, drop",
				Validators: []validator.String{
					stringvalidator.OneOf(business.AclRuleActionAccept, business.AclRuleActionDrop),
				},
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "acl 规则是否启用，支持更新。默认启用",
				Default:     booldefault.StaticBool(true),
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "acl规则描述，支持更新",
				Validators: []validator.String{
					validator2.Desc(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "acl 规则id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *CtyunAclRule) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunAclRuleConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
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
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunAclRule) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunAclRuleConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunAclRule) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunAclRuleConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunAclRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunAclRule) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunAclRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunAclRule) create(ctx context.Context, config *CtyunAclRuleConfig) error {
	params := &ctvpc.CtvpcCreateAclRuleRequest{
		ClientToken: uuid.NewString(),
		RegionID:    config.RegionID.ValueString(),
		AclID:       config.AclID.ValueString(),
	}

	var rules []*ctvpc.CtvpcCreateAclRuleRulesRequest
	var rule ctvpc.CtvpcCreateAclRuleRulesRequest
	rule.Direction = config.Direction.ValueString()
	rule.Priority = config.Priority.ValueInt32()
	rule.Protocol = config.Protocol.ValueString()
	rule.IpVersion = config.IpVersion.ValueString()
	if !config.DestinationPort.IsNull() {
		rule.DestinationPort = config.DestinationPort.ValueStringPointer()
	}
	if !config.SourcePort.IsNull() {
		rule.SourcePort = config.SourcePort.ValueStringPointer()
	}
	rule.SourceIpAddress = config.SourceIpAddress.ValueString()
	rule.DestinationIpAddress = config.DestinationIpAddress.ValueString()
	rule.Action = config.Action.ValueString()
	rule.Enabled = map[bool]string{true: business.AclRuleEnable, false: business.AclRuleDisable}[config.Enabled.ValueBool()]
	if !config.Description.IsNull() {
		rule.Description = config.Description.ValueString()
	}
	rules = append(rules, &rule)
	params.Rules = rules
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateAclRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建acl规则失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	err = c.getRuleID(ctx, config)
	if err != nil {
		return err

	}
	return nil
}

func (c *CtyunAclRule) getAndMerge(ctx context.Context, config *CtyunAclRuleConfig) error {
	ingressDetail, egressDetail, err := c.getRuleDetail(ctx, config)
	if err != nil {
		return err
	}
	if ingressDetail == nil {
		config.Direction = types.StringValue(*egressDetail.Direction)
		config.Priority = types.Int32Value(egressDetail.Priority)
		config.Protocol = types.StringValue(*egressDetail.Protocol)
		config.IpVersion = types.StringValue(*egressDetail.IpVersion)
		if config.Protocol.ValueString() == business.AclRuleProtocolTCP ||
			config.Protocol.ValueString() == business.AclRuleProtocolUDP {
			config.DestinationPort = types.StringValue(*egressDetail.DestinationPort)
			config.SourcePort = types.StringValue(*egressDetail.SourcePort)
		}
		config.SourceIpAddress = types.StringValue(*egressDetail.SourceIpAddress)
		config.DestinationIpAddress = types.StringValue(*egressDetail.DestinationIpAddress)
		config.Action = types.StringValue(*egressDetail.Action)
		config.Enabled = types.BoolValue(map[string]bool{business.AclRuleEnable: true, business.AclRuleDisable: false}[*egressDetail.Enabled])
		config.Description = types.StringValue(*egressDetail.Description)
	} else {
		config.Direction = types.StringValue(*ingressDetail.Direction)
		config.Priority = types.Int32Value(ingressDetail.Priority)
		config.Protocol = types.StringValue(*ingressDetail.Protocol)
		config.IpVersion = types.StringValue(*ingressDetail.IpVersion)
		if config.Protocol.ValueString() == business.AclRuleProtocolTCP ||
			config.Protocol.ValueString() == business.AclRuleProtocolUDP {
			config.DestinationPort = types.StringValue(*ingressDetail.DestinationPort)
			config.SourcePort = types.StringValue(*ingressDetail.SourcePort)
		}
		config.SourceIpAddress = types.StringValue(*ingressDetail.SourceIpAddress)
		config.DestinationIpAddress = types.StringValue(*ingressDetail.DestinationIpAddress)
		config.Action = types.StringValue(*ingressDetail.Action)
		config.Enabled = types.BoolValue(map[string]bool{business.AclRuleEnable: true, business.AclRuleDisable: false}[*ingressDetail.Enabled])
		config.Description = types.StringValue(*ingressDetail.Description)
	}
	return nil
}

func (c *CtyunAclRule) getRuleID(ctx context.Context, config *CtyunAclRuleConfig) error {
	// 通过获取列表获取，rule id
	ruleList, err := c.getRuleList(ctx, config)
	if err != nil {
		return err
	}
	// 若刚刚创建的规则为入规则，遍历入规则列表
	if config.Direction.ValueString() == business.AclDirectionIngress {
		ingressRuleList := ruleList.InRules
		for _, ingressRule := range ingressRuleList {
			same := c.ingressCheckSame(ingressRule, config)
			if same {
				config.ID = types.StringValue(*ingressRule.AclRuleID)
				break
			}
		}
	} else if config.Direction.ValueString() == business.AclDirectionEgress {
		egressRuleList := ruleList.OutRules
		for _, egressRule := range egressRuleList {
			same := c.egressCheckSame(egressRule, config)
			if same {
				config.ID = types.StringValue(*egressRule.AclRuleID)
			}
		}
	} else {
		err = fmt.Errorf("direction 取值有误！当前值为%s", config.Direction.ValueString())
		return err
	}
	return nil
}

func (c *CtyunAclRule) getRuleList(ctx context.Context, config *CtyunAclRuleConfig) (*ctvpc.CtvpcListAclRuleReturnObjResponse, error) {
	params := &ctvpc.CtvpcListAclRuleRequest{
		RegionID: config.RegionID.ValueString(),
		AclID:    config.AclID.ValueString(),
	}
	if !config.ProjectID.IsUnknown() && !config.ProjectID.IsNull() {
		params.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListAclRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取acl规则列表失败（acl id=%s），接口返回nil，请联系研发确认问题原因！", config.AclID.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp.ReturnObj, nil
}

func (c *CtyunAclRule) ingressCheckSame(rule *ctvpc.CtvpcListAclRuleReturnObjInRulesResponse, config *CtyunAclRuleConfig) bool {
	if *rule.Direction != config.Direction.ValueString() ||
		rule.Priority != config.Priority.ValueInt32() ||
		*rule.Protocol != config.Protocol.ValueString() ||
		*rule.IpVersion != config.IpVersion.ValueString() ||
		*rule.SourceIpAddress != config.SourceIpAddress.ValueString() ||
		*rule.DestinationIpAddress != config.DestinationIpAddress.ValueString() ||
		*rule.Action != config.Action.ValueString() ||
		*rule.Enabled != map[bool]string{true: business.AclRuleEnable, false: business.AclRuleDisable}[config.Enabled.ValueBool()] {
		return false
	}
	if config.Protocol.ValueString() == business.AclRuleProtocolTCP ||
		config.Protocol.ValueString() == business.AclRuleProtocolUDP {
		if *rule.DestinationPort != config.DestinationPort.ValueString() ||
			*rule.SourcePort != config.SourcePort.ValueString() {
			return false
		}
	}
	return true
}

func (c *CtyunAclRule) egressCheckSame(rule *ctvpc.CtvpcListAclRuleReturnObjOutRulesResponse, config *CtyunAclRuleConfig) bool {
	if *rule.Direction != config.Direction.ValueString() ||
		rule.Priority != config.Priority.ValueInt32() ||
		*rule.Protocol != config.Protocol.ValueString() ||
		*rule.IpVersion != config.IpVersion.ValueString() ||
		*rule.DestinationPort != config.DestinationPort.ValueString() ||
		*rule.SourcePort != config.SourcePort.ValueString() ||
		*rule.SourceIpAddress != config.SourceIpAddress.ValueString() ||
		*rule.DestinationIpAddress != config.DestinationIpAddress.ValueString() ||
		*rule.Action != config.Action.ValueString() ||
		*rule.Enabled != map[bool]string{true: business.AclRuleEnable, false: business.AclRuleDisable}[config.Enabled.ValueBool()] {
		return false
	}
	return true
}

func (c *CtyunAclRule) getRuleDetail(ctx context.Context, config *CtyunAclRuleConfig) (*ctvpc.CtvpcListAclRuleReturnObjInRulesResponse, *ctvpc.CtvpcListAclRuleReturnObjOutRulesResponse, error) {
	ruleList, err := c.getRuleList(ctx, config)
	if err != nil {
		return nil, nil, err
	}
	if config.Direction.ValueString() == business.AclDirectionIngress {
		ingressList := ruleList.InRules
		for _, ingressRule := range ingressList {
			if *ingressRule.AclRuleID == config.ID.ValueString() {
				return ingressRule, nil, nil
			}
		}
	} else if config.Direction.ValueString() == business.AclDirectionEgress {
		egressList := ruleList.OutRules
		for _, egressRule := range egressList {
			if *egressRule.AclRuleID == config.ID.ValueString() {
				return nil, egressRule, nil
			}
		}
	} else {
		err = fmt.Errorf("direction 取值有误！当前值为%s", config.Direction.ValueString())
		return nil, nil, err
	}
	return nil, nil, nil
}

func (c *CtyunAclRule) update(ctx context.Context, state *CtyunAclRuleConfig, plan *CtyunAclRuleConfig) error {
	uuidStr := uuid.NewString()
	params := &ctvpc.CtvpcUpdateAclRuleAttributeRequest{
		ClientToken: &uuidStr,
		RegionID:    state.RegionID.ValueString(),
		AclID:       state.AclID.ValueString(),
		Rules:       nil,
	}
	var rules []*ctvpc.CtvpcUpdateAclRuleAttributeRulesRequest
	var rule ctvpc.CtvpcUpdateAclRuleAttributeRulesRequest
	rule.AclRuleID = state.ID.ValueString()
	rule.Direction = plan.Direction.ValueString()
	rule.Priority = plan.Priority.ValueInt32()
	rule.Protocol = plan.Protocol.ValueString()
	rule.IpVersion = plan.IpVersion.ValueString()
	if !plan.DestinationPort.IsNull() {
		rule.DestinationPort = plan.DestinationPort.ValueStringPointer()
	}
	if !plan.SourcePort.IsNull() {
		rule.SourcePort = plan.SourcePort.ValueStringPointer()
	}
	rule.SourceIpAddress = plan.SourceIpAddress.ValueString()
	rule.DestinationIpAddress = plan.DestinationIpAddress.ValueString()
	rule.Action = plan.Action.ValueString()
	rule.Enabled = map[bool]string{true: business.AclRuleEnable, false: business.AclRuleDisable}[plan.Enabled.ValueBool()]
	if !plan.Description.IsNull() {
		rule.Description = plan.Description.ValueStringPointer()
	}
	rules = append(rules, &rule)
	params.Rules = rules
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcUpdateAclRuleAttributeApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新acl规则失败（acl_id =%s,acl_rule_id =%s），接口返回nil，请联系研发确认问题原因！", state.AclID.ValueString(), plan.ID)
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	return nil
}

func (c *CtyunAclRule) delete(ctx context.Context, config CtyunAclRuleConfig) error {
	aclRuleIDs := []string{config.ID.ValueString()}
	params := &ctvpc.CtvpcDeleteAclRuleRequest{
		ClientToken:   uuid.NewString(),
		RegionID:      config.RegionID.ValueString(),
		AclID:         config.AclID.ValueString(),
		AclRuleIDList: aclRuleIDs,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteAclRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除acl规则失败（acl_id =%s,acl_rule_id =%s），接口返回nil，请联系研发确认问题原因！", config.AclID.ValueString(), config.ID)
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	return nil
}

type CtyunAclRuleConfig struct {
	RegionID             types.String `tfsdk:"region_id"`
	ProjectID            types.String `tfsdk:"project_id"`
	AclID                types.String `tfsdk:"acl_id"`
	Direction            types.String `tfsdk:"direction"`
	Priority             types.Int32  `tfsdk:"priority"`
	Protocol             types.String `tfsdk:"protocol"`
	IpVersion            types.String `tfsdk:"ip_version"`
	DestinationPort      types.String `tfsdk:"destination_port"`
	SourcePort           types.String `tfsdk:"source_port"`
	SourceIpAddress      types.String `tfsdk:"source_ip_address"`
	DestinationIpAddress types.String `tfsdk:"destination_ip_address"`
	Action               types.String `tfsdk:"action"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	Description          types.String `tfsdk:"description"`
	ID                   types.String `tfsdk:"id"`
}
