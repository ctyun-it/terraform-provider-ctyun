package resource

import (
	"context"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctvpc"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "terraform-provider-ctyun/internal/extend/terraform/validator"
)

func NewCtyunSubnet() resource.Resource {
	return &ctyunSubnet{}
}

type ctyunSubnet struct {
	meta       *common.CtyunMetadata
	vpcService *business.VpcService
}

func (c *ctyunSubnet) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_subnet"
}

func (c *ctyunSubnet) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026755/10197656**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "id",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "支持字母、中文、数字，下划线以及-，中文/英文字母开头，长度 2-32",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 32),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z\u4e00-\u9fa5][0-9a-zA-Z_\u4e00-\u9fa5-]+$"), "子网名称不符合规则"),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "vpcId",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "描述，长度最大为128",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(128),
				},
			},
			"dns": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "子网dns列表, 最多同时支持4个dns地址",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.SizeAtMost(4),
					setvalidator.ValueStringsAre(validator2.Ip()),
				},
			},
			"cidr": schema.StringAttribute{
				Required:    true,
				Description: "网段，取值范围：10.0.0.0/8~10.255.255.0/24或者172.16.0.0/12~172.31.255.0/24或者192.168.0.0/16~192.168.255.0/24。约束：必须是cidr格式，例如:192.168.0.0/16",
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
				Description: "是否开启IPv6网段，false：不开启，true: 开启，默认为不开启false，注意：在子网内开启IPv6网段时，必须保证所在vpc也启用了开启IPv6网段",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Default: booldefault.StaticBool(false),
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "子网类型，common：普通子网，cbm：裸金属子网，默认为普通子网common",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString(business.SubnetTypeCommon),
				Validators: []validator.String{
					stringvalidator.OneOf(business.SubnetTypes...),
				},
			},
			"gateway_ip": schema.StringAttribute{
				Computed:    true,
				Description: "子网网关",
			},
			"ipv4_start": schema.StringAttribute{
				Computed:    true,
				Description: "子网网段起始ip",
			},
			"ipv4_end": schema.StringAttribute{
				Computed:    true,
				Description: "子网网段结束ip",
			},
			"ipv6_start": schema.StringAttribute{
				Computed:    true,
				Description: "子网内可用的起始ipv6地址",
			},
			"ipv6_end": schema.StringAttribute{
				Computed:    true,
				Description: "子网内可用的结束ipv6地址",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目id，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
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

func (c *ctyunSubnet) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunSubnetConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.checkCreate(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	regionId := plan.RegionId.ValueString()
	projectId := plan.ProjectId.ValueString()
	var dnsList []string
	var strings []types.String
	plan.Dns.ElementsAs(ctx, &strings, true)
	for _, dns := range strings {
		dnsList = append(dnsList, dns.ValueString())
	}
	resp, err := c.meta.Apis.CtVpcApis.SubnetCreateApi.Do(ctx, c.meta.Credential, &ctvpc.SubnetCreateRequest{
		RegionId:    regionId,
		ProjectId:   projectId,
		ClientToken: uuid.NewString(),
		VpcId:       plan.VpcId.ValueString(),
		Name:        plan.Name.ValueString(),
		Cidr:        plan.Cidr.ValueString(),
		Description: plan.Description.ValueString(),
		DnsList:     dnsList,
		EnableIpv6:  plan.EnableIpv6.ValueBool(),
		SubnetType:  plan.Type.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	plan.Id = types.StringValue(resp.SubnetId)
	plan.RegionId = types.StringValue(regionId)
	plan.ProjectId = types.StringValue(projectId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, ctyunRequestError := c.getAndMergeSubnet(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunSubnet) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunSubnetConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeSubnet(ctx, state)
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

func (c *ctyunSubnet) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state CtyunSubnetConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	var plan CtyunSubnetConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	dnsList := []string{}
	var strings []types.String
	plan.Dns.ElementsAs(ctx, &strings, true)
	for _, t := range strings {
		dnsList = append(dnsList, t.ValueString())
	}
	_, err := c.meta.Apis.CtVpcApis.SubnetUpdateApi.Do(ctx, c.meta.Credential, &ctvpc.SubnetUpdateRequest{
		SubnetId:    state.Id.ValueString(),
		RegionId:    state.RegionId.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		DnsList:     dnsList,
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	instance, ctyunRequestError := c.getAndMergeSubnet(ctx, state)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunSubnet) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunSubnetConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtVpcApis.SubnetDeleteApi.Do(ctx, c.meta.Credential, &ctvpc.SubnetDeleteRequest{
		ClientToken: uuid.NewString(),
		RegionId:    state.RegionId.ValueString(),
		SubnetId:    state.Id.ValueString(),
		ProjectId:   state.ProjectId.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [subnetId],[vpcId],[regionId]
func (c *ctyunSubnet) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunSubnetConfig
	var subnetId, vpcId, regionId string
	err := terraform_extend.Split(request.ID, &subnetId, &vpcId, &regionId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.Id = types.StringValue(subnetId)
	cfg.VpcId = types.StringValue(vpcId)
	cfg.RegionId = types.StringValue(regionId)

	instance, err := c.getAndMergeSubnet(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunSubnet) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.vpcService = business.NewVpcService(meta)
}

// getAndMergeSubnet 查询子网
func (c *ctyunSubnet) getAndMergeSubnet(ctx context.Context, cfg CtyunSubnetConfig) (*CtyunSubnetConfig, error) {
	resp, err := c.meta.Apis.CtVpcApis.SubnetQueryApi.Do(ctx, c.meta.Credential, &ctvpc.SubnetQueryRequest{
		RegionId:    cfg.RegionId.ValueString(),
		ProjectId:   cfg.ProjectId.ValueString(),
		ClientToken: uuid.NewString(),
		SubnetId:    cfg.Id.ValueString(),
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiSubnetNotFound {
			return nil, nil
		}
		return nil, err
	}

	dl := []types.String{}
	for _, dns := range resp.DnsList {
		dl = append(dl, types.StringValue(dns))
	}
	dnsList, _ := types.SetValueFrom(ctx, types.StringType, dl)

	subnetType, err2 := business.SubnetTypeMap.ToOriginalScene(resp.Type, business.SubnetTypeMapScene1)
	if err2 != nil {
		return nil, err2
	}

	cfg.Id = types.StringValue(resp.SubnetId)
	cfg.Name = types.StringValue(resp.Name)
	cfg.VpcId = types.StringValue(resp.VpcId)
	cfg.Cidr = types.StringValue(resp.Cidr)
	cfg.Description = types.StringValue(resp.Description)
	cfg.Dns = dnsList
	cfg.EnableIpv6 = types.BoolValue(resp.EnableIpv6)
	cfg.Type = types.StringValue(subnetType.(string))
	cfg.GatewayIp = types.StringValue(resp.Gateway)
	cfg.Ipv4Start = types.StringValue(resp.Start)
	cfg.Ipv4End = types.StringValue(resp.End)
	cfg.Ipv6Start = types.StringValue(resp.Ipv6Start)
	cfg.Ipv6End = types.StringValue(resp.Ipv6End)
	return &cfg, nil
}

// checkCreate 校验创建动作是否能执行
func (c *ctyunSubnet) checkCreate(ctx context.Context, plan CtyunSubnetConfig) error {
	return c.vpcService.MustExist(ctx, plan.VpcId.ValueString(), plan.RegionId.ValueString(), plan.ProjectId.ValueString())
}

type CtyunSubnetConfig struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	VpcId       types.String `tfsdk:"vpc_id"`
	Cidr        types.String `tfsdk:"cidr"`
	Description types.String `tfsdk:"description"`
	Dns         types.Set    `tfsdk:"dns"`
	EnableIpv6  types.Bool   `tfsdk:"enable_ipv6"`
	Type        types.String `tfsdk:"type"`
	GatewayIp   types.String `tfsdk:"gateway_ip"`
	Ipv4Start   types.String `tfsdk:"ipv4_start"`
	Ipv4End     types.String `tfsdk:"ipv4_end"`
	Ipv6Start   types.String `tfsdk:"ipv6_start"`
	Ipv6End     types.String `tfsdk:"ipv6_end"`
	ProjectId   types.String `tfsdk:"project_id"`
	RegionId    types.String `tfsdk:"region_id"`
}
