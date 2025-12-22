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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

type CtyunOceanfsPermissionGroupAssociation struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunOceanfsPermissionGroupAssociation) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_oceanfs_permission_group_association"
}

func (c *CtyunOceanfsPermissionGroupAssociation) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunOceanfsPermissionGroupAssociation() resource.Resource {
	return &CtyunOceanfsPermissionGroupAssociation{}
}

func (c *CtyunOceanfsPermissionGroupAssociation) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunOceanfsPermissionGroupAssociationConfig
	var regionId, permissionGroupId, sfsId, vpcId, subnetId string
	err = terraform_extend.Split(request.ID, &regionId, &permissionGroupId, &sfsId, &vpcId, &subnetId)
	if err != nil {
		return
	}
	config.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", vpcId, sfsId, regionId))
	config.RegionID = types.StringValue(regionId)
	config.SfsUID = types.StringValue(sfsId)
	config.PermissionGroupFuid = types.StringValue(permissionGroupId)
	config.VpcID = types.StringValue(vpcId)
	config.SubnetID = types.StringValue(subnetId)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunOceanfsPermissionGroupAssociation) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10088966/10332853",
		Attributes: map[string]schema.Attribute{
			"permission_group_id": schema.StringAttribute{
				Required:    true,
				Description: "oceanfs 权限组ID",
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
			"oceanfs_id": schema.StringAttribute{
				Required:    true,
				Description: "文件系统ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "虚拟私有云ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"is_vpce": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "文件系统绑定VPC时是否自动创建VPC终端节点。开启后本服务将为您创建免费的VPC终端节点（VPCE），连接文件存储服务。创建VPCE后将返回该VPC专属的挂载地址，通常需要1~3分钟。取值：\ntrue：创建VPC终端节点（推荐）\nfalse：不创建VPC终端节点\n注：物理机必须通过VPCE专属挂载地址访问文件系统，其它计算服务如云主机、容器为非必须",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Description: "子网ID",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "oceanfs与权限组绑定id",
			},
		},
	}
}

func (c *CtyunOceanfsPermissionGroupAssociation) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunOceanfsPermissionGroupAssociationConfig
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

func (c *CtyunOceanfsPermissionGroupAssociation) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunOceanfsPermissionGroupAssociationConfig
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

func (c *CtyunOceanfsPermissionGroupAssociation) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunOceanfsPermissionGroupAssociationConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunOceanfsPermissionGroupAssociationConfig
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

func (c *CtyunOceanfsPermissionGroupAssociation) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunOceanfsPermissionGroupAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunOceanfsPermissionGroupAssociation) create(ctx context.Context, config *CtyunOceanfsPermissionGroupAssociationConfig) error {
	params := &oceanfs.OceanfsVpcBindPermissionRequest{
		PermissionGroupFuid: config.PermissionGroupFuid.ValueString(),
		RegionID:            config.RegionID.ValueString(),
		SfsUID:              config.SfsUID.ValueString(),
		VpcID:               config.VpcID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsVpcBindPermissionApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("权限组绑定vpc失败(permission_group_id = %s, vpc_id = %s)，接口返回为nil。请与研发联系确认问题原因", config.PermissionGroupFuid.ValueString(), config.VpcID.ValueString())
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

	config.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", config.VpcID.ValueString(), config.SfsUID.ValueString(), config.RegionID.ValueString()))
	return nil
}

func (c *CtyunOceanfsPermissionGroupAssociation) getAndMerge(ctx context.Context, config *CtyunOceanfsPermissionGroupAssociationConfig) error {
	resp, err := c.getSfsVpcList(ctx, config)
	if err != nil {
		return err
	}
	if len(resp.ReturnObj.List) < 1 {
		err = fmt.Errorf("未查询到vpc信息(SfsUID=%s)，请检查参数是否正确", config.SfsUID.ValueString())
		return err
	}
	for _, item := range resp.ReturnObj.List {
		if *item.VpcID == config.VpcID.ValueString() && *item.PermissionGroupFuid == config.PermissionGroupFuid.ValueString() {
			return nil
		}
	}
	return fmt.Errorf("权限组绑定vpc未成功！")
}

func (c *CtyunOceanfsPermissionGroupAssociation) getSfsVpcList(ctx context.Context, config *CtyunOceanfsPermissionGroupAssociationConfig) (*oceanfs.OceanfsListVpcPermissionResponse, error) {
	params := &oceanfs.OceanfsListVpcPermissionRequest{
		RegionID: config.RegionID.ValueString(),
		SfsUID:   config.SfsUID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsListVpcPermissionApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取vpc列表失败(oceanfs_id = %s)，接口返回为nil。请与研发联系确认问题原因", config.SfsUID.ValueString())
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

func (c *CtyunOceanfsPermissionGroupAssociation) update(ctx context.Context, state *CtyunOceanfsPermissionGroupAssociationConfig, plan *CtyunOceanfsPermissionGroupAssociationConfig) error {
	params := &oceanfs.OceanfsVpcChangePermissionRequest{
		PermissionGroupFuid: plan.PermissionGroupFuid.ValueString(),
		RegionID:            state.RegionID.ValueString(),
		SfsUID:              state.SfsUID.ValueString(),
		VpcID:               plan.VpcID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsVpcChangePermissionApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新权限组绑定vpc失败(permission_group_id = %s, vpc_id = %s)，接口返回为nil。请与研发联系确认问题原因", state.PermissionGroupFuid.ValueString(), plan.VpcID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	// 绑定后需要轮询下
	err = c.bindLoop(ctx, plan)
	state.PermissionGroupFuid = plan.PermissionGroupFuid
	return nil
}

func (c *CtyunOceanfsPermissionGroupAssociation) delete(ctx context.Context, config CtyunOceanfsPermissionGroupAssociationConfig) error {
	params := &oceanfs.OceanfsVpcUnbindPermissionRequest{
		PermissionGroupFuid: config.PermissionGroupFuid.ValueString(),
		RegionID:            config.RegionID.ValueString(),
		SfsUID:              config.SfsUID.ValueString(),
		VpcID:               config.VpcID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsVpcUnbindPermissionApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("解绑权限组绑定vpc失败(permission_group_id = %s, vpc_id = %s)，接口返回为nil。请与研发联系确认问题原因", config.PermissionGroupFuid.ValueString(), config.VpcID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	return nil
}

func (c *CtyunOceanfsPermissionGroupAssociation) bindLoop(ctx context.Context, config *CtyunOceanfsPermissionGroupAssociationConfig) error {
	var err error
	retryer, err := business.NewRetryer(time.Second*10, 60)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.getSfsVpcList(ctx, config)
			if err2 != nil {
				err = err2
				return false
			}
			returnObj := resp.ReturnObj

			for _, association := range returnObj.List {
				if *association.VpcID == config.VpcID.ValueString() && *association.PermissionGroupFuid == config.PermissionGroupFuid.ValueString() {
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

type CtyunOceanfsPermissionGroupAssociationConfig struct {
	PermissionGroupFuid types.String `tfsdk:"permission_group_id"`
	RegionID            types.String `tfsdk:"region_id"`
	SfsUID              types.String `tfsdk:"oceanfs_id"`
	VpcID               types.String `tfsdk:"vpc_id"`
	IsVpce              types.Bool   `tfsdk:"is_vpce"`
	SubnetID            types.String `tfsdk:"subnet_id"`
	ID                  types.String `tfsdk:"id"`
}
