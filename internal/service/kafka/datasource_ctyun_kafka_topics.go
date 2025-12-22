package kafka

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgkafka "github.com/ctyun-it/terraform-provider-ctyun/internal/core/kafka"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunKafkaTopics{}
	_ datasource.DataSourceWithConfigure = &ctyunKafkaTopics{}
)

type ctyunKafkaTopics struct {
	meta *common.CtyunMetadata
}

func NewCtyunKafkaTopics() datasource.DataSource {
	return &ctyunKafkaTopics{}
}

func (c *ctyunKafkaTopics) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_kafka_topics"
}

type CtyunKafkaTopicsModel struct {
	InstanceId   types.String `tfsdk:"instance_id"`
	Name         types.String `tfsdk:"name"`
	PartitionNum types.Int32  `tfsdk:"partition_num"`
	Factor       types.Int32  `tfsdk:"factor"`
}

type CtyunKafkaTopicsConfig struct {
	Name       types.String            `tfsdk:"name"`
	Labels     types.String            `tfsdk:"labels"`
	InstanceId types.String            `tfsdk:"instance_id"`
	PageNum    types.Int32             `tfsdk:"page_num"`
	PageSize   types.Int32             `tfsdk:"page_size"`
	Total      types.Int32             `tfsdk:"total"`
	Topics     []CtyunKafkaTopicsModel `tfsdk:"topics"`
	RegionId   types.String            `tfsdk:"region_id"`

	MaxPartitions      types.Int32 `tfsdk:"max_partitions"`       // 分区数量上限
	RemainPartitions   types.Int32 `tfsdk:"remain_partitions"`    // 剩余分区数量
	TopicMaxPartitions types.Int32 `tfsdk:"topic_max_partitions"` // 主题最大分区数量

}

func (c *ctyunKafkaTopics) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029624/10144604`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "主题名称，模糊查询",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "实例ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"labels": schema.StringAttribute{
				Optional:    true,
				Description: "标签",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"page_num": schema.Int32Attribute{
				Optional:    true,
				Description: "分页中的页数，默认1，范围1-40000",
				Validators: []validator.Int32{
					int32validator.Between(1, 40000),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "分页中的每页大小，默认10，范围1-40000",
				Validators: []validator.Int32{
					int32validator.Between(1, 40000),
				},
			},
			"total": schema.Int32Attribute{
				Computed:    true,
				Description: "主题总数",
			},
			"max_partitions": schema.Int32Attribute{
				Computed:    true,
				Description: "分区数量上限",
			},
			"remain_partitions": schema.Int32Attribute{
				Computed:    true,
				Description: "剩余分区数量",
			},
			"topic_max_partitions": schema.Int32Attribute{
				Computed:    true,
				Description: "主题最大分区数量",
			},

			"topics": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "实例ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "主题名称",
						},
						"partition_num": schema.Int32Attribute{
							Computed:    true,
							Description: "主题分区数量",
						},
						"factor": schema.Int32Attribute{
							Computed:    true,
							Description: "主题副本数",
						},
					},
				},
				Description: "主题列表",
			},
		},
	}
}

func (c *ctyunKafkaTopics) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunKafkaTopicsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := c.meta.GetExtraIfEmpty(config.RegionId.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	config.RegionId = types.StringValue(regionId)
	// 组装请求体
	params := &ctgkafka.CtgkafkaTopicQueryV3Request{
		RegionId:   config.RegionId.ValueString(),
		ProdInstId: config.InstanceId.ValueString(),
		TopicName:  config.Name.ValueString(),
		Labels:     config.Labels.ValueString(),
		PageNum:    config.PageNum.ValueInt32(),
		PageSize:   config.PageSize.ValueInt32(),
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaTopicQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s ", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回值
	config.Topics = []CtyunKafkaTopicsModel{}
	config.Total = types.Int32Value(resp.ReturnObj.Total)
	config.RegionId = types.StringValue(params.RegionId)

	config.MaxPartitions = types.Int32Value(resp.ReturnObj.MaxPartitions)
	config.RemainPartitions = types.Int32Value(resp.ReturnObj.RemainPartitions)
	config.TopicMaxPartitions = types.Int32Value(resp.ReturnObj.TopicMaxPartitions)

	// 解析主题列表
	for _, data := range resp.ReturnObj.Data {
		item := CtyunKafkaTopicsModel{
			InstanceId:   types.StringValue(data.ProdInstId),
			Name:         types.StringValue(data.Name),
			PartitionNum: types.Int32Value(data.PartitionNum),
			Factor:       types.Int32Value(data.Factor),
		}

		config.Topics = append(config.Topics, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunKafkaTopics) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
