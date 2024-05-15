package datasource

import (
	"context"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctiam"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/common"
)

func NewCtyunIamUserGroups() datasource.DataSource {
	return &ctyunIamUserGroups{}
}

type ctyunIamUserGroups struct {
	meta *common.CtyunMetadata
}

func (c *ctyunIamUserGroups) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_user_groups"
}

func (c *ctyunIamUserGroups) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10345725/10355805**`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "用户组名称，为模糊搜索",
			},
			"page_size": schema.Int64Attribute{
				Required:    true,
				Description: "每页显示数量，取值范围1-1000",
				Validators: []validator.Int64{
					int64validator.Between(1, 1000),
				},
			},
			"page_no": schema.Int64Attribute{
				Required:    true,
				Description: "当前页码",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"groups": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "用户组id",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "用户组名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "用户组信息",
						},
						"user_count": schema.Int64Attribute{
							Computed:    true,
							Description: "用户数量",
						},
					},
				},
				Computed:    true,
				Description: "用户组列表",
			},
		},
	}
}

func (c *ctyunIamUserGroups) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CtyunIamUserGroupsConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}
	response, err := c.meta.Apis.CtIamApis.UserGroupQueryApi.Do(ctx, c.meta.Credential, &ctiam.UserGroupQueryRequest{
		GroupName: config.Name.ValueString(),
		PageNum:   config.PageNo.ValueInt64(),
		PageSize:  config.PageSize.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	var infos []CtyunIamUserGroupInfo
	for _, info := range response.Result {
		infos = append(infos, CtyunIamUserGroupInfo{
			Id:          types.StringValue(info.Id),
			Name:        types.StringValue(info.GroupName),
			UserCount:   types.Int64Value(info.UserCount),
			Description: types.StringValue(info.GroupIntro),
		})
	}
	config.Result = infos
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunIamUserGroups) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	mete := req.ProviderData.(*common.CtyunMetadata)
	c.meta = mete
}

type CtyunIamUserGroupsConfig struct {
	Name     types.String            `tfsdk:"name"`
	Result   []CtyunIamUserGroupInfo `tfsdk:"groups"`
	PageSize types.Int64             `tfsdk:"page_size"`
	PageNo   types.Int64             `tfsdk:"page_no"`
}

type CtyunIamUserGroupInfo struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	UserCount   types.Int64  `tfsdk:"user_count"`
}
