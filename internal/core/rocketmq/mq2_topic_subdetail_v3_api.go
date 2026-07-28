package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2TopicSubdetailV3Api
/* 查看Topic的订阅信息 */
type Mq2TopicSubdetailV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2TopicSubdetailV3Api(client *core.CtyunClient) *Mq2TopicSubdetailV3Api {
	return &Mq2TopicSubdetailV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/topic/subDetail",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2TopicSubdetailV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2TopicSubdetailV3Request) (*Mq2TopicSubdetailV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("topicName", req.TopicName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2TopicSubdetailV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2TopicSubdetailV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	TopicName  string `json:"topicName"`  /*  topic名字  */
}

type Mq2TopicSubdetailV3Response struct {
	StatusCode *string                               `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                               `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2TopicSubdetailV3ReturnObjResponse `json:"returnObj"`  /*  返回对象。  */
	Error      *string                               `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2TopicSubdetailV3ReturnObjResponse struct {
	Data *Mq2TopicSubdetailV3ReturnObjDataResponse `json:"data"` /*  主题订阅信息  */
}

type Mq2TopicSubdetailV3ReturnObjDataResponse struct {
	SubscriptionDataList []*Mq2TopicSubdetailV3ReturnObjDataSubscriptionDataListResponse `json:"subscriptionDataList"` /*  订阅数据列表  */
	Topic                *string                                                         `json:"topic"`                /*  主题名称  */
	ProdInstId           *string                                                         `json:"prodInstId"`           /*  实例id  */
}

type Mq2TopicSubdetailV3ReturnObjDataSubscriptionDataListResponse struct {
	ConsumerGroup *string `json:"consumerGroup"` /*  消费者组  */
	SubString     *string `json:"subString"`     /*  订阅表达式  */
}
