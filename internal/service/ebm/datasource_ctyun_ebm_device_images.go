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

type ctyunEbmDeviceImages struct {
	meta *common.CtyunMetadata
}

func NewCtyunEbmDeviceImages() datasource.DataSource {
	return &ctyunEbmDeviceImages{}
}

func (c *ctyunEbmDeviceImages) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebm_device_images"
}

type CtyunEbmDeviceImagesModel struct {
	NameZh     types.String                `tfsdk:"name_zh"`
	Format     types.String                `tfsdk:"format"`
	ImageType  types.String                `tfsdk:"image_type"`
	IsShared   types.Bool                  `tfsdk:"is_shared"`
	Version    types.String                `tfsdk:"version"`
	ImageUUID  types.String                `tfsdk:"image_uuid"`
	NameEn     types.String                `tfsdk:"name_en"`
	LayoutType types.String                `tfsdk:"layout_type"`
	Os         CtyunEbmDeviceImagesOsModel `tfsdk:"os"`
}

type CtyunEbmDeviceImagesOsModel struct {
	Uuid         types.String `tfsdk:"uuid"`
	SuperUser    types.String `tfsdk:"super_user"`
	Platform     types.String `tfsdk:"platform"`
	Version      types.String `tfsdk:"version"`
	Architecture types.String `tfsdk:"architecture"`
	NameEn       types.String `tfsdk:"name_en"`
	Bits         types.Int32  `tfsdk:"bits"`
	OsType       types.String `tfsdk:"os_type"`
	NameZh       types.String `tfsdk:"name_zh"`
}

type CtyunEbmDeviceImagesConfig struct {
	RegionID   types.String                `tfsdk:"region_id"`
	AzName     types.String                `tfsdk:"az_name"`
	DeviceType types.String                `tfsdk:"device_type"`
	ImageType  types.String                `tfsdk:"image_type"`
	ImageUUID  types.String                `tfsdk:"image_uuid"`
	OsType     types.String                `tfsdk:"os_type"`
	Images     []CtyunEbmDeviceImagesModel `tfsdk:"images"`
}

func (c *ctyunEbmDeviceImages) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027724/10173844**`,
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
			"image_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "镜像类型，可选择：private(私有镜像)、standard(标准镜像)、shared(共享镜像)",
				Validators: []validator.String{
					stringvalidator.OneOf("private", "standard", "shared"),
				},
			},
			"image_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "物理机镜像UUID",
			},
			"os_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "操作系统类型，取值范围：linux，windows",
			},
			"images": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name_zh": schema.StringAttribute{
							Computed:    true,
							Description: "中文名称",
						},
						"format": schema.StringAttribute{
							Computed:    true,
							Description: "规格;包括squashfs,qcow2",
						},
						"image_type": schema.StringAttribute{
							Computed:    true,
							Description: "镜像类型;包括standard,private,shared，默认为standard",
						},
						"is_shared": schema.BoolAttribute{
							Computed:    true,
							Description: "镜像是否共享",
						},
						"version": schema.StringAttribute{
							Computed:    true,
							Description: "版本",
						},
						"image_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "镜像uuid",
						},
						"name_en": schema.StringAttribute{
							Computed:    true,
							Description: "英文名称",
						},
						"layout_type": schema.StringAttribute{
							Computed:    true,
							Description: "布局类型;包括lvm,direct",
						},
						"os": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"uuid": schema.StringAttribute{
									Computed:    true,
									Description: "操作系统uuid",
								},
								"super_user": schema.StringAttribute{
									Computed:    true,
									Description: "超级管理员",
								},
								"platform": schema.StringAttribute{
									Computed:    true,
									Description: "平台",
								},
								"version": schema.StringAttribute{
									Computed:    true,
									Description: "版本",
								},
								"architecture": schema.StringAttribute{
									Computed:    true,
									Description: "支持的机器类型",
								},
								"name_en": schema.StringAttribute{
									Computed:    true,
									Description: "英文名称",
								},
								"bits": schema.Int64Attribute{
									Computed:    true,
									Description: "比特数",
								},
								"os_type": schema.StringAttribute{
									Computed:    true,
									Description: "操作系统类别",
								},
								"name_zh": schema.StringAttribute{
									Computed:    true,
									Description: "中文名称",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunEbmDeviceImages) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunEbmDeviceImagesConfig
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
	params := &ctebm.EbmImageListRequest{
		RegionID:   regionId,
		AzName:     azName,
		DeviceType: config.DeviceType.ValueString(),
	}
	if !config.ImageType.IsNull() {
		imageType := config.ImageType.ValueString()
		params.ImageType = &imageType
	}
	if !config.ImageUUID.IsNull() {
		imageUUID := config.ImageUUID.ValueString()
		params.ImageUUID = &imageUUID
	}
	if !config.OsType.IsNull() {
		osType := config.OsType.ValueString()
		params.OsType = &osType
	}
	// 调用API
	resp, err := c.meta.Apis.CtEbmApis.EbmImageListApi.Do(ctx, c.meta.SdkCredential, params)
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
	images := []CtyunEbmDeviceImagesModel{}
	for _, image := range resp.ReturnObj.Results {
		item := CtyunEbmDeviceImagesModel{
			Format:     utils.SecStringValue(image.Format),
			ImageType:  utils.SecStringValue(image.ImageType),
			Version:    utils.SecStringValue(image.Version),
			ImageUUID:  utils.SecStringValue(image.ImageUUID),
			LayoutType: utils.SecStringValue(image.LayoutType),
			NameEn:     utils.SecStringValue(image.NameEn),
			NameZh:     utils.SecStringValue(image.NameZh),
			IsShared:   utils.SecBoolValue(image.IsShared),
			Os:         CtyunEbmDeviceImagesOsModel{},
		}
		if image.Os != nil {
			item.Os.OsType = utils.SecStringValue(image.Os.OsType)
			item.Os.Version = utils.SecStringValue(image.Os.Version)
			item.Os.Architecture = utils.SecStringValue(image.Os.Architecture)
			item.Os.Uuid = utils.SecStringValue(image.Os.Uuid)
			item.Os.NameEn = utils.SecStringValue(image.Os.NameEn)
			item.Os.NameZh = utils.SecStringValue(image.Os.NameZh)
			item.Os.Bits = types.Int32Value(image.Os.Bits)
			item.Os.Platform = utils.SecStringValue(image.Os.Platform)
			item.Os.SuperUser = utils.SecStringValue(image.Os.SuperUser)
		}
		images = append(images, item)
	}
	config.RegionID = types.StringValue(regionId)
	config.AzName = types.StringValue(azName)
	config.Images = images
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunEbmDeviceImages) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
