package elb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &CtyunElbHealthCheck{}
	_ resource.ResourceWithConfigure   = &CtyunElbHealthCheck{}
	_ resource.ResourceWithImportState = &CtyunElbHealthCheck{}
)

type CtyunElbHealthCheck struct {
	meta *common.CtyunMetadata
}

func NewCtyunElbHealthCheck() resource.Resource {
	return &CtyunElbHealthCheck{}
}

func (c *CtyunElbHealthCheck) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunElbHealthCheck) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_health_check"
}

func (c *CtyunElbHealthCheck) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026756/10032101`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 32),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·！@#￥%……&*（） —— -+={}\\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 128),
					validator2.Desc(),
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "健康检查协议。取值范围：TCP、UDP、HTTP",
				Validators: []validator.String{
					stringvalidator.OneOf(business.HealthCheckProtocols...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"timeout": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "健康检查响应的最大超时时间，取值范围：2-60秒，默认为2秒，支持更新",
				Default:     int32default.StaticInt32(2),
				Validators: []validator.Int32{
					int32validator.Between(2, 60),
				},
			},
			"interval": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "负载均衡进行健康检查的时间间隔，取值范围：1-20940秒，默认为5秒，支持更新",
				Default:     int32default.StaticInt32(5),
				Validators: []validator.Int32{
					int32validator.Between(1, 20940),
				},
			},
			"max_retry": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "最大重试次数，取值范围：1-10次，默认为2次，支持更新",
				Default:     int32default.StaticInt32(2),
				Validators: []validator.Int32{
					int32validator.Between(1, 10),
				},
			},
			"http_method": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "仅当protocol为HTTP时必填且生效，HTTP请求的方法默认GET，{GET/HEAD/POST/PUT/DELETE/TRACE/OPTIONS/CONNECT/PATCH}，支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.HttpMethods...),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.HealthCheckProtocolHTTP),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.HealthCheckProtocolUDP),
						types.StringValue(business.HealthCheckProtocolTCP),
					),
				},
			},
			"http_url_path": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "仅当protocol为HTTP时必填且生效,支持的最大字符长度：80，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(80),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.HealthCheckProtocolHTTP),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("protocol"),
						types.StringValue(business.HealthCheckProtocolUDP),
						types.StringValue(business.HealthCheckProtocolTCP),
					),
				},
			},
			"http_expected_codes": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "利用逗号分割，仅当protocol为HTTP时必填且生效,支持http_2xx/http_3xx/http_4xx/http_5xx，一个或者多个的列表, 当 protocol 为 HTTP 时, 不填默认为 http_2xx，支持更新",
				Validators: []validator.Set{
					validator2.AlsoRequiresEqualSet(
						path.MatchRoot("protocol"),
						types.StringValue(business.HealthCheckProtocolHTTP),
					),
					validator2.ConflictsWithEqualSet(
						path.MatchRoot("protocol"),
						types.StringValue(business.HealthCheckProtocolUDP),
						types.StringValue(business.HealthCheckProtocolTCP),
					),
				},
			},
			"protocol_port": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "健康检查端口 1 - 65535，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "健康检查ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"status": schema.Int32Attribute{
				Computed:    true,
				Description: "状态 1 - UP, 0 - DOWN",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
			},
		},
	}
}

func (c *CtyunElbHealthCheck) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunElbHealthCheckConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建健康检查
	err = c.createHealthCheck(ctx, &plan)
	if err != nil {
		return
	}
	// 创建后反查创建后的nat信息
	err = c.getAndMergeHealthCheck(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (c *CtyunElbHealthCheck) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunElbHealthCheckConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeHealthCheck(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunElbHealthCheck) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置
	var plan CtyunElbHealthCheckConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunElbHealthCheckConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新健康检查信息
	err = c.updateHealthCheck(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeHealthCheck(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (c *CtyunElbHealthCheck) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunElbHealthCheckConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	params := &ctelb.CtelbDeleteHealthCheckRequest{
		ClientToken:   uuid.NewString(),
		RegionID:      state.RegionID.ValueString(),
		HealthCheckID: state.ID.ValueString(),
	}
	// 同步接口，无需轮训
	// 同步接口，无需轮训
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbDeleteHealthCheckApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
}
func (c *CtyunElbHealthCheck) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunElbHealthCheck) createHealthCheck(ctx context.Context, plan *CtyunElbHealthCheckConfig) (err error) {
	params := &ctelb.CtelbCreateHealthCheckRequest{
		ClientToken: uuid.NewString(),
		RegionID:    plan.RegionID.ValueString(),
		Name:        plan.Name.ValueString(),
		Protocol:    plan.Protocol.ValueString(),
		Timeout:     plan.Timeout.ValueInt32(),
		Interval:    plan.Interval.ValueInt32(),
		MaxRetry:    plan.MaxRetry.ValueInt32(),
	}
	if !plan.Description.IsNull() {
		params.Description = plan.Description.ValueString()
	}
	if plan.Timeout.ValueInt32() > 0 {
		params.Timeout = plan.Timeout.ValueInt32()
	}
	if plan.Interval.ValueInt32() > 0 {
		params.Interval = plan.Interval.ValueInt32()
	}
	if plan.MaxRetry.ValueInt32() > 0 {
		params.MaxRetry = plan.MaxRetry.ValueInt32()
	}
	if !plan.HttpMethod.IsNull() {
		params.HttpMethod = plan.HttpMethod.ValueString()
	}
	if !plan.HttpUrlPath.IsNull() {
		params.HttpUrlPath = plan.HttpUrlPath.ValueString()
	}
	httpExpectedCodes := []string{"http_2xx"}
	if len(plan.HttpExpectedCodes.Elements()) > 0 {
		diagnostics := plan.HttpExpectedCodes.ElementsAs(ctx, &httpExpectedCodes, false)
		if diagnostics.HasError() {
			return
		}
		params.HttpExpectedCodes = httpExpectedCodes
	}
	if !plan.ProtocolPort.IsNull() && !plan.ProtocolPort.IsUnknown() {
		params.ProtocolPort = plan.ProtocolPort.ValueInt32()
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbCreateHealthCheckApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	//解析返回值
	if resp.ReturnObj.ID == "" {
		err = fmt.Errorf("返回的健康检查ID为空！")
		return
	}
	plan.ID = types.StringValue(resp.ReturnObj.ID)

	return
}

func (c *CtyunElbHealthCheck) updateHealthCheck(ctx context.Context, state *CtyunElbHealthCheckConfig, plan *CtyunElbHealthCheckConfig) (err error) {
	params := &ctelb.CtelbUpdateHealthCheckRequest{
		ClientToken:   uuid.NewString(),
		RegionID:      state.RegionID.ValueString(),
		HealthCheckID: state.ID.ValueString(),
	}
	if plan.Name.ValueString() != "" && !plan.Name.Equal(state.Name) {
		params.Name = plan.Name.ValueString()
	}
	if plan.Description.ValueString() != "" && !plan.Description.Equal(state.Description) {
		params.Description = plan.Description.ValueString()
	}
	if !plan.Timeout.IsNull() {
		params.Timeout = plan.Timeout.ValueInt32()
	}
	if !plan.MaxRetry.IsNull() {
		params.MaxRetry = plan.MaxRetry.ValueInt32()
	}
	if !plan.Interval.IsNull() {
		params.Interval = plan.Interval.ValueInt32()
	}
	if plan.HttpMethod.ValueString() != "" && !plan.HttpMethod.Equal(state.HttpMethod) {
		params.HttpMethod = plan.HttpMethod.ValueString()
	}
	if plan.HttpUrlPath.ValueString() != "" && !plan.HttpUrlPath.Equal(state.HttpUrlPath) {
		params.HttpUrlPath = plan.HttpUrlPath.ValueString()
	}

	httpExpectedCodes := []string{"http_2xx"}
	if state.Protocol.ValueString() == business.HealthCheckProtocolHTTP {
		if len(plan.HttpExpectedCodes.Elements()) > 0 {
			var diagnostics diag.Diagnostics
			diagnostics = plan.HttpExpectedCodes.ElementsAs(ctx, &params.HttpExpectedCodes, false)
			if diagnostics.HasError() {
				return
			}
		} else {
			params.HttpExpectedCodes = httpExpectedCodes
		}
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateHealthCheckApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}

	return
}

func (c *CtyunElbHealthCheck) getAndMergeHealthCheck(ctx context.Context, plan *CtyunElbHealthCheckConfig) (err error) {
	params := &ctelb.CtelbShowHealthCheckRequest{
		RegionID:      plan.RegionID.ValueString(),
		HealthCheckID: plan.ID.ValueString(),
	}
	// 请求查看健康检查详情接口
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbShowHealthCheckApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析详情
	plan.Status = types.Int32Value(resp.ReturnObj.Status)
	plan.CreateTime = types.StringValue(resp.ReturnObj.CreateTime)
	plan.Description = types.StringValue(resp.ReturnObj.Description)
	plan.Name = types.StringValue(resp.ReturnObj.Name)
	plan.Protocol = types.StringValue(resp.ReturnObj.Protocol)
	plan.ProtocolPort = types.Int32Value(resp.ReturnObj.ProtocolPort)
	plan.Timeout = types.Int32Value(resp.ReturnObj.Timeout)
	plan.Interval = types.Int32Value(resp.ReturnObj.Interval)
	plan.MaxRetry = types.Int32Value(resp.ReturnObj.MaxRetry)
	plan.HttpMethod = types.StringValue(resp.ReturnObj.HttpMethod)
	plan.HttpUrlPath = types.StringValue(resp.ReturnObj.HttpUrlPath)
	// HttpExpectedCodes解析
	var HttpExpectedCodesList []string
	HttpExpectedCodesList = strings.Split(resp.ReturnObj.HttpExpectedCodes, " ")
	var HttpExpectedCodes []types.String
	for _, HttpExpectedCode := range HttpExpectedCodesList {
		HttpExpectedCodes = append(HttpExpectedCodes, types.StringValue(HttpExpectedCode))
	}
	var diagnostics diag.Diagnostics
	plan.HttpExpectedCodes, diagnostics = types.SetValueFrom(ctx, types.StringType, HttpExpectedCodes)
	if diagnostics.HasError() {
		return
	}

	return
}

type CtyunElbHealthCheckConfig struct {
	RegionID          types.String `tfsdk:"region_id"`           //区域ID
	Name              types.String `tfsdk:"name"`                //	唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32
	Description       types.String `tfsdk:"description"`         //	支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128
	Protocol          types.String `tfsdk:"protocol"`            //健康检查协议。取值范围：TCP、UDP、HTTP
	Timeout           types.Int32  `tfsdk:"timeout"`             //健康检查响应的最大超时时间，取值范围：2-60秒，默认为2秒
	Interval          types.Int32  `tfsdk:"interval"`            //负载均衡进行健康检查的时间间隔，取值范围：1-20940秒，默认为5秒
	MaxRetry          types.Int32  `tfsdk:"max_retry"`           //最大重试次数，取值范围：1-10次，默认为2次
	HttpMethod        types.String `tfsdk:"http_method"`         //仅当protocol为HTTP时必填且生效,HTTP请求的方法默认GET，{GET/HEAD/POST/PUT/DELETE/TRACE/OPTIONS/CONNECT/PATCH}
	HttpUrlPath       types.String `tfsdk:"http_url_path"`       //仅当protocol为HTTP时必填且生效,默认为'/',支持的最大字符长度：80
	HttpExpectedCodes types.Set    `tfsdk:"http_expected_codes"` //仅当protocol为HTTP时必填且生效,支持http_2xx/http_3xx/http_4xx/http_5xx，一个或者多个的列表, 当 protocol 为 HTTP 时, 不填默认为 http_2xx
	ProtocolPort      types.Int32  `tfsdk:"protocol_port"`       //健康检查端口 1 - 65535
	ID                types.String `tfsdk:"id"`                  //健康检查ID
	ProjectID         types.String `tfsdk:"project_id"`          //	项目ID
	Status            types.Int32  `tfsdk:"status"`              //状态 1 表示 UP, 0 表示 DOWN
	CreateTime        types.String `tfsdk:"create_time"`         //	创建时间，为UTC格式
}
