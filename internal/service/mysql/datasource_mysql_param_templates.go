package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunMysqlParamTemplates{}
	_ datasource.DataSourceWithConfigure = &CtyunMysqlParamTemplates{}
)

type CtyunMysqlParamTemplates struct {
	meta *common.CtyunMetadata
}

func NewCtyunMysqlParamTemplates() *CtyunMysqlParamTemplates {
	return &CtyunMysqlParamTemplates{}
}
func (c *CtyunMysqlParamTemplates) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMysqlParamTemplates) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_param_templates"
}

func (c *CtyunMysqlParamTemplates) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10098794",
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
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "参数模板名称",
			},
			"engine": schema.StringAttribute{
				Optional:    true,
				Description: "mysql版本",
				Validators: []validator.String{
					stringvalidator.OneOf("5.7", "8.0"),
				},
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "页码，默认1",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页记录数，默认10",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
				},
			},
			"param_templates": schema.ListNestedAttribute{
				Computed:    true,
				Description: "参数实例列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "参数模板名称",
						},
						"engine": schema.StringAttribute{
							Computed:    true,
							Description: "mysql实例版本",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"restart": schema.BoolAttribute{
							Computed:    true,
							Description: "修改该参数是否需要重启",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "参数模板描述",
						},
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "参数模板id",
						},
						"is_default": schema.Int32Attribute{
							Computed:    true,
							Description: "属否是默认参参数组",
						},
						"user_id": schema.Int64Attribute{
							Computed:    true,
							Description: "用户id",
						},
						"parameters": schema.ListNestedAttribute{
							Computed:    true,
							Description: "参数信息",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"parameter_name": schema.StringAttribute{
										Computed:    true,
										Description: "参数名",
									},
									"restart": schema.BoolAttribute{
										Computed:    true,
										Description: "修改配置是否需要重启",
									},
									"parameter_value": schema.StringAttribute{
										Computed:    true,
										Description: "参数值，支持更新",
									},
									"permit_value": schema.StringAttribute{
										Computed:    true,
										Description: "允许取值范围",
									},
									"description": schema.StringAttribute{
										Computed:    true,
										Description: "参数描述",
									},
									"parameter_id": schema.Int64Attribute{
										Computed:    true,
										Description: "参数id",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (c *CtyunMysqlParamTemplates) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunMysqlParamTemplatesConfig
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
	err = c.getParamTemplates(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlParamTemplates) getParamTemplates(ctx context.Context, config *CtyunMysqlParamTemplatesConfig) error {
	params := mysql.TeledbGetParameterTemplateListRequest{
		PageNow:  1,
		PageSize: 10,
	}
	if !config.Name.IsNull() {
		params.ParameterGroupName = config.Name.ValueStringPointer()
	}
	if !config.Engine.IsNull() {
		params.ParameterGroupName = config.Engine.ValueStringPointer()
	}
	header := mysql.TeledbGetParameterTemplateListRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	if !config.PageNo.IsNull() {
		params.PageNow = config.PageNo.ValueInt32()
	}
	if !config.PageSize.IsNull() {
		params.PageSize = config.PageSize.ValueInt32()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetParameterTemplateListApi.Do(ctx, c.meta.Credential, &params, &header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("查询参数模板(name=%s)列表失败，接口返回nil。请联系研发确认问题原因！", config.Name.ValueString())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}

	var parameterTemplates []MysqlParameterTemplateModel
	for _, item := range resp.ReturnObj.List {
		var parameterTemplate MysqlParameterTemplateModel
		parameterTemplate.Name = types.StringValue(item.ParameterGroupName)
		parameterTemplate.Engine = types.StringValue(item.MysqlEngine)
		parameterTemplate.ID = types.Int64Value(item.ID)
		parameterTemplate.Description = types.StringValue(item.Description)
		parameterTemplate.IsDefault = types.Int32Value(item.IsDefault)
		parameterTemplate.UserId = types.Int64Value(item.UserId)
		parameterTemplate.Restart = types.BoolValue(item.Restart)
		parameterTemplate.CreateTime = types.StringValue(utils.FromUnixToUTC(item.CreateTime))

		detail, err2 := c.getParameterTemplateDetail(ctx, config, item.ID)
		if err2 != nil {
			return err2
		}
		var parameters []CtyunMysqlParameterModel
		for _, parameterItem := range detail {
			var parameter CtyunMysqlParameterModel
			parameter.ParameterID = types.Int64Value(parameterItem.ID)
			parameter.ParameterName = types.StringValue(parameterItem.ParameterName)
			parameter.ParameterValue = types.StringValue(parameterItem.ParameterValue)
			parameter.Description = types.StringValue(parameterItem.Description)
			parameter.PermitValue = types.StringValue(parameterItem.PermitValue)
			if parameterItem.Restart == "1" {
				parameter.Restart = types.BoolValue(true)
			} else {
				parameter.Restart = types.BoolValue(false)
			}
			parameters = append(parameters, parameter)
		}
		parameterTemplate.Parameters = parameters
		parameterTemplates = append(parameterTemplates, parameterTemplate)
	}
	config.ParamTemplates = parameterTemplates
	return nil
}

func (c *CtyunMysqlParamTemplates) getParameterTemplateDetail(ctx context.Context, config *CtyunMysqlParamTemplatesConfig, templateID int64) ([]mysql.TeledbGetParameterTemplateDetailResponseReturnObjDeatil, error) {
	var pageNo, pageSize int32
	var parameters []mysql.TeledbGetParameterTemplateDetailResponseReturnObjDeatil
	pageNo = 1
	pageSize = 100
	resp, err := c.getParameterListByPage(ctx, config, pageNo, pageSize, templateID)
	if err != nil {
		return nil, err
	}
	pages := resp.ReturnObj.Pages
	for pageNo <= pages {
		parameters = append(parameters, resp.ReturnObj.List...)
		pageNo++
		if pageNo > pages {
			break
		}
		resp, err = c.getParameterListByPage(ctx, config, pageNo, pageSize, templateID)
		if err != nil {
			return nil, err
		}
	}
	return parameters, nil
}

func (c *CtyunMysqlParamTemplates) getParameterListByPage(ctx context.Context, config *CtyunMysqlParamTemplatesConfig, pageNo int32, pageSize int32, templateID int64) (*mysql.TeledbGetParameterTemplateDetailResponse, error) {
	params := &mysql.TeledbGetParameterTemplateDetailRequest{
		ID:       templateID,
		PageNow:  pageNo,
		PageSize: pageSize,
	}
	header := &mysql.TeledbGetParameterTemplateDetailRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetParameterTemplateDetailApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询参数模板详情(id=%d)失败，接口返回nil。请联系研发确认问题原因！", templateID)
		return nil, err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp, nil
}

type MysqlParameterTemplateModel struct {
	Name        types.String               `tfsdk:"name"`
	Engine      types.String               `tfsdk:"engine"`
	CreateTime  types.String               `tfsdk:"create_time"`
	Restart     types.Bool                 `tfsdk:"restart"`
	Description types.String               `tfsdk:"description"`
	ID          types.Int64                `tfsdk:"id"`
	IsDefault   types.Int32                `tfsdk:"is_default"`
	UserId      types.Int64                `tfsdk:"user_id"`
	Parameters  []CtyunMysqlParameterModel `tfsdk:"parameters"`
}

type CtyunMysqlParamTemplatesConfig struct {
	ProjectID      types.String                  `tfsdk:"project_id"`
	RegionID       types.String                  `tfsdk:"region_id"`
	Name           types.String                  `tfsdk:"name"`
	Engine         types.String                  `tfsdk:"engine"`
	PageNo         types.Int32                   `tfsdk:"page_no"`
	PageSize       types.Int32                   `tfsdk:"page_size"`
	ParamTemplates []MysqlParameterTemplateModel `tfsdk:"param_templates"`
}
