package datasource

import (
	"context"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctimage"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
)

func NewCtyunImages() datasource.DataSource {
	return &ctyunImages{}
}

type ctyunImages struct {
	meta *common.CtyunMetadata
}

func (c *ctyunImages) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_images"
}

func (c *ctyunImages) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027726**`,
		Attributes: map[string]schema.Attribute{
			"visibility": schema.StringAttribute{
				Required:    true,
				Description: "镜像可见类型：private：私有镜像，public：公共镜像，shared：共享镜像，safe：安全产品镜像，app：甄选应用镜像",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ImageVisibilities...),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "镜像名称，模糊查询",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区id，如果不填则默认使用provider ctyun中的az_name或环境变量中的CTYUN_AZ_NAME",
			},
			"page_size": schema.Int64Attribute{
				Required:    true,
				Description: "每页显示数量，取值范围1-50",
				Validators: []validator.Int64{
					int64validator.Between(1, 50),
				},
			},
			"page_no": schema.Int64Attribute{
				Required:    true,
				Description: "当前页码",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"images": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "镜像id",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "镜像名称",
						},
						"os_type": schema.StringAttribute{
							Computed:    true,
							Description: "系统类型",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunImages) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CtyunImagesConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	regionId := c.meta.GetExtraIfEmpty(config.RegionId.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		msg := "regionId不能为空"
		resp.Diagnostics.AddError(msg, msg)
		return
	}
	azName := c.meta.GetExtraIfEmpty(config.AzName.ValueString(), common.ExtraAzName)

	visibility, err := business.ImageVisibilityMap.FromOriginalScene(config.Visibility.ValueString(), business.ImageVisibilityMapScene1)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	imageListResponse, err := c.meta.Apis.CtImageApis.ImageListApi.Do(ctx, c.meta.Credential, &ctimage.ImageListRequest{
		RegionId:     regionId,
		AzName:       azName,
		Visibility:   visibility.(int),
		QueryContent: config.Name.ValueString(),
		PageNo:       int(config.PageNo.ValueInt64()),
		PageSize:     int(config.PageSize.ValueInt64()),
	})
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	var images []CtyunImageImagesConfig
	for _, ig := range imageListResponse.Images {
		images = append(images, CtyunImageImagesConfig{
			Name:        types.StringValue(ig.ImageName),
			Id:          types.StringValue(ig.ImageId),
			OsType:      types.StringValue(ig.OsType),
			Description: types.StringValue(ig.Description),
		})
	}
	config.Images = images
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (c *ctyunImages) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunImageImagesConfig struct {
	Id          types.String `tfsdk:"id"`          // 镜像id
	Name        types.String `tfsdk:"name"`        // 镜像名称
	OsType      types.String `tfsdk:"os_type"`     // 系统类型
	Description types.String `tfsdk:"description"` // 描述
}

type CtyunImagesConfig struct {
	Visibility types.String             `tfsdk:"visibility"`
	Name       types.String             `tfsdk:"name"`
	RegionId   types.String             `tfsdk:"region_id"`
	AzName     types.String             `tfsdk:"az_name"`
	Images     []CtyunImageImagesConfig `tfsdk:"images"`
	PageSize   types.Int64              `tfsdk:"page_size"`
	PageNo     types.Int64              `tfsdk:"page_no"`
}
