package ccse

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ccse2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ccse"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunCcseImages{}
	_ datasource.DataSourceWithConfigure = &ctyunCcseImages{}
)

type ctyunCcseImages struct {
	meta *common.CtyunMetadata
}

func NewCtyunCcseImages() datasource.DataSource {
	return &ctyunCcseImages{}
}

func (c *ctyunCcseImages) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ccse_images"
}

type CtyunCcseImagesConfig struct {
	RegionID     types.String            `tfsdk:"region_id"`
	AzName       types.String            `tfsdk:"az_name"`
	ProjectID    types.String            `tfsdk:"project_id"`
	FlavorName   types.String            `tfsdk:"flavor_name"`
	InstanceType types.String            `tfsdk:"instance_type"`
	Images       []*CtyunCcseImagesModel `tfsdk:"images"`
}

type CtyunCcseImagesModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	OsType       types.String `tfsdk:"os_type"`
	OsDistro     types.String `tfsdk:"os_distro"`
	OsVersion    types.String `tfsdk:"os_version"`
	Visibility   types.String `tfsdk:"visibility"`
	Architecture types.String `tfsdk:"architecture"`
}

func (c *ctyunCcseImages) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10083472/10656137`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"az_name": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "可用区",
			},
			"project_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "企业项目ID",
			},
			"flavor_name": schema.StringAttribute{
				Required:    true,
				Description: "规格名称",
			},
			"instance_type": schema.StringAttribute{
				Required:    true,
				Description: "查询镜像类型，支持ecs（云主机）、ebm（裸金属）",
				Validators: []validator.String{
					stringvalidator.OneOf(business.CcseSlaveInstanceTypeEcs, business.CcseSlaveInstanceTypeEbm),
				},
			},
			"images": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "镜像唯一标识符",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "镜像名称",
						},
						"os_type": schema.StringAttribute{
							Computed:    true,
							Description: "操作系统类型",
						},
						"os_distro": schema.StringAttribute{
							Computed:    true,
							Description: "操作系统发行版",
						},
						"os_version": schema.StringAttribute{
							Computed:    true,
							Description: "操作系统版本号",
						},
						"visibility": schema.StringAttribute{
							Computed:    true,
							Description: "镜像可见性",
						},
						"architecture": schema.StringAttribute{
							Computed:    true,
							Description: "支持的CPU架构",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunCcseImages) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunCcseImagesConfig
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

	azName := c.meta.GetExtraIfEmpty(config.AzName.ValueString(), common.ExtraAzName)
	if azName == "" && config.InstanceType.ValueString() == business.CcseSlaveInstanceTypeEbm {
		err = fmt.Errorf("查询物理机镜像时azName不能为空")
		return
	}
	config.AzName = types.StringValue(azName)

	projectID := c.meta.GetExtraIfEmpty(config.ProjectID.ValueString(), common.ExtraProjectId)
	config.ProjectID = types.StringValue(projectID)

	err = c.getImagesAndMerge(ctx, &config)
	if err != nil {
		return
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunCcseImages) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunCcseImages) getImagesAndMerge(ctx context.Context, config *CtyunCcseImagesConfig) (err error) {
	// 组装请求体
	params := &ccse2.CcseGetPublicImageListRequest{
		RegionId:   config.RegionID.ValueString(),
		FlavorName: config.FlavorName.ValueString(),
		VmType:     config.InstanceType.ValueString(),
		AzName:     config.AzName.ValueString(),
		ProjectId:  config.ProjectID.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCcseApis.CcseGetPublicImageListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	}
	config.Images = []*CtyunCcseImagesModel{}
	for _, image := range resp.ReturnObj {
		item := &CtyunCcseImagesModel{
			ID:           types.StringValue(image.ImageID),
			Name:         types.StringValue(image.ImageName),
			OsType:       types.StringValue(image.OsType),
			OsDistro:     types.StringValue(image.OsDistro),
			OsVersion:    types.StringValue(image.OsVersion),
			Visibility:   types.StringValue(image.Visibility),
			Architecture: types.StringValue(image.Architecture),
		}
		config.Images = append(config.Images, item)
	}
	return
}
