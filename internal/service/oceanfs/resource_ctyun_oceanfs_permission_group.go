package oceanfs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/oceanfs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type CtyunOceanfsPermissionGroup struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunOceanfsPermissionGroup) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_oceanfs_permission_group"
}

func (c *CtyunOceanfsPermissionGroup) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunOceanfsPermissionGroup() resource.Resource {
	return &CtyunOceanfsPermissionGroup{}
}

func (c *CtyunOceanfsPermissionGroup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunOceanfsPermissionGroupConfig
	var ID, regionID string
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		ID = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &ID, &regionID)
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
	config.ID = types.StringValue(ID)
	config.RegionID = types.StringValue(regionID)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunOceanfsPermissionGroup) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10088966/10332853",
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
				Description: "权限组名称。同一资源池下权限组名称不能重复。文件系统实例名称只能由数字、连字符（-）、字母组成，不能以数字和连字符（-）开头、且不能以连字符（-）结尾，长度2~255字符。字母不区分大小写。",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 255),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "权限组描述信息，支持更新。支持中英文，长度为0-128字符",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 128),
					validator2.Desc(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "权限组id",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间，为UTC格式",
			},
		},
	}
}

func (c *CtyunOceanfsPermissionGroup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunOceanfsPermissionGroupConfig
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

func (c *CtyunOceanfsPermissionGroup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunOceanfsPermissionGroupConfig
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

func (c *CtyunOceanfsPermissionGroup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunOceanfsPermissionGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunOceanfsPermissionGroupConfig
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

func (c *CtyunOceanfsPermissionGroup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunOceanfsPermissionGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunOceanfsPermissionGroup) create(ctx context.Context, config *CtyunOceanfsPermissionGroupConfig) error {
	params := &oceanfs.OceanfsNewPermissionGroupRequest{
		RegionID:            config.RegionID.ValueString(),
		PermissionGroupName: config.Name.ValueString(),
		NetworkType:         "private_network",
	}
	if !config.Description.IsNull() {
		params.PermissionGroupDescription = config.Description.ValueString()
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsNewPermissionGroupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建海量文件服务的权限组失败， 接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}

	// 通过查询列表，获取id
	permissions, err := c.getPermissionGroupList(ctx, config)
	if err != nil {
		return err
	}

	for _, permission := range permissions {
		if permission.PermissionGroupName == config.Name.ValueString() {
			config.ID = types.StringValue(permission.PermissionGroupFuid)
			break
		}
	}
	if config.ID.IsUnknown() || config.ID.IsNull() {
		err = fmt.Errorf("未查询到name=%s的海量文件存储权限组！", config.Name.ValueString())
		return err
	}
	return nil
}

func (c *CtyunOceanfsPermissionGroup) getPermissionGroupList(ctx context.Context, config *CtyunOceanfsPermissionGroupConfig) ([]*oceanfs.OceanfsListPermissionGroupReturnObjListItemResponse, error) {
	var pageNo, pageSize int32
	pageNo = 1
	pageSize = 50
	resp, err := c.requestOceanfsPermissionGroupList(ctx, config, pageNo, pageSize)
	if err != nil {
		return nil, err
	}
	totalCount := resp.ReturnObj.TotalCount
	totalPageNo := pageNo
	if totalCount >= pageSize {
		totalPageNo = totalCount/pageSize + 1
	}

	var groupList []*oceanfs.OceanfsListPermissionGroupReturnObjListItemResponse
	for pageNo <= totalPageNo {
		for _, group := range resp.ReturnObj.List {
			groupList = append(groupList, group)
		}
		pageNo++
		if pageNo > totalPageNo {
			break
		}
		resp, err = c.requestOceanfsPermissionGroupList(ctx, config, pageNo, pageSize)
		if err != nil {
			return nil, err
		}
	}
	return groupList, nil
}

func (c *CtyunOceanfsPermissionGroup) requestOceanfsPermissionGroupList(ctx context.Context, config *CtyunOceanfsPermissionGroupConfig, pageNo int32, pageSize int32) (*oceanfs.OceanfsListPermissionGroupResponse, error) {
	params := &oceanfs.OceanfsListPermissionGroupRequest{
		RegionID: config.RegionID.ValueString(),
		PageSize: pageSize,
		PageNo:   pageNo,
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsListPermissionGroupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取海量文件服务的权限组列表失败， 接口返回nil，请联系研发确认问题原因！")
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	}
	return resp, nil
}

func (c *CtyunOceanfsPermissionGroup) getAndMerge(ctx context.Context, config *CtyunOceanfsPermissionGroupConfig) error {
	params := &oceanfs.OceanfsListPermissionGroupRequest{
		RegionID:            config.RegionID.ValueString(),
		PermissionGroupFuid: config.ID.ValueString(),
		PageSize:            10,
		PageNo:              1,
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsListPermissionGroupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("获取海量文件服务的权限组列表失败， 接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	if len(resp.ReturnObj.List) > 1 {
		err = fmt.Errorf("通过ID=%s查询海量文件服务的权限组列表，接口返回多个结果！", config.ID.ValueString())
		return err
	}
	if len(resp.ReturnObj.List) == 0 {
		err = fmt.Errorf("通过ID=%s查询海量文件服务的权限组列表，接口返回空结果！", config.ID.ValueString())
		return err
	}
	returnObj := resp.ReturnObj.List[0]
	config.Name = types.StringValue(returnObj.PermissionGroupName)
	config.Description = types.StringValue(returnObj.PermissionGroupDescription)
	config.CreateTime = types.StringValue(returnObj.CreateTime)
	config.UpdateTime = types.StringValue(returnObj.UpdateTime)
	return nil
}

func (c *CtyunOceanfsPermissionGroup) update(ctx context.Context, state *CtyunOceanfsPermissionGroupConfig, plan *CtyunOceanfsPermissionGroupConfig) error {
	params := &oceanfs.OceanfsModifyPermissionGroupRequest{
		PermissionGroupFuid: state.ID.ValueString(),
		RegionID:            state.RegionID.ValueString(),
		PermissionGroupName: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() && !plan.Description.Equal(state.Description) {
		params.PermissionGroupDescription = plan.Description.ValueString()
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsModifyPermissionGroupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("修改海量文件服务的权限组失败， 接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	return nil
}

func (c *CtyunOceanfsPermissionGroup) delete(ctx context.Context, config CtyunOceanfsPermissionGroupConfig) error {
	params := &oceanfs.OceanfsDeletePermissionGroupRequest{
		RegionID:            config.RegionID.ValueString(),
		PermissionGroupFuid: config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsDeletePermissionGroupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除海量文件服务的权限组失败(id=%s)， 接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	return nil
}

type CtyunOceanfsPermissionGroupConfig struct {
	RegionID    types.String `tfsdk:"region_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ID          types.String `tfsdk:"id"`
	CreateTime  types.String `tfsdk:"create_time"`
	UpdateTime  types.String `tfsdk:"update_time"`
}
