package crs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/crs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunCrsOpensourceImages{}
	_ datasource.DataSourceWithConfigure = &ctyunCrsOpensourceImages{}
)

type ctyunCrsOpensourceImages struct {
	meta *common.CtyunMetadata
}

func NewCtyunCrsOpensourceImages() datasource.DataSource {
	return &ctyunCrsOpensourceImages{}
}

func (c *ctyunCrsOpensourceImages) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_crs_opensource_images"
}

type CtyunCrsOpensourceImagesModel struct {
	NamespaceName    types.String `tfsdk:"namespace_name"`
	RepositoryName   types.String `tfsdk:"repository_name"`
	RepositoryID     types.String `tfsdk:"repository_id"`
	ImageUrl         types.String `tfsdk:"image_url"`
	ImageUrlInternal types.String `tfsdk:"image_url_internal"`
	Category         []string     `tfsdk:"category"`
	Architecture     []string     `tfsdk:"architecture"`
	Os               []string     `tfsdk:"os"`
}

type CtyunCrsOpensourceImagesConfig struct {
	RegionID       types.String                    `tfsdk:"region_id"`
	RepositoryName types.String                    `tfsdk:"repository_name"`
	Category       types.String                    `tfsdk:"category"`
	Architecture   types.String                    `tfsdk:"architecture"`
	PageNo         types.Int32                     `tfsdk:"page_no"`
	PageSize       types.Int32                     `tfsdk:"page_size"`
	Images         []CtyunCrsOpensourceImagesModel `tfsdk:"images"`
}

func (c *ctyunCrsOpensourceImages) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10007018/10007025`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"repository_name": schema.StringAttribute{
				Optional:    true,
				Description: "镜像仓库名称",
			},
			"category": schema.StringAttribute{
				Optional:    true,
				Description: "根据分类标签筛选（base：基础镜像，os：操作系统，ai：AI，middleware：中间件，storage：存储，network：网络，other：其他），支持根据多个标签筛选（多个标签使用,分隔，筛选结果为满足任一标签的镜像），如果为空表示包含所有分类",
			},
			"architecture": schema.StringAttribute{
				Optional:    true,
				Description: "根据架构筛选（arm64，amd64），支持根据多个标签筛选（多个标签使用,分隔，筛选结果为满足任一标签的镜像），如果为空表示包含所有架构",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "页码，取值范围：正整数（≥1），注：默认值为1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "每页记录数目，取值范围：[1,50]，注：默认值为10",
			},
			"images": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"namespace_name": schema.StringAttribute{
							Description: "命名空间名称",
							Computed:    true,
						},
						"repository_name": schema.StringAttribute{
							Description: "镜像仓库名称",
							Computed:    true,
						},
						"repository_id": schema.StringAttribute{
							Description: "镜像仓库ID",
							Computed:    true,
						},
						"image_url": schema.StringAttribute{
							Description: "镜像公开URL",
							Computed:    true,
						},
						"image_url_internal": schema.StringAttribute{
							Description: "镜像内部URL",
							Computed:    true,
						},
						"category": schema.ListAttribute{
							Description: "支持标签列表",
							ElementType: types.StringType,
							Computed:    true,
						},
						"architecture": schema.ListAttribute{
							Description: "支持架构列表",
							ElementType: types.StringType,
							Computed:    true,
						},
						"os": schema.ListAttribute{
							Description: "支持系统列表",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		}}
}

func (c *ctyunCrsOpensourceImages) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunCrsOpensourceImagesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	// 组装请求体
	params := &crs.CrsListOpenSourceRepositoryV2Request{
		RegionId:       regionId,
		RepositoryName: config.RepositoryName.ValueStringPointer(),
		Category:       config.Category.ValueStringPointer(),
		Architecture:   config.Architecture.ValueStringPointer(),
		PageNum:        config.PageNo.ValueInt32Pointer(),
		PageSize:       config.PageSize.ValueInt32Pointer(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCrsApis.CrsListOpenSourceRepositoryV2Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回值
	images := []CtyunCrsOpensourceImagesModel{}
	for _, image := range resp.ReturnObj.Records {
		item := CtyunCrsOpensourceImagesModel{
			NamespaceName:    utils.SecStringValue(image.NamespaceName),
			RepositoryName:   utils.SecStringValue(image.RepositoryName),
			RepositoryID:     utils.SecStringValue(image.RepositoryId),
			ImageUrl:         utils.SecStringValue(image.ImageUrl),
			ImageUrlInternal: utils.SecStringValue(image.ImageUrlInternal),
			Category:         utils.StrPointerArrayToStrArray(image.Category),
			Architecture:     utils.StrPointerArrayToStrArray(image.Architecture),
			Os:               utils.StrPointerArrayToStrArray(image.Os),
		}

		images = append(images, item)
	}
	config.RegionID = types.StringValue(regionId)
	config.Images = images
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunCrsOpensourceImages) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
