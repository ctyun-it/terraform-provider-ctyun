package mongodb

import ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"

type Apis struct {
	MongodbDestroyApi                 *MongodbDestroyApi
	MongodbCreateApi                  *MongodbCreateApi
	MongodbGetListApi                 *MongodbGetListApi
	MongodbQueryDetailApi             *MongodbQueryDetailApi
	MongodbRefundApi                  *MongodbRefundApi
	MongodbUpgradeApi                 *MongodbUpgradeApi
	MongodbBindEipApi                 *MongodbBindEipApi
	MongodbUnbindEipApi               *MongodbUnbindEipApi
	MongodbUpdateSecurityGroupApi     *MongodbUpdateSecurityGroupApi
	MongodbUpdateInstanceNameApi      *MongodbUpdateInstanceNameApi
	MongodbUpdatePortApi              *MongodbUpdatePortApi
	MongodbBoundEipListApi            *MongodbBoundEipListApi
	TeledbGetAvailabilityZone         *TeledbGetAvailabilityZone
	MongodbDescribeAccountsApi        *MongodbDescribeAccountsApi
	MongodbCreateAccountApi           *MongodbCreateAccountApi
	MongodbUpdateAccountPasswordApi   *MongodbUpdateAccountPasswordApi
	MongodbDeleteAccountApi           *MongodbDeleteAccountApi
	MongodbModifyAccountPermissionApi *MongodbModifyAccountPermissionApi
	MongodbCreateIpWhitelistApi       *MongodbCreateIpWhitelistApi
	MongodbUpdateIpWhitelistApi       *MongodbUpdateIpWhitelistApi
	MongodbDeleteIpWhitelistApi       *MongodbDeleteIpWhitelistApi
	MongodbDescribeIpWhitelistApi     *MongodbDescribeIpWhitelistApi
	MongodbCreateBackupApi            *MongodbCreateBackupApi            // 添加手动备份API
	MongodbDeleteBackupApi            *MongodbDeleteBackupApi            // 添加删除备份API
	MongodbDescribeBackupsApi         *MongodbDescribeBackupsApi         // 添加查询备份列表API
	MongodbCreateParamTemplateApi     *MongodbCreateParamTemplateApi     // 添加创建参数组API
	MongodbDeleteParamTemplateApi     *MongodbDeleteParamTemplateApi     // 添加删除参数组API
	MongodbUpdateParamTemplateDescApi *MongodbUpdateParamTemplateDescApi // 添加修改参数组描述API
	MongodbDescribeParamTemplatesApi  *MongodbDescribeParamTemplatesApi  // 添加查询参数组列表API
	MongodbUpdatePasswordApi          *MongodbUpdatePasswordApi          // 重置实例中root账号的密码
	MongodbQueryInstNodesApi          *MongodbQueryInstNodesApi          // 查询实例节点信息
	MongodbRestartDbApi               *MongodbRestartDbApi               // 重启实例
	MongodbGetIDByOrderApi            *MongodbGetIDByOrderApi
}

func NewApis(client *ctyunsdk.CtyunClient) *Apis {
	builder := ctyunsdk.NewApiHookBuilder()
	for _, hook := range client.Config.ApiHooks {
		builder.AddHooks(hook)
	}

	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointMongodbTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointMongodbTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentProd, EndPointMongodbProd)
	return &Apis{
		MongodbDestroyApi:                 NewMongodbDestroyApi(client),
		MongodbCreateApi:                  NewMongodbCreateApi(client),
		MongodbGetListApi:                 NewMongodbGetListApi(client),
		MongodbQueryDetailApi:             NewMongodbQueryDetailApi(client),
		MongodbRefundApi:                  NewMongodbRefundApi(client),
		MongodbUpgradeApi:                 NewMongodbUpgradeApi(client),
		MongodbBindEipApi:                 NewMongodbBindEipApi(client),
		MongodbUnbindEipApi:               NewMongodbUnbindEipApi(client),
		MongodbUpdateSecurityGroupApi:     NewMongodbUpdateSecurityGroupApi(client),
		MongodbUpdateInstanceNameApi:      NewMongodbUpdateInstanceNameApi(client),
		MongodbUpdatePortApi:              NewMongodbUpdatePortApi(client),
		MongodbBoundEipListApi:            NewMongodbBoundEipListApi(client),
		TeledbGetAvailabilityZone:         NewTeledbGetAvailabilityZone(client),
		MongodbDescribeAccountsApi:        NewMongodbDescribeAccountsApi(client),
		MongodbCreateAccountApi:           NewMongodbCreateAccountApi(client),
		MongodbUpdateAccountPasswordApi:   NewMongodbUpdateAccountPasswordApi(client),
		MongodbDeleteAccountApi:           NewMongodbDeleteAccountApi(client),
		MongodbModifyAccountPermissionApi: NewMongodbModifyAccountPermissionApi(client),
		MongodbCreateIpWhitelistApi:       NewMongodbCreateIpWhitelistApi(client),
		MongodbUpdateIpWhitelistApi:       NewMongodbUpdateIpWhitelistApi(client),
		MongodbDeleteIpWhitelistApi:       NewMongodbDeleteIpWhitelistApi(client),
		MongodbDescribeIpWhitelistApi:     NewMongodbDescribeIpWhitelistApi(client),
		MongodbCreateBackupApi:            NewMongodbCreateBackupApi(client),            // 初始化手动备份API
		MongodbDeleteBackupApi:            NewMongodbDeleteBackupApi(client),            // 初始化删除备份API
		MongodbDescribeBackupsApi:         NewMongodbDescribeBackupsApi(client),         // 初始化查询备份列表API
		MongodbCreateParamTemplateApi:     NewMongodbCreateParamTemplateApi(client),     // 初始化创建参数组API
		MongodbDeleteParamTemplateApi:     NewMongodbDeleteParamTemplateApi(client),     // 初始化删除参数组API
		MongodbUpdateParamTemplateDescApi: NewMongodbUpdateParamTemplateDescApi(client), // 初始化修改参数组描述API
		MongodbDescribeParamTemplatesApi:  NewMongodbDescribeParamTemplatesApi(client),  // 初始化查询参数组列表API
		MongodbUpdatePasswordApi:          NewMongodbUpdatePasswordApi(client),          // 初始化重置实例中root账号的密码API
		MongodbQueryInstNodesApi:          NewMongodbQueryInstNodesApi(client),          // 初始化查询实例节点信息API
		MongodbRestartDbApi:               NewMongodbRestartDbApi(client),               // 初始化重启实例API
		MongodbGetIDByOrderApi:            NewMongodbGetIDByOrderApi(client),
	}
}
