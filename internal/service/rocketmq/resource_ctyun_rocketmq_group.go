package rocketmq

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/rocketmq"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
)

var (
	_ resource.Resource                = &ctyunRocketmqGroup{}
	_ resource.ResourceWithConfigure   = &ctyunRocketmqGroup{}
	_ resource.ResourceWithImportState = &ctyunRocketmqGroup{}
)

type ctyunRocketmqGroup struct {
	meta *common.CtyunMetadata
	name string
}

func NewCtyunRocketmqGroup() resource.Resource {
	return &ctyunRocketmqGroup{}
}

func (c *ctyunRocketmqGroup) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rocketmq_group"
	c.name = response.TypeName
}

type CtyunRocketmqGroupConfig struct {
	ID                     types.String `tfsdk:"id"`
	InstanceID             types.String `tfsdk:"instance_id"`
	Name                   types.String `tfsdk:"name"`
	RegionID               types.String `tfsdk:"region_id"`
	ConsumeEnable          types.Bool   `tfsdk:"consume_enable"`           // 是否允许消费
	ConsumeBroadcastEnable types.Bool   `tfsdk:"consume_broadcast_enable"` // 是否允许广播消费
	FirstConsumeMechanism  types.Int32  `tfsdk:"first_consume_mechanism"`  // 首次消费机制
	PullMechanism          types.Int32  `tfsdk:"pull_mechanism"`           // 拉取机制
	RetryMaxTimes          types.Int32  `tfsdk:"retry_max_times"`          // 最大重试次数
	Remark                 types.String `tfsdk:"remark"`                   // 备注
	BrokerNames            types.Set    `tfsdk:"broker_names"`             // Broker 节点名称列表
}

func (c *ctyunRocketmqGroup) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理 RocketMQ 订阅组", "分布式消息服务 RocketMQ", "https://www.ctyun.cn/document/10000118/10001967"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "订阅组 ID",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "订阅组名称，需符合 MQ 订阅组命名规范（字母、数字、下划线组合，长度 1-64 字符），不可重复",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 64),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-zA-Z_-]+$`), "订阅组名称不符合规则"),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池 ID，如果不填则默认使用 provider ctyun 中的 region_id 或环境变量中的 CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "RocketMQ 实例 ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"consume_enable": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否允许消费，true 表示允许消费，false 表示禁止消费，默认 true 支持更新",
				Default:     booldefault.StaticBool(true),
			},
			"consume_broadcast_enable": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否开启广播消费，true 表示开启广播消费，false 表示关闭广播消费，默认 true",
				Default:     booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"first_consume_mechanism": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "首次消费机制，0:从最小位点开始消费，1:从最新位点开始消费，默认 0",
				Default:     int32default.StaticInt32(1),
				Validators: []validator.Int32{
					int32validator.OneOf(1, 2, 3),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"pull_mechanism": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "拉取机制，0:服务器主动推送，1:客户端拉取，默认 1",
				Default:     int32default.StaticInt32(1),
				Validators: []validator.Int32{
					int32validator.OneOf(1),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"retry_max_times": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "最大重试次数，消费失败后最多重试次数，默认 16",
				Default:     int32default.StaticInt32(16),
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "备注信息",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"broker_names": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "关联的 Broker 节点名称列表，需填写实际部署的 Broker 节点名称，不可为空数组",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (c *ctyunRocketmqGroup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunRocketmqGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.checkBeforeCreate(ctx, plan)
	if err != nil {
		return
	}
	err = c.create(ctx, &plan)
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

func (c *ctyunRocketmqGroup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRocketmqGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) || strings.Contains(err.Error(), "不存在") {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRocketmqGroup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunRocketmqGroupConfig
	var state CtyunRocketmqGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 更新可在线修改的字段
	err = c.update(ctx, &plan, &state)
	if err != nil {
		return
	}

	// 更新后重新读取资源信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunRocketmqGroup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRocketmqGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 销毁
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunRocketmqGroup) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunRocketmqGroup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例：%s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [name],[instance_id],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunRocketmqGroupConfig
	var name, instanceID, regionID string
	// 根据分隔符数量判断是否输入了 regionID
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &name, &instanceID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &name, &instanceID, &regionID)
		if err != nil {
			return
		}
	}
	if name == "" {
		err = fmt.Errorf("name 不能为空")
		return
	}
	if instanceID == "" {
		err = fmt.Errorf("instance_id 不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("region_id 不能为空")
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.Name = types.StringValue(name)
	cfg.InstanceID = types.StringValue(instanceID)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *ctyunRocketmqGroup) checkBeforeCreate(ctx context.Context, plan CtyunRocketmqGroupConfig) (err error) {
	name, instanceID, regionID := plan.Name.ValueString(), plan.InstanceID.ValueString(), plan.RegionID.ValueString()
	// 确保订阅组名称没有被占用
	exist, err := c.checkGroupExist(ctx, name, instanceID, regionID)
	if err != nil {
		return
	}
	if exist {
		return fmt.Errorf("rocketmq group %s already exist", name)
	}
	return
}

// create 创建
func (c *ctyunRocketmqGroup) create(ctx context.Context, plan *CtyunRocketmqGroupConfig) (err error) {
	// 获取 broker_names 列表
	var brokerNames []string
	plan.BrokerNames.ElementsAs(ctx, &brokerNames, false)

	params := &rocketmq.Mq2GroupCreateV3Request{
		RegionId:       plan.RegionID.ValueString(),
		ProdInstId:     plan.InstanceID.ValueString(),
		BrokerNameList: brokerNames,
		Remark:         plan.Remark.ValueString(),
		SubscriptionGroupConfig: &rocketmq.Mq2GroupCreateV3SubscriptionGroupConfigRequest{
			ConsumeEnable:          plan.ConsumeEnable.ValueBool(),
			FirstConsumeMechanism:  plan.FirstConsumeMechanism.ValueInt32(),
			GroupName:              plan.Name.ValueString(),
			PullMechanism:          plan.PullMechanism.ValueInt32(),
			ConsumeBroadcastEnable: plan.ConsumeBroadcastEnable.ValueBool(),
			RetryMaxTimes:          plan.RetryMaxTimes.ValueInt32(),
		},
	}

	resp, err := c.meta.Apis.RocketmqApis.Mq2GroupCreateV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		msg := ""
		if resp.Message != nil {
			msg = *resp.Message
		}
		err = fmt.Errorf("API return error. Message: %s", msg)
		return
	}
	return
}

// delete 删除
func (c *ctyunRocketmqGroup) delete(ctx context.Context, plan CtyunRocketmqGroupConfig) (err error) {
	params := &rocketmq.Mq2GroupDeleteV3Request{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
		GroupName:  plan.Name.ValueString(),
	}

	resp, err := c.meta.Apis.RocketmqApis.Mq2GroupDeleteV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		msg := ""
		if resp.Message != nil {
			msg = *resp.Message
		}
		err = fmt.Errorf("API return error. Message: %s", msg)
		return
	}
	return
}

// update 更新订阅组配置
func (c *ctyunRocketmqGroup) update(ctx context.Context, plan, state *CtyunRocketmqGroupConfig) (err error) {
	err = c.updateConsumeEnableIfNeeded(ctx, plan, state)
	if err != nil {
		return
	}

	err = c.checkConfigFields(*plan, *state)
	if err != nil {
		return
	}

	return
}

// updateConsumeEnableIfNeeded 如果 ConsumeEnable 字段变更则更新
func (c *ctyunRocketmqGroup) updateConsumeEnableIfNeeded(ctx context.Context, plan, state *CtyunRocketmqGroupConfig) (err error) {
	if plan.ConsumeEnable.Equal(state.ConsumeEnable) {
		return nil
	}

	return c.updateConsumeEnable(ctx, plan)
}

// checkConfigFields 检查配置字段是否变更，若变更则返回错误
func (c *ctyunRocketmqGroup) checkConfigFields(plan, state CtyunRocketmqGroupConfig) error {
	if plan.FirstConsumeMechanism.Equal(state.FirstConsumeMechanism) &&
		plan.PullMechanism.Equal(state.PullMechanism) &&
		plan.RetryMaxTimes.Equal(state.RetryMaxTimes) &&
		plan.ConsumeBroadcastEnable.Equal(state.ConsumeBroadcastEnable) {
		return nil
	}
	return fmt.Errorf("以下字段不支持更新: first_consume_mechanism, pull_mechanism, retry_max_times, consume_broadcast_enable")
}

// updateConsumeEnable 更新订阅组读取权限
func (c *ctyunRocketmqGroup) updateConsumeEnable(ctx context.Context, plan *CtyunRocketmqGroupConfig) (err error) {
	// 获取 broker_names 列表
	var brokerNames []string
	plan.BrokerNames.ElementsAs(ctx, &brokerNames, false)

	params := &rocketmq.Mq2GroupUpdatepermV3Request{
		RegionId:       plan.RegionID.ValueString(),
		ProdInstId:     plan.InstanceID.ValueString(),
		GroupName:      plan.Name.ValueString(),
		BrokerNameList: brokerNames,
		Enable:         plan.ConsumeEnable.ValueBool(),
		Remark:         plan.Remark.ValueString(),
	}

	resp, err := c.meta.Apis.RocketmqApis.Mq2GroupUpdatepermV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		msg := ""
		if resp.Message != nil {
			msg = *resp.Message
		}
		err = fmt.Errorf("API return error. Message: %s", msg)
		return
	}
	return
}

// checkGroupExist 检查订阅组是否存在
func (c *ctyunRocketmqGroup) checkGroupExist(ctx context.Context, name, instanceID, regionID string) (exist bool, err error) {
	params := &rocketmq.Mq2GroupListV3Request{
		RegionId:   regionID,
		ProdInstId: instanceID,
	}

	resp, err := c.meta.Apis.RocketmqApis.Mq2GroupListV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		msg := ""
		if resp.Message != nil {
			msg = *resp.Message
		}
		err = fmt.Errorf("API return error. Message: %s", msg)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	for _, group := range resp.ReturnObj.Rows {
		if group.GroupName == name {
			exist = true
			return
		}
	}
	return
}

// getGroupByName 根据名称查询订阅组
func (c *ctyunRocketmqGroup) getGroupByName(ctx context.Context, plan CtyunRocketmqGroupConfig) (group *rocketmq.Mq2GroupListV3ReturnObjRowsResponse, err error) {
	params := &rocketmq.Mq2GroupListV3Request{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
	}

	resp, err := c.meta.Apis.RocketmqApis.Mq2GroupListV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		msg := ""
		if resp.Message != nil {
			msg = *resp.Message
		}
		err = fmt.Errorf("API return error. Message: %s", msg)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if len(resp.ReturnObj.Rows) == 0 {
		err = common.ResourceNotExistError
		return
	}

	// 查找匹配的订阅组
	for _, g := range resp.ReturnObj.Rows {
		if g.GroupName == plan.Name.ValueString() {
			group = g
			return
		}
	}

	err = common.ResourceNotExistError
	return
}

// getAndMerge 从远端查询并合并数据
func (c *ctyunRocketmqGroup) getAndMerge(ctx context.Context, plan *CtyunRocketmqGroupConfig) (err error) {
	group, err := c.getGroupByName(ctx, *plan)
	if err != nil {
		return
	}

	plan.Name = types.StringValue(group.GroupName)
	plan.ConsumeEnable = types.BoolValue(group.ConsumeEnable)
	plan.FirstConsumeMechanism = types.Int32Value(group.FirstConsumeMechanism)
	plan.PullMechanism = types.Int32Value(group.PullMechanism)
	plan.RetryMaxTimes = types.Int32Value(group.RetryMaxTimes)

	plan.ConsumeBroadcastEnable = types.BoolValue(group.ConsumeBroadcastEnable)
	plan.Remark = types.StringValue(group.Remark)

	brokerNames := strings.Split(group.BrokerName, ",")
	brokerNamesValue := make([]attr.Value, len(brokerNames))
	for i, name := range brokerNames {
		brokerNamesValue[i] = types.StringValue(strings.TrimSpace(name))
	}
	plan.BrokerNames = types.SetValueMust(types.StringType, brokerNamesValue)

	// 设置 ID
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s",
		plan.Name.ValueString(),
		plan.InstanceID.ValueString(),
		plan.RegionID.ValueString()))
	return
}
