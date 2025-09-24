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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunScalings{}
	_ datasource.DataSourceWithConfigure = &CtyunScalings{}
)

type CtyunScalings struct {
	meta *common.CtyunMetadata
}

func NewCtyunScalings() datasource.DataSource {
	return &CtyunScalings{}
}

func (c *CtyunScalings) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_scalings"
}

func (c *CtyunScalings) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunScalings) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027725**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id",
			},
			"id": schema.Int64Attribute{
				Optional:    true,
				Description: "伸缩组id",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "项目id",
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
			"scaling_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"config_list": schema.SetAttribute{
							ElementType: types.Int32Type,
							Computed:    true,
							Description: "伸缩组配置ID列表",
						},
						"health_period": schema.Int32Attribute{
							Computed:    true,
							Description: "健康检查时间间隔",
						},
						"max_count": schema.Int32Attribute{
							Computed:    true,
							Description: "最大云主机数",
						},
						"min_count": schema.Int32Attribute{
							Computed:    true,
							Description: "最小云主机数",
						},
						"expected_count": schema.Int32Attribute{
							Computed:    true,
							Description: "期望云主机数",
						},
						"move_out_strategy": schema.StringAttribute{
							Computed:    true,
							Description: "实例移出策略",
						},
						"create_date": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"group_id": schema.Int64Attribute{
							Computed:    true,
							Description: "伸缩组ID",
						},
						"update_date": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间",
						},
						"health_mode": schema.StringAttribute{
							Computed:    true,
							Description: "健康检查方式",
						},
						"use_lb": schema.Int32Attribute{
							Computed:    true,
							Description: "是否使用负载均衡",
						},
						"subnet_id_list": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "子网ID列表",
						},
						"vpc_cidr": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟私有云网段",
						},
						"status": schema.Int32Attribute{
							Computed:    true,
							Description: "伸缩组状态",
						},
						"vpc_name": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟私有云名称",
						},
						"instance_count": schema.Int32Attribute{
							Computed:    true,
							Description: "伸缩组包含云主机数量",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "企业项目ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "伸缩组名称",
						},
						"security_group_id_list": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "多可用区资源池安全组ID列表",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟私有云ID",
						},
						"az_strategy": schema.StringAttribute{
							Computed:    true,
							Description: "扩容策略类型",
						},
						"delete_protection": schema.BoolAttribute{
							Computed:    true,
							Description: "是否开启伸缩组保护",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunScalings) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunScalingsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	if regionId == "" {
		err = errors.New("region id 为空")
		return
	}
	params := &scaling.ScalingGroupListRequest{
		RegionID: regionId,
		PageNo:   1,
		PageSize: 10,
	}

	if !config.ID.IsNull() && !config.ID.IsUnknown() {
		params.GroupID = config.ID.ValueInt64()
	}

	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		params.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = errors.New("获取弹性伸缩列表信息返回nil，请稍后重试或联系研发人员！")
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var scalingList []CtyunScalingsModel
	// 处理查询结果
	for _, scalingItem := range resp.ReturnObj.ScalingGroups {
		var scalingInfo CtyunScalingsModel
		var diags diag.Diagnostics

		scalingInfo.ConfigList, diags = types.SetValueFrom(ctx, types.Int32Type, scalingItem.ConfigList)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		scalingInfo.HealthPeriod = types.Int32Value(scalingItem.HealthPeriod)
		scalingInfo.MaxCount = types.Int32Value(scalingItem.MaxCount)
		scalingInfo.MinCount = types.Int32Value(scalingItem.MinCount)
		scalingInfo.ExpectedCount = types.Int32Value(scalingItem.ExpectedCount)
		scalingInfo.MoveOutStrategy = types.StringValue(business.ScalingMoveOutStrategyDictRev[scalingItem.MoveOutStrategy])
		scalingInfo.CreateDate = types.StringValue(scalingItem.CreateDate)
		scalingInfo.GroupID = types.Int64Value(scalingItem.GroupID)
		scalingInfo.UpdateDate = types.StringValue(scalingItem.UpdateDate)
		scalingInfo.HealthMode = types.StringValue(business.ScalingHealthModeDictRev[scalingItem.HealthMode])
		scalingInfo.UseLb = types.Int32Value(scalingItem.UseLb)

		scalingInfo.SubnetIDList, diags = types.ListValueFrom(ctx, types.StringType, scalingItem.SubnetIDList)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		scalingInfo.VpcCidr = types.StringValue(scalingItem.VpcCidr)
		scalingInfo.Status = types.Int32Value(scalingItem.Status)
		scalingInfo.VpcName = types.StringValue(scalingItem.VpcName)
		scalingInfo.InstanceCount = types.Int32Value(scalingItem.InstanceCount)
		scalingInfo.ProjectID = types.StringValue(scalingItem.ProjectIDEcs)
		scalingInfo.Name = types.StringValue(scalingItem.Name)
		scalingInfo.SecurityGroupIDList, diags = types.SetValueFrom(ctx, types.StringType, scalingItem.SecurityGroupIDList)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return
		}
		scalingInfo.VpcID = types.StringValue(scalingItem.VpcID)
		scalingInfo.AzStrategy = types.StringValue(business.ScalingAzStrategyDictRev[scalingItem.AzStrategy])
		scalingInfo.DeleteProtection = types.BoolValue(*scalingItem.DeleteProtection)
		scalingList = append(scalingList, scalingInfo)
	}
	config.ScalingList = scalingList
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunScalingsModel struct {
	ConfigList          types.Set    `tfsdk:"config_list"`            // 伸缩组配置ID列表
	HealthPeriod        types.Int32  `tfsdk:"health_period"`          // 健康检查时间间隔
	MaxCount            types.Int32  `tfsdk:"max_count"`              // 最大云主机数
	MinCount            types.Int32  `tfsdk:"min_count"`              // 最小云主机数
	ExpectedCount       types.Int32  `tfsdk:"expected_count"`         // 期望云主机数
	MoveOutStrategy     types.String `tfsdk:"move_out_strategy"`      // 实例移出策略
	CreateDate          types.String `tfsdk:"create_date"`            // 创建时间
	GroupID             types.Int64  `tfsdk:"group_id"`               // 伸缩组ID
	UpdateDate          types.String `tfsdk:"update_date"`            // 更新时间
	HealthMode          types.String `tfsdk:"health_mode"`            // 健康检查方式
	UseLb               types.Int32  `tfsdk:"use_lb"`                 // 是否使用负载均衡
	SubnetIDList        types.List   `tfsdk:"subnet_id_list"`         // 子网ID列表
	VpcCidr             types.String `tfsdk:"vpc_cidr"`               // 虚拟私有云网段
	Status              types.Int32  `tfsdk:"status"`                 // 伸缩组状态
	VpcName             types.String `tfsdk:"vpc_name"`               // 虚拟私有云名称
	InstanceCount       types.Int32  `tfsdk:"instance_count"`         // 伸缩组包含云主机数量
	ProjectID           types.String `tfsdk:"project_id"`             // 企业项目ID
	Name                types.String `tfsdk:"name"`                   // 伸缩组名称
	SecurityGroupIDList types.Set    `tfsdk:"security_group_id_list"` // 多可用区资源池安全组ID列表
	VpcID               types.String `tfsdk:"vpc_id"`                 // 虚拟私有云ID
	AzStrategy          types.String `tfsdk:"az_strategy"`            // 扩容策略类型
	DeleteProtection    types.Bool   `tfsdk:"delete_protection"`      // 是否开启伸缩组保护
}

type CtyunScalingsConfig struct {
	RegionID    types.String         `tfsdk:"region_id"`    // 资源池id
	ID          types.Int64          `tfsdk:"id"`           // 伸缩组id
	ProjectID   types.String         `tfsdk:"project_id"`   // 项目id
	PageNo      types.Int32          `tfsdk:"page_no"`      // 页码
	PageSize    types.Int32          `tfsdk:"page_size"`    // 分页查询时设置的每页行数，取值范围:[1~100]，默认值为10
	ScalingList []CtyunScalingsModel `tfsdk:"scaling_list"` // 弹性伸缩列表
}
