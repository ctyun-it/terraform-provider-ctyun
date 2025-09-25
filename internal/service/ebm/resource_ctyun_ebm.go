package ebm

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebm"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
	"time"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
)

var (
	_ resource.Resource                = &ctyunEbm{}
	_ resource.ResourceWithConfigure   = &ctyunEbm{}
	_ resource.ResourceWithImportState = &ctyunEbm{}
)

type ctyunEbm struct {
	meta *common.CtyunMetadata
}

func NewCtyunEbm() resource.Resource {
	return &ctyunEbm{}
}

func (c *ctyunEbm) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebm"
}

type CtyunEbmConfig struct {
	ID                   types.String `tfsdk:"id"`
	InstanceID           types.String `tfsdk:"instance_id"`
	RegionID             types.String `tfsdk:"region_id"`
	AzName               types.String `tfsdk:"az_name"`
	DeviceType           types.String `tfsdk:"device_type"`
	InstanceName         types.String `tfsdk:"instance_name"`
	Name                 types.String `tfsdk:"name"`
	Hostname             types.String `tfsdk:"hostname"`
	ImageUUID            types.String `tfsdk:"image_uuid"`
	ActualImageID        types.String `tfsdk:"actual_image_id"`
	Password             types.String `tfsdk:"password"`
	ProjectID            types.String `tfsdk:"project_id"`
	SystemDiskType       types.String `tfsdk:"system_disk_type"`
	SystemDiskSize       types.Int32  `tfsdk:"system_disk_size"`
	SystemDiskID         types.String `tfsdk:"system_disk_id"`
	SystemVolumeRaidUUID types.String `tfsdk:"system_volume_raid_uuid"`
	DataVolumeRaidUUID   types.String `tfsdk:"data_volume_raid_uuid"`
	VpcID                types.String `tfsdk:"vpc_id"`
	EipID                types.String `tfsdk:"eip_id"`
	EipAddress           types.String `tfsdk:"eip_address"`
	SecurityGroupIDs     types.Set    `tfsdk:"security_group_ids"`
	UserData             types.String `tfsdk:"user_data"`
	KeyPairName          types.String `tfsdk:"key_pair_name"`
	AutoRenew            types.Bool   `tfsdk:"auto_renew"`
	CycleCount           types.Int64  `tfsdk:"cycle_count"`
	CycleType            types.String `tfsdk:"cycle_type"`
	MasterOrderID        types.String `tfsdk:"master_order_id"`
	Status               types.String `tfsdk:"status"`
	FixedIP              types.String `tfsdk:"fixed_ip"`
	SubnetID             types.String `tfsdk:"subnet_id"`
	PortID               types.String `tfsdk:"port_id"`
	InterfaceID          types.String `tfsdk:"interface_id"`
}

func (c *ctyunEbm) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027724`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "名称",
			},
			"instance_id": schema.StringAttribute{
				Computed:    true,
				Description: "物理机UUID，值与id相同",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订单id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称",
				Default:     defaults.AcquireFromGlobalString(common.ExtraAzName, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"device_type": schema.StringAttribute{
				Required:    true,
				Description: "物理机套餐类型，可通过ctyun_ebm_device_types查询",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"instance_name": schema.StringAttribute{
				Required:    true,
				Description: "物理机名称，长度为2-31位，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 31),
				},
			},
			"hostname": schema.StringAttribute{
				Required:    true,
				Description: "hostname，linux系统2到63位长度；windows系统2-15位长度；<br/>允许使用大小写字母、数字、连字符'-'、点号'.'，不能连续使用'-'或者'.'，'-'和'.'不能用于开头或结尾，不能仅使用数字，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 63),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-\.]*[a-zA-Z0-9]$`), "hostname必须以字母开头，以字母或数字结尾"),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9\-\.]*$`), "hostname只能包含字母、数字、连字符和点号"),
					stringvalidator.RegexMatches(regexp.MustCompile(`^.*[a-zA-Z].*$`), "hostname不能仅使用数字"),
				},
			},
			"image_uuid": schema.StringAttribute{
				Required:    true,
				Description: "物理机镜像id，可通过ctyun_ebm_device_images查询",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"actual_image_id": schema.StringAttribute{
				Computed:    true,
				Description: "实际镜像id，重装、集群纳管等操作会导致actual_image_id与image_id不同",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password": schema.StringAttribute{
				Sensitive:   true,
				Optional:    true,
				Computed:    true,
				Description: "密码(必须包含大小写字母和（一个数字或者特殊字符）长度8到30位)，未传入有效的keyName时必须传入password，支持更新",
				Validators: []validator.String{
					validator2.EbmPassword(),
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("key_pair_name"),
					}...),
				},
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
			"system_volume_raid_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "本地系统盘raid类型，如果有本地盘则必填，可通过ctyun_ebm_device_raids查询",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("system_disk_type")),
				},
			},
			"data_volume_raid_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "本地数据盘raid类型，如果有本地盘则必填，可通过ctyun_ebm_device_raids查询",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "主网卡虚拟私有云ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"eip_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "弹性公网IP的ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.EipValidate(),
				},
			},
			"eip_address": schema.StringAttribute{
				Computed:    true,
				Description: "弹性公网IP的地址",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_group_ids": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				Description: "主网卡安全组ID，套餐smart_nic_exist为true可支持安全组。创建弹性裸金属必须传入安全组ID，标准裸金属不支持传入安全组ID",
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(validator2.SecurityGroupValidate()),
				},
			},
			"system_disk_type": schema.StringAttribute{
				Optional:    true,
				Description: "系统盘类型，sata：普通IO，sas：高IO，ssd：超高IO",
				Validators: []validator.String{
					stringvalidator.OneOf(business.EbmDiskTypes...),
					stringvalidator.ConflictsWith(path.MatchRoot("system_volume_raid_uuid")),
					stringvalidator.AlsoRequires(path.MatchRoot("system_disk_size")),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"system_disk_size": schema.Int32Attribute{
				Optional:    true,
				Description: "系统盘大小，单位为G，取值范围：[100, 2048]，当前不支持公网",
				Validators: []validator.Int32{
					int32validator.Between(100, 2048),
					int32validator.ConflictsWith(path.MatchRoot("system_volume_raid_uuid")),
					int32validator.AlsoRequires(path.MatchRoot("system_disk_type")),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"system_disk_id": schema.StringAttribute{
				Computed:    true,
				Description: "系统盘的id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "主网卡的子网id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"fixed_ip": schema.StringAttribute{
				Computed:    true,
				Description: "加入子网后的ip地址",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"port_id": schema.StringAttribute{
				Computed:    true,
				Description: "主网卡PORT UUID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"interface_id": schema.StringAttribute{
				Computed:    true,
				Description: "主网卡UUID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_data": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "用户自定义数据，需要以Base64方式编码，Base64编码后的长度限制为1-16384字符",
				Default:     stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 16384),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key_pair_name": schema.StringAttribute{
				Optional:    true,
				Description: "密钥对名词，和password只能传其中之一",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("password"),
					}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auto_renew": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否自动续订，默认非自动续订，当cycle_type不等于on_demand时才可填写。",
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
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围:[on_demand=按需,month=按月,year=按年]",
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypeOnDemand, business.OrderCycleTypeYear, business.OrderCycleTypeMonth),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_count": schema.Int64Attribute{
				Optional:    true,
				Description: "订购时长，最长订购周期为60个月（5年）；非按需时必填",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					validator2.AlsoRequiresEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeYear),
						types.StringValue(business.OrderCycleTypeMonth),
					),
					validator2.ConflictsWithEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
					validator2.CycleCount(1, 11, 1, 5),
				},
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "物理机状态，支持running（开机）和stopped（关机），默认running，支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(
						business.EbmStatusRunning,
						business.EbmStatusStopped,
					),
				},
				Default: stringdefault.StaticString(business.EbmStatusRunning),
			},
		},
	}
}

func (c *ctyunEbm) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEbmConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建前检查
	err = c.checkBeforeCreateInstance(ctx, plan)
	if err != nil {
		return
	}
	// 创建
	returnObj, err := c.createInstance(ctx, plan)
	if err != nil {
		return
	}

	// 先保存订单号
	masterOrderId := *returnObj.MasterOrderID
	plan.MasterOrderID = types.StringValue(masterOrderId)
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 根据订单号轮询查资源的uuid
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	loop, err := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId, 600)
	if err != nil {
		return
	}
	plan.InstanceID = types.StringValue(loop.Uuid[0])
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建机器后状态默认为启动状态，可根据用户要求的状态，去执行对应的操作，比如关机
	err = c.handleInstance(ctx, plan, business.EbmStatusRunning, plan.Status.ValueString())
	if err != nil {
		return
	}
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *ctyunEbm) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbmConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 通过订单号同步
	if !c.acquireAndSetIdIfOrderNotFinished(ctx, &state, response) {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "instance is not found") {
			// 查下主网卡是否存在
			var exist bool
			portID := state.PortID.ValueString()
			exist, err = business.NewPortService(c.meta).Exist(ctx, portID, state.RegionID.ValueString())
			if err != nil {
				return
			}
			// 主网卡存在，则监听到主网卡删除为止
			if exist {
				err = c.checkAfterDelete(ctx, state)
				if err != nil {
					return
				}
			}
			// 主网卡不存在，清理state
			response.State.RemoveResource(ctx)
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEbm) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunEbmConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunEbmConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 处理开关机
	err = c.handleInstance(ctx, state, state.Status.ValueString(), plan.Status.ValueString())
	if err != nil {
		return
	}
	// 修改实例名称
	err = c.updateInstanceName(ctx, state, plan)
	if err != nil {
		return
	}
	// 修改密码或主机名
	err = c.updatePasswordOrHostname(ctx, state, plan)
	if err != nil {
		return
	}
	state.Password = plan.Password
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

func (c *ctyunEbm) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbmConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 关机
	err = c.handleInstance(ctx, state, state.Status.ValueString(), business.EbmStatusStopped)
	if err != nil {
		return
	}
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
	err = c.destroy(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDelete(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunEbm) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [instanceID],[regionID],[azName]
func (c *ctyunEbm) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEbmConfig
	var instanceUUID, regionID, azName string
	err = terraform_extend.Split(request.ID, &instanceUUID, &regionID, &azName)
	if err != nil {
		return
	}

	plan.InstanceID = types.StringValue(instanceUUID)
	plan.AzName = types.StringValue(azName)
	plan.RegionID = types.StringValue(regionID)
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

// createInstance 创建物理机
func (c *ctyunEbm) createInstance(ctx context.Context, plan CtyunEbmConfig) (returnObj ctebm.EbmCreateInstanceV4plusReturnObjResponse, err error) {
	regionID := plan.RegionID.ValueString()
	projectID := plan.ProjectID.ValueString()
	azName := plan.AzName.ValueString()
	password := plan.Password.ValueString()
	systemVolumeRaidUUID := plan.SystemVolumeRaidUUID.ValueString()
	dataVolumeRaidUUID := plan.DataVolumeRaidUUID.ValueString()
	userData := plan.UserData.ValueString()
	keyName := plan.KeyPairName.ValueString()
	securityGroupIDs, _ := c.buildSecGroupList(ctx, plan)
	securityGroupStr := strings.Join(securityGroupIDs, ",")
	params := &ctebm.EbmCreateInstanceV4plusRequest{
		RegionID:        regionID,
		AzName:          azName,
		DeviceType:      plan.DeviceType.ValueString(),
		InstanceName:    plan.InstanceName.ValueString(),
		Hostname:        plan.Hostname.ValueString(),
		ImageUUID:       plan.ImageUUID.ValueString(),
		VpcID:           plan.VpcID.ValueString(),
		ProjectID:       &projectID,
		AutoRenewStatus: map[bool]int32{false: 0, true: 1}[plan.AutoRenew.ValueBool()],
		ClientToken:     uuid.NewString(),
		OrderCount:      1,
		SecurityGroupID: &securityGroupStr,
		NetworkCardList: []*ctebm.EbmCreateInstanceV4plusNetworkCardListRequest{{Master: true, SubnetID: plan.SubnetID.ValueString()}},
	}

	if plan.EipID.ValueString() != "" {
		params.PublicIP = plan.EipID.ValueStringPointer()
		params.ExtIP = business.EbmExtIpUseExist
	} else {
		params.ExtIP = business.EbmExtIpNotUse
	}

	if password != "" {
		params.Password = &password
	} else if keyName != "" {
		params.KeyName = &keyName
	} else {
		err = fmt.Errorf("password or keyname is empty")
	}
	if userData != "" {
		params.UserData = &userData
	}
	if systemVolumeRaidUUID != "" {
		params.SystemVolumeRaidUUID = &systemVolumeRaidUUID
	}
	if dataVolumeRaidUUID != "" {
		params.DataVolumeRaidUUID = &dataVolumeRaidUUID
	}
	if !plan.SystemDiskType.IsNull() {
		params.DiskList = []*ctebm.EbmCreateInstanceV4plusDiskListRequest{
			{
				DiskType: "system",
				Size:     plan.SystemDiskSize.ValueInt32(),
				RawType:  strings.ToUpper(plan.SystemDiskType.ValueString()),
			},
		}
	}

	switch plan.CycleType.ValueString() {
	case business.OrderCycleTypeOnDemand:
		params.InstanceChargeType = business.EbmOrderOnDemand
	case business.OrderCycleTypeMonth, business.OrderCycleTypeYear:
		params.InstanceChargeType = business.EbmOrderOnCycle
		params.CycleType = strings.ToUpper(plan.CycleType.ValueString())
		params.CycleCount = int32(plan.CycleCount.ValueInt64())
	}

	resp, err := c.meta.Apis.CtEbmApis.EbmCreateInstanceV4plusApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	returnObj = *resp.ReturnObj
	return
}

// checkBeforeCreateInstance 创建前检查
func (c *ctyunEbm) checkBeforeCreateInstance(ctx context.Context, plan CtyunEbmConfig) error {
	// 确保当前虚拟私有云存在，且子网与虚拟私有云存在对应关系
	vpc := plan.VpcID.ValueString()
	subnets, err := business.NewVpcService(c.meta).GetVpcSubnet(ctx, vpc, plan.RegionID.ValueString(), plan.ProjectID.ValueString())
	if err != nil {
		return err
	}
	// 查询套餐
	deviceTypeConfig, err := c.getDeviceTypeConfig(ctx, plan)
	if err != nil {
		return err
	}

	subnetID := plan.SubnetID.ValueString()
	if subnet, ok := subnets[subnetID]; !ok {
		return fmt.Errorf("子网 %s 不属于 %s", subnetID, vpc)
	} else if *deviceTypeConfig.SmartNicExist && subnet.Type != business.SubnetTypeCommonInt {
		return fmt.Errorf("该套餐 %s 为弹性裸金属, 必须使用普通子网", plan.DeviceType.ValueString())
	} else if !*deviceTypeConfig.SmartNicExist && subnet.Type != business.SubnetTypeEbmInt {
		return fmt.Errorf("该套餐 %s 为标准裸金属, 必须使用裸金属子网", plan.DeviceType.ValueString())
	}
	// 弹性裸金属必须有安全组id，标准裸金属一定不能有安全组id
	secGroup, err := c.buildSecGroupList(ctx, plan)
	if err != nil {
		return err
	}
	if *deviceTypeConfig.SmartNicExist && len(secGroup) == 0 {
		return fmt.Errorf("该套餐 %s 为弹性裸金属，必须传递安全组ID", plan.DeviceType.ValueString())
	}
	if !*deviceTypeConfig.SmartNicExist && len(secGroup) != 0 {
		return fmt.Errorf("该套餐 %s 为标准裸金属，不能传递安全组ID", plan.DeviceType.ValueString())
	}
	// 安全组必须存在
	for _, g := range secGroup {
		err = business.NewSecurityGroupService(c.meta).MustExist(ctx, g, plan.RegionID.ValueString())
		if err != nil {
			return err
		}
	}

	// 校验eip
	if plan.EipID.ValueString() != "" {
		err = business.NewEipService(c.meta).MustExist(ctx, plan.EipID.ValueString(), plan.RegionID.ValueString())
		if err != nil {
			return err
		}
	}

	// 高级版必须关联云硬盘
	if *deviceTypeConfig.CloudBoot && plan.SystemDiskType.IsNull() {
		return fmt.Errorf("该套餐 %s 需要从云硬盘启动，必须关联云硬盘", plan.DeviceType.ValueString())
	}
	if deviceTypeConfig.SystemVolumeAmount > 0 && plan.SystemVolumeRaidUUID.ValueString() == "" {
		return fmt.Errorf("该套餐 %s 必须传递本地系统盘ID", plan.DeviceType.ValueString())
	}
	if deviceTypeConfig.DataVolumeAmount > 0 && plan.DataVolumeRaidUUID.ValueString() == "" {
		return fmt.Errorf("该套餐 %s 必须传递本地数据盘ID", plan.DeviceType.ValueString())
	}

	// 检查库存
	enough, err := c.checkStock(ctx, plan)
	if err != nil {
		return err
	} else if !enough {
		return fmt.Errorf("该套餐 %s 库存不足", plan.DeviceType.ValueString())
	}

	// 检查镜像
	available, err := c.checkImage(ctx, plan)
	if err != nil {
		return err
	} else if !available {
		return fmt.Errorf("该套餐 %s 不能使用镜像 %s", plan.DeviceType.ValueString(), plan.ImageUUID.ValueString())
	}

	return nil
}

// checkImage 检查镜像是否可用
func (c *ctyunEbm) checkImage(ctx context.Context, plan CtyunEbmConfig) (available bool, err error) {
	deviceType := plan.DeviceType.ValueString()
	imageUUID := plan.ImageUUID.ValueString()
	params := &ctebm.EbmImageListRequest{
		RegionID:   plan.RegionID.ValueString(),
		AzName:     plan.AzName.ValueString(),
		DeviceType: deviceType,
		ImageUUID:  &imageUUID,
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmImageListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	available = len(resp.ReturnObj.Results) > 0
	return
}

// checkStock 获取库存
func (c *ctyunEbm) checkStock(ctx context.Context, plan CtyunEbmConfig) (enough bool, err error) {
	deviceType := plan.DeviceType.ValueString()
	params := &ctebm.EbmDeviceStockListRequest{
		RegionID:   plan.RegionID.ValueString(),
		AzName:     plan.AzName.ValueString(),
		DeviceType: &deviceType,
		Count:      1,
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmDeviceStockListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	enough = *resp.ReturnObj.Results[0].Success
	return
}

// getDeviceTypeConfig 查询套餐详情
func (c *ctyunEbm) getDeviceTypeConfig(ctx context.Context, plan CtyunEbmConfig) (result ctebm.EbmDeviceTypeListReturnObjResultsResponse, err error) {
	deviceType := plan.DeviceType.ValueString()
	params := &ctebm.EbmDeviceTypeListRequest{
		RegionID:   plan.RegionID.ValueString(),
		AzName:     plan.AzName.ValueString(),
		DeviceType: &deviceType,
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmDeviceTypeListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	if len(resp.ReturnObj.Results) == 0 {
		err = fmt.Errorf("未查询到该套餐 %s", deviceType)
		return
	}
	result = *resp.ReturnObj.Results[0]
	return
}

// buildSecGroup 构建安全组列表
func (c *ctyunEbm) buildSecGroupList(ctx context.Context, plan CtyunEbmConfig) (secGroupIDs []string, err error) {
	if plan.SecurityGroupIDs.IsNull() {
		return
	}
	// 处理安全组集合
	diags := plan.SecurityGroupIDs.ElementsAs(ctx, &secGroupIDs, true) // 第二个参数为是否忽略未知值
	if diags.HasError() {
		err = fmt.Errorf("invalid security group ids")
		return
	}
	return
}

// handleInstance 操作机器，开机或关机
func (c *ctyunEbm) handleInstance(ctx context.Context, plan CtyunEbmConfig, currentStatus string, targetStatus string) (err error) {
	for i := 0; i < 100; i++ {
		if currentStatus == targetStatus {
			return
		}
		switch currentStatus {
		case business.EbmStatusStopped: // 当前是关机，目标肯定是开机
			return c.startInstance(ctx, plan)
		case business.EbmStatusRunning: // 当前是开机，目标则是关机
			return c.stopInstance(ctx, plan)
		default:
			// 查当前状态，并等待
			time.Sleep(30 * time.Second)
			instance, err := c.getEbm(ctx, plan)
			if err != nil {
				return err
			}
			currentStatus = strings.ToLower(utils.SecString(instance.EbmState))
		}
	}
	return errors.New("操作机器状态失败，请检查实例状态")
}

// startInstance 启动物理机
func (c *ctyunEbm) startInstance(ctx context.Context, plan CtyunEbmConfig) (err error) {
	resp, err := c.meta.Apis.CtEbmApis.EbmStartInstanceApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmStartInstanceRequest{
		RegionID:     plan.RegionID.ValueString(),
		AzName:       plan.AzName.ValueString(),
		InstanceUUID: plan.InstanceID.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var executeSuccessFlag bool
	var status string
	retryer, _ := business.NewRetryer(time.Second*10, 60)
	retryer.Start(
		func(currentTime int) bool {
			status, err = c.getInstanceStatus(ctx, plan)
			if err != nil {
				return false
			}
			switch status {
			case business.EbmStatusRunning:
				// 执行成功
				executeSuccessFlag = true
				return false
			default:
				return true
			}
		},
	)
	if err != nil {
		return err
	}
	if !executeSuccessFlag {
		return errors.New("执行开启ebm动作时，ebm状态异常：status")
	}
	return
}

// stopInstance 关闭物理机
func (c *ctyunEbm) stopInstance(ctx context.Context, plan CtyunEbmConfig) (err error) {
	resp, err := c.meta.Apis.CtEbmApis.EbmStopInstanceApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmStopInstanceRequest{
		RegionID:     plan.RegionID.ValueString(),
		AzName:       plan.AzName.ValueString(),
		InstanceUUID: plan.InstanceID.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var executeSuccessFlag bool
	var status string
	retryer, _ := business.NewRetryer(time.Second*10, 60)
	retryer.Start(
		func(currentTime int) bool {
			status, err = c.getInstanceStatus(ctx, plan)
			if err != nil {
				return false
			}
			switch status {
			case business.EbmStatusStopped:
				// 执行成功
				executeSuccessFlag = true
				return false
			default: // 其他状态持续轮询
				return true
			}
		})
	if err != nil {
		return err
	}
	if !executeSuccessFlag {
		return errors.New("执行关闭ebm动作时，ebm状态异常，当前状态：" + status)
	}
	return
}

// getAndMerge 查询并何必
func (c *ctyunEbm) getAndMerge(ctx context.Context, cfg *CtyunEbmConfig) (err error) {
	instance, err := c.getEbm(ctx, *cfg)
	if err != nil {
		return
	}
	cfg.InstanceID = utils.SecStringValue(instance.InstanceUUID)
	cfg.RegionID = utils.SecStringValue(instance.RegionID)
	cfg.AzName = utils.SecStringValue(instance.AzName)
	cfg.DeviceType = utils.SecStringValue(instance.DeviceType)
	cfg.InstanceName = utils.SecStringValue(instance.DisplayName)
	cfg.Name = cfg.InstanceName
	cfg.Hostname = utils.SecStringValue(instance.InstanceName)
	cfg.ActualImageID = utils.SecStringValue(instance.ImageID)
	cfg.VpcID = utils.SecStringValue(instance.VpcID)
	cfg.Status = utils.SecLowerStringValue(instance.EbmState)

	eipAddress := utils.SecString(instance.PublicIP)
	cfg.EipAddress = types.StringValue(eipAddress)
	if eipAddress != "" {
		eip, err := business.NewEipService(c.meta).GetEipByAddress(ctx, eipAddress, cfg.RegionID.ValueString())
		if err != nil {
			return err
		}
		cfg.EipID = utils.SecStringValue(eip.ID)
	} else {
		cfg.EipID = types.StringValue("")
	}

	for _, card := range instance.Interfaces {
		master := utils.SecBoolValue(card.Master)
		if master.ValueBool() && len(card.SecurityGroups) > 0 {
			var secGroups []string
			for _, g := range card.SecurityGroups {
				secGroups = append(secGroups, utils.SecString(g.SecurityGroupID))
			}
			cfg.SecurityGroupIDs, _ = types.SetValueFrom(ctx, types.StringType, secGroups)
		}
		if master.ValueBool() {
			cfg.InterfaceID = utils.SecStringValue(card.InterfaceUUID)
			cfg.PortID = utils.SecStringValue(card.PortUUID)
			cfg.FixedIP = utils.SecStringValue(card.Ipv4)
			cfg.SubnetID = utils.SecStringValue(card.SubnetUUID)
		}
	}

	cfg.SystemDiskID = types.StringValue("")
	for _, diskId := range instance.AttachedVolumes {
		diskInfo, err := business.NewEbsService(c.meta).GetEbsInfo(ctx, *diskId, cfg.RegionID.ValueString())
		if err != nil {
			return err
		}
		if diskInfo.IsSystemVolume {
			cfg.SystemDiskType = types.StringValue(strings.ToLower(diskInfo.DiskType))
			cfg.SystemDiskSize = types.Int32Value(int32(diskInfo.DiskSize))
			cfg.SystemDiskID = types.StringValue(diskInfo.DiskID)
		}
	}

	cfg.ID = cfg.InstanceID

	return nil
}

// getEbm 查询ebm信息
func (c *ctyunEbm) getEbm(ctx context.Context, cfg CtyunEbmConfig) (instance *ctebm.EbmDescribeInstanceV4plusReturnObjResponse, err error) {
	resp, err := c.meta.Apis.CtEbmApis.EbmDescribeInstanceV4plusApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmDescribeInstanceV4plusRequest{
		RegionID:     cfg.RegionID.ValueString(),
		InstanceUUID: cfg.InstanceID.ValueString(),
		AzName:       cfg.AzName.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	instance = resp.ReturnObj
	return
}

// acquireIdIfOrderNotFinished 重新获取id，如果前订单状态有问题需要重新轮询
// 返回值：数据是否有效
func (c *ctyunEbm) acquireAndSetIdIfOrderNotFinished(ctx context.Context, state *CtyunEbmConfig, response *resource.ReadResponse) bool {
	id := state.InstanceID.ValueString()
	masterOrderId := state.MasterOrderID.ValueString()
	if id != "" {
		// 数据是完整的，无需处理
		return true
	}
	if state.MasterOrderID.ValueString() == "" {
		// 没有受理的订购单id，数据是不可恢复的，直接把当前状态移除并且返回
		response.State.RemoveResource(ctx)
		return false
	}
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	resp, err := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId)
	if err != nil || len(resp.Uuid) == 0 {
		// 报错了，或者受理没有返回数据的情况，那么意思是这个单子并没有开通出来，此时数据无法恢复
		response.State.RemoveResource(ctx)
		return false
	}

	// 成功把id恢复出来
	state.InstanceID = types.StringValue(resp.Uuid[0])
	response.State.Set(ctx, state)
	return true
}

// updateInstanceName 更新实例名称
func (c *ctyunEbm) updateInstanceName(ctx context.Context, state CtyunEbmConfig, plan CtyunEbmConfig) (err error) {
	// 判断名字是否相同
	if plan.InstanceName.Equal(state.InstanceName) {
		return
	}

	name := plan.InstanceName.ValueString()
	resp, err := c.meta.Apis.CtEbmApis.EbmUpdateInstanceApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmUpdateInstanceRequest{
		RegionID:     state.RegionID.ValueString(),
		AzName:       state.AzName.ValueString(),
		DisplayName:  &name,
		InstanceUUID: state.InstanceID.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
	}
	return
}

// updatePasswordOrHostname 修改密码或主机名
func (c *ctyunEbm) updatePasswordOrHostname(ctx context.Context, state CtyunEbmConfig, plan CtyunEbmConfig) (err error) {
	if state.Password.Equal(plan.Password) && state.Hostname.Equal(plan.Hostname) {
		return
	}
	// 修改前需要检查机器状态是否是关机
	err = c.checkBeforeUpdatePasswordOrHostname(ctx, state)
	if err != nil {
		return
	}
	// 修改密码
	err = c.updatePassword(ctx, state, plan)
	if err != nil {
		return
	}
	// 修改主机名
	err = c.updateHostname(ctx, state, plan)
	if err != nil {
		return
	}

	return
}

// updatePassword 修改密码
func (c *ctyunEbm) updatePassword(ctx context.Context, state CtyunEbmConfig, plan CtyunEbmConfig) (err error) {
	if state.Password.Equal(plan.Password) {
		return
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmResetPasswordApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmResetPasswordRequest{
		RegionID:     state.RegionID.ValueString(),
		AzName:       state.AzName.ValueString(),
		InstanceUUID: state.InstanceID.ValueString(),
		NewPassword:  plan.Password.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
	}
	// 通过状态检查是否修改完成
	err = c.checkAfterUpdatePasswordOrHostname(ctx, state)
	if err != nil {
		return
	}
	// 关机
	err = c.stopInstance(ctx, state)
	if err != nil {
		return
	}
	return
}

// updatePassword 修改主机名
func (c *ctyunEbm) updateHostname(ctx context.Context, state CtyunEbmConfig, plan CtyunEbmConfig) (err error) {
	if state.Hostname.Equal(plan.Hostname) {
		return
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmResetHostnameApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmResetHostnameRequest{
		RegionID:     state.RegionID.ValueString(),
		AzName:       state.AzName.ValueString(),
		InstanceUUID: state.InstanceID.ValueString(),
		Hostname:     plan.Hostname.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
	}
	// 通过状态检查是否修改完成
	err = c.checkAfterUpdatePasswordOrHostname(ctx, state)
	if err != nil {
		return
	}
	// 关机
	err = c.stopInstance(ctx, state)
	if err != nil {
		return
	}
	return
}

// checkBeforeUpdatePasswordOrHostname 修改密码或主机名前对机器状态做检查
func (c *ctyunEbm) checkBeforeUpdatePasswordOrHostname(ctx context.Context, state CtyunEbmConfig) error {
	var executeSuccessFlag bool
	var status string
	var err error
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			status, err = c.getInstanceStatus(ctx, state)
			if err != nil {
				return false
			}
			switch status {
			case business.EbmStatusStopping, business.EbmStatusResettingPassword, business.EbmStatusResettingHostname:
				return true
			case business.EbmStatusStopped:
				executeSuccessFlag = true
				return false
			default:
				return false
			}
		})
	if err != nil {
		return err
	}
	if !executeSuccessFlag {
		return errors.New("修改物理机密码或hostname失败，请确认物理机状态，修改密码或hostname必须先关机，当前状态：" + status)
	}
	return nil
}

// checkAfterUpdatePasswordOrHostname 修改后检查机器状态
func (c *ctyunEbm) checkAfterUpdatePasswordOrHostname(ctx context.Context, state CtyunEbmConfig) error {
	var executeSuccessFlag bool
	var status string
	var err error
	var cnt int
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			status, err = c.getInstanceStatus(ctx, state)
			if err != nil {
				return false
			}
			switch status {
			case business.EbmStatusResettingPassword, business.EbmStatusResettingHostname:
				return true
			case business.EbmStatusRunning:
				executeSuccessFlag = true
				return false
			case business.EbmStatusStopped:
				cnt++
				if cnt > 3 {
					return false
				}
				return true
			default:
				return false
			}
		})
	if err != nil {
		return err
	}
	if !executeSuccessFlag {
		return errors.New("修改物理机密码或hostname后，检查失败，请确认物理机状态：" + status)
	}
	return nil
}

// getInstanceStatus 查询物理机状态
func (c *ctyunEbm) getInstanceStatus(ctx context.Context, state CtyunEbmConfig) (status string, err error) {
	return business.NewEbmService(c.meta).GetEbmStatus(
		ctx,
		state.InstanceID.ValueString(),
		state.RegionID.ValueString(),
		state.AzName.ValueString(),
	)
}

// delete 删除物理机
func (c *ctyunEbm) delete(ctx context.Context, state CtyunEbmConfig) (err error) {
	resp, err := c.meta.Apis.CtEbmApis.EbmDeleteInstanceApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmDeleteInstanceRequest{
		RegionID:     state.RegionID.ValueString(),
		AzName:       state.AzName.ValueString(),
		InstanceUUID: state.InstanceID.ValueString(),
		ClientToken:  uuid.NewString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	err = helper.RefundLoop(ctx, c.meta.Credential, *resp.ReturnObj.MasterOrderID)
	if err != nil {
		return
	}
	return
}

// destroy 销毁包周期
func (c *ctyunEbm) destroy(ctx context.Context, state CtyunEbmConfig) (err error) {
	if state.CycleType.ValueString() == business.OnDemandCycleType {
		return nil
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmDestroyInstanceApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmDestroyInstanceRequest{
		RegionID:     state.RegionID.ValueString(),
		AzName:       state.AzName.ValueString(),
		InstanceUUID: state.InstanceID.ValueString(),
		ClientToken:  uuid.NewString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	err = helper.RefundLoop(ctx, c.meta.Credential, *resp.ReturnObj.MasterOrderID)
	if err != nil {
		return
	}
	return
}

// checkAfterDelete 删除后检查
func (c *ctyunEbm) checkAfterDelete(ctx context.Context, state CtyunEbmConfig) (err error) {
	portID := state.PortID.ValueString()
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	var exist bool
	retryer.Start(
		func(currentTime int) bool {
			exist, err = business.NewPortService(c.meta).Exist(ctx, portID, state.RegionID.ValueString())
			if err != nil {
				return false
			}
			if !exist {
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if exist {
		err = fmt.Errorf("裸金属 %s 的主网卡 %s 残留", state.InstanceID.ValueString(), portID)
	}
	return
}
