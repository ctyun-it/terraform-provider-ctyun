package kafka

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgkafka "github.com/ctyun-it/terraform-provider-ctyun/internal/core/kafka"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunKafkaConsumerGroups{}
	_ datasource.DataSourceWithConfigure = &ctyunKafkaConsumerGroups{}
)

type ctyunKafkaConsumerGroups struct {
	meta *common.CtyunMetadata
}

func NewCtyunKafkaConsumerGroups() datasource.DataSource {
	return &ctyunKafkaConsumerGroups{}
}

func (c *ctyunKafkaConsumerGroups) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_kafka_consumer_groups"
}

type CtyunKafkaConsumerGroupsModel struct {
	ID            types.Int32  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Ctime         types.String `tfsdk:"create_time"`
	State         types.String `tfsdk:"state"`
	CoordinatorId types.Int32  `tfsdk:"coordinator_id"`
}

type CtyunKafkaConsumerGroupsConfig struct {
	Name           types.String                    `tfsdk:"name"`
	RegionID       types.String                    `tfsdk:"region_id"`
	InstanceId     types.String                    `tfsdk:"instance_id"`
	PageNum        types.Int32                     `tfsdk:"page_no"`
	PageSize       types.Int32                     `tfsdk:"page_size"`
	Total          types.Int32                     `tfsdk:"total"`
	ConsumerGroups []CtyunKafkaConsumerGroupsModel `tfsdk:"consumer_groups"`
}

func (c *ctyunKafkaConsumerGroups) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029624/10145103`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "消费组名称，模糊查询",
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
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "分页中的页数，默认1，范围1-40000",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "分页中的每页大小，默认10，范围1-40000",
			},
			"total": schema.Int32Attribute{
				Computed:    true,
				Description: "消费组总数",
			},
			"consumer_groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int32Attribute{
							Computed:    true,
							Description: "消费组ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "消费组名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "消费组描述",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "消费组状态",
						},
						"coordinator_id": schema.Int32Attribute{
							Computed:    true,
							Description: "协调器编号",
						},
					},
				},
				Description: "消费组列表",
			},
		},
	}
}

func (c *ctyunKafkaConsumerGroups) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunKafkaConsumerGroupsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	config.RegionID = types.StringValue(regionId)
	// 组装请求体
	params := &ctgkafka.CtgkafkaConsumerGroupQueryV3Request{
		RegionId:   config.RegionID.ValueString(),
		ProdInstId: config.InstanceId.ValueString(),
		GroupName:  config.Name.ValueString(),
	}
	if config.PageNum.ValueInt32() > 0 {
		params.PageNum = fmt.Sprintf("%d", config.PageNum.ValueInt32())
	}
	if config.PageSize.ValueInt32() > 0 {
		params.PageSize = fmt.Sprintf("%d", config.PageSize.ValueInt32())
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaConsumerGroupQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
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
	config.ConsumerGroups = []CtyunKafkaConsumerGroupsModel{}
	config.Total = types.Int32Value(resp.ReturnObj.Total)
	for _, data := range resp.ReturnObj.Data {
		item := CtyunKafkaConsumerGroupsModel{
			ID:            types.Int32Value(data.Id),
			Name:          types.StringValue(data.Name),
			Description:   types.StringValue(data.Description),
			Ctime:         types.StringValue(utils.ConvertToUTCZ(utils.Layout2, data.Ctime)),
			State:         types.StringValue(data.State),
			CoordinatorId: types.Int32Value(data.CoordinatorId),
		}

		config.ConsumerGroups = append(config.ConsumerGroups, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunKafkaConsumerGroups) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
