package ebs

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
//	ctebsbackup "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebsbackup"
//	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
//	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
//	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
//	"github.com/google/uuid"
//	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
//	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
//	"github.com/hashicorp/terraform-plugin-framework/attr"
//	"github.com/hashicorp/terraform-plugin-framework/resource"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
//	"github.com/hashicorp/terraform-plugin-framework/types"
//	"time"
//)
//
///*
//云硬盘备份存储库
//*/
//
//func NewCtyunEbsBackupRepo() resource.Resource {
//	return &ctyunEbsBackupRepo{}
//}
//
//type ctyunEbsBackupRepo struct {
//	meta *common.CtyunMetadata
//}
//
//func (c *ctyunEbsBackupRepo) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
//	response.TypeName = request.ProviderTypeName + "_ebs_backup_repo"
//}
//
//type CtyunEbsBackupRepoConfig struct {
//	Id              types.String `tfsdk:"id"`
//	RegionID        types.String `tfsdk:"region_id"`
//	ProjectID       types.String `tfsdk:"project_id"`
//	RepositoryName  types.String `tfsdk:"name"`
//	CycleCount      types.Int64  `tfsdk:"cycle_count"`
//	CycleType       types.String `tfsdk:"cycle_type"`
//	Size            types.Int64  `tfsdk:"size"`
//	AutoRenewStatus types.Int64  `tfsdk:"auto_renew_status"`
//	OnDemand        types.String `tfsdk:"on_demand"`
//
//	Status     types.String  `tfsdk:"status"`
//	FreeSize   types.Float64 `tfsdk:"free_size"`
//	UsedSize   types.Int64   `tfsdk:"used_size"`
//	CreatedAt  types.Int32   `tfsdk:"created_at"`
//	ExpiredAt  types.Int32   `tfsdk:"expired_at"`
//	Expired    types.Bool    `tfsdk:"expired"`
//	Freeze     types.Bool    `tfsdk:"freeze"`
//	Paas       types.Bool    `tfsdk:"paas"`
//	BackupList types.List    `tfsdk:"backup_list"`
//}
//
//func (c *ctyunEbsBackupRepo) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
//	response.Schema = schema.Schema{
//		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026752/10212971**`,
//		Attributes: map[string]schema.Attribute{
//			"id": schema.StringAttribute{
//				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
//				Computed:      true,
//				Description:   "云硬盘备份存储库id",
//			},
//			"region_id": schema.StringAttribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//				Validators: []validator.String{
//					stringvalidator.UTF8LengthAtLeast(1),
//				},
//				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
//			},
//			"project_id": schema.StringAttribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看创建企业项目了解如何创建企业项目。注：默认值为\"0\"",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//				Validators: []validator.String{
//					validator2.Project(),
//				},
//				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
//			},
//			"name": schema.StringAttribute{
//				Required:    true,
//				Description: "云硬盘备份存储库名称，满足以下规则：长度为2-63字符，头尾不支持输入空格。支持更新",
//				Validators: []validator.String{
//					stringvalidator.LengthBetween(2, 63),
//				},
//			},
//			"cycle_count": schema.Int64Attribute{
//				Required:    true,
//				Description: "订购时长，该参数需要与cycleType一同使用。注：最长订购周期为60个月（5年）",
//				Validators: []validator.Int64{
//					int64validator.Between(1, 60),
//				},
//				PlanModifiers: []planmodifier.Int64{
//					int64planmodifier.RequiresReplace(),
//				},
//			},
//			"cycle_type": schema.StringAttribute{
//				Required:    true,
//				Description: "订购周期类型，取值范围：MONTH：按月，YEAR：按年。最长订购周期为5年",
//				Validators: []validator.String{
//					stringvalidator.OneOf("MONTH", "YEAR"),
//				},
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"size": schema.Int64Attribute{
//				Optional:    true,
//				Description: "云硬盘备份存储库容量，单位GB，取值范围：[100-1024000]，默认值100。支持更新",
//				Validators: []validator.Int64{
//					int64validator.Between(100, 1024000),
//				},
//			},
//			"auto_renew_status": schema.Int64Attribute{
//				Optional:    true,
//				Description: "是否自动续订，取值范围：0（不续费），1（自动续费）。注：按月购买，自动续订周期为1个月；按年购买，自动续订周期为1年",
//				Validators: []validator.Int64{
//					int64validator.OneOf(0, 1),
//				},
//				PlanModifiers: []planmodifier.Int64{
//					int64planmodifier.RequiresReplace(),
//				},
//			},
//			"on_demand": schema.StringAttribute{
//				Optional:    true,
//				Description: "是否按需下单，取值范围：true：是，false：否，默认为false",
//				Validators: []validator.String{
//					stringvalidator.OneOf("true", "false"),
//				},
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//
//			// 返回字段
//			"status": schema.StringAttribute{
//				Computed:    true,
//				Description: "云硬盘存储库状态，expired: 已到期，active: 可用",
//			},
//			"free_size": schema.Float64Attribute{
//				Computed:    true,
//				Description: "云硬盘存储库剩余大小，单位GB(废弃该字段)",
//			},
//			"used_size": schema.Int64Attribute{
//				Computed:    true,
//				Description: "云硬盘存储库使用大小，单位Byte",
//			},
//			"created_at": schema.Int32Attribute{
//				Computed:    true,
//				Description: "创建时间",
//			},
//			"expired_at": schema.Int32Attribute{
//				Computed:    true,
//				Description: "到期时间",
//			},
//			"expired": schema.BoolAttribute{
//				Computed:    true,
//				Description: "存储库是否到期",
//			},
//			"freeze": schema.BoolAttribute{
//				Computed:    true,
//				Description: "是否冻结",
//			},
//			"paas": schema.BoolAttribute{
//				Computed:    true,
//				Description: "是否支持PAAS",
//			},
//			"backup_list": schema.ListAttribute{
//				ElementType: types.StringType,
//				Computed:    true,
//				Description: "存储库下可用的备份ID列表",
//			},
//		},
//	}
//}
//
//func (c *ctyunEbsBackupRepo) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var plan CtyunEbsBackupRepoConfig
//	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//
//	// 实际创建
//	err = c.create(ctx, plan)
//	if err != nil {
//		return
//	}
//	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
//
//	//轮询状态为可用状态
//	err = c.StartedLoop(ctx, &plan)
//	if err != nil {
//		return
//	}
//	// 查询信息
//	err = c.getAndMerge(ctx, &plan)
//	if err != nil {
//		return
//	}
//
//	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
//}
//
//// getAndMerge 查询
//func (c *ctyunEbsBackupRepo) getAndMerge(ctx context.Context, cfg *CtyunEbsBackupRepoConfig) (err error) {
//	params := &ctebsbackup.EbsbackupListBackupRepoRequest{
//		RegionID:     cfg.RegionID.ValueString(),
//		RepositoryID: cfg.Id.ValueString(),
//	}
//	// 调用API
//	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupListBackupRepoApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode == common.ErrorStatusCode {
//		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	} else if resp.ReturnObj.TotalCount != 1 {
//		err = common.InvalidReturnObjResultsError
//		return
//	}
//
//	//资源返回内容更新
//	result := resp.ReturnObj.RepositoryList[0]
//	cfg.RepositoryName = types.StringValue(result.RepositoryName)
//	cfg.Status = types.StringValue(result.Status)
//	cfg.Size = types.Int64Value(int64(result.Size))
//	cfg.FreeSize = types.Float64Value(float64(result.FreeSize))
//	cfg.UsedSize = types.Int64Value(int64(result.UsedSize))
//	cfg.CreatedAt = types.Int32Value(result.CreatedTime)
//	cfg.ExpiredAt = types.Int32Value(result.ExpiredTime)
//	cfg.Expired = types.BoolValue(*result.Expired)
//	cfg.Freeze = types.BoolValue(*result.Freeze)
//	cfg.Paas = types.BoolValue(*result.Paas)
//	// 处理备份列表
//	backupList := make([]attr.Value, len(result.BackupList))
//	for i, backupID := range result.BackupList {
//		backupList[i] = types.StringValue(backupID)
//	}
//	cfg.BackupList, _ = types.ListValue(types.StringType, backupList)
//
//	return
//}
//
//func (c *ctyunEbsBackupRepo) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var state CtyunEbsBackupRepoConfig
//	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// 查询远端
//	err = c.getAndMerge(ctx, &state)
//	if err != nil {
//		return
//	}
//
//	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
//}
//
//func (c *ctyunEbsBackupRepo) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	// tf文件中的
//	var plan CtyunEbsBackupRepoConfig
//	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// state中的
//	var state CtyunEbsBackupRepoConfig
//	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	if plan.Size != state.Size {
//		// 更新
//		err = c.updateSize(ctx, plan, state)
//		if err != nil {
//			return
//		}
//	}
//	if plan.RepositoryName != state.RepositoryName {
//		// 更新
//		err = c.updateName(ctx, plan, state)
//	}
//
//	// 查询远端信息
//	err = c.getAndMerge(ctx, &state)
//	if err != nil {
//		return
//	}
//	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
//}
//
//func (c *ctyunEbsBackupRepo) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var state CtyunEbsBackupRepoConfig
//	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// 删除
//	err = c.delete(ctx, state)
//	if err != nil {
//		return
//	}
//}
//
//func (c *ctyunEbsBackupRepo) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
//	if request.ProviderData == nil {
//		return
//	}
//	meta := request.ProviderData.(*common.CtyunMetadata)
//	c.meta = meta
//}
//
//// create 创建
//func (c *ctyunEbsBackupRepo) create(ctx context.Context, plan CtyunEbsBackupRepoConfig) (err error) {
//
//	params := &ctebsbackup.EbsbackupCreateRepoRequest{
//		RegionID:        plan.RegionID.ValueString(),
//		ProjectID:       plan.ProjectID.ValueString(),
//		RepositoryName:  plan.RepositoryName.ValueString(),
//		CycleCount:      int32(plan.CycleCount.ValueInt64()),
//		CycleType:       plan.CycleType.ValueString(),
//		Size:            int32(plan.Size.ValueInt64()),
//		AutoRenewStatus: int32(plan.AutoRenewStatus.ValueInt64()),
//		OnDemand:        plan.OnDemand.ValueString(),
//		ClientToken:     uuid.NewString(),
//	}
//
//	// 创建实例
//	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupCreateRepoApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode == common.ErrorStatusCode {
//		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//
//	var masterOrderId string
//
//	if resp.ReturnObj.MasterOrderID != "" {
//		masterOrderId = resp.ReturnObj.MasterOrderID
//	}
//
//	// 轮询订单状态
//	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
//	loop, err2 := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId)
//	if err2 != nil {
//		return err2
//	}
//	// 最后设置id
//	id := loop.Uuid[0]
//	plan.Id = types.StringValue(id)
//
//	return
//}
//
//// updateSize 扩容云硬盘备份存储库
//func (c *ctyunEbsBackupRepo) updateSize(ctx context.Context, plan, state CtyunEbsBackupRepoConfig) (err error) {
//	if plan.Size.Equal(state.Size) {
//		return
//	}
//	params := &ctebsbackup.EbsbackupResizeRepoRequest{
//		RegionID:     plan.RegionID.ValueString(),
//		RepositoryID: state.Id.ValueString(),
//		Size:         int32(plan.Size.ValueInt64()),
//		ClientToken:  uuid.NewString(),
//	}
//	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupResizeRepoApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode == common.ErrorStatusCode {
//		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//
//	var masterOrderId string
//	if resp.ReturnObj.MasterOrderID != "" {
//		masterOrderId = resp.ReturnObj.MasterOrderID
//	}
//	// 轮询订单状态
//	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
//	_, err2 := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId)
//	if err2 != nil {
//		return err2
//	}
//	return
//}
//
//func (c *ctyunEbsBackupRepo) updateName(ctx context.Context, plan, state CtyunEbsBackupRepoConfig) (err error) {
//	if plan.Size.Equal(state.Size) {
//		return
//	}
//	params := &ctebsbackup.EbsbackupUpdateEbsBackupRepoRequest{
//		RegionID:       plan.RegionID.ValueString(),
//		RepositoryID:   state.Id.ValueString(),
//		RepositoryName: plan.RepositoryName.ValueString(),
//	}
//	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupUpdateEbsBackupRepoApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode == common.ErrorStatusCode {
//		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
//		return
//	} else if resp.ReturnObj == nil {
//		err = common.InvalidReturnObjError
//		return
//	}
//
//	return
//}
//
//// delete 删除
//func (c *ctyunEbsBackupRepo) delete(ctx context.Context, plan CtyunEbsBackupRepoConfig) (err error) {
//	params := &ctebsbackup.EbsbackupDeleteRepoRequest{
//		RegionID:     plan.RegionID.ValueString(),
//		RepositoryID: plan.Id.ValueString(),
//		ClientToken:  uuid.NewString(),
//	}
//	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupDeleteRepoApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode == common.ErrorStatusCode {
//		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
//		return
//	}
//
//	var masterOrderId string
//	if resp.ReturnObj.MasterOrderID != "" {
//		masterOrderId = resp.ReturnObj.MasterOrderID
//	}
//	// 轮询订单状态
//	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
//	_, err2 := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId)
//	if err2 != nil {
//		return err2
//	}
//	return
//}
//
//func (c *ctyunEbsBackupRepo) StartedLoop(ctx context.Context, state *CtyunEbsBackupRepoConfig, loopCount ...int) (err error) {
//	count := 30
//	if len(loopCount) > 0 {
//		count = loopCount[0]
//	}
//	retryer, err := business.NewRetryer(time.Second*30, count)
//	if err != nil {
//		return
//	}
//	result := retryer.Start(
//		func(currentTime int) bool {
//			// 获取实例详情
//			params := &ctebsbackup.EbsbackupListBackupRepoRequest{
//				RegionID:     state.RegionID.ValueString(),
//				RepositoryID: state.Id.ValueString(),
//			}
//			// 调用API
//			resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupListBackupRepoApi.Do(ctx, c.meta.SdkCredential, params)
//			if err != nil {
//				return false
//			} else if resp.StatusCode == common.ErrorStatusCode {
//				err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
//				return false
//			} else if resp.ReturnObj == nil {
//				err = common.InvalidReturnObjError
//				return false
//			} else if resp.ReturnObj.CurrentCount == 0 {
//				err = fmt.Errorf("no details found for ID: %s", state.Id.ValueString())
//				return false
//			}
//
//			runningStatus := resp.ReturnObj.RepositoryList[0].Status
//			if runningStatus == business.EcsBackupRepoStatusActive {
//				return false
//			}
//			return true
//		},
//	)
//	if result.ReturnReason == business.ReachMaxLoopTime {
//		return errors.New("轮询已达最大次数，资源仍未到达启动状态！")
//	}
//	return
//}
//
//func (c *ctyunEbsBackupRepo) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var cfg CtyunEbsBackupRepoConfig
//	var id, regionID string
//	err = terraform_extend.Split(request.ID, &id, &regionID)
//	if err != nil {
//		return
//	}
//	cfg.RegionID = types.StringValue(regionID)
//	cfg.Id = types.StringValue(id)
//	// 查询远端
//	err = c.getAndMerge(ctx, &cfg)
//	if err != nil {
//		return
//	}
//
//	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
//}
