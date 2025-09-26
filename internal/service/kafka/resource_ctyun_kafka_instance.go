package kafka

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgkafka "github.com/ctyun-it/terraform-provider-ctyun/internal/core/kafka"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &ctyunKafkaInstance{}
	_ resource.ResourceWithConfigure   = &ctyunKafkaInstance{}
	_ resource.ResourceWithImportState = &ctyunKafkaInstance{}
)

type ctyunKafkaInstance struct {
	meta       *common.CtyunMetadata
	vpcService *business.VpcService
	sgService  *business.SecurityGroupService
}

func NewCtyunKafkaInstance() resource.Resource {
	return &ctyunKafkaInstance{}
}

func (c *ctyunKafkaInstance) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_kafka_instance"
}

type CtyunKafkaInstanceConfig struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	MasterOrderID       types.String `tfsdk:"master_order_id"`
	RegionID            types.String `tfsdk:"region_id"`
	ProjectID           types.String `tfsdk:"project_id"`             /*  企业项目ID(默认值：0)。您可以通过 <a href="https://www.ctyun.cn/document/10017248/10017965">查看企业项目资源</a> 获取企业项目ID。  */
	InstanceName        types.String `tfsdk:"instance_name"`          /*  实例名称。<br>规则：长度4~40个字符，大小写字母开头，只能包含大小写字母、数字及分隔符(-)，大小写字母或数字结尾，实例名称不可重复。  */
	EngineVersion       types.String `tfsdk:"engine_version"`         /*  实例的引擎版本，默认为3.6。<li>2.8：2.8.x的引擎版本<li>3.6：3.6.x的引擎版本  */
	SpecName            types.String `tfsdk:"spec_name"`              /*  实例的规格类型，资源池所具备的规格可通过查询产品规格接口获取，默认可选如下：<br>计算增强型的规格可选为：<li>kafka.2u4g.cluster<li>kafka.4u8g.cluster<li>kafka.8u16g.cluster<li>kafka.12u24g.cluster<li>kafka.16u32g.cluster<li>kafka.24u48g.cluster<li>kafka.32u64g.cluster<li>kafka.48u96g.cluster<li>kafka.64u128g.cluster <br>海光-计算增强型的规格可选为：<li>kafka.hg.2u4g.cluster<li>kafka.hg.4u8g.cluster<li>kafka.hg.8u16g.cluster<li>kafka.hg.16u32g.cluster<li>kafka.hg.32u64g.cluster <br>鲲鹏-计算增强型的规格可选为：<li>kafka.kp.2u4g.cluster<li>kafka.kp.4u8g.cluster<li>kafka.kp.8u16g.cluster<li>kafka.kp.16u32g.cluster<li>kafka.kp.32u64g.cluster  */
	NodeNum             types.Int32  `tfsdk:"node_num"`               /*  节点数。单机版为1个，集群版3~50个。  */
	ZoneList            types.Set    `tfsdk:"zone_list"`              /*  实例所在可用区信息。只能填一个（单可用区）或三个（多可用区），可用区信息可调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87&isNormal=1&vid=81">资源池可用区查询</a>API接口查询。  */
	DiskType            types.String `tfsdk:"disk_type"`              /*  磁盘类型，资源池所具备的磁盘类型可通过查询产品规格接口获取，默认取值：<li>SAS：高IO<li>SSD：超高IO<li>FAST-SSD：极速型SSD  */
	DiskSize            types.Int32  `tfsdk:"disk_size"`              /*  单个节点的磁盘存储空间，单位为GB，存储空间取值范围100GB ~ 10000，并且为100的倍数。实例总存储空间为diskSize * nodeNum。  */
	VpcID               types.String `tfsdk:"vpc_id"`                 /*  VPC网络ID。获取方法如下：<li>方法一：登录网络控制台界面，在虚拟私有云的详情页面查找VPC ID。<li>方法二：您可以通过 <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4814&data=94&vid=88">查询VPC列表</a> vpcID字段获取。  */
	SubnetID            types.String `tfsdk:"subnet_id"`              /*  VPC子网ID。获取方法如下：<li>方法一：登录网络控制台界面，单击VPC下的子网，进入子网详情页面，查找子网ID。<li>方法二：您可以通过 <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8659&data=94&vid=88">查询子网列表</a> subnetID字段获取。  */
	SecurityGroupID     types.String `tfsdk:"security_group_id"`      /*  安全组ID。获取方法如下：<li>方法一：登录网络控制台界面，在安全组的详情页面查找安全组ID。<li>方法二：您可以通过 <a href="https://eop.ctyun.cn/ebp/searchCtapi/ctApiDebug?product=18&api=4817&vid=88">查询用户安全组列表</a> id字段获取。  */
	EnableIpv6          types.Bool   `tfsdk:"enable_ipv6"`            /*  是否启用IPv6，默认为false。<li>true：启用IPv6。<li>false：不启用IPv6，默认值。  */
	PlainPort           types.Int32  `tfsdk:"plain_port"`             /*  公共接入点(PLAINTEXT)端口，范围在8000到9100之间，默认为8090。  */
	SaslPort            types.Int32  `tfsdk:"sasl_port"`              /*  安全接入点(SASL_PLAINTEXT)端口，范围在8000到9100之间，默认为8092。  */
	SslPort             types.Int32  `tfsdk:"ssl_port"`               /*  SSL接入点(SASL_SSL)端口，范围在8000到9100之间，默认为8098。  */
	HttpPort            types.Int32  `tfsdk:"http_port"`              /*  HTTP接入点端口，范围在8000到9100之间，默认为8082。  */
	RetentionHours      types.Int32  `tfsdk:"retention_hours"`        /*  实例消息保留时长，默认为72小时，可选1~10000小时。  */
	CycleType           types.String `tfsdk:"cycle_type"`             /*  按需: on_demand, 包月：month */
	CycleCount          types.Int32  `tfsdk:"cycle_count"`            /*  付费周期，单位为月，取值：1~6,12,24,36。  */
	AutoRenew           types.Bool   `tfsdk:"auto_renew"`             /*  过期是否自动续订。，默认为false。<li>true：自动续订。<li>false：不自动续订，默认值。  */
	AutoRenewCycleCount types.Int32  `tfsdk:"auto_renew_cycle_count"` /*  自动续订时间长，当autoRenewStatus为true时必填，取值：1~6。  */
}

func (c *ctyunKafkaInstance) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029624/10030700`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "名称",
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "主订单号",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"instance_name": schema.StringAttribute{
				Required:    true,
				Description: "实例名称，长度4~40个字符，大小写字母开头，只能包含大小写字母、数字及分隔符(-)，大小写字母或数字结尾，实例名称不可重复，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 40),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z0-9](?:[a-zA-Z0-9]|[-][a-zA-Z0-9])+$"), "实例名称不符合规则"),
				},
			},
			"engine_version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "实例引擎版本，支持2.8和3.6，默认3.6",
				Validators: []validator.String{
					stringvalidator.OneOf(business.KafkaVersion28, business.KafkaVersion36),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString(business.KafkaVersion36),
			},
			"spec_name": schema.StringAttribute{
				Required:    true,
				Description: "实例的规格类型，建议使用ctyun_kafka_specs查看，也可查看<a href=\"https://www.ctyun.cn/document/10029624/10030704\">产品规格说明</a>，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"node_num": schema.Int32Attribute{
				Required:    true,
				Description: "节点数。单机版为1个，集群版3~50个，支持更新，不支持缩容",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},
			"zone_list": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "实例所在可用区信息，只能传一个或三个可用区，可通过ctyun_regions查看",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(stringvalidator.UTF8LengthAtLeast(1)),
				},
			},
			"disk_type": schema.StringAttribute{
				Required:    true,
				Description: "磁盘类型，建议使用ctyun_kafka_specs查看，通常支持SAS、SSD、FAST-SSD",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"disk_size": schema.Int32Attribute{
				Required:    true,
				Description: "单个节点的磁盘存储空间，单位为GB，存储空间取值范围100GB ~ 10000，并且为100的倍数。实例总存储空间为diskSize * nodeNum，支持更新，不支持缩容",
				Validators: []validator.Int32{
					int32validator.Between(100, 10000),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "虚拟私有云ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"security_group_id": schema.StringAttribute{
				Required:    true,
				Description: "安全组ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SecurityGroupValidate(),
				},
			},
			"enable_ipv6": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否启用IPv6，默认为false",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"plain_port": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Description: "公共接入点(PLAINTEXT)端口，范围在8000到9100之间，默认为8090",
				Validators: []validator.Int32{
					int32validator.Between(8000, 9100),
				},
				Default: int32default.StaticInt32(8090),
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"sasl_port": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Description: "安全接入点(SASL_PLAINTEXT)端口，范围在8000到9100之间，默认为8092",
				Validators: []validator.Int32{
					int32validator.Between(8000, 9100),
				},
				Default: int32default.StaticInt32(8092),
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"ssl_port": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Description: "SSL接入点(SASL_SSL)端口，范围在8000到9100之间，默认为8098。",
				Validators: []validator.Int32{
					int32validator.Between(8000, 9100),
				},
				Default: int32default.StaticInt32(8098),
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"http_port": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Description: "HTTP接入点端口，范围在8000到9100之间，默认为8082",
				Validators: []validator.Int32{
					int32validator.Between(8000, 9100),
				},
				Default: int32default.StaticInt32(8082),
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"retention_hours": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "实例消息保留时长，单位小时。默认为72小时，可选1~10000小时，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
				Default: int32default.StaticInt32(72),
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围：month：按月，on_demand：按需。当此值为month时，cycle_count为必填",
				Validators: []validator.String{
					stringvalidator.OneOf("month", "on_demand"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_count": schema.Int32Attribute{
				Optional:    true,
				Description: "订购时长，该参数在cycle_type为month时才生效，当cycle_type=month，支持传递1、2、3、4、5、6、12、24、36",
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeMonth),
					),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
					int32validator.OneOf(1, 2, 3, 5, 6, 7, 12, 24, 36),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"auto_renew": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否自动续订，默认非自动续订，当cycle_type不等于on_demand时才可填写",
				Default:     booldefault.StaticBool(false),
				Validators: []validator.Bool{
					validator2.ConflictsWithEqualBool(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
				},
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"auto_renew_cycle_count": schema.Int32Attribute{
				Optional:    true,
				Description: "自动续订时长，当且仅当auto_renew为true时填写。支持自动续订范围：1-6月",
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("auto_renew"),
						types.BoolValue(true),
					),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("auto_renew"),
						types.BoolValue(false),
					),
					int32validator.Between(1, 6),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (c *ctyunKafkaInstance) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunKafkaInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeCreate(ctx, plan)
	if err != nil {
		return
	}
	// 创建
	masterOrderID, err := c.create(ctx, plan)
	if err != nil {
		return
	}
	plan.MasterOrderID = types.StringValue(masterOrderID)
	// 创建后检查
	id, err := c.checkAfterCreate(ctx, plan)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(id)

	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunKafkaInstance) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunKafkaInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "已经销毁") {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunKafkaInstance) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunKafkaInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunKafkaInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkBeforeUpdate(ctx, plan, state)
	if err != nil {
		return
	}
	// 更新
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunKafkaInstance) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunKafkaInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getByNameOrID(ctx, state)
	if err != nil || instance == nil {
		return
	}
	// 如果状态不是已退订状态，则执行退订
	if instance.Status != business.KafkaStatusUnsubscribed {
		// 退订
		err = c.unsubscribe(ctx, state)
		if err != nil {
			return
		}
		err = c.checkAfterUnsubscribe(ctx, state)
		if err != nil {
			return
		}
		time.Sleep(120 * time.Second)
	}
	// 销毁
	err = c.destroy(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDestroy(ctx, state)
	if err != nil {
		return
	}

	response.Diagnostics.AddWarning("删除Kakfa集群成功", "集群退订后，若立即删除子网或安全组可能会失败，需要等待底层资源释放")
}

func (c *ctyunKafkaInstance) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.vpcService = business.NewVpcService(meta)
	c.sgService = business.NewSecurityGroupService(meta)
}

// 导入命令：terraform import [配置标识].[导入配置名称] [id],[regionID]
func (c *ctyunKafkaInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunKafkaInstanceConfig
	var id, regionID string
	err = terraform_extend.Split(request.ID, &id, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.ID = types.StringValue(id)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// checkBeforeCreate 创建前检查
func (c *ctyunKafkaInstance) checkBeforeCreate(ctx context.Context, plan CtyunKafkaInstanceConfig) (err error) {
	regionID, projectID := plan.RegionID.ValueString(), plan.ProjectID.ValueString()
	vpc, subnetID, sgID := plan.VpcID.ValueString(), plan.SubnetID.ValueString(), plan.SecurityGroupID.ValueString()
	subnets, err := c.vpcService.GetVpcSubnet(ctx, vpc, regionID, projectID)
	if err != nil {
		return err
	}
	_, exist := subnets[subnetID]
	if !exist {
		err = fmt.Errorf("子网不存在")
		return err
	}
	err = c.sgService.MustExistInVpc(ctx, vpc, sgID, regionID)
	if err != nil {
		return err
	}
	err = c.checkSpecParams(ctx, plan)
	if err != nil {
		return err
	}
	return
}

// checkSpecParams 检查规格参数
func (c *ctyunKafkaInstance) checkSpecParams(ctx context.Context, plan CtyunKafkaInstanceConfig) (err error) {
	nodeNum := plan.NodeNum.ValueInt32()
	specName := plan.SpecName.ValueString()
	diskType := plan.DiskType.ValueString()

	if strings.HasSuffix(specName, "single") && nodeNum != 1 {
		return fmt.Errorf("单机版实例节点数必须为1")
	} else if strings.HasSuffix(specName, "cluster") && nodeNum < 3 {
		return fmt.Errorf("集群版实例节点数必须大于等于3")
	}
	// 组装请求体
	params := &ctgkafka.CtgkafkaProdDetailRequest{
		RegionId: plan.RegionID.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaProdDetailApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	var skuRes ctgkafka.CtgkafkaProdDetailReturnObjResponseSkuResItem
	var skuDisk ctgkafka.CtgkafkaProdDetailReturnObjResponseSkuDiskItem
	for _, s := range resp.ReturnObj.Data.Series {
		for _, p := range s.Sku {
			if p.ProdName == "集群版" && plan.NodeNum.ValueInt32() >= 3 {
				skuRes = p.ResItem
				skuDisk = p.DiskItem
				break
			} else if p.ProdName == "单机版" && plan.NodeNum.ValueInt32() == 1 {
				skuRes = p.ResItem
				skuDisk = p.DiskItem
				break
			}
		}
	}

	var specAvailable bool
	for _, r := range skuRes.ResItems {
		for _, s := range r.Spec {
			if s.SpecName == specName {
				specAvailable = true
				break
			}
		}
		if specAvailable {
			break
		}
	}
	if !specAvailable {
		return fmt.Errorf("本资源池不支持 %s", specName)
	}

	var diskAvailable bool
	for _, d := range skuDisk.ResItems {
		if d == diskType {
			diskAvailable = true
			break
		}
	}
	if !diskAvailable {
		return fmt.Errorf("本资源池不支持 %s", diskType)
	}

	return
}

// create 创建
func (c *ctyunKafkaInstance) create(ctx context.Context, plan CtyunKafkaInstanceConfig) (masterOrderID string, err error) {
	switch plan.CycleType.ValueString() {
	case business.OrderCycleTypeMonth:
		return c.createPrePayOrder(ctx, plan)
	case business.OrderCycleTypeOnDemand:
		return c.createPostPayOrder(ctx, plan)
	}
	return
}

// createPrePayOrder 创建包年包月
func (c *ctyunKafkaInstance) createPrePayOrder(ctx context.Context, plan CtyunKafkaInstanceConfig) (masterOrderID string, err error) {
	params := &ctgkafka.CtgkafkaCreateOrderRequest{
		RegionId:            plan.RegionID.ValueString(),
		ProjectId:           plan.ProjectID.ValueString(),
		CycleCnt:            plan.CycleCount.ValueInt32(),
		ClusterName:         plan.InstanceName.ValueString(),
		EngineVersion:       plan.EngineVersion.ValueString(),
		SpecName:            plan.SpecName.ValueString(),
		NodeNum:             plan.NodeNum.ValueInt32(),
		DiskType:            plan.DiskType.ValueString(),
		DiskSize:            plan.DiskSize.ValueInt32(),
		VpcId:               plan.VpcID.ValueString(),
		SubnetId:            plan.SubnetID.ValueString(),
		SecurityGroupId:     plan.SecurityGroupID.ValueString(),
		EnableIpv6:          plan.EnableIpv6.ValueBoolPointer(),
		PlainPort:           plan.PlainPort.ValueInt32(),
		SaslPort:            plan.SaslPort.ValueInt32(),
		SslPort:             plan.SslPort.ValueInt32(),
		HttpPort:            plan.HttpPort.ValueInt32(),
		RetentionHours:      plan.RetentionHours.ValueInt32(),
		AutoRenewStatus:     plan.AutoRenew.ValueBoolPointer(),
		AutoRenewCycleCount: plan.AutoRenewCycleCount.ValueInt32(),
	}

	var zoneList []string
	var str []types.String
	plan.ZoneList.ElementsAs(ctx, &str, true)
	for _, s := range str {
		zoneList = append(zoneList, s.ValueString())
	}
	params.ZoneList = zoneList

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaCreateOrderApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	masterOrderID = resp.ReturnObj.Data.NewOrderId
	return
}

// createPostPayOrder 创建按需
func (c *ctyunKafkaInstance) createPostPayOrder(ctx context.Context, plan CtyunKafkaInstanceConfig) (masterOrderID string, err error) {
	params := &ctgkafka.CtgkafkaCreatePostPayOrderRequest{
		RegionId:        plan.RegionID.ValueString(),
		ProjectId:       plan.ProjectID.ValueString(),
		ClusterName:     plan.InstanceName.ValueString(),
		EngineVersion:   plan.EngineVersion.ValueString(),
		SpecName:        plan.SpecName.ValueString(),
		NodeNum:         plan.NodeNum.ValueInt32(),
		DiskType:        plan.DiskType.ValueString(),
		DiskSize:        plan.DiskSize.ValueInt32(),
		VpcId:           plan.VpcID.ValueString(),
		SubnetId:        plan.SubnetID.ValueString(),
		SecurityGroupId: plan.SecurityGroupID.ValueString(),
		EnableIpv6:      plan.EnableIpv6.ValueBoolPointer(),
		PlainPort:       plan.PlainPort.ValueInt32(),
		SaslPort:        plan.SaslPort.ValueInt32(),
		SslPort:         plan.SslPort.ValueInt32(),
		HttpPort:        plan.HttpPort.ValueInt32(),
		RetentionHours:  plan.RetentionHours.ValueInt32(),
	}

	var zoneList []string
	var strings []types.String
	plan.ZoneList.ElementsAs(ctx, &strings, true)
	for _, s := range strings {
		zoneList = append(zoneList, s.ValueString())
	}
	params.ZoneList = zoneList

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaCreatePostPayOrderApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	masterOrderID = resp.ReturnObj.Data.NewOrderId
	return
}

// getAndMerge 从远端查询
func (c *ctyunKafkaInstance) getAndMerge(ctx context.Context, plan *CtyunKafkaInstanceConfig) (err error) {
	instance, err := c.getByNameOrID(ctx, *plan)
	if err != nil {
		return
	}
	if instance == nil {
		return fmt.Errorf("%s 已经销毁", plan.ID.ValueString())
	}

	plan.InstanceName = types.StringValue(instance.InstanceName)
	plan.Name = plan.InstanceName
	if len(instance.Version) >= 3 {
		plan.EngineVersion = types.StringValue(instance.Version[:3])
	}
	plan.SpecName = types.StringValue(instance.Specifications)
	plan.NodeNum = types.Int32Value(int32(len(instance.NodeList)))

	plan.DiskType = types.StringValue(instance.DiskType)
	plan.DiskSize = types.Int32Value(utils.StringToInt32Must(instance.Space))
	plan.VpcID = types.StringValue(instance.VpcId)
	plan.SubnetID = types.StringValue(instance.SubnetId)

	plan.EnableIpv6 = types.BoolValue(map[int32]bool{1: true, 0: false}[instance.Ipv6Enable])
	if len(instance.NodeList) > 0 {
		plan.PlainPort = types.Int32Value(utils.StringToInt32Must(instance.NodeList[0].VpcPort))
		plan.SaslPort = types.Int32Value(utils.StringToInt32Must(instance.NodeList[0].SaslPort))
		plan.SslPort = types.Int32Value(utils.StringToInt32Must(instance.NodeList[0].ListenNodePort))
		plan.HttpPort = types.Int32Value(utils.StringToInt32Must(instance.NodeList[0].HttpPort))
	}

	config, err := c.getInstanceConfig(ctx, *plan)
	if err != nil {
		return
	}
	plan.RetentionHours = types.Int32Value(utils.StringToInt32Must(config["log.retention.hours"].Value))
	if plan.ZoneList.IsNull() {
		plan.ZoneList = types.SetNull(types.StringType)
	}

	return
	// 下列字段没有地方获取
	//CycleType
	//CycleCount
	//AutoRenew
	//AutoRenewCycleCount
	//SecurityGroupID
}

func (c *ctyunKafkaInstance) checkBeforeUpdate(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	err = c.checkSpecParams(ctx, plan)
	if err != nil {
		return
	}
	if state.NodeNum.ValueInt32() == 1 && !plan.NodeNum.Equal(state.NodeNum) {
		return fmt.Errorf("单机版实例不能进行节点扩缩容操作")
	}
	if strings.HasSuffix(plan.SpecName.ValueString(), "single") && strings.HasSuffix(state.SpecName.ValueString(), "cluster") ||
		strings.HasSuffix(plan.SpecName.ValueString(), "cluster") && strings.HasSuffix(state.SpecName.ValueString(), "single") {
		return fmt.Errorf("不支持单机版和集群版互相变更")
	}
	instance, err := c.getByNameOrID(ctx, state)
	if err != nil {
		return
	}
	if instance.Status != business.KafkaStatusRunning {
		return fmt.Errorf("请在实例处于运行中状态时再进行更新操作")
	}

	return nil
}

// update 更新
func (c *ctyunKafkaInstance) update(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	err = c.updateRetentionHours(ctx, plan, state)
	if err != nil {
		return
	}
	err = c.updateName(ctx, plan, state)
	if err != nil {
		return
	}
	err = c.updateDiskSize(ctx, plan, state)
	if err != nil {
		return
	}
	err = c.updateNodeNum(ctx, plan, state)
	if err != nil {
		return
	}
	err = c.updateSpec(ctx, plan, state)
	if err != nil {
		return
	}
	return
}

// updateDiskSize 更新磁盘大小
func (c *ctyunKafkaInstance) updateDiskSize(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	if plan.DiskSize.Equal(state.DiskSize) {
		return
	}
	if plan.DiskSize.ValueInt32() > state.DiskSize.ValueInt32() {
		err = c.diskExtend(ctx, plan, state)
	} else {
		err = fmt.Errorf("目前不支持磁盘缩容")
		//err = c.diskShrink(ctx, plan, state)
	}
	if err != nil {
		return
	}
	return c.checkAfterUpdateDiskSize(ctx, plan, state)
}

//// diskExtend 磁盘缩容
//func (c *ctyunKafkaInstance) diskShrink(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
//	params := &ctgkafka.CtgkafkaDiskShrinkRequest{
//		RegionId:   state.RegionID.ValueString(),
//		ProdInstId: state.ID.ValueString(),
//		DiskSize:   plan.DiskSize.String(),
//	}
//	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaDiskShrinkApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//	return
//}

// diskExtend 磁盘扩容
func (c *ctyunKafkaInstance) diskExtend(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	autoPay := true
	params := &ctgkafka.CtgkafkaDiskExtendRequest{
		RegionId:       state.RegionID.ValueString(),
		ProdInstId:     state.ID.ValueString(),
		DiskExtendSize: plan.DiskSize.String(),
		AutoPay:        &autoPay,
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaDiskExtendApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// checkAfterUpdateDiskSize 检查磁盘大小是否变更成功
func (c *ctyunKafkaInstance) checkAfterUpdateDiskSize(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	var executeSuccessFlag bool
	var successCnt int
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *ctgkafka.CtgkafkaInstQueryReturnObjDataResponse
			instance, err = c.getByNameOrID(ctx, state)
			if err != nil {
				return false
			}
			if utils.StringToInt32Must(instance.Space) != plan.DiskSize.ValueInt32() || instance.Status != business.KafkaStatusRunning {
				return true
			}
			successCnt++
			if successCnt < 3 {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("磁盘变配时间过长")
	}
	return
}

// updateNodeNum 更新节点数量
func (c *ctyunKafkaInstance) updateNodeNum(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	if plan.NodeNum.Equal(state.NodeNum) {
		return
	}
	if plan.NodeNum.ValueInt32() > state.NodeNum.ValueInt32() {
		err = c.nodeExtend(ctx, plan, state)
	} else {
		err = fmt.Errorf("目前不支持节点缩容")
		//err = c.nodeShrink(ctx, plan, state)
	}
	if err != nil {
		return
	}
	return c.checkAfterUpdateNodeNum(ctx, plan, state)
}

// nodeShrink 节点缩容
//func (c *ctyunKafkaInstance) nodeShrink(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
//	params := &ctgkafka.CtgkafkaNodeShrinkRequest{
//		RegionId:   state.RegionID.ValueString(),
//		ProdInstId: state.ID.ValueString(),
//		NodeNum:    plan.DiskSize.String(),
//	}
//	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaNodeShrinkApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCodeString {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//	return
//}

// nodeExtend 节点扩容
func (c *ctyunKafkaInstance) nodeExtend(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	autoPay := true
	params := &ctgkafka.CtgkafkaNodeExtendRequest{
		RegionId:      state.RegionID.ValueString(),
		ProdInstId:    state.ID.ValueString(),
		ExtendNodeNum: plan.NodeNum.ValueInt32(),
		AutoPay:       &autoPay,
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaNodeExtendApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// checkAfterUpdateNodeNum 检查节点数量是否变更成功
func (c *ctyunKafkaInstance) checkAfterUpdateNodeNum(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	var executeSuccessFlag bool
	var successCnt int
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *ctgkafka.CtgkafkaInstQueryReturnObjDataResponse
			instance, err = c.getByNameOrID(ctx, state)
			if err != nil {
				return false
			}
			if len(instance.NodeList) != int(plan.NodeNum.ValueInt32()) || instance.Status != business.KafkaStatusRunning {
				return true
			}
			successCnt++
			if successCnt < 3 {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("节点数量变配时间过长")
	}
	return
}

// updateSpec 更新规格
func (c *ctyunKafkaInstance) updateSpec(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	if plan.SpecName.Equal(state.SpecName) {
		return
	}
	planU, err := c.parseSpec(plan.SpecName.ValueString())
	if err != nil {
		return
	}
	stateU, err := c.parseSpec(state.SpecName.ValueString())
	if err != nil {
		return
	}
	if planU > stateU {
		err = c.specExtend(ctx, plan, state)
	} else {
		err = c.specShrink(ctx, plan, state)
	}
	if err != nil {
		return
	}
	return c.checkAfterUpdateSpec(ctx, plan, state)
}

// parseSpec 从规格名称解析cpu
func (c *ctyunKafkaInstance) parseSpec(s string) (u int, err error) {
	re := regexp.MustCompile(`(\d+)u(\d+)g`)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 3 {
		err = fmt.Errorf("invalid format: %s", s)
		return
	}

	if _, err = fmt.Sscanf(matches[1], "%d", &u); err != nil {
		return
	}
	return
}

// specShrink 规格缩容
func (c *ctyunKafkaInstance) specShrink(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	params := &ctgkafka.CtgkafkaSpecShrinkRequest{
		RegionId:   state.RegionID.ValueString(),
		ProdInstId: state.ID.ValueString(),
		SpecName:   plan.SpecName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaSpecShrinkApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// specExtend 规格扩容
func (c *ctyunKafkaInstance) specExtend(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	autoPay := true
	params := &ctgkafka.CtgkafkaSpecExtendRequest{
		RegionId:   state.RegionID.ValueString(),
		ProdInstId: state.ID.ValueString(),
		SpecName:   plan.SpecName.ValueString(),
		AutoPay:    &autoPay,
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaSpecExtendApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// checkAfterUpdateSpec 检查规格是否变更成功
func (c *ctyunKafkaInstance) checkAfterUpdateSpec(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	var executeSuccessFlag bool
	var successCnt int
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *ctgkafka.CtgkafkaInstQueryReturnObjDataResponse
			instance, err = c.getByNameOrID(ctx, state)
			if err != nil {
				return false
			}
			if instance.Specifications != plan.SpecName.ValueString() || instance.Status != business.KafkaStatusRunning {
				return true
			}
			successCnt++
			if successCnt < 3 {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("规格变配时间过长")
	}
	return
}

// updateName 更新实例名称
func (c *ctyunKafkaInstance) updateName(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	if plan.InstanceName.Equal(state.InstanceName) {
		return
	}
	params := &ctgkafka.CtgkafkaInstancesModifyNameV3Request{
		RegionId:     state.RegionID.ValueString(),
		ProdInstId:   state.ID.ValueString(),
		InstanceName: plan.InstanceName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaInstancesModifyNameV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		return fmt.Errorf("API return error. Message: %s", resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	} else if resp.ReturnObj.Data != "modify success" {
		return fmt.Errorf("API return error. Data: %s", resp.ReturnObj.Data)
	}
	return
}

// updateRetentionHours 更新实例消息保留时长
func (c *ctyunKafkaInstance) updateRetentionHours(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	if plan.RetentionHours.Equal(state.RetentionHours) {
		return
	}
	params := &ctgkafka.CtgkafkaUpdateInstanceConfigRequest{
		RegionId:   state.RegionID.ValueString(),
		ProdInstId: state.ID.ValueString(),
		StaticConfigs: []*ctgkafka.CtgkafkaUpdateInstanceConfigStaticConfigsRequest{
			{Name: "log.retention.hours", Value: plan.RetentionHours.String()},
		},
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaUpdateInstanceConfigApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		return fmt.Errorf("API return error. Message: %s", resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	} else if resp.ReturnObj.Data != "modify success" {
		return fmt.Errorf("API return error. Data: %s", resp.ReturnObj.Data)
	}
	// 更新后需要重启
	return c.reboot(ctx, plan, state)
}

// reboot 重启实例
func (c *ctyunKafkaInstance) reboot(ctx context.Context, plan, state CtyunKafkaInstanceConfig) (err error) {
	params := &ctgkafka.CtgkafkaInstancesRestartV3Request{
		RegionId:   state.RegionID.ValueString(),
		ProdInstId: state.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaInstancesRestartV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 等待重启完成
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 60)
	retryer.Start(
		func(currentTime int) bool {
			var instance *ctgkafka.CtgkafkaInstQueryReturnObjDataResponse
			instance, err = c.getByNameOrID(ctx, state)
			if err != nil {
				return false
			}
			if instance.Status != business.KafkaStatusRunning {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("重启时间过长")
	}
	return
}

// unsubscribe 退订
func (c *ctyunKafkaInstance) unsubscribe(ctx context.Context, plan CtyunKafkaInstanceConfig) (err error) {
	params := &ctgkafka.CtgkafkaUnsubscribeInstV3Request{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaUnsubscribeInstV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// unsubscribe 退订后检查
func (c *ctyunKafkaInstance) checkAfterUnsubscribe(ctx context.Context, plan CtyunKafkaInstanceConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *ctgkafka.CtgkafkaInstQueryReturnObjDataResponse
			instance, err = c.getByNameOrID(ctx, plan)
			if err != nil {
				return false
			}
			if instance.Status != business.KafkaStatusUnsubscribed {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("退订时间过长")
	}
	return
}

// destroy 销毁
func (c *ctyunKafkaInstance) destroy(ctx context.Context, plan CtyunKafkaInstanceConfig) (err error) {
	params := &ctgkafka.CtgkafkaInstanceDeleteV3Request{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaInstanceDeleteV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// unsubscribe 销毁后检查
func (c *ctyunKafkaInstance) checkAfterDestroy(ctx context.Context, plan CtyunKafkaInstanceConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *ctgkafka.CtgkafkaInstQueryReturnObjDataResponse
			instance, err = c.getByNameOrID(ctx, plan)
			if err != nil {
				return false
			}
			if instance != nil {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("销毁时间过长")
	}
	return
}

// checkAfterCreate 创建后检查
func (c *ctyunKafkaInstance) checkAfterCreate(ctx context.Context, plan CtyunKafkaInstanceConfig) (id string, err error) {
	var executeSuccessFlag bool
	var successCnt int
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *ctgkafka.CtgkafkaInstQueryReturnObjDataResponse
			instance, err = c.getByNameOrID(ctx, plan)
			if err != nil {
				return false
			}
			if instance == nil || instance.Status != business.KafkaStatusRunning || instance.ProdInstId == "" {
				return true
			}
			// 等待订单完成
			successCnt++
			if successCnt < 3 {
				return true
			}
			id = instance.ProdInstId
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("创建时间过长")
	}
	return
}

// getByNameOrID 根据ID或名称查询集群
func (c *ctyunKafkaInstance) getByNameOrID(ctx context.Context, plan CtyunKafkaInstanceConfig) (instance *ctgkafka.CtgkafkaInstQueryReturnObjDataResponse, err error) {
	params := &ctgkafka.CtgkafkaInstQueryRequest{
		RegionId:       plan.RegionID.ValueString(),
		OuterProjectId: plan.ProjectID.ValueString(),
	}

	if plan.ID.ValueString() != "" {
		params.ProdInstId = plan.ID.ValueString()
	} else if plan.InstanceName.ValueString() != "" {
		params.Name = plan.InstanceName.ValueString()
		e := true
		params.ExactMatchName = &e
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaInstQueryApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	if len(resp.ReturnObj.Data) > 0 {
		instance = resp.ReturnObj.Data[0]
		if instance == nil {
			err = common.InvalidReturnObjResultsError
		}
	}
	return
}

// getInstanceConfig 获取实例配置
func (c *ctyunKafkaInstance) getInstanceConfig(ctx context.Context, plan CtyunKafkaInstanceConfig) (attr map[string]*ctgkafka.CtgkafkaGetInstanceConfigReturnObjDataResponse, err error) {
	params := &ctgkafka.CtgkafkaGetInstanceConfigRequest{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaGetInstanceConfigApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	attr = map[string]*ctgkafka.CtgkafkaGetInstanceConfigReturnObjDataResponse{}
	for _, d := range resp.ReturnObj.Data {
		attr[d.Name] = d
	}
	return
}
