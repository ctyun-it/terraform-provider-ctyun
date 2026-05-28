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
	validators "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CtyunFunctionAlias{}
	_ resource.ResourceWithConfigure   = &CtyunFunctionAlias{}
	_ resource.ResourceWithImportState = &CtyunFunctionAlias{}
)

func NewCtyunFunctionAlias() resource.Resource {
	return &CtyunFunctionAlias{}
}

type CtyunFunctionAlias struct {
	meta *common.CtyunMetadata
	name string
}

type CtyunFunctionAliasConfig struct {
	ID           types.String `tfsdk:"id"`
	FunctionName types.String `tfsdk:"function_name"`
	AliasName    types.String `tfsdk:"alias_name"`
	VersionID    types.String `tfsdk:"version_id"`
	Description  types.String `tfsdk:"description"`
	//GrayType      types.Int32  `tfsdk:"gray_type"`
	GrayVersionID types.String `tfsdk:"gray_version_id"`
	GrayWeight    types.Int32  `tfsdk:"gray_weight"`
	RegionID      types.String `tfsdk:"region_id"`
	CreateTime    types.String `tfsdk:"create_time"`
}

func (c *CtyunFunctionAlias) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function_alias"
	c.name = resp.TypeName
}

func (c *CtyunFunctionAlias) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理函数别名", "函数计算（FaaS）", "https://www.ctyun.cn/document/10006234/10827404"),
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
				Description: "函数名称，函数必须存在",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 64),
				},
			},
			"alias_name": schema.StringAttribute{
				Required:    true,
				Description: "别名名称。只能包含字母、数字和中划线。只能字母开头，字母数字结尾。长度在 2~44 之间",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 44),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$|^[a-zA-Z]$`), "别名必须以字母开头，只能包含字母、数字和中划线"),
				},
			},
			"version_id": schema.StringAttribute{
				Required:    true,
				Description: "主版本 ID 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "关于别名的描述 支持更新",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(512),
				},
			},
			//"gray_type": schema.Int32Attribute{
			//	Optional:    true,
			//	Computed:    true,
			//	Description: "灰度类型，当前支持：1、按百分比随机灰度",
			//	Validators: []validator.Int32{
			//		int32validator.OneOf(1),
			//	},
			//},
			"gray_version_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "灰度版本 ID 支持更新",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					validators.AlsoRequiresEqualString(path.MatchRoot("gray_weight")),
				},
			},
			"gray_weight": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "切流的比例。假设值为 5%，函数计算会将 5% 的流量到打到灰度版本，95% 的流量打到主版本。范围是 [0-100] 支持更新",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int32{
					int32validator.Between(0, 100),
					validators.AlsoRequiresEqualInt32(path.MatchRoot("gray_version_id")),
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
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间",
			},
		},
	}
}

func (c *CtyunFunctionAlias) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunFunctionAlias) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunFunctionAliasConfig
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *CtyunFunctionAlias) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunFunctionAliasConfig
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

func (c *CtyunFunctionAlias) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunFunctionAliasConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state CtyunFunctionAliasConfig
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
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *CtyunFunctionAlias) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunFunctionAliasConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunFunctionAlias) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例：%s 失败：%s", c.name, req.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [alias_name],[function_name],<region_id>", c.name)
			resp.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunFunctionAliasConfig
	var aliasName, functionName, regionID string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(req.ID, common.ImportSeparator) == 1 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(req.ID, &aliasName, &functionName)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(req.ID, &aliasName, &functionName, &regionID)
		if err != nil {
			return
		}
	}

	if functionName == "" {
		err = fmt.Errorf("function_name不能为空")
		return
	}

	if aliasName == "" {
		err = fmt.Errorf("alias_name不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}
	config.FunctionName = types.StringValue(functionName)
	config.AliasName = types.StringValue(aliasName)
	config.ID = types.StringValue(req.ID)
	config.RegionID = types.StringValue(regionID)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (c *CtyunFunctionAlias) create(ctx context.Context, plan *CtyunFunctionAliasConfig) (err error) {
	createReq := &cf.CfCreateAliasRequest{
		FunctionName: plan.FunctionName.ValueString(),
		RegionId:     plan.RegionID.ValueString(),
		AliasName:    plan.AliasName.ValueString(),
		VersionId:    plan.VersionID.ValueString(),
		Description:  plan.Description.ValueStringPointer(),
	}

	// 处理灰度配置
	if !plan.GrayVersionID.IsNull() && !plan.GrayVersionID.IsUnknown() &&
		!plan.GrayWeight.IsNull() && !plan.GrayWeight.IsUnknown() {

		createReq.Gray = &cf.CfCreateAliasGrayRequest{
			VersionId: plan.GrayVersionID.ValueString(),
			RawType:   1,
			Config: &cf.CfCreateAliasGrayConfigRequest{
				Weight: plan.GrayWeight.ValueInt32(),
			},
		}
	}

	createResp, err := c.meta.Apis.SdkCtCfApis.CfCreateAliasApi.Do(ctx, c.meta.SdkCredential, createReq)
	if err != nil {
		return
	} else if *createResp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s Description:", *createResp.Message)
		return
	} else if createResp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return
}

func (c *CtyunFunctionAlias) update(ctx context.Context, plan *CtyunFunctionAliasConfig, state CtyunFunctionAliasConfig) (err error) {
	updateReq := &cf.CfUpdateAliasRequest{
		FunctionName: plan.FunctionName.ValueString(),
		RegionId:     plan.RegionID.ValueString(),
		AliasName:    plan.AliasName.ValueString(),
		VersionId:    plan.VersionID.ValueString(),
	}

	// 只有当描述发生变化时才传递
	if !plan.Description.Equal(state.Description) {
		updateReq.Description = plan.Description.ValueStringPointer()
	}

	// Gray 配置
	if !plan.GrayVersionID.IsNull() && !plan.GrayVersionID.IsUnknown() &&
		!plan.GrayWeight.IsNull() && !plan.GrayWeight.IsUnknown() {
		updateReq.Gray = &cf.CfUpdateAliasGrayRequest{
			RawType:   1,
			VersionId: plan.GrayVersionID.ValueString(),
			Config: &cf.CfUpdateAliasGrayConfigRequest{
				Weight: plan.GrayWeight.ValueInt32(),
			},
		}
	} else if plan.GrayVersionID.IsNull() || plan.GrayVersionID.IsUnknown() ||
		plan.GrayWeight.IsNull() || plan.GrayWeight.IsUnknown() {
		updateReq.Gray = nil
	}

	updateResp, err := c.meta.Apis.SdkCtCfApis.CfUpdateAliasApi.Do(ctx, c.meta.SdkCredential, updateReq)
	if err != nil {
		return
	} else if *updateResp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *updateResp.Message, *updateResp.Message)
		return
	} else if updateResp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return
}

func (c *CtyunFunctionAlias) delete(ctx context.Context, state *CtyunFunctionAliasConfig) (err error) {
	deleteReq := &cf.CfDeleteAliasRequest{
		FunctionName: state.FunctionName.ValueString(),
		RegionId:     state.RegionID.ValueString(),
		AliasName:    state.AliasName.ValueString(),
	}

	deleteResp, err := c.meta.Apis.SdkCtCfApis.CfDeleteAliasApi.Do(ctx, c.meta.SdkCredential, deleteReq)
	if err != nil {
		return
	} else if *deleteResp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *deleteResp.Message, *deleteResp.Message)
		return
	}

	return
}

func (c *CtyunFunctionAlias) getAndMerge(ctx context.Context, config *CtyunFunctionAliasConfig) (err error) {
	getReq := &cf.CfGetAliasRequest{
		FunctionName: config.FunctionName.ValueString(),
		RegionId:     config.RegionID.ValueString(),
		AliasName:    config.AliasName.ValueString(),
	}

	getResp, err := c.meta.Apis.SdkCtCfApis.CfGetAliasApi.Do(ctx, c.meta.SdkCredential, getReq)
	if err != nil {
		return
	} else if *getResp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", *getResp.Message)
		return
	} else if getResp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	config.updateFromAPI(getResp.ReturnObj)

	// 设置 ID
	config.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", config.AliasName.ValueString(), config.FunctionName.ValueString(), config.RegionID.ValueString()))

	return
}

// updateFromAPI 从 API 响应更新状态
func (cfg *CtyunFunctionAliasConfig) updateFromAPI(returnObj *cf.CfGetAliasReturnObjResponse) {
	cfg.AliasName = types.StringPointerValue(returnObj.AliasName)
	cfg.VersionID = types.StringPointerValue(returnObj.VersionId)
	cfg.Description = types.StringPointerValue(returnObj.Description)
	cfg.CreateTime = types.StringPointerValue(returnObj.CreateTime)

	// 只更新 API 返回的灰度配置字段

	if returnObj.GrayVersionId != nil {
		cfg.GrayVersionID = types.StringPointerValue(returnObj.GrayVersionId)
	}

	// 处理灰度配置
	if returnObj.AliasGrayConfig != nil && returnObj.AliasGrayConfig.Weight != nil {
		cfg.GrayWeight = types.Int32PointerValue(returnObj.AliasGrayConfig.Weight)
	} else if returnObj.AliasGrayConfig != nil {
		cfg.GrayWeight = types.Int32Null()
	}
}
