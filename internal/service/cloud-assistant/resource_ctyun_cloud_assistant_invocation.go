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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &CtyunCloudAssistantInvocation{}
	_ resource.ResourceWithConfigure = &CtyunCloudAssistantInvocation{}
)

type CtyunCloudAssistantInvocation struct {
	meta                  *common.CtyunMetadata
	cloudAssistantService *business.CloudAssistantService

	name string
}

func NewCtyunCloudAssistantInvocation() resource.Resource {
	return &CtyunCloudAssistantInvocation{}
}

func (c *CtyunCloudAssistantInvocation) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.cloudAssistantService = business.NewCloudAssistantService(c.meta)
}

func (c *CtyunCloudAssistantInvocation) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_cloud_assistant_invocation"
	c.name = response.TypeName
}

func (c *CtyunCloudAssistantInvocation) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("触发云助手命令执行并查询执行结果", "云助手", "https://www.ctyun.cn/document/10026730/10764038"),
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
				Description: "待执行命令的弹性云主机ID列表，使用英文逗号分割",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"command_id": schema.StringAttribute{
				Required:    true,
				Description: "命令 ID（使用已有命令时传入）。若执行未保存的命令可以使用ctyun__cloud_assistant_run_command",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"working_directory": schema.StringAttribute{
				Optional:    true,
				Description: "命令运行目录",
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
			"parameter": schema.MapAttribute{
				Optional:    true,
				Description: "自定义参数的默认取值",
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
						"invocation_record_status": schema.StringAttribute{
							Computed:    true,
							Description: "单个命令执行任务的总状态，取值范围:Pending：未执行，当有云主机中命令状态为Pending，则总的执行状态为未执行；\nRunning：运行中，有云主机中命令进程为运行中，则总的执行状态为运行中；\nFinished：已完成。所有云主机命令进程全部完成执行；\nFailed：执行失败，有云主机中命令进程为执行失败，则总的状态为Failed。",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunCloudAssistantInvocation) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunCloudAssistantInvocationConfig
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

func (c *CtyunCloudAssistantInvocation) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunCloudAssistantInvocationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunCloudAssistantInvocation) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunCloudAssistantInvocation) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	return
}

func (c *CtyunCloudAssistantInvocation) getInvocationResult(ctx context.Context, config *CtyunCloudAssistantInvocationConfig) ([]*ctecs.CtecsCloudAssistantDescribeInvocationResultsReturnObjResultsResponse, error) {
	req := &ctecs.CtecsCloudAssistantDescribeInvocationResultsRequest{
		RegionID:  config.RegionID.ValueString(),
		InvokedID: config.ID.ValueString(),
		CommandID: config.CommandID.ValueString(),
		PageNo:    1,
		PageSize:  100,
	}

	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsCloudAssistantDescribeInvocationResultsApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取脚本执行状态失败(id=%s, command_id =%s)，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString(), config.CommandID.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	results := resp.ReturnObj.Results
	if results != nil {
		return results, nil
	}

	return nil, common.InvalidReturnObjError
}

func (c *CtyunCloudAssistantInvocation) waitForCompletion(ctx context.Context, config *CtyunCloudAssistantInvocationConfig) (err error) {
	retryer, _ := business.NewRetryer(time.Second*10, 60)
	retryer.Start(
		func(currentTime int) bool {
			results, err2 := c.getInvocationResult(ctx, config)
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

func (c *CtyunCloudAssistantInvocation) create(ctx context.Context, plan *CtyunCloudAssistantInvocationConfig) error {
	var req = &ctecs.CtecsInvokeCommandRequest{
		RegionID:    plan.RegionID.ValueString(),
		InstanceIDs: plan.InstanceIDs.ValueString(),
		CommandID:   plan.CommandID.ValueString(),
		Timeout:     plan.Timeout.ValueInt32(),
	}

	if !plan.WorkingDirectory.IsNull() && !plan.WorkingDirectory.IsUnknown() {
		req.WorkingDirectory = plan.WorkingDirectory.ValueString()
	}

	if !plan.Parameter.IsNull() && !plan.Parameter.IsUnknown() {
		var inputParameter map[string]string
		diags := plan.Parameter.ElementsAs(ctx, &inputParameter, false)
		if diags.HasError() {
			err := fmt.Errorf(diags[0].Detail())
			return err
		}
		req.Parameter = &inputParameter
	}

	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsInvokeCommandApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("触发云助手命令失败，command_id=%s，接口返回nil,请联系研发确认问题原因！", plan.CommandID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	plan.ID = types.StringValue(resp.ReturnObj.InvokedID)
	err = c.waitForCompletion(ctx, plan)
	if err != nil {
		return err
	}
	return nil
}

func (c *CtyunCloudAssistantInvocation) getAndMerge(ctx context.Context, config *CtyunCloudAssistantInvocationConfig) error {
	results, err := c.getInvocationResult(ctx, config)
	if err != nil {
		return err
	}
	var resultModels []InvocationResultModel
	for _, result := range results {
		resultModels = append(resultModels, InvocationResultModel{
			InstanceID:       types.StringValue(result.InstanceID),
			InvocationStatus: types.StringValue(result.InvocationStatus),
			Output:           types.StringValue(result.Output),
			ExitCode:         types.Int32Value(result.ExitCode),
			ErrorInfo:        types.StringValue(result.ErrorInfo),
			CreateTime:       types.StringValue(result.CreateTime),
			UpdateTime:       types.StringValue(result.UpdateTime),
		})
	}
	resultList, diags := types.ListValueFrom(ctx, utils.StructToTFObjectTypes(InvocationResultModel{}), resultModels)
	if diags.HasError() {
		err = errors.New(diags[0].Detail())
		return err
	}
	config.Results = resultList
	return nil
}

type CtyunCloudAssistantInvocationConfig struct {
	ID               types.String `tfsdk:"id"`
	RegionID         types.String `tfsdk:"region_id"`
	InstanceIDs      types.String `tfsdk:"instance_ids"`
	CommandID        types.String `tfsdk:"command_id"`
	WorkingDirectory types.String `tfsdk:"working_directory"`
	Timeout          types.Int32  `tfsdk:"timeout"`
	Parameter        types.Map    `tfsdk:"parameter"`
	Results          types.List   `tfsdk:"results"`
}

type InvocationResultModel struct {
	InstanceID             types.String `tfsdk:"instance_id"`
	InvocationStatus       types.String `tfsdk:"invocation_status"`
	Output                 types.String `tfsdk:"output"`
	ExitCode               types.Int32  `tfsdk:"exit_code"`
	ErrorInfo              types.String `tfsdk:"error_info"`
	InvocationRecordStatus types.String `tfsdk:"invocation_record_status"`
	CreateTime             types.String `tfsdk:"create_time"`
	UpdateTime             types.String `tfsdk:"update_time"`
}
