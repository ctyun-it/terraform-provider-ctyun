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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &CtyunSubnetAssociationAcl{}
	_ resource.ResourceWithConfigure   = &CtyunSubnetAssociationAcl{}
	_ resource.ResourceWithImportState = &CtyunSubnetAssociationAcl{}
)

type CtyunSubnetAssociationAcl struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunSubnetAssociationAcl) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_subnet_association_acl"
}

func (c *CtyunSubnetAssociationAcl) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunSubnetAssociationAcl() resource.Resource {
	return &CtyunSubnetAssociationAcl{}
}

func (c *CtyunSubnetAssociationAcl) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[subnetId],[aclId],[projectId],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunSubnetAssociationAclConfig
	var ID, subnetId, aclId, projectId, regionId string
	// 根据分隔符数量判断是否输入了regionID,projectId
	if strings.Count(request.ID, common.ImportSeparator) == 2 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		projectId = c.meta.GetExtraIfEmpty(projectId, common.ExtraProjectId)
		err = terraform_extend.Split(request.ID, &ID, &subnetId, &aclId)
		if err != nil {
			return
		}
	} else if strings.Count(request.ID, common.ImportSeparator) == 3 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &ID, &subnetId, &aclId, &projectId)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &ID, &subnetId, &aclId, &projectId, &regionId)
		if err != nil {
			return
		}
	}

	if ID == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}

	if subnetId == "" {
		err = fmt.Errorf("subnetId不能为空")
		return
	}
	if aclId == "" {
		err = fmt.Errorf("aclId不能为空")
		return
	}

	config.ID = types.StringValue(ID)
	config.SubnetID = types.StringValue(subnetId)
	config.AclID = types.StringValue(aclId)
	if projectId != "" {
		config.ProjectID = types.StringValue(projectId)
	}
	config.RegionID = types.StringValue(regionId)

	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunSubnetAssociationAcl) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10028591",
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
				Description: "acl_id，支持更新。acl列表可以通过data.ctyun_acls查询",
				Validators: []validator.String{
					validator2.AclID(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "subnet_id，subnet列表可能通过data.ctyun_subnets查询，不支持更新",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *CtyunSubnetAssociationAcl) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunSubnetAssociationAclConfig
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

func (c *CtyunSubnetAssociationAcl) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSubnetAssociationAclConfig
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

func (c *CtyunSubnetAssociationAcl) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunSubnetAssociationAclConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunSubnetAssociationAclConfig
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

func (c *CtyunSubnetAssociationAcl) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunSubnetAssociationAclConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunSubnetAssociationAcl) create(ctx context.Context, config *CtyunSubnetAssociationAclConfig) error {
	uuidStr := uuid.NewString()
	params := &ctvpc.CtvpcReplaceSubnetAclRequest{
		ClientToken: &uuidStr,
		RegionID:    config.RegionID.ValueString(),
		SubnetID:    config.SubnetID.ValueString(),
		AclID:       config.AclID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		params.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcReplaceSubnetAclApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("子网(id=%s)绑定acl(id=%s)失败，接口返回nil，请联系研发确认问题原因", config.SubnetID.ValueString(), config.AclID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return err
	}
	config.ID = types.StringValue(uuidStr)
	return nil
}

func (c *CtyunSubnetAssociationAcl) getAndMerge(ctx context.Context, config *CtyunSubnetAssociationAclConfig) error {
	// 查询acl详情，确认绑定情况
	resp, err := c.getAclDetail(ctx, config)
	if err != nil {
		return err
	}
	subnetIds := resp.ReturnObj.SubnetIDs
	for _, subnetId := range subnetIds {
		if config.SubnetID.ValueString() == *subnetId {
			return nil
		}
	}
	return fmt.Errorf("subnet 和 acl未绑定")
}

func (c *CtyunSubnetAssociationAcl) getAclDetail(ctx context.Context, config *CtyunSubnetAssociationAclConfig) (*ctvpc.CtvpcShowAclResponse, error) {
	params := &ctvpc.CtvpcShowAclRequest{
		RegionID: config.RegionID.ValueString(),
		AclID:    config.AclID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		params.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowAclApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取acl详情失败，接口返回nil，请联系研发确认问题原因！")
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp, nil
}

func (c *CtyunSubnetAssociationAcl) delete(ctx context.Context, config CtyunSubnetAssociationAclConfig) error {
	params := &ctvpc.CtvpcDisassociateSubnetAclRequest{
		RegionID: config.RegionID.ValueString(),
		SubnetID: config.SubnetID.ValueString(),
		AclID:    config.AclID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		params.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDisassociateSubnetAclApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("subnet(id=%s) 与acl (id=%s)解绑失败，接口返回nil，请联系研发确认问题原因！", config.SubnetID.ValueString(), config.AclID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	err = c.checkDelete(ctx, config)
	if err != nil {
		return err
	}
	return nil
}

func (c *CtyunSubnetAssociationAcl) checkDelete(ctx context.Context, config CtyunSubnetAssociationAclConfig) error {
	// 查询acl详情，确认绑定情况
	resp, err := c.getAclDetail(ctx, &config)
	if err != nil {
		return err
	}
	subnetIds := resp.ReturnObj.SubnetIDs
	for _, subnetId := range subnetIds {
		if config.SubnetID.ValueString() == *subnetId {
			return fmt.Errorf("subnet(id=%s) 与acl (id=%s)解绑失败！", config.SubnetID.ValueString(), config.AclID.ValueString())
		}
	}
	return nil
}

func (c *CtyunSubnetAssociationAcl) update(ctx context.Context, state *CtyunSubnetAssociationAclConfig, plan *CtyunSubnetAssociationAclConfig) error {
	if !plan.AclID.Equal(state.AclID) {
		err := c.create(ctx, plan)
		if err != nil {
			return err
		}
		state.AclID = plan.AclID
	}
	return nil
}

type CtyunSubnetAssociationAclConfig struct {
	RegionID  types.String `tfsdk:"region_id"`
	AclID     types.String `tfsdk:"acl_id"`
	SubnetID  types.String `tfsdk:"subnet_id"`
	ProjectID types.String `tfsdk:"project_id"`
	ID        types.String `tfsdk:"id"`
}
