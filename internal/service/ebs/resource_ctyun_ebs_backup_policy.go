package ebs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctebs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebsbackup"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
)

/*
云硬盘备份策略
*/

func NewCtyunEbsBackupPolicy() resource.Resource {
	return &ctyunEbsBackupPolicy{}
}

type ctyunEbsBackupPolicy struct {
	meta *common.CtyunMetadata
}

func (c *ctyunEbsBackupPolicy) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs_backup_policy"
}

type CtyunEbsBackupPolicyConfig struct {
	Id                    types.String                      `tfsdk:"id"`
	RegionID              types.String                      `tfsdk:"region_id"`
	ProjectID             types.String                      `tfsdk:"project_id"`
	Name                  types.String                      `tfsdk:"name"`
	CycleType             types.String                      `tfsdk:"cycle_type"`
	CycleDay              types.Int64                       `tfsdk:"cycle_day"`
	CycleWeek             types.String                      `tfsdk:"cycle_week"`
	Time                  types.String                      `tfsdk:"time"`
	Status                types.Int64                       `tfsdk:"status"`
	RetentionType         types.String                      `tfsdk:"retention_type"`
	RetentionDay          types.Int64                       `tfsdk:"retention_day"`
	RetentionNum          types.Int64                       `tfsdk:"retention_num"`
	FullBackupInterval    types.Int32                       `tfsdk:"full_backup_interval"`
	AdvRetentionStatus    types.Bool                        `tfsdk:"adv_retention_status"`
	AdvRetention          *CtyunEbsBackupPolicyAdvRetention `tfsdk:"adv_retention"`
	ResourceIDs           types.String                      `tfsdk:"resource_ids"`
	RepositoryList        types.List                        `tfsdk:"repository_list"`
	RemainFirstOfCurMonth types.Bool                        `tfsdk:"remain_first_of_cur_month"`
}

type CtyunEbsBackupPolicyAdvRetention struct {
	AdvDay   types.Int64 `tfsdk:"adv_day"`
	AdvWeek  types.Int64 `tfsdk:"adv_week"`
	AdvMonth types.Int64 `tfsdk:"adv_month"`
	AdvYear  types.Int64 `tfsdk:"adv_year"`
}

func (c *ctyunEbsBackupPolicy) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026752/10037448`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "云硬盘备份策略id",
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
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看创建企业项目了解如何创建企业项目。注：默认值为\"0\"",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "云硬盘备份策略名称。满足以下规则：只能由数字、英文字母、中划线-、下划线_、点.组成，长度为2-64字符。注：在所有资源池不可重复。支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 64),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[0-9a-zA-Z\-_\.]+$`),
						"只能由数字、英文字母、中划线-、下划线_、点.组成",
					),
				},
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "云硬盘备份周期类型，取值范围：day（按天备份）week（按星期备份）。支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf("day", "week"),
				},
			},
			"cycle_day": schema.Int64Attribute{
				Optional:    true,
				Description: "备份周期（天），取值范围：[1, 30]，默认值为1。注：只有cycleType为day时有效。支持更新",
				Validators: []validator.Int64{
					int64validator.Between(1, 30),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"cycle_week": schema.StringAttribute{
				Optional:    true,
				Description: "备份周期（星期），星期取值范围：0~6（代表周几，其中0为周日），默认值是0。注：只有cycleType为week时有效；如果一周有多天备份，以逗号隔开（如周日周三进行快照，则填写\"0,3\"）。支持更新",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[0-6,]+$`),
						"星期格式错误",
					),
				},
			},
			"time": schema.StringAttribute{
				Required:    true,
				Description: "备份整点时间，时间取值范围：0~23。注：如果一天内多个时间节点备份，以逗号隔开（如11点15点进行快照，则填写\"11,15\"），默认值12 。支持更新",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[0-9,]+$`),
						"时间格式错误",
					),
				},
			},
			"status": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "备份策略状态，是否启用，取值范围：0（不启用），1（启用）。注：默认值0（不启用）。支持更新",
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"retention_type": schema.StringAttribute{
				Required:    true,
				Description: "云硬盘备份保留类型，取值范围：date（按时间保存），num（按数量保存），all（永久保存）。支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf("date", "num", "all"),
				},
			},
			"retention_day": schema.Int64Attribute{
				Optional:    true,
				Description: "云硬盘备份保留天数，单位为天，取值范围：[1, 99999] ，默认值1。注：只有retentionType为date时有效。支持更新",
				Validators: []validator.Int64{
					int64validator.Between(1, 99999),
				},
			},
			"remain_first_of_cur_month": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "是否保留每个月第一个备份，在retentionType为num时可设置，默认false，支持更新",
			},
			"retention_num": schema.Int64Attribute{
				Optional:    true,
				Description: "云硬盘备份保留数量，取值范围：[1, 99999]，默认值1。注：只有retentionType为num时有效。支持更新",
				Validators: []validator.Int64{
					int64validator.Between(1, 99999),
				},
			},
			"full_backup_interval": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "是否启用周期性全量备份。-1代表不开启，默认为-1；取值范围为[-1,100]，即每执行n次增量备份后，执行一次全量备份；若传入为0，代表每一次均为全量备份。支持更新",
				Validators: []validator.Int32{
					int32validator.Between(-1, 100),
				},
			},
			"adv_retention_status": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "是否开启高级保留策略，false（不启用），true(启用)，默认值为false。需校验云硬盘备份保留类型（retentionType），若保留类型为按数量保存（num），可开启高级保留策略；若保留类型为date（按时间保存）或all（永久保存），不可开启高级保留策略。支持更新",
			},
			"resource_ids": schema.StringAttribute{
				Computed:    true,
				Description: "策略已绑定的云硬盘ID，以逗号分隔",
			},
			"repository_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"repository_id": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘备份库ID",
						},
						"repository_name": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘备份库名称",
						},
					},
				},
				Description: "策略已绑定的云硬盘备份库列表",
			},
		},
		Blocks: map[string]schema.Block{
			"adv_retention": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"adv_day": schema.Int64Attribute{
						Optional:    true,
						Description: "保留n天内，每天最新的一个备份，n为输入的数字。单位为天，取值范围：[0, 100]，默认值0 支持更新",
						Validators: []validator.Int64{
							int64validator.Between(0, 100),
						},
					},
					"adv_week": schema.Int64Attribute{
						Optional:    true,
						Description: "保留n周内，每周最新的一个备份，n为输入的数字。单位为周，取值范围：[0, 100]，默认值0 支持更新",
						Validators: []validator.Int64{
							int64validator.Between(0, 100),
						},
					},
					"adv_month": schema.Int64Attribute{
						Optional:    true,
						Description: "保留n月内，每月最新的一个备份，n为输入的数字。单位为月，取值范围：[0, 100]，默认值0 支持更新",
						Validators: []validator.Int64{
							int64validator.Between(0, 100),
						},
					},
					"adv_year": schema.Int64Attribute{
						Optional:    true,
						Description: "保留n年内，每年最新的一个备份，n为输入的数字。单位为年，取值范围：[0, 100]，默认值0 支持更新",
						Validators: []validator.Int64{
							int64validator.Between(0, 100),
						},
					},
				},
				Description: "高级保留策略内容，只有retentionType为num且advRetentionStatus为true才生效。支持更新",
			},
		},
	}
}

func (c *ctyunEbsBackupPolicy) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEbsBackupPolicyConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 校验创建动作的前置条件
	err = c.checkName(ctx, plan)
	if err != nil {
		return
	}

	// 实际创建
	id, err := c.create(ctx, plan)
	if err != nil {
		return
	}
	plan.Id = types.StringValue(id)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)

	// 查询信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

// checkName 校验名称是否重复
func (c *ctyunEbsBackupPolicy) checkName(ctx context.Context, plan CtyunEbsBackupPolicyConfig) (err error) {
	params := &ctebs2.EbsbackupListBackupPolicyRequest{
		RegionID:   plan.RegionID.ValueString(),
		PolicyName: plan.Name.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupListBackupPolicyApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj != nil && resp.ReturnObj.TotalCount > 0 {
		// 如果查询结果存在，说明名称已存在，返回错误
		err = fmt.Errorf("name '%s' already exists", plan.Name.ValueString())
		return
	}
	return
}

// getAndMerge 查询
func (c *ctyunEbsBackupPolicy) getAndMerge(ctx context.Context, cfg *CtyunEbsBackupPolicyConfig) (err error) {
	params := &ctebs2.EbsbackupListBackupPolicyRequest{
		RegionID: cfg.RegionID.ValueString(),
	}

	// 优先使用策略ID查询
	if !cfg.Id.IsNull() && cfg.Id.ValueString() != "" {
		params.PolicyID = cfg.Id.ValueString()
	} else {
		// 回退到名称查询（仅用于创建场景）
		params.PolicyName = cfg.Name.ValueString()
	}
	// 调用API
	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupListBackupPolicyApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.ReturnObj.CurrentCount != 1 {
		err = common.InvalidReturnObjResultsError
		return
	}

	//资源返回内容更新
	result := resp.ReturnObj.PolicyList[0]
	cfg.Id = types.StringValue(result.PolicyID)
	cfg.Name = types.StringValue(result.PolicyName)
	cfg.Status = types.Int64Value(int64(result.Status))
	cfg.CycleType = types.StringValue(result.CycleType)
	if result.CycleDay != 0 {
		cfg.CycleDay = types.Int64Value(int64(result.CycleDay))
	}
	if result.CycleWeek != "" {
		cfg.CycleWeek = types.StringValue(result.CycleWeek)
	}
	if result.Time != "" {
		cfg.Time = types.StringValue(result.Time)
	}
	cfg.RetentionType = types.StringValue(result.RetentionType)
	if result.RetentionDay != 0 {
		cfg.RetentionDay = types.Int64Value(int64(result.RetentionDay))
	}
	if result.RetentionNum != 0 {
		cfg.RetentionNum = types.Int64Value(int64(result.RetentionNum))
	}
	cfg.FullBackupInterval = types.Int32Value(result.FullBackupInterval)
	cfg.AdvRetentionStatus = types.BoolPointerValue(result.AdvRetentionStatus)
	cfg.AdvRetention = convertFromApiAdvRetention(result.AdvRetention)
	cfg.ResourceIDs = types.StringValue(result.BindedDiskIDs)
	cfg.RemainFirstOfCurMonth = types.BoolPointerValue(result.RemainFirstOfCurMonth)
	// 处理备份库列表
	if len(result.RepositoryList) > 0 {
		repositoryList := make([]repositoryListModel, len(result.RepositoryList))
		for i, repo := range result.RepositoryList {
			repoItem := repositoryListModel{
				RepositoryID:   types.StringValue(repo.RepositoryID),
				RepositoryName: types.StringValue(repo.RepositoryName),
			}
			repositoryList[i] = repoItem
		}

		// 转换为 types.List
		objectType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"repository_id":   types.StringType,
				"repository_name": types.StringType,
			},
		}

		repoObjects := make([]attr.Value, len(repositoryList))
		for i, repo := range repositoryList {
			obj, objDiags := types.ObjectValue(
				map[string]attr.Type{
					"repository_id":   types.StringType,
					"repository_name": types.StringType,
				},
				map[string]attr.Value{
					"repository_id":   repo.RepositoryID,
					"repository_name": repo.RepositoryName,
				},
			)
			if objDiags.HasError() {
				err = fmt.Errorf("failed to create repository object: %v", objDiags.Errors())
				return
			}
			repoObjects[i] = obj
		}

		repoList, listDiags := types.ListValue(objectType, repoObjects)
		if listDiags.HasError() {
			err = fmt.Errorf("failed to create repository list: %v", listDiags.Errors())
			return
		}
		cfg.RepositoryList = repoList
	} else {
		// 创建空列表
		objectType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"repository_id":   types.StringType,
				"repository_name": types.StringType,
			},
		}
		emptyList, listDiags := types.ListValue(objectType, []attr.Value{})
		if listDiags.HasError() {
			err = fmt.Errorf("failed to create empty repository list: %v", listDiags.Errors())
			return
		}
		cfg.RepositoryList = emptyList
	}

	return
}

func (c *ctyunEbsBackupPolicy) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbsBackupPolicyConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEbsBackupPolicy) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunEbsBackupPolicyConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunEbsBackupPolicyConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	} else {
		state.Name = types.StringValue(plan.Name.ValueString())
	}
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEbsBackupPolicy) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbsBackupPolicyConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunEbsBackupPolicy) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// create 创建
func (c *ctyunEbsBackupPolicy) create(ctx context.Context, plan CtyunEbsBackupPolicyConfig) (id string, err error) {

	params := &ctebs2.EbsbackupCreateEbsBackupPolicyRequest{
		RegionID:              plan.RegionID.ValueString(),
		ProjectID:             plan.ProjectID.ValueString(),
		PolicyName:            plan.Name.ValueString(),
		CycleType:             plan.CycleType.ValueString(),
		CycleDay:              int32(plan.CycleDay.ValueInt64()),
		CycleWeek:             plan.CycleWeek.ValueString(),
		Time:                  plan.Time.ValueString(),
		Status:                int32(plan.Status.ValueInt64()),
		RetentionType:         plan.RetentionType.ValueString(),
		RetentionDay:          int32(plan.RetentionDay.ValueInt64()),
		RetentionNum:          int32(plan.RetentionNum.ValueInt64()),
		FullBackupInterval:    plan.FullBackupInterval.ValueInt32(),
		AdvRetentionStatus:    aws.Bool(plan.AdvRetentionStatus.ValueBool()),
		AdvRetention:          convertToApiAdvRetention(plan.AdvRetention),
		RemainFirstOfCurMonth: aws.Bool(plan.RemainFirstOfCurMonth.ValueBool()),
	}

	// 创建实例
	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupCreateEbsBackupPolicyApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return
}

// update 修改启用停用云硬盘备份策略
func (c *ctyunEbsBackupPolicy) update(ctx context.Context, plan, state CtyunEbsBackupPolicyConfig) (err error) {
	params := &ctebs2.EbsbackupUpdateEbsBackupPolicyRequest{
		PolicyID:              state.Id.ValueString(),
		RegionID:              state.RegionID.ValueString(),
		PolicyName:            plan.Name.ValueString(),
		CycleType:             plan.CycleType.ValueString(),
		CycleDay:              int32(plan.CycleDay.ValueInt64()),
		CycleWeek:             plan.CycleWeek.ValueString(),
		Time:                  plan.Time.ValueString(),
		RetentionType:         plan.RetentionType.ValueString(),
		RetentionDay:          int32(plan.RetentionDay.ValueInt64()),
		RetentionNum:          int32(plan.RetentionNum.ValueInt64()),
		RemainFirstOfCurMonth: aws.Bool(plan.RemainFirstOfCurMonth.ValueBool()),
		FullBackupInterval:    plan.FullBackupInterval.ValueInt32(),
		AdvRetentionStatus:    aws.Bool(plan.AdvRetentionStatus.ValueBool()),
		AdvRetention:          convertToApiAdvRetentionUpdate(plan.AdvRetention),
	}
	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupUpdateEbsBackupPolicyApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return
}

// delete 删除
func (c *ctyunEbsBackupPolicy) delete(ctx context.Context, plan CtyunEbsBackupPolicyConfig) (err error) {
	params := &ctebs2.EbsbackupDeleteEbsBackupPolicyRequest{
		RegionID:  plan.RegionID.ValueString(),
		PolicyIDs: plan.Id.ValueString(),
	}
	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupDeleteEbsBackupPolicyApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}

	return
}

func (c *ctyunEbsBackupPolicy) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunEbsBackupPolicyConfig
	var id, regionID string
	err = terraform_extend.Split(request.ID, &id, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.Id = types.StringValue(id)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// convertToApiAdvRetention 将本地结构体转换为API结构体
func convertToApiAdvRetention(local *CtyunEbsBackupPolicyAdvRetention) *ctebs2.EbsbackupCreateEbsBackupPolicyAdvRetentionRequest {
	if local == nil {
		return nil
	}

	return &ctebs2.EbsbackupCreateEbsBackupPolicyAdvRetentionRequest{
		AdvDay:   int32(local.AdvDay.ValueInt64()),
		AdvWeek:  int32(local.AdvWeek.ValueInt64()),
		AdvMonth: int32(local.AdvMonth.ValueInt64()),
		AdvYear:  int32(local.AdvYear.ValueInt64()),
	}
}

// convertToApiAdvRetention 将本地结构体转换为API结构体
func convertToApiAdvRetentionUpdate(local *CtyunEbsBackupPolicyAdvRetention) *ctebs2.EbsbackupUpdateEbsBackupPolicyAdvRetentionRequest {
	if local == nil {
		return nil
	}

	return &ctebs2.EbsbackupUpdateEbsBackupPolicyAdvRetentionRequest{
		AdvDay:   int32(local.AdvDay.ValueInt64()),
		AdvWeek:  int32(local.AdvWeek.ValueInt64()),
		AdvMonth: int32(local.AdvMonth.ValueInt64()),
		AdvYear:  int32(local.AdvYear.ValueInt64()),
	}
}

// convertFromApiAdvRetention 将API结构体转换为本地结构体
func convertFromApiAdvRetention(api *ctebs2.EbsbackupListBackupPolicyReturnObjPolicyListAdvRetentionResponse) *CtyunEbsBackupPolicyAdvRetention {
	if api == nil {
		return nil
	}

	result := &CtyunEbsBackupPolicyAdvRetention{}

	// 只有当值不为0时才设置具体的值，否则保持为null
	if api.AdvDay != 0 {
		result.AdvDay = types.Int64Value(int64(api.AdvDay))
	} else {
		result.AdvDay = types.Int64Null()
	}

	if api.AdvWeek != 0 {
		result.AdvWeek = types.Int64Value(int64(api.AdvWeek))
	} else {
		result.AdvWeek = types.Int64Null()
	}

	if api.AdvMonth != 0 {
		result.AdvMonth = types.Int64Value(int64(api.AdvMonth))
	} else {
		result.AdvMonth = types.Int64Null()
	}

	if api.AdvYear != 0 {
		result.AdvYear = types.Int64Value(int64(api.AdvYear))
	} else {
		result.AdvYear = types.Int64Null()
	}
	return result
}
