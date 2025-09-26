package ebm

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebm"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &ctyunEbm{}
	_ resource.ResourceWithConfigure   = &ctyunEbm{}
	_ resource.ResourceWithImportState = &ctyunEbm{}
)

type ctyunEbmInterface struct {
	meta *common.CtyunMetadata
}

func NewCtyunEbmInterface() resource.Resource {
	return &ctyunEbmInterface{}
}

func (c *ctyunEbmInterface) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebm_interface"
}

type CtyunEbmInterfaceConfig struct {
	RegionID         types.String `tfsdk:"region_id"`
	AzName           types.String `tfsdk:"az_name"`
	InstanceID       types.String `tfsdk:"instance_id"`
	SubnetID         types.String `tfsdk:"subnet_id"`
	SecurityGroupIDs types.Set    `tfsdk:"security_group_ids"`
	Ipv4             types.String `tfsdk:"ipv4"`
	InterfaceID      types.String `tfsdk:"interface_id"`
	ID               types.String `tfsdk:"id"`

	secGroupParams []string
}

func (c *ctyunEbmInterface) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027724/10040142`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"interface_id": schema.StringAttribute{
				Computed:    true,
				Description: "网卡ID",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称",
				Default:     defaults.AcquireFromGlobalString(common.ExtraAzName, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "物理机UUID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"security_group_ids": schema.SetAttribute{
				Description: "安全组ID，支持更新",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(validator2.SecurityGroupValidate()),
				},
			},
			"ipv4": schema.StringAttribute{
				Required:    true,
				Description: "IPV4地址",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					validator2.Ip(),
				},
			},
		},
	}
}

func (c *ctyunEbmInterface) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEbmInterfaceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 处理安全组集合
	diags := plan.SecurityGroupIDs.ElementsAs(ctx, &plan.secGroupParams, true) // 第二个参数为是否忽略未知值
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
	// 创建前检查
	err = c.checkBeforeCreate(ctx, plan)
	if err != nil {
		return
	}
	// 创建
	err = c.createInterface(ctx, plan)
	if err != nil {
		return
	}
	// 创建后检查
	err = c.checkAfterCreate(ctx, &plan)
	if err != nil {
		return
	}
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *ctyunEbmInterface) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbmInterfaceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEbmInterface) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunEbmInterfaceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunEbmInterfaceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 处理安全组集合
	diags := plan.SecurityGroupIDs.ElementsAs(ctx, &plan.secGroupParams, true) // 第二个参数为是否忽略未知值
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
	diags = state.SecurityGroupIDs.ElementsAs(ctx, &state.secGroupParams, true) // 第二个参数为是否忽略未知值
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
	// 更新
	err = c.updateInterface(ctx, plan, state)
	if err != nil {
		return
	}
	plan.InterfaceID = state.InterfaceID
	err = c.checkAfterUpdate(ctx, plan)
	if err != nil {
		return
	}
	// 查询远端信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *ctyunEbmInterface) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbmInterfaceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.deleteInterface(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDelete(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunEbmInterface) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [instance_id],[interface_id],[regionID],[azName]
func (c *ctyunEbmInterface) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEbmInterfaceConfig
	var instanceID, interfaceID, regionID, azName string
	err = terraform_extend.Split(request.ID, &instanceID, &interfaceID, &regionID, &azName)
	if err != nil {
		return
	}
	plan.InterfaceID = types.StringValue(interfaceID)
	plan.InstanceID = types.StringValue(instanceID)
	plan.RegionID = types.StringValue(regionID)
	plan.AzName = types.StringValue(azName)

	// 查询远端
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

// checkBeforeCreate 创建前检查
func (c *ctyunEbmInterface) checkBeforeCreate(ctx context.Context, plan CtyunEbmInterfaceConfig) (err error) {
	// 查询物理机
	instance, err := business.NewEbmService(c.meta).GetEbmInfo(
		ctx,
		plan.InstanceID.ValueString(),
		plan.RegionID.ValueString(),
		plan.AzName.ValueString(),
	)
	if err != nil {
		return
	}
	smart := utils.SecBool(instance.DeviceDetail.SmartNicExist)
	if !smart {
		return fmt.Errorf("物理机 %s 必须是弹性裸金属", utils.SecString(instance.InstanceUUID))
	}
	vpcID := utils.SecString(instance.VpcID)
	// 检查配置的每个安全组是否存在
	secService := business.NewSecurityGroupService(c.meta)
	for _, id := range plan.secGroupParams {
		err = secService.MustExistInVpc(ctx, vpcID, id, plan.RegionID.ValueString())
		if err != nil {
			return
		}
	}

	return nil
}

// createInterface 创建网卡
func (c *ctyunEbmInterface) createInterface(ctx context.Context, plan CtyunEbmInterfaceConfig) (err error) {
	params := &ctebm.EbmAddNicRequest{
		RegionID:       plan.RegionID.ValueString(),
		AzName:         plan.AzName.ValueString(),
		InstanceUUID:   plan.InstanceID.ValueString(),
		SubnetUUID:     plan.SubnetID.ValueString(),
		SecurityGroups: strings.Join(plan.secGroupParams, ","),
		Ipv4:           plan.Ipv4.ValueStringPointer(),
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmAddNicApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		return fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	}
	return nil
}

// checkAfterCreate 创建后检查，因为创建是异步的。
// 用 子网ID + ipv4 唯一确定是否创建成功。
func (c *ctyunEbmInterface) checkAfterCreate(ctx context.Context, plan *CtyunEbmInterfaceConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			instance, err := business.NewEbmService(c.meta).GetEbmInfo(
				ctx,
				plan.InstanceID.ValueString(),
				plan.RegionID.ValueString(),
				plan.AzName.ValueString(),
			)
			if err != nil {
				return false
			}
			for _, nic := range instance.Interfaces {
				if plan.Ipv4.ValueString() == utils.SecString(nic.Ipv4) &&
					plan.SubnetID.ValueString() == utils.SecString(nic.SubnetUUID) {
					plan.InterfaceID = utils.SecStringValue(nic.InterfaceUUID)
					executeSuccessFlag = true
					return false
				}
			}
			return true
		})
	if !executeSuccessFlag {
		return fmt.Errorf("物理机 %s 网卡创建失败", plan.InstanceID.ValueString())
	}
	return nil
}

// getAndMerge 查询
func (c *ctyunEbmInterface) getAndMerge(ctx context.Context, plan *CtyunEbmInterfaceConfig) (err error) {
	// 查询物理机
	instanceID := plan.InstanceID.ValueString()
	interfaceID := plan.InterfaceID.ValueString()
	regionID := plan.RegionID.ValueString()
	azName := plan.AzName.ValueString()

	instance, err := business.NewEbmService(c.meta).GetEbmInfo(
		ctx,
		instanceID,
		regionID,
		azName,
	)
	if err != nil {
		return
	}
	var secGroupIDs []string
	for _, nic := range instance.Interfaces {
		if interfaceID == utils.SecString(nic.InterfaceUUID) {
			plan.SubnetID = utils.SecStringValue(nic.SubnetUUID)
			plan.Ipv4 = utils.SecLowerStringValue(nic.Ipv4)

			for _, g := range nic.SecurityGroups {
				secGroupIDs = append(secGroupIDs, utils.SecString(g.SecurityGroupID))
			}
			ids, diags := types.SetValueFrom(ctx, types.StringType, secGroupIDs)
			if diags.HasError() {
				err = fmt.Errorf("从OpenAPI查询到了非预期的安全组")
				return
			}
			plan.secGroupParams = secGroupIDs
			plan.SecurityGroupIDs = ids
			plan.ID = plan.InterfaceID
			return
		}
	}
	err = fmt.Errorf("未找到目标网卡 %s", plan.InterfaceID.ValueString())
	return
}

// updateInterface 更新网卡的安全组
func (c *ctyunEbmInterface) updateInterface(ctx context.Context, plan, state CtyunEbmInterfaceConfig) (err error) {
	// 计算出需要删除哪些，需要更新哪些
	addSecGroup, delSecGroup := utils.DifferenceStrArray(plan.secGroupParams, state.secGroupParams)
	if len(addSecGroup) == 0 && len(delSecGroup) == 0 {
		return
	}
	params := &ctebm.EbmUpdateSecurityGroupRequest{
		RegionID:               state.RegionID.ValueString(),
		AzName:                 state.AzName.ValueString(),
		InstanceUUID:           state.InstanceID.ValueString(),
		InterfaceUUID:          state.InterfaceID.ValueString(),
		AddSecurityGroupIDList: strings.Join(addSecGroup, ","),
		DelSecurityGroupIDList: strings.Join(delSecGroup, ","),
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmUpdateSecurityGroupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		return fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	}
	return
}

// checkAfterUpdate 更新后检查是否更新成功
func (c *ctyunEbmInterface) checkAfterUpdate(ctx context.Context, plan CtyunEbmInterfaceConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			instance, err := business.NewEbmService(c.meta).GetEbmInfo(
				ctx,
				plan.InstanceID.ValueString(),
				plan.RegionID.ValueString(),
				plan.AzName.ValueString(),
			)
			if err != nil {
				return false
			}
			for _, nic := range instance.Interfaces {
				if plan.InterfaceID.ValueString() == utils.SecString(nic.InterfaceUUID) {
					var secGroupIDs []string
					for _, g := range nic.SecurityGroups {
						secGroupIDs = append(secGroupIDs, utils.SecString(g.SecurityGroupID))
					}
					if utils.AreStringSlicesEqual(secGroupIDs, plan.secGroupParams) { // 目标网卡的安全组已经更新完成
						executeSuccessFlag = true
						return false
					}
				}
			}
			return true
		})
	if !executeSuccessFlag {
		return fmt.Errorf("物理机 %s 网卡更新失败", plan.InstanceID.ValueString())
	}
	return nil
}

// deleteInterface 删除物理机网卡
func (c *ctyunEbmInterface) deleteInterface(ctx context.Context, plan CtyunEbmInterfaceConfig) (err error) {
	params := &ctebm.EbmRemoveNicRequest{
		RegionID:      plan.RegionID.ValueString(),
		AzName:        plan.AzName.ValueString(),
		InstanceUUID:  plan.InstanceID.ValueString(),
		InterfaceUUID: plan.InterfaceID.ValueString(),
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmRemoveNicApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		return fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	}
	return
}

// checkAfterDelete 删除后检查是否删除成功
func (c *ctyunEbmInterface) checkAfterDelete(ctx context.Context, plan CtyunEbmInterfaceConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			instance, err := business.NewEbmService(c.meta).GetEbmInfo(
				ctx,
				plan.InstanceID.ValueString(),
				plan.RegionID.ValueString(),
				plan.AzName.ValueString(),
			)
			if err != nil {
				return false
			}
			for _, nic := range instance.Interfaces {
				if plan.InterfaceID.ValueString() == utils.SecString(nic.InterfaceUUID) {
					return true
				}
			}
			executeSuccessFlag = true
			return false
		})
	if !executeSuccessFlag {
		return fmt.Errorf("物理机 %s 网卡删除失败", plan.InstanceID.ValueString())
	}
	return nil
}
