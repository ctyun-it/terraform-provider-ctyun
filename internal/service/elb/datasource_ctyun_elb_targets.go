package elb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunElbTargets{}
	_ datasource.DataSourceWithConfigure = &ctyunElbTargets{}
)

type ctyunElbTargets struct {
	meta *common.CtyunMetadata
}

func NewCtyunElbTargets() datasource.DataSource {
	return &ctyunElbTargets{}
}

func (c *ctyunElbTargets) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunElbTargets) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_targets"
}

func (c *ctyunElbTargets) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026756/10196689**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"target_group_id": schema.StringAttribute{
				Optional:    true,
				Description: "后端服务组ID",
			},
			"ids": schema.StringAttribute{
				Optional:    true,
				Description: "后端服务ID列表，以,分隔",
			},
			"elb_targets": schema.ListNestedAttribute{
				Computed:    true,
				Description: "后端主机列表",
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
							Description: "后端服务ID",
						},
						"target_group_id": schema.StringAttribute{
							Computed:    true,
							Description: "后端服务组ID",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"instance_type": schema.StringAttribute{
							Computed:    true,
							Description: "实例类型: VM / BM",
							Validators: []validator.String{
								stringvalidator.OneOf(business.ElbTargetType...),
							},
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "实例ID",
						},
						"protocol_port": schema.Int32Attribute{
							Computed:    true,
							Description: "协议端口",
						},
						"weight": schema.Int32Attribute{
							Computed:    true,
							Description: "权重",
						},
						"health_check_status": schema.StringAttribute{
							Computed:    true,
							Description: "IPv4的健康检查状态: offline / online / unknown",
							Validators: []validator.String{
								stringvalidator.OneOf(business.ElbTargetIpStatus...),
							},
						},
						"health_check_status_ipv6": schema.StringAttribute{
							Computed:    true,
							Description: "IPv6的健康检查状态: offline / online / unknown",
							Validators: []validator.String{
								stringvalidator.OneOf(business.ElbTargetIpStatus...),
							},
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "状态: DOWN / ACTIVE",
							Validators: []validator.String{
								stringvalidator.OneOf(business.ElbRuleStatus...),
							},
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

func (c *ctyunElbTargets) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunElbTargetsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		msg := "regionID不能为空"
		response.Diagnostics.AddError(msg, msg)
		return
	}
	params := &ctelb.CtelbListTargetRequest{
		RegionID: regionId,
	}
	if !config.TargetGroupID.IsNull() {
		params.TargetGroupID = config.TargetGroupID.ValueString()
	}
	if !config.IDs.IsNull() {
		params.IDs = config.IDs.ValueString()
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbListTargetApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回值
	var elbTargets []ElbTargetModel
	for _, elbTargetItem := range resp.ReturnObj {
		var elbTarget ElbTargetModel
		elbTarget.RegionID = types.StringValue(elbTargetItem.RegionID)
		elbTarget.AzName = types.StringValue(elbTargetItem.AzName)
		elbTarget.ProjectID = types.StringValue(elbTargetItem.ProjectID)
		elbTarget.ID = types.StringValue(elbTargetItem.ID)
		elbTarget.TargetGroupID = types.StringValue(elbTargetItem.TargetGroupID)
		elbTarget.Description = types.StringValue(elbTargetItem.Description)
		elbTarget.InstanceType = types.StringValue(elbTargetItem.InstanceType)
		elbTarget.InstanceID = types.StringValue(elbTargetItem.InstanceID)
		elbTarget.ProtocolPort = types.Int32Value(elbTargetItem.ProtocolPort)
		elbTarget.Weight = types.Int32Value(elbTargetItem.Weight)
		elbTarget.HealthCheckStatus = types.StringValue(elbTargetItem.HealthCheckStatus)
		elbTarget.HealthCheckStatusIpv6 = types.StringValue(elbTargetItem.HealthCheckStatusIpv6)
		elbTarget.Status = types.StringValue(elbTargetItem.Status)
		elbTarget.CreatedTime = types.StringValue(elbTargetItem.CreatedTime)
		elbTarget.UpdatedTime = types.StringValue(elbTargetItem.UpdatedTime)
		elbTargets = append(elbTargets, elbTarget)
	}
	config.ElbTargets = elbTargets
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunElbTargetsConfig struct {
	RegionID      types.String     `tfsdk:"region_id"`       //区域ID
	TargetGroupID types.String     `tfsdk:"target_group_id"` //	后端服务组ID
	IDs           types.String     `tfsdk:"ids"`             //后端服务ID列表，以,分隔
	ElbTargets    []ElbTargetModel `tfsdk:"elb_targets"`
}

type ElbTargetModel struct {
	RegionID              types.String `tfsdk:"region_id"`                //区域ID
	AzName                types.String `tfsdk:"az_name"`                  //可用区名称
	ProjectID             types.String `tfsdk:"project_id"`               //项目ID
	ID                    types.String `tfsdk:"id"`                       //后端服务ID
	TargetGroupID         types.String `tfsdk:"target_group_id"`          //后端服务组ID
	Description           types.String `tfsdk:"description"`              //描述
	InstanceType          types.String `tfsdk:"instance_type"`            //实例类型: VM / BM
	InstanceID            types.String `tfsdk:"instance_id"`              //实例ID
	ProtocolPort          types.Int32  `tfsdk:"protocol_port"`            //	协议端口
	Weight                types.Int32  `tfsdk:"weight"`                   //权重
	HealthCheckStatus     types.String `tfsdk:"health_check_status"`      //IPv4的健康检查状态: offline / online / unknown
	HealthCheckStatusIpv6 types.String `tfsdk:"health_check_status_ipv6"` //IPv6的健康检查状态: offline / online / unknown
	Status                types.String `tfsdk:"status"`                   //状态: DOWN / ACTIVE
	CreatedTime           types.String `tfsdk:"created_time"`             //创建时间，为UTC格式
	UpdatedTime           types.String `tfsdk:"updated_time"`             //更新时间，为UTC格式
}
