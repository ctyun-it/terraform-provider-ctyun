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
	_ datasource.DataSource              = &ctyunKafkaUsers{}
	_ datasource.DataSourceWithConfigure = &ctyunKafkaUsers{}
)

type ctyunKafkaUsers struct {
	meta *common.CtyunMetadata
}

func NewCtyunKafkaUsers() datasource.DataSource {
	return &ctyunKafkaUsers{}
}

func (c *ctyunKafkaUsers) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_kafka_users"
}

type CtyunKafkaUsersModel struct {
	Id          types.Int32  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreateTime  types.String `tfsdk:"create_time"`
}

type CtyunKafkaUsersConfig struct {
	Name       types.String           `tfsdk:"name"`
	InstanceId types.String           `tfsdk:"instance_id"`
	PageNum    types.String           `tfsdk:"page_num"`
	PageSize   types.String           `tfsdk:"page_size"`
	Total      types.Int32            `tfsdk:"total"`
	Users      []CtyunKafkaUsersModel `tfsdk:"users"`
	RegionId   types.String           `tfsdk:"region_id"`
}

func (c *ctyunKafkaUsers) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029624/10145597`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "用户名称，模糊查询",
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
				Description: "用户总数",
			},

			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int32Attribute{
							Computed:    true,
							Description: "用户ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "用户名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "用户描述",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
					},
				},
				Description: "用户列表",
			},
		},
	}
}

func (c *ctyunKafkaUsers) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunKafkaUsersConfig
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
	params := &ctgkafka.CtgkafkaSaslUserQueryV3Request{
		RegionId:   config.RegionId.ValueString(),
		ProdInstId: config.InstanceId.ValueString(),
		Username:   config.Name.ValueString(),
		PageNum:    config.PageNum.ValueString(),
		PageSize:   config.PageSize.ValueString(),
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaSaslUserQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
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
	config.Users = []CtyunKafkaUsersModel{}
	config.Total = types.Int32Value(resp.ReturnObj.Total)

	// 解析用户列表
	for _, data := range resp.ReturnObj.Data {
		item := CtyunKafkaUsersModel{
			Id:          types.Int32Value(data.Id),
			Name:        types.StringValue(data.Username),
			Description: types.StringValue(data.Description),
			CreateTime:  types.StringValue(utils.ConvertToUTCZ(utils.Layout2, data.Ctime)),
		}

		config.Users = append(config.Users, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunKafkaUsers) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
