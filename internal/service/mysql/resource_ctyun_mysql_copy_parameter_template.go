package mysql

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CtyunMysqlCopyParamTemplate{}
	_ resource.ResourceWithConfigure   = &CtyunMysqlCopyParamTemplate{}
	_ resource.ResourceWithImportState = &CtyunMysqlCopyParamTemplate{}
)

type CtyunMysqlCopyParamTemplate struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMysqlCopyParamTemplate) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_copy_param_template"
}
func NewCtyunMysqlCopyParamTemplate() resource.Resource {
	return &CtyunMysqlCopyParamTemplate{}
}

func (c *CtyunMysqlCopyParamTemplate) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMysqlCopyParamTemplate) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	return
}

func (c *CtyunMysqlCopyParamTemplate) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10035290",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "模板参数名",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_parameter_template_id": schema.Int64Attribute{
				Required:    true,
				Description: "源参数模板id",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "参数模板描述",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.Desc(),
				},
			},
		},
	}
}

func (c *CtyunMysqlCopyParamTemplate) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMysqlCopyParamTemplateConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 开始创建参数模板
	err = c.CtyunMysqlCopyParamTemplateConfig(ctx, &plan)
	if err != nil {
		return
	}

	// 复制后，确认复制的模板时候存在
	err = c.checkCopySuccess(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlCopyParamTemplate) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	return
}

func (c *CtyunMysqlCopyParamTemplate) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunMysqlCopyParamTemplate) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	return
}

func (c *CtyunMysqlCopyParamTemplate) CtyunMysqlCopyParamTemplateConfig(ctx context.Context, config *CtyunMysqlCopyParamTemplateConfig) error {
	params := &mysql.TeledbCopyParameterTemplateRequest{
		SourceParameterGroupId: config.sourceParameterTemplateID.ValueInt64(),
		ParameterGroupName:     config.Name.ValueString(),
	}
	if !config.Description.IsNull() {
		params.ParameterGroupDesc = config.Description.ValueStringPointer()
	}
	header := &mysql.TeledbCopyParameterTemplateRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() || !config.ProjectID.IsUnknown() {
		config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbCopyParameterTemplateApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("复制id=%d的参数模板失败，接口返回nil。请联系研发确认问题原因！", config.sourceParameterTemplateID.ValueInt64())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	return nil
}

func (c *CtyunMysqlCopyParamTemplate) checkCopySuccess(ctx context.Context, config *CtyunMysqlCopyParamTemplateConfig) error {
	params := mysql.TeledbGetParameterTemplateListRequest{
		ParameterGroupName: config.Name.ValueStringPointer(),
		PageNow:            1,
		PageSize:           100,
	}

	header := mysql.TeledbGetParameterTemplateListRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
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
	if len(resp.ReturnObj.List) < 1 {
		err = fmt.Errorf("未查询到参数模板(name=%s)列表", config.Name.ValueString())
		return err
	}
	if len(resp.ReturnObj.List) > 1 {
		err = fmt.Errorf("查询到多条参数模板(name=%s)列表,查询结果为：%#v", config.Name.ValueString(), resp.ReturnObj.List)
		return err
	}
	return nil
}

type CtyunMysqlCopyParamTemplateConfig struct {
	RegionID                  types.String `tfsdk:"region_id"`
	ProjectID                 types.String `tfsdk:"project_id"`
	sourceParameterTemplateID types.Int64  `tfsdk:"source_parameter_template_id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
}
