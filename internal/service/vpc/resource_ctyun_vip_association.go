package vpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
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
	meta        *common.CtyunMetadata
	name        string
	eipService  *business.EipService
	portService *business.PortService
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
				Description: "资源唯一标识",
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
				Optional:           true,
				DeprecationMessage: "废弃字段，请不要指定",
				Description:        "资源类型",
			},
			"instance_id": schema.StringAttribute{
				Optional:           true,
				DeprecationMessage: "废弃字段，请不要指定",
				Description:        "实例ID",
			},
			"network_interface_id": schema.StringAttribute{
				Optional:    true,
				Description: "弹性网卡ID, 与弹性IP ID有且只能有一个",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("floating_id")),
					stringvalidator.ConflictsWith(path.MatchRoot("floating_id")),
				},
			},
			"floating_id": schema.StringAttribute{
				Optional:    true,
				Description: "弹性IP ID，与弹性网卡ID有且只能有一个",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("network_interface_id")),
					stringvalidator.ConflictsWith(path.MatchRoot("network_interface_id")),
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
	c.portService = business.NewPortService(c.meta)
	c.eipService = business.NewEipService(c.meta)
}

func (c *ctyunVipAssociation) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError("绑定虚拟IP失败", err.Error())
		}
	}()
	var plan CtyunVipAssociationConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkBefore(ctx, &plan)
	if err != nil {
		return
	}
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunVipAssociation) checkBefore(ctx context.Context, plan *CtyunVipAssociationConfig) error {
	eipID, portID, regionID := plan.FloatingId.ValueString(), plan.NetworkInterfaceId.ValueString(), plan.RegionId.ValueString()
	if eipID != "" {
		plan.resourceType = "NETWORK"
		return c.eipService.MustExist(ctx, eipID, regionID)
	}
	port, err := c.portService.GetPortDetail(ctx, portID, regionID)
	if err != nil {
		return err
	}
	plan.resourceType = utils.SecString(port.InstanceType)
	if plan.resourceType != "VM" && plan.resourceType != "BM" {
		return fmt.Errorf("%s 绑定的必须是云主机或物理机", portID)
	}
	if plan.resourceType == "BM" {
		plan.resourceType = "PM"
	}
	plan.instanceID = utils.SecString(port.InstanceID)
	if plan.instanceID == "" {
		return fmt.Errorf("%s 未查询到绑定的实例id", portID)
	}
	return nil
}

func (c *ctyunVipAssociation) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVipAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.getAndMerge(ctx, &state)
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
	var plan CtyunVipAssociationConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunVipAssociation) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError("解绑虚拟IP失败", err.Error())
		}
	}()
	var state CtyunVipAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkBefore(ctx, &state)
	if err != nil {
		return
	}
	err = c.delete(ctx, &state)
	if err != nil {
		return
	}
}
func (c *ctyunVipAssociation) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			detail := fmt.Sprintf("导入命令：\n"+
				"terraform import [%s].[导入配置名称] [vip_id],[network_interface_id],<region_id>\n"+
				"或\n"+
				"terraform import [%s].[导入配置名称] [vip_id],[floating_id],<region_id>\n", c.name, c.name)
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var state CtyunVipAssociationConfig
	var vipId, resourceID, regionId string
	if strings.Count(request.ID, common.ImportSeparator) == 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &vipId, &resourceID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &vipId, &resourceID, &regionId)
		if err != nil {
			return
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
	if resourceID == "" {
		err = fmt.Errorf("network_interface_id或floating_id不能为空")
		return
	}
	state.VipId = types.StringValue(vipId)
	state.RegionId = types.StringValue(regionId)
	if strings.HasPrefix(resourceID, "eip") {
		state.FloatingId = types.StringValue(resourceID)
	} else if strings.HasPrefix(resourceID, "port") {
		state.NetworkInterfaceId = types.StringValue(resourceID)
	} else {
		err = fmt.Errorf("输入的 %s 不是弹性IP的ID或弹性网卡ID", resourceID)
		return
	}
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

// create 绑定虚拟IP
func (c *ctyunVipAssociation) create(ctx context.Context, plan *CtyunVipAssociationConfig) (err error) {
	regionId := plan.RegionId.ValueString()
	// 准备请求参数
	req := &ctvpc.CtvpcBindHavipRequest{
		ClientToken:  uuid.NewString(),
		RegionID:     regionId,
		ResourceType: plan.resourceType,
		HaVipID:      plan.VipId.ValueString(),
	}
	switch plan.resourceType {
	case "NETWORK":
		req.FloatingID = plan.FloatingId.ValueStringPointer()
	case "VM", "BM":
		req.NetworkInterfaceID = plan.NetworkInterfaceId.ValueStringPointer()
		req.InstanceID = &plan.instanceID
	default:
		return fmt.Errorf("未知资源类型 %s", plan.resourceType)
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
	plan.Id = plan.VipId
	return
}

// getAndMerge 查询虚拟IP绑定信息并合并状态
func (c *ctyunVipAssociation) getAndMerge(ctx context.Context, state *CtyunVipAssociationConfig) (err error) {
	regionId := state.RegionId.ValueString()
	// 通过查询HaVip详情来确认绑定状态
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowHavipApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcShowHavipRequest{
		RegionID: regionId,
		HaVipID:  state.VipId.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		if utils.SecString(resp.ErrorCode) == common.OpenapiHavipNotFound {
			err = common.ResourceNotExistError
		} else {
			err = fmt.Errorf("API return error. Message: %s", utils.SecString(resp.Message))
		}
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	state.Id = state.VipId
	// 检查绑定信息是否存在
	returnObj := resp.ReturnObj
	var found bool
	for _, network := range returnObj.NetworkInfo {
		if network.EipID != nil && *network.EipID == state.FloatingId.ValueString() {
			found = true
			break
		}
	}
	if found {
		return
	}
	for _, port := range returnObj.BindPorts {
		if port.PortID != nil && *port.PortID == state.NetworkInterfaceId.ValueString() {
			found = true
			break
		}
	}
	if !found {
		return common.ResourceNotExistError
	}
	return
}

// delete 解绑虚拟IP
func (c *ctyunVipAssociation) delete(ctx context.Context, state *CtyunVipAssociationConfig) (err error) {
	// 获取region_id，如果未提供则从provider中获取
	regionId := state.RegionId.ValueString()
	// 准备请求参数
	req := &ctvpc.CtvpcUnbindHavipRequest{
		ClientToken:  uuid.NewString(),
		RegionID:     regionId,
		ResourceType: state.resourceType,
		HaVipID:      state.VipId.ValueString(),
	}
	switch state.resourceType {
	case "NETWORK":
		req.FloatingID = state.FloatingId.ValueStringPointer()
	case "VM", "BM":
		req.NetworkInterfaceID = state.NetworkInterfaceId.ValueString()
		req.InstanceID = &state.instanceID
	default:
		return fmt.Errorf("未知资源类型 %s", state.resourceType)
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

	resourceType string
	instanceID   string
}
