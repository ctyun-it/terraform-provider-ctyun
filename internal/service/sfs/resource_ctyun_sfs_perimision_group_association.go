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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

type ctyunSfsPermissionGroupAssociation struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
	vpcService    *business.VpcService
}

func (c *ctyunSfsPermissionGroupAssociation) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_sfs_permission_group_association"
}

func (c *ctyunSfsPermissionGroupAssociation) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.vpcService = business.NewVpcService(c.meta)
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunSfsPermissionGroupAssociation() resource.Resource {
	return &ctyunSfsPermissionGroupAssociation{}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [id],[regionId]
func (c *ctyunSfsPermissionGroupAssociation) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunSfsPermissionGroupAssociationConfig
	var vpcID, regionId, sfsUid string
	err = terraform_extend.Split(request.ID, &vpcID, &sfsUid, &regionId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	cfg.RegionID = types.StringValue(regionId)
	cfg.VpcID = types.StringValue(vpcID)
	cfg.SfsUID = types.StringValue(sfsUid)

	err = c.getAndMergeSfsPermissionGroupAssociation(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *ctyunSfsPermissionGroupAssociation) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027350/10192625**`,
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
			"permission_group_fuid": schema.StringAttribute{
				Required:    true,
				Description: "权限组ID，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"sfs_uid": schema.StringAttribute{
				Required:    true,
				Description: "弹性文件系统唯一ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "vpcID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"vpc_name": schema.StringAttribute{
				Computed:    true,
				Description: "vpc名称",
			},
			"vpc_cidr": schema.StringAttribute{
				Computed:    true,
				Description: "vpc cidr",
			},
			"permission_group_name": schema.StringAttribute{
				Computed:    true,
				Description: "权限组名称",
			},
			"permission_group_description": schema.StringAttribute{
				Computed:    true,
				Description: "权限组描述",
			},
			"permission_group_is_default": schema.BoolAttribute{
				Computed:    true,
				Description: "是否为默认权限组",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID",
			},
		},
	}
}

func (c *ctyunSfsPermissionGroupAssociation) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunSfsPermissionGroupAssociationConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkBeforeBind(ctx, plan)
	if err != nil {
		return
	}
	err = c.bindSfsPermissionGroupAssociation(ctx, plan)
	if err != nil {
		return
	}
	err = c.getAndMergeSfsPermissionGroupAssociation(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunSfsPermissionGroupAssociation) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSfsPermissionGroupAssociationConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergeSfsPermissionGroupAssociation(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunSfsPermissionGroupAssociation) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取 plan -tf文件中配置
	var plan CtyunSfsPermissionGroupAssociationConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunSfsPermissionGroupAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.updateSfsPermissionGroupAssociation(ctx, plan, state)
	if err != nil {
		return
	}
	// 更新远端数据，并同步本地state
	err = c.getAndMergeSfsPermissionGroupAssociation(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunSfsPermissionGroupAssociation) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunSfsPermissionGroupAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &sfs.SfsSfsUnbindVpcSfsRequest{
		PermissionGroupFuid: state.PermissionGroupFuid.ValueString(),
		RegionID:            state.RegionID.ValueString(),
		SfsUID:              state.SfsUID.ValueString(),
		VpcID:               state.VpcID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsUnbindVpcSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("解绑弹性文件服务（id=%s）的权限组（id=%s）失败，接口返回为nil。请与研发联系确认问题原因", state.SfsUID.ValueString(), state.PermissionGroupFuid.ValueString())
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

func (c *ctyunSfsPermissionGroupAssociation) checkBeforeBind(ctx context.Context, plan CtyunSfsPermissionGroupAssociationConfig) error {
	vpc, regionID := plan.VpcID.ValueString(), plan.RegionID.ValueString()
	subnets, err := c.vpcService.GetVpcSubnet(ctx, vpc, regionID, "")
	if err != nil {
		return err
	}
	if len(subnets) == 0 {
		return fmt.Errorf("%s 必须有子网", vpc)
	}
	return nil
}

func (c *ctyunSfsPermissionGroupAssociation) bindSfsPermissionGroupAssociation(ctx context.Context, config CtyunSfsPermissionGroupAssociationConfig) error {
	params := &sfs.SfsSfsBindVpcSfsRequest{
		PermissionGroupFuid: config.PermissionGroupFuid.ValueString(),
		RegionID:            config.RegionID.ValueString(),
		SfsUID:              config.SfsUID.ValueString(),
		VpcID:               config.VpcID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsBindVpcSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("弹性文件系统（id=%s）绑定权限组（id=%s）失败，接口返回nil。请与研发联系确认问题原因。", config.SfsUID.ValueString(), config.PermissionGroupFuid.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	// 绑定后需要轮询下
	err = c.bindLoop(ctx, config)
	if err != nil {
		return err
	}
	return nil
}

func (c *ctyunSfsPermissionGroupAssociation) getAndMergeSfsPermissionGroupAssociation(ctx context.Context, config *CtyunSfsPermissionGroupAssociationConfig) error {

	resp, err := c.requestSfsVpcList(ctx, *config)
	if err != nil {
		return err
	}

	returnObj := resp.ReturnObj

	for _, association := range returnObj.List {
		if association.VpcID == config.VpcID.ValueString() {
			config.PermissionGroupIsDefault = types.BoolValue(*association.PermissionGroupIsDefault)
			config.PermissionGroupDescription = types.StringValue(association.PermissionGroupDescription)
			config.PermissionGroupName = types.StringValue(association.PermissionGroupName)
			config.PermissionGroupFuid = types.StringValue(association.PermissionGroupFuid)
			config.VpcCidr = types.StringValue(association.VpcCidr)
			config.VpcName = types.StringValue(association.VpcName)
			break
		}
	}
	config.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", config.VpcID.ValueString(), config.SfsUID.ValueString(), config.RegionID.ValueString()))
	return nil
}

func (c *ctyunSfsPermissionGroupAssociation) requestSfsVpcList(ctx context.Context, config CtyunSfsPermissionGroupAssociationConfig) (*sfs.SfsSfsListVpcSfsResponse, error) {
	params := &sfs.SfsSfsListVpcSfsRequest{
		RegionID: config.RegionID.ValueString(),
		SfsUID:   config.SfsUID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsListVpcSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询文件系统（id=%s）下的vpc列表、vpc绑定的权限组列表表失败， 接口返回nil。请与研发联系确认问题原因。 ", config.SfsUID.ValueString())
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

func (c *ctyunSfsPermissionGroupAssociation) updateSfsPermissionGroupAssociation(ctx context.Context, plan, state CtyunSfsPermissionGroupAssociationConfig) error {
	if plan.PermissionGroupFuid.Equal(state.PermissionGroupFuid) {
		return nil
	}

	params := &sfs.SfsSfsChangeVpcSfsRequest{
		PermissionGroupFuid: plan.PermissionGroupFuid.ValueString(),
		RegionID:            state.RegionID.ValueString(),
		SfsUID:              state.SfsUID.ValueString(),
		VpcID:               state.VpcID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsChangeVpcSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("换绑弹性文件服务（id=%s）的权限组（id=%s）失败，接口返回为nil。请与研发联系确认问题原因", plan.SfsUID.ValueString(), plan.PermissionGroupFuid.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	// 轮询确认更新完成
	err = c.bindLoop(ctx, plan)
	if err != nil {
		return err
	}
	return nil
}

func (c *ctyunSfsPermissionGroupAssociation) bindLoop(ctx context.Context, config CtyunSfsPermissionGroupAssociationConfig) error {
	var err error
	retryer, err := business.NewRetryer(time.Second*10, 60)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.requestSfsVpcList(ctx, config)
			if err2 != nil {
				err = err2
				return false
			}
			returnObj := resp.ReturnObj

			for _, association := range returnObj.List {
				if association.VpcID == config.VpcID.ValueString() && association.PermissionGroupFuid == config.PermissionGroupFuid.ValueString() {
					return false
				}
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return fmt.Errorf("轮询已达最大次数，sfs(id=%s)绑定vpc(%s)仍未创建成功！", config.SfsUID.ValueString(), config.VpcID.ValueString())
	}
	return err
}

type CtyunSfsPermissionGroupAssociationConfig struct {
	RegionID                   types.String `tfsdk:"region_id"`
	PermissionGroupFuid        types.String `tfsdk:"permission_group_fuid"`
	SfsUID                     types.String `tfsdk:"sfs_uid"`
	VpcID                      types.String `tfsdk:"vpc_id"`
	VpcName                    types.String `tfsdk:"vpc_name"`
	VpcCidr                    types.String `tfsdk:"vpc_cidr"`
	PermissionGroupName        types.String `tfsdk:"permission_group_name"`
	PermissionGroupDescription types.String `tfsdk:"permission_group_description"`
	PermissionGroupIsDefault   types.Bool   `tfsdk:"permission_group_is_default"`
	ID                         types.String `tfsdk:"id"`
}
