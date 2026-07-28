package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2ConsumerAccumulateV3Api
/* 查询订阅组的消息消费堆积 */
type Mq2ConsumerAccumulateV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2ConsumerAccumulateV3Api(client *core.CtyunClient) *Mq2ConsumerAccumulateV3Api {
	return &Mq2ConsumerAccumulateV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/consumer/accumulate",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2ConsumerAccumulateV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2ConsumerAccumulateV3Request) (*Mq2ConsumerAccumulateV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("groupName", req.GroupName)
	ctReq.AddParam("topicName", req.TopicName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2ConsumerAccumulateV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2ConsumerAccumulateV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	GroupName  string `json:"groupName"`  /*  订阅组名字  */
	TopicName  string `json:"topicName"`  /*  Topic名字  */
}

type Mq2ConsumerAccumulateV3Response struct {
	StatusCode *string                                   `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                                   `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2ConsumerAccumulateV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                                   `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2ConsumerAccumulateV3ReturnObjResponse struct {
	Total *int32                                          `json:"total"` /*  总记录数。  */
	Rows  []*Mq2ConsumerAccumulateV3ReturnObjRowsResponse `json:"rows"`  /*  消费组消息统计信息列表。  */
}

type Mq2ConsumerAccumulateV3ReturnObjRowsResponse struct {
	Diff              *int32  `json:"diff"`              /*  差值  */
	Total             *int32  `json:"total"`             /*  消息总数  */
	ConsumedTotal     *int32  `json:"consumedTotal"`     /*  已消费消息总数  */
	RetryNums         *int32  `json:"retryNums"`         /*  重试次数  */
	GroupName         *string `json:"groupName"`         /*  消费组名称  */
	BrokerName        *string `json:"brokerName"`        /*  Broker名称  */
	TopicName         *string `json:"topicName"`         /*  主题名称  */
	OutTps            *int32  `json:"outTps"`            /*  出消息TPS  */
	OutTotalMsgToday  *int32  `json:"outTotalMsgToday"`  /*  今日出消息总数  */
	OffSetDiff        *int32  `json:"offSetDiff"`        /*  偏移量差值  */
	NotAckMsgNum      *int32  `json:"notAckMsgNum"`      /*  未确认消息数  */
	NotAckMsgTime     *int32  `json:"notAckMsgTime"`     /*  未确认消息时长  */
	NotConsumeMsgTime *int32  `json:"notConsumeMsgTime"` /*  未消费消息时长  */
}
