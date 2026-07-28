package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2TraceV3Api
/* 查看消息消费结果 */
type Mq2TraceV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2TraceV3Api(client *core.CtyunClient) *Mq2TraceV3Api {
	return &Mq2TraceV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/message/trace",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2TraceV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2TraceV3Request) (*Mq2TraceV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("msgId", req.MsgId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("groupName", req.GroupName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2TraceV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2TraceV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	MsgId      string `json:"msgId"`      /*  消息ID  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	GroupName  string `json:"groupName"`  /*  订阅组名字  */
}

type Mq2TraceV3Response struct {
	StatusCode *string                      `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                      `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2TraceV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                      `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2TraceV3ReturnObjResponse struct {
	MsgId             *string                                         `json:"msgId"`             /*  消息ID  */
	ClusterName       *string                                         `json:"clusterName"`       /*  集群名称  */
	ConsumeStatusList []*Mq2TraceV3ReturnObjConsumeStatusListResponse `json:"consumeStatusList"` /*  消费状态列表  */
}

type Mq2TraceV3ReturnObjConsumeStatusListResponse struct {
	Object1 *string `json:"object1"` /*  消费组名称  */
	Object2 *string `json:"object2"` /*  消费状态  */
}
