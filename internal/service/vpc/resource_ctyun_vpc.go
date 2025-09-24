package vpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
)

func NewCtyunVpc() resource.Resource {
	return &ctyunVpc{}
}

type ctyunVpc struct {
	meta *common.CtyunMetadata
}

func (c *ctyunVpc) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc"
}

func (c *ctyunVpc) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026755**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "id",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "虚拟私有云名称。取值范围：2-32，支持数字、字母、中文、_(下划线)、-（中划线）。约束：同一个租户下的名称不能重复。(中文/英文字母开头)，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 32),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z\\x{4e00}-\\x{9fa5}][0-9a-zA-Z_\\x{4e00}-\\x{9fa5}-]+$"), "虚拟私有云名称不符合规则"),
				},
			},
			"cidr": schema.StringAttribute{
				Required:    true,
				Description: "VPC的网段。建议您使用 192.168.0.0/16、172.16.0.0/12、10.0.0.0/8 三个RFC标准私网网段及其子网作为专有网络的主IPv4网段，网段掩码有效范围为8-28位",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.Cidr(),
				},
			},
			"enable_ipv6": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否开启IPv6网段。false：不开启，true: 开启，默认为不开启false",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Default: booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "描述，长度最大为128，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(128),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
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

func (c *ctyunVpc) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunVpcConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := plan.RegionId.ValueString()
	projectId := plan.ProjectId.ValueString()
	resp, err := c.meta.Apis.CtVpcApis.VpcCreateApi.Do(ctx, c.meta.Credential, &ctvpc.VpcCreateRequest{
		RegionId:    regionId,
		ProjectId:   projectId,
		ClientToken: uuid.NewString(),
		Name:        plan.Name.ValueString(),
		Cidr:        plan.Cidr.ValueString(),
		Description: plan.Description.ValueString(),
		EnableIpv6:  plan.EnableIpv6.ValueBool(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	plan.Id = types.StringValue(resp.VpcId)
	plan.RegionId = types.StringValue(regionId)
	plan.ProjectId = types.StringValue(projectId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, ctyunRequestError := c.getAndMergeVpc(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunVpc) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunVpcConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeVpc(ctx, state)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunVpc) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state CtyunVpcConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	var plan CtyunVpcConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtVpcApis.VpcUpdateApi.Do(ctx, c.meta.Credential, &ctvpc.VpcUpdateRequest{
		VpcId:       state.Id.ValueString(),
		RegionId:    state.RegionId.ValueString(),
		ProjectId:   state.ProjectId.ValueString(),
		ClientToken: uuid.NewString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	instance, ctyunRequestError := c.getAndMergeVpc(ctx, state)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunVpc) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunVpcConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtVpcApis.VpcDeleteApi.Do(ctx, c.meta.Credential, &ctvpc.VpcDeleteRequest{
		VpcId:       state.Id.ValueString(),
		RegionId:    state.RegionId.ValueString(),
		ProjectId:   state.ProjectId.ValueString(),
		ClientToken: uuid.NewString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [vpcId],[regionId],[projectId]
func (c *ctyunVpc) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunVpcConfig
	var vpcId, regionId, projectId string
	err := terraform_extend.Split(request.ID, &vpcId, &regionId, &projectId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.Id = types.StringValue(vpcId)
	cfg.RegionId = types.StringValue(regionId)
	cfg.ProjectId = types.StringValue(projectId)

	instance, err := c.getAndMergeVpc(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunVpc) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// getVpc 查询vpc
func (c *ctyunVpc) getAndMergeVpc(ctx context.Context, cfg CtyunVpcConfig) (*CtyunVpcConfig, error) {
	resp, err := c.meta.Apis.CtVpcApis.VpcQueryApi.Do(ctx, c.meta.Credential, &ctvpc.VpcQueryRequest{
		RegionId:    cfg.RegionId.ValueString(),
		ProjectId:   cfg.ProjectId.ValueString(),
		ClientToken: uuid.NewString(),
		VpcId:       cfg.Id.ValueString(),
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiVpcNotFound {
			return nil, nil
		}
		return nil, err
	}
	cfg.Id = types.StringValue(resp.VpcId)
	cfg.Name = types.StringValue(resp.Name)
	cfg.Description = types.StringValue(resp.Description)
	cfg.Cidr = types.StringValue(resp.Cidr)
	cfg.EnableIpv6 = types.BoolValue(resp.Ipv6Enabled)
	return &cfg, nil
}

type CtyunVpcConfig struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Cidr        types.String `tfsdk:"cidr"`
	EnableIpv6  types.Bool   `tfsdk:"enable_ipv6"`
	Description types.String `tfsdk:"description"`
	ProjectId   types.String `tfsdk:"project_id"`
	RegionId    types.String `tfsdk:"region_id"`
}
