package sfs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/sfs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
)

type ctyunSfsPermissionGroup struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *ctyunSfsPermissionGroup) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_sfs_permission_group"
}

func (c *ctyunSfsPermissionGroup) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunSfsPermissionGroup() resource.Resource {
	return &ctyunSfsPermissionGroup{}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [id],[regionId]
func (c *ctyunSfsPermissionGroup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunSfsPermissionGroupConfig
	var ID, regionId string
	err := terraform_extend.Split(request.ID, &ID, &regionId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.ID = types.StringValue(ID)
	cfg.RegionID = types.StringValue(regionId)

	err = c.getAndMergeSfsPermissionGroup(ctx, &cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *ctyunSfsPermissionGroup) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027350/10192622`,
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "权限组名称。名称不能重复，长度为2-63字符，只能由数字、字母(区分大小写)、-组成，不能以数字和-开头、且不能以-结尾，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 63),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z]([a-zA-Z0-9-]*[a-zA-Z0-9])?$`),
						"权限组名称只能由数字、字母(区分大小写)、-组成，不能以数字和-开头、且不能以-结尾",
					),
				},
			},
			//"network_type": schema.StringAttribute{
			//	Required:    true,
			//	Description: "权限组网络类型：private_network（专有网络）",
			//},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "描述信息。长度为0-128字符，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 128),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "权限组Fuid",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sfs_count": schema.Int32Attribute{
				Computed:    true,
				Description: "绑定的文件系统个数",
			},
			"permission_rule_count": schema.Int32Attribute{
				Computed:    true,
				Description: "权限组规则个数",
			},
			"permission_group_is_default": schema.BoolAttribute{
				Computed:    true,
				Description: "是否为默认权限组",
			},
		},
	}
}

func (c *ctyunSfsPermissionGroup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunSfsPermissionGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.createSfsPermissionGroup(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后反查创建后的证书信息
	err = c.getAndMergeSfsPermissionGroup(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunSfsPermissionGroup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSfsPermissionGroupConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeSfsPermissionGroup(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunSfsPermissionGroup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取 plan -tf文件中配置
	var plan CtyunSfsPermissionGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunSfsPermissionGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.updateSfsPermissionGroup(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端数据，并同步本地state
	err = c.getAndMergeSfsPermissionGroup(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunSfsPermissionGroup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunSfsPermissionGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &sfs.SfsSfsDeletePermissionSfsRequest{
		RegionID:            state.RegionID.ValueString(),
		PermissionGroupFuid: state.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsDeletePermissionSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("id为%s的权限组删除失败，接口返回nil。请与研发联系确认原因。", state.ID.ValueString())
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

func (c *ctyunSfsPermissionGroup) createSfsPermissionGroup(ctx context.Context, config *CtyunSfsPermissionGroupConfig) error {
	params := &sfs.SfsSfsNewPermissionSfsRequest{
		RegionID:            config.RegionID.ValueString(),
		PermissionGroupName: config.Name.ValueString(),
		NetworkType:         "private_network",
	}
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		params.PermissionGroupDescription = config.Description.ValueString()
	}
	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsNewPermissionSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = errors.New("创建弹性文件权限组失败，接口返回nil。请与研发联系确认问题原因。")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}

	// 通过查询列表，获取id
	permissionGroupList, err := c.getSfsPermissionGroupList(ctx, config)
	for _, permissionGroup := range permissionGroupList {
		if permissionGroup.PermissionGroupName == config.Name.ValueString() {
			config.ID = types.StringValue(permissionGroup.PermissionGroupFuid)
			break
		}
	}
	return nil
}

func (c *ctyunSfsPermissionGroup) getAndMergeSfsPermissionGroup(ctx context.Context, config *CtyunSfsPermissionGroupConfig) error {
	params := &sfs.SfsSfsListPermissionSfsRequest{
		RegionID:            config.RegionID.ValueString(),
		PermissionGroupFuid: config.ID.ValueString(),
		PageSize:            50,
		PageNo:              1,
	}
	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsListPermissionSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("查询全量权限组列表失败， 接口返回nil。请与研发联系确认问题原因。 ")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}

	if len(resp.ReturnObj.List) > 1 {
		err = fmt.Errorf("弹性文件查询有误，通过id=%s查询，结果大于一条，当前条数为%d。", config.ID.ValueString(), len(resp.ReturnObj.List))
		return err
	} else if len(resp.ReturnObj.List) == 0 {
		err = fmt.Errorf("弹性文件查询有误，通过id=%s查询，结果为空。", config.ID.ValueString())
		return err
	}
	groupItem := resp.ReturnObj.List[0]
	config.Name = types.StringValue(groupItem.PermissionGroupName)
	config.Description = types.StringValue(groupItem.PermissionGroupDescription)
	config.PermissionRuleCount = types.Int32Value(groupItem.PermissionRuleCount)
	config.SfsCount = types.Int32Value(groupItem.SfsCount)
	config.PermissionGroupIsDefault = types.BoolValue(*groupItem.PermissionGroupIsDefault)
	return nil
}

func (c *ctyunSfsPermissionGroup) getSfsPermissionGroupList(ctx context.Context, config *CtyunSfsPermissionGroupConfig) ([]*sfs.SfsSfsListPermissionSfsReturnObjListResponse, error) {
	var pageNo, pageSize int32
	pageNo = 1
	pageSize = 50
	resp, err := c.requestSfsPermissionGroupList(ctx, config, pageNo, pageSize)
	if err != nil {
		return nil, err
	}
	totalCount := resp.ReturnObj.TotalCount
	totalPageNo := pageNo
	if totalCount >= pageSize {
		totalPageNo = totalCount/pageSize + 1
	}

	var groupList []*sfs.SfsSfsListPermissionSfsReturnObjListResponse
	for pageNo <= totalPageNo {
		for _, group := range resp.ReturnObj.List {
			groupList = append(groupList, group)
		}
		pageNo++
		if pageNo > totalPageNo {
			break
		}
		resp, err = c.requestSfsPermissionGroupList(ctx, config, pageNo, pageSize)
		if err != nil {
			return nil, err
		}
	}
	return groupList, nil
}

func (c *ctyunSfsPermissionGroup) requestSfsPermissionGroupList(ctx context.Context, config *CtyunSfsPermissionGroupConfig, pageNo int32, pageSize int32) (*sfs.SfsSfsListPermissionSfsResponse, error) {
	params := &sfs.SfsSfsListPermissionSfsRequest{
		RegionID: config.RegionID.ValueString(),
		PageSize: pageSize,
		PageNo:   pageNo,
	}
	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsListPermissionSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询全量权限组列表失败， 接口返回nil。请与研发联系确认问题原因。 ")
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

func (c *ctyunSfsPermissionGroup) updateSfsPermissionGroup(ctx context.Context, state *CtyunSfsPermissionGroupConfig, plan *CtyunSfsPermissionGroupConfig) error {
	params := &sfs.SfsSfsModifyPermissionSfsRequest{
		PermissionGroupFuid: state.ID.ValueString(),
		RegionID:            state.RegionID.ValueString(),
		PermissionGroupName: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		params.PermissionGroupDescription = plan.Description.ValueString()
	}
	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsModifyPermissionSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新权限组（id=%s）失败，接口返回为nil。请与研发联系确认问题原因", state.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	return nil
}

type CtyunSfsPermissionGroupConfig struct {
	RegionID types.String `tfsdk:"region_id"`
	Name     types.String `tfsdk:"name"`
	//NetworkType types.String `tfsdk:"network_type"`
	Description              types.String `tfsdk:"description"`
	ID                       types.String `tfsdk:"id"`
	SfsCount                 types.Int32  `tfsdk:"sfs_count"`
	PermissionRuleCount      types.Int32  `tfsdk:"permission_rule_count"`
	PermissionGroupIsDefault types.Bool   `tfsdk:"permission_group_is_default"`
}
