package vpc

import (
	"context"
	"fmt"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	AvaliableFlowPackageStatus = "有效"
)

var (
	_ datasource.DataSource              = &ctyunFlowPackages{}
	_ datasource.DataSourceWithConfigure = &ctyunFlowPackages{}
)

type ctyunFlowPackages struct {
	meta *common.CtyunMetadata
}

func NewCtyunFlowPackages() datasource.DataSource {
	return &ctyunFlowPackages{}
}

func (c *ctyunFlowPackages) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_flow_packages"
}

type CtyunFlowPackagesModel struct {
	ID               types.String  `tfsdk:"id"`
	Status           types.String  `tfsdk:"status"`
	CycleType        types.String  `tfsdk:"cycle_type"`
	EffectiveTime    types.String  `tfsdk:"effective_time"`
	ExpireTime       types.String  `tfsdk:"expire_time"`
	PackageName      types.String  `tfsdk:"package_name"`
	TotalVolume      types.Float64 `tfsdk:"total_volume"`
	LeftVolume       types.Float64 `tfsdk:"left_volume"`
	MasterResourceID types.String  `tfsdk:"master_resource_id"`
}

type CtyunFlowPackagesConfig struct {
	RegionID     types.String             `tfsdk:"region_id"`
	FlowPackages []CtyunFlowPackagesModel `tfsdk:"flow_packages"`
}

func (c *ctyunFlowPackages) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("查询流量包列表", "弹性IP（Elastic IP，EIP）", "https://www.ctyun.cn/document/10026753/10032118"),
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"flow_packages": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "流量包ID",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "状态，可能的取值：初始、有效、退订、过期、销毁",
						},
						"cycle_type": schema.StringAttribute{
							Computed:    true,
							Description: "支持的取值：包小时、包天、包周、包月、包年",
						},
						"effective_time": schema.StringAttribute{
							Computed:    true,
							Description: "生效时间",
						},
						"expire_time": schema.StringAttribute{
							Computed:    true,
							Description: "过期时间",
						},
						"package_name": schema.StringAttribute{
							Computed:    true,
							Description: "套餐名",
						},
						"total_volume": schema.Float64Attribute{
							Computed:    true,
							Description: "总流量",
						},
						"left_volume": schema.Float64Attribute{
							Computed:    true,
							Description: "剩余流量",
						},
						"master_resource_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源id",
						},
					},
				},
			},
		}}
}

func (c *ctyunFlowPackages) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunFlowPackagesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}
	config.RegionID = types.StringValue(regionId)
	// 组装请求体
	params := &ctvpc.CtvpcListFlowPackageRequest{
		RegionID: regionId,
	}

	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListFlowPackageApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if utils.SecString(resp.ErrorCode) == common.OpenapiFlowPackageNotFound || resp.ReturnObj == nil {
		return
	}

	for _, v := range resp.ReturnObj {
		if utils.SecString(v.Status) != AvaliableFlowPackageStatus {
			continue
		}
		item := CtyunFlowPackagesModel{
			ID:               utils.SecStringValue(v.Id),
			Status:           utils.SecStringValue(v.Status),
			CycleType:        utils.SecStringValue(v.CycleType),
			EffectiveTime:    utils.SecStringValue(v.EffectiveTime),
			ExpireTime:       utils.SecStringValue(v.ExpireTime),
			PackageName:      utils.SecStringValue(v.PackageName),
			TotalVolume:      types.Float64Value(v.TotalVolumn),
			LeftVolume:       types.Float64Value(v.LeftVolumn),
			MasterResourceID: utils.SecStringValue(v.MasterResourceBundleId),
		}
		config.FlowPackages = append(config.FlowPackages, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunFlowPackages) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
