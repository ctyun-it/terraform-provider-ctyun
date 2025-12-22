package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunPostgresqlCollationTimeZone{}
	_ datasource.DataSourceWithConfigure = &ctyunPostgresqlCollationTimeZone{}
)

type ctyunPostgresqlCollationTimeZone struct {
	meta *common.CtyunMetadata
}

func NewCtyunPostgresqlCollationTimeZone() datasource.DataSource {
	return &ctyunPostgresqlCollationTimeZone{}
}
func (c *ctyunPostgresqlCollationTimeZone) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunPostgresqlCollationTimeZone) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_collation_time_zone"
}

func (c *ctyunPostgresqlCollationTimeZone) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10034019/10159978",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，默认使用provider配置",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "项目ID",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "PostgreSQL实例ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"standard_time_offset": schema.StringAttribute{
				Computed:    true,
				Description: "世界协调时间偏移。由世界协调时间UTC+时区差组成，格式：(UTC+HH:mm)",
			},
			"time_zone": schema.StringAttribute{
				Computed:    true,
				Description: "时区",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "时区描述",
			},
			"collations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"coll_name": schema.StringAttribute{
							Computed:    true,
							Description: "排序规则名字（在每一个名字空间和编码中唯一）",
						},
						"collen_coding": schema.StringAttribute{
							Computed:    true,
							Description: "该排序规则可应用的编码，为空表示它可用于任何编码",
						},
						"coll_collate": schema.StringAttribute{
							Computed:    true,
							Description: "该排序规则对象的LC_COLLATE",
						},
						"coll_type": schema.StringAttribute{
							Computed:    true,
							Description: "该排序规则对象的LC_CTYPE",
						},
					},
				},
				Description: "排序规则列表",
			},
		},
	}
}

func (c *ctyunPostgresqlCollationTimeZone) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunPostgresqlCollationTimeZoneConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	config.RegionID = types.StringValue(regionId)

	params := &pgsql.PgsqlGetCollationTimeZoneRequest{
		ProdInstId: config.InstID.ValueString(),
	}
	header := &pgsql.PgsqlGetCollationTimeZoneRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGetCollationTimeZoneApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询postgresql排序规则和时区失败，接口返回nil，请联系研发确认问题原因！")
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	config.StandardTimeOffset = types.StringValue(resp.ReturnObj.StandardTimeOffset)
	config.TimeZone = types.StringValue(resp.ReturnObj.TimeZone)
	config.Description = types.StringValue(resp.ReturnObj.Description)
	var collations []CollationModel
	for _, collationItem := range resp.ReturnObj.Collations {
		var collation CollationModel
		collation.CollName = types.StringValue(collationItem.CollName)
		collation.CollenCoding = types.StringValue(collationItem.CollenCoding)
		collation.CollCollate = types.StringValue(collationItem.CollCollate)
		collation.CollType = types.StringValue(collationItem.CollCtype)
		collations = append(collations, collation)
	}
	config.Collations = collations
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CollationModel struct {
	CollName     types.String `tfsdk:"coll_name"`
	CollenCoding types.String `tfsdk:"collen_coding"`
	CollCollate  types.String `tfsdk:"coll_collate"`
	CollType     types.String `tfsdk:"coll_type"`
}

type CtyunPostgresqlCollationTimeZoneConfig struct {
	RegionID           types.String     `tfsdk:"region_id"`
	ProjectID          types.String     `tfsdk:"project_id"`
	InstID             types.String     `tfsdk:"instance_id"`
	StandardTimeOffset types.String     `tfsdk:"standard_time_offset"`
	TimeZone           types.String     `tfsdk:"time_zone"`
	Collations         []CollationModel `tfsdk:"collations"`
	Description        types.String     `tfsdk:"description"`
}
