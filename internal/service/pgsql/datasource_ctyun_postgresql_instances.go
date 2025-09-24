package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ datasource.DataSource              = &ctyunPgsqlInstances{}
	_ datasource.DataSourceWithConfigure = &ctyunPgsqlInstances{}
)

type ctyunPgsqlInstances struct {
	meta *common.CtyunMetadata
}

func NewCtyunPgsqlInstances() datasource.DataSource {
	return &ctyunPgsqlInstances{}
}
func (c *ctyunPgsqlInstances) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunPgsqlInstances) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_instances"
}

func (c *ctyunPgsqlInstances) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10034019/10153165**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "项目ID",
			},
			"page_num": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "当前页码。默认:1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "页大小，范围1-500。默认:20",
				Validators: []validator.Int32{
					int32validator.Between(1, 500),
				},
			},
			"prod_inst_name": schema.StringAttribute{
				Optional:    true,
				Description: "实例名称，支持模糊匹配",
			},
			"label_name": schema.StringAttribute{
				Optional:    true,
				Description: "标签名称（一级标签）",
			},
			"label_value": schema.StringAttribute{
				Optional:    true,
				Description: "标签值（二级标签）",
			},
			"prod_inst_id": schema.StringAttribute{
				Optional:    true,
				Description: "实例ID",
			},
			"instance_type": schema.StringAttribute{
				Optional:    true,
				Description: "实例类型（primary/readonly）",
				Validators: []validator.String{
					stringvalidator.OneOf("primary", "readonly"),
				},
			},
			"pgsql_instances": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "实例创建时间",
						},
						"prod_db_engine": schema.StringAttribute{
							Computed:    true,
							Description: "数据库引擎类型",
						},
						"prod_inst_id": schema.StringAttribute{
							Computed:    true,
							Description: "实例唯一ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "实例名称",
						},
						"prod_running_status": schema.Int32Attribute{
							Computed:    true,
							Description: "运行状态: 0(运行中),1(重启中),2(备份中),3(恢复中),1001(已停止),1006(复失败),1007(VIP不可用),1008(GATEWAY不可用),1009(主库不可用),1010(备库不可用),1021(实例维护中),2000(开通中),2002(已退订),2005(扩容中),2011(冻结)",
							Validators: []validator.Int32{
								int32validator.OneOf(business.PgsqlProdRunningStatus...),
							},
						},
						"alive": schema.Int32Attribute{
							Computed:    true,
							Description: "实例存活状态: 0(存活), -1(异常)",
							Validators: []validator.Int32{
								int32validator.Between(0, 1),
							},
						},
						"prod_order_status": schema.Int32Attribute{
							Computed:    true,
							Description: "订单状态: 0(正常),1(冻结),2(删除),3(操作中),4(失败),2005(扩容中)",
							Validators: []validator.Int32{
								int32validator.OneOf(business.PgsqlProdOrderStatus...),
							},
						},
						"prod_type": schema.Int32Attribute{
							Computed:    true,
							Description: "部署方式: 0(单机部署),1(主备部署)",
							Validators: []validator.Int32{
								int32validator.Between(0, 1),
							},
						},
						"read_port": schema.Int32Attribute{
							Computed:    true,
							Description: "读连接端口号",
							Validators: []validator.Int32{
								int32validator.Between(1, 65535),
							},
						},
						"vip": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟IP地址",
						},
						"write_port": schema.Int32Attribute{
							Computed:    true,
							Description: "写连接端口号",
						},
						"readonly_instance_ids": schema.StringAttribute{
							Computed:    true,
							Description: "关联的只读实例ID列表",
						},
						"instance_type": schema.StringAttribute{
							Computed:    true,
							Description: "实例类型: primary(主实例), readonly(只读实例)",
						},
						"tool_type": schema.Int32Attribute{
							Computed:    true,
							Description: "备份工具类型: 1(pg_baseback), 2(pgbackrest), 3(s3)",
							Validators: []validator.Int32{
								int32validator.Between(1, 3),
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunPgsqlInstances) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunPgsqlInstancesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	params := &pgsql.PgsqlListRequest{}

	if config.PageNum.ValueInt32() == 0 {
		config.PageNum = types.Int32Value(1)
		params.PageNum = 1
	}
	if config.PageSize.ValueInt32() == 0 {
		config.PageSize = types.Int32Value(20)
		params.PageSize = 20
	}
	if !config.ProdInstID.IsNull() && !config.ProdInstID.IsUnknown() {
		params.ProdInstId = config.ProdInstID.ValueStringPointer()
	}
	if !config.LabelName.IsNull() && !config.LabelName.IsUnknown() {
		params.LabelName = config.LabelName.ValueStringPointer()
	}
	if !config.LabelValue.IsNull() && !config.LabelValue.IsUnknown() {
		params.LabelValue = config.LabelValue.ValueStringPointer()
	}
	if !config.ProdInstName.IsNull() && !config.ProdInstName.IsUnknown() {
		params.ProdInstName = config.ProdInstName.ValueStringPointer()
	}
	if config.InstanceType.IsNull() && !config.InstanceType.IsUnknown() {
		params.InstanceType = config.InstanceType.ValueStringPointer()
	}
	headers := &pgsql.PgsqlListRequestHeader{
		RegionID: regionId,
	}
	if config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		headers.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlListApi.Do(ctx, c.meta.Credential, params, headers)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询反馈空指针，请稍后尝试！")
		return
	} else if resp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s ", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回值
	var ctyunPgsqlInstanceInfoModel []CtyunPgsqlInstanceInfoModel
	instances := resp.ReturnObj.List
	for _, instance := range instances {
		pgsqlInstance := CtyunPgsqlInstanceInfoModel{}
		pgsqlInstance.CreateTime = types.StringValue(instance.CreateTime)
		pgsqlInstance.ProdDbEngine = types.StringValue(instance.ProdDbEngine)
		pgsqlInstance.ProdInstId = types.StringValue(instance.ProdInstId)
		pgsqlInstance.Name = types.StringValue(instance.ProdInstName)
		pgsqlInstance.ProdRunningStatus = types.Int32Value(instance.ProdRunningStatus)
		pgsqlInstance.Alive = types.Int32Value(instance.Alive)
		pgsqlInstance.ProdOrderStatus = types.Int32Value(instance.ProdOrderStatus)
		pgsqlInstance.ProdType = types.Int32Value(instance.ProdType)
		pgsqlInstance.ReadPort = types.Int32Value(instance.ReadPort)
		pgsqlInstance.Vip = types.StringValue(instance.Vip)
		pgsqlInstance.WritePort = types.Int32Value(instance.WritePort)
		pgsqlInstance.InstanceType = types.StringValue(instance.InstanceType)
		pgsqlInstance.ToolType = types.Int32Value(instance.ToolType)
		if len(instance.ReadonlyInstnaceIds) > 0 {
			pgsqlInstance.ReadonlyInstanceIds = types.StringValue(strings.Join(instance.ReadonlyInstnaceIds, ","))
		} else {
			pgsqlInstance.ReadonlyInstanceIds = types.StringValue("")
		}
		ctyunPgsqlInstanceInfoModel = append(ctyunPgsqlInstanceInfoModel, pgsqlInstance)
	}
	config.PgsqlInstances = ctyunPgsqlInstanceInfoModel
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunPgsqlInstancesConfig struct {
	RegionID       types.String                  `tfsdk:"region_id"`      // 区域id
	ProjectID      types.String                  `tfsdk:"project_id"`     // 项目id
	PageNum        types.Int32                   `tfsdk:"page_num"`       // 当前页（必填）
	PageSize       types.Int32                   `tfsdk:"page_size"`      // 页大小，范围1-500（必填）
	ProdInstName   types.String                  `tfsdk:"prod_inst_name"` // 实例名称，支持模糊匹配（可选）
	LabelName      types.String                  `tfsdk:"label_name"`     // 标签名称（一级标签）（可选）
	LabelValue     types.String                  `tfsdk:"label_value"`    // 标签值（二级标签）（可选）
	ProdInstID     types.String                  `tfsdk:"prod_inst_id"`   // 实例id（可选）
	InstanceType   types.String                  `tfsdk:"instance_type"`  // 实例类型（primary/readonly）（可选）
	PgsqlInstances []CtyunPgsqlInstanceInfoModel `tfsdk:"pgsql_instances"`
}

type CtyunPgsqlInstanceInfoModel struct {
	CreateTime          types.String `tfsdk:"create_time"`           // 创建时间
	ProdDbEngine        types.String `tfsdk:"prod_db_engine"`        // 数据库实例引擎
	ProdInstId          types.String `tfsdk:"prod_inst_id"`          // 实例ID
	Name                types.String `tfsdk:"name"`                  // 实例名称
	ProdRunningStatus   types.Int32  `tfsdk:"prod_running_status"`   // 运行状态代码
	Alive               types.Int32  `tfsdk:"alive"`                 // 实例存活状态
	ProdOrderStatus     types.Int32  `tfsdk:"prod_order_status"`     // 订单状态代码
	ProdType            types.Int32  `tfsdk:"prod_type"`             // 实例部署方式
	ReadPort            types.Int32  `tfsdk:"read_port"`             // 读端口
	Vip                 types.String `tfsdk:"vip"`                   // 虚拟IP地址
	WritePort           types.Int32  `tfsdk:"write_port"`            // 写端口
	ReadonlyInstanceIds types.String `tfsdk:"readonly_instance_ids"` // 只读实例ID列表,用逗号分割
	InstanceType        types.String `tfsdk:"instance_type"`         // 实例类型
	ToolType            types.Int32  `tfsdk:"tool_type"`             // 备份工具类型
}
