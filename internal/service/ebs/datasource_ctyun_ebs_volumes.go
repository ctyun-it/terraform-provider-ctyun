package ebs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctebs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunEbsVolumes{}
	_ datasource.DataSourceWithConfigure = &ctyunEbsVolumes{}
)

type ctyunEbsVolumes struct {
	meta *common.CtyunMetadata
}

func NewCtyunEbsVolumes() datasource.DataSource {
	return &ctyunEbsVolumes{}
}

func (c *ctyunEbsVolumes) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs_volumes"
}

type CtyunEbsVolumesModel struct {
	ID              types.String                 `tfsdk:"id"`
	Name            types.String                 `tfsdk:"name"`
	Mode            types.String                 `tfsdk:"mode"`
	Type            types.String                 `tfsdk:"type"`
	Size            types.Int32                  `tfsdk:"size"`
	Status          types.String                 `tfsdk:"status"`
	IsEncrypt       types.Bool                   `tfsdk:"is_encrypt"`
	CreateTime      types.Int64                  `tfsdk:"create_time"`
	UpdateTime      types.Int64                  `tfsdk:"update_time"`
	ExpireTime      types.Int64                  `tfsdk:"expire_time"`
	IsSystemVolume  types.Bool                   `tfsdk:"is_system_volume"`
	IsPackaged      types.Bool                   `tfsdk:"is_packaged"`
	InstanceName    types.String                 `tfsdk:"instance_name"`
	InstanceID      types.String                 `tfsdk:"instance_id"`
	InstanceStatus  types.String                 `tfsdk:"instance_status"`
	MultiAttach     types.Bool                   `tfsdk:"multi_attach"`
	ProjectID       types.String                 `tfsdk:"project_id"`
	KmsUUID         types.String                 `tfsdk:"kms_uuid"`
	RegionID        types.String                 `tfsdk:"region_id"`
	AzName          types.String                 `tfsdk:"az_name"`
	DiskFreeze      types.Bool                   `tfsdk:"disk_freeze"`
	ProvisionedIops types.Int32                  `tfsdk:"provisioned_iops"`
	Attachments     []CtyunEbsVolumesAttachments `tfsdk:"attachments"`
}

type CtyunEbsVolumesAttachments struct {
	InstanceID   types.String `tfsdk:"instance_id"`
	AttachmentID types.String `tfsdk:"attachment_id"`
	Device       types.String `tfsdk:"device"`
}

type CtyunEbsVolumesConfig struct {
	RegionID  types.String           `tfsdk:"region_id"`
	AzName    types.String           `tfsdk:"az_name"`
	ProjectID types.String           `tfsdk:"project_id"`
	PageNo    types.Int32            `tfsdk:"page_no"`
	DiskID    types.String           `tfsdk:"disk_id"`
	DiskName  types.String           `tfsdk:"disk_name"`
	PageSize  types.Int32            `tfsdk:"page_size"`
	Volumes   []CtyunEbsVolumesModel `tfsdk:"volumes"`
}

func (c *ctyunEbsVolumes) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027696/10027930**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区id，如果不填则默认使用provider ctyun中的az_name或环境变量中的CTYUN_AZ_NAME",
			},
			"disk_id": schema.StringAttribute{
				Optional:    true,
				Description: "磁盘id",
			},
			"disk_name": schema.StringAttribute{
				Optional:    true,
				Description: "磁盘名称",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "页码，取值范围：正整数（≥1），注：默认值为1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "每页记录数目，取值范围：[1,300]，注：默认值为10",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID",
			},
			"volumes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "磁盘名",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "磁盘ID",
						},
						"size": schema.Int64Attribute{
							Computed:    true,
							Description: "磁盘大小（GB）",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "磁盘规格类型SATA/SAS/SSD-genric/SSD/FAST-SSD/XSSD-XSSD-1/XSSD-2",
						},
						"mode": schema.StringAttribute{
							Computed:    true,
							Description: "磁盘模式。VBD/ISCSI/FCSAN",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "参考云硬盘使用状态",
						},
						"create_time": schema.Int64Attribute{
							Computed:    true,
							Description: "创建时刻，epoch时戳，精度毫秒",
						},
						"update_time": schema.Int64Attribute{
							Computed:    true,
							Description: "更新时刻，epoch时戳，精度毫秒",
						},
						"expire_time": schema.Int64Attribute{
							Computed:    true,
							Description: "过期时刻，epoch时戳，精度毫秒",
						},
						"is_system_volume": schema.BoolAttribute{
							Computed:    true,
							Description: "只有为系统盘时才返回该字段",
						},
						"is_packaged": schema.BoolAttribute{
							Computed:    true,
							Description: "是否是云主机成套资源",
						},
						"instance_name": schema.StringAttribute{
							Computed:    true,
							Description: "绑定的云主机名，有挂载时才返回",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "绑定云主机resourceUUID，有挂载时才返回",
						},
						"instance_status": schema.StringAttribute{
							Computed:    true,
							Description: "云主机状态参考云主机状态，有挂载时才返回",
						},
						"multi_attach": schema.BoolAttribute{
							Computed:    true,
							Description: "是否共享云硬盘",
						},
						"attachments": schema.ListNestedAttribute{
							Computed:    true,
							Description: "挂载信息。如果是共享挂载云硬盘，有多项，无挂载时不返回该字段",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"instance_id": schema.StringAttribute{
										Computed:    true,
										Description: "绑定主机实例id",
									},
									"attachment_id": schema.StringAttribute{
										Computed:    true,
										Description: "挂载id",
									},
									"device": schema.StringAttribute{
										Computed:    true,
										Description: "挂载设备名",
									},
								},
							},
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源所属企业项目ID",
						},
						"is_encrypt": schema.BoolAttribute{
							Computed:    true,
							Description: "是否加密盘",
						},
						"kms_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "加密盘密钥UUID，是加密盘时才返回",
						},
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源池ID",
						},
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "多可用区下的可用区名字，非多可用区不返回该字段",
						},
						"disk_freeze": schema.BoolAttribute{
							Computed:    true,
							Description: "是否冻结",
						},
						"provisioned_iops": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "XSSD类型盘的预配置iops，未配置返回0，其他类型盘不返回",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunEbsVolumes) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunEbsVolumesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	projectID := c.meta.GetExtraIfEmpty(config.ProjectID.ValueString(), common.ExtraProjectId)
	config.RegionID = types.StringValue(regionId)
	config.ProjectID = types.StringValue(projectID)

	if config.DiskID.ValueString() != "" {
		err = c.getByID(ctx, &config)
	} else if config.DiskName.ValueString() != "" {
		err = c.getByName(ctx, &config)
	} else {
		err = c.getByPage(ctx, &config)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunEbsVolumes) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunEbsVolumes) getByID(ctx context.Context, config *CtyunEbsVolumesConfig) (err error) {
	// 组装请求体
	params := &ctebs2.EbsQueryEbsByIDRequest{
		RegionID: config.RegionID.ValueStringPointer(),
		DiskID:   config.DiskID.ValueString(),
	}

	// 调用API
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsQueryEbsByIDApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回值
	volumes := []CtyunEbsVolumesModel{}
	disk := resp.ReturnObj
	item := CtyunEbsVolumesModel{
		ID:              utils.SecStringValue(disk.DiskID),
		Name:            utils.SecStringValue(disk.DiskName),
		Size:            types.Int32Value(disk.DiskSize),
		Type:            utils.SecStringValue(disk.DiskType),
		Mode:            utils.SecStringValue(disk.DiskMode),
		Status:          utils.SecStringValue(disk.DiskStatus),
		CreateTime:      types.Int64Value(disk.CreateTime),
		UpdateTime:      types.Int64Value(disk.UpdateTime),
		ExpireTime:      types.Int64Value(disk.ExpireTime),
		IsSystemVolume:  utils.SecBoolValue(disk.IsSystemVolume),
		IsPackaged:      utils.SecBoolValue(disk.IsPackaged),
		InstanceName:    utils.SecStringValue(disk.InstanceName),
		InstanceID:      utils.SecStringValue(disk.InstanceID),
		InstanceStatus:  utils.SecStringValue(disk.InstanceStatus),
		MultiAttach:     utils.SecBoolValue(disk.MultiAttach),
		ProjectID:       utils.SecStringValue(disk.ProjectID),
		IsEncrypt:       utils.SecBoolValue(disk.IsEncrypt),
		KmsUUID:         utils.SecStringValue(disk.KmsUUID),
		RegionID:        utils.SecStringValue(disk.RegionID),
		AzName:          utils.SecStringValue(disk.AzName),
		DiskFreeze:      utils.SecBoolValue(disk.DiskFreeze),
		ProvisionedIops: types.Int32Value(disk.ProvisionedIops),
		Attachments:     []CtyunEbsVolumesAttachments{},
	}
	for _, a := range disk.Attachments {
		item.Attachments = append(item.Attachments, CtyunEbsVolumesAttachments{
			InstanceID:   utils.SecStringValue(a.InstanceID),
			AttachmentID: utils.SecStringValue(a.AttachmentID),
			Device:       utils.SecStringValue(a.Device),
		})
	}
	volumes = append(volumes, item)
	config.Volumes = volumes
	return
}

func (c *ctyunEbsVolumes) getByName(ctx context.Context, config *CtyunEbsVolumesConfig) (err error) {
	// 组装请求体
	params := &ctebs2.EbsQueryEbsByNameRequest{
		RegionID: config.RegionID.ValueString(),
		DiskName: config.DiskName.ValueString(),
	}

	// 调用API
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsQueryEbsByNameApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回值
	volumes := []CtyunEbsVolumesModel{}
	disk := resp.ReturnObj
	item := CtyunEbsVolumesModel{
		ID:              utils.SecStringValue(disk.DiskID),
		Name:            utils.SecStringValue(disk.DiskName),
		Size:            types.Int32Value(disk.DiskSize),
		Type:            utils.SecStringValue(disk.DiskType),
		Mode:            utils.SecStringValue(disk.DiskMode),
		Status:          utils.SecStringValue(disk.DiskStatus),
		CreateTime:      types.Int64Value(disk.CreateTime),
		UpdateTime:      types.Int64Value(disk.UpdateTime),
		ExpireTime:      types.Int64Value(disk.ExpireTime),
		IsSystemVolume:  utils.SecBoolValue(disk.IsSystemVolume),
		IsPackaged:      utils.SecBoolValue(disk.IsPackaged),
		InstanceName:    utils.SecStringValue(disk.InstanceName),
		InstanceID:      utils.SecStringValue(disk.InstanceID),
		InstanceStatus:  utils.SecStringValue(disk.InstanceStatus),
		MultiAttach:     utils.SecBoolValue(disk.MultiAttach),
		ProjectID:       utils.SecStringValue(disk.ProjectID),
		IsEncrypt:       utils.SecBoolValue(disk.IsEncrypt),
		KmsUUID:         utils.SecStringValue(disk.KmsUUID),
		RegionID:        utils.SecStringValue(disk.RegionID),
		AzName:          utils.SecStringValue(disk.AzName),
		DiskFreeze:      utils.SecBoolValue(disk.DiskFreeze),
		ProvisionedIops: types.Int32Value(disk.ProvisionedIops),
		Attachments:     []CtyunEbsVolumesAttachments{},
	}
	for _, a := range disk.Attachments {
		item.Attachments = append(item.Attachments, CtyunEbsVolumesAttachments{
			InstanceID:   utils.SecStringValue(a.InstanceID),
			AttachmentID: utils.SecStringValue(a.AttachmentID),
			Device:       utils.SecStringValue(a.Device),
		})
	}
	volumes = append(volumes, item)
	config.Volumes = volumes
	return
}

func (c *ctyunEbsVolumes) getByPage(ctx context.Context, config *CtyunEbsVolumesConfig) (err error) {
	// 组装请求体
	params := &ctebs2.EbsQueryEbsListRequest{
		RegionID: config.RegionID.ValueString(),
	}
	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	if pageNo > 0 {
		params.PageNo = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	if config.ProjectID.ValueString() != "" {
		params.ProjectID = config.ProjectID.ValueStringPointer()
	}

	// 调用API
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsQueryEbsListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回值
	volumes := []CtyunEbsVolumesModel{}
	for _, disk := range resp.ReturnObj.DiskList {
		if disk == nil {
			continue
		}
		item := CtyunEbsVolumesModel{
			ID:              utils.SecStringValue(disk.DiskID),
			Name:            utils.SecStringValue(disk.DiskName),
			Size:            types.Int32Value(disk.DiskSize),
			Type:            utils.SecStringValue(disk.DiskType),
			Mode:            utils.SecStringValue(disk.DiskMode),
			Status:          utils.SecStringValue(disk.DiskStatus),
			CreateTime:      types.Int64Value(disk.CreateTime),
			UpdateTime:      types.Int64Value(disk.UpdateTime),
			ExpireTime:      types.Int64Value(disk.ExpireTime),
			IsSystemVolume:  utils.SecBoolValue(disk.IsSystemVolume),
			IsPackaged:      utils.SecBoolValue(disk.IsPackaged),
			InstanceName:    utils.SecStringValue(disk.InstanceName),
			InstanceID:      utils.SecStringValue(disk.InstanceID),
			InstanceStatus:  utils.SecStringValue(disk.InstanceStatus),
			MultiAttach:     utils.SecBoolValue(disk.MultiAttach),
			ProjectID:       utils.SecStringValue(disk.ProjectID),
			IsEncrypt:       utils.SecBoolValue(disk.IsEncrypt),
			KmsUUID:         utils.SecStringValue(disk.KmsUUID),
			RegionID:        utils.SecStringValue(disk.RegionID),
			AzName:          utils.SecStringValue(disk.AzName),
			DiskFreeze:      utils.SecBoolValue(disk.DiskFreeze),
			ProvisionedIops: types.Int32Value(disk.ProvisionedIops),
			Attachments:     []CtyunEbsVolumesAttachments{},
		}
		for _, a := range disk.Attachments {
			item.Attachments = append(item.Attachments, CtyunEbsVolumesAttachments{
				InstanceID:   utils.SecStringValue(a.InstanceID),
				AttachmentID: utils.SecStringValue(a.AttachmentID),
				Device:       utils.SecStringValue(a.Device),
			})
		}
		volumes = append(volumes, item)
	}
	config.Volumes = volumes
	return
}
