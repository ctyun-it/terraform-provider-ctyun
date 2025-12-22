package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunMysqlCharacterSet{}
	_ datasource.DataSourceWithConfigure = &ctyunMysqlCharacterSet{}
)

type ctyunMysqlCharacterSet struct {
	meta *common.CtyunMetadata
}

func NewCtyunMysqlCharacterSet() datasource.DataSource {
	return &ctyunMysqlCharacterSet{}
}
func (c *ctyunMysqlCharacterSet) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunMysqlCharacterSet) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_character_set"
}

func (c *ctyunMysqlCharacterSet) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10140487",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "MySQL实例ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "项目ID",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，默认使用provider配置",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"character_set": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"charset": schema.StringAttribute{
							Computed:    true,
							Description: "字符集名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "字符集描述",
						},
						"max_len": schema.Int64Attribute{
							Computed:    true,
							Description: "最大字节长度",
						},
						"default_collation": schema.StringAttribute{
							Computed:    true,
							Description: "默认校对规则",
						},
					},
				},
				Description: "MySQL支持的字符集列表",
			},
		},
	}
}

func (c *ctyunMysqlCharacterSet) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunMysqlCharacterSetConfig
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

	params := &mysql.TeledbGetCharacterSetRequest{
		OuterProdInstId: config.InstID.ValueString(),
	}
	header := &mysql.TeledbGetCharacterSetRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetCharacterSetApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp == nil {
		err = errors.New("查询mysql支持字符集失败")
		return
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var MysqlCharacterSet []CtyunMysqlDatabaseCharacterModel
	for _, characterItem := range resp.ReturnObj.Data {
		var character CtyunMysqlDatabaseCharacterModel
		character.Charset = types.StringValue(characterItem.Charset)
		character.DefaultCollation = types.StringValue(characterItem.DefaultCollation)
		character.MaxLen = types.Int32Value(characterItem.MaxLen)
		character.Description = types.StringValue(characterItem.Description)
		MysqlCharacterSet = append(MysqlCharacterSet, character)
	}
	config.MysqlCharacterSet = MysqlCharacterSet
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunMysqlDatabaseCharacterModel struct {
	Charset          types.String `tfsdk:"charset"`
	Description      types.String `tfsdk:"description"`
	MaxLen           types.Int32  `tfsdk:"max_len"`
	DefaultCollation types.String `tfsdk:"default_collation"`
}

type CtyunMysqlCharacterSetConfig struct {
	InstID            types.String                       `tfsdk:"instance_id"`
	ProjectID         types.String                       `tfsdk:"project_id"`
	RegionID          types.String                       `tfsdk:"region_id"`
	MysqlCharacterSet []CtyunMysqlDatabaseCharacterModel `tfsdk:"character_set"`
}
