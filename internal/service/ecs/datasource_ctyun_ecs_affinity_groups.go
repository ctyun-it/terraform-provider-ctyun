package ecs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctecs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunEcsAffinityGroups{}
	_ datasource.DataSourceWithConfigure = &ctyunEcsAffinityGroups{}
)

type ctyunEcsAffinityGroups struct {
	meta *common.CtyunMetadata
}

func NewCtyunEcsAffinityGroups() datasource.DataSource {
	return &ctyunEcsAffinityGroups{}
}

func (c *ctyunEcsAffinityGroups) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ecs_affinity_groups"
}

type CtyunEcsAffinityGroupsModel struct {
	AffinityGroupID     types.String `tfsdk:"affinity_group_id"`
	AffinityGroupName   types.String `tfsdk:"affinity_group_name"`
	AffinityGroupPolicy types.String `tfsdk:"affinity_group_policy"`
	CreatedTime         types.String `tfsdk:"created_time"`
	UpdatedTime         types.String `tfsdk:"updated_time"`
}

type CtyunEcsAffinityGroupsConfig struct {
	RegionID        types.String                  `tfsdk:"region_id"`
	AffinityGroupID types.String                  `tfsdk:"affinity_group_id"`
	PageNo          types.Int32                   `tfsdk:"page_no"`
	PageSize        types.Int32                   `tfsdk:"page_size"`
	Groups          []CtyunEcsAffinityGroupsModel `tfsdk:"groups"`
}

func (c *ctyunEcsAffinityGroups) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026730/10597687`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"affinity_group_id": schema.StringAttribute{
				Optional:    true,
				Description: "云主机组ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "页码，取值范围：正整数（≥1），注：默认值为1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页记录数目，取值范围：[1,50]，注：默认值为10",
			},
			"groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"affinity_group_id": schema.StringAttribute{
							Computed:    true,
							Description: "云主机组ID",
						},
						"affinity_group_name": schema.StringAttribute{
							Computed:    true,
							Description: "云主机组名称",
						},
						"affinity_group_policy": schema.StringAttribute{
							Computed:    true,
							Description: "云主机组策略",
						},
						"created_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"updated_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间",
						},
					},
				}},
		},
	}
}

func (c *ctyunEcsAffinityGroups) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunEcsAffinityGroupsConfig
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
	config.Groups = []CtyunEcsAffinityGroupsModel{}
	// 组装请求体
	params := &ctecs2.CtecsListAffinityGroupV41Request{
		RegionID:        regionId,
		AffinityGroupID: config.AffinityGroupID.ValueString(),
		PageNo:          config.PageNo.ValueInt32(),
		PageSize:        config.PageSize.ValueInt32(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsListAffinityGroupV41Api.Do(ctx, c.meta.SdkCredential, params)
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
	for _, group := range resp.ReturnObj.Results {
		item := CtyunEcsAffinityGroupsModel{
			AffinityGroupID:   types.StringValue(group.AffinityGroupID),
			AffinityGroupName: types.StringValue(group.AffinityGroupName),
			CreatedTime:       types.StringValue(group.CreatedTime),
			UpdatedTime:       types.StringValue(group.UpdatedTime),
		}
		if group.AffinityGroupPolicy != nil {
			item.AffinityGroupPolicy = types.StringValue(group.AffinityGroupPolicy.PolicyTypeName)
		}
		config.Groups = append(config.Groups, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunEcsAffinityGroups) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
