package cloud_assistant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CtyunCloudAssistantCommand{}
	_ resource.ResourceWithConfigure   = &CtyunCloudAssistantCommand{}
	_ resource.ResourceWithImportState = &CtyunCloudAssistantCommand{}
)

type CtyunCloudAssistantCommand struct {
	meta *common.CtyunMetadata
	name string
}

func NewCtyunCloudAssistantCommand() resource.Resource {
	return &CtyunCloudAssistantCommand{}
}

func (c *CtyunCloudAssistantCommand) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunCloudAssistantCommand) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_cloud_assistant_command"
	c.name = response.TypeName
}

func (c *CtyunCloudAssistantCommand) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理云助手命令", "云助手", "https://www.ctyun.cn/document/10026730/10764038"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "命令 ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池 ID，如果不填则默认使用 provider ctyun 中的 region_id 或环境变量中的 CTYUN_REGION_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"command_name": schema.StringAttribute{
				Required:    true,
				Description: "命令名称，长度不超过128个字符",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 128),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "命令描述，长度不超过512个字符",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(512),
				},
			},
			"command_type": schema.StringAttribute{
				Required:    true,
				Description: "命令类型，可选值：Shell, Bat, PowerShell, Python",
				Validators: []validator.String{
					stringvalidator.OneOf("shell", "bat", "powershell", "python"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"command_content": schema.StringAttribute{
				Required:    true,
				Description: "加密后的命令内容，base64编码长度不可超过24KB",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(32768),
				},
			},
			"working_directory": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "命令在实例中运行目录，Linux默认/root，Windows默认C:\\Windows\\System32",
			},
			"timeout": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "命令超时时间，默认60秒",
				Default:     int32default.StaticInt32(60),
			},
			"enabled_parameter": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否启用自定义参数",
				Default:     booldefault.StaticBool(false),
			},
			"parameter": schema.MapAttribute{
				Optional:    true,
				Description: "命令参数，最多支持20个参数",
				ElementType: types.StringType,
				Validators: []validator.Map{
					mapvalidator.SizeAtMost(20),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间",
			},
		},
	}
}

func (c *CtyunCloudAssistantCommand) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunCloudAssistantCommandConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.create(ctx, &plan)
	if err != nil {
		return
	}

	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *CtyunCloudAssistantCommand) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunCloudAssistantCommandConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunCloudAssistantCommand) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取tf文件中配置
	var plan CtyunCloudAssistantCommandConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunCloudAssistantCommandConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, state, plan)
	if err != nil {
		return
	}

	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *CtyunCloudAssistantCommand) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunCloudAssistantCommandConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	req := &ctecs.CtecsDeleteCommandRequest{
		RegionID:  state.RegionID.ValueString(),
		CommandID: state.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsDeleteCommandApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("删除云助手命令失败(id=%s)，接口返回nil，请联系研发确认问题原因！", state.ID.ValueString())
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

func (c *CtyunCloudAssistantCommand) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入云助手命令: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import %s.[导入配置名称] <command_id>,<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()

	var config CtyunCloudAssistantCommandConfig
	var ID, regionId string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		ID = request.ID
	} else {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &ID, &regionId)
		if err != nil {
			return
		}
	}

	if ID == "" {
		err = fmt.Errorf("id不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}
	config.ID = types.StringValue(ID)
	config.RegionID = types.StringValue(regionId)

	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunCloudAssistantCommand) getAndMerge(ctx context.Context, config *CtyunCloudAssistantCommandConfig) error {

	req := &ctecs.CtecsGetCommandRequest{
		RegionID:  config.RegionID.ValueString(),
		CommandID: config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsGetCommandApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("获取云助手命令(id=%s)接口失败，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		if strings.Contains(resp.Message, "not exist") || strings.Contains(resp.Description, "不存在") {
			return common.ResourceNotExistError
		}
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError

	}

	obj := resp.ReturnObj
	config.CommandName = types.StringValue(obj.CommandName)
	config.CommandType = types.StringValue(strings.ToLower(obj.CommandType))
	config.CommandContent = types.StringValue(obj.CommandContent)

	config.Description = types.StringValue(obj.Description)
	config.WorkingDirectory = types.StringValue(obj.WorkingDirectory)
	config.Timeout = types.Int32Value(obj.Timeout)
	if obj.EnabledParameter != nil {
		config.EnabledParameter = types.BoolValue(*obj.EnabledParameter)
	} else {
		config.EnabledParameter = types.BoolValue(false)
	}
	config.CreateTime = types.StringValue(obj.CreateTime)
	config.UpdateTime = types.StringValue(obj.UpdateTime)
	// 处理参数， 因为接口返回字符串，需要先进行序列化操作
	paramStr := obj.DefaultParameter
	if paramStr != "" {
		params, err2 := c.parseParameter(paramStr)
		if err2 != nil {
			return err2
		}
		if params != nil && len(params) > 0 {
			parameter, diags := types.MapValueFrom(ctx, types.StringType, params)
			if diags.HasError() {
				err = errors.New(diags[0].Detail())
				return err
			}
			config.Parameter = parameter
		}
	} else {
		config.Parameter = types.MapNull(types.StringType)
	}
	return nil
}

func (c *CtyunCloudAssistantCommand) create(ctx context.Context, plan *CtyunCloudAssistantCommandConfig) error {
	req := &ctecs.CtecsCreateCommandRequest{
		RegionID:         plan.RegionID.ValueString(),
		CommandName:      plan.CommandName.ValueString(),
		Description:      plan.Description.ValueString(),
		CommandType:      plan.CommandType.ValueString(),
		CommandContent:   utils.GenerateCryptoContent(plan.CommandContent.ValueString()),
		WorkingDirectory: plan.WorkingDirectory.ValueString(),
		Timeout:          plan.Timeout.ValueInt32(),
		EnabledParameter: plan.EnabledParameter.ValueBoolPointer(),
	}
	var inputParameter map[string]string
	diags := plan.Parameter.ElementsAs(ctx, &inputParameter, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}

	if plan.EnabledParameter.ValueBool() && len(inputParameter) > 0 {
		params := make([]*ctecs.CtecsCreateCommandDefaultParameterRequest, 0)
		for key, value := range inputParameter {
			item := &ctecs.CtecsCreateCommandDefaultParameterRequest{
				Key:   key,
				Value: value,
			}
			params = append(params, item)
		}
		req.DefaultParameter = params
	}

	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsCreateCommandApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建云助手命令失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}

	if resp.ReturnObj != nil {
		plan.ID = types.StringValue(resp.ReturnObj.CommandID)
	}
	return nil
}

func (c *CtyunCloudAssistantCommand) update(ctx context.Context, state CtyunCloudAssistantCommandConfig, plan CtyunCloudAssistantCommandConfig) error {
	req := &ctecs.CtecsModifyCommandRequest{
		RegionID:         plan.RegionID.ValueString(),
		CommandID:        state.ID.ValueString(),
		CommandName:      plan.CommandName.ValueString(),
		Description:      plan.Description.ValueString(),
		CommandType:      plan.CommandType.ValueString(),
		CommandContent:   utils.GenerateCryptoContent(plan.CommandContent.ValueString()),
		WorkingDirectory: plan.WorkingDirectory.ValueString(),
		Timeout:          plan.Timeout.ValueInt32(),
	}

	enabledParam := plan.EnabledParameter.ValueBool()
	req.EnabledParameter = &enabledParam

	var inputParameter map[string]string
	diags := plan.Parameter.ElementsAs(ctx, &inputParameter, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}

	if enabledParam && len(inputParameter) > 0 {
		params := make([]*ctecs.CtecsModifyCommandDefaultParameterRequest, 0)
		for key, value := range inputParameter {
			item := &ctecs.CtecsModifyCommandDefaultParameterRequest{
				Key:   key,
				Value: value,
			}
			params = append(params, item)
		}
		req.DefaultParameter = params
	}

	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsModifyCommandApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("修改云助手失败（id=%s），接口返回nil，请联系研发确认问题原因！", plan.ID.ValueString())
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	return nil
}

func (c *CtyunCloudAssistantCommand) parseParameter(params string) (map[string]string, error) {
	// 2. 定义切片（list）
	var itemList []DefaultParameterModel

	// 3. 反序列化（字符串 → Go 对象）
	err := json.Unmarshal([]byte(params), &itemList)
	if err != nil {
		fmt.Println("解析失败:", err)
		return nil, err
	}
	paramMap := make(map[string]string)
	for _, item := range itemList {
		paramMap[item.Key] = item.Value
	}
	return paramMap, nil
}

func splitStr(s string, sep string, maxParts int) []string {
	result := make([]string, 0)
	current := ""
	count := 0
	for _, char := range s {
		if string(char) == sep && count < maxParts-1 {
			result = append(result, current)
			current = ""
			count++
		} else {
			current += string(char)
		}
	}
	result = append(result, current)
	return result
}

type CtyunCloudAssistantCommandConfig struct {
	ID               types.String `tfsdk:"id"`
	RegionID         types.String `tfsdk:"region_id"`
	CommandName      types.String `tfsdk:"command_name"`
	Description      types.String `tfsdk:"description"`
	CommandType      types.String `tfsdk:"command_type"`
	CommandContent   types.String `tfsdk:"command_content"`
	WorkingDirectory types.String `tfsdk:"working_directory"`
	Timeout          types.Int32  `tfsdk:"timeout"`
	EnabledParameter types.Bool   `tfsdk:"enabled_parameter"`
	Parameter        types.Map    `tfsdk:"parameter"`
	CreateTime       types.String `tfsdk:"create_time"`
	UpdateTime       types.String `tfsdk:"update_time"`
}

type DefaultParameterModel struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Value       string `json:"value"`
}
