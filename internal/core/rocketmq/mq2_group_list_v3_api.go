package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2GroupListV3Api
/* 查询订阅组列表 */
type Mq2GroupListV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2GroupListV3Api(client *core.CtyunClient) *Mq2GroupListV3Api {
	return &Mq2GroupListV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/group/list",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2GroupListV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2GroupListV3Request) (*Mq2GroupListV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2GroupListV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2GroupListV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
}

type Mq2GroupListV3Response struct {
	StatusCode int32                            `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                          `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2GroupListV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                          `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2GroupListV3ReturnObjResponse struct {
	Rows  []*Mq2GroupListV3ReturnObjRowsResponse `json:"rows"`  /*  消费组信息列表  */
	Total *int32                                 `json:"total"` /*  总记录数。  */
}

type Mq2GroupListV3ReturnObjRowsResponse struct {
	GroupName                      string                                           `json:"groupName"`                      /*  消费组名称  */
	ConsumeEnable                  bool                                             `json:"consumeEnable"`                  /*  是否允许消费  */
	ConsumeFromMinEnable           bool                                             `json:"consumeFromMinEnable"`           /*  是否从最小位点消费  */
	ConsumeBroadcastEnable         bool                                             `json:"consumeBroadcastEnable"`         /*  是否允许广播消费  */
	RetryQueueNums                 int32                                            `json:"retryQueueNums"`                 /*  重试队列数量  */
	RetryMaxTimes                  int32                                            `json:"retryMaxTimes"`                  /*  最大重试次数  */
	BrokerId                       int32                                            `json:"brokerId"`                       /*  Broker ID  */
	WhichBrokerWhenConsumeSlowly   int32                                            `json:"whichBrokerWhenConsumeSlowly"`   /*  消费慢时指定的Broker ID  */
	NotifyConsumerIdsChangedEnable bool                                             `json:"notifyConsumerIdsChangedEnable"` /*  是否开启消费者ID变更通知  */
	PullMechanism                  int32                                            `json:"pullMechanism"`                  /*  拉取机制  */
	FirstConsumeMechanism          int32                                            `json:"firstConsumeMechanism"`          /*  首次消费机制  */
	SubscribeMap                   *Mq2GroupListV3ReturnObjRowsSubscribeMapResponse `json:"subscribeMap"`                   /*  订阅关系  */
	Remark                         string                                           `json:"remark"`                         /*  备注  */
	BrokerName                     string                                           `json:"brokerName"`                     /*  Broker名称  */
	RawType                        int32                                            `json:"type"`                           /*  类型  */
	AllowdAllTopics                bool                                             `json:"allowdAllTopics"`                /*  是否允许所有主题  */
}

type Mq2GroupListV3ReturnObjRowsSubscribeMapResponse struct{}
