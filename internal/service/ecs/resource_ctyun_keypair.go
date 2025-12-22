package ecs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctecs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
)

var (
	_ resource.Resource                = &ctyunKeypair{}
	_ resource.ResourceWithConfigure   = &ctyunKeypair{}
	_ resource.ResourceWithImportState = &ctyunKeypair{}
)

func NewCtyunKeypair() resource.Resource {
	return &ctyunKeypair{}
}

type ctyunKeypair struct {
	meta *common.CtyunMetadata
}

func (c *ctyunKeypair) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_keypair"
}

func (c *ctyunKeypair) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026730/10230554`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "密钥对的id",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "密钥对名称。只能由数字、字母、-组成，不能以数字和-开头、以-结尾，且长度为2-63字符",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 63),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$"), "不满足密钥对名称要求"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"public_key": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "公钥，填写时会导入密钥对，不填写时会创建密钥对",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"private_key": schema.StringAttribute{
				Computed:    true,
				Description: "私钥，创建密钥对场景下才有值",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"finger_print": schema.StringAttribute{
				Computed:    true,
				Description: "密钥对的指纹，采用MD5信息摘要算法",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
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
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
		},
	}
}

func (c *ctyunKeypair) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunKeypairConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 实际创建
	if plan.PublicKey.ValueString() == "" {
		err = c.createKeyPair(ctx, &plan)
	} else {
		err = c.importKeyPair(ctx, plan)
		plan.PrivateKey = types.StringNull()
	}

	err = c.getAndMergeKeypair(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunKeypair) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunKeypairConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.getAndMergeKeypair(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			response.State.RemoveResource(ctx)
			err = nil
		} else {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (c *ctyunKeypair) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (c *ctyunKeypair) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunKeypairConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := ctecs2.CtecsDeleteKeypairV41Request{
		RegionID:    state.RegionId.ValueString(),
		KeyPairName: state.Name.ValueString(),
	}
	_, err := c.meta.Apis.SdkCtEcsApis.CtecsDeleteKeypairV41Api.Do(ctx, c.meta.SdkCredential, &params)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

func (c *ctyunKeypair) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [keyPairName],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunKeypairConfig
	var keyPairName, regionId string
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		keyPairName = request.ID
	} else {

		err = terraform_extend.Split(request.ID, &keyPairName, &regionId)
		if err != nil {
			return
		}
	}
	if keyPairName == "" {
		err = fmt.Errorf("keyPairName不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}

	cfg.Name = types.StringValue(keyPairName)
	cfg.RegionId = types.StringValue(regionId)

	err = c.getAndMergeKeypair(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *ctyunKeypair) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// getAndMergeKeypair 查询密钥对
func (c *ctyunKeypair) getAndMergeKeypair(ctx context.Context, plan *CtyunKeypairConfig) (err error) {
	params := ctecs2.CtecsDetailsKeypairV41Request{
		RegionID:    plan.RegionId.ValueString(),
		KeyPairName: plan.Name.ValueString(),
		ProjectID:   plan.ProjectId.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsDetailsKeypairV41Api.Do(ctx, c.meta.SdkCredential, &params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if len(resp.ReturnObj.Results) == 0 {
		err = common.ResourceNotExistError
		return
	}

	keypairResponse := resp.ReturnObj.Results[0]
	plan.PublicKey = types.StringValue(keypairResponse.PublicKey)
	plan.Name = types.StringValue(keypairResponse.KeyPairName)
	plan.FingerPrint = types.StringValue(keypairResponse.FingerPrint)
	plan.Id = types.StringValue(keypairResponse.KeyPairID)
	return
}

type CtyunKeypairConfig struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	PublicKey   types.String `tfsdk:"public_key"`
	FingerPrint types.String `tfsdk:"finger_print"`
	ProjectId   types.String `tfsdk:"project_id"`
	RegionId    types.String `tfsdk:"region_id"`
	PrivateKey  types.String `tfsdk:"private_key"`
}

// createKeyPari 创建密钥对
func (c *ctyunKeypair) createKeyPair(ctx context.Context, plan *CtyunKeypairConfig) (error error) {
	params := ctecs2.CtecsCreateKeypairV41Request{
		RegionID:    plan.RegionId.ValueString(),
		KeyPairName: plan.Name.ValueString(),
		ProjectID:   plan.ProjectId.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsCreateKeypairV41Api.Do(ctx, c.meta.SdkCredential, &params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	plan.PublicKey = types.StringValue(resp.ReturnObj.PublicKey)
	plan.PrivateKey = types.StringValue(resp.ReturnObj.PrivateKey)
	return
}

// importKeyPair 导入密钥对
func (c *ctyunKeypair) importKeyPair(ctx context.Context, plan CtyunKeypairConfig) (err error) {
	params := ctecs2.CtecsImportKeypairV41Request{
		RegionID:    plan.RegionId.ValueString(),
		KeyPairName: plan.Name.ValueString(),
		PublicKey:   plan.PublicKey.ValueString(),
		ProjectID:   plan.ProjectId.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsImportKeypairV41Api.Do(ctx, c.meta.SdkCredential, &params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}
