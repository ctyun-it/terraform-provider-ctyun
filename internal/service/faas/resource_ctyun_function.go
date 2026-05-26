package faas

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/cf"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &CtyunFunction{}
	_ resource.ResourceWithConfigure   = &CtyunFunction{}
	_ resource.ResourceWithImportState = &CtyunFunction{}
)

func NewCtyunFunction() resource.Resource {
	return &CtyunFunction{}
}

type CtyunFunction struct {
	meta *common.CtyunMetadata
	name string
}

type CtyunFunctionConfig struct {
	ID                         types.String  `tfsdk:"id"`
	Name                       types.String  `tfsdk:"name"`
	CreateType                 types.Int32   `tfsdk:"create_type"`
	RuntimeRuntime             types.String  `tfsdk:"runtime_runtime"`
	RuntimeHandleType          types.String  `tfsdk:"runtime_handle_type"`
	RuntimeHandler             types.String  `tfsdk:"runtime_handler"`
	RuntimeExecuteTimeout      types.Int32   `tfsdk:"runtime_execute_timeout"`
	RuntimeInstanceConcurrency types.Int32   `tfsdk:"runtime_instance_concurrency"`
	ContainerTimeZone          types.String  `tfsdk:"container_time_zone"`
	ContainerDiskSize          types.Int32   `tfsdk:"container_disk_size"`
	ContainerMemorySize        types.Int32   `tfsdk:"container_memory_size"`
	ContainerCpu               types.Float64 `tfsdk:"container_cpu"`
	ContainerListenPort        types.Int32   `tfsdk:"container_listen_port"`
	ContainerMaxScale          types.Int32   `tfsdk:"container_max_scale"`
	ContainerFastStart         types.Int32   `tfsdk:"container_fast_start"`
	ContainerImage             types.String  `tfsdk:"container_image"`
	ContainerRunCommand        types.String  `tfsdk:"container_run_command"`
	Description                types.String  `tfsdk:"description"`
	Environment                types.Map     `tfsdk:"environment"`
	CodeContent                types.String  `tfsdk:"code_content"`
	CodeBucket                 types.String  `tfsdk:"code_bucket"`
	CodeKey                    types.String  `tfsdk:"code_key"`
	CodeVersion                types.String  `tfsdk:"code_version"`
	Role                       types.String  `tfsdk:"role"`
	RegionID                   types.String  `tfsdk:"region_id"`
	FunctionID                 types.String  `tfsdk:"function_id"`
	CreatedAt                  types.String  `tfsdk:"created_at"`
	UpdatedAt                  types.String  `tfsdk:"updated_at"`
	Status                     types.String  `tfsdk:"status"`
}

func (c *CtyunFunction) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function"
	c.name = resp.TypeName
}

func (c *CtyunFunction) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理函数", "函数工作流（FunctionGraph）", "https://www.ctyun.cn/document/10355289"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "函数名称。长度限制：2~64 个字符，支持数字、字母、下划线、连字符，以字母开头",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 64),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]*$"), "函数名称必须以字母开头，只能包含字母、数字、下划线、连字符"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"create_type": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "创建函数的类型 1:标准运行时 2:自定义运行时 3:自定义镜像",
				Default:     int32default.StaticInt32(1),
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"runtime_runtime": schema.StringAttribute{
				Required:    true,
				Description: "运行时类型，如 python3.9 支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
					//TODO 等函数计算 文档添加可用的枚举值
				},
			},
			"runtime_handle_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "请求处理程序类型（标准运行时必填），如 http/event",
				Default:     stringdefault.StaticString("event"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"runtime_handler": schema.StringAttribute{
				Required:    true,
				Description: "函数执行的入口（标准运行时必填）支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 128),
				},
			},
			"runtime_execute_timeout": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "执行超时时间，默认 10 秒 支持更新",
				Default:     int32default.StaticInt32(10),
			},
			"runtime_instance_concurrency": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "实例最大并发度 支持更新",
				Default:     int32default.StaticInt32(1),
			},
			"container_time_zone": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "时区，默认 UTC 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"container_disk_size": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "磁盘规格 (Mb)，默认 512，范围 512,10240,20480 支持更新",
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
					int32validator.OneOf(512, 10240, 20480),
				},
			},
			"container_memory_size": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "内存规格 (Mb)，默认 128，范围 128~32768 支持更新",
				Validators: []validator.Int32{
					int32validator.AtLeast(128),
					int32validator.AtMost(32768),
				},
			},
			"container_cpu": schema.Float64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "CPU 规格 (vCPU)，默认 0.1，范围 0.05~12.0 支持更新",
				Validators: []validator.Float64{
					float64validator.AtLeast(0.05),
					float64validator.AtMost(12.0),
				},
			},
			"container_listen_port": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "监听端口，默认 8080 支持更新",
				Default:     int32default.StaticInt32(8080),
			},
			"container_max_scale": schema.Int32Attribute{
				Optional:    true,
				Description: "并发实例数上限",
			},
			"container_fast_start": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "启动加速，默认为 0，1 表示使用启动加速 ",
				Default:     int32default.StaticInt32(0),
				Validators:  []validator.Int32{
					//	TODO 只有当 runtime_handle_type 是java 类型才支持该字段 ,暂时先不管该字段 默认为0
				},
			},
			"container_run_command": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "函数服务启动命令 支持更新",
			},
			"container_image": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "函数服务镜像 支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "函数描述，长度最大为 512，支持更新",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(512),
				},
			},
			"environment": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "函数环境变量，支持更新",
			},
			"code_content": schema.StringAttribute{
				Optional:    true,
				Description: "内联代码内容或 Zip 文件的 Base64 编码",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"code_bucket": schema.StringAttribute{
				Optional:    true,
				Description: "OBS 桶名称",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"code_key": schema.StringAttribute{
				Optional:    true,
				Description: "OBS 对象键（文件路径）",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"code_version": schema.StringAttribute{
				Optional:    true,
				Description: "代码版本号",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Optional:    true,
				Description: "角色 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池 ID，如果不填则默认使用 provider ctyun 中的 region_id 或环境变量中的 CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"function_id": schema.StringAttribute{
				Computed:    true,
				Description: "函数唯一标识",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "函数状态",
			},
		},
	}
}

func (c *CtyunFunction) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunFunction) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunFunctionConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
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

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	err = c.waitForDeploymentSuccess(ctx, &plan, 5*time.Minute)
	if err != nil {
		resp.Diagnostics.AddWarning("函数部署超时警告", err.Error())
	}

}

func (c *CtyunFunction) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunFunctionConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			err = nil
			resp.State.RemoveResource(ctx)
		}
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunFunction) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunFunctionConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state CtyunFunctionConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &plan, state)
	if err != nil {
		return
	}

	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	err = c.waitForDeploymentSuccess(ctx, &plan, 5*time.Minute)
	if err != nil {
		resp.Diagnostics.AddWarning("函数部署超时警告", err.Error())
	}
}

func (c *CtyunFunction) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunFunctionConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunFunction) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例：%s 失败：%s", c.name, req.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [name],<region_id>", c.name)
			resp.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunFunctionConfig
	var name, regionID string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(req.ID, common.ImportSeparator) == 0 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(req.ID, &name)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(req.ID, &name, &regionID)
		if err != nil {
			return
		}
	}

	if name == "" {
		err = fmt.Errorf("name不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}
	config.Name = types.StringValue(name)
	config.RegionID = types.StringValue(regionID)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (c *CtyunFunction) create(ctx context.Context, plan *CtyunFunctionConfig) (err error) {
	req := &cf.CfCreateFunctionV2024Request{
		RegionId:     plan.RegionID.ValueString(),
		FunctionName: plan.Name.ValueString(),
		CreateType:   plan.CreateType.ValueInt32(),
		Container: &cf.CfCreateFunctionV2024ContainerRequest{
			TimeZone:   plan.ContainerTimeZone.ValueString(),
			DiskSize:   plan.ContainerDiskSize.ValueInt32(),
			MemorySize: plan.ContainerMemorySize.ValueInt32(),
			Cpu:        plan.ContainerCpu.ValueFloat64(),
			ListenPort: plan.ContainerListenPort.ValueInt32(),
		},
		Runtime: &cf.CfCreateFunctionV2024RuntimeRequest{
			Runtime:             plan.RuntimeRuntime.ValueStringPointer(),
			HandleType:          plan.RuntimeHandleType.ValueStringPointer(),
			ExecuteTimeout:      plan.RuntimeExecuteTimeout.ValueInt32Pointer(),
			Handler:             plan.RuntimeHandler.ValueStringPointer(),
			InstanceConcurrency: plan.RuntimeInstanceConcurrency.ValueInt32Pointer(),
		},
	}

	req.Container.MaxScale = plan.ContainerMaxScale.ValueInt32Pointer()
	req.Container.FastStart = plan.ContainerFastStart.ValueInt32Pointer()
	req.Container.RunCommand = plan.ContainerRunCommand.ValueStringPointer()
	req.Description = plan.Description.ValueStringPointer()

	envVars, err := utils.TypesMapToStringMap(ctx, plan.Environment)
	if err != nil {
		return fmt.Errorf("failed to convert environment variables: %w", err)
	}
	if envVars != nil {
		req.Container.EnvironmentVariables = envVars
	}

	req.Code = &cf.CfCreateFunctionV2024CodeRequest{
		OssBucketName: plan.CodeBucket.ValueStringPointer(),
		OssObjectName: plan.CodeKey.ValueStringPointer(),
		ZipFile:       plan.CodeContent.ValueStringPointer(),
	}
	req.Role = plan.Role.ValueStringPointer()

	resp, err := c.meta.Apis.SdkCtCfApis.CfCreateFunctionV2024Api.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if utils.SecInt32(resp.StatusCode) != 0 {
		err = fmt.Errorf("API return error. Message: %s Description:", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	plan.FunctionID = types.StringPointerValue(resp.ReturnObj.FunctionId)
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s", plan.RegionID.ValueString(), plan.Name.ValueString()))

	return
}

func (c *CtyunFunction) update(ctx context.Context, plan *CtyunFunctionConfig, state CtyunFunctionConfig) (err error) {

	req := &cf.CfUpdateFunctionRequest{
		FunctionName: plan.Name.ValueString(),
		RegionId:     plan.RegionID.ValueString(),
		UpdateType:   "all",
	}

	hasChanges := false

	if !plan.Description.Equal(state.Description) {
		req.Description = plan.Description.ValueStringPointer()
		hasChanges = true
	}

	containerChanged := plan.ContainerTimeZone != state.ContainerTimeZone ||
		plan.ContainerDiskSize != state.ContainerDiskSize ||
		plan.ContainerMemorySize != state.ContainerMemorySize ||
		plan.ContainerCpu != state.ContainerCpu

	environmentChanged := !plan.Environment.Equal(state.Environment)

	if containerChanged || environmentChanged {
		if req.Container == nil {
			req.Container = &cf.CfUpdateFunctionContainerRequest{}
		}

		req.Container.TimeZone = plan.ContainerTimeZone.ValueStringPointer()

		req.Container.DiskSize = plan.ContainerDiskSize.ValueInt32Pointer()

		req.Container.MemorySize = plan.ContainerMemorySize.ValueInt32Pointer()

		req.Container.Cpu = plan.ContainerCpu.ValueFloat64Pointer()
		req.Container.FastStart = plan.ContainerFastStart.ValueInt32Pointer()
		req.Container.ListenPort = plan.ContainerListenPort.ValueInt32()
		req.Container.RunCommand = plan.ContainerRunCommand.ValueStringPointer()
		req.Container.MaxScale = plan.ContainerMaxScale.ValueInt32Pointer()
		if plan.ContainerImage.ValueString() == "" {
			req.Container.Image = state.ContainerImage.ValueStringPointer()
		} else {
			req.Container.Image = plan.ContainerImage.ValueStringPointer()
		}

		if environmentChanged {
			envVars, err := utils.TypesMapToStringMap(ctx, plan.Environment)
			if err != nil {
				return fmt.Errorf("failed to convert environment variables: %w", err)
			}
			if envVars != nil {
				req.Container.EnvironmentVariables = envVars
			}
		}

		hasChanges = true
	}

	runtimeChanged := plan.RuntimeRuntime != state.RuntimeRuntime ||
		plan.RuntimeHandleType != state.RuntimeHandleType ||
		plan.RuntimeHandler != state.RuntimeHandler ||
		plan.RuntimeExecuteTimeout != state.RuntimeExecuteTimeout ||
		plan.RuntimeInstanceConcurrency != state.RuntimeInstanceConcurrency

	if runtimeChanged {
		if req.Runtime == nil {
			req.Runtime = &cf.CfUpdateFunctionRuntimeRequest{}
		}

		req.Runtime.Runtime = plan.RuntimeRuntime.ValueStringPointer()
		req.Runtime.HandleType = plan.RuntimeHandleType.ValueStringPointer()
		req.Runtime.Handler = plan.RuntimeHandler.ValueStringPointer()
		req.Runtime.ExecuteTimeout = plan.RuntimeExecuteTimeout.ValueInt32Pointer()
		req.Runtime.InstanceConcurrency = plan.RuntimeInstanceConcurrency.ValueInt32Pointer()
		hasChanges = true
	}

	if !hasChanges {
		return
	}

	resp, err := c.meta.Apis.SdkCtCfApis.CfUpdateFunctionApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if utils.SecInt32(resp.StatusCode) != 0 {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return
}

func (c *CtyunFunction) delete(ctx context.Context, state *CtyunFunctionConfig) (err error) {
	req := &cf.CfDeleteFunctionV2024Request{
		FunctionName: state.Name.ValueString(),
		RegionId:     state.RegionID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtCfApis.CfDeleteFunctionV2024Api.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if utils.SecInt32(resp.StatusCode) != 0 {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Message)
		return
	}

	return
}

func (c *CtyunFunction) getAndMerge(ctx context.Context, config *CtyunFunctionConfig) (err error) {
	functionName := config.Name.ValueString()
	if functionName == "" {
		return fmt.Errorf("function name is required")
	}

	req := &cf.CfGetFunctionV2024Request{
		FunctionName: functionName,
		RegionId:     config.RegionID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtCfApis.CfGetFunctionV2024Api.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if utils.SecInt32(resp.StatusCode) != 0 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	returnObj := resp.ReturnObj

	config.Name = types.StringPointerValue(returnObj.FunctionName)
	config.FunctionID = types.StringPointerValue(returnObj.FunctionId)
	config.CreateType = types.Int32PointerValue(returnObj.CreateType)
	config.Description = types.StringPointerValue(returnObj.Description)
	if returnObj.Container != nil {
		config.ContainerTimeZone = types.StringPointerValue(returnObj.Container.TimeZone)
		config.ContainerDiskSize = types.Int32PointerValue(returnObj.Container.DiskSize)
		config.ContainerMemorySize = types.Int32PointerValue(returnObj.Container.MemorySize)
		config.ContainerCpu = types.Float64PointerValue(returnObj.Container.Cpu)
		config.ContainerListenPort = types.Int32PointerValue(returnObj.Container.ListenPort)
		config.ContainerMaxScale = types.Int32PointerValue(returnObj.Container.MaxScale)

		config.ContainerFastStart = types.Int32PointerValue(returnObj.Container.FastStart)

		config.ContainerRunCommand = types.StringPointerValue(returnObj.Container.RunCommand)
		config.ContainerImage = types.StringValue(returnObj.Container.Image)
		if len(returnObj.Container.EnvironmentVariables) > 0 {
			envMap, diags := types.MapValueFrom(ctx, types.StringType, returnObj.Container.EnvironmentVariables)
			if diags.HasError() {
				return fmt.Errorf("failed to convert environment variables")
			}
			config.Environment = envMap
		} else {
			config.Environment = types.MapNull(types.StringType)
		}
	}

	if returnObj.Runtime != nil {
		config.RuntimeRuntime = types.StringPointerValue(returnObj.Runtime.Runtime)
		config.RuntimeHandleType = types.StringPointerValue(returnObj.Runtime.HandleType)
		config.RuntimeHandler = types.StringPointerValue(returnObj.Runtime.Handler)
		config.RuntimeExecuteTimeout = types.Int32PointerValue(returnObj.Runtime.ExecuteTimeout)
		config.RuntimeInstanceConcurrency = types.Int32PointerValue(returnObj.Runtime.InstanceConcurrency)
	}

	if returnObj.DeployInfo != nil {
		config.Status = types.StringPointerValue(returnObj.DeployInfo.Status)
		if returnObj.DeployInfo.TaskBegin != nil {
			config.CreatedAt = types.StringValue(strconv.FormatInt(int64(*returnObj.DeployInfo.TaskBegin), 10))
		}
		if returnObj.DeployInfo.TaskEnd != nil {
			config.UpdatedAt = types.StringValue(strconv.FormatInt(int64(*returnObj.DeployInfo.TaskEnd), 10))
		}
	}

	config.ID = types.StringValue(fmt.Sprintf("%s,%s", config.Name.ValueString(), config.RegionID.ValueString()))

	return
}

func (c *CtyunFunction) waitForDeploymentSuccess(ctx context.Context, config *CtyunFunctionConfig, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for function %s deployment to succeed...", config.Name.ValueString())

	interval := 10 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for deployment")
		default:
		}

		req := &cf.CfGetFunctionV2024Request{
			FunctionName: config.Name.ValueString(),
			RegionId:     config.RegionID.ValueString(),
		}

		resp, err := c.meta.Apis.SdkCtCfApis.CfGetFunctionV2024Api.Do(ctx, c.meta.SdkCredential, req)
		if err != nil {
			log.Printf("[WARN] Error querying function %s status: %s", config.Name.ValueString(), err.Error())
			time.Sleep(interval)
			continue
		}

		if resp.ReturnObj == nil || resp.ReturnObj.DeployInfo == nil {
			log.Printf("[DEBUG] Function %s deploy info not available yet", config.Name.ValueString())
			time.Sleep(interval)
			continue
		}

		status := resp.ReturnObj.DeployInfo.Status
		if status == nil {
			log.Printf("[DEBUG] Function %s status is nil, continuing to wait", config.Name.ValueString())
			time.Sleep(interval)
			continue
		}

		log.Printf("[DEBUG] Function %s current status: %s", config.Name.ValueString(), *status)

		if *status == "success" {
			log.Printf("[INFO] Function %s deployed successfully", config.Name.ValueString())
			return nil
		}

		if *status == "failed" {
			errMsg := ""
			if resp.ReturnObj.DeployInfo.ErrMsg != nil {
				errMsg = *resp.ReturnObj.DeployInfo.ErrMsg
			}
			return fmt.Errorf("function deployment failed with status: %s, error: %s", *status, errMsg)
		}

		time.Sleep(interval)
	}

	return fmt.Errorf("timeout waiting for function %s deployment after %v, current status may still be: %s",
		config.Name.ValueString(), timeout, "deploying")
}
