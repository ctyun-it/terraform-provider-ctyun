package scaling

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/scaling"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunScalingEcsList{}
	_ datasource.DataSourceWithConfigure = &CtyunScalingEcsList{}
)

type CtyunScalingEcsList struct {
	meta *common.CtyunMetadata
}

func (c *CtyunScalingEcsList) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func NewCtyunScalingEcsList() datasource.DataSource {
	return &CtyunScalingEcsList{}
}

func (c *CtyunScalingEcsList) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_scaling_ecs_list"
}

func (c *CtyunScalingEcsList) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027725/10216515`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"group_id": schema.Int64Attribute{
				Required:    true,
				Description: "伸缩组ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "页码",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页行数，取值范围:[1~100]，默认值为10",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
				},
			},
			"ecs_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "云主机ID",
						},
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源池ID",
						},
						"group_id": schema.Int64Attribute{
							Computed:    true,
							Description: "伸缩组ID",
						},
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "实例所在的可用区",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "企业项目ID",
						},
						"create_date": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "实例ID",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "伸缩活动状态",
						},
						"instance_name": schema.StringAttribute{
							Computed:    true,
							Description: "云主机名称",
						},
						"execution_mode": schema.StringAttribute{
							Computed: true,
							Description: "执行方式，取值范围：1：自动执行策略，2：手动执行策略，3：手动移入实例，4：手动移出实例，5：新建伸缩组满足最小数，6：修改伸缩组满足最大最小限制，" +
								"7：健康检查移入，8：健康检查移出",
						},
						"health_status": schema.StringAttribute{
							Computed:    true,
							Description: "健康检查状态，取值范围：1：正常，2：异常，3：初始化",
						},
						"config_name": schema.StringAttribute{
							Computed:    true,
							Description: "伸缩配置名称",
						},
						"config_id": schema.Int64Attribute{
							Computed:    true,
							Description: "伸缩配置ID",
						},
						"active_id": schema.Int64Attribute{
							Computed:    true,
							Description: "伸缩活动ID",
						},
						"protect_status": schema.StringAttribute{
							Computed:    true,
							Description: "保护状态，取值范围：1：已保护，2：未保护",
						},
						"join_date": schema.StringAttribute{
							Computed:    true,
							Description: "加入时间",
						},
					},
				},
				Description: "伸缩组ECS实例列表",
			},
		},
	}
}

func (c *CtyunScalingEcsList) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunScalingEcsListConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	if regionId == "" {
		err = errors.New("region id 为空")
		return
	}
	params := &scaling.ScalingGroupQueryInstanceListRequest{
		RegionID: regionId,
		GroupID:  config.GroupID.ValueInt64(),
		PageNo:   1,
		PageSize: 100,
	}
	if !config.PageNo.IsNull() {
		params.PageNo = config.PageNo.ValueInt32()
	}
	if !config.PageSize.IsNull() {
		params.PageSize = config.PageSize.ValueInt32()
	}

	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupQueryInstanceListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询弹性伸缩组的列表失败，接口返回nil，伸缩组id：%d", config.GroupID.ValueInt64())
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	instanceList := resp.ReturnObj.InstanceList
	var ecsList []CtyunScalingEcsListModel
	for _, instance := range instanceList {
		var ecs CtyunScalingEcsListModel
		ecs.InstanceID = types.StringValue(instance.InstanceID)
		ecs.RegionID = types.StringValue(instance.RegionID)
		ecs.AzName = types.StringValue(instance.AzName)
		ecs.ProjectID = types.StringValue(instance.ProjectIDEcs)
		ecs.ID = types.Int64Value(int64(instance.Id))
		ecs.Status = types.StringValue(business.ScalingActivityStatusMapRev[instance.Status])
		ecs.InstanceName = types.StringValue(instance.InstanceName)
		ecs.ExecutionMode = types.StringValue(business.ExecutionModeToString[instance.ExecutionMode])
		ecs.HealthStatus = types.StringValue(business.HealthStatusCodeToStr[instance.HealthStatus])
		ecs.ConfigName = types.StringValue(instance.ConfigName)
		ecs.ConfigID = types.Int64Value(instance.ConfigID)
		ecs.ActiveID = types.Int64Value(int64(instance.ActiveID))
		ecs.ProtectStatus = types.StringValue(business.ProtectStatusCodeToStr[instance.ProtectStatus])
		ecs.JoinDate = types.StringValue(instance.JoinDate)
		ecsList = append(ecsList, ecs)
	}
	config.EcsList = ecsList
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

}

// CtyunScalingActivityInstance 伸缩活动实例
type CtyunScalingEcsListModel struct {
	InstanceID    types.String `tfsdk:"instance_id"`    // 云主机ID
	RegionID      types.String `tfsdk:"region_id"`      // 资源池ID
	GroupID       types.Int64  `tfsdk:"group_id"`       // 伸缩组ID
	AzName        types.String `tfsdk:"az_name"`        // 实例所在的可用区
	ProjectID     types.String `tfsdk:"project_id"`     // 企业项目ID
	CreateDate    types.String `tfsdk:"create_date"`    // 创建时间
	ID            types.Int64  `tfsdk:"id"`             // 实例ID
	Status        types.String `tfsdk:"status"`         // 伸缩活动状态
	InstanceName  types.String `tfsdk:"instance_name"`  // 云主机名称
	ExecutionMode types.String `tfsdk:"execution_mode"` // 执行方式
	HealthStatus  types.String `tfsdk:"health_status"`  // 健康检查状态
	ConfigName    types.String `tfsdk:"config_name"`    // 伸缩配置名称
	ConfigID      types.Int64  `tfsdk:"config_id"`      // 伸缩配置ID
	ActiveID      types.Int64  `tfsdk:"active_id"`      // 伸缩活动ID
	ProtectStatus types.String `tfsdk:"protect_status"` // 保护状态
	JoinDate      types.String `tfsdk:"join_date"`      // 加入时间
}

type CtyunScalingEcsListConfig struct {
	RegionID types.String               `tfsdk:"region_id"` // 资源池ID
	GroupID  types.Int64                `tfsdk:"group_id"`  // 伸缩组ID
	PageNo   types.Int32                `tfsdk:"page_no"`   // 页码
	PageSize types.Int32                `tfsdk:"page_size"` // 分页查询时设置的每页行数，取值范围:[1~100]，默认值为10
	EcsList  []CtyunScalingEcsListModel `tfsdk:"ecs_list"`  // ecs列表
}
