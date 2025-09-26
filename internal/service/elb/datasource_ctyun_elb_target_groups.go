package elb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunElbTargetGroups{}
	_ datasource.DataSourceWithConfigure = &CtyunElbTargetGroups{}
)

type CtyunElbTargetGroups struct {
	meta *common.CtyunMetadata
}

func NewCtyunElbTargetGroups() datasource.DataSource {
	return &CtyunElbTargetGroups{}
}

func (c *CtyunElbTargetGroups) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_target_groups"

}

func (c *CtyunElbTargetGroups) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026756/10155289`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"ids": schema.StringAttribute{
				Optional:    true,
				Description: "后端服务组ID列表，以,分隔",
			},
			"vpc_id": schema.StringAttribute{
				Optional:    true,
				Description: "vpcID",
			},
			"health_check_id": schema.StringAttribute{
				Optional:    true,
				Description: "健康检查ID",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "后端服务组名称",
			},
			"target_groups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "后端服务组列表",
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
							Description: "后端服务组ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "后端服务组名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "vpc ID",
						},
						"health_check_id": schema.StringAttribute{
							Computed:    true,
							Description: "健康检查ID",
						},
						"algorithm": schema.StringAttribute{
							Computed:    true,
							Description: "调度算法",
						},
						"session_sticky_mode": schema.StringAttribute{
							Computed:    true,
							Description: "会话保持模式，支持取值：CLOSE（关闭）、INSERT（插入）、REWRITE（重写）",
						},
						"cookie_expire": schema.Int32Attribute{
							Computed:    true,
							Description: "cookie过期时间",
						},
						"rewrite_cookie_name": schema.StringAttribute{
							Computed:    true,
							Description: "cookie重写名称",
						},
						"source_ip_timeout": schema.Int32Attribute{
							Computed:    true,
							Description: "源IP会话保持超时时间",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "状态: DOWN / ACTIVE",
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

func (c *CtyunElbTargetGroups) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunElbTargetGroupsConfig
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

	params := &ctelb.CtelbListTargetGroupRequest{
		ClientToken: uuid.NewString(),
		RegionID:    regionId,
	}
	if !config.IDs.IsNull() {
		params.IDs = config.IDs.ValueString()
	}
	if !config.VpcID.IsNull() {
		params.VpcID = config.VpcID.ValueString()
	}
	if !config.HealthCheckID.IsNull() {
		params.HealthCheckID = config.HealthCheckID.ValueString()
	}
	if !config.Name.IsNull() {
		params.Name = config.Name.ValueString()
	}

	// 请求查看后端主机组列表接口
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbListTargetGroupApi.Do(ctx, c.meta.SdkCredential, params)
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
	var targetGroups []CtyunTargetGroupModel
	for _, targetGroupItem := range resp.ReturnObj {
		var targetGroup CtyunTargetGroupModel
		targetGroup.RegionID = types.StringValue(targetGroupItem.RegionID)
		targetGroup.ID = types.StringValue(targetGroupItem.ID)
		targetGroup.Name = types.StringValue(targetGroupItem.Name)
		targetGroup.AzName = types.StringValue(targetGroupItem.AzName)
		targetGroup.ProjectID = types.StringValue(targetGroupItem.ProjectID)
		targetGroup.Description = types.StringValue(targetGroupItem.Description)
		targetGroup.VpcID = types.StringValue(targetGroupItem.VpcID)
		targetGroup.HealthCheckID = types.StringValue(targetGroupItem.HealthCheckID)
		targetGroup.Algorithm = types.StringValue(targetGroupItem.Algorithm)
		targetGroup.Status = types.StringValue(targetGroupItem.Status)
		targetGroup.CreatedTime = types.StringValue(targetGroupItem.CreatedTime)
		targetGroup.UpdatedTime = types.StringValue(targetGroupItem.UpdatedTime)
		targetGroup.SessionStickyMode = types.StringValue(targetGroupItem.SessionSticky.SessionStickyMode)
		targetGroup.CookieExpire = types.Int32Value(targetGroupItem.SessionSticky.CookieExpire)
		targetGroup.RewriteCookieName = types.StringValue(targetGroupItem.SessionSticky.RewriteCookieName)
		targetGroup.SourceIpTimeout = types.Int32Value(targetGroupItem.SessionSticky.SourceIpTimeout)
		targetGroups = append(targetGroups, targetGroup)
	}
	config.TargetGroups = targetGroups

	// 保存后端主机组列表
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunElbTargetGroups) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunElbTargetGroupsConfig struct {
	RegionID      types.String            `tfsdk:"region_id"`       //	区域ID
	IDs           types.String            `tfsdk:"ids"`             //后端服务组ID列表，以,分隔
	VpcID         types.String            `tfsdk:"vpc_id"`          //vpc ID
	HealthCheckID types.String            `tfsdk:"health_check_id"` //健康检查ID
	Name          types.String            `tfsdk:"name"`            //后端服务组名称
	TargetGroups  []CtyunTargetGroupModel `tfsdk:"target_groups"`   //
}

type CtyunTargetGroupModel struct {
	RegionID          types.String `tfsdk:"region_id"`           //	区域ID
	AzName            types.String `tfsdk:"az_name"`             //可用区名称
	ProjectID         types.String `tfsdk:"project_id"`          //项目ID
	ID                types.String `tfsdk:"id"`                  //后端服务组ID
	Name              types.String `tfsdk:"name"`                //后端服务组名称
	Description       types.String `tfsdk:"description"`         //描述
	VpcID             types.String `tfsdk:"vpc_id"`              //vpc ID
	HealthCheckID     types.String `tfsdk:"health_check_id"`     //健康检查ID
	Algorithm         types.String `tfsdk:"algorithm"`           //调度算法
	SessionStickyMode types.String `tfsdk:"session_sticky_mode"` //会话保持模式，支持取值：CLOSE（关闭）、INSERT（插入）、REWRITE（重写）
	CookieExpire      types.Int32  `tfsdk:"cookie_expire"`       //cookie过期时间
	RewriteCookieName types.String `tfsdk:"rewrite_cookie_name"` //cookie重写名称
	SourceIpTimeout   types.Int32  `tfsdk:"source_ip_timeout"`   //源IP会话保持超时时间
	Status            types.String `tfsdk:"status"`              //状态: DOWN / ACTIVE
	CreatedTime       types.String `tfsdk:"created_time"`        //创建时间，为UTC格式
	UpdatedTime       types.String `tfsdk:"updated_time"`        //更新时间，为UTC格式
}
