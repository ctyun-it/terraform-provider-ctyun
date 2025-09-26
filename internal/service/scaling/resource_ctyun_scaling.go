package scaling

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/scaling"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
	"time"
)

type ctyunScaling struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *ctyunScaling) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_scaling_group"
}

func (c *ctyunScaling) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunScaling() resource.Resource {
	return &ctyunScaling{}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [id],[regionId],[projectId]
func (c *ctyunScaling) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var cfg CtyunScalingConfig
	var ID, regionId, projectId, vpcId string
	err = terraform_extend.Split(request.ID, &ID, &regionId, &projectId, &vpcId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	id, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return
	}
	cfg.ID = types.Int64Value(id)
	cfg.RegionID = types.StringValue(regionId)
	cfg.ProjectID = types.StringValue(projectId)
	cfg.VpcID = types.StringValue(vpcId)

	cfg.AddInstanceUUIDList = types.SetNull(types.StringType)
	cfg.RemoveInstanceUUIDList = types.SetNull(types.StringType)

	err = c.getAndMergeScaling(ctx, &cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *ctyunScaling) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：：https://www.ctyun.cn/document/10027725`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
			},
			"security_group_id_list": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "安全组ID列表。支持更新",
			},
			//"recovery_mode": schema.Int64Attribute{
			//	Required:    true,
			//	Description: "实例回收模式: 1-释放模式, 2-停机回收模式",
			//	Validators: []validator.Int64{
			//		int64validator.OneOf(1, 2),
			//	},
			//},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "伸缩组名称。支持更新。",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"health_mode": schema.StringAttribute{
				Required:    true,
				Description: "健康检查方式：server-云服务器健康检查，lb-弹性负载均衡健康检查。支持更新，当选择lb（弹性负载均衡健康检查）时，use_lb=1",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ScalingHealthMode...),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("use_lb"),
						types.Int32Value(1),
					),
				},
			},
			"subnet_id_list": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "子网ID列表，支持一主多辅。最多支持5个。支持更新，但是更新阶段status必须为2（停用）",
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 5), // 最大支持5个网卡信息
				},
			},
			"move_out_strategy": schema.StringAttribute{
				Required:    true,
				Description: "实例移出策略：earlier_config-较早创建的配置较早创建，later_config-较晚创建的配置较晚创建，earlier_vm-较早创建的云主机，later_vm-较晚创建的云主机。支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ScalingMoveOutStrategy...),
				},
			},
			// todo 做校验
			"use_lb": schema.Int32Attribute{
				Required:    true,
				Description: "是否使用负载均衡：1-是，2-否。status=disable时支持更新",
				Validators: []validator.Int32{
					int32validator.OneOf(1, 2),
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
			"min_count": schema.Int32Attribute{
				Required:    true,
				Description: "最小云主机数，取值范围：[0,50]。支持更新",
				Validators: []validator.Int32{
					int32validator.Between(0, 50),
					validator2.ScalingCountValidate(),
				},
			},
			"max_count": schema.Int32Attribute{
				Required:    true,
				Description: "最大云主机数，取值范围：[min_count,2147483647]。支持更新",
				Validators: []validator.Int32{
					validator2.ScalingCountValidate(),
				},
			},
			"expected_count": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "期望云主机数，非多可用区资源池不支持该参数。支持更新，若未填写。expected_count默认=min_count",
				Validators: []validator.Int32{
					validator2.ScalingCountValidate(),
				},
			},
			"health_period": schema.Int32Attribute{
				Required:    true,
				Description: "健康检查时间间隔（周期），单位：秒，取值范围：[300,10080]。支持更新",
				Validators: []validator.Int32{
					int32validator.Between(300, 10080),
				},
			},
			// todo status=disable时支持更新
			"lb_list": schema.ListNestedAttribute{
				Optional: true,
				Computed: false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"port": schema.Int32Attribute{
							Required:    true,
							Description: "端口号，当status=disable时支持更新",
						},
						"lb_id": schema.StringAttribute{
							Required:    true,
							Description: "负载均衡ID，当status=disable时支持更新",
						},
						"weight": schema.Int32Attribute{
							Required:    true,
							Description: "权重，当status=disable时支持更新",
						},
						"host_group_id": schema.StringAttribute{
							Required:    true,
							Description: "后端主机组ID，当status=disable时支持更新",
						},
					},
				},
				Description: "负载均衡列表，use_lb=1时必填。当status=disable时支持更新",
				Validators: []validator.List{
					validator2.AlsoRequiresEqualList(
						path.MatchRoot("use_lb"),
						types.Int32Value(1),
					),
				},
			},
			"config_list": schema.SetAttribute{
				ElementType: types.Int32Type,
				Optional:    true,
				Description: "伸缩配置ID列表，最大支持传入10个伸缩配置。支持更新",
				Validators: []validator.Set{
					setvalidator.SizeAtMost(10),
				},
			},
			"az_strategy": schema.StringAttribute{
				Required:    true,
				Description: "扩容策略类型：uniform_distribution-均衡分布，priority_distribution-优先级分布。支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ScalingAzStrategy...),
				},
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "伸缩组ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "伸缩组状态。取值范围：enable 或 disable，支持更新。可以用于控制伸缩组的状态更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ScalingControlStatus...),
				},
			},
			"delete_protection": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "控制伸缩组保护，开启伸缩组保护，不可删除该伸缩组。取值范围：enable 或 disable，支持更新。可以用于控制伸缩组保护的开启/关闭",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ScalingControlProtectionStatus...),
				},
			},
			"add_instance_uuid_list": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "需要手动添加的云主机uuid列表。伸缩组内云主机清单可以根据data.ctyun_scaling_ecs_list获取。支持更新。",
			},
			"remove_instance_uuid_list": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "需要删除手动/自动加入伸缩组的云主机uuid列表。伸缩组内云主机清单可以根据data.ctyun_scaling_ecs_list获取。支持更新。",
			},
			//"protect_status": schema.StringAttribute{
			//	Optional:    true,
			//	Computed:    true,
			//	Default:     stringdefault.StaticString(business.ProtectStatusUnprotectedStr),
			//	Description: "云主机保护状态（仅对手动添加的机器做处理），设置了保护状态的云主机实例，在伸缩组进行缩容活动时将不会被移出。disable-关闭云主机保护，enable-开启云主机保护。支持更新",
			//	Validators: []validator.String{
			//		stringvalidator.OneOf(business.ScalingPolicyStatuses...),
			//	},
			//},
			"is_destroy": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "移除时是否销毁，仅当移除云主机时生效（对手动添加的机器做处理），true-ecs移出伸缩组时销毁， false-ecs移出伸缩组时不销毁，支持更新",
			},
			"real_count": schema.Int32Attribute{
				Computed:    true,
				Description: "当前的云主机数量，直接通过接口获取，一般为expected_count + 手动添加（+）/移除（-）云主机数量。",
			},
		},
	}
}

func (c *ctyunScaling) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunScalingConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	//创建前检查,检查创建参数有效性
	isValid, err := c.checkBeforeScaling(ctx, plan)
	if !isValid || err != nil {
		return
	}
	err = c.createScaling(ctx, &plan)
	if err != nil {
		return
	}
	// 创建后，通过轮询，查询组内机器时候满足最小预期
	err = c.createLoop(ctx, &plan, 10)
	if err != nil {
		return
	}
	// 判断是否需要手动添加云主机
	err = c.manualAddInstance(ctx, &plan)
	if err != nil {
		return
	}
	// 添加完成后，确认云主机数量
	err = c.checkAfterAddEcs(ctx, &plan)
	if err != nil {
		return
	}
	// 开启云主机保护
	// 若创建时候，就需要开启/关闭云主机保护
	//err = c.updateProtectStatus(ctx, &plan, &plan)
	//if err != nil {
	//	return
	//}
	// 创建后反查创建后的证书信息
	err = c.getAndMergeScaling(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunScaling) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunScalingConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeScaling(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "NotExists") || strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunScaling) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取 plan -tf文件中配置
	var plan CtyunScalingConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunScalingConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 控制伸缩组开关
	err = c.controlScaling(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 控制伸缩组保护开关
	err = c.controlScalingProtection(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新基本信息
	err = c.updateScaling(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 修改弹性组机器列表
	err = c.updateInstanceByUUIDList(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 控制云主机保护开关
	//err = c.updateProtectStatus(ctx, &state, &plan)
	//if err != nil {
	//	return
	//}
	// 更新远端数据，并同步本地state
	err = c.getAndMergeScaling(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunScaling) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunScalingConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	force := true
	params := &scaling.ScalingGroupDeleteRequest{
		RegionID: state.RegionID.ValueString(),
		GroupID:  state.ID.ValueInt64(),
		Force:    &force,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupDeleteApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = errors.New("删除弹性伸缩实例失败，接口返回nil。请稍后再试！")
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	time.Sleep(30 * time.Second)
}

func (c *ctyunScaling) checkBeforeScaling(ctx context.Context, plan CtyunScalingConfig) (bool, error) {
	return true, nil
}

func (c *ctyunScaling) createScaling(ctx context.Context, config *CtyunScalingConfig) error {

	if config.ExpectedCount.IsNull() || config.ExpectedCount.IsUnknown() {
		config.ExpectedCount = config.MinCount
	}

	params := &scaling.ScalingGroupCreateRequest{
		RegionID:        config.RegionID.ValueString(),
		RecoveryMode:    1,
		Name:            config.Name.ValueString(),
		HealthMode:      business.ScalingHealthModeDict[config.HealthMode.ValueString()],
		MoveOutStrategy: business.ScalingMoveOutStrategyDict[config.MoveOutStrategy.ValueString()],
		UseLb:           config.UseLb.ValueInt32(),
		VpcID:           config.VpcID.ValueString(),
		MinCount:        config.MinCount.ValueInt32(),
		MaxCount:        config.MaxCount.ValueInt32(),
		ExpectedCount:   config.ExpectedCount.ValueInt32Pointer(),
		HealthPeriod:    config.HealthPeriod.ValueInt32(),
	}
	// 判断资源池是否为多AZ
	//zones, err2 := c.regionService.GetZonesByRegionID(ctx, config.RegionID.ValueString())
	//if err2 != nil {
	//	return err2
	//}
	//isNaz := false
	//if len(zones) > 1 {
	//	isNaz = true
	//}

	// securityGroupIDList， expectedCount 非多az不传
	if !config.SecurityGroupIDList.IsNull() && !config.SecurityGroupIDList.IsUnknown() {
		var securityGroupIDList []string
		diags := config.SecurityGroupIDList.ElementsAs(ctx, &securityGroupIDList, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		params.SecurityGroupIDList = securityGroupIDList
	}

	if !config.AzStrategy.IsNull() && !config.AzStrategy.IsUnknown() {
		params.AzStrategy = business.ScalingAzStrategyDict[config.AzStrategy.ValueString()]
	}

	if !config.SubnetIDList.IsNull() && !config.SubnetIDList.IsUnknown() {
		var subnetIDList []string
		diags := config.SubnetIDList.ElementsAs(ctx, &subnetIDList, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		params.SubnetIDList = subnetIDList
	}

	if !config.LbList.IsNull() && !config.LbList.IsUnknown() {
		var lbList []CtyunLbInfoModel
		var paramLbList []*scaling.ScalingGroupCreateLbListRequest
		diags := config.LbList.ElementsAs(ctx, &lbList, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		for _, lbItem := range lbList {
			var lbInfo scaling.ScalingGroupCreateLbListRequest
			lbInfo.Port = lbItem.Port.ValueInt32()
			lbInfo.HostGroupID = lbItem.HostGroupID.ValueString()
			lbInfo.LbID = lbItem.LbID.ValueString()
			lbInfo.Weight = lbItem.Weight.ValueInt32()
			paramLbList = append(paramLbList, &lbInfo)
		}
		params.LbList = paramLbList
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		params.ProjectID = config.ProjectID.ValueString()
	}
	if !config.ConfigList.IsNull() && !config.ConfigList.IsUnknown() {
		var configList []int32
		diags := config.ConfigList.ElementsAs(ctx, &configList, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		params.ConfigList = configList
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupCreateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		return errors.New("创建弹性伸缩服务组时，返回为nil。请稍微重试")
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}
	config.ID = types.Int64Value(resp.ReturnObj.GroupID)
	return nil
}

func (c *ctyunScaling) getAndMergeScaling(ctx context.Context, config *CtyunScalingConfig) error {
	respDetail, err := c.getScalingDetail(ctx, config)

	if err != nil {
		return err
	}
	scalingDetail := respDetail.ReturnObj.ScalingGroups[0]
	config.Name = types.StringValue(scalingDetail.Name)
	config.MinCount = types.Int32Value(scalingDetail.MinCount)
	config.MaxCount = types.Int32Value(scalingDetail.MaxCount)
	//  todo 考虑是否需要同步
	//config.ExpectedCount = types.Int32Value(scalingDetail.ExpectedCount)
	config.RealCount = types.Int32Value(scalingDetail.InstanceCount)
	config.UseLb = types.Int32Value(scalingDetail.UseLb)
	config.HealthPeriod = types.Int32Value(scalingDetail.HealthPeriod)
	var diags diag.Diagnostics
	config.SecurityGroupIDList, diags = types.SetValueFrom(ctx, types.StringType, scalingDetail.SecurityGroupIDList)
	if diags.HasError() {
		return errors.New(diags[0].Detail())
	}
	// 处理lb列表
	if scalingDetail.UseLb == 1 {
		var lbList []CtyunLbInfoModel
		lbListResp, err2 := c.getLbList(ctx, config)
		if err2 != nil {
			return err2
		}
		lbListReturnList := lbListResp.ReturnObj.LoadBalancers
		for _, lbItem := range lbListReturnList {
			var lbInfo CtyunLbInfoModel
			lbInfo.Port = types.Int32Value(lbItem.Port)
			lbInfo.LbID = types.StringValue(lbItem.LbID)
			lbInfo.Weight = types.Int32Value(lbItem.Weight)
			lbInfo.HostGroupID = types.StringValue(lbItem.HostGroupID)
			lbList = append(lbList, lbInfo)
		}

		lbObj := utils.StructToTFObjectTypes(CtyunLbInfoModel{})
		config.LbList, diags = types.ListValueFrom(ctx, lbObj, lbList)
		if diags.HasError() {
			return errors.New(diags[0].Detail())
		}
	}

	// 处理subnetIDList
	config.SubnetIDList, diags = types.SetValueFrom(ctx, types.StringType, scalingDetail.SubnetIDList)
	if diags.HasError() {
		return errors.New(diags[0].Detail())
	}
	config.MoveOutStrategy = types.StringValue(business.ScalingMoveOutStrategyDictRev[scalingDetail.MoveOutStrategy])
	config.HealthMode = types.StringValue(business.ScalingHealthModeDictRev[scalingDetail.HealthMode])
	config.ConfigList, diags = types.SetValueFrom(ctx, types.Int64Type, scalingDetail.ConfigList)
	config.AzStrategy = types.StringValue(business.ScalingAzStrategyDictRev[scalingDetail.AzStrategy])
	config.Status = types.StringValue(business.ScalingControlStatusDictRev[scalingDetail.Status])
	config.DeleteProtection = types.StringValue(business.ScalingControlProtectionDictRev[*scalingDetail.DeleteProtection])
	return nil
}

func (c *ctyunScaling) getScalingDetail(ctx context.Context, config *CtyunScalingConfig) (*scaling.ScalingGroupListResponse, error) {
	params := &scaling.ScalingGroupListRequest{
		RegionID: config.RegionID.ValueString(),
		GroupID:  config.ID.ValueInt64(),
		PageNo:   1,
		PageSize: 10,
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		params.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = errors.New("获取弹性伸缩列表信息返回nil，请稍后重试或联系研发人员！")
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	if len(resp.ReturnObj.ScalingGroups) > 1 {
		err = fmt.Errorf("根据groupid: %d 获取的弹性伸缩详情返回多个实例。具体如下:%#v\n", config.ID.ValueInt64(), resp.ReturnObj.ScalingGroups)
		return nil, err
	}
	return resp, nil
}

func (c *ctyunScaling) getLbList(ctx context.Context, config *CtyunScalingConfig) (*scaling.ScalingGroupQueryLoadBalancerListResponse, error) {
	params := &scaling.ScalingGroupQueryLoadBalancerListRequest{
		RegionID: config.RegionID.ValueString(),
		GroupID:  config.ID.ValueInt64(),
		PageNo:   1,
		PageSize: 100,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupQueryLoadBalancerListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = errors.New("获取伸缩组的负载均衡器列表失败，接口返回nil，请稍后再试！")
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp, nil
}

func (c *ctyunScaling) updateScaling(ctx context.Context, state *CtyunScalingConfig, plan *CtyunScalingConfig) error {
	//1. 伸缩组为非启用状态、实例数为0且没有正在进行的伸缩活动时，可以修改伸缩组的子网。
	//2. 伸缩组为非启用状态且没有正在进行的伸缩活动时，可以修改伸缩组关联的负载均衡器。
	// 判断是否需要更新ExpectedCount
	if plan.ExpectedCount.IsNull() || plan.ExpectedCount.IsUnknown() {
		plan.ExpectedCount = plan.MinCount
	} else {
		state.ExpectedCount = plan.ExpectedCount
	}

	params := &scaling.ScalingGroupUpdateRequest{
		RegionID:        state.RegionID.ValueString(),
		GroupID:         state.ID.ValueInt64(),
		Name:            plan.Name.ValueString(),
		MinCount:        plan.MinCount.ValueInt32(),
		MaxCount:        plan.MaxCount.ValueInt32(),
		ExpectedCount:   plan.ExpectedCount.ValueInt32Pointer(),
		HealthMode:      business.ScalingHealthModeDict[plan.HealthMode.ValueString()],
		HealthPeriod:    plan.HealthPeriod.ValueInt32(),
		MoveOutStrategy: business.ScalingMoveOutStrategyDict[plan.MoveOutStrategy.ValueString()],
	}

	if !plan.LbList.IsNull() && !plan.LbList.IsUnknown() && !plan.LbList.Equal(state.LbList) {
		// lb的更新需要status = disable
		detail, err := c.getScalingDetail(ctx, state)
		if err != nil {
			return err
		}
		status := detail.ReturnObj.ScalingGroups[0].Status
		if status == business.ScalingControlStatusStart {
			return errors.New("伸缩组状态为启用状态，暂不支持修改lb列表")
		}
		params.UseLb = plan.UseLb.ValueInt32Pointer()
		// 必须启用lb
		if plan.UseLb.ValueInt32() == 1 {
			var lbList []CtyunLbInfoModel
			var paramLbList []*scaling.ScalingGroupUpdateLbListRequest
			diags := plan.LbList.ElementsAs(ctx, &lbList, true)
			if diags.HasError() {
				return errors.New(diags[0].Detail())
			}
			for _, lbItem := range lbList {
				var lbInfo scaling.ScalingGroupUpdateLbListRequest
				lbInfo.Port = lbItem.Port.ValueInt32()
				lbInfo.HostGroupID = lbItem.HostGroupID.ValueString()
				lbInfo.Id = lbItem.LbID.ValueString()
				lbInfo.Weight = lbItem.Weight.ValueInt32()
				paramLbList = append(paramLbList, &lbInfo)
			}
			params.LbList = paramLbList
		}
	}

	// 判断SecurityGroupIDList 是否需要更新
	if !plan.SecurityGroupIDList.IsNull() && !plan.SecurityGroupIDList.Equal(state.SecurityGroupIDList) {
		var securityGroupIDList []string
		diags := plan.SecurityGroupIDList.ElementsAs(ctx, &securityGroupIDList, true)
		if diags.HasError() {
			return errors.New(diags[0].Detail())
		}
		params.SecurityGroupIDList = securityGroupIDList
	}
	//  判断SubnetIDList是否需要更新
	if !plan.SubnetIDList.IsNull() && !plan.SubnetIDList.Equal(state.SubnetIDList) {
		// 先判断伸缩组状态，如果状态为停用可以更新。 status = 2
		detail, err := c.getScalingDetail(ctx, state)
		if err != nil {
			return err
		}
		status := detail.ReturnObj.ScalingGroups[0].Status
		if status == business.ScalingControlStatusStart {
			return errors.New("伸缩组状态为启用状态，暂不支持修改子网列表")
		}
		var subnetIDList []string
		diags := plan.SubnetIDList.ElementsAs(ctx, &subnetIDList, true)
		if diags.HasError() {
			return errors.New(diags[0].Detail())
		}
		params.SubnetIDList = subnetIDList
	}
	// 判断configList是否需要更新
	if !plan.ConfigList.IsNull() && !plan.ConfigList.Equal(state.ConfigList) {
		var configList []int32
		diags := plan.ConfigList.ElementsAs(ctx, &configList, true)
		if diags.HasError() {
			return errors.New(diags[0].Detail())
		}
		params.ConfigList = configList
	}

	// 判断AzStrategy是否需要更新
	if !plan.AzStrategy.IsNull() && !plan.AzStrategy.Equal(state.AzStrategy) {
		params.AzStrategy = business.ScalingAzStrategyDict[plan.AzStrategy.ValueString()]
	}

	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupUpdateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		return errors.New("更新弹性伸缩配置失败，接口返回nil。请稍后重试！")
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}
	return nil
}

func (c *ctyunScaling) controlScaling(ctx context.Context, state *CtyunScalingConfig, plan *CtyunScalingConfig) error {
	// 判断开关，或ControlStatus字段为空
	if plan.Status.IsNull() || plan.Status.IsUnknown() {
		return nil
	}
	// 启用伸缩组
	if plan.Status.ValueString() == business.ScalingControlStatusStartStr {
		startParams := &scaling.ScalingGroupEnableRequest{
			RegionID: state.RegionID.ValueString(),
			GroupID:  state.ID.ValueInt64(),
		}
		resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupEnableApi.Do(ctx, c.meta.SdkCredential, startParams)
		if err != nil {
			return err
		} else if resp == nil {
			return errors.New("启用伸缩组失败，接口返回nil。请稍后重试！")
		} else if resp.StatusCode != common.NormalStatusCode {
			return fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		}
	} else if plan.Status.ValueString() == business.ScalingControlStatusStopStr {
		stopParams := &scaling.ScalingGroupDisableRequest{
			RegionID: state.RegionID.ValueString(),
			GroupID:  state.ID.ValueInt64(),
		}
		resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupDisableApi.Do(ctx, c.meta.SdkCredential, stopParams)
		if err != nil {
			return err
		} else if resp == nil {
			return errors.New("停用伸缩组失败，接口返回nil。请稍后重试！")
		} else if resp.StatusCode != common.NormalStatusCode {
			return fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		}
	}
	// 确认伸缩组状态是否变更成功
	flag, err := c.checkStatusAfterUpdate(ctx, state, plan)
	if err != nil {
		return err
	}
	if !flag {
		return errors.New("更新伸缩组状态失败！")
	}
	return nil
}

func (c *ctyunScaling) checkStatusAfterUpdate(ctx context.Context, state *CtyunScalingConfig, plan *CtyunScalingConfig) (bool, error) {
	// 获取伸缩组详情
	respDetail, err := c.getScalingDetail(ctx, state)
	if err != nil {
		return false, err
	}
	status := respDetail.ReturnObj.ScalingGroups[0].Status
	if status == business.ScalingControlStatusDict[plan.Status.ValueString()] {
		return true, nil
	}
	return false, nil
}

func (c *ctyunScaling) controlScalingProtection(ctx context.Context, state *CtyunScalingConfig, plan *CtyunScalingConfig) error {
	if plan.DeleteProtection.IsNull() || plan.DeleteProtection.IsUnknown() {
		return nil
	}
	if state.DeleteProtection.ValueString() == business.ScalingControlProtectionEnableStr {
		enableParams := &scaling.ScalingGroupEnableProtectionRequest{
			RegionID: state.RegionID.ValueString(),
			GroupID:  state.ID.ValueInt64(),
		}
		resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupEnableProtectionApi.Do(ctx, c.meta.SdkCredential, enableParams)
		if err != nil {
			return err
		} else if resp == nil {
			return errors.New("开启伸缩组保护失败，接口返回nil。请稍后再试！")
		} else if resp.StatusCode != common.NormalStatusCode {
			return fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		}
	} else if state.DeleteProtection.ValueString() == business.ScalingControlProtectionDisableStr {
		disableParams := &scaling.ScalingGroupDisableProtectionRequest{
			RegionID: state.RegionID.ValueString(),
			GroupID:  state.ID.ValueInt64(),
		}
		resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupDisableProtectionApi.Do(ctx, c.meta.SdkCredential, disableParams)
		if err != nil {
			return err
		} else if resp == nil {
			return errors.New("关闭伸缩组保护失败，接口返回nil。请稍后再试！")
		} else if resp.StatusCode != common.NormalStatusCode {
			return fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		}
	}

	// 查询判断伸缩组保护是否开启/关闭
	flag, err := c.checkProtectionAfterUpdate(ctx, state, plan)
	if err != nil {
		return err
	}
	if !flag {
		return errors.New("伸缩组保护状态更新失败！")
	}
	return nil
}

func (c *ctyunScaling) checkProtectionAfterUpdate(ctx context.Context, state *CtyunScalingConfig, plan *CtyunScalingConfig) (bool, error) {
	// 获取弹性伸缩详情
	respDetail, err := c.getScalingDetail(ctx, state)
	if err != nil {
		return false, err
	}
	protectionStatus := respDetail.ReturnObj.ScalingGroups[0].DeleteProtection
	if *protectionStatus == business.ScalingControlProtectionDict[plan.DeleteProtection.ValueString()] {
		return true, nil
	}
	return false, nil
}

func (c *ctyunScaling) createLoop(ctx context.Context, config *CtyunScalingConfig, loopCount ...int) error {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	var err error
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}

	result := retryer.Start(
		func(currentTime int) bool {
			// 获取当前机器列表，判断当前伸缩组内机器是否符合预期
			// 判断跳出轮询的条件：
			// 1. 当伸缩组未填写expected count， 机器组数量 = min count即可
			// 2. 当伸缩组填写expected count， 机器组数量 = expected count即可
			instanceList, err2 := c.getInstanceListByGroupID(ctx, config)
			if err2 != nil {
				err = err2
				return false
			}
			if int32(len(instanceList)) == config.ExpectedCount.ValueInt32() {
				return false
			}
			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return fmt.Errorf("轮询已达最大次数，弹性组（id:%d）伸缩实例数量仍未达到预期%d台！", config.ID.ValueInt64(), config.ExpectedCount.ValueInt32())
	}
	return err
}

// 根据弹性伸缩组id获取云主机全量列表
func (c *ctyunScaling) getInstanceListByGroupID(ctx context.Context, config *CtyunScalingConfig) ([]*scaling.ScalingGroupQueryInstanceListReturnObjInstanceListResponse, error) {
	var pageSize, pageNo int32
	pageSize = 100
	pageNo = 1
	resp, err := c.requestInstanceListByGroup(ctx, config, pageNo, pageSize)
	if err != nil {
		return nil, err
	}

	totalCount := resp.ReturnObj.TotalCount
	totalPageNo := pageNo
	if totalCount > pageSize {
		totalPageNo = totalCount/pageSize + 1
	}
	var instances []*scaling.ScalingGroupQueryInstanceListReturnObjInstanceListResponse
	for pageNo <= totalPageNo {
		instanceList := resp.ReturnObj.InstanceList
		for _, instance := range instanceList {
			instances = append(instances, instance)
		}
		pageNo++
		if pageNo > totalPageNo {
			break
		}
		resp, err = c.requestInstanceListByGroup(ctx, config, pageNo, pageSize)
		if err != nil {
			return nil, err
		}
	}
	return instances, nil
}

func (c *ctyunScaling) requestInstanceListByGroup(ctx context.Context, config *CtyunScalingConfig, pageNo, pageSize int32) (*scaling.ScalingGroupQueryInstanceListResponse, error) {
	params := &scaling.ScalingGroupQueryInstanceListRequest{
		RegionID: config.RegionID.ValueString(),
		GroupID:  config.ID.ValueInt64(),
		PageNo:   pageNo,
		PageSize: pageSize,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupQueryInstanceListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询group id为%d下的云主机列表失败，接口范围nil。请联系研发，或稍后重试！", config.ID.ValueInt64())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp, nil
}

func (c *ctyunScaling) manualAddInstance(ctx context.Context, config *CtyunScalingConfig) error {
	if config.AddInstanceUUIDList.IsNull() || config.AddInstanceUUIDList.IsUnknown() {
		return nil
	}
	var instanceUUIDList []string
	diags := config.AddInstanceUUIDList.ElementsAs(ctx, &instanceUUIDList, true)
	if diags.HasError() {
		err := errors.New(diags[0].Detail())
		return err
	}
	// 添加云主机前进行合法校验，校验两个方面：
	// 1）expected_count + 手动添加云主机个数 <= max_count
	// 2）添加的机器是否已经存在于弹性伸缩机器组内
	isValid, err := c.checkBeforeManualAddEcs(ctx, config)
	if err != nil {
		return err
	}
	// 如果合法再进行添加
	if isValid {
		err = c.addScalingEcsList(ctx, instanceUUIDList, config)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
checkBeforeManualAddEcs
// 添加云主机前进行合法校验，校验两个方面：
// 1）expected_count + 手动添加云主机个数 <= max_count
// 2）添加的机器是否已经存在于弹性伸缩机器组内
*/
func (c *ctyunScaling) checkBeforeManualAddEcs(ctx context.Context, config *CtyunScalingConfig) (bool, error) {
	// check 1: expected_count + 手动添加云主机个数 <= max_count
	// 校验最大实例数量
	// 获取scaling group 详情
	params := &scaling.ScalingGroupListRequest{
		RegionID: config.RegionID.ValueString(),
		GroupID:  config.ID.ValueInt64(),
		PageNo:   1,
		PageSize: 10,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return false, err
	} else if resp == nil {
		err = errors.New("获取弹性伸缩列表信息返回nil，请稍后重试或联系研发人员！")
		return false, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return false, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return false, err
	}
	if len(resp.ReturnObj.ScalingGroups) > 1 {
		err = fmt.Errorf("根据groupid: %d 获取的弹性伸缩详情返回多个实例。具体如下:%#v\n", config.ID.ValueInt64(), resp.ReturnObj.ScalingGroups)
		return false, err
	}
	// 确认maxCount
	maxCount := resp.ReturnObj.ScalingGroups[0].MaxCount
	// 获取原本弹性组下有多少台机器
	instanceList, err := c.getInstanceListByGroupID(ctx, config)
	if err != nil {
		return false, err
	}
	var instanceUUIDList []string
	diags := config.AddInstanceUUIDList.ElementsAs(ctx, &instanceUUIDList, true)
	if diags.HasError() {
		err = errors.New(diags[0].Detail())
		return false, err
	}
	if len(instanceList)+len(instanceUUIDList) > int(maxCount) {
		err = fmt.Errorf("弹性伸缩组（id：%d）的max_count=%d。当前伸缩组下有云主机%d台，若手动移入%d台，将移入失败！", config.ID.ValueInt64(), maxCount, len(instanceList), len(instanceUUIDList))
		return false, err
	}

	// check 2:
	// 确认待添加云主机列表不在弹性组云主机清单内
	// params : 待添加的机器列表， 伸缩组机器列表
	intersection := c.FindIntersection(instanceUUIDList, instanceList)
	// 不能有并集
	if len(intersection) > 0 {
		err = fmt.Errorf("待添加机器列表与弹性组内云主机重叠，待添加的列表中已被添加至弹性组的云主机：%s", strings.Join(intersection, ","))
		return false, err
	}
	return true, nil
}

func (c *ctyunScaling) addScalingEcsList(ctx context.Context, instanceUUIDList []string, config *CtyunScalingConfig) error {
	if instanceUUIDList == nil || len(instanceUUIDList) == 0 {
		return nil
	}
	params := &scaling.ScalingGroupInstanceMoveInRequest{
		RegionID: config.RegionID.ValueString(),
		GroupID:  config.ID.ValueInt64(),
	}

	params.InstanceUUIDList = instanceUUIDList
	// 轮询确认伸缩组无伸缩动作后再进行操作
	err := c.updateInstanceBeforeLoop(ctx, config, 60)
	if err != nil {
		return err
	}

	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupInstanceMoveInApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		return fmt.Errorf("手动添加云主机失败，接口返回nil，伸缩组id：%d。需添加的云主机id列表为：%s", config.ID.ValueInt64(), strings.Join(instanceUUIDList, ", "))
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
	}
	return nil
}

func (c *ctyunScaling) checkAfterAddEcs(ctx context.Context, config *CtyunScalingConfig) error {
	instanceList, err := c.getInstanceListByGroupID(ctx, config)
	if err != nil {
		return err
	}
	var instanceUUIDList []string
	var planInstanceUUIDList []string
	for _, instance := range instanceList {
		if instance.ExecutionMode == business.ExecutionModeManualAddInstances {
			instanceUUIDList = append(instanceUUIDList, instance.InstanceID)
		}
	}
	diags := config.AddInstanceUUIDList.ElementsAs(ctx, &planInstanceUUIDList, true)
	if diags.HasError() {
		err = errors.New(diags[0].Detail())
		return err
	}
	if len(instanceUUIDList) == len(planInstanceUUIDList) {
		return nil
	} else {
		return fmt.Errorf("弹性组（%d）手动插入云主机失败，当前手动已成功添加云主机列表为：%s，预期需要添加的云主机列表为：%s", config.ID.ValueInt64(),
			strings.Join(planInstanceUUIDList, ", "), strings.Join(instanceUUIDList, ", "))
	}
}

//func (c *ctyunScaling) updateProtectStatus(ctx context.Context, state *CtyunScalingConfig, plan *CtyunScalingConfig) error {
//	if !plan.ProtectStatus.IsNull() && !plan.InstanceUUIDList.IsNull() && !plan.InstanceUUIDList.IsUnknown() {
//		var instanceUUIDs []string
//		diags := plan.InstanceUUIDList.ElementsAs(ctx, &instanceUUIDs, true)
//		if diags.HasError() {
//			err := errors.New(diags[0].Detail())
//			return err
//		}
//
//		instanceIds, err := c.getInstanceAssocIdByUUID(ctx, state, instanceUUIDs)
//		if err != nil {
//			return err
//		}
//		// 关闭云主机保护
//		if plan.ProtectStatus.ValueString() == business.StatusDisabledStr {
//			err = c.disableProtectEcs(ctx, state, instanceIds, instanceUUIDs)
//			if err != nil {
//				return err
//			}
//		} else if plan.ProtectStatus.ValueString() == business.StatusEnabledStr {
//			// 开启云主机保护
//			err = c.enableProtectEcs(ctx, state, instanceIds, instanceUUIDs)
//			if err != nil {
//				return err
//			}
//		}
//	}
//	state.ProtectStatus = plan.ProtectStatus
//	return nil
//}

func (c *ctyunScaling) getInstanceAssocIdByUUID(ctx context.Context, state *CtyunScalingConfig, UUIDs []string) ([]int32, error) {
	instanceMap, err := c.getInstanceMap(ctx, state)
	if err != nil {
		return nil, err
	}
	var instanceIdList []int32
	for _, uuid := range UUIDs {
		ecsInfo := instanceMap[uuid]
		if ecsInfo == nil {
			continue
		}
		instanceIdList = append(instanceIdList, ecsInfo.Id)
	}
	return instanceIdList, nil
}

func (c *ctyunScaling) getInstanceMap(ctx context.Context, config *CtyunScalingConfig) (map[string]*scaling.ScalingGroupQueryInstanceListReturnObjInstanceListResponse, error) {
	var pageNo, pageSize int32
	pageNo = 1
	pageSize = 100
	pageEndNo := pageNo
	ecsListResp, err := c.requestInstanceListByGroup(ctx, config, pageNo, pageSize)
	if err != nil {
		return nil, err
	}

	totalCount := ecsListResp.ReturnObj.TotalCount

	if totalCount > pageSize {
		pageEndNo = totalCount / pageSize
	}
	// 先获取所有ecs列表，并设置成map
	instanceMap := make(map[string]*scaling.ScalingGroupQueryInstanceListReturnObjInstanceListResponse)
	for pageNo <= pageEndNo {
		ecsList := ecsListResp.ReturnObj.InstanceList
		for _, ecs := range ecsList {
			instanceId := ecs.InstanceID
			instanceMap[instanceId] = ecs
		}

		pageNo++
		if pageNo > pageEndNo {
			break
		}
		ecsListResp, err = c.requestInstanceListByGroup(ctx, config, pageNo, pageSize)
		if err != nil {
			return nil, err
		}
	}
	return instanceMap, nil
}

func (c *ctyunScaling) updateInstanceByUUIDList(ctx context.Context, state *CtyunScalingConfig, plan *CtyunScalingConfig) error {

	// 比对list， 挑出需要删除/新增的机器
	//var stateInstanceList []string
	instanceList, err := c.getInstanceListByGroupID(ctx, state)
	if err != nil {
		return err
	}
	if !plan.AddInstanceUUIDList.IsNull() && !state.AddInstanceUUIDList.Equal(plan.AddInstanceUUIDList) {
		var addInstanceUUIDList []string
		diags := plan.AddInstanceUUIDList.ElementsAs(ctx, &addInstanceUUIDList, true)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return err
		}
		// 确认待添加云主机未在伸缩组内
		addIntersection := c.FindIntersection(addInstanceUUIDList, instanceList)
		if len(addIntersection) > 0 {
			err = fmt.Errorf("待添加的云主机中有部分已加入伸缩组，列表为：%s", strings.Join(addInstanceUUIDList, ", "))
			return err
		}

		// 轮询确认是否可以更新
		err = c.updateInstanceBeforeLoop(ctx, state, 60)
		if err != nil {
			return err
		}

		// 处理新增
		err = c.addScalingEcsList(ctx, addInstanceUUIDList, state)
		if err != nil {
			return err
		}
	}
	// 处理删除列表
	// 1. 确认待移除的云主机都包含在此伸缩组内
	// 2. 区分是否存在自动伸缩的机器？是否需要销毁
	if !plan.RemoveInstanceUUIDList.IsNull() && !state.RemoveInstanceUUIDList.Equal(plan.RemoveInstanceUUIDList) {
		var removeInstanceUUIDList []string
		diags := plan.RemoveInstanceUUIDList.ElementsAs(ctx, &removeInstanceUUIDList, true)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return err
		}
		// 确认待添加云主机未在伸缩组内
		removeIntersection := c.FindIntersection(removeInstanceUUIDList, instanceList)
		if len(removeIntersection) != len(removeInstanceUUIDList) {
			err = fmt.Errorf("待删除的云主机中有部分未加入伸缩组，符合删除条件列表为：%s", strings.Join(removeIntersection, ", "))
			return err
		}
		// 轮询确认是否可以更新
		err = c.updateInstanceBeforeLoop(ctx, state, 60)
		if err != nil {
			return err
		}
		// 处理删除
		// 查看is_destroy是否为true。若为true挑选出自动伸缩的列表
		if plan.IsDestroy.ValueBool() {
			// 挑出自动伸缩云主机列表
			autoAddEcs, manualAddEcs := c.getAutoAndManualEcs(removeInstanceUUIDList, instanceList)
			err = c.destroyAutoEcsRequest(ctx, autoAddEcs, state)
			if err != nil {
				return err
			}
			err = c.removeEcsRequest(ctx, manualAddEcs, state, plan)
			if err != nil {
				return err
			}
		} else {
			err = c.removeEcsRequest(ctx, removeInstanceUUIDList, state, plan)
			if err != nil {
				return err
			}
		}
	}
	state.AddInstanceUUIDList = plan.AddInstanceUUIDList
	state.RemoveInstanceUUIDList = plan.RemoveInstanceUUIDList
	state.IsDestroy = plan.IsDestroy
	return nil
}

//func (c *ctyunScaling) getDiffInstanceList(state, plan []string) (toAdd, toRemove []string) {
//	// 使用 map 快速查找差异
//	planSet := make(map[string]bool)
//	stateSet := make(map[string]bool)
//
//	// 填充计划集
//	for _, item := range plan {
//		planSet[item] = true
//	}
//
//	// 填充状态集
//	for _, item := range state {
//		stateSet[item] = true
//	}
//
//	// 找出需要新增的项目（在 plan 中但不在 state 中）
//	for _, item := range plan {
//		if !stateSet[item] {
//			toAdd = append(toAdd, item)
//		}
//	}
//
//	// 找出需要删除的项目（在 state 中但不在 plan 中）
//	for _, item := range state {
//		if !planSet[item] {
//			toRemove = append(toRemove, item)
//		}
//	}
//	return toAdd, toRemove
//}

func (c *ctyunScaling) removeEcsRequest(ctx context.Context, instanceUUIDs []string, config *CtyunScalingConfig, plan *CtyunScalingConfig) error {
	if instanceUUIDs == nil || len(instanceUUIDs) <= 0 {
		return nil
	}
	// 移除不释放
	params := &scaling.ScalingGroupInstanceMoveOutRequest{
		RegionID:       config.RegionID.ValueString(),
		GroupID:        config.ID.ValueInt64(),
		InstanceIDList: instanceUUIDs,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupInstanceMoveOutApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("ecs移除失败，接口返回nil。ecs uuid 列表：%s", strings.Join(instanceUUIDs, ", "))
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	//}
	return nil
}

func (c *ctyunScaling) destroyAutoEcsRequest(ctx context.Context, instanceUUIDs []string, config *CtyunScalingConfig) error {
	// 移除并释放
	instanceIds, err := c.getInstanceAssocIdByUUID(ctx, config, instanceUUIDs)
	if err != nil {
		return err
	}
	params := &scaling.ScalingGroupInstanceMoveOutReleaseRequest{
		RegionID:       config.RegionID.ValueString(),
		GroupID:        config.ID.ValueInt64(),
		InstanceIDList: instanceIds,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupInstanceMoveOutReleaseApi.Do(ctx, c.meta.SdkCredential, params)

	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("ecs移除并释放失败，接口返回nil。ecs uuid 列表：%s", strings.Join(instanceUUIDs, ", "))
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	return nil
}

//func (c *ctyunScaling) checkUpdate(ctx context.Context, state *CtyunScalingConfig) (bool, error) {
//	params := &scaling.ScalingGroupCheckRequest{
//		RegionID: state.RegionID.ValueString(),
//		GroupID:  state.ID.ValueInt64(),
//	}
//	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupCheckApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return false, err
//	} else if resp.StatusCode != common.NormalStatusCode {
//		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
//		return false, err
//	}
//	flag := resp.ReturnObj.Result
//	if flag == 2 {
//		return false, nil
//	} else if flag == 1 {
//		return true, nil
//	}
//	return true, nil
//}

//func (c *ctyunScaling) updateBeforeLoop(ctx context.Context, state *CtyunScalingConfig, loopCount ...int) error {
//	var err error
//	count := 60
//	if len(loopCount) > 0 {
//		count = loopCount[0]
//	}
//	retryer, err := business.NewRetryer(time.Second*30, count)
//	if err != nil {
//		return err
//	}
//	result := retryer.Start(
//		func(currentTime int) bool {
//			isUpdate, err2 := c.checkUpdate(ctx, state)
//			if err2 != nil {
//				err = err2
//				return false
//			}
//			if isUpdate {
//				return false
//			}
//			return true
//		})
//	if result.ReturnReason == business.ReachMaxLoopTime {
//		return errors.New("轮询已达最大次数，弹性文件系统仍未扩容成功！")
//	}
//	return err
//}

func (c *ctyunScaling) updateInstanceBeforeLoop(ctx context.Context, state *CtyunScalingConfig, loopCount ...int) error {
	var err error
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			activitiesList, err2 := c.getScalingActivitiesList(ctx, state)
			if err2 != nil {
				err = err2
				return false
			}
			for _, activity := range activitiesList {
				executionResult := activity.ExecutionResult
				if executionResult == 0 {
					return true
				}
			}
			return false
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，弹性文件系统仍未扩容成功！")
	}
	return err
}

func (c *ctyunScaling) getScalingActivitiesList(ctx context.Context, state *CtyunScalingConfig) ([]*scaling.ScalingQueryActivitiesListReturnObjActiveListResponse, error) {
	params := &scaling.ScalingQueryActivitiesListRequest{
		RegionID: state.RegionID.ValueString(),
		GroupID:  state.ID.ValueInt64(),
		PageNo:   1,
		PageSize: 100,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingQueryActivitiesListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp.ReturnObj.ActiveList, nil
}

func (c *ctyunScaling) FindIntersection(instanceUUIDList []string, scalingInstanceList []*scaling.ScalingGroupQueryInstanceListReturnObjInstanceListResponse) []string {
	// 使用map记录第一个数组的元素
	set := make(map[string]bool)
	for _, item := range instanceUUIDList {
		set[item] = true
	}

	var intersection []string
	// 遍历第二个数组，检查元素是否在第一个数组中存在
	for _, item := range scalingInstanceList {
		instanceID := item.InstanceID
		if set[instanceID] {
			intersection = append(intersection, instanceID)
			set[instanceID] = false // 标记已添加，避免重复
		}
	}
	return intersection
}

func (c *ctyunScaling) getAutoAndManualEcs(removeEcsList []string, scalingEcsList []*scaling.ScalingGroupQueryInstanceListReturnObjInstanceListResponse) ([]string, []string) {
	ecsMap := make(map[string]int32)
	for _, ecs := range scalingEcsList {
		ecsMap[ecs.InstanceID] = ecs.ExecutionMode
	}

	var autoInstanceID, manualInstanceID []string
	for _, ecs := range removeEcsList {
		executionMode := ecsMap[ecs]
		if executionMode == business.ExecutionModeAutoStrategy || executionMode == business.ExecutionModeManualStrategy {
			autoInstanceID = append(autoInstanceID, ecs)
		} else if executionMode == business.ExecutionModeManualAddInstances {
			manualInstanceID = append(manualInstanceID, ecs)
		}
	}
	return autoInstanceID, manualInstanceID
}

type CtyunScalingConfig struct {
	RegionID            types.String `tfsdk:"region_id"`              // 资源池ID
	SecurityGroupIDList types.Set    `tfsdk:"security_group_id_list"` // 安全组ID列表
	Name                types.String `tfsdk:"name"`                   // 伸缩组名称
	HealthMode          types.String `tfsdk:"health_mode"`            // 健康检查方式
	SubnetIDList        types.Set    `tfsdk:"subnet_id_list"`         // 子网ID列表
	MoveOutStrategy     types.String `tfsdk:"move_out_strategy"`      // 实例移出策略
	UseLb               types.Int32  `tfsdk:"use_lb"`                 // 是否使用负载均衡
	VpcID               types.String `tfsdk:"vpc_id"`                 // 虚拟私有云ID
	MinCount            types.Int32  `tfsdk:"min_count"`              // 最小云主机数
	MaxCount            types.Int32  `tfsdk:"max_count"`              // 最大云主机数
	ExpectedCount       types.Int32  `tfsdk:"expected_count"`         // 期望云主机数
	RealCount           types.Int32  `tfsdk:"real_count"`             // 当前云主机数
	HealthPeriod        types.Int32  `tfsdk:"health_period"`          // 健康检查时间间隔
	LbList              types.List   `tfsdk:"lb_list"`                // 负载均衡列表
	ProjectID           types.String `tfsdk:"project_id"`             // 企业项目ID
	ConfigList          types.Set    `tfsdk:"config_list"`            // 伸缩配置ID列表
	AzStrategy          types.String `tfsdk:"az_strategy"`            // 扩容策略类型
	ID                  types.Int64  `tfsdk:"id"`                     // 伸缩组ID
	Status              types.String `tfsdk:"status"`                 // 伸缩组状态
	DeleteProtection    types.String `tfsdk:"delete_protection"`      // 控制伸缩组保护开关
	//InstanceUUIDList       types.Set    `tfsdk:"instance_uuid_list"`        // 云主机ID列表
	AddInstanceUUIDList    types.Set `tfsdk:"add_instance_uuid_list"`    // 需要手动添加的云主机列表
	RemoveInstanceUUIDList types.Set `tfsdk:"remove_instance_uuid_list"` // 需要手动移除的云主机列表
	//ProtectStatus          types.String `tfsdk:"protect_status"`            // 保护状态。1：已保护。2：未保护。
	IsDestroy types.Bool `tfsdk:"is_destroy"` // 移除时是否销毁
}

// LbInfo 负载均衡信息
type CtyunLbInfoModel struct {
	Port        types.Int32  `tfsdk:"port"`          // 端口号
	LbID        types.String `tfsdk:"lb_id"`         // 负载均衡ID
	Weight      types.Int32  `tfsdk:"weight"`        // 权重
	HostGroupID types.String `tfsdk:"host_group_id"` // 后端主机组ID
}
