package sdwan

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/sdwan"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &CtyunSdwanAclRules{}
)

func NewCtyunSdwanAclRules() datasource.DataSource {
	return &CtyunSdwanAclRules{}
}

type CtyunSdwanAclRules struct {
	meta *common.CtyunMetadata
}

type CtyunSdwanAclRulesConfig struct {
	AclID types.String `tfsdk:"acl_id"`
	Rules []RuleInfo   `tfsdk:"rules"`
	ID    types.String `tfsdk:"id"`
}

type RuleInfo struct {
	ID           types.String `tfsdk:"id"`
	Direction    types.String `tfsdk:"direction"`
	Action       types.String `tfsdk:"action"`
	Protocol     types.String `tfsdk:"protocol"`
	SrcCidr      types.String `tfsdk:"src_cidr"`
	DstCidr      types.String `tfsdk:"dst_cidr"`
	SrcPortRange types.String `tfsdk:"src_port_range"`
	DstPortRange types.String `tfsdk:"dst_port_range"`
	Priority     types.Int32  `tfsdk:"priority"`
	Status       types.String `tfsdk:"status"`
}

func (c *CtyunSdwanAclRules) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdwan_acl_rules"
}

func (c *CtyunSdwanAclRules) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10035786/10035852`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"acl_id": schema.StringAttribute{
				Required:    true,
				Description: "ACL ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"rules": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "规则ID",
						},
						"direction": schema.StringAttribute{
							Computed:    true,
							Description: "控制方向",
						},
						"action": schema.StringAttribute{
							Computed:    true,
							Description: "策略类型",
						},
						"protocol": schema.StringAttribute{
							Computed:    true,
							Description: "协议类型",
						},
						"src_cidr": schema.StringAttribute{
							Computed:    true,
							Description: "源网段",
						},
						"dst_cidr": schema.StringAttribute{
							Computed:    true,
							Description: "目的网段",
						},
						"src_port_range": schema.StringAttribute{
							Computed:    true,
							Description: "源端口范围",
						},
						"dst_port_range": schema.StringAttribute{
							Computed:    true,
							Description: "目的端口范围",
						},
						"priority": schema.Int32Attribute{
							Computed:    true,
							Description: "优先级",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "状态",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunSdwanAclRules) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunSdwanAclRules) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunSdwanAclRulesConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := &sdwan.SdwanGetSdwanAclRuleRequest{
		PageNo:   1,
		PageSize: 1000,
		AclID:    plan.AclID.ValueStringPointer(),
	}

	response, err := c.meta.Apis.SdkSdwanApis.SdwanGetSdwanAclRuleApi.Do(ctx, c.meta.SdkCredential, request)
	if err != nil {
		return
	} else if response.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *response.Message)
		return
	} else if response.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 正确处理API返回的规则列表
	var rules []RuleInfo
	for _, rule := range response.ReturnObj.Result {
		ruleInfo := RuleInfo{
			ID:        types.StringValue(*rule.AclRuleID),
			Direction: types.StringValue(*rule.Direction),
			Action:    types.StringValue(*rule.Action),
			Protocol:  types.StringValue(*rule.Protocol),
			SrcCidr:   types.StringValue(*rule.SrcCidr),
			DstCidr:   types.StringValue(*rule.DstCidr),
			Priority:  types.Int32Value(rule.Priority),
			Status:    types.StringValue(*rule.Status),
		}

		// 处理可选字段
		if rule.SrcPortRange != nil {
			ruleInfo.SrcPortRange = types.StringValue(*rule.SrcPortRange)
		}
		if rule.DstPortRange != nil {
			ruleInfo.DstPortRange = types.StringValue(*rule.DstPortRange)
		}

		rules = append(rules, ruleInfo)
	}

	plan.Rules = rules
	plan.ID = plan.AclID

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
