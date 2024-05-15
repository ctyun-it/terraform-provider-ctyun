package resource

import (
	"context"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctiam"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
)

func NewCtyunPolicy() resource.Resource {
	return &ctyunPolicy{}
}

type ctyunPolicy struct {
	meta *common.CtyunMetadata
}

func (c *ctyunPolicy) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_iam_policy"
}

func (c *ctyunPolicy) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10345725/10390484**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "绑定关系id",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "策略的名称，长度最大为64",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(64),
				},
			},
			"range": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "策略范围，region：资源池级别，global：全局级别，默认为全局级别global",
				Default:     stringdefault.StaticString(business.PolicyRangeGlobal),
				Validators: []validator.String{
					stringvalidator.OneOf(business.PolicyRanges...),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "策略描述，长度最大为128",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(128),
				},
			},
			"content": schema.SingleNestedAttribute{
				Required:    true,
				Description: "权限控制的对象",
				Attributes: map[string]schema.Attribute{
					"version": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "权限控制的版本号，默认为1.1",
						Default:     stringdefault.StaticString("1.1"),
						Validators: []validator.String{
							stringvalidator.OneOf("1.1"),
						},
					},
					"statement": schema.SetNestedAttribute{
						Required:    true,
						Description: "权限控制的对象信息，必填，且数量至少为1",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"action": schema.SetAttribute{
									Required:    true,
									ElementType: types.StringType,
									Description: "对应权限点的code，必填，项目至少为1个，详见ctyun_iam_authorities中的code属性",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(stringvalidator.UTF8LengthAtMost(128)),
									},
								},
								"effect": schema.StringAttribute{
									Required:    true,
									Description: "对应的权限策略动作，allow：允许，deny：拒绝",
									Validators: []validator.String{
										stringvalidator.OneOf(business.PolicyEffects...),
									},
								},
								"resource": schema.SetAttribute{
									Optional:    true,
									Computed:    true,
									ElementType: types.StringType,
									Description: "资源池级别的维度，当权限点为资源池级别时候才生效，不填默认写*",
									Default:     setdefault.StaticValue(types.SetValueMust(basetypes.StringType{}, []attr.Value{types.StringValue("*")})),
								},
							},
						},
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
				},
			},
		},
	}
}

func (c *ctyunPolicy) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunPolicyConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	policyRange, err := business.PolicyRangeMap.FromOriginalScene(plan.Range.ValueString(), business.PolicyRangeMapScene1)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	var ppRequest []ctiam.PolicyCreatePolicyContentStatementRequest
	for _, st := range plan.Content.Statement {
		effect, err := business.PolicyEffectMap.FromOriginalScene(st.Effect.ValueString(), business.PolicyEffectMapScene1)
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
		var res, act []string
		st.Resource.ElementsAs(ctx, &res, true)
		st.Action.ElementsAs(ctx, &act, true)
		ppRequest = append(ppRequest, ctiam.PolicyCreatePolicyContentStatementRequest{
			Resource: res,
			Action:   act,
			Effect:   effect.(string),
		})
	}

	resp, err := c.meta.Apis.CtIamApis.PolicyCreateApi.Do(ctx, c.meta.Credential, &ctiam.PolicyCreateRequest{
		PolicyName:        plan.Name.ValueString(),
		PolicyRange:       policyRange.(int),
		PolicyDescription: plan.Description.ValueString(),
		PolicyContent: ctiam.PolicyCreatePolicyContentRequest{
			Version:   plan.Content.Version.ValueString(),
			Statement: ppRequest,
		},
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	plan.Id = types.StringValue(resp.PolicyId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, ctyunRequestError := c.getAndMergeIamPolicy(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunPolicy) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunPolicyConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeIamPolicy(ctx, state)
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

func (c *ctyunPolicy) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan CtyunPolicyConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var state CtyunPolicyConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	policyRange, err := business.PolicyRangeMap.FromOriginalScene(plan.Range.ValueString(), business.PolicyRangeMapScene1)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	var ppRequest []ctiam.PolicyUpdatePolicyContentStatementRequest
	for _, st := range plan.Content.Statement {
		effect, err := business.PolicyEffectMap.FromOriginalScene(st.Effect.ValueString(), business.PolicyEffectMapScene1)
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
		var res, act []string
		st.Resource.ElementsAs(ctx, &res, true)
		st.Action.ElementsAs(ctx, &act, true)
		ppRequest = append(ppRequest, ctiam.PolicyUpdatePolicyContentStatementRequest{
			Resource: res,
			Action:   act,
			Effect:   effect.(string),
		})
	}

	_, err2 := c.meta.Apis.CtIamApis.PolicyUpdateApi.Do(ctx, c.meta.Credential, &ctiam.PolicyUpdateRequest{
		PolicyId:          state.Id.ValueString(),
		PolicyName:        plan.Name.ValueString(),
		PolicyRange:       policyRange.(int),
		PolicyDescription: plan.Description.ValueString(),
		PolicyContent: ctiam.PolicyUpdatePolicyContentRequest{
			Version:   plan.Content.Version.ValueString(),
			Statement: ppRequest,
		},
	})
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}

	instance, ctyunRequestError := c.getAndMergeIamPolicy(ctx, state)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunPolicy) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunPolicyConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtIamApis.PolicyDeleteApi.Do(ctx, c.meta.Credential, &ctiam.PolicyDeleteRequest{
		PolicyId: state.Id.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [policyId]
func (c *ctyunPolicy) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunPolicyConfig
	var policyId string
	err := terraform_extend.Split(request.ID, &policyId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.Id = types.StringValue(policyId)

	instance, err := c.getAndMergeIamPolicy(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunPolicy) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// getAndMergeIamPolicy 查询策略
func (c *ctyunPolicy) getAndMergeIamPolicy(ctx context.Context, cfg CtyunPolicyConfig) (*CtyunPolicyConfig, error) {
	resp, err := c.meta.Apis.CtIamApis.PolicyGetApi.Do(ctx, c.meta.Credential, &ctiam.PolicyGetRequest{
		PolicyId: cfg.Id.ValueString(),
	})
	if err != nil {
		return nil, err
	}
	if resp.Id == "" {
		return nil, nil
	}

	content := resp.PolicyContent
	statement := []CtyunPolicyStatementConfig{}
	for _, v := range content.Statement {
		act := []types.String{}
		res := []types.String{}
		for _, a := range v.Action {
			act = append(act, types.StringValue(a))
		}
		for _, s := range v.Resource {
			res = append(res, types.StringValue(s))
		}
		action, _ := types.SetValueFrom(ctx, types.StringType, act)
		ress, _ := types.SetValueFrom(ctx, types.StringType, res)
		effect, err := business.PolicyEffectMap.ToOriginalScene(v.Effect, business.PolicyEffectMapScene1)
		if err != nil {
			return nil, err
		}
		statement = append(statement, CtyunPolicyStatementConfig{
			Effect:   types.StringValue(effect.(string)),
			Action:   action,
			Resource: ress,
		})
	}

	policyRange, err2 := business.PolicyRangeMap.ToOriginalScene(resp.PolicyRange, business.PolicyRangeMapScene1)
	if err2 != nil {
		return nil, err2
	}
	cfg.Name = types.StringValue(resp.PolicyName)
	cfg.Range = types.StringValue(policyRange.(string))
	cfg.Description = types.StringValue(resp.PolicyDescription)
	cfg.Content = CtyunPolicyContentConfig{
		Version:   types.StringValue(content.Version),
		Statement: statement,
	}
	return &cfg, nil
}

type CtyunPolicyStatementConfig struct {
	Effect   types.String `tfsdk:"effect"`
	Action   types.Set    `tfsdk:"action"`
	Resource types.Set    `tfsdk:"resource"`
}

type CtyunPolicyContentConfig struct {
	Version   types.String                 `tfsdk:"version"`
	Statement []CtyunPolicyStatementConfig `tfsdk:"statement"`
}

type CtyunPolicyConfig struct {
	Id          types.String             `tfsdk:"id"`
	Name        types.String             `tfsdk:"name"`
	Range       types.String             `tfsdk:"range"`
	Content     CtyunPolicyContentConfig `tfsdk:"content"`
	Description types.String             `tfsdk:"description"`
}
