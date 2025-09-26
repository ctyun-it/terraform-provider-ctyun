package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunSecurityGroups{}
	_ datasource.DataSourceWithConfigure = &ctyunSecurityGroups{}
)

type ctyunSecurityGroups struct {
	meta *common.CtyunMetadata
}

func NewCtyunSecurityGroups() datasource.DataSource {
	return &ctyunSecurityGroups{}
}

func (c *ctyunSecurityGroups) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_security_groups"
}

type CtyunSecurityGroupRule struct {
	Direction       types.String `tfsdk:"direction"`
	Priority        types.Int32  `tfsdk:"priority"`
	Ethertype       types.String `tfsdk:"ethertype"`
	Protocol        types.String `tfsdk:"protocol"`
	Range           types.String `tfsdk:"range"`
	DestCidrIp      types.String `tfsdk:"dest_cidr_ip"`
	Description     types.String `tfsdk:"description"`
	CreateTime      types.String `tfsdk:"create_time"`
	RuleID          types.String `tfsdk:"rule_id"`
	SecurityGroupID types.String `tfsdk:"security_group_id"`
	Action          types.String `tfsdk:"action"`
}

type CtyunSecurityGroupsModel struct {
	Name                  types.String             `tfsdk:"name"`
	SecurityGroupID       types.String             `tfsdk:"security_group_id"`
	VmNum                 types.Int32              `tfsdk:"vm_num"`
	Origin                types.String             `tfsdk:"origin"`
	VpcName               types.String             `tfsdk:"vpc_name"`
	VpcID                 types.String             `tfsdk:"vpc_id"`
	CreationTime          types.String             `tfsdk:"creation_time"`
	Description           types.String             `tfsdk:"description"`
	SecurityGroupRuleList []CtyunSecurityGroupRule `tfsdk:"security_group_rule_list"`
}

type CtyunSecurityGroupsConfig struct {
	RegionID        types.String `tfsdk:"region_id"`
	VpcID           types.String `tfsdk:"vpc_id"`
	PageNo          types.Int32  `tfsdk:"page_no"`
	PageSize        types.Int32  `tfsdk:"page_size"`
	InstanceID      types.String `tfsdk:"instance_id"`
	SecurityGroupID types.String `tfsdk:"security_group_id"`

	CurrentCount   types.Int32                `tfsdk:"current_count"`
	TotalCount     types.Int32                `tfsdk:"total_count"`
	TotalPage      types.Int32                `tfsdk:"total_page"`
	SecurityGroups []CtyunSecurityGroupsModel `tfsdk:"security_groups"`
}

func (c *ctyunSecurityGroups) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10028520`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"vpc_id": schema.StringAttribute{
				Optional:    true,
				Description: "虚拟私有云ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码，默认值为1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "分页查询时每页的行数，最大值为50，默认值为10。",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "实例ID",
			},
			"security_group_id": schema.StringAttribute{
				Optional:    true,
				Description: "安全组ID",
			},
			"current_count": schema.Int32Attribute{
				Computed:    true,
				Description: "分页查询时每页的行数。",
			},
			"total_count": schema.Int32Attribute{
				Computed:    true,
				Description: "总数。",
			},
			"total_page": schema.Int32Attribute{
				Computed:    true,
				Description: "总页数。",
			},
			"security_groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "安全组名称",
						},
						"security_group_id": schema.StringAttribute{
							Computed:    true,
							Description: "安全组id",
						},
						"vm_num": schema.Int32Attribute{
							Computed:    true,
							Description: "相关云主机",
						},
						"origin": schema.StringAttribute{
							Computed:    true,
							Description: "表示是否是默认安全组",
						},
						"vpc_name": schema.StringAttribute{
							Computed:    true,
							Description: "vpc名称",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "安全组所属的专有网络。",
						},
						"creation_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "安全组描述信息。",
						},
						"security_group_rule_list": schema.ListNestedAttribute{
							Computed:    true,
							Description: "安全组规则信息",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"direction": schema.StringAttribute{
										Computed:    true,
										Description: "出方向-egress、入方向-ingress",
									},
									"priority": schema.Int64Attribute{
										Computed:    true,
										Description: "优先级:0~100",
									},
									"ethertype": schema.StringAttribute{
										Computed:    true,
										Description: "IP类型:IPv4、IPv6",
									},
									"protocol": schema.StringAttribute{
										Computed:    true,
										Description: "协议:ANY、TCP、UDP、ICMP、ICMP6",
									},
									"range": schema.StringAttribute{
										Computed:    true,
										Description: "接口范围/ICMP类型:1-65535",
									},
									"dest_cidr_ip": schema.StringAttribute{
										Computed:    true,
										Description: "远端地址:0.0.0.0/0",
									},
									"description": schema.StringAttribute{
										Computed:    true,
										Description: "安全组规则描述信息。",
									},
									"create_time": schema.StringAttribute{
										Computed:    true,
										Description: "创建时间，UTC时间。",
									},
									"rule_id": schema.StringAttribute{
										Computed:    true,
										Description: "唯一标识ID",
									},
									"security_group_id": schema.StringAttribute{
										Computed:    true,
										Description: "安全组ID",
									},
									"action": schema.StringAttribute{
										Computed:    true,
										Description: "拒绝策略:允许-accept拒绝-drop",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunSecurityGroups) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunSecurityGroupsConfig
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

	params := &ctvpc.CtvpcNewQuerySecurityGroupsRequest{
		RegionID: regionId,
	}
	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	instanceId := config.InstanceID.ValueString()
	securityGroupId := config.SecurityGroupID.ValueString()
	vpcId := config.VpcID.ValueString()
	if pageNo > 0 {
		params.PageNo = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	if instanceId != "" {
		params.InstanceID = &instanceId
	}
	if vpcId != "" {
		params.VpcID = &vpcId
	}
	if securityGroupId != "" {
		params.QueryContent = &securityGroupId
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcNewQuerySecurityGroupsApi.Do(ctx, c.meta.SdkCredential, params)
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
	config.SecurityGroups = []CtyunSecurityGroupsModel{}
	config.TotalPage = types.Int32Value(resp.ReturnObj.TotalPage)
	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
	config.CurrentCount = types.Int32Value(resp.ReturnObj.CurrentCount)
	for _, s := range resp.ReturnObj.SecurityGroups {
		item := CtyunSecurityGroupsModel{
			Name:            utils.SecStringValue(s.SecurityGroupName),
			SecurityGroupID: utils.SecStringValue(s.Id),
			VmNum:           types.Int32Value(s.VmNum),
			Origin:          utils.SecStringValue(s.Origin),
			VpcName:         utils.SecStringValue(s.VpcName),
			VpcID:           utils.SecStringValue(s.VpcID),
			CreationTime:    utils.SecStringValue(s.CreationTime),
			Description:     utils.SecStringValue(s.Description),
		}
		for _, r := range s.SecurityGroupRuleList {
			rule := CtyunSecurityGroupRule{
				Direction:       utils.SecStringValue(r.Direction),
				Priority:        types.Int32Value(r.Priority),
				Ethertype:       utils.SecStringValue(r.Ethertype),
				Protocol:        utils.SecStringValue(r.Protocol),
				Range:           utils.SecStringValue(r.RawRange),
				DestCidrIp:      utils.SecStringValue(r.DestCidrIp),
				Description:     utils.SecStringValue(r.Description),
				CreateTime:      utils.SecStringValue(r.CreateTime),
				RuleID:          utils.SecStringValue(r.Id),
				SecurityGroupID: utils.SecStringValue(r.SecurityGroupID),
				Action:          utils.SecStringValue(r.Action),
			}
			item.SecurityGroupRuleList = append(item.SecurityGroupRuleList, rule)
		}
		config.SecurityGroups = append(config.SecurityGroups, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunSecurityGroups) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
