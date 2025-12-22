package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &CtyunPostgresqlBackup{}
	_ resource.ResourceWithConfigure   = &CtyunPostgresqlBackup{}
	_ resource.ResourceWithImportState = &CtyunPostgresqlBackup{}
)

type CtyunPostgresqlBackup struct {
	meta *common.CtyunMetadata
}

func (c *CtyunPostgresqlBackup) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_backup"
}
func NewCtyunPostgresqlBackup() resource.Resource {
	return &CtyunPostgresqlBackup{}
}

func (c *CtyunPostgresqlBackup) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPostgresqlBackup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [name],[instanceID],[projectID],[regionID]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunPostgresqlBackupConfig
	var name, regionId, projectId, instId string

	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		projectId = c.meta.GetExtraIfEmpty(projectId, common.ExtraProjectId)
		err = terraform_extend.Split(request.ID, &name, &instId)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &name, &instId, &projectId, &regionId)
		if err != nil {
			return
		}
	}
	if name == "" {
		err = fmt.Errorf("name 不能为空")
		return
	}
	if instId == "" {
		err = fmt.Errorf("instID 不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID 不能为空")
		return
	}
	cfg.RegionID = types.StringValue(regionId)
	cfg.ProjectID = types.StringValue(projectId)
	cfg.Name = types.StringValue(name)
	cfg.InstID = types.StringValue(instId)
	err = c.getAndMergePostgresqlBackup(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)

}

func (c *CtyunPostgresqlBackup) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10034019/10160072",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "MySQL实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "备份集名称，不可重复",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "备份描述",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
					validator2.Desc(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "备份集ID",
			},
			"backup_type": schema.StringAttribute{
				Computed:    true,
				Description: "备份集类型，auto：自动备份，manual：手动备份，recovery：恢复备份",
			},
			"backup_result": schema.StringAttribute{
				Computed:    true,
				Description: "备份结果，ing：运行中，success：备份成功，fail：备份失败",
			},
			"start_time": schema.StringAttribute{
				Computed:    true,
				Description: "备份开始时间，时间格式为utc",
			},
			"end_time": schema.StringAttribute{
				Computed:    true,
				Description: "备份结束时间，时间格式为utc",
			},
		},
	}
}

func (c *CtyunPostgresqlBackup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunPostgresqlBackupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 开始创建备份集
	err = c.CreatePostgresqlBackup(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取mysql详情
	err = c.getAndMergePostgresqlBackup(ctx, &plan)
	if err != nil {
		return
	}
	//plan.ID = types.StringValue(plan.BackupName.ValueString())
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPostgresqlBackup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPostgresqlBackupConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergePostgresqlBackup(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPostgresqlBackup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunPostgresqlBackup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunPostgresqlBackupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.backupIngLoop(ctx, &config)
	if err != nil {
		return
	}

	params := &pgsql.PgsqlDeleteBackupRequest{
		ProdInstId: config.InstID.ValueString(),
		BackupId:   config.ID.ValueInt64(),
	}
	header := &pgsql.PgsqlDeleteBackupRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlDeleteBackupApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("删除postgresql实例(id=%s)的备份集(name=%s, id=%d)失败，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString(), config.Name.ValueString(), config.ID.ValueInt64())
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
}

func (c *CtyunPostgresqlBackup) CreatePostgresqlBackup(ctx context.Context, config *CtyunPostgresqlBackupConfig) error {
	params := &pgsql.PgsqlCreateBackupRequest{
		ProdInstId: config.InstID.ValueString(),
		BackupName: config.Name.ValueString(),
	}
	if !config.Description.IsNull() {
		params.Desc = config.Description.ValueStringPointer()
	}
	header := &pgsql.PgsqlCreateBackupRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlCreateBackupApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例(id=%s)创建备份集失败，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}

	// 创建后，获取id
	listResp, err := c.getBackupDetail(ctx, config)
	if err != nil {
		return err
	}

	config.ID = types.Int64Value(listResp.ReturnObj.List[0].Id)

	return nil
}

func (c *CtyunPostgresqlBackup) getBackupDetail(ctx context.Context, config *CtyunPostgresqlBackupConfig) (*pgsql.PgsqlGetBackupListResponse, error) {
	params := &pgsql.PgsqlGetBackupListRequest{
		ProdInstId: config.InstID.ValueString(),
		PageNum:    1,
		PageSize:   10,
		BackupName: config.Name.ValueStringPointer(),
	}
	header := &pgsql.PgsqlGetBackupListRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGetBackupListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询postgresql实例(id=%s)的备份集(name=%s)失败，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString(), config.Name.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	if len(resp.ReturnObj.List) > 1 {
		err = fmt.Errorf("postgresql实例(id=%s)中存在重名备份集(name=%s)", config.InstID.ValueString(), config.Name.ValueString())
		return nil, err
	} else if len(resp.ReturnObj.List) == 0 {
		err = fmt.Errorf("postgresql实例(id=%s)中不存在名为%s的备份集", config.InstID.ValueString(), config.Name.ValueString())
		return nil, err
	}
	return resp, nil
}

func (c *CtyunPostgresqlBackup) getAndMergePostgresqlBackup(ctx context.Context, config *CtyunPostgresqlBackupConfig) error {
	detailList, err := c.getBackupDetail(ctx, config)
	if err != nil {
		return err
	}
	detail := detailList.ReturnObj.List[0]
	config.ID = types.Int64Value(detail.Id)
	config.BackupType = types.StringValue(business.PgsqlBackupTypeMapConv[detail.Type])
	config.BackupResult = types.StringValue(business.PgsqlBackupResultMapConv[detail.Result])
	config.StartTime = types.StringValue(utils.FromBJTimeToUTCZ(detail.StartTime))
	config.EndTime = types.StringValue(utils.FromBJTimeToUTCZ(detail.EndTime))
	return nil

}

func (c *CtyunPostgresqlBackup) backupIngLoop(ctx context.Context, config *CtyunPostgresqlBackupConfig, loopCount ...int) error {
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
			resp, err2 := c.getBackupDetail(ctx, config)
			if err2 != nil {
				err = err2
				return false
			}
			backupList := resp.ReturnObj.List

			backupRecordList := backupList[0]
			taskStatus := backupRecordList.Result
			switch taskStatus {
			case business.PgsqlBackupResultING:
				return true
			case business.PgsqlBackupResultING1:
				return true
			case business.PgsqlBackupResultSuccess:
				return false
			case business.PgsqlBackupResultFail:
				return false
			default:
				return false
			}
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，备份任务仍未备份成功！")
	}
	return err
}

type CtyunPostgresqlBackupConfig struct {
	RegionID     types.String `tfsdk:"region_id"`
	ProjectID    types.String `tfsdk:"project_id"`
	InstID       types.String `tfsdk:"instance_id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	ID           types.Int64  `tfsdk:"id"`
	BackupType   types.String `tfsdk:"backup_type"`
	BackupResult types.String `tfsdk:"backup_result"`
	StartTime    types.String `tfsdk:"start_time"`
	EndTime      types.String `tfsdk:"end_time"`
}
