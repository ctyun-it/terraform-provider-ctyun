package cloud_assistant

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &CtyunCloudAssistantRunCommand{}
	_ resource.ResourceWithConfigure = &CtyunCloudAssistantRunCommand{}
)

type CtyunCloudAssistantRunCommand struct {
	meta                  *common.CtyunMetadata
	cloudAssistantService *business.CloudAssistantService
	name                  string
}

func NewCtyunCloudAssistantRunCommand() resource.Resource {
	return &CtyunCloudAssistantRunCommand{}
}

func (c *CtyunCloudAssistantRunCommand) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.cloudAssistantService = business.NewCloudAssistantService(c.meta)
}

func (c *CtyunCloudAssistantRunCommand) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_cloud_assistant_run_command"
	c.name = response.TypeName
}

func (c *CtyunCloudAssistantRunCommand) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("在云主机中执行一段Shell、PowerShell、Bat或Python类型的脚本命令", "云助手", "https://www.ctyun.cn/document/10026730/10764038"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源标识，invokedID。可以通过此id查询到执行结果",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池 ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"instance_ids": schema.StringAttribute{
				Required:    true,
				Description: "待执行命令的弹性云主机ID列表，使用英文逗号分割，一次最多不超过100台",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 100),
				},
			},
			"command_name": schema.StringAttribute{
				Optional:    true,
				Description: "命令名称，长度不超过128个字符",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 128),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "命令描述，长度不超过512个字符",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(512),
				},
			},
			"command_type": schema.StringAttribute{
				Optional:    true,
				Description: "命令类型：shell、bat、powershell、python",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("shell", "bat", "powershell", "python"),
				},
			},
			"command_content": schema.StringAttribute{
				Optional:    true,
				Description: "加密后的命令内容，长度不可超过24KB",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"working_directory": schema.StringAttribute{
				Optional:    true,
				Description: "命令在实例中运行目录，Linux默认/tmp，Windows默认C:\\Windows\\System32",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"timeout": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "命令超时时间，默认60秒",
				Default:     int32default.StaticInt32(60),
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"save_command": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否保存命令，默认false",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"enabled_parameter": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否启用自定义参数，默认false",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"parameter": schema.MapAttribute{
				Optional:    true,
				Description: "自定义参数，key和value均只支持string",
				ElementType: types.StringType,
			},
			"results": schema.ListNestedAttribute{
				Computed:    true,
				Description: "各实例的命令执行结果",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "云主机 ID",
						},
						"invocation_status": schema.StringAttribute{
							Computed:    true,
							Description: "命令运行状态：Pending/Running/Success/Failed",
						},
						"output": schema.StringAttribute{
							Computed:    true,
							Description: "命令执行输出信息",
						},
						"exit_code": schema.Int32Attribute{
							Computed:    true,
							Description: "命令退出码",
						},
						"error_info": schema.StringAttribute{
							Computed:    true,
							Description: "命令执行失败原因",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "命令执行创建时间",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "命令执行完成时间",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunCloudAssistantRunCommand) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunCloudAssistantRunCommandConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.runCommand(ctx, &plan)
	if err != nil {
		return
	}

	err = c.waitForCompletion(ctx, &plan)
	if err != nil {
		return
	}

	err = c.getAndMergeResults(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *CtyunCloudAssistantRunCommand) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunCloudAssistantRunCommandConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.getAndMergeResults(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunCloudAssistantRunCommand) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunCloudAssistantRunCommand) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	return
}

func (c *CtyunCloudAssistantRunCommand) runCommand(ctx context.Context, config *CtyunCloudAssistantRunCommandConfig) error {
	req := &ctecs.CtecsCloudAssistantRunCommandRequest{
		RegionID:         config.RegionID.ValueString(),
		InstanceIDs:      config.InstanceIDs.ValueString(),
		Timeout:          config.Timeout.ValueInt32(),
		SaveCommand:      config.SaveCommand.ValueBoolPointer(),
		EnabledParameter: config.EnabledParameter.ValueBoolPointer(),
		CommandType:      config.CommandType.ValueString(),
		CommandContent:   utils.GenerateCryptoContent(config.CommandContent.ValueString()),
	}

	if !config.CommandName.IsNull() && !config.CommandName.IsUnknown() {
		req.CommandName = config.CommandName.ValueString()
	}
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		req.Description = config.Description.ValueString()
	}

	if !config.WorkingDirectory.IsNull() && !config.WorkingDirectory.IsUnknown() {
		req.WorkingDirectory = config.WorkingDirectory.ValueString()
	}

	if config.EnabledParameter.ValueBool() && !config.Parameter.IsNull() && !config.Parameter.IsUnknown() {
		var inputParameter map[string]string
		diags := config.Parameter.ElementsAs(ctx, &inputParameter, false)
		if diags.HasError() {
			return fmt.Errorf(diags[0].Detail())
		}
		commandParams := make([]*ctecs.CtecsCloudAssistantRunCommandDefaultParameterRequest, 0)
		for key, value := range inputParameter {
			commandParams = append(commandParams, &ctecs.CtecsCloudAssistantRunCommandDefaultParameterRequest{
				Key:   key,
				Value: value,
			})
		}
		req.DefaultParameter = commandParams
	}

	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsCloudAssistantRunCommandApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return err
	} else if resp == nil {
		return fmt.Errorf("执行云助手命令失败，接口返回nil，请联系研发确认问题原因！")
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}

	config.ID = types.StringValue(resp.ReturnObj.InvokedID)

	return nil
}

func (c *CtyunCloudAssistantRunCommand) waitForCompletion(ctx context.Context, config *CtyunCloudAssistantRunCommandConfig) (err error) {
	retryer, _ := business.NewRetryer(time.Second*10, 60)
	retryer.Start(
		func(currentTime int) bool {
			results, err2 := c.getInvocationResults(ctx, config)
			if err2 != nil {
				err = err2
				return false
			}
			for _, result := range results {
				if result.InvocationStatus != "Success" && result.InvocationStatus != "Failed" {
					return true
				}
			}
			return false
		})
	return err
}

func (c *CtyunCloudAssistantRunCommand) getInvocationResults(ctx context.Context, config *CtyunCloudAssistantRunCommandConfig) ([]*ctecs.CtecsCloudAssistantDescribeInvocationResultsReturnObjResultsResponse, error) {
	req := &ctecs.CtecsCloudAssistantDescribeInvocationResultsRequest{
		RegionID:  config.RegionID.ValueString(),
		InvokedID: config.ID.ValueString(),
		PageNo:    1,
		PageSize:  100,
	}

	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsCloudAssistantDescribeInvocationResultsApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return nil, err
	} else if resp == nil || resp.StatusCode != common.NormalStatusCode {
		return nil, fmt.Errorf("获取脚本执行状态失败，API返回异常")
	} else if resp.ReturnObj == nil || resp.ReturnObj.Results == nil {
		return nil, common.InvalidReturnObjError
	}

	return resp.ReturnObj.Results, nil
}

func (c *CtyunCloudAssistantRunCommand) getAndMergeResults(ctx context.Context, config *CtyunCloudAssistantRunCommandConfig) error {
	results, err := c.getInvocationResults(ctx, config)
	if err != nil {
		return err
	}

	var resultModels []RunCommandResultModel
	for _, result := range results {
		resultModels = append(resultModels, RunCommandResultModel{
			InstanceID:       types.StringValue(result.InstanceID),
			InvocationStatus: types.StringValue(result.InvocationStatus),
			Output:           types.StringValue(result.Output),
			ExitCode:         types.Int32Value(result.ExitCode),
			ErrorInfo:        types.StringValue(result.ErrorInfo),
			CreateTime:       types.StringValue(result.CreateTime),
			UpdateTime:       types.StringValue(result.UpdateTime),
		})
	}

	resultList, diags := types.ListValueFrom(ctx, utils.StructToTFObjectTypes(RunCommandResultModel{}), resultModels)
	if diags.HasError() {
		return errors.New(diags[0].Detail())
	}
	config.Results = resultList
	return nil
}

type CtyunCloudAssistantRunCommandConfig struct {
	ID               types.String `tfsdk:"id"`
	RegionID         types.String `tfsdk:"region_id"`
	InstanceIDs      types.String `tfsdk:"instance_ids"`
	CommandName      types.String `tfsdk:"command_name"`
	Description      types.String `tfsdk:"description"`
	CommandType      types.String `tfsdk:"command_type"`
	CommandContent   types.String `tfsdk:"command_content"`
	WorkingDirectory types.String `tfsdk:"working_directory"`
	Timeout          types.Int32  `tfsdk:"timeout"`
	SaveCommand      types.Bool   `tfsdk:"save_command"`
	EnabledParameter types.Bool   `tfsdk:"enabled_parameter"`
	Parameter        types.Map    `tfsdk:"parameter"`
	Results          types.List   `tfsdk:"results"`
}

type RunCommandParamModel struct {
	Key         types.String `tfsdk:"key"`
	Description types.String `tfsdk:"description"`
	Value       types.String `tfsdk:"value"`
}

type RunCommandResultModel struct {
	InstanceID       types.String `tfsdk:"instance_id"`
	InvocationStatus types.String `tfsdk:"invocation_status"`
	Output           types.String `tfsdk:"output"`
	ExitCode         types.Int32  `tfsdk:"exit_code"`
	ErrorInfo        types.String `tfsdk:"error_info"`
	CreateTime       types.String `tfsdk:"create_time"`
	UpdateTime       types.String `tfsdk:"update_time"`
}
