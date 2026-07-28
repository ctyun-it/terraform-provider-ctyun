package pgsql

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	explanmodifier "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/planmodifier"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
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
	_ resource.Resource                = &CtyunPgsqlParamTemplate{}
	_ resource.ResourceWithConfigure   = &CtyunPgsqlParamTemplate{}
	_ resource.ResourceWithImportState = &CtyunPgsqlParamTemplate{}
)

type CtyunPgsqlParamTemplate struct {
	meta *common.CtyunMetadata
	name string
}

func (c *CtyunPgsqlParamTemplate) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_param_template"
	c.name = response.TypeName
}
func NewCtyunPgsqlParamTemplate() resource.Resource {
	return &CtyunPgsqlParamTemplate{}
}

func (c *CtyunPgsqlParamTemplate) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPgsqlParamTemplate) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {

	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [id],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunPgsqlParameterTemplateConfig
	var id, regionID string
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		id = request.ID
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &id, &regionID)
		if err != nil {
			return
		}
	}
	if id == "" {
		err = fmt.Errorf("id不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}

	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		err = fmt.Errorf("转换失败: %v\n", err)
		return
	}
	config.ID = types.Int64Value(num)
	config.RegionID = types.StringValue(regionID)
	err = c.getAndMergePostgresqlParameterTemplate(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunPgsqlParamTemplate) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理PostgreSQL参数模板", "关系数据库PostgreSQL版", "https://www.ctyun.cn/document/10034019/10166169"),
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
				Optional:           true,
				DeprecationMessage: "废弃字段，请不要指定",
				Description:        "企业项目ID",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "数据库参数模板名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
				},
			},
			"source_template_id": schema.Int64Attribute{
				Optional:    true,
				Description: "参考的参数模板ID，创建时必填。可以根据data.ctyun_postgresql_param_templates查询",
				PlanModifiers: []planmodifier.Int64{
					explanmodifier.NullIgnoreInt64(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "参数模板的描述，支持更新，若不为空，则长度限制：1-1024",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 1024),
				},
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "参数模板id",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"template_parameters": schema.MapAttribute{
				Optional:    true,
				Description: "postgresql模板参数列表，创建参数模板时不可传，更新阶段可传，支持更新。可修改每个参数值，无法新增参数或删除",
				ElementType: types.StringType,
			},
		},
	}
}

func (c *CtyunPgsqlParamTemplate) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunPgsqlParameterTemplateConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 开始创建参数模板
	err = c.CreatePostgresqlParameterTemplate(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取mysql详情
	err = c.getAndMergePostgresqlParameterTemplate(ctx, &plan)
	if err != nil {
		return
	}

	//plan.ID = types.StringValue(plan.BackupName.ValueString())
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPgsqlParamTemplate) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPgsqlParameterTemplateConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergePostgresqlParameterTemplate(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPgsqlParamTemplate) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunPgsqlParameterTemplateConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunPgsqlParameterTemplateConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	if !plan.SourceTemplateId.IsUnknown() && !plan.SourceTemplateId.IsNull() && state.SourceTemplateId.IsNull() {
		state.SourceTemplateId = plan.SourceTemplateId
		response.Diagnostics.AddWarning("source_template_id的更新仅写入状态文件", "在import时，状态文件中source_template_id为null，允许用模板中的值进行一次更新，该更新不触发远程调用")
	}

	err = c.updatePgsqlParameters(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergePostgresqlParameterTemplate(ctx, &state)
	if err != nil {
		return
	}
	state.ProjectID = plan.ProjectID
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPgsqlParamTemplate) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunPgsqlParameterTemplateConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	params := &pgsql.PgsqlDeleteParameterTemplateRequest{
		TemplateId: config.ID.ValueInt64(),
	}
	header := &pgsql.PgsqlDeleteParameterTemplateRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlDeleteParameterTemplateApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("删除参数模板(name=%s)失败，接口返回nil，请联系研发确认问题原因！", config.Name.ValueString())
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return
	}
	return
}

func (c *CtyunPgsqlParamTemplate) CreatePostgresqlParameterTemplate(ctx context.Context, config *CtyunPgsqlParameterTemplateConfig) error {
	if config.SourceTemplateId.IsNull() || config.SourceTemplateId.IsUnknown() {
		err := fmt.Errorf("创建postgresql 参数模板时，source_template_id必填！")
		return err
	}
	params := &pgsql.PgsqlCreateParameterTemplateRequest{
		SourceTemplateId: config.SourceTemplateId.ValueInt64(),
		Name:             config.Name.ValueString(),
	}
	if !config.Description.IsNull() {
		params.Description = config.Description.ValueStringPointer()
	}
	header := &pgsql.PgsqlCreateParameterTemplateRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlCreateParameterTemplateApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建参数模板失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	config.ID = types.Int64Value(resp.ReturnObj.TemplateId)
	return nil
}

func (c *CtyunPgsqlParamTemplate) getAndMergePostgresqlParameterTemplate(ctx context.Context, config *CtyunPgsqlParameterTemplateConfig) error {
	detail, err := c.getPgsqlParameterTemplateList(ctx, config)
	if err != nil {
		return err
	}
	config.ID = types.Int64Value(detail.PgTemplateId)
	config.Description = types.StringValue(detail.Description)
	config.Name = types.StringValue(detail.Name)
	if config.TemplateParameters.IsNull() {
		config.TemplateParameters = types.MapNull(types.StringType)
	}
	return nil
}

func (c *CtyunPgsqlParamTemplate) getTemplateParametersValue(ctx context.Context, config *CtyunPgsqlParameterTemplateConfig) (map[string]string, error) {
	templateParameterMap := make(map[string]string)
	params := &pgsql.PgsqlGetParameterTemplateDetailRequest{
		TemplateId: config.ID.ValueInt64(),
	}
	header := &pgsql.PgsqlGetParameterTemplateDetailRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGetParameterTemplateDetailApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询模板id=%d的参数信息失败，接口返回nil，请联系研发确认问题原因！", config.ID.ValueInt64())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return nil, err
	}
	for _, parameters := range resp.ReturnObj {
		parameterName := parameters.ParameterName
		parameterValue := parameters.ParameterValue
		if _, exist := templateParameterMap[parameterName]; !exist {
			templateParameterMap[parameterName] = parameterValue
		}
	}
	return templateParameterMap, nil
}

func (c *CtyunPgsqlParamTemplate) updatePgsqlParameters(ctx context.Context, state *CtyunPgsqlParameterTemplateConfig, plan *CtyunPgsqlParameterTemplateConfig) error {

	if !plan.Description.IsNull() && !plan.Description.Equal(state.Description) {
		err := c.updatePgsqlParameterTemplate(ctx, state, plan)
		if err != nil {
			return err
		}
		state.Description = plan.Description
	}

	oldParameters, err := c.getTemplateParametersValue(ctx, state)
	if err != nil {
		return err
	}
	updateParameters, err := utils.TypesMapToStringMap(ctx, plan.TemplateParameters)
	if err != nil {
		return err
	}
	if len(updateParameters) > 0 {
		params := &pgsql.PgsqlUpdateParameterTemplateRequest{
			TemplateId: state.ID.ValueInt64(),
		}
		header := &pgsql.PgsqlUpdateParameterTemplateRequestHeader{
			RegionID: state.RegionID.ValueString(),
		}

		var updateParameterObjs []pgsql.ParameterObj
		for parameterName, parameterValue := range updateParameters {
			var parameterObj pgsql.ParameterObj
			if _, exists := oldParameters[parameterName]; !exists {
				err = fmt.Errorf("参数名称：%s无法更新，参数列表中查询到该参数，请确认后进行更新操作！", parameterName)
				return err
			}
			parameterObj.Name = parameterName
			parameterObj.Value = parameterValue
			updateParameterObjs = append(updateParameterObjs, parameterObj)
		}
		params.Params = updateParameterObjs
		resp, err2 := c.meta.Apis.SdkCtPgsqlApis.PgsqlUpdateParameterTemplateApi.Do(ctx, c.meta.Credential, params, header)
		if err2 != nil {
			return err2
		} else if resp == nil {
			err = fmt.Errorf("更新参数模板状态失败，模板id=%d", state.ID.ValueInt64())
			return err
		}
		state.TemplateParameters = plan.TemplateParameters
	}
	return nil
}

func (c *CtyunPgsqlParamTemplate) getPgsqlParameterTemplateList(ctx context.Context, config *CtyunPgsqlParameterTemplateConfig) (*pgsql.ParameterTemplateInfo, error) {
	params := &pgsql.PgsqlGetParameterTemplateListRequest{
		PageNow:  1,
		PageSize: 500,
	}
	header := &pgsql.PgsqlGetParameterTemplateListRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.Name.IsNull() && config.Name.IsUnknown() {
		params.Name = config.Name.ValueStringPointer()
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
	if resp.ReturnObj.List == nil || len(resp.ReturnObj.List) == 0 {
		err = common.ResourceNotExistError
		return nil, err
	}
	if len(resp.ReturnObj.List) > 1 {
		for _, item := range resp.ReturnObj.List {
			if !config.ID.IsNull() && !config.ID.IsUnknown() && config.ID.ValueInt64() == item.PgTemplateId {
				return &item, nil
			}
			if item.Name == config.Name.ValueString() && item.Description == config.Description.ValueString() {
				return &item, nil
			}
		}
	}
	return &resp.ReturnObj.List[0], nil
}

func (c *CtyunPgsqlParamTemplate) updatePgsqlParameterTemplate(ctx context.Context, state *CtyunPgsqlParameterTemplateConfig, plan *CtyunPgsqlParameterTemplateConfig) error {
	params := &pgsql.PgsqlUpdateParameterTemplateRemarkRequest{
		TemplateId:  state.ID.ValueInt64(),
		Description: plan.Description.ValueString(),
	}
	header := &pgsql.PgsqlUpdateParameterTemplateRemarkRequestHeader{
		RegionID: state.RegionID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlUpdateParameterTemplateRemarkApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新参数模板(id=%d，name=%s)备注失败，接口返回nil，请联系研发确认问题原因！", state.ID.ValueInt64(), state.Name.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	return nil
}

type CtyunPgsqlParameterTemplateConfig struct {
	RegionID           types.String `tfsdk:"region_id"`
	ProjectID          types.String `tfsdk:"project_id"`
	Name               types.String `tfsdk:"name"`
	SourceTemplateId   types.Int64  `tfsdk:"source_template_id"`
	Description        types.String `tfsdk:"description"`
	ID                 types.Int64  `tfsdk:"id"`
	TemplateParameters types.Map    `tfsdk:"template_parameters"`
}
