package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &CtyunMysqlInstance{}
	_ resource.ResourceWithConfigure   = &CtyunMysqlInstance{}
	_ resource.ResourceWithImportState = &CtyunMysqlInstance{}
)

type CtyunMysqlInstance struct {
	meta         *common.CtyunMetadata
	ecsService   *business.EcsService
	mysqlService *business.MysqlService
}

func (c *CtyunMysqlInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunMysqlInstance) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ecsService = business.NewEcsService(c.meta)
	c.mysqlService = business.NewMysqlService(c.meta)
}

func NewCtyunMysqlInstance() resource.Resource {
	return &CtyunMysqlInstance{}
}

func (c *CtyunMysqlInstance) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_instance"
}

func (c *CtyunMysqlInstance) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10033813/10134365**`,
		Attributes: map[string]schema.Attribute{
			"flavor_name": schema.StringAttribute{
				Required:    true,
				Description: "规格名称，形如c7.2xlarge.4，可从data.ctyun_mysql_specs查询支持的规格，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围：month：按月，on_demand：按需。当此值为month时，cycle_count为必填",
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypeOnDemand, business.OrderCycleTypeMonth),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_count": schema.Int32Attribute{
				Optional:    true,
				Description: "订购时长，该参数当且仅当在cycle_type为month时填写，支持传递1-36",
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeMonth),
					),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
					int32validator.Between(1, 36),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"auto_renew": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否自动续订，默认非自动续订，当cycle_type不等于on_demand时才可填写，当cycle_count<12，到期自动续订1个月，当cycle_count>=12，到期自动续订12个月",
				Default:     booldefault.StaticBool(false),
				Validators: []validator.Bool{
					validator2.ConflictsWithEqualBool(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
				},
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id,如果不填这默认使用provider ctyun总region_id 或者环境变量",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "虚拟私有云Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"security_group_id": schema.StringAttribute{
				Required:    true,
				Description: "安全组Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SecurityGroupValidate(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "实例名称（长度在 4 到 64个字符，必须以字母开头，不区分大小写，可以包含字母、数字、中划线或下划线，不能包含其他特殊字符）",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(4, 64),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z][0-9a-zA-Z_-]+$"), "终端节点服务名称不符合规则"),
				},
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "实例密码，密码为8-26位，需为字母、数字和特殊字符~!@#$%^*_-+{[]}:,.?/的组合，区分大小写",
				Validators: []validator.String{
					validator2.DBPassword(
						8,
						26,
						3,
						"MYSQL",
						"~!@#$%^*_-+{[]}:,.?/",
					),
				},
			},
			"prod_id": schema.StringAttribute{
				Required:    true,
				Description: "产品id，支持更新。取值范围：Single57（单实例5.7版本）, Single80（单实例8.0版本）, MasterSlave57（一主一备5.7版本）, MasterSlave80（一主一备8.0版本）, Master2Slave57（一主两备5.7版本）, Master2Slave80（一主两备8.0版本）。在更新时，不支持prod_id（节点）和prod_performance_spec（规格）同时更新。",
				Validators: []validator.String{
					stringvalidator.OneOf(business.MysqlProdIds...),
				},
			},
			"storage_type": schema.StringAttribute{
				Required:    true,
				Description: "存储类型: SSD=超高IO、SATA=普通IO、SAS=高IO、SSD-genric=通用型SSD、FAST-SSD=极速型SSD",
				Validators: []validator.String{
					stringvalidator.OneOf(business.StorageType...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_space": schema.Int32Attribute{
				Required:    true,
				Description: "存储空间(单位:G，范围100,32768)，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(100, 32768),
				},
			},
			"backup_storage_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "备份空间磁盘存储类型：SSD=超高IO、SATA=普通IO、SAS=高IO",
				Validators: []validator.String{
					stringvalidator.OneOf(business.StorageTypeSSD, business.StorageTypeSATA, business.StorageTypeSAS),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString(business.StorageTypeSATA),
			},
			"backup_storage_space": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "备份存储空间(单位:G，范围100,32768)，若storage_space和backup_storage_space都不为空，优先升配备份节点存储空间，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(100, 32768),
				},
				Default: int32default.StaticInt32(100),
			},
			"availability_zone_info": schema.ListNestedAttribute{
				Optional:    true,
				Description: "可用区信息，需要根据prod_id而定。创建阶段，需要指定master和slave的所在az。例：若一主一备，需要传参：[｛'availability_zone_name':'xxxx', 'availability_zone_count':1,node_type:'master'｝,｛'availability_zone_name':'xxxx', 'availability_zone_count':1,node_type:'slave'｝]；在更新阶段，仅需要填写扩容部分的AZ信息。例：将单节点扩容至1主2备，[{'availability_zone_name':'xxxx', 'availability_zone_count':2,node_type:'slave'}]",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"availability_zone_name": schema.StringAttribute{
							Required:    true,
							Description: "资源池可用区名称",
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
						"availability_zone_count": schema.Int32Attribute{
							Required:    true,
							Description: "该AZ内存在的实例节点数量",
							Validators: []validator.Int32{
								int32validator.Between(1, 16),
							},
						},
						"node_type": schema.StringAttribute{
							Required:    true,
							Description: "表示分布AZ的节点类型，master/slave",
							Validators: []validator.String{
								stringvalidator.OneOf("master", "slave"),
							},
						},
					},
				},
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订单id",
			},
			"inst_id": schema.StringAttribute{
				Computed:    true,
				Description: "实例Id",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"prod_running_status": schema.Int32Attribute{
				Computed:    true,
				Description: "0.正常 1.重启中 2.备份中 3.恢复中 4.修改参数中 5.应用参数组中 6.扩容预处理中 7.扩容预处理完成 8.修改端口中 9.迁移中 10.重置密码中 11.修改数据复制方式中 12.缩容预处理中 13.缩容预处理完成 15.内核小版本升级 17.迁移可用区中 18.修改备份配置中 20.停止中 21.已停止 22.启动中 26.白名单配置中",
				Validators: []validator.Int32{
					int32validator.Between(0, 26),
				},
			},
			"vip": schema.StringAttribute{
				Computed:    true,
				Description: "虚拟IP地址",
			},
			"write_port": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "写数据端口，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(0, 65535),
				},
			},
			"read_port": schema.StringAttribute{
				Computed:    true,
				Description: "读端口",
			},
			"prod_db_engine": schema.StringAttribute{
				Computed:    true,
				Description: "数据库引擎",
			},
			"eip": schema.StringAttribute{
				Computed:    true,
				Description: "弹性ip",
			},
			"eip_status": schema.Int32Attribute{
				Computed:    true,
				Description: "弹性ip状态 0->unbind，1->bind,2->binding",
			},
			"ssl_status": schema.Int32Attribute{
				Computed:    true,
				Description: "Ssl状态 0->off，1->on",
			},
			"new_mysql_version": schema.StringAttribute{
				Computed:    true,
				Description: "mysql版本",
			},
			"audit_log_status": schema.Int32Attribute{
				Computed:    true,
				Description: "日志审计开关",
			},
			"inst_release_protection_status": schema.Int32Attribute{
				Computed:    true,
				Description: "实例释放保护开关 1:on,0:off",
			},
			"pause_enable": schema.BoolAttribute{
				Computed:    true,
				Description: "是否允许暂停",
			},
			"mysql_port": schema.StringAttribute{
				Computed:    true,
				Description: "数据库端口",
			},
			"security_group_status": schema.Int32Attribute{
				Computed:    true,
				Description: "安全组状态 0->normal, 1->changing, 2->deleted",
			},
			"running_control": schema.StringAttribute{
				Optional:    true,
				Description: "控制是否暂停，启用和重启实例，支持更新，取值范围：freeze, unfreeze, restart",
				Validators: []validator.String{
					stringvalidator.OneOf("freeze", "unfreeze", "restart"),
				},
			},
			"prod_order_status": schema.Int32Attribute{
				Computed:    true,
				Description: "0.正常 1.欠费暂停 2.已注销 3.创建中 4.施工失败 5.到期退订状态 6.新增的状态-openApi暂停 7.创建完成等待变更单 8.待注销 9.手动暂停 10.手动退订",
			},
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "实例Id，同inst_id",
			},
		},
	}
}

func (c *CtyunMysqlInstance) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMysqlInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建前检查
	err = c.checkSpec(ctx, &plan)
	if err != nil {
		return
	}
	// 开始创建
	err = c.CreateMysqlInstance(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取mysql详情
	err = c.getAndMergeMysqlInstance(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlInstance) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunMysqlInstanceConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergeMysqlInstance(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlInstance) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置
	var plan CtyunMysqlInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunMysqlInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	if !plan.Password.Equal(state.Password) {
		err = fmt.Errorf("数据库密码暂时不支持修改")
		return
	}

	// 校验规格
	err = c.checkSpec(ctx, &plan)
	if err != nil {
		return
	}
	err = c.updateMysqlInstance(ctx, &state, &plan)
	if err != nil {
		return
	}
	state.FlavorName = plan.FlavorName
	time.Sleep(30 * time.Second)
	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeMysqlInstance(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlInstance) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunMysqlInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 确保主机在退订之前是处于running状态
	err = c.StartedLoop(ctx, &state)
	if err != nil {
		return
	}

	err = c.refund(ctx, state)
	if err != nil {
		return
	}
	// 轮询确认时候退订成功
	err = c.refundLoop(ctx, state)
	if err != nil {
		return
	}
	time.Sleep(30 * time.Second)
	err = c.destroy(ctx, state)
	if err != nil {
		return
	}
	err = c.destroyLoop(ctx, state)
	if err != nil {
		return
	}
	response.Diagnostics.AddWarning("删除MySql集群成功", "集群退订后，若立即删除子网或安全组可能会失败，需要等待底层资源释放")
}

// CreateMysqlInstance 创建mysql实例
func (c *CtyunMysqlInstance) CreateMysqlInstance(ctx context.Context, config *CtyunMysqlInstanceConfig) (err error) {
	cycleType := config.CycleType.ValueString()
	params := &mysql.TeledbCreateRequest{
		BillMode:        business.MysqlBillMode[cycleType],
		RegionId:        config.RegionID.ValueString(),
		ProdVersion:     business.MysqlProdVersionDict[config.ProdID.ValueString()],
		VpcId:           config.VpcID.ValueString(),
		HostType:        config.hostType,
		SubnetId:        config.SubnetID.ValueString(),
		SecurityGroupId: config.SecurityGroupID.ValueString(),
		Name:            config.Name.ValueString(),
		Period:          config.CycleCount.ValueInt32(),
		Count:           1,
		ProdId:          business.MysqlProdIdDict[config.ProdID.ValueString()],
		CpuType:         business.MysqlCpuTypeDict[config.cpuType],
		OsType:          business.MysqlOSTypeDict[config.osType],
	}
	if !config.Password.IsNull() && !config.Password.IsUnknown() {
		password := business.Encode(config.Password.ValueString())
		params.Password = password
	}
	if cycleType == business.OnDemandCycleType {
		params.AutoRenewStatus = 0
	} else {
		params.AutoRenewStatus = map[bool]int32{true: 1, false: 0}[config.AutoRenew.ValueBool()]
	}

	header := &mysql.TeledbCreateRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}

	var MysqlNodeInfos []mysql.MysqlNodeInfoListRequest

	mysqlNodeInfo := mysql.MysqlNodeInfoListRequest{}
	mysqlNodeInfo.NodeType = business.NodeTypeDict[config.ProdID.ValueString()]
	mysqlNodeInfo.InstSpec = business.MysqlInstanceSeriesDict[config.instanceSeries]
	mysqlNodeInfo.StorageType = config.StorageType.ValueString()
	mysqlNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()
	mysqlNodeInfo.ProdPerformanceSpec = config.prodPerformanceSpec
	mysqlNodeInfo.BackupStorageType = config.BackupStorageType.ValueString()
	mysqlNodeInfo.BackupStorageSpace = config.BackupStorageSpace.ValueInt32()
	mysqlNodeInfo.Disks = 1
	// 处理availabilityZoneInfo可用区信息

	var availabilityZoneInfos []mysql.AvailabilityZoneInfoRequest
	err = c.generateAzInfos(ctx, config, &availabilityZoneInfos)
	if err != nil {
		return
	}

	mysqlNodeInfo.AvailabilityZoneInfo = availabilityZoneInfos
	MysqlNodeInfos = append(MysqlNodeInfos, mysqlNodeInfo)
	params.MysqlNodeInfoList = MysqlNodeInfos

	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbCreateApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 保存orderId
	if resp.ReturnObj.Data.NewOrderId == nil {
		err = errors.New("订单id为空，创建有误！")
		return
	}

	config.MasterOrderID = types.StringValue(*resp.ReturnObj.Data.NewOrderId)
	return
}

func (c *CtyunMysqlInstance) getAndMergeMysqlInstance(ctx context.Context, config *CtyunMysqlInstanceConfig) (err error) {
	// 若实例id为空，可能是因为实例刚创建，需要通过查询列表获取
	if config.InstID.ValueString() == "" {
		mysqlListParams := &mysql.TeledbGetListRequest{
			PageNow:      1,
			PageSize:     100,
			ProdInstName: config.Name.ValueStringPointer(),
		}
		mysqlListHeaders := &mysql.TeledbGetListHeaders{
			RegionID: config.RegionID.ValueString(),
		}
		if config.ProjectID.ValueString() != "" {
			mysqlListHeaders.ProjectID = config.ProjectID.ValueStringPointer()
		}

		resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbGetListApi.Do(ctx, c.meta.Credential, mysqlListParams, mysqlListHeaders)
		if err2 != nil {
			err = err2
			return
		}
		if len(resp.ReturnObj.List) > 1 {
			err = errors.New("实例名重复！")
			return
		} else if len(resp.ReturnObj.List) < 1 {
			//若根据name查询不到机器，可能存在还未创建好的情况，需要轮询
			resp, err = c.ListLoop(ctx, mysqlListParams, mysqlListHeaders, 60)
			if err != nil {
				return
			}
			if len(resp.ReturnObj.List) != 1 {
				err = errors.New("未查询该实例mysql，mysql name:" + config.Name.ValueString())
				return
			}
		}
		config.InstID = types.StringValue(resp.ReturnObj.List[0].OuterProdInstId)
		config.ID = config.InstID
		// 确认资源是否开通完成
		// 若暂未开通完成，需要轮询等待
		if resp.ReturnObj.List[0].ProdOrderStatus != business.MysqlOrderStatusStarted {
			err = c.CreateLoop(ctx, mysqlListParams, mysqlListHeaders)
			if err != nil {
				return err
			}
		}
	}
	// 获取实例详情
	if config.InstID.ValueString() == "" {
		err = errors.New("查询实例详情时，实例 ID为空")
		return err
	}
	detailParams := &mysql.TeledbQueryDetailRequest{
		OuterProdInstId: config.InstID.ValueString(),
	}
	detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		detailHeaders.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
	if err != nil {
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 处理实例详情
	returnOjb := resp.ReturnObj
	config.ProdRunningStatus = types.Int32Value(returnOjb.ProdRunningStatus)
	config.ProdOrderStatus = types.Int32Value(returnOjb.ProdOrderStatus)
	config.Vip = types.StringValue(returnOjb.Vip)
	config.ReadPort = types.StringValue(returnOjb.ReadPort)
	config.ProdDbEngine = types.StringValue(returnOjb.ProdDbEngine)
	config.EIP = types.StringValue(returnOjb.EIP)
	config.EipStatus = types.Int32Value(returnOjb.EIPStatus)
	config.SSlStatus = types.Int32Value(returnOjb.SSlStatus)
	config.NewMysqlVersion = types.StringValue(returnOjb.NewMysqlVersion)
	config.AuditLogStatus = types.Int32Value(returnOjb.AuditLogStatus)
	config.InstReleaseProtectionStatus = types.Int32Value(returnOjb.InstReleaseProtectionStatus)
	config.PauseEnable = types.BoolValue(returnOjb.PauseEnable)
	config.MysqlPort = types.StringValue(returnOjb.MysqlPort)
	config.SecurityGroupStatus = types.Int32Value(returnOjb.SecurityGroupStatus)
	config.Name = types.StringValue(returnOjb.ProdInstName)
	writePort, err := strconv.ParseInt(returnOjb.WritePort, 10, 32)
	if err != nil {
		return
	}
	config.WritePort = types.Int32Value(int32(writePort))

	// 更新disk， 主机配置相关信息
	config.ProdID = types.StringValue(business.MysqlProdIdRevDict[returnOjb.ProdId])

	config.StorageSpace = types.Int32Value(returnOjb.DiskSize)
	config.BackupStorageSpace = types.Int32Value(returnOjb.BackupDiskSize)
	return
}

func (c *CtyunMysqlInstance) CreateLoop(ctx context.Context, ListParams *mysql.TeledbGetListRequest, ListHeaders *mysql.TeledbGetListHeaders, loopCount ...int) (err error) {

	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbGetListApi.Do(ctx, c.meta.Credential, ListParams, ListHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}

			status := resp.ReturnObj.List[0].ProdOrderStatus
			switch status {
			case business.MysqlOrderStatusStarted:
				return false
			case business.MysqlOrderStatusCreating:
				return true
			case business.MysqlOrderStatusWaiting:
				return true
			case business.MysqlRunningStatusBackup:
				return true
			default:
				// 在开通的时候，其他状态是异常的，因此抛出异常，并跳出轮询
				err = errors.New("mysql创建状态有误： " + fmt.Sprintf("%d", status))
				return false
			}
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未创建成功！")
	}
	return
}

func (c *CtyunMysqlInstance) ListLoop(ctx context.Context, params *mysql.TeledbGetListRequest, headers *mysql.TeledbGetListHeaders, loopCount ...int) (*mysql.TeledbGetListResponse, error) {
	var err error
	var response *mysql.TeledbGetListResponse
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return nil, err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbGetListApi.Do(ctx, c.meta.Credential, params, headers)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}

			if len(resp.ReturnObj.List) > 1 {
				err = fmt.Errorf("查询到多条为名为%s的记录！", *params.ProdInstName)
				return false
			}
			if len(resp.ReturnObj.List) == 1 {
				response = resp
				return false
			}
			// 未查询到，继续轮询
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return nil, errors.New("轮询已达最大次数，资源仍未创建或查询到！")
	}
	return response, nil
}

func (c *CtyunMysqlInstance) UpgradeLoop(ctx context.Context, state *CtyunMysqlInstanceConfig, plan *CtyunMysqlInstanceConfig, loopCount ...int) (err error) {

	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			// 获取实例详情
			detailParams := &mysql.TeledbQueryDetailRequest{
				OuterProdInstId: state.InstID.ValueString(),
			}
			detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
				InstID:   state.InstID.ValueString(),
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			runningStatus := resp.ReturnObj.ProdRunningStatus
			orderStatus := resp.ReturnObj.ProdOrderStatus
			// 若符合预期，跳出循环，扩容成功
			if resp.ReturnObj.ProdId == business.MysqlProdIdDict[plan.ProdID.ValueString()] && resp.ReturnObj.DiskSize == plan.StorageSpace.ValueInt32() && resp.ReturnObj.MachineSpec == plan.prodPerformanceSpec {
				//若备份磁盘空间不为空，且预期的分配磁盘空间与远端磁盘备份空间不相同，则继续轮询
				if plan.BackupStorageSpace.ValueInt32() != 0 && plan.BackupStorageSpace.ValueInt32() != resp.ReturnObj.BackupDiskSize {
					return true
				}
				if runningStatus == business.MysqlRunningStatusStarted && orderStatus == business.MysqlOrderStatusStarted {
					return false
				} else {
					return true
				}
			}
			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未升级成功！")
	}
	return
}

func (c *CtyunMysqlInstance) RunningStatusLoop(ctx context.Context, config *CtyunMysqlInstanceConfig, runningStatus int32, orderStatus int32, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			mysqlListParams := &mysql.TeledbGetListRequest{
				PageNow:      1,
				PageSize:     100,
				ProdInstName: config.Name.ValueStringPointer(),
			}
			mysqlListHeaders := &mysql.TeledbGetListHeaders{
				RegionID: config.RegionID.ValueString(),
			}
			if config.ProjectID.ValueString() != "" {
				mysqlListHeaders.ProjectID = config.ProjectID.ValueStringPointer()
			}

			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbGetListApi.Do(ctx, c.meta.Credential, mysqlListParams, mysqlListHeaders)
			if err2 != nil {
				err = err2
				return false
			}
			if len(resp.ReturnObj.List) > 1 {
				err = errors.New("实例名重复！")
				return false
			} else if len(resp.ReturnObj.List) < 1 {
				//若根据name查询不到机器，可能存在还未创建好的情况，需要轮询
				resp, err = c.ListLoop(ctx, mysqlListParams, mysqlListHeaders, 60)
				if err != nil {
					return false
				}
				if len(resp.ReturnObj.List) != 1 {
					err = errors.New("未查询该实例mysql，mysql name:" + config.Name.ValueString())
					return false
				}
			}

			currentRunningStatus := resp.ReturnObj.List[0].ProdRunningStatus
			currentOrderStatus := resp.ReturnObj.List[0].ProdOrderStatus
			if currentOrderStatus == orderStatus && currentRunningStatus == runningStatus {
				return false
			}
			return true

		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍完成状态更新！")
	}
	return
}

func (c *CtyunMysqlInstance) updateInfoLoop(ctx context.Context, state *CtyunMysqlInstanceConfig, plan *CtyunMysqlInstanceConfig, loopCount ...int) (err error) {

	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			// 获取实例详情
			detailParams := &mysql.TeledbQueryDetailRequest{
				OuterProdInstId: state.InstID.ValueString(),
			}
			detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
				InstID:   state.InstID.ValueString(),
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			//status := resp.ReturnObj.ProdRunningStatus
			// 跳出轮询条件如下：
			// 当state.name = plan.name，并且write_port无须更新时
			// 当state.name = plan.name，且write_port符合预期时
			if resp.ReturnObj.ProdInstName == plan.Name.ValueString() {
				if plan.WritePort.ValueInt32() == 0 {
					return false
				} else {
					if resp.ReturnObj.WritePort == fmt.Sprintf("%d", plan.WritePort.ValueInt32()) {
						return false
					} else {
						return true
					}
				}
			}
			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未更新成功！")
	}
	return
}

func (c *CtyunMysqlInstance) StartedLoop(ctx context.Context, state *CtyunMysqlInstanceConfig, loopCount ...int) (err error) {
	count := 30
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	var cnt int
	result := retryer.Start(
		func(currentTime int) bool {
			// 获取实例详情
			detailParams := &mysql.TeledbQueryDetailRequest{
				OuterProdInstId: state.InstID.ValueString(),
			}
			detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
				InstID:   state.InstID.ValueString(),
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			runningStatus := resp.ReturnObj.ProdRunningStatus
			orderStatus := resp.ReturnObj.ProdOrderStatus
			// 若变配前，发现数据库已冻结，将其恢复
			if orderStatus == business.MysqlOrderStatusPause {
				err = c.startMysqlInstance(ctx, state, nil)
				if err != nil {
					return false
				}
			}
			if runningStatus == business.MysqlRunningStatusStarted && orderStatus == business.MysqlRunningStatusStarted {
				// 有三次是start，才认为状态正常
				cnt++
				if cnt > 3 {
					return false
				}
			}
			if orderStatus == business.MysqlOrderStatusPause {
				err = errors.New("订单处于暂停状态，不可进行变更操作")
				return false
			}
			if runningStatus == business.MysqlRunningStatusStopping || runningStatus == business.MysqlRunningStatusStopped {
				err = errors.New("主机处于关机状态，不可进行变更操作")
				return false
			}

			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未到达启动状态！")
	}
	return
}

// refund 退订
func (c *CtyunMysqlInstance) refund(ctx context.Context, state CtyunMysqlInstanceConfig) (err error) {
	deleteParams := &mysql.TeledbRefundRequest{
		InstId: state.InstID.ValueString(),
	}
	deleteHeader := &mysql.TeledbRefundRequestHeader{}
	if state.ProjectID.ValueString() != "" {
		deleteHeader.ProjectID = state.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbRefundApi.Do(ctx, c.meta.Credential, deleteParams, deleteHeader)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// refundLoop 退订后检查
func (c *CtyunMysqlInstance) refundLoop(ctx context.Context, state CtyunMysqlInstanceConfig, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			params := &mysql.TeledbGetListRequest{
				PageNow:      1,
				PageSize:     100,
				ProdInstName: state.Name.ValueStringPointer(),
			}
			headers := &mysql.TeledbGetListHeaders{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				headers.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbGetListApi.Do(ctx, c.meta.Credential, params, headers)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			// 若查询列表已经查询不到，资源已经销毁
			if len(resp.ReturnObj.List) == 0 {
				return false
			}
			status := resp.ReturnObj.List[0].ProdOrderStatus
			switch status {
			case business.MysqlOrderStatusDestroy:
				return false
			case business.MysqlOrderStatusDestroyed:
				return false
			case business.MysqlOrderStatusStarted:
				return true
			case business.MysqlOrderStatusPause:
				return true
			default:
				err = errors.New("退订状态有误，当前状态为：" + fmt.Sprintf("%d", status))
				return false
			}
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未退订成功！")
	}
	return
}

// destroy 销毁
func (c *CtyunMysqlInstance) destroy(ctx context.Context, state CtyunMysqlInstanceConfig) (err error) {
	deleteParams := &mysql.TeledbDestroyRequest{
		InstId: state.InstID.ValueString(),
	}
	deleteHeader := &mysql.TeledbDestroyRequestHeader{}
	if state.ProjectID.ValueString() != "" {
		deleteHeader.ProjectID = state.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbDestroyApi.Do(ctx, c.meta.Credential, deleteParams, deleteHeader)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// destroyLoop 销毁后检查
func (c *CtyunMysqlInstance) destroyLoop(ctx context.Context, state CtyunMysqlInstanceConfig, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			params := &mysql.TeledbGetListRequest{
				PageNow:      1,
				PageSize:     100,
				ProdInstName: state.Name.ValueStringPointer(),
			}
			headers := &mysql.TeledbGetListHeaders{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				headers.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbGetListApi.Do(ctx, c.meta.Credential, params, headers)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			// 若查询列表已经查询不到，资源已经销毁
			if len(resp.ReturnObj.List) == 0 {
				return false
			}
			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未退订成功！")
	}
	return
}

func (c *CtyunMysqlInstance) updateMysqlInstance(ctx context.Context, state *CtyunMysqlInstanceConfig, plan *CtyunMysqlInstanceConfig) (err error) {
	if state.InstID.ValueString() == "" {
		err = errors.New("变配实例时，实例ID为空！")
		return err
	}
	// 修改实例名称
	if plan.Name.ValueString() != "" && state.Name.ValueString() != plan.Name.ValueString() {
		updateNameParams := &mysql.TeledbUpdateInstanceNameRequest{
			OuterProdInstID:     state.InstID.ValueString(),
			InstanceDescription: plan.Name.ValueString(),
		}
		updatedNameHeaders := &mysql.TeledbUpdateInstanceNameRequestHeader{
			InstID:   state.InstID.ValueString(),
			RegionID: state.RegionID.ValueString(),
		}
		if state.ProjectID.ValueString() != "" {
			updatedNameHeaders.ProjectID = state.ProjectID.ValueStringPointer()
		}
		resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbUpdateInstanceNameApi.Do(ctx, c.meta.Credential, updateNameParams, updatedNameHeaders)
		if err2 != nil {
			err = err2
			return
		} else if resp.StatusCode != 0 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
	}

	// 修改实例写端口
	if plan.WritePort.ValueInt32() != 0 && state.WritePort.ValueInt32() != plan.WritePort.ValueInt32() {
		// 更新之前需要确定主机状态必须为started
		err = c.StartedLoop(ctx, state)
		if err != nil {
			return
		}
		updateWritePortParams := &mysql.TeledbUpdateWritePortRequest{
			OuterProdInstId: state.InstID.ValueString(),
			WritePort:       fmt.Sprintf("%d", plan.WritePort.ValueInt32()),
		}
		updateWritePortHeaders := &mysql.TeledbUpdateWritePortRequestHeader{
			InstID:   state.InstID.ValueString(),
			RegionID: state.RegionID.ValueString(),
		}
		if state.ProjectID.ValueString() != "" {
			updateWritePortHeaders.ProjectID = state.ProjectID.ValueString()
		}
		resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbUpdateWritePortApi.Do(ctx, c.meta.Credential, updateWritePortParams, updateWritePortHeaders)
		if err2 != nil {
			return err2
		} else if resp.StatusCode != 0 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
	}
	// 轮询基础信息是否修改成功
	err = c.updateInfoLoop(ctx, state, plan)
	if err != nil {
		return
	}
	nodeType := business.NodeTypeDict[plan.ProdID.ValueString()]
	upgradeParams := &mysql.TeledbUpgradeRequest{
		InstId:   state.InstID.ValueString(),
		NodeType: &nodeType,
	}
	upgradeHeader := &mysql.TeledbUpgradeRequestHeader{}
	if plan.ProjectID.ValueString() != "" {
		upgradeHeader.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	// 若BackupStorageSpace不为空，触发备节点扩容存储空间
	if plan.BackupStorageSpace.ValueInt32() != 0 && state.BackupStorageSpace.ValueInt32() != plan.BackupStorageSpace.ValueInt32() {
		upgradeParams.DiskVolume = plan.BackupStorageSpace.ValueInt32Pointer()
		backupNodeType := business.PgsqlStorageTypeBackUp
		upgradeParams.NodeType = &backupNodeType

		err = c.upgradeMysqlStorage(ctx, state, plan, upgradeParams, upgradeHeader)
		if err != nil {
			return
		}
		upgradeParams.DiskVolume = nil
		upgradeParams.NodeType = &nodeType
	}

	// 若StorageSpace不为空，触发主节点扩容存储空间
	if plan.StorageSpace.ValueInt32() != 0 && state.StorageSpace.ValueInt32() != plan.StorageSpace.ValueInt32() {
		upgradeParams.DiskVolume = plan.StorageSpace.ValueInt32Pointer()

		err = c.upgradeMysqlStorage(ctx, state, plan, upgradeParams, upgradeHeader)
		if err != nil {
			return
		}
		upgradeParams.DiskVolume = nil
		upgradeParams.NodeType = &nodeType
	}

	// 扩容云数据库实例
	// 若plan.ProdPerformanceSpec不为空,且state和plan的ProdPerformanceSpec不一致，触发规格扩容
	if !plan.FlavorName.Equal(state.FlavorName) {
		if !plan.ProdID.Equal(state.ProdID) {
			err = errors.New("实例节点和规格(prod_id, flavor_name)不可同时变更")
			return
		}
		upgradeParams.ProdPerformanceSpec = &plan.prodPerformanceSpec
	}
	// 若plan.prodId不为空,且state和plan的prodId不一致，触发实例类型扩容
	if !plan.ProdID.IsNull() && state.ProdID.ValueString() != plan.ProdID.ValueString() {
		prodId := business.MysqlProdIdDict[plan.ProdID.ValueString()]
		upgradeParams.ProdId = &prodId
	}

	// 若实例扩容或更新ProdID---从单节点升级至，一主一备、一主两备。需要补充AZ信息
	if upgradeParams.ProdPerformanceSpec != nil || upgradeParams.ProdId != nil {
		var upgradeAzList []mysql.AvailabilityZoneInfo
		// 若AZ info 不为空，直接填写用户的输入
		if !plan.AvailabilityZoneInfo.IsNull() && !plan.AvailabilityZoneInfo.IsUnknown() {
			if !state.AvailabilityZoneInfo.IsNull() && state.AvailabilityZoneInfo.IsUnknown() && !plan.AvailabilityZoneInfo.Equal(state.AvailabilityZoneInfo) {
				err = errors.New("未变配实例规格或者实例节点时，az info不可修改！")
				return err
			}
			var azInfoList []AvailabilityZoneModel

			diag := plan.AvailabilityZoneInfo.ElementsAs(ctx, &azInfoList, true)
			if diag.HasError() {
				return
			}

			for _, azInfoItem := range azInfoList {
				azInfo := mysql.AvailabilityZoneInfo{
					AvailabilityZoneName:  azInfoItem.AvailabilityZoneName.ValueString(),
					AvailabilityZoneCount: azInfoItem.AvailabilityZoneCount.ValueInt32(),
				}
				upgradeAzList = append(upgradeAzList, azInfo)
			}
		} else {
			// 若az info 为空，直接生成
			err = c.getUpgradeAzInfo(ctx, state, plan, upgradeParams, &upgradeAzList)
			if err != nil {
				return
			} else if len(upgradeAzList) <= 0 {
				err = errors.New("mysql生成 az列表失败，可能存在问题：mysql实例暂不支持降配操作， 或查询资源池az信息失败等情况，可联系研发人员确定")
				return
			}
		}
		upgradeParams.AzList = upgradeAzList
	}
	// 若ProdPerformanceSpec, DiskVolume或者ProdId不为空时候，触发变配
	if upgradeParams.ProdPerformanceSpec != nil || upgradeParams.ProdId != nil {
		// 更新之前需要确定主机状态必须为started
		err = c.StartedLoop(ctx, state)
		if err != nil {
			return
		}
		resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbUpgradeApi.Do(ctx, c.meta.Credential, upgradeParams, upgradeHeader)
		if err2 != nil {
			err = err2
			return
		} else if resp.StatusCode != 200 {
			err = fmt.Errorf("API return error. Message: %s Error: %s", resp.Message, resp.Error)
			return
		}
		// 扩容后，轮循请求实例详情，确认已经完成升配
		err = c.UpgradeLoop(ctx, state, plan)
		if err != nil {
			return
		}
		// 扩容完成后，同步state AvailabilityZoneModel 状态
		if !plan.AvailabilityZoneInfo.IsNull() {
			state.AvailabilityZoneInfo = plan.AvailabilityZoneInfo
		}
	}

	// 启动实例
	if state.ProdOrderStatus.ValueInt32() == business.MysqlOrderStatusPause && plan.RunningControl.ValueString() == "unfreeze" {
		err = c.startMysqlInstance(ctx, state, plan)
		if err != nil {
			return
		}
	}

	// 停止实例
	if plan.RunningControl.ValueString() == "freeze" {
		// 进行重启、停止实例时，确保实例处于started状态
		err = c.StartedLoop(ctx, state)
		if err != nil {
			return
		}
		pauseParams := &mysql.TeledbStopRequest{
			OuterProdInstId: state.InstID.ValueString(),
		}
		pauseHeader := &mysql.TeledbStopRequestHeader{
			InstID:   state.InstID.ValueString(),
			RegionID: state.RegionID.ValueString(),
		}
		if state.ProjectID.ValueString() != "" {
			pauseHeader.ProjectID = state.ProjectID.ValueString()
		}
		resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbStopApi.Do(ctx, c.meta.Credential, pauseParams, pauseHeader)
		if err2 != nil {
			err = err2
			return err
		} else if resp.StatusCode != 0 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
		// 轮询验证，是否已停止，停止状态下，验证订单状态，预期=6
		err = c.RunningStatusLoop(ctx, state, business.MysqlRunningStatusStarted, business.MysqlOrderStatusPause, 60)
		if err != nil {
			return
		}
	}

	// 重启实例
	if plan.RunningControl.ValueString() == "restart" {
		// 进行重启、关机实例时，确保实例处于started状态
		err = c.StartedLoop(ctx, state)
		if err != nil {
			return
		}
		restartParams := &mysql.TeledbRestartRequest{
			OuterProdInstId: state.InstID.ValueString(),
		}
		restartHeader := &mysql.TeledbRestartRequestHeader{
			InstID:   state.InstID.ValueString(),
			RegionID: state.RegionID.ValueString(),
		}
		if state.ProjectID.ValueString() != "" {
			restartHeader.ProjectID = state.ProjectID.ValueString()
		}
		resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbRestartApi.Do(ctx, c.meta.Credential, restartParams, restartHeader)
		if err2 != nil {
			err = err2
			return err
		} else if resp.StatusCode != 0 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}
		//轮询验证，是否已完成重启
		err = c.RunningStatusLoop(ctx, state, business.MysqlRunningStatusStarted, business.MysqlOrderStatusStarted, 60)
		if err != nil {
			return
		}
	}
	state.RunningControl = plan.RunningControl
	return
}

func (c *CtyunMysqlInstance) startMysqlInstance(ctx context.Context, state *CtyunMysqlInstanceConfig, plan *CtyunMysqlInstanceConfig) (err error) {
	startParams := &mysql.TeledbStartRequest{
		OuterProdInstId: state.InstID.ValueString(),
	}
	startHeaders := &mysql.TeledbStartRequestHeader{
		InstID:   state.InstID.ValueString(),
		RegionID: state.RegionID.ValueString(),
	}
	if state.ProjectID.ValueString() != "" {
		startHeaders.ProjectID = state.ProjectID.ValueString()
	}
	resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbStartApi.Do(ctx, c.meta.Credential, startParams, startHeaders)
	if err2 != nil {
		err = err2
		return
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	// 轮询验证，是否已启动
	err = c.RunningStatusLoop(ctx, state, business.MysqlRunningStatusStarted, business.MysqlOrderStatusStarted, 60)
	if err != nil {
		return
	}
	return
}

func (c *CtyunMysqlInstance) generateAzInfos(ctx context.Context, config *CtyunMysqlInstanceConfig, availabilityZoneInfos *[]mysql.AvailabilityZoneInfoRequest) (err error) {
	if config.AvailabilityZoneInfo.IsNull() || config.AvailabilityZoneInfo.IsUnknown() {
		// 		 		1AZ 			2AZ				3个以上AZ
		// 单实例        AZ1  			AZ2				AZ3
		// 1主1备		AZ1				AZ1,AZ2			AZ1,2个AZ2
		// 1主2备		AZ1  			AZ1,2个AZ2      AZ1, AZ2, AZ3
		// 1. 判断实例类型，确认需要几个节点
		nodeNum := business.MysqlNodeNumDict[config.ProdID.ValueString()]
		// 2. 获取az信息
		var regionAzList []mysql.TeledbGetAvailabilityZoneResponseReturnObjData
		regionAzList, err = c.getAzInfoByRegion(ctx, config)

		if len(regionAzList) < 1 {
			err = errors.New("该资源池AZ信息获取为空，无法直接分配节点AZ信息")
		}
		// 定义一个az信息遍历下标
		idx := 0

		// 3. 生成master结点
		var masterAzInfo mysql.AvailabilityZoneInfoRequest
		masterAzInfo.AvailabilityZoneCount = 1
		masterAzInfo.AvailabilityZoneName = regionAzList[idx].AvailabilityZoneName
		masterAzInfo.NodeType = "master"
		*availabilityZoneInfos = append(*availabilityZoneInfos, masterAzInfo)
		nodeNum = nodeNum - 1

		// 4. 判断实例类型是否为1主1备或，1主2备。若是，则继续生成
		// 若AzNum为1 或 2，count = nodeNum -1
		// 若AzNum为3， 分两个azInfo存储
		if len(regionAzList) > 1 {
			idx = idx + 1
		}
		if nodeNum >= 1 {
			var slaveAzInfo mysql.AvailabilityZoneInfoRequest
			slaveAzInfo.AvailabilityZoneName = regionAzList[idx].AvailabilityZoneName
			slaveAzInfo.AvailabilityZoneCount = 1
			slaveAzInfo.NodeType = "slave"
			nodeNum = nodeNum - 1
			if nodeNum >= 1 {
				if len(regionAzList) >= 3 {
					*availabilityZoneInfos = append(*availabilityZoneInfos, slaveAzInfo)
					idx = idx + 1
					slaveAzInfo.AvailabilityZoneName = regionAzList[idx].AvailabilityZoneName
					slaveAzInfo.AvailabilityZoneCount = 1
					slaveAzInfo.NodeType = "slave"
				} else {
					slaveAzInfo.AvailabilityZoneCount = 2
				}
				*availabilityZoneInfos = append(*availabilityZoneInfos, slaveAzInfo)

			} else {
				*availabilityZoneInfos = append(*availabilityZoneInfos, slaveAzInfo)
			}
		}
	} else {
		var availabilityZoneInfoList []AvailabilityZoneModel
		diag := config.AvailabilityZoneInfo.ElementsAs(ctx, &availabilityZoneInfoList, true)
		if diag.HasError() {
			return
		}
		for _, availabilityZoneInfoItem := range availabilityZoneInfoList {
			availabilityZoneInfo := mysql.AvailabilityZoneInfoRequest{}
			availabilityZoneInfo.AvailabilityZoneName = availabilityZoneInfoItem.AvailabilityZoneName.ValueString()
			availabilityZoneInfo.AvailabilityZoneCount = availabilityZoneInfoItem.AvailabilityZoneCount.ValueInt32()
			availabilityZoneInfo.NodeType = availabilityZoneInfoItem.NodeType.ValueString()
			*availabilityZoneInfos = append(*availabilityZoneInfos, availabilityZoneInfo)
		}
	}

	return
}

func (c *CtyunMysqlInstance) getUpgradeAzInfo(ctx context.Context, state *CtyunMysqlInstanceConfig, plan *CtyunMysqlInstanceConfig, upgradeParams *mysql.TeledbUpgradeRequest, azInfoList *[]mysql.AvailabilityZoneInfo) (err error) {
	// 1.获取控制台上该实例目前AZ分布
	// 获取实例详情
	if state.InstID.ValueString() == "" {
		err = errors.New("查询实例详情时，实例 ID为空")
		return err
	}
	detailParams := &mysql.TeledbQueryDetailRequest{
		OuterProdInstId: state.InstID.ValueString(),
	}
	detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
		InstID:   state.InstID.ValueString(),
		RegionID: state.RegionID.ValueString(),
	}
	if !state.ProjectID.IsNull() {
		detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
	if err != nil {
		return err
	} else if resp == nil {
		err = errors.New("对Mysql实例扩容时，查询实例详情返回为nil，扩容失败，请稍后重试")
		return
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	nodeDist := resp.ReturnObj.AzInfoList
	// 2.判断是规格扩容还是节点扩容
	if upgradeParams.ProdPerformanceSpec != nil {
		// 3.规格扩容直接输入实例AZ分布情况
		err = c.getNodeDist(ctx, azInfoList, nodeDist)
		return
	} else if upgradeParams.ProdId != nil {
		// 4.节点扩容需要获取az信息，确定需要增加的节点数。
		stateNodeNum := business.MysqlNodeNumDict[state.ProdID.ValueString()]
		planNodeNum := business.MysqlNodeNumDict[plan.ProdID.ValueString()]
		addNodeNum := planNodeNum - stateNodeNum
		if addNodeNum <= 0 {
			// 如果需要增加的节点数小于等于0，无需操作
			return
		}
		err = c.getAddNodeDist(ctx, azInfoList, nodeDist, state, int(addNodeNum))
	}
	return
}

func (c *CtyunMysqlInstance) getNodeDist(_ context.Context, azInfoList *[]mysql.AvailabilityZoneInfo, nodeDist []mysql.AzInfo) (err error) {
	nodeMap := make(map[string]int32)
	// 计算节点类型，az内节点出现频率
	// 例： 1主2备
	// az1 = 3
	// az1 = 1 或者 az2 = 2
	for _, nodeDistItem := range nodeDist {
		var azInfo mysql.AvailabilityZoneInfo
		azInfo.AvailabilityZoneName = nodeDistItem.AzId
		key := nodeDistItem.AzId
		if _, exists := nodeMap[key]; exists {
			nodeMap[key] = nodeMap[key] + 1
		} else {
			nodeMap[key] = 1
		}
	}
	// 将map转成 mysql扩容的参数
	for key, value := range nodeMap {
		var azInfo mysql.AvailabilityZoneInfo
		azInfo.AvailabilityZoneCount = value
		azInfo.AvailabilityZoneName = key
		*azInfoList = append(*azInfoList, azInfo)
	}
	return
}

func (c *CtyunMysqlInstance) getAddNodeDist(ctx context.Context, azInfoList *[]mysql.AvailabilityZoneInfo, nodeDist []mysql.AzInfo, state *CtyunMysqlInstanceConfig, addNodeNum int) (err error) {
	// 定义map,存放az-节点数分布
	nodeMap := make(map[string]int32)
	addNodeMap := make(map[string]int32)
	defaultAzId := ""
	// 获取该资源池az列表
	regionAzList, err := c.getAzInfoByRegion(ctx, state)
	if err != nil {
		return err
	}
	// 对map进行初始化
	for _, AzInfo := range regionAzList {
		if defaultAzId == "" {
			defaultAzId = AzInfo.AvailabilityZoneId
		}
		if _, exist := nodeMap[AzInfo.AvailabilityZoneId]; !exist {
			nodeMap[AzInfo.AvailabilityZoneId] = 0
		}
	}
	// 统计每个az的节点数
	for _, nodeDistItem := range nodeDist {
		//var azInfo mysql.AvailabilityZoneInfo
		//azInfo.AvailabilityZoneName = nodeDistItem.AzId
		key := nodeDistItem.AzId
		if _, exists := nodeMap[key]; exists {
			nodeMap[key] = nodeMap[key] + 1
		} else {
			nodeMap[key] = 1
		}
	}
	// 根据需要增加的节点数，每次选取，最小的value值的az
	for i := 0; i < addNodeNum; i++ {
		minNodeNum := int32(math.MaxInt32)
		minAzName := ""
		for key, value := range nodeMap {
			if value < minNodeNum {
				minNodeNum = value
				minAzName = key
			}
		}
		if minAzName == "" {
			minAzName = defaultAzId
		}
		if _, exists := addNodeMap[minAzName]; exists {
			addNodeMap[minAzName] = addNodeMap[minAzName] + 1
		} else {
			addNodeMap[minAzName] = 1
		}
		nodeMap[minAzName] = nodeMap[minAzName] + 1
	}
	// 将map转成 mysql扩容的参数
	for key, value := range addNodeMap {
		var azInfo mysql.AvailabilityZoneInfo
		azInfo.AvailabilityZoneCount = value
		azInfo.AvailabilityZoneName = key
		*azInfoList = append(*azInfoList, azInfo)
	}
	return
}

func (c *CtyunMysqlInstance) getAzInfoByRegion(ctx context.Context, config *CtyunMysqlInstanceConfig) (regionAzList []mysql.TeledbGetAvailabilityZoneResponseReturnObjData, err error) {
	params := &mysql.TeledbGetAvailabilityZoneRequest{
		RegionId: config.RegionID.ValueString(),
	}
	header := &mysql.TeledbGetAvailabilityZoneRequestHeader{}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbGetAvailabilityZone.Do(ctx, c.meta.Credential, params, header)
	if err2 != nil {
		err = err2
		return
	} else if resp == nil {
		err = errors.New("查询该资源池AZ信息时，返回为nil。请稍后再试")
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj.Data == nil {
		err = common.InvalidReturnObjError
		return
	}
	regionAzList = resp.ReturnObj.Data
	if regionAzList == nil || len(regionAzList) == 0 {
		err = errors.New("查询该资源池AZ信息时，返回为空。请稍后再试")
		return
	}

	return
}

func (c *CtyunMysqlInstance) UpgradeStorageLoop(ctx context.Context, state *CtyunMysqlInstanceConfig, plan *CtyunMysqlInstanceConfig, NodeType string) (err error) {
	count := 60
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			// 获取实例详情
			detailParams := &mysql.TeledbQueryDetailRequest{
				OuterProdInstId: state.InstID.ValueString(),
			}
			detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
				InstID:   state.InstID.ValueString(),
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			runningStatus := resp.ReturnObj.ProdRunningStatus
			orderStatus := resp.ReturnObj.ProdOrderStatus
			// 若符合预期，跳出循环，扩容成功
			if runningStatus == business.MysqlRunningStatusStarted && orderStatus == business.MysqlOrderStatusStarted {
				if NodeType == business.PgsqlStorageTypeMaster {
					if resp.ReturnObj.DiskSize == plan.StorageSpace.ValueInt32() {
						return false
					}
				} else if NodeType == business.PgsqlStorageTypeBackUp {
					if plan.BackupStorageSpace.ValueInt32() != 0 && plan.BackupStorageSpace.ValueInt32() == resp.ReturnObj.BackupDiskSize {
						return false
					}
				}
			}
			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未升级成功！")
	}
	return
}

func (c *CtyunMysqlInstance) upgradeMysqlStorage(ctx context.Context, state *CtyunMysqlInstanceConfig, plan *CtyunMysqlInstanceConfig, upgradeParams *mysql.TeledbUpgradeRequest, upgradeHeader *mysql.TeledbUpgradeRequestHeader) (err error) {
	err = c.StartedLoop(ctx, state)
	if err != nil {
		return
	}
	resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbUpgradeApi.Do(ctx, c.meta.Credential, upgradeParams, upgradeHeader)
	if err2 != nil {
		err = err2
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("update storage failed, API return error. Message: %s Error: %s", resp.Message, resp.Error)
		return
	}
	// 扩容后，轮循请求实例详情，确认已经完成升配
	err = c.UpgradeStorageLoop(ctx, state, plan, *upgradeParams.NodeType)
	return
}

// checkSpec 检查规格
func (c *CtyunMysqlInstance) checkSpec(ctx context.Context, plan *CtyunMysqlInstanceConfig) error {
	// 先根据spec_name调用云主机规格接口
	_, err := c.ecsService.GetFlavorByName(ctx, plan.FlavorName.ValueString(), plan.RegionID.ValueString())
	if err != nil {
		return err
	}

	f := strings.Split(plan.FlavorName.ValueString(), ".")
	hostType := strings.ToUpper(f[0])
	plan.instanceSeries = string(hostType[0]) // S、M 或 C
	if len(hostType) > 2 {
		plan.instanceSeries = hostType
	}
	// 再调用数据库规格接口
	mysqlFlavor, err := c.mysqlService.GetFlavorByProdIdAndFlavorName(
		ctx,
		plan.ProdID.ValueString(),
		plan.FlavorName.ValueString(),
		plan.RegionID.ValueString(),
		plan.instanceSeries,
	)
	if err != nil {
		return err
	}
	plan.prodPerformanceSpec = mysqlFlavor.ProdPerformanceSpec
	plan.hostType = mysqlFlavor.Generation

	// 映射关系
	if strings.HasPrefix(plan.hostType, "K") { // 鲲鹏
		plan.cpuType = "KunPeng"
	} else if strings.HasPrefix(plan.hostType, "H") { // 海光
		plan.cpuType = "Hygon"
	} else if strings.HasPrefix(plan.hostType, "F") {
		plan.cpuType = "Phytium"
	} else {
		plan.cpuType = "Intel"
	}
	plan.osType = "ctyunos"
	return nil
}

type CtyunMysqlInstanceConfig struct {
	CycleType                   types.String `tfsdk:"cycle_type"`                     // 计费模式： 支持on_demand和month
	RegionID                    types.String `tfsdk:"region_id"`                      // 资源池Id
	VpcID                       types.String `tfsdk:"vpc_id"`                         // 虚拟私有云Id
	FlavorName                  types.String `tfsdk:"flavor_name"`                    // 规格名称
	SubnetID                    types.String `tfsdk:"subnet_id"`                      // 子网Id
	SecurityGroupID             types.String `tfsdk:"security_group_id"`              // 安全组
	Name                        types.String `tfsdk:"name"`                           // 集群名称
	Password                    types.String `tfsdk:"password"`                       // 管理员密码（RSA公钥加密）
	CycleCount                  types.Int32  `tfsdk:"cycle_count"`                    // 购买时长：单位月（范围：1-12，24，36）
	AutoRenew                   types.Bool   `tfsdk:"auto_renew"`                     // 自动续订状态
	ProdID                      types.String `tfsdk:"prod_id"`                        // 产品id
	MasterOrderID               types.String `tfsdk:"master_order_id"`                // 订单id
	InstID                      types.String `tfsdk:"inst_id"`                        // 实例id
	ProjectID                   types.String `tfsdk:"project_id"`                     // 项目id
	ProdRunningStatus           types.Int32  `tfsdk:"prod_running_status"`            // 以查询实例列表为主，0.正常 1.重启中 2.备份中 3.恢复中 4.修改参数中 5.应用参数组中 6&7.扩容中 8.修改端口中 9.迁移中 10.重置密码中
	Vip                         types.String `tfsdk:"vip"`                            // 虚拟IP地址
	WritePort                   types.Int32  `tfsdk:"write_port"`                     // 写数据端口
	ReadPort                    types.String `tfsdk:"read_port"`                      // 读端口
	ProdDbEngine                types.String `tfsdk:"prod_db_engine"`                 // 数据库引擎
	EIP                         types.String `tfsdk:"eip"`                            // 弹性ip
	EipStatus                   types.Int32  `tfsdk:"eip_status"`                     // 弹性ip状态 0->unbind，1->bind,2->binding
	SSlStatus                   types.Int32  `tfsdk:"ssl_status"`                     // Ssl状态 0->off，1->on
	NewMysqlVersion             types.String `tfsdk:"new_mysql_version"`              // mysql版本
	AuditLogStatus              types.Int32  `tfsdk:"audit_log_status"`               // 日志审计开关
	InstReleaseProtectionStatus types.Int32  `tfsdk:"inst_release_protection_status"` // 实例释放保护开关 1:on,0:off
	PauseEnable                 types.Bool   `tfsdk:"pause_enable"`                   // 是否允许暂停
	MysqlPort                   types.String `tfsdk:"mysql_port"`                     // 数据库端口
	SecurityGroupStatus         types.Int32  `tfsdk:"security_group_status"`          // 安全组状态 0->normal, 1->changing, 2->deleted
	StorageType                 types.String `tfsdk:"storage_type"`                   // 存储类型：SSD, SATA, SAS, SSD-genric, FAST-SSD
	StorageSpace                types.Int32  `tfsdk:"storage_space"`                  // 存储空间（单位：GB，范围100到32768）
	BackupStorageType           types.String `tfsdk:"backup_storage_type"`            // 备份存储空间磁盘类型
	BackupStorageSpace          types.Int32  `tfsdk:"backup_storage_space"`           // 备份节点，存储空间扩容使用
	AvailabilityZoneInfo        types.List   `tfsdk:"availability_zone_info"`         // 可用区信息
	RunningControl              types.String `tfsdk:"running_control"`                //
	ProdOrderStatus             types.Int32  `tfsdk:"prod_order_status"`
	ID                          types.String `tfsdk:"id"` // 实例id

	osType              string
	cpuType             string
	prodPerformanceSpec string
	hostType            string
	instanceSeries      string
}

type AvailabilityZoneModel struct {
	AvailabilityZoneName  types.String `tfsdk:"availability_zone_name"`  // 资源池可用区名称
	AvailabilityZoneCount types.Int32  `tfsdk:"availability_zone_count"` // 资源池可用区总数
	NodeType              types.String `tfsdk:"node_type"`               // 表示分布AZ的节点类型，master/slave
}

type UpdatedAZModel struct {
	AvailabilityZoneName  types.String `tfsdk:"availability_zone_name"`  // 资源池可用区名称
	AvailabilityZoneCount types.Int32  `tfsdk:"availability_zone_count"` // 资源池可用区总数
}
