package ebs

//
//import (
//	"context"
//	"fmt"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
//	ctebsbackup "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebsbackup"
//	"github.com/hashicorp/terraform-plugin-framework/datasource"
//	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
//	"github.com/hashicorp/terraform-plugin-framework/types"
//)
//
//var (
//	_ datasource.DataSource              = &ctyunEbsBackupRepos{}
//	_ datasource.DataSourceWithConfigure = &ctyunEbsBackupRepos{}
//)
//
//type ctyunEbsBackupRepos struct {
//	meta *common.CtyunMetadata
//}
//
//func NewCtyunEbsBackupRepos() datasource.DataSource {
//	return &ctyunEbsBackupRepos{}
//}
//
//func (c *ctyunEbsBackupRepos) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
//	response.TypeName = request.ProviderTypeName + "_ebs_backup_repos"
//}
//
//type ctyunEbsBackupReposModel struct {
//	RegionID     types.String   `tfsdk:"region_id"`
//	RepositoryID types.String   `tfsdk:"id"`
//	ProjectID    types.String   `tfsdk:"project_id"`
//	Name         types.String   `tfsdk:"name"`
//	Status       types.String   `tfsdk:"status"`
//	Size         types.Int64    `tfsdk:"size"`
//	FreeSize     types.Float64  `tfsdk:"free_size"`
//	UsedSize     types.Int64    `tfsdk:"used_size"`
//	CreatedAt    types.Int32    `tfsdk:"created_at"`
//	ExpiredAt    types.Int32    `tfsdk:"expired_at"`
//	Expired      types.Bool     `tfsdk:"expired"`
//	Freeze       types.Bool     `tfsdk:"freeze"`
//	Paas         types.Bool     `tfsdk:"paas"`
//	BackupList   []types.String `tfsdk:"backup_list"`
//	BackupCount  types.Int64    `tfsdk:"backup_count"`
//}
//
//type ctyunEbsBackupReposConfig struct {
//	RegionID     types.String               `tfsdk:"region_id"`
//	ProjectID    types.String               `tfsdk:"project_id"`
//	Name         types.String               `tfsdk:"name"`
//	RepositoryID types.String               `tfsdk:"id"`
//	Status       types.String               `tfsdk:"status"`
//	PageNo       types.Int32                `tfsdk:"page_no"`
//	PageSize     types.Int32                `tfsdk:"page_size"`
//	HideExpire   types.Bool                 `tfsdk:"hide_expire"`
//	QueryContent types.String               `tfsdk:"query_content"`
//	Asc          types.Bool                 `tfsdk:"asc"`
//	Sort         types.String               `tfsdk:"sort"`
//	BackupRepos  []ctyunEbsBackupReposModel `tfsdk:"backup_repos"`
//	TotalCount   types.Int32                `tfsdk:"total_count"`
//	CurrentCount types.Int32                `tfsdk:"current_count"`
//}
//
//func (c *ctyunEbsBackupRepos) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
//	response.Schema = schema.Schema{
//		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026752/10212971`,
//		Attributes: map[string]schema.Attribute{
//			"region_id": schema.StringAttribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
//			},
//			"project_id": schema.StringAttribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。",
//			},
//			"name": schema.StringAttribute{
//				Optional:    true,
//				Description: "云硬盘备份存储库名称，满足以下规则：只能由数字、字母、-组成，不能以数字和-开头、且不能以-结尾，长度为2-63字符",
//			},
//			"id": schema.StringAttribute{
//				Optional:    true,
//				Description: "云硬盘备份存储库ID",
//			},
//			"status": schema.StringAttribute{
//				Optional:    true,
//				Description: "存储库状态",
//			},
//			"page_no": schema.Int32Attribute{
//				Optional:    true,
//				Description: "页码，取值范围：正整数（≥1），注：默认值为1",
//			},
//			"page_size": schema.Int32Attribute{
//				Optional:    true,
//				Description: "每页记录数目，取值范围：[1~50]，默认值：10，单页最大记录不超过50",
//			},
//			"hide_expire": schema.BoolAttribute{
//				Optional:    true,
//				Description: "是否隐藏过期的云硬盘备份存储库",
//			},
//			"query_content": schema.StringAttribute{
//				Optional:    true,
//				Description: "模糊查询，目前仅支持备份存储库名称的过滤",
//			},
//			"asc": schema.BoolAttribute{
//				Optional:    true,
//				Description: "和sort配合使用，是否是升序排列",
//			},
//			"sort": schema.StringAttribute{
//				Optional:    true,
//				Description: "和asc配合使用，指定用于排序的字段。可选字段：createdTime：创建时间。expiredTime：过期时间。size：存储库空间大小。freeSize：存储库剩余空间。usedSize：存储库已使用空间大小。repositoryName：存储库名称。",
//			},
//			"total_count": schema.Int32Attribute{
//				Computed:    true,
//				Description: "云硬盘备份存储库总数",
//			},
//			"current_count": schema.Int32Attribute{
//				Computed:    true,
//				Description: "当前页云硬盘备份存储库数",
//			},
//			"backup_repos": schema.ListNestedAttribute{
//				Computed: true,
//				NestedObject: schema.NestedAttributeObject{
//					Attributes: map[string]schema.Attribute{
//						"region_id": schema.StringAttribute{
//							Computed:    true,
//							Description: "资源池ID",
//						},
//						"id": schema.StringAttribute{
//							Computed:    true,
//							Description: "存储库ID",
//						},
//						"project_id": schema.StringAttribute{
//							Computed:    true,
//							Description: "企业项目ID",
//						},
//						"name": schema.StringAttribute{
//							Computed:    true,
//							Description: "存储库名称",
//						},
//						"status": schema.StringAttribute{
//							Computed:    true,
//							Description: "云硬盘存储库状态，expired: 已到期，active: 可用， master_order_creating: 主订单未完成，freezing: 已冻结",
//						},
//						"size": schema.Int64Attribute{
//							Computed:    true,
//							Description: "云硬盘存储库总容量，单位GB",
//						},
//						"free_size": schema.Float64Attribute{
//							Computed:    true,
//							Description: "云硬盘存储库剩余大小，单位GB(废弃该字段)",
//						},
//						"used_size": schema.Int64Attribute{
//							Computed:    true,
//							Description: "云硬盘存储库使用大小，单位Byte",
//						},
//						"created_at": schema.StringAttribute{
//							Computed:    true,
//							Description: "创建时间",
//						},
//						"expired_at": schema.StringAttribute{
//							Computed:    true,
//							Description: "到期时间",
//						},
//						"expired": schema.BoolAttribute{
//							Computed:    true,
//							Description: "存储库是否到期",
//						},
//						"freeze": schema.BoolAttribute{
//							Computed:    true,
//							Description: "是否冻结",
//						},
//						"paas": schema.BoolAttribute{
//							Computed:    true,
//							Description: "是否支持PAAS",
//						},
//						"backup_list": schema.ListAttribute{
//							ElementType: types.StringType,
//							Computed:    true,
//							Description: "存储库下可用的备份ID列表",
//						},
//					},
//				}},
//		},
//	}
//}
//
//func (c *ctyunEbsBackupRepos) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var config ctyunEbsBackupReposConfig
//	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
//	if regionId == "" {
//		err = fmt.Errorf("regionId不能为空")
//		return
//	}
//	config.RegionID = types.StringValue(regionId)
//	config.BackupRepos = []ctyunEbsBackupReposModel{}
//	// 组装请求体
//	params := &ctebsbackup.EbsbackupListBackupRepoRequest{
//		RegionID:       config.RegionID.ValueString(),
//		ProjectID:      config.ProjectID.ValueString(),
//		RepositoryName: config.Name.ValueString(),
//		RepositoryID:   config.RepositoryID.ValueString(),
//		Status:         config.Status.ValueString(),
//		PageNo:         config.PageNo.ValueInt32(),
//		PageSize:       config.PageSize.ValueInt32(),
//		Asc:            config.HideExpire.ValueBoolPointer(),
//		Sort:           config.Sort.ValueString(),
//		HideExpire:     config.HideExpire.ValueBoolPointer(),
//		QueryContent:   config.QueryContent.ValueString(),
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
//	}
//	// 解析返回值
//	for _, repo := range resp.ReturnObj.RepositoryList {
//		// 处理备份列表
//		backupList := make([]types.String, len(repo.BackupList))
//		for i, backupID := range repo.BackupList {
//			backupList[i] = types.StringValue(backupID)
//		}
//
//		item := ctyunEbsBackupReposModel{
//			RegionID:     types.StringValue(repo.RegionID),
//			RepositoryID: types.StringValue(repo.RepositoryID),
//			ProjectID:    types.StringValue(repo.ProjectID),
//			Name:         types.StringValue(repo.RepositoryName),
//			Status:       types.StringValue(repo.Status),
//			Size:         types.Int64Value(int64(repo.Size)),
//			FreeSize:     types.Float64Value(float64(repo.FreeSize)),
//			UsedSize:     types.Int64Value(int64(repo.UsedSize)),
//			CreatedAt:    types.Int32Value(repo.CreatedTime),
//			ExpiredAt:    types.Int32Value(repo.ExpiredTime),
//			Expired:      types.BoolValue(*repo.Expired),
//			Freeze:       types.BoolValue(*repo.Freeze),
//			Paas:         types.BoolValue(*repo.Paas),
//			BackupList:   backupList,
//		}
//
//		config.BackupRepos = append(config.BackupRepos, item)
//	}
//
//	// 设置分页信息
//	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
//	config.CurrentCount = types.Int32Value(resp.ReturnObj.CurrentCount)
//	// 保存到state
//	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
//}
//
//func (c *ctyunEbsBackupRepos) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
//	if request.ProviderData == nil {
//		return
//	}
//	meta := request.ProviderData.(*common.CtyunMetadata)
//	c.meta = meta
//}
