package common

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctiam"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCtyunServices() datasource.DataSource {
	return &ctyunServices{}
}

type ctyunServices struct {
	meta *common.CtyunMetadata
}

func (c *ctyunServices) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_services"
}

func (c *ctyunServices) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `**服务和产品**`,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "服务/产品类型，region：资源池级云服务，global：全局级云服务，默认为global",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ServiceTypes...),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "服务名称（中文）",
			},
			"services": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "服务id",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "服务名称",
						},
						"code": schema.StringAttribute{
							Computed:    true,
							Description: "服务编码",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "服务描述",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunServices) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CtyunServicesConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := &ctiam.ServiceListRequest{
		ServiceName: config.Name.ValueString(),
	}

	serviceType, err := business.ServiceTypeMap.FromOriginalScene(config.Type.ValueString(), business.ServiceTypeMapScene1)
	if err == nil {
		request.ServiceType = serviceType.(int)
	}

	response, err := c.meta.Apis.CtIamApis.ServiceListApi.Do(ctx, c.meta.Credential, request)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	var services []CtyunServicesServicesConfig
	for _, s := range response.ServiceList {
		services = append(services, CtyunServicesServicesConfig{
			Id:          types.Int64Value(int64(s.Id)),
			Name:        types.StringValue(s.MainServiceName),
			Code:        types.StringValue(s.ServiceCode),
			Description: types.StringValue(s.ServiceDesc),
		})
	}
	config.Services = services
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (c *ctyunServices) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunServicesServicesConfig struct {
	Id          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Code        types.String `tfsdk:"code"`
	Description types.String `tfsdk:"description"`
}

type CtyunServicesConfig struct {
	Type     types.String                  `tfsdk:"type"`
	Name     types.String                  `tfsdk:"name"`
	Services []CtyunServicesServicesConfig `tfsdk:"services"`
}
