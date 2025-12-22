package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunPgsqlParamTemplates{}
	_ datasource.DataSourceWithConfigure = &CtyunPgsqlParamTemplates{}
)

type CtyunPgsqlParamTemplates struct {
	meta *common.CtyunMetadata
}

func NewCtyunPgsqlParamTemplates() *CtyunPgsqlParamTemplates {
	return &CtyunPgsqlParamTemplates{}
}
func (c *CtyunPgsqlParamTemplates) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPgsqlParamTemplates) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_param_templates"
}

func (c *CtyunPgsqlParamTemplates) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10034019/10166169",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "区域ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "项目ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "分页页码，默认为1",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页记录数，默认为10",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "模板名称过滤条件",
			},
			"parameter_templates": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "参数模板ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "模板名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "模板描述",
						},
						"version": schema.StringAttribute{
							Computed:    true,
							Description: "数据库版本",
						},
						"modify": schema.BoolAttribute{
							Computed:    true,
							Description: "是否允许修改",
						},
						"update_timestamp": schema.StringAttribute{
							Computed:    true,
							Description: "最后更新时间戳",
						},
						"parameters": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"description": schema.StringAttribute{
										Computed:    true,
										Description: "参数描述",
									},
									"unit": schema.StringAttribute{
										Computed:    true,
										Description: "参数单位",
									},
									"min_val": schema.StringAttribute{
										Computed:    true,
										Description: "最小值",
									},
									"max_val": schema.StringAttribute{
										Computed:    true,
										Description: "最大值",
									},
									"enum_values": schema.SetAttribute{
										Computed:    true,
										ElementType: types.StringType,
										Description: "枚举值列表",
									},
									"restart": schema.Int32Attribute{
										Computed:    true,
										Description: "是否需要重启生效（0-否 1-是）",
									},
									"param_name": schema.StringAttribute{
										Computed:    true,
										Description: "参数名称",
									},
									"param_value": schema.StringAttribute{
										Computed:    true,
										Description: "参数值",
									},
									"value_type": schema.StringAttribute{
										Computed:    true,
										Description: "值类型（int/string/bool等）",
									},
								},
							},
							Description: "参数列表",
						},
					},
				},
				Description: "参数模板列表",
			},
		},
	}
}

func (c *CtyunPgsqlParamTemplates) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunPostgresqlParameterTemplatesConfig
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
	templateList, err := c.getPgsqlParameterTemplateList(ctx, config)
	var parameterTemplates []ParameterTemplateInfo
	for _, templateItem := range templateList {
		var template ParameterTemplateInfo
		template.ID = types.Int64Value(templateItem.PgTemplateId)
		template.Name = types.StringValue(templateItem.Name)
		template.Description = types.StringValue(templateItem.Description)
		template.Modify = types.BoolValue(templateItem.Modify)
		template.Version = types.StringValue(templateItem.Version)
		template.UpdateTimestamp = types.StringValue(templateItem.UpdateTimestamp)
		parametersList, err2 := c.getTemplateDetail(ctx, config, templateItem.PgTemplateId)
		if err2 != nil {
			return
		}
		var templateParameters []ParameterInfo
		for _, parameterItem := range parametersList {
			var parameterInfo ParameterInfo
			parameterInfo.ParamName = types.StringValue(parameterItem.ParameterName)
			parameterInfo.ParamValue = types.StringValue(parameterItem.ParameterValue)
			parameterInfo.MaxVal = types.StringValue(parameterItem.MaxVal)
			parameterInfo.MinVal = types.StringValue(parameterItem.MinVal)
			parameterInfo.Unit = types.StringValue(parameterItem.Unit)
			parameterInfo.Description = types.StringValue(parameterItem.Description)
			parameterInfo.Restart = types.Int32Value(parameterItem.Restart)
			parameterInfo.ValueType = types.StringValue(parameterItem.ValueType)
			parameterInfo.Unit = types.StringValue(parameterItem.Unit)
			enumValues, diags := types.SetValueFrom(ctx, types.StringType, parameterItem.EnumValues)
			if diags != nil {
				err = fmt.Errorf(diags[0].Detail())
				return
			}
			parameterInfo.EnumValues = enumValues
			templateParameters = append(templateParameters, parameterInfo)
		}

		parameterTemplates = append(parameterTemplates, template)
	}
	config.ParameterTemplates = parameterTemplates
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPgsqlParamTemplates) getPgsqlParameterTemplateList(ctx context.Context, config CtyunPostgresqlParameterTemplatesConfig) ([]pgsql.ParameterTemplateInfo, error) {
	params := &pgsql.PgsqlGetParameterTemplateListRequest{
		PageNow:  1,
		PageSize: 10,
	}
	header := &pgsql.PgsqlGetParameterTemplateListRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.PageSize.IsNull() {
		params.PageSize = config.PageSize.ValueInt32()
	}
	if !config.PageNo.IsNull() {
		params.PageNow = config.PageNo.ValueInt32()
	}
	if !config.Name.IsNull() {
		params.Name = config.Name.ValueStringPointer()
	}

	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGetParameterTemplateListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询name=%s的参数模板失败，接口返回nil，请联系研发确认问题原因！", config.Name.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return nil, err
	}

	return resp.ReturnObj.List, nil
}

func (c *CtyunPgsqlParamTemplates) getTemplateDetail(ctx context.Context, config CtyunPostgresqlParameterTemplatesConfig, templateID int64) ([]pgsql.PgsqlGetParameterTemplateDetailResponseReturnObj, error) {

	params := &pgsql.PgsqlGetParameterTemplateDetailRequest{
		TemplateId: templateID,
	}
	header := &pgsql.PgsqlGetParameterTemplateDetailRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGetParameterTemplateDetailApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询模板id=%d的参数信息失败，接口返回nil，请联系研发确认问题原因！", templateID)
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return nil, err
	}
	return resp.ReturnObj, nil

}

type ParameterInfo struct {
	Description types.String `tfsdk:"description"`
	Unit        types.String `tfsdk:"unit"`
	MinVal      types.String `tfsdk:"min_val"`
	MaxVal      types.String `tfsdk:"max_val"`
	EnumValues  types.Set    `tfsdk:"enum_values"`
	Restart     types.Int32  `tfsdk:"restart"`
	ParamName   types.String `tfsdk:"param_name"`
	ParamValue  types.String `tfsdk:"param_value"`
	ValueType   types.String `tfsdk:"value_type"`
}

type ParameterTemplateInfo struct {
	ID              types.Int64     `tfsdk:"id"`
	Name            types.String    `tfsdk:"name"`
	Description     types.String    `tfsdk:"description"`
	Version         types.String    `tfsdk:"version"`
	Modify          types.Bool      `tfsdk:"modify"`
	UpdateTimestamp types.String    `tfsdk:"update_timestamp"`
	Parameters      []ParameterInfo `tfsdk:"parameters"`
}

type CtyunPostgresqlParameterTemplatesConfig struct {
	RegionID           types.String            `tfsdk:"region_id"`
	ProjectID          types.String            `tfsdk:"project_id"`
	PageNo             types.Int32             `tfsdk:"page_no"`
	PageSize           types.Int32             `tfsdk:"page_size"`
	Name               types.String            `tfsdk:"name"`
	ParameterTemplates []ParameterTemplateInfo `tfsdk:"parameter_templates"`
}
