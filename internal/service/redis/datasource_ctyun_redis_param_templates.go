package redis

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgdcs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunRedisParamTemplates{}
	_ datasource.DataSourceWithConfigure = &ctyunRedisParamTemplates{}
)

type ctyunRedisParamTemplates struct {
	meta *common.CtyunMetadata
}

func NewCtyunRedisParamTemplates() datasource.DataSource {
	return &ctyunRedisParamTemplates{}
}

func (c *ctyunRedisParamTemplates) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_param_templates"
}

type CtyunRedisParamTemplateModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CacheMode   types.String `tfsdk:"cache_mode"`
	SysTemplate types.Bool   `tfsdk:"sys_template"`
}

type CtyunRedisParamTemplatesConfig struct {
	RegionId types.String                   `tfsdk:"region_id"`
	Type     types.String                   `tfsdk:"type"`
	PageNum  types.Int32                    `tfsdk:"page_num"`
	PageSize types.Int32                    `tfsdk:"page_size"`
	Total    types.Int32                    `tfsdk:"total"`
	Size     types.Int32                    `tfsdk:"size"`
	List     []CtyunRedisParamTemplateModel `tfsdk:"list"`
}

func (c *ctyunRedisParamTemplates) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029420/10156164`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "模板类型 sys：系统模板 custom：自定义模板",
				Validators: []validator.String{
					stringvalidator.OneOf("sys", "custom"),
				},
			},
			"page_num": schema.Int32Attribute{
				Optional:    true,
				Description: "页码，默认1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页记录数，默认值：10",
			},
			"total": schema.Int32Attribute{
				Computed:    true,
				Description: "总数",
			},
			"size": schema.Int32Attribute{
				Computed:    true,
				Description: "本次返回数量",
			},
			"list": schema.ListNestedAttribute{
				Computed:    true,
				Description: "参数模板列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "参数记录ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "参数名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "参数描述",
						},
						"cache_mode": schema.StringAttribute{
							Computed:    true,
							Description: "适合的实例架构版本 ORIGINAL_67：Redis 6.0/7.0类型 ORIGINAL_5：Redis 5.0类型",
						},
						"sys_template": schema.BoolAttribute{
							Computed:    true,
							Description: "是否为系统模板 true：系统模板 false：自定义模板",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRedisParamTemplates) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRedisParamTemplatesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := c.meta.GetExtraIfEmpty(config.RegionId.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	config.RegionId = types.StringValue(regionId)

	// 组装请求体
	params := &ctgdcs2.Dcs2DescribeRedisTemplateRequest{
		RegionId: regionId,
		RawType:  config.Type.ValueString(),
		PageNum:  config.PageNum.ValueInt32(),
		PageSize: config.PageSize.ValueInt32(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeRedisTemplateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 设置分页信息
	config.Total = types.Int32Value(resp.ReturnObj.Total)
	config.Size = types.Int32Value(resp.ReturnObj.Size)

	// 解析返回值
	config.List = []CtyunRedisParamTemplateModel{}
	for _, template := range resp.ReturnObj.List {
		item := CtyunRedisParamTemplateModel{
			ID:          types.StringValue(template.Id),
			Name:        types.StringValue(template.Name),
			Description: types.StringValue(template.Description),
			CacheMode:   types.StringValue(template.CacheMode),
		}
		if template.SysTemplate != nil {
			item.SysTemplate = types.BoolValue(*template.SysTemplate)
		} else {
			item.SysTemplate = types.BoolNull()
		}
		config.List = append(config.List, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRedisParamTemplates) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
