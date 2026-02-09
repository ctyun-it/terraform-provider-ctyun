package vpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	_ resource.Resource                = &ctyunVipAssociation{}
	_ resource.ResourceWithConfigure   = &ctyunVipAssociation{}
	_ resource.ResourceWithImportState = &ctyunVipAssociation{}
)

func NewCtyunVipAssociation() resource.Resource {
	return &ctyunVipAssociation{}
}

type ctyunVipAssociation struct {
	meta *common.CtyunMetadata
	name string
}

func (c *ctyunVipAssociation) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vip_association"
	c.name = response.TypeName
}

func (c *ctyunVipAssociation) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理虚拟IP的绑定", "VIP", "https://www.ctyun.cn/document/10026730/10224288"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识，格式为 region_id:vip_id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"vip_id": schema.StringAttribute{
				Required:    true,
				Description: "高可用虚IP的ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"resource_type": schema.StringAttribute{
				Required:    true,
				Description: "绑定的实例类型，VM 表示虚拟机ECS, PM 表示裸金属, NETWORK 表示弹性 IP",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("VM", "PM", "NETWORK"),
				},
			},
			"network_interface_id": schema.StringAttribute{
				Optional:    true,
				Description: "虚拟网卡ID, 该网卡属于instance_id, 当 resource_type 为 VM / PM 时，必填",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("resource_type"),
						types.StringValue("PM"),
						types.StringValue("VM"),
					),
				},
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "ECS示例ID，当 resource_type 为 VM / PM 时，必填",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("resource_type"),
						types.StringValue("PM"),
						types.StringValue("VM"),
					),
				},
			},
			"floating_id": schema.StringAttribute{
				Optional:    true,
				Description: "弹性IP ID，当 resource_type 为 NETWORK 时，必填",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("resource_type"),
						types.StringValue("NETWORK"),
					),
				},
			},
		},
	}
}

func (c *ctyunVipAssociation) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunVipAssociation) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunVipAssociationConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.create(ctx, &plan)
	if err != nil {
		response.Diagnostics.AddError(
			"绑定虚拟IP失败",
			err.Error(),
		)
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunVipAssociation) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunVipAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (c *ctyunVipAssociation) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// 虚拟IP绑定资源不支持更新操作，直接返回
	var plan CtyunVipAssociationConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunVipAssociation) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunVipAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.delete(ctx, &state)
	if err != nil {
		response.Diagnostics.AddError(
			"解绑虚拟IP失败",
			err.Error(),
		)
		return
	}
}
func (c *ctyunVipAssociation) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	var hasErrorOccurred = false
	defer func() {
		if err != nil && !hasErrorOccurred {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [vip_id],[(instance_id:network_interface_id)/floating_id],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var state CtyunVipAssociationConfig
	var vipId, info, regionId, instanceId, networkInterfaceId string
	if strings.Count(request.ID, common.ImportSeparator) == 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &vipId, &info)
		if err != nil {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [vip_id],[(instance_id:network_interface_id)/floating_id],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	} else {
		err = terraform_extend.Split(request.ID, &vipId, &info, &regionId)
		if err != nil {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [vip_id],[(instance_id:network_interface_id)/floating_id],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}

	if vipId == "" {
		err = fmt.Errorf("vip_id不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}
	if info == "" {
		err = fmt.Errorf("info不能为空")
		return
	}
	state.VipId = types.StringValue(vipId)
	state.RegionId = types.StringValue(regionId)
	if strings.Count(info, ":") == 1 {
		err = terraform_extend.SplitComma(info, &instanceId, &networkInterfaceId)
		if err != nil {
			return
		}

		state.InstanceId = types.StringValue(instanceId)
		// 在 ImportState 方法中添加判断逻辑
		if strings.HasPrefix(instanceId, "ss-") {
			// 如果以 ss- 开头，可以认为是某种特定类型的资源
			// 例如设置 resourceType 为 "VM" 或其他适当的类型
			state.ResourceType = types.StringValue("PM")
		} else {
			state.ResourceType = types.StringValue("VM")
		}
		state.NetworkInterfaceId = types.StringValue(networkInterfaceId)
		err = c.getAndMerge(ctx, &state)
		if err != nil {
			hasErrorOccurred = true
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [vipId],[instance_id:network_interface_id],[regionId]"
			response.Diagnostics.AddError(title, detail)
			return
		}
	} else {
		state.ResourceType = types.StringValue("NETWORK")
		state.FloatingId = types.StringValue(info)
		err = c.getAndMerge(ctx, &state)
		if err != nil {
			hasErrorOccurred = true
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [vipId],[floating_id],[regionId]"
			response.Diagnostics.AddError(title, detail)
			return
		}
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

// create 绑定虚拟IP
func (c *ctyunVipAssociation) create(ctx context.Context, plan *CtyunVipAssociationConfig) (err error) {
	// 获取region_id，如果未提供则从provider中获取
	regionId := plan.RegionId.ValueString()
	if regionId == "" {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
	}

	// 准备请求参数
	req := &ctvpc.CtvpcBindHavipRequest{
		ClientToken:  uuid.NewString(),
		RegionID:     regionId,
		ResourceType: plan.ResourceType.ValueString(),
		HaVipID:      plan.VipId.ValueString(),
	}

	// 设置可选参数
	if !plan.NetworkInterfaceId.IsNull() {
		networkInterfaceId := plan.NetworkInterfaceId.ValueString()
		req.NetworkInterfaceID = &networkInterfaceId
	}

	if !plan.InstanceId.IsNull() {
		instanceId := plan.InstanceId.ValueString()
		req.InstanceID = &instanceId
	}

	if !plan.FloatingId.IsNull() {
		floatingId := plan.FloatingId.ValueString()
		req.FloatingID = &floatingId
	}

	// 调用API绑定HaVip
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcBindHavipApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 等待绑定完成
	for {
		status := resp.ReturnObj.Status
		if status != nil && *status == "done" {
			break
		}

		// 如果状态是 in_progress，则继续轮询
		if status != nil && *status == "in_progress" {
			time.Sleep(5 * time.Second)
			resp, err = c.meta.Apis.SdkCtVpcApis.CtvpcBindHavipApi.Do(ctx, c.meta.SdkCredential, req)
			if err != nil {
				return
			} else if resp.StatusCode != common.NormalStatusCode {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return
			}
			continue
		}

		// 其他状态则报错
		message := "未知状态"
		if resp.ReturnObj.Message != nil {
			message = *resp.ReturnObj.Message
		}
		return fmt.Errorf("绑定失败: %s", message)
	}

	// 设置资源ID
	plan.Id = types.StringValue(plan.VipId.ValueString())
	plan.RegionId = types.StringValue(regionId)

	return nil
}

// getAndMerge 查询虚拟IP绑定信息并合并状态
func (c *ctyunVipAssociation) getAndMerge(ctx context.Context, state *CtyunVipAssociationConfig) (err error) {
	// 获取region_id，如果未提供则从provider中获取
	regionId := state.RegionId.ValueString()
	if regionId == "" {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
	}

	// 通过查询HaVip详情来确认绑定状态
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowHavipApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcShowHavipRequest{
		RegionID: regionId,
		HaVipID:  state.VipId.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 检查绑定信息是否存在
	returnObj := resp.ReturnObj
	found := false
	// 如果是网络类型绑定，检查NetworkInfo
	if state.ResourceType.ValueString() == "NETWORK" {
		if len(returnObj.NetworkInfo) == 0 {
			return
		}

		for _, network := range returnObj.NetworkInfo {
			if network.EipID != nil && *network.EipID == state.FloatingId.ValueString() {
				found = true

			}
		}
		if !found {
			return common.ResourceNotExistError
		}
	} else {
		// 如果是VM/PM类型绑定，检查InstanceInfo和BindPorts
		if len(returnObj.InstanceInfo) == 0 || len(returnObj.BindPorts) == 0 {
			return
		}
		for _, instance := range returnObj.InstanceInfo {
			if instance.Id != nil && *instance.Id == state.InstanceId.ValueString() {
				found = true
				break
			}
		}
		if !found {
			return common.ResourceNotExistError
		}
	}

	// 更新状态
	state.Id = types.StringValue(state.VipId.ValueString())
	state.RegionId = types.StringValue(regionId)

	return
}

// delete 解绑虚拟IP
func (c *ctyunVipAssociation) delete(ctx context.Context, state *CtyunVipAssociationConfig) (err error) {
	// 获取region_id，如果未提供则从provider中获取
	regionId := state.RegionId.ValueString()
	if regionId == "" {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
	}

	// 准备请求参数
	req := &ctvpc.CtvpcUnbindHavipRequest{
		ClientToken:  uuid.NewString(),
		RegionID:     regionId,
		ResourceType: state.ResourceType.ValueString(),
		HaVipID:      state.VipId.ValueString(),
	}

	// 设置必需参数
	if !state.NetworkInterfaceId.IsNull() {
		req.NetworkInterfaceID = state.NetworkInterfaceId.ValueString()
	}

	if !state.InstanceId.IsNull() {
		req.InstanceID = state.InstanceId.ValueStringPointer()

	}

	if !state.FloatingId.IsNull() {

		req.FloatingID = state.FloatingId.ValueStringPointer()
	}

	// 调用API解绑HaVip
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcUnbindHavipApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 等待解绑完成
	for {
		status := resp.ReturnObj.Status
		if status != nil && *status == "done" {
			break
		}

		// 如果状态是 in_progress，则继续轮询
		if status != nil && *status == "in_progress" {
			time.Sleep(5 * time.Second)
			resp, err = c.meta.Apis.SdkCtVpcApis.CtvpcUnbindHavipApi.Do(ctx, c.meta.SdkCredential, req)
			if err != nil {
				return
			} else if resp.StatusCode != common.NormalStatusCode {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return
			}

			continue
		}

		// 其他状态则报错
		message := "未知状态"
		if resp.ReturnObj.Message != nil {
			message = *resp.ReturnObj.Message
		}
		return fmt.Errorf("解绑失败: %s", message)
	}

	return nil
}

type CtyunVipAssociationConfig struct {
	Id                 types.String `tfsdk:"id"`
	RegionId           types.String `tfsdk:"region_id"`
	VipId              types.String `tfsdk:"vip_id"`
	ResourceType       types.String `tfsdk:"resource_type"`
	NetworkInterfaceId types.String `tfsdk:"network_interface_id"`
	InstanceId         types.String `tfsdk:"instance_id"`
	FloatingId         types.String `tfsdk:"floating_id"`
}
