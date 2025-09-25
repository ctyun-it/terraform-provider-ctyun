package elb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ datasource.DataSource              = &ctyunElbHealthChecks{}
	_ datasource.DataSourceWithConfigure = &ctyunElbHealthChecks{}
)

type ctyunElbHealthChecks struct {
	meta *common.CtyunMetadata
}

func NewCtyunElbHealthChecks() datasource.DataSource {
	return &ctyunElbHealthChecks{}
}

func (c *ctyunElbHealthChecks) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunElbHealthChecks) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_health_checks"
}

func (c *ctyunElbHealthChecks) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026756/10032101`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"ids": schema.StringAttribute{
				Optional:    true,
				Description: "健康检查ID列表，用逗号分隔",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "健康检查名称, 只能由数字，字母，-组成不能以数字和-开头，最大长度32\t",
			},
			"health_checks": schema.ListNestedAttribute{
				Computed:    true,
				Description: "健康检查返回列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源池ID",
						},
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "可用区名称",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "项目ID",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "健康检查ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "健康检查名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"protocol": schema.StringAttribute{
							Computed:    true,
							Description: "健康检查协议: TCP / UDP / HTTP",
							Validators: []validator.String{
								stringvalidator.OneOf(business.HealthCheckProtocols...),
							},
						},
						"protocol_port": schema.Int32Attribute{
							Computed:    true,
							Description: "健康检查端口",
							Validators: []validator.Int32{
								int32validator.Between(1, 65535),
							},
						},
						"timeout": schema.Int32Attribute{
							Computed:    true,
							Description: "健康检查响应的最大超时时间",
							Validators: []validator.Int32{
								int32validator.Between(2, 60),
							},
						},
						"integererval": schema.Int32Attribute{
							Computed:    true,
							Description: "负载均衡进行健康检查的时间间隔",
							Validators: []validator.Int32{
								int32validator.Between(1, 20940),
							},
						},
						"max_retry": schema.Int32Attribute{
							Computed:    true,
							Description: "最大重试次数",
							Validators: []validator.Int32{
								int32validator.Between(1, 10),
							},
						},
						"http_method": schema.StringAttribute{
							Computed:    true,
							Description: "HTTP请求的方法",
							Validators: []validator.String{
								stringvalidator.OneOf(business.HttpMethods...),
							},
						},
						"http_url_path": schema.StringAttribute{
							Computed:    true,
							Description: "HTTP请求url路径",
							Validators: []validator.String{
								stringvalidator.LengthAtMost(80),
							},
						},
						"http_expected_codes": schema.StringAttribute{
							Computed:    true,
							Description: "HTTP预期码",
						},
						"status": schema.Int32Attribute{
							Computed:    true,
							Description: "状态 1 表示 UP, 0 表示 DOWN",
							Validators: []validator.Int32{
								int32validator.Between(0, 1),
							},
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

func (c *ctyunElbHealthChecks) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunElbHealthChecks
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	params := &ctelb.CtelbListHealthCheckRequest{
		ClientToken: uuid.NewString(),
		RegionID:    regionId,
	}
	var ids []string
	if !config.IDs.IsNull() {
		ids = strings.Split(config.IDs.ValueString(), ",")
		params.IDs = ids
	}
	if !config.Name.IsNull() {
		params.Name = config.Name.ValueString()
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbListHealthCheckApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回值
	var healthChecks []HealthCheckModel
	for _, healthCheckItem := range resp.ReturnObj {
		var healthCheck HealthCheckModel
		healthCheck.RegionID = types.StringValue(healthCheckItem.RegionID)
		healthCheck.AzName = types.StringValue(healthCheckItem.Name)
		healthCheck.ID = types.StringValue(healthCheckItem.ID)
		healthCheck.ProjectID = types.StringValue(healthCheckItem.ProjectID)
		healthCheck.Name = types.StringValue(healthCheckItem.Name)
		healthCheck.Description = types.StringValue(healthCheckItem.Description)
		healthCheck.Protocol = types.StringValue(healthCheckItem.Protocol)
		healthCheck.ProtocolPort = types.Int32Value(healthCheckItem.ProtocolPort)
		healthCheck.Timeout = types.Int32Value(healthCheckItem.Timeout)
		healthCheck.Integererval = types.Int32Value(healthCheckItem.Integererval)
		healthCheck.MaxRetry = types.Int32Value(healthCheckItem.MaxRetry)
		healthCheck.HttpMethod = types.StringValue(healthCheckItem.HttpMethod)
		healthCheck.HttpUrlPath = types.StringValue(healthCheckItem.HttpUrlPath)
		healthCheck.HttpExpectedCodes = types.StringValue(healthCheckItem.HttpExpectedCodes)
		healthCheck.Status = types.Int32Value(healthCheckItem.Status)
		healthCheck.CreateTime = types.StringValue(healthCheckItem.CreateTime)
		healthChecks = append(healthChecks, healthCheck)
	}
	config.HealthChecks = healthChecks
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunElbHealthChecks struct {
	RegionID     types.String       `tfsdk:"region_id"` //区域ID
	IDs          types.String       `tfsdk:"ids"`       //健康检查ID列表
	Name         types.String       `tfsdk:"name"`      //健康检查名称, 只能由数字，字母，-组成不能以数字和-开头，最大长度32
	HealthChecks []HealthCheckModel `tfsdk:"health_checks"`
}

type HealthCheckModel struct {
	RegionID          types.String `tfsdk:"region_id"`           //区域ID
	AzName            types.String `tfsdk:"az_name"`             //可用区名称
	ProjectID         types.String `tfsdk:"project_id"`          //项目ID
	ID                types.String `tfsdk:"id"`                  //健康检查ID
	Name              types.String `tfsdk:"name"`                //健康检查名称
	Description       types.String `tfsdk:"description"`         //描述
	Protocol          types.String `tfsdk:"protocol"`            //健康检查协议: TCP / UDP / HTTP
	ProtocolPort      types.Int32  `tfsdk:"protocol_port"`       //健康检查端口
	Timeout           types.Int32  `tfsdk:"timeout"`             //健康检查响应的最大超时时间
	Integererval      types.Int32  `tfsdk:"integererval"`        //负载均衡进行健康检查的时间间隔
	MaxRetry          types.Int32  `tfsdk:"max_retry"`           //最大重试次数
	HttpMethod        types.String `tfsdk:"http_method"`         //HTTP请求的方法
	HttpUrlPath       types.String `tfsdk:"http_url_path"`       //HTTP请求url路径
	HttpExpectedCodes types.String `tfsdk:"http_expected_codes"` //HTTP预期码
	Status            types.Int32  `tfsdk:"status"`              //状态 1 表示 UP, 0 表示 DOWN
	CreateTime        types.String `tfsdk:"create_time"`         //创建时间，为UTC格式
}
