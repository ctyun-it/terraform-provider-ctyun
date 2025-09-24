package elb

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &CtyunElbRule{}
	_ resource.ResourceWithConfigure   = &CtyunElbRule{}
	_ resource.ResourceWithImportState = &CtyunElbRule{}
)

type CtyunElbRule struct {
	meta *common.CtyunMetadata
}

func NewCtyunElbRule() resource.Resource {
	return &CtyunElbRule{}
}

func (c *CtyunElbRule) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunElbRule) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunElbRule) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_rule"
}

func (c *CtyunElbRule) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026756/10032110**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池Id，默认使用provider ctyun总region_id 或者环境变量",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"listener_id": schema.StringAttribute{
				Required:    true,
				Description: "监听器listener Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:'{},./;'[,]·~！@#￥%……&*（） —— -+={}，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					validator2.Desc(),
				},
			},
			"conditions": schema.ListNestedAttribute{
				Required:    true,
				Description: "匹配规则数据，支持更新",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"condition_type": schema.StringAttribute{
							Required:    true,
							Description: "匹配规则类型。取值范围：server_name（服务名称）、url_path（匹配路径），支持更新",
							Validators: []validator.String{
								stringvalidator.OneOf(business.ElbRuleConditionTypes...),
							},
						},
						"condition_server_name": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "服务名称,格式为：xxx.xxx结构，不支持下划线'_'。当type = server_name填写，支持更新",
							Validators: []validator.String{
								validator2.AlsoRequiresEqualString(
									path.MatchRelative().AtParent().AtName("condition_type"),
									types.StringValue(business.ElbRuleConditionTypeServerName),
								),
								validator2.ConflictsWithEqualString(
									path.MatchRelative().AtParent().AtName("condition_type"),
									types.StringValue(business.ElbRuleConditionTypeUrlPath),
								),
							},
						},
						"condition_url_paths": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "匹配路径。当type = url_path填写，支持更新",
							Validators: []validator.String{
								validator2.AlsoRequiresEqualString(
									path.MatchRelative().AtParent().AtName("condition_type"),
									types.StringValue(business.ElbRuleConditionTypeUrlPath),
								),
								validator2.ConflictsWithEqualString(
									path.MatchRelative().AtParent().AtName("condition_type"),
									types.StringValue(business.ElbRuleConditionTypeServerName),
								),
							},
						},
						"condition_match_type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "匹配类型。取值范围：ABSOLUTE，PREFIX，REG，支持更新",
							Validators: []validator.String{
								stringvalidator.OneOf(business.ElbRuleMatchTypes...),
								validator2.AlsoRequiresEqualString(
									path.MatchRelative().AtParent().AtName("condition_type"),
									types.StringValue(business.ElbRuleConditionTypeUrlPath),
								),
								validator2.ConflictsWithEqualString(
									path.MatchRelative().AtParent().AtName("condition_type"),
									types.StringValue(business.ElbRuleConditionTypeServerName),
								),
							},
						},
					},
				},
			},
			"action_type": schema.StringAttribute{
				Required:    true,
				Description: "默认规则动作类型。取值范围：forward、redirect，支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ElbRuleActionType...),
				},
			},
			"action_target_groups": schema.ListNestedAttribute{
				Optional:    true,
				Description: "后端服务组，支持更新",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"target_group_id": schema.StringAttribute{
							Required:    true,
							Description: "后端服务组ID，支持更新",
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
						"weight": schema.Int32Attribute{
							Optional:    true,
							Computed:    true,
							Description: "权重，取值范围：1-256。默认为100，支持更新",
							Default:     int32default.StaticInt32(100),
							Validators: []validator.Int32{
								int32validator.Between(1, 256),
							},
						},
					},
				},
			},
			"action_redirect_listener_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "重定向监听器ID，当action_type = redirect时，必填，支持更新",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("action_type"),
						types.StringValue(business.ElbRuleActionTypeRedirect),
					),
				},
			},
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "转发规则 ID",
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				// az时候有必要设定默认值
				Default: defaults.AcquireFromGlobalString(common.ExtraAzName, false),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"load_balancer_id": schema.StringAttribute{
				Computed:    true,
				Description: "负载均衡Id",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "状态: ACTIVE / DOWN",
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
			},
			"updated_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间，为UTC格式",
			},
		},
	}
}

func (c *CtyunElbRule) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunElbRuleConfig

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 开始创建
	err = c.createElbRule(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后反查创建后的Rule信息
	err = c.getAndMergeRule(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunElbRule) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunElbRuleConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 确认该rule是否异常
	if !c.acquireAndSetIdIfOrderNotFinished(ctx, &state, response) {
		return
	}
	//查询远端并同步state
	err = c.getAndMergeRule(ctx, &state)

	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunElbRule) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置
	var plan CtyunElbRuleConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunElbRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
	}

	// 更新rule信息
	err = c.updateElbRule(ctx, state, plan)
	if err != nil {
		return
	}
	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeRule(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (c *CtyunElbRule) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunElbRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	params := &ctelb.CtelbDeleteRuleRequest{
		ClientToken: uuid.NewString(),
		RegionID:    state.RegionID.ValueString(),
		PolicyID:    state.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbDeleteRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}
}

func (c *CtyunElbRule) createElbRule(ctx context.Context, plan *CtyunElbRuleConfig) (err error) {
	//配置创建接口所需请求参数
	if plan.ListenerID.IsNull() {
		err = errors.New("创建转发规则时，ListenerID不能为空")
		return
	}

	params := &ctelb.CtelbCreateRuleRequest{
		ClientToken: uuid.NewString(),
		RegionID:    plan.RegionID.ValueString(),
		ListenerID:  plan.ListenerID.ValueString(),
		Priority:    100, //目前不支持配置此参数,只取默认值100
	}
	// 构建condition参数体
	var conditionList []ConditionsModel
	var conditions []*ctelb.CtelbCreateRuleConditionsRequest
	if plan.Conditions.IsNull() {
		err = errors.New("创建转发规则时，conditions不能为空")
		return
	}
	// 将types.list->[]ConditionsModel
	diags := plan.Conditions.ElementsAs(ctx, &conditionList, false)
	if diags.HasError() {
		return
	}
	for _, conditionItem := range conditionList {
		condition := &ctelb.CtelbCreateRuleConditionsRequest{}
		if conditionItem.ConditionType.IsNull() {
			err = errors.New("创建转发规则时，condition_type不能为空")
			return
		}
		condition.RawType = conditionItem.ConditionType.ValueString()
		if conditionItem.ConditionType.ValueString() == business.ElbRuleConditionTypeServerName {
			// 若conditionType = server_name,传递serverName参数
			condition.ServerNameConfig = &ctelb.CtelbCreateRuleConditionsServerNameConfigRequest{}
			if conditionItem.ConditionServerName.IsNull() {
				err = errors.New("当condition type为server_name, server name必填")
				return
			}
			condition.ServerNameConfig.ServerName = conditionItem.ConditionServerName.ValueString()
		} else if conditionItem.ConditionType.ValueString() == business.ElbRuleConditionTypeUrlPath {
			condition.UrlPathConfig = &ctelb.CtelbCreateRuleConditionsUrlPathConfigRequest{}
			if !conditionItem.ConditionMatchType.IsNull() {
				err = errors.New("当condition type = url_path时，urlPaths和matchType必填！")
				return
			}
			if !conditionItem.ConditionUrlPaths.IsNull() {
				err = errors.New("当condition type = url_path时，urlPaths和matchType必填！")
				return
			}
			condition.UrlPathConfig.MatchType = conditionItem.ConditionMatchType.ValueString()
			condition.UrlPathConfig.UrlPaths = conditionItem.ConditionUrlPaths.ValueString()
		} else {
			err = errors.New("condition type 取值有误！")
			return
		}
		conditions = append(conditions, condition)
	}
	params.Conditions = conditions

	// 构建Action请求体
	action := &ctelb.CtelbCreateRuleActionRequest{}
	if plan.ActionType.IsNull() {
		err = errors.New("创建转发规则时，action type不能为空")
	}
	action.RawType = plan.ActionType.ValueString()
	if plan.ActionType.ValueString() == business.ElbRuleActionTypeRedirect {
		if plan.ActionType.ValueString() == business.ElbRuleActionTypeRedirect && plan.ActionRedirectListenerID.IsNull() {
			err = errors.New("创建转发规则时，若action type = redirect, redirectListenerID不能为空")
			return
		}
		if !plan.ActionRedirectListenerID.IsNull() {
			action.RedirectListenerID = plan.ActionRedirectListenerID.ValueString()
		}

	} else if plan.ActionType.ValueString() == business.ElbRuleActionTypeForward {
		// 构建action.forwardConfig请求体
		action.ForwardConfig = &ctelb.CtelbCreateRuleActionForwardConfigRequest{}
		var targetGroupList []TargetGroupModel
		var targetGroups []*ctelb.CtelbCreateRuleActionForwardConfigTargetGroupsRequest
		diags = plan.ActionTargetGroups.ElementsAs(ctx, &targetGroupList, false)
		if diags.HasError() {
			return
		}
		for _, targetGroupItem := range targetGroupList {
			targetGroup := &ctelb.CtelbCreateRuleActionForwardConfigTargetGroupsRequest{}
			if targetGroupItem.TargetGroupID.IsNull() {
				err = errors.New("创建转发规则时，targetGroupID不能为空")
				return
			}
			targetGroup.TargetGroupID = targetGroupItem.TargetGroupID.ValueString()
			if !targetGroupItem.Weight.IsNull() {
				targetGroup.Weight = targetGroupItem.Weight.ValueInt32()
			}
			targetGroups = append(targetGroups, targetGroup)
		}
		action.ForwardConfig.TargetGroups = targetGroups
	}
	params.Action = action
	// 调用创建接口
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbCreateRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}
	var ids []string
	// 获取规则id
	for _, returnOjbItem := range resp.ReturnObj {
		ids = append(ids, returnOjbItem.ID)
	}

	idsTmp := strings.Join(ids, ",")
	plan.ID = types.StringValue(idsTmp)
	return
}

func (c *CtyunElbRule) getAndMergeRule(ctx context.Context, plan *CtyunElbRuleConfig) (err error) {
	//查看rule详情
	params := &ctelb.CtelbShowRuleRequest{
		RegionID: plan.RegionID.ValueString(),
		PolicyID: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbShowRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}
	returnObj := resp.ReturnObj
	//解析rule明细，将远端明细合并至本地
	plan.LoadBalancerID = types.StringValue(returnObj.LoadBalancerID)
	plan.Status = types.StringValue(returnObj.Status)
	plan.CreatedTime = types.StringValue(returnObj.CreatedTime)
	plan.UpdatedTime = types.StringValue(returnObj.UpdatedTime)

	// 合并conditions
	conditionList := returnObj.Conditions
	var conditions []ConditionsModel
	for _, conditionItem := range conditionList {
		var condition ConditionsModel
		condition.ConditionType = types.StringValue(conditionItem.RawType)
		condition.ConditionServerName = types.StringValue(conditionItem.ServerNameConfig.ServerName)
		condition.ConditionUrlPaths = types.StringValue(conditionItem.UrlPathConfig.UrlPaths)
		condition.ConditionMatchType = types.StringValue(conditionItem.UrlPathConfig.MatchType)
		conditions = append(conditions, condition)
	}
	plan.Conditions, _ = types.ListValueFrom(ctx, utils.StructToTFObjectTypes(ConditionsModel{}), conditions)

	// 处理action
	plan.ActionType = types.StringValue(returnObj.Action.RawType)
	plan.ActionRedirectListenerID = types.StringValue(returnObj.Action.RedirectListenerID)

	targetGroupList := returnObj.Action.ForwardConfig.TargetGroups
	var targetGroups []TargetGroupModel
	for _, targetGroupItem := range targetGroupList {
		var targetGroup TargetGroupModel
		targetGroup.TargetGroupID = types.StringValue(targetGroupItem.TargetGroupID)
		targetGroup.Weight = types.Int32Value(targetGroupItem.Weight)
		targetGroups = append(targetGroups, targetGroup)
	}
	plan.ActionTargetGroups, _ = types.ListValueFrom(ctx, utils.StructToTFObjectTypes(TargetGroupModel{}), targetGroups)
	return
}

func (c *CtyunElbRule) updateElbRule(ctx context.Context, state CtyunElbRuleConfig, plan CtyunElbRuleConfig) (err error) {
	params := &ctelb.CtelbUpdateRuleRequest{
		ClientToken: uuid.NewString(),
		RegionID:    state.RegionID.ValueString(),
		PolicyID:    state.ID.ValueString(),
		Conditions:  nil,
		Action:      nil,
	}
	// 处理condition更新值
	var conditionList []ConditionsModel
	var conditions []*ctelb.CtelbUpdateRuleConditionsRequest
	// 将types.list->[]ConditionsModel
	diags := plan.Conditions.ElementsAs(ctx, &conditionList, false)
	if diags.HasError() {
		return
	}
	if len(conditionList) > 0 {
		for _, conditionItem := range conditionList {
			condition := &ctelb.CtelbUpdateRuleConditionsRequest{}
			if conditionItem.ConditionType.IsNull() {
				err = errors.New("更新转发规则时，condition_type不能为空")
				return
			}

			condition.RawType = conditionItem.ConditionType.ValueString()
			if conditionItem.ConditionType.ValueString() == business.ElbRuleConditionTypeServerName {
				condition.ServerNameConfig = &ctelb.CtelbUpdateRuleConditionsServerNameConfigRequest{}

				if conditionItem.ConditionServerName.IsNull() {
					err = errors.New("当condition type为server_name, server name必填")
				}
				condition.ServerNameConfig.ServerName = conditionItem.ConditionServerName.ValueString()
			} else if conditionItem.ConditionType.ValueString() == business.ElbRuleConditionTypeUrlPath {
				condition.UrlPathConfig = &ctelb.CtelbUpdateRuleConditionsUrlPathConfigRequest{}

				if conditionItem.ConditionMatchType.IsNull() || conditionItem.ConditionUrlPaths.IsNull() {
					err = errors.New("当condition type = url_path时，urlPaths和matchType必填！")
					return
				}
				condition.UrlPathConfig.MatchType = conditionItem.ConditionMatchType.ValueString()
				condition.UrlPathConfig.UrlPaths = conditionItem.ConditionUrlPaths.ValueString()
			} else {
				err = errors.New("condition type 取值有误！")
				return
			}

			conditions = append(conditions, condition)
		}
		params.Conditions = conditions
	}

	// 处理action更新值

	action := &ctelb.CtelbUpdateRuleActionRequest{}
	if plan.ActionType.IsNull() {
		err = errors.New("修改转发规则时，action type不能为空")
	}
	action.RawType = plan.ActionType.ValueString()
	if plan.ActionType.ValueString() == business.ElbRuleActionTypeRedirect {
		if plan.ActionType.ValueString() == business.ElbRuleActionTypeRedirect && plan.ActionRedirectListenerID.IsNull() {
			err = errors.New("修改转发规则时，若action type = redirect, redirectListenerID不能为空")
			return
		}
		if !plan.ActionRedirectListenerID.IsNull() {
			action.RedirectListenerID = plan.ActionRedirectListenerID.ValueString()
		}
	} else if plan.ActionType.ValueString() == business.ElbRuleActionTypeForward {
		action.ForwardConfig = &ctelb.CtelbUpdateRuleActionForwardConfigRequest{}
		// 构建action.forwardConfig请求体
		var targetGroupList []TargetGroupModel
		var targetGroups []*ctelb.CtelbUpdateRuleActionForwardConfigTargetGroupsRequest
		diags = plan.ActionTargetGroups.ElementsAs(ctx, &targetGroupList, false)
		if diags.HasError() {
			return
		}
		for _, targetGroupItem := range targetGroupList {
			var targetGroup ctelb.CtelbUpdateRuleActionForwardConfigTargetGroupsRequest
			if targetGroupItem.TargetGroupID.IsNull() {
				err = errors.New("修改转发规则时，targetGroupID不能为空")
				return
			}
			targetGroup.TargetGroupID = targetGroupItem.TargetGroupID.ValueString()
			if !targetGroupItem.Weight.IsNull() {
				targetGroup.Weight = targetGroupItem.Weight.ValueInt32()
			}
			targetGroups = append(targetGroups, &targetGroup)
		}
		action.ForwardConfig.TargetGroups = targetGroups
	}

	params.Action = action

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}
	return
}

func (c *CtyunElbRule) acquireAndSetIdIfOrderNotFinished(ctx context.Context, state *CtyunElbRuleConfig, response *resource.ReadResponse) bool {
	if state.ID.IsNull() {
		// 该rule没有id，为非法id。移除当前状态并返回
		response.State.RemoveResource(ctx)
		return false
	}
	return true
}

type CtyunElbRuleConfig struct {
	RegionID                 types.String `tfsdk:"region_id"`                   //区域ID
	ListenerID               types.String `tfsdk:"listener_id"`                 //监听器ID
	Description              types.String `tfsdk:"description"`                 //支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:'{},./;'[,]·~！@#￥%……&*（） —— -+={},
	Conditions               types.List   `tfsdk:"conditions"`                  //匹配规则数据
	ActionType               types.String `tfsdk:"action_type"`                 //默认规则动作类型。取值范围：forward、redirect、deny(目前暂不支持配置为deny)
	ActionTargetGroups       types.List   `tfsdk:"action_target_groups"`        //后端服务组
	ActionRedirectListenerID types.String `tfsdk:"action_redirect_listener_id"` //重定向监听器ID，当type为redirect时，此字段必填
	ID                       types.String `tfsdk:"id"`                          //转发规则 ID
	AzName                   types.String `tfsdk:"az_name"`                     //可用区名称
	ProjectID                types.String `tfsdk:"project_id"`                  //	项目ID
	LoadBalancerID           types.String `tfsdk:"load_balancer_id"`            //负载均衡ID
	Status                   types.String `tfsdk:"status"`                      //状态: ACTIVE / DOWN
	CreatedTime              types.String `tfsdk:"created_time"`                //创建时间，为UTC格式
	UpdatedTime              types.String `tfsdk:"updated_time"`                //更新时间，为UTC格式
}

type ConditionsModel struct {
	ConditionType       types.String `tfsdk:"condition_type"`        //类型。取值范围：server_name（服务名称）、url_path（匹配路径）
	ConditionServerName types.String `tfsdk:"condition_server_name"` //服务名称
	ConditionUrlPaths   types.String `tfsdk:"condition_url_paths"`   //匹配路径
	ConditionMatchType  types.String `tfsdk:"condition_match_type"`  //匹配类型。取值范围：ABSOLUTE，PREFIX，REG
}

type TargetGroupModel struct {
	TargetGroupID types.String `tfsdk:"target_group_id"` //后端服务组ID
	Weight        types.Int32  `tfsdk:"weight"`          //权重，取值范围：1-256。默认为100
}
