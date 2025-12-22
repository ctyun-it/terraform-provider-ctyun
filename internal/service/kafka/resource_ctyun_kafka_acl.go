package kafka

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgkafka "github.com/ctyun-it/terraform-provider-ctyun/internal/core/kafka"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
)

var (
	_ resource.Resource                = &ctyunKafkaAcl{}
	_ resource.ResourceWithConfigure   = &ctyunKafkaAcl{}
	_ resource.ResourceWithImportState = &ctyunKafkaAcl{}
)

type ctyunKafkaAcl struct {
	meta       *common.CtyunMetadata
	vpcService *business.VpcService
	sgService  *business.SecurityGroupService
}

func NewCtyunKafkaAcl() resource.Resource {
	return &ctyunKafkaAcl{}
}

func (c *ctyunKafkaAcl) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_kafka_acl"
}

type CtyunKafkaAclConfig struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	InstanceId  types.String `tfsdk:"instance_id"`
	RegionId    types.String `tfsdk:"region_id"`
	UseNewTopic types.Bool   `tfsdk:"use_new_topic"`
	Topics      types.Set    `tfsdk:"topics"`
	Rules       types.Set    `tfsdk:"rules"`
	rules       []CtyunKafkaAclRule
}

type CtyunKafkaAclRule struct {
	Permission types.String `tfsdk:"permission"`
	UserName   types.String `tfsdk:"user_name"`
	Ip         types.String `tfsdk:"ip"`
	Operation  types.String `tfsdk:"operation"`
}

func (c *ctyunKafkaAcl) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029624/11078051`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "资源唯一标识符",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "策略名称，规则如下：\n以英文字母、数字、下划线开头，且只能由英文字母、数字、句点、中划线、下划线组成。\n长度3-64。\n名称不可重复。",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 64),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9_.-]*$`),
						"必须以英文字母、数字、下划线开头，只能包含英文字母、数字、句点、中划线、下划线",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "实例ID。",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"use_new_topic": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "是否应用到新增主题，默认不应用，支持更新",
			},
			"topics": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "匹配的topic列表",
			},
			"rules": schema.SetNestedAttribute{
				Description: "ACL规则",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"permission": schema.StringAttribute{
							Required:    true,
							Description: "权限，ALLOW:允许，DENY:拒绝，默认：ALLOW 支持更新",
							Validators: []validator.String{
								stringvalidator.OneOf("ALLOW", "DENY"),
							},
						},
						"user_name": schema.StringAttribute{
							Required:    true,
							Description: "用户名，必须是已经集群中创建了的用户，支持更新",
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
						"ip": schema.StringAttribute{
							Optional:    true,
							Description: "ip或网段，多个用半角分号分开，默认*，表示所有ip 支持更新",
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
						"operation": schema.StringAttribute{
							Required:    true,
							Description: "操作，READ:消费，WRITE:生产 支持更新",
							Validators: []validator.String{
								stringvalidator.OneOf("READ", "WRITE"),
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunKafkaAcl) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunKafkaAclConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.calcRules(ctx, &plan)
	if err != nil {
		return
	}
	// 创建
	err = c.create(ctx, plan)
	if err != nil {
		return
	}

	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunKafkaAcl) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunKafkaAclConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunKafkaAcl) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunKafkaAclConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunKafkaAclConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 更新
	err = c.updateAutoMatch(ctx, plan)
	if err != nil {
		return
	}
	state.UseNewTopic = plan.UseNewTopic
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunKafkaAcl) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunKafkaAclConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 销毁
	err = c.destroy(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunKafkaAcl) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.vpcService = business.NewVpcService(meta)
	c.sgService = business.NewSecurityGroupService(meta)
}

func (c *ctyunKafkaAcl) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [instanceId],[name],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunKafkaAclConfig
	var instanceId, regionID, name string

	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) == 1 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &instanceId, &name)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &instanceId, &name, &regionID)
		if err != nil {
			return
		}
	}

	if instanceId == "" {
		err = fmt.Errorf("instanceId不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	if name == "" {
		err = fmt.Errorf("name不能为空")
		return
	}

	cfg.RegionId = types.StringValue(regionID)
	cfg.InstanceId = types.StringValue(instanceId)
	cfg.Name = types.StringValue(name)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// calcRule 将types.Set类型的Rule转换为[]CtyunKafkaAclRule
func (c *ctyunKafkaAcl) calcRules(ctx context.Context, plan *CtyunKafkaAclConfig) (err error) {
	if plan.Rules.IsUnknown() || plan.Rules.IsNull() {
		return
	}
	plan.rules = []CtyunKafkaAclRule{}
	diags := plan.Rules.ElementsAs(ctx, &plan.rules, false)
	if diags.HasError() {
		err = fmt.Errorf(diags.Errors()[0].Detail())
		return
	}
	return
}

// create 创建
func (c *ctyunKafkaAcl) create(ctx context.Context, plan CtyunKafkaAclConfig) (err error) {
	params := &ctgkafka.CtgkafkaAclStrategyCreateRequest{
		RegionId:    plan.RegionId.ValueString(),
		ProdInstId:  plan.InstanceId.ValueString(),
		Name:        plan.Name.ValueString(),
		UseNewTopic: map[bool]string{true: "1", false: "2"}[plan.UseNewTopic.ValueBool()],
	}

	if len(plan.rules) > 0 {
		params.Rules = make([]*ctgkafka.CtgkafkaAclStrategyCreateRulesRequest, 0, len(plan.rules))
		for _, rule := range plan.rules {
			params.Rules = append(params.Rules, &ctgkafka.CtgkafkaAclStrategyCreateRulesRequest{
				Permission: rule.Permission.ValueString(),
				UserName:   rule.UserName.ValueString(),
				Ip:         rule.Ip.ValueString(),
				Operation:  rule.Operation.ValueString(),
			})
		}
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaAclStrategyCreateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj != "create success" {
		err = common.InvalidReturnObjError
		return
	}

	return
}

// updateAutoMatch 更新自动匹配设置
func (c *ctyunKafkaAcl) updateAutoMatch(ctx context.Context, plan CtyunKafkaAclConfig) (err error) {
	params := &ctgkafka.CtgkafkaAclStrategyTurnAutoMatchRequest{
		RegionId:    plan.RegionId.ValueString(),
		ProdInstId:  plan.InstanceId.ValueString(),
		Name:        plan.Name.ValueString(),
		UseNewTopic: map[bool]string{true: "1", false: "2"}[plan.UseNewTopic.ValueBool()],
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaAclStrategyTurnAutoMatchApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		return fmt.Errorf("API return error. Message: %s", resp.Message)
	} else if resp.ReturnObj == "" {
		return common.InvalidReturnObjError
	}
	return
}

// destroy 销毁
func (c *ctyunKafkaAcl) destroy(ctx context.Context, plan CtyunKafkaAclConfig) (err error) {
	params := &ctgkafka.CtgkafkaAclStrategyDeleteRequest{
		RegionId:   plan.RegionId.ValueString(),
		ProdInstId: plan.InstanceId.ValueString(),
		Name:       plan.Name.ValueString(),
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaAclStrategyDeleteApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj != "delete success" {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// getAndMerge 从远端查询
func (c *ctyunKafkaAcl) getAndMerge(ctx context.Context, plan *CtyunKafkaAclConfig) (err error) {
	params := &ctgkafka.CtgkafkaAclStrategyDetailRequest{
		RegionId:   plan.RegionId.ValueString(),
		ProdInstId: plan.InstanceId.ValueString(),
		Name:       plan.Name.ValueString(),
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaAclStrategyDetailApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 设置基本属性
	plan.Id = types.StringValue(fmt.Sprintf("%s,%s,%s", plan.InstanceId.ValueString(), plan.RegionId.ValueString(), plan.Name.ValueString()))
	if resp.ReturnObj.TopicNum > 0 {
		// 设置topics
		topicsSet, diags := types.SetValueFrom(ctx, types.StringType, resp.ReturnObj.Topics)
		if diags.HasError() {
			err = fmt.Errorf("failed to set topics: %v", diags)
			return
		}
		plan.Topics = topicsSet
	} else {
		plan.Topics = types.SetNull(types.StringType)
	}

	// 设置rules
	if len(resp.ReturnObj.Rules) > 0 {
		rules := make([]CtyunKafkaAclRule, 0, len(resp.ReturnObj.Rules))
		for _, rule := range resp.ReturnObj.Rules {
			// 处理可选字段 ip，如果为空则设置为 Null
			ipValue := types.StringNull()
			if rule.Ip != "" && rule.Ip != "*" {
				ipValue = types.StringValue(rule.Ip)
			}

			rules = append(rules, CtyunKafkaAclRule{
				Permission: types.StringValue(rule.Permission),
				UserName:   types.StringValue(rule.UserName),
				Ip:         ipValue,
				Operation:  types.StringValue(rule.Operation),
			})
		}

		rulesObjType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"permission": types.StringType,
				"user_name":  types.StringType,
				"ip":         types.StringType,
				"operation":  types.StringType,
			},
		}

		rulesValue, diags := types.SetValueFrom(ctx, rulesObjType, rules)
		if diags.HasError() {
			err = fmt.Errorf("failed to set rules: %v", diags)
			return
		}
		plan.Rules = rulesValue
		plan.rules = rules
	} else {
		plan.Rules = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"permission": types.StringType,
				"user_name":  types.StringType,
				"ip":         types.StringType,
				"operation":  types.StringType,
			},
		})
	}

	return
}
