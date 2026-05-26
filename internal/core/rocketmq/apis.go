package rocketmq

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "mq2"

type Apis struct {
	RocketmqInstQueryV3Api *RocketmqInstQueryV3Api

	RocketmqInstanceNameV3Api *RocketmqInstanceNameV3Api

	RocketmqSourceBindingV3Api *RocketmqSourceBindingV3Api

	RocketmqQueryInstanceApi *RocketmqQueryInstanceApi

	RocketmqInstQueryDetailV3Api *RocketmqInstQueryDetailV3Api

	RocketmqDiskExtendApi      *RocketmqDiskExtendApi
	RocketmqNodeExtendApi      *RocketmqNodeExtendApi
	RocketmqSpecExtendApi      *RocketmqSpecExtendApi
	RocketmqSpecShrinkApi      *RocketmqSpecShrinkApi
	RocketmqTransToPrePaidApi  *RocketmqTransToPrePaidApi
	RocketmqTransChargeTypeApi *RocketmqTransChargeTypeApi
	RocketmqProdDetailApi      *RocketmqProdDetailApi

	RocketmqUnsubscribeInstApi            *RocketmqUnsubscribeInstApi
	RocketmqRenewApi                      *RocketmqRenewApi
	RocketmqInstanceDeleteApi             *RocketmqInstanceDeleteApi
	RocketmqCreatePostPayOrderApi         *RocketmqCreatePostPayOrderApi
	Mq2UnsubscribeInstApi                 *Mq2UnsubscribeInstApi
	Mq2QueryDlqbyTimeApi                  *Mq2QueryDlqbyTimeApi
	Mq2QueryDlqbyIdApi                    *Mq2QueryDlqbyIdApi
	Mq2TopicResetTimespanApi              *Mq2TopicResetTimespanApi
	Mq2ResetOffsetApi                     *Mq2ResetOffsetApi
	Mq2PushApi                            *Mq2PushApi
	Mq2UpdateGroupPermApi                 *Mq2UpdateGroupPermApi
	Mq2DeleteGroupApi                     *Mq2DeleteGroupApi
	Mq2CreateGroupApi                     *Mq2CreateGroupApi
	Mq2SpuInstApi                         *Mq2SpuInstApi
	Mq2DeleteTopicApi                     *Mq2DeleteTopicApi
	Mq2CreateTopicApi                     *Mq2CreateTopicApi
	Mq2UpdateInstanceNameApi              *Mq2UpdateInstanceNameApi
	Mq2ProdResPoolsApi                    *Mq2ProdResPoolsApi
	Mq2AvailableRegionApi                 *Mq2AvailableRegionApi
	Mq2AccumulateInfoApi                  *Mq2AccumulateInfoApi
	Mq2CreateInstApi                      *Mq2CreateInstApi
	Mq2InstanceQueryPriceApi              *Mq2InstanceQueryPriceApi
	Mq2QueryPriceForNodeExtendApi         *Mq2QueryPriceForNodeExtendApi
	Mq2InstanceQueryPriceForDiskExtendApi *Mq2InstanceQueryPriceForDiskExtendApi
	Mq2DiskExtendApi                      *Mq2DiskExtendApi
	Mq2QueryPriceForSpecExtendApi         *Mq2QueryPriceForSpecExtendApi
	Mq2SpecExtendApi                      *Mq2SpecExtendApi
	Mq2QueryPriceForRenewApi              *Mq2QueryPriceForRenewApi
	Mq2InstanceUnsubscribeInstApi         *Mq2InstanceUnsubscribeInstApi
	Mq2DlqResendApi                       *Mq2DlqResendApi
	Mq2TopicCreateV3Api                   *Mq2TopicCreateV3Api
	Mq2TopicUpdateApi                     *Mq2TopicUpdateApi
	Mq2TopicDeleteV3Api                   *Mq2TopicDeleteV3Api
	Mq2GroupCreateV3Api                   *Mq2GroupCreateV3Api
	Mq2GroupOutputtpsV3Api                *Mq2GroupOutputtpsV3Api
	Mq2GroupUpdatepermV3Api               *Mq2GroupUpdatepermV3Api
	Mq2ConsumerResetoffsetV3Api           *Mq2ConsumerResetoffsetV3Api
	Mq2GroupDeleteV3Api                   *Mq2GroupDeleteV3Api
	Mq2InstanceQueryProdApi               *Mq2InstanceQueryProdApi
	Mq2InstanceUpdatenameV3Api            *Mq2InstanceUpdatenameV3Api
	Mq2MessagePushV3Api                   *Mq2MessagePushV3Api
	Mq2TopicStatusV3Api                   *Mq2TopicStatusV3Api
	Mq2BaseInfoV3Api                      *Mq2BaseInfoV3Api
	Mq2InstanceInfoApi                    *Mq2InstanceInfoApi
	Mq2InstanceListApi                    *Mq2InstanceListApi
	Mq2GroupSubDetailV3Api                *Mq2GroupSubDetailV3Api
	Mq2GroupListV3Api                     *Mq2GroupListV3Api
	Mq2GroupSubDetailApi                  *Mq2GroupSubDetailApi
	Mq2ConsumerGroupOutputtpsApi          *Mq2ConsumerGroupOutputtpsApi
	Mq2GroupListApi                       *Mq2GroupListApi
	Mq2MessageGetTraceApi                 *Mq2MessageGetTraceApi
	Mq2QueryByTimeV3Api                   *Mq2QueryByTimeV3Api
	Mq2QueryByKeyV3Api                    *Mq2QueryByKeyV3Api
	Mq2QueryByMsgIdV3Api                  *Mq2QueryByMsgIdV3Api
	Mq2TraceV3Api                         *Mq2TraceV3Api
	Mq2QueryByTimeApi                     *Mq2QueryByTimeApi
	Mq2QueryByKeyApi                      *Mq2QueryByKeyApi
	Mq2QueryByMsgIdApi                    *Mq2QueryByMsgIdApi
	Mq2TraceApi                           *Mq2TraceApi
	Mq2TopicInputtpsV3Api                 *Mq2TopicInputtpsV3Api
	Mq2TopicSubdetailV3Api                *Mq2TopicSubdetailV3Api
	Mq2TopicListV3Api                     *Mq2TopicListV3Api
	Mq2TopicInputTpsApi                   *Mq2TopicInputTpsApi
	Mq2TopicSubDetailApi                  *Mq2TopicSubDetailApi
	Mq2TopicListApi                       *Mq2TopicListApi
	Mq2TopicStatusApi                     *Mq2TopicStatusApi
	Mq2CreatePostPayOrderApi              *Mq2CreatePostPayOrderApi
	Mq2ConsumerAccumulateV3Api            *Mq2ConsumerAccumulateV3Api
	Mq2ConsumerTimeSpanV3Api              *Mq2ConsumerTimeSpanV3Api
	Mq2ConsumerConnectionV3Api            *Mq2ConsumerConnectionV3Api
	Mq2ConsumerStatusV3Api                *Mq2ConsumerStatusV3Api
	Mq2GetByTimeV3Api                     *Mq2GetByTimeV3Api
	Mq2QueryDlqbyIdV3Api                  *Mq2QueryDlqbyIdV3Api
	Mq2NodeExtendApi                      *Mq2NodeExtendApi
	Mq2QueryPriceForUnsubscribeApi        *Mq2QueryPriceForUnsubscribeApi
	Mq2ListV3Api                          *Mq2ListV3Api
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		RocketmqInstQueryV3Api:                NewRocketmqInstQueryV3Api(client),
		RocketmqInstanceNameV3Api:             NewRocketmqInstanceNameV3Api(client),
		RocketmqSourceBindingV3Api:            NewRocketmqSourceBindingV3Api(client),
		RocketmqQueryInstanceApi:              NewRocketmqQueryInstanceApi(client),
		RocketmqInstQueryDetailV3Api:          NewRocketmqInstQueryDetailV3Api(client),
		RocketmqCreatePostPayOrderApi:         NewRocketmqCreatePostPayOrderApi(client),
		RocketmqDiskExtendApi:                 NewRocketmqDiskExtendApi(client),
		RocketmqNodeExtendApi:                 NewRocketmqNodeExtendApi(client),
		RocketmqSpecExtendApi:                 NewRocketmqSpecExtendApi(client),
		RocketmqSpecShrinkApi:                 NewRocketmqSpecShrinkApi(client),
		RocketmqTransToPrePaidApi:             NewRocketmqTransToPrePaidApi(client),
		RocketmqTransChargeTypeApi:            NewRocketmqTransChargeTypeApi(client),
		RocketmqProdDetailApi:                 NewRocketmqProdDetailApi(client),
		RocketmqUnsubscribeInstApi:            NewRocketmqUnsubscribeInstApi(client),
		RocketmqRenewApi:                      NewRocketmqRenewApi(client),
		RocketmqInstanceDeleteApi:             NewRocketmqInstanceDeleteApi(client),
		Mq2UnsubscribeInstApi:                 NewMq2UnsubscribeInstApi(client),
		Mq2QueryDlqbyTimeApi:                  NewMq2QueryDlqbyTimeApi(client),
		Mq2QueryDlqbyIdApi:                    NewMq2QueryDlqbyIdApi(client),
		Mq2TopicResetTimespanApi:              NewMq2TopicResetTimespanApi(client),
		Mq2ResetOffsetApi:                     NewMq2ResetOffsetApi(client),
		Mq2PushApi:                            NewMq2PushApi(client),
		Mq2UpdateGroupPermApi:                 NewMq2UpdateGroupPermApi(client),
		Mq2DeleteGroupApi:                     NewMq2DeleteGroupApi(client),
		Mq2CreateGroupApi:                     NewMq2CreateGroupApi(client),
		Mq2SpuInstApi:                         NewMq2SpuInstApi(client),
		Mq2DeleteTopicApi:                     NewMq2DeleteTopicApi(client),
		Mq2CreateTopicApi:                     NewMq2CreateTopicApi(client),
		Mq2UpdateInstanceNameApi:              NewMq2UpdateInstanceNameApi(client),
		Mq2ProdResPoolsApi:                    NewMq2ProdResPoolsApi(client),
		Mq2AvailableRegionApi:                 NewMq2AvailableRegionApi(client),
		Mq2AccumulateInfoApi:                  NewMq2AccumulateInfoApi(client),
		Mq2CreateInstApi:                      NewMq2CreateInstApi(client),
		Mq2InstanceQueryPriceApi:              NewMq2InstanceQueryPriceApi(client),
		Mq2QueryPriceForNodeExtendApi:         NewMq2QueryPriceForNodeExtendApi(client),
		Mq2InstanceQueryPriceForDiskExtendApi: NewMq2InstanceQueryPriceForDiskExtendApi(client),
		Mq2DiskExtendApi:                      NewMq2DiskExtendApi(client),
		Mq2QueryPriceForSpecExtendApi:         NewMq2QueryPriceForSpecExtendApi(client),
		Mq2SpecExtendApi:                      NewMq2SpecExtendApi(client),
		Mq2QueryPriceForRenewApi:              NewMq2QueryPriceForRenewApi(client),
		Mq2InstanceUnsubscribeInstApi:         NewMq2InstanceUnsubscribeInstApi(client),
		Mq2DlqResendApi:                       NewMq2DlqResendApi(client),
		Mq2TopicCreateV3Api:                   NewMq2TopicCreateV3Api(client),
		Mq2TopicUpdateApi:                     NewMq2TopicUpdateApi(client),
		Mq2TopicDeleteV3Api:                   NewMq2TopicDeleteV3Api(client),
		Mq2GroupCreateV3Api:                   NewMq2GroupCreateV3Api(client),
		Mq2GroupOutputtpsV3Api:                NewMq2GroupOutputtpsV3Api(client),
		Mq2GroupUpdatepermV3Api:               NewMq2GroupUpdatepermV3Api(client),
		Mq2ConsumerResetoffsetV3Api:           NewMq2ConsumerResetoffsetV3Api(client),
		Mq2GroupDeleteV3Api:                   NewMq2GroupDeleteV3Api(client),
		Mq2InstanceQueryProdApi:               NewMq2InstanceQueryProdApi(client),
		Mq2InstanceUpdatenameV3Api:            NewMq2InstanceUpdatenameV3Api(client),
		Mq2MessagePushV3Api:                   NewMq2MessagePushV3Api(client),
		Mq2TopicStatusV3Api:                   NewMq2TopicStatusV3Api(client),
		Mq2BaseInfoV3Api:                      NewMq2BaseInfoV3Api(client),
		Mq2InstanceInfoApi:                    NewMq2InstanceInfoApi(client),
		Mq2InstanceListApi:                    NewMq2InstanceListApi(client),
		Mq2GroupSubDetailV3Api:                NewMq2GroupSubDetailV3Api(client),
		Mq2GroupListV3Api:                     NewMq2GroupListV3Api(client),
		Mq2GroupSubDetailApi:                  NewMq2GroupSubDetailApi(client),
		Mq2ConsumerGroupOutputtpsApi:          NewMq2ConsumerGroupOutputtpsApi(client),
		Mq2GroupListApi:                       NewMq2GroupListApi(client),
		Mq2MessageGetTraceApi:                 NewMq2MessageGetTraceApi(client),
		Mq2QueryByTimeV3Api:                   NewMq2QueryByTimeV3Api(client),
		Mq2QueryByKeyV3Api:                    NewMq2QueryByKeyV3Api(client),
		Mq2QueryByMsgIdV3Api:                  NewMq2QueryByMsgIdV3Api(client),
		Mq2TraceV3Api:                         NewMq2TraceV3Api(client),
		Mq2QueryByTimeApi:                     NewMq2QueryByTimeApi(client),
		Mq2QueryByKeyApi:                      NewMq2QueryByKeyApi(client),
		Mq2QueryByMsgIdApi:                    NewMq2QueryByMsgIdApi(client),
		Mq2TraceApi:                           NewMq2TraceApi(client),
		Mq2TopicInputtpsV3Api:                 NewMq2TopicInputtpsV3Api(client),
		Mq2TopicSubdetailV3Api:                NewMq2TopicSubdetailV3Api(client),
		Mq2TopicListV3Api:                     NewMq2TopicListV3Api(client),
		Mq2TopicInputTpsApi:                   NewMq2TopicInputTpsApi(client),
		Mq2TopicSubDetailApi:                  NewMq2TopicSubDetailApi(client),
		Mq2TopicListApi:                       NewMq2TopicListApi(client),
		Mq2TopicStatusApi:                     NewMq2TopicStatusApi(client),
		Mq2CreatePostPayOrderApi:              NewMq2CreatePostPayOrderApi(client),
		Mq2ConsumerAccumulateV3Api:            NewMq2ConsumerAccumulateV3Api(client),
		Mq2ConsumerTimeSpanV3Api:              NewMq2ConsumerTimeSpanV3Api(client),
		Mq2ConsumerConnectionV3Api:            NewMq2ConsumerConnectionV3Api(client),
		Mq2ConsumerStatusV3Api:                NewMq2ConsumerStatusV3Api(client),
		Mq2GetByTimeV3Api:                     NewMq2GetByTimeV3Api(client),
		Mq2QueryDlqbyIdV3Api:                  NewMq2QueryDlqbyIdV3Api(client),
		Mq2NodeExtendApi:                      NewMq2NodeExtendApi(client),
		Mq2QueryPriceForUnsubscribeApi:        NewMq2QueryPriceForUnsubscribeApi(client),
		Mq2ListV3Api:                          NewMq2ListV3Api(client),
	}
}
