package redis

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgdcs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &ctyunRedisParamTemplate{}
	_ resource.ResourceWithConfigure   = &ctyunRedisParamTemplate{}
	_ resource.ResourceWithImportState = &ctyunRedisParamTemplate{}
)

type ctyunRedisParamTemplate struct {
	meta *common.CtyunMetadata
}

func NewCtyunRedisParamTemplate() resource.Resource {
	return &ctyunRedisParamTemplate{}
}

func (c *ctyunRedisParamTemplate) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_param_template"
}

type CtyunRedisParamTemplateConfig struct {
	ID           types.String `tfsdk:"id"`
	RegionId     types.String `tfsdk:"region_id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	CacheMode    types.String `tfsdk:"cache_mode"`
	SysTemplate  types.Bool   `tfsdk:"sys_template"`
	Params       types.Set    `tfsdk:"params"`
	ParamsReturn types.Set    `tfsdk:"params_return"`
}

// ParamInputModel 用于输入参数，只包含基本字段
type ParamInputModel struct {
	ParamName    types.String `tfsdk:"param_name"`
	CurrentValue types.String `tfsdk:"current_value"`
}

// ParamReturnModel 用于返回参数，包含所有字段
type ParamReturnModel struct {
	ParamName    types.String `tfsdk:"param_name"`
	CurrentValue types.String `tfsdk:"current_value"`
	Description  types.String `tfsdk:"description"`
	ValueRange   types.String `tfsdk:"value_range"`
	DefaultValue types.String `tfsdk:"default_value"`
	NeedRestart  types.Bool   `tfsdk:"need_restart"`
}

func (c *ctyunRedisParamTemplate) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029420/10156164`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "模板ID",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "参数模板名称 支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "参数模板描述 支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"cache_mode": schema.StringAttribute{
				Required:    true,
				Description: "适合的实例架构版本 ORIGINAL_67：Redis 6.0/7.0类型 ORIGINAL_5：Redis 5.0类型 CLASSIC：经典版",
				Validators: []validator.String{
					stringvalidator.OneOf("ORIGINAL_67", "ORIGINAL_5", "CLASSIC"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sys_template": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否为系统模板 true：系统模板 false：自定义模板",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"params": schema.SetNestedAttribute{
				Description: "输入的参数列表",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"param_name": schema.StringAttribute{
							Required:    true,
							Description: "参数名称",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
						"current_value": schema.StringAttribute{
							Required:    true,
							Description: "目标值 支持更新",
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
					},
				},
			},
			"params_return": schema.SetNestedAttribute{
				Description: "返回的参数列表",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"param_name": schema.StringAttribute{
							Computed:    true,
							Description: "参数名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "参数描述",
						},
						"current_value": schema.StringAttribute{
							Computed:    true,
							Description: "当前值",
						},
						"value_range": schema.StringAttribute{
							Computed:    true,
							Description: "参数范围",
						},
						"default_value": schema.StringAttribute{
							Computed:    true,
							Description: "默认值",
						},
						"need_restart": schema.BoolAttribute{
							Computed:    true,
							Description: "是否需要重启",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRedisParamTemplate) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunRedisParamTemplateConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建参数模板
	err = c.createParamTemplate(ctx, plan)
	if err != nil {
		return
	}
	id, err := c.getIdByQueryList(ctx, plan)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(id)
	err = c.getAndMerge(ctx, nil, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunRedisParamTemplate) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRedisParamTemplateConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMerge(ctx, nil, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRedisParamTemplate) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunRedisParamTemplateConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// state中的
	var state CtyunRedisParamTemplateConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 更新参数模板
	err = c.updateParamTemplate(ctx, plan, state)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRedisParamTemplate) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRedisParamTemplateConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 删除参数模板
	err = c.deleteParamTemplate(ctx, state)
	if err != nil {
		return
	}
	response.State.RemoveResource(ctx)
}

func (c *ctyunRedisParamTemplate) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunRedisParamTemplate) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [templateID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()

	var cfg CtyunRedisParamTemplateConfig

	var templateId, regionId string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		templateId = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &templateId, &regionId)
		if err != nil {
			return
		}
	}

	if templateId == "" {
		err = fmt.Errorf("templateID不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	cfg.RegionId = types.StringValue(regionId)
	cfg.ID = types.StringValue(templateId)
	// 查询远端
	err = c.getAndMerge(ctx, nil, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// createParamTemplate 创建参数模板
func (c *ctyunRedisParamTemplate) createParamTemplate(ctx context.Context, plan CtyunRedisParamTemplateConfig) (err error) {
	// 构建创建参数模板请求
	template := &ctgdcs2.Dcs2CreateRedisTemplateTemplateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		CacheMode:   plan.CacheMode.ValueString(),
		SysTemplate: plan.SysTemplate.ValueBoolPointer(),
	}

	// 处理参数列表
	var params []*ctgdcs2.Dcs2CreateRedisTemplateParamsRequest
	if !plan.Params.IsNull() && !plan.Params.IsUnknown() {
		var paramModels []ParamInputModel
		diags := plan.Params.ElementsAs(ctx, &paramModels, false)
		if diags.HasError() {
			err = fmt.Errorf("failed to parse params: %v", diags.Errors())
			return
		}

		for _, paramModel := range paramModels {
			params = append(params, &ctgdcs2.Dcs2CreateRedisTemplateParamsRequest{
				ParamName:    paramModel.ParamName.ValueString(),
				CurrentValue: paramModel.CurrentValue.ValueString(),
			})
		}
	}

	paramsReq := &ctgdcs2.Dcs2CreateRedisTemplateRequest{
		RegionId: plan.RegionId.ValueString(),
		Template: template,
		Params:   params,
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2CreateRedisTemplateApi.Do(ctx, c.meta.SdkCredential, paramsReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// updateParamTemplate 更新参数模板
func (c *ctyunRedisParamTemplate) updateParamTemplate(ctx context.Context, plan, state CtyunRedisParamTemplateConfig) (err error) {
	// 构建更新参数模板请求
	template := &ctgdcs2.Dcs2EditRedisTemplateTemplateRequest{
		Id:          state.ID.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		CacheMode:   state.CacheMode.ValueString(),
		SysTemplate: plan.SysTemplate.ValueBoolPointer(),
	}

	// 处理参数列表
	var params []*ctgdcs2.Dcs2EditRedisTemplateParamsRequest
	if !plan.Params.IsNull() && !plan.Params.IsUnknown() {
		var paramModels []ParamInputModel
		diags := plan.Params.ElementsAs(ctx, &paramModels, false)
		if diags.HasError() {
			err = fmt.Errorf("failed to parse params: %v", diags.Errors())
			return
		}

		for _, paramModel := range paramModels {
			params = append(params, &ctgdcs2.Dcs2EditRedisTemplateParamsRequest{
				ParamName:    paramModel.ParamName.ValueString(),
				CurrentValue: paramModel.CurrentValue.ValueString(),
			})
		}
	}

	paramsReq := &ctgdcs2.Dcs2EditRedisTemplateRequest{
		RegionId: plan.RegionId.ValueString(),
		Template: template,
		Params:   params,
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2EditRedisTemplateApi.Do(ctx, c.meta.SdkCredential, paramsReq)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// deleteParamTemplate 删除参数模板
func (c *ctyunRedisParamTemplate) deleteParamTemplate(ctx context.Context, state CtyunRedisParamTemplateConfig) (err error) {
	params := &ctgdcs2.Dcs2DeleteRedisTemplateRequest{
		RegionId:   state.RegionId.ValueString(),
		TemplateId: state.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DeleteRedisTemplateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// getAndMerge 从远端查询参数模板信息
func (c *ctyunRedisParamTemplate) getAndMerge(ctx context.Context, plan, state *CtyunRedisParamTemplateConfig) (err error) {
	// 调用API查询参数模板详情
	params := &ctgdcs2.Dcs2DescribeRedisTemplateDetailRequest{
		RegionId:   state.RegionId.ValueString(),
		TemplateId: state.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeRedisTemplateDetailApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil || resp.ReturnObj.Template == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 更新state中的信息
	state.Name = types.StringValue(resp.ReturnObj.Template.Name)
	state.Description = types.StringValue(resp.ReturnObj.Template.Description)
	state.CacheMode = types.StringValue(resp.ReturnObj.Template.CacheMode)
	if resp.ReturnObj.Template.SysTemplate != nil {
		state.SysTemplate = types.BoolValue(*resp.ReturnObj.Template.SysTemplate)
	}

	// 处理完整的参数列表（用于params_return）- 包含所有参数信息
	if len(resp.ReturnObj.Params) > 0 {
		paramModels := make([]ParamReturnModel, 0, len(resp.ReturnObj.Params))
		for _, param := range resp.ReturnObj.Params {
			paramModel := ParamReturnModel{
				ParamName:    types.StringValue(param.ParamName),
				CurrentValue: types.StringValue(param.CurrentValue),
			}

			// 设置可选字段
			if param.Description != "" {
				paramModel.Description = types.StringValue(param.Description)
			} else {
				paramModel.Description = types.StringNull()
			}

			if param.ValueRange != "" {
				paramModel.ValueRange = types.StringValue(param.ValueRange)
			} else {
				paramModel.ValueRange = types.StringNull()
			}

			if param.DefaultValue != "" {
				paramModel.DefaultValue = types.StringValue(param.DefaultValue)
			} else {
				paramModel.DefaultValue = types.StringNull()
			}

			if param.NeedRestart != nil {
				paramModel.NeedRestart = types.BoolValue(*param.NeedRestart)
			} else {
				paramModel.NeedRestart = types.BoolNull()
			}

			paramModels = append(paramModels, paramModel)
		}

		paramObjType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"param_name":    types.StringType,
				"current_value": types.StringType,
				"description":   types.StringType,
				"value_range":   types.StringType,
				"default_value": types.StringType,
				"need_restart":  types.BoolType,
			},
		}

		paramsValue, diags := types.SetValueFrom(ctx, paramObjType, paramModels)
		if diags.HasError() {
			err = fmt.Errorf("failed to set params_return: %v", diags)
			return
		}
		state.ParamsReturn = paramsValue
	} else {
		state.ParamsReturn = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"param_name":    types.StringType,
				"current_value": types.StringType,
				"description":   types.StringType,
				"value_range":   types.StringType,
				"default_value": types.StringType,
				"need_restart":  types.BoolType,
			},
		})
	}
	if plan != nil && !plan.Params.IsNull() && !plan.Params.IsUnknown() {
		state.Params = plan.Params
	}
	// 处理参数列表（用于params）- 只包含state中已有的参数
	if len(resp.ReturnObj.Params) > 0 && !state.Params.IsNull() && !state.Params.IsUnknown() {
		// 获取state中指定的参数名
		var stateParamModels []ParamInputModel
		diags := state.Params.ElementsAs(ctx, &stateParamModels, false)
		if diags.HasError() {
			err = fmt.Errorf("failed to parse state params: %v", diags.Errors())
			return
		}

		// 创建参数名到参数值的映射
		stateParamMap := make(map[string]string)
		for _, paramModel := range stateParamModels {
			stateParamMap[paramModel.ParamName.ValueString()] = paramModel.CurrentValue.ValueString()
		}

		// 过滤API返回的参数，只保留state中已有的参数
		var filteredParamModels []ParamInputModel
		for _, param := range resp.ReturnObj.Params {
			// 只保留state中已有的参数
			if _, exists := stateParamMap[param.ParamName]; exists {
				filteredParamModels = append(filteredParamModels, ParamInputModel{
					ParamName:    types.StringValue(param.ParamName),
					CurrentValue: types.StringValue(param.CurrentValue),
				})
			}
		}

		// 更新state中的参数列表
		if len(filteredParamModels) > 0 {
			paramObjType := types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"param_name":    types.StringType,
					"current_value": types.StringType,
				},
			}

			paramsValue, diags := types.SetValueFrom(ctx, paramObjType, filteredParamModels)
			if diags.HasError() {
				err = fmt.Errorf("failed to set params: %v", diags)
				return
			}
			state.Params = paramsValue
		} else {
			state.Params = types.SetNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"param_name":    types.StringType,
					"current_value": types.StringType,
				},
			})
		}
	} else {
		// 如果state中没有指定参数或者为空，则设置为正确的空值类型
		state.Params = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"param_name":    types.StringType,
				"current_value": types.StringType,
			},
		})
	}
	return
}

// getIdByQueryList 从远端查询参数模板id
func (c *ctyunRedisParamTemplate) getIdByQueryList(ctx context.Context, state CtyunRedisParamTemplateConfig) (id string, err error) {
	// 组装请求体
	params := &ctgdcs2.Dcs2DescribeRedisTemplateRequest{
		RegionId: state.RegionId.ValueString(),
		RawType:  "custom",
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

	// 遍历返回的模板列表，查找匹配的模板并设置ID
	for _, template := range resp.ReturnObj.List {
		// 根据模板名称匹配来设置ID
		if template.Name == state.Name.ValueString() {
			id = template.Id
			return
		}
	}
	err = common.InvalidReturnObjError
	return
}
