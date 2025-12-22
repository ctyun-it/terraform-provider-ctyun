package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var (
	_ resource.Resource              = &CtyunMysqlBackupCancel{}
	_ resource.ResourceWithConfigure = &CtyunMysqlBackupCancel{}
)

type CtyunMysqlBackupCancel struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMysqlBackupCancel) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_backup_cancel"
}
func NewCtyunMysqlBackupCancel() resource.Resource {
	return &CtyunMysqlBackupCancel{}
}

func (c *CtyunMysqlBackupCancel) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMysqlBackupCancel) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10098797",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "mysql实例id",
				Validators: []validator.String{
					stringvalidator.LengthBetween(32, 32),
				},
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
				Validators: []validator.String{
					validator2.Project(),
				},
			},
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
			"backup_record_id": schema.Int64Attribute{
				Required:    true,
				Description: "mysql备份记录id，可以通过data.ctyun_mysql_backups获取",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				//Validators: []validator.String{
				//	stringvalidator.UTF8LengthAtLeast(1),
				//},
			},
		},
	}
}

func (c *CtyunMysqlBackupCancel) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunBackupCancelConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 取消备份任务
	err = c.CancelBackupRecord(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlBackupCancel) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	return
}

func (c *CtyunMysqlBackupCancel) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunMysqlBackupCancel) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	return
}

func (c *CtyunMysqlBackupCancel) CancelBackupRecord(ctx context.Context, config *CtyunBackupCancelConfig) error {
	params := &mysql.TeledbCancelBackupRequest{
		BackupRecordId: config.BackupRecordId.ValueInt64(),
	}
	header := &mysql.TeledbCancelBackupRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbCancelBackupApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("终止备份任务失败，mysql实例id=%s，record_id=%d", config.InstID.ValueString(), config.BackupRecordId.ValueInt64())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("cancel backup error, API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}

	err = c.cancelLoop(ctx, config, 60)
	if err != nil {
		return err
	}
	return nil
}

func (c *CtyunMysqlBackupCancel) cancelLoop(ctx context.Context, config *CtyunBackupCancelConfig, loopCount ...int) error {
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
			resp, err2 := c.getBackupRecordDetail(ctx, config)
			if err2 != nil {
				err = err2
				return false
			}
			detail := resp.ReturnObj
			switch detail.TaskStatus {
			case business.MysqlBackupTaskStatusCancel:
				return false
			case business.MysqlBackupTaskStatusWaitStart:
				return true
			case business.MysqlBackupTaskStatusSubmit:
				return true
			case business.MysqlBackupTaskStatusSuccess:
				err = fmt.Errorf("备份任务状态为成功，取消备份失败，mysql实例id=%s，record_id=%d", config.InstID.ValueString(), config.BackupRecordId.ValueInt64())
				return false
			default:
				err = fmt.Errorf("备份任务状态为失败，取消备份失败，mysql实例id=%s，record_id=%d", config.InstID.ValueString(), config.BackupRecordId.ValueInt64())
				return false
			}
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，备份任务仍未取消成功！")
	}
	return err
}

func (c *CtyunMysqlBackupCancel) getBackupRecordDetail(ctx context.Context, config *CtyunBackupCancelConfig) (*mysql.TeledbGetBackupRecordDetailResponse, error) {
	params := &mysql.TeledbGetBackupRecordDetailRequest{
		Id: fmt.Sprintf("%d", config.BackupRecordId.ValueInt64()),
	}
	header := &mysql.TeledbGetBackupRecordDetailRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetBackupRecordDetailApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询备份任务(record_id=%d)信息失败，接口返回nil。具体原因请联系研发确定。", config.BackupRecordId.ValueInt64())
		return nil, err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	}
	if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}

	return resp, nil
}

type CtyunBackupCancelConfig struct {
	InstID         types.String `tfsdk:"instance_id"`
	ProjectID      types.String `tfsdk:"project_id"`
	RegionID       types.String `tfsdk:"region_id"`
	BackupRecordId types.Int64  `tfsdk:"backup_record_id"`
}
