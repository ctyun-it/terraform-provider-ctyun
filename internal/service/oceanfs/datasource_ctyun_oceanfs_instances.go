package oceanfs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/oceanfs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CtyunOceanfsInstances struct {
	meta *common.CtyunMetadata
}

func NewCtyunOceanfsInstances() datasource.DataSource {
	return &CtyunOceanfsInstances{}
}

func (c *CtyunOceanfsInstances) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunOceanfsInstances) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_oceanfs_instances"
}

func (c *CtyunOceanfsInstances) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10088966/10115906",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "区域ID",
			},
			"project_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "项目ID",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页包含的元素个数，默认为10",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的分页页码，默认为1",
			},
			"instances": schema.ListNestedAttribute{
				Computed:    true,
				Description: "OceanFS实例列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "文件存储名称",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "文件存储唯一标识",
						},
						"size": schema.Int32Attribute{
							Computed:    true,
							Description: "文件存储容量大小(GB)",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "文件存储类型",
						},
						"protocol": schema.StringAttribute{
							Computed:    true,
							Description: "文件存储协议",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "文件存储状态",
						},
						"used_size": schema.Int32Attribute{
							Computed:    true,
							Description: "文件系统已使用容量，单位MB",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间，为UTC格式",
						},
						"expire_time": schema.StringAttribute{
							Computed:    true,
							Description: "到期时间，为UTC格式，按需时为空",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "项目ID",
						},
						"on_demand": schema.BoolAttribute{
							Computed:    true,
							Description: "是否按需计费",
						},
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "区域ID",
						},
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "可用区名称",
						},
						"share_path": schema.StringAttribute{
							Computed:    true,
							Description: "NFS文件系统用于Linux操作系统及IPv4挂载访问的挂载地址。\n注：不可用于挂载物理机（包括标准裸金属、弹性裸金属），可以用于挂载云主机、容器等除物理机以外的计算服务。不可用于专线访问",
						},
						"share_path_v6": schema.StringAttribute{
							Computed:    true,
							Description: "NFS文件系统用于Linux操作系统及IPv6挂载访问的挂载地址。\n注：不可用于挂载物理机（包括标准裸金属、弹性裸金属），可以用于挂载云主机、容器等除物理机以外的计算服务。不可用于专线访问",
						},
						"windows_share_path": schema.StringAttribute{
							Computed:    true,
							Description: "CIFS文件系统用于Windows操作系统IPv4挂载访问的挂载地址。\n注：不可用于挂载物理机（包括标准裸金属、弹性裸金属），可以用于挂载云主机、容器等除物理机以外的计算服务。不可用于专线访问",
						},
						"windows_share_path_v6": schema.StringAttribute{
							Computed:    true,
							Description: "CIFS文件系统用于Windows操作系统IPv6挂载访问的挂载地址。\n注：不可用于挂载物理机（包括标准裸金属、弹性裸金属），可以用于挂载云主机、容器等除物理机以外的计算服务。不可用于专线访问",
						},
						"mount_count": schema.Int32Attribute{
							Computed:    true,
							Description: "文件系统绑定的VPC数量",
						},
						"ceph_id": schema.StringAttribute{
							Computed:    true,
							Description: "监控实例ID。仅用于云监控服务",
						},
						"used_size_charge": schema.BoolAttribute{
							Computed:    true,
							Description: "是否为按实际使用量付费资源",
						},
						"vpce_share_path": schema.ListNestedAttribute{
							Computed:    true,
							Description: "VPC终端节点（VPCE）专属挂载地址",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"vpc_id": schema.StringAttribute{
										Computed:    true,
										Description: "VPC ID",
									},
									"vpc_name": schema.StringAttribute{
										Computed:    true,
										Description: "VPC名称",
									},
									"share_path": schema.StringAttribute{
										Computed:    true,
										Description: "NFS文件系统用于Linux操作系统及IPv4挂载访问的挂载地址，可用于物理机（弹性/标准裸金属）、容器、云主机、专线访问、HPC集群等各种计算服务访问文件存储",
									},
									"share_path_v6": schema.StringAttribute{
										Computed:    true,
										Description: "NFS文件系统用于Linux操作系统及IPv6挂载访问的挂载地址，可用于物理机（弹性/标准裸金属）、容器、云主机、专线访问、HPC集群等各种计算服务访问文件存储",
									},
									"windows_share_path": schema.StringAttribute{
										Computed:    true,
										Description: "CIFS文件系统用于Windows操作系统IPv4挂载访问的挂载地址，可用于物理机（弹性/标准裸金属）、容器、云主机、专线访问、HPC集群等各种计算服务访问文件存储",
									},
									"windows_share_path_v6": schema.StringAttribute{
										Computed:    true,
										Description: "CIFS文件系统用于Windows操作系统IPv6挂载访问的挂载地址，可用于物理机（弹性/标准裸金属）、容器、云主机、专线访问、HPC集群等各种计算服务访问文件存储",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (c *CtyunOceanfsInstances) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunOceanfsInstancesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	if regionId == "" {
		err = errors.New("region id 为空")
		return
	}
	config.RegionID = types.StringValue(regionId)
	params := &oceanfs.OceanfsListSfsRequest{
		RegionID: config.RegionID.ValueString(),
		PageSize: 10,
		PageNo:   1,
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		params.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsListSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询海量文件存储列表失败，接口返回nil，请联系研发确认问题原因！")
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", resp.Message)
		return
	}

	var oceanfsInstances []CtyunOceanfsInfoModel
	for _, item := range resp.ReturnObj.List {
		var info CtyunOceanfsInfoModel
		info.SfsName = types.StringValue(item.SfsName)
		info.SfsUID = types.StringValue(item.SfsUID)
		info.SfsSize = types.Int32Value(item.SfsSize)
		info.SfsType = types.StringValue(item.SfsType)
		info.SfsProtocol = types.StringValue(item.SfsProtocol)
		info.SfsStatus = types.StringValue(item.SfsStatus)
		info.UsedSize = types.Int32Value(item.UsedSize)
		info.CreateTime = types.StringValue(utils.FromUnixToUTC(item.CreateTime))
		info.UpdateTime = types.StringValue(utils.FromUnixToUTC(item.UpdateTime))
		info.ExpireTime = types.StringValue(utils.FromUnixToUTC(item.ExpireTime))
		info.ProjectID = types.StringValue(item.ProjectID)
		info.OnDemand = types.BoolValue(false)
		info.RegionID = types.StringValue(item.RegionID)
		info.AzName = types.StringValue(item.AzName)
		info.SharePath = types.StringValue(item.SharePath)
		info.WindowsSharePath = types.StringValue(item.WindowsSharePath)
		info.WindowsSharePathV6 = types.StringValue(item.WindowsSharePathV6)
		info.MountCount = types.Int32Value(item.MountCount)
		info.CephID = types.StringValue(item.CephID)
		info.UsedSizeCharge = types.BoolValue(item.UsedSizeCharge)
		var vpcesharePath []VpceSharePathModel
		for _, v := range item.VpceSharePath {
			vpcesharePath = append(vpcesharePath, VpceSharePathModel{
				VpcID:              types.StringValue(v.VpcID),
				VpcName:            types.StringValue(v.VpcName),
				SharePath:          types.StringValue(v.SharePath),
				SharePathV6:        types.StringValue(v.SharePathV6),
				WindowsSharePath:   types.StringValue(v.WindowsSharePath),
				WindowsSharePathV6: types.StringValue(v.WindowsSharePathV6),
			})
		}
		info.VpceSharePath = vpcesharePath
		oceanfsInstances = append(oceanfsInstances, info)
	}
	config.OceanfsInstances = oceanfsInstances
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunOceanfsInfoModel struct {
	SfsName            types.String         `tfsdk:"name"`
	SfsUID             types.String         `tfsdk:"id"`
	SfsSize            types.Int32          `tfsdk:"size"`
	SfsType            types.String         `tfsdk:"type"`
	SfsProtocol        types.String         `tfsdk:"protocol"`
	SfsStatus          types.String         `tfsdk:"status"`
	UsedSize           types.Int32          `tfsdk:"used_size"`
	CreateTime         types.String         `tfsdk:"create_time"`
	UpdateTime         types.String         `tfsdk:"update_time"`
	ExpireTime         types.String         `tfsdk:"expire_time"`
	ProjectID          types.String         `tfsdk:"project_id"`
	OnDemand           types.Bool           `tfsdk:"on_demand"`
	RegionID           types.String         `tfsdk:"region_id"`
	AzName             types.String         `tfsdk:"az_name"`
	SharePath          types.String         `tfsdk:"share_path"`
	SharePathV6        types.String         `tfsdk:"share_path_v6"`
	WindowsSharePath   types.String         `tfsdk:"windows_share_path"`
	WindowsSharePathV6 types.String         `tfsdk:"windows_share_path_v6"`
	MountCount         types.Int32          `tfsdk:"mount_count"`
	CephID             types.String         `tfsdk:"ceph_id"`
	UsedSizeCharge     types.Bool           `tfsdk:"used_size_charge"`
	VpceSharePath      []VpceSharePathModel `tfsdk:"vpce_share_path"`
}
type CtyunOceanfsInstancesConfig struct {
	RegionID         types.String            `tfsdk:"region_id"`
	ProjectID        types.String            `tfsdk:"project_id"`
	PageSize         types.Int32             `tfsdk:"page_size"`
	PageNo           types.Int32             `tfsdk:"page_no"`
	OceanfsInstances []CtyunOceanfsInfoModel `tfsdk:"instances"`
}

type VpceSharePathModel struct {
	VpcID              types.String `tfsdk:"vpc_id"`
	VpcName            types.String `tfsdk:"vpc_name"`
	SharePath          types.String `tfsdk:"share_path"`
	SharePathV6        types.String `tfsdk:"share_path_v6"`
	WindowsSharePath   types.String `tfsdk:"windows_share_path"`
	WindowsSharePathV6 types.String `tfsdk:"windows_share_path_v6"`
}
