package ccse

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ccse2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ccse"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ datasource.DataSource              = &ctyunCcseNodePools{}
	_ datasource.DataSourceWithConfigure = &ctyunCcseNodePools{}
)

type ctyunCcseNodePools struct {
	meta *common.CtyunMetadata
}

func NewCtyunCcseNodePools() datasource.DataSource {
	return &ctyunCcseNodePools{}
}

func (c *ctyunCcseNodePools) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ccse_node_pools"
}

type CtyunCcseNodePoolsModel struct {
	ID                       types.String            `tfsdk:"id"`
	NodePoolName             types.String            `tfsdk:"name"`
	CycleCount               types.Int64             `tfsdk:"cycle_count"`
	CycleType                types.String            `tfsdk:"cycle_type"`
	AutoRenew                types.Bool              `tfsdk:"auto_renew"`
	VisibilityPostHostScript types.String            `tfsdk:"visibility_post_host_script"`
	VisibilityHostScript     types.String            `tfsdk:"visibility_host_script"`
	InstanceType             types.String            `tfsdk:"instance_type"`
	MirrorID                 types.String            `tfsdk:"mirror_id"`
	MirrorName               types.String            `tfsdk:"mirror_name"`
	ItemDefName              types.String            `tfsdk:"item_def_name"`
	SysDisk                  CtyunCcseNodePoolDisk   `tfsdk:"sys_disk"`
	DataDisks                []CtyunCcseNodePoolDisk `tfsdk:"data_disks"`
	MaxPodNum                types.Int32             `tfsdk:"max_pod_num"`
}

type CtyunCcseNodePoolsConfig struct {
	ClusterID    types.String              `tfsdk:"cluster_id"`
	RegionID     types.String              `tfsdk:"region_id"`
	PageNo       types.Int32               `tfsdk:"page_no"`
	PageSize     types.Int32               `tfsdk:"page_size"`
	NodePoolName types.String              `tfsdk:"name"`
	Records      []CtyunCcseNodePoolsModel `tfsdk:"records"`
}

func (c *ctyunCcseNodePools) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10083472/10318452**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "集群ID",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "节点池名称",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页数据量大小",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},

			"records": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "节点池id",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "节点池名称",
						},
						"cycle_type": schema.StringAttribute{
							Computed:    true,
							Description: "订购周期类型，取值范围：month：按月，year：按年、on_demand：按需",
						},
						"cycle_count": schema.Int64Attribute{
							Computed:    true,
							Description: "订购时长",
						},
						"auto_renew": schema.BoolAttribute{
							Computed:    true,
							Description: "是否自动续订",
						},
						"visibility_post_host_script": schema.StringAttribute{
							Computed:    true,
							Description: "部署后执行自定义脚本，base64编码",
						},
						"visibility_host_script": schema.StringAttribute{
							Computed:    true,
							Description: "部署前执行自定义脚本，base64编码",
						},
						"instance_type": schema.StringAttribute{
							Computed:    true,
							Description: "实例类型，支持ecs（云主机）、ebm（裸金属）",
						},
						"mirror_id": schema.StringAttribute{
							Computed:    true,
							Description: "镜像id",
						},
						"mirror_name": schema.StringAttribute{
							Computed:    true,
							Description: "镜像名称",
						},
						"item_def_name": schema.StringAttribute{
							Computed:    true,
							Description: "规格名称",
						},
						"max_pod_num": schema.Int32Attribute{
							Optional:    true,
							Description: "最大pod数,默认110",
						},
						"sys_disk": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "系统盘",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed:    true,
									Description: "系统盘规格",
								},
								"size": schema.Int32Attribute{
									Computed:    true,
									Description: "系统盘大小，单位为G",
								},
							},
						},
						"data_disks": schema.ListNestedAttribute{
							Optional:    true,
							Description: "数据盘",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Computed:    true,
										Description: "系统盘规格",
									},
									"size": schema.Int32Attribute{
										Computed:    true,
										Description: "系统盘大小，单位为G",
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

func (c *ctyunCcseNodePools) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunCcseNodePoolsConfig
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
	// 组装请求体
	params := &ccse2.CcseListNodePoolsRequest{
		RegionId:  regionId,
		ClusterId: config.ClusterID.ValueString(),
		PageNow:   1,
		PageSize:  10,
	}
	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	name := config.NodePoolName.ValueString()
	if pageNo > 0 {
		params.PageNow = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	if name != "" {
		params.NodePoolName = name
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCcseApis.CcseListNodePoolsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.ReturnObj == nil {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	// 解析返回值
	config.Records = []CtyunCcseNodePoolsModel{}
	for _, p := range resp.ReturnObj.Records {
		item := CtyunCcseNodePoolsModel{
			ID:                       types.StringValue(p.Id),
			NodePoolName:             types.StringValue(p.NodePoolName),
			AutoRenew:                types.BoolValue(map[int32]bool{0: false, 1: true}[p.AutoRenewStatus]),
			VisibilityPostHostScript: types.StringValue(p.VisibilityPostHostScript),
			VisibilityHostScript:     types.StringValue(p.VisibilityHostScript),
			MirrorID:                 types.StringValue(p.ImageUuid),
			MirrorName:               types.StringValue(p.ImageName),
			ItemDefName:              types.StringValue(p.VmSpecName),
			MaxPodNum:                types.Int32Value(p.MaxPodNum),
		}
		switch p.BillMode {
		case "1":
			item.CycleType = types.StringValue(strings.ToLower(p.CycleType))
			item.CycleCount = types.Int64Value(int64(p.CycleCount))
		case "2":
			item.CycleType = types.StringValue(business.OnDemandCycleType)
		}
		if strings.HasPrefix(p.VmSpecName, "physical") {
			item.InstanceType = types.StringValue("ebm")
		} else {
			item.InstanceType = types.StringValue("ecs")
		}
		for _, disk := range p.DataDisks {
			item.DataDisks = append(item.DataDisks, CtyunCcseNodePoolDisk{
				Size: types.Int32Value(int32(disk.Size)),
				Type: types.StringValue(disk.DiskSpecName),
			})
		}
		item.SysDisk.Type = types.StringValue(p.SysDiskType)
		item.SysDisk.Size = types.Int32Value(p.SysDiskSize)
		config.Records = append(config.Records, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunCcseNodePools) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
