package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2GroupOutputtpsV3Api
/* 查询订阅组一段时间内消费tps信息 */
type Mq2GroupOutputtpsV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2GroupOutputtpsV3Api(client *core.CtyunClient) *Mq2GroupOutputtpsV3Api {
	return &Mq2GroupOutputtpsV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/group/outputTps",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2GroupOutputtpsV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2GroupOutputtpsV3Request) (*Mq2GroupOutputtpsV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("topicName", req.TopicName)
	ctReq.AddParam("brokerName", req.BrokerName)
	ctReq.AddParam("groupName", req.GroupName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2GroupOutputtpsV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2GroupOutputtpsV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	TopicName  string `json:"topicName"`  /*  topic名字  */
	BrokerName string `json:"brokerName"` /*  Broker名字  */
	GroupName  string `json:"groupName"`  /*  订阅组名字  */
}

type Mq2GroupOutputtpsV3Response struct {
	StatusCode *string                               `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                               `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2GroupOutputtpsV3ReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例'里面的注释  */
	Error      *string                               `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2GroupOutputtpsV3ReturnObjResponse struct {
	Data *Mq2GroupOutputtpsV3ReturnObjDataResponse `json:"data"` /*  主题TPS统计数据  */
}

type Mq2GroupOutputtpsV3ReturnObjDataResponse struct {
	ProdInstId *string                                              `json:"prodInstId"` /*  实例id  */
	TopicName  *string                                              `json:"topicName"`  /*  主题名称  */
	GroupName  *string                                              `json:"groupName"`  /*  消费组名称  */
	TpsVoList  []*Mq2GroupOutputtpsV3ReturnObjDataTpsVoListResponse `json:"tpsVoList"`  /*  TPS统计列表  */
}

type Mq2GroupOutputtpsV3ReturnObjDataTpsVoListResponse struct {
	Tps     *float64 `json:"tps"`     /*  消息吞吐量（条/秒）  */
	CrtTime *int64   `json:"crtTime"` /*  统计时间戳（毫秒）  */
}
