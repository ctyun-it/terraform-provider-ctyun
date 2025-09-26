package ebs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctebs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebs"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunEbsSnapshotPolicies{}
	_ datasource.DataSourceWithConfigure = &ctyunEbsSnapshotPolicies{}
)

type ctyunEbsSnapshotPolicies struct {
	meta *common.CtyunMetadata
}

func NewCtyunEbsSnapshotPolicies() datasource.DataSource {
	return &ctyunEbsSnapshotPolicies{}
}

func (c *ctyunEbsSnapshotPolicies) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs_snapshot_policies"
}

type ctyunEbsSnapshotPoliciesModel struct {
	Id                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	RepeatWeekdays           types.String `tfsdk:"repeat_weekdays"`
	RepeatTimes              types.String `tfsdk:"repeat_times"`
	RetentionTime            types.Int32  `tfsdk:"retention_time"`
	ProjectId                types.String `tfsdk:"project_id"`
	SnapshotPolicyStatus     types.String `tfsdk:"snapshot_policy_status"`
	BoundDiskNum             types.Int32  `tfsdk:"bound_disk_num"`
	SnapshotPolicyCreateTime types.String `tfsdk:"snapshot_policy_create_time"`
}

type ctyunEbsSnapshotPoliciesConfig struct {
	RegionID                 types.String                    `tfsdk:"region_id"`
	PolicyName               types.String                    `tfsdk:"name"`
	PolicyID                 types.String                    `tfsdk:"id"`
	SnapshotPolicies         []ctyunEbsSnapshotPoliciesModel `tfsdk:"snapshot_policies"`
	SnapshotPolicyTotalCount types.Int32                     `tfsdk:"snapshot_policy_total_count"`
}

func (c *ctyunEbsSnapshotPolicies) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027696/10118840`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘自动快照策略名称。满足以下规则：只能由数字、英文字母、中划线-、下划线_、点.组成，长度为2-64字符",
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "云硬盘自动快照策略ID，32字符",
			},
			"snapshot_policy_total_count": schema.Int32Attribute{
				Computed:    true,
				Description: "云硬盘自动快照策略总数",
			},
			"snapshot_policies": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘自动快照策略ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "云硬盘自动快照策略名称",
						},
						"repeat_weekdays": schema.StringAttribute{
							Computed:    true,
							Description: "创建快照的重复日期，0-6分别代表周日-周六，多个日期用英文逗号隔开。支持更新",
						},
						"repeat_times": schema.StringAttribute{
							Computed:    true,
							Description: "创建快照的重复时间，0-23分别代表零点-23点，多个时间用英文逗号隔开。支持更新",
						},
						"retention_time": schema.Int32Attribute{
							Computed:    true,
							Description: "创建快照的保留时间，输入范围为[-1，1-65535]，-1代表永久保留。单位为天。支持更新",
							Validators: []validator.Int32{
								int32validator.Between(-1, 65535),
							},
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "企业项目ID",
						},
						"snapshot_policy_status": schema.StringAttribute{
							Computed:    true,
							Description: "自动快照策略状态，取值范围：activated:启用，deactivated：停用",
						},
						"bound_disk_num": schema.Int32Attribute{
							Computed:    true,
							Description: "关联云硬盘的数量",
						},
						"snapshot_policy_create_time": schema.StringAttribute{
							Computed:    true,
							Description: "策略创建时间",
						},
					},
				}},
		},
	}
}

func (c *ctyunEbsSnapshotPolicies) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config ctyunEbsSnapshotPoliciesConfig
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
	config.SnapshotPolicies = []ctyunEbsSnapshotPoliciesModel{}
	// 组装请求体
	// 获取实例详情
	params := &ctebs2.EbsQueryPolicyEbsSnapRequest{
		RegionID:           config.RegionID.ValueString(),
		SnapshotPolicyID:   config.PolicyID.ValueStringPointer(),
		SnapshotPolicyName: config.PolicyName.ValueStringPointer(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsQueryPolicyEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
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
	for _, policy := range resp.ReturnObj.SnapshotPolicyList {

		item := ctyunEbsSnapshotPoliciesModel{
			Id:                       types.StringValue(*policy.SnapshotPolicyID),
			Name:                     types.StringValue(*policy.SnapshotPolicyName),
			SnapshotPolicyStatus:     types.StringValue(*policy.SnapshotPolicyStatus),
			BoundDiskNum:             types.Int32Value(policy.BoundDiskNum),
			SnapshotPolicyCreateTime: types.StringValue(*policy.SnapshotPolicyCreateTime),
			ProjectId:                types.StringValue(*policy.ProjectID),
			RepeatWeekdays:           types.StringValue(*policy.RepeatWeekdays),
			RepeatTimes:              types.StringValue(*policy.RepeatTimes),
			RetentionTime:            types.Int32Value(policy.RetentionTime),
		}

		config.SnapshotPolicies = append(config.SnapshotPolicies, item)
	}
	config.SnapshotPolicyTotalCount = types.Int32Value(resp.ReturnObj.SnapshotPolicyTotalCount)
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunEbsSnapshotPolicies) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
