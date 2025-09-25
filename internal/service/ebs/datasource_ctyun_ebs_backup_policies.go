package ebs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctebs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebsbackup"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunEbsBackupPolicies{}
	_ datasource.DataSourceWithConfigure = &ctyunEbsBackupPolicies{}
)

type ctyunEbsBackupPolicies struct {
	meta *common.CtyunMetadata
}

func NewCtyunEbsBackupPolicies() datasource.DataSource {
	return &ctyunEbsBackupPolicies{}
}

func (c *ctyunEbsBackupPolicies) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs_backup_policies"
}

type ctyunEbsBackupPoliciesModel struct {
	RegionID              types.String          `tfsdk:"region_id"`
	PolicyID              types.String          `tfsdk:"id"`
	PolicyName            types.String          `tfsdk:"name"`
	Status                types.Int64           `tfsdk:"status"`
	CycleType             types.String          `tfsdk:"cycle_type"`
	CycleDay              types.Int64           `tfsdk:"cycle_day"`
	CycleWeek             types.String          `tfsdk:"cycle_week"`
	Time                  types.String          `tfsdk:"time"`
	RetentionType         types.String          `tfsdk:"retention_type"`
	RetentionNum          types.Int64           `tfsdk:"retention_num"`
	RetentionDay          types.Int64           `tfsdk:"retention_day"`
	ResourceCount         types.Int64           `tfsdk:"resource_count"`
	ResourceIDs           types.String          `tfsdk:"resource_ids"`
	ProjectID             types.String          `tfsdk:"project_id"`
	FullBackupInterval    types.Int64           `tfsdk:"full_backup_interval"`
	AdvRetentionStatus    types.Bool            `tfsdk:"adv_retention_status"`
	RemainFirstOfCurMonth types.Bool            `tfsdk:"remain_first_of_cur_month"`
	RepositoryList        []repositoryListModel `tfsdk:"repository_list"`
}

type ctyunEbsBackupPoliciesConfig struct {
	RegionID       types.String                  `tfsdk:"region_id"`
	ProjectID      types.String                  `tfsdk:"project_id"`
	PolicyName     types.String                  `tfsdk:"name"`
	PolicyID       types.String                  `tfsdk:"id"`
	PageNo         types.Int32                   `tfsdk:"page_no"`
	PageSize       types.Int32                   `tfsdk:"page_size"`
	BackupPolicies []ctyunEbsBackupPoliciesModel `tfsdk:"backup_policies"`
}

type repositoryListModel struct {
	RepositoryID   types.String `tfsdk:"repository_id"`
	RepositoryName types.String `tfsdk:"repository_name"`
}

func (c *ctyunEbsBackupPolicies) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026752/10628749`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘备份策略名称。满足以下规则：只能由数字、英文字母、中划线-、下划线_、点.组成，长度为2-64字符",
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘备份策略ID，32字符",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "页码，取值范围：正整数（≥1），注：默认值为1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页记录数目，取值范围：[1~50]，默认值：10，单页最大记录不超过50",
			},
			"backup_policies": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源池ID",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘备份策略ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘备份策略名称",
						},
						"status": schema.Int64Attribute{
							Computed:    true,
							Description: "是否启用策略，取值范围：0：停用，1：启用",
						},
						"cycle_type": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘备份周期类型，取值范围：day：按天备份，week：按星期备份",
						},
						"cycle_day": schema.Int64Attribute{
							Computed:    true,
							Description: "只有cycleType为day时返回备份周期值",
						},
						"cycle_week": schema.StringAttribute{
							Computed:    true,
							Description: "只有cycleType为week时返回备份周期，取值范围：0-6代表星期日-星期六，如果一周有多天备份，以逗号隔开",
						},
						"time": schema.StringAttribute{
							Computed:    true,
							Description: "备份整点时间，取值范围：0-23，如果一天内多个时间节点备份，以逗号隔开",
						},
						"retention_type": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘备份保留类型，取值范围：date：按时间保留，num：按数量保留，all：永久保留",
						},
						"retention_num": schema.Int64Attribute{
							Computed:    true,
							Description: "只有retentionType为num时返回保留数量值",
						},
						"retention_day": schema.Int64Attribute{
							Computed:    true,
							Description: "只有retentionType为date时返回保留天数值",
						},
						"resource_count": schema.Int64Attribute{
							Computed:    true,
							Description: "策略已绑定的云硬盘数量",
						},
						"resource_ids": schema.StringAttribute{
							Computed:    true,
							Description: "策略已绑定的云硬盘ID，以逗号分隔",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "企业项目ID",
						},
						"remain_first_of_cur_month": schema.BoolAttribute{
							Computed:    true,
							Description: "是否保留每个月第一个备份，在retentionType为num时可设置，默认false",
						},
						"full_backup_interval": schema.Int64Attribute{
							Computed:    true,
							Description: "是否启用周期性全量备份。-1代表不开启，默认为-1；取值范围为[-1,100]，即每执行n次增量备份后，执行一次全量备份；若传入为0，代表每一次均为全量备份。",
						},
						"adv_retention_status": schema.BoolAttribute{
							Computed:    true,
							Description: "是否开启高级保留策略，false（不启用），true(启用)，默认值为false。需校验云硬盘备份保留类型（retentionType），若保留类型为按数量保存（num），可开启高级保留策略；若保留类型为date（按时间保存）或all（永久保存），不可开启高级保留策略。",
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
				}},
		},
	}
}

func (c *ctyunEbsBackupPolicies) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config ctyunEbsBackupPoliciesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	config.RegionID = types.StringValue(regionId)
	config.BackupPolicies = []ctyunEbsBackupPoliciesModel{}
	// 组装请求体
	params := &ctebs2.EbsbackupListBackupPolicyRequest{
		RegionID:   config.RegionID.ValueString(),
		ProjectID:  config.ProjectID.ValueString(),
		PolicyName: config.PolicyName.ValueString(),
		PolicyID:   config.PolicyID.ValueString(),
		PageNo:     config.PageNo.ValueInt32(),
		PageSize:   config.PageSize.ValueInt32(),
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
	}
	// 解析返回值
	for _, policy := range resp.ReturnObj.PolicyList {
		// 处理备份库列表
		repositoryList := make([]repositoryListModel, len(policy.RepositoryList))
		for i, repo := range policy.RepositoryList {
			repoItem := repositoryListModel{
				RepositoryID:   types.StringValue(repo.RepositoryID),
				RepositoryName: types.StringValue(repo.RepositoryName),
			}
			repositoryList[i] = repoItem
		}
		item := ctyunEbsBackupPoliciesModel{
			RegionID:              types.StringValue(policy.RegionID),
			PolicyID:              types.StringValue(policy.PolicyID),
			PolicyName:            types.StringValue(policy.PolicyName),
			Status:                types.Int64Value(int64(policy.Status)),
			CycleType:             types.StringValue(policy.CycleType),
			CycleDay:              types.Int64Value(int64(policy.CycleDay)),
			CycleWeek:             types.StringValue(policy.CycleWeek),
			Time:                  types.StringValue(policy.Time),
			RetentionType:         types.StringValue(policy.RetentionType),
			RetentionNum:          types.Int64Value(int64(policy.RetentionNum)),
			RetentionDay:          types.Int64Value(int64(policy.RetentionDay)),
			ResourceCount:         types.Int64Value(int64(policy.BindedDiskCount)),
			ResourceIDs:           types.StringValue(policy.BindedDiskIDs),
			ProjectID:             types.StringValue(policy.ProjectID),
			FullBackupInterval:    types.Int64Value(int64(policy.FullBackupInterval)),
			AdvRetentionStatus:    types.BoolPointerValue(policy.AdvRetentionStatus),
			RepositoryList:        repositoryList,
			RemainFirstOfCurMonth: types.BoolPointerValue(policy.RemainFirstOfCurMonth),
		}

		config.BackupPolicies = append(config.BackupPolicies, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunEbsBackupPolicies) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
