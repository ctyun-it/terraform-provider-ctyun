package ec

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CtyunExpressConnectRegionPeer struct {
	meta *common.CtyunMetadata
}

func (c *CtyunExpressConnectRegionPeer) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ec_region_peer"
}

func (c *CtyunExpressConnectRegionPeer) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta

}

func NewCtyunExpressConnectRegionPeer() resource.Resource {
	return &CtyunExpressConnectRegionPeer{}
}

func (c *CtyunExpressConnectRegionPeer) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[ecId],[packetId],[srcCgwId]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunExpressConnectRegionPeerConfig

	var id, ecId, packetId, srcCgwId string

	err = terraform_extend.Split(request.ID, &id, &ecId, &packetId, &srcCgwId)
	if err != nil {
		return
	}

	if id == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if ecId == "" {
		err = fmt.Errorf("ecId不能为空")
		return
	}
	if packetId == "" {
		err = fmt.Errorf("packetId不能为空")
		return
	}
	if srcCgwId == "" {
		err = fmt.Errorf("srcCgwId不能为空")
		return
	}

	config.ID = types.StringValue(id)
	config.EcID = types.StringValue(ecId)
	config.PacketID = types.StringValue(packetId)
	config.SrcCgwID = types.StringValue(srcCgwId)

	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}
func (c *CtyunExpressConnectRegionPeer) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026763/10038250",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "连接名称",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ec_id": schema.StringAttribute{
				Required:    true,
				Description: "云间高速ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(36, 36),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"src_cgw_id": schema.StringAttribute{
				Required:    true,
				Description: "本端网关ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(36, 36),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dst_cgw_id": schema.StringAttribute{
				Required:    true,
				Description: "对端网关ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(36, 36),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"packet_id": schema.StringAttribute{
				Required:    true,
				Description: "带宽包ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(36, 36),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rate": schema.Int32Attribute{
				Required:    true,
				Description: "带宽值（MB） 支持更新",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"route_learn": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(1),
				Description: "是否开启两端云网关路由自动学习，0：不开启，1：开启。默认开启",
				Validators: []validator.Int32{
					int32validator.OneOf(0, 1),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "跨域连接ID",
			},
			"src_region_id": schema.StringAttribute{
				Computed:    true,
				Description: "本端资源池ID",
			},
			"dst_region_id": schema.StringAttribute{
				Computed:    true,
				Description: "对端资源池ID",
			},
			"peer_type": schema.Int32Attribute{
				Computed:    true,
				Description: "互通类型，1：境内，2: 跨境（中国大陆-亚太），3: 境外（亚太），4: 定制",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间，为UTC格式",
			},
		},
	}
}

func (c *CtyunExpressConnectRegionPeer) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunExpressConnectRegionPeerConfig
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

func (c *CtyunExpressConnectRegionPeer) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunExpressConnectRegionPeerConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunExpressConnectRegionPeer) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunExpressConnectRegionPeerConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunExpressConnectRegionPeerConfig
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

func (c *CtyunExpressConnectRegionPeer) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunExpressConnectRegionPeerConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunExpressConnectRegionPeer) create(ctx context.Context, config *CtyunExpressConnectRegionPeerConfig) error {
	params := &ec.EcCreateRegionPeerRequest{
		PeerName:   config.Name.ValueString(),
		EcID:       config.EcID.ValueString(),
		SrcCgwID:   config.SrcCgwID.ValueString(),
		DstCgwID:   config.DstCgwID.ValueString(),
		PacketID:   config.PacketID.ValueString(),
		Rate:       config.Rate.ValueInt32(),
		RouteLearn: config.RouteLearn.ValueInt32Pointer(),
	}
	resp, err := c.meta.Apis.SdkEcApis.EcCreateRegionPeerApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建跨域连接失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return err
	} else if resp.ReturnObj == nil || resp.ReturnObj.PeerID == "" {
		err = common.InvalidReturnObjError
		return err
	}
	config.ID = types.StringValue(resp.ReturnObj.PeerID)
	return nil
}

func (c *CtyunExpressConnectRegionPeer) getAndMerge(ctx context.Context, config *CtyunExpressConnectRegionPeerConfig) error {
	regionPeerList, err := c.getRegionPeerList(ctx, config)
	if err != nil {
		return err
	}

	for _, regionPeer := range regionPeerList {
		if *regionPeer.PeerID == config.ID.ValueString() {
			config.Name = types.StringValue(*regionPeer.PeerName)
			config.EcID = types.StringValue(*regionPeer.EcID)
			config.SrcCgwID = types.StringValue(*regionPeer.SrcCgwID)
			config.DstCgwID = types.StringValue(*regionPeer.DstCgwID)
			config.PacketID = types.StringValue(*regionPeer.PacketID)
			config.Rate = types.Int32Value(*regionPeer.Rate)
			config.SrcRegionID = types.StringValue(*regionPeer.SrcDcID)
			config.DstRegionID = types.StringValue(*regionPeer.DstDcID)
			config.PeerType = types.Int32Value(*regionPeer.PeerType)
			config.UpdateTime = types.StringValue(utils.FromBJTimeToUTCZ(*regionPeer.UpdateDate))
			return nil
		}
	}
	return common.ResourceNotExistError
}

func (c *CtyunExpressConnectRegionPeer) getRegionPeerList(ctx context.Context, config *CtyunExpressConnectRegionPeerConfig) ([]*ec.EcRegionPeerListReturnObjResultsResponse, error) {
	params := &ec.EcRegionPeerListRequest{
		EcID:     config.EcID.ValueString(),
		PacketID: config.PacketID.ValueStringPointer(),
		CgwID:    config.SrcCgwID.ValueStringPointer(),
	}
	resp, err := c.meta.Apis.SdkEcApis.EcRegionPeerListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询云间高速(id=%s)的跨域连接列表失败，接口返回nil，请联系研发确认问题原因！", config.EcID.ValueString())
		return nil, err
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp.ReturnObj.Results, nil
}

func (c *CtyunExpressConnectRegionPeer) update(ctx context.Context, state *CtyunExpressConnectRegionPeerConfig, plan *CtyunExpressConnectRegionPeerConfig) error {
	if plan.Rate.Equal(state.Rate) {
		return nil
	}
	params := &ec.EcRegionPeerUpdateRequest{
		PeerID: state.ID.ValueString(),
		Rate:   plan.Rate.ValueInt32(),
	}
	resp, err := c.meta.Apis.SdkEcApis.EcRegionPeerUpdateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新跨域连接id=%s带宽大小失败，接口返回nil，请联系研发确认问题原因！", state.ID.ValueString())
		return err
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	return nil
}

func (c *CtyunExpressConnectRegionPeer) delete(ctx context.Context, config CtyunExpressConnectRegionPeerConfig) error {
	params := ec.EcDeleteRegionPeerRequest{
		PeerID: config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkEcApis.EcDeleteRegionPeerApi.Do(ctx, c.meta.SdkCredential, &params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除跨域连接失败(id=%s)，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return err
	}
	return nil
}

type CtyunExpressConnectRegionPeerConfig struct {
	Name        types.String `tfsdk:"name"`
	EcID        types.String `tfsdk:"ec_id"`
	SrcCgwID    types.String `tfsdk:"src_cgw_id"`
	DstCgwID    types.String `tfsdk:"dst_cgw_id"`
	PacketID    types.String `tfsdk:"packet_id"`
	Rate        types.Int32  `tfsdk:"rate"`
	RouteLearn  types.Int32  `tfsdk:"route_learn"`
	ID          types.String `tfsdk:"id"`
	SrcRegionID types.String `tfsdk:"src_region_id"`
	DstRegionID types.String `tfsdk:"dst_region_id"`
	PeerType    types.Int32  `tfsdk:"peer_type"`
	UpdateTime  types.String `tfsdk:"update_time"`
}
