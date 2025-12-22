package acl

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CtyunAclRules struct {
	meta *common.CtyunMetadata
}

func NewCtyunAclRules() datasource.DataSource {
	return &CtyunAclRules{}
}
func (c *CtyunAclRules) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunAclRules) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_acl_rules"
}

func (c *CtyunAclRules) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10028588",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，默认使用provider配置",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"acl_id": schema.StringAttribute{
				Required:    true,
				Description: "ACL ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "项目ID",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "ACL名称过滤条件",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "ACL描述过滤条件",
			},
			"vpc_id": schema.StringAttribute{
				Computed:    true,
				Description: "VPC ID过滤条件",
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "是否启用",
			},
			"in_policy_id": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "入方向策略ID列表",
			},
			"out_policy_id": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "出方向策略ID列表",
			},
			"subnet_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "绑定的子网ID",
			},
			"in_rules": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"acl_rule_id": schema.StringAttribute{
							Computed:    true,
							Description: "ACL规则ID",
						},
						"direction": schema.StringAttribute{
							Computed:    true,
							Description: "规则方向（ingress/egress）",
						},
						"priority": schema.Int64Attribute{
							Computed:    true,
							Description: "规则优先级",
						},
						"protocol": schema.StringAttribute{
							Computed:    true,
							Description: "协议类型（tcp/udp/icmp/all）",
						},
						"ip_version": schema.StringAttribute{
							Computed:    true,
							Description: "IP版本（ipv4/ipv6）",
						},
						"destination_port": schema.StringAttribute{
							Computed:    true,
							Description: "目标端口范围",
						},
						"source_port": schema.StringAttribute{
							Computed:    true,
							Description: "源端口范围",
						},
						"source_ip_address": schema.StringAttribute{
							Computed:    true,
							Description: "源IP地址范围",
						},
						"destination_ip_address": schema.StringAttribute{
							Computed:    true,
							Description: "目标IP地址范围",
						},
						"action": schema.StringAttribute{
							Computed:    true,
							Description: "动作（accept/drop）",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "是否启用",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "规则描述",
						},
					},
				},
				Description: "入方向规则列表",
			},
			"out_rules": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"acl_rule_id": schema.StringAttribute{
							Computed:    true,
							Description: "ACL规则ID",
						},
						"direction": schema.StringAttribute{
							Computed:    true,
							Description: "规则方向（ingress/egress）",
						},
						"priority": schema.Int32Attribute{
							Computed:    true,
							Description: "规则优先级",
						},
						"protocol": schema.StringAttribute{
							Computed:    true,
							Description: "协议类型（tcp/udp/icmp/all）",
						},
						"ip_version": schema.StringAttribute{
							Computed:    true,
							Description: "IP版本（ipv4/ipv6）",
						},
						"destination_port": schema.StringAttribute{
							Computed:    true,
							Description: "目标端口范围",
						},
						"source_port": schema.StringAttribute{
							Computed:    true,
							Description: "源端口范围",
						},
						"source_ip_address": schema.StringAttribute{
							Computed:    true,
							Description: "源IP地址范围",
						},
						"destination_ip_address": schema.StringAttribute{
							Computed:    true,
							Description: "目标IP地址范围",
						},
						"action": schema.StringAttribute{
							Computed:    true,
							Description: "动作（accept/drop）",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "是否启用",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "规则描述",
						},
					},
				},
				Description: "出方向规则列表",
			},
		},
	}
}

func (c *CtyunAclRules) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunAclRulesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	config.RegionID = types.StringValue(regionId)
	rule, err := c.getRuleList(ctx, &config)
	if err != nil {
		return
	}
	config.Name = types.StringValue(*rule.Name)
	config.Description = types.StringValue(*rule.Description)
	config.VpcID = types.StringValue(*rule.VpcID)
	config.Enabled = types.BoolValue(map[string]bool{business.AclRuleEnable: true, business.AclRuleDisable: false}[*rule.Enabled])
	config.InPolicyID = rule.InPolicyID
	config.OutPolicyID = rule.OutPolicyID
	config.SubnetIDs = rule.SubnetIDs
	var inRules []CtyunAclRuleModel
	for _, ruleItem := range rule.InRules {
		var inRule CtyunAclRuleModel
		inRule.AclRuleID = types.StringValue(*ruleItem.AclRuleID)
		inRule.Direction = types.StringValue(*ruleItem.Direction)
		inRule.Priority = types.Int32Value(ruleItem.Priority)
		inRule.IpVersion = types.StringValue(*ruleItem.IpVersion)
		inRule.DestinationPort = types.StringValue(*ruleItem.DestinationPort)
		inRule.SourcePort = types.StringValue(*ruleItem.SourcePort)
		inRule.SourceIpAddress = types.StringValue(*ruleItem.SourceIpAddress)
		inRule.DestinationIpAddress = types.StringValue(*ruleItem.DestinationIpAddress)
		inRule.Action = types.StringValue(*ruleItem.Action)
		inRule.Enabled = types.BoolValue(map[string]bool{business.AclRuleEnable: true, business.AclRuleDisable: false}[*ruleItem.Enabled])
		inRule.Description = types.StringValue(*ruleItem.Description)
		inRules = append(inRules, inRule)
	}
	config.InRules = inRules

	var outRules []CtyunAclRuleModel
	for _, ruleItem := range rule.OutRules {
		var outRule CtyunAclRuleModel
		outRule.AclRuleID = types.StringValue(*ruleItem.AclRuleID)
		outRule.Direction = types.StringValue(*ruleItem.Direction)
		outRule.Priority = types.Int32Value(ruleItem.Priority)
		outRule.IpVersion = types.StringValue(*ruleItem.IpVersion)
		outRule.DestinationPort = types.StringValue(*ruleItem.DestinationPort)
		outRule.SourcePort = types.StringValue(*ruleItem.SourcePort)
		outRule.SourceIpAddress = types.StringValue(*ruleItem.SourceIpAddress)
		outRule.DestinationIpAddress = types.StringValue(*ruleItem.DestinationIpAddress)
		outRule.Action = types.StringValue(*ruleItem.Action)
		outRule.Enabled = types.BoolValue(map[string]bool{business.AclRuleEnable: true, business.AclRuleDisable: false}[*ruleItem.Enabled])
		outRule.Description = types.StringValue(*ruleItem.Description)
		outRules = append(outRules, outRule)
	}
	config.OutRules = outRules
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *CtyunAclRules) getRuleList(ctx context.Context, config *CtyunAclRulesConfig) (*ctvpc.CtvpcListAclRuleReturnObjResponse, error) {
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

type CtyunAclRuleModel struct {
	AclRuleID            types.String `tfsdk:"acl_rule_id"`
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
}

type CtyunAclRulesConfig struct {
	RegionID    types.String        `tfsdk:"region_id"`
	AclID       types.String        `tfsdk:"acl_id"`
	ProjectID   types.String        `tfsdk:"project_id"`
	Name        types.String        `tfsdk:"name"`
	Description types.String        `tfsdk:"description"`
	VpcID       types.String        `tfsdk:"vpc_id"`
	Enabled     types.Bool          `tfsdk:"enabled"`
	InPolicyID  []string            `tfsdk:"in_policy_id"`
	OutPolicyID []string            `tfsdk:"out_policy_id"`
	InRules     []CtyunAclRuleModel `tfsdk:"in_rules"`
	OutRules    []CtyunAclRuleModel `tfsdk:"out_rules"`
	SubnetIDs   []string            `tfsdk:"subnet_ids"`
}
