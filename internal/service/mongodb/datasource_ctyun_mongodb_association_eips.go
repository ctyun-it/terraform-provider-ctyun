package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunMongodbAssociationEips{}
	_ datasource.DataSourceWithConfigure = &CtyunMongodbAssociationEips{}
)

type CtyunMongodbAssociationEips struct {
	meta *common.CtyunMetadata
}

func NewCtyunMongodbAssociationEips() *CtyunMongodbAssociationEips {
	return &CtyunMongodbAssociationEips{}
}
func (c *CtyunMongodbAssociationEips) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMongodbAssociationEips) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mongodb_association_eips"
}

func (c *CtyunMongodbAssociationEips) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10034467/10183412`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id",
			},
			"inst_id": schema.StringAttribute{
				Optional:    true,
				Description: "实例id",
			},
			"eip_id": schema.StringAttribute{
				Optional:    true,
				Description: "弹性ip唯一标识",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "项目id",
			},
			"eips": schema.ListNestedAttribute{
				Computed:    true,
				Description: "绑定eip列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"eip_id": schema.StringAttribute{
							Computed:    true,
							Description: "弹性ip唯一标识",
						},
						"eip": schema.StringAttribute{
							Computed:    true,
							Description: "弹性IP",
						},
						"bind_status": schema.Int32Attribute{
							Computed:    true,
							Description: "0-未绑定，1-已绑定",
							Validators: []validator.Int32{
								int32validator.Between(0, 1),
							},
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "状态标识：ACTIVE=已使用，DOWN=未使用，ERROR=中间状态-异常，UPDATING=中间状态-更新中，BANDING_OR_UNBANGDING=中间状态-绑定或解绑中，DELETING=中间状态-删除中，DELETED=中间状态-已删除",
							Validators: []validator.String{
								stringvalidator.OneOf(business.MysqlBindEipStatus...),
							},
						},
						"band_width": schema.Int32Attribute{
							Computed:    true,
							Description: "加入的共享带宽，单位M",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunMongodbAssociationEips) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunMongodbAssociationEipsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	// 查询绑定eip 列表
	params := &mongodb.MongodbBoundEipListRequest{
		RegionID: regionId,
	}
	if config.EipID.ValueString() != "" {
		params.EipID = config.EipID.ValueStringPointer()
	}
	if config.InstID.ValueString() != "" {
		params.InstID = config.InstID.ValueStringPointer()
	}
	headers := &mongodb.MongodbBoundEipListRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		headers.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkMongodbApis.MongodbBoundEipListApi.Do(ctx, c.meta.Credential, params, headers)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s ", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回的绑定eip列表
	returnObj := resp.ReturnObj.Data
	var eips []BoundEipModel
	for _, eipItem := range returnObj {
		var eip BoundEipModel
		eip.EipID = types.StringValue(eipItem.EipID)
		eip.Eip = types.StringValue(eipItem.Eip)
		eip.BindStatus = types.Int32Value(eipItem.BindStatus)
		eip.Status = types.StringValue(eipItem.Status)
		eip.BandWidth = types.Int32Value(eipItem.BandWidth)
		eips = append(eips, eip)
	}
	config.Eips = eips
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunMongodbAssociationEipsConfig struct {
	RegionID  types.String    `tfsdk:"region_id"`
	InstID    types.String    `tfsdk:"inst_id"`
	EipID     types.String    `tfsdk:"eip_id"`
	ProjectID types.String    `tfsdk:"project_id"`
	Eips      []BoundEipModel `tfsdk:"eips"`
}

type BoundEipModel struct {
	EipID      types.String `tfsdk:"eip_id"`
	Eip        types.String `tfsdk:"eip"`
	BindStatus types.Int32  `tfsdk:"bind_status"`
	Status     types.String `tfsdk:"status"`
	BandWidth  types.Int32  `tfsdk:"band_width"`
}
