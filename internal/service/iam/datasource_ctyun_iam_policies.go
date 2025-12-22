package iam

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctiam"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCtyunIamPolicies() datasource.DataSource {
	return &CtyunIamPolicies{}
}

type CtyunIamPolicies struct {
	meta       *common.CtyunMetadata
	iamService *business.IamService
}

func (c *CtyunIamPolicies) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_policies"
}

func (c *CtyunIamPolicies) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10345725/10390484`,
		Attributes: map[string]schema.Attribute{
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页显示数量",
				Validators: []validator.Int32{
					int32validator.AlsoRequires(path.MatchRoot("page_no")),
					int32validator.Between(1, 1000),
				},
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "当前页码",
				Validators: []validator.Int32{
					int32validator.AlsoRequires(path.MatchRoot("page_size")),
					int32validator.AtLeast(1),
				},
			},
			"policy_id": schema.StringAttribute{
				Optional:    true,
				Description: "策略ID，传递时分页参数不生效",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"policies": schema.ListNestedAttribute{
				Computed:    true,
				Description: "策略列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "策略ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "策略名称",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "策略类型，system：系统策略，custom：自定义策略",
						},
						"range": schema.StringAttribute{
							Computed:    true,
							Description: "策略范围，region：资源池级别，global：全局",
						},
						"content": schema.StringAttribute{
							Computed:    true,
							Description: "策略内容",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "策略描述",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunIamPolicies) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunIamPoliciesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	if config.PolicyID.ValueString() == "" {
		if config.PageNo.ValueInt32() == 0 {
			err = fmt.Errorf("分页参数和策略ID必须要有其一")
			return
		}
		err = c.listByPage(ctx, &config)
	} else {
		err = c.getByID(ctx, &config)
	}

	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *CtyunIamPolicies) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.iamService = business.NewIamService(meta)
}

// getByID 通过ID查询
func (c *CtyunIamPolicies) getByID(ctx context.Context, config *CtyunIamPoliciesConfig) (err error) {
	params := &ctiam.CtiamGetPolicyByIdRequest{
		PolicyId: config.PolicyID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtIamApis.CtiamGetPolicyByIdApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	config.Policies = []CtyunIamPoliciesModel{}
	policy := resp.ReturnObj
	item := CtyunIamPoliciesModel{
		PolicyID:    utils.SecStringValue(policy.Id),
		Name:        utils.SecStringValue(policy.PolicyName),
		Type:        types.StringValue(map[int32]string{1: business.PolicyTypeSystem, 2: business.PolicyTypeCustom}[policy.PolicyType]),
		Range:       types.StringValue(map[int32]string{1: business.PolicyRangeRegion, 2: business.PolicyRangeGlobal}[policy.PolicyRange]),
		Content:     utils.SecStringValue(policy.PolicyContent),
		Description: utils.SecStringValue(policy.PolicyDescription),
	}
	config.Policies = append(config.Policies, item)
	return
}

// listByPage 分页查询
func (c *CtyunIamPolicies) listByPage(ctx context.Context, config *CtyunIamPoliciesConfig) (err error) {
	params := &ctiam.CtiamQueryPolicyRequest{
		PageNum:  config.PageNo.ValueInt32(),
		PageSize: config.PageSize.ValueInt32(),
	}
	resp, err := c.meta.Apis.SdkCtIamApis.CtiamQueryPolicyApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	config.Policies = []CtyunIamPoliciesModel{}
	for _, policy := range resp.ReturnObj.List {
		item := CtyunIamPoliciesModel{
			PolicyID:    utils.SecStringValue(policy.Id),
			Name:        utils.SecStringValue(policy.PolicyName),
			Type:        types.StringValue(map[string]string{"1": business.PolicyTypeSystem, "2": business.PolicyTypeCustom}[utils.SecString(policy.PolicyType)]),
			Range:       types.StringValue(map[string]string{"1": business.PolicyRangeRegion, "2": business.PolicyRangeGlobal}[utils.SecString(policy.PolicyRange)]),
			Content:     utils.SecStringValue(policy.PolicyContent),
			Description: utils.SecStringValue(policy.PolicyDescription),
		}
		config.Policies = append(config.Policies, item)
	}
	return
}

type CtyunIamPoliciesConfig struct {
	PageSize types.Int32             `tfsdk:"page_size"`
	PageNo   types.Int32             `tfsdk:"page_no"`
	PolicyID types.String            `tfsdk:"policy_id"`
	Policies []CtyunIamPoliciesModel `tfsdk:"policies"`
}

type CtyunIamPoliciesModel struct {
	PolicyID    types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`        /*  策略名称  */
	Type        types.String `tfsdk:"type"`        /*  策略类型  */
	Range       types.String `tfsdk:"range"`       /*  策略范围  */
	Content     types.String `tfsdk:"content"`     /*  策略内容  */
	Description types.String `tfsdk:"description"` /*  策略描述  */
}
