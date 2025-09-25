package nat

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
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
	"time"
)

var (
	_ resource.Resource                = &ctyunDnatResource{}
	_ resource.ResourceWithConfigure   = &ctyunDnatResource{}
	_ resource.ResourceWithImportState = &ctyunDnatResource{}
)

type ctyunDnatResource struct {
	meta *common.CtyunMetadata
}

func (c *ctyunDnatResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_nat_dnat"
}

func NewCtyunDnatResource() resource.Resource {
	return &ctyunDnatResource{}
}

func (c *ctyunDnatResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026759/10166499`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID，同dnat_id",
			},
			"dnat_id": schema.StringAttribute{
				Computed:    true,
				Description: "Dnat规则的id",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池Id，默认使用provider ctyun总region_id 或者环境变量",
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
			"external_id": schema.StringAttribute{
				Required:    true,
				Description: "弹性IP的ID，形如eip-xxxxx，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"external_port": schema.Int32Attribute{
				Required:    true,
				Description: "弹性IP公网端口，1 - 1024，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(1, 1024),
				},
			},
			"internal_port": schema.Int32Attribute{
				Required:    true,
				Description: "主机内网端口，1 - 65535，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "协议：tcp/udp，支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.DNatProtocols...),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "描述，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"dnat_type": schema.StringAttribute{
				Required:    true,
				Description: "dnat规则类型，支持传递instance或custom，支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf("instance", "custom"),
				},
			},
			"server_type": schema.StringAttribute{
				Optional:    true,
				Description: "服务器类型，当且仅当dnat_type为instance时必填，支持：VM / BM，支持更新",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("dnat_type"),
						types.StringValue(business.VirtualMachineTypeCloud),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("dnat_type"),
						types.StringValue(business.VirtualMachineTypeCustom),
					),
					stringvalidator.OneOf(business.ServerTypeVM, business.ServerTypeBM),
				},
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "云主机或物理机实例ID，当且仅当dnat_type为instance时必填，支持更新",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("dnat_type"),
						types.StringValue(business.VirtualMachineTypeCloud),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("dnat_type"),
						types.StringValue(business.VirtualMachineTypeCustom),
					),
				},
			},
			"internal_ip": schema.StringAttribute{
				Optional:    true,
				Description: "内部IP，当且仅当dnat_type为custom时必填，支持更新",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("dnat_type"),
						types.StringValue(business.VirtualMachineTypeCustom),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("dnat_type"),
						types.StringValue(business.VirtualMachineTypeCloud),
					),
					validator2.Ip(),
				},
			},
			"external_ip": schema.StringAttribute{
				Computed:    true,
				Description: "弹性公网IP地址",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "运行状态: ACTIVE / FREEZING / CREATING",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间",
			},
			"ip_expire_time": schema.StringAttribute{
				Computed:    true,
				Description: "ip到期时间",
			},
		},
	}
}

func (c *ctyunDnatResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunDnatConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建前检查
	err = c.checkBeforeCreateDnat(ctx, plan)
	if err != nil {
		return
	}

	// 创建DNAT
	id, err := c.createDnat(ctx, plan)
	if err != nil {
		return
	}

	// 如果loopResponse不为空，则表示创建成功,保存dnat id和状态
	plan.DNatID = types.StringValue(id)

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 反查信息
	err = c.getAndMergeDnat(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)

}

func (c *ctyunDnatResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunDnatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeDnat(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(request.State.Set(ctx, &state)...)
}

func (c *ctyunDnatResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf 文件中的
	var plan CtyunDnatConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// state中的
	var state CtyunDnatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新dnat规则
	err = c.updateDNat(ctx, state, plan)
	if err != nil {
		return
	}
	state.DnatType = plan.DnatType
	state.InstanceID = plan.InstanceID
	state.ServerType = plan.ServerType
	state.InternalIP = plan.InternalIP
	// 查询远端信息
	err = c.getAndMergeDnat(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunDnatResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunDnatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 删除操作
	// 1. 定义删除参数
	params := &ctvpc.CtvpcDeleteDnatEntryRequest{
		RegionID:    state.RegionID.ValueString(),
		DNatID:      state.DNatID.ValueString(),
		ClientToken: uuid.NewString(),
	}
	// 2.调用删除方法
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteDnatEntryApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}

	// 轮询检测时候彻底删除
	err = c.DeleteLoop(ctx, state)
	if err != nil {
		return
	}

}

func (c *ctyunDnatResource) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)

	c.meta = meta
}

func (c *ctyunDnatResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var cfg CtyunDnatConfig
	var id, ngID, regionID string
	err = terraform_extend.Split(request.ID, &id, &ngID, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.NatGatewayID = types.StringValue(ngID)
	cfg.DNatID = types.StringValue(id)
	err = c.getAndMergeDnat(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *ctyunDnatResource) getAndMergeDnat(ctx context.Context, cfg *CtyunDnatConfig) (err error) {
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowDnatEntryApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcShowDnatEntryRequest{
		RegionID:     cfg.RegionID.ValueString(),
		NatGatewayID: cfg.NatGatewayID.ValueString(),
		DNatID:       cfg.DNatID.ValueString(),
	})
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	dnat := resp.ReturnObj
	cfg.DNatID = utils.SecStringValue(dnat.DNatID)
	cfg.ID = utils.SecStringValue(dnat.DNatID)
	cfg.CreatedAt = utils.SecStringValue(dnat.CreationTime)
	cfg.Description = utils.SecStringValue(dnat.Description)
	cfg.IpExpireTime = utils.SecStringValue(dnat.IpExpireTime)
	cfg.ExternalIP = utils.SecStringValue(dnat.ExternalIp)
	cfg.ExternalID = utils.SecStringValue(dnat.ExternalID)
	cfg.Protocol = utils.SecStringValue(dnat.Protocol)
	cfg.State = utils.SecStringValue(dnat.State)
	cfg.ExternalPort = types.Int32Value(dnat.ExternalPort)
	cfg.InternalPort = types.Int32Value(dnat.InternalPort)
	switch cfg.DnatType.ValueString() {
	case business.VirtualMachineTypeCloud:
		cfg.InstanceID = utils.SecStringValue(dnat.VirtualMachineID)
	case business.VirtualMachineTypeCustom:
		cfg.InternalIP = utils.SecStringValue(dnat.InternalIp)
	}

	return nil
}

// checkBeforeCreateDnat 创建dnat之前进行检查
func (c *ctyunDnatResource) checkBeforeCreateDnat(ctx context.Context, plan CtyunDnatConfig) (err error) {
	_, err = business.NewNatService(c.meta).GetNatByID(ctx, plan.NatGatewayID.ValueString(), plan.RegionID.ValueString())
	if err != nil {
		err = fmt.Errorf("校验nat网关失败" + err.Error())
		return
	}
	return
}

// createDnat 创建dnat规则
func (c *ctyunDnatResource) createDnat(ctx context.Context, plan CtyunDnatConfig) (id string, err error) {
	// 定义创建dnat规则的请求参数
	params := &ctvpc.CtvpcCreateDnatEntryRequest{
		RegionID:           plan.RegionID.ValueString(),
		NatGatewayID:       plan.NatGatewayID.ValueString(),
		ExternalID:         plan.ExternalID.ValueString(),
		ExternalPort:       plan.ExternalPort.ValueInt32(),
		VirtualMachineType: map[string]int32{business.VirtualMachineTypeCustom: 2, business.VirtualMachineTypeCloud: 1}[plan.DnatType.ValueString()],
		InternalPort:       plan.InternalPort.ValueInt32(),
		Protocol:           plan.Protocol.ValueString(),
		ClientToken:        uuid.NewString(),
		Description:        plan.Description.ValueStringPointer(),
	}

	switch plan.DnatType.ValueString() {
	case business.VirtualMachineTypeCloud:
		params.VirtualMachineID = plan.InstanceID.ValueStringPointer()
		params.ServerType = plan.ServerType.ValueStringPointer()
	case business.VirtualMachineTypeCustom:
		params.InternalIp = plan.InternalIP.ValueStringPointer()
	}

	// SDK接口：ctvpc_create_dnat_entry_api.go
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateDnatEntryApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil || resp.ReturnObj.Dnat == nil || resp.ReturnObj.Dnat.DnatID == nil {
		err = common.InvalidReturnObjError
		return
	}
	id = *resp.ReturnObj.Dnat.DnatID
	return
}

func (c *ctyunDnatResource) updateLoop(ctx context.Context, state *CtyunDnatConfig, updatedParams *ctvpc.CtvpcUpdateDnatEntryAttributeRequest, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*5, count)
	if err != nil {
		return
	}
	result := retryer.Start(func(currentTime int) bool {
		resp, err2 := c.meta.Apis.SdkCtVpcApis.CtvpcShowDnatEntryApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcShowDnatEntryRequest{
			RegionID:     state.RegionID.ValueString(),
			NatGatewayID: state.NatGatewayID.ValueString(),
			DNatID:       state.DNatID.ValueString(),
		})
		if err2 != nil {
			err = err2
			return false
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return false
		} else if resp.ReturnObj == nil {
			err = common.InvalidReturnObjError
			return false
		}
		return false
	})
	time.Sleep(5 * time.Second)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未更新!Dnat: " + state.DNatID.ValueString())
	}
	return
}

func (c *ctyunDnatResource) DeleteLoop(ctx context.Context, state CtyunDnatConfig) (err error) {
	var respErr error
	retryer, err := business.NewRetryer(time.Second*5, 60)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowDnatEntryApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcShowDnatEntryRequest{
				RegionID:     state.RegionID.ValueString(),
				NatGatewayID: state.NatGatewayID.ValueString(),
				DNatID:       state.DNatID.ValueString(),
			})
			if err != nil {
				respErr = err
				return false
			}
			// 如果返回为空了，说明已经删除成功
			if resp.ReturnObj == nil {
				return false
			} else {
				// 如果仍能查询到dnat信息，说明仍未删除完成
				return true
			}
		},
	)
	time.Sleep(5 * time.Second)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未删除!Dnat: " + state.DNatID.ValueString())
	}
	return respErr
}

func (c *ctyunDnatResource) updateDNat(ctx context.Context, state CtyunDnatConfig, plan CtyunDnatConfig) (err error) {
	params := &ctvpc.CtvpcUpdateDnatEntryAttributeRequest{
		ClientToken:        uuid.NewString(),
		RegionID:           state.RegionID.ValueString(),
		DNatID:             state.DNatID.ValueString(),
		Protocol:           plan.Protocol.ValueString(),
		VirtualMachineType: map[string]int32{business.VirtualMachineTypeCloud: 1, business.VirtualMachineTypeCustom: 2}[plan.DnatType.ValueString()],
		ExternalPort:       plan.ExternalPort.ValueInt32(),
		ExternalID:         plan.ExternalID.ValueStringPointer(),
		InternalIp:         plan.InternalIP.ValueStringPointer(),
		InternalPort:       plan.InternalPort.ValueInt32(),
		Description:        plan.Description.ValueStringPointer(),
		ServerType:         plan.ServerType.ValueStringPointer(),
		VirtualMachineID:   plan.InstanceID.ValueStringPointer(),
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcUpdateDnatEntryAttributeApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	}
	// 轮询确认已经更新完毕
	err = c.updateLoop(ctx, &state, params, 50)
	if err != nil {
		return
	}
	return
}

// contains 方法用于判断value时候被包含在list中，区分大小写
func (c *ctyunDnatResource) contains(value string, list []string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

type CtyunDnatConfig struct {
	ID           types.String `tfsdk:"id"`
	RegionID     types.String `tfsdk:"region_id"`      //区域id
	NatGatewayID types.String `tfsdk:"nat_gateway_id"` //要查询的私网NAT的ID
	DNatID       types.String `tfsdk:"dnat_id"`        //DNAT规则的ID
	ExternalID   types.String `tfsdk:"external_id"`    //中转IP ID
	ExternalIP   types.String `tfsdk:"external_ip"`    //中转IP
	ExternalPort types.Int32  `tfsdk:"external_port"`  //外部端口
	InternalIP   types.String `tfsdk:"internal_ip"`    //内部IP
	InternalPort types.Int32  `tfsdk:"internal_port"`  //内部端口
	Protocol     types.String `tfsdk:"protocol"`       //协议: tcp/udp
	State        types.String `tfsdk:"state"`          //DNAT状态: running代表运行中, freeze代表已冻结, expired代表已到期
	Description  types.String `tfsdk:"description"`    //描述
	InstanceID   types.String `tfsdk:"instance_id"`
	DnatType     types.String `tfsdk:"dnat_type"`
	ServerType   types.String `tfsdk:"server_type"`    //当 virtualMachineType 为 1 时，serverType 必传，支持: VM / BM （仅支持大写）
	CreatedAt    types.String `tfsdk:"created_at"`     //创建时间
	IpExpireTime types.String `tfsdk:"ip_expire_time"` //ip到期时间
}

type DnatLoopCreateResponse struct {
	DNatID types.String
	Status types.String
}
