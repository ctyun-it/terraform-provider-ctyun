package peer_connection

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &CtyunVpcPeerConnectionAttach{}
	_ resource.ResourceWithConfigure = &CtyunVpcPeerConnectionAttach{}
)

type CtyunVpcPeerConnectionAttach struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunVpcPeerConnectionAttach) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc_peer_connection_attach"
}

func (c *CtyunVpcPeerConnectionAttach) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunVpcPeerConnectionAttach() resource.Resource {
	return &CtyunVpcPeerConnectionAttach{}
}

func (c *CtyunVpcPeerConnectionAttach) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026760/10037761",
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
			"peer_connection_id": schema.StringAttribute{
				Required:    true,
				Description: "对接连接id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"operation": schema.StringAttribute{
				Required:    true,
				Description: "同意或拒绝取消范围：[enable, disable]",
				Validators: []validator.String{
					stringvalidator.OneOf("enable", "disable"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "id",
			},
		},
	}
}

func (c *CtyunVpcPeerConnectionAttach) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunVpcPeerConnectionAttachConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunVpcPeerConnectionAttach) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	return
}

func (c *CtyunVpcPeerConnectionAttach) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunVpcPeerConnectionAttach) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	return
}

func (c *CtyunVpcPeerConnectionAttach) create(ctx context.Context, config *CtyunVpcPeerConnectionAttachConfig) error {
	if config.Operation.ValueString() == "enable" {
		err := c.agree(ctx, config)
		if err != nil {
			return err
		}
	} else if config.Operation.ValueString() == "disable" {
		err := c.reject(ctx, config)
		if err != nil {
			return err
		}
	} else {
		err := fmt.Errorf("operation 非法输入，取消范围：[enable, disable]")
		return err
	}
	config.ID = types.StringValue(uuid.NewString())
	return nil
}

func (c *CtyunVpcPeerConnectionAttach) agree(ctx context.Context, config *CtyunVpcPeerConnectionAttachConfig) error {
	params := &ctvpc.CtvpcAgreeVpcPeerRequestRequest{
		InstanceID: config.PeerConnectionID.ValueString(),
		RegionID:   config.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcAgreeVpcPeerRequestApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("同意建立对等连接失败，对等连接id=%s，接口返回nil，请联系研发确认问题原因！", config.PeerConnectionID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return err
	}
	return nil
}

func (c *CtyunVpcPeerConnectionAttach) reject(ctx context.Context, config *CtyunVpcPeerConnectionAttachConfig) error {
	params := &ctvpc.CtvpcRejectVpcPeerRequestRequest{
		InstanceID: config.PeerConnectionID.ValueString(),
		RegionID:   config.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcRejectVpcPeerRequestApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("拒绝建立对等连接失败，对等连接id=%s，接口返回nil，请联系研发确认问题原因！", config.PeerConnectionID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return err
	}
	return nil
}

type CtyunVpcPeerConnectionAttachConfig struct {
	RegionID         types.String `tfsdk:"region_id"`
	PeerConnectionID types.String `tfsdk:"peer_connection_id"`
	Operation        types.String `tfsdk:"operation"`
	ID               types.String `tfsdk:"id"`
}
