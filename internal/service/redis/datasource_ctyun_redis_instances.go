package redis

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
)

var (
	_ datasource.DataSource              = &ctyunRedisInstances{}
	_ datasource.DataSourceWithConfigure = &ctyunRedisInstances{}
)

type ctyunRedisInstances struct {
	meta *common.CtyunMetadata
}

func NewCtyunRedisInstances() datasource.DataSource {
	return &ctyunRedisInstances{}
}

func (c *ctyunRedisInstances) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_instances"
}

type CtyunRedisInstancesModel struct {
	ID              types.String `tfsdk:"id"`
	InstanceName    types.String `tfsdk:"instance_name"`
	AzName          types.String `tfsdk:"az_name"`
	SecondaryAzName types.String `tfsdk:"secondary_az_name"`
	EngineVersion   types.String `tfsdk:"engine_version"`
	DataDiskType    types.String `tfsdk:"data_disk_type"`
	ShardMemSize    types.Int32  `tfsdk:"shard_mem_size"`
	ShardCount      types.Int32  `tfsdk:"shard_count"`
	CopiesCount     types.Int32  `tfsdk:"copies_count"`
	EipAddress      types.String `tfsdk:"eip_address"`
	CreateTime      types.String `tfsdk:"create_time"`
}

type CtyunRedisInstancesConfig struct {
	RegionID     types.String               `tfsdk:"region_id"`
	PageNo       types.Int32                `tfsdk:"page_no"`
	PageSize     types.Int32                `tfsdk:"page_size"`
	InstanceName types.String               `tfsdk:"instance_name"`
	ProjectID    types.String               `tfsdk:"project_id"`
	Instances    []CtyunRedisInstancesModel `tfsdk:"instances"`
}

func (c *ctyunRedisInstances) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10029420/11030280**`,
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
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "主可用区",
						},
						"secondary_az_name": schema.StringAttribute{
							Computed:    true,
							Description: "备可用区",
						},
						"engine_version": schema.StringAttribute{
							Computed:    true,
							Description: "Redis引擎版本",
						},
						"data_disk_type": schema.StringAttribute{
							Computed:    true,
							Description: "磁盘类型",
						},
						"shard_mem_size": schema.Int32Attribute{
							Computed:    true,
							Description: "分片规格",
						},
						"shard_count": schema.Int32Attribute{
							Computed:    true,
							Description: "分片数量",
						},
						"copies_count": schema.Int32Attribute{
							Computed:    true,
							Description: "副本数量",
						},
						"instance_name": schema.StringAttribute{
							Computed:    true,
							Description: "实例名称",
						},
						"eip_address": schema.StringAttribute{
							Computed:    true,
							Description: "绑定的弹性IP地址",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时长",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRedisInstances) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRedisInstancesConfig
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
	projectID := c.meta.GetExtraIfEmpty(config.ProjectID.ValueString(), common.ExtraProjectId)
	config.ProjectID = types.StringValue(projectID)

	// 组装请求体
	params := &dcs2.Dcs2DescribeInstancesRequest{
		RegionId: regionId,
	}
	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	instanceName := config.InstanceName.ValueString()
	if instanceName != "" {
		params.InstanceName = instanceName
	}
	if projectID != "" {
		params.ProjectId = projectID
	}
	if pageNo > 0 {
		params.PageIndex = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeInstancesApi.Do(ctx, c.meta.SdkCredential, params)
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
	config.Instances = []CtyunRedisInstancesModel{}
	for _, instance := range resp.ReturnObj.Rows {
		item := CtyunRedisInstancesModel{}
		if len(instance.AzList) > 0 {
			item.AzName = types.StringValue(instance.AzList[0].AzEngName)
		}
		if len(instance.AzList) > 1 {
			item.SecondaryAzName = types.StringValue(instance.AzList[1].AzEngName)
		}

		item.EngineVersion = types.StringValue(instance.EngineVersion)
		item.DataDiskType = types.StringValue(instance.DataDiskType)
		shardMemSize, _ := strconv.Atoi(instance.ShardMemSize)
		if shardMemSize == 0 {
			shardMemSize, _ = strconv.Atoi(instance.Capacity)
		}
		item.ShardMemSize = types.Int32Value(int32(shardMemSize))
		shardCount, _ := strconv.Atoi(instance.ShardCount)
		item.ShardCount = types.Int32Value(int32(shardCount))
		copiesCount, _ := strconv.Atoi(instance.CopiesCount)
		item.CopiesCount = types.Int32Value(int32(copiesCount))
		item.InstanceName = types.StringValue(instance.InstanceName)
		item.CreateTime = types.StringValue(instance.CreateTime)
		if instance.ElasticIpBind == 1 {
			item.EipAddress = types.StringValue(instance.ElasticIp)
		} else {
			item.EipAddress = types.StringValue("")
		}
		item.ID = types.StringValue(instance.ProdInstId)
		config.Instances = append(config.Instances, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRedisInstances) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
