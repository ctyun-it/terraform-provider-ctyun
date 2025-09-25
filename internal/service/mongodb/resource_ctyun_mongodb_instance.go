package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
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
	_ resource.Resource                = &CtyunMongodbInstance{}
	_ resource.ResourceWithConfigure   = &CtyunMongodbInstance{}
	_ resource.ResourceWithImportState = &CtyunMongodbInstance{}
)

type CtyunMongodbInstance struct {
	meta           *common.CtyunMetadata
	ecsService     *business.EcsService
	mongodbService *business.MongodbService
}

func (c *CtyunMongodbInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {

}

func (c *CtyunMongodbInstance) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ecsService = business.NewEcsService(c.meta)
	c.mongodbService = business.NewMongodbService(c.meta)
}

func NewCtyunMongodbInstance() resource.Resource {
	return &CtyunMongodbInstance{}
}

func (c *CtyunMongodbInstance) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mongodb_instance"
}

func (c *CtyunMongodbInstance) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10034467/10089535`,
		Attributes: map[string]schema.Attribute{
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
				Description: "区域id,如果不填这默认使用provider ctyun总region_id 或者环境变量",
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
			"flavor_name": schema.StringAttribute{
				Required:    true,
				Description: "规格名称，形如c7.2xlarge.4，可从data.ctyun_mongodb_specs查询支持的规格。支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
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
				Description: "实例名称（长度在 4 到 64个字符，必须以字母开头，不区分大小写，可以包含字母、数字、中划线或下划线，不能包含其他特殊字符），支持更新。",
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 64),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]*$"), "实例名称不符合规范"),
				},
			},
			// 实现一个validator方法
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "实例密码，长度为8~26个字符，必须包含大写字母、小写字母、数字和特殊字符~!@#%^*_=+",
				Validators: []validator.String{
					validator2.DBPassword(
						8,
						26,
						4,
						"MongoDB",
						"~!@#%^*_=+",
					),
				},
			},
			"prod_id": schema.StringAttribute{
				Required:    true,
				Description: "产品id，开通时用于确定开通单机/集群版/副本集和版本，支持更新。取值范围包括：Single34（3.4单机版）,Single40（4.0单机版）,Replica3R34（3.4副本集三副本）,Replica3R40（4.0副本集三副本）,Replica5R34（3.4副本集五副本）,Replica5R40（4.0副本集五副本）,Replica7R34（3.4副本集七副本）,Replica7R40（4.0副本集七副本）,Cluster34（3.4集群版）,Cluster40（4.0集群版）,Single42（4.2单机版）,Replica3R42（4.2副本集三副本）,Replica5R42（4.2副本集五副本）,Replica7R42（4.2副本集七副本）,Cluster42（4.2集群版）,Single50（5.0单机版）,Replica3R50（5.0副本集三副本）,Replica5R50（5.0副本集五副本）,Replica7R50（5.0副本集七副本）,Cluster50（5.0集群版）,Cluster60（6.0集群版）,Replica3R60（6.0副本集三副本）,Replica5R60（6.0副本集五副本）,Replica7R60（6.0副本集七副本）,Single60（6.0单机版）",
				Validators: []validator.String{
					stringvalidator.OneOf(business.MongodbProdIDs...),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订单id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"read_port": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "读端口,支持更新。若需要更新读取端口时可填，取值范围：1~65535",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"host_ip": schema.StringAttribute{
				Computed:    true,
				Description: "主机ip",
			},
			"innodb_buffer_pool_size": schema.StringAttribute{
				Computed:    true,
				Description: "缓存池大小",
			},
			"innodb_thread_concurrency": schema.Int64Attribute{
				Computed:    true,
				Description: "线程数",
			},
			"prod_running_status": schema.Int32Attribute{
				Computed:    true,
				Description: "实例运行状态: 0->运行正常, 1->重启中, 2-备份操作中, 3->恢复操作中,4->转换ssl,5->异常,6->修改参数组中,7->已冻结,8->已注销,9->施工中,10->施工失败,11->扩容中,12->主备切换中",
			},
			"prod_running_status_desc": schema.StringAttribute{
				Computed:    true,
				Description: "实例运行状态解释字段",
			},
			"eip_id": schema.StringAttribute{
				Computed:    true,
				Description: "eip Id",
			},
			"is_upgrade_back_up": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "磁盘扩容时候会使用,是否主磁盘与备磁盘一起扩容，支持更新。该参数仅在升配主存储空间时生效，且需要注意is_upgrade_back_up=ture时，待升配的磁盘空间必须大于现磁盘空间（包括备份空间）。取值范围：true-主备同时扩容； false-主备不同时扩容。默认为false",
			},
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "mongodb实例id",
			},
			"storage_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("SSD"),
				Description: "存储类型，默认为SSD。取值范围：SSD=超高IO, SAS=高IO, SATA=普通IO，SSD-genric=通用型SSD",
				Validators: []validator.String{
					stringvalidator.OneOf(business.MongodbStorageType...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_space": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(100),
				Description: "存储空间(单位:G)，默认为100GB，支持更新。取值范围：10-6144，backup节点为单个shard的容量乘以shard的个数",
				Validators: []validator.Int32{
					int32validator.Between(10, 6144),
				},
			},
			"availability_zone_info": schema.ListNestedAttribute{
				Optional:    true,
				Description: "mongodb实例节点指定可用区字段，选填。若不填写，将按节点个数均匀分布到各个可用区上。若需要填写可参考提供的examples",
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
							Description: "资源池可用区总数（开通集群版--nodeType为mongos时范围为[2,16]，nodeType为shard时,shard数量取值范围[2,16]，每一个shard对应3个availabilityZoneCount, 例：nodeType: shard且要开通shard数 量为3时，availabilityZoneCount:9 ；nodeType为config时节点默认为3即availabilityZoneCount: 3）",
							Validators: []validator.Int32{
								int32validator.Between(1, 16),
							},
						},
						"node_type": schema.StringAttribute{
							Required:    true,
							Description: "master:主节点、mongos:mongos节点、shard:shard节点 、config:config节点（存储类型storageType与shard节点一致，存储空间storageSpace为单个shard的storageSpace）、 backup:备份机(存储类型storageType与shard 节点一致，存储空间storageSpace为shard节点数量乘以单个shard的storageSpace)",
							Validators: []validator.String{
								stringvalidator.OneOf(business.MongodbNodeType...),
							},
						},
					},
				},
			},
			"shard_num": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "shard节点数量，mongodb为集群版需填写，支持更新。默认为2，取值范围：2~32",
				Default:     int32default.StaticInt32(2),
				Validators: []validator.Int32{
					int32validator.Between(2, 32),
				},
			},
			"mongos_num": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "mongos节点数量，mongodb为集群版需填写，支持更新。默认为2，取值范围：2~32",
				Default:     int32default.StaticInt32(2),
				Validators: []validator.Int32{
					int32validator.Between(2, 32),
				},
			},
			"backup_storage_space": schema.Int32Attribute{
				Computed:    true,
				Description: "backup节点磁盘空间，当前不支持指定。默认与存储空间相同",
			},
			"backup_storage_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "backup节点存储类型，取值范围：SATA, SAS, SSD, OS（对象存储）。若不填写，默认为云硬盘（SSD）",
				Validators: []validator.String{
					stringvalidator.OneOf(business.StorageTypeSATA, business.StorageTypeSAS, business.StorageTypeSSD, business.BackupStorageTypeOS),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"upgrade_node_type": schema.StringAttribute{
				Optional:    true,
				Description: "当实例为集群版，若升配mongos、shard节点个数时可填写，支持更新。取值范围：shard, mongos",
				Validators: []validator.String{
					stringvalidator.OneOf("shard", "mongos"),
				},
			},
		},
	}
}

func (c *CtyunMongodbInstance) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMongodbInstanceConfig
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
	err = c.CreateMongodbInstance(ctx, &plan)
	if err != nil {
		return
	}
	// 如果ReadPort不为空，需要暂存下来，否则merge操作后会被默认端口覆盖
	var updatedReadPort int32
	var needUpatePort bool
	if !plan.ReadPort.IsNull() && !plan.ReadPort.IsUnknown() {
		updatedReadPort = plan.ReadPort.ValueInt32()
		needUpatePort = true
	}
	// 创建完成后，同步云端信息
	err = c.getAndMergeMongodbInstance(ctx, &plan)
	if err != nil {
		return
	}
	// 确保实例创建成功后，判断port是否需要指定
	if needUpatePort {
		plan.ReadPort = types.Int32Value(updatedReadPort)
		err = c.updateReadPort(ctx, &plan, &plan)
		if err != nil {
			return
		}
	}

	err = c.getAndMergeMongodbInstance(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMongodbInstance) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunMongodbInstanceConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergeMongodbInstance(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMongodbInstance) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置
	var plan CtyunMongodbInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunMongodbInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	if !plan.Password.Equal(state.Password) {
		err = fmt.Errorf("数据库密码暂时不支持修改")
		return
	}

	// 通过flavor_name获取cpu，memory等规格信息
	err = c.checkSpec(ctx, &plan)
	if err != nil {
		return
	}
	err = c.updateMongodbInstance(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeMongodbInstance(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMongodbInstance) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunMongodbInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	deleteParams := &mongodb.MongodbRefundRequest{
		InstId: state.ID.ValueString(),
	}
	deleteHeader := &mongodb.MongodbRefundRequestHeader{}
	if state.ProjectID.ValueString() != "" {
		deleteHeader.ProjectID = state.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkMongodbApis.MongodbRefundApi.Do(ctx, c.meta.Credential, deleteParams, deleteHeader)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	// 轮询确认时候退订成功
	err = c.DeleteLoop(ctx, &state, 60)
	if err != nil {
		return
	}
	//time.Sleep(30 * time.Second)
	//err = c.destroy(ctx, state)
	//if err != nil {
	//	return
	//}
	//err = c.destroyLoop(ctx, state)
	response.Diagnostics.AddWarning("删除MongoDB集群成功", "集群退订后，若立即删除子网或安全组可能会失败，需要等待底层资源释放")
}

func (c *CtyunMongodbInstance) CreateMongodbInstance(ctx context.Context, config *CtyunMongodbInstanceConfig) (err error) {

	cycleType := config.CycleType.ValueString()
	params := &mongodb.MongodbCreateRequest{
		BillMode:          business.MysqlBillMode[cycleType],
		RegionId:          config.RegionID.ValueString(),
		VpcId:             config.VpcID.ValueString(),
		HostType:          config.hostType,
		SubnetId:          config.SubnetID.ValueString(),
		SecurityGroupId:   config.SecurityGroupID.ValueString(),
		Name:              config.Name.ValueString(),
		Period:            config.CycleCount.ValueInt32(),
		Count:             1,
		ProdId:            business.MongodbProdIDDict[config.ProdID.ValueString()],
		MysqlNodeInfoList: nil,
	}
	if config.BackupStorageType.ValueString() == business.BackupStorageTypeOS {
		osStr := strings.ToLower(config.BackupStorageType.ValueString())
		params.BackupStorageType = &osStr
	}
	password := business.Encode(config.Password.ValueString())
	params.Password = password
	//config.Password = types.StringValue(password)
	if cycleType == business.OnDemandCycleType {
		params.AutoRenewStatus = 0
	} else {
		params.AutoRenewStatus = map[bool]int32{true: 1, false: 0}[config.AutoRenew.ValueBool()]
	}

	var mongodbNodeInfoListRequest []mongodb.MongodbNodeInfoListRequest
	// 获取az信息
	if strings.Contains(config.ProdID.ValueString(), "Single") {
		// 处理单节点nodeInfoList
		err2 := c.getSingleNodeInfo(ctx, config, &mongodbNodeInfoListRequest)
		if err2 != nil {
			err = err2
			return
		}
	} else if strings.Contains(config.ProdID.ValueString(), "Replica") {
		// 处理副本级nodeInfoList
		err2 := c.getReplicaNodeInfo(ctx, config, &mongodbNodeInfoListRequest)
		if err2 != nil {
			err = err2
			return
		}
	} else if strings.Contains(config.ProdID.ValueString(), "Cluster") {
		// 处理集群版本nodeInfoList
		err2 := c.getClusterNodeInfo(ctx, config, &mongodbNodeInfoListRequest)
		if err2 != nil {
			err = err2
			return
		}
	}
	params.MysqlNodeInfoList = mongodbNodeInfoListRequest

	header := &mongodb.MongodbCreateRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkMongodbApis.MongodbCreateApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}
	if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 保存newOrderId
	config.MasterOrderID = utils.SecStringValue(resp.ReturnObj.Data.NewOrderId)
	return
}

func (c *CtyunMongodbInstance) getAndMergeMongodbInstance(ctx context.Context, config *CtyunMongodbInstanceConfig) (err error) {

	listParams := &mongodb.MongodbGetListRequest{
		PageNow:      1,
		PageSize:     100,
		ProdInstName: config.Name.ValueStringPointer(),
	}
	listHeader := &mongodb.MongodbGetListHeaders{
		RegionID: config.RegionID.ValueString(),
	}
	if config.ProjectID.ValueString() != "" {
		listHeader.ProjectID = config.ProjectID.ValueStringPointer()
	}
	// 若实例id为空，实例刚刚创建，还未查询到id，需要轮询列表获取实例信息
	if config.ID.ValueString() == "" {
		resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, listParams, listHeader)
		if err2 != nil {
			err = err2
			return
		}
		if len(resp.ReturnObj.List) > 1 {
			err = errors.New("实例名重复！")
			return
		} else if len(resp.ReturnObj.List) < 1 {
			//若根据name查询不到机器，可能存在还未创建好的情况，需要轮询
			resp, err = c.ListLoop(ctx, listParams, listHeader, 60)
			if err != nil {
				return
			} else if resp == nil {
				err = errors.New("获取mongodb列表信息，返回Nil")
				return
			}

			if len(resp.ReturnObj.List) != 1 {
				err = errors.New("未查询该实例mysql，mysql name:" + config.Name.ValueString())
				return
			}
			// 查询到实例后，保存id
			if resp.ReturnObj.List[0].ProdInstId == "" {
				err = errors.New("实例创建后，实例id仍为空")
				return
			}
			config.ID = types.StringValue(resp.ReturnObj.List[0].ProdInstId)
		}
	}
	// 确认实例id不为空后,分两步查询实例详情
	// 1）轮询实例状态，确认已经正常运行，并获取实例部门详情：读端口、缓冲池信息和安全组信息
	listResp, err := c.RunningLoop(ctx, listParams, listHeader, 60)
	if err != nil {
		return err
	} else if listResp == nil {
		err = fmt.Errorf("列表查询返回为nil")
		return
	} else if listResp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s", *listResp.Message)
		return
	} else if listResp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	listReturnObj := listResp.ReturnObj.List[0]

	if config.ID.ValueString() == "" {
		err = errors.New("查询实例详情时，实例id为空")
	}
	// 2）查询实例详情，获取allowBeMaster信息和eip id信息

	detailReturnObj, err := c.getMongoDetailInfo(ctx, config)
	if err != nil {
		return
	}

	port, err := strconv.ParseInt(detailReturnObj.Port, 10, 32)
	if err != nil {
		return
	}
	config.ReadPort = types.Int32Value(int32(port))
	config.InnodbBufferPoolSize = types.StringValue(listReturnObj.InnodbBufferPoolSize)
	config.InnodbThreadConcurrency = types.Int64Value(listReturnObj.InnodbThreadConcurrency)
	config.ProdRunningStatus = types.Int32Value(listReturnObj.ProdRunningStatus)
	config.ProdRunningStatusDesc = types.StringValue(business.MongodbStatusDescDict[listReturnObj.ProdRunningStatus])
	config.EipID = types.StringValue(detailReturnObj.NodeInfoVOS[0].OuterElasticIpId)
	config.Name = types.StringValue(listReturnObj.ProdInstName)
	config.SecurityGroupID = types.StringValue(listReturnObj.SecurityGroupId)
	//prodID, err := strconv.ParseInt(listReturnObj.ProdId, 10, 64)
	//if err != nil {
	//	return
	//}
	//config.ProdID = types.StringValue(business.MongodbProdIDRevDict[prodID])
	config.HostIp = types.StringValue(detailReturnObj.Host)
	config.StorageSpace = types.Int32Value(detailReturnObj.DiskSize)
	if detailReturnObj.Backup != nil {
		backupSize := strings.TrimSuffix(detailReturnObj.Backup.Size, "G")
		backupStorageSpace, err2 := strconv.ParseInt(backupSize, 10, 32)
		if err2 != nil {
			err = err2
			return
		}
		config.BackupStorageSpace = types.Int32Value(int32(backupStorageSpace))

		//if strings.Contains(config.ProdID.ValueString(), "Cluster") {
		//	backupStorageSpaceAvg := int32(backupStorageSpace) / config.ShardNum.ValueInt32()
		//	config.BackupStorageSpace = types.Int32Value(backupStorageSpaceAvg)
		//} else {
		//	config.BackupStorageSpace = types.Int32Value(int32(backupStorageSpace))
		//}
	} else {
		config.BackupStorageSpace = types.Int32Value(0)
	}
	return
}

func (c *CtyunMongodbInstance) updateMongodbInstance(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig) (err error) {
	if state.ID.ValueString() == "" {
		err = errors.New("在变配实例过程中， 实例id为空")
		return
	}
	// prod_id（节点） 和 flavor不可同时升级
	if !plan.FlavorName.Equal(state.FlavorName) && !plan.ProdID.Equal(state.ProdID) {
		err = fmt.Errorf("mongodb flavor_name（cpu 内存规格） 和 prod id（节点） 不可同时更新")
		return err
	}
	// 修改实例名称
	if plan.Name.ValueString() != "" && state.Name.ValueString() != plan.Name.ValueString() {
		// 修改实例前，确定实例状态为running
		_, err = c.PreCheckUpdateLoop(ctx, state)
		if err != nil {
			return
		}
		updateNameParams := &mongodb.MongodbUpdateInstanceNameRequest{
			ProdInstId:   state.ID.ValueString(),
			ProdInstName: plan.Name.ValueString(),
		}
		updateNameHeader := &mongodb.MongodbUpdateInstanceNameRequestHeader{
			RegionID: state.RegionID.ValueString(),
		}
		resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbUpdateInstanceNameApi.Do(ctx, c.meta.Credential, updateNameParams, updateNameHeader)
		if err2 != nil {
			err = err2
			return
		} else if resp.StatusCode != 800 {
			err = fmt.Errorf("API return error. Message: %s", *resp.Message)
			return
		}
	}

	// 变更完name，确认是否变更完成
	err = c.PostCheckUpdateNameAndSecurityGroupLoop(ctx, state, plan, 60)
	if err != nil {
		return
	}
	// 更新name，因为read-merge阶段，需要查询列表，查询列表通过name，如果name更新了，查询不到导致异常
	state.Name = types.StringValue(plan.Name.ValueString())
	// 修改实例端口
	if plan.ReadPort.ValueInt32() != 0 && state.ReadPort.ValueInt32() != plan.ReadPort.ValueInt32() {
		// 修改实例前，确定实例状态为running
		err = c.updateReadPort(ctx, state, plan)
		if err != nil {
			return
		}
	}

	// 实例扩容
	// 扩容磁盘
	err = c.upgradeStorage(ctx, state, plan)
	if err != nil {
		return
	}

	// 扩容规格
	err = c.upgradeSpec(ctx, state, plan)
	if err != nil {
		return
	}

	// 扩容节点
	// 单集群不支持扩容，cluster集群支持扩容
	// 集群版支持扩容shard数量和mongos数量
	err = c.upgradeNode(ctx, state, plan)
	if err != nil {
		return
	}
	state.AvailabilityZoneInfo = plan.AvailabilityZoneInfo
	state.FlavorName = plan.FlavorName
	if !plan.UpgradeNodeType.IsNull() {
		state.UpgradeNodeType = plan.UpgradeNodeType
	}
	return
}

func (c *CtyunMongodbInstance) PreCheckUpdateLoop(ctx context.Context, state *CtyunMongodbInstanceConfig, loopCount ...int) (ListResp *mongodb.MongodbGetListResponse, err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	syncCount := 3
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return nil, err
	}
	listParams := &mongodb.MongodbGetListRequest{
		PageNow:      1,
		PageSize:     100,
		ProdInstName: state.Name.ValueStringPointer(),
	}
	listHeader := &mongodb.MongodbGetListHeaders{
		RegionID: state.RegionID.ValueString(),
	}
	if state.ProjectID.ValueString() != "" {
		listHeader.ProjectID = state.ProjectID.ValueStringPointer()
	}

	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, listParams, listHeader)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			if len(resp.ReturnObj.List) != 1 {
				err = errors.New("实例name不唯一，有误！")
				return false
			}
			if resp.ReturnObj.List[0].ProdRunningStatus == business.MongodbRunningStatusStarted && resp.ReturnObj.List[0].ProdOrderStatus == business.MongodbOrderStatusStarted {
				if syncCount > 0 {
					syncCount--
					return true
				}
				ListResp = resp
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return nil, errors.New("轮询已达最大次数，实例仍未运行成功！")
	}
	return
}

func (c *CtyunMongodbInstance) ListLoop(ctx context.Context, params *mongodb.MongodbGetListRequest, header *mongodb.MongodbGetListHeaders, loopCount ...int) (*mongodb.MongodbGetListResponse, error) {
	var err error
	var response *mongodb.MongodbGetListResponse
	count := 120
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return nil, err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, params, header)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
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
		return nil, errors.New("轮询已达最大次数，实例仍未创建或查询到！")
	}
	return response, nil
}

func (c *CtyunMongodbInstance) RunningLoop(ctx context.Context, params *mongodb.MongodbGetListRequest, header *mongodb.MongodbGetListHeaders, loopCount ...int) (*mongodb.MongodbGetListResponse, error) {
	var err error
	var response *mongodb.MongodbGetListResponse
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
			resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, params, header)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			if len(resp.ReturnObj.List) <= 0 {
				err = common.InvalidReturnObjError
				return false
			}

			runningStatus := resp.ReturnObj.List[0].ProdRunningStatus
			// 若实例状态已经运行正常，跳出轮询
			if runningStatus == business.MongodbRunningStatusStarted {
				response = resp
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return nil, errors.New("轮询已达最大次数，实例仍未启动！")
	}
	return response, nil
}

func (c *CtyunMongodbInstance) DeleteLoop(ctx context.Context, config *CtyunMongodbInstanceConfig, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	listParams := &mongodb.MongodbGetListRequest{
		PageNow:      1,
		PageSize:     100,
		ProdInstName: config.Name.ValueStringPointer(),
	}
	listHeader := &mongodb.MongodbGetListHeaders{
		RegionID: config.RegionID.ValueString(),
	}
	if config.ProjectID.ValueString() != "" {
		listHeader.ProjectID = config.ProjectID.ValueStringPointer()
	}

	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, listParams, listHeader)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			}
			if len(resp.ReturnObj.List) == 0 {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，实例仍未删除成功！")
	}
	return
}

func (c *CtyunMongodbInstance) PostCheckUpdateNameAndSecurityGroupLoop(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	listParams := &mongodb.MongodbGetListRequest{
		PageNow:      1,
		PageSize:     100,
		ProdInstName: plan.Name.ValueStringPointer(),
	}
	listHeader := &mongodb.MongodbGetListHeaders{
		RegionID: state.RegionID.ValueString(),
	}
	if state.ProjectID.ValueString() != "" {
		listHeader.ProjectID = state.ProjectID.ValueStringPointer()
	}

	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, listParams, listHeader)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			}
			if len(resp.ReturnObj.List) != 1 {
				err = errors.New("根据name查询实例数量有误！")
				return false
			}
			updatedName := resp.ReturnObj.List[0].ProdInstName
			updatedSecurityGroupID := resp.ReturnObj.List[0].SecurityGroupId
			flagName := true
			flagSecurityGroup := true
			if plan.Name.ValueString() != "" {
				if updatedName != plan.Name.ValueString() {
					flagName = false
				}
			}
			if plan.SecurityGroupID.ValueString() != "" {
				if updatedSecurityGroupID != plan.SecurityGroupID.ValueString() {
					flagSecurityGroup = false
				}
			}
			if flagName && flagSecurityGroup {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，实例仍未删除成功！")
	}
	return
}

func (c *CtyunMongodbInstance) UpdatePortLoop(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	tolerateCount := 5
	result := retryer.Start(
		func(currentTime int) bool {
			// 查询详情
			detailResp, err2 := c.getMongoDetailInfo(ctx, state)
			if err2 != nil {
				if tolerateCount <= 0 {
					err = err2
					return false
				}
				tolerateCount--
				return true
			}
			// 若云端port信息与预期相符，退出轮询
			if detailResp.Port == fmt.Sprintf("%d", plan.ReadPort.ValueInt32()) {
				return false
			}
			// 继续轮询
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，实例端口仍未更新成功！")
	}
	return
}

func (c *CtyunMongodbInstance) UpgradeLoop(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig, planNodeInfoList []NodeInfoListModel, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	tolerateCount := 30

	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			detailParams := &mongodb.MongodbQueryDetailRequest{
				ProdInstId: state.ID.ValueString(),
			}
			detailHeader := &mongodb.MongodbQueryDetailRequestHeaders{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeader.ProjectID = state.ProjectID.ValueStringPointer()
			}

			listParams := &mongodb.MongodbGetListRequest{
				PageNow:      1,
				PageSize:     100,
				ProdInstName: state.Name.ValueStringPointer(),
			}
			listHeader := &mongodb.MongodbGetListHeaders{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				listHeader.ProjectID = state.ProjectID.ValueStringPointer()
			}

			listResp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, listParams, listHeader)
			if err2 != nil {
				err = err2
				return false

			} else if listResp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *listResp.Message)
				return false
			} else if listResp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}

			detailResp, err2 := c.meta.Apis.SdkMongodbApis.MongodbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeader)
			if err2 != nil {
				if tolerateCount <= 0 {
					err2 = err
					return false
				}
				tolerateCount--
				return true
			} else if detailResp.StatusCode != 800 {
				if tolerateCount <= 0 {
					err = fmt.Errorf("API return error. Message: %s", *detailResp.Message)
					return false
				}
				tolerateCount--
				return true
			} else if detailResp.ReturnObj == nil {
				if tolerateCount <= 0 {
					err = common.InvalidReturnObjError
					return false
				}
				tolerateCount--
				return true
			}
			// 验证配置
			specFlag := true
			machineSpec := detailResp.ReturnObj.MachineSpec
			if planNodeInfoList[0].ProdPerformanceSpec.ValueString() != machineSpec {
				specFlag = false
			}
			// 验证prodID
			prodIDFlag := true
			prodID := listResp.ReturnObj.List[0].ProdId
			if plan.ProdID.ValueString() != "" && prodID != fmt.Sprintf("%d", business.MongodbProdIDDict[plan.ProdID.ValueString()]) {
				prodIDFlag = false
			}
			if specFlag && prodIDFlag {
				return false
			}
			// 继续轮询
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，实例端口仍未更新成功！")
	}
	return

}

func (c *CtyunMongodbInstance) UpgradeStorageLoop(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig, planNodeInfoList NodeInfoListModel, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	tolerateCount := 5
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			detailParams := &mongodb.MongodbQueryDetailRequest{
				ProdInstId: state.ID.ValueString(),
			}
			detailHeader := &mongodb.MongodbQueryDetailRequestHeaders{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeader.ProjectID = state.ProjectID.ValueStringPointer()
			}

			listParams := &mongodb.MongodbGetListRequest{
				PageNow:      1,
				PageSize:     100,
				ProdInstName: state.Name.ValueStringPointer(),
			}
			listHeader := &mongodb.MongodbGetListHeaders{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				listHeader.ProjectID = state.ProjectID.ValueStringPointer()
			}

			listResp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, listParams, listHeader)
			if err2 != nil {
				err = err2
				return false
			} else if listResp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *listResp.Message)
				return false
			} else if listResp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			// 若实例在升级过程中，直接继续轮询
			if listResp.ReturnObj.List[0].ProdRunningStatus != business.MongodbRunningStatusStarted || listResp.ReturnObj.List[0].ProdOrderStatus != business.MongodbOrderStatusStarted {
				return true
			}

			detailResp, err2 := c.meta.Apis.SdkMongodbApis.MongodbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeader)
			if err2 != nil {
				if tolerateCount <= 0 {
					err2 = err
					return false
				}
				tolerateCount--
				return true
			} else if detailResp.StatusCode != 800 {
				if tolerateCount <= 0 {
					err = fmt.Errorf("API return error. Message: %s", *detailResp.Message)
					return false
				}
				tolerateCount--
				return true
			} else if detailResp.ReturnObj == nil {
				if tolerateCount <= 0 {
					err = common.InvalidReturnObjError
					return false
				}
				tolerateCount--
				return true
			}
			// 验证扩容结果（磁盘空间，备份磁盘空间，配置，shard数和prodId）
			masterDiskFlag := true
			backupDiskFlag := true
			// 验证master磁盘空间
			if planNodeInfoList.NodeType.ValueString() == "master" {
				diskSize := detailResp.ReturnObj.DiskSize
				if planNodeInfoList.StorageSpace.ValueInt32() != 0 && planNodeInfoList.StorageSpace.ValueInt32() != diskSize {
					masterDiskFlag = false
				}
			}
			if plan.IsUpgradeBackUp.ValueBool() || planNodeInfoList.NodeType.ValueString() == "backup" {
				if plan.BackupStorageType.ValueString() == business.BackupStorageTypeOS {
					backupDiskFlag = true
				}
				if detailResp.ReturnObj.Backup != nil {
					diskSize := detailResp.ReturnObj.Backup.Size[:len(detailResp.ReturnObj.Backup.Size)-1]
					expectedDiskSize := planNodeInfoList.StorageSpace.ValueInt32()
					if strings.Contains(state.ProdID.ValueString(), "Cluster") {
						expectedDiskSize = planNodeInfoList.StorageSpace.ValueInt32() * state.ShardNum.ValueInt32()
					}
					if planNodeInfoList.StorageSpace.ValueInt32() != 0 && fmt.Sprintf("%d", expectedDiskSize) != diskSize {
						backupDiskFlag = false
					}
				}
			}
			if masterDiskFlag && backupDiskFlag {
				return false
			}
			// 继续轮询
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，实例端口仍未更新成功！")
	}
	return
}

func (c *CtyunMongodbInstance) updateReadPort(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig) (err error) {
	_, err2 := c.PreCheckUpdateLoop(ctx, state, 60)
	if err2 != nil {
		err = err2
		return
	}
	updateParams := &mongodb.MongodbUpdatePortRequest{
		ProdInstId: state.ID.ValueString(),
		NewPort:    fmt.Sprintf("%d", plan.ReadPort.ValueInt32()),
	}
	updateHeader := &mongodb.MongodbUpdatePortRequestHeader{
		RegionID: state.RegionID.ValueString(),
	}
	if state.ProjectID.ValueString() != "" {
		updateHeader.ProjectID = state.ProjectID.ValueStringPointer()
	}
	resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbUpdatePortApi.Do(ctx, c.meta.Credential, updateParams, updateHeader)
	if err2 != nil {
		err = err2
		return
	} else if resp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	// 轮询确认端口更新完成
	err = c.UpdatePortLoop(ctx, state, plan, 60)
	return
}

func (c *CtyunMongodbInstance) generateAzInfo(ctx context.Context, config *CtyunMongodbInstanceConfig, prodType string, nodeType string) (AzInfoList []mongodb.AvailabilityZoneInfoRequest, err error) {
	// 获取az列表
	azList, err := c.getRegionAzInfoList(ctx, config)
	if err != nil {
		return nil, err
	}
	azNum := len(azList)

	if azNum <= 0 {
		err = errors.New("未查询到该资源池az信息，请稍后重试，或者手动填写az信息进行创建")
		return
	}
	// 因此1个az和2个az分配规则相同，只有3个及以上az的资源池有区别
	// mongodb分配节点规则：主节点与备用节点1、2需完全相同，或者完全不相同
	// 副本集和集群节点，还需要生成backup

	if prodType == "single" || nodeType == "backup" {
		var azInfo mongodb.AvailabilityZoneInfoRequest
		azInfo.NodeType = nodeType
		azInfo.AvailabilityZoneName = azList[0].AvailabilityZoneName
		azInfo.AvailabilityZoneCount = 1
		AzInfoList = append(AzInfoList, azInfo)
		return
	} else if prodType == "replica" {
		config.replicaNum = business.MongodbReplicaNodeNum[config.ProdID.ValueString()]

		if azNum >= 3 {
			distNodeNum := [3]int32{
				(int32(config.replicaNum) + 2) / 3,
				(int32(config.replicaNum) + 1) / 3,
				int32(config.replicaNum) / 3,
			}
			//nodeDist := business.MongodbReplicaNodeDistMap[config.replicaNum]
			// 有3个az，节点可以平均分摊在各个az下
			for idx, azItem := range azList {
				if len(AzInfoList) >= 3 {
					break
				}
				var azInfo mongodb.AvailabilityZoneInfoRequest
				azInfo.NodeType = nodeType
				azInfo.AvailabilityZoneCount = distNodeNum[idx]
				if azInfo.AvailabilityZoneCount <= 0 {
					continue
				}
				azInfo.AvailabilityZoneName = azItem.AvailabilityZoneName
				AzInfoList = append(AzInfoList, azInfo)
			}
		} else if azNum == 2 {
			distNodeNum := []int32{
				(int32(config.replicaNum) + 1) / 2,
				int32(config.replicaNum) / 2,
			}
			for idx, azItem := range azList {
				if len(AzInfoList) >= 2 {
					break
				}
				var azInfo mongodb.AvailabilityZoneInfoRequest
				azInfo.NodeType = nodeType
				azInfo.AvailabilityZoneName = azItem.AvailabilityZoneName
				azInfo.AvailabilityZoneCount = distNodeNum[idx]
				if azInfo.AvailabilityZoneCount <= 0 {
					continue
				}
				AzInfoList = append(AzInfoList, azInfo)
			}
		} else {
			var azInfo mongodb.AvailabilityZoneInfoRequest
			azInfo.NodeType = nodeType
			azInfo.AvailabilityZoneName = azList[0].AvailabilityZoneName
			azInfo.AvailabilityZoneCount = business.MongodbReplicaNodeDistMap[config.replicaNum]
			AzInfoList = append(AzInfoList, azInfo)
		}
		return
	} else if prodType == "cluster" {

		var azInfo mongodb.AvailabilityZoneInfoRequest
		// 默认为config的数量
		var nodeNum int32
		nodeNum = 3
		if nodeType == "mongos" {
			nodeNum = business.MongodbClusterNodeBaseNumMap[nodeType] * config.MongosNum.ValueInt32()
		} else if nodeType == "shard" {
			nodeNum = business.MongodbClusterNodeBaseNumMap[nodeType] * config.ShardNum.ValueInt32()
		}
		if azNum >= 3 {
			// 先计算每个AZ的节点数量
			distNodeNum := [3]int32{
				(int32(nodeNum) + 2) / 3,
				(int32(nodeNum) + 1) / 3,
				int32(nodeNum) / 3,
			}
			for idx, azItem := range azList {
				if len(AzInfoList) >= 3 {
					break
				}
				azInfo.NodeType = nodeType
				azInfo.AvailabilityZoneName = azItem.AvailabilityZoneName
				azInfo.AvailabilityZoneCount = distNodeNum[idx]
				if azInfo.AvailabilityZoneCount <= 0 {
					continue
				}
				AzInfoList = append(AzInfoList, azInfo)
			}
		} else if azNum == 2 {
			distNodeNum := []int32{
				(int32(nodeNum) + 1) / 2,
				int32(nodeNum) / 2,
			}
			for idx, azItem := range azList {
				if len(AzInfoList) >= 3 {
					break
				}
				azInfo.NodeType = nodeType
				azInfo.AvailabilityZoneName = azItem.AvailabilityZoneName
				azInfo.AvailabilityZoneCount = distNodeNum[idx]
				if azInfo.AvailabilityZoneCount <= 0 {
					continue
				}
				AzInfoList = append(AzInfoList, azInfo)
			}
		} else {
			// 处理单节点情况
			azInfo.NodeType = nodeType
			azInfo.AvailabilityZoneName = azList[0].AvailabilityZoneName
			azInfo.AvailabilityZoneCount = nodeNum
			AzInfoList = append(AzInfoList, azInfo)
		}
		return
	} else {
		err = errors.New("mongodb数据库类型有误")
		return
	}
}

func (c *CtyunMongodbInstance) getSingleNodeInfo(ctx context.Context, config *CtyunMongodbInstanceConfig, mongoNodeInfoList *[]mongodb.MongodbNodeInfoListRequest) (err error) {
	var mongoMasterNodeInfo mongodb.MongodbNodeInfoListRequest
	mongoMasterNodeInfo.NodeType = "s"
	mongoMasterNodeInfo.InstSpec = "1"
	mongoMasterNodeInfo.StorageType = config.StorageType.ValueString()
	mongoMasterNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()
	mongoMasterNodeInfo.Disks = 1
	mongoMasterNodeInfo.ProdPerformanceSpec = config.prodPerformanceSpec
	var mongoBackupNodeInfo mongodb.MongodbNodeInfoListRequest

	if config.BackupStorageType.ValueString() != business.MongodbBackupStorageTypeOS {
		mongoBackupNodeInfo.NodeType = "backup"
		mongoBackupNodeInfo.InstSpec = "1"
		mongoBackupNodeInfo.StorageType = config.BackupStorageType.ValueString()
		mongoBackupNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()
	}

	// 处理azInfo，若azInfo不为空，用用户输入的azInfo
	if !config.AvailabilityZoneInfo.IsNull() && !config.AvailabilityZoneInfo.IsUnknown() {
		var azZoneInfoList []AvailabilityZoneModel
		var azZoneInfo []mongodb.AvailabilityZoneInfoRequest
		var backupAzZoneInfo []mongodb.AvailabilityZoneInfoRequest
		diags := config.AvailabilityZoneInfo.ElementsAs(ctx, &azZoneInfoList, true)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		for _, azInfoItem := range azZoneInfoList {
			azZone := mongodb.AvailabilityZoneInfoRequest{
				AvailabilityZoneName:  azInfoItem.AvailabilityZoneName.ValueString(),
				AvailabilityZoneCount: azInfoItem.AvailabilityZoneCount.ValueInt32(),
				NodeType:              azInfoItem.NodeType.ValueString(),
			}
			if azZone.NodeType == "backup" {
				backupAzZoneInfo = append(backupAzZoneInfo, azZone)
			} else if azZone.NodeType == "master" {
				azZoneInfo = append(azZoneInfo, azZone)
			}
		}
		mongoMasterNodeInfo.AvailabilityZoneInfo = azZoneInfo

		if config.BackupStorageType.ValueString() != business.MongodbBackupStorageTypeOS {
			mongoBackupNodeInfo.AvailabilityZoneInfo = backupAzZoneInfo
		}
	} else {
		// 若azInfo为空，则生成az信息
		azInfo, err2 := c.generateAzInfo(ctx, config, "single", "master")
		if err2 != nil {
			err = err2
			return
		}
		mongoMasterNodeInfo.AvailabilityZoneInfo = azInfo

		// backup节点
		azInfo, err2 = c.generateAzInfo(ctx, config, "replica", "backup")
		if err2 != nil {
			err = err2
			return
		}
		mongoBackupNodeInfo.AvailabilityZoneInfo = azInfo
	}

	*mongoNodeInfoList = append(*mongoNodeInfoList, mongoMasterNodeInfo)

	if config.BackupStorageType.ValueString() != business.MongodbBackupStorageTypeOS {
		*mongoNodeInfoList = append(*mongoNodeInfoList, mongoBackupNodeInfo)
	}
	return
}

func (c *CtyunMongodbInstance) getReplicaNodeInfo(ctx context.Context, config *CtyunMongodbInstanceConfig, mongoNodeInfoList *[]mongodb.MongodbNodeInfoListRequest) (err error) {
	// 副本集需要一个master节点，和一个backup节点
	// master节点
	var mongoMasterNodeInfo mongodb.MongodbNodeInfoListRequest
	mongoMasterNodeInfo.NodeType = "ms"
	mongoMasterNodeInfo.InstSpec = "1"
	mongoMasterNodeInfo.Disks = 1
	mongoMasterNodeInfo.StorageType = config.StorageType.ValueString()
	mongoMasterNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()
	mongoMasterNodeInfo.ProdPerformanceSpec = config.prodPerformanceSpec

	// backup节点
	var mongoBackupNodeInfo mongodb.MongodbNodeInfoListRequest
	if config.BackupStorageType.ValueString() != business.MongodbBackupStorageTypeOS {
		mongoBackupNodeInfo.NodeType = "backup"
		mongoBackupNodeInfo.InstSpec = "1"
		mongoBackupNodeInfo.StorageType = config.BackupStorageType.ValueString()
		mongoBackupNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()
	}

	// 处理azInfo，若azInfo不为空，用用户输入的azInfo
	if !config.AvailabilityZoneInfo.IsNull() && !config.AvailabilityZoneInfo.IsUnknown() {
		var azZoneInfoList []AvailabilityZoneModel
		var masterAzZoneInfo []mongodb.AvailabilityZoneInfoRequest
		var backupAzZoneInfo []mongodb.AvailabilityZoneInfoRequest
		diags := config.AvailabilityZoneInfo.ElementsAs(ctx, &azZoneInfoList, true)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		for _, azInfoItem := range azZoneInfoList {
			azZone := mongodb.AvailabilityZoneInfoRequest{
				AvailabilityZoneName:  azInfoItem.AvailabilityZoneName.ValueString(),
				AvailabilityZoneCount: azInfoItem.AvailabilityZoneCount.ValueInt32(),
				NodeType:              azInfoItem.NodeType.ValueString(),
			}
			if azInfoItem.NodeType.ValueString() == "master" {
				masterAzZoneInfo = append(masterAzZoneInfo, azZone)
			} else if azInfoItem.NodeType.ValueString() == "backup" {
				backupAzZoneInfo = append(backupAzZoneInfo, azZone)
			}
		}
		mongoMasterNodeInfo.AvailabilityZoneInfo = masterAzZoneInfo
		if config.BackupStorageType.ValueString() != business.MongodbBackupStorageTypeOS {
			mongoBackupNodeInfo.AvailabilityZoneInfo = backupAzZoneInfo
		}
	} else {
		// 若azInfo为空，则生成az信息
		// master节点
		azInfo, err2 := c.generateAzInfo(ctx, config, "replica", "master")
		if err2 != nil {
			err = err2
			return
		}
		mongoMasterNodeInfo.AvailabilityZoneInfo = azInfo

		// backup节点
		azInfo, err2 = c.generateAzInfo(ctx, config, "replica", "backup")
		if err2 != nil {
			err = err2
			return
		}
		mongoBackupNodeInfo.AvailabilityZoneInfo = azInfo
	}
	*mongoNodeInfoList = append(*mongoNodeInfoList, mongoMasterNodeInfo)
	if config.BackupStorageType.ValueString() != business.MongodbBackupStorageTypeOS {
		*mongoNodeInfoList = append(*mongoNodeInfoList, mongoBackupNodeInfo)
	}
	return
}

func (c *CtyunMongodbInstance) getClusterNodeInfo(ctx context.Context, config *CtyunMongodbInstanceConfig, mongoNodeInfoList *[]mongodb.MongodbNodeInfoListRequest) (err error) {
	// 副本集需要一个mongos, shard, config节点
	// mongos节点
	var mongoMongosNodeInfo mongodb.MongodbNodeInfoListRequest
	mongoMongosNodeInfo.NodeType = "mongos"
	mongoMongosNodeInfo.InstSpec = "1"
	mongoMongosNodeInfo.Disks = 1
	mongoMongosNodeInfo.StorageType = config.StorageType.ValueString()
	mongoMongosNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()
	mongoMongosNodeInfo.ProdPerformanceSpec = config.prodPerformanceSpec
	// shard节点
	var mongoShardNodeInfo mongodb.MongodbNodeInfoListRequest
	mongoShardNodeInfo.NodeType = "shard"
	mongoShardNodeInfo.InstSpec = "1"
	mongoShardNodeInfo.Disks = 1
	mongoShardNodeInfo.StorageType = config.StorageType.ValueString()
	mongoShardNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()
	mongoShardNodeInfo.ProdPerformanceSpec = config.prodPerformanceSpec
	// config节点, config节点配置固定
	var mongoConfigNodeInfo mongodb.MongodbNodeInfoListRequest
	mongoConfigNodeInfo.NodeType = "config"
	mongoConfigNodeInfo.InstSpec = "1"
	mongoConfigNodeInfo.Disks = 1
	mongoConfigNodeInfo.StorageType = config.StorageType.ValueString()
	mongoConfigNodeInfo.StorageSpace = config.StorageSpace.ValueInt32()
	mongoConfigNodeInfo.ProdPerformanceSpec = "2C4G"
	// backup节点
	var mongoBackupNodeInfo mongodb.MongodbNodeInfoListRequest

	if config.BackupStorageType.ValueString() != business.MongodbBackupStorageTypeOS {
		mongoBackupNodeInfo.NodeType = "backup"
		mongoBackupNodeInfo.InstSpec = "1"
		mongoBackupNodeInfo.StorageType = config.BackupStorageType.ValueString()
		// backup节点磁盘空间 = shard数量*每个shard磁盘空间
		mongoBackupNodeInfo.StorageSpace = config.StorageSpace.ValueInt32() * config.ShardNum.ValueInt32()
	}
	// 处理azInfo，若azInfo不为空，用用户输入的azInfo
	if !config.AvailabilityZoneInfo.IsNull() && !config.AvailabilityZoneInfo.IsUnknown() {
		var azZoneInfoList []AvailabilityZoneModel
		var mongosAzZoneInfo []mongodb.AvailabilityZoneInfoRequest
		var shardAzZoneInfo []mongodb.AvailabilityZoneInfoRequest
		var configAzZoneInfo []mongodb.AvailabilityZoneInfoRequest
		var backupAzZoneInfo []mongodb.AvailabilityZoneInfoRequest
		diags := config.AvailabilityZoneInfo.ElementsAs(ctx, &azZoneInfoList, true)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		for _, azInfoItem := range azZoneInfoList {
			azZone := mongodb.AvailabilityZoneInfoRequest{
				AvailabilityZoneName:  azInfoItem.AvailabilityZoneName.ValueString(),
				AvailabilityZoneCount: azInfoItem.AvailabilityZoneCount.ValueInt32(),
				NodeType:              azInfoItem.NodeType.ValueString(),
			}
			if azInfoItem.NodeType.ValueString() == "mongos" {
				mongosAzZoneInfo = append(mongosAzZoneInfo, azZone)
			} else if azInfoItem.NodeType.ValueString() == "backup" {
				backupAzZoneInfo = append(backupAzZoneInfo, azZone)
			} else if azInfoItem.NodeType.ValueString() == "config" {
				configAzZoneInfo = append(configAzZoneInfo, azZone)
			} else if azInfoItem.NodeType.ValueString() == "shard" {
				shardAzZoneInfo = append(shardAzZoneInfo, azZone)
			}
		}
		mongoMongosNodeInfo.AvailabilityZoneInfo = mongosAzZoneInfo
		mongoShardNodeInfo.AvailabilityZoneInfo = shardAzZoneInfo
		mongoConfigNodeInfo.AvailabilityZoneInfo = configAzZoneInfo
		if config.BackupStorageType.ValueString() != business.MongodbBackupStorageTypeOS {
			mongoBackupNodeInfo.AvailabilityZoneInfo = backupAzZoneInfo
		}
	} else {
		// 若azInfo为空，则生成az信息
		// mongos节点
		azInfo, err2 := c.generateAzInfo(ctx, config, "cluster", "mongos")
		if err2 != nil {
			err = err2
			return
		}
		mongoMongosNodeInfo.AvailabilityZoneInfo = azInfo
		// shard
		azInfo, err2 = c.generateAzInfo(ctx, config, "cluster", "shard")
		if err2 != nil {
			err = err2
			return
		}
		mongoShardNodeInfo.AvailabilityZoneInfo = azInfo

		// config
		azInfo, err2 = c.generateAzInfo(ctx, config, "cluster", "config")
		if err2 != nil {
			err = err2
			return
		}
		mongoConfigNodeInfo.AvailabilityZoneInfo = azInfo

		// backup
		azInfo, err2 = c.generateAzInfo(ctx, config, "cluster", "backup")
		if err2 != nil {
			err = err2
			return
		}
		mongoBackupNodeInfo.AvailabilityZoneInfo = azInfo
	}

	*mongoNodeInfoList = append(*mongoNodeInfoList, mongoMongosNodeInfo)
	*mongoNodeInfoList = append(*mongoNodeInfoList, mongoShardNodeInfo)
	*mongoNodeInfoList = append(*mongoNodeInfoList, mongoConfigNodeInfo)
	if config.BackupStorageType.ValueString() != business.MongodbBackupStorageTypeOS {
		*mongoNodeInfoList = append(*mongoNodeInfoList, mongoBackupNodeInfo)
	}
	return
}

func (c *CtyunMongodbInstance) upgradeStorage(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig) (err error) {
	// 若plan阶段存储空间与state阶段不一致，触发更新
	if !plan.StorageSpace.IsNull() && !plan.StorageSpace.Equal(state.StorageSpace) {
		// 确定实例处于running状态
		_, err = c.PreCheckUpdateLoop(ctx, state, 60)
		nodeType := "master"
		updateParams := &mongodb.MongodbUpgradeRequest{
			InstId:          state.ID.ValueString(),
			DiskVolume:      plan.StorageSpace.ValueInt32Pointer(),
			IsUpgradeBackup: plan.IsUpgradeBackUp.ValueBoolPointer(),
			NodeType:        &nodeType,
		}
		updateHeader := &mongodb.MongodbUpgradeRequestHeader{}
		if state.ProjectID.ValueString() != "" {
			updateHeader.ProjectID = state.ProjectID.ValueStringPointer()
		}
		resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbUpgradeApi.Do(ctx, c.meta.Credential, updateParams, updateHeader)
		if err2 != nil {
			err = err2
			return
		} else if resp == nil {
			err = errors.New("当进行磁盘升配操作中，执行结果返回为nil。请确认未升配成功后，再重试。")
			return
		} else if resp.StatusCode != 200 {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
			return
		}

		var planNodeInfo NodeInfoListModel
		planNodeInfo.NodeType = types.StringValue(nodeType)
		planNodeInfo.StorageSpace = types.Int32Value(plan.StorageSpace.ValueInt32())

		// 轮询确认是否已扩容完成
		err = c.UpgradeStorageLoop(ctx, state, plan, planNodeInfo, 60)
		if err != nil {
			return
		}
	}
	return
}

func (c *CtyunMongodbInstance) getMongoDetailInfo(ctx context.Context, config *CtyunMongodbInstanceConfig) (detail *mongodb.DetailRespReturnObj, err error) {
	detailParams := &mongodb.MongodbQueryDetailRequest{
		ProdInstId: config.ID.ValueString(),
	}
	detailHeader := &mongodb.MongodbQueryDetailRequestHeaders{
		RegionID: config.RegionID.ValueString(),
	}
	if config.ProjectID.ValueString() != "" {
		detailHeader.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkMongodbApis.MongodbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeader)
	if err != nil {
		return
	} else if resp == nil {
		err = errors.New("获取mongodb实例为nil，请稍后再试！")
		return
	} else if resp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	detail = resp.ReturnObj
	return
}

func (c *CtyunMongodbInstance) getRegionAzInfoList(ctx context.Context, state *CtyunMongodbInstanceConfig) (azList []mongodb.TeledbGetAvailabilityZoneResponseReturnObjData, err error) {
	params := &mongodb.TeledbGetAvailabilityZoneRequest{
		RegionId: state.RegionID.ValueString(),
	}
	header := &mongodb.TeledbGetAvailabilityZoneRequestHeader{}
	// 1. 获取az信息
	resp, err := c.meta.Apis.SdkMongodbApis.TeledbGetAvailabilityZone.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = errors.New("查询az信息时返回为nil，请稍后再试")
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj.Data == nil {
		err = common.InvalidReturnObjError
		return
	}
	azList = resp.ReturnObj.Data
	return
}

func (c *CtyunMongodbInstance) checkSpec(ctx context.Context, plan *CtyunMongodbInstanceConfig) error {
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
	mysqlFlavor, err := c.mongodbService.GetMongodbFlavorByProdIdAndFlavorName(
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

func (c *CtyunMongodbInstance) upgradeSpec(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig) error {
	if plan.FlavorName.Equal(state.FlavorName) || plan.FlavorName.ValueString() == state.FlavorName.ValueString() {
		return nil
	}

	// 获取mongodb类型
	mongodbType := c.getMongodbType(state)
	if mongodbType == "" {
		return errors.New("prod_id 有误，请确认后再进行升配规格操作")
	}

	updateParams := &mongodb.MongodbUpgradeRequest{
		InstId: state.ID.ValueString(),
	}
	//fmt.Println(updateParams)
	updateHeader := &mongodb.MongodbUpgradeRequestHeader{}
	if state.ProjectID.ValueString() != "" {
		updateHeader.ProjectID = state.ProjectID.ValueStringPointer()
	}

	var azInfo []mongodb.AvailabilityZoneInfo
	// 默认为单机版node type
	nodeType := "s"
	if mongodbType == business.MongodbProdTypeReplica {
		nodeType = "ms"
	}
	// mongodb 副本级和单节点扩容规格类似
	if mongodbType == business.MongodbProdTypeSingle || mongodbType == business.MongodbProdTypeReplica {
		// 确认是否需要扩容
		// 若plan spec和state spec 相同，无需变配
		if plan.prodPerformanceSpec == state.prodPerformanceSpec {
			return nil
		}
		// 获取az节点规格
		err := c.getNodeInfo(ctx, state, plan, mongodbType, &azInfo, nodeType)
		if err != nil {
			return err
		}
	} else if mongodbType == business.MongodbProdTypeCluster {
		err := c.getNodeInfo(ctx, state, plan, mongodbType, &azInfo, plan.UpgradeNodeType.ValueString())
		if err != nil {
			return err
		}
	}

	spec := plan.prodPerformanceSpec
	updateParams.ProdPerformanceSpec = &spec
	updateParams.AzList = azInfo
	// 调用升配接口
	resp, err := c.meta.Apis.SdkMongodbApis.MongodbUpgradeApi.Do(ctx, c.meta.Credential, updateParams, updateHeader)
	if err != nil {
		return err
	} else if resp == nil {
		return errors.New("升配返回为nil，可联系研发确认具体原因")
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	// 轮询确保接口执行后进入升配状态
	err = c.afterUpdateSpecLoop(ctx, state)
	if err != nil {
		return err
	}
	return nil
}

func (c *CtyunMongodbInstance) getMongodbType(state *CtyunMongodbInstanceConfig) string {
	prodId := state.ProdID.ValueString()
	if strings.Contains(prodId, "Single") {
		return business.MongodbProdTypeSingle
	} else if strings.Contains(prodId, "Replica") {
		return business.MongodbProdTypeReplica
	} else if strings.Contains(prodId, "Cluster") {
		return business.MongodbProdTypeCluster
	}
	return ""
}

// 通过mongo类型，获取mongo各个节点的信息
func (c *CtyunMongodbInstance) getNodeInfo(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig, mongodbType string, azInfo *[]mongodb.AvailabilityZoneInfo, nodeType string) error {
	var err error
	// 若az info为空，或者产品类型也需要更新的时候则
	//
	if plan.AvailabilityZoneInfo.IsNull() || !plan.AvailabilityZoneInfo.IsNull() && !plan.ProdID.Equal(state.ProdID) {
		nodeDist := make(map[string]int32)
		// 获取节点AZ分布
		nodeDist, err = c.getNodeDist(ctx, state, mongodbType, nodeType)
		if err != nil {
			return err
		}
		for az, nodeNum := range nodeDist {
			var azItem mongodb.AvailabilityZoneInfo
			azItem.AvailabilityZoneCount = nodeNum
			azItem.AvailabilityZoneName = az
			azItem.NodeType = &nodeType
			if azItem.AvailabilityZoneCount <= 0 {
				continue
			}
			*azInfo = append(*azInfo, azItem)
		}
	} else {
		var azModelList []AvailabilityZoneModel
		diag := plan.AvailabilityZoneInfo.ElementsAs(ctx, &azModelList, true)
		if diag.HasError() {
			err = errors.New(diag[0].Detail())
			return err
		}
		for _, azModelItem := range azModelList {
			var azItem mongodb.AvailabilityZoneInfo
			azItem.AvailabilityZoneCount = azModelItem.AvailabilityZoneCount.ValueInt32()
			azItem.AvailabilityZoneName = azModelItem.AvailabilityZoneName.ValueString()
			azItem.NodeType = azModelItem.NodeType.ValueStringPointer()
			if azItem.AvailabilityZoneCount <= 0 {
				continue
			}
			*azInfo = append(*azInfo, azItem)
		}
	}

	return nil
}

func (c *CtyunMongodbInstance) getNodeDist(ctx context.Context, state *CtyunMongodbInstanceConfig, mongodbType string, nodeType string) (map[string]int32, error) {
	nodeDist := make(map[string]int32)

	// 获取实例详情，单机版和副本集通过查询实例详情。集群版直接生成AZ分布信息
	if mongodbType == business.MongodbProdTypeSingle || mongodbType == business.MongodbProdTypeReplica {
		mongoDetailInfo, err := c.getMongoDetailInfo(ctx, state)
		if err != nil {
			return nil, err
		}
		nodeInfoList := mongoDetailInfo.NodeInfoVOS
		for _, nodeInfo := range nodeInfoList {
			azId := *nodeInfo.AzId
			if _, exists := nodeDist[azId]; exists {
				nodeDist[azId] = nodeDist[azId] + 1
			} else {
				nodeDist[azId] = 1
			}
		}
	} else if mongodbType == business.MongodbProdTypeCluster {
		// 获取az列表
		azList, err := c.getRegionAzInfoList(ctx, state)
		if err != nil {
			return nil, err
		}
		azNum := len(azList)
		if azNum <= 0 {
			return nil, errors.New("获取az信息失败，接口返回长度为0")
		}
		// 集群版直接重新生成一个az分布
		var nodeNum int32
		nodeNum = 3
		if nodeType == business.MongodbNodeTypeMongos {
			nodeNum = business.MongodbClusterNodeBaseNumMap[nodeType] * state.MongosNum.ValueInt32()
		} else if nodeType == business.MongodbNodeTypeShard {
			nodeNum = business.MongodbClusterNodeBaseNumMap[nodeType] * state.ShardNum.ValueInt32()
		}
		//if nodeType == business.MongodbNodeTypeMongos {
		//} else if nodeType == business.MongodbNodeTypeShard {
		//	nodeNum = business.MongodbClusterNodeBaseNumMap[nodeType] * state.ShardNum.ValueInt32()
		//}
		if azNum >= 3 {
			// 先计算每个AZ的节点数量
			distNodeNum := []int32{
				(int32(nodeNum) + 2) / 3,
				(int32(nodeNum) + 1) / 3,
				int32(nodeNum) / 3,
			}
			for idx, azItem := range azList {
				azId := azItem.AvailabilityZoneId
				if _, exists := nodeDist[azId]; exists {
					continue
				} else {
					nodeDist[azId] = distNodeNum[idx]
				}
			}
		} else if azNum == 2 {
			// 处理2AZ节点
			distNodeNum := []int32{
				(int32(nodeNum) + 1) / 2,
				int32(nodeNum) / 2,
			}
			for idx, azItem := range azList {
				azId := azItem.AvailabilityZoneId
				if _, exists := nodeDist[azId]; exists {
					continue
				} else {
					nodeDist[azId] = distNodeNum[idx]
				}
			}
		} else {
			// 处理单节点情况
			azId := azList[0].AvailabilityZoneId
			nodeDist[azId] = nodeNum
		}
	}
	return nodeDist, nil
}

func (c *CtyunMongodbInstance) upgradeNode(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig) error {
	// 先确认mongodb实例类型
	// 扩容节点
	// 单集群不支持扩容，cluster集群支持扩容
	// 集群版支持扩容shard数量和mongos数量
	var err error
	// 如果plan和state阶段的prod id一致，无需变化

	mongodbType := c.getMongodbType(state)

	var azInfo []mongodb.AvailabilityZoneInfo
	nodeType := plan.UpgradeNodeType.ValueString()
	if mongodbType == business.MongodbProdTypeReplica {
		if state.ProdID.Equal(plan.ProdID) || state.ProdID.ValueString() == plan.ProdID.ValueString() {
			return nil
		}
		nodeType = "master"
	} else if mongodbType == business.MongodbProdTypeCluster {
		if !plan.UpgradeNodeType.IsNull() {
			if plan.UpgradeNodeType.ValueString() == business.MongodbNodeTypeShard {
				if plan.ShardNum.Equal(state.ShardNum) {
					return nil
				}
			} else if plan.UpgradeNodeType.ValueString() == business.MongodbNodeTypeMongos {
				if plan.MongosNum.Equal(state.MongosNum) {
					return nil
				}
			}
		} else {
			return nil
		}
	}
	upgradeParams := mongodb.MongodbUpgradeRequest{
		InstId:   state.ID.ValueString(),
		NodeType: &nodeType,
	}

	upgradeHeader := mongodb.MongodbUpgradeRequestHeader{}
	if !state.ProjectID.IsNull() {
		upgradeHeader.ProjectID = state.ProjectID.ValueStringPointer()
	}

	if mongodbType == business.MongodbProdTypeSingle {
		return nil
	} else if mongodbType == business.MongodbProdTypeReplica {
		// 根据prodid判断需要扩容几个副本
		// 计算出需要扩容的节点数量，根据原来节点分布情况。平摊新增节点分布
		err = c.getUpgradeReplicaAzList(ctx, state, plan, &azInfo)
		if err != nil {
			return err
		}
		prodId := business.MongodbProdIDDict[plan.ProdID.ValueString()]
		upgradeParams.ProdId = &prodId
	} else if mongodbType == business.MongodbProdTypeCluster {
		// 根据shard_num 和 mongos_num 判断
		// 集群版本扩容需要填写全部的az信息，因此直接生成一个新节点分布azList
		if plan.UpgradeNodeType.ValueString() == business.MongodbNodeTypeShard {
			if plan.ShardNum.IsNull() || plan.ShardNum.IsUnknown() {
				return errors.New("shard num为空")
			}
			if plan.ShardNum.ValueInt32() <= state.ShardNum.ValueInt32() {
				return errors.New("shard num有误，plan阶段shard num <= 原shard num")
			}
		} else if plan.UpgradeNodeType.ValueString() == business.MongodbNodeTypeMongos {
			if plan.MongosNum.IsNull() || plan.MongosNum.IsUnknown() {
				return errors.New("mongos num为空")
			}
			if plan.MongosNum.ValueInt32() <= state.MongosNum.ValueInt32() {
				return errors.New("mongos num有误，plan阶段mongos num <= 原mongos num")
			}
		}
		err = c.getUpgradeClusterAzList(ctx, state, plan, &azInfo)
		if err != nil {
			return err
		}
		state.MongosNum = plan.MongosNum
		state.ShardNum = plan.ShardNum
	}
	upgradeParams.AzList = azInfo
	resp, err := c.meta.Apis.SdkMongodbApis.MongodbUpgradeApi.Do(ctx, c.meta.Credential, &upgradeParams, &upgradeHeader)
	if err != nil {
		return err
	} else if resp == nil {
		return errors.New("升配失败，返回为nil。请联系研发确认")
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	err = c.afterUpdateSpecLoop(ctx, state)
	state.ProdID = plan.ProdID
	return nil
}

func (c *CtyunMongodbInstance) getUpgradeReplicaAzList(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig, azInfoList *[]mongodb.AvailabilityZoneInfo) error {
	if plan.AvailabilityZoneInfo.IsNull() {
		// 定义一个map，用于存放新增的节点，和az分布
		addNodeMap := make(map[string]int32)
		// 计算需要升配的节点数量
		stateNodeNum := business.MongodbReplicaNodeNum[state.ProdID.ValueString()]
		planNodeNum := business.MongodbReplicaNodeNum[plan.ProdID.ValueString()]
		addNum := planNodeNum - stateNodeNum
		if addNum <= 0 {
			return errors.New("plan阶段 prodID有误")
		}
		// 获取原节点分布
		nodeDist, err := c.getNodeDist(ctx, state, business.MongodbProdTypeReplica, "ms")
		if err != nil {
			return err
		}
		if len(nodeDist) <= 0 {
			return errors.New("原各个副本-az信息获取为空")
		}
		i := int32(0)
		for ; i < addNum; i++ {
			// 遍历寻找最少node节点的az
			minNodeNum := int32(math.MaxInt32)
			minAzId := ""
			for az, count := range nodeDist {
				if count < minNodeNum {
					minNodeNum = count
					minAzId = az
				}
			}
			if minAzId == "" {
				return errors.New("原各个副本-az信息获取为空")
			}
			if _, exist := addNodeMap[minAzId]; !exist {
				addNodeMap[minAzId] = 1
			} else {
				addNodeMap[minAzId] += 1
			}
			nodeDist[minAzId] += 1
		}

		// 生成azInfoList
		for az, count := range addNodeMap {
			var azInfoItem mongodb.AvailabilityZoneInfo
			nodeType := "ms"
			azInfoItem.NodeType = &nodeType
			azInfoItem.AvailabilityZoneCount = count
			azInfoItem.AvailabilityZoneName = az
			*azInfoList = append(*azInfoList, azInfoItem)
		}
	} else {
		// 如果用户自行输入，使用用户输入的az 信息
		var azZoneInfoList []AvailabilityZoneModel
		diags := plan.AvailabilityZoneInfo.ElementsAs(ctx, &azZoneInfoList, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		for _, azItem := range azZoneInfoList {
			var az mongodb.AvailabilityZoneInfo
			az.AvailabilityZoneName = azItem.AvailabilityZoneName.ValueString()
			az.AvailabilityZoneCount = azItem.AvailabilityZoneCount.ValueInt32()
			az.NodeType = azItem.NodeType.ValueStringPointer()
			*azInfoList = append(*azInfoList, az)
		}
	}
	return nil
}

func (c *CtyunMongodbInstance) getUpgradeClusterAzList(ctx context.Context, state *CtyunMongodbInstanceConfig, plan *CtyunMongodbInstanceConfig, azInfoList *[]mongodb.AvailabilityZoneInfo) error {
	if plan.AvailabilityZoneInfo.IsNull() {
		// 确定升级的节点类型
		azList, err := c.getRegionAzInfoList(ctx, state)
		if err != nil {
			return err
		}
		azNum := len(azList)

		if plan.UpgradeNodeType.ValueString() == business.MongodbNodeTypeShard {
			addShardNum := plan.ShardNum.ValueInt32() - state.ShardNum.ValueInt32()
			if addShardNum < 0 {
				err = fmt.Errorf("输入的shard_num 有误，实例当前shard_num=%d，需要升配至%d，暂不支持降配。", state.ShardNum, plan.ShardNum)
				return err
			} else if addShardNum == 0 {
				return nil
			}
			state.ShardNum = plan.ShardNum
			var nodeNum int32
			nodeNum = addShardNum * 3

			nodeType := business.MongodbNodeTypeShard
			var azInfo mongodb.AvailabilityZoneInfo
			if azNum >= 3 {
				// 先计算每个AZ的节点数量
				distNodeNum := [3]int32{
					(int32(nodeNum) + 2) / 3,
					(int32(nodeNum) + 1) / 3,
					int32(nodeNum) / 3,
				}
				for idx, azItem := range azList {
					azInfo.NodeType = &nodeType
					azInfo.AvailabilityZoneName = azItem.AvailabilityZoneName
					azInfo.AvailabilityZoneCount = distNodeNum[idx]
					if azInfo.AvailabilityZoneCount <= 0 {
						continue
					}
					*azInfoList = append(*azInfoList, azInfo)
				}
			} else if azNum == 2 {
				distNodeNum := []int32{
					(int32(nodeNum) + 1) / 2,
					int32(nodeNum) / 2,
				}
				for idx, azItem := range azList {
					azInfo.NodeType = &nodeType
					azInfo.AvailabilityZoneName = azItem.AvailabilityZoneName
					azInfo.AvailabilityZoneCount = distNodeNum[idx]
					if azInfo.AvailabilityZoneCount <= 0 {
						continue
					}
					*azInfoList = append(*azInfoList, azInfo)
				}
			} else {
				// 处理单节点情况
				azInfo.NodeType = &nodeType
				azInfo.AvailabilityZoneName = azList[0].AvailabilityZoneName
				azInfo.AvailabilityZoneCount = nodeNum
				*azInfoList = append(*azInfoList, azInfo)
			}
		} else if plan.UpgradeNodeType.ValueString() == business.MongodbNodeTypeMongos {
			addMongosNum := plan.MongosNum.ValueInt32() - state.MongosNum.ValueInt32()
			if addMongosNum < 0 {
				err = fmt.Errorf("输入的mongos_num 有误，实例当前mongos_num=%d，需要升配至%d，暂不支持降配。", state.MongosNum, plan.MongosNum)
				return err
			} else if addMongosNum == 0 {
				return nil
			}

			state.MongosNum = plan.MongosNum
			var nodeNum int32
			nodeNum = addMongosNum

			nodeType := business.MongodbNodeTypeMongos
			var azInfo mongodb.AvailabilityZoneInfo
			if azNum >= 3 {
				// 先计算每个AZ的节点数量
				distNodeNum := [3]int32{
					(int32(nodeNum) + 2) / 3,
					(int32(nodeNum) + 1) / 3,
					int32(nodeNum) / 3,
				}
				for idx, azItem := range azList {
					azInfo.NodeType = &nodeType
					azInfo.AvailabilityZoneName = azItem.AvailabilityZoneName
					azInfo.AvailabilityZoneCount = distNodeNum[idx]
					if azInfo.AvailabilityZoneCount <= 0 {
						continue
					}
					*azInfoList = append(*azInfoList, azInfo)
				}
			} else if azNum == 2 {
				distNodeNum := []int32{
					(int32(nodeNum) + 1) / 2,
					int32(nodeNum) / 2,
				}
				for idx, azItem := range azList {
					azInfo.NodeType = &nodeType
					azInfo.AvailabilityZoneName = azItem.AvailabilityZoneName
					azInfo.AvailabilityZoneCount = distNodeNum[idx]
					if azInfo.AvailabilityZoneCount <= 0 {
						continue
					}
					*azInfoList = append(*azInfoList, azInfo)
				}
			} else {
				// 处理单节点情况
				azInfo.NodeType = &nodeType
				azInfo.AvailabilityZoneName = azList[0].AvailabilityZoneName
				azInfo.AvailabilityZoneCount = nodeNum
				*azInfoList = append(*azInfoList, azInfo)
			}
		} else {
			return errors.New("输入的升级节点类型有误，仅支持mongos节点和shard节点升级")
		}
	} else {
		// 如果用户自行输入，使用用户输入的az 信息
		var azZoneInfoList []AvailabilityZoneModel
		diags := plan.AvailabilityZoneInfo.ElementsAs(ctx, &azZoneInfoList, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		for _, azItem := range azZoneInfoList {
			var az mongodb.AvailabilityZoneInfo
			az.AvailabilityZoneName = azItem.AvailabilityZoneName.ValueString()
			az.AvailabilityZoneCount = azItem.AvailabilityZoneCount.ValueInt32()
			az.NodeType = azItem.NodeType.ValueStringPointer()
			*azInfoList = append(*azInfoList, az)
		}
	}
	return nil
}

func (c *CtyunMongodbInstance) afterUpdateSpecLoop(ctx context.Context, state *CtyunMongodbInstanceConfig, loopCount ...int) error {
	var err error
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	params := &mongodb.MongodbGetListRequest{
		PageNow:      1,
		PageSize:     100,
		ProdInstName: state.Name.ValueStringPointer(),
	}
	header := &mongodb.MongodbGetListHeaders{
		RegionID: state.RegionID.ValueString(),
	}
	if !state.ProjectID.IsNull() && !state.ProjectID.IsUnknown() {
		header.ProjectID = state.ProjectID.ValueStringPointer()
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbGetListApi.Do(ctx, c.meta.Credential, params, header)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			runningStatus := resp.ReturnObj.List[0].ProdRunningStatus
			// 若实例状态已经运行正常，跳出轮询
			if runningStatus != business.MongodbRunningStatusStarted {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，实例更新后仍未进入更新状态！")
	}
	return nil
}

type CtyunMongodbInstanceConfig struct {
	CycleType               types.String `tfsdk:"cycle_type"`                // 计费模式： 1是包周期，2是按需
	RegionID                types.String `tfsdk:"region_id"`                 // 资源池Id
	VpcID                   types.String `tfsdk:"vpc_id"`                    // 虚拟私有云Id
	FlavorName              types.String `tfsdk:"flavor_name"`               // 规格名称
	SubnetID                types.String `tfsdk:"subnet_id"`                 // 子网Id
	SecurityGroupID         types.String `tfsdk:"security_group_id"`         // 安全组
	Name                    types.String `tfsdk:"name"`                      // 集群名称
	Password                types.String `tfsdk:"password"`                  // 管理员密码（RSA公钥加密）
	CycleCount              types.Int32  `tfsdk:"cycle_count"`               // 购买时长：单位月（范围：1-36）
	AutoRenew               types.Bool   `tfsdk:"auto_renew"`                // 自动续订状态（0-不自动续订，1-自动续订）
	ProdID                  types.String `tfsdk:"prod_id"`                   // 产品id
	ProjectID               types.String `tfsdk:"project_id"`                // 项目ID
	MasterOrderID           types.String `tfsdk:"master_order_id"`           // 订单ID
	ID                      types.String `tfsdk:"id"`                        // 实例ID
	ReadPort                types.Int32  `tfsdk:"read_port"`                 // 读端口
	InnodbBufferPoolSize    types.String `tfsdk:"innodb_buffer_pool_size"`   // 缓存池大小
	InnodbThreadConcurrency types.Int64  `tfsdk:"innodb_thread_concurrency"` // 线程数
	ProdRunningStatus       types.Int32  `tfsdk:"prod_running_status"`       // 实例运行状态: 0->运行正常, 1->重启中, 2-备份操作中,3->恢复操作中,4->转换ssl,5->异常,6->修改参数组中,7->已冻结,8->已注销,9->施工中,10->施工失败,11->扩容中,12->主备切换中
	ProdRunningStatusDesc   types.String `tfsdk:"prod_running_status_desc"`  // prod_running_status的解释版本
	EipID                   types.String `tfsdk:"eip_id"`                    // eip id
	IsUpgradeBackUp         types.Bool   `tfsdk:"is_upgrade_back_up"`        // DDS模块磁盘扩容时候会使用 是否主磁盘与备磁盘一起扩容
	HostIp                  types.String `tfsdk:"host_ip"`                   // 主机ip
	StorageType             types.String `tfsdk:"storage_type"`              // 存储类型
	StorageSpace            types.Int32  `tfsdk:"storage_space"`             // 存储空间
	AvailabilityZoneInfo    types.List   `tfsdk:"availability_zone_info"`    // 节点可用区信息
	ShardNum                types.Int32  `tfsdk:"shard_num"`                 // 当实例为集群版，shard数量
	MongosNum               types.Int32  `tfsdk:"mongos_num"`                // 当实例为集群版，mongos节点数量
	BackupStorageSpace      types.Int32  `tfsdk:"backup_storage_space"`      // 备用节点磁盘空间
	BackupStorageType       types.String `tfsdk:"backup_storage_type"`       // 备份节点存储类型
	UpgradeNodeType         types.String `tfsdk:"upgrade_node_type"`         // 集群版mongodb升配规格时，

	prodPerformanceSpec string
	instanceSeries      string
	hostType            string
	osType              string
	cpuType             string
	replicaNum          int32
}

type NodeInfoListModel struct {
	NodeType             types.String `tfsdk:"node_type"`              // 实例类型：master 或 readNode
	InstanceSeries       types.String `tfsdk:"instance_series"`        // 实例规格（实例类型，1=通用型，2=计算增强型，3=内存优化型，4=直通（未用到））
	StorageType          types.String `tfsdk:"storage_type"`           // 存储类型：SSD, SATA, SAS, SSD-genric, FAST-SSD
	StorageSpace         types.Int32  `tfsdk:"storage_space"`          // 存储空间（单位：GB，范围100到32768）
	ProdPerformanceSpec  types.String `tfsdk:"prod_performance_spec"`  // 规格（例：4C8G）
	AvailabilityZoneInfo types.List   `tfsdk:"availability_zone_info"` // 可用区信息
}

type AvailabilityZoneModel struct {
	AvailabilityZoneName  types.String `tfsdk:"availability_zone_name"`  // 资源池可用区名称
	AvailabilityZoneCount types.Int32  `tfsdk:"availability_zone_count"` // 资源池可用区总数
	NodeType              types.String `tfsdk:"node_type"`               // 表示分布AZ的节点类型，master/slave/readNode
}
