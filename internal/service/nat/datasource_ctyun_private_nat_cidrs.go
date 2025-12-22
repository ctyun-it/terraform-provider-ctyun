package nat

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctnat"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunPrivateNatCidrsDatasource{}
	_ datasource.DataSourceWithConfigure = &ctyunPrivateNatCidrsDatasource{}
)

type ctyunPrivateNatCidrsDatasource struct {
	meta *common.CtyunMetadata
}

func NewCtyunPrivateNatCidrs() datasource.DataSource {
	return &ctyunPrivateNatCidrsDatasource{}
}

func (c *ctyunPrivateNatCidrsDatasource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_private_nat_cidrs"
}

func (c *ctyunPrivateNatCidrsDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026759/10166345`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填这默认使用provider ctyun总region_id 或者环境变量",
			},
			"nat_gateway_id": schema.StringAttribute{
				Required:    true,
				Description: "要查询的私网NAT网关的ID",
			},
			"cidrs": schema.ListNestedAttribute{
				Computed:    true,
				Description: "中转网段列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "中转网段名称",
						},
						"cidr": schema.StringAttribute{
							Computed:    true,
							Description: "对应网段",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunPrivateNatCidrsDatasource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunPrivateNatCidrsConfig
	// 读取请求信息
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	// region_id不能为空
	if regionId == "" {
		msg := "regionID不能为空"
		response.Diagnostics.AddError(msg, msg)
		return
	}
	natGatewayId := config.NatGatewayID.ValueString()
	params := &ctnat.CtnatListPrivatenatCidrsRequest{
		RegionID:     regionId,
		NatGatewayID: natGatewayId,
		PageNumber:   1,
		PageSize:     50,
	}

	resp, err := c.meta.Apis.SdkCtNatApis.CtnatListPrivatenatCidrsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var cidrs []CtyunPrivateNatCidrModel
	for _, cidr := range resp.ReturnObj {
		cidrItem := CtyunPrivateNatCidrModel{
			Name: types.StringValue(cidr.Name),
			Cidr: types.StringValue(cidr.Cidr),
		}
		cidrs = append(cidrs, cidrItem)
	}

	config.RegionID = types.StringValue(regionId)
	config.NatGatewayID = types.StringValue(natGatewayId)
	config.Cidrs = cidrs
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunPrivateNatCidrsDatasource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunPrivateNatCidrsConfig struct {
	RegionID     types.String               `tfsdk:"region_id"`
	NatGatewayID types.String               `tfsdk:"nat_gateway_id"`
	Cidrs        []CtyunPrivateNatCidrModel `tfsdk:"cidrs"`
}

type CtyunPrivateNatCidrModel struct {
	Name types.String `tfsdk:"name"` // 中转网段名称
	Cidr types.String `tfsdk:"cidr"` // 对应网段
}
