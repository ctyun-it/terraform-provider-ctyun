package iam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCtyunIamUserAks() datasource.DataSource {
	return &CtyunIamUserAks{}
}

type CtyunIamUserAks struct {
	meta       *common.CtyunMetadata
	iamService *business.IamService
}

func (c *CtyunIamUserAks) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_user_aks"
}

func (c *CtyunIamUserAks) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10345725/10355289`,
		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{
				Required:    true,
				Description: "用户ID",
			},
			"ak": schema.StringAttribute{
				Optional:    true,
				Description: "用户AK",
			},
			"list": schema.SetNestedAttribute{
				Computed:    true,
				Description: "AK列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ak": schema.StringAttribute{
							Computed:    true,
							Description: "用户AK",
						},
						"sk": schema.StringAttribute{
							Computed:    true,
							Description: "用户SK",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "密钥状态",
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

func (c *CtyunIamUserAks) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunIamUserAksConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	aks, err := c.iamService.QueryAkList(ctx, config.UserID.ValueString())
	if err != nil {
		return
	}
	var sk string
	config.List = []CtyunIamUserAksModel{}
	for _, a := range aks {
		if config.AK.ValueString() == "" || config.AK.ValueString() == utils.SecString(a.AccessKey) {
			item := CtyunIamUserAksModel{
				AK:         utils.SecStringValue(a.AccessKey),
				Enabled:    types.BoolValue(map[string]bool{business.AkEnabled: true, business.AkDisabled: false}[utils.SecString(a.Status)]),
				CreateTime: types.StringValue(utils.FromUnixToUTC(a.CreatedTime)),
			}
			sk, err = c.iamService.DecryptSK(utils.SecString(a.SecretKey), c.meta.SdkCredential.GetAccessKey())
			if err != nil {
				return
			}
			item.SK = types.StringValue(sk)
			config.List = append(config.List, item)
		}
	}
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *CtyunIamUserAks) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.iamService = business.NewIamService(meta)
}

type CtyunIamUserAksConfig struct {
	UserID types.String           `tfsdk:"user_id"`
	AK     types.String           `tfsdk:"ak"`
	List   []CtyunIamUserAksModel `tfsdk:"list"`
}

type CtyunIamUserAksModel struct {
	AK         types.String `tfsdk:"ak"`
	SK         types.String `tfsdk:"sk"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	CreateTime types.String `tfsdk:"create_time"`
}
