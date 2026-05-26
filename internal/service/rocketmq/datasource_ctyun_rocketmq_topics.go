package rocketmq

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/rocketmq"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
)

var (
	_ datasource.DataSource              = &ctyunRocketmqTopics{}
	_ datasource.DataSourceWithConfigure = &ctyunRocketmqTopics{}
)

type ctyunRocketmqTopics struct {
	meta *common.CtyunMetadata
}

func NewCtyunRocketmqTopics() datasource.DataSource {
	return &ctyunRocketmqTopics{}
}

func (c *ctyunRocketmqTopics) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rocketmq_topics"
}

type CtyunRocketmqTopicModel struct {
	Name           types.String `tfsdk:"name"`             // 主题名称
	ReadQueueNums  types.Int32  `tfsdk:"read_queue_nums"`  // 读队列数量
	WriteQueueNums types.Int32  `tfsdk:"write_queue_nums"` // 写队列数量
	Perm           types.Int32  `tfsdk:"perm"`             // 权限控制值
	Order          types.Bool   `tfsdk:"order"`            // 是否为顺序消息
	Remark         types.String `tfsdk:"remark"`           // 备注
	BrokerName     types.String `tfsdk:"broker_name"`      // Broker 名称
	MessageType    types.String `tfsdk:"message_type"`     // 消息类型
}

type CtyunRocketmqTopicsConfig struct {
	RegionID   types.String              `tfsdk:"region_id"`
	InstanceID types.String              `tfsdk:"instance_id"`
	Name       types.String              `tfsdk:"name"`
	Topics     []CtyunRocketmqTopicModel `tfsdk:"topics"`
}

func (c *ctyunRocketmqTopics) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("查询 RocketMQ 主题列表", "分布式消息服务 RocketMQ", "https://www.ctyun.cn/document/10000118/10001967"),
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池 ID",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "RocketMQ 实例 ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "主题名称，支持模糊查询",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 64),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-zA-Z_-]+$`), "主题名称不符合规则"),
				},
			},
			"topics": schema.ListNestedAttribute{
				Description: "主题列表",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "主题名称",
							Computed:    true,
						},
						"read_queue_nums": schema.Int32Attribute{
							Description: "读队列数量",
							Computed:    true,
						},
						"write_queue_nums": schema.Int32Attribute{
							Description: "写队列数量",
							Computed:    true,
						},
						"perm": schema.Int32Attribute{
							Description: "权限控制值（2=只读，4=只写，6=读写）",
							Computed:    true,
						},
						"order": schema.BoolAttribute{
							Description: "是否为顺序消息",
							Computed:    true,
						},
						"remark": schema.StringAttribute{
							Description: "备注信息",
							Computed:    true,
						},
						"broker_name": schema.StringAttribute{
							Description: "Broker 名称",
							Computed:    true,
						},
						"message_type": schema.StringAttribute{
							Description: "消息类型（NORMAL-普通消息）",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRocketmqTopics) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRocketmqTopicsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("region_id 不能为空")
		return
	}

	config.RegionID = types.StringValue(regionId)
	// 组装请求体
	params := &rocketmq.Mq2TopicListV3Request{
		RegionId:   regionId,
		ProdInstId: config.InstanceID.ValueString(),
	}
	if config.Name.ValueString() != "" {
		params.TopicName = config.Name.ValueString()
	}
	// 调用 API
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
	}

	config.Topics = []CtyunRocketmqTopicModel{}
	// 解析返回值
	for _, e := range resp.ReturnObj.Rows {
		item := CtyunRocketmqTopicModel{}
		item.Name = utils.SecStringValue(e.TopicName)
		item.ReadQueueNums = types.Int32PointerValue(e.ReadQueueNums)
		item.WriteQueueNums = types.Int32PointerValue(e.WriteQueueNums)
		item.Perm = types.Int32PointerValue(e.Perm)
		item.Order = utils.SecBoolValue(e.Order)
		item.Remark = types.StringValue(e.Remark)
		item.BrokerName = types.StringPointerValue(e.BrokerName)
		item.MessageType = types.StringPointerValue(e.MessageType)
		config.Topics = append(config.Topics, item)
	}

	// 保存到 state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRocketmqTopics) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
