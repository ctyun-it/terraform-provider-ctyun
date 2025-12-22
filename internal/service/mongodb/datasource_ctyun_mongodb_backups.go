package mongodb

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &CtyunMongodbBackups{}
	_ datasource.DataSourceWithConfigure = &CtyunMongodbBackups{}
)

func NewCtyunMongodbBackups() datasource.DataSource {
	return &CtyunMongodbBackups{}
}

// CtyunMongodbBackups defines the data source implementation.
type CtyunMongodbBackups struct {
	meta *common.CtyunMetadata
}

// CtyunMongodbBackupsConfig describes the data source data model.
type CtyunMongodbBackupsConfig struct {
	ID         types.String              `tfsdk:"id"`
	RegionID   types.String              `tfsdk:"region_id"`
	ProjectID  types.String              `tfsdk:"project_id"`
	InstanceID types.String              `tfsdk:"instance_id"`
	BackupType types.String              `tfsdk:"type"`
	PageNo     types.Int32               `tfsdk:"page_no"`
	PageSize   types.Int32               `tfsdk:"page_size"`
	Backups    []CtyunMongodbBackupModel `tfsdk:"backups"`
}

type CtyunMongodbBackupModel struct {
	BackupID          types.String `tfsdk:"id"`
	BackupName        types.String `tfsdk:"name"`
	BackupMethod      types.String `tfsdk:"method"`
	BackupType        types.String `tfsdk:"type"`
	BackupStatus      types.String `tfsdk:"status"`
	BackupStartTime   types.String `tfsdk:"start_time"`
	BackupEndTime     types.String `tfsdk:"end_time"`
	BackupSize        types.Int64  `tfsdk:"size"`
	Description       types.String `tfsdk:"description"`
	BackupTriggerType types.String `tfsdk:"trigger_type"`
}

func (d *CtyunMongodbBackups) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_backups"
}

func (d *CtyunMongodbBackups) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10034467/10089535`,

		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "项目id",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "实例ID",
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "备份ID",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "备份类型",
				Validators: []validator.String{
					stringvalidator.OneOf("full", "incremental"),
				},
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "页码",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "每页条数",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
					int32validator.AtMost(100),
				},
			},
			"backups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "备份列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "备份ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "备份名称",
						},
						"method": schema.StringAttribute{
							Computed:    true,
							Description: "备份方式",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "备份类型",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "备份状态",
						},
						"start_time": schema.StringAttribute{
							Computed:    true,
							Description: "备份开始时间",
						},
						"end_time": schema.StringAttribute{
							Computed:    true,
							Description: "备份结束时间",
						},
						"size": schema.Int64Attribute{
							Computed:    true,
							Description: "备份大小",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "备份描述",
						},
						"trigger_type": schema.StringAttribute{
							Computed:    true,
							Description: "备份触发类型",
						},
					},
				},
			},
		},
	}
}

func (d *CtyunMongodbBackups) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	d.meta = meta
}

func (d *CtyunMongodbBackups) Read(ctx context.Context, req datasource.ReadRequest, response *datasource.ReadResponse) {
	var data CtyunMongodbBackupsConfig

	// Read Terraform configuration data into the model
	response.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 设置默认分页参数
	pageNo := int32(1)
	if !data.PageNo.IsNull() {
		pageNo = data.PageNo.ValueInt32()
	}

	pageSize := int32(1000)
	if !data.PageSize.IsNull() {
		pageSize = data.PageSize.ValueInt32()
	}

	// 查询备份列表
	describeReq := &mongodb.MongodbDescribeBackupsRequest{
		ProdInstId: data.InstanceID.ValueString(),
		PageNow:    pageNo,
		PageSize:   pageSize,
	}

	//if !data.BackupID.IsNull() {
	//	describeReq.BackupId = data.BackupID.ValueStringPointer()
	//}

	if !data.BackupType.IsNull() {
		describeReq.BackupType = data.BackupType.ValueStringPointer()
	}
	regionId := d.meta.GetExtraIfEmpty(data.RegionID.ValueString(), common.ExtraRegionId)

	header := &mongodb.MongodbDescribeBackupsRequestHeaders{
		RegionID: regionId,
	}

	if !data.ProjectID.IsNull() {
		header.ProjectID = data.ProjectID.ValueStringPointer()
	}

	resp, err := d.meta.Apis.SdkMongodbApis.MongodbDescribeBackupsApi.Do(ctx, d.meta.Credential, describeReq, header)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s ", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 转换备份信息
	var backups []CtyunMongodbBackupModel
	for _, item := range resp.ReturnObj.List {
		backup := CtyunMongodbBackupModel{
			BackupID:          types.StringValue(fmt.Sprintf("%d", item.BackupId)),
			BackupName:        types.StringValue(item.BackupName),
			BackupMethod:      types.StringValue(item.BackupMethod),
			BackupType:        types.StringValue(item.BackupType),
			BackupStatus:      types.StringValue(item.BackupStatus),
			BackupStartTime:   types.StringValue(item.BackupStartTime),
			BackupEndTime:     types.StringValue(item.BackupEndTime),
			BackupSize:        types.Int64Value(item.BackupSize),
			BackupTriggerType: types.StringValue(item.BackupTriggerType),
		}

		if item.Description != nil {
			backup.Description = types.StringValue(*item.Description)
		}

		backups = append(backups, backup)
	}

	data.Backups = backups
	data.ID = types.StringValue(fmt.Sprintf("%s:%s", data.RegionID.ValueString(), data.InstanceID.ValueString()))

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
