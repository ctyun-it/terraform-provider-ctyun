package iam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCtyunIamUsers() datasource.DataSource {
	return &CtyunIamUsers{}
}

type CtyunIamUsers struct {
	meta       *common.CtyunMetadata
	iamService *business.IamService
}

func (c *CtyunIamUsers) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_users"
}

func (c *CtyunIamUsers) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10345725/10355289`,
		Attributes: map[string]schema.Attribute{
			"page_size": schema.Int32Attribute{
				Required:    true,
				Description: "每页显示数量",
				Validators: []validator.Int32{
					int32validator.Between(1, 1000),
				},
			},
			"page_no": schema.Int32Attribute{
				Required:    true,
				Description: "当前页码",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"users": schema.ListNestedAttribute{
				Computed:    true,
				Description: "用户列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"user_id": schema.StringAttribute{
							Computed:    true,
							Description: "用户ID",
						},
						"email": schema.StringAttribute{
							Computed:    true,
							Description: "用户邮箱",
						},
						"account_id": schema.StringAttribute{
							Computed:    true,
							Description: "账户ID",
						},
						"phone": schema.StringAttribute{
							Computed:    true,
							Description: "手机号",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "用户名",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"is_root": schema.BoolAttribute{
							Computed:    true,
							Description: "是否主用户",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "用户状态",
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
								},
							},
							Computed:    true,
							Description: "关联的用户组列表",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunIamUsers) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunIamUsersConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	resp, err := c.iamService.QueryUserList(ctx, config.PageNo.ValueInt32(), config.PageSize.ValueInt32())
	if err != nil {
		return
	}

	config.Users = []CtyunIamUsersModel{}
	for _, user := range resp.Result {
		item := CtyunIamUsersModel{
			UserID:     utils.SecStringValue(user.UserId),
			Email:      utils.SecStringValue(user.LoginEmail),
			Phone:      utils.SecStringValue(user.MobilePhone),
			Name:       utils.SecStringValue(user.UserName),
			AccountID:  utils.SecStringValue(user.AccountId),
			IsRoot:     types.BoolValue(map[int32]bool{1: true, 0: false}[user.IsRoot]),
			Enabled:    types.BoolValue(map[int32]bool{0: true, 1: false}[user.Prohibit]),
			CreateTime: types.StringValue(utils.FromUnixToUTC(user.CreateDate)),
		}
		for _, group := range user.Groups {
			item.Groups = append(item.Groups, CtyunIamUsersGroup{
				ID:          utils.SecStringValue(group.Id),
				Name:        utils.SecStringValue(group.GroupName),
				Description: utils.SecStringValue(group.GroupIntro),
			})
		}
		config.Users = append(config.Users, item)
	}
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *CtyunIamUsers) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.iamService = business.NewIamService(meta)
}

type CtyunIamUsersConfig struct {
	PageSize types.Int32          `tfsdk:"page_size"`
	PageNo   types.Int32          `tfsdk:"page_no"`
	Users    []CtyunIamUsersModel `tfsdk:"users"`
}

// CtyunIamUsersModel 对应 Schema 根节点结构
type CtyunIamUsersModel struct {
	UserID     types.String         `tfsdk:"user_id"`
	Email      types.String         `tfsdk:"email"`
	AccountID  types.String         `tfsdk:"account_id"`
	Phone      types.String         `tfsdk:"phone"`
	Name       types.String         `tfsdk:"name"`
	CreateTime types.String         `tfsdk:"create_time"`
	IsRoot     types.Bool           `tfsdk:"is_root"`
	Enabled    types.Bool           `tfsdk:"enabled"`
	Groups     []CtyunIamUsersGroup `tfsdk:"groups"`
}

type CtyunIamUsersGroup struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}
