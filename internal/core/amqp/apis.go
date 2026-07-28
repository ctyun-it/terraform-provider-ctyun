package amqp

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "amqp"

type Apis struct {
	AmqpInstQueryV3Api          *AmqpInstQueryV3Api
	AmqpMetadataV3Api           *AmqpMetadataV3Api
	AmqpInstanceNameV3Api       *AmqpInstanceNameV3Api
	AmqpVhostQueryV3Api         *AmqpVhostQueryV3Api
	AmqpVhostCreateV3Api        *AmqpVhostCreateV3Api
	AmqpExchangeCreateV3Api     *AmqpExchangeCreateV3Api
	AmqpExchangeQueryV3Api      *AmqpExchangeQueryV3Api
	AmqpExchangeDeleteV3Api     *AmqpExchangeDeleteV3Api
	AmqpQueueCreateV3Api        *AmqpQueueCreateV3Api
	AmqpQueueQueryV3Api         *AmqpQueueQueryV3Api
	AmqpConsumerV3Api           *AmqpConsumerV3Api
	AmqpQueueDeleteV3Api        *AmqpQueueDeleteV3Api
	AmqpBindingCreateV3Api      *AmqpBindingCreateV3Api
	AmqpBindingQueryV3Api       *AmqpBindingQueryV3Api
	AmqpSourceBindingV3Api      *AmqpSourceBindingV3Api
	AmqpQueryBinding3Api        *AmqpQueryBinding3Api
	AmqpQueryQueueBindingApi    *AmqpQueryQueueBindingApi
	AmqpQueryExchangeBindedApi  *AmqpQueryExchangeBindedApi
	AmqpQueryExchangeBindingApi *AmqpQueryExchangeBindingApi
	AmqpQueryBindingApi         *AmqpQueryBindingApi
	AmqpCreateBindingApi        *AmqpCreateBindingApi
	AmqpQueryVhostApi           *AmqpQueryVhostApi
	AmqpDeleteVhostApi          *AmqpDeleteVhostApi
	AmqpCreateVhostApi          *AmqpCreateVhostApi
	AmqpQueryExchangeApi        *AmqpQueryExchangeApi
	AmqpDeleteExchangeApi       *AmqpDeleteExchangeApi
	AmqpCreateExchangeApi       *AmqpCreateExchangeApi
	AmqpQueryQueueConsumerApi   *AmqpQueryQueueConsumerApi
	AmqpQueryQueueApi           *AmqpQueryQueueApi
	AmqpDeleteQueueApi          *AmqpDeleteQueueApi
	AmqpCreateQueueApi          *AmqpCreateQueueApi
	AmqpMetadataApi             *AmqpMetadataApi
	AmqpQueryInstanceApi        *AmqpQueryInstanceApi
	AmqpDeleteBindingApi        *AmqpDeleteBindingApi
	AmqpBindingDeleteV3Api      *AmqpBindingDeleteV3Api
	AmqpVhostDeleteV3Api        *AmqpVhostDeleteV3Api
	AmqpInstQueryDetailV3Api    *AmqpInstQueryDetailV3Api
	AmqpChangeInstanceNameApi   *AmqpChangeInstanceNameApi
	AmqpCreatePostPayOrderApi   *AmqpCreatePostPayOrderApi
	AmqpCreateOrderApi          *AmqpCreateOrderApi
	AmqpDiskExtendApi           *AmqpDiskExtendApi
	AmqpNodeExtendApi           *AmqpNodeExtendApi
	AmqpSpecExtendApi           *AmqpSpecExtendApi
	AmqpSpecShrinkApi           *AmqpSpecShrinkApi
	AmqpTransToPrePaidApi       *AmqpTransToPrePaidApi
	AmqpTransChargeTypeApi      *AmqpTransChargeTypeApi
	AmqpProdDetailApi           *AmqpProdDetailApi
	AmqpCanExpandProdApi        *AmqpCanExpandProdApi
	AmqpUnsubscribeInstApi      *AmqpUnsubscribeInstApi
	AmqpRenewApi                *AmqpRenewApi
	AmqpInstanceDeleteApi       *AmqpInstanceDeleteApi
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		AmqpInstQueryV3Api:          NewAmqpInstQueryV3Api(client),
		AmqpMetadataV3Api:           NewAmqpMetadataV3Api(client),
		AmqpInstanceNameV3Api:       NewAmqpInstanceNameV3Api(client),
		AmqpVhostQueryV3Api:         NewAmqpVhostQueryV3Api(client),
		AmqpVhostCreateV3Api:        NewAmqpVhostCreateV3Api(client),
		AmqpExchangeCreateV3Api:     NewAmqpExchangeCreateV3Api(client),
		AmqpExchangeQueryV3Api:      NewAmqpExchangeQueryV3Api(client),
		AmqpExchangeDeleteV3Api:     NewAmqpExchangeDeleteV3Api(client),
		AmqpQueueCreateV3Api:        NewAmqpQueueCreateV3Api(client),
		AmqpQueueQueryV3Api:         NewAmqpQueueQueryV3Api(client),
		AmqpConsumerV3Api:           NewAmqpConsumerV3Api(client),
		AmqpQueueDeleteV3Api:        NewAmqpQueueDeleteV3Api(client),
		AmqpBindingCreateV3Api:      NewAmqpBindingCreateV3Api(client),
		AmqpBindingQueryV3Api:       NewAmqpBindingQueryV3Api(client),
		AmqpSourceBindingV3Api:      NewAmqpSourceBindingV3Api(client),
		AmqpQueryBinding3Api:        NewAmqpQueryBinding3Api(client),
		AmqpQueryQueueBindingApi:    NewAmqpQueryQueueBindingApi(client),
		AmqpQueryExchangeBindedApi:  NewAmqpQueryExchangeBindedApi(client),
		AmqpQueryExchangeBindingApi: NewAmqpQueryExchangeBindingApi(client),
		AmqpQueryBindingApi:         NewAmqpQueryBindingApi(client),
		AmqpCreateBindingApi:        NewAmqpCreateBindingApi(client),
		AmqpQueryVhostApi:           NewAmqpQueryVhostApi(client),
		AmqpDeleteVhostApi:          NewAmqpDeleteVhostApi(client),
		AmqpCreateVhostApi:          NewAmqpCreateVhostApi(client),
		AmqpQueryExchangeApi:        NewAmqpQueryExchangeApi(client),
		AmqpDeleteExchangeApi:       NewAmqpDeleteExchangeApi(client),
		AmqpCreateExchangeApi:       NewAmqpCreateExchangeApi(client),
		AmqpQueryQueueConsumerApi:   NewAmqpQueryQueueConsumerApi(client),
		AmqpQueryQueueApi:           NewAmqpQueryQueueApi(client),
		AmqpDeleteQueueApi:          NewAmqpDeleteQueueApi(client),
		AmqpCreateQueueApi:          NewAmqpCreateQueueApi(client),
		AmqpMetadataApi:             NewAmqpMetadataApi(client),
		AmqpQueryInstanceApi:        NewAmqpQueryInstanceApi(client),
		AmqpDeleteBindingApi:        NewAmqpDeleteBindingApi(client),
		AmqpBindingDeleteV3Api:      NewAmqpBindingDeleteV3Api(client),
		AmqpVhostDeleteV3Api:        NewAmqpVhostDeleteV3Api(client),
		AmqpInstQueryDetailV3Api:    NewAmqpInstQueryDetailV3Api(client),
		AmqpChangeInstanceNameApi:   NewAmqpChangeInstanceNameApi(client),
		AmqpCreatePostPayOrderApi:   NewAmqpCreatePostPayOrderApi(client),
		AmqpCreateOrderApi:          NewAmqpCreateOrderApi(client),
		AmqpDiskExtendApi:           NewAmqpDiskExtendApi(client),
		AmqpNodeExtendApi:           NewAmqpNodeExtendApi(client),
		AmqpSpecExtendApi:           NewAmqpSpecExtendApi(client),
		AmqpSpecShrinkApi:           NewAmqpSpecShrinkApi(client),
		AmqpTransToPrePaidApi:       NewAmqpTransToPrePaidApi(client),
		AmqpTransChargeTypeApi:      NewAmqpTransChargeTypeApi(client),
		AmqpProdDetailApi:           NewAmqpProdDetailApi(client),
		AmqpCanExpandProdApi:        NewAmqpCanExpandProdApi(client),
		AmqpUnsubscribeInstApi:      NewAmqpUnsubscribeInstApi(client),
		AmqpRenewApi:                NewAmqpRenewApi(client),
		AmqpInstanceDeleteApi:       NewAmqpInstanceDeleteApi(client),
	}
}
