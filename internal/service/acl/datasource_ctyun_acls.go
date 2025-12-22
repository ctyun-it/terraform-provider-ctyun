package acl

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunAcls{}
	_ datasource.DataSourceWithConfigure = &ctyunAcls{}
)

type ctyunAcls struct {
	meta *common.CtyunMetadata
}

func NewCtyunAcls() datasource.DataSource {
	return &ctyunAcls{}
}
func (c *ctyunAcls) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunAcls) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_acls"
}

func (c *ctyunAcls) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10028583",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，默认使用provider配置",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "ACL ID，精确查询",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "项目ID",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "ACL名称，查询条件之一",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "分页页码，默认为1",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页记录数，默认为10。取值范围：1~50",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},
			"acls": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ACL ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "ACL名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "ACL描述",
						},
						"apply_to_public_lb": schema.BoolAttribute{
							Computed:    true,
							Description: "是否应用到公网负载均衡",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "所属VPC ID",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "是否启用",
						},
						"in_policy_id": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "入方向策略ID集合",
						},
						"out_policy_id": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "出方向策略ID集合",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间，为UTC格式",
						},
						"subnet_ids": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "关联的子网ID集合",
						},
					},
				},
				Description: "ACL列表",
			},
		},
	}
}

func (c *ctyunAcls) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunAclsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	config.RegionID = types.StringValue(regionId)
	params := &ctvpc.CtvpcListAclRequest{
		RegionID: config.RegionID.ValueString(),
		PageNo:   1,
		PageSize: 10,
	}
	if !config.ID.IsNull() {
		params.AclID = config.ID.ValueStringPointer()
	}
	if !config.ProjectID.IsUnknown() && !config.ProjectID.IsNull() {
		params.ProjectID = config.ProjectID.ValueStringPointer()
	}
	if !config.Name.IsUnknown() && !config.Name.IsNull() {
		params.Name = config.Name.ValueStringPointer()
	}
	if !config.PageNo.IsNull() {
		params.PageNo = config.PageNo.ValueInt32()
	}
	if !config.PageSize.IsNull() {
		params.PageSize = config.PageSize.ValueInt32()
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListAclApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("获取acl列表失败，接口返回nil，请联系研发确认问题原因！")
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var acls []CtyunAclInfoModel
	for _, aclItem := range resp.ReturnObj.Acls {
		var acl CtyunAclInfoModel
		acl.ID = types.StringValue(*aclItem.AclID)
		acl.Name = types.StringValue(*aclItem.Name)
		acl.Enabled = types.BoolValue(map[string]bool{business.AclEnable: true, business.AclDisable: true}[utils.SecString(aclItem.Enabled)])
		acl.Description = types.StringValue(*aclItem.Description)
		acl.ApplyToPublicLb = types.BoolValue(*aclItem.ApplyToPublicLb)
		acl.CreateTime = types.StringValue(*aclItem.CreatedAt)
		acl.UpdateTime = types.StringValue(*aclItem.UpdatedAt)

		var diags diag.Diagnostics
		acl.OutPolicyID, diags = types.SetValueFrom(ctx, types.StringType, aclItem.OutPolicyID)
		if diags.HasError() {
			err = fmt.Errorf(diags[0].Detail())
			return
		}
		acl.InPolicyID, diags = types.SetValueFrom(ctx, types.StringType, aclItem.InPolicyID)
		if diags.HasError() {
			err = fmt.Errorf(diags[0].Detail())
			return
		}
		acl.SubnetIDs, diags = types.SetValueFrom(ctx, types.StringType, aclItem.SubnetIDs)
		if diags.HasError() {
			err = fmt.Errorf(diags[0].Detail())
			return
		}
		acls = append(acls, acl)
	}
	config.Acls = acls
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunAclInfoModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	ApplyToPublicLb types.Bool   `tfsdk:"apply_to_public_lb"`
	VpcID           types.String `tfsdk:"vpc_id"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	InPolicyID      types.Set    `tfsdk:"in_policy_id"`
	OutPolicyID     types.Set    `tfsdk:"out_policy_id"`
	CreateTime      types.String `tfsdk:"create_time"`
	UpdateTime      types.String `tfsdk:"update_time"`
	SubnetIDs       types.Set    `tfsdk:"subnet_ids"`
}

type CtyunAclsConfig struct {
	RegionID  types.String        `tfsdk:"region_id"`
	ID        types.String        `tfsdk:"id"`
	ProjectID types.String        `tfsdk:"project_id"`
	Name      types.String        `tfsdk:"name"`
	PageNo    types.Int32         `tfsdk:"page_no"`
	PageSize  types.Int32         `tfsdk:"page_size"`
	Acls      []CtyunAclInfoModel `tfsdk:"acls"`
}
