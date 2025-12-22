package oceanfs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/oceanfs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type CtyunOceanfsPermissionRule struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunOceanfsPermissionRule) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_oceanfs_permission_rule"
}

func (c *CtyunOceanfsPermissionRule) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunOceanfsPermissionRule() resource.Resource {
	return &CtyunOceanfsPermissionRule{}
}

func (c *CtyunOceanfsPermissionRule) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[permissionGroupID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunOceanfsPermissionRuleConfig
	var ID, regionID, permissionGroupID string
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &ID, &permissionGroupID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &ID, &permissionGroupID, &regionID)
		if err != nil {
			return
		}
	}
	if ID == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	if permissionGroupID == "" {
		err = fmt.Errorf("permissionGroupID不能为空")
		return
	}
	config.ID = types.StringValue(ID)
	config.RegionID = types.StringValue(regionID)
	config.PermissionGroupFuid = types.StringValue(permissionGroupID)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunOceanfsPermissionRule) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10088966/10332853",
		Attributes: map[string]schema.Attribute{
			"permission_group_id": schema.StringAttribute{
				Required:    true,
				Description: "权限组ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
				Description: "授权地址，支持更新。可填写单个 IP 或者单个网段，支持IPv4和IPv6两种网络类型。默认来访地址为*表示允许所有",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"rw_permission": schema.StringAttribute{
				Required:    true,
				Description: "读写权限类型。取值：rw：读写权限。对应IP下的计算服务可以对文件系统进行读写操作；ro：只读权限。对应IP下的计算服务可以对文件系统只有只读权限",
				Validators: []validator.String{
					stringvalidator.OneOf("ro", "rw"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"priority": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(1),
				Description: "规则优先级，支持更新。可选范围为1-400，默认值为1，即最高优先级。",
				Validators: []validator.Int32{
					int32validator.Between(1, 400),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "权限组规则id",
			},
		},
	}
}

func (c *CtyunOceanfsPermissionRule) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunOceanfsPermissionRuleConfig
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

func (c *CtyunOceanfsPermissionRule) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunOceanfsPermissionRuleConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "未找到") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunOceanfsPermissionRule) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunOceanfsPermissionRuleConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunOceanfsPermissionRuleConfig
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

func (c *CtyunOceanfsPermissionRule) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunOceanfsPermissionRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunOceanfsPermissionRule) create(ctx context.Context, config *CtyunOceanfsPermissionRuleConfig) error {
	params := &oceanfs.OceanfsNewPermissionRuleRequest{
		PermissionGroupFuid:    config.PermissionGroupFuid.ValueString(),
		RegionID:               config.RegionID.ValueString(),
		AuthAddr:               config.AuthAddr.ValueString(),
		RwPermission:           config.RwPermission.ValueString(),
		UserPermission:         "no_root_squash",
		PermissionRulePriority: config.PermissionRulePriority.ValueInt32(),
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsNewPermissionRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建权限规则失败，接口返回nil。请与研发联系确认问题原因。")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
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
		if *rule.AuthAddr == config.AuthAddr.ValueString() && rule.PermissionGroupFuid == config.PermissionGroupFuid.ValueString() {
			config.ID = types.StringValue(rule.PermissionRuleFuid)
			break
		}
	}
	return nil
}

func (c *CtyunOceanfsPermissionRule) getRuleList(ctx context.Context, config *CtyunOceanfsPermissionRuleConfig) ([]*oceanfs.OceanfsListPermissionRuleReturnObjItemResponse, error) {
	var ruleList []*oceanfs.OceanfsListPermissionRuleReturnObjItemResponse
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

func (c *CtyunOceanfsPermissionRule) getRuleDetail(ctx context.Context, config *CtyunOceanfsPermissionRuleConfig, ruleType string, pageSize int32, pageNo int32) (*oceanfs.OceanfsListPermissionRuleResponse, error) {
	params := &oceanfs.OceanfsListPermissionRuleRequest{
		RegionID:            config.RegionID.ValueString(),
		PermissionGroupFuid: config.PermissionGroupFuid.ValueString(),
		PageSize:            pageSize,
		PageNo:              pageNo,
	}
	if ruleType == "detail" {
		params.PermissionGroupFuid = config.PermissionGroupFuid.ValueString()
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsListPermissionRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询权限组列表失败，接口返回nil。请与研发联系确认问题原因。")
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

func (c *CtyunOceanfsPermissionRule) getAndMerge(ctx context.Context, config *CtyunOceanfsPermissionRuleConfig) error {
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
	returnObj := resp.ReturnObj.List[0]
	config.AuthAddr = types.StringValue(*returnObj.AuthAddr)
	config.RwPermission = types.StringValue(returnObj.RwPermission)
	config.PermissionRulePriority = types.Int32Value(returnObj.PermissionRulePriority)
	return nil
}

func (c *CtyunOceanfsPermissionRule) update(ctx context.Context, state *CtyunOceanfsPermissionRuleConfig, plan *CtyunOceanfsPermissionRuleConfig) error {
	params := &oceanfs.OceanfsModifyPermissionRuleRequest{
		PermissionRuleFuid:     state.ID.ValueString(),
		RegionID:               state.RegionID.ValueString(),
		AuthAddr:               plan.AuthAddr.ValueString(),
		RwPermission:           plan.RwPermission.ValueString(),
		UserPermission:         "no_root_squash",
		PermissionRulePriority: plan.PermissionRulePriority.ValueInt32(),
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsModifyPermissionRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新权限规则失败(id=%s)，接口返回nil。请与研发联系确认问题原因。", state.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	return nil
}

func (c *CtyunOceanfsPermissionRule) delete(ctx context.Context, config CtyunOceanfsPermissionRuleConfig) error {
	params := &oceanfs.OceanfsDeletePermissionRuleRequest{
		RegionID:           config.RegionID.ValueString(),
		PermissionRuleFuid: config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsDeletePermissionRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除权限规则失败(id=%s)，接口返回nil。请与研发联系确认问题原因。", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	return nil
}

type CtyunOceanfsPermissionRuleConfig struct {
	PermissionGroupFuid    types.String `tfsdk:"permission_group_id"`
	RegionID               types.String `tfsdk:"region_id"`
	AuthAddr               types.String `tfsdk:"auth_addr"`
	RwPermission           types.String `tfsdk:"rw_permission"`
	PermissionRulePriority types.Int32  `tfsdk:"priority"`
	ID                     types.String `tfsdk:"id"`
}
