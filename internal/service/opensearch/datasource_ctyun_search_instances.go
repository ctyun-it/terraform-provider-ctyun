package opensearch

import (
	"context"
	"fmt"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/opensearch"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource = &CtyunSearchInstances{}
)

func NewCtyunSearchInstances() datasource.DataSource {
	return &CtyunSearchInstances{}
}

type CtyunSearchInstances struct {
	meta *common.CtyunMetadata
}

type CtyunSearchInstancesConfig struct {
	ID        types.String              `tfsdk:"id"`
	RegionID  types.String              `tfsdk:"region_id"`
	PageNo    types.Int64               `tfsdk:"page_no"`
	PageSize  types.Int64               `tfsdk:"page_size"`
	Instances []CtyunSearchInstanceInfo `tfsdk:"instances"`
	Total     types.Int64               `tfsdk:"total"`
}

type CtyunSearchInstanceInfo struct {
	ID          types.String `tfsdk:"id"`
	ClusterName types.String `tfsdk:"name"`
	Status      types.String `tfsdk:"status"`
	CreateTime  types.String `tfsdk:"create_time"`
}

func (c *CtyunSearchInstances) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_search_instances"
}

func (c *CtyunSearchInstances) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("查询OpenSearch实例列表", "天翼云OpenSearch服务", "https://www.ctyun.cn/document/10026730/10040008"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "数据源唯一标识",
			},
			"region_id": schema.StringAttribute{
				Required:    true,
				Description: "资源池ID",
			},
			"page_no": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "页码，默认为1",
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "每页数量，默认为10",
			},
			"total": schema.Int64Attribute{
				Computed:    true,
				Description: "总实例数",
			},
			"instances": schema.ListNestedAttribute{
				Computed:    true,
				Description: "实例列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "实例ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "实例名称",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "实例状态",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunSearchInstances) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunSearchInstances) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSearchInstancesConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 设置默认值
	pageIndex := int(state.PageNo.ValueInt64())
	if pageIndex <= 0 {
		pageIndex = 1
	}

	pageSize := int(state.PageSize.ValueInt64())
	if pageSize <= 0 {
		pageSize = 10
	}

	tflog.Info(ctx, "查询 OpenSearch 实例列表", map[string]interface{}{
		"region_id": state.RegionID.ValueString(),
		"page_no":   pageIndex,
		"page_size": pageSize,
	})

	// 构造查询请求
	listReq := &opensearch.ListInstancesRequest{
		RegionID:    state.RegionID.ValueString(),
		PageIndex:   pageIndex,
		PageSize:    pageSize,
		ClusterType: 1, // 默认为 OpenSearch
	}

	response, err := c.meta.Apis.SdkOpensearchApis.OpensearchListInstancesApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		return
	} else if response.StatusCode != 200 {
		err = fmt.Errorf("API return error. StatusCode: %d, Message: %s", response.StatusCode, response.Message)
		return
	} else if response.ReturnObj.Records == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析结果
	var instanceList []CtyunSearchInstanceInfo
	for _, record := range response.ReturnObj.Records {
		instanceInfo := CtyunSearchInstanceInfo{
			ID:          types.StringValue(record.ClusterID),
			ClusterName: types.StringValue(record.ClusterName),
			Status:      types.StringValue(record.ClusterState),
			CreateTime:  types.StringValue(fmt.Sprintf("%d", record.CreateTime)),
		}
		instanceList = append(instanceList, instanceInfo)
	}

	state.Instances = instanceList
	state.Total = types.Int64Value(int64(response.ReturnObj.Total))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
