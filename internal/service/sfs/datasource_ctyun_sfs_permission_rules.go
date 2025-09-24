package sfs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/sfs"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunSfsPermissionRules{}
	_ datasource.DataSourceWithConfigure = &CtyunSfsPermissionRules{}
)

type CtyunSfsPermissionRules struct {
	meta *common.CtyunMetadata
}

func NewCtyunSfsPermissionRules() datasource.DataSource {
	return &CtyunSfsPermissionRules{}
}

func (c *CtyunSfsPermissionRules) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunSfsPermissionRules) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_sfs_permission_rules"
}

func (c *CtyunSfsPermissionRules) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "天翼云SFS（文件存储）权限组规则管理，支持权限组规则创建、更新和删除。具体文档可参考：https://www.ctyun.cn/document/10027350/10192622",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
			},
			"permission_group_fuid": schema.StringAttribute{
				Optional:    true,
				Description: "权限组Fuid，permission_group_fuid和permission_rule_fuid至少存在一个",
			},
			"permission_rule_fuid": schema.StringAttribute{
				Optional:    true,
				Description: "权限组规则Fuid，permissionGroupFuid和permissionRuleFuid至少存在一个",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页记录数目，取值范围：[1, 50]，注：默认值为10",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "页码，取值范围：正整数（≥1），注：默认值为1",
			},
			"permission_rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "权限组规则列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"permission_rule_fuid": schema.StringAttribute{
							Computed:    true,
							Description: "权限组规则ID",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间。UTC时间",
						},
						"user_id": schema.StringAttribute{
							Computed:    true,
							Description: "租户ID",
						},
						"permission_group_id": schema.StringAttribute{
							Computed:    true,
							Description: "权限组底层ID",
						},
						"permission_group_fuid": schema.StringAttribute{
							Computed:    true,
							Description: "权限组Fuid",
						},
						"permission_rule_id": schema.StringAttribute{
							Computed:    true,
							Description: "权限组规则底层ID",
						},
						"auth_addr": schema.StringAttribute{
							Computed:    true,
							Description: "授权地址，可用于区分子网及具体虚机等",
						},
						"rw_permission": schema.StringAttribute{
							Computed:    true,
							Description: "读写权限控制",
						},
						"user_permission": schema.StringAttribute{
							Computed:    true,
							Description: "用户权限",
						},
						"permission_rule_priority": schema.Int32Attribute{
							Computed:    true,
							Description: "优先级",
						},
						"permission_rule_is_default": schema.BoolAttribute{
							Computed:    true,
							Description: "是否为默认规则",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunSfsPermissionRules) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunSfsPermissionRulesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	if regionId == "" {
		err = errors.New("region id 为空")
		return
	}
	params := &sfs.SfsSfsListPermissionRuleSfsRequest{
		RegionID: regionId,
		PageSize: 10,
		PageNo:   1,
	}
	if !config.PermissionGroupFuid.IsNull() && !config.PermissionGroupFuid.IsUnknown() {
		params.PermissionGroupFuid = config.PermissionGroupFuid.ValueString()
	}
	if !config.PermissionRuleFuid.IsNull() && !config.PermissionRuleFuid.IsUnknown() {
		params.PermissionRuleFuid = config.PermissionRuleFuid.ValueString()
	}
	if params.PermissionGroupFuid == "" && params.PermissionRuleFuid == "" {
		err = fmt.Errorf("permission_group_fuid和permission_rule_fuid至少有一个不为空")
		return
	}
	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsListPermissionRuleSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询权限组列表失败，接口返回nil。请与研发联系确认问题原因。")
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析权限组规则列表
	var permissionRules []PermissionRuleModel
	for _, ruleItem := range resp.ReturnObj.List {
		var rule PermissionRuleModel
		rule.PermissionRuleFuid = types.StringValue(ruleItem.PermissionRuleFuid)
		rule.UpdateTime = types.StringValue(ruleItem.UpdateTime)
		rule.UserID = types.StringValue(ruleItem.UserID)
		rule.PermissionGroupID = types.StringValue(ruleItem.PermissionGroupID)
		rule.PermissionGroupFuid = types.StringValue(ruleItem.PermissionGroupFuid)
		rule.PermissionRuleID = types.StringValue(ruleItem.PermissionRuleID)
		rule.AuthAddr = types.StringValue(ruleItem.AuthAddr)
		rule.RwPermission = types.StringValue(ruleItem.RwPermission)
		rule.UserPermission = types.StringValue(ruleItem.UserPermission)
		rule.PermissionRulePriority = types.Int32Value(ruleItem.PermissionRulePriority)
		rule.PermissionRuleIsDefault = types.BoolValue(*ruleItem.PermissionRuleIsDefault)
		permissionRules = append(permissionRules, rule)
	}
	config.PermissionRules = permissionRules
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type PermissionRuleModel struct {
	PermissionRuleFuid      types.String `tfsdk:"permission_rule_fuid"`
	UpdateTime              types.String `tfsdk:"update_time"`
	UserID                  types.String `tfsdk:"user_id"`
	PermissionGroupID       types.String `tfsdk:"permission_group_id"`
	PermissionGroupFuid     types.String `tfsdk:"permission_group_fuid"`
	PermissionRuleID        types.String `tfsdk:"permission_rule_id"`
	AuthAddr                types.String `tfsdk:"auth_addr"`
	RwPermission            types.String `tfsdk:"rw_permission"`
	UserPermission          types.String `tfsdk:"user_permission"`
	PermissionRulePriority  types.Int32  `tfsdk:"permission_rule_priority"`
	PermissionRuleIsDefault types.Bool   `tfsdk:"permission_rule_is_default"`
}

type CtyunSfsPermissionRulesConfig struct {
	RegionID            types.String          `tfsdk:"region_id"`
	PermissionGroupFuid types.String          `tfsdk:"permission_group_fuid"`
	PermissionRuleFuid  types.String          `tfsdk:"permission_rule_fuid"`
	PageSize            types.Int32           `tfsdk:"page_size"`
	PageNo              types.Int32           `tfsdk:"page_no"`
	PermissionRules     []PermissionRuleModel `tfsdk:"permission_rules"`
}
