package provider

import (
	"context"
	"errors"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-core"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctebs"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctecs"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctiam"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctimage"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctvpc"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"slices"
	"strings"
	"terraform-provider-ctyun/internal/common"
	dataSource2 "terraform-provider-ctyun/internal/datasource"
	sdk_extend "terraform-provider-ctyun/internal/extend/sdk"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
	resource2 "terraform-provider-ctyun/internal/resource"
	"terraform-provider-ctyun/internal/utils"
)

func NewCtyunProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CtyunProvider{
			version: version,
		}
	}
}

type CtyunProvider struct {
	version string
}

func (c *CtyunProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ctyun"
	resp.Version = c.version
}

func (c *CtyunProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema.Attributes = map[string]schema.Attribute{
		"ak": schema.StringAttribute{
			Optional:    true,
			Description: "身份信息ak",
		},
		"sk": schema.StringAttribute{
			Optional:    true,
			Description: "身份信息sk",
		},
		"env": schema.StringAttribute{
			Optional:    true,
			Description: "环境类型env，可选值为：dev：开发环境、test：测试环境、prod：生产环境，默认为生产环境prod",
		},
		"region_id": schema.StringAttribute{
			Optional:    true,
			Description: "资源区域id",
		},
		"az_name": schema.StringAttribute{
			Optional:    true,
			Description: "可用区id，如果是3.0资源池，则此值无需填写；如果是4.0资源池，则填写选用的az_name",
		},
		"project_id": schema.StringAttribute{
			Optional:    true,
			Description: "企业项目id，不填则使用用户默认的企业项目",
		},
		"console_url": schema.StringAttribute{
			Optional:    true,
			Description: "请求分发地址，仅供测试使用，需配合inspect_url_keywords一起使用",
			Validators: []validator.String{
				stringvalidator.AlsoRequires(path.MatchRoot("inspect_url_keywords")),
			},
		},
		"inspect_url_keywords": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "请求拦截的地址，仅供测试使用，如果填入*则表示拦截所有请求，需配合console_url一起使用",
			Validators: []validator.Set{
				setvalidator.AlsoRequires(path.MatchRoot("console_url")),
			},
		},
	}
}

func (c *CtyunProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "开始配置ctyun")

	// 先在环境变量中获取对应的值
	ak := utils.AcquireEnvParam(common.EnvKeyAk)
	sk := utils.AcquireEnvParam(common.EnvKeySk)
	regionId := utils.AcquireEnvParam(common.EnvKeyRegionId)
	azName := utils.AcquireEnvParam(common.EnvKeyAzName)
	env := utils.AcquireEnvParam(common.EnvKeyEnv)
	projectId := utils.AcquireEnvParam(common.EnvKeyProjectId)
	consoleUrl := utils.AcquireEnvParam(common.EnvKeyConsoleUrl)
	inspectUrlKeywords := strings.Split(utils.AcquireEnvParam(common.EnvKeyInspectUrlKeywords), ",")

	// 在配置文件中获取对应的值，并覆盖
	var cfg CtyunProviderConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !cfg.Ak.IsUnknown() && !cfg.Ak.IsNull() {
		ak = cfg.Ak.ValueString()
	}
	if !cfg.Sk.IsUnknown() && !cfg.Sk.IsNull() {
		sk = cfg.Sk.ValueString()
	}
	if !cfg.RegionId.IsUnknown() && !cfg.RegionId.IsNull() {
		regionId = cfg.RegionId.ValueString()
	}
	if !cfg.AzName.IsUnknown() && !cfg.AzName.IsNull() {
		azName = cfg.AzName.ValueString()
	}
	if !cfg.Env.IsUnknown() && !cfg.Env.IsNull() {
		env = cfg.Env.ValueString()
	}
	if !cfg.ProjectId.IsUnknown() && !cfg.ProjectId.IsNull() {
		projectId = cfg.ProjectId.ValueString()
	}
	if !cfg.ConsoleUrl.IsUnknown() && !cfg.ConsoleUrl.IsNull() {
		consoleUrl = cfg.ConsoleUrl.String()
	}
	if !cfg.InspectUrlKeywords.IsUnknown() && !cfg.InspectUrlKeywords.IsNull() {
		cfg.InspectUrlKeywords.ElementsAs(ctx, &inspectUrlKeywords, true)
	}

	// 校验，设置默认的值
	if ak == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("ak"),
			"未找到ak，请在环境变量中配置CTYUN_AK或者在配置文件中设置ak",
			"未找到ak，请在环境变量中配置CTYUN_AK或者在配置文件中设置ak",
		)
		return
	} else {
		err := ctyunsdk.CheckAk(ak)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("ak"),
				err.Error(),
				err.Error(),
			)
			return
		}
	}
	if sk == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("sk"),
			"未找到sk，请在环境变量中配置CTYUN_SK或者在配置文件中设置sk",
			"未找到sk，请在环境变量中配置CTYUN_SK或者在配置文件中设置sk",
		)
		return
	} else {
		err := ctyunsdk.CheckSk(sk)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("sk"),
				err.Error(),
				err.Error(),
			)
			return
		}
	}
	environment := ctyunsdk.Environment(env)
	if env == "" {
		// 默认环境为prod
		env = ctyunsdk.EnvironmentProd
	} else if !slices.Contains(ctyunsdk.Environments, environment) {
		resp.Diagnostics.AddAttributeError(
			path.Root("env"),
			"env配置错误",
			"env配置错误",
		)
		return
	}

	// 设置对应的参数上下文
	ctx = tflog.SetField(ctx, "ak", ak)
	ctx = tflog.SetField(ctx, "regionId", regionId)
	ctx = tflog.SetField(ctx, "azName", azName)
	ctx = tflog.SetField(ctx, "projectId", projectId)
	ctx = tflog.SetField(ctx, "env", env)

	// 构造config
	config := &ctyunsdk.CtyunClientConfig{
		HttpHooks: []ctyunsdk.HttpHook{
			ctyunsdk.AddUserAgentHttpHook{},
			sdk_extend.LogHttpHook{},
		},
	}

	// consoleUrl不为空的情况
	if consoleUrl != "" && len(inspectUrlKeywords) != 0 && env != ctyunsdk.EnvironmentProd {
		hook := sdk_extend.NewAddConsoleUrlHook(consoleUrl, inspectUrlKeywords...)
		config.HttpHooks = append([]ctyunsdk.HttpHook{hook}, config.HttpHooks...)
	}

	var httpClient *http.Client
	switch environment {
	case ctyunsdk.EnvironmentDev:
		httpClient = ctyunsdk.ClientTest()
	case ctyunsdk.EnvironmentTest:
		httpClient = ctyunsdk.ClientTest()
	case ctyunsdk.EnvironmentProd:
		httpClient = ctyunsdk.ClientProd()
	default:
		httpClient = ctyunsdk.ClientProd()
	}
	config.Client = httpClient

	// 构造client
	client, err := ctyunsdk.NewCtyunClient(environment, config)
	if err != nil {
		err := errors.New("创建ctyun client失败：" + err.Error())
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// 构造私秘钥信息
	credential, err := ctyunsdk.NewCredential(ak, sk)
	if err != nil {
		err := errors.New("创建私秘钥失败：" + err.Error())
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	extra := map[string]string{
		common.ExtraRegionId:  regionId,
		common.ExtraAzName:    azName,
		common.ExtraProjectId: projectId,
	}

	// 填充对应的内容信息
	common.InitCtyunMetadata(
		&common.Apis{
			CtEbsApis:   ctebs.NewApis(client),
			CtEcsApis:   ctecs.NewApis(client),
			CtIamApis:   ctiam.NewApis(client),
			CtImageApis: ctimage.NewApis(client),
			CtVpcApis:   ctvpc.NewApis(client),
		},
		*credential,
		extra,
	)
	metadata := common.AcquireCtyunMetadata()
	resp.ResourceData = metadata
	resp.DataSourceData = metadata
}

func (c *CtyunProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return c.buildDataSource(dataSource2.NewCtyunRegions(), dataSource2.NewCtyunServices(), dataSource2.NewCtyunAuthorities(), dataSource2.NewCtyunImages(), dataSource2.NewCtyunEcsFlavors(), dataSource2.NewCtyunIamUserGroups())
}

func (c *CtyunProvider) Resources(_ context.Context) []func() resource.Resource {
	return c.buildResource(resource2.NewCtyunPolicy(), resource2.NewCtyunVpc(), resource2.NewCtyunSecurityGroup(), resource2.NewCtyunSecurityGroupRule(), resource2.NewCtyunEip(), resource2.NewCtyunEipAssociation(), resource2.NewCtyunSubnet(), resource2.NewCtyunEcs(), resource2.NewCtyunIamUser(), resource2.NewCtyunIamUserGroup(), resource2.NewCtyunIdp(), resource2.NewCtyunEbs(), resource2.NewCtyunEbsAssociation(), resource2.NewCtyunImage(), resource2.NewCtyunImageAssociationUser(), resource2.NewCtyunKeypair(), resource2.NewCtyunBandwidth(), resource2.NewCtyunBandwidthAssociationEip(), resource2.NewCtyunPolicyAssociationUserGroup(), resource2.NewCtyunPolicyAssociationUser(), resource2.NewCtyunEnterpriseProject(), resource2.NewCtyunEnterpriseProjectAssociationUserGroup())
}

// buildDataSource 构建datasource
func (c *CtyunProvider) buildDataSource(datasources ...datasource.DataSource) []func() datasource.DataSource {
	advices := terraform_extend.NewAopAdvices()
	var result []func() datasource.DataSource
	for _, dataSource := range datasources {
		result = append(result, terraform_extend.WrapDataSource(dataSource, advices))
	}
	return result
}

// buildResource 构建resource
func (c *CtyunProvider) buildResource(resources ...resource.Resource) []func() resource.Resource {
	advices := terraform_extend.NewAopAdvices()
	var result []func() resource.Resource
	for _, res := range resources {
		result = append(result, terraform_extend.WrapResource(res, advices))
	}
	return result
}

type CtyunProviderConfig struct {
	Ak                 types.String `tfsdk:"ak"`
	Sk                 types.String `tfsdk:"sk"`
	Env                types.String `tfsdk:"env"`
	RegionId           types.String `tfsdk:"region_id"`
	AzName             types.String `tfsdk:"az_name"`
	ProjectId          types.String `tfsdk:"project_id"`
	ConsoleUrl         types.String `tfsdk:"console_url"`
	InspectUrlKeywords types.Set    `tfsdk:"inspect_url_keywords"`
}
