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
	_ resource.Resource                = &ctyunRocketmqTopic{}
	_ resource.ResourceWithConfigure   = &ctyunRocketmqTopic{}
	_ resource.ResourceWithImportState = &ctyunRocketmqTopic{}
)

type ctyunRocketmqTopic struct {
	meta *common.CtyunMetadata
	name string
}

func NewCtyunRocketmqTopic() resource.Resource {
	return &ctyunRocketmqTopic{}
}

func (c *ctyunRocketmqTopic) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rocketmq_topic"
	c.name = response.TypeName
}

type CtyunRocketmqTopicConfig struct {
	ID          types.String `tfsdk:"id"`
	InstanceID  types.String `tfsdk:"instance_id"`
	Name        types.String `tfsdk:"name"`
	RegionID    types.String `tfsdk:"region_id"`
	QueueNums   types.Int32  `tfsdk:"queue_nums"`   // 写队列数量
	Order       types.Bool   `tfsdk:"order"`        // 是否为顺序消息
	Perm        types.Int32  `tfsdk:"perm"`         // 权限控制值
	Remark      types.String `tfsdk:"remark"`       // 备注
	BrokerNames types.Set    `tfsdk:"broker_names"` // Broker 节点名称列表
}

func (c *ctyunRocketmqTopic) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理 RocketMQ 主题", "分布式消息服务 RocketMQ", "https://www.ctyun.cn/document/10000118/10001967"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "主题 ID",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "主题名称，需符合 MQ 主题命名规范（字母、数字、下划线组合，长度 1-64 字符），不可重复",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 64),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-zA-Z_-]+$`), "主题名称不符合规则"),
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
			"queue_nums": schema.Int32Attribute{
				Required:    true,
				Description: "主题的读写队列数量，用于控制消息读写入的并发能力，需为正整数",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"order": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "标识主题是否为顺序消息队列，true 表示顺序队列，false 表示普通队列，默认 false",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"perm": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "主题的权限控制值，固定可选值为 2（只读）、4（只写）、6（读写），默认推荐 6 支持更新",
				Default:     int32default.StaticInt32(6),
				Validators: []validator.Int32{
					int32validator.OneOf(2, 4, 6),
				},
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "备注信息 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
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

func (c *ctyunRocketmqTopic) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunRocketmqTopicConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkBeforeCreate(ctx, plan)
	if err != nil {
		return
	}
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

func (c *ctyunRocketmqTopic) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRocketmqTopicConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRocketmqTopic) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// RocketMQ Topic 仅 perm remark 字段支持在线更新
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunRocketmqTopicConfig
	var state CtyunRocketmqTopicConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 检查是否只有 perm 字段变更
	permChanged := !plan.Perm.Equal(state.Perm) || !plan.Remark.Equal(state.Remark)

	if permChanged {
		err = c.updatePerm(ctx, &plan)
		if err != nil {
			return
		}
	}

	// 更新后重新读取资源信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunRocketmqTopic) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRocketmqTopicConfig
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

func (c *ctyunRocketmqTopic) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunRocketmqTopic) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例：%s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [instance_id],[name],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunRocketmqTopicConfig
	var name, instanceID, regionID string
	// 根据分隔符数量判断是否输入了 regionID
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &instanceID, &name)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &instanceID, &name, &regionID)
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

func (c *ctyunRocketmqTopic) checkBeforeCreate(ctx context.Context, plan CtyunRocketmqTopicConfig) (err error) {
	name, instanceID, regionID := plan.Name.ValueString(), plan.InstanceID.ValueString(), plan.RegionID.ValueString()
	// 确保主题名称没有被占用
	exist, err := c.checkTopicExist(ctx, name, instanceID, regionID)
	if err != nil {
		return
	}
	if exist {
		return fmt.Errorf("rocketmq topic %s already exist", name)
	}
	return
}

// create 创建
func (c *ctyunRocketmqTopic) create(ctx context.Context, plan CtyunRocketmqTopicConfig) (err error) {
	// 获取 broker_names 列表
	var brokerNames []string
	plan.BrokerNames.ElementsAs(ctx, &brokerNames, false)

	params := &rocketmq.Mq2TopicCreateV3Request{
		RegionId:             plan.RegionID.ValueString(),
		ProdInstId:           plan.InstanceID.ValueString(),
		BrokerNameList:       brokerNames,
		WriteQueueNums:       plan.QueueNums.ValueInt32(),
		ReadQueueNums:        plan.QueueNums.ValueInt32(),
		Order:                plan.Order.ValueBool(),
		Perm:                 plan.Perm.ValueInt32(),
		AllowdConsumerGroups: []string{}, // 空数组表示不限制消费者组订阅
		TopicName:            plan.Name.ValueString(),
	}
	if !plan.Remark.IsUnknown() && !plan.Remark.IsNull() && plan.Remark.ValueString() != "" {
		params.Remark = plan.Remark.ValueString()
	}
	resp, err := c.meta.Apis.RocketmqApis.Mq2TopicCreateV3Api.Do(ctx, c.meta.SdkCredential, params)
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
func (c *ctyunRocketmqTopic) delete(ctx context.Context, plan CtyunRocketmqTopicConfig) (err error) {
	params := &rocketmq.Mq2TopicDeleteV3Request{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
		TopicName:  plan.Name.ValueString(),
	}

	resp, err := c.meta.Apis.RocketmqApis.Mq2TopicDeleteV3Api.Do(ctx, c.meta.SdkCredential, params)
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

// updatePerm 更新 Topic的读写模式
func (c *ctyunRocketmqTopic) updatePerm(ctx context.Context, plan *CtyunRocketmqTopicConfig) (err error) {
	params := &rocketmq.Mq2TopicUpdateRequest{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
		Topic:      plan.Name.ValueString(),
		Perm:       plan.Perm.ValueInt32(),
		Remark:     plan.Remark.ValueString(),
	}

	resp, err := c.meta.Apis.RocketmqApis.Mq2TopicUpdateApi.Do(ctx, c.meta.SdkCredential, params)
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

// checkTopicExist 检查主题是否存在
func (c *ctyunRocketmqTopic) checkTopicExist(ctx context.Context, name, instanceID, regionID string) (exist bool, err error) {
	params := &rocketmq.Mq2TopicListV3Request{
		RegionId:   regionID,
		ProdInstId: instanceID,
		TopicName:  name,
	}

	resp, err := c.meta.Apis.RocketmqApis.Mq2TopicListV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		msg := utils.SecString(resp.Message)
		err = fmt.Errorf("API return error. Message: %s", msg)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	for _, topic := range resp.ReturnObj.Rows {
		if utils.SecString(topic.TopicName) == name {
			exist = true
			return
		}
	}
	return
}

// getTopicByName 根据名称查询主题
func (c *ctyunRocketmqTopic) getTopicByName(ctx context.Context, plan CtyunRocketmqTopicConfig) (topic *rocketmq.Mq2TopicListV3ReturnObjRowsResponse, err error) {
	params := &rocketmq.Mq2TopicListV3Request{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
		TopicName:  plan.Name.ValueString(),
	}

	resp, err := c.meta.Apis.RocketmqApis.Mq2TopicListV3Api.Do(ctx, c.meta.SdkCredential, params)
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
	topic = resp.ReturnObj.Rows[0]
	return
}

// getAndMerge 从远端查询
func (c *ctyunRocketmqTopic) getAndMerge(ctx context.Context, plan *CtyunRocketmqTopicConfig) (err error) {
	topic, err := c.getTopicByName(ctx, *plan)
	if err != nil {
		return
	}

	plan.Name = utils.SecStringValue(topic.TopicName)
	plan.QueueNums = types.Int32PointerValue(topic.ReadQueueNums)
	plan.Order = types.BoolPointerValue(topic.Order)
	plan.Perm = types.Int32PointerValue(topic.Perm)
	plan.Remark = types.StringValue(topic.Remark)

	// BrokerNames 已经在创建时设置，且 RequiresReplace，直接使用原有值

	brokerNames := strings.Split(*topic.BrokerName, ",")
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
