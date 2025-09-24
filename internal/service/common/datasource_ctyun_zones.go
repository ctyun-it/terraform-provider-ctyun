package common

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctecs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ctyunZones struct {
	meta *common.CtyunMetadata
}

func NewCtyunZones() datasource.DataSource {
	return &ctyunZones{}
}

func (c *ctyunZones) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zones"
}

func (c *ctyunZones) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `**资源池可用区**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"zones": schema.ListAttribute{
				Computed:    true,
				Description: "可用区编码，即az_name",
				ElementType: types.StringType,
			},
		},
	}
}

func (c *ctyunZones) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunZonesConfig
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
	params := &ctecs2.CtecsQueryZonesInRegionV41Request{
		RegionID: regionId,
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsQueryZonesInRegionV41Api.Do(ctx, c.meta.SdkCredential, params)
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
	var azList []string
	for _, az := range resp.ReturnObj.ZoneList {
		azList = append(azList, az.Name)
	}
	config.Zones, _ = types.ListValueFrom(ctx, types.StringType, azList)
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunZones) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunZonesConfig struct {
	RegionID types.String `tfsdk:"region_id"`
	Zones    types.List   `tfsdk:"zones"`
}
