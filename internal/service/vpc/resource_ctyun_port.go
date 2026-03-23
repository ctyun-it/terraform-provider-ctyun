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
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ctyunNetworkInterface{}
	_ resource.ResourceWithConfigure   = &ctyunNetworkInterface{}
	_ resource.ResourceWithImportState = &ctyunNetworkInterface{}
)

func NewCtyunNetworkInterface() resource.Resource {
	return &ctyunNetworkInterface{}
}

type ctyunNetworkInterface struct {
	meta *common.CtyunMetadata
	name string
}

func (c *ctyunNetworkInterface) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_port"
}

type CtyunNetworkInterfaceConfig struct {
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	RegionId                types.String `tfsdk:"region_id"`
	SubnetId                types.String `tfsdk:"subnet_id"`
	PrimaryIpAddress        types.String `tfsdk:"primary_ip_address"`
	SecurityGroupIds        types.Set    `tfsdk:"security_group_ids"`
	SecondaryPrivateIpCount types.Int32  `tfsdk:"secondary_private_ip_count"`
	SecondaryPrivateIps     types.Set    `tfsdk:"secondary_private_ips"`
	Ipv6AddressCount        types.Int32  `tfsdk:"ipv6_address_count"`
	Ipv6Addresses           types.List   `tfsdk:"ipv6_addresses"`
	NetworkInterfaceId      types.String `tfsdk:"port_id"`
	MacAddress              types.String `tfsdk:"mac_address"`
	InstanceId              types.String `tfsdk:"instance_id"`
	InstanceType            types.String `tfsdk:"instance_type"`
	Status                  types.String `tfsdk:"status"`
}

func (c *ctyunNetworkInterface) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
<<<<<<< HEAD
		MarkdownDescription: utils.FormatDesc("管理弹性网卡", "弹性网卡（Elastic Network Interface）", "https://www.ctyun.cn/document/10026730/10225195"),
=======
		MarkdownDescription: utils.FormatDesc("管理弹性网卡", "PORT", "https://www.ctyun.cn/document/10026730/10225195"),
>>>>>>> a527f4c5a75df1829a2bf97af49cf86513680061
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "网卡ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "弹性网卡名称。支持拉丁字母、中文、数字，下划线，连字符，中文/英文字母开头，不能以http:/https:开头，长度2-32  支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 32),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z\u4e00-\u9fa5][0-9a-zA-Z_\u4e00-\u9fa5-]*[0-9a-zA-Z_\u4e00-\u9fa5]$"), "弹性网卡名称不符合规则"),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "弹性网卡描述。支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}|《》？：“”【】、；‘'，。、，不能以http:/https:开头，长度0-128 支持更新",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				//Validators: []validator.String{
				//	stringvalidator.UTF8LengthAtMost(128),
				//	validator2.Desc(),
				//},
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
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
					//validator2.SubnetValidate(),
				},
			},
			"primary_ip_address": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "主私有IP地址，如果不指定则自动分配",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					validator2.Ip(),
				},
			},
			"security_group_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "安全组ID列表，最多支持10个",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
					setplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(10),
				},
			},
			"secondary_private_ip_count": schema.Int32Attribute{
				Optional:    true,
				Description: "辅助私有IP地址数量，指定私有IP地址数量，让系统为您自动创建IP地址，最多支持10个",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
					int32validator.AtMost(10),
				},
			},
			"secondary_private_ips": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "辅助私有IP地址列表，指定私有IP地址，不能和secondary_private_ip_count同时指定，最多支持10个",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
					setplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(10),
				},
			},
			"ipv6_address_count": schema.Int32Attribute{
				Optional:    true,
				Description: "IPv6地址数量，指定IPv6地址数量，让系统为您自动创建IPv6地址，最多支持10个",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
					int32validator.AtMost(10),
				},
			},
			"ipv6_addresses": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "IPv6地址列表，指定IPv6地址，不能和ipv6_address_count同时指定，最多支持10个",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
					listplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(10),
				},
			},
			"port_id": schema.StringAttribute{
				Computed:    true,
				Description: "网卡ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mac_address": schema.StringAttribute{
				Computed:    true,
				Description: "MAC地址",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Computed:    true,
				Description: "绑定的实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_type": schema.StringAttribute{
				Computed:    true,
				Description: "实例类型",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "网卡状态",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *ctyunNetworkInterface) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunNetworkInterfaceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.createNetworkInterface(ctx, &plan)
	if err != nil {
		return
	}

	// 查询网卡详细信息并更新状态
	err = c.getAndMergePort(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *ctyunNetworkInterface) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunNetworkInterfaceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 更新状态
	err = c.getAndMergePort(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (c *ctyunNetworkInterface) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan, state CtyunNetworkInterfaceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkBeforeUpdate(ctx, plan, state)
	if err != nil {
		return
	}
	err = c.updateNetworkInterface(ctx, &plan, &state)
	if err != nil {
		return
	}
	// 查询网卡详细信息并更新状态
	err = c.getAndMergePort(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *ctyunNetworkInterface) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunNetworkInterfaceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	_, err = c.getNetworkInterfaceById(ctx, &state)
	if err != nil {
		return
	}
	err = c.delete(ctx, state)
	if err != nil {
		return
	}

}

// ImportState 导入资源状态
func (c *ctyunNetworkInterface) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [id],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunNetworkInterfaceConfig
	var id, regionId string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) == 0 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		id = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &id, &regionId)
		if err != nil {
			return
		}
	}

	if id == "" {
		err = fmt.Errorf("id不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}
	cfg.RegionId = types.StringValue(regionId)
	cfg.NetworkInterfaceId = types.StringValue(id)

	// 查询网卡详细信息并更新状态
	err = c.getAndMergePort(ctx, &cfg)
	if err != nil {
		return
	}

	// 设置导入的属性
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)

}

// Configure 配置资源
func (c *ctyunNetworkInterface) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
func (c *ctyunNetworkInterface) createNetworkInterface(ctx context.Context, plan *CtyunNetworkInterfaceConfig) (err error) {
	// 根据传入的不同参数确定创建方式
	regionId := plan.RegionId.ValueString()
	// 构造创建请求参数
	createReq := &ctvpc.CtvpcCreatePortRequest{
		ClientToken: uuid.NewString(),
		RegionID:    regionId,
		SubnetID:    plan.SubnetId.ValueString(),
	}

	// 处理主私有IP地址
	if !plan.PrimaryIpAddress.IsNull() && !plan.PrimaryIpAddress.IsUnknown() {
		primaryIp := plan.PrimaryIpAddress.ValueString()
		createReq.PrimaryPrivateIp = &primaryIp
	}

	// 处理安全组ID列表
	if !plan.SecurityGroupIds.IsNull() && len(plan.SecurityGroupIds.Elements()) > 0 {
		var sgIds []string
		plan.SecurityGroupIds.ElementsAs(ctx, &sgIds, false)
		sgIdPtrs := make([]*string, len(sgIds))
		for i, sgId := range sgIds {
			sgIdPtrs[i] = &sgId
		}
		createReq.SecurityGroupIds = sgIdPtrs
	}

	// 处理辅助私有IP地址数量
	if !plan.SecondaryPrivateIpCount.IsNull() {
		createReq.SecondaryPrivateIpCount = plan.SecondaryPrivateIpCount.ValueInt32()
	}

	// 处理辅助私有IP地址列表
	if !plan.SecondaryPrivateIps.IsNull() && len(plan.SecondaryPrivateIps.Elements()) > 0 {
		var secondaryIps []string
		plan.SecondaryPrivateIps.ElementsAs(ctx, &secondaryIps, false)
		secondaryIpPtrs := make([]*string, len(secondaryIps))
		for i, ip := range secondaryIps {
			secondaryIpPtrs[i] = &ip
		}
		createReq.SecondaryPrivateIps = secondaryIpPtrs
	}

	// 处理IPv6地址数量
	if !plan.Ipv6AddressCount.IsNull() {
		ipv6Count := plan.Ipv6AddressCount.ValueInt32()
		if ipv6Count > 0 {
			ipv6Addresses := make([]*string, ipv6Count)
			for i := int32(0); i < ipv6Count; i++ {
				ipv6Addresses[i] = nil // 表示自动分配
			}
			createReq.Ipv6Addresses = ipv6Addresses
		}
	}

	// 处理IPv6地址列表
	if !plan.Ipv6Addresses.IsNull() && len(plan.Ipv6Addresses.Elements()) > 0 {
		var ipv6Addresses []string
		plan.Ipv6Addresses.ElementsAs(ctx, &ipv6Addresses, false)
		ipv6AddrPtrs := make([]*string, len(ipv6Addresses))
		for i, addr := range ipv6Addresses {
			ipv6AddrPtrs[i] = &addr
		}
		createReq.Ipv6Addresses = ipv6AddrPtrs
	}

	// 处理名称和描述
	if !plan.Name.IsNull() {
		name := plan.Name.ValueString()
		createReq.Name = &name
	}
	if !plan.Description.IsNull() {
		description := plan.Description.ValueString()
		createReq.Description = &description
	}

	// 调用API创建弹性网卡
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreatePortApi.Do(ctx, c.meta.SdkCredential, createReq)
	if err != nil {
		return
	}
	if resp.ReturnObj == nil || resp.ReturnObj.NetworkInterfaceID == nil {
		return fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	}
	// 更新计划中的所有字段
	plan.Id = types.StringPointerValue(resp.ReturnObj.NetworkInterfaceID)
	plan.NetworkInterfaceId = types.StringPointerValue(resp.ReturnObj.NetworkInterfaceID)
	return
}
func (c *ctyunNetworkInterface) updateNetworkInterface(ctx context.Context, plan, state *CtyunNetworkInterfaceConfig) (err error) {
	updateReq := &ctvpc.CtvpcUpdatePortRequest{
		ClientToken:        uuid.NewString(),
		RegionID:           plan.RegionId.ValueString(),
		NetworkInterfaceID: plan.Id.ValueString(),
	}

	// 处理名称
	if !plan.Name.IsNull() {
		updateReq.Name = plan.Name.ValueStringPointer()
	}
	// 处理描述
	if !plan.Description.IsNull() {
		updateReq.Description = plan.Description.ValueStringPointer()
	}
	// 处理安全组ID列表
	if !plan.SecurityGroupIds.IsNull() && len(plan.SecurityGroupIds.Elements()) > 0 {
		var sgIds []string
		plan.SecurityGroupIds.ElementsAs(ctx, &sgIds, false)
		sgIdPtrs := make([]*string, len(sgIds))
		for i, sgId := range sgIds {
			sgIdPtrs[i] = &sgId
		}
		updateReq.SecurityGroupIDs = sgIdPtrs
	}
	// 调用API更新网卡属性
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcUpdatePortApi.Do(ctx, c.meta.SdkCredential, updateReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}
	return
}

// getNetworkInterface 获取网卡详细信息
func (c *ctyunNetworkInterface) getNetworkInterfaceById(ctx context.Context, plan *CtyunNetworkInterfaceConfig) (networkInterface *ctvpc.CtvpcShowPortReturnObjResponse, err error) {
	req := &ctvpc.CtvpcShowPortRequest{
		RegionID:           plan.RegionId.ValueString(),
		NetworkInterfaceID: plan.NetworkInterfaceId.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowPortApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		msg := utils.SecString(resp.Message)
		if strings.Contains(msg, "is not exists") {
			err = common.ResourceNotExistError
		} else {
			err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		}
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	networkInterface = resp.ReturnObj
	return
}

// getAndMergePort 查询网卡信息并合并到资源配置中
func (c *ctyunNetworkInterface) getAndMergePort(ctx context.Context, plan *CtyunNetworkInterfaceConfig) (err error) {
	networkInterface, err := c.getNetworkInterfaceById(ctx, plan)
	if err != nil {
		return
	}
	// 更新计划中的所有字段
	plan.Id = types.StringPointerValue(networkInterface.NetworkInterfaceID)
	plan.NetworkInterfaceId = types.StringPointerValue(networkInterface.NetworkInterfaceID)
	plan.Name = types.StringPointerValue(networkInterface.NetworkInterfaceName)
	plan.Description = types.StringPointerValue(networkInterface.Description)
	plan.MacAddress = types.StringPointerValue(networkInterface.MacAddress)
	plan.SubnetId = types.StringPointerValue(networkInterface.SubnetID)
	plan.PrimaryIpAddress = types.StringPointerValue(networkInterface.PrimaryPrivateIp)
	plan.InstanceId = types.StringPointerValue(networkInterface.InstanceID)
	plan.InstanceType = types.StringPointerValue(networkInterface.InstanceType)

	// 设置状态
	if networkInterface.AdminStatus != nil {
		plan.Status = types.StringValue(*networkInterface.AdminStatus)
	} else {
		plan.Status = types.StringValue("UNKNOWN")
	}

	// 设置安全组ID
	if networkInterface.SecurityGroupIds != nil {
		sgIds := make([]attr.Value, len(networkInterface.SecurityGroupIds))
		for i, sgId := range networkInterface.SecurityGroupIds {
			if sgId != nil {
				sgIds[i] = types.StringValue(*sgId)
			}
		}
		plan.SecurityGroupIds, _ = types.SetValue(types.StringType, sgIds)
	} else {
		// 如果没有安全组ID，确保字段被正确初始化为空集合
		plan.SecurityGroupIds, _ = types.SetValue(types.StringType, []attr.Value{})
	}

	// 设置辅助私有IP
	if networkInterface.SecondaryPrivateIps != nil {
		secondaryIps := make([]attr.Value, len(networkInterface.SecondaryPrivateIps))
		for i, ip := range networkInterface.SecondaryPrivateIps {
			if ip != nil {
				secondaryIps[i] = types.StringValue(*ip)
			}
		}
		plan.SecondaryPrivateIps, _ = types.SetValue(types.StringType, secondaryIps)
	} else {
		// 如果没有辅助私有IP，确保字段被正确初始化为空集合
		plan.SecondaryPrivateIps, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if len(networkInterface.SecondaryPrivateIps) > 0 {
		plan.SecondaryPrivateIpCount = types.Int32Value(int32(len(networkInterface.SecondaryPrivateIps)))

	}
	// 设置IPv6地址
	if networkInterface.Ipv6Addresses != nil {
		ipv6Addrs := make([]attr.Value, len(networkInterface.Ipv6Addresses))
		for i, addr := range networkInterface.Ipv6Addresses {
			if addr != nil {
				ipv6Addrs[i] = types.StringValue(*addr)
			}
		}
		plan.Ipv6Addresses, _ = types.ListValue(types.StringType, ipv6Addrs)
	} else {
		// 如果没有IPv6地址，确保字段被正确初始化为空列表
		plan.Ipv6Addresses, _ = types.ListValue(types.StringType, []attr.Value{})
	}

	return
}
func (c *ctyunNetworkInterface) delete(ctx context.Context, state CtyunNetworkInterfaceConfig) (err error) {
	// 构造删除请求参数
	deleteReq := &ctvpc.CtvpcDeletePortRequest{
		ClientToken:        uuid.NewString(),
		RegionID:           state.RegionId.ValueString(),
		NetworkInterfaceID: state.NetworkInterfaceId.ValueString(),
	}

	// 调用API删除弹性网卡
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeletePortApi.Do(ctx, c.meta.SdkCredential, deleteReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}
	return
}

func (c *ctyunNetworkInterface) checkBeforeUpdate(ctx context.Context, plan CtyunNetworkInterfaceConfig, state CtyunNetworkInterfaceConfig) (err error) {
	return
}
