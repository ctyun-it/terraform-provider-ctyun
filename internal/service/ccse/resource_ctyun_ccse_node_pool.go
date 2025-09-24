package ccse

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ccse2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ccse"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &ctyunCcseNodePool{}
	_ resource.ResourceWithConfigure   = &ctyunCcseNodePool{}
	_ resource.ResourceWithImportState = &ctyunCcseNodePool{}
)

type ctyunCcseNodePool struct {
	meta *common.CtyunMetadata
}

func NewCtyunCcseNodePool() resource.Resource {
	return &ctyunCcseNodePool{}
}

func (c *ctyunCcseNodePool) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ccse_node_pool"
}

type CtyunCcseNodePoolConfig struct {
	ID                       types.String              `tfsdk:"id"`
	ClusterID                types.String              `tfsdk:"cluster_id"`
	RegionID                 types.String              `tfsdk:"region_id"`
	NodePoolName             types.String              `tfsdk:"name"`
	CycleCount               types.Int64               `tfsdk:"cycle_count"`
	CycleType                types.String              `tfsdk:"cycle_type"`
	AutoRenew                types.Bool                `tfsdk:"auto_renew"`
	VisibilityPostHostScript types.String              `tfsdk:"visibility_post_host_script"`
	VisibilityHostScript     types.String              `tfsdk:"visibility_host_script"`
	InstanceType             types.String              `tfsdk:"instance_type"`
	MirrorID                 types.String              `tfsdk:"mirror_id"`
	MirrorName               types.String              `tfsdk:"mirror_name"`
	MirrorType               types.Int32               `tfsdk:"mirror_type"`
	Password                 types.String              `tfsdk:"password"`
	KeyPairName              types.String              `tfsdk:"key_pair_name"`
	UseAffinityGroup         types.Bool                `tfsdk:"use_affinity_group"`
	AffinityGroupID          types.String              `tfsdk:"affinity_group_id"`
	ItemDefName              types.String              `tfsdk:"item_def_name"`
	SysDisk                  *CtyunCcseNodePoolDisk    `tfsdk:"sys_disk"`
	DataDisks                []CtyunCcseNodePoolDisk   `tfsdk:"data_disks"`
	MaxPodNum                types.Int32               `tfsdk:"max_pod_num"`
	NodeNum                  types.Int32               `tfsdk:"node_num"`
	AzInfos                  []CtyunCcseNodePoolAzInfo `tfsdk:"az_infos"`
}

type CtyunCcseNodePoolAzInfo struct {
	AzName types.String `tfsdk:"az_name"`
}

type CtyunCcseNodePoolDisk struct {
	Type types.String `tfsdk:"type"`
	Size types.Int32  `tfsdk:"size"`
}

func (c *ctyunCcseNodePool) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10083472/10318452**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
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
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "集群ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "节点池名称，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围：month：按月，year：按年、on_demand：按需。当此值为month或者year时，cycle_count为必填，支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypes...),
				},
			},
			"cycle_count": schema.Int64Attribute{
				Optional:    true,
				Description: "订购时长，该参数在cycle_type为month或year时才生效，当cycle_type=month，支持订购1-11个月；当cycle_type=year，支持订购1-5年，支持更新",
				Validators: []validator.Int64{
					validator2.AlsoRequiresEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeMonth),
						types.StringValue(business.OrderCycleTypeYear),
					),
					validator2.ConflictsWithEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
					validator2.CycleCount(1, 11, 1, 5),
				},
			},
			"auto_renew": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否自动续订，默认非自动续订，当cycle_type不等于on_demand时才可填写，支持更新",
				Default:     booldefault.StaticBool(false),
				Validators: []validator.Bool{
					validator2.ConflictsWithEqualBool(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
				}},
			"visibility_post_host_script": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "部署后执行自定义脚本，base64编码，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"visibility_host_script": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "部署前执行自定义脚本，base64编码，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"instance_type": schema.StringAttribute{
				Required:    true,
				Description: "实例类型，支持ecs（云主机）、ebm（裸金属）",
				Validators: []validator.String{
					stringvalidator.OneOf(business.CcseSlaveInstanceTypeEcs, business.CcseSlaveInstanceTypeEbm),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mirror_id": schema.StringAttribute{
				Optional:    true,
				Description: "镜像id，实例为ecs类型必填，可查看<a href=\"https://www.ctyun.cn/document/10083472/11004475\">节点规格和节点镜像</a>",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("instance_type"),
						types.StringValue(business.CcseSlaveInstanceTypeEcs),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("instance_type"),
						types.StringValue(business.CcseSlaveInstanceTypeEbm),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mirror_name": schema.StringAttribute{
				Optional:    true,
				Description: "镜像名称，实例为ebm类型必填，可查看<a href=\"https://www.ctyun.cn/document/10083472/11004475\">节点规格和节点镜像</a>",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("instance_type"),
						types.StringValue(business.CcseSlaveInstanceTypeEbm),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("instance_type"),
						types.StringValue(business.CcseSlaveInstanceTypeEcs),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mirror_type": schema.Int32Attribute{
				Required:    true,
				Description: "镜像类型，支持传0（私有），1（公有），可查看<a href=\"https://www.ctyun.cn/document/10026730/10030151\">镜像概述</a>",
				Validators: []validator.Int32{
					int32validator.Between(0, 1),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"item_def_name": schema.StringAttribute{
				Required:    true,
				Description: "实例规格名称，使用至少4C8G以上的规格，云主机规格通过ctyun_ecs_flavors查询，裸金属规格通过ctyun_ebm_device_types查询",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"key_pair_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "密钥对名称，与password有且只能有一个",
				Default:     stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("password"),
					),
				},
			},

			"password": schema.StringAttribute{
				Optional:    true,
				Description: "用户密码，与key_pair_name有且只能有一个，需要满足以下规则：长度在8～30个字符；必须包含大写字母、小写字母、数字以及特殊符号中的三项；特殊符号可选：()`~!@#$%^&*_-+=|{}[]:;'<>,.?/\\且不能以斜线号/开头",
				Validators: []validator.String{
					stringvalidator.Any(
						stringvalidator.All(
							validator2.AlsoRequiresEqualString(
								path.MatchRoot("instance_type"),
								types.StringValue(business.CcseSlaveInstanceTypeEcs),
							),
							validator2.EcsPassword(),
						),

						stringvalidator.All(
							validator2.AlsoRequiresEqualString(
								path.MatchRoot("instance_type"),
								types.StringValue(business.CcseSlaveInstanceTypeEbm),
							),
							validator2.EbmPassword(),
						),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("key_pair_name"),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Sensitive: true,
			},
			"use_affinity_group": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否使用主机组，默认不使用",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"affinity_group_id": schema.StringAttribute{
				Optional:    true,
				Description: "云主机组id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"max_pod_num": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "最大pod数, 默认110",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"sys_disk": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "系统盘信息",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    true,
						Description: "系统盘类型，支持SATA、SAS、SSD，支持更新",
						Validators: []validator.String{
							stringvalidator.OneOf(business.CcseDiskTypes...),
						},
					},
					"size": schema.Int32Attribute{
						Required:    true,
						Description: "系统盘大小，单位为G，支持范围40-2040，支持更新",
						Validators: []validator.Int32{
							int32validator.Between(40, 2040),
						},
					},
				},
			},
			"node_num": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "节点数，不填则默认为0，支持更新，创建节点池后只能增加不能减少",
				Default:     int32default.StaticInt32(0),
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			"az_infos": schema.ListNestedAttribute{
				Required: true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				Description: "可用区信息，支持的可用区可通过ctyun_regions查询",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"az_name": schema.StringAttribute{
							Required:    true,
							Description: "可用区编码",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
					},
				},
			},
			"data_disks": schema.ListNestedAttribute{
				Optional:    true,
				Description: "数据盘信息，支持更新",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:    true,
							Description: "数据盘类型，支持SATA、SAS、SSD，支持更新",
							Validators: []validator.String{
								stringvalidator.OneOf(business.CcseDiskTypes...),
							},
						},
						"size": schema.Int32Attribute{
							Optional:    true,
							Description: "数据盘大小，单位为G，支持范围10-20000，支持更新",
							Validators: []validator.Int32{
								int32validator.Between(10, 20000),
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunCcseNodePool) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunCcseNodePoolConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建
	err = c.create(ctx, plan)
	if err != nil {
		return
	}
	// 创建后检查
	id, err := c.checkAfterCreate(ctx, plan)
	if err != nil {
		return
	}

	plan.ID = types.StringValue(id)
	// 扩容
	planB := plan
	planB.NodeNum = types.Int32Value(0)
	err = c.scaleUp(ctx, plan, planB)
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

func (c *ctyunCcseNodePool) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunCcseNodePoolConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunCcseNodePool) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunCcseNodePoolConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunCcseNodePoolConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}
	// 扩容
	err = c.scaleUp(ctx, plan, state)
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

func (c *ctyunCcseNodePool) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunCcseNodePoolConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
	//response.State.RemoveResource(ctx)
}

func (c *ctyunCcseNodePool) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [id],[clusterID],[regionID]
func (c *ctyunCcseNodePool) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunCcseNodePoolConfig
	var id, clusterID, regionID string
	err = terraform_extend.Split(request.ID, &id, &clusterID, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.ClusterID = types.StringValue(clusterID)
	cfg.ID = types.StringValue(id)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// checkAfterCreate 创建后检查
func (c *ctyunCcseNodePool) checkAfterCreate(ctx context.Context, plan CtyunCcseNodePoolConfig) (id string, err error) {
	params := &ccse2.CcseListNodePoolsRequest{
		ClusterId:    plan.ClusterID.ValueString(),
		RegionId:     plan.RegionID.ValueString(),
		NodePoolName: plan.NodePoolName.ValueString(),
		PageSize:     1,
		PageNow:      1,
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseListNodePoolsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if len(resp.ReturnObj.Records) == 0 {
		err = common.InvalidReturnObjResultsError
	}
	id = resp.ReturnObj.Records[0].Id
	return
}

// create 创建
func (c *ctyunCcseNodePool) create(ctx context.Context, plan CtyunCcseNodePoolConfig) (err error) {
	params := &ccse2.CcseCreateNodePoolRequest{
		ClusterId:                plan.ClusterID.ValueString(),
		RegionId:                 plan.RegionID.ValueString(),
		NodePoolName:             plan.NodePoolName.ValueString(),
		AutoRenewStatus:          map[bool]int32{false: 0, true: 1}[plan.AutoRenew.ValueBool()],
		VisibilityPostHostScript: plan.VisibilityPostHostScript.ValueString(),
		VisibilityHostScript:     plan.VisibilityHostScript.ValueString(),
		UseAffinityGroup:         plan.UseAffinityGroup.ValueBoolPointer(),
		AffinityGroupUuid:        plan.AffinityGroupID.ValueString(),
		MaxPodNum:                plan.MaxPodNum.ValueInt32(),
		ImageType:                plan.MirrorType.ValueInt32(),
	}
	if plan.SysDisk != nil {
		params.SysDiskType = plan.SysDisk.Type.ValueString()
		params.SysDiskSize = plan.SysDisk.Size.ValueInt32()
	}

	switch plan.InstanceType.ValueString() {
	case business.CcseSlaveInstanceTypeEcs:
		params.ImageUuid = plan.MirrorID.ValueString()
	case business.CcseSlaveInstanceTypeEbm:
		params.ImageName = plan.MirrorName.ValueString()
	}

	if plan.Password.ValueString() == "" && plan.KeyPairName.ValueString() == "" {
		err = fmt.Errorf("password和key_pair_name两者不能都为空")
		return
	} else if plan.Password.ValueString() != "" {
		params.LoginType = "password"
		params.EcsPasswd = plan.Password.ValueString()
	} else {
		params.LoginType = "secretPair"
		params.KeyName = plan.KeyPairName.ValueString()
		// 查keypair_id
		params.KeyPairId, err = business.NewKeyPairService(c.meta).GetKeyPairID(
			ctx,
			plan.KeyPairName.ValueString(),
			plan.RegionID.ValueString(),
			"",
		)
		if err != nil {
			return
		}
	}

	for _, az := range plan.AzInfos {
		params.AzInfo = append(params.AzInfo, &ccse2.CcseCreateNodePoolAzInfoRequest{AzName: az.AzName.ValueString()})
	}

	// 处理订单
	switch plan.CycleType.ValueString() {
	case business.OnDemandCycleType:
		params.BillMode = "2"
	case business.MonthCycleType:
		params.BillMode = "1"
		params.CycleType = "MONTH"
		params.CycleCount = int32(plan.CycleCount.ValueInt64())
	case business.YearCycleType:
		params.BillMode = "1"
		params.CycleType = "YEAR"
		params.CycleCount = int32(plan.CycleCount.ValueInt64())
	}

	// 处理规格
	switch plan.InstanceType.ValueString() {
	case business.CcseSlaveInstanceTypeEcs:
		flavorName := plan.ItemDefName.ValueString()
		flavor, err := business.NewEcsService(c.meta).GetFlavorByName(ctx, flavorName, plan.RegionID.ValueString())
		if err != nil {
			return err
		}
		params.Cpu = int32(flavor.FlavorCpu)
		params.Memory = int32(flavor.FlavorRam)
		params.VmSpecName = flavorName
		params.VmType = flavor.FlavorType
	case business.CcseSlaveInstanceTypeEbm:
		deviceType := plan.ItemDefName.ValueString()
		flavor, err := business.NewEbmService(c.meta).GetDeviceType(ctx, deviceType, plan.RegionID.ValueString(), plan.AzInfos[0].AzName.ValueString())
		if err != nil {
			return err
		}
		params.Cpu = flavor.CpuAmount
		params.Memory = flavor.MemAmount
		params.VmSpecName = deviceType
		params.VmType = deviceType
	}

	for _, disk := range plan.DataDisks {
		params.DataDisks = append(params.DataDisks, &ccse2.CcseCreateNodePoolDataDisksRequest{
			DiskSpecName: disk.Type.ValueString(),
			Size:         disk.Size.ValueInt32(),
		})
	}

	resp, err := c.meta.Apis.SdkCcseApis.CcseCreateNodePoolApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return
}

// 扩容
func (c *ctyunCcseNodePool) scaleUp(ctx context.Context, plan, state CtyunCcseNodePoolConfig) (err error) {
	if plan.NodeNum.Equal(state.NodeNum) {
		return
	}
	if plan.NodeNum.ValueInt32() < state.NodeNum.ValueInt32() {
		err = fmt.Errorf("不支持减少节点数")
		return
	}

	params := &ccse2.CcseScaleUpNodePoolRequest{
		ClusterId:  state.ClusterID.ValueString(),
		NodePoolId: state.ID.ValueString(),
		RegionId:   state.RegionID.ValueString(),
		Num:        plan.NodeNum.ValueInt32() - state.NodeNum.ValueInt32(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseScaleUpNodePoolApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	plan.ID = state.ID
	return c.checkAfterScaleUp(ctx, plan)
}

// checkAfterScaleUp 扩容后检查
func (c *ctyunCcseNodePool) checkAfterScaleUp(ctx context.Context, plan CtyunCcseNodePoolConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var pool *ccse2.CcseGetNodePoolReturnObjResponse
			pool, err = c.getNodePoolByID(ctx, plan)
			if err != nil {
				return false
			}
			if pool.NormalNodeNum < plan.NodeNum.ValueInt32() {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("扩容时间过长")
		return
	}
	return
}

// 缩容
func (c *ctyunCcseNodePool) scaleDown(ctx context.Context, state CtyunCcseNodePoolConfig) (err error) {
	pool, err := c.getNodePoolByID(ctx, state)
	if err != nil {
		return
	}
	var nodes []string
	for _, node := range pool.Nodes {
		nodes = append(nodes, node.NodeName)
	}

	if len(nodes) == 0 {
		return
	}

	params := &ccse2.CcseScaleDownNodePoolRequest{
		ClusterId:  state.ClusterID.ValueString(),
		NodePoolId: state.ID.ValueString(),
		RegionId:   state.RegionID.ValueString(),
		NodeNames:  nodes,
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseScaleDownNodePoolApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return c.checkAfterScaleDown(ctx, state)
}

// checkAfterScaleUp 扩容后检查
func (c *ctyunCcseNodePool) checkAfterScaleDown(ctx context.Context, plan CtyunCcseNodePoolConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var pool *ccse2.CcseGetNodePoolReturnObjResponse
			pool, err = c.getNodePoolByID(ctx, plan)
			if err != nil {
				return false
			}
			if pool.NodeTotalNum > 0 {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("扩容时间过长")
		return
	}
	return
}

// getNodePoolByID 根据ID查询节点池
func (c *ctyunCcseNodePool) getNodePoolByID(ctx context.Context, plan CtyunCcseNodePoolConfig) (pool *ccse2.CcseGetNodePoolReturnObjResponse, err error) {
	params := &ccse2.CcseGetNodePoolRequest{
		ClusterId:  plan.ClusterID.ValueString(),
		RegionId:   plan.RegionID.ValueString(),
		NodePoolId: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseGetNodePoolApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	pool = resp.ReturnObj
	return
}

// getAndMerge 从远端查询
func (c *ctyunCcseNodePool) getAndMerge(ctx context.Context, plan *CtyunCcseNodePoolConfig) (err error) {
	pool, err := c.getNodePoolByID(ctx, *plan)
	if err != nil {
		return
	}
	plan.NodePoolName = types.StringValue(pool.NodePoolName)
	plan.MirrorType = types.Int32Value(pool.ImageType)
	plan.NodeNum = types.Int32Value(pool.NormalNodeNum)
	listParams := &ccse2.CcseListNodePoolsRequest{
		ClusterId:    plan.ClusterID.ValueString(),
		RegionId:     plan.RegionID.ValueString(),
		NodePoolName: plan.NodePoolName.ValueString(),
		PageNow:      1,
		PageSize:     10,
	}

	listResp, err := c.meta.Apis.SdkCcseApis.CcseListNodePoolsApi.Do(ctx, c.meta.SdkCredential, listParams)
	if err != nil {
		return
	} else if listResp.ReturnObj == nil {
		err = fmt.Errorf("API return error. Message: %s", listResp.Message)
		return
	}
	// 解析返回值
	records := listResp.ReturnObj.Records
	if len(records) == 0 {
		return
	}
	p := records[0]
	if plan.SysDisk != nil {
		plan.SysDisk.Size = types.Int32Value(p.SysDiskSize)
		plan.SysDisk.Type = types.StringValue(p.SysDiskType)
	}
	plan.VisibilityPostHostScript = types.StringValue(p.VisibilityPostHostScript)
	plan.VisibilityHostScript = types.StringValue(p.VisibilityHostScript)
	plan.MaxPodNum = types.Int32Value(p.MaxPodNum)
	plan.AutoRenew = types.BoolValue(map[int32]bool{0: false, 1: true}[p.AutoRenewStatus])
	plan.ItemDefName = types.StringValue(p.VmSpecName)
	plan.KeyPairName = types.StringValue(p.KeyName)

	switch plan.InstanceType.ValueString() {
	case business.CcseSlaveInstanceTypeEcs:
		plan.MirrorID = types.StringValue(p.ImageUuid)
	case business.CcseSlaveInstanceTypeEbm:
		plan.MirrorName = types.StringValue(p.ImageName)
	}

	switch p.BillMode {
	case "1":
		plan.CycleType = types.StringValue(strings.ToLower(p.CycleType))
		plan.CycleCount = types.Int64Value(int64(p.CycleCount))
	case "2":
		plan.CycleType = types.StringValue(business.OnDemandCycleType)
	}
	if strings.HasPrefix(p.VmSpecName, "physical") {
		plan.InstanceType = types.StringValue(business.CcseSlaveInstanceTypeEbm)
	} else {
		plan.InstanceType = types.StringValue(business.CcseSlaveInstanceTypeEcs)
	}
	plan.DataDisks = nil
	for _, disk := range p.DataDisks {
		plan.DataDisks = append(plan.DataDisks, CtyunCcseNodePoolDisk{
			Size: types.Int32Value(int32(disk.Size)),
			Type: types.StringValue(disk.DiskSpecName),
		})
	}

	return
}

// update 更新
func (c *ctyunCcseNodePool) update(ctx context.Context, plan, state CtyunCcseNodePoolConfig) (err error) {
	params := &ccse2.CcseUpdateNodePoolRequest{
		ClusterId:                state.ClusterID.ValueString(),
		NodePoolId:               state.ID.ValueString(),
		RegionId:                 state.RegionID.ValueString(),
		NodePoolName:             plan.NodePoolName.ValueString(),
		AutoRenewStatus:          map[bool]int32{false: 0, true: 1}[plan.AutoRenew.ValueBool()],
		VisibilityPostHostScript: plan.VisibilityPostHostScript.ValueString(),
		VisibilityHostScript:     plan.VisibilityHostScript.ValueString(),
		SysDiskType:              plan.SysDisk.Type.ValueString(),
		SysDiskSize:              plan.SysDisk.Size.ValueInt32(),
	}

	for _, disk := range plan.DataDisks {
		params.DataDisks = append(params.DataDisks, &ccse2.CcseUpdateNodePoolDataDisksRequest{
			DiskSpecName: disk.Type.ValueString(),
			Size:         disk.Size.ValueInt32(),
		})
	}
	// 处理订单
	switch plan.CycleType.ValueString() {
	case business.OnDemandCycleType:
		params.BillMode = "2"
	case business.MonthCycleType:
		params.BillMode = "1"
		params.CycleType = "MONTH"
		params.CycleCount = int32(plan.CycleCount.ValueInt64())
	case business.YearCycleType:
		params.BillMode = "1"
		params.CycleType = "YEAR"
		params.CycleCount = int32(plan.CycleCount.ValueInt64())
	}

	resp, err := c.meta.Apis.SdkCcseApis.CcseUpdateNodePoolApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}

	return
}

// delete 删除
func (c *ctyunCcseNodePool) delete(ctx context.Context, plan CtyunCcseNodePoolConfig) (err error) {
	err = c.scaleDown(ctx, plan)
	if err != nil {
		return
	}
	params := &ccse2.CcseDeleteNodePoolRequest{
		ClusterId:  plan.ClusterID.ValueString(),
		RegionId:   plan.RegionID.ValueString(),
		NodePoolId: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseDeleteNodePoolApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}
