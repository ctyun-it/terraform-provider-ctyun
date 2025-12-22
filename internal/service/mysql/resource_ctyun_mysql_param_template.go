package mysql

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

var (
	_ resource.Resource                = &CtyunMysqlParamTemplate{}
	_ resource.ResourceWithConfigure   = &CtyunMysqlParamTemplate{}
	_ resource.ResourceWithImportState = &CtyunMysqlParamTemplate{}
)

type CtyunMysqlParamTemplate struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMysqlParamTemplate) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_param_template"
}
func NewCtyunMysqlParamTemplate() resource.Resource {
	return &CtyunMysqlParamTemplate{}
}

func (c *CtyunMysqlParamTemplate) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMysqlParamTemplate) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[projectID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunMysqlParamTemplateConfig
	var id, regionID, projectID string
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		projectID = c.meta.GetExtraIfEmpty(projectID, common.ExtraProjectId)
		id = request.ID
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &id, &projectID, &regionID)
		if err != nil {
			return
		}
	}
	if id == "" {
		err = fmt.Errorf("id不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}

	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		err = fmt.Errorf("转换失败: %v\n", err)
		return
	}
	config.ID = types.Int64Value(num)
	config.RegionID = types.StringValue(regionID)
	config.ProjectID = types.StringValue(projectID)
	err = c.getAndMergeMysqlParameterTemplate(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunMysqlParamTemplate) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10098794",
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
			"engine": schema.StringAttribute{
				Required:    true,
				Description: "数据库版本，取值：5.7, 8.0",
				Validators: []validator.String{
					stringvalidator.OneOf("5.7", "8.0"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "参数模板id",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"template_parameters": schema.MapAttribute{
				Optional:    true,
				Description: "mysql模板参数列表，创建参数模板时不可传，更新阶段可传，支持更新。可修改每个参数值，无法新增参数或删除",
				ElementType: types.StringType,
				Validators: []validator.Map{
					mapvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

func (c *CtyunMysqlParamTemplate) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMysqlParamTemplateConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 开始创建参数模板
	err = c.CreateMysqlParameterTemplate(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取mysql详情
	err = c.getAndMergeMysqlParameterTemplate(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlParamTemplate) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunMysqlParamTemplateConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergeMysqlParameterTemplate(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlParamTemplate) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunMysqlParamTemplateConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunMysqlParamTemplateConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.updateMysqlParameterTemplate(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeMysqlParameterTemplate(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlParamTemplate) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunMysqlParamTemplateConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.deleteMysqlParameterTemplate(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunMysqlParamTemplate) CreateMysqlParameterTemplate(ctx context.Context, config *CtyunMysqlParamTemplateConfig) error {
	if !config.TemplateParameters.IsNull() {
		err := fmt.Errorf("创建阶段的参数不允许填写！")
		return err
	}

	params := &mysql.TeledbCreateParameterTemplateRequest{
		ParameterGroupName: config.Name.ValueString(),
		Engine:             config.Engine.ValueString(),
	}
	if !config.Description.IsNull() {
		params.Description = config.Description.ValueStringPointer()
	}
	header := &mysql.TeledbCreateParameterTemplateRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbCreateParameterTemplateApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建参数模板失败，接口返回nil。请与研发联系确认问题原因！")
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}

	return nil
}

func (c *CtyunMysqlParamTemplate) getAndMergeMysqlParameterTemplate(ctx context.Context, config *CtyunMysqlParamTemplateConfig) error {
	// 如果id为空，通过查询list获取
	if config.ID.IsNull() || config.ID.IsUnknown() {
		templateList, err := c.getIDByParameterTemplateName(ctx, config)
		if err != nil {
			return err
		}
		config.ID = types.Int64Value(templateList[0].ID)
	}
	if config.TemplateParameters.IsNull() {
		config.TemplateParameters = types.MapNull(types.StringType)
	}
	return nil
}

func (c *CtyunMysqlParamTemplate) getParameterTemplateDetail(ctx context.Context, config *CtyunMysqlParamTemplateConfig) ([]mysql.TeledbGetParameterTemplateDetailResponseReturnObjDeatil, error) {
	var pageNo, pageSize int32
	var parameters []mysql.TeledbGetParameterTemplateDetailResponseReturnObjDeatil
	pageNo = 1
	pageSize = 100
	resp, err := c.getParameterListByPage(ctx, config, pageNo, pageSize)
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
		resp, err = c.getParameterListByPage(ctx, config, pageNo, pageSize)
		if err != nil {
			return nil, err
		}
	}
	return parameters, nil
}

func (c *CtyunMysqlParamTemplate) getParameterListByPage(ctx context.Context, config *CtyunMysqlParamTemplateConfig, pageNo int32, pageSize int32) (*mysql.TeledbGetParameterTemplateDetailResponse, error) {
	params := &mysql.TeledbGetParameterTemplateDetailRequest{
		ID:       config.ID.ValueInt64(),
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
		err = fmt.Errorf("查询参数模板详情(id=%d)失败，接口返回nil。请联系研发确认问题原因！", config.ID.ValueInt64())
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

func (c *CtyunMysqlParamTemplate) getIDByParameterTemplateName(ctx context.Context, config *CtyunMysqlParamTemplateConfig) ([]mysql.ParameterTemplateInfo, error) {
	params := mysql.TeledbGetParameterTemplateListRequest{
		ParameterGroupName: config.Name.ValueStringPointer(),
		EngineVersion:      config.Engine.ValueStringPointer(),
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
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询参数模板(name=%s)列表失败，接口返回nil。请联系研发确认问题原因！", config.Name.ValueString())
		return nil, err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	if len(resp.ReturnObj.List) < 1 {
		err = fmt.Errorf("未查询到参数模板(name=%s)列表", config.Name.ValueString())
		return nil, err
	}
	if len(resp.ReturnObj.List) > 1 {
		err = fmt.Errorf("查询到多条参数模板(name=%s)列表,查询结果为：%#v", config.Name.ValueString(), resp.ReturnObj.List)
		return nil, err
	}
	return resp.ReturnObj.List, nil
}

func (c *CtyunMysqlParamTemplate) updateMysqlParameterTemplate(ctx context.Context, state *CtyunMysqlParamTemplateConfig, plan *CtyunMysqlParamTemplateConfig) error {
	// 支持更新模板中参数
	// 若plan阶段的template_parameters有值，说明需要更新
	if plan.TemplateParameters.IsNull() {
		return nil
	}

	// 对比state和plan参数，确认哪些参数value发生变化
	//updateParameterList, stateParameterMap, err := c.compareParameters(ctx, state, plan)
	//if err != nil {
	//	return err
	//}
	oldParameters, err := c.getOldParameterValue(ctx, state)
	if err != nil {
		return err
	}
	updateParameters, err := utils.TypesMapToStringMap(ctx, plan.TemplateParameters)
	if err != nil {
		return err
	}
	// 若待更新参数 > 0 触发参数值更新操作
	if len(updateParameters) > 0 {
		params := &mysql.TeledbUpdateParameterTemplateRequest{
			ID: state.ID.ValueInt64(),
		}
		var updateParameterObjs []mysql.ParameterObj
		for parameterName, parameterValue := range updateParameters {
			var parameterObj mysql.ParameterObj
			parameterObj.ParameterName = parameterName
			parameterObj.ParameterValue = parameterValue
			if parameter, exists := oldParameters[parameterName]; exists {
				parameterObj.OldValue = parameter.ParameterValue
				parameterObj.Restart = parameter.Restart
			} else {
				err = fmt.Errorf("参数名称：%s无法更新，参数列表中查询到该参数，请确认后进行更新操作！", parameterName)
				return err
			}
			updateParameterObjs = append(updateParameterObjs, parameterObj)
		}
		params.Value = updateParameterObjs
		header := &mysql.TeledbUpdateParameterTemplateRequestHeader{
			RegionID: state.RegionID.ValueString(),
		}
		if !state.ProjectID.IsNull() {
			header.ProjectID = state.ProjectID.ValueStringPointer()
		}
		resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbUpdateParameterTemplateApi.Do(ctx, c.meta.Credential, params, header)
		if err2 != nil {
			return err2
		} else if resp == nil {
			err = fmt.Errorf("")
		} else if resp.StatusCode != 0 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return err
		}
	}
	state.TemplateParameters = plan.TemplateParameters
	return nil
}

func (c *CtyunMysqlParamTemplate) deleteMysqlParameterTemplate(ctx context.Context, config CtyunMysqlParamTemplateConfig) error {
	// 调用删除接口
	params := &mysql.TeledbDeleteParameterTemplateRequest{
		ID: config.ID.ValueInt64(),
	}
	header := &mysql.TeledbDeleteParameterTemplateRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbDeleteParameterTemplateApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("参数id=%d的参数模板失败，接口返回nil。请联系研发确认问题原因", config.ID.ValueInt64())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	// 调用列表接口，确认是否已经删除
	_, err = c.getIDByParameterTemplateName(ctx, &config)
	if err != nil {
		if strings.Contains(err.Error(), "未查询到参数模板") {
			return nil
		}
		return err
	}
	return nil
}

func (c *CtyunMysqlParamTemplate) compareParameters(ctx context.Context, state *CtyunMysqlParamTemplateConfig, plan *CtyunMysqlParamTemplateConfig) ([]CtyunMysqlParameterModel, map[int64]CtyunMysqlParameterModel, error) {
	// 定义一个更新列表
	var updateParameterList []CtyunMysqlParameterModel
	// 定义两个map，int64: object，
	stateParameterMap, err := c.initParameterMap(ctx, state)
	if err != nil {
		return nil, nil, err
	}
	planParameterMap, err := c.initParameterMap(ctx, plan)
	if err != nil {
		return nil, nil, err
	}
	if len(stateParameterMap) != len(planParameterMap) {
		err = fmt.Errorf("参数不支持增加或删除，原参数个数：%d，更新参数个数为%d", len(stateParameterMap), len(planParameterMap))
		return nil, nil, err
	}
	// 对比
	// 遍历plan列表与state做对比
	for parameterID, planParameterItem := range planParameterMap {
		stateParameter, ok := stateParameterMap[parameterID]
		if !ok {
			err = fmt.Errorf("参数不支持增加，parameter_name=%s为新增参数，state阶段不存在。", planParameterItem.ParameterName)
		}
		flag, err2 := c.compareDetail(stateParameter, planParameterItem)
		if err2 != nil {
			return nil, nil, err2
		}
		if flag {
			updateParameterList = append(updateParameterList, planParameterItem)
		}
	}
	return updateParameterList, stateParameterMap, nil
}

func (c *CtyunMysqlParamTemplate) initParameterMap(ctx context.Context, state *CtyunMysqlParamTemplateConfig) (map[int64]CtyunMysqlParameterModel, error) {
	var parameterMap map[int64]CtyunMysqlParameterModel
	parameterMap = make(map[int64]CtyunMysqlParameterModel)
	var parameters []CtyunMysqlParameterModel
	diags := state.TemplateParameters.ElementsAs(ctx, &parameters, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return nil, err
	}
	for _, parameter := range parameters {
		id := parameter.ParameterID.ValueInt64()
		if _, exist := parameterMap[id]; !exist {
			parameterMap[id] = parameter
		}
	}
	return parameterMap, nil
}

func (c *CtyunMysqlParamTemplate) compareDetail(stateParameter CtyunMysqlParameterModel, planParameter CtyunMysqlParameterModel) (bool, error) {
	var flag bool
	flag = false
	if stateParameter.ParameterName.ValueString() != planParameter.ParameterName.ValueString() {
		err := fmt.Errorf("parameter_name不支持修改，state阶段parameter_name=%s, plan阶段parameter_name=%s", stateParameter.ParameterName.ValueString(), planParameter.ParameterName.ValueString())
		return flag, err
	}
	if !stateParameter.Restart.Equal(planParameter.Restart) {
		err := fmt.Errorf("restart不支持修改")
		return flag, err
	}
	if !stateParameter.PermitValue.Equal(planParameter.PermitValue) {
		err := fmt.Errorf("permit_value不支持更新")
		return flag, err
	}
	if !stateParameter.Description.Equal(planParameter.Description) {
		err := fmt.Errorf("description不支持更新")
		return flag, err
	}
	if !stateParameter.ParameterValue.Equal(planParameter.ParameterValue) {
		flag = true
		return flag, nil
	}
	return flag, nil
}

func (c *CtyunMysqlParamTemplate) getOldParameterValue(ctx context.Context, state *CtyunMysqlParamTemplateConfig) (map[string]mysql.TeledbGetParameterTemplateDetailResponseReturnObjDeatil, error) {
	parameters, err := c.getParameterTemplateDetail(ctx, state)
	if err != nil {
		return nil, err
	}
	oldParameters := make(map[string]mysql.TeledbGetParameterTemplateDetailResponseReturnObjDeatil)
	for _, parameter := range parameters {
		if _, exist := oldParameters[parameter.ParameterName]; !exist {
			oldParameters[parameter.ParameterName] = parameter
		}
	}
	return oldParameters, nil
}

type CtyunMysqlParamTemplateConfig struct {
	RegionID           types.String `tfsdk:"region_id"`
	ProjectID          types.String `tfsdk:"project_id"`
	Name               types.String `tfsdk:"name"`
	Engine             types.String `tfsdk:"engine"`
	Description        types.String `tfsdk:"description"`
	ID                 types.Int64  `tfsdk:"id"`
	TemplateParameters types.Map    `tfsdk:"template_parameters"`
}
type CtyunMysqlParameterModel struct {
	ParameterName  types.String `tfsdk:"parameter_name"`
	Restart        types.Bool   `tfsdk:"restart"`
	ParameterValue types.String `tfsdk:"parameter_value"`
	PermitValue    types.String `tfsdk:"permit_value"`
	Description    types.String `tfsdk:"description"`
	ParameterID    types.Int64  `tfsdk:"parameter_id"`
}
