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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &CtyunFunctionVersion{}
	_ resource.ResourceWithConfigure   = &CtyunFunctionVersion{}
	_ resource.ResourceWithImportState = &CtyunFunctionVersion{}
)

func NewCtyunFunctionVersion() resource.Resource {
	return &CtyunFunctionVersion{}
}

type CtyunFunctionVersion struct {
	meta *common.CtyunMetadata
	name string
}

type CtyunFunctionVersionConfig struct {
	ID           types.String `tfsdk:"id"`
	FunctionName types.String `tfsdk:"function_name"`
	VersionID    types.String `tfsdk:"version_id"`
	Description  types.String `tfsdk:"description"`
	RegionID     types.String `tfsdk:"region_id"`
	CreateTime   types.String `tfsdk:"create_time"`
}

func (c *CtyunFunctionVersion) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function_version"
	c.name = resp.TypeName
}

func (c *CtyunFunctionVersion) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理函数版本", "函数工作流（FunctionGraph）", "https://www.ctyun.cn/document/10355289"),
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
			"version_id": schema.StringAttribute{
				Computed:    true,
				Description: "版本 ID，发布版本后自动递增",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "版本描述",
				Default:     stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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

func (c *CtyunFunctionVersion) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunFunctionVersion) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunFunctionVersionConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	regionId := c.meta.GetExtraIfEmpty(plan.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("region_id 不能为空！")
		return
	}
	plan.RegionID = types.StringValue(regionId)

	err = c.create(ctx, &plan)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *CtyunFunctionVersion) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunFunctionVersionConfig
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

func (c *CtyunFunctionVersion) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (c *CtyunFunctionVersion) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunFunctionVersionConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunFunctionVersion) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例：%s 失败：%s", c.name, req.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [function_name],[version_id],<region_id>", c.name)
			resp.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunFunctionVersionConfig
	var functionName, versionID, regionID string

	if req.ID == "" {
		err = fmt.Errorf("导入 ID 不能为空")
		return
	}

	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(req.ID, common.ImportSeparator) == 1 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(req.ID, &functionName, &versionID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(req.ID, &functionName, &versionID, &regionID)
		if err != nil {
			return
		}
	}

	if functionName == "" {
		err = fmt.Errorf("function_name 不能为空")
		return
	}

	if versionID == "" {
		err = fmt.Errorf("version_id 不能为空")
		return
	}

	config.FunctionName = types.StringValue(functionName)
	config.VersionID = types.StringValue(versionID)
	config.RegionID = types.StringValue(regionID)
	config.ID = types.StringValue(fmt.Sprintf("%s,%s", functionName, versionID))

	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (c *CtyunFunctionVersion) create(ctx context.Context, plan *CtyunFunctionVersionConfig) (err error) {
	createReq := &cf.CfPublishFunctionVersionRequest{
		FunctionName: plan.FunctionName.ValueString(),
		RegionId:     plan.RegionID.ValueString(),
		Description:  plan.Description.ValueString(),
	}

	createResp, err := c.meta.Apis.SdkCtCfApis.CfPublishFunctionVersionApi.Do(ctx, c.meta.SdkCredential, createReq)
	if err != nil {
		return
	} else if *createResp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", *createResp.Message)
		return
	} else if createResp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	plan.VersionID = types.StringPointerValue(createResp.ReturnObj.VersionId)
	if createResp.ReturnObj.CreateTime != nil {
		createTimeStr := strconv.FormatInt(int64(*createResp.ReturnObj.CreateTime), 10)
		plan.CreateTime = types.StringValue(createTimeStr)
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s", plan.FunctionName.ValueString(), plan.VersionID.ValueString()))

	return
}

func (c *CtyunFunctionVersion) delete(ctx context.Context, state *CtyunFunctionVersionConfig) (err error) {
	deleteReq := &cf.CfDeleteFunctionVersionRequest{
		FunctionName: state.FunctionName.ValueString(),
		VersionId:    state.VersionID.ValueString(),
		RegionId:     state.RegionID.ValueString(),
	}

	deleteResp, err := c.meta.Apis.SdkCtCfApis.CfDeleteFunctionVersionApi.Do(ctx, c.meta.SdkCredential, deleteReq)
	if err != nil {
		return
	} else if *deleteResp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", *deleteResp.Message)
		return
	}

	return
}

func (c *CtyunFunctionVersion) getAndMerge(ctx context.Context, config *CtyunFunctionVersionConfig) (err error) {
	listReq := &cf.CfListFunctionVersionsRequest{
		FunctionName: config.FunctionName.ValueString(),
		RegionId:     config.RegionID.ValueString(),
	}

	listResp, err := c.meta.Apis.SdkCtCfApis.CfListFunctionVersionsApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		return
	} else if *listResp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", *listResp.Message)
		return
	} else if listResp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	found := false
	for _, version := range listResp.ReturnObj.Data {
		if version.VersionId != nil && *version.VersionId == config.VersionID.ValueString() {
			config.Description = types.StringPointerValue(version.Description)

			if version.CreateTime != nil {
				createTime, parseErr := parseCreateTime(*version.CreateTime)
				if parseErr == nil {
					config.CreateTime = types.StringValue(createTime)
				} else {
					config.CreateTime = types.StringPointerValue(version.CreateTime)
				}
			} else {
				config.CreateTime = types.StringNull()
			}

			found = true
			break
		}
	}

	if !found {
		err = common.ResourceNotExistError
		return
	}

	return
}

func parseCreateTime(timeStr string) (string, error) {
	if timeStr == "" {
		return "", fmt.Errorf("empty time string")
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return "", err
	}
	unixTime := t.Unix()
	return strconv.FormatInt(unixTime, 10), nil
}
