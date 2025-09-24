package vpce

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ datasource.DataSource              = &ctyunVpces{}
	_ datasource.DataSourceWithConfigure = &ctyunVpces{}
)

type ctyunVpces struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpces() datasource.DataSource {
	return &ctyunVpces{}
}

func (c *ctyunVpces) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpces"
}

type CtyunVpcesModel struct {
	ID                types.String `tfsdk:"id"`
	EndpointServiceID types.String `tfsdk:"endpoint_service_id"`
	Type              types.String `tfsdk:"type"`
	Name              types.String `tfsdk:"name"`
	VpcID             types.String `tfsdk:"vpc_id"`
	SubnetID          types.String `tfsdk:"subnet_id"`
	SubnetIP          types.String `tfsdk:"subnet_ip"`
	WhitelistCidr     types.Set    `tfsdk:"whitelist_cidr"`
	Status            types.Int32  `tfsdk:"status"`
	CreatedTime       types.String `tfsdk:"created_time"`
	UpdatedTime       types.String `tfsdk:"updated_time"`
}

type CtyunVpcesConfig struct {
	RegionID   types.String `tfsdk:"region_id"`
	PageNo     types.Int32  `tfsdk:"page_no"`
	PageSize   types.Int32  `tfsdk:"page_size"`
	EndpointID types.String `tfsdk:"endpoint_id"`

	CurrentCount types.Int32       `tfsdk:"current_count"`
	TotalCount   types.Int32       `tfsdk:"total_count"`
	TotalPage    types.Int32       `tfsdk:"total_page"`
	Vpces        []CtyunVpcesModel `tfsdk:"vpces"`
}

func (c *ctyunVpces) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10042658/10217121**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"endpoint_id": schema.StringAttribute{
				Optional:    true,
				Description: "终端节点id",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码，默认值为1,推荐使用该字段",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页数据量大小，取值1-50",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},
			"current_count": schema.Int32Attribute{
				Computed:    true,
				Description: "分页查询时每页的行数。",
			},
			"total_count": schema.Int32Attribute{
				Computed:    true,
				Description: "总数。",
			},
			"total_page": schema.Int32Attribute{
				Computed:    true,
				Description: "总页数。",
			},
			"vpces": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:    true,
							Description: "终端节点ID",
						},
						"endpoint_service_id": schema.StringAttribute{
							Required:    true,
							Description: "终端节点服务ID",
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "接口还是反向，interface:接口，reverse:反向",
						},
						"name": schema.StringAttribute{
							Required:    true,
							Description: "终端节点名称",
						},
						"vpc_id": schema.StringAttribute{
							Required:    true,
							Description: "所属的专有网络id",
						},
						"subnet_id": schema.StringAttribute{
							Required:    true,
							Description: "子网ID",
						},
						"subnet_ip": schema.StringAttribute{
							Required:    true,
							Description: "子网IP",
						},
						"whitelist_cidr": schema.SetAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "白名单",
						},
						"status": schema.Int32Attribute{
							Required:    true,
							Description: "endpoint状态,1表示已链接，2表示未链接",
						},
						"created_time": schema.StringAttribute{
							Required:    true,
							Description: "创建时间",
						},
						"updated_time": schema.StringAttribute{
							Required:    true,
							Description: "更新时间",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunVpces) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunVpcesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	config.RegionID = types.StringValue(regionId)
	// 组装请求体
	params := &ctvpc.CtvpcNewEndpointsListRequest{
		RegionID: regionId,
	}
	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	id := config.EndpointID.ValueString()
	if pageNo > 0 {
		params.PageNo = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	if id != "" {
		params.EndpointID = &id
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcNewEndpointsListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回值
	config.Vpces = []CtyunVpcesModel{}
	config.TotalPage = types.Int32Value(resp.ReturnObj.TotalPage)
	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
	config.CurrentCount = types.Int32Value(resp.ReturnObj.CurrentCount)
	for _, e := range resp.ReturnObj.Endpoints {
		item := CtyunVpcesModel{
			ID:                utils.SecStringValue(e.ID),
			EndpointServiceID: utils.SecStringValue(e.EndpointServiceID),
			Type:              utils.SecStringValue(e.RawType),
			VpcID:             utils.SecStringValue(e.VpcID),
			Name:              utils.SecStringValue(e.Name),
			Status:            types.Int32Value(e.Status),
			CreatedTime:       utils.SecStringValue(e.CreatedTime),
			UpdatedTime:       utils.SecStringValue(e.UpdatedTime),
		}
		if e.EndpointObj != nil {
			item.SubnetID = utils.SecStringValue(e.EndpointObj.SubnetID)
			item.SubnetIP = utils.SecStringValue(e.EndpointObj.Ip)
		}
		whiteList := utils.SecString(e.Whitelist)
		if len(whiteList) > 0 {
			t := strings.Split(whiteList, ",")
			item.WhitelistCidr, _ = types.SetValueFrom(ctx, types.StringType, t)
		} else {
			item.WhitelistCidr = types.SetNull(types.StringType)
		}
		config.Vpces = append(config.Vpces, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunVpces) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
