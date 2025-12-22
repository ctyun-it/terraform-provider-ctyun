package sdwan

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/sdwan"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var (
	_ resource.Resource                = &CtyunSdwanAcl{}
	_ resource.ResourceWithConfigure   = &CtyunSdwanAcl{}
	_ resource.ResourceWithImportState = &CtyunSdwanAcl{}
)

func NewCtyunSdwanAcl() resource.Resource {
	return &CtyunSdwanAcl{}
}

type CtyunSdwanAcl struct {
	meta *common.CtyunMetadata
}

type CtyunSdwanAclConfig struct {
	ID        types.String `tfsdk:"id"`
	ProjectID types.String `tfsdk:"project_id"`
	Name      types.String `tfsdk:"name"`
	Rules     types.List   `tfsdk:"rules"`
}

type SdwanAclRule struct {
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

func (c *CtyunSdwanAcl) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdwan_acl"
}

func (c *CtyunSdwanAcl) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.Project(),
					stringvalidator.LengthAtLeast(1),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "访问控制名称 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"rules": schema.ListNestedAttribute{
				Required:    true,
				Description: "ACL规则列表",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"direction": schema.StringAttribute{
							Required:    true,
							Description: "控制方向，取值范围: in(入方向), out(出方向)",
							Validators: []validator.String{
								stringvalidator.OneOf("in", "out"),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"protocol": schema.StringAttribute{
							Required:    true,
							Description: "协议类型，取值范围: udp(UDP), icmp(ICMP), all(ALL), tcp(TCP)",
							Validators: []validator.String{
								stringvalidator.OneOf("udp", "icmp", "all", "tcp"),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"ip_version": schema.StringAttribute{
							Required:    true,
							Description: "IP协议版本，取值范围: IPv4(IPv4), IPv6(IPv6)",
							Validators: []validator.String{
								stringvalidator.OneOf("IPv4", "IPv6"),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"dst_cidr": schema.StringAttribute{
							Required:    true,
							Description: "目的网段",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"dst_port_range": schema.StringAttribute{
							Required:    true,
							Description: "目的端口范围（例如1-200， -1/-1为默认值，表示1-65535）",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"priority": schema.Int32Attribute{
							Required:    true,
							Description: "优先级",
							Validators: []validator.Int32{
								int32validator.Between(1, 100),
							},
							PlanModifiers: []planmodifier.Int32{
								int32planmodifier.RequiresReplace(),
							},
						},
						"action": schema.StringAttribute{
							Required:    true,
							Description: "策略类型，取值范围: allow(允许), deny(拒绝)",
							Validators: []validator.String{
								stringvalidator.OneOf("allow", "deny"),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"src_cidr": schema.StringAttribute{
							Required:    true,
							Description: "源网段",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"src_port_range": schema.StringAttribute{
							Required:    true,
							Description: "源端口范围（例如1-200， -1/-1为默认值，表示1-65535）",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
		},
	}
}

func (c *CtyunSdwanAcl) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunSdwanAcl) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunSdwanAclConfig
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

func (c *CtyunSdwanAcl) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSdwanAclConfig
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

func (c *CtyunSdwanAcl) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// SD-WAN ACL 不支持更新操作，直接返回
	var plan CtyunSdwanAclConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state CtyunSdwanAclConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &plan, &state)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *CtyunSdwanAcl) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSdwanAclConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunSdwanAcl) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID]"
			resp.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunSdwanAclConfig

	ID := req.ID
	if ID == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}

	config.ID = types.StringValue(ID)

	// 查询远端
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}

	// 导入时不设置 rules 字段，保持其为未知状态
	config.Rules = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"direction":      types.StringType,
			"protocol":       types.StringType,
			"ip_version":     types.StringType,
			"dst_cidr":       types.StringType,
			"dst_port_range": types.StringType,
			"priority":       types.Int32Type,
			"action":         types.StringType,
			"src_cidr":       types.StringType,
			"src_port_range": types.StringType,
		},
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (c *CtyunSdwanAcl) create(ctx context.Context, plan *CtyunSdwanAclConfig) (err error) {
	var rules []*sdwan.SdwanCreateSdwanAclRulesRequest

	// 遍历规则列表
	for _, ruleValue := range plan.Rules.Elements() {
		// 将规则值转换为对象
		ruleObj := ruleValue.(types.Object)
		attributes := ruleObj.Attributes()

		rules = append(rules, &sdwan.SdwanCreateSdwanAclRulesRequest{
			Direction:    attributes["direction"].(types.String).ValueString(),
			Protocol:     attributes["protocol"].(types.String).ValueString(),
			IpVersion:    attributes["ip_version"].(types.String).ValueString(),
			DstCidr:      attributes["dst_cidr"].(types.String).ValueString(),
			DstPortRange: attributes["dst_port_range"].(types.String).ValueString(),
			Priority:     attributes["priority"].(types.Int32).ValueInt32(),
			Action:       attributes["action"].(types.String).ValueString(),
			SrcCidr:      attributes["src_cidr"].(types.String).ValueString(),
			SrcPortRange: attributes["src_port_range"].(types.String).ValueString(),
		})
	}

	createReq := &sdwan.SdwanCreateSdwanAclRequest{
		AclName: plan.Name.ValueString(),
		Rules:   rules,
	}
	if plan.ProjectID.IsNull() || plan.ProjectID.IsUnknown() || plan.ProjectID.ValueString() == "" {
		createReq.ProjectID = "0"
	} else {
		createReq.ProjectID = plan.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkSdwanApis.SdwanCreateSdwanAclApi.Do(ctx, c.meta.SdkCredential, createReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}

	// 等待创建完成
	time.Sleep(3 * time.Second)
	return
}

func (c *CtyunSdwanAcl) getAndMerge(ctx context.Context, plan *CtyunSdwanAclConfig) (err error) {
	listReq := &sdwan.SdwanGetSdwanAclRequest{
		PageNo:   1,
		PageSize: 1000,
	}
	if plan.ID.ValueStringPointer() != nil {
		listReq.AclID = plan.ID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkSdwanApis.SdwanGetSdwanAclApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}

	// 查找对应的ACL
	var found bool
	if plan.ID.ValueString() != "" {
		for _, aclItem := range resp.ReturnObj.Result {
			if aclItem.AclID != nil && *aclItem.AclID == plan.ID.ValueString() {
				plan.Name = types.StringValue(*aclItem.Name)
				found = true
				break
			}
		}

		if !found {
			return common.ResourceNotExistError
		}
		return
	} else if plan.Name.ValueString() != "" {
		for _, aclItem := range resp.ReturnObj.Result {
			if aclItem.Name != nil && *aclItem.Name == plan.Name.ValueString() {
				plan.ID = types.StringValue(*aclItem.AclID)
				found = true
				break
			}
		}
		if !found {
			return common.ResourceNotExistError
		}
		return
	}

	return
}

func (c *CtyunSdwanAcl) delete(ctx context.Context, state *CtyunSdwanAclConfig) (err error) {
	deleteReq := &sdwan.SdwanDeleteSdwanAclRequest{
		AclID: state.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkSdwanApis.SdwanDeleteSdwanAclApi.Do(ctx, c.meta.SdkCredential, deleteReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}

	// 等待删除完成
	time.Sleep(3 * time.Second)

	return
}

func (c *CtyunSdwanAcl) update(ctx context.Context, plan *CtyunSdwanAclConfig, state *CtyunSdwanAclConfig) (err error) {
	if plan.Name.ValueString() != state.Name.ValueString() {
		updateReq := &sdwan.SdwanUpdateSdwanAclRequest{
			AclID:   plan.ID.ValueString(),
			AclName: plan.Name.ValueString(),
		}

		resp, err := c.meta.Apis.SdkSdwanApis.SdwanUpdateSdwanAclApi.Do(ctx, c.meta.SdkCredential, updateReq)
		if err != nil {
			return err
		} else if resp.StatusCode != common.NormalStatusCode {
			return fmt.Errorf("API return error. Message: %s", *resp.Message)
		}
	} else if !plan.Rules.Equal(state.Rules) {
		// 规则发生变化，需要找出新增、删除和更新的规则
		planRules, err := c.extractRulesFromList(ctx, plan.Rules)
		if err != nil {
			return err
		}

		stateRules, err := c.extractRulesFromList(ctx, state.Rules)
		if err != nil {
			return err
		}

		// 找出需要新增、删除和更新的规则
		toAdd, toDelete, toUpdate := c.diffRules(planRules, stateRules)

		// 执行删除操作
		if len(toDelete) > 0 {
			err = c.deleteRules(ctx, plan.ID.ValueString(), toDelete)
			if err != nil {
				return err
			}
		}

		// 执行新增操作
		if len(toAdd) > 0 {
			err = c.addRules(ctx, plan.ID.ValueString(), toAdd)
			if err != nil {
				return err
			}
		}

		// 执行更新操作
		if len(toUpdate) > 0 {
			err = c.updateRules(ctx, plan.ID.ValueString(), toUpdate)
			if err != nil {
				return err
			}
		}
	}

	// 等待更新完成
	time.Sleep(3 * time.Second)
	return
}

// extractRulesFromList 从types.List中提取规则
func (c *CtyunSdwanAcl) extractRulesFromList(ctx context.Context, rulesList types.List) ([]*SdwanAclRule, error) {
	var rules []*SdwanAclRule

	if rulesList.IsNull() || rulesList.IsUnknown() {
		return rules, nil
	}

	for _, ruleValue := range rulesList.Elements() {
		// 将规则值转换为对象
		ruleObj, ok := ruleValue.(types.Object)
		if !ok {
			return nil, fmt.Errorf("failed to convert rule value to object")
		}

		attributes := ruleObj.Attributes()

		rule := &SdwanAclRule{
			Direction:    attributes["direction"].(types.String),
			Protocol:     attributes["protocol"].(types.String),
			IpVersion:    attributes["ip_version"].(types.String),
			DstCidr:      attributes["dst_cidr"].(types.String),
			DstPortRange: attributes["dst_port_range"].(types.String),
			Priority:     attributes["priority"].(types.Int32),
			Action:       attributes["action"].(types.String),
			SrcCidr:      attributes["src_cidr"].(types.String),
			SrcPortRange: attributes["src_port_range"].(types.String),
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

// diffRules 比较计划和状态中的规则，找出需要新增、删除和更新的规则
func (c *CtyunSdwanAcl) diffRules(planRules, stateRules []*SdwanAclRule) ([]*SdwanAclRule, []*SdwanAclRule, []*SdwanAclRule) {
	var toAdd, toDelete, toUpdate []*SdwanAclRule

	// 创建state规则的映射，用于快速查找
	stateRuleMap := make(map[string]*SdwanAclRule)
	for _, rule := range stateRules {
		key := c.ruleKey(rule)
		stateRuleMap[key] = rule
	}

	// 创建plan规则的映射，用于快速查找
	planRuleMap := make(map[string]*SdwanAclRule)
	for _, rule := range planRules {
		key := c.ruleKey(rule)
		planRuleMap[key] = rule
	}

	// 找出需要新增的规则（存在于plan但不存在于state）
	for key, rule := range planRuleMap {
		if _, exists := stateRuleMap[key]; !exists {
			toAdd = append(toAdd, rule)
		}
	}

	// 找出需要删除的规则（存在于state但不存在于plan）
	for key, rule := range stateRuleMap {
		if _, exists := planRuleMap[key]; !exists {
			toDelete = append(toDelete, rule)
		}
	}

	// 找出需要更新的规则（同时存在于plan和state，但内容不同）
	for key, planRule := range planRuleMap {
		if stateRule, exists := stateRuleMap[key]; exists {
			// 如果规则存在但内容不同，则需要更新
			if !c.rulesEqual(planRule, stateRule) {
				toUpdate = append(toUpdate, planRule)
			}
		}
	}

	return toAdd, toDelete, toUpdate
}

// ruleKey 生成规则的唯一键
func (c *CtyunSdwanAcl) ruleKey(rule *SdwanAclRule) string {
	// 使用方向、源CIDR、目的CIDR、协议作为规则的唯一标识
	return fmt.Sprintf("%s_%s_%s_%s",
		rule.Direction.ValueString(),
		rule.SrcCidr.ValueString(),
		rule.DstCidr.ValueString(),
		rule.Protocol.ValueString())
}

// rulesEqual 比较两个规则是否相等
func (c *CtyunSdwanAcl) rulesEqual(rule1, rule2 *SdwanAclRule) bool {
	return rule1.Direction.Equal(rule2.Direction) &&
		rule1.Protocol.Equal(rule2.Protocol) &&
		rule1.IpVersion.Equal(rule2.IpVersion) &&
		rule1.DstCidr.Equal(rule2.DstCidr) &&
		rule1.DstPortRange.Equal(rule2.DstPortRange) &&
		rule1.Priority.Equal(rule2.Priority) &&
		rule1.Action.Equal(rule2.Action) &&
		rule1.SrcCidr.Equal(rule2.SrcCidr) &&
		rule1.SrcPortRange.Equal(rule2.SrcPortRange)
}

// deleteRules 删除规则
func (c *CtyunSdwanAcl) deleteRules(ctx context.Context, aclID string, rulesToDelete []*SdwanAclRule) error {
	// 首先获取现有的规则以获得规则ID
	listReq := &sdwan.SdwanGetSdwanAclRuleRequest{
		PageNo:   1,
		PageSize: 1000,
		AclID:    &aclID,
	}

	resp, err := c.meta.Apis.SdkSdwanApis.SdwanGetSdwanAclRuleApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}

	// 创建规则映射以便快速查找规则ID
	ruleIDMap := make(map[string]string)
	for _, rule := range resp.ReturnObj.Result {
		if rule.AclRuleID != nil && rule.Direction != nil && rule.SrcCidr != nil && rule.DstCidr != nil && rule.Protocol != nil {
			key := fmt.Sprintf("%s_%s_%s_%s", *rule.Direction, *rule.SrcCidr, *rule.DstCidr, *rule.Protocol)
			ruleIDMap[key] = *rule.AclRuleID
		}
	}

	// 收集要删除的规则ID
	var ruleIDsToDelete []*string
	for _, rule := range rulesToDelete {
		key := c.ruleKey(rule)
		if ruleID, exists := ruleIDMap[key]; exists {
			ruleIDsToDelete = append(ruleIDsToDelete, &ruleID)
		}
	}

	// 执行删除操作
	if len(ruleIDsToDelete) > 0 {
		deleteReq := &sdwan.SdwanDeleteSdwanAclRuleRequest{
			AclID:       aclID,
			DeleteRules: ruleIDsToDelete,
		}

		deleteResp, err := c.meta.Apis.SdkSdwanApis.SdwanDeleteSdwanAclRuleApi.Do(ctx, c.meta.SdkCredential, deleteReq)
		if err != nil {
			return err
		} else if deleteResp.StatusCode != common.NormalStatusCode {
			return fmt.Errorf("API return error. Message: %s", *deleteResp.Message)
		}
	}

	return nil
}

// addRules 添加新规则
func (c *CtyunSdwanAcl) addRules(ctx context.Context, aclID string, rulesToAdd []*SdwanAclRule) error {
	var addRules []*sdwan.SdwanCreateSdwanAclRuleAddRulesRequest

	for _, rule := range rulesToAdd {
		addRules = append(addRules, &sdwan.SdwanCreateSdwanAclRuleAddRulesRequest{
			Direction:    rule.Direction.ValueStringPointer(),
			Protocol:     rule.Protocol.ValueStringPointer(),
			IpVersion:    rule.IpVersion.ValueStringPointer(),
			DstCidr:      rule.DstCidr.ValueStringPointer(),
			DstPortRange: rule.DstPortRange.ValueStringPointer(),
			Priority:     rule.Priority.ValueInt32(),
			Action:       rule.Action.ValueStringPointer(),
			SrcCidr:      rule.SrcCidr.ValueStringPointer(),
			SrcPortRange: rule.SrcPortRange.ValueStringPointer(),
		})
	}

	if len(addRules) > 0 {
		createReq := &sdwan.SdwanCreateSdwanAclRuleRequest{
			AclID:    aclID,
			AddRules: addRules,
		}

		resp, err := c.meta.Apis.SdkSdwanApis.SdwanCreateSdwanAclRuleApi.Do(ctx, c.meta.SdkCredential, createReq)
		if err != nil {
			return err
		} else if resp.StatusCode != common.NormalStatusCode {
			return fmt.Errorf("API return error. Message: %s", *resp.Message)
		}
	}

	return nil
}

// updateRules 更新现有规则
func (c *CtyunSdwanAcl) updateRules(ctx context.Context, aclID string, rulesToUpdate []*SdwanAclRule) error {
	// 首先获取现有的规则以获得规则ID
	listReq := &sdwan.SdwanGetSdwanAclRuleRequest{
		PageNo:   1,
		PageSize: 1000,
		AclID:    &aclID,
	}

	resp, err := c.meta.Apis.SdkSdwanApis.SdwanGetSdwanAclRuleApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}

	// 创建规则映射以便快速查找规则ID
	ruleIDMap := make(map[string]string)
	for _, rule := range resp.ReturnObj.Result {
		if rule.AclRuleID != nil && rule.Direction != nil && rule.SrcCidr != nil && rule.DstCidr != nil && rule.Protocol != nil {
			key := fmt.Sprintf("%s_%s_%s_%s", *rule.Direction, *rule.SrcCidr, *rule.DstCidr, *rule.Protocol)
			ruleIDMap[key] = *rule.AclRuleID
		}
	}

	// 准备更新规则列表
	var updateRules []*sdwan.SdwanUpdateSdwanAclRuleUpdateRulesRequest

	for _, rule := range rulesToUpdate {
		key := c.ruleKey(rule)
		if ruleID, exists := ruleIDMap[key]; exists {
			updateRules = append(updateRules, &sdwan.SdwanUpdateSdwanAclRuleUpdateRulesRequest{
				AclRuleID:    &ruleID,
				Direction:    rule.Direction.ValueStringPointer(),
				Protocol:     rule.Protocol.ValueStringPointer(),
				IpVersion:    rule.IpVersion.ValueStringPointer(),
				DstCidr:      rule.DstCidr.ValueStringPointer(),
				DstPortRange: rule.DstPortRange.ValueStringPointer(),
				Priority:     rule.Priority.ValueInt32(),
				Action:       rule.Action.ValueStringPointer(),
				SrcCidr:      rule.SrcCidr.ValueStringPointer(),
				SrcPortRange: rule.SrcPortRange.ValueStringPointer(),
			})
		}
	}

	// 执行更新操作
	if len(updateRules) > 0 {
		updateReq := &sdwan.SdwanUpdateSdwanAclRuleRequest{
			AclID:       aclID,
			UpdateRules: updateRules,
		}

		updateResp, err := c.meta.Apis.SdkSdwanApis.SdwanUpdateSdwanAclRuleApi.Do(ctx, c.meta.SdkCredential, updateReq)
		if err != nil {
			return err
		} else if updateResp.StatusCode != common.NormalStatusCode {
			return fmt.Errorf("API return error. Message: %s", *updateResp.Message)
		}
	}

	return nil
}
