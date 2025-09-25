package ebm

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebm"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunEbmDeviceTypes{}
	_ datasource.DataSourceWithConfigure = &ctyunEbmDeviceTypes{}
)

type ctyunEbmDeviceRaids struct {
	meta *common.CtyunMetadata
}

func NewCtyunEbmDeviceRaids() datasource.DataSource {
	return &ctyunEbmDeviceRaids{}
}

func (c *ctyunEbmDeviceRaids) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ebm_device_raids"
}

type CtyunEbmDeviceRaidsModel struct {
	DeviceType    types.String `tfsdk:"device_type"`
	VolumeType    types.String `tfsdk:"volume_type"`
	Uuid          types.String `tfsdk:"uuid"`
	NameEn        types.String `tfsdk:"name_en"`
	NameZh        types.String `tfsdk:"name_zh"`
	VolumeDetail  types.String `tfsdk:"volume_detail"`
	DescriptionEn types.String `tfsdk:"description_en"`
	DescriptionZh types.String `tfsdk:"description_zh"`
}

type CtyunEbmDeviceRaidsConfig struct {
	RegionID   types.String               `tfsdk:"region_id"`
	AzName     types.String               `tfsdk:"az_name"`
	DeviceType types.String               `tfsdk:"device_type"`
	VolumeType types.String               `tfsdk:"volume_type"`
	Raids      []CtyunEbmDeviceRaidsModel `tfsdk:"raids"`
}

func (c *ctyunEbmDeviceRaids) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027724/10166084**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称",
			},
			"device_type": schema.StringAttribute{
				Required:    true,
				Description: "套餐类型",
			},
			"volume_type": schema.StringAttribute{
				Required:    true,
				Description: "磁盘类型（system/data）",
				Validators: []validator.String{
					stringvalidator.OneOf("system", "data"),
				},
			},
			"raids": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_type": schema.StringAttribute{
							Computed:    true,
							Description: "套餐类型",
						},
						"volume_type": schema.StringAttribute{
							Computed:    true,
							Description: "磁盘类型",
						},
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "raid_uuid",
						},
						"name_en": schema.StringAttribute{
							Computed:    true,
							Description: "raid英文名称",
						},
						"name_zh": schema.StringAttribute{
							Computed:    true,
							Description: "raid中文名称",
						},
						"volume_detail": schema.StringAttribute{
							Computed:    true,
							Description: "对应套餐磁盘描述",
						},
						"description_en": schema.StringAttribute{
							Computed:    true,
							Description: "raid英文介绍",
						},
						"description_zh": schema.StringAttribute{
							Computed:    true,
							Description: "raid中文介绍",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunEbmDeviceRaids) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunEbmDeviceRaidsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	azName := c.meta.GetExtraIfEmpty(config.AzName.ValueString(), common.ExtraAzName)
	if azName == "" {
		err = fmt.Errorf("azName不能为空")
		return
	}
	// 组装请求体
	params := &ctebm.EbmRaidTypeListRequest{
		RegionID:   regionId,
		AzName:     azName,
		DeviceType: config.DeviceType.ValueString(),
		VolumeType: config.VolumeType.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.CtEbmApis.EbmRaidTypeListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回值
	raids := []CtyunEbmDeviceRaidsModel{}
	for _, raid := range resp.ReturnObj.Results {
		item := CtyunEbmDeviceRaidsModel{
			DeviceType:    utils.SecStringValue(raid.DeviceType),
			VolumeType:    utils.SecStringValue(raid.VolumeType),
			Uuid:          utils.SecStringValue(raid.Uuid),
			NameEn:        utils.SecStringValue(raid.NameEn),
			NameZh:        utils.SecStringValue(raid.NameZh),
			VolumeDetail:  utils.SecStringValue(raid.VolumeDetail),
			DescriptionEn: utils.SecStringValue(raid.DescriptionEn),
			DescriptionZh: utils.SecStringValue(raid.DescriptionZh),
		}
		raids = append(raids, item)
	}
	config.RegionID = types.StringValue(regionId)
	config.AzName = types.StringValue(azName)
	config.Raids = raids
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunEbmDeviceRaids) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
