package hpfs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/hpfs"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunHpfsClusters{}
	_ datasource.DataSourceWithConfigure = &ctyunHpfsClusters{}
)

type ctyunHpfsClusters struct {
	meta *common.CtyunMetadata
}

func (c *ctyunHpfsClusters) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunHpfsClusters) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_hpfs_clusters"
}

func NewCtyunHpfsClusters() datasource.DataSource {
	return &ctyunHpfsClusters{}
}

func (c *ctyunHpfsClusters) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10088932/10090437**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id",
			},
			"sfs_type": schema.StringAttribute{
				Optional:    true,
				Description: "类型，hpfs_perf(HPC性能型)",
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Description: "多可用区下的可用区名字",
			},
			"ebm_device_type": schema.StringAttribute{
				Optional:    true,
				Description: "裸金属设备规格",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的分页页码，默认值为1",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页包含的元素个数范围(1-50)，默认值为10",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},
			"hpfs_clusters": schema.ListNestedAttribute{
				Computed:    true,
				Description: "hpfs集群列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cluster_name": schema.StringAttribute{
							Computed:    true,
							Description: "集群名称",
						},
						"remaining_status": schema.BoolAttribute{
							Computed:    true,
							Description: "该集群是否可以售卖",
						},
						"sfs_type": schema.StringAttribute{
							Computed:    true,
							Description: "集群的存储类型",
						},
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "多可用区下的可用区名字",
						},
						"sfs_protocol": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "集群支持的协议列表",
						},
						"baselines": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "集群支持的性能基线列表（仅当资源池支持性能基线时返回）",
						},
						"network_type": schema.StringAttribute{
							Computed:    true,
							Description: "集群的网络类型（tcp/o2ib）",
						},
						"ebm_device_types": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "集群支持的裸金属设备规格列表",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunHpfsClusters) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunHpfsClustersConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	if regionId == "" {
		err = errors.New("region id 为空")
		return
	}
	params := &hpfs.HpfsListClusterRequest{
		RegionID: regionId,
		PageNo:   1,
		PageSize: 10,
	}
	if !config.SfsType.IsNull() {
		params.SfsType = config.SfsType.ValueString()
	}
	if !config.AzName.IsNull() {
		params.AzName = config.AzName.ValueString()
	}
	if !config.EbmDeviceType.IsNull() {
		params.EbmDeviceType = config.EbmDeviceType.ValueString()
	}
	if !config.PageSize.IsNull() && config.PageSize.ValueInt32() != 0 {
		params.PageSize = config.PageSize.ValueInt32()
	}
	if !config.PageNo.IsNull() && config.PageNo.ValueInt32() != 0 {
		params.PageNo = config.PageNo.ValueInt32()
	}

	resp, err := c.meta.Apis.SdkHpfsApis.HpfsListClusterApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = errors.New("查询hpfs 集群列表失败，返回为nil")
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var hpfsClusterList []CtyunHpfsClusterModel
	clusterList := resp.ReturnObj.ClusterList
	for _, clusterItem := range clusterList {
		var cluster CtyunHpfsClusterModel
		cluster.ClusterName = types.StringValue(clusterItem.ClusterName)
		cluster.RemainingStatus = types.BoolValue(*clusterItem.RemainingStatus)
		cluster.SfsType = types.StringValue(clusterItem.StorageType)
		cluster.AzName = types.StringValue(clusterItem.AzName)
		cluster.NetworkType = types.StringValue(clusterItem.NetworkType)
		protocolType, diags := types.SetValueFrom(ctx, types.StringType, clusterItem.ProtocolType)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		cluster.SfsProtocol = protocolType
		baselines, diags := types.SetValueFrom(ctx, types.StringType, clusterItem.Baselines)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		cluster.Baselines = baselines

		ebmDeviceTypes, diags := types.SetValueFrom(ctx, types.StringType, clusterItem.EbmDeviceTypes)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		cluster.EbmDeviceTypes = ebmDeviceTypes

		hpfsClusterList = append(hpfsClusterList, cluster)
	}
	config.HpfsClusters = hpfsClusterList
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		err = errors.New(response.Diagnostics[0].Detail())
		return
	}
}

type CtyunHpfsClusterModel struct {
	ClusterName     types.String `tfsdk:"cluster_name"`     // 集群名称
	RemainingStatus types.Bool   `tfsdk:"remaining_status"` // 是否可以售卖
	SfsType         types.String `tfsdk:"sfs_type"`         // 集群的存储类型
	AzName          types.String `tfsdk:"az_name"`          // 多可用区下的可用区名字
	SfsProtocol     types.Set    `tfsdk:"sfs_protocol"`     // 支持的协议列表
	Baselines       types.Set    `tfsdk:"baselines"`        // 性能基线列表
	NetworkType     types.String `tfsdk:"network_type"`     // 集群的网络类型
	EbmDeviceTypes  types.Set    `tfsdk:"ebm_device_types"` // 裸金属设备规格列表
}

type CtyunHpfsClustersConfig struct {
	RegionID      types.String            `tfsdk:"region_id"`       // 资源池 ID
	SfsType       types.String            `tfsdk:"sfs_type"`        // 文件系统类型
	AzName        types.String            `tfsdk:"az_name"`         // 可用区名称
	EbmDeviceType types.String            `tfsdk:"ebm_device_type"` // 裸金属设备规格
	PageNo        types.Int32             `tfsdk:"page_no"`         // 分页页码
	PageSize      types.Int32             `tfsdk:"page_size"`       // 每页元素数量
	HpfsClusters  []CtyunHpfsClusterModel `tfsdk:"hpfs_clusters"`   // hpfs cluster列表
}
