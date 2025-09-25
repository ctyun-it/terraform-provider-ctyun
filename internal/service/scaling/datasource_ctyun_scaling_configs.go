package scaling

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/scaling"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunScalingConfigs{}
	_ datasource.DataSourceWithConfigure = &CtyunScalingConfigs{}
)

type CtyunScalingConfigs struct {
	meta *common.CtyunMetadata
}

func NewCtyunScalingConfigs() datasource.DataSource {
	return &CtyunScalingConfigs{}
}

func (c *CtyunScalingConfigs) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunScalingConfigs) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_scaling_configs"
}

func (c *CtyunScalingConfigs) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027725/10241446`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id",
			},
			"id": schema.Int64Attribute{
				Optional:    true,
				Description: "伸缩配置id",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页包含的元素个数范围(1-50)，默认值为10",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的分页页码，默认值为1",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"scaling_config_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int32Attribute{
							Computed:    true,
							Description: "伸缩配置ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "伸缩配置名称",
						},
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源池ID",
						},
						"visibility": schema.StringAttribute{
							Computed:    true,
							Description: "镜像类型: public-公有镜像, private-私有镜像",
						},
						"image_name": schema.StringAttribute{
							Computed:    true,
							Description: "镜像名称",
						},
						"image_id": schema.StringAttribute{
							Computed:    true,
							Description: "镜像ID",
						},
						"cpu": schema.Int32Attribute{
							Computed:    true,
							Description: "CPU核数",
						},
						"memory": schema.Int32Attribute{
							Computed:    true,
							Description: "内存大小(GB)",
						},
						"flavor_name": schema.StringAttribute{
							Computed:    true,
							Description: "规格名称",
						},
						"os_type": schema.StringAttribute{
							Computed:    true,
							Description: "镜像系统类型: Linux/Windows",
						},
						"volumes": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"volume_type": schema.StringAttribute{
										Computed:    true,
										Description: "磁盘类型",
									},
									"volume_size": schema.Int32Attribute{
										Computed:    true,
										Description: "磁盘大小(GB)",
									},
									"disk_mode": schema.StringAttribute{
										Computed:    true,
										Description: "磁盘模式",
									},
									"flag": schema.StringAttribute{
										Computed:    true,
										Description: "磁盘类型: OS-系统盘, DATA-数据盘",
									},
								},
							},
							Description: "磁盘类型和大小列表",
						},
						"use_floatings": schema.StringAttribute{
							Computed:    true,
							Description: "是否使用弹性IP",
						},
						"bandwidth": schema.Int32Attribute{
							Computed:    true,
							Description: "弹性IP带宽",
						},
						"login_mode": schema.StringAttribute{
							Computed:    true,
							Description: "登录方式",
						},
						"username": schema.StringAttribute{
							Computed:    true,
							Description: "用户名",
						},
						"password": schema.StringAttribute{
							Computed:    true,
							Sensitive:   true,
							Description: "密码",
						},
						"key_pair_id": schema.StringAttribute{
							Computed:    true,
							Description: "密钥对ID",
						},
						"tags": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Computed:    true,
										Description: "标签键",
									},
									"value": schema.StringAttribute{
										Computed:    true,
										Description: "标签值",
									},
								},
							},
							Description: "标签集",
						},
						"az_names": schema.StringAttribute{
							Computed:    true,
							Description: "可用区列表",
						},
						"monitor_service": schema.BoolAttribute{
							Computed:    true,
							Description: "是否开启详细监控",
						},
					},
				},
				Description: "伸缩配置列表",
			},
		},
	}
}

func (c *CtyunScalingConfigs) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunScalingConfigsModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	if regionId == "" {
		err = errors.New("region id 为空")
		return
	}
	params := &scaling.ScalingConfigListRequest{
		RegionID: regionId,
		PageNo:   1,
		PageSize: 10,
	}

	if !config.PageNo.IsNull() {
		params.PageNo = config.PageNo.ValueInt32()
	}
	if !config.PageSize.IsNull() {
		params.PageSize = config.PageSize.ValueInt32()
	}
	if !config.ID.IsNull() {
		params.ConfigID = config.ID.ValueInt64()
	}

	resp, err := c.meta.Apis.SdkScalingApis.ScalingConfigListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("获取id为 %d 的伸缩配置详情失败，接口返回nil。请稍后重试！", config.ID.ValueInt64())
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var diags diag.Diagnostics
	var scalingConfigList []CtyunScalingConfigInfo
	for _, configItem := range resp.ReturnObj {
		var scalingConfig CtyunScalingConfigInfo
		scalingConfig.ID = types.Int32Value(configItem.ConfigID)
		scalingConfig.Name = types.StringValue(configItem.Name)
		scalingConfig.RegionID = types.StringValue(configItem.RegionID)
		scalingConfig.Visibility = types.StringValue(business.ScalingVisibilityDictRev[configItem.Visibility])
		scalingConfig.ImageName = types.StringValue(configItem.ImageName)
		scalingConfig.ImageID = types.StringValue(configItem.ImageID)
		scalingConfig.Cpu = types.Int32Value(configItem.Cpu)
		scalingConfig.Memory = types.Int32Value(configItem.Memory)
		scalingConfig.FlavorName = types.StringValue(configItem.SpecName)
		scalingConfig.OsType = types.StringValue(business.ScalingOsTypeDictRev[configItem.OsType])
		scalingConfig.UseFloatings = types.StringValue(business.ScalingUseFloatingsDictRev[configItem.UseFloatings])
		scalingConfig.AzNames = types.StringValue(configItem.AzNames)
		if scalingConfig.UseFloatings.ValueString() == business.ScalingUseFloatingsAutoStr {
			scalingConfig.BandWidth = types.Int32Value(configItem.Bandwidth)
		} else {
			scalingConfig.BandWidth = types.Int32Value(0)
		}
		scalingConfig.LoginMode = types.StringValue(business.ScalingLoginModeDictRev[configItem.LoginMode])
		scalingConfig.Username = types.StringValue(configItem.Username)
		scalingConfig.MonitorService = types.BoolValue(*configItem.MonitorService)

		// 处理Volumes
		var volumeList []CtyunVolumesModel
		for _, volumeItem := range configItem.Volumes {
			var volume CtyunVolumesModel
			volume.VolumeSize = types.Int32Value(volumeItem.VolumeSize)
			volume.VolumeType = types.StringValue(volumeItem.VolumeType)
			volume.DiskMode = types.StringValue(volumeItem.DiskMode)
			volume.Flag = types.StringValue(business.ScalingVolumeFlagDictRev[volumeItem.Flag])
			volumeList = append(volumeList, volume)
		}
		scalingConfig.Volumes, diags = types.ListValueFrom(ctx, utils.StructToTFObjectTypes(CtyunVolumesModel{}), volumeList)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		// 处理tag
		var tagList []CtyunTagModel
		for _, tagItem := range configItem.Tags {
			var tag CtyunTagModel
			tag.Key = types.StringValue(tagItem.Key)
			tag.Value = types.StringValue(tagItem.Value)
			tagList = append(tagList, tag)
		}
		scalingConfig.Tags, diags = types.ListValueFrom(ctx, utils.StructToTFObjectTypes(CtyunTagModel{}), tagList)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}

		scalingConfigList = append(scalingConfigList, scalingConfig)
	}

	config.ScalingConfigList = scalingConfigList
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

}

type CtyunScalingConfigInfo struct {
	ID             types.Int32  `tfsdk:"id"`              // 伸缩配置ID
	Name           types.String `tfsdk:"name"`            // 伸缩配置名称
	RegionID       types.String `tfsdk:"region_id"`       // 资源池ID
	Visibility     types.String `tfsdk:"visibility"`      // 镜像类型。 取值范围：  1：公有镜像; 0：私有镜像
	ImageName      types.String `tfsdk:"image_name"`      // 镜像名称
	ImageID        types.String `tfsdk:"image_id"`        // 镜像ID
	Cpu            types.Int32  `tfsdk:"cpu"`             // CPU核数
	Memory         types.Int32  `tfsdk:"memory"`          // 内存，单位：G
	FlavorName     types.String `tfsdk:"flavor_name"`     // 规格名称
	OsType         types.String `tfsdk:"os_type"`         // 镜像系统类型。 取值范围： Linux ; Windows
	Volumes        types.List   `tfsdk:"volumes"`         // 磁盘类型和大小列表
	UseFloatings   types.String `tfsdk:"use_floatings"`   // 是否使用弹性IP
	BandWidth      types.Int32  `tfsdk:"bandwidth"`       // 弹性IP带宽
	LoginMode      types.String `tfsdk:"login_mode"`      // 登录方式
	Username       types.String `tfsdk:"username"`        // 用户名
	Tags           types.List   `tfsdk:"tags"`            // 标签集
	AzNames        types.String `tfsdk:"az_names"`        // 可用区列表
	MonitorService types.Bool   `tfsdk:"monitor_service"` // 是否开启详细监控
	KeyPairID      types.String `tfsdk:"key_pair_id"`     //
	Password       types.String `tfsdk:"password"`
}

type CtyunScalingConfigsModel struct {
	RegionID          types.String             `tfsdk:"region_id"` // 资源池ID
	ID                types.Int64              `tfsdk:"id"`        // 伸缩配置ID
	PageSize          types.Int32              `tfsdk:"page_size"` // 每页包含的元素个数范围(1-50)，默认值为10
	PageNo            types.Int32              `tfsdk:"page_no"`   // 列表的分页页码，默认值为1
	ScalingConfigList []CtyunScalingConfigInfo `tfsdk:"scaling_config_list"`
}
