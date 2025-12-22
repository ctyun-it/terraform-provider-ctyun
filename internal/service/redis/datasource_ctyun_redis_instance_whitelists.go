package redis

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgdcs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunRedisInstanceWhitelists{}
	_ datasource.DataSourceWithConfigure = &ctyunRedisInstanceWhitelists{}
)

type ctyunRedisInstanceWhitelists struct {
	meta *common.CtyunMetadata
}

func NewCtyunRedisInstanceWhitelists() datasource.DataSource {
	return &ctyunRedisInstanceWhitelists{}
}

func (c *ctyunRedisInstanceWhitelists) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_instance_whitelists"
}

type CtyunRedisInstanceWhitelistModel struct {
	Name types.String `tfsdk:"name"`
	Ip   types.String `tfsdk:"ip"`
}

type CtyunRedisInstanceWhitelistsConfig struct {
	RegionId   types.String                       `tfsdk:"region_id"`
	InstanceId types.String                       `tfsdk:"instance_id"`
	Rows       []CtyunRedisInstanceWhitelistModel `tfsdk:"rows"`
}

func (c *ctyunRedisInstanceWhitelists) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029420/10398174`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "实例ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"rows": schema.ListNestedAttribute{
				Computed:    true,
				Description: "白名单分组列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "白名单分组名",
						},
						"ip": schema.StringAttribute{
							Computed:    true,
							Description: "白名单集合",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRedisInstanceWhitelists) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRedisInstanceWhitelistsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := c.meta.GetExtraIfEmpty(config.RegionId.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	config.RegionId = types.StringValue(regionId)

	instanceId := config.InstanceId.ValueString()
	if instanceId == "" {
		err = fmt.Errorf("instanceId不能为空")
		return
	}

	// 组装请求体
	params := &ctgdcs2.Dcs2DescribeSecurityIpsRequest{
		RegionId:   regionId,
		ProdInstId: instanceId,
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeSecurityIpsApi.Do(ctx, c.meta.SdkCredential, params)
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
	config.Rows = []CtyunRedisInstanceWhitelistModel{}
	for _, whitelist := range resp.ReturnObj.Rows {
		item := CtyunRedisInstanceWhitelistModel{
			Name: types.StringValue(whitelist.Group),
			Ip:   types.StringValue(whitelist.Ip),
		}
		config.Rows = append(config.Rows, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRedisInstanceWhitelists) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
