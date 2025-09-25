package kafka

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgkafka "github.com/ctyun-it/terraform-provider-ctyun/internal/core/kafka"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunKafkaInstances{}
	_ datasource.DataSourceWithConfigure = &ctyunKafkaInstances{}
)

type ctyunKafkaInstances struct {
	meta *common.CtyunMetadata
}

func NewCtyunKafkaInstances() datasource.DataSource {
	return &ctyunKafkaInstances{}
}

func (c *ctyunKafkaInstances) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_kafka_instances"
}

type CtyunKafkaInstancesModel struct {
	ID            types.String `tfsdk:"id"`
	Status        types.Int32  `tfsdk:"status"`
	StatusDesc    types.String `tfsdk:"status_desc"`
	ProjectID     types.String `tfsdk:"project_id"`
	InstanceName  types.String `tfsdk:"instance_name"`
	EngineVersion types.String `tfsdk:"engine_version"`
	SpecName      types.String `tfsdk:"spec_name"`
	NodeNum       types.Int32  `tfsdk:"node_num"`
	DiskType      types.String `tfsdk:"disk_type"`
	DiskSize      types.Int32  `tfsdk:"disk_size"`
	VpcID         types.String `tfsdk:"vpc_id"`
	SubnetID      types.String `tfsdk:"subnet_id"`
	EnableIpv6    types.Bool   `tfsdk:"enable_ipv6"`
}

type CtyunKafkaInstancesConfig struct {
	RegionID     types.String               `tfsdk:"region_id"`
	PageNo       types.Int32                `tfsdk:"page_no"`
	PageSize     types.Int32                `tfsdk:"page_size"`
	InstanceName types.String               `tfsdk:"instance_name"`
	InstanceID   types.String               `tfsdk:"instance_id"`
	ProjectID    types.String               `tfsdk:"project_id"`
	Instances    []CtyunKafkaInstancesModel `tfsdk:"instances"`
}

func (c *ctyunKafkaInstances) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029624/10030700`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"project_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "企业项目ID",
			},
			"instance_name": schema.StringAttribute{
				Optional:    true,
				Description: "实例名称",
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "实例ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页数据量大小，支持范围1-50",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},
			"instances": schema.ListNestedAttribute{
				Computed:    true,
				Description: "实例列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "实例id",
						},
						"status": schema.Int32Attribute{
							Computed:    true,
							Description: "状态",
						},
						"status_desc": schema.StringAttribute{
							Computed:    true,
							Description: "状态描述",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "企业项目",
						},
						"instance_name": schema.StringAttribute{
							Computed:    true,
							Description: "实例名称",
						},
						"engine_version": schema.StringAttribute{
							Computed:    true,
							Description: "实例引擎版本",
						},
						"spec_name": schema.StringAttribute{
							Computed:    true,
							Description: "实例的规格类型",
						},
						"node_num": schema.Int32Attribute{
							Computed:    true,
							Description: "节点数。单机版为1个，集群版3~50个",
						},
						"disk_type": schema.StringAttribute{
							Computed:    true,
							Description: "磁盘类型",
						},
						"disk_size": schema.Int32Attribute{
							Computed:    true,
							Description: "单个节点的磁盘存储空间，单位为GB，存储空间取值范围100GB ~ 10000，并且为100的倍数。实例总存储空间为diskSize * nodeNum",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "关联的vpcID",
						},
						"subnet_id": schema.StringAttribute{
							Computed:    true,
							Description: "子网ID",
						},
						"enable_ipv6": schema.BoolAttribute{
							Computed:    true,
							Description: "是否启用IPv6，默认为false",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunKafkaInstances) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunKafkaInstancesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}

	projectID := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraProjectId)
	config.RegionID = types.StringValue(regionId)
	config.ProjectID = types.StringValue(projectID)

	// 组装请求体
	params := &ctgkafka.CtgkafkaInstQueryRequest{
		RegionId:       regionId,
		OuterProjectId: projectID,
		PageSize:       config.PageSize.ValueInt32(),
		PageNum:        config.PageNo.ValueInt32(),
		ProdInstId:     config.InstanceID.ValueString(),
		Name:           config.InstanceName.ValueString(),
	}
	if config.InstanceName.ValueString() != "" {
		e := true
		params.ExactMatchName = &e
	}
	// 调用API
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaInstQueryApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回值
	config.Instances = []CtyunKafkaInstancesModel{}
	for _, instance := range resp.ReturnObj.Data {
		item := CtyunKafkaInstancesModel{}
		item.Status = types.Int32Value(instance.Status)
		item.StatusDesc = types.StringValue(
			map[int32]string{
				1:   "运行中",
				2:   "已过期",
				3:   "已注销",
				4:   "变更中",
				5:   "已退订",
				6:   "开通中",
				7:   "已取消",
				8:   "已停止",
				9:   "弹性IP处理中",
				10:  "重启中",
				11:  "重启失败",
				12:  "升级中",
				13:  "已欠费",
				101: "开通失败",
			}[instance.Status],
		)
		item.ProjectID = types.StringValue(instance.OuterProjectId)
		item.InstanceName = types.StringValue(instance.InstanceName)
		if len(instance.Version) >= 3 {
			item.EngineVersion = types.StringValue(instance.Version[:3])
		}
		item.SpecName = types.StringValue(instance.Specifications)
		item.NodeNum = types.Int32Value(int32(len(instance.NodeList)))

		item.DiskType = types.StringValue(instance.DiskType)
		item.DiskSize = types.Int32Value(utils.StringToInt32Must(instance.Space))
		item.VpcID = types.StringValue(instance.VpcId)
		item.SubnetID = types.StringValue(instance.SubnetId)

		item.EnableIpv6 = types.BoolValue(map[int32]bool{1: true, 0: false}[instance.Ipv6Enable])

		config.Instances = append(config.Instances, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunKafkaInstances) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
