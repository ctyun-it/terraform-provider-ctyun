package nat

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctnat"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &ctyunPrivateDnatResource{}
	_ resource.ResourceWithConfigure   = &ctyunPrivateDnatResource{}
	_ resource.ResourceWithImportState = &ctyunPrivateDnatResource{}
)

type ctyunPrivateDnatResource struct {
	meta *common.CtyunMetadata
}

func NewCtyunPrivateDnatResource() resource.Resource {
	return &ctyunPrivateDnatResource{}
}

func (c *ctyunPrivateDnatResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_private_nat_dnat"
}

func (c *ctyunPrivateDnatResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026759/10378362`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "ID，同dnat_id",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"dnat_id": schema.StringAttribute{
				Computed:      true,
				Description:   "PrivateDnat规则的id",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，默认使用provider ctyun总region_id 或者环境变量",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"nat_gateway_id": schema.StringAttribute{
				Required:    true,
				Description: "NAT网关Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"external_ip": schema.StringAttribute{
				Required:    true,
				Description: "中转IP 支持更新",
				Validators: []validator.String{
					validator2.Ip(),
				},
			},
			"external_port": schema.Int32Attribute{
				Required:    true,
				Description: "对外的端口（1-65535） 支持更新",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"internal_port": schema.Int32Attribute{
				Required:    true,
				Description: "对应的内部端口（1-65535）支持更新",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"internal_ip": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "对应的内部IP(和port_id二选一) 支持更新",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("port_id")),
				},
			},
			"port_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "对应的网卡ID(和internal_ip二选一) 支持更新",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("internal_ip")),
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "协议: tcp, udp 支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf("tcp", "udp"),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "DNAT描述支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·！@#￥%……&*（） —— -+={}\\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 128),
					validator2.Desc(),
					validator2.DescNotStartWithHttp(),
				},
			},
			"state": schema.StringAttribute{
				Computed:      true,
				Description:   "运行状态: running代表运行中, freeze代表已冻结, expired代表已到期",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"create_time": schema.StringAttribute{
				Computed:      true,
				Description:   "创建时间，为UTC格式",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"port_name": schema.StringAttribute{
				Computed:      true,
				Description:   "网卡名称",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"device_id": schema.StringAttribute{
				Computed:      true,
				Description:   "网卡对应的设备ID",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (c *ctyunPrivateDnatResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunPrivateDnatConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 确保可选字段有正确的零值而不是null，避免后续状态不一致
	if plan.PortID.IsNull() {
		plan.PortID = types.StringValue("")
	}
	if plan.InternalIP.IsNull() {
		plan.InternalIP = types.StringValue("")
	}
	if plan.Description.IsNull() {
		plan.Description = types.StringValue("")
	}

	// 创建
	resp, err := c.meta.Apis.SdkCtNatApis.CtnatCreatePrivatenatDnatApi.Do(ctx, c.meta.SdkCredential, &ctnat.CtnatCreatePrivatenatDnatRequest{
		RegionID:     plan.RegionID.ValueString(),
		NatGatewayID: plan.NatGatewayID.ValueString(),
		ExternalIP:   plan.ExternalIP.ValueString(),
		ExternalPort: plan.ExternalPort.ValueInt32(),
		InternalPort: plan.InternalPort.ValueInt32(),
		InternalIP:   plan.InternalIP.ValueString(),
		PortID:       plan.PortID.ValueString(),
		Protocol:     plan.Protocol.ValueString(),
		Description:  plan.Description.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}

	// 设置ID
	if resp.ReturnObj != nil {
		plan.DnatID = types.StringValue(resp.ReturnObj.DnatID)
		plan.ID = types.StringValue(resp.ReturnObj.DnatID)
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建后反查创建后的dnat信息
	err = c.getAndMergePrivateDnat(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunPrivateDnatResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunPrivateDnatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergePrivateDnat(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(request.State.Set(ctx, &state)...)
}

func (c *ctyunPrivateDnatResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取tf文件中配置
	var plan CtyunPrivateDnatConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunPrivateDnatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 确保必需的字段存在
	if plan.DnatID.IsNull() || plan.DnatID.ValueString() == "" {
		plan.DnatID = state.DnatID
	}

	if plan.NatGatewayID.IsNull() || plan.NatGatewayID.ValueString() == "" {
		plan.NatGatewayID = state.NatGatewayID
	}

	// 更新dnat信息
	resp, err := c.meta.Apis.SdkCtNatApis.CtnatModifyPrivatenatDnatApi.Do(ctx, c.meta.SdkCredential, &ctnat.CtnatModifyPrivatenatDnatRequest{
		RegionID:     plan.RegionID.ValueString(),
		DnatID:       plan.DnatID.ValueString(),
		NatGatewayID: plan.NatGatewayID.ValueString(),
		ExternalIP:   plan.ExternalIP.ValueString(),
		ExternalPort: plan.ExternalPort.ValueInt32(),
		InternalIP:   plan.InternalIP.ValueString(),
		PortID:       plan.PortID.ValueString(),
		InternalPort: plan.InternalPort.ValueInt32(),
		Protocol:     plan.Protocol.ValueString(),
		Description:  plan.Description.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergePrivateDnat(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *ctyunPrivateDnatResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunPrivateDnatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 删除私网DNAT规则
	resp, err := c.meta.Apis.SdkCtNatApis.CtnatDeletePrivatenatDnatApi.Do(ctx, c.meta.SdkCredential, &ctnat.CtnatDeletePrivatenatDnatRequest{
		RegionID: state.RegionID.ValueString(),
		DnatID:   state.DnatID.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
}

func (c *ctyunPrivateDnatResource) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)

	c.meta = meta
}

func (c *ctyunPrivateDnatResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[natGateWayID],[projectID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunPrivateDnatConfig
	var ID, projectID, regionID, natGateWayID string
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		projectID = c.meta.GetExtraIfEmpty(projectID, common.ExtraProjectId)
		err = terraform_extend.Split(request.ID, &ID, &natGateWayID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &ID, &natGateWayID, &projectID, &regionID)
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
	if natGateWayID == "" {
		err = fmt.Errorf("natGateWayID不能为空")
	}
	config.ID = types.StringValue(ID)
	config.DnatID = types.StringValue(ID)
	config.RegionID = types.StringValue(regionID)
	config.NatGatewayID = types.StringValue(natGateWayID)
	err = c.getAndMergePrivateDnat(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *ctyunPrivateDnatResource) getAndMergePrivateDnat(ctx context.Context, plan *CtyunPrivateDnatConfig) (err error) {
	resp, err := c.meta.Apis.SdkCtNatApis.CtnatQueryPrivatenatDnatApi.Do(ctx, c.meta.SdkCredential, &ctnat.CtnatQueryPrivatenatDnatRequest{
		RegionID:     plan.RegionID.ValueString(),
		NatGatewayID: plan.NatGatewayID.ValueString(),
		PageNo:       1,
		PageSize:     50,
	})
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}

	// 查找匹配的DNAT规则
	var targetDnat *ctnat.CtnatQueryPrivatenatDnatReturnObjResponse
	for _, dnat := range resp.ReturnObj {
		if dnat.DnatID == plan.DnatID.ValueString() {
			targetDnat = dnat
			break
		}
	}

	// 如果没找到，尝试通过其他属性匹配
	if targetDnat == nil && !plan.DnatID.IsNull() && plan.DnatID.ValueString() != "" {
		for _, dnat := range resp.ReturnObj {
			// 根据请求参数匹配
			if dnat.ExternalIP == plan.ExternalIP.ValueString() &&
				dnat.ExternalPort == plan.ExternalPort.ValueInt32() &&
				dnat.InternalPort == plan.InternalPort.ValueInt32() &&
				dnat.Protocol == plan.Protocol.ValueString() {
				// 如果internal_ip或port_id匹配其中一个
				if (plan.InternalIP.ValueString() != "" && dnat.InternalIP == plan.InternalIP.ValueString()) ||
					(plan.PortID.ValueString() != "" && dnat.PortID == plan.PortID.ValueString()) {
					targetDnat = dnat
					break
				}
			}
		}
	}

	if targetDnat == nil {
		err = fmt.Errorf("private dnat rule not found")
		return
	}

	plan.DnatID = types.StringValue(targetDnat.DnatID)
	plan.ID = types.StringValue(targetDnat.DnatID)
	plan.ExternalIP = types.StringValue(targetDnat.ExternalIP)
	plan.ExternalPort = types.Int32Value(targetDnat.ExternalPort)
	plan.InternalPort = types.Int32Value(targetDnat.InternalPort)

	// 只有当internal_ip在配置中已设置或不为null时，才更新该值
	if !plan.InternalIP.IsNull() || targetDnat.InternalIP != "" {
		plan.InternalIP = types.StringValue(targetDnat.InternalIP)
	} else {
		plan.InternalIP = types.StringValue("")
	}

	// 处理可能为null的字段
	if targetDnat.PortID == "" && plan.PortID.IsNull() {
		plan.PortID = types.StringValue("")
	} else {
		plan.PortID = types.StringValue(targetDnat.PortID)
	}

	plan.PortName = types.StringValue(targetDnat.PortName)
	plan.DeviceID = types.StringValue(targetDnat.DeviceID)
	plan.Protocol = types.StringValue(targetDnat.Protocol)

	// 处理描述字段
	if targetDnat.Description == "" && plan.Description.IsNull() {
		plan.Description = types.StringValue("")
	} else {
		plan.Description = types.StringValue(targetDnat.Description)
	}

	plan.State = types.StringValue(targetDnat.State)
	plan.CreatedAt = types.StringValue(targetDnat.CreatedAt)

	return nil
}

type CtyunPrivateDnatConfig struct {
	ID           types.String `tfsdk:"id"`
	RegionID     types.String `tfsdk:"region_id"`      //区域id
	NatGatewayID types.String `tfsdk:"nat_gateway_id"` //NAT网关Id
	DnatID       types.String `tfsdk:"dnat_id"`        /*  DNAT规则的ID  */
	ExternalIP   types.String `tfsdk:"external_ip"`    /*  中转IP  */
	ExternalPort types.Int32  `tfsdk:"external_port"`  /*  外部端口  */
	InternalIP   types.String `tfsdk:"internal_ip"`    /*  内部IP  */
	InternalPort types.Int32  `tfsdk:"internal_port"`  /*  内部端口  */
	PortID       types.String `tfsdk:"port_id"`        /*  对应的网卡ID  */
	PortName     types.String `tfsdk:"port_name"`      /*  网卡名称  */
	DeviceID     types.String `tfsdk:"device_id"`      /*  网卡对应的设备ID  */
	Protocol     types.String `tfsdk:"protocol"`       /*  协议: tcp/udp  */
	State        types.String `tfsdk:"state"`          /*  DNAT状态: running代表运行中, freeze代表已冻结, expired代表已到期  */
	CreatedAt    types.String `tfsdk:"create_time"`    /*  创建时间  */
	Description  types.String `tfsdk:"description"`    /*  描述  */
}
