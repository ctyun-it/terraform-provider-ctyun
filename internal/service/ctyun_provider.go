package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ccse2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ccse"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/crs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebm"
	ctebs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebsbackup"
	ctecs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	ctvpc2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/amqp"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctebs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctiam"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctimage"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
	mongodb2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	pgsql2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctzos"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	hpfs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/hpfs"
	ctgkafka "github.com/ctyun-it/terraform-provider-ctyun/internal/core/kafka"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/scaling"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/sfs"
	sdk_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/sdk"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/ccse"
	common2 "github.com/ctyun-it/terraform-provider-ctyun/internal/service/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/ebm"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/ebs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/ecs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/elb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/hpfs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/iam"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/image"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/kafka"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/mongodb"
	mysql2 "github.com/ctyun-it/terraform-provider-ctyun/internal/service/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/nat"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/pgsql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/rabbitmq"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/redis"
	scaling2 "github.com/ctyun-it/terraform-provider-ctyun/internal/service/scaling"
	sfs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/service/sfs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/vpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/vpce"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service/zos"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"slices"
	"strings"
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
			Description: "身份信息AK",
			Sensitive:   true,
		},
		"sk": schema.StringAttribute{
			Optional:    true,
			Description: "身份信息SK",
			Sensitive:   true,
		},
		"env": schema.StringAttribute{
			Optional:    true,
			Description: "环境类型env，可选值为：dev：开发环境、test：测试环境、prod：生产环境，默认为生产环境prod",
			Validators: []validator.String{
				stringvalidator.OneOf(ctyunsdk.EnvironmentDev, ctyunsdk.EnvironmentTest, ctyunsdk.EnvironmentProd),
			},
		},
		"region_id": schema.StringAttribute{
			Optional:    true,
			Description: "资源池ID",
		},
		"az_name": schema.StringAttribute{
			Optional:    true,
			Description: "可用区英文，填写选用资源池的az_name",
		},
		"project_id": schema.StringAttribute{
			Optional:    true,
			Description: "企业项目ID，不填则使用用户默认的企业项目",
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

	// 默认环境为prod
	if env == "" {
		env = ctyunsdk.EnvironmentProd
	}

	environment := ctyunsdk.Environment(env)
	if !slices.Contains(ctyunsdk.Environments, environment) {
		resp.Diagnostics.AddAttributeError(
			path.Root("env"),
			"env配置错误",
			"env配置错误，可选：prod、test",
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
			sdk_extend.MetricHttpHook{},
			sdk_extend.LogHttpHook{},
		},
	}
	coreConfig := &core.CtyunClientConfig{
		HttpHooks: []core.HttpHook{
			ctyunsdk.AddUserAgentHttpHook{},
			sdk_extend.MetricHttpHook{},
			sdk_extend.LogHttpHook{},
		},
	}

	// consoleUrl不为空的情况
	if consoleUrl != "" && len(inspectUrlKeywords) != 0 && env != ctyunsdk.EnvironmentProd {
		hook := sdk_extend.NewAddConsoleUrlHook(consoleUrl, inspectUrlKeywords...)
		config.HttpHooks = append([]ctyunsdk.HttpHook{hook}, config.HttpHooks...)
	}

	var httpClient *http.Client
	var endpointUrl string
	switch environment {
	case ctyunsdk.EnvironmentDev:
		httpClient = ctyunsdk.ClientTest()
		endpointUrl = "https://%s-global.ctapi-test.ctyun.cn:21443"
	case ctyunsdk.EnvironmentTest:
		httpClient = ctyunsdk.ClientTest()
		endpointUrl = "https://%s-global.ctapi-test.ctyun.cn:21443"
	case ctyunsdk.EnvironmentProd:
		httpClient = ctyunsdk.ClientProd()
		endpointUrl = "https://%s-global.ctapi.ctyun.cn"
	default:
		httpClient = ctyunsdk.ClientProd()
		endpointUrl = "https://%s-global.ctapi.ctyun.cn"
	}
	config.Client = httpClient
	coreConfig.Client = httpClient
	// 构造client
	client, err := ctyunsdk.NewCtyunClient(environment, config)
	if err != nil {
		err := errors.New("创建ctyun client失败：" + err.Error())
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	coreClient := core.NewCtyunClient(coreConfig)

	// 构造私秘钥信息
	credential, err := ctyunsdk.NewCredential(ak, sk)
	if err != nil {
		err := errors.New("创建私秘钥失败：" + err.Error())
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	SdkCredential := core.NewCredential(ak, sk)

	extra := map[string]string{
		common.ExtraRegionId:  regionId,
		common.ExtraAzName:    azName,
		common.ExtraProjectId: projectId,
	}

	// 填充对应的内容信息
	common.InitCtyunMetadata(
		&common.Apis{
			CtEbsApis:       ctebs.NewApis(client),
			CtEcsApis:       ctecs.NewApis(client),
			CtIamApis:       ctiam.NewApis(client),
			CtImageApis:     ctimage.NewApis(client),
			CtVpcApis:       ctvpc.NewApis(client),
			CtEbmApis:       ctebm.NewApis(fmt.Sprintf(endpointUrl, ctebm.EndpointName), coreClient),
			SdkCtEbsApis:    ctebs2.NewApis(fmt.Sprintf(endpointUrl, ctebs2.EndpointName), coreClient),
			SdkCtEcsApis:    ctecs2.NewApis(fmt.Sprintf(endpointUrl, ctecs2.EndpointName), coreClient),
			SdkCtVpcApis:    ctvpc2.NewApis(fmt.Sprintf(endpointUrl, ctvpc2.EndpointName), coreClient),
			SdkCtZosApis:    ctzos.NewApis(fmt.Sprintf(endpointUrl, ctzos.EndpointName), coreClient),
			SdkCcseApis:     ccse2.NewApis(fmt.Sprintf(endpointUrl, ccse2.EndpointName), coreClient),
			SdkDcs2Apis:     dcs2.NewApis(fmt.Sprintf(endpointUrl, dcs2.EndpointName), coreClient),
			SdkCtElbApis:    ctelb.NewApis(fmt.Sprintf(endpointUrl, ctelb.EndpointName), coreClient),
			SdkCtMysqlApis:  mysql.NewApis(client),
			SdkKafkaApis:    ctgkafka.NewApis(fmt.Sprintf(endpointUrl, ctgkafka.EndpointName), coreClient),
			SdkAmqpApis:     amqp.NewApis(client),
			SdkCrsApis:      crs.NewApis(fmt.Sprintf(endpointUrl, crs.EndpointName), coreClient),
			SdkCtPgsqlApis:  pgsql2.NewApis(client),
			SdkMongodbApis:  mongodb2.NewApis(client),
			SdkHpfsApis:     hpfs2.NewApis(fmt.Sprintf(endpointUrl, hpfs2.EndpointName), coreClient),
			CtEbsBackupApis: ctebsbackup.NewApis(fmt.Sprintf(endpointUrl, ctebsbackup.EndpointName), coreClient),
			SdkScalingApis:  scaling.NewApis(fmt.Sprintf(endpointUrl, scaling.EndpointName), coreClient),
			SdkSfsApi:       sfs.NewApis(fmt.Sprintf(endpointUrl, sfs.EndpointName), coreClient),
		},
		*credential,
		*SdkCredential,
		extra,
	)
	metadata := common.AcquireCtyunMetadata()
	resp.ResourceData = metadata
	resp.DataSourceData = metadata
}

func (c *CtyunProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return c.buildDataSource(
		common2.NewCtyunRegions(),
		common2.NewCtyunServices(),
		iam.NewCtyunAuthorities(),
		image.NewCtyunImages(),
		ecs.NewCtyunEcsFlavors(),
		iam.NewCtyunIamUserGroups(),
		ebm.NewCtyunEbmDeviceTypes(),
		ebm.NewCtyunEbms(),
		ebm.NewCtyunEbmDeviceRaids(),
		ebm.NewCtyunEbmDeviceImages(),
		nat.NewCtyunNats(),
		elb.NewElbLoadBalancers(),
		nat.NewCtyunSNats(),
		nat.NewCtyunDNats(),
		ebs.NewCtyunEbsVolumes(),
		ebs.NewCtyunEbsSnapshots(),
		ebs.NewCtyunEbsBackups(),
		//ebs.NewCtyunEbsBackupRepos(),
		ebs.NewCtyunEbsBackupPolicies(),
		ebs.NewCtyunEbsSnapshotPolicies(),
		ecs.NewCtyunEcsInstances(),
		ecs.NewCtyunEcsAffinityGroups(),
		ecs.NewCtyunEcsSnapshots(),
		//ecs.NewCtyunEcsBackupRepos(),
		ecs.NewCtyunEcsBackups(),
		ecs.NewCtyunEcsBackupPolicies(),
		vpc.NewCtyunVpcs(),
		vpc.NewCtyunSubnets(),
		vpc.NewCtyunSecurityGroups(),
		vpc.NewCtyunVpcRouteTables(),
		vpc.NewCtyunVpcRouteTableRules(),
		vpc.NewCtyunEips(),
		vpc.NewCtyunBandwidths(),
		vpce.NewCtyunVpces(),
		vpce.NewCtyunVpceServices(),
		vpce.NewCtyunVpceServiceTransitIPs(),
		vpce.NewCtyunVpceServiceReverseRules(),
		zos.NewCtyunZosBuckets(),
		zos.NewCtyunZosBucketObjects(),
		ccse.NewCtyunCcseClusters(),
		ccse.NewCtyunCcseNodePools(),
		redis.NewCtyunRedisSpecs(),
		redis.NewCtyunRedisInstances(),
		elb.NewCtyunElbHealthChecks(),
		elb.NewCtyunElbTargetGroups(),
		elb.NewCtyunElbAcls(),
		elb.NewCtyunElbTargets(),
		elb.NewElbCertificates(),
		elb.NewElbListeners(),
		elb.NewCtyunElbRules(),
		mysql2.NewCtyunMysqlInstances(),
		mysql2.NewCtyunMysqlAssociationEips(),
		mysql2.NewCtyunMysqlSpecs(),
		kafka.NewCtyunKafkaInstances(),
		kafka.NewCtyunKafkaSpecs(),
		rabbitmq.NewCtyunRabbitmqInstances(),
		rabbitmq.NewCtyunRabbitmqSpecs(),
		ccse.NewCtyunCcsePluginMarket(),
		pgsql.NewCtyunPgsqlInstances(),
		pgsql.NewCtyunPgsqlSpecs(),
		common2.NewCtyunZones(),
		mysql2.NewCtyunMysqlWhiteLists(),
		mongodb.NewCtyunMongodbInstances(),
		hpfs.NewCtyunHpfsInstances(),
		hpfs.NewCtyunHpfsClusters(),
		scaling2.NewCtyunScalings(),
		scaling2.NewCtyunScalingConfigs(),
		scaling2.NewCtyunScalingActivities(),
		scaling2.NewCtyunScalingEcsList(),
		scaling2.NewCtyunScalingPolicies(),
		sfs2.NewCtyunSfsInstances(),
		sfs2.NewCtyunSfsPermissionRules(),
		ccse.NewCtyunCcseTemplateMarket(),
		mongodb.NewCtyunMongodbSpecs(),
		mongodb.NewCtyunMongodbAssociationEips(),
	)
}

func (c *CtyunProvider) Resources(_ context.Context) []func() resource.Resource {
	return c.buildResource(
		iam.NewCtyunPolicy(),
		vpc.NewCtyunVpc(),
		vpc.NewCtyunSecurityGroup(),
		vpc.NewCtyunSecurityGroupRule(),
		vpc.NewCtyunEip(),
		vpc.NewCtyunEipAssociation(),
		vpc.NewCtyunSubnet(),
		ecs.NewCtyunEcs(),
		iam.NewCtyunIamUser(),
		iam.NewCtyunIamUserGroup(),
		iam.NewCtyunIdp(),
		ebs.NewCtyunEbs(),
		ebs.NewCtyunEbsAssociation(),
		ebs.NewCtyunEbsSnapshot(),
		ebs.NewCtyunEbsBackup(),
		//ebs.NewCtyunEbsBackupRepo(),
		ebs.NewCtyunEbsBackupPolicy(),
		ebs.NewCtyunEcsBackupPolicyBindDisks(),
		ebs.NewCtyunEbsBackupPolicyBindRepo(),
		ebs.NewCtyunEbsSnapshotPolicy(),
		ebs.NewCtyunEbsSnapshotPolicyAssociation(),
		image.NewCtyunImage(),
		image.NewCtyunImageAssociationUser(),
		ecs.NewCtyunKeypair(),
		vpc.NewCtyunBandwidth(),
		vpc.NewCtyunBandwidthAssociationEip(),
		iam.NewCtyunPolicyAssociationUserGroup(),
		iam.NewCtyunPolicyAssociationUser(),
		iam.NewCtyunEnterpriseProject(),
		iam.NewCtyunEnterpriseProjectAssociationUserGroup(),
		nat.NewCtyunNatResource(),
		nat.NewCtyunSnatResource(),
		nat.NewCtyunDnatResource(),
		ebm.NewCtyunEbm(),
		ebm.NewCtyunEbmInterface(),
		ebm.NewCtyunEbmAssociationEbs(),
		ecs.NewCtyunEcsAffinityGroup(),
		ecs.NewCtyunEcsAffinityGroupAssociation(),
		vpc.NewCtyunVpcRouteTable(),
		vpc.NewCtyunVpcRouteTableRule(),
		vpce.NewCtyunVpce(),
		vpce.NewCtyunVpceService(),
		vpce.NewCtyunVpceServiceTransitIP(),
		vpce.NewCtyunVpceServiceReverseRule(),
		vpce.NewCtyunVpceServiceConnection(),
		zos.NewCtyunZosBucket(),
		zos.NewCtyunZosBucketObject(),
		ccse.NewCtyunCcseCluster(),
		ccse.NewCtyunCcseNodePool(),
		redis.NewCtyunRedisInstance(),
		redis.NewCtyunRedisAssociationEip(),
		elb.NewCtyunElbLoadBalancer(),
		elb.NewCtyunElbHealthCheck(),
		elb.NewCtyunElbTargetGroup(),
		elb.NewCtyunElbAcl(),
		elb.NewCtyunElbTarget(),
		elb.NewCtyunElbCertificate(),
		elb.NewCtyunElbListener(),
		elb.NewCtyunElbRule(),
		mysql2.NewCtyunMysqlInstance(),
		mysql2.NewCtyunMysqlAssociationEip(),
		kafka.NewCtyunKafkaInstance(),
		ccse.NewCtyunCcsePlugin(),
		pgsql.NewCtyunPostgresqlInstance(),
		rabbitmq.NewCtyunRabbitmqInstance(),
		pgsql.NewCtyunMysqlAssociationEip(),
		mongodb.NewCtyunMongodbInstance(),
		mysql2.NewCtyunMysqlWhiteList(),
		ecs.NewCtyunEcsSnapshot(),
		hpfs.NewCtyunHpfsInstance(),
		scaling2.NewCtyunScaling(),
		scaling2.NewCtyunScalingConfig(),
		scaling2.NewCtyunScalingPolicy(),
		sfs2.NewCtyunSfs(),
		sfs2.NewCtyunSfsPermissionGroup(),
		sfs2.NewCtyunSfsPermissionGroupAssociation(),
		sfs2.NewCtyunSfsPermissionGroupRule(),
		//ecs.NewCtyunEcsBackupRepo(),
		ecs.NewCtyunEcsBackup(),
		ecs.NewCtyunEcsBackupPolicy(),
		ecs.NewCtyunEcsBackupPolicyBindInstances(),
		ecs.NewCtyunEcsBackupPolicyBindRepo(),
		ccse.NewCtyunCcseNodeAssociation(),
		//ccse.NewCtyunCcseTemplateInstance(),
		scaling2.NewCtyunScalingEcsProtection(),
		mongodb.NewCtyunMongodbAssociationEip(),
	)
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

//func (c *CtyunProvider) buildMysqlConfiguration() *mysql.Configuration {
//	// 创建一个新的 API 客户端实例
//	cfg := mysql.NewConfiguration()
//	cfg.Scheme = "https"
//	tr := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//	httpClient := &http.Client{Transport: tr}
//	cfg.HTTPClient = httpClient
//	// 设置 API 服务器的基础 URL
//	cfg.Servers = mysql.ServerConfigurations{
//		{
//			URL: "https://ctdas-global.ctapi.ctyun.cn/teledb-mysql",
//		},
//	}
//
//	return cfg
//}

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

func GetTestAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"ctyun": providerserver.NewProtocol6WithError(NewCtyunProvider("test")()),
	}
}
