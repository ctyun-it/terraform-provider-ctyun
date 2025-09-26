package ecs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctecs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunEcsSnapshots{}
	_ datasource.DataSourceWithConfigure = &ctyunEcsSnapshots{}
)

type ctyunEcsSnapshots struct {
	meta *common.CtyunMetadata
}

func NewCtyunEcsSnapshots() datasource.DataSource {
	return &ctyunEcsSnapshots{}
}

func (c *ctyunEcsSnapshots) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ecs_snapshots"
}

type ctyunEcsSnapshotsModel struct {
	SnapshotID   types.String `tfsdk:"id"`
	SnapshotName types.String `tfsdk:"name"`
	InstanceID   types.String `tfsdk:"instance_id"`
	InstanceName types.String `tfsdk:"instance_name"`
}

type ctyunEcsSnapshotsConfig struct {
	RegionID       types.String             `tfsdk:"region_id"`
	InstanceID     types.String             `tfsdk:"instance_id"`
	SnapshotID     types.String             `tfsdk:"id"`
	SnapshotName   types.String             `tfsdk:"name"`
	SnapshotStatus types.String             `tfsdk:"snapshot_status"`
	PageNo         types.Int32              `tfsdk:"page_no"`
	PageSize       types.Int32              `tfsdk:"page_size"`
	Snapshots      []ctyunEcsSnapshotsModel `tfsdk:"snapshots"`
}

func (c *ctyunEcsSnapshots) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026730/10335345**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "云主机快照ID",
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "云主机快照名称",
			},
			"snapshot_status": schema.StringAttribute{
				Optional:    true,
				Description: "快照状态",
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "云主机ID",
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "页码，取值范围：正整数（≥1），注：默认值为1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页记录数目，取值范围：[1,50]，注：默认值为10",
			},
			"snapshots": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "云主机快照ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "云主机快照名称",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "云主机ID",
						},
						"instance_name": schema.StringAttribute{
							Computed:    true,
							Description: "云主机名称",
						},
					},
				}},
		},
	}
}

func (c *ctyunEcsSnapshots) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config ctyunEcsSnapshotsConfig
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
	config.Snapshots = []ctyunEcsSnapshotsModel{}
	// 组装请求体
	params := &ctecs2.CtecsQuerySnapshotListV41Request{
		RegionID:       config.RegionID.ValueString(),
		SnapshotName:   config.SnapshotName.ValueString(),
		SnapshotID:     config.SnapshotID.ValueString(),
		SnapshotStatus: config.SnapshotStatus.ValueString(),
		InstanceID:     config.InstanceID.ValueString(),
		PageNo:         config.PageNo.ValueInt32(),
		PageSize:       config.PageSize.ValueInt32(),
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsQuerySnapshotListV41Api.Do(ctx, c.meta.SdkCredential, params)
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
	for _, snapshot := range resp.ReturnObj.Results {
		item := ctyunEcsSnapshotsModel{
			SnapshotID:   types.StringValue(snapshot.SnapshotID),
			SnapshotName: types.StringValue(snapshot.SnapshotName),
			InstanceID:   types.StringValue(snapshot.InstanceID),
			InstanceName: types.StringValue(snapshot.InstanceName),
		}

		config.Snapshots = append(config.Snapshots, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunEcsSnapshots) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
