package ec

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ctyunExpressConnects struct {
	meta *common.CtyunMetadata
}

func NewCtyunExpressConnects() datasource.DataSource {
	return &ctyunExpressConnects{}
}

func (c *ctyunExpressConnects) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_express_connects"
}

func (c *ctyunExpressConnects) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026763/10038220`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "云间高速实例ID",
			},
			"query_content": schema.StringAttribute{
				Optional:    true,
				Description: "模糊匹配，支持ecID,ecName, ecDescription三个属性",
			},
			"page_no": schema.Int64Attribute{
				Optional:    true,
				Description: "页码，从1开始",
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "每页记录数目",
			},
			"express_connects": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "云间高速实例ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述信息",
						},
						"status": schema.Int64Attribute{
							Computed:    true,
							Description: "运行状态，取值范围: 1:不可用 2:可用",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"vrf": schema.Int64Attribute{
							Computed:    true,
							Description: "云间高速vrf信息",
						},
						"email": schema.StringAttribute{
							Computed:    true,
							Description: "email",
						},
						"vpc_count": schema.Int64Attribute{
							Computed:    true,
							Description: "添加vpc网络实例的数量",
						},
						"cgw_count": schema.Int64Attribute{
							Computed:    true,
							Description: "云网关数量",
						},
						"cda_count": schema.Int64Attribute{
							Computed:    true,
							Description: "专线数量",
						},
						"sdwan_count": schema.Int64Attribute{
							Computed:    true,
							Description: "sdwan实例数量",
						},
						"vpn_count": schema.Int64Attribute{
							Computed:    true,
							Description: "vpn实例数量",
						},
						"eds_count": schema.Int64Attribute{
							Computed:    true,
							Description: "云桌面网络实例数量",
						},
						"project": schema.StringAttribute{
							Computed:    true,
							Description: "企业项目",
						},
						"packet_status": schema.Int64Attribute{
							Computed:    true,
							Description: "带宽包状态，取值范围: 1:已购买 0:未购买",
						},
						"packet_rate": schema.Int64Attribute{
							Computed:    true,
							Description: "带宽包总量",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunExpressConnects) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CtyunExpressConnectsConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := &ec.EcEcListRequest{}

	// 设置查询参数
	if !config.Id.IsNull() {
		id := config.Id.ValueString()
		request.EcID = &id
	}

	if !config.QueryContent.IsNull() {
		queryContent := config.QueryContent.ValueString()
		request.QueryContent = &queryContent
	}

	pageNo := int32(1)
	if !config.PageNo.IsNull() {
		pageNo = int32(config.PageNo.ValueInt64())
	}
	request.PageNo = &pageNo

	pageSize := int32(10)
	if !config.PageSize.IsNull() {
		pageSize = int32(config.PageSize.ValueInt64())
	}
	request.PageSize = &pageSize

	response, err := c.meta.Apis.SdkEcApis.EcEcListApi.Do(ctx, c.meta.SdkCredential, request)
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

	var expressConnects []CtyunExpressConnectsExpressConnectsConfig
	for _, r := range response.ReturnObj.Results {
		if r.EcID == nil {
			resp.Diagnostics.AddError("API return error", "EcID is nil")
			return
		}

		if r.EcName == nil {
			resp.Diagnostics.AddError("API return error", "EcName is nil")
			return
		}

		if r.Status == nil {
			resp.Diagnostics.AddError("API return error", "Status is nil")
			return
		}

		if r.CreateDate == nil {
			resp.Diagnostics.AddError("API return error", "CreateDate is nil")
			return
		}

		if r.Vrf == nil {
			resp.Diagnostics.AddError("API return error", "Vrf is nil")
			return
		}

		if r.Email == nil {
			resp.Diagnostics.AddError("API return error", "Email is nil")
			return
		}

		if r.VpcCount == nil {
			resp.Diagnostics.AddError("API return error", "VpcCount is nil")
			return
		}

		if r.CgwCount == nil {
			resp.Diagnostics.AddError("API return error", "CgwCount is nil")
			return
		}

		if r.CdaCount == nil {
			resp.Diagnostics.AddError("API return error", "CdaCount is nil")
			return
		}

		if r.SdwanCount == nil {
			resp.Diagnostics.AddError("API return error", "SdwanCount is nil")
			return
		}

		if r.VpnCount == nil {
			resp.Diagnostics.AddError("API return error", "VpnCount is nil")
			return
		}

		if r.EdsCount == nil {
			resp.Diagnostics.AddError("API return error", "EdsCount is nil")
			return
		}

		if r.Project == nil {
			resp.Diagnostics.AddError("API return error", "Project is nil")
			return
		}

		if r.PacketStatus == nil {
			resp.Diagnostics.AddError("API return error", "PacketStatus is nil")
			return
		}

		if r.PacketRate == nil {
			resp.Diagnostics.AddError("API return error", "PacketRate is nil")
			return
		}

		expressConnect := CtyunExpressConnectsExpressConnectsConfig{
			Id:           types.StringValue(*r.EcID),
			Name:         types.StringValue(*r.EcName),
			Status:       types.Int64Value(int64(*r.Status)),
			CreateTime:   types.StringValue(*r.CreateDate),
			Vrf:          types.Int64Value(int64(*r.Vrf)),
			Email:        types.StringValue(*r.Email),
			VpcCount:     types.Int64Value(int64(*r.VpcCount)),
			CgwCount:     types.Int64Value(int64(*r.CgwCount)),
			CdaCount:     types.Int64Value(int64(*r.CdaCount)),
			SdwanCount:   types.Int64Value(int64(*r.SdwanCount)),
			VpnCount:     types.Int64Value(int64(*r.VpnCount)),
			EdsCount:     types.Int64Value(int64(*r.EdsCount)),
			Project:      types.StringValue(*r.Project),
			PacketStatus: types.Int64Value(int64(*r.PacketStatus)),
			PacketRate:   types.Int64Value(int64(*r.PacketRate)),
		}

		if r.EcDescription != nil {
			expressConnect.Description = types.StringValue(*r.EcDescription)
		}

		expressConnects = append(expressConnects, expressConnect)
	}

	config.ExpressConnects = expressConnects
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (c *ctyunExpressConnects) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunExpressConnectsExpressConnectsConfig struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Status       types.Int64  `tfsdk:"status"`
	CreateTime   types.String `tfsdk:"create_time"`
	Vrf          types.Int64  `tfsdk:"vrf"`
	Email        types.String `tfsdk:"email"`
	VpcCount     types.Int64  `tfsdk:"vpc_count"`
	CgwCount     types.Int64  `tfsdk:"cgw_count"`
	CdaCount     types.Int64  `tfsdk:"cda_count"`
	SdwanCount   types.Int64  `tfsdk:"sdwan_count"`
	VpnCount     types.Int64  `tfsdk:"vpn_count"`
	EdsCount     types.Int64  `tfsdk:"eds_count"`
	Project      types.String `tfsdk:"project"`
	PacketStatus types.Int64  `tfsdk:"packet_status"`
	PacketRate   types.Int64  `tfsdk:"packet_rate"`
}

type CtyunExpressConnectsConfig struct {
	Id              types.String                                `tfsdk:"id"`
	QueryContent    types.String                                `tfsdk:"query_content"`
	PageNo          types.Int64                                 `tfsdk:"page_no"`
	PageSize        types.Int64                                 `tfsdk:"page_size"`
	ExpressConnects []CtyunExpressConnectsExpressConnectsConfig `tfsdk:"express_connects"`
}
