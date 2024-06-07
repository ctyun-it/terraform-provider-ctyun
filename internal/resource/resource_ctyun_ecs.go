package resource

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strconv"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctimage"
	defaults2 "terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
	"terraform-provider-ctyun/internal/utils"
	"time"
)

func NewCtyunEcs() resource.Resource {
	return &ctyunEcs{}
}

type ctyunEcs struct {
	meta                 *common.CtyunMetadata
	ecsService           *business.EcsService
	ebsService           *business.EbsService
	securityGroupService *business.SecurityGroupService
	keyPairService       *business.KeyPairService
	imageService         *business.ImageService
	vpcService           *business.VpcService
}

func (c *ctyunEcs) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ecs"
}

func (c *ctyunEcs) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026730**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "id",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "云主机名称，不可以使用已存在的云主机名称，长度2-63，允许使用大小写字符和数字和-，不能以-开始和结尾以及连续出现，不能使用纯数字",
				Validators: []validator.String{
					validator2.EcsName(),
				},
			},
			"flavor_id": schema.StringAttribute{
				Required:    true,
				Description: "规格id，请用ctyun_ecs_flavors查询具体id",
			},
			"image_id": schema.StringAttribute{
				Required:    true,
				Description: "镜像id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"system_disk_type": schema.StringAttribute{
				Required:    true,
				Description: "系统盘类型，sata：普通IO，sas：高IO，ssd：超高IO，ssd-genric：通用型SSD，fast-ssd：极速型SSD",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.EbsDiskTypes...),
				},
			},
			"system_disk_size": schema.Int64Attribute{
				Required:    true,
				Description: "系统盘大小，单位为G，取值范围：[40, 32768]",
				Validators: []validator.Int64{
					int64validator.Between(40, 32768),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "虚拟私有云id，在多可用区类型资源池下，vpcID通常为“vpc-”开头，非多可用区类型资源池vpcID为uuid格式",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "主网卡的子网id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"fixed_ip": schema.StringAttribute{
				Computed:    true,
				Description: "主网卡固定子网的ip地址，不填则自动返回加入子网后的ip地址",
				Validators: []validator.String{
					validator2.Ip(),
				},
			},
			"security_group_ids": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "安全组id列表，在多可用区类型资源池下，安全组ID通常以“sg-”开头，非多可用区类型资源池安全组ID为uuid格式；默认使用默认安全组，无默认安全组情况下请填写该参数",
			},
			"key_pair_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "密钥对名称",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 63),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$"), "不满足密钥对名称要求"),
				},
				Default: stringdefault.StaticString(""),
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Description: "用户密码，满足以下规则：长度在8～30个字符；必须包含大写字母、小写字母、数字以及特殊符号中的三项；特殊符号可选：()`~!@#$%^&*_-+=|{}[]:;'<>,.?/\\且不能以斜线号/开头",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(8, 30),
					validator2.EcsPassword(),
				},
				Sensitive: true,
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围：month：按月，year：按年、on_demand：按需。当此值为month或者year时，cycle_count为必填",
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypes...),
				},
			},
			"cycle_count": schema.Int64Attribute{
				Optional:    true,
				Description: "订购时长，该参数在cycle_type为month或year时才生效，当cycleType=month，支持续订1-11个月；当cycleType=year，支持续订1-5年",
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
				Description: "是否自动续订，true：自动续订，false：不自动续订；注意：此参数在包周期情况下才有效；当为包周期时此值默认为true",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Default: booldefault.StaticBool(true),
				Validators: []validator.Bool{
					validator2.ConflictsWithEqualBool(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
				},
			},
			"default_security_group_id": schema.StringAttribute{
				Computed:    true,
				Description: "默认加入安全组id",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "云主机状态，取值范围：backingup: 备份中，creating: 创建中，expired: 已到期，freezing: 冻结中，rebuild: 重装，restarting: 重启中，running: 运行中，starting: 开机中，stopped: 已关机，stopping: 关机中，error: 错误，snapshotting: 快照创建中，unsubscribed: 包周期已退订，unsubscribing: 包周期退订中",
			},
			"expire_time": schema.StringAttribute{
				Computed:    true,
				Description: "到期时间",
			},
			"system_disk_id": schema.StringAttribute{
				Computed:    true,
				Description: "系统盘的id",
			},
			// "user_data": schema.StringAttribute{
			// 	Optional:    true,
			// 	Computed:    true,
			// 	Description: "用户自定义数据，需要以Base64方式编码，Base64编码后的长度限制为1-16384字符。注：非多可用区类型资源池暂不支持该参数",
			// 	Validators: []validator.String{
			// 		stringvalidator.UTF8LengthBetween(1, 16384),
			// 	},
			// },
			// "monitor_service": schema.BoolAttribute{
			// 	Optional:    true,
			// 	Computed:    true,
			// 	Description: "监控参数，支持通过该参数指定云主机在创建后是否开启详细监控，false：不开启，true：开启。若指定该参数为true或不指定该参数，云主机内默认开启最新详细监控服务。若指定该参数为false，默认公共镜像不开启最新监控服务；私有镜像使用镜像中保留的监控服务。说明：仅部分资源池支持",
			// 	Default:     booldefault.StaticBool(true),
			// },
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订购的受理单id",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目id，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区id，如果不填则默认使用provider ctyun中的az_name或环境变量中的CTYUN_AZ_NAME",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraAzName, false),
			},
		},
	}
}

func (c *ctyunEcs) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunEcsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 校验创建动作的前置条件
	err := c.checkCreate(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// 实际创建
	err = c.createInstance(ctx, &plan)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// 创建机器后默认为启动状态，可以直接绑定keypair，注意如果绑定失败了，这里不能抛出错误，因为实际的云主机已经创建出来了
	_ = c.bindKeyPair(ctx, plan)

	// 查询信息
	instance, err := c.getAndMergeEcs(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}

	// 修复bug，因为创建的时候，后端会将实例自动加入到到某个特定的安全组中，如果直接返回会导致terraform报错，因此要把多余的安全组给过滤掉
	instance.DefaultSecurityGroupId = c.getAndRemoveSecurityGroups(ctx, plan, instance)

	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEcs) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunEcsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	if !c.acquireAndSetIdIfOrderNotFinished(ctx, &state, response) {
		return
	}
	instance, err := c.getAndMergeEcs(ctx, state)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEcs) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state CtyunEcsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	var plan CtyunEcsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 修改基础信息
	err := c.updateInstanceInfo(ctx, state, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// 修改硬盘大小
	err2 := c.updateSystemDisk(ctx, state, plan)
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}

	// 修改密码
	err2 = c.updatePassword(ctx, state, plan)
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}

	// 修改规格
	err2 = c.updateFlavor(ctx, state, plan)
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}

	// 按需转包，包转按需
	err2 = c.changePayType(ctx, state, plan)
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}

	// 更新安全组
	err2 = c.updateSecurityGroup(ctx, state, plan)
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}

	// 更新密钥
	err2 = c.updateKeyPair(ctx, state, plan)
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}

	// 反查信息
	instance, err2 := c.getAndMergeEcs(ctx, state)
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}
	instance.Password = plan.Password
	instance.CycleType = plan.CycleType
	instance.CycleCount = plan.CycleCount
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEcs) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunEcsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 先关机，因为销毁是默认用户意识到资料销毁的动作，所以直接关机是ok的
	err := c.closeInstance(ctx, state.Id.ValueString(), state.RegionId.ValueString())
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// 解绑对应的安全设组
	_ = c.leaveSecurityGroups(ctx, state)

	// 退订操作
	resp, err := c.meta.Apis.CtEcsApis.EcsUnsubscribeInstanceApi.Do(ctx, c.meta.Credential, &ctecs.EcsUnsubscribeInstanceRequest{
		RegionId:    state.RegionId.ValueString(),
		InstanceId:  state.Id.ValueString(),
		ClientToken: uuid.NewString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	err = helper.RefundLoop(ctx, c.meta.Credential, resp.MasterOrderId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

func (c *ctyunEcs) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ecsService = business.NewEcsService(meta)
	c.ebsService = business.NewEbsService(meta)
	c.securityGroupService = business.NewSecurityGroupService(meta)
	c.keyPairService = business.NewKeyPairService(meta)
	c.imageService = business.NewImageService(meta)
	c.vpcService = business.NewVpcService(meta)
}

// createInstance 创建实例
func (c *ctyunEcs) createInstance(ctx context.Context, plan *CtyunEcsConfig) error {
	// 镜像类型参数
	imageResponse, err := c.meta.Apis.CtImageApis.ImageDetailApi.Do(ctx, c.meta.Credential, &ctimage.ImageDetailRequest{
		RegionId: plan.RegionId.ValueString(),
		ImageId:  plan.ImageId.ValueString(),
	})
	if err != nil {
		return err
	}
	imageVisibility, err2 := business.ImageVisibilityMap.FromOriginalScene(imageResponse.Images[0].Visibility, business.ImageVisibilityMapScene1)
	if err2 != nil {
		return err2
	}

	// 是否按需参数
	onDemand := plan.CycleType.ValueString() == business.OrderCycleTypeOnDemand

	// 订购周期类型参数
	cycleType, err2 := business.OrderCycleTypeMap.FromOriginalScene(plan.CycleType.ValueString(), business.OrderCycleTypeMapScene1)
	if err2 != nil {
		return err2
	}

	// 自定续订参数
	autoRenewStatus := 0
	if plan.AutoRenew.ValueBool() {
		autoRenewStatus = 1
	}

	// 系统盘类型参数
	diskType, err2 := business.EbsDiskTypeMap.FromOriginalScene(plan.SystemDiskType.ValueString(), business.EbsDiskTypeMapScene1)
	if err2 != nil {
		return err2
	}

	var securityGroupIds []types.String
	var sgIds []string
	plan.SecurityGroupIds.ElementsAs(ctx, &securityGroupIds, true)
	for _, id := range securityGroupIds {
		sgIds = append(sgIds, id.ValueString())
	}

	regionId := plan.RegionId.ValueString()
	azName := plan.AzName.ValueString()
	projectId := plan.ProjectId.ValueString()

	// 创建ecs实例
	resp, err2 := c.meta.Apis.CtEcsApis.EcsCreateInstanceApi.Do(ctx, c.meta.Credential, &ctecs.EcsCreateInstanceRequest{
		RegionId:        regionId,
		AzName:          azName,
		ProjectId:       projectId,
		ClientToken:     uuid.NewString(),
		InstanceName:    plan.Name.ValueString(),
		DisplayName:     plan.Name.ValueString(),
		FlavorId:        plan.FlavorId.ValueString(),
		ImageType:       imageVisibility.(int),
		ImageId:         plan.ImageId.ValueString(),
		BootDiskType:    diskType.(string),
		BootDiskSize:    int(plan.SystemDiskSize.ValueInt64()),
		VpcId:           plan.VpcId.ValueString(),
		OnDemand:        onDemand,
		ExtIp:           "0",
		UserPassword:    plan.Password.ValueString(),
		CycleCount:      int(plan.CycleCount.ValueInt64()),
		CycleType:       cycleType.(string),
		AutoRenewStatus: autoRenewStatus,
		NetworkCardList: []ctecs.EcsCreateInstanceNetworkCardListRequest{
			{
				SubnetId: plan.SubnetId.ValueString(),
				FixedIp:  plan.FixedIp.ValueString(),
				IsMaster: true,
			},
		},
		SecGroupList: sgIds,
		// UserData:       plan.UserData.ValueString(),
		// MonitorService: plan.MonitorService.ValueBool(),
	})
	if err2 != nil {
		return err2
	}

	// 先设置重要的属性
	plan.RegionId = types.StringValue(regionId)
	plan.AzName = types.StringValue(azName)
	plan.ProjectId = types.StringValue(projectId)
	plan.MasterOrderId = types.StringValue(resp.MasterOrderId)

	// 轮询订单状态
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	loop, err2 := helper.OrderLoop(ctx, c.meta.Credential, resp.MasterOrderId)
	if err2 != nil {
		return err2
	}

	// 最后设置id
	id := loop.Uuid[0]
	plan.Id = types.StringValue(id)

	// 等待云主机状态为运行中的状态
	_ = c.waitInstanceStatusFor(ctx, id, regionId, business.EcsStatusRunning)
	return nil
}

// updateInstanceInfo 更新主机的部分信息
func (c *ctyunEcs) updateInstanceInfo(ctx context.Context, state CtyunEcsConfig, plan CtyunEcsConfig) error {
	if state.Name.Equal(plan.Name) {
		return nil
	}
	_, err := c.meta.Apis.CtEcsApis.EcsBatchUpdateInstancesApi.Do(ctx, c.meta.Credential, &ctecs.EcsBatchUpdateInstancesRequest{
		RegionId: state.RegionId.ValueString(),
		UpdateInfo: []ctecs.EcsBatchUpdateInstancesUpdateInfoRequest{
			{
				InstanceId:  state.Id.ValueString(),
				DisplayName: plan.Name.ValueString(),
			},
		},
	})
	return err
}

// checkInstanceStatus 校验云主机状态必须为目标状态
func (c *ctyunEcs) checkInstanceStatus(ctx context.Context, id string, regionId, targetStatus string) bool {
	currentStatus, err := c.getInstanceStatus(ctx, id, regionId)
	return err == nil && currentStatus == targetStatus
}

// changePayType 变更付费模式
func (c *ctyunEcs) changePayType(ctx context.Context, state CtyunEcsConfig, plan CtyunEcsConfig) error {
	if plan.CycleType.Equal(state.CycleType) {
		return nil
	}
	// 变更付费模式前必须为关机状态
	if !c.checkInstanceStatus(ctx, state.Id.ValueString(), state.RegionId.ValueString(), business.EcsStatusStopped) {
		return errors.New("变更云主机付费模式，请先将云主机关机")
	}

	cycleType := plan.CycleType.ValueString()
	if cycleType == business.OrderCycleTypeMonth || cycleType == business.OrderCycleTypeYear {
		// 按需转包
		err := c.onDemandToCycle(ctx, state.Id.ValueString(), state.RegionId.ValueString(), plan.CycleType.ValueString(), int(plan.CycleCount.ValueInt64()))
		if err != nil {
			return err
		}
	} else if cycleType == business.OrderCycleTypeOnDemand {
		// 包转按需
		err := c.cycleToOnDemand(ctx, state.Id.ValueString(), state.RegionId.ValueString())
		if err != nil {
			return err
		}
	}
	return nil
}

// cycleToOnDemand 包转按需
func (c *ctyunEcs) cycleToOnDemand(ctx context.Context, id, regionId string) (err error) {
	// 首先进行对主机实例进行打标处理
	tagResp, err := c.meta.Apis.CtEcsApis.EcsTagOnDemandApi.Do(ctx, c.meta.Credential, &ctecs.EcsTagOnDemandRequest{
		ClientToken: uuid.NewString(),
		RegionId:    regionId,
		InstanceIds: []string{id},
	})
	if err != nil {
		return err
	}

	// 轮询订单打标状态
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	_, err = helper.OrderLoop(ctx, c.meta.Credential, tagResp.OrderInfo[0].OrderId)
	if err != nil {
		return err
	}

	terminateCycleResp, err := c.meta.Apis.CtEcsApis.EcsTerminateCycleApi.Do(ctx, c.meta.Credential, &ctecs.EcsTerminateCycleRequest{
		ClientToken: uuid.NewString(),
		RegionId:    regionId,
		InstanceIds: []string{id},
	})
	if err != nil {
		return err
	}

	// 轮询包周期终止订单状态
	_, err2 := helper.OrderLoop(ctx, c.meta.Credential, terminateCycleResp.OrderInfo[0].OrderId)
	return err2
}

// onDemandToCycle 按需转包
func (c *ctyunEcs) onDemandToCycle(ctx context.Context, id, regionId, cycleType string, cycleCount int) error {
	// 按需转包
	cycleTypeParam, err := business.OrderCycleTypeMap.FromOriginalScene(cycleType, business.OrderCycleTypeMapScene1)
	if err != nil {
		return err
	}

	resp, err := c.meta.Apis.CtEcsApis.EcsChangeToCycleApi.Do(ctx, c.meta.Credential, &ctecs.EcsChangeToCycleRequest{
		ClientToken: uuid.NewString(),
		RegionId:    regionId,
		InstanceIds: []string{id},
		CycleType:   cycleTypeParam.(string),
		CycleCount:  cycleCount,
	})
	if err != nil {
		return err
	}

	// 轮询订单状态
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	_, err2 := helper.OrderLoop(ctx, c.meta.Credential, resp.OrderInfo[0].OrderId)
	return err2
}

// handleInstance 操作机器
func (c *ctyunEcs) handleInstance(ctx context.Context, id, regionId, targetStatus string) error {
	switch targetStatus {
	case business.EcsStatusStopped:
		return c.closeInstance(ctx, id, regionId)
	case business.EcsStatusRunning:
		return c.startInstance(ctx, id, regionId)
	}
	return errors.New("操作机器状态失败，请检查实例状态")
}

// closeInstance 关机
func (c *ctyunEcs) closeInstance(ctx context.Context, id, regionId string) error {
	currentStatus, err := c.getInstanceStatus(ctx, id, regionId)
	if err != nil {
		return err
	}
	// 已经是关机的状态了
	if currentStatus == business.EcsStatusStopped {
		return nil
	}
	if currentStatus != business.EcsStatusRunning {
		return errors.New("当前机器状态异常，请稍后重试或在控制台进行操作")
	}

	executeSuccessFlag := false
	// 关机的情况
	_, err2 := c.meta.Apis.CtEcsApis.EcsStopInstanceApi.Do(ctx, c.meta.Credential, &ctecs.EcsStopInstanceRequest{
		RegionId:   regionId,
		InstanceId: id,
		Force:      false,
	})
	if err2 != nil {
		// 已经是开机的情况，直接返回
		if err2.ErrorCode() == common.EcsInstanceStatusNotRunning {
			return nil
		}
		return err2
	}

	// 轮询关机状态
	retryer, _ := business.NewRetryer(time.Second*5, 20)
	retryer.Start(
		func(currentTime int) bool {
			status, err3 := c.getInstanceStatus(ctx, id, regionId)
			if err3 != nil {
				return false
			}
			switch status {
			case business.EcsStatusStopping:
				// 执行中
				return true
			case business.EcsStatusStopped:
				// 执行成功
				executeSuccessFlag = true
				return false
			default:
				// 默认为执行失败
				return false
			}
		},
	)

	if !executeSuccessFlag {
		return errors.New("执行关闭ecs动作时，查询ecs状态异常")
	}
	return nil
}

// startInstance 开机
func (c *ctyunEcs) startInstance(ctx context.Context, id, regionId string) error {
	currentStatus, err := c.getInstanceStatus(ctx, id, regionId)
	if err != nil {
		return err
	}
	// 已经是开机的状态了
	if currentStatus == business.EcsStatusRunning {
		return nil
	}
	if currentStatus != business.EcsStatusStopped {
		return errors.New("当前机器状态异常，请稍后重试或在控制台进行操作")
	}

	executeSuccessFlag := false
	// 开机的情况
	_, err2 := c.meta.Apis.CtEcsApis.EcsStartInstanceApi.Do(ctx, c.meta.Credential, &ctecs.EcsStartInstanceRequest{
		RegionId:   regionId,
		InstanceId: id,
		Force:      false,
	})
	if err2 != nil {
		// 已经是关机的情况，直接返回
		if err2.ErrorCode() == common.EcsInstanceStatusNotStopped {
			return nil
		}
		return err2
	}

	// 轮询开机状态
	retryer, _ := business.NewRetryer(time.Second*5, 20)
	retryer.Start(
		func(currentTime int) bool {
			status, err3 := c.getInstanceStatus(ctx, id, regionId)
			if err3 != nil {
				return false
			}
			switch status {
			case business.EcsStatusStarting:
				// 执行中
				return true
			case business.EcsStatusRunning:
				// 执行成功
				executeSuccessFlag = true
				return false
			default:
				// 默认为执行失败
				return false
			}
		},
	)

	if !executeSuccessFlag {
		return errors.New("执行开启ecs动作时，查询ecs状态异常")
	}
	return nil
}

// getInstanceStatus 获取云主机状态信息
func (c *ctyunEcs) getInstanceStatus(ctx context.Context, id, regionId string) (string, error) {
	resp, err := c.meta.Apis.CtEcsApis.EcsInstanceDetailsApi.Do(ctx, c.meta.Credential, &ctecs.EcsInstanceDetailsRequest{
		RegionId:   regionId,
		InstanceId: id,
	})
	if err != nil {
		return "", err
	}
	return resp.InstanceStatus, err
}

// getAndRemoveSecurityGroups 获取并删除对应安全组
func (c *ctyunEcs) getAndRemoveSecurityGroups(ctx context.Context, plan CtyunEcsConfig, target *CtyunEcsConfig) types.String {
	var securityGroupIds []types.String
	plan.SecurityGroupIds.ElementsAs(ctx, &securityGroupIds, true)
	mapping := make(map[string]struct{})
	for _, id := range securityGroupIds {
		mapping[id.ValueString()] = struct{}{}
	}

	newSecurityGroupIds := []types.String{}
	var targetSecurityGroupIds []types.String
	target.SecurityGroupIds.ElementsAs(ctx, &targetSecurityGroupIds, true)
	var defaultSecurityGroupId types.String
	for _, id := range targetSecurityGroupIds {
		_, ok := mapping[id.ValueString()]
		if ok {
			newSecurityGroupIds = append(newSecurityGroupIds, id)
		} else {
			defaultSecurityGroupId = id
		}
	}
	sgs, _ := types.SetValueFrom(ctx, types.StringType, newSecurityGroupIds)
	target.SecurityGroupIds = sgs
	return defaultSecurityGroupId
}

// joinSecurityGroups 加入安全组
func (c *ctyunEcs) joinSecurityGroups(ctx context.Context, plan CtyunEcsConfig) error {
	var securityGroupIds []types.String
	plan.SecurityGroupIds.ElementsAs(ctx, &securityGroupIds, true)
	if len(securityGroupIds) == 0 {
		return nil
	}
	for _, id := range securityGroupIds {
		_, err := c.meta.Apis.CtEcsApis.EcsJoinSecurityGroupApi.Do(ctx, c.meta.Credential, &ctecs.EcsJoinSecurityGroupRequest{
			RegionId:        plan.RegionId.ValueString(),
			SecurityGroupId: id.ValueString(),
			InstanceId:      plan.Id.ValueString(),
			Action:          "joinSecurityGroup",
		})
		if err != nil {
			return errors.New("加入安全组：" + id.ValueString() + "失败：" + err.Error())
		}
	}
	return nil
}

// leaveSecurityGroups 离开安全组
func (c *ctyunEcs) leaveSecurityGroups(ctx context.Context, state CtyunEcsConfig) error {
	var securityGroupIds []types.String
	state.SecurityGroupIds.ElementsAs(ctx, &securityGroupIds, true)
	if len(securityGroupIds) == 0 {
		return nil
	}
	for _, id := range securityGroupIds {
		_, err := c.meta.Apis.CtEcsApis.EcsLeaveSecurityGroupApi.Do(ctx, c.meta.Credential, &ctecs.EcsLeaveSecurityGroupRequest{
			RegionId:        state.RegionId.ValueString(),
			SecurityGroupId: id.ValueString(),
			InstanceId:      state.Id.ValueString(),
		})
		if err != nil {
			return errors.New("离开安全组：" + id.ValueString() + "失败：" + err.Error())
		}
	}
	return nil
}

// bindKeyPair 绑定密钥对
func (c *ctyunEcs) waitInstanceStatusFor(ctx context.Context, id, regionId, statusFor string) error {
	retryer, _ := business.NewRetryer(time.Second*5, 12)
	result := retryer.Start(func(currentTime int) bool {
		return !c.checkInstanceStatus(ctx, id, regionId, statusFor)
	})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("查询等待云主机状态：" + statusFor + "超时")
	}
	return nil
}

// bindKeyPair 绑定密钥对
func (c *ctyunEcs) bindKeyPair(ctx context.Context, plan CtyunEcsConfig) error {
	if plan.KeyPairName.ValueString() == "" {
		return nil
	}
	// 绑定密钥对
	_, err := c.meta.Apis.CtEcsApis.KeypairAttachApi.Do(ctx, c.meta.Credential, &ctecs.KeypairAttachRequest{
		RegionId:    plan.RegionId.ValueString(),
		KeyPairName: plan.KeyPairName.ValueString(),
		InstanceId:  plan.Id.ValueString(),
	})
	if err != nil {
		return err
	}
	return nil
}

// updateFlavor 更新云主机实例规格
func (c *ctyunEcs) updateFlavor(ctx context.Context, state CtyunEcsConfig, plan CtyunEcsConfig) error {
	if state.FlavorId.Equal(plan.FlavorId) {
		return nil
	}

	// 更新云主机前必须为关机状态
	if !c.checkInstanceStatus(ctx, state.Id.ValueString(), state.RegionId.ValueString(), business.EcsStatusStopped) {
		return errors.New("变更云主机配置规格，请先将云主机关机")
	}

	// 校验规格必须存在
	err := c.ecsService.FlavorMustExist(ctx, plan.FlavorId.ValueString(), state.RegionId.ValueString(), state.AzName.ValueString())
	if err != nil {
		return err
	}

	// 更新云主机规格
	resp, err := c.meta.Apis.CtEcsApis.EcsUpdateFlavorSpecApi.Do(ctx, c.meta.Credential, &ctecs.EcsUpdateFlavorSpecRequest{
		RegionId:    state.RegionId.ValueString(),
		ClientToken: uuid.NewString(),
		InstanceId:  state.Id.ValueString(),
		FlavorId:    plan.FlavorId.ValueString(),
	})
	if err != nil {
		return err
	}

	// 轮询订单
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	_, e := helper.OrderLoop(ctx, c.meta.Credential, resp.MasterOrderId)
	if e != nil {
		return e
	}

	return nil
}

// updateKeyPair 更新密钥对
func (c *ctyunEcs) updateKeyPair(ctx context.Context, state CtyunEcsConfig, plan CtyunEcsConfig) error {
	if state.KeyPairName.Equal(plan.KeyPairName) {
		return nil
	}
	// 变更密钥对前必须为开机状态
	if !c.checkInstanceStatus(ctx, state.Id.ValueString(), state.RegionId.ValueString(), business.EcsStatusRunning) {
		return errors.New("变更云主机密钥对，请先将云主机开机")
	}
	// 先校验变更的密钥对必须存在
	if !plan.KeyPairName.Equal(state.KeyPairName) && plan.KeyPairName.ValueString() != "" {
		err := c.keyPairService.MustExist(ctx, plan.KeyPairName.ValueString(), state.RegionId.ValueString(), state.ProjectId.ValueString())
		if err != nil {
			return err
		}
	}

	if state.KeyPairName.ValueString() != "" {
		// 解绑旧的密钥对
		_, err := c.meta.Apis.CtEcsApis.KeypairDetachApi.Do(ctx, c.meta.Credential, &ctecs.KeypairDetachRequest{
			RegionId:    state.RegionId.ValueString(),
			KeyPairName: state.KeyPairName.ValueString(),
			InstanceId:  state.Id.ValueString(),
		})
		if err != nil {
			return err
		}
	}
	if plan.KeyPairName.ValueString() != "" {
		// 绑定新的密钥对
		_, err := c.meta.Apis.CtEcsApis.KeypairAttachApi.Do(ctx, c.meta.Credential, &ctecs.KeypairAttachRequest{
			RegionId:    state.RegionId.ValueString(),
			KeyPairName: plan.KeyPairName.ValueString(),
			InstanceId:  state.Id.ValueString(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// updateSecurityGroup 更新安全组
func (c *ctyunEcs) updateSecurityGroup(ctx context.Context, state CtyunEcsConfig, plan CtyunEcsConfig) error {
	var mapping = make(map[string]struct{})
	var securityGroups []types.String
	state.SecurityGroupIds.ElementsAs(ctx, &securityGroups, true)
	for _, group := range securityGroups {
		mapping[group.ValueString()] = struct{}{}
	}

	// 过滤出需要新加入的安全组id
	var joinGroupIds []string
	plan.SecurityGroupIds.ElementsAs(ctx, &securityGroups, true)
	for _, group := range securityGroups {
		groupStr := group.ValueString()
		_, ok := mapping[groupStr]
		if ok {
			delete(mapping, groupStr)
		} else {
			// 先校验安全组必须存在
			err := c.securityGroupService.MustExist(ctx, groupStr, state.RegionId.ValueString())
			if err != nil {
				return err
			}
			joinGroupIds = append(joinGroupIds, groupStr)
		}
	}

	// 实际加入安全组
	for _, id := range joinGroupIds {
		_, err := c.meta.Apis.CtEcsApis.EcsJoinSecurityGroupApi.Do(ctx, c.meta.Credential, &ctecs.EcsJoinSecurityGroupRequest{
			RegionId:        state.RegionId.ValueString(),
			SecurityGroupId: id,
			InstanceId:      state.Id.ValueString(),
			Action:          "joinSecurityGroup",
		})
		if err != nil {
			return err
		}
	}

	// 剩余的是离开的安全组
	for key := range mapping {
		_, err := c.meta.Apis.CtEcsApis.EcsLeaveSecurityGroupApi.Do(ctx, c.meta.Credential, &ctecs.EcsLeaveSecurityGroupRequest{
			RegionId:        state.RegionId.ValueString(),
			SecurityGroupId: key,
			InstanceId:      state.Id.ValueString(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// getAndMergeEcs 查询ecs
func (c *ctyunEcs) getAndMergeEcs(ctx context.Context, cfg CtyunEcsConfig) (*CtyunEcsConfig, error) {
	regionId := cfg.RegionId.ValueString()

	resp, err := c.meta.Apis.CtEcsApis.EcsInstanceDetailsApi.Do(ctx, c.meta.Credential, &ctecs.EcsInstanceDetailsRequest{
		RegionId:   regionId,
		InstanceId: cfg.Id.ValueString(),
	})
	if err != nil {
		// 实例已经被退订的情况
		if err.ErrorCode() == common.EcsInstanceNotFound {
			return nil, nil
		}
		return nil, err
	}

	// 基础信息
	cfg.Id = types.StringValue(resp.InstanceId)
	cfg.Name = types.StringValue(resp.DisplayName)
	cfg.FlavorId = types.StringValue(resp.Flavor.FlavorId)
	cfg.ImageId = types.StringValue(resp.Image.ImageId)
	cfg.VpcId = types.StringValue(resp.VpcId)
	cfg.Status = types.StringValue(resp.InstanceStatus)
	cfg.ExpireTime = types.StringValue(utils.FromRFC3339ToLocal(resp.ExpiredTime))
	cfg.KeyPairName = types.StringValue(resp.KeypairName)

	// 填充安全组信息
	sgs := []types.String{}
	for _, sg := range resp.SecGroupList {
		// 如果存在默认的用户组，要判断一下返回的是否为默认的用户组，如果是默认的就把它排除掉
		if !cfg.DefaultSecurityGroupId.IsNull() && !cfg.DefaultSecurityGroupId.IsUnknown() {
			if sg.SecurityGroupId == cfg.DefaultSecurityGroupId.ValueString() {
				continue
			}
		}
		sgs = append(sgs, types.StringValue(sg.SecurityGroupId))
	}
	securityGroupIds, _ := types.SetValueFrom(ctx, types.StringType, sgs)
	cfg.SecurityGroupIds = securityGroupIds

	// 填充主网卡信息
	for _, nc := range resp.NetworkCardList {
		if nc.IsMaster {
			cfg.SubnetId = types.StringValue(nc.SubnetId)
			cfg.FixedIp = types.StringValue(nc.Ipv4Address)
		}
	}

	// 密钥对信息
	if resp.KeypairName != "" {
		cfg.KeyPairName = types.StringValue(resp.KeypairName)
	}

	// 查询系统盘，填补其信息
	ecsVolumeResponse, err := c.meta.Apis.CtEcsApis.EcsVolumeListApi.Do(ctx, c.meta.Credential, &ctecs.EcsVolumeListRequest{
		RegionId:   regionId,
		InstanceId: cfg.Id.ValueString(),
		PageNo:     1,
		PageSize:   50,
	})
	if err != nil {
		return nil, err
	}
	var vs []ctecs.EcsVolumeListResultsResponse
	for _, v := range ecsVolumeResponse.Results {
		if v.DiskType == "系统盘" {
			vs = append(vs, v)
		}
	}
	if len(vs) != 1 {
		return nil, errors.New("查询系统盘信息发生错误，查询到系统盘数量" + strconv.Itoa(len(vs)))
	}
	result := ecsVolumeResponse.Results[0]
	diskType, err2 := business.EbsDiskTypeMap.ToOriginalScene(result.DiskDataType, business.EbsDiskTypeMapScene1)
	if err2 != nil {
		return nil, err2
	}
	cfg.SystemDiskType = types.StringValue(diskType.(string))
	cfg.SystemDiskSize = types.Int64Value(int64(result.DiskSize))
	cfg.SystemDiskId = types.StringValue(result.DiskId)

	return &cfg, nil
}

// checkCreate 校验创建动作是否能执行
func (c *ctyunEcs) checkCreate(ctx context.Context, plan CtyunEcsConfig) error {
	// 镜像必须存在
	err := c.imageService.MustExist(ctx, plan.ImageId.ValueString(), plan.RegionId.ValueString())
	if err != nil {
		return err
	}

	// vpc必须存在
	err = c.vpcService.MustExist(ctx, plan.VpcId.ValueString(), plan.RegionId.ValueString(), plan.ProjectId.ValueString())
	if err != nil {
		return err
	}

	// 安全组必须存在
	var securityGroupIds []types.String
	plan.SecurityGroupIds.ElementsAs(ctx, &securityGroupIds, true)
	for _, id := range securityGroupIds {
		err := c.securityGroupService.MustExist(ctx, id.ValueString(), plan.RegionId.ValueString())
		if err != nil {
			return err
		}
	}

	// 密钥对必须存在
	err = c.keyPairService.MustExist(ctx, plan.KeyPairName.ValueString(), plan.RegionId.ValueString(), plan.ProjectId.ValueString())
	if err != nil {
		return err
	}

	// 云主机规格必须存在
	err = c.ecsService.FlavorMustExist(ctx, plan.FlavorId.ValueString(), plan.RegionId.ValueString(), plan.AzName.ValueString())
	if err != nil {
		return err
	}

	return nil
}

// updateSystemDisk 更新系统盘
func (c *ctyunEcs) updateSystemDisk(ctx context.Context, state CtyunEcsConfig, plan CtyunEcsConfig) error {
	if state.SystemDiskSize.Equal(plan.SystemDiskSize) {
		return nil
	}
	// 先校验关机状态，注意这个动作必须让用户自我决定执行
	if !c.checkInstanceStatus(ctx, state.Id.ValueString(), state.RegionId.ValueString(), business.EcsStatusStopped) {
		return errors.New("变更云主机系统盘大小，请先将云主机关机")
	}
	return c.ebsService.UpdateSize(ctx, state.SystemDiskId.ValueString(), state.RegionId.ValueString(), int(state.SystemDiskSize.ValueInt64()), int(plan.SystemDiskSize.ValueInt64()))

}

// updatePassword 修改密码
func (c *ctyunEcs) updatePassword(ctx context.Context, state CtyunEcsConfig, plan CtyunEcsConfig) error {
	if state.Password.Equal(plan.Password) {
		return nil
	}
	// 先校验关机状态，注意这个动作必须让用户自我决定执行
	if !c.checkInstanceStatus(ctx, state.Id.ValueString(), state.RegionId.ValueString(), business.EcsStatusRunning) {
		return errors.New("修改云主机密码，请先将云主机开机")
	}
	_, err := c.meta.Apis.CtEcsApis.EcsResetPasswordApi.Do(ctx, c.meta.Credential, &ctecs.EcsResetPasswordRequest{
		RegionId:    state.RegionId.ValueString(),
		InstanceId:  state.Id.ValueString(),
		NewPassword: plan.Password.ValueString(),
	})
	return err
}

// acquireIdIfOrderNotFinished 重新获取id，如果前订单状态有问题需要重新轮询
// 返回值：数据是否有效
func (c *ctyunEcs) acquireAndSetIdIfOrderNotFinished(ctx context.Context, state *CtyunEcsConfig, response *resource.ReadResponse) bool {
	id := state.Id.ValueString()
	masterOrderId := state.MasterOrderId.ValueString()
	if id != "" {
		// 数据是完整的，无需处理
		return true
	}
	if state.MasterOrderId.ValueString() == "" {
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
	id = resp.Uuid[0]

	// 成功把id恢复出来
	state.Id = types.StringValue(id)
	response.State.Set(ctx, state)
	return true
}

type CtyunEcsConfig struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	FlavorId               types.String `tfsdk:"flavor_id"`
	ImageId                types.String `tfsdk:"image_id"`
	SystemDiskType         types.String `tfsdk:"system_disk_type"`
	SystemDiskSize         types.Int64  `tfsdk:"system_disk_size"`
	VpcId                  types.String `tfsdk:"vpc_id"`
	SecurityGroupIds       types.Set    `tfsdk:"security_group_ids"`
	KeyPairName            types.String `tfsdk:"key_pair_name"`
	Password               types.String `tfsdk:"password"`
	CycleCount             types.Int64  `tfsdk:"cycle_count"`
	CycleType              types.String `tfsdk:"cycle_type"`
	AutoRenew              types.Bool   `tfsdk:"auto_renew"`
	SubnetId               types.String `tfsdk:"subnet_id"`
	FixedIp                types.String `tfsdk:"fixed_ip"`
	DefaultSecurityGroupId types.String `tfsdk:"default_security_group_id"`
	Status                 types.String `tfsdk:"status"`
	ExpireTime             types.String `tfsdk:"expire_time"`
	SystemDiskId           types.String `tfsdk:"system_disk_id"`
	// UserData               types.String `tfsdk:"user_data"`
	// MonitorService         types.Bool   `tfsdk:"monitor_service"`
	MasterOrderId types.String `tfsdk:"master_order_id"`
	ProjectId     types.String `tfsdk:"project_id"`
	RegionId      types.String `tfsdk:"region_id"`
	AzName        types.String `tfsdk:"az_name"`
}
