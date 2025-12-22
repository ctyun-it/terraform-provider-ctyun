package kafka

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgkafka "github.com/ctyun-it/terraform-provider-ctyun/internal/core/kafka"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"regexp"
	"strings"
)

var (
	_ resource.Resource                = &ctyunKafkaTopic{}
	_ resource.ResourceWithConfigure   = &ctyunKafkaTopic{}
	_ resource.ResourceWithImportState = &ctyunKafkaTopic{}
)

type ctyunKafkaTopic struct {
	meta       *common.CtyunMetadata
	vpcService *business.VpcService
	sgService  *business.SecurityGroupService
}

func NewCtyunKafkaTopic() resource.Resource {
	return &ctyunKafkaTopic{}
}

func (c *ctyunKafkaTopic) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_kafka_topic"
}

type CtyunKafkaTopicConfig struct {
	Id                          types.String `tfsdk:"id"`
	TopicName                   types.String `tfsdk:"name"`
	InstanceId                  types.String `tfsdk:"instance_id"`
	RegionId                    types.String `tfsdk:"region_id"`
	PartitionNum                types.Int32  `tfsdk:"partition_num"`
	FactorNum                   types.Int32  `tfsdk:"factor_num"`
	PartitionCapacity           types.Int32  `tfsdk:"partition_capacity"`
	RetentionTime               types.Int32  `tfsdk:"retention_time"`
	MinReplicas                 types.Int32  `tfsdk:"min_replicas"`
	MaxMessage                  types.Int32  `tfsdk:"max_message"`
	NeedFlush                   types.Bool   `tfsdk:"need_flush"`
	TimestampType               types.String `tfsdk:"timestamp_type"`
	Description                 types.String `tfsdk:"description"`
	CleanupPolicy               types.String `tfsdk:"cleanup_policy"`
	UncleanLeaderElectionEnable types.Bool   `tfsdk:"unclean_leader_election_enable"`
	SegmentMs                   types.Int64  `tfsdk:"segment_ms"`
	SegmentBytes                types.Int64  `tfsdk:"segment_bytes"`
	RemoteStorageEnable         types.Bool   `tfsdk:"remote_storage_enable"`
	LocalRetentionMs            types.Int64  `tfsdk:"local_retention_ms"`

	PartitionList   types.List `tfsdk:"partition_list"`   // 分区详情列表
	GroupSubscribed types.List `tfsdk:"group_subscribed"` // 订阅该主题的消费组列表

}

type CtyunKafkaTopicPartitionDetail struct {
	TopicName   types.String                                                                  `tfsdk:"topic_name"`
	PartitionId types.Int32                                                                   `tfsdk:"partition_id"`
	Offsets     *ctgkafka.CtgkafkaGetTopicDetailsReturnObjDataPartitionListOffsetsResponse    `tfsdk:"offsets"`
	Replicas    []*ctgkafka.CtgkafkaGetTopicDetailsReturnObjDataPartitionListReplicasResponse `tfsdk:"replicas"`
}

var (
	offsetAttrTypes = map[string]attr.Type{
		"total":            types.Int64Type,
		"begin":            types.Int64Type,
		"end":              types.Int64Type,
		"update_timestamp": types.Int64Type,
		"hw":               types.Int64Type,
	}

	replicaAttrTypes = map[string]attr.Type{
		"broker_id": types.Int32Type,
		"is_leader": types.BoolType,
		"in_sync":   types.BoolType,
		"size":      types.Int64Type,
		"lag":       types.Int64Type,
	}

	partitionDetailAttrTypes = map[string]attr.Type{
		"topic_name":   types.StringType,
		"partition_id": types.Int32Type,
		"offsets":      types.ObjectType{AttrTypes: offsetAttrTypes},
		"replicas":     types.ListType{ElemType: types.ObjectType{AttrTypes: replicaAttrTypes}},
	}
)

func (c *ctyunKafkaTopic) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029624/10144604`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "资源唯一标识符",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "主题名称，英文字母、数字、下划线开头，且只能由英文字母、数字、中划线、下划线组成，长度为3-64个字符。",
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
			"partition_num": schema.Int32Attribute{
				Required:    true,
				Description: "分区数，取值范围[1, min(100, 实例剩余分区数量)]，实例剩余分区数量=实例分区上限-所有主题分区数之和。支持更新",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"factor_num": schema.Int32Attribute{
				Optional:    true,
				Description: "副本数，取值范围[1, 3]，单机版默认值1，集群版默认值3。",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"partition_capacity": schema.Int32Attribute{
				Optional:    true,
				Description: "分区容量限制，单位GB，取值-1或范围[1, 100]。-1表示无限制，默认值-1。支持更新",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"retention_time": schema.Int32Attribute{
				Optional:    true,
				Description: "消息保留时长，单位毫秒，取值-1或范围[3600000, 315360000000]，单位毫秒，-1表示永久保留。 默认值259200000。支持更新",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"min_replicas": schema.Int32Attribute{
				Optional:    true,
				Description: "最小同步副本数，需小于等于factorNum，单机版默认值1，集群版默认值min(2, factorNum)。支持更新",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"max_message": schema.Int32Attribute{
				Optional:    true,
				Description: "最大消息大小，单位字节，取值范围[1, 10485760]， 默认值1048588。支持更新",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"need_flush": schema.BoolAttribute{
				Optional:    true,
				Description: "是否同步刷盘。<br><li>true：是<br><li>false：否<br><li>默认值false 支持更新",
			},
			"timestamp_type": schema.StringAttribute{
				Optional:    true,
				Description: "消息时间戳类型。<br><li>CreateTime<br><li>LogAppendTime<br><li>默认值CreateTime 支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf("CreateTime", "LogAppendTime"),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "主题描述，规则如下：<br><li>不能以+,-,@,= 特殊字符开头。 <br><li>长度不能大于200。支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 200),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[^+\-@=].*$`),
						"不能以+,-,@,=特殊字符开头",
					),
				},
			},
			"cleanup_policy": schema.StringAttribute{
				Optional:    true,
				Description: "日志保留策略。<br><li>delete<br><li>compact<br><li>默认值：delete。支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf("delete", "compact"),
				},
			},
			"unclean_leader_election_enable": schema.BoolAttribute{
				Optional:    true,
				Description: "是否允许不同步的副本参与leader选举。<br><li>false<br><li>true<br><li>默认值：false。支持更新",
			},
			"segment_ms": schema.Int64Attribute{
				Optional:    true,
				Description: "日志滚动时间，单位ms。 取值范围[86400000, 7776000000]，默认值：259200000 支持更新",
				Validators: []validator.Int64{
					int64validator.Between(86400000, 7776000000),
				},
			},
			"segment_bytes": schema.Int64Attribute{
				Optional:    true,
				Description: "分片大小，单位byte。 取值范围[268435456, 10737418240]，默认值：1073741824 支持更新",
				Validators: []validator.Int64{
					int64validator.Between(268435456, 10737418240),
				},
			},
			"remote_storage_enable": schema.BoolAttribute{
				Optional:    true,
				Description: "是否开启对象存储。<br><li>true：是<br><li>false：否<br><li>默认值false 支持更新",
			},
			"local_retention_ms": schema.Int64Attribute{
				Optional:    true,
				Description: "本地保留时长，单位ms。 取值范围[180000, 315360000000] 支持更新",
				Validators: []validator.Int64{
					int64validator.Between(180000, 315360000000),
				},
			},

			"group_subscribed": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "订阅该主题的消费组名称列表",
			},

			"partition_list": schema.ListNestedAttribute{
				Description: "分区详情列表",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"topic_name": schema.StringAttribute{
							Computed:    true,
							Description: "主题名称",
						},
						"partition_id": schema.Int32Attribute{
							Computed:    true,
							Description: "分区ID",
						},

						"offsets": schema.SingleNestedAttribute{
							Description: "分区偏移量信息",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"total": schema.Int64Attribute{
									Computed:    true,
									Description: "分区消息总数",
								},
								"begin": schema.Int64Attribute{
									Computed:    true,
									Description: "分区leader副本的最大偏移量",
								},
								"end": schema.Int64Attribute{
									Computed:    true,
									Description: "分区leader副本的最小偏移量",
								},
								"update_timestamp": schema.Int64Attribute{
									Computed:    true,
									Description: "分区最近写入消息的毫秒时间戳",
								},
								"hw": schema.Int64Attribute{
									Computed:    true,
									Description: "分区消息高水位线，所有副本均已确认写入的最大偏移量",
								},
							},
						},
						"replicas": schema.ListNestedAttribute{
							Description: "副本信息",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"broker_id": schema.Int32Attribute{
										Computed:    true,
										Description: "Broker节点ID",
									},
									"is_leader": schema.BoolAttribute{
										Computed:    true,
										Description: "是否是主副本",
									},
									"in_sync": schema.BoolAttribute{
										Computed:    true,
										Description: "副本是否处于同步状态",
									},
									"size": schema.Int64Attribute{
										Computed:    true,
										Description: "副本消息大小，单位字节",
									},
									"lag": schema.Int64Attribute{
										Computed:    true,
										Description: "该副本当前落后hw的消息数",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunKafkaTopic) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunKafkaTopicConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
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

func (c *ctyunKafkaTopic) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunKafkaTopicConfig
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

func (c *ctyunKafkaTopic) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunKafkaTopicConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunKafkaTopicConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 更新
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunKafkaTopic) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunKafkaTopicConfig
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

func (c *ctyunKafkaTopic) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.vpcService = business.NewVpcService(meta)
	c.sgService = business.NewSecurityGroupService(meta)
}

func (c *ctyunKafkaTopic) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [instanceId],[topicName],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunKafkaTopicConfig
	var instanceId, regionID, topicName string

	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) == 1 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &instanceId, &topicName)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &instanceId, &topicName, &regionID)
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
	if topicName == "" {
		err = fmt.Errorf("topicName不能为空")
		return
	}

	cfg.RegionId = types.StringValue(regionID)
	cfg.InstanceId = types.StringValue(instanceId)
	cfg.TopicName = types.StringValue(topicName)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// create 创建
func (c *ctyunKafkaTopic) create(ctx context.Context, plan CtyunKafkaTopicConfig) (err error) {
	params := &ctgkafka.CtgkafkaTopicCreateV3Request{
		RegionId:                    plan.RegionId.ValueString(),
		ProdInstId:                  plan.InstanceId.ValueString(),
		TopicName:                   plan.TopicName.ValueString(),
		PartitionNum:                plan.PartitionNum.ValueInt32(),
		FactorNum:                   plan.FactorNum.ValueInt32(),
		PartitionCapacity:           plan.PartitionCapacity.ValueInt32(),
		RetentionTime:               plan.RetentionTime.ValueInt32(),
		MinReplicas:                 plan.MinReplicas.ValueInt32(),
		MaxMessage:                  plan.MaxMessage.ValueInt32(),
		NeedFlush:                   plan.NeedFlush.ValueBoolPointer(),
		TimestampType:               plan.TimestampType.ValueString(),
		Description:                 plan.Description.ValueString(),
		CleanupPolicy:               plan.CleanupPolicy.ValueString(),
		UncleanLeaderElectionEnable: plan.UncleanLeaderElectionEnable.ValueBoolPointer(),
		SegmentMs:                   plan.SegmentMs.ValueInt64(),
		SegmentBytes:                plan.SegmentBytes.ValueInt64(),
		RemoteStorageEnable:         plan.RemoteStorageEnable.ValueBoolPointer(),
		LocalRetentionMs:            plan.LocalRetentionMs.ValueInt64(),
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaTopicCreateV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return
}

// update 更新
func (c *ctyunKafkaTopic) update(ctx context.Context, plan, state CtyunKafkaTopicConfig) (err error) {
	params := &ctgkafka.CtgkafkaUpdateTopicRequest{
		RegionId:   state.RegionId.ValueString(),
		ProdInstId: state.InstanceId.ValueString(),
		TopicName:  plan.TopicName.ValueString(),
	}

	// 只有当字段有值时才设置
	if plan.PartitionNum.ValueInt32() != state.PartitionNum.ValueInt32() {
		params.PartitionNum = plan.PartitionNum.ValueInt32()
	}
	if !plan.PartitionCapacity.IsNull() {
		params.PartitionCapacity = plan.PartitionCapacity.ValueInt32()
	}
	if !plan.RetentionTime.IsNull() {
		params.RetentionTime = plan.RetentionTime.ValueInt32()
	}
	if !plan.MaxMessage.IsNull() {
		params.MaxMessage = plan.MaxMessage.ValueInt32()
	}
	if !plan.TimestampType.IsNull() {
		params.TimestampType = plan.TimestampType.ValueString()
	}
	if !plan.CleanupPolicy.IsNull() {
		params.CleanupPolicy = plan.CleanupPolicy.ValueString()
	}
	if !plan.UncleanLeaderElectionEnable.IsNull() {
		params.UncleanLeaderElectionEnable = plan.UncleanLeaderElectionEnable.ValueBoolPointer()
	}
	if !plan.SegmentMs.IsNull() {
		params.SegmentMs = plan.SegmentMs.ValueInt64()
	}
	if !plan.SegmentBytes.IsNull() {
		params.SegmentBytes = plan.SegmentBytes.ValueInt64()
	}
	if !plan.RemoteStorageEnable.IsNull() {
		params.RemoteStorageEnable = plan.RemoteStorageEnable.ValueBoolPointer()
	}
	if !plan.LocalRetentionMs.IsNull() {
		params.LocalRetentionMs = plan.LocalRetentionMs.ValueInt64()
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaUpdateTopicApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		return fmt.Errorf("API return error. Message: %s", resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	} else if resp.ReturnObj.Data != "update success" {
		return fmt.Errorf("API return error. Data: %s", resp.ReturnObj.Data)
	}
	return
}

// destroy 销毁
func (c *ctyunKafkaTopic) destroy(ctx context.Context, plan CtyunKafkaTopicConfig) (err error) {
	params := &ctgkafka.CtgkafkaTopicDeleteV3Request{
		RegionId:   plan.RegionId.ValueString(),
		ProdInstId: plan.InstanceId.ValueString(),
		TopicName:  plan.TopicName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaTopicDeleteV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// getAndMerge 从远端查询
func (c *ctyunKafkaTopic) getAndMerge(ctx context.Context, plan *CtyunKafkaTopicConfig) (err error) {
	params := &ctgkafka.CtgkafkaGetTopicDetailsRequest{
		RegionId:   plan.RegionId.ValueString(),
		ProdInstId: plan.InstanceId.ValueString(),
		TopicName:  plan.TopicName.ValueString(),
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaGetTopicDetailsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != nil && *resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil || resp.ReturnObj.Data == nil {
		err = common.InvalidReturnObjError
		return
	}

	topicData := resp.ReturnObj.Data

	// 设置ID（使用组合键）
	id := fmt.Sprintf("%s,%s,%s", plan.InstanceId.ValueString(), plan.RegionId.ValueString(), plan.TopicName.ValueString())
	plan.Id = types.StringValue(id)

	// 设置主题名称
	plan.TopicName = types.StringValue(*topicData.TopicName)
	// 设置分区数量
	plan.PartitionNum = types.Int32Value(int32(len(topicData.PartitionList)))
	// 设置订阅该主题的消费组列表
	if len(topicData.GroupSubscribed) > 0 {
		groupSubscribedStrs := make([]attr.Value, len(topicData.GroupSubscribed))
		for i, group := range topicData.GroupSubscribed {
			groupSubscribedStrs[i] = types.StringValue(*group)
		}
		plan.GroupSubscribed, _ = types.ListValue(types.StringType, groupSubscribedStrs)
	} else {
		// 如果没有订阅的消费组，设置为空列表
		plan.GroupSubscribed, _ = types.ListValue(types.StringType, []attr.Value{})
	}

	// 设置分区详情列表
	if len(topicData.PartitionList) > 0 {
		partitionDetailObjects := make([]attr.Value, len(topicData.PartitionList))
		for i, partition := range topicData.PartitionList {
			// 处理偏移量信息
			var offsetsObject basetypes.ObjectValue
			if partition.Offsets != nil {
				offsetsAttrs := map[string]attr.Value{
					"total":            types.Int64Value(partition.Offsets.Total),
					"begin":            types.Int64Value(partition.Offsets.Begin),
					"end":              types.Int64Value(partition.Offsets.End),
					"update_timestamp": types.Int64Value(partition.Offsets.UpdateTime),
					"hw":               types.Int64Value(partition.Offsets.Hw),
				}
				offsetsObject, _ = types.ObjectValue(offsetAttrTypes, offsetsAttrs)
			} else {
				// 如果没有偏移量信息，创建空对象
				offsetsObject = types.ObjectNull(offsetAttrTypes)
			}

			// 处理副本信息
			var replicaList basetypes.ListValue
			if len(partition.Replicas) > 0 {
				replicaObjects := make([]attr.Value, len(partition.Replicas))
				for j, replica := range partition.Replicas {
					replicaAttrs := map[string]attr.Value{
						"broker_id": types.Int32Value(replica.BrokerId),
						"is_leader": types.BoolPointerValue(replica.IsLeader),
						"in_sync":   types.BoolPointerValue(replica.InSync),
						"size":      types.Int64Value(replica.Size),
						"lag":       types.Int64Value(replica.Lag),
					}
					replicaObject, _ := types.ObjectValue(replicaAttrTypes, replicaAttrs)
					replicaObjects[j] = replicaObject
				}
				replicaList, _ = types.ListValue(types.ObjectType{AttrTypes: replicaAttrTypes}, replicaObjects)
			} else {
				// 如果没有副本信息，创建空列表
				replicaList = types.ListNull(types.ObjectType{AttrTypes: replicaAttrTypes})
			}

			// 构建分区详情对象
			partitionDetailAttrs := map[string]attr.Value{
				"topic_name":   types.StringValue(*partition.TopicName),
				"partition_id": types.Int32Value(partition.PartitionId),
				"offsets":      offsetsObject,
				"replicas":     replicaList,
			}
			partitionDetailObject, _ := types.ObjectValue(partitionDetailAttrTypes, partitionDetailAttrs)
			partitionDetailObjects[i] = partitionDetailObject
		}

		// 设置分区详情列表
		plan.PartitionList, _ = types.ListValue(
			types.ObjectType{AttrTypes: partitionDetailAttrTypes},
			partitionDetailObjects,
		)
	} else {
		// 如果没有分区信息，设置为空列表
		plan.PartitionList, _ = types.ListValue(
			types.ObjectType{AttrTypes: partitionDetailAttrTypes},
			[]attr.Value{},
		)
	}

	return
}
