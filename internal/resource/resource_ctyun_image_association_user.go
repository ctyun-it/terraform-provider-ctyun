package resource

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctimage"
	defaults2 "terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
)

func NewCtyunImageAssociationUser() resource.Resource {
	return &ctyunImageAssociationUser{}
}

type ctyunImageAssociationUser struct {
	meta         *common.CtyunMetadata
	imageService *business.ImageService
}

func (c *ctyunImageAssociationUser) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_image_association_user"
}

func (c *ctyunImageAssociationUser) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027726**`,
		Attributes: map[string]schema.Attribute{
			"image_id": schema.StringAttribute{
				Required:    true,
				Description: "要共享的私有镜像id",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "类型，share：表示将私有镜像分享给其他人，receive：表示接收或拒绝来自其他分享的私有镜像",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.ImageAssociationUserTypes...),
				},
			},
			"user_email": schema.StringAttribute{
				Optional:    true,
				Description: "共享镜像的接收者，此值为对应账号的邮箱，当type为share时此值必填",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(path.MatchRoot("type"), types.StringValue(business.ImageAssociationUserTypeShare)),
					validator2.ConflictsWithEqualStrings(path.MatchRoot("type"), types.StringValue(business.ImageAssociationUserTypeReceive)),
					validator2.Email(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
		},
	}
}

func (c *ctyunImageAssociationUser) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunImageAssociationUserConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.check(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	regionId := plan.RegionId.ValueString()
	switch plan.Type.ValueString() {
	case business.ImageAssociationUserTypeShare:
		// 判断镜像是否已经分享过了，如果分享过了就报错
		business.RollingPage(func(param business.RollingPageParam) business.RollingPageResult {
			resp, err := c.meta.Apis.CtImageApis.ImageShareListApi.Do(ctx, c.meta.Credential, &ctimage.ImageShareListRequest{
				RegionId: plan.RegionId.ValueString(),
				ImageId:  plan.ImageId.ValueString(),
				PageSize: 1,
				PageNo:   param.CurrentPage,
			})
			if err != nil {
				response.Diagnostics.AddError(err.Error(), err.Error())
				return business.NoNextRollingPage
			}
			for _, image := range resp.Images {
				if image.DestinationUser == plan.UserEmail.ValueString() {
					msg := "镜像已经分享给" + plan.UserEmail.ValueString() + "，请勿再次分享"
					response.Diagnostics.AddError(msg, msg)
					return business.NoNextRollingPage
				}
			}
			return business.RollingPageResult{
				TotalCount:       resp.TotalCount,
				CurrentPageCount: len(resp.Images),
				ExecuteContinue:  true,
			}
		})
		if response.Diagnostics.HasError() {
			return
		}

		_, err = c.meta.Apis.CtImageApis.ImageShareCreateApi.Do(ctx, c.meta.Credential, &ctimage.ImageShareCreateRequest{
			RegionId:        regionId,
			DestinationUser: plan.UserEmail.ValueString(),
			ImageId:         plan.ImageId.ValueString(),
		})
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
	case business.ImageAssociationUserTypeReceive:
		_, err = c.meta.Apis.CtImageApis.ImageShareAcceptApi.Do(ctx, c.meta.Credential, &ctimage.ImageShareAcceptRequest{
			RegionId: regionId,
			ImageId:  plan.ImageId.ValueString(),
		})
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunImageAssociationUser) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {

}

func (c *ctyunImageAssociationUser) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {

}

func (c *ctyunImageAssociationUser) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunImageAssociationUserConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	var err error
	switch state.Type.ValueString() {
	case business.ImageAssociationUserTypeShare:
		// 判断是否已经接受了镜像，如果接受了，那么就不能取消了
		business.RollingPage(func(param business.RollingPageParam) business.RollingPageResult {
			resp, err := c.meta.Apis.CtImageApis.ImageShareListApi.Do(ctx, c.meta.Credential, &ctimage.ImageShareListRequest{
				RegionId: state.RegionId.ValueString(),
				ImageId:  state.ImageId.ValueString(),
				PageSize: 1,
				PageNo:   param.CurrentPage,
			})
			if err != nil {
				response.Diagnostics.AddError(err.Error(), err.Error())
				return business.NoNextRollingPage
			}
			for _, image := range resp.Images {
				if image.DestinationUser == state.UserEmail.ValueString() {
					if image.Status == business.ImageStatusAccepted {
						msg := "镜像已经被接受，取消分享前请联系" + state.UserEmail.ValueString() + "执行拒绝动作"
						response.Diagnostics.AddError(msg, msg)
					}
					return business.NoNextRollingPage
				}
			}
			return business.RollingPageResult{
				TotalCount:       resp.TotalCount,
				CurrentPageCount: len(resp.Images),
				ExecuteContinue:  true,
			}
		})
		if response.Diagnostics.HasError() {
			return
		}

		_, err = c.meta.Apis.CtImageApis.ImageShareDeleteApi.Do(ctx, c.meta.Credential, &ctimage.ImageShareDeleteRequest{
			RegionId:        state.RegionId.ValueString(),
			DestinationUser: state.UserEmail.ValueString(),
			ImageId:         state.ImageId.ValueString(),
		})
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
	case business.ImageAssociationUserTypeReceive:
		_, err = c.meta.Apis.CtImageApis.ImageShareRejectApi.Do(ctx, c.meta.Credential, &ctimage.ImageShareRejectRequest{
			RegionId: state.RegionId.ValueString(),
			ImageId:  state.ImageId.ValueString(),
		})
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
	}
}

func (c *ctyunImageAssociationUser) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.imageService = business.NewImageService(meta)
}

// check 校验
func (c *ctyunImageAssociationUser) check(ctx context.Context, cfg CtyunImageAssociationUserConfig) error {
	return c.imageService.MustExist(ctx, cfg.ImageId.ValueString(), cfg.RegionId.ValueString())
}

type CtyunImageAssociationUserConfig struct {
	ImageId   types.String `tfsdk:"image_id"`
	Type      types.String `tfsdk:"type"`
	UserEmail types.String `tfsdk:"user_email"`
	RegionId  types.String `tfsdk:"region_id"`
}
