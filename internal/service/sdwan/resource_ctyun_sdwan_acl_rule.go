package sdwan

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/sdwan"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var (
	_ resource.Resource                = &CtyunSdwanAclRule{}
	_ resource.ResourceWithConfigure   = &CtyunSdwanAclRule{}
	_ resource.ResourceWithImportState = &CtyunSdwanAclRule{}
)

func NewCtyunSdwanAclRule() resource.Resource {
	return &CtyunSdwanAclRule{}
}

type CtyunSdwanAclRule struct {
	meta *common.CtyunMetadata
}

type CtyunSdwanAclRuleConfig struct {
	ID           types.String `tfsdk:"id"`
	AclID        types.String `tfsdk:"acl_id"`
	Direction    types.String `tfsdk:"direction"`
	Protocol     types.String `tfsdk:"protocol"`
	IpVersion    types.String `tfsdk:"ip_version"`
	DstCidr      types.String `tfsdk:"dst_cidr"`
	DstPortRange types.String `tfsdk:"dst_port_range"`
	Priority     types.Int32  `tfsdk:"priority"`
	Action       types.String `tfsdk:"action"`
	SrcCidr      types.String `tfsdk:"src_cidr"`
	SrcPortRange types.String `tfsdk:"src_port_range"`
}

func (c *CtyunSdwanAclRule) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdwan_acl_rule"
}

func (c *CtyunSdwanAclRule) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10035786/10035852`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"acl_id": schema.StringAttribute{
				Required:    true,
				Description: "ACL ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"direction": schema.StringAttribute{
				Required:    true,
				Description: "控制方向，取值范围: in(入方向), out(出方向)  支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf("in", "out"),
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "协议类型，取值范围: udp(UDP), icmp(ICMP), all(ALL), tcp(TCP)  支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf("udp", "icmp", "all", "tcp"),
				},
			},
			"ip_version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "IP协议版本，取值范围: IPv4(IPv4), IPv6(IPv6)  支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf("IPv4", "IPv6"),
				},
			},
			"dst_cidr": schema.StringAttribute{
				Required:    true,
				Description: "目的网段 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"dst_port_range": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "目的端口范围（例如1-200， -1/-1为默认值，表示1-65535） 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"priority": schema.Int32Attribute{
				Required:    true,
				Description: "优先级 支持更新",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"action": schema.StringAttribute{
				Required:    true,
				Description: "策略类型，取值范围: allow(允许), deny(拒绝)  支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf("allow", "deny"),
				},
			},
			"src_cidr": schema.StringAttribute{
				Required:    true,
				Description: "源网段 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"src_port_range": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "源端口范围（例如1-200， -1/-1为默认值，表示1-65535） 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (c *CtyunSdwanAclRule) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunSdwanAclRule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunSdwanAclRuleConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
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
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunSdwanAclRule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSdwanAclRuleConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 查询远端确认资源是否存在
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunSdwanAclRule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunSdwanAclRuleConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &plan)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunSdwanAclRule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSdwanAclRuleConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunSdwanAclRule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [aclID],[ID]"
			resp.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunSdwanAclRuleConfig
	var aclID, aclRuleID string
	err = terraform_extend.Split(req.ID, &aclID, &aclRuleID)
	if err != nil {
		return
	}

	if aclRuleID == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if aclID == "" {
		err = fmt.Errorf("aclID不能为空")
		return
	}

	config.AclID = types.StringValue(aclID)
	config.ID = types.StringValue(aclRuleID)

	// 查询远端
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (c *CtyunSdwanAclRule) create(ctx context.Context, plan *CtyunSdwanAclRuleConfig) (err error) {
	createReq := &sdwan.SdwanCreateSdwanAclRuleRequest{
		AclID: plan.AclID.ValueString(),
		AddRules: []*sdwan.SdwanCreateSdwanAclRuleAddRulesRequest{
			{
				Direction:    plan.Direction.ValueStringPointer(),
				Protocol:     plan.Protocol.ValueStringPointer(),
				IpVersion:    plan.IpVersion.ValueStringPointer(),
				DstCidr:      plan.DstCidr.ValueStringPointer(),
				DstPortRange: plan.DstPortRange.ValueStringPointer(),
				Priority:     plan.Priority.ValueInt32(),
				Action:       plan.Action.ValueStringPointer(),
				SrcCidr:      plan.SrcCidr.ValueStringPointer(),
				SrcPortRange: plan.SrcPortRange.ValueStringPointer(),
			},
		},
	}

	resp, err := c.meta.Apis.SdkSdwanApis.SdwanCreateSdwanAclRuleApi.Do(ctx, c.meta.SdkCredential, createReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}

	// 等待创建完成
	time.Sleep(3 * time.Second)
	return
}

func (c *CtyunSdwanAclRule) getAndMerge(ctx context.Context, plan *CtyunSdwanAclRuleConfig) (err error) {
	listReq := &sdwan.SdwanGetSdwanAclRuleRequest{
		PageNo:   1,
		PageSize: 1000,
		AclID:    plan.AclID.ValueStringPointer(),
	}

	resp, err := c.meta.Apis.SdkSdwanApis.SdwanGetSdwanAclRuleApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}

	// 查找对应的ACL Rule
	var found bool

	// 通过ID查找规则
	if plan.ID.ValueString() != "" {
		for _, rule := range resp.ReturnObj.Result {
			if rule.AclRuleID != nil && plan.ID == types.StringValue(*rule.AclRuleID) {
				plan.Direction = types.StringValue(*rule.Direction)
				plan.Protocol = types.StringValue(*rule.Protocol)

				//if rule.IpVersion != nil {
				//	plan.IpVersion = types.StringValue(*rule.IpVersion)
				//}
				plan.DstCidr = types.StringValue(*rule.DstCidr)
				if rule.DstPortRange != nil {
					plan.DstPortRange = types.StringValue(*rule.DstPortRange)
				}
				plan.Priority = types.Int32Value(rule.Priority)
				plan.Action = types.StringValue(*rule.Action)
				plan.SrcCidr = types.StringValue(*rule.SrcCidr)
				if rule.SrcPortRange != nil {
					plan.SrcPortRange = types.StringValue(*rule.SrcPortRange)
				}
				plan.ID = types.StringValue(*rule.AclRuleID)
				found = true
				break
			}
		}
	} else {
		for _, rule := range resp.ReturnObj.Result {
			if plan.DstCidr == types.StringValue(*rule.DstCidr) && plan.SrcCidr == types.StringValue(*rule.SrcCidr) && plan.Protocol == types.StringValue(*rule.Protocol) {
				plan.Direction = types.StringValue(*rule.Direction)
				plan.Protocol = types.StringValue(*rule.Protocol)

				//if rule.IpVersion != nil {
				//	plan.IpVersion = types.StringValue(*rule.IpVersion)
				//}
				plan.DstCidr = types.StringValue(*rule.DstCidr)
				if rule.DstPortRange != nil {
					plan.DstPortRange = types.StringValue(*rule.DstPortRange)
				}
				plan.Priority = types.Int32Value(rule.Priority)
				plan.Action = types.StringValue(*rule.Action)
				plan.SrcCidr = types.StringValue(*rule.SrcCidr)
				if rule.SrcPortRange != nil {
					plan.SrcPortRange = types.StringValue(*rule.SrcPortRange)
				}
				plan.ID = types.StringValue(*rule.AclRuleID)
				found = true
				break
			}
		}
	}

	if !found {
		return common.ResourceNotExistError
	}
	return

}

func (c *CtyunSdwanAclRule) update(ctx context.Context, plan *CtyunSdwanAclRuleConfig) (err error) {
	updateReq := &sdwan.SdwanUpdateSdwanAclRuleRequest{
		AclID: plan.AclID.ValueString(),
		UpdateRules: []*sdwan.SdwanUpdateSdwanAclRuleUpdateRulesRequest{
			{
				AclRuleID:    plan.ID.ValueStringPointer(),
				Direction:    plan.Direction.ValueStringPointer(),
				Protocol:     plan.Protocol.ValueStringPointer(),
				IpVersion:    plan.IpVersion.ValueStringPointer(),
				DstCidr:      plan.DstCidr.ValueStringPointer(),
				DstPortRange: plan.DstPortRange.ValueStringPointer(),
				Priority:     plan.Priority.ValueInt32(),
				Action:       plan.Action.ValueStringPointer(),
				SrcCidr:      plan.SrcCidr.ValueStringPointer(),
				SrcPortRange: plan.SrcPortRange.ValueStringPointer(),
			},
		},
	}

	resp, err := c.meta.Apis.SdkSdwanApis.SdwanUpdateSdwanAclRuleApi.Do(ctx, c.meta.SdkCredential, updateReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}

	// 等待更新完成
	time.Sleep(3 * time.Second)

	return
}

func (c *CtyunSdwanAclRule) delete(ctx context.Context, state *CtyunSdwanAclRuleConfig) (err error) {
	deleteReq := &sdwan.SdwanDeleteSdwanAclRuleRequest{
		AclID:       state.AclID.ValueString(),
		DeleteRules: []*string{state.ID.ValueStringPointer()},
	}

	resp, err := c.meta.Apis.SdkSdwanApis.SdwanDeleteSdwanAclRuleApi.Do(ctx, c.meta.SdkCredential, deleteReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}

	// 等待删除完成
	time.Sleep(3 * time.Second)

	return
}
