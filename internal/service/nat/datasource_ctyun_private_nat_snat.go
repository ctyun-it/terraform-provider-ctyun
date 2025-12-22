package nat

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctnat"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunPrivateSNats{}
	_ datasource.DataSourceWithConfigure = &ctyunPrivateSNats{}
)

type ctyunPrivateSNats struct {
	meta *common.CtyunMetadata
}

func (c *ctyunPrivateSNats) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func NewCtyunPrivateSNats() datasource.DataSource {
	return &ctyunPrivateSNats{}
}

func (c *ctyunPrivateSNats) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_private_nat_snats"
}

func (c *ctyunPrivateSNats) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026759/10166268`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填这默认使用provider ctyun总region_id 或者环境变量",
			},
			"nat_gateway_id": schema.StringAttribute{
				Required:    true,
				Description: "私网NAT网关ID",
			},
			"snat_id": schema.StringAttribute{
				Optional:    true,
				Description: "snat id，选填",
			},
			"page_number": schema.Int64Attribute{
				Optional:    true,
				Description: "列表的页码，默认值为1",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "分页查询时每页的行数，最大值为50，默认值为10。",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(50),
				},
			},
			"snats": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"snat_id": schema.StringAttribute{
							Computed:    true,
							Description: "snat id",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"source_cidr": schema.StringAttribute{
							Computed:    true,
							Description: "源地址段",
						},
						"source_vpc_name": schema.StringAttribute{
							Computed:    true,
							Description: "源vpc名称",
						},
						"source_subnet_id": schema.StringAttribute{
							Computed:    true,
							Description: "源Subnet的ID",
						},
						"source_subnet_name": schema.StringAttribute{
							Computed:    true,
							Description: "源Subnet名称",
						},
						"addresses": schema.ListAttribute{
							Computed:    true,
							Description: "中转IP地址",
							ElementType: types.StringType,
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "SNAT状态: running代表运行中, freeze代表已冻结, expired代表已到期",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunPrivateSNats) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunPrivateSNatsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	// region_id必不能为空
	if regionId == "" {
		msg := "regionID不能为空"
		response.Diagnostics.AddError(msg, msg)
		return
	}
	natGatewayId := config.NatGateWayID.ValueString()
	snatId := config.SNatID.ValueString()
	// 分页信息先预留
	pageNumber := c.ParseIntIfEmpty(config.PageNumber, types.Int64Value(1))
	pageSize := c.ParseIntIfEmpty(config.PageSize, types.Int64Value(10))

	params := &ctnat.CtnatQueryPrivatenatSnatRequest{
		RegionID:     regionId,
		NatGatewayID: natGatewayId,
		SnatID:       snatId,
		PageNo:       int32(pageNumber.ValueInt64()),
		PageSize:     int32(pageSize.ValueInt64()),
	}

	resp, err := c.meta.Apis.SdkCtNatApis.CtnatQueryPrivatenatSnatApi.Do(ctx, c.meta.SdkCredential, params)
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
	var snats []CtyunPrivateSNatsModel
	for _, snat := range resp.ReturnObj {
		addressesList, diags := types.ListValueFrom(ctx, types.StringType, snat.Addresses)
		if diags.HasError() {
			err = fmt.Errorf("failed to convert addresses to list")
			return
		}

		snatItem := CtyunPrivateSNatsModel{
			SnatID:           types.StringValue(snat.SnatID),
			Description:      types.StringValue(snat.Description),
			SourceCIDR:       types.StringValue(snat.SrcCIDR),
			SourceVpcName:    types.StringValue(snat.SrcVpcName),
			SourceSubnetID:   types.StringValue(snat.SrcSubnetID),
			SourceSubnetName: types.StringValue(snat.SrcSubnetName),
			Addresses:        addressesList,
			State:            types.StringValue(snat.State),
		}
		snats = append(snats, snatItem)
	}

	config.Snats = snats
	config.RegionID = types.StringValue(regionId)
	config.NatGateWayID = types.StringValue(natGatewayId)
	config.SNatID = types.StringValue(snatId)
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

// ParseIntIfEmpty 自定义方法，用于判断types.Int64类型字段是否为空，若为空则返回默认值。
func (c *ctyunPrivateSNats) ParseIntIfEmpty(value types.Int64, defaultValue types.Int64) types.Int64 {
	if value.IsNull() {
		return defaultValue
	}
	return value
}

type CtyunPrivateSNatsConfig struct {
	RegionID     types.String             `tfsdk:"region_id"`      //区域id
	NatGateWayID types.String             `tfsdk:"nat_gateway_id"` //要查询的私网NAT网关的ID。
	SNatID       types.String             `tfsdk:"snat_id"`        // snat id
	PageNumber   types.Int64              `tfsdk:"page_number"`    //	列表的页码，默认值为1。
	PageSize     types.Int64              `tfsdk:"page_size"`      //分页查询时每页的行数，最大值为50，默认值为10。
	Snats        []CtyunPrivateSNatsModel `tfsdk:"snats"`
}

type CtyunPrivateSNatsModel struct {
	SnatID           types.String `tfsdk:"snat_id"`            //snat id
	Description      types.String `tfsdk:"description"`        //描述信息
	SourceCIDR       types.String `tfsdk:"source_cidr"`        //源地址段
	SourceVpcName    types.String `tfsdk:"source_vpc_name"`    //源vpc名称
	SourceSubnetID   types.String `tfsdk:"source_subnet_id"`   //源Subnet的ID
	SourceSubnetName types.String `tfsdk:"source_subnet_name"` //源Subnet名称
	Addresses        types.List   `tfsdk:"addresses"`          //中转IP地址
	State            types.String `tfsdk:"state"`              //SNAT状态: running代表运行中, freeze代表已冻结, expired代表已到期
}
