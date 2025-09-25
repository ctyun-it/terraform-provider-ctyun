package amqp

import (
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

type Apis struct {
	AmqpInstancesQueryProdApi          *AmqpInstancesQueryProdApi
	AmqpInstancesQueryApi              *AmqpInstancesQueryApi
	AmqpInstancesCreatePostPayOrderApi *AmqpInstancesCreatePostPayOrderApi
	AmqpInstancesCreatePrePayOrderApi  *AmqpInstancesCreatePrePayOrderApi
	AmqpInstancesDiskExtendApi         *AmqpInstancesDiskExtendApi
	AmqpInstancesNodeExtendApi         *AmqpInstancesNodeExtendApi
	AmqpInstancesSpecExtendApi         *AmqpInstancesSpecExtendApi
	AmqpInstancesUnsubscribeInstApi    *AmqpInstancesUnsubscribeInstApi
	AmqpInstancesInstanceNameApi       *AmqpInstancesInstanceNameApi
	AmqpInstancesQueryDetailApi        *AmqpInstancesQueryDetailApi
	AmqpInstanceDeleteApi              *AmqpInstanceDeleteApi
	AmqpProdDetailApi                  *AmqpProdDetailApi
}

func NewApis(client *ctyunsdk.CtyunClient) *Apis {
	builder := ctyunsdk.NewApiHookBuilder()
	for _, hook := range client.Config.ApiHooks {
		builder.AddHooks(hook)
	}

	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentProd, EndPointProd)
	return &Apis{
		AmqpInstancesQueryProdApi:          NewAmqpInstancesQueryProdApi(client),
		AmqpInstancesQueryApi:              NewAmqpInstancesQueryApi(client),
		AmqpInstancesCreatePostPayOrderApi: NewAmqpInstancesCreatePostPayOrderApi(client),
		AmqpInstancesCreatePrePayOrderApi:  NewAmqpInstancesCreatePrePayOrderApi(client),
		AmqpInstancesDiskExtendApi:         NewAmqpInstancesDiskExtendApi(client),
		AmqpInstancesNodeExtendApi:         NewAmqpInstancesNodeExtendApi(client),
		AmqpInstancesSpecExtendApi:         NewAmqpInstancesSpecExtendApi(client),
		AmqpInstancesUnsubscribeInstApi:    NewAmqpInstancesUnsubscribeInstApi(client),
		AmqpInstancesInstanceNameApi:       NewAmqpInstancesInstanceNameApi(client),
		AmqpInstancesQueryDetailApi:        NewAmqpInstancesQueryDetailApi(client),
		AmqpInstanceDeleteApi:              NewAmqpInstanceDeleteApi(client),
		AmqpProdDetailApi:                  NewAmqpProdDetailApi(client),
	}
}
