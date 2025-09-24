package rabbitmq

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/amqp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunRabbitmqInstances{}
	_ datasource.DataSourceWithConfigure = &ctyunRabbitmqInstances{}
)

type ctyunRabbitmqInstances struct {
	meta *common.CtyunMetadata
}

func NewCtyunRabbitmqInstances() datasource.DataSource {
	return &ctyunRabbitmqInstances{}
}

func (c *ctyunRabbitmqInstances) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rabbitmq_instances"
}

type CtyunRabbitmqInstancesModel struct {
	ID           types.String `tfsdk:"id"`
	Prod         types.String `tfsdk:"prod"`
	EngineType   types.String `tfsdk:"engine_type"`   // 引擎类型
	BillMode     types.String `tfsdk:"bill_mode"`     // 账单
	ExpireTime   types.String `tfsdk:"expire_time"`   // 过期时间
	CreateTime   types.String `tfsdk:"create_time"`   // 创建时间
	InstanceName types.String `tfsdk:"instance_name"` // 实例名称
	Status       types.Int32  `tfsdk:"status"`        // 状态
	StatusDesc   types.String `tfsdk:"status_desc"`   // 状态描述
}

type CtyunRabbitmqInstancesConfig struct {
	RegionID   types.String                  `tfsdk:"region_id"`
	PageNo     types.Int32                   `tfsdk:"page_no"`
	PageSize   types.Int32                   `tfsdk:"page_size"`
	InstanceID types.String                  `tfsdk:"instance_id"`
	Instances  []CtyunRabbitmqInstancesModel `tfsdk:"instances"`
}

func (c *ctyunRabbitmqInstances) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10000118/10001967**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页数据量大小",
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "实例ID",
			},
			"instances": schema.ListNestedAttribute{
				Description: "实例列表",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "id",
						},
						"prod": schema.StringAttribute{
							Computed:    true,
							Description: "规格",
						},
						"engine_type": schema.StringAttribute{
							Computed:    true,
							Description: "引擎类型",
						},
						"bill_mode": schema.StringAttribute{
							Computed:    true,
							Description: "账单",
						},
						"expire_time": schema.StringAttribute{
							Computed:    true,
							Description: "过期时间",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"instance_name": schema.StringAttribute{
							Computed:    true,
							Description: "实例名称",
						},
						"status": schema.Int32Attribute{
							Computed:    true,
							Description: "状态",
						},
						"status_desc": schema.StringAttribute{
							Computed:    true,
							Description: "状态描述",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRabbitmqInstances) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRabbitmqInstancesConfig
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
	if config.InstanceID.ValueString() != "" {
		err = c.getByID(ctx, &config)
	} else {
		err = c.getByPage(ctx, &config)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRabbitmqInstances) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// getByID 通过id查询
func (c *ctyunRabbitmqInstances) getByID(ctx context.Context, config *CtyunRabbitmqInstancesConfig) (err error) {

	// 组装请求体
	params := &amqp.AmqpInstancesQueryDetailRequest{
		RegionId:   config.RegionID.ValueString(),
		ProdInstId: config.InstanceID.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpInstancesQueryDetailApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	r := resp.ReturnObj.Data
	config.Instances = []CtyunRabbitmqInstancesModel{}
	// 解析返回值
	item := CtyunRabbitmqInstancesModel{
		ID:           types.StringValue(r.ProdInstId),
		BillMode:     types.StringValue(map[string]string{"1": "包年包月", "2": "按需计费"}[r.BillMode]),
		Prod:         types.StringValue(r.Prod),
		EngineType:   types.StringValue(r.EngineType),
		ExpireTime:   types.StringValue(r.ExpireTime),
		CreateTime:   types.StringValue(r.CreateTime),
		InstanceName: types.StringValue(r.ClusterName),
		Status:       types.Int32Value(r.Status),
		StatusDesc: types.StringValue(map[int32]string{
			1: "运行中", 3: "已注销", 4: "已退订", 5: "变更中", 6: "创建中",
		}[r.Status]),
	}
	config.Instances = append(config.Instances, item)
	return

}

// getByPage 无id时查询
func (c *ctyunRabbitmqInstances) getByPage(ctx context.Context, config *CtyunRabbitmqInstancesConfig) (err error) {
	// 组装请求体
	params := &amqp.AmqpInstancesQueryRequest{
		RegionId: config.RegionID.ValueString(),
		PageNum:  1,
		PageSize: 10,
	}
	if config.PageNo.ValueInt32() > 0 {
		params.PageNum = config.PageNo.ValueInt32()
	}
	if config.PageSize.ValueInt32() > 0 {
		params.PageSize = config.PageSize.ValueInt32()
	}
	// 调用API
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpInstancesQueryApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	config.Instances = []CtyunRabbitmqInstancesModel{}
	// 解析返回值
	for _, r := range resp.ReturnObj.Data {
		item := CtyunRabbitmqInstancesModel{
			ID:           types.StringValue(r.ProdInstId),
			BillMode:     types.StringValue(map[string]string{"1": "包年包月", "2": "按需计费"}[r.BillMode]),
			Prod:         types.StringValue(r.Prod),
			EngineType:   types.StringValue(r.EngineType),
			ExpireTime:   types.StringValue(r.ExpireTime),
			CreateTime:   types.StringValue(r.CreateTime),
			InstanceName: types.StringValue(r.ClusterName),
			Status:       types.Int32Value(r.Status),
			StatusDesc: types.StringValue(map[int32]string{
				1: "运行中", 3: "已注销", 4: "已退订", 5: "变更中", 6: "创建中",
			}[r.Status]),
		}
		config.Instances = append(config.Instances, item)
	}
	return
}
