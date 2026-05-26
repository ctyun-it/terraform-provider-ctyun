package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2MessageGetTraceApi
/* 根据msgId查询消息消息轨迹 */
type Mq2MessageGetTraceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2MessageGetTraceApi(client *core.CtyunClient) *Mq2MessageGetTraceApi {
	return &Mq2MessageGetTraceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/message/getTrace",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2MessageGetTraceApi) Do(ctx context.Context, credential core.Credential, req *Mq2MessageGetTraceRequest) (*Mq2MessageGetTraceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("msgId", req.MsgId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2MessageGetTraceResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2MessageGetTraceRequest struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID   */
	MsgId      string `json:"msgId"`      /*  消息ID  */
}

type Mq2MessageGetTraceResponse struct {
	StatusCode *string                              `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                              `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2MessageGetTraceReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例'里面的注释  */
	Error      *string                              `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2MessageGetTraceReturnObjResponse struct {
	Data *Mq2MessageGetTraceReturnObjDataResponse `json:"data"` /*  消息生产消费轨迹数据对象  */
}

type Mq2MessageGetTraceReturnObjDataResponse struct {
	Topic             *string                                                `json:"topic"`             /*  消息主题名  */
	Tags              *string                                                `json:"tags"`              /*  标签  */
	Keys              *string                                                `json:"keys"`              /*  消息索引值  */
	GroupName         *string                                                `json:"groupName"`         /*  生产者名  */
	MsgType           *string                                                `json:"msgType"`           /*  消息类型  */
	MsgId             *string                                                `json:"msgId"`             /*  消息id  */
	CostTime          *string                                                `json:"costTime"`          /*  耗时  */
	Status            *string                                                `json:"status"`            /*  发送状态  */
	PubTime           *string                                                `json:"pubTime"`           /*  发送时间戳  */
	ArrivedServerTime *string                                                `json:"arrivedServerTime"` /*  消息到达服务端时间戳  */
	ClientHost        *string                                                `json:"clientHost"`        /*  客户端主机地址  */
	SubTraceList      []*Mq2MessageGetTraceReturnObjDataSubTraceListResponse `json:"subTraceList"`      /*  消息堆积列表  */
}

type Mq2MessageGetTraceReturnObjDataSubTraceListResponse struct {
	GroupName    *string                                                            `json:"groupName"`    /*  消费组名  */
	TraceLogList []*Mq2MessageGetTraceReturnObjDataSubTraceListTraceLogListResponse `json:"traceLogList"` /*  消费轨迹日志列表  */
}

type Mq2MessageGetTraceReturnObjDataSubTraceListTraceLogListResponse struct {
	MsgId      *string `json:"msgId"`      /*  消息id  */
	ClientHost *string `json:"clientHost"` /*  客户端主机地址  */
	CostTime   *string `json:"costTime"`   /*  消费耗时（毫秒）  */
	SubTime    *string `json:"subTime"`    /*  消费时间戳  */
	Status     *string `json:"status"`     /*  消费状态  CONSUME_SUCCESS-消费成功;RECONSUME_LATER - 稍后重试  */
}
