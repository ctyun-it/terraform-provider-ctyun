package elb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunElbRules{}
	_ datasource.DataSourceWithConfigure = &ctyunElbRules{}
)

type ctyunElbRules struct {
	meta *common.CtyunMetadata
}

func NewCtyunElbRules() datasource.DataSource {
	return &ctyunElbRules{}
}
func (c *ctyunElbRules) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunElbRules) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_rules"

}

func (c *ctyunElbRules) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026756/10032110**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"ids": schema.StringAttribute{
				Optional:    true,
				Description: "转发规则ID列表，以,分隔",
			},
			"load_balancer_id": schema.StringAttribute{
				Optional:    true,
				Description: "负载均衡实例ID",
			},
			"elb_rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "elb转发规则",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源池ID",
						},
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "可用区名称",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "项目ID",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "转发规则ID",
						},
						"load_balancer_id": schema.StringAttribute{
							Computed:    true,
							Description: "负载均衡ID",
						},
						"listener_id": schema.StringAttribute{
							Computed:    true,
							Description: "监听器ID",
						},
						//"priority": schema.Int32Attribute{
						//	Computed:    true,
						//	Description: "优先级",
						//},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"conditions": schema.ListNestedAttribute{
							Computed:    true,
							Description: "匹配规则数据",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"condition_type": schema.StringAttribute{
										Computed:    true,
										Description: "类型。取值范围：server_name（服务名称）、url_path（匹配路径）",
										Validators: []validator.String{
											stringvalidator.OneOf(business.ElbRuleConditionTypes...),
										},
									},
									"server_name": schema.StringAttribute{
										Computed:    true,
										Description: "服务名称",
									},
									"url_paths": schema.StringAttribute{
										Computed:    true,
										Description: "匹配路径",
									},
									"match_type": schema.StringAttribute{
										Computed:    true,
										Description: "匹配类型。取值范围：ABSOLUTE，PREFIX，REG",
										Validators: []validator.String{
											stringvalidator.OneOf(business.ElbRuleMatchTypes...),
										},
									},
								},
							},
						},
						"action_type": schema.StringAttribute{
							Computed:    true,
							Description: "默认规则动作类型",
						},
						"action_target_groups": schema.ListNestedAttribute{
							Computed:    true,
							Description: "后端服务组",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"target_group_id": schema.StringAttribute{
										Computed:    true,
										Description: "后端服务组ID",
									},
									"weight": schema.Int32Attribute{
										Computed:    true,
										Description: "权重",
									},
								},
							},
						},
						"action_redirect_listener_id": schema.StringAttribute{
							Computed:    true,
							Description: "重定向监听器ID",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "状态: ACTIVE / DOWN",
							Validators: []validator.String{
								stringvalidator.OneOf(business.ElbRuleStatus...),
							},
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
				},
			},
		},
	}
}

func (c *ctyunElbRules) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunElbRulesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		msg := "regionID不能为空"
		response.Diagnostics.AddError(msg, msg)
		return
	}

	params := &ctelb.CtelbListQueryRequest{
		RegionID: regionId,
	}

	if !config.IDs.IsNull() {
		params.IDs = config.IDs.ValueString()
	}
	if !config.LoadBalancerID.IsNull() {
		params.LoadBalancerID = config.LoadBalancerID.ValueString()
	}

	// 请求查看转发规则列表
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbListQueryApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回结果
	var elbRules []ElbRuelModel
	elbRuleList := resp.ReturnObj
	for _, rule := range elbRuleList {
		var elbRule ElbRuelModel
		elbRule.RegionID = types.StringValue(rule.RegionID)
		elbRule.AzName = types.StringValue(rule.AzName)
		elbRule.ProjectID = types.StringValue(rule.ProjectID)
		elbRule.ID = types.StringValue(rule.ID)
		elbRule.LoadBalancerID = types.StringValue(rule.LoadBalancerID)
		elbRule.Description = types.StringValue(rule.Description)
		elbRule.ListenerID = types.StringValue(rule.ListenerID)
		elbRule.Status = types.StringValue(rule.Status)
		elbRule.CreatedTime = types.StringValue(rule.CreatedTime)
		elbRule.UpdatedTime = types.StringValue(rule.UpdatedTime)

		// 解析Conditions
		var conditions []ConditionModel
		for _, conditionItem := range rule.Conditions {
			var condition ConditionModel
			condition.ConditionType = types.StringValue(conditionItem.RawType)
			condition.ServerName = types.StringValue(conditionItem.ServerNameConfig.ServerName)
			condition.UrlPaths = types.StringValue(conditionItem.UrlPathConfig.UrlPaths)
			condition.MatchType = types.StringValue(conditionItem.UrlPathConfig.MatchType)
			conditions = append(conditions, condition)
		}
		elbRule.Conditions = conditions
		// 解析Action
		elbRule.ActionType = types.StringValue(rule.Action.RawType)
		elbRule.ActionRedirectListenerID = types.StringValue(rule.Action.RedirectListenerID)
		var targetGroups []TargetGroupModel
		for _, targetItem := range rule.Action.ForwardConfig.TargetGroups {
			var targetGroup TargetGroupModel
			targetGroup.TargetGroupID = types.StringValue(targetItem.TargetGroupID)
			targetGroup.Weight = types.Int32Value(targetItem.Weight)
			targetGroups = append(targetGroups, targetGroup)
		}
		elbRule.ActionTargetGroups = targetGroups
		elbRules = append(elbRules, elbRule)
	}
	config.ElbRules = elbRules
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunElbRulesConfig struct {
	RegionID       types.String   `tfsdk:"region_id"`        //区域ID
	IDs            types.String   `tfsdk:"ids"`              //转发规则ID列表，以,分隔
	LoadBalancerID types.String   `tfsdk:"load_balancer_id"` //负载均衡实例ID
	ElbRules       []ElbRuelModel `tfsdk:"elb_rules"`
}

type ElbRuelModel struct {
	RegionID       types.String `tfsdk:"region_id"`        //区域ID
	AzName         types.String `tfsdk:"az_name"`          //可用区名称
	ProjectID      types.String `tfsdk:"project_id"`       //项目ID
	ID             types.String `tfsdk:"id"`               //转发规则ID
	LoadBalancerID types.String `tfsdk:"load_balancer_id"` //负载均衡ID
	ListenerID     types.String `tfsdk:"listener_id"`      //监听器ID
	//Priority                 types.Int32        `tfsdk:"priority"`                    //优先级
	Description              types.String       `tfsdk:"description"`                 //描述
	Conditions               []ConditionModel   `tfsdk:"conditions"`                  //匹配规则数据
	ActionType               types.String       `tfsdk:"action_type"`                 //默认规则动作类型
	ActionTargetGroups       []TargetGroupModel `tfsdk:"action_target_groups"`        //后端服务组
	ActionRedirectListenerID types.String       `tfsdk:"action_redirect_listener_id"` //重定向监听器ID
	Status                   types.String       `tfsdk:"status"`                      //状态: ACTIVE / DOWN
	CreatedTime              types.String       `tfsdk:"created_time"`                //创建时间，为UTC格式
	UpdatedTime              types.String       `tfsdk:"updated_time"`                //更新时间，为UTC格式
}

type ConditionModel struct {
	ConditionType types.String `tfsdk:"condition_type"` //类型。取值范围：server_name（服务名称）、url_path（匹配路径）
	ServerName    types.String `tfsdk:"server_name"`    //服务名称
	UrlPaths      types.String `tfsdk:"url_paths"`      //匹配路径
	MatchType     types.String `tfsdk:"match_type"`     //匹配类型。取值范围：ABSOLUTE，PREFIX，REG
}
