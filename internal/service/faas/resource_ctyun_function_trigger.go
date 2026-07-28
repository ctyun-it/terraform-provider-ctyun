package faas

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/cf"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CtyunFunctionTrigger{}
	_ resource.ResourceWithConfigure   = &CtyunFunctionTrigger{}
	_ resource.ResourceWithImportState = &CtyunFunctionTrigger{}
)

func NewCtyunFunctionTrigger() resource.Resource {
	return &CtyunFunctionTrigger{}
}

type CtyunFunctionTrigger struct {
	meta *common.CtyunMetadata
	name string
}

type CtyunFunctionTriggerConfig struct {
	ID           types.String `tfsdk:"id"`
	FunctionName types.String `tfsdk:"function_name"`
	TriggerName  types.String `tfsdk:"trigger_name"`
	TriggerType  types.String `tfsdk:"trigger_type"`
	EventData    types.String `tfsdk:"event_data"`
	Version      types.String `tfsdk:"version"`
	Enable       types.Bool   `tfsdk:"enable"`
	RegionID     types.String `tfsdk:"region_id"`
	Status       types.Int32  `tfsdk:"status"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

func (c *CtyunFunctionTrigger) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function_trigger"
	c.name = resp.TypeName
}

func (c *CtyunFunctionTrigger) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理函数触发器", "函数计算（FaaS）", "https://www.ctyun.cn/document/10006234/10528260"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"function_name": schema.StringAttribute{
				Required:    true,
				Description: "函数名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 64),
				},
			},
			"trigger_name": schema.StringAttribute{
				Required:    true,
				Description: "触发器名称。只能包含字母、数字和中划线。只能字母开头，字母数字结尾。长度在 3-63 之间",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(3, 63),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$"), "触发器名称必须以字母开头，字母或数字结尾，只能包含字母、数字和中划线"),
				},
			},
			"trigger_type": schema.StringAttribute{
				Required:    true,
				Description: "触发器类型。schedule: 定时触发器，http: Http 触发器",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("schedule", "http"),
				},
			},
			"event_data": schema.StringAttribute{
				Required:    true,
				Description: "触发器事件配置，JSON 格式字符串。不同触发器类型配置不同。 支持更新",
				Validators: []validator.String{
					validator2.StringIsJson(),
				},
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "版本或别名，版本包括 1,2,...这样的普通版本，和特殊版本 LATEST 支持更新",
				Default:     stringdefault.StaticString("LATEST"),
			},
			"enable": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否启用触发器。true：启用，false：禁用 支持更新",
				Default:     booldefault.StaticBool(true),
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
			"status": schema.Int32Attribute{
				Computed:    true,
				Description: "触发器状态。1：启用；2：禁用；3：系统禁用",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间",
			},
		},
	}
}

func (c *CtyunFunctionTrigger) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunFunctionTrigger) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunFunctionTriggerConfig
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
}

func (c *CtyunFunctionTrigger) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunFunctionTriggerConfig
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

func (c *CtyunFunctionTrigger) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunFunctionTriggerConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state CtyunFunctionTriggerConfig
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
}

func (c *CtyunFunctionTrigger) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunFunctionTriggerConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunFunctionTrigger) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例：%s 失败：%s", c.name, req.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [name], [functionName],<region_id>", c.name)
			resp.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunFunctionTriggerConfig
	var name, functionName, regionID string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(req.ID, common.ImportSeparator) == 1 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(req.ID, &name, &functionName)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(req.ID, &name, &functionName, &regionID)
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
	config.TriggerName = types.StringValue(name)
	config.FunctionName = types.StringValue(functionName)
	config.RegionID = types.StringValue(regionID)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (c *CtyunFunctionTrigger) create(ctx context.Context, plan *CtyunFunctionTriggerConfig) (err error) {
	req := &cf.CfCreateTriggerV2024Request{
		FunctionName:  plan.FunctionName.ValueString(),
		RegionId:      plan.RegionID.ValueString(),
		TriggerName:   plan.TriggerName.ValueString(),
		TriggerType:   plan.TriggerType.ValueString(),
		TriggerConfig: plan.EventData.ValueString(),
		Version:       plan.Version.ValueString(),
	}

	enable := plan.Enable.ValueBool()
	req.Enable = &enable

	resp, err := c.meta.Apis.SdkCtCfApis.CfCreateTriggerV2024Api.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp.StatusCode != common.FaasNormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", plan.FunctionName.ValueString(), plan.TriggerName.ValueString(), plan.RegionID.ValueString()))

	return
}

func (c *CtyunFunctionTrigger) update(ctx context.Context, plan *CtyunFunctionTriggerConfig, state CtyunFunctionTriggerConfig) (err error) {
	hasChanges := false

	if plan.EventData != state.EventData || plan.Enable != state.Enable || plan.Version != state.Version {
		req := &cf.CfUpdateTriggerV2024Request{
			FunctionName:  plan.FunctionName.ValueString(),
			RegionId:      plan.RegionID.ValueString(),
			TriggerName:   plan.TriggerName.ValueString(),
			TriggerConfig: plan.EventData.ValueString(),
			Version:       plan.Version.ValueString(),
			Enable:        plan.Enable.ValueBool(),
		}

		resp, err := c.meta.Apis.SdkCtCfApis.CfUpdateTriggerV2024Api.Do(ctx, c.meta.SdkCredential, req)
		if err != nil {
			return err
		} else if resp.StatusCode != common.FaasNormalStatusCode {
			return fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Message)
		} else if resp.ReturnObj == nil {
			return common.InvalidReturnObjError
		}
		hasChanges = true
	}
	if !hasChanges {
		return nil
	}
	return
}

func (c *CtyunFunctionTrigger) delete(ctx context.Context, state *CtyunFunctionTriggerConfig) (err error) {
	req := &cf.CfDeleteTriggerV2024Request{
		FunctionName: state.FunctionName.ValueString(),
		TriggerName:  state.TriggerName.ValueString(),
		RegionId:     state.RegionID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtCfApis.CfDeleteTriggerV2024Api.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if *resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Message)
		return
	}

	return
}

func (c *CtyunFunctionTrigger) getAndMerge(ctx context.Context, config *CtyunFunctionTriggerConfig) (err error) {
	functionName := config.FunctionName.ValueString()
	triggerName := config.TriggerName.ValueString()

	if functionName == "" || triggerName == "" {
		idParts := config.ID.ValueString()
		if idParts != "" {
			idPartsArray := strings.Split(idParts, ",")
			if len(idPartsArray) >= 3 {
				functionName = idPartsArray[0]
				triggerName = idPartsArray[1]
			}
		}
	}

	if functionName == "" || triggerName == "" {
		return fmt.Errorf("function_name and trigger_name are required")
	}

	req := &cf.CfGetTriggerRequest{
		FunctionName: functionName,
		TriggerName:  triggerName,
		RegionId:     config.RegionID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtCfApis.CfGetTriggerApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp.StatusCode != common.FaasNormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	returnObj := resp.ReturnObj

	config.FunctionName = types.StringValue(functionName)
	config.TriggerName = types.StringPointerValue(returnObj.TriggerName)
	config.TriggerType = types.StringPointerValue(returnObj.TriggerType)
	if config.EventData.ValueString() == "" {
		config.EventData = types.StringPointerValue(returnObj.TriggerConfig)
	}
	config.Version = types.StringPointerValue(returnObj.Version)
	config.Status = types.Int32PointerValue(returnObj.Status)
	config.CreatedAt = types.StringPointerValue(returnObj.CreatedAt)
	config.UpdatedAt = types.StringPointerValue(returnObj.UpdatedAt)

	if returnObj.Status != nil {
		config.Enable = types.BoolValue(*returnObj.Status == 1)
	}

	config.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", functionName, triggerName, config.RegionID.ValueString()))

	return
}
