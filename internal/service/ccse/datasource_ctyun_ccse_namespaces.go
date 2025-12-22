package ccse

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ccse2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ccse"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunCcseNamespaces{}
	_ datasource.DataSourceWithConfigure = &ctyunCcseNamespaces{}
)

type ctyunCcseNamespaces struct {
	meta *common.CtyunMetadata
}

func NewCtyunCcseNamespaces() datasource.DataSource {
	return &ctyunCcseNamespaces{}
}

func (c *ctyunCcseNamespaces) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ccse_namespaces"
}

type CtyunCcseNamespacesConfig struct {
	RegionID   types.String `tfsdk:"region_id"`
	ClusterID  types.String `tfsdk:"cluster_id"`
	Label      types.String `tfsdk:"label"`
	Field      types.String `tfsdk:"field"`
	ValuesYaml types.String `tfsdk:"values_yaml"`
}

func (c *ctyunCcseNamespaces) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10083472/10656137`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "集群ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(32, 32),
				},
			},
			"label": schema.StringAttribute{
				Optional:    true,
				Description: "Kubernetes labelSelector，可通过label过滤资源；label之间通过“,”分隔，特殊符号要转义为url编码",
			},
			"field": schema.StringAttribute{
				Optional:    true,
				Description: "Kubernetes fieldSelector，可通过field过滤资源；field之间通过“,”分隔，特殊符号要转义为url编码",
			},
			"values_yaml": schema.StringAttribute{
				Computed:    true,
				Description: "命名空间配置",
			},
		},
	}
}

func (c *ctyunCcseNamespaces) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunCcseNamespacesConfig
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
	err = c.getYamlAndMerge(ctx, &config)
	if err != nil {
		return
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunCcseNamespaces) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunCcseNamespaces) getYamlAndMerge(ctx context.Context, config *CtyunCcseNamespacesConfig) (err error) {
	// 组装请求体
	params := &ccse2.CcseListNamespaceV2P2Request{
		ClusterName:   config.ClusterID.ValueString(),
		RegionId:      config.RegionID.ValueString(),
		LabelSelector: config.Label.ValueString(),
		FieldSelector: config.Field.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCcseApis.CcseListNamespaceV2P2Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	}
	config.ValuesYaml = types.StringValue(resp.ReturnObj)
	return
}
