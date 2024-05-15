package datasource

import (
	"context"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctecs"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/common"
)

type ctyunRegions struct {
	meta *common.CtyunMetadata
}

func NewCtyunRegions() datasource.DataSource {
	return &ctyunRegions{}
}

func (c *ctyunRegions) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_regions"
}

func (c *ctyunRegions) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "资源池名称",
			},
			"regions": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "资源池id",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "资源池名称",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "资源池类型",
						},
						"zones": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "可用区",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRegions) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CtyunRegionsConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	response, err := c.meta.Apis.CtEcsApis.RegionListApi.Do(ctx, c.meta.Credential, &ctecs.RegionListRequest{
		RegionName: config.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	var regions []CtyunRegionsRegionsConfig
	for _, r := range response.RegionList {
		rawZones := []string{}
		rawZones = append(rawZones, r.ZoneList...)
		zones, _ := types.SetValueFrom(ctx, types.StringType, rawZones)
		regions = append(regions, CtyunRegionsRegionsConfig{
			Id:    types.StringValue(r.RegionId),
			Name:  types.StringValue(r.RegionName),
			Type:  types.StringValue(r.RegionType),
			Zones: zones,
		})
	}
	config.Regions = regions
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (c *ctyunRegions) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunRegionsRegionsConfig struct {
	Id    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Type  types.String `tfsdk:"type"`
	Zones types.Set    `tfsdk:"zones"`
}

type CtyunRegionsConfig struct {
	Name    types.String                `tfsdk:"name"`
	Regions []CtyunRegionsRegionsConfig `tfsdk:"regions"`
}
