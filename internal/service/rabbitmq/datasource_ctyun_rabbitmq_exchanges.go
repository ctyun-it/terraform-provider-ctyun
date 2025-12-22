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
	_ datasource.DataSource              = &ctyunRabbitmqExchanges{}
	_ datasource.DataSourceWithConfigure = &ctyunRabbitmqExchanges{}
)

type ctyunRabbitmqExchanges struct {
	meta *common.CtyunMetadata
}

func NewCtyunRabbitmqExchanges() datasource.DataSource {
	return &ctyunRabbitmqExchanges{}
}

func (c *ctyunRabbitmqExchanges) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rabbitmq_exchanges"
}

type CtyunRabbitmqExchangesModel struct {
	Name         types.String `tfsdk:"name"`
	Vhost        types.String `tfsdk:"vhost"`
	Type         types.String `tfsdk:"type"`
	Durable      types.Bool   `tfsdk:"durable"`
	Internal     types.Bool   `tfsdk:"internal"`
	AutoDelete   types.Bool   `tfsdk:"auto_delete"`
	XDelayedType types.String `tfsdk:"x_delayed_type"`
}

type CtyunRabbitmqExchangesConfig struct {
	RegionID   types.String                  `tfsdk:"region_id"`
	InstanceID types.String                  `tfsdk:"instance_id"`
	Vhost      types.String                  `tfsdk:"vhost"`
	Name       types.String                  `tfsdk:"name"`
	Exchanges  []CtyunRabbitmqExchangesModel `tfsdk:"exchanges"`
}

func (c *ctyunRabbitmqExchanges) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10000118/10001967`,
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
			"vhost": schema.StringAttribute{
				Optional:    true,
				Description: "vhost名称",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 128),
					stringvalidator.Any(
						stringvalidator.RegexMatches(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "vhost名称不符合规则"),
						stringvalidator.OneOf("/"),
					),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "交换器名称（模糊匹配）",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 128),
					//stringvalidator.RegexMatches(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "交换器名称不符合规则"),
				},
			},
			"exchanges": schema.ListNestedAttribute{
				Description: "交换器列表",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "交换器名称",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "交换器类型",
							Computed:    true,
						},
						"vhost": schema.StringAttribute{
							Description: "vhost",
							Computed:    true,
						},
						"durable": schema.BoolAttribute{
							Description: "是否持久化",
							Computed:    true,
						},
						"internal": schema.BoolAttribute{
							Description: "是否内置",
							Computed:    true,
						},
						"auto_delete": schema.BoolAttribute{
							Description: "是否自动删除",
							Computed:    true,
						},
						"x_delayed_type": schema.StringAttribute{
							Description: "当type为x-delayed-message时有值",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRabbitmqExchanges) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRabbitmqExchangesConfig
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
	params := &amqp.AmqpExchangeQueryV3Request{
		RegionId:   regionId,
		ProdInstId: config.InstanceID.ValueString(),
	}
	if config.Name.ValueString() != "" {
		params.Name = config.Name.ValueString()
	}
	if config.Vhost.ValueString() != "" {
		params.Vhost = config.Vhost.ValueString()
	}
	// 调用API
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpExchangeQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	config.Exchanges = []CtyunRabbitmqExchangesModel{}
	// 解析返回值
	for _, e := range resp.ReturnObj.Data.Items {
		item := CtyunRabbitmqExchangesModel{
			Name:         types.StringValue(e.Name),
			Vhost:        types.StringValue(e.Vhost),
			Type:         types.StringValue(e.RawType),
			Durable:      types.BoolValue(e.Durable),
			Internal:     types.BoolValue(e.Internal),
			AutoDelete:   types.BoolValue(e.Auto_delete),
			XDelayedType: types.StringValue(e.Argument.XDelayedType),
		}
		config.Exchanges = append(config.Exchanges, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRabbitmqExchanges) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
