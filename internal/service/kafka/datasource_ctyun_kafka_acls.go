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
	_ datasource.DataSource              = &ctyunKafkaAcls{}
	_ datasource.DataSourceWithConfigure = &ctyunKafkaAcls{}
)

type ctyunKafkaAcls struct {
	meta *common.CtyunMetadata
}

func NewCtyunKafkaAcls() datasource.DataSource {
	return &ctyunKafkaAcls{}
}

func (c *ctyunKafkaAcls) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_kafka_acls"
}

type CtyunKafkaAclsModel struct {
	Id          types.Int32  `tfsdk:"id"`
	StrategyId  types.String `tfsdk:"strategy_id"`
	Name        types.String `tfsdk:"name"`
	UseNewTopic types.Int32  `tfsdk:"use_new_topic"`
	CreateTime  types.String `tfsdk:"create_time"`
}

type CtyunKafkaAclsConfig struct {
	Name       types.String          `tfsdk:"name"`
	InstanceId types.String          `tfsdk:"instance_id"`
	PageNum    types.String          `tfsdk:"page_num"`
	PageSize   types.String          `tfsdk:"page_size"`
	Total      types.Int32           `tfsdk:"total"`
	Acls       []CtyunKafkaAclsModel `tfsdk:"acls"`
	RegionId   types.String          `tfsdk:"region_id"`
}

func (c *ctyunKafkaAcls) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029624/10145597`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "策略名称，模糊查询",
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
			"page_num": schema.StringAttribute{
				Optional:    true,
				Description: "分页中的页数，默认1，范围1-40000",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},

			"page_size": schema.StringAttribute{
				Optional:    true,
				Description: "分页中的每页大小，默认10，范围1-40000",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"total": schema.Int32Attribute{
				Computed:    true,
				Description: "策略总数",
			},

			"acls": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int32Attribute{
							Computed:    true,
							Description: "策略ID",
						},
						"strategy_id": schema.StringAttribute{
							Computed:    true,
							Description: "策略唯一ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "策略名称",
						},
						"use_new_topic": schema.Int32Attribute{
							Computed:    true,
							Description: "是否应用到新增主题，1：是，2：否",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
					},
				},
				Description: "ACL策略列表",
			},
		},
	}
}

func (c *ctyunKafkaAcls) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunKafkaAclsConfig
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
	params := &ctgkafka.CtgkafkaAclStrategyListRequest{
		RegionId:   config.RegionId.ValueString(),
		ProdInstId: config.InstanceId.ValueString(),
		Name:       config.Name.ValueString(),
		PageNum:    config.PageNum.ValueString(),
		PageSize:   config.PageSize.ValueString(),
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaAclStrategyListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil || resp.ReturnObj.Data == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回值
	config.Acls = []CtyunKafkaAclsModel{}
	config.Total = types.Int32Value(resp.ReturnObj.Total)

	// 解析ACL策略列表
	for _, data := range resp.ReturnObj.Data {
		item := CtyunKafkaAclsModel{
			Id:          types.Int32Value(data.Id),
			StrategyId:  types.StringValue(data.StrategyId),
			Name:        types.StringValue(data.Name),
			UseNewTopic: types.Int32Value(data.UseNewTopic),
			CreateTime:  types.StringValue(utils.ConvertToUTCZ(utils.Layout2, data.CreateTime)),
		}

		config.Acls = append(config.Acls, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunKafkaAcls) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
