package redis

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunRedisBackups{}
	_ datasource.DataSourceWithConfigure = &ctyunRedisBackups{}
)

type ctyunRedisBackups struct {
	meta *common.CtyunMetadata
}

func NewCtyunRedisBackups() datasource.DataSource {
	return &ctyunRedisBackups{}
}

func (c *ctyunRedisBackups) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_backups"
}

type CtyunRedisBackupModel struct {
	Name           types.String `tfsdk:"name"`
	CreateTime     types.String `tfsdk:"create_time"`
	Status         types.String `tfsdk:"status"`
	RecoveryStatus types.String `tfsdk:"recovery_status"`
	Type           types.Int32  `tfsdk:"type"`
	Remark         types.String `tfsdk:"remark"`
}

type CtyunRedisBackupsConfig struct {
	RegionID   types.String            `tfsdk:"region_id"`
	InstanceId types.String            `tfsdk:"instance_id"`
	Name       types.String            `tfsdk:"name"`
	Rows       []CtyunRedisBackupModel `tfsdk:"rows"`
}

func (c *ctyunRedisBackups) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029420/10142282`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "实例ID",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "备份名，用于精确查询特定备份",
			},
			"rows": schema.ListNestedAttribute{
				Computed:    true,
				Description: "备份列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "备份名",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间（格式：yyyy-MM-dd HH:mm:ss）",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "节点状态: success(成功), processing(进行中), fail(失败)",
						},
						"recovery_status": schema.StringAttribute{
							Computed:    true,
							Description: "备份恢复状态: success(成功), processing(进行中), fail(失败), create(备份点创建)",
						},
						"type": schema.Int32Attribute{
							Computed:    true,
							Description: "备份类型: 0(手动备份), 1(自动备份)",
						},
						"remark": schema.StringAttribute{
							Computed:    true,
							Description: "备注信息",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRedisBackups) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRedisBackupsConfig
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
	instanceId := config.InstanceId.ValueString()
	if instanceId == "" {
		err = fmt.Errorf("instanceId不能为空")
		return
	}

	// 组装请求体
	params := &dcs2.Dcs2DescribeBackupsRequest{
		RegionId:    regionId,
		ProdInstId:  instanceId,
		RestoreName: config.Name.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeBackupsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回值
	config.Rows = []CtyunRedisBackupModel{}
	for _, backup := range resp.ReturnObj.Rows {
		item := CtyunRedisBackupModel{
			Name:           types.StringValue(backup.RestoreName),
			CreateTime:     types.StringValue(backup.CreateTime),
			Status:         types.StringValue(backup.Status),
			RecoveryStatus: types.StringValue(backup.RecoveryStatus),
			Type:           types.Int32Value(backup.RawType),
			Remark:         types.StringValue(backup.Remark),
		}
		config.Rows = append(config.Rows, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRedisBackups) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
