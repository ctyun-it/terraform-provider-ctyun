package nat

// ctyun_nats datasource
//

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunNats{}
	_ datasource.DataSourceWithConfigure = &ctyunNats{}
)

type ctyunNats struct {
	meta *common.CtyunMetadata
}

func NewCtyunNats() datasource.DataSource {
	return &ctyunNats{}
}

func (c *ctyunNats) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_nats"
}

func (c *ctyunNats) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026759/10033140`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填这默认使用provider ctyun总region_id 或者环境变量",
			},
			"nat_gateway_id": schema.StringAttribute{
				Optional:    true,
				Description: "要查询的NAT网关的ID",
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
							Description: "nat 网关 id",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "nat 网关名字",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "nat网关描述",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟私有云 id",
						},
						"status": schema.Int32Attribute{
							Computed:    true,
							Description: "nat 网关状态: 0 表示创建中，2 表示运行中，3 表示冻结",
							Validators: []validator.Int32{
								int32validator.Between(0, 3),
							},
						},
						"nat_gateway_id": schema.StringAttribute{
							Computed:    true,
							Description: "nat 网关 id,与上面的id重复",
						},
						"zone_id": schema.StringAttribute{
							Computed:    true,
							Description: "可用区 ID",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "NAT网关运行状态: running 表示运行中, creating 表示创建中, expired 表示已过期, freeze 表示已冻结",
							Validators: []validator.String{
								stringvalidator.OneOf(business.NatStates...),
							},
						},
						"vpc_name": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟私有云名字",
						},
						"expire_time": schema.StringAttribute{
							Computed:    true,
							Description: "过期时间",
						},
						"creation_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "项目 ID",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunNats) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunNatsConfig
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

	params := &ctvpc.CtvpcListNatGatewaysRequest{
		RegionID:   regionId,
		PageNumber: pageNumber,
		PageNo:     pageNo,
		PageSize:   pageSize,
	}

	if natGatewayId != "" {
		params.NatGatewayID = &natGatewayId
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListNatGatewaysApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回值
	var nats []CtyunNatsModel
	for _, natObj := range resp.ReturnObj {
		natItem := CtyunNatsModel{
			ID:           utils.SecStringValue(natObj.ID),
			Name:         utils.SecStringValue(natObj.Name),
			Description:  utils.SecStringValue(natObj.Description),
			Status:       types.Int32Value(natObj.Status), // 还需要测一下，如果types.Int32Value传参为空的话结果如何？
			NatGatewayID: utils.SecStringValue(natObj.NatGatewayID),
			ZoneID:       utils.SecStringValue(natObj.ZoneID),
			State:        utils.SecStringValue(natObj.State),
			VpcID:        utils.SecStringValue(natObj.VpcID),
			VpcName:      utils.SecStringValue(natObj.VpcName),
			ExpireTime:   utils.SecStringValue(natObj.ExpireTime),
			CreationTime: utils.SecStringValue(natObj.CreationTime),
			ProjectID:    utils.SecStringValue(natObj.ProjectID),
		}
		nats = append(nats, natItem)
	}
	config.RegionID = types.StringValue(regionId)
	config.Nats = nats
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunNats) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// ParseIntIfEmpty 自定义方法，用于判断types.Int64类型字段是否为空，若为空则返回默认值。
func (c *ctyunNats) ParseInt32IfEmpty(value types.Int32, defaultValue int32) int32 {
	if value.IsNull() {
		return defaultValue
	}
	return value.ValueInt32()
}

type CtyunNatsConfig struct {
	RegionID     types.String     `tfsdk:"region_id"`      //区域id
	NatGatewayID types.String     `tfsdk:"nat_gateway_id"` //要查询的NAT网关的ID。
	PageNumber   types.Int32      `tfsdk:"page_number"`    //	列表的页码，默认值为1。
	PageNo       types.Int32      `tfsdk:"page_no"`        //列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃
	PageSize     types.Int32      `tfsdk:"page_size"`      //分页查询时每页的行数，最大值为50，默认值为10。
	Nats         []CtyunNatsModel `tfsdk:"nats"`           // 获取的nat列表
}

type CtyunNatsModel struct {
	ID           types.String `tfsdk:"id"`             //nat网关id
	Name         types.String `tfsdk:"name"`           //nat网关名称
	Description  types.String `tfsdk:"description"`    //ctvpc 网关描述
	Status       types.Int32  `tfsdk:"status"`         //ctvpc 网关状态: 0 表示创建中，2 表示运行中，3 表示冻结
	NatGatewayID types.String `tfsdk:"nat_gateway_id"` //ctvpc 网关 id
	ZoneID       types.String `tfsdk:"zone_id"`        //可用区 ID
	State        types.String `tfsdk:"state"`          //NAT网关运行状态: running 表示运行中, creating 表示创建中, expired 表示已过期, freeze 表示已冻结
	VpcID        types.String `tfsdk:"vpc_id"`         //虚拟私有云 id
	VpcName      types.String `tfsdk:"vpc_name"`       //虚拟私有云名字
	ExpireTime   types.String `tfsdk:"expire_time"`    //过期时间
	CreationTime types.String `tfsdk:"creation_time"`  //创建时间
	ProjectID    types.String `tfsdk:"project_id"`     //项目 ID
}
