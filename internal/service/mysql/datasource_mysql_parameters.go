package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var (
	_ datasource.DataSource              = &CtyunMysqlParameters{}
	_ datasource.DataSourceWithConfigure = &CtyunMysqlParameters{}
)

type CtyunMysqlParameters struct {
	meta *common.CtyunMetadata
}

func NewCtyunMysqlParameters() *CtyunMysqlParameters {
	return &CtyunMysqlParameters{}
}
func (c *CtyunMysqlParameters) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMysqlParameters) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_parameters"
}

func (c *CtyunMysqlParameters) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10035295",
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
				Description: "mysql实例id",
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
			"parameters": schema.ListNestedAttribute{
				Computed:    true,
				Description: "实例参数列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"parameter_name": schema.StringAttribute{
							Computed:    true,
							Description: "参数名",
						},
						"value_type": schema.StringAttribute{
							Computed:    true,
							Description: "参数类型",
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
	}
}

func (c *CtyunMysqlParameters) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunMysqlParametersConfig
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
	parameterList, err := c.getParameters(ctx, &config)
	if err != nil {
		return
	}

	var parameters []MysqlParametersModel
	for _, parameterItem := range parameterList {
		var parameter MysqlParametersModel
		parameter.ParameterName = types.StringValue(parameterItem.ParameterName)
		parameter.ParameterValue = types.StringValue(parameterItem.ParameterValue)
		parameter.ValueType = types.StringValue(parameterItem.ValueType)
		parameter.Restart = types.BoolValue(parameterItem.Restart == "1")
		parameter.Description = types.StringValue(parameterItem.Description)
		parameter.PermitValue = types.StringValue(parameterItem.PermitValue)
		parameter.ParameterID = types.Int64Value(parameterItem.ID)
		parameters = append(parameters, parameter)
	}
	config.Parameters = parameters
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlParameters) getUnixTime(timestamp int64) string {
	isoTime := time.Unix(timestamp, 0).UTC().Format("2006-01-02T15:04:05Z")
	return isoTime
}

func (c *CtyunMysqlParameters) getParameters(ctx context.Context, config *CtyunMysqlParametersConfig) ([]mysql.TeledbGetRdsParameterTemplateDetailResponseReturnObjDetail, error) {
	params := &mysql.TeledbGetRdsParameterTemplateDetailRequest{
		OuterProdInstId: config.InstID.ValueString(),
		PageNow:         1,
		PageSize:        10,
	}
	header := &mysql.TeledbGetRdsParameterTemplateDetailRequestHeader{
		RegionID: config.RegionID.ValueString(),
		InstID:   config.InstID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() && config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	if !config.PageNo.IsNull() {
		params.PageNow = config.PageNo.ValueInt32()
	}
	if !config.PageSize.IsNull() {
		params.PageSize = config.PageSize.ValueInt32()
	}

	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetRdsParameterTemplateDetailApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询参数模板列表失败，接口返回nil。请联系研发确认问题原因！")
		return nil, err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp.ReturnObj.List, nil
}

type MysqlParametersModel struct {
	ParameterName  types.String `tfsdk:"parameter_name"`
	ValueType      types.String `tfsdk:"value_type"`
	Restart        types.Bool   `tfsdk:"restart"`
	ParameterValue types.String `tfsdk:"parameter_value"`
	PermitValue    types.String `tfsdk:"permit_value"`
	Description    types.String `tfsdk:"description"`
	ParameterID    types.Int64  `tfsdk:"parameter_id"`
}

type CtyunMysqlParametersConfig struct {
	ProjectID  types.String           `tfsdk:"project_id"`
	RegionID   types.String           `tfsdk:"region_id"`
	InstID     types.String           `tfsdk:"instance_id"`
	PageNo     types.Int32            `tfsdk:"page_no"`
	PageSize   types.Int32            `tfsdk:"page_size"`
	Parameters []MysqlParametersModel `tfsdk:"parameters"`
}
