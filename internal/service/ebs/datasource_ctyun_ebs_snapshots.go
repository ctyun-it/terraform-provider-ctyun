package ebs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctebs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebs"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunEbsSnapshots{}
	_ datasource.DataSourceWithConfigure = &ctyunEbsSnapshots{}
)

type ctyunEbsSnapshots struct {
	meta *common.CtyunMetadata
}

func NewCtyunEbsSnapshots() datasource.DataSource {
	return &ctyunEbsSnapshots{}
}

func (c *ctyunEbsSnapshots) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs_snapshots"
}

type ctyunEbsSnapshotsModel struct {
	SnapshotID       types.String `tfsdk:"id"`
	SnapshotName     types.String `tfsdk:"name"`
	DiskID           types.String `tfsdk:"disk_id"`
	VolumeName       types.String `tfsdk:"volume_name"`
	AvailabilityZone types.String `tfsdk:"availability_zone"`
	SnapshotStatus   types.String `tfsdk:"snapshot_status"`
	CreateTime       types.String `tfsdk:"create_time"`
	AzName           types.String `tfsdk:"az_name"`
	DeleteTime       types.String `tfsdk:"delete_time"`
	Description      types.String `tfsdk:"description"`
	ExpireTime       types.String `tfsdk:"expire_time"`
	Freezing         types.Bool   `tfsdk:"freezing"`
	IsEncryted       types.Bool   `tfsdk:"is_encryted"`
	IsMaz            types.Bool   `tfsdk:"is_maz"`
	RegionID         types.String `tfsdk:"region_id"`
	IsTalkOrder      types.Bool   `tfsdk:"is_talk_order"`
	RetentionPolicy  types.String `tfsdk:"retention_policy"`
	RetentionTime    types.Int64  `tfsdk:"retention_time"`
	SnapshotType     types.String `tfsdk:"snapshot_type"`
	UpdateTime       types.String `tfsdk:"update_time"`
	VolumeAttr       types.String `tfsdk:"volume_attr"`
	VolumeSize       types.Int64  `tfsdk:"volume_size"`
	VolumeSource     types.String `tfsdk:"volume_source"`
	VolumeStatus     types.String `tfsdk:"volume_status"`
	VolumeType       types.String `tfsdk:"volume_type"`
}

type ctyunEbsSnapshotsConfig struct {
	RegionID        types.String             `tfsdk:"region_id"`
	DiskID          types.String             `tfsdk:"disk_id"`
	SnapshotID      types.String             `tfsdk:"id"`
	SnapshotName    types.String             `tfsdk:"name"`
	SnapshotStatus  types.String             `tfsdk:"snapshot_status"`
	SnapshotType    types.String             `tfsdk:"snapshot_type"`
	VolumeAttr      types.String             `tfsdk:"volume_attr"`
	RetentionPolicy types.String             `tfsdk:"retention_policy"`
	PageNo          types.Int64              `tfsdk:"page_no"`
	PageSize        types.Int64              `tfsdk:"page_size"`
	Snapshots       []ctyunEbsSnapshotsModel `tfsdk:"snapshots"`
}

func (c *ctyunEbsSnapshots) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026730/10335345**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘快照ID",
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘快照名称",
			},
			"disk_id": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘ID",
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
			"snapshot_status": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘快照状态。取值为：available：可用，freezing：冻结，creating：创建中，deleting：删除中，rollbacking：回滚中，cloning：从快照创建云硬盘中，error：错误",
			},
			"snapshot_type": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘快照创建类型。取值为：manu：手动，timer：自动",
			},
			"volume_attr": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘属性。取值为：data：数据盘，system：系统盘",
			},
			"retention_policy": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘快照保留策略。取值为：forever：永久保留，custom：自定义保留天数",
			},
			"page_no": schema.Int64Attribute{
				Optional:    true,
				Description: "页码，取值范围：正整数（≥1），注：默认值为1",
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "每页记录数目，取值范围：[1,50]，注：默认值为10",
			},
			"snapshots": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘快照ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘快照名称",
						},
						"disk_id": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘ID",
						},
						"volume_name": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘名称",
						},
						"availability_zone": schema.StringAttribute{
							Computed:    true,
							Description: "可用区名称",
						},
						"snapshot_status": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘快照状态，取值为：creating/deleting/rollbacking/cloning/available/error，分别对应创建中/删除中/回滚中/从快照创建云硬盘中/可用/错误",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "可用区名称",
						},
						"delete_time": schema.StringAttribute{
							Computed:    true,
							Description: "删除时间",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述信息",
						},
						"expire_time": schema.StringAttribute{
							Computed:    true,
							Description: "过期时间",
						},
						"freezing": schema.BoolAttribute{
							Computed:    true,
							Description: "是否被冻结",
						},
						"is_encryted": schema.BoolAttribute{
							Computed:    true,
							Description: "是否是加密盘",
						},
						"is_maz": schema.BoolAttribute{
							Computed:    true,
							Description: "是否跨AZ",
						},
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源池ID",
						},
						"is_talk_order": schema.BoolAttribute{
							Computed:    true,
							Description: "是否是按需计费资源",
						},
						"retention_policy": schema.StringAttribute{
							Computed:    true,
							Description: "快照保留策略，取值为：custom：自定义保留天数，forever：永久保留",
						},
						"retention_time": schema.Int64Attribute{
							Computed:    true,
							Description: "快照保留时间",
						},
						"snapshot_type": schema.StringAttribute{
							Computed:    true,
							Description: "快照类型，取值为：manu：手动，timer：自动",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间",
						},
						"volume_attr": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘属性，取值为：data：数据盘，system：系统盘",
						},
						"volume_size": schema.Int64Attribute{
							Computed:    true,
							Description: "云硬盘大小",
						},
						"volume_source": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘来源，如果为空，则是普通云硬盘，如果不为空，则是由快照创建而来，显示来源快照ID",
						},
						"volume_status": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘的状态，请参考云硬盘使用状态",
						},
						"volume_type": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘类型，取值为：SATA：普通IO，SAS：高IO，SSD：超高IO，FAST-SSD：极速型SSD，XSSD-0、XSSD-1、XSSD-2：X系列云硬盘",
						},
					},
				}},
		},
	}
}

func (c *ctyunEbsSnapshots) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config ctyunEbsSnapshotsConfig
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
	config.Snapshots = []ctyunEbsSnapshotsModel{}
	// 组装请求体
	params := &ctebs2.EbsListEbsSnapRequest{
		RegionID:        config.RegionID.ValueString(),
		SnapshotName:    stringPtr(config.SnapshotName.ValueString()),
		SnapshotID:      stringPtr(config.SnapshotID.ValueString()),
		DiskID:          stringPtr(config.DiskID.ValueString()),
		SnapshotStatus:  stringPtr(config.SnapshotStatus.ValueString()),
		SnapshotType:    stringPtr(config.SnapshotType.ValueString()),
		VolumeAttr:      stringPtr(config.VolumeAttr.ValueString()),
		RetentionPolicy: stringPtr(config.RetentionPolicy.ValueString()),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsListEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
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
	for _, snapshot := range resp.ReturnObj.SnapshotList {
		item := ctyunEbsSnapshotsModel{
			SnapshotID:       types.StringValue(stringValue(snapshot.SnapshotID)),
			SnapshotName:     types.StringValue(stringValue(snapshot.SnapshotName)),
			DiskID:           types.StringValue(stringValue(snapshot.DiskID)),
			VolumeName:       types.StringValue(stringValue(snapshot.VolumeName)),
			AvailabilityZone: types.StringValue(stringValue(snapshot.AvailabilityZone)),
			SnapshotStatus:   types.StringValue(stringValue(snapshot.SnapshotStatus)),
			CreateTime:       types.StringValue(stringValue(snapshot.CreateTime)),
			AzName:           types.StringValue(stringValue(snapshot.AzName)),
			DeleteTime:       types.StringValue(stringValue(snapshot.DeleteTime)),
			Description:      types.StringValue(stringValue(snapshot.Description)),
			ExpireTime:       types.StringValue(stringValue(snapshot.ExpireTime)),
			Freezing:         types.BoolValue(boolValue(snapshot.Freezing)),
			IsEncryted:       types.BoolValue(boolValue(snapshot.IsEncryted)),
			IsMaz:            types.BoolValue(boolValue(snapshot.IsMaz)),
			RegionID:         types.StringValue(stringValue(snapshot.RegionID)),
			IsTalkOrder:      types.BoolValue(boolValue(snapshot.IsTalkOrder)),
			RetentionPolicy:  types.StringValue(stringValue(snapshot.RetentionPolicy)),
			RetentionTime:    types.Int64Value(snapshot.RetentionTime),
			SnapshotType:     types.StringValue(stringValue(snapshot.SnapshotType)),
			UpdateTime:       types.StringValue(stringValue(snapshot.UpdateTime)),
			VolumeAttr:       types.StringValue(stringValue(snapshot.VolumeAttr)),
			VolumeSize:       types.Int64Value(snapshot.VolumeSize),
			VolumeSource:     types.StringValue(stringValue(snapshot.VolumeSource)),
			VolumeStatus:     types.StringValue(stringValue(snapshot.VolumeStatus)),
			VolumeType:       types.StringValue(stringValue(snapshot.VolumeType)),
		}

		config.Snapshots = append(config.Snapshots, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

// stringPtr 返回指向字符串值的指针
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// boolValue 返回布尔指针指向的值，如果指针为nil则返回false
func boolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func (c *ctyunEbsSnapshots) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
