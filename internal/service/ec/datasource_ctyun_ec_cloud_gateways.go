package ec

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ctyunEcCloudGateways struct {
	meta *common.CtyunMetadata
}

func NewCtyunEcCloudGateways() datasource.DataSource {
	return &ctyunEcCloudGateways{}
}

func (c *ctyunEcCloudGateways) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ec_cloud_gateways"
}

func (c *ctyunEcCloudGateways) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026763/10038220`,
		Attributes: map[string]schema.Attribute{
			"ec_id": schema.StringAttribute{
				Required:    true,
				Description: "云间高速实例ID",
			},
			"cgw_id": schema.StringAttribute{
				Optional:    true,
				Description: "云网关实例ID",
			},
			"query_content": schema.StringAttribute{
				Optional:    true,
				Description: "模糊匹配，支持cgwID,cgwName,cgwDescription三个属性",
			},
			"region": schema.Int64Attribute{
				Optional:    true,
				Description: "地域信息，不填默认查询全部 取值如下 1：中国大陆 2:亚太",
			},
			"cloud_gateways": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "云网关实例ID",
						},
						"region": schema.Int64Attribute{
							Computed:    true,
							Description: "地域信息，不填默认查询全部 取值如下 1：中国大陆 2:亚太",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "云网关名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "云网关描述",
						},
						"ec_id": schema.StringAttribute{
							Computed:    true,
							Description: "云间高速实例ID",
						},
						"dc_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源池ID信息",
						},
						"dc_type": schema.StringAttribute{
							Computed:    true,
							Description: "资源池类型，取值范围: 'CNP':CNP资源池 'MAZ':MAZ资源池",
						},
						"dc_name": schema.StringAttribute{
							Computed:    true,
							Description: "资源池名称",
						},
						"rtb_count": schema.Int64Attribute{
							Computed:    true,
							Description: "路由表数量",
						},
						"route_count": schema.Int64Attribute{
							Computed:    true,
							Description: "路由数量",
						},
						"policy_count": schema.Int64Attribute{
							Computed:    true,
							Description: "流量策略数量",
						},
						"ins_count": schema.Int64Attribute{
							Computed:    true,
							Description: "连接实例数量",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"rtb_id": schema.StringAttribute{
							Computed:    true,
							Description: "云网关默认路由表ID",
						},
						"has_monitor": schema.BoolAttribute{
							Computed:    true,
							Description: "是否支持监控 true:支持，false:不支持",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunEcCloudGateways) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CtyunEcCloudGatewaysConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := &ec.EcEcListGatewayRequest{
		EcID: config.EcID.ValueString(),
	}

	// 设置查询参数
	if !config.CgwID.IsNull() {
		cgwID := config.CgwID.ValueString()
		request.CgwID = &cgwID
	}

	if !config.QueryContent.IsNull() {
		queryContent := config.QueryContent.ValueString()
		request.QueryContent = &queryContent
	}

	if !config.Region.IsNull() {
		region := fmt.Sprintf("%d", config.Region.ValueInt64())
		request.Region = &region
	}

	response, err := c.meta.Apis.SdkEcApis.EcEcListGatewayApi.Do(ctx, c.meta.SdkCredential, request)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	} else if response.StatusCode == nil {
		resp.Diagnostics.AddError("API return error", "StatusCode is nil")
		return
	} else if *response.StatusCode != common.NormalStatusCode {
		resp.Diagnostics.AddError("API return error", fmt.Sprintf("Message: %s", func() string {
			if response.Message != nil {
				return *response.Message
			}
			return "unknown error"
		}()))
		return
	} else if response.ReturnObj == nil {
		resp.Diagnostics.AddError("API return error", "ReturnObj is nil")
		return
	}

	var cloudGateways []CtyunEcCloudGatewaysCloudGatewaysConfig
	for _, r := range response.ReturnObj.Results {
		if r.CgwID == nil {
			resp.Diagnostics.AddError("API return error", "CgwID is nil")
			return
		}

		if r.CgwName == nil {
			resp.Diagnostics.AddError("API return error", "CgwName is nil")
			return
		}

		if r.EcID == nil {
			resp.Diagnostics.AddError("API return error", "EcID is nil")
			return
		}

		if r.DcID == nil {
			resp.Diagnostics.AddError("API return error", "RegionID is nil")
			return
		}

		if r.DcType == nil {
			resp.Diagnostics.AddError("API return error", "DcType is nil")
			return
		}

		if r.DcName == nil {
			resp.Diagnostics.AddError("API return error", "DcName is nil")
			return
		}

		cloudGateway := CtyunEcCloudGatewaysCloudGatewaysConfig{
			Id:     types.StringValue(*r.CgwID),
			Name:   types.StringValue(*r.CgwName),
			EcID:   types.StringValue(*r.EcID),
			DcID:   types.StringValue(*r.DcID),
			DcType: types.StringValue(*r.DcType),
			DcName: types.StringValue(*r.DcName),
		}

		if r.Region != nil {
			cloudGateway.Region = types.Int64Value(*r.Region)
		}

		if r.CgwDescription != nil {
			cloudGateway.Description = types.StringValue(*r.CgwDescription)
		}

		if r.RtbCnt != nil {
			cloudGateway.RtbCount = types.Int64Value(int64(*r.RtbCnt))
		}

		if r.RouteCnt != nil {
			cloudGateway.RouteCount = types.Int64Value(int64(*r.RouteCnt))
		}

		if r.PolicyCount != nil {
			cloudGateway.PolicyCount = types.Int64Value(int64(*r.PolicyCount))
		}

		if r.InsCnt != nil {
			cloudGateway.InsCount = types.Int64Value(int64(*r.InsCnt))
		}

		if r.CreateDate != nil {
			cloudGateway.CreateTime = types.StringValue(utils.FromBJTimeToUTCZ(*r.CreateDate))
		}

		if r.DefaultRtbID != nil {
			cloudGateway.DefaultRtbID = types.StringValue(*r.DefaultRtbID)
		}

		if r.HasMonitor != nil {
			cloudGateway.HasMonitor = types.BoolValue(*r.HasMonitor)
		}

		cloudGateways = append(cloudGateways, cloudGateway)
	}

	config.CloudGateways = cloudGateways
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (c *ctyunEcCloudGateways) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunEcCloudGatewaysCloudGatewaysConfig struct {
	Id           types.String `tfsdk:"id"`
	Region       types.Int64  `tfsdk:"region"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	EcID         types.String `tfsdk:"ec_id"`
	DcID         types.String `tfsdk:"dc_id"`
	DcType       types.String `tfsdk:"dc_type"`
	DcName       types.String `tfsdk:"dc_name"`
	RtbCount     types.Int64  `tfsdk:"rtb_count"`
	RouteCount   types.Int64  `tfsdk:"route_count"`
	PolicyCount  types.Int64  `tfsdk:"policy_count"`
	InsCount     types.Int64  `tfsdk:"ins_count"`
	CreateTime   types.String `tfsdk:"create_time"`
	DefaultRtbID types.String `tfsdk:"rtb_id"`
	HasMonitor   types.Bool   `tfsdk:"has_monitor"`
}

type CtyunEcCloudGatewaysConfig struct {
	EcID          types.String                              `tfsdk:"ec_id"`
	CgwID         types.String                              `tfsdk:"cgw_id"`
	QueryContent  types.String                              `tfsdk:"query_content"`
	Region        types.Int64                               `tfsdk:"region"`
	CloudGateways []CtyunEcCloudGatewaysCloudGatewaysConfig `tfsdk:"cloud_gateways"`
}
