package ebs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctebs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctebs"
	sdkPlanmodifier "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/sdk/planmodifier"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	explanmodifier "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/planmodifier"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
)

var (
	_ resource.Resource                = &ctyunEbs{}
	_ resource.ResourceWithConfigure   = &ctyunEbs{}
	_ resource.ResourceWithImportState = &ctyunEbs{}
)

type ctyunEbs struct {
	meta       *common.CtyunMetadata
	name       string
	ebsService *business.EbsService
}

func NewCtyunEbs() resource.Resource {
	return &ctyunEbs{}
}

func (c *ctyunEbs) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs"
	c.name = response.TypeName
}

func (c *ctyunEbs) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理云硬盘", "云硬盘（CT-EVS，Elastic Volume Service）", "https://www.ctyun.cn/document/10027696"),
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "磁盘命名，单账户单资源池下，命名需唯一，长度为2-64个字符，仅允许英文字母、数字及特殊字符._-，不能以特殊字符开头，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 64),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z][0-9a-zA-Z_-]+$"), "磁盘名称不符合规则"),
				},
			},
			"mode": schema.StringAttribute{
				Required:    true,
				Description: "磁盘模式，vbd，iscsi，fcsan",
				Validators: []validator.String{
					stringvalidator.OneOf(business.EbsDiskModes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "磁盘类型，SATA：普通IO，SAS：高IO，SSD：超高IO，SSD-genric：通用型SSD，FAST-SSD：极速型SSD，不支持ISCSI模式；XSSD-0、XSSD-1、XSSD-2：X系列云硬盘，不支持加密，不支持ISCSI模式或FCSAN模式",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					sdkPlanmodifier.EbsDiskTypeNormalize(business.EbsDiskTypeMap.FromOriginalScene, business.EbsDiskTypeMapScene1),
				},
				Validators: []validator.String{
					stringvalidator.Any(
						stringvalidator.OneOf(business.EbsDiskTypes...),
						stringvalidator.OneOf(business.EbsDiskTypesUpper...),
					),
				},
			},
			"size": schema.Int64Attribute{
				Required:    true,
				Description: "磁盘大小（单位GB），超高IO/高IO/极速型SSD/普通IO：取值范围[10, 32768]；XSSD-0：10GB-65536GB；XSSD-1：20GB-65536GB；XSSD-2：512GB-65536GB 支持更新（不支持缩容）",
				Validators: []validator.Int64{
					int64validator.Between(10, 65536),
				},
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围：month：按月，year：按年、on_demand：按需。当此值为month或者year时，cycle_count为必填",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypes...),
				},
			},
			"cycle_count": schema.Int64Attribute{
				Optional:    true,
				Description: "订购时长，该参数在cycle_type为month或year时才生效，当cycle_type=month，支持订购1-11个月；当cycle_type=year，支持订购1-5年",
				PlanModifiers: []planmodifier.Int64{
					explanmodifier.RequiresReplaceUnlessDependencyEqualsInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
				},
				Validators: []validator.Int64{
					validator2.AlsoRequiresEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeMonth),
						types.StringValue(business.OrderCycleTypeYear),
					),
					validator2.CycleCount(1, 11, 1, 5),
				},
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订购的受理单id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "磁盘id",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "云硬盘使用状态，deleting：删除中，creating：资源创建中，detaching：解绑中，detached：未绑定云主机，attaching：绑定中，attached：已绑定，extending：扩容中，error：错误状态，backup：备份中，backupRestoring：从备份恢复中，expired：包周期已结束，freezing：按需计费，处于冻结状态，可能账户受限或余额不足，available：可用，in-use：已挂载云主机，resizing：扩容中",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expire_time": schema.StringAttribute{
				Computed:    true,
				Description: "到期时间，为UTC格式，按需时为空",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"multi_attach": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否共享云硬盘",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"encrypted": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否加密盘； 共享盘、ISCSI模式磁盘、极速型SSD类型盘、XSSD系列盘不支持加密",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"kms_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "加密盘密钥UUID，是加密盘时才返回",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					explanmodifier.Project(),
				},
				Validators: []validator.String{
					validator2.Project(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
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
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区id，如果不填则默认使用provider ctyun中的az_name或环境变量中的CTYUN_AZ_NAME",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraAzName, true),
			},

			"provisioned_iops": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "XSSD类型云硬盘的预配置IOPS值，最小值为1，最大值计算公式为“min(单盘最大IOPS，500*容量) - 基础性能IOPS”。 其他类型磁盘不支持此参数 具体取值范围如下：\n\t●XSSD-0：（基础IOPS（min{1800+12×容量， 10000}） + 预配置IOPS） ≤ min{500×容量，100000}\n\t●XSSD-1：（基础IOPS（min{1800+50×容量， 50000}） + 预配置IOPS） ≤ min{500×容量，100000}\n\t●XSSD-2：（基础IOPS（min{3000+50×容量， 100000}） + 预配置IOPS） ≤ min{500×容量，1000000}  */  支持更新",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"delete_snap_with_ebs": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "设置快照是否随云硬盘删除，true表示随盘删除，false表示不随盘删除",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"image_id": schema.StringAttribute{
				Optional:    true,
				Description: "镜像ID，如果用镜像创建，只支持数据盘的私有镜像和共享镜像，所创建的数据盘的所在地域要与镜像源一致，容量不可小于镜像对应的磁盘容量。从镜像创建的数据盘不支持加密、ISCSI和FCSAN高级配置。",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"backup_id": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘备份ID参数，有以下限制：从备份创建盘仅支持VBD模式；新盘容量不能小于备份源盘容量；不支持配置加密属性（自动与备份源盘保持一致）；备份状态必须是可用。",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"labels": schema.ListNestedAttribute{
				Optional:    true,
				Description: "设置云硬盘标签，实际绑定标签的结果请查询云硬盘详情的labels返回值是否如预期。",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required:    true,
							Description: "标签的key值，长度不能超过32个字符。",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "标签的value值，长度不能超过32个字符。",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunEbs) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunEbsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := plan.RegionId.ValueString()
	projectId := plan.ProjectId.ValueString()
	azName := plan.AzName.ValueString()
	onDemand := plan.CycleType.ValueString() == business.OrderCycleTypeOnDemand

	diskMode, err := business.EbsDiskModeMap.FromOriginalScene(plan.Mode.ValueString(), business.EbsDiskModeMapScene1)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	diskType := plan.Type.ValueString()

	// 构建标签请求
	var labels []*ctebs2.EbsNewEbsLabelsRequest
	if plan.Labels != nil {
		for _, label := range plan.Labels {
			labels = append(labels, &ctebs2.EbsNewEbsLabelsRequest{
				Key:   label.Key.ValueString(),
				Value: label.Value.ValueString(),
			})
		}
	}

	params := &ctebs2.EbsNewEbsRequest{
		ClientToken:     uuid.NewString(),
		RegionID:        regionId,
		MultiAttach:     plan.MultiAttach.ValueBoolPointer(),
		IsEncrypt:       plan.Encrypted.ValueBoolPointer(),
		KmsUUID:         plan.KmsUuid.ValueString(),
		ProjectID:       projectId,
		DiskMode:        diskMode.(string),
		DiskType:        diskType,
		DiskName:        plan.Name.ValueString(),
		DiskSize:        plan.Size.ValueInt64(),
		OnDemand:        &onDemand,
		CycleType:       plan.CycleType.ValueString(),
		CycleCount:      int32(plan.CycleCount.ValueInt64()),
		ImageID:         plan.ImageId.ValueString(),
		AzName:          azName,
		ProvisionedIops: plan.ProvisionedIops.ValueInt64(),
		Labels:          labels,
		BackupID:        plan.BackupId.ValueString(),
	}
	if !plan.DeleteSnapWithEbs.IsUnknown() && !plan.DeleteSnapWithEbs.IsNull() {
		params.DeleteSnapWithEbs = plan.DeleteSnapWithEbs.ValueBoolPointer()
	}
	resp, err2 := c.meta.Apis.SdkCtEbsApis.EbsNewEbsApi.Do(ctx, c.meta.SdkCredential, params)
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}
	var id, masterOrderId string
	if resp.StatusCode == common.ErrorStatusCode && resp.ErrorCode != common.EbsOrderInProgress {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	if resp.ReturnObj != nil {
		masterOrderId = resp.ReturnObj.MasterOrderID
		if resp.ReturnObj.Resources != nil && len(resp.ReturnObj.Resources) > 0 {
			id = resp.ReturnObj.Resources[0].DiskID
		}
	}
	// 轮询结果
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	loop, err := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	if loop != nil && len(loop.Uuid) > 0 {
		id = loop.Uuid[0]
	}

	plan.Id = types.StringValue(id)
	plan.RegionId = types.StringValue(regionId)
	plan.ProjectId = types.StringValue(projectId)
	plan.AzName = types.StringValue(azName)
	plan.MasterOrderId = types.StringValue(masterOrderId)

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, ctyunRequestError := c.getAndMergeEbs(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
		return
	} else if plan.CycleType.ValueString() == business.OrderCycleTypeOnDemand {
		instance.CycleCount = plan.CycleCount
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEbs) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunEbsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	if !c.acquireAndSetIdIfOrderNotFinished(ctx, &state, response) {
		return
	}
	instance, err := c.getAndMergeEbs(ctx, state)
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

func (c *ctyunEbs) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan CtyunEbsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var state CtyunEbsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 判断名字是否相同
	if !plan.Name.Equal(state.Name) {
		_, err := c.meta.Apis.CtEbsApis.EbsChangeNameApi.Do(ctx, c.meta.Credential, &ctebs.EbsChangeNameRequest{
			RegionId: state.RegionId.ValueString(),
			DiskId:   state.Id.ValueString(),
			DiskName: plan.Name.ValueString(),
		})
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
	}
	// 判断硬盘大小是否相同，不同要走修改ebs接口
	err := c.ebsService.UpdateSize(ctx, state.Id.ValueString(), state.RegionId.ValueString(), int(state.Size.ValueInt64()), int(plan.Size.ValueInt64()))
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// 如果有IOPS相关字段，需要在这里添加IOPS更新逻辑
	if !plan.ProvisionedIops.Equal(state.ProvisionedIops) && !plan.ProvisionedIops.IsUnknown() {
		err := c.updateIops(ctx, plan, state)
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
		state.ProvisionedIops = plan.ProvisionedIops
	}

	// 如果删除策略字段变更，需要更新删除策略
	if !plan.DeleteSnapWithEbs.Equal(state.DeleteSnapWithEbs) && !plan.DeleteSnapWithEbs.IsUnknown() {
		err := c.updateDeletePolicy(ctx, plan, state)
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
		state.DeleteSnapWithEbs = plan.DeleteSnapWithEbs
	}

	instance, ctyunRequestError := c.getAndMergeEbs(ctx, state)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
		return
	} else if plan.CycleType.ValueString() == business.OrderCycleTypeOnDemand {
		instance.CycleCount = plan.CycleCount
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunEbs) updateIops(ctx context.Context, plan, state CtyunEbsConfig) (err error) {
	regionId := state.RegionId.ValueString()
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsUpdateIopsEbsApi.Do(ctx, c.meta.SdkCredential, &ctebs2.EbsUpdateIopsEbsRequest{
		ProvisionedIops: int32(plan.ProvisionedIops.ValueInt64()),
		DiskID:          state.Id.ValueString(),
		RegionID:        &regionId,
	})

	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", utils.SecString(resp.Message), utils.SecString(resp.Description))
		return
	}
	return
}

func (c *ctyunEbs) updateDeletePolicy(ctx context.Context, plan, state CtyunEbsConfig) (err error) {
	regionId := state.RegionId.ValueString()
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsSetDeletePolicyEbsApi.Do(ctx, c.meta.SdkCredential, &ctebs2.EbsSetDeletePolicyEbsRequest{
		RegionID:          regionId,
		DiskID:            state.Id.ValueString(),
		DeleteSnapWithEbs: plan.DeleteSnapWithEbs.ValueBool(),
	})

	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", utils.SecString(resp.Message), utils.SecString(resp.Description))
		return
	}

	return
}
func (c *ctyunEbs) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunEbsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	resp, err := c.meta.Apis.CtEbsApis.EbsDeleteApi.Do(ctx, c.meta.Credential, &ctebs.EbsDeleteRequest{
		RegionId:    state.RegionId.ValueString(),
		DiskId:      state.Id.ValueString(),
		ClientToken: uuid.NewString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	err2 := helper.RefundLoop(ctx, c.meta.Credential, resp.MasterOrderId)
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}
}

func (c *ctyunEbs) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ebsService = business.NewEbsService(meta)
}

func (c *ctyunEbs) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import %s.[导入配置名称] [id],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunEbsConfig

	var ID, regionId string
	// 根据分隔符数量判断是否输入了regionId
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		ID = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &ID, &regionId)
		if err != nil {
			return
		}
	}

	if ID == "" {
		err = fmt.Errorf("id不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}
	config.Id = types.StringValue(ID)
	config.RegionId = types.StringValue(regionId)
	instance, err := c.getAndMergeEbs(ctx, config)
	if err != nil {
		return
	}
	if instance == nil {
		err = common.ResourceNotExistError
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

// getAndMergeEbs 查询ebs
func (c *ctyunEbs) getAndMergeEbs(ctx context.Context, cfg CtyunEbsConfig) (*CtyunEbsConfig, error) {
	regionId := cfg.RegionId.ValueString()
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsQueryEbsByIDApi.Do(ctx, c.meta.SdkCredential, &ctebs2.EbsQueryEbsByIDRequest{
		RegionID: regionId, // 修正类型，API需要string而非*string
		DiskID:   cfg.Id.ValueString(),
	})
	if err != nil {
		return nil, err
	} else if resp == nil {
		return nil, common.InvalidReturnObjError
	} else if resp.StatusCode != common.NormalStatusCode {
		if resp.ErrorCode == common.EbsEbsInfoDataDamaged || resp.ErrorCode == common.EbsEbsInfoNotExists || resp.ReturnObj == nil {
			// 磁盘没了
			return nil, nil
		}
		return nil, fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
	}
	obj := resp.ReturnObj
	if utils.SecBool(obj.IsSystemVolume) {
		return nil, fmt.Errorf("不支持系统盘")
	}
	diskMode, err2 := business.EbsDiskModeMap.ToOriginalScene(obj.DiskMode, business.EbsDiskModeMapScene1)
	if err2 != nil {
		return nil, err2
	}
	cfg.Type = types.StringValue(obj.DiskType)
	cfg.Name = types.StringValue(obj.DiskName)
	cfg.Id = types.StringValue(obj.DiskID)
	cfg.Size = types.Int64Value(obj.DiskSize)
	cfg.Mode = types.StringValue(diskMode.(string))
	cfg.Status = types.StringValue(obj.DiskStatus)
	cfg.ExpireTime = types.StringValue(utils.FromUnixToUTC(obj.ExpireTime))
	cfg.CreateTime = types.StringValue(utils.FromUnixToUTC(obj.CreateTime))
	cfg.AzName = types.StringValue(obj.AzName)
	cfg.ProjectId = types.StringValue(obj.ProjectID)

	// 处理可选的布尔字段
	if obj.MultiAttach != nil {
		cfg.MultiAttach = types.BoolValue(*obj.MultiAttach)
	} else {
		cfg.MultiAttach = types.BoolValue(false)
	}

	if obj.IsEncrypt != nil {
		cfg.Encrypted = types.BoolValue(*obj.IsEncrypt)
	} else {
		cfg.Encrypted = types.BoolValue(false)
	}

	cfg.KmsUuid = types.StringValue(obj.KmsUUID)

	// 在 getAndMergeEbs 方法中正确设置周期相关字段
	if obj.OnDemand != nil && *obj.OnDemand {
		cfg.CycleType = types.StringValue("on_demand")
	} else {
		cfg.CycleType = types.StringValue(obj.CycleType)
	}

	// 正确处理 CycleCount
	if obj.CycleCount > 0 {
		cfg.CycleCount = types.Int64Value(int64(obj.CycleCount))
	}

	// 处理IOPS字段
	//如果ProvisionedIops为0，那么就不设置ProvisionedIops
	if obj.ProvisionedIops > 0 {
		cfg.ProvisionedIops = types.Int64Value(obj.ProvisionedIops)
	} else {
		cfg.ProvisionedIops = types.Int64Null()
	}

	// 处理删除快照策略字段
	if obj.DeleteSnapWithEbs == "true" {
		cfg.DeleteSnapWithEbs = types.BoolValue(true)
	} else {
		cfg.DeleteSnapWithEbs = types.BoolValue(false)
	}

	return &cfg, nil

}

// getMasterOrderIdIfOrderInProgress 获取masterOrderId
func (c *ctyunEbs) getMasterOrderIdIfOrderInProgress(err ctyunsdk.CtyunRequestError) (string, error) {
	resp := struct {
		MasterOrderId string `json:"masterOrderID"`
		MasterOrderNo string `json:"masterOrderNO"`
	}{}
	if err.CtyunResponse() == nil {
		return "", err
	}
	_, err = err.CtyunResponse().ParseByStandardModel(&resp)
	if err != nil {
		return "", err
	}
	return resp.MasterOrderId, err
}

// acquireIdIfOrderNotFinished 重新获取id，如果前订单状态有问题需要重新轮询
// 返回值：数据是否有效
func (c *ctyunEbs) acquireAndSetIdIfOrderNotFinished(ctx context.Context, state *CtyunEbsConfig, response *resource.ReadResponse) bool {
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

	// 成功把id恢复出来
	state.Id = types.StringValue(resp.Uuid[0])
	response.State.Set(ctx, state)
	return true
}

type CtyunEbsConfig struct {
	Name              types.String `tfsdk:"name"`
	Mode              types.String `tfsdk:"mode"`
	Type              types.String `tfsdk:"type"`
	Size              types.Int64  `tfsdk:"size"`
	CycleType         types.String `tfsdk:"cycle_type"`
	CycleCount        types.Int64  `tfsdk:"cycle_count"`
	MasterOrderId     types.String `tfsdk:"master_order_id"`
	Id                types.String `tfsdk:"id"`          // 磁盘ID
	Status            types.String `tfsdk:"status"`      // 云硬盘使用状态 deleting/creating/detaching，具体请参考云硬盘使用状态
	ExpireTime        types.String `tfsdk:"expire_time"` // 过期时刻
	CreateTime        types.String `tfsdk:"create_time"`
	MultiAttach       types.Bool   `tfsdk:"multi_attach"` // 是否共享云硬盘
	Encrypted         types.Bool   `tfsdk:"encrypted"`    // 是否加密盘
	KmsUuid           types.String `tfsdk:"kms_uuid"`     // 加密盘密钥UUID，是加密盘时才返回
	ProjectId         types.String `tfsdk:"project_id"`
	RegionId          types.String `tfsdk:"region_id"`
	AzName            types.String `tfsdk:"az_name"`
	ProvisionedIops   types.Int64  `tfsdk:"provisioned_iops"` // 预配置IOPS值
	DeleteSnapWithEbs types.Bool   `tfsdk:"delete_snap_with_ebs"`
	ImageId           types.String `tfsdk:"image_id"`  // 镜像ID
	BackupId          types.String `tfsdk:"backup_id"` // 云硬盘备份ID
	Labels            []Label      `tfsdk:"labels"`    // 云硬盘标签
}

type Label struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}
