package rocketmq

import (
	"context"
	"fmt"
	"regexp"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/rocketmq"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunRocketmqGroups{}
	_ datasource.DataSourceWithConfigure = &ctyunRocketmqGroups{}
)

type ctyunRocketmqGroups struct {
	meta *common.CtyunMetadata
}

func NewCtyunRocketmqGroups() datasource.DataSource {
	return &ctyunRocketmqGroups{}
}

func (c *ctyunRocketmqGroups) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rocketmq_groups"
}

type CtyunRocketmqGroupModel struct {
	Name                           types.String `tfsdk:"name"`                               // 订阅组名称
	ConsumeEnable                  types.Bool   `tfsdk:"consume_enable"`                     // 是否允许消费
	FirstConsumeMechanism          types.Int32  `tfsdk:"first_consume_mechanism"`            // 首次消费机制
	PullMechanism                  types.Int32  `tfsdk:"pull_mechanism"`                     // 拉取机制
	RetryQueueNums                 types.Int32  `tfsdk:"retry_queue_nums"`                   // 重试队列数量
	RetryMaxTimes                  types.Int32  `tfsdk:"retry_max_times"`                    // 最大重试次数
	BrokerId                       types.Int32  `tfsdk:"broker_id"`                          // Broker ID
	WhichBrokerWhenConsumeSlowly   types.Int32  `tfsdk:"which_broker_when_consume_slowly"`   // 消费慢时指定的 Broker ID
	NotifyConsumerIdsChangedEnable types.Bool   `tfsdk:"notify_consumer_ids_changed_enable"` // 是否开启消费者 ID 变更通知
	Remark                         types.String `tfsdk:"remark"`                             // 备注
	BrokerName                     types.String `tfsdk:"broker_name"`                        // Broker 名称
}

type CtyunRocketmqGroupsConfig struct {
	RegionID   types.String              `tfsdk:"region_id"`
	InstanceID types.String              `tfsdk:"instance_id"`
	Name       types.String              `tfsdk:"name"`
	Groups     []CtyunRocketmqGroupModel `tfsdk:"groups"`
}

func (c *ctyunRocketmqGroups) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("查询 RocketMQ 订阅组列表", "分布式消息服务RocketMQ", "https://www.ctyun.cn/document/10000114/10143635"),
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
				Description: "订阅组名称，支持模糊查询",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 64),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-zA-Z_-]+$`), "订阅组名称不符合规则"),
				},
			},
			"groups": schema.ListNestedAttribute{
				Description: "订阅组列表",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "订阅组名称",
							Computed:    true,
						},
						"consume_enable": schema.BoolAttribute{
							Description: "是否允许消费",
							Computed:    true,
						},
						"first_consume_mechanism": schema.Int32Attribute{
							Description: "首次消费机制（0:从最小位点开始消费，1:从最新位点开始消费）",
							Computed:    true,
						},
						"pull_mechanism": schema.Int32Attribute{
							Description: "拉取机制（0:服务器主动推送，1:客户端拉取）",
							Computed:    true,
						},
						"retry_queue_nums": schema.Int32Attribute{
							Description: "重试队列数量",
							Computed:    true,
						},
						"retry_max_times": schema.Int32Attribute{
							Description: "最大重试次数",
							Computed:    true,
						},
						"broker_id": schema.Int32Attribute{
							Description: "Broker ID",
							Computed:    true,
						},
						"which_broker_when_consume_slowly": schema.Int32Attribute{
							Description: "消费慢时指定的 Broker ID",
							Computed:    true,
						},
						"notify_consumer_ids_changed_enable": schema.BoolAttribute{
							Description: "是否开启消费者 ID 变更通知",
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
					},
				},
			},
		},
	}
}

func (c *ctyunRocketmqGroups) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRocketmqGroupsConfig
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
	params := &rocketmq.Mq2GroupListV3Request{
		RegionId:   regionId,
		ProdInstId: config.InstanceID.ValueString(),
	}
	//if config.Name.ValueString() != "" {
	//	params.GroupName = config.Name.ValueString()
	//}

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

	config.Groups = []CtyunRocketmqGroupModel{}
	for _, e := range resp.ReturnObj.Rows {
		if config.Name.ValueString() != "" && config.Name.ValueString() == e.GroupName {
			item := CtyunRocketmqGroupModel{}
			item.Name = types.StringValue(e.GroupName)
			item.ConsumeEnable = types.BoolValue(e.ConsumeEnable)
			item.FirstConsumeMechanism = types.Int32Value(e.FirstConsumeMechanism)
			item.PullMechanism = types.Int32Value(e.PullMechanism)
			item.RetryQueueNums = types.Int32Value(e.RetryQueueNums)
			item.RetryMaxTimes = types.Int32Value(e.RetryMaxTimes)
			item.BrokerId = types.Int32Value(e.BrokerId)
			item.WhichBrokerWhenConsumeSlowly = types.Int32Value(e.WhichBrokerWhenConsumeSlowly)
			item.NotifyConsumerIdsChangedEnable = types.BoolValue(e.NotifyConsumerIdsChangedEnable)
			item.Remark = types.StringValue(e.Remark)
			item.BrokerName = types.StringValue(e.BrokerName)
			config.Groups = append(config.Groups, item)
		}

	}

	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRocketmqGroups) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
