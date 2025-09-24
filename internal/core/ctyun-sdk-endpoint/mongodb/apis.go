package mongodb

import ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"

type Apis struct {
	MongodbDestroyApi             *MongodbDestroyApi
	MongodbCreateApi              *MongodbCreateApi
	MongodbGetListApi             *MongodbGetListApi
	MongodbQueryDetailApi         *MongodbQueryDetailApi
	MongodbRefundApi              *MongodbRefundApi
	MongodbUpgradeApi             *MongodbUpgradeApi
	MongodbBindEipApi             *MongodbBindEipApi
	MongodbUnbindEipApi           *MongodbUnbindEipApi
	MongodbUpdateSecurityGroupApi *MongodbUpdateSecurityGroupApi
	MongodbUpdateInstanceNameApi  *MongodbUpdateInstanceNameApi
	MongodbUpdatePortApi          *MongodbUpdatePortApi
	MongodbBoundEipListApi        *MongodbBoundEipListApi
	TeledbGetAvailabilityZone     *TeledbGetAvailabilityZone
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
		MongodbDestroyApi:             NewMongodbDestroyApi(client),
		MongodbCreateApi:              NewMongodbCreateApi(client),
		MongodbGetListApi:             NewMongodbGetListApi(client),
		MongodbQueryDetailApi:         NewMongodbQueryDetailApi(client),
		MongodbRefundApi:              NewMongodbRefundApi(client),
		MongodbUpgradeApi:             NewMongodbUpgradeApi(client),
		MongodbBindEipApi:             NewMongodbBindEipApi(client),
		MongodbUnbindEipApi:           NewMongodbUnbindEipApi(client),
		MongodbUpdateSecurityGroupApi: NewMongodbUpdateSecurityGroupApi(client),
		MongodbUpdateInstanceNameApi:  NewMongodbUpdateInstanceNameApi(client),
		MongodbUpdatePortApi:          NewMongodbUpdatePortApi(client),
		MongodbBoundEipListApi:        NewMongodbBoundEipListApi(client),
		TeledbGetAvailabilityZone:     NewTeledbGetAvailabilityZone(client),
	}
}
