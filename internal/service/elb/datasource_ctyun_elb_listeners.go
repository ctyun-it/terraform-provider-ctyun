package elb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunElbListeners{}
	_ datasource.DataSourceWithConfigure = &ctyunElbListeners{}
)

type ctyunElbListeners struct {
	meta *common.CtyunMetadata
}

func NewElbListeners() datasource.DataSource {
	return &ctyunElbListeners{}
}

func (c *ctyunElbListeners) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunElbListeners) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_listeners"
}

func (c *ctyunElbListeners) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026756/10140276**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "企业项目ID，默认为'0'",
			},
			"ids": schema.StringAttribute{
				Optional:    true,
				Description: "监听器ID列表，以','分隔",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "监听器名称",
			},
			"load_balancer_id": schema.StringAttribute{
				Optional:    true,
				Description: "负载均衡实例ID",
			},
			"access_control_id": schema.StringAttribute{
				Optional:    true,
				Description: "访问控制ID",
			},
			"listeners": schema.ListNestedAttribute{
				Computed:    true,
				Description: "监听器列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源池ID",
						},
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "可用区名称",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "项目ID",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "监听器ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "监听器名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"load_balancer_id": schema.StringAttribute{
							Computed:    true,
							Description: "负载均衡实例ID",
						},
						"protocol": schema.StringAttribute{
							Computed:    true,
							Description: "监听协议: TCP / UDP / HTTP / HTTPS",
						},
						"protocol_port": schema.Int32Attribute{
							Computed:    true,
							Description: "监听端口",
						},
						"certificate_id": schema.StringAttribute{
							Computed:    true,
							Description: "证书ID",
						},
						"ca_enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "是否开启双向认证",
						},
						"client_certificate_id": schema.StringAttribute{
							Computed:    true,
							Description: "双向认证的证书ID",
						},
						"default_action_type": schema.StringAttribute{
							Computed:    true,
							Description: "默认规则动作类型: forward / redirect",
						},
						"forward_config": schema.ListNestedAttribute{
							Computed:    true,
							Description: "转发配置",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"target_group_id": schema.StringAttribute{
										Required:    true,
										Description: "后端服务组ID",
									},
									"weight": schema.Int32Attribute{
										Optional:    true,
										Description: "权重，取值范围：1-256。默认为100",
										Validators: []validator.Int32{
											int32validator.Between(1, 256),
										},
									},
								},
							},
						},
						"redirect_listener_id": schema.StringAttribute{
							Computed:    true,
							Description: "重定向监听器ID",
						},
						"access_control_id": schema.StringAttribute{
							Computed:    true,
							Description: "访问控制ID",
						},
						"access_control_type": schema.StringAttribute{
							Computed:    true,
							Description: "访问控制类型: Close / White / Black",
						},
						"forwarded_for_enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "是否开启x forward for功能",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "监听器状态: DOWN / ACTIVE",
						},
						"created_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"updated_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间，为UTC格式",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunElbListeners) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunElbListenersConfig
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	params := &ctelb.CtelbListListenerRequest{
		ClientToken: uuid.NewString(),
		RegionID:    regionId,
	}

	if !config.ProjectID.IsNull() {
		params.ProjectID = config.ProjectID.ValueString()
	}
	if !config.IDs.IsNull() {
		params.IDs = config.IDs.ValueString()
	}
	if !config.Name.IsNull() {
		params.Name = config.Name.ValueString()
	}
	if !config.LoadBalancerID.IsNull() {
		params.LoadBalancerID = config.LoadBalancerID.ValueString()
	}
	if !config.AccessControlID.IsNull() {
		params.AccessControlID = config.AccessControlID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbListListenerApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回值
	var listeners []CtyunElbListenersDetailModel
	for _, listenerItem := range resp.ReturnObj {
		var listener CtyunElbListenersDetailModel
		listener.RegionID = types.StringValue(listenerItem.RegionID)
		listener.AzName = types.StringValue(listenerItem.AzName)
		listener.ID = types.StringValue(listenerItem.ID)
		listener.Name = types.StringValue(listenerItem.Name)
		listener.Description = types.StringValue(listenerItem.Description)
		listener.LoadBalancerID = types.StringValue(listenerItem.LoadBalancerID)
		listener.Protocol = types.StringValue(listenerItem.Protocol)
		listener.ProtocolPort = types.Int32Value(listenerItem.ProtocolPort)
		listener.CertificateID = types.StringValue(listenerItem.CertificateID)
		listener.AccessControlID = types.StringValue(listenerItem.AccessControlID)
		listener.AccessControlType = types.StringValue(listenerItem.AccessControlType)
		listener.ForwardedForEnabled = types.BoolValue(*listenerItem.ForwardedForEnabled)
		listener.Status = types.StringValue(listenerItem.Status)
		listener.CreatedTime = types.StringValue(listenerItem.CreatedTime)
		listener.UpdatedTime = types.StringValue(listenerItem.UpdatedTime)
		// 处理defaultAction
		listener.DefaultActionType = types.StringValue(listenerItem.DefaultAction.RawType)

		var targetGroups []TargetGroupModel
		if listenerItem.DefaultAction.ForwardConfig.TargetGroups != nil && len(listenerItem.DefaultAction.ForwardConfig.TargetGroups) > 0 {
			for _, targetGroupItem := range listenerItem.DefaultAction.ForwardConfig.TargetGroups {
				var targetConfig TargetGroupModel
				targetConfig.TargetGroupID = types.StringValue(targetGroupItem.TargetGroupID)
				targetConfig.Weight = types.Int32Value(targetGroupItem.Weight)
				targetGroups = append(targetGroups, targetConfig)
			}
		}

		listener.ForwardConfig = targetGroups
		listener.RedirectListenerID = types.StringValue(listenerItem.DefaultAction.RedirectListenerID)

		listeners = append(listeners, listener)
	}
	config.Listeners = listeners
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunElbListenersConfig struct {
	RegionID        types.String                   `tfsdk:"region_id"`
	ProjectID       types.String                   `tfsdk:"project_id"`
	IDs             types.String                   `tfsdk:"ids"`
	Name            types.String                   `tfsdk:"name"`
	LoadBalancerID  types.String                   `tfsdk:"load_balancer_id"`
	AccessControlID types.String                   `tfsdk:"access_control_id"`
	Listeners       []CtyunElbListenersDetailModel `tfsdk:"listeners"`
}

type CtyunElbListenersDetailModel struct {
	RegionID            types.String       `tfsdk:"region_id"`
	AzName              types.String       `tfsdk:"az_name"`
	ProjectID           types.String       `tfsdk:"project_id"`
	ID                  types.String       `tfsdk:"id"`
	Name                types.String       `tfsdk:"name"`
	Description         types.String       `tfsdk:"description"`
	LoadBalancerID      types.String       `tfsdk:"load_balancer_id"`
	Protocol            types.String       `tfsdk:"protocol"`
	ProtocolPort        types.Int32        `tfsdk:"protocol_port"`
	CertificateID       types.String       `tfsdk:"certificate_id"`
	CaEnabled           types.Bool         `tfsdk:"ca_enabled"` //是否开启双向认证
	ClientCertificateID types.String       `tfsdk:"client_certificate_id"`
	DefaultActionType   types.String       `tfsdk:"default_action_type"`
	ForwardConfig       []TargetGroupModel `tfsdk:"forward_config"`
	RedirectListenerID  types.String       `tfsdk:"redirect_listener_id"`
	AccessControlID     types.String       `tfsdk:"access_control_id"`
	AccessControlType   types.String       `tfsdk:"access_control_type"`
	ForwardedForEnabled types.Bool         `tfsdk:"forwarded_for_enabled"`
	Status              types.String       `tfsdk:"status"`
	CreatedTime         types.String       `tfsdk:"created_time"`
	UpdatedTime         types.String       `tfsdk:"updated_time"`
}
