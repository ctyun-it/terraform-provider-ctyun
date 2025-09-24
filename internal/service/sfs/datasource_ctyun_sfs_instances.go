package sfs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/sfs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunSfsInstances{}
	_ datasource.DataSourceWithConfigure = &CtyunSfsInstances{}
)

type CtyunSfsInstances struct {
	meta *common.CtyunMetadata
}

func NewCtyunSfsInstances() datasource.DataSource {
	return &CtyunSfsInstances{}
}

func (c *CtyunSfsInstances) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunSfsInstances) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_sfs_instances"
}

func (c *CtyunSfsInstances) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027350**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "企业项目ID，默认为0",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页包含的元素个数，默认为10",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的分页页码，默认为1",
			},
			"sfs_list": schema.ListNestedAttribute{
				Computed:    true,
				Description: "弹性文件存储列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sfs_name": schema.StringAttribute{
							Computed:    true,
							Description: "弹性文件系统名称",
						},
						"sfs_uid": schema.StringAttribute{
							Computed:    true,
							Description: "弹性文件系统唯一ID",
						},
						"sfs_size": schema.Int32Attribute{
							Computed:    true,
							Description: "文件系统大小（GB）",
						},
						"sfs_type": schema.StringAttribute{
							Computed:    true,
							Description: "文件系统类型，capacity(标准型)或performance(性能型)",
						},
						"sfs_protocol": schema.StringAttribute{
							Computed:    true,
							Description: "挂载协议，nfs或cifs",
						},
						"sfs_status": schema.StringAttribute{
							Computed:    true,
							Description: "文件系统状态，creating/available/unusable/expired/fail",
						},
						"used_size": schema.Int32Attribute{
							Computed:    true,
							Description: "已使用大小（MB）",
						},
						"create_time": schema.Int64Attribute{
							Computed:    true,
							Description: "创建时间，Unix时间戳（毫秒）",
						},
						"update_time": schema.Int64Attribute{
							Computed:    true,
							Description: "更新时间，Unix时间戳（毫秒）",
						},
						"expire_time": schema.Int64Attribute{
							Computed:    true,
							Description: "过期时间，Unix时间戳（毫秒）",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源所属企业项目ID",
						},
						"is_encrypt": schema.BoolAttribute{
							Computed:    true,
							Description: "是否为加密盘",
						},
						"kms_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "加密盘密钥UUID",
						},
						"vpce_share_path": schema.ListNestedAttribute{
							Computed:    true,
							Description: "VPCE共享路径信息，仅适用于v4.0资源池",
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
										Description: "Linux主机共享路径",
									},
									"share_path_v6": schema.StringAttribute{
										Computed:    true,
										Description: "Linux主机IPv6共享路径",
									},
									"windows_share_path": schema.StringAttribute{
										Computed:    true,
										Description: "Windows主机共享路径",
									},
									"windows_share_path_v6": schema.StringAttribute{
										Computed:    true,
										Description: "Windows主机IPv6共享路径",
									},
								},
							},
						},
						"on_demand": schema.BoolAttribute{
							Computed:    true,
							Description: "是否按需订购",
						},
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源池ID",
						},
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "多可用区下的可用区名称",
						},
						"share_path": schema.StringAttribute{
							Computed:    true,
							Description: "Linux主机共享路径",
						},
						"share_path_v6": schema.StringAttribute{
							Computed:    true,
							Description: "Linux主机IPv6共享路径",
						},
						"windows_share_path": schema.StringAttribute{
							Computed:    true,
							Description: "Windows主机共享路径",
						},
						"windows_share_path_v6": schema.StringAttribute{
							Computed:    true,
							Description: "Windows主机IPv6共享路径",
						},
						"mount_count": schema.Int32Attribute{
							Computed:    true,
							Description: "挂载点数量",
						},
						"ceph_id": schema.StringAttribute{
							Computed:    true,
							Description: "Ceph底层ID",
						},
						"phy_share_path": schema.StringAttribute{
							Computed:    true,
							Description: "Linux物理机共享路径，仅适用于v3.0资源池",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunSfsInstances) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunSfsInstancesConfig
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
	params := &sfs.SfsSfsInfoSfsRequest{
		RegionID: regionId,
		PageSize: 10,
		PageNo:   1,
	}
	if !config.ProjectID.IsNull() {
		params.ProjectID = config.ProjectID.ValueString()
	}
	if !config.PageNo.IsNull() && !config.PageNo.IsUnknown() {
		params.PageNo = config.PageNo.ValueInt32()
	}

	if !config.PageSize.IsNull() && !config.PageSize.IsUnknown() {
		params.PageSize = config.PageSize.ValueInt32()
	}

	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsInfoSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询region_id=%s的sfs列表失败，返回接口为nil。请与研发联系确认问题原因。", regionId)
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回的list
	var sfsList []CtyunSfsInfoModel
	sfsListResp := resp.ReturnObj.List
	for _, sfsItem := range sfsListResp {
		var sfsInfo CtyunSfsInfoModel
		sfsInfo.SfsName = types.StringValue(sfsItem.SfsName)
		sfsInfo.SfsUID = types.StringValue(sfsItem.SfsUID)
		sfsInfo.SfsSize = types.Int32Value(sfsItem.SfsSize)
		sfsInfo.SfsType = types.StringValue(sfsItem.SfsType)
		sfsInfo.SfsProtocol = types.StringValue(sfsItem.SfsProtocol)
		sfsInfo.SfsStatus = types.StringValue(sfsItem.SfsStatus)
		sfsInfo.UsedSize = types.Int32Value(sfsItem.UsedSize)
		sfsInfo.CreateTime = types.Int64Value(sfsItem.CreateTime)
		sfsInfo.UpdateTime = types.Int64Value(sfsItem.UpdateTime)
		sfsInfo.ExpireTime = types.Int64Value(sfsItem.ExpireTime)
		sfsInfo.ProjectID = types.StringValue(sfsItem.ProjectID)
		sfsInfo.IsEncrypt = types.BoolValue(sfsItem.IsEncrypt)
		sfsInfo.KmsUUID = types.StringValue(sfsItem.KmsUUID)
		sfsInfo.OnDemand = types.BoolValue(sfsItem.OnDemand)
		sfsInfo.RegionID = types.StringValue(sfsItem.RegionID)
		sfsInfo.AzName = types.StringValue(sfsItem.AzName)
		sfsInfo.SharePath = types.StringValue(sfsItem.SharePath)
		sfsInfo.SharePathV6 = types.StringValue(sfsItem.SharePathV6)
		sfsInfo.WindowsSharePath = types.StringValue(sfsItem.WindowsSharePath)
		sfsInfo.WindowsSharePathV6 = types.StringValue(sfsItem.WindowsSharePathV6)
		sfsInfo.MountCount = types.Int32Value(sfsItem.MountCount)
		sfsInfo.CephID = types.StringValue(sfsItem.CephID)
		sfsInfo.PhySharePath = types.StringValue(sfsItem.PhySharePath)

		var sharePathList []CtyunSfsVpceSharePathModel
		for _, pathItem := range sfsItem.VpceSharePath {
			var path CtyunSfsVpceSharePathModel
			path.VpcID = types.StringValue(pathItem.VpcID)
			path.VpcName = types.StringValue(pathItem.VpcName)
			path.SharePath = types.StringValue(pathItem.SharePath)
			path.SharePathV6 = types.StringValue(pathItem.SharePathV6)
			path.WindowsSharePath = types.StringValue(pathItem.WindowsSharePath)
			path.WindowsSharePathV6 = types.StringValue(pathItem.WindowsSharePathV6)
			sharePathList = append(sharePathList, path)
		}
		var diags diag.Diagnostics
		sfsInfo.VpceSharePath, diags = types.ListValueFrom(ctx, utils.StructToTFObjectTypes(CtyunSfsVpceSharePathModel{}), &sharePathList)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		sfsList = append(sfsList, sfsInfo)
	}
	config.SfsList = sfsList
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunSfsInfoModel struct {
	SfsName            types.String `tfsdk:"sfs_name"`
	SfsUID             types.String `tfsdk:"sfs_uid"`
	SfsSize            types.Int32  `tfsdk:"sfs_size"`
	SfsType            types.String `tfsdk:"sfs_type"`
	SfsProtocol        types.String `tfsdk:"sfs_protocol"`
	SfsStatus          types.String `tfsdk:"sfs_status"`
	UsedSize           types.Int32  `tfsdk:"used_size"`
	CreateTime         types.Int64  `tfsdk:"create_time"`
	UpdateTime         types.Int64  `tfsdk:"update_time"`
	ExpireTime         types.Int64  `tfsdk:"expire_time"`
	ProjectID          types.String `tfsdk:"project_id"`
	IsEncrypt          types.Bool   `tfsdk:"is_encrypt"`
	KmsUUID            types.String `tfsdk:"kms_uuid"`
	VpceSharePath      types.List   `tfsdk:"vpce_share_path"` // List of CtyunSfsVpceSharePath
	OnDemand           types.Bool   `tfsdk:"on_demand"`
	RegionID           types.String `tfsdk:"region_id"`
	AzName             types.String `tfsdk:"az_name"`
	SharePath          types.String `tfsdk:"share_path"`
	SharePathV6        types.String `tfsdk:"share_path_v6"`
	WindowsSharePath   types.String `tfsdk:"windows_share_path"`
	WindowsSharePathV6 types.String `tfsdk:"windows_share_path_v6"`
	MountCount         types.Int32  `tfsdk:"mount_count"`
	CephID             types.String `tfsdk:"ceph_id"`
	PhySharePath       types.String `tfsdk:"phy_share_path"`
}

type CtyunSfsVpceSharePathModel struct {
	VpcID              types.String `tfsdk:"vpc_id"`
	VpcName            types.String `tfsdk:"vpc_name"`
	SharePath          types.String `tfsdk:"share_path"`
	SharePathV6        types.String `tfsdk:"share_path_v6"`
	WindowsSharePath   types.String `tfsdk:"windows_share_path"`
	WindowsSharePathV6 types.String `tfsdk:"windows_share_path_v6"`
}

type CtyunSfsInstancesConfig struct {
	RegionID  types.String        `tfsdk:"region_id"`
	ProjectID types.String        `tfsdk:"project_id"`
	PageSize  types.Int32         `tfsdk:"page_size"`
	PageNo    types.Int32         `tfsdk:"page_no"`
	SfsList   []CtyunSfsInfoModel `tfsdk:"sfs_list"`
}
