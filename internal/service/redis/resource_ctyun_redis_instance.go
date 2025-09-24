package redis

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &ctyunRedisInstance{}
	_ resource.ResourceWithConfigure   = &ctyunRedisInstance{}
	_ resource.ResourceWithImportState = &ctyunRedisInstance{}
)

type ctyunRedisInstance struct {
	meta       *common.CtyunMetadata
	vpcService *business.VpcService
	sgService  *business.SecurityGroupService
}

func NewCtyunRedisInstance() resource.Resource {
	return &ctyunRedisInstance{}
}

func (c *ctyunRedisInstance) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_instance"
}

type CtyunRedisInstanceConfig struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	MasterOrderID       types.String `tfsdk:"master_order_id"`
	RegionID            types.String `tfsdk:"region_id"`
	ProjectID           types.String `tfsdk:"project_id"`
	CycleCount          types.Int32  `tfsdk:"cycle_count"`
	CycleType           types.String `tfsdk:"cycle_type"`        // on_demand 和 month
	AzName              types.String `tfsdk:"az_name"`           /*  主可用区名称，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解可用区<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=17764&isNormal=1&vid=270">查询可用区信息</a> name字段  */
	SecondaryAzName     types.String `tfsdk:"secondary_az_name"` /*  备可用区名称(双/多副本建议填写)<br>默认与主可用区相同  */
	EngineVersion       types.String `tfsdk:"engine_version"`    /*  Redis引擎版本<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表SeriesInfo中的engineTypeItems(引擎版本可选值)  */
	Version             types.String `tfsdk:"version"`           /*  版本类型。<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表SeriesInfo中的version值<br>可选值：<li>BASIC：基础版<li>PLUS：增强版<li>Classic：经典版  */
	Edition             types.String `tfsdk:"edition"`           /*  实例类型<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a>  使用表SeriesInfo中的seriesCode值  */
	HostType            types.String `tfsdk:"host_type"`         /*  主机类型<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表resItems中resType==ecs的items(主机类型可选值)  */
	DataDiskType        types.String `tfsdk:"data_disk_type"`
	ShardMemSize        types.Int32  `tfsdk:"shard_mem_size"` /*  单分片内存(GB)<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表SeriesInfo中shardMemSizeItems(单分片内存可选值)，若shardMemSizeItems为空则无需填写  */
	ShardCount          types.Int32  `tfsdk:"shard_count"`
	CopiesCount         types.Int32  `tfsdk:"copies_count"`           /*  副本数量，取值范围2-6。<li>OriginalMultipleReadLvs：必填</li><li>StandardDual/DirectCluster/ClusterOriginalProxy：选填</li><li>其他实例类型：无需填写</li>  */
	InstanceName        types.String `tfsdk:"instance_name"`          /*  实例名称<li>字母开头</li><li>可包含字母/数字/中划线</li><li>长度1-39<li>实例名称不可重复</li>  */
	VpcID               types.String `tfsdk:"vpc_id"`                 /*  虚拟私有云ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028310">产品定义-虚拟私有云</a>来了解虚拟私有云<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4814&data=94&vid=88">查询VPC列表</a> vpcID字段。<br><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4811&data=94&vid=88">创建VPC</a>  */
	SubnetID            types.String `tfsdk:"subnet_id"`              /*  子网ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10098380">基本概念</a>来查找子网的相关定义<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8659&data=94&vid=88">查询子网列表</a> subnetID字段。  */
	SecurityGroupID     types.String `tfsdk:"security_group_id"`      /*  安全组ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关信息<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/searchCtapi/ctApiDebug?product=18&api=4817&vid=88">查询用户安全组列表</a> id字段。  */
	Password            types.String `tfsdk:"password"`               /*  实例密码<li>长度8-26字符</li><li>必须同时包含大写字母、小写字母、数字、英文格式特殊符号(@%^*_+!$-=.) 中的三种类型</li><li>不能有空格</li>  */
	AutoRenew           types.Bool   `tfsdk:"auto_renew"`             /*  自动续费开关<li>true：开启</li><li>false：关闭(默认)</li>  */
	AutoRenewCycleCount types.Int32  `tfsdk:"auto_renew_cycle_count"` /*  自动续费周期(月)<br>autoRenew=true时必填，可选：1-6,12,24,36  */
	MaintenanceTime     types.String `tfsdk:"maintenance_time"`
	ProtectionStatus    types.Bool   `tfsdk:"protection_status"`
}

func (c *ctyunRedisInstance) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10029420/10029727**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed:    true,
				Description: "ID",
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
				},
				Validators: []validator.String{
					validator2.Project(),
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
					int32validator.OneOf(1, 2, 3, 4, 5, 6, 12, 24, 36),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "主可用区，如果不填则默认使用provider ctyun中的az_name或环境变量中的CTYUN_AZ_NAME",
				Default:     defaults.AcquireFromGlobalString(common.ExtraAzName, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"secondary_az_name": schema.StringAttribute{
				Optional:    true,
				Description: "备可用区",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "版本类型，SeriesInfo中的version值，支持BASIC和PLUS，默认BASIC",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.RedisVersionBasic, business.RedisVersionPlus),
				},
				Default: stringdefault.StaticString(business.RedisVersionBasic),
			},
			"edition": schema.StringAttribute{
				Required:    true,
				Description: "实例类型，SeriesInfo中的seriesCode值，可参考<a href=\"https://www.ctyun.cn/document/10029420/11030280\">产品规格说明</a>",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.RedisEdition...),
				},
			},
			"engine_version": schema.StringAttribute{
				Required:    true,
				Description: "Redis引擎版本，SeriesInfo中的engineTypeItems(引擎版本可选值)，当version取值为BASIC时，版本号取值：5.0，6.0，7.0，当version取值为PLUS，版本号取值：6.0，7.0，支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.RedisEngineVersion...),
				},
			},
			"data_disk_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "磁盘类型，支持SAS和SSD，默认SAS",
				Validators: []validator.String{
					stringvalidator.OneOf(business.RedisDiskTypeSas, business.RedisDiskTypeSsd),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString(business.RedisDiskTypeSas),
			},
			"host_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "主机类型，默认S，X86取值：S：通用型、C：计算增强型、M：内存型、HS：海光通用型、HC：海光计算增强型，ARM取值：KS：鲲鹏通用型、KC：鲲鹏计算增强型",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.RedisHostType...),
				},
				Default: stringdefault.StaticString(business.RedisHostTypeS),
			},
			"shard_mem_size": schema.Int32Attribute{
				Required:    true,
				Description: "分片规格，当version取值为BASIC，取值：1、2、4、8、16、32、64，当version取值为PLUS时，取值：8、16、32、64",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					int32validator.OneOf(1, 2, 4, 8, 16, 32, 64),
				},
			},
			"shard_count": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Description: "分片数量，当edition取值为DirectClusterSingle时: 3~256。当edition取值为DirectCluster时: 3~256。当edition取值为ClusterOriginalProxy时: 3~64。当edition取其他值时不填。",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("edition"),
						types.StringValue(business.RedisEditionDirectClusterSingle),
						types.StringValue(business.RedisEditionDirectCluster),
						types.StringValue(business.RedisEditionClusterOriginalProxy),
					),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("edition"),
						types.StringValue(business.RedisEditionStandardSingle),
						types.StringValue(business.RedisEditionStandardDual),
						types.StringValue(business.RedisEditionOriginalMultipleReadLvs),
					),
				},
			},
			"copies_count": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Description: "副本数量，当edition取值为OriginalMultipleReadLvs/StandardDual/DirectCluster/ClusterOriginalProxy时必填（取值范围2-6），当edition取其他值时不填。",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("edition"),
						types.StringValue(business.RedisEditionOriginalMultipleReadLvs),
						types.StringValue(business.RedisEditionStandardDual),
						types.StringValue(business.RedisEditionDirectCluster),
						types.StringValue(business.RedisEditionClusterOriginalProxy),
					),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("edition"),
						types.StringValue(business.RedisEditionStandardSingle),
						types.StringValue(business.RedisEditionDirectClusterSingle),
					),
					int32validator.Between(2, 6),
				},
			},
			"instance_name": schema.StringAttribute{
				Required:    true,
				Description: "实例名称，大小写字母开头。只能包含大小写字母、数字及分隔符(-)。大小写字母或数字结尾。长度4~40个字符。实例名称不可重复。",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(4, 40),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$"), "不满足实例名称要求"),
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
			"password": schema.StringAttribute{
				Required:    true,
				Description: "实例密码。长度8-26字符。必须同时包含大写字母、小写字母、数字、英文格式特殊符号(@%^*_+!$-=.)中的三种类型。不能有空格。支持更新",
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(8, 26),
					validator2.RedisPassword(),
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
					boolplanmodifier.RequiresReplace(),
				},
			},
			"auto_renew_cycle_count": schema.Int32Attribute{
				Optional:    true,
				Description: "自动续订时长，单位月，支持1, 2, 3, 5, 6, 7, 12, 24, 36",
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("auto_renew"),
						types.BoolValue(true),
					),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("auto_renew"),
						types.BoolValue(false),
					),
					int32validator.OneOf(1, 2, 3, 5, 6, 7, 12, 24, 36),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"maintenance_time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "实例维护时间窗口，总时长必须为2小时，默认：00:00-02:00，支持更新",
				Default:     stringdefault.StaticString("00:00-02:00"),
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^([0-1][0-9]|2[0-3]):[0-5][0-9]-([0-1][0-9]|2[0-3]):[0-5][0-9]$"), "时间格式错误"),
				},
			},
			"protection_status": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "退订保护开关，默认为不保护，支持更新",
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (c *ctyunRedisInstance) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunRedisInstanceConfig
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
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)

	err = c.updateAttr(ctx, plan)
	if err != nil {
		return
	}
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunRedisInstance) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRedisInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "can't find") {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRedisInstance) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunRedisInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunRedisInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}
	state.Password = plan.Password
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRedisInstance) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRedisInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	instance, err := c.getByName(ctx, state)
	if err != nil || instance == nil {
		return
	}
	// 如果状态不是已退订状态，则执行退订
	if instance.Status != business.RedisStatusUnsubscribed {
		// 退订
		err = c.unsubscribe(ctx, state)
		if err != nil {
			return
		}
		err = c.checkAfterUnsubscribe(ctx, state)
		if err != nil {
			return
		}
		time.Sleep(60 * time.Second)
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

	response.Diagnostics.AddWarning("删除Redis集群成功", "集群退订后，若立即删除子网或安全组可能会失败，需要等待底层资源释放")
}

func (c *ctyunRedisInstance) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.vpcService = business.NewVpcService(meta)
	c.sgService = business.NewSecurityGroupService(meta)
}

// 导入命令：terraform import [配置标识].[导入配置名称] [id],[regionID]
func (c *ctyunRedisInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunRedisInstanceConfig
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
func (c *ctyunRedisInstance) checkBeforeCreate(ctx context.Context, plan CtyunRedisInstanceConfig) (err error) {
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
func (c *ctyunRedisInstance) checkSpecParams(ctx context.Context, plan CtyunRedisInstanceConfig) (err error) {
	copiesCount := plan.ShardCount.ValueInt32()
	shardCount := plan.ShardCount.ValueInt32()

	switch plan.Edition.ValueString() {
	case business.RedisEditionDirectClusterSingle, business.RedisEditionDirectCluster:
		if shardCount < 3 || shardCount > 256 {
			return fmt.Errorf("edition为DirectClusterSingle或DirectCluster，shard_count需要在3-256")
		}
	case "ClusterOriginalP":
		if shardCount < 3 || shardCount > 64 {
			return fmt.Errorf("edition为ClusterOriginalP，shard_count需要在3-64")
		}
	}
	if shardCount == 0 {
		shardCount = 1
	}
	if copiesCount == 0 {
		copiesCount = 1
	}

	// 组装请求体
	params := &dcs2.Dcs2DescribeAvailableResourceRequest{
		RegionId: plan.RegionID.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeAvailableResourceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	var available bool
	for _, spec := range resp.ReturnObj.SeriesInfoList {
		if spec.Version == plan.Version.ValueString() && spec.SeriesCode == plan.Edition.ValueString() {
			var engineOK bool
			for _, engine := range spec.EngineTypeItems {
				if engine == plan.EngineVersion.ValueString() {
					engineOK = true
					break
				}
			}
			if !engineOK {
				return fmt.Errorf("engine_version不在合法取值范围内")
			}

			var memSizeOK bool
			var memSize []string
			if spec.ShardMemSizeItems != nil {
				memSize = spec.ShardMemSizeItems
			} else {
				memSize = spec.MemSizeItems
			}
			for _, ms := range memSize {
				s := fmt.Sprintf("%d", plan.ShardMemSize.ValueInt32())
				if s == ms {
					memSizeOK = true
					break
				}
			}
			if !memSizeOK {
				return fmt.Errorf("shard_mem_size不在合法取值范围内")
			}

			var dataDiskTypeOK, hostType bool
			for _, items := range spec.ResItems {
				switch items.ResType {
				case "ebs":
					for _, item := range items.Items {
						if item == plan.DataDiskType.ValueString() {
							dataDiskTypeOK = true
							break
						}
					}
				case "hostType":
					for _, item := range items.Items {
						if item == plan.HostType.ValueString() {
							hostType = true
							break
						}
					}
				}
			}
			if !dataDiskTypeOK {
				return fmt.Errorf("请指定正确的data_disk_type")
			}
			if !hostType {
				return fmt.Errorf("请指定正确的host_type")
			}
			available = true
			break
		}
	}
	if !available {
		err = fmt.Errorf("未找到对应规格，请确认version和edition")
	}
	return
}

// create 创建
func (c *ctyunRedisInstance) create(ctx context.Context, plan CtyunRedisInstanceConfig) (masterOrderID string, err error) {
	autoPay := true
	params := &dcs2.Dcs2CreateInstanceRequest{
		RegionId:          plan.RegionID.ValueString(),
		ProjectID:         plan.ProjectID.ValueString(),
		ZoneName:          plan.AzName.ValueString(),
		SecondaryZoneName: plan.SecondaryAzName.ValueString(),
		EngineVersion:     plan.EngineVersion.ValueString(),
		Version:           plan.Version.ValueString(),
		Edition:           plan.Edition.ValueString(),
		HostType:          plan.HostType.ValueString(),
		DataDiskType:      plan.DataDiskType.ValueString(),
		ShardCount:        plan.ShardCount.ValueInt32(),
		CopiesCount:       plan.CopiesCount.ValueInt32(),
		InstanceName:      plan.InstanceName.ValueString(),
		VpcId:             plan.VpcID.ValueString(),
		SubnetId:          plan.SubnetID.ValueString(),
		Secgroups:         plan.SecurityGroupID.ValueString(),
		Password:          plan.Password.ValueString(),
		AutoPay:           &autoPay,
	}

	if plan.CycleType.ValueString() == business.OnDemandCycleType {
		params.ChargeType = "PostPaid"
	} else {
		params.ChargeType = "PrePaid"
		params.Period = plan.CycleCount.ValueInt32()
	}
	if plan.AutoRenew.ValueBool() {
		params.AutoRenew = plan.AutoRenew.ValueBoolPointer()
		params.AutoRenewPeriod = fmt.Sprintf("%d", plan.AutoRenewCycleCount.ValueInt32())
	}
	if plan.ShardMemSize.ValueInt32() > 0 {
		params.ShardMemSize = fmt.Sprintf("%d", plan.ShardMemSize.ValueInt32())
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2CreateInstanceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	masterOrderID = resp.ReturnObj.NewOrderId
	return
}

// getAndMerge 从远端查询
func (c *ctyunRedisInstance) getAndMerge(ctx context.Context, plan *CtyunRedisInstanceConfig) (err error) {
	id, regionID := plan.ID.ValueString(), plan.RegionID.ValueString()
	params := &dcs2.Dcs2DescribeInstancesOverviewRequest{
		RegionId:   regionID,
		ProdInstId: id,
	}
	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeInstancesOverviewApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	instance := resp.ReturnObj.UserInfo

	if len(instance.AzList) > 0 {
		plan.AzName = types.StringValue(instance.AzList[0].AzEngName)
	}
	if len(instance.AzList) > 1 {
		plan.SecondaryAzName = types.StringValue(instance.AzList[1].AzEngName)
	}

	plan.MaintenanceTime = types.StringValue(instance.MaintenanceTime)
	plan.ProtectionStatus = utils.SecBoolValue(instance.ProtectionStatus)
	plan.EngineVersion = types.StringValue(instance.EngineVersion)
	plan.DataDiskType = types.StringValue(instance.DataDiskType)
	shardMemSize, _ := strconv.Atoi(instance.ShardMemSize)
	if shardMemSize == 0 {
		shardMemSize, _ = strconv.Atoi(instance.Capacity)
	}
	plan.ShardMemSize = types.Int32Value(int32(shardMemSize))
	shardCount, _ := strconv.Atoi(instance.ShardCount)
	plan.ShardCount = types.Int32Value(int32(shardCount))
	copiesCount, _ := strconv.Atoi(instance.CopiesCount)
	plan.CopiesCount = types.Int32Value(int32(copiesCount))
	plan.InstanceName = types.StringValue(instance.InstanceName)
	plan.Name = plan.InstanceName
	for _, p := range instance.PaasInstAttrs {
		switch p.AttrKey {
		case "vpcUuid":
			plan.VpcID = types.StringValue(p.AttrVal)
		case "subnetUuid":
			plan.SubnetID = types.StringValue(p.AttrVal)
		case "securityGroupUuid":
			plan.SecurityGroupID = types.StringValue(p.AttrVal)
		case "autoRenewStatus":
			plan.AutoRenew = types.BoolValue(map[string]bool{"false": false, "true": true}[p.AttrVal])
		case "projectId":
			plan.ProjectID = types.StringValue(p.AttrVal)
		case "autoRenewPeriod":
		}
	}

	return
}

// update 更新
func (c *ctyunRedisInstance) update(ctx context.Context, plan, state CtyunRedisInstanceConfig) (err error) {
	plan.ID = state.ID
	if !plan.MaintenanceTime.Equal(state.MaintenanceTime) || !plan.ProtectionStatus.Equal(state.ProtectionStatus) {
		err = c.updateAttr(ctx, plan)
		if err != nil {
			return
		}
	}
	err = c.updatePassword(ctx, plan, state)
	if err != nil {
		return
	}
	err = c.updateEngineVersion(ctx, plan, state)
	if err != nil {
		return
	}
	return
}

// unsubscribe 退订
func (c *ctyunRedisInstance) unsubscribe(ctx context.Context, plan CtyunRedisInstanceConfig) (err error) {
	id, regionID := plan.ID.ValueString(), plan.RegionID.ValueString()
	params := &dcs2.Dcs2DeleteInstanceRequest{
		RegionId:   regionID,
		ProdInstId: id,
	}
	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DeleteInstanceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// destroy 销毁
func (c *ctyunRedisInstance) destroy(ctx context.Context, plan CtyunRedisInstanceConfig) (err error) {
	id, regionID := plan.ID.ValueString(), plan.RegionID.ValueString()
	params := &dcs2.Dcs2DestroyInstanceRequest{
		RegionId:   regionID,
		ProdInstId: id,
	}
	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DestroyInstanceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// getByName 根据名称查询集群
func (c *ctyunRedisInstance) getByName(ctx context.Context, plan CtyunRedisInstanceConfig) (instance *dcs2.Dcs2DescribeInstancesReturnObjRowsResponse, err error) {
	params := &dcs2.Dcs2DescribeInstancesRequest{
		RegionId:     plan.RegionID.ValueString(),
		ProjectId:    plan.ProjectID.ValueString(),
		InstanceName: plan.InstanceName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeInstancesApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	if len(resp.ReturnObj.Rows) > 0 {
		instance = resp.ReturnObj.Rows[0]
		if instance == nil {
			err = common.InvalidReturnObjResultsError
		}
	}
	return
}

// getCycleByName 根据名称查询回收站
func (c *ctyunRedisInstance) getCycleByName(ctx context.Context, plan CtyunRedisInstanceConfig) (instance *dcs2.Dcs2DescribeCycleBinInstancesReturnObjRowsResponse, err error) {
	params := &dcs2.Dcs2DescribeCycleBinInstancesRequest{
		RegionId:     plan.RegionID.ValueString(),
		ProjectId:    plan.ProjectID.ValueString(),
		InstanceName: plan.InstanceName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeCycleBinInstancesApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	if len(resp.ReturnObj.Rows) > 0 {
		instance = resp.ReturnObj.Rows[0]
		if instance == nil {
			err = common.InvalidReturnObjResultsError
		}
	}
	return
}

// checkAfterCreate 创建后检查
func (c *ctyunRedisInstance) checkAfterCreate(ctx context.Context, plan CtyunRedisInstanceConfig) (id string, err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *dcs2.Dcs2DescribeInstancesReturnObjRowsResponse
			instance, err = c.getByName(ctx, plan)
			if err != nil {
				return false
			}
			// 确认失败了
			if instance != nil && instance.Status == business.RedisStatusActivationFailed {
				err = fmt.Errorf("%s 开通失败", plan.Name.ValueString())
				return false
			}
			if instance == nil || instance.Status != business.RedisStatusRunning || instance.ProdInstId == "" {
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

// checkAfterUnsubscribe 退订后检查
func (c *ctyunRedisInstance) checkAfterUnsubscribe(ctx context.Context, plan CtyunRedisInstanceConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *dcs2.Dcs2DescribeInstancesReturnObjRowsResponse
			instance, err = c.getByName(ctx, plan)
			if err != nil {
				return false
			}
			if instance != nil && instance.Status != business.RedisStatusUnsubscribed {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("删除时间过长")
	}
	return
}

// checkAfterDestroy 销毁后检查
func (c *ctyunRedisInstance) checkAfterDestroy(ctx context.Context, plan CtyunRedisInstanceConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *dcs2.Dcs2DescribeCycleBinInstancesReturnObjRowsResponse
			instance, err = c.getCycleByName(ctx, plan)
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

// updateAttr 更新保护开关和维护时间
func (c *ctyunRedisInstance) updateAttr(ctx context.Context, plan CtyunRedisInstanceConfig) (err error) {
	params := &dcs2.Dcs2ModifyInstanceAttributeRequest{
		RegionId:         plan.RegionID.ValueString(),
		ProdInstId:       plan.ID.ValueString(),
		ProtectionStatus: plan.ProtectionStatus.ValueBoolPointer(),
		MaintenanceTime:  plan.MaintenanceTime.ValueString(),
	}
	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2ModifyInstanceAttributeApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	}
	return
}

// updatePassword 更新密码
func (c *ctyunRedisInstance) updatePassword(ctx context.Context, plan, state CtyunRedisInstanceConfig) (err error) {
	if plan.Password.Equal(state.Password) {
		return
	}
	params := &dcs2.Dcs2ResetInstancePasswordRequest{
		RegionId:    plan.RegionID.ValueString(),
		ProdInstId:  plan.ID.ValueString(),
		NewPassword: plan.Password.ValueString(),
	}
	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2ResetInstancePasswordApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	}
	return
}

// updateEngineVersion 升级引擎大版本
func (c *ctyunRedisInstance) updateEngineVersion(ctx context.Context, plan, state CtyunRedisInstanceConfig) (err error) {
	if plan.EngineVersion.Equal(state.EngineVersion) {
		return
	}
	if plan.EngineVersion.ValueString() < state.EngineVersion.ValueString() {
		return fmt.Errorf("仅支持升级引擎版本")
	}
	params := &dcs2.Dcs2ModifyInstanceMajorVersionRequest{
		RegionId:      plan.RegionID.ValueString(),
		ProdInstId:    plan.ID.ValueString(),
		EngineVersion: plan.EngineVersion.ValueString(),
	}
	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2ModifyInstanceMajorVersionApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	}
	err = c.checkAfterUpdateEngineVersion(ctx, plan, state)
	if err != nil {
		return
	}
	return
}

// checkAfterUpdateEngineVersion 检查引擎版本升级是否成功
func (c *ctyunRedisInstance) checkAfterUpdateEngineVersion(ctx context.Context, plan, state CtyunRedisInstanceConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 60)
	retryer.Start(
		func(currentTime int) bool {
			var instance *dcs2.Dcs2DescribeInstancesReturnObjRowsResponse
			instance, err = c.getByName(ctx, plan)
			if err != nil {
				return false
			}
			if instance == nil {
				err = fmt.Errorf("%s 该实例已经不存在", plan.ID.ValueString())
				return false
			}
			if instance.EngineVersion != plan.EngineVersion.ValueString() {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("引擎版本升级时间过长")
	}
	return
}
