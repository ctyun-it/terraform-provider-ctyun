package iam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctiam"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strconv"
)

func NewCtyunIdp() resource.Resource {
	return &ctyunIdp{}
}

type ctyunIdp struct {
	meta *common.CtyunMetadata
}

func (c *ctyunIdp) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_idp"
}

func (c *ctyunIdp) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10345725/10390452`,
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "id",
			},
			"account_id": schema.StringAttribute{
				Computed:    true,
				Description: "账号id",
			},
			"file": schema.StringAttribute{
				Required:    true,
				Description: "联邦登录文件",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"file_name": schema.StringAttribute{
				Required:    true,
				Description: "文件名称（需携带后缀）",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(regexp.MustCompile(`^.+\..+$`), "文件名需要携带后缀"),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "身份提供商名称",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "类型，virtual：虚拟用户SSO，iam：IAM用户SSO",
				Validators: []validator.String{
					stringvalidator.OneOf(business.IdpTypes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString(business.IdpTypeIam),
			},
			"protocol": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "协议类型，saml：SAML协议，oidc：OIDC协议，不填默认为SAML协议",
				Validators: []validator.String{
					stringvalidator.OneOf(business.IdpProtocols...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString(business.IdpProtocolSaml),
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "描述",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (c *ctyunIdp) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunIdpConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	requestType, err := business.IdpTypeMap.FromOriginalScene(plan.Type.ValueString(), business.IdpTypeMapScene1)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	requestProtocol, err := business.IdpProtocolMap.FromOriginalScene(plan.Protocol.ValueString(), business.IdpProtocolMapScene1)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	resp, err := c.meta.Apis.CtIamApis.IdpCreateApi.Do(ctx, c.meta.Credential, &ctiam.IdpCreateRequest{
		Name:     plan.Name.ValueString(),
		Type:     requestType.(int),
		Protocol: requestProtocol.(int),
		Remark:   plan.Description.ValueString(),
		FileName: plan.FileName.ValueString(),
		File:     []byte(plan.File.ValueString()),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	plan.Id = types.Int64Value(resp.Id)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, ctyunRequestError := c.getAndMergeIdp(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunIdp) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunIdpConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	instance, err := c.getAndMergeIdp(ctx, state)
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

func (c *ctyunIdp) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state CtyunIdpConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	var plan CtyunIdpConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)

	_, err := c.meta.Apis.CtIamApis.IdpUpdateApi.Do(ctx, c.meta.Credential, &ctiam.IdpUpdateRequest{
		Id:       state.Id.ValueInt64(),
		Remark:   plan.Description.ValueString(),
		FileName: plan.FileName.ValueString(),
		File:     []byte(plan.File.ValueString()),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	instance, ctyunRequestError := c.getAndMergeIdp(ctx, state)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunIdp) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunIdpConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtIamApis.IdpDeleteApi.Do(ctx, c.meta.Credential, &ctiam.IdpDeleteRequest{
		Id: state.Id.ValueInt64(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [idpId]
func (c *ctyunIdp) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunIdpConfig
	var idpId string
	err := terraform_extend.Split(request.ID, &idpId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	value, err := strconv.ParseInt(idpId, 10, 64)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.Id = types.Int64Value(value)
	instance, err := c.getAndMergeIdp(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunIdp) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// getAndMergeIdp 查询idp
func (c *ctyunIdp) getAndMergeIdp(ctx context.Context, cfg CtyunIdpConfig) (*CtyunIdpConfig, error) {
	resp, err := c.meta.Apis.CtIamApis.IdpListApi.Do(ctx, c.meta.Credential, &ctiam.IdpListRequest{
		Id: cfg.Id.ValueInt64(),
	})
	if err != nil {
		return nil, err
	}
	// 被删除的状态，status为0
	if resp.Status == 0 {
		return nil, nil
	}

	responseType, err2 := business.IdpTypeMap.ToOriginalScene(resp.Type, business.IdpTypeMapScene1)
	if err2 != nil {
		return nil, err2
	}
	responseProtocol, err2 := business.IdpProtocolMap.ToOriginalScene(resp.Protocol, business.IdpProtocolMapScene1)
	if err2 != nil {
		return nil, err2
	}

	cfg.Id = types.Int64Value(resp.Id)
	cfg.AccountId = types.StringValue(resp.AccountId)
	cfg.File = types.StringValue(resp.MetadataDocument)
	cfg.FileName = types.StringValue(resp.FileName)
	cfg.Name = types.StringValue(resp.Name)
	cfg.Type = types.StringValue(responseType.(string))
	cfg.Protocol = types.StringValue(responseProtocol.(string))
	cfg.Description = types.StringValue(resp.Remark)
	return &cfg, nil
}

type CtyunIdpConfig struct {
	Id          types.Int64  `tfsdk:"id"`
	AccountId   types.String `tfsdk:"account_id"`
	File        types.String `tfsdk:"file"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Protocol    types.String `tfsdk:"protocol"`
	Description types.String `tfsdk:"description"`
	FileName    types.String `tfsdk:"file_name"`
}
