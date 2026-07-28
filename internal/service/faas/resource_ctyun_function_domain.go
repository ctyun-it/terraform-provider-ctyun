package faas

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/cf"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CtyunFunctionDomain{}
	_ resource.ResourceWithConfigure   = &CtyunFunctionDomain{}
	_ resource.ResourceWithImportState = &CtyunFunctionDomain{}
)

func NewCtyunFunctionDomain() resource.Resource {
	return &CtyunFunctionDomain{}
}

type CtyunFunctionDomain struct {
	meta *common.CtyunMetadata
	name string
}

type CtyunFunctionDomainConfig struct {
	ID              types.String `tfsdk:"id"`
	DomainName      types.String `tfsdk:"domain_name"`
	Protocol        types.String `tfsdk:"protocol"`
	Description     types.String `tfsdk:"description"`
	CertName        types.String `tfsdk:"cert_name"`
	CertCertificate types.String `tfsdk:"cert_certificate"`
	CertPrivateKey  types.String `tfsdk:"cert_private_key"`
	RouteFuncName   types.String `tfsdk:"route_function_name"`
	RoutePath       types.String `tfsdk:"route_path"`
	RouteQualifier  types.String `tfsdk:"route_qualifier"`
	RouteEnableJwt  types.Int32  `tfsdk:"route_enable_jwt"`
	RouteMethods    types.List   `tfsdk:"route_methods"`
	RouteFuncID     types.Int32  `tfsdk:"route_function_id"`
	RouteFuncUnique types.String `tfsdk:"route_function_unique_name"`
	AuthType        types.String `tfsdk:"auth_type"`
	AuthJwtJwks     types.String `tfsdk:"auth_jwt_jwks"`
	CnameCheck      types.Bool   `tfsdk:"cname_check"`
	CnameValid      types.Bool   `tfsdk:"cname_valid"`
	DomainStatus    types.String `tfsdk:"domain_status"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
	RegionID        types.String `tfsdk:"region_id"`
}

func (c *CtyunFunctionDomain) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function_domain"
	c.name = resp.TypeName
}

func (c *CtyunFunctionDomain) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理自定义域名", "函数计算（FaaS）", "https://www.ctyun.cn/document/10006234/10532586"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_name": schema.StringAttribute{
				Required:    true,
				Description: "自定义域名。要求：由多个以点分隔的字符串组成，可包含字母、数字中划线、下划线，单个字符串不超过 63 个字符，域名总长度不超过 254 个字符 支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 254),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9][-a-zA-Z0-9_]*(\.[a-zA-Z0-9][-a-zA-Z0-9_]*)*$`), "域名格式不正确"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "协议类型，可选值：HTTP、HTTPS。默认为 HTTP 支持更新",
				Default:     stringdefault.StaticString("HTTP"),
				Validators: []validator.String{
					stringvalidator.OneOf("HTTP", "HTTPS"),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "描述，长度最大为 512，支持更新",
				Default:     stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(512),
				},
			},
			"cname_check": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否检查 CNAME 配置，默认为 false",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"cname_valid": schema.BoolAttribute{
				Computed:    true,
				Description: "CNAME 是否有效",
			},
			"domain_status": schema.StringAttribute{
				Computed:    true,
				Description: "域名备案状态",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池 ID，如果不填则默认使用 provider ctyun 中的 region_id 或环境变量中的 CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"cert_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "证书名称",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(128),
				},
			},
			"cert_certificate": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "证书内容",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cert_private_key": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "私钥内容",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"route_function_name": schema.StringAttribute{
				Optional:    true,
				Description: "函数名称",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 64),
				},
			},
			"route_path": schema.StringAttribute{
				Optional:    true,
				Description: "请求路径",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"route_qualifier": schema.StringAttribute{
				Optional:    true,
				Description: "函数版本",
			},
			"route_enable_jwt": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "是否启用 JWT，0-不启用，1-启用",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"route_methods": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "请求方法列表",
			},
			"route_function_id": schema.Int32Attribute{
				Computed:    true,
				Description: "函数 ID",
			},
			"route_function_unique_name": schema.StringAttribute{
				Computed:    true,
				Description: "函数唯一名称",
			},
			"auth_type": schema.StringAttribute{
				Optional:    true,
				Description: "认证类型",
			},
			"auth_jwt_jwks": schema.StringAttribute{
				Optional:    true,
				Description: "JWK 字符串",
			},
		},
	}
}

func (c *CtyunFunctionDomain) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunFunctionDomain) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunFunctionDomainConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.create(ctx, &plan)
	if err != nil {
		return
	}

	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunFunctionDomain) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunFunctionDomainConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			err = nil
			resp.State.RemoveResource(ctx)
		}
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunFunctionDomain) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunFunctionDomainConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state CtyunFunctionDomainConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &plan, &state)
	if err != nil {
		return
	}

	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunFunctionDomain) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunFunctionDomainConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunFunctionDomain) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例：%s 失败：%s", c.name, req.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [domain_name],<region_id>", c.name)
			resp.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunFunctionDomainConfig
	var domainName, regionID string

	if strings.Count(req.ID, common.ImportSeparator) == 0 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		domainName = req.ID
	} else {
		err = terraform_extend.Split(req.ID, &domainName, &regionID)
		if err != nil {
			return
		}
	}

	if domainName == "" {
		err = fmt.Errorf("domain_name 不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("region_id 不能为空")
		return
	}
	config.DomainName = types.StringValue(domainName)
	config.RegionID = types.StringValue(regionID)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (c *CtyunFunctionDomain) create(ctx context.Context, plan *CtyunFunctionDomainConfig) (err error) {
	req := &cf.CfCreateCustomDomainRequest{
		RegionId:   plan.RegionID.ValueString(),
		DomainName: plan.DomainName.ValueString(),
	}

	if !plan.Protocol.IsNull() && !plan.Protocol.IsUnknown() {
		protocol := plan.Protocol.ValueString()
		req.Protocol = &protocol
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		description := plan.Description.ValueString()
		req.Description = &description
	}

	if !plan.CertName.IsNull() && !plan.CertCertificate.IsNull() && !plan.CertPrivateKey.IsNull() {
		req.CertConfig = &cf.CfCreateCertConfig{
			CertName:    plan.CertName.ValueString(),
			Certificate: plan.CertCertificate.ValueString(),
			PrivateKey:  plan.CertPrivateKey.ValueString(),
		}
	}

	if !plan.RouteFuncName.IsNull() && !plan.RoutePath.IsNull() {
		methods := make([]string, 0)
		for _, method := range plan.RouteMethods.Elements() {
			if method, ok := method.(types.String); ok && !method.IsNull() && !method.IsUnknown() {
				methods = append(methods, method.ValueString())
			}
		}

		pathConfig := cf.CfCreatePathConfig{
			FunctionName:       plan.RouteFuncName.ValueString(),
			Path:               plan.RoutePath.ValueString(),
			EnableJwt:          plan.RouteEnableJwt.ValueInt32(),
			FunctionUniqueName: plan.RouteFuncUnique.ValueString(),
		}

		if !plan.RouteQualifier.IsNull() && !plan.RouteQualifier.IsUnknown() {
			pathConfig.Qualifier = plan.RouteQualifier.ValueString()
		}

		if !plan.RouteFuncID.IsNull() && !plan.RouteFuncID.IsUnknown() {
			pathConfig.FunctionId = plan.RouteFuncID.ValueInt32()
		}

		req.RouteConfig = &cf.CfCreateRouteConfig{
			Routes: []cf.CfCreatePathConfig{pathConfig},
		}
	}

	if !plan.AuthType.IsNull() {
		req.AuthConfig = &cf.CfCreateAuthConfig{
			AuthType: plan.AuthType.ValueString(),
		}

		if !plan.AuthJwtJwks.IsNull() {
			req.AuthConfig.JwtConfig = &cf.CfCreateJwtConfig{
				Jwks: plan.AuthJwtJwks.ValueString(),
			}
		}
	}

	resp, err := c.meta.Apis.SdkCtCfApis.CfCreateCustomDomainApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if *resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	plan.ID = types.StringValue(plan.DomainName.ValueString())
	return
}

func (c *CtyunFunctionDomain) update(ctx context.Context, plan *CtyunFunctionDomainConfig, state *CtyunFunctionDomainConfig) (err error) {
	req := &cf.CfUpdateCustomDomainRequest{
		DomainName: plan.DomainName.ValueString(),
		RegionId:   plan.RegionID.ValueString(),
	}

	// 只有当协议发生变化时才传递
	if !plan.Protocol.Equal(state.Protocol) && !plan.Protocol.IsNull() && !plan.Protocol.IsUnknown() {
		protocol := plan.Protocol.ValueString()
		req.Protocol = &protocol
	}

	// 只有当描述发生变化时才传递
	if !plan.Description.Equal(state.Description) {
		description := plan.Description.ValueString()
		req.Description = &description
	}

	// 只有当证书配置发生变化时才传递
	if !plan.CertName.Equal(state.CertName) || !plan.CertCertificate.Equal(state.CertCertificate) || !plan.CertPrivateKey.Equal(state.CertPrivateKey) {
		if !plan.CertName.IsNull() && !plan.CertCertificate.IsNull() && !plan.CertPrivateKey.IsNull() {
			req.CertConfig = &cf.CfUpdateCertConfig{
				CertName:    plan.CertName.ValueString(),
				Certificate: plan.CertCertificate.ValueString(),
				PrivateKey:  plan.CertPrivateKey.ValueString(),
			}
		}
	}

	// 只有当路由配置发生变化时才传递
	if !plan.RouteFuncName.Equal(state.RouteFuncName) || !plan.RoutePath.Equal(state.RoutePath) ||
		!plan.RouteQualifier.Equal(state.RouteQualifier) || !plan.RouteEnableJwt.Equal(state.RouteEnableJwt) ||
		!plan.RouteFuncUnique.Equal(state.RouteFuncUnique) {
		if !plan.RouteFuncName.IsNull() && !plan.RoutePath.IsNull() {
			pathConfig := cf.CfUpdatePathConfig{
				FunctionName:       plan.RouteFuncName.ValueString(),
				Path:               plan.RoutePath.ValueString(),
				EnableJwt:          plan.RouteEnableJwt.ValueInt32(),
				FunctionUniqueName: plan.RouteFuncUnique.ValueString(),
			}

			if !plan.RouteQualifier.IsNull() && !plan.RouteQualifier.IsUnknown() {
				qualifier := plan.RouteQualifier.ValueString()
				pathConfig.Qualifier = &qualifier
			}

			req.RouteConfig = &cf.CfUpdateRouteConfig{
				Routes: []cf.CfUpdatePathConfig{pathConfig},
			}
		}
	}

	// 只有当认证配置发生变化时才传递
	if !plan.AuthType.Equal(state.AuthType) || !plan.AuthJwtJwks.Equal(state.AuthJwtJwks) {
		if !plan.AuthType.IsNull() {
			req.AuthConfig = &cf.CfUpdateAuthConfig{
				AuthType: plan.AuthType.ValueString(),
			}

			if !plan.AuthJwtJwks.IsNull() {
				req.AuthConfig.JwtConfig = &cf.CfUpdateJwtConfig{
					Jwks: plan.AuthJwtJwks.ValueString(),
				}
			}
		}
	}

	resp, err := c.meta.Apis.SdkCtCfApis.CfUpdateCustomDomainApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if *resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return
}

func (c *CtyunFunctionDomain) delete(ctx context.Context, state *CtyunFunctionDomainConfig) (err error) {
	req := &cf.CfDeleteCustomDomainRequest{
		DomainName: state.DomainName.ValueString(),
		RegionId:   state.RegionID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtCfApis.CfDeleteCustomDomainApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if *resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}

	return
}

func (c *CtyunFunctionDomain) getAndMerge(ctx context.Context, config *CtyunFunctionDomainConfig) (err error) {
	req := &cf.CfGetCustomDomainRequest{
		DomainName: config.DomainName.ValueString(),
		RegionId:   config.RegionID.ValueString(),
		CnameCheck: config.CnameCheck.ValueBool(),
	}

	resp, err := c.meta.Apis.SdkCtCfApis.CfGetCustomDomainApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if *resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	returnObj := resp.ReturnObj

	config.DomainName = types.StringPointerValue(returnObj.DomainName)
	config.Protocol = types.StringPointerValue(returnObj.Protocol)
	config.Description = types.StringPointerValue(returnObj.Description)
	config.CnameValid = types.BoolPointerValue(returnObj.CnameValid)
	config.DomainStatus = types.StringPointerValue(returnObj.DomainStatus)
	config.CreatedAt = types.StringPointerValue(returnObj.CreatedAt)
	config.UpdatedAt = types.StringPointerValue(returnObj.UpdatedAt)

	// 处理 CertConfig
	if returnObj.CertConfig != nil {
		config.CertName = types.StringValue(utils.SecString(returnObj.CertConfig.CertName))
		config.CertCertificate = types.StringValue(utils.SecString(returnObj.CertConfig.Certificate))
		config.CertPrivateKey = types.StringValue(utils.SecString(returnObj.CertConfig.PrivateKey))
	}

	// 处理 RouteConfig
	if returnObj.RouteConfig != nil && len(returnObj.RouteConfig.Routes) > 0 {
		routeItem := returnObj.RouteConfig.Routes[0]
		methods := make([]string, 0)
		if routeItem.Methods != nil {
			for _, method := range routeItem.Methods {
				if method != nil {
					methods = append(methods, *method)
				}
			}
		}
		config.RouteFuncName = types.StringPointerValue(routeItem.FunctionName)
		config.RoutePath = types.StringPointerValue(routeItem.Path)
		config.RouteQualifier = types.StringPointerValue(routeItem.Qualifier)
		config.RouteEnableJwt = types.Int32PointerValue(routeItem.EnableJwt)
		config.RouteMethods, _ = types.ListValueFrom(ctx, types.StringType, methods)
		config.RouteFuncID = types.Int32PointerValue(routeItem.FunctionId)
		config.RouteFuncUnique = types.StringPointerValue(routeItem.FunctionUniqueName)
	} else {
		// 确保在 API 未返回路由配置时清除相关字段
		config.RouteFuncName = types.StringNull()
		config.RoutePath = types.StringNull()
		config.RouteQualifier = types.StringNull()
		config.RouteEnableJwt = types.Int32Null()
		config.RouteMethods, _ = types.ListValueFrom(ctx, types.StringType, []string{})
		config.RouteFuncID = types.Int32Null()
		config.RouteFuncUnique = types.StringNull()
	}

	// 处理 AuthConfig
	if returnObj.AuthConfig != nil {
		config.AuthType = types.StringPointerValue(returnObj.AuthConfig.AuthType)
		if returnObj.AuthConfig.AuthType != nil && *returnObj.AuthConfig.AuthType == "" {
			config.AuthType = types.StringNull()
		}
		if returnObj.AuthConfig.JwtConfig != nil {
			config.AuthJwtJwks = types.StringPointerValue(returnObj.AuthConfig.JwtConfig.Jwks)
		}
	} else {
		config.AuthType = types.StringNull()
		config.AuthJwtJwks = types.StringNull()
	}

	config.ID = types.StringValue(config.DomainName.ValueString())
	return
}
