package rabbitmq

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/amqp"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
)

var (
	_ datasource.DataSource              = &ctyunRabbitmqVhosts{}
	_ datasource.DataSourceWithConfigure = &ctyunRabbitmqVhosts{}
)

type ctyunRabbitmqVhosts struct {
	meta *common.CtyunMetadata
}

func NewCtyunRabbitmqVhosts() datasource.DataSource {
	return &ctyunRabbitmqVhosts{}
}

func (c *ctyunRabbitmqVhosts) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rabbitmq_vhosts"
}

type CtyunRabbitmqVhostsModel struct {
	Name                types.String  `tfsdk:"name"`
	PublishBytesRate    types.Float64 `tfsdk:"publish_bytes_rate"`
	DeliverBytesRate    types.Float64 `tfsdk:"deliver_bytes_rate"`
	PublishMessagesRate types.Float64 `tfsdk:"publish_messages_rate"`
	DeliverMessagesRate types.Float64 `tfsdk:"deliver_messages_rate"`
}

type CtyunRabbitmqVhostsConfig struct {
	RegionID   types.String               `tfsdk:"region_id"`
	InstanceID types.String               `tfsdk:"instance_id"`
	Vhosts     []CtyunRabbitmqVhostsModel `tfsdk:"vhosts"`
	Name       types.String               `tfsdk:"name"`
}

func (c *ctyunRabbitmqVhosts) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10000118/10220893`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "rabbitMq实例ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "vhost名称，只能包含字母，数字，短横线-和下划线_",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 128),
					stringvalidator.Any(
						stringvalidator.RegexMatches(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "vhost名称不符合规则"),
						stringvalidator.OneOf("/"),
					),
				},
			},
			"vhosts": schema.ListNestedAttribute{
				Description: "vhost列表",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "vhost名称",
							Computed:    true,
						},
						"publish_bytes_rate": schema.Float64Attribute{
							Description: "消息生产流量（字节/秒）",
							Computed:    true,
						},
						"deliver_bytes_rate": schema.Float64Attribute{
							Description: "消息消费流量（字节/秒）",
							Computed:    true,
						},
						"publish_messages_rate": schema.Float64Attribute{
							Description: "消息生产速率（条/秒）",
							Computed:    true,
						},
						"deliver_messages_rate": schema.Float64Attribute{
							Description: "消息消费速率（条/秒）",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRabbitmqVhosts) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRabbitmqVhostsConfig
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
	params := &amqp.AmqpVhostQueryV3Request{
		RegionId:   regionId,
		ProdInstId: config.InstanceID.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpVhostQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	config.Vhosts = []CtyunRabbitmqVhostsModel{}
	// 解析返回值
	for _, vhost := range resp.ReturnObj.Data.VhostsDetail {
		item := CtyunRabbitmqVhostsModel{
			Name:                types.StringValue(vhost.Name),
			PublishBytesRate:    types.Float64Value(vhost.PublishBytesRate),
			DeliverBytesRate:    types.Float64Value(vhost.DeliverBytesRate),
			PublishMessagesRate: types.Float64Value(vhost.PublishMessagesRate),
			DeliverMessagesRate: types.Float64Value(vhost.DeliverMessagesRate),
		}
		if config.Name.ValueString() == "" {
			config.Vhosts = append(config.Vhosts, item)
		} else if vhost.Name == config.Name.ValueString() {
			config.Vhosts = append(config.Vhosts, item)
			break
		}
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRabbitmqVhosts) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
