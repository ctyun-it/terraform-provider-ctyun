package ecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
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
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026730/10230554**`,
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
				Required:    true,
				Description: "公钥",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"finger_print": schema.StringAttribute{
				Computed:    true,
				Description: "密钥对的指纹，采用MD5信息摘要算法",
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
	var plan CtyunKeypairConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := plan.RegionId.ValueString()
	projectId := plan.ProjectId.ValueString()
	createResponse, err := c.meta.Apis.CtEcsApis.KeypairImportApi.Do(ctx, c.meta.Credential, &ctecs.KeypairImportRequest{
		RegionId:    regionId,
		KeyPairName: plan.Name.ValueString(),
		PublicKey:   plan.PublicKey.ValueString(),
		ProjectId:   projectId,
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	plan.Name = types.StringValue(createResponse.KeyPairName)
	plan.RegionId = types.StringValue(regionId)
	plan.ProjectId = types.StringValue(projectId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, ctyunRequestError := c.getAndMergeKeypair(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunKeypair) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunKeypairConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeKeypair(ctx, state)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunKeypair) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (c *ctyunKeypair) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunKeypairConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtEcsApis.KeypairDeleteApi.Do(ctx, c.meta.Credential, &ctecs.KeypairDeleteRequest{
		RegionId:    state.RegionId.ValueString(),
		KeyPairName: state.Name.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [keyPairName],[regionId]
func (c *ctyunKeypair) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunKeypairConfig
	var keyPairName, regionId string
	err := terraform_extend.Split(request.ID, &keyPairName, &regionId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.Name = types.StringValue(keyPairName)
	cfg.RegionId = types.StringValue(regionId)

	instance, err := c.getAndMergeKeypair(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunKeypair) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// getAndMergeKeypair 查询密钥对
func (c *ctyunKeypair) getAndMergeKeypair(ctx context.Context, cfg CtyunKeypairConfig) (*CtyunKeypairConfig, error) {
	listResponse, err := c.meta.Apis.CtEcsApis.KeypairDetailApi.Do(ctx, c.meta.Credential, &ctecs.KeypairDetailRequest{
		RegionId:    cfg.RegionId.ValueString(),
		KeyPairName: cfg.Name.ValueString(),
		PageNo:      1,
		PageSize:    1,
	})
	if err != nil {
		return nil, err
	}
	if len(listResponse.Results) == 0 {
		return nil, nil
	}

	keypairResponse := listResponse.Results[0]
	cfg.PublicKey = types.StringValue(keypairResponse.PublicKey)
	cfg.Name = types.StringValue(keypairResponse.KeyPairName)
	cfg.FingerPrint = types.StringValue(keypairResponse.FingerPrint)
	cfg.Id = types.StringValue(keypairResponse.KeyPairId)
	return &cfg, nil
}

type CtyunKeypairConfig struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	PublicKey   types.String `tfsdk:"public_key"`
	FingerPrint types.String `tfsdk:"finger_print"`
	ProjectId   types.String `tfsdk:"project_id"`
	RegionId    types.String `tfsdk:"region_id"`
}
