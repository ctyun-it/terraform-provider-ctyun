package elb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ datasource.DataSource              = &ctyunElbAcls{}
	_ datasource.DataSourceWithConfigure = &ctyunElbAcls{}
)

type ctyunElbAcls struct {
	meta *common.CtyunMetadata
}

func (c *ctyunElbAcls) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunElbAcls) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_acls"
}

func (c *ctyunElbAcls) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026756/10032777**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"ids": schema.StringAttribute{
				Optional:    true,
				Description: "访问控制ID列表, 逗号分割",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "访问控制名称,只能由数字，字母，-组成不能以数字和-开头，最大长度32",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(32),
				},
			},
			"acls": schema.ListNestedAttribute{
				Computed:    true,
				Description: "访问控制列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "可用区名称",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "项目ID",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "访问控制ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "访问控制名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"source_ips": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "IP地址的集合或者CIDR",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunElbAcls) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunElbAclsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	params := &ctelb.CtelbListAccessControlRequest{
		RegionID: regionId,
		Name:     "",
	}
	if config.IDs.ValueString() != "" {
		params.IDs = strings.Split(config.IDs.ValueString(), ",")
	}
	if config.Name.ValueString() != "" {
		params.Name = config.Name.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbListAccessControlApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var acls []CtyunElbAclsModel
	var diags diag.Diagnostics
	for _, aclItem := range resp.ReturnObj {
		var acl CtyunElbAclsModel
		acl.AzName = types.StringValue(aclItem.AzName)
		acl.ProjectID = types.StringValue(aclItem.ProjectID)
		acl.ID = types.StringValue(aclItem.ID)
		acl.Name = types.StringValue(aclItem.Name)
		acl.Description = types.StringValue(aclItem.Description)
		acl.CreateTime = types.StringValue(aclItem.CreateTime)
		acl.SourceIps, diags = types.SetValueFrom(ctx, types.StringType, aclItem.SourceIps)
		if diags.HasError() {
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

func NewCtyunElbAcls() datasource.DataSource {
	return &ctyunElbAcls{}
}

type CtyunElbAclsConfig struct {
	RegionID types.String        `tfsdk:"region_id"` //区域ID
	IDs      types.String        `tfsdk:"ids"`       //访问控制ID列表
	Name     types.String        `tfsdk:"name"`      //访问控制名称,只能由数字，字母，-组成不能以数字和-开头，最大长度32
	Acls     []CtyunElbAclsModel `tfsdk:"acls"`      //访问控制列表
}

type CtyunElbAclsModel struct {
	AzName      types.String `tfsdk:"az_name"`     //可用区名称
	ProjectID   types.String `tfsdk:"project_id"`  //项目ID
	ID          types.String `tfsdk:"id"`          //访问控制ID
	Name        types.String `tfsdk:"name"`        //访问控制名称
	Description types.String `tfsdk:"description"` //描述
	SourceIps   types.Set    `tfsdk:"source_ips"`  //IP地址的集合或者CIDR
	CreateTime  types.String `tfsdk:"create_time"` //创建时间，为UTC格式
}
