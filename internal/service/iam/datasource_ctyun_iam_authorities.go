package iam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctiam"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCtyunAuthorities() datasource.DataSource {
	return &ctyunIamAuthorities{}
}

type ctyunIamAuthorities struct {
	meta *common.CtyunMetadata
}

func (c *ctyunIamAuthorities) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_authorities"
}

func (c *ctyunIamAuthorities) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10345725/10409363**`,
		Attributes: map[string]schema.Attribute{
			"service_id": schema.Int64Attribute{
				Required:    true,
				Description: "服务id，可以用ctyun_services进行查询",
			},
			"authorities": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "权限点名称",
						},
						"code": schema.StringAttribute{
							Computed:    true,
							Description: "权限点编码",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "权限点描述",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunIamAuthorities) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CtyunAuthoritiesConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	response, err := c.meta.Apis.CtIamApis.AuthorityListApi.Do(ctx, c.meta.Credential, &ctiam.AuthorityListRequest{
		ServiceId: int(config.ServiceId.ValueInt64()),
	})
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	var authorities []CtyunPolicyAuthoritiesConfig
	for _, s := range response.AuthorityList {
		authorities = append(authorities, CtyunPolicyAuthoritiesConfig{
			Name:        types.StringValue(s.Name),
			Code:        types.StringValue(s.Code),
			Description: types.StringValue(s.Description),
		})
	}
	config.Authorities = authorities
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (c *ctyunIamAuthorities) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunPolicyAuthoritiesConfig struct {
	Name        types.String `tfsdk:"name"`
	Code        types.String `tfsdk:"code"`
	Description types.String `tfsdk:"description"`
}

type CtyunAuthoritiesConfig struct {
	ServiceId   types.Int64                    `tfsdk:"service_id"`
	Authorities []CtyunPolicyAuthoritiesConfig `tfsdk:"authorities"`
}
