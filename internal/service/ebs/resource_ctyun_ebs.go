package ebs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctebs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctebs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	ebsService *business.EbsService
}

func NewCtyunEbs() resource.Resource {
	return &ctyunEbs{}
}

func (c *ctyunEbs) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs"
}

func (c *ctyunEbs) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027696`,
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
				Description: "磁盘类型，sata：普通IO，sas：高IO，ssd：超高IO，ssd-genric：通用型SSD，fast-ssd：极速型SSD，不支持ISCSI模式；XSSD-0、XSSD-1、XSSD-2：X系列云硬盘，不支持加密，不支持ISCSI模式或FCSAN模式",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.EbsDiskTypes...),
				},
			},
			"size": schema.Int64Attribute{
				Required:    true,
				Description: "磁盘大小，单位GB，超高IO/高IO/极速型SSD/普通IO：取值范围[10, 32768]；XSSD-0：10GB-65536GB；XSSD-1：20GB-65536GB；XSSD-2：512GB-65536GB 支持更新（不支持缩容）",
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
					int64planmodifier.RequiresReplace(),
				},
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
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订购的受理单id",
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
			},
			"multi_attach": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否共享云硬盘",
			},
			"encrypted": schema.BoolAttribute{
				Computed:    true,
				Description: "是否加密盘； 共享盘、ISCSI模式磁盘、极速型SSD类型盘、XSSD系列盘不支持加密",
			},
			"kms_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "加密盘密钥UUID，是加密盘时才返回",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
				Default: defaults2.AcquireFromGlobalString(common.ExtraAzName, false),
			},

			"provisioned_iops": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "XSSD类型云硬盘的预配置IOPS值，最小值为1，最大值计算公式为“min(单盘最大IOPS，500*容量) - 基础性能IOPS”。 其他类型磁盘不支持此参数 具体取值范围如下：\n\t●XSSD-0：（基础IOPS（min{1800+12×容量， 10000}） + 预配置IOPS） ≤ min{500×容量，100000}\n\t●XSSD-1：（基础IOPS（min{1800+50×容量， 50000}） + 预配置IOPS） ≤ min{500×容量，100000}\n\t●XSSD-2：（基础IOPS（min{3000+50×容量， 100000}） + 预配置IOPS） ≤ min{500×容量，1000000}  */  支持更新",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"delete_snap_with_ebs": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "设置快照是否随云硬盘删除，true表示随盘删除，false表示不随盘删除",
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
	diskType, err := business.EbsDiskTypeMap.FromOriginalScene(plan.Type.ValueString(), business.EbsDiskTypeMapScene1)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// 构建标签请求
	var labels []*ctebs2.EbsNewEbsLabelsRequest
	if plan.Labels != nil {
		var labels []*ctebs2.EbsNewEbsLabelsRequest
		for _, label := range plan.Labels {
			labels = append(labels, &ctebs2.EbsNewEbsLabelsRequest{
				Key:   label.Key.ValueString(),
				Value: label.Value.ValueString(),
			})
		}
	}

	resp, err2 := c.meta.Apis.SdkCtEbsApis.EbsNewEbsApi.Do(ctx, c.meta.SdkCredential, &ctebs2.EbsNewEbsRequest{
		ClientToken:       uuid.NewString(),
		RegionID:          regionId,
		MultiAttach:       plan.MultiAttach.ValueBoolPointer(),
		IsEncrypt:         plan.Encrypted.ValueBoolPointer(),
		KmsUUID:           plan.KmsUuid.ValueString(),
		ProjectID:         projectId,
		DiskMode:          diskMode.(string),
		DiskType:          diskType.(string),
		DiskName:          plan.Name.ValueString(),
		DiskSize:          plan.Size.ValueInt64(),
		OnDemand:          &onDemand,
		CycleType:         plan.CycleType.ValueString(),
		CycleCount:        int32(plan.CycleCount.ValueInt64()),
		ImageID:           plan.ImageId.ValueString(),
		AzName:            azName,
		ProvisionedIops:   plan.ProvisionedIops.ValueInt64(),
		DeleteSnapWithEbs: plan.DeleteSnapWithEbs.ValueBoolPointer(),
		Labels:            labels,
		BackupID:          plan.BackupId.ValueString(),
	})

	var id, masterOrderId string
	if err2 == nil {

		if resp.StatusCode == common.ErrorStatusCode && resp.ErrorCode != common.EbsOrderInProgress {
			err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
		if resp.ReturnObj != nil && resp.ReturnObj.Resources != nil && len(resp.ReturnObj.Resources) > 0 {
			id = resp.ReturnObj.Resources[0].DiskID
		}
		masterOrderId = resp.ReturnObj.MasterOrderID
	}

	// 判断返回信息是否需要轮询
	if resp.ErrorCode != common.EbsOrderInProgress {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}

	// 轮询结果
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	loop, err := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	id = loop.Uuid[0]

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
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
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
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
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
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[project_id],[az_name],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunEbsConfig

	var ID, projectId, azName, regionId string
	// 根据分隔符数量判断是否输入了regionId,projectId,azName
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		projectId = c.meta.GetExtraIfEmpty(projectId, common.ExtraProjectId)
		azName = c.meta.GetExtraIfEmpty(azName, common.ExtraAzName)
		ID = request.ID
	} else if strings.Count(request.ID, common.ImportSeparator) == 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		projectId = c.meta.GetExtraIfEmpty(projectId, common.ExtraProjectId)
		err = terraform_extend.Split(request.ID, &ID, &azName)
		if err != nil {
			return
		}
	} else if strings.Count(request.ID, common.ImportSeparator) == 2 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &ID, &projectId, &azName)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &ID, &projectId, &azName, &regionId)
		if err != nil {
			return
		}
	}

	if ID == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	config.Id = types.StringValue(ID)
	config.RegionId = types.StringValue(regionId)
	config.ProjectId = types.StringValue(projectId)
	config.AzName = types.StringValue(azName)

	instance, err := c.getAndMergeEbs(ctx, config)
	if err != nil {
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
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
	diskMode, err2 := business.EbsDiskModeMap.ToOriginalScene(obj.DiskMode, business.EbsDiskModeMapScene1)
	if err2 != nil {
		return nil, err2
	}
	diskType, err2 := business.EbsDiskTypeMap.ToOriginalScene(obj.DiskType, business.EbsDiskTypeMapScene1)
	if err2 != nil {
		return nil, err2
	}
	cfg.Name = types.StringValue(obj.DiskName)
	cfg.Id = types.StringValue(obj.DiskID)
	cfg.Size = types.Int64Value(obj.DiskSize)
	cfg.Type = types.StringValue(diskType.(string))
	cfg.Mode = types.StringValue(diskMode.(string))
	cfg.Status = types.StringValue(obj.DiskStatus)
	cfg.ExpireTime = types.StringValue(utils.FromUnixToUTC(obj.ExpireTime))
	cfg.CreateTime = types.StringValue(utils.FromUnixToUTC(obj.CreateTime))

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
	} else {
		cfg.CycleCount = types.Int64Null()
	}

	// 处理IOPS字段
	cfg.ProvisionedIops = types.Int64Value(obj.ProvisionedIops)

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
