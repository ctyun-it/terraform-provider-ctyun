package nat

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctnat"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunPrivateNats{}
	_ datasource.DataSourceWithConfigure = &ctyunPrivateNats{}
)

type ctyunPrivateNats struct {
	meta *common.CtyunMetadata
}

func NewCtyunPrivateNats() datasource.DataSource {
	return &ctyunPrivateNats{}
}

func (c *ctyunPrivateNats) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_private_nats"
}

func (c *ctyunPrivateNats) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026759/10033140",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填这默认使用provider ctyun总region_id 或者环境变量",
			},
			"nat_gateway_id": schema.StringAttribute{
				Optional:    true,
				Description: "要查询的私有NAT网关的ID",
			},
			"page_number": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码，默认值为1",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "分页查询时每页的行数，最大值为50，默认值为10。",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
					int32validator.AtMost(50),
				},
			},
			"nats": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "私有nat 网关 id",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "私有nat 网关名字",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "私有nat网关描述",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟私有云 id",
						},
						"subnet_id": schema.StringAttribute{
							Computed:    true,
							Description: "子网 id",
						},
						"subnet_name": schema.StringAttribute{
							Computed:    true,
							Description: "子网名称",
						},
						"nat_gateway_id": schema.StringAttribute{
							Computed:    true,
							Description: "私有nat 网关 id,与上面的id重复",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "私有NAT网关运行状态: running 表示运行中, freeze 表示冻结, expired 表示已过期",
						},
						"spec": schema.StringAttribute{
							Computed:    true,
							Description: "规格取值: small, medium, large, xlarge",
						},
						"vpc_name": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟私有云名字",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "项目 ID",
						},
						"project_name": schema.StringAttribute{
							Computed:    true,
							Description: "项目名称",
						},
						"az_id": schema.StringAttribute{
							Computed:    true,
							Description: "可用区ID",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunPrivateNats) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// ParseInt32IfEmpty 自定义方法，用于判断types.Int32类型字段是否为空，若为空则返回默认值。
func (c *ctyunPrivateNats) ParseInt32IfEmpty(value types.Int32, defaultValue int32) int32 {
	if value.IsNull() {
		return defaultValue
	}
	return value.ValueInt32()
}

func (c *ctyunPrivateNats) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunPrivateNatsConfig
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

	natGatewayId := config.NatGatewayID.ValueString()

	// todo: pageNo和pageSize判空兜底方案
	pageNo := c.ParseInt32IfEmpty(config.PageNo, 1)

	pageNumber := pageNo

	pageSize := c.ParseInt32IfEmpty(config.PageSize, 10)

	params := &ctnat.CtnatListPrivatenatRequest{
		RegionID:   regionId,
		PageNumber: pageNumber,
		PageNo:     pageNo,
		PageSize:   pageSize,
	}

	if natGatewayId != "" {
		params.NatGatewayID = natGatewayId
	}

	resp, err := c.meta.Apis.SdkCtNatApis.CtnatListPrivatenatApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}

	// 解析返回值，即使ReturnObj为nil也要正确处理
	var nats []CtyunPrivateNatsModel
	if resp.ReturnObj != nil {
		for _, natObj := range resp.ReturnObj {
			natItem := CtyunPrivateNatsModel{
				ID:           types.StringValue(natObj.NatGatewayID),
				Name:         types.StringValue(natObj.Name),
				Description:  types.StringValue(natObj.Description),
				VpcID:        types.StringValue(natObj.VpcID),
				SubnetID:     types.StringValue(natObj.SubnetID),
				SubnetName:   types.StringValue(natObj.SubnetName),
				NatGatewayID: types.StringValue(natObj.NatGatewayID),
				State:        types.StringValue(natObj.State),
				Spec:         types.StringValue(natObj.Spec),
				VpcName:      types.StringValue(natObj.VpcName),
				ProjectID:    types.StringValue(natObj.ProjectID),
				ProjectName:  types.StringValue(natObj.ProjectName),
				AzID:         types.StringValue(natObj.AzID),
				CreateDate:   types.StringValue(natObj.CreateDate),
			}
			nats = append(nats, natItem)
		}
	}

	config.RegionID = types.StringValue(regionId)
	config.Nats = nats
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunPrivateNatsConfig struct {
	RegionID     types.String            `tfsdk:"region_id"`      //区域id
	NatGatewayID types.String            `tfsdk:"nat_gateway_id"` //要查询的私有NAT网关的ID。
	PageNumber   types.Int32             `tfsdk:"page_number"`    //	列表的页码，默认值为1。
	PageNo       types.Int32             `tfsdk:"page_no"`        //列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃
	PageSize     types.Int32             `tfsdk:"page_size"`      //分页查询时每页的行数，最大值为50，默认值为10。
	Nats         []CtyunPrivateNatsModel `tfsdk:"nats"`           // 获取的私有nat列表
}

type CtyunPrivateNatsModel struct {
	ID           types.String `tfsdk:"id"`             //私有nat网关id
	Name         types.String `tfsdk:"name"`           //私有nat网关名称
	Description  types.String `tfsdk:"description"`    //私有nat网关描述
	VpcID        types.String `tfsdk:"vpc_id"`         //虚拟私有云 id
	SubnetID     types.String `tfsdk:"subnet_id"`      //子网 id
	SubnetName   types.String `tfsdk:"subnet_name"`    //子网名称
	NatGatewayID types.String `tfsdk:"nat_gateway_id"` //私有nat网关 id
	State        types.String `tfsdk:"state"`          //私有NAT网关运行状态: running 表示运行中, freeze 表示冻结, expired 表示已过期
	Spec         types.String `tfsdk:"spec"`           //规格取值: small, medium, large, xlarge
	VpcName      types.String `tfsdk:"vpc_name"`       //虚拟私有云名字
	ProjectID    types.String `tfsdk:"project_id"`     //项目 ID
	ProjectName  types.String `tfsdk:"project_name"`   //项目名称
	AzID         types.String `tfsdk:"az_id"`          //可用区ID
	CreateDate   types.String `tfsdk:"create_time"`    //创建时间
}
