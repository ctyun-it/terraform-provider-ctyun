package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunMysqlRecoverableTimePoints{}
	_ datasource.DataSourceWithConfigure = &ctyunMysqlRecoverableTimePoints{}
)

type ctyunMysqlRecoverableTimePoints struct {
	meta *common.CtyunMetadata
}

func NewCtyunMysqlRecoverableTimePoints() datasource.DataSource {
	return &ctyunMysqlRecoverableTimePoints{}
}
func (c *ctyunMysqlRecoverableTimePoints) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunMysqlRecoverableTimePoints) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_recoverable_time_points"
}

func (c *ctyunMysqlRecoverableTimePoints) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10098797",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "mysql 实例id",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "项目ID",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "区域ID",
			},
			"backup_time_points": schema.ListNestedAttribute{
				Computed:    true,
				Description: "可恢复时间点",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start_time": schema.StringAttribute{
							Computed:    true,
							Description: "开始时间",
						},
						"end_time": schema.StringAttribute{
							Computed:    true,
							Description: "结束时间",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunMysqlRecoverableTimePoints) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunMysqlRecoverableTimePointsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	config.RegionID = types.StringValue(regionId)
	params := &mysql.TeledbGetRecoverableTimeRangesRequest{
		OuterProdInstId: config.InstID.ValueString(),
	}
	header := &mysql.TeledbGetRecoverableTimeRangesRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetRecoverableTimeRangesApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("获取备份可恢复时间点失败，接口返回nil。请联系研发确认问题原因。")
		return
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s ", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var backupTimePoints []BackupTimePointModel
	timeRanges := resp.ReturnObj
	for _, timeRange := range timeRanges.Data {
		var backupTimePoint BackupTimePointModel
		backupTimePoint.StartTimestamp = types.StringValue(utils.FromBJTimeToUTCZ(timeRange.StartTimestamp))
		backupTimePoint.EndTimestamp = types.StringValue(utils.FromBJTimeToUTCZ(timeRange.EndTimestamp))
		backupTimePoints = append(backupTimePoints, backupTimePoint)
	}
	config.BackupTimePoints = backupTimePoints
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type BackupTimePointModel struct {
	StartTimestamp types.String `tfsdk:"start_time"`
	EndTimestamp   types.String `tfsdk:"end_time"`
}

type CtyunMysqlRecoverableTimePointsConfig struct {
	InstID           types.String           `tfsdk:"instance_id"`
	ProjectID        types.String           `tfsdk:"project_id"`
	RegionID         types.String           `tfsdk:"region_id"`
	BackupTimePoints []BackupTimePointModel `tfsdk:"backup_time_points"`
}
