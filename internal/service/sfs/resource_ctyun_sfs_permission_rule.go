package sfs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/sfs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ctyunSfsPermissionGroupRule struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *ctyunSfsPermissionGroupRule) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_sfs_permission_rule"
}

func (c *ctyunSfsPermissionGroupRule) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunSfsPermissionGroupRule() resource.Resource {
	return &ctyunSfsPermissionGroupRule{}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [id],[regionId]
func (c *ctyunSfsPermissionGroupRule) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunSfsPermissionGroupRuleConfig
	var ID, regionId, permissionGroupFuid string
	err := terraform_extend.Split(request.ID, &ID, &regionId, &permissionGroupFuid)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.ID = types.StringValue(ID)
	cfg.RegionID = types.StringValue(regionId)
	cfg.PermissionGroupFuid = types.StringValue(permissionGroupFuid)

	err = c.getAndMergeSfsPermissionGroupRule(ctx, &cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *ctyunSfsPermissionGroupRule) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027350/10192622`,
		Attributes: map[string]schema.Attribute{
			"permission_group_fuid": schema.StringAttribute{
				Required:    true,
				Description: "权限组FUID标识",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
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
			"auth_addr": schema.StringAttribute{
				Required:    true,
				Description: "授权地址。支持IPv4和IPv6两种网络类型，可填写单个IP或者单个网段。同一权限组内，授权地址不能重复格式。ipv4格式为： 192.168.1.0/24，ipv6格式为：0000:0000:0000:0000:0000:0000:0000:0000/0。支持更新",
				Validators: []validator.String{
					validator2.AuthAddr(),
				},
			},
			"rw_permission": schema.StringAttribute{
				Required:    true,
				Description: "读写权限，可选值: 'rw' (读写), 'ro' (只读)，支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"rw", "ro"}...),
				},
			},
			//"user_permission": schema.StringAttribute{
			//	Required:    true,
			//	Description: "用户权限，可选值: 'no_root_squash', 'root_squash'",
			//	Validators: []validator.String{
			//		stringvalidator.OneOf([]string{"no_root_squash", "root_squash"}...),
			//	},
			//},
			"permission_rule_priority": schema.Int32Attribute{
				Required:    true,
				Description: "规则优先级(数值越小优先级越高),有效范围为1-400。当同一个权限组内单个 IP 与网段中包含的 IP 的权限有冲突时，会生效优先级高的规则。注：优先级不可重复，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(1, 400),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "权限组规则唯一标识",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间。UTC时间",
			},
		},
	}
}

func (c *ctyunSfsPermissionGroupRule) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunSfsPermissionGroupRuleConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.createSfsPermissionRule(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后反查创建的信息
	err = c.getAndMergeSfsPermissionGroupRule(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunSfsPermissionGroupRule) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSfsPermissionGroupRuleConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeSfsPermissionGroupRule(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunSfsPermissionGroupRule) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取 plan -tf文件中配置
	var plan CtyunSfsPermissionGroupRuleConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunSfsPermissionGroupRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.updateSfsPermissionRule(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端数据，并同步本地state
	err = c.getAndMergeSfsPermissionGroupRule(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunSfsPermissionGroupRule) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunSfsPermissionGroupRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &sfs.SfsSfsDeletePermissionRuleSfsRequest{
		RegionID:           state.RegionID.ValueString(),
		PermissionRuleFuid: state.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsDeletePermissionRuleSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("删除弹性文件服务权限组规格id=%s失败，接口返回nil。请与研发联系确认问题原因。", state.ID.ValueString())
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

func (c *ctyunSfsPermissionGroupRule) createSfsPermissionRule(ctx context.Context, config *CtyunSfsPermissionGroupRuleConfig) error {
	params := &sfs.SfsSfsNewPermissionRuleSfsRequest{
		PermissionGroupFuid:    config.PermissionGroupFuid.ValueString(),
		RegionID:               config.RegionID.ValueString(),
		AuthAddr:               config.AuthAddr.ValueString(),
		RwPermission:           config.RwPermission.ValueString(),
		UserPermission:         business.SfsPermissionGroupRuleUserPermissionNoRootSquash,
		PermissionRulePriority: config.PermissionRulePriority.ValueInt32(),
	}
	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsNewPermissionRuleSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建弹性文件服务（sfs）权限组规则失败，接口返回nil。请与研发联系确认问题原因。")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	// 通过查询权限组列表，获取权限组规则id
	ruleList, err := c.getRuleList(ctx, config)
	if err != nil {
		return err
	}
	if len(ruleList) <= 0 {
		err = fmt.Errorf("未查询到权限组id=%s，权限组规则id=%s详情", config.PermissionGroupFuid.ValueString(), config.ID.ValueString())
		return err
	}
	for _, rule := range ruleList {
		if rule.AuthAddr == config.AuthAddr.ValueString() && rule.PermissionGroupFuid == config.PermissionGroupFuid.ValueString() {
			config.ID = types.StringValue(rule.PermissionRuleFuid)
			break
		}
	}
	return nil
}

func (c *ctyunSfsPermissionGroupRule) getAndMergeSfsPermissionGroupRule(ctx context.Context, config *CtyunSfsPermissionGroupRuleConfig) error {
	resp, err := c.getRuleDetail(ctx, config, "detail", 10, 1)
	if err != nil {
		return err
	}
	if len(resp.ReturnObj.List) > 1 {
		err = fmt.Errorf("查询权限组id=%s，权限组规则id=%s详情，返回信息条数>1。", config.PermissionGroupFuid.ValueString(), config.ID.ValueString())
		return err
	}
	if len(resp.ReturnObj.List) == 0 {
		err = fmt.Errorf("未查询到权限组id=%s，权限组规则id=%s详情", config.PermissionGroupFuid.ValueString(), config.ID.ValueString())
		return err
	}
	rule := resp.ReturnObj.List[0]
	config.AuthAddr = types.StringValue(rule.AuthAddr)
	config.RwPermission = types.StringValue(rule.RwPermission)
	config.PermissionRulePriority = types.Int32Value(rule.PermissionRulePriority)
	config.UpdateTime = types.StringValue(rule.UpdateTime)
	return nil
}

func (c *ctyunSfsPermissionGroupRule) getRuleDetail(ctx context.Context, config *CtyunSfsPermissionGroupRuleConfig, ruleType string, pageSize, pageNo int32) (*sfs.SfsSfsListPermissionRuleSfsResponse, error) {
	params := sfs.SfsSfsListPermissionRuleSfsRequest{
		RegionID:            config.RegionID.ValueString(),
		PermissionGroupFuid: config.PermissionGroupFuid.ValueString(),
		PageSize:            pageSize,
		PageNo:              pageNo,
	}
	if ruleType == "detail" {
		params.PermissionRuleFuid = config.ID.ValueString()
	}

	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsListPermissionRuleSfsApi.Do(ctx, c.meta.SdkCredential, &params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询权限组id=%s列表的详情失败，接口返回nil。请与研发联系确认问题原因。", config.PermissionGroupFuid.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp, nil
}

func (c *ctyunSfsPermissionGroupRule) updateSfsPermissionRule(ctx context.Context, state *CtyunSfsPermissionGroupRuleConfig, plan *CtyunSfsPermissionGroupRuleConfig) error {
	params := &sfs.SfsSfsModifyPermissionRuleSfsRequest{
		PermissionRuleFuid:     state.ID.ValueString(),
		RegionID:               state.RegionID.ValueString(),
		AuthAddr:               plan.AuthAddr.ValueString(),
		RwPermission:           plan.RwPermission.ValueString(),
		UserPermission:         business.SfsPermissionGroupRuleUserPermissionNoRootSquash,
		PermissionRulePriority: plan.PermissionRulePriority.ValueInt32(),
	}
	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsModifyPermissionRuleSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新弹性文件服务权限组规则（id=%s）失败，接口返回nil。请与研发联系确认问题原因。", state.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	return nil
}

func (c *ctyunSfsPermissionGroupRule) getRuleList(ctx context.Context, config *CtyunSfsPermissionGroupRuleConfig) ([]*sfs.SfsSfsListPermissionRuleSfsReturnObjListResponse, error) {
	var ruleList []*sfs.SfsSfsListPermissionRuleSfsReturnObjListResponse
	var pageSize, pageNo, totalPageNo int32
	pageSize = 50
	pageNo = 1
	totalPageNo = 1
	resp, err := c.getRuleDetail(ctx, config, "list", pageSize, pageNo)
	if err != nil {
		return nil, err
	}
	totalCount := resp.ReturnObj.TotalCount
	if pageSize < totalCount {
		totalPageNo = totalCount/pageSize + 1
	}
	for pageNo <= totalPageNo {
		for _, rule := range resp.ReturnObj.List {
			ruleList = append(ruleList, rule)
		}
		pageNo++
		if pageNo > totalPageNo {
			break
		}
		resp, err = c.getRuleDetail(ctx, config, "list", pageSize, pageNo)
		if err != nil {
			return nil, err
		}
	}
	return ruleList, nil
}

type CtyunSfsPermissionGroupRuleConfig struct {
	PermissionGroupFuid types.String `tfsdk:"permission_group_fuid"`
	RegionID            types.String `tfsdk:"region_id"`
	AuthAddr            types.String `tfsdk:"auth_addr"`
	RwPermission        types.String `tfsdk:"rw_permission"`
	//UserPermission         types.String `hcl:"user_permission"`
	PermissionRulePriority types.Int32  `tfsdk:"permission_rule_priority"`
	ID                     types.String `tfsdk:"id"`
	UpdateTime             types.String `tfsdk:"update_time"`
}
