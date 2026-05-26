package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2TopicInputtpsV3Api
/* 查询主题一段时间内消息写入tps信息 */
type Mq2TopicInputtpsV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2TopicInputtpsV3Api(client *core.CtyunClient) *Mq2TopicInputtpsV3Api {
	return &Mq2TopicInputtpsV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/topic/inputTps",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2TopicInputtpsV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2TopicInputtpsV3Request) (*Mq2TopicInputtpsV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("topicName", req.TopicName)
	ctReq.AddParam("brokerName", req.BrokerName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2TopicInputtpsV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2TopicInputtpsV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	TopicName  string `json:"topicName"`  /*  Topic名字  */
	BrokerName string `json:"brokerName"` /*  Broker名字  */
}

type Mq2TopicInputtpsV3Response struct {
	StatusCode *string                              `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                              `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2TopicInputtpsV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                              `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2TopicInputtpsV3ReturnObjResponse struct {
	Data *Mq2TopicInputtpsV3ReturnObjDataResponse `json:"data"` /*  主题TPS统计数据  */
}

type Mq2TopicInputtpsV3ReturnObjDataResponse struct {
	ProdInstId *string                                             `json:"prodInstId"` /*  实例id  */
	TopicName  *string                                             `json:"topicName"`  /*  主题名称  */
	TpsVoList  []*Mq2TopicInputtpsV3ReturnObjDataTpsVoListResponse `json:"tpsVoList"`  /*  TPS统计列表  */
}

type Mq2TopicInputtpsV3ReturnObjDataTpsVoListResponse struct {
	Tps     *float64 `json:"tps"`     /*  消息吞吐量（条/秒）  */
	CrtTime *int64   `json:"crtTime"` /*  统计时间戳（毫秒）  */
}
