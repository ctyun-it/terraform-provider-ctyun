package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunMongodbInstances{}
	_ datasource.DataSourceWithConfigure = &ctyunMongodbInstances{}
)

type ctyunMongodbInstances struct {
	meta *common.CtyunMetadata
}

func (c *ctyunMongodbInstances) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mongodb_instances"
}
func NewCtyunMongodbInstances() datasource.DataSource {
	return &ctyunMongodbInstances{}
}
func (c *ctyunMongodbInstances) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunMongodbInstances) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10034467/10089535**`,
		Attributes: map[string]schema.Attribute{
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "当前页，不传默认为1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "页大小，不传默认为10",
			},
			"res_db_engine": schema.StringAttribute{
				Optional:    true,
				Description: "版本号",
			},
			"prod_inst_name": schema.StringAttribute{
				Optional:    true,
				Description: "实例名称",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "项目id",
			},
			"label_ids": schema.StringAttribute{
				Optional:    true,
				Description: "标签id",
			},
			"mongodb_instances": schema.ListNestedAttribute{
				Computed:    true,
				Description: "mongodb实例列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"prod_order_status": schema.Int32Attribute{
							Description: "订单状态: 0->订单正常, 1->订单冻结, 2->订单注销, 3->施工中, 4->施工失败",
							Computed:    true,
							Validators: []validator.Int32{
								int32validator.Between(0, 4),
							},
						},
						"subnet_id": schema.StringAttribute{
							Description: "子网ID",
							Computed:    true,
						},
						"maintain_time": schema.StringAttribute{
							Description: "可维护时间",
							Computed:    true,
						},
						"subnet": schema.StringAttribute{
							Description: "子网名称",
							Computed:    true,
						},
						"log_status": schema.BoolAttribute{
							Description: "实例日志审计状态",
							Computed:    true,
						},
						"order_id": schema.Int64Attribute{
							Description: "订单ID",
							Computed:    true,
						},
						"net_name": schema.StringAttribute{
							Description: "专有网络",
							Computed:    true,
						},
						"version_num": schema.StringAttribute{
							Description: "版本号",
							Computed:    true,
						},
						"security_group_id": schema.StringAttribute{
							Description: "安全组ID",
							Computed:    true,
						},
						"parameter_configsvr_group_used": schema.StringAttribute{
							Description: "参数配置",
							Computed:    true,
						},
						"disk_size": schema.Int32Attribute{
							Description: "存储空间大小",
							Computed:    true,
						},
						"tpl_name": schema.StringAttribute{
							Description: "模板名称",
							Computed:    true,
						},
						"prod_inst_set_name": schema.StringAttribute{
							Description: "实例对应的SET名",
							Computed:    true,
						},
						"released": schema.BoolAttribute{
							Description: "实例是否已被释放",
							Computed:    true,
						},
						"security_group": schema.StringAttribute{
							Description: "安全组",
							Computed:    true,
						},
						"prod_type": schema.Int32Attribute{
							Description: "实例类型: 0:单机, 2:副本集(三节点), 4:副本集(五节点), 6:副本集(七节点), 10:分片集群",
							Computed:    true,
							Validators: []validator.Int32{
								int32validator.OneOf(0, 2, 4, 6, 10),
							},
						},
						"expire_time": schema.Int64Attribute{
							Description: "到期时间",
							Computed:    true,
						},
						"prod_inst_id": schema.StringAttribute{
							Description: "实例ID",
							Required:    true, // 通常作为数据源的唯一标识
						},
						"project_name": schema.StringAttribute{
							Description: "企业项目名称",
							Computed:    true,
						},
						"project_id": schema.StringAttribute{
							Description: "企业项目ID",
							Computed:    true,
						},
						"destroyed_time": schema.StringAttribute{
							Description: "实例销毁时间",
							Computed:    true,
						},
						"prod_inst_flag": schema.StringAttribute{
							Description: "实例标识（格式：实例ID 实例名称）",
							Computed:    true,
						},
						"prod_db_engine": schema.StringAttribute{
							Description: "数据库引擎版本",
							Computed:    true,
						},
						"bill_mode": schema.Int32Attribute{
							Description: "计费模式: 1:包周期计费, 2:按需计费",
							Computed:    true,
							Validators: []validator.Int32{
								int32validator.OneOf(1, 2),
							},
						},
						"prod_id": schema.StringAttribute{
							Description: "产品ID",
							Computed:    true,
						},
						"restore_time": schema.StringAttribute{
							Description: "实例恢复时间",
							Computed:    true,
						},
						"prod_running_status": schema.Int32Attribute{
							Description: "实例运行状态: 0->运行正常, 1->重启中, 2->备份操作中, 3->恢复操作中, 4->转换ssl, 5->异常, 6->修改参数组中, 7->已冻结, 8->已注销, 9->施工中, 10->施工失败, 11->扩容中, 12->主备切换中",
							Computed:    true,
							Validators: []validator.Int32{
								int32validator.Between(0, 12),
							},
						},
						"disk_used": schema.StringAttribute{
							Description: "磁盘使用情况",
							Computed:    true,
						},
						"parameter_group_used": schema.StringAttribute{
							Description: "参数组名称（参数组版本）",
							Computed:    true,
						},
						"vpc_id": schema.StringAttribute{
							Description: "VPC网络ID",
							Computed:    true,
						},
						"innodb_thread_concurrency": schema.Int64Attribute{
							Description: "线程数",
							Computed:    true,
						},
						"disk_type": schema.StringAttribute{
							Description: "存储类型",
							Computed:    true,
						},
						"prod_bill_type": schema.Int32Attribute{
							Description: "计费类型: 0:按月计费, 1:按天计费, 2:按年计费, 3:按流量计费",
							Computed:    true,
							Validators: []validator.Int32{
								int32validator.Between(0, 3),
							},
						},
						"machine_spec": schema.StringAttribute{
							Description: "CPU内存规格",
							Computed:    true,
						},
						"prod_inst_name": schema.StringAttribute{
							Description: "实例名称",
							Computed:    true,
						},
						"innodb_buffer_pool_size": schema.StringAttribute{
							Description: "缓存池大小",
							Computed:    true,
						},
						"used_space": schema.StringAttribute{
							Description: "已使用空间",
							Computed:    true,
							Optional:    true, // 指针类型字段设为 Optional
						},
						"user_id": schema.Int64Attribute{
							Description: "用户ID",
							Computed:    true,
						},
						"prod_bill_time": schema.Int64Attribute{
							Description: "购买时长",
							Computed:    true,
						},
						"destroyed": schema.BoolAttribute{
							Description: "实例是否已经被销毁",
							Computed:    true,
						},
						"create_time": schema.Int64Attribute{
							Description: "创建时间",
							Computed:    true,
						},
						"tenant_id": schema.Int64Attribute{
							Description: "租户ID",
							Computed:    true,
						},
						"outer_id": schema.StringAttribute{
							Description: "外部产品ID",
							Computed:    true,
						},
						"tpl_code": schema.StringAttribute{
							Description: "模板编码",
							Computed:    true,
						},
						"read_port": schema.Int32Attribute{
							Description: "读端口",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (c *ctyunMongodbInstances) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunMongodbInstancesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	if config.PageNo.ValueInt32() == 0 {
		config.PageNo = types.Int32Value(1)
	}
	if config.PageSize.ValueInt32() == 0 {
		config.PageSize = types.Int32Value(10)
	}
	params := &mongodb.MongodbGetListRequest{
		PageNow:      config.PageNo.ValueInt32(),
		PageSize:     config.PageSize.ValueInt32(),
		ResDbEngine:  nil,
		ProdInstName: nil,
		LabelIds:     nil,
	}
	if !config.ResDbEngine.IsNull() {
		params.ResDbEngine = config.ResDbEngine.ValueStringPointer()
	}
	if !config.ProdInstName.IsNull() {
		params.ProdInstName = config.ProdInstName.ValueStringPointer()
	}
	if !config.LabelIds.IsNull() {
		params.LabelIds = config.LabelIds.ValueStringPointer()
	}

	header := &mongodb.MongodbGetListHeaders{
		ProjectID: nil,
		RegionID:  regionId,
	}
	if config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s ", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var mongodbInstances []CtyunMongodbInstanceModel
	for _, mongodbItem := range resp.ReturnObj.List {
		var mongodbInst CtyunMongodbInstanceModel
		mongodbInst.ProdOrderStatus = types.Int32Value(mongodbItem.ProdOrderStatus)
		mongodbInst.SubNetID = types.StringValue(mongodbItem.SubNetID)
		mongodbInst.MaintainTime = types.StringValue(mongodbItem.MaintainTime)
		mongodbInst.Subnet = types.StringValue(mongodbItem.Subnet)
		mongodbInst.LogStatus = types.BoolValue(mongodbItem.LogStatus)
		mongodbInst.OrderId = types.Int64Value(mongodbItem.OrderId)
		mongodbInst.NetName = types.StringValue(mongodbItem.NetName)
		mongodbInst.VersionNum = utils.SecStringValue(mongodbItem.VersionNum)
		mongodbInst.SecurityGroupId = types.StringValue(mongodbItem.SecurityGroupId)
		mongodbInst.ParameterConfigsvrGroupUsed = utils.SecStringValue(mongodbItem.ParameterConfigsvrGroupUsed)
		mongodbInst.DiskSize = types.Int32Value(mongodbItem.DiskSize)
		mongodbInst.TplName = types.StringValue(mongodbItem.TplName)
		mongodbInst.ProdInstSetName = types.StringValue(mongodbItem.ProdInstSetName)
		mongodbInst.Released = types.BoolValue(mongodbItem.Released)
		mongodbInst.SecurityGroup = types.StringValue(mongodbItem.SecurityGroup)
		mongodbInst.ProdType = types.Int32Value(mongodbItem.ProdType)
		mongodbInst.ExpireTime = types.Int64Value(mongodbItem.ExpireTime)
		mongodbInst.ProdInstId = types.StringValue(mongodbItem.ProdInstId)
		mongodbInst.ProjectName = types.StringValue(mongodbItem.ProjectName)
		mongodbInst.ProjectId = types.StringValue(mongodbItem.ProjectId)
		mongodbInst.DestroyedTime = types.StringValue(mongodbItem.DestroyedTime)
		mongodbInst.ProdInstFlag = types.StringValue(mongodbItem.ProdInstFlag)
		mongodbInst.ProdDbEngine = types.StringValue(mongodbItem.ProdDbEngine)
		mongodbInst.BillMode = types.Int32Value(mongodbItem.BillMode)
		mongodbInst.ProdId = types.StringValue(mongodbItem.ProdId)
		mongodbInst.RestoreTime = types.StringValue(mongodbItem.RestoreTime)
		mongodbInst.ProdRunningStatus = types.Int32Value(mongodbItem.ProdRunningStatus)
		mongodbInst.DiskUsed = utils.SecStringValue(mongodbItem.DiskUsed)
		mongodbInst.ParameterGroupUsed = types.StringValue(mongodbItem.ParameterGroupUsed)
		mongodbInst.VpcId = types.StringValue(mongodbItem.VpcId)
		mongodbInst.InnodbThreadConcurrency = types.Int64Value(mongodbItem.InnodbThreadConcurrency)
		mongodbInst.DiskType = types.StringValue(mongodbItem.DiskType)
		mongodbInst.ProdBillType = types.Int32Value(mongodbItem.ProdBillType)
		mongodbInst.MachineSpec = types.StringValue(mongodbItem.MachineSpec)
		mongodbInst.ProdInstName = types.StringValue(mongodbItem.ProdInstName)
		mongodbInst.InnodbBufferPoolSize = types.StringValue(mongodbItem.InnodbBufferPoolSize)
		mongodbInst.UsedSpace = utils.SecStringValue(mongodbItem.UsedSpace)
		mongodbInst.UserId = types.Int64Value(mongodbItem.UserId)
		mongodbInst.ProdBillTime = types.Int32Value(mongodbItem.ProdBillTime)
		mongodbInst.Destroyed = types.BoolValue(mongodbItem.Destroyed)
		mongodbInst.CreateTime = types.Int64Value(mongodbItem.CreateTime)
		mongodbInst.TenantId = types.Int64Value(mongodbItem.TenantId)
		mongodbInst.OuterId = types.StringValue(mongodbItem.OuterId)
		mongodbInst.TplCode = types.StringValue(mongodbItem.TplCode)
		if mongodbItem.ReadPort != nil {
			mongodbInst.ReadPort = types.Int32Value(*mongodbItem.ReadPort)
		}
		mongodbInstances = append(mongodbInstances, mongodbInst)
	}
	config.MongodbInstances = mongodbInstances
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunMongodbInstanceModel struct {
	ProdOrderStatus             types.Int32  `tfsdk:"prod_order_status"`              //0->订单正常,1->订单冻结,2->订单注销,3->施工中,4->施工失败
	SubNetID                    types.String `tfsdk:"subnet_id"`                      //子网ID
	MaintainTime                types.String `tfsdk:"maintain_time"`                  //可维护时间
	Subnet                      types.String `tfsdk:"subnet"`                         //子网
	LogStatus                   types.Bool   `tfsdk:"log_status"`                     //实例日志审计状态
	OrderId                     types.Int64  `tfsdk:"order_id"`                       //订单ID
	NetName                     types.String `tfsdk:"net_name"`                       //专有网络
	VersionNum                  types.String `tfsdk:"version_num"`                    //版本号
	SecurityGroupId             types.String `tfsdk:"security_group_id"`              //安全组ID
	ParameterConfigsvrGroupUsed types.String `tfsdk:"parameter_configsvr_group_used"` //参数配置
	DiskSize                    types.Int32  `tfsdk:"disk_size"`                      //存储空间大小
	TplName                     types.String `tfsdk:"tpl_name"`                       //模板名称
	ProdInstSetName             types.String `tfsdk:"prod_inst_set_name"`             //实例对应的SET名
	Released                    types.Bool   `tfsdk:"released"`                       //实例是否已被释放
	SecurityGroup               types.String `tfsdk:"security_group"`                 //安全组
	ProdType                    types.Int32  `tfsdk:"prod_type"`                      //0:单机,2:副本集(三节点),4:副本集(五节点),6:副本集(七节点),10:分片集群
	ExpireTime                  types.Int64  `tfsdk:"expire_time"`                    //到期时间
	ProdInstId                  types.String `tfsdk:"prod_inst_id"`                   //实例id
	ProjectName                 types.String `tfsdk:"project_name"`                   //企业项目名称
	ProjectId                   types.String `tfsdk:"project_id"`                     //企业项目id
	DestroyedTime               types.String `tfsdk:"destroyed_time"`                 //实例销毁时间
	ProdInstFlag                types.String `tfsdk:"prod_inst_flag"`                 //规定为"实例ID 实例名称"
	ProdDbEngine                types.String `tfsdk:"prod_db_engine"`                 //dds数据库产品的版本
	BillMode                    types.Int32  `tfsdk:"bill_mode"`                      //1:包周期计费,2:按需计费
	ProdId                      types.String `tfsdk:"prod_id"`                        //产品表示
	RestoreTime                 types.String `tfsdk:"restore_time"`                   //实例恢复时间
	ProdRunningStatus           types.Int32  `tfsdk:"prod_running_status"`            //实例运行状态:0->运行正常,1->重启中,2->备份操作中,3->恢复操作中,4->转换ssl,5->异常,6->修改参数组中,7->已冻结,8->已注销,9->施工中,10->施工失败,11->扩容中,12->主备切换中
	DiskUsed                    types.String `tfsdk:"disk_used"`                      //磁盘空间
	ParameterGroupUsed          types.String `tfsdk:"parameter_group_used"`           //参数组名称，标明参数组的版本
	VpcId                       types.String `tfsdk:"vpc_id"`                         //vpc网络ID
	InnodbThreadConcurrency     types.Int64  `tfsdk:"innodb_thread_concurrency"`      //线程数
	DiskType                    types.String `tfsdk:"disk_type"`                      //存储类型
	ProdBillType                types.Int32  `tfsdk:"prod_bill_type"`                 //0:按月计费,1:按天计费,2:按年计费,3:按流量计费
	MachineSpec                 types.String `tfsdk:"machine_spec"`                   //CPU内存规格
	ProdInstName                types.String `tfsdk:"prod_inst_name"`                 //实例名称
	InnodbBufferPoolSize        types.String `tfsdk:"innodb_buffer_pool_size"`        //缓存池大小
	UsedSpace                   types.String `tfsdk:"used_space"`                     //已使用空间
	UserId                      types.Int64  `tfsdk:"user_id"`                        //用户id
	ProdBillTime                types.Int32  `tfsdk:"prod_bill_time"`                 //购买时长
	Destroyed                   types.Bool   `tfsdk:"destroyed"`                      //实例是否已经被销毁
	CreateTime                  types.Int64  `tfsdk:"create_time"`                    //创建时间
	TenantId                    types.Int64  `tfsdk:"tenant_id"`                      //租户id
	OuterId                     types.String `tfsdk:"outer_id"`                       //产品ID
	TplCode                     types.String `tfsdk:"tpl_code"`                       //模板编码
	ReadPort                    types.Int32  `tfsdk:"read_port"`                      //读端口
}

type CtyunMongodbInstancesConfig struct {
	PageNo           types.Int32                 `tfsdk:"page_no"`        // 当前页，不传默认为1
	PageSize         types.Int32                 `tfsdk:"page_size"`      // 页大小，不传默认为10
	ResDbEngine      types.String                `tfsdk:"res_db_engine"`  // 版本号
	ProdInstName     types.String                `tfsdk:"prod_inst_name"` // 实例名称
	LabelIds         types.String                `tfsdk:"label_ids"`      // 标签id
	RegionID         types.String                `tfsdk:"region_id"`      // 资源池id
	ProjectID        types.String                `tfsdk:"project_id"`     // 项目id
	MongodbInstances []CtyunMongodbInstanceModel `tfsdk:"mongodb_instances"`
}
