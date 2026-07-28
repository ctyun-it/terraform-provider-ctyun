package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2QueryByMsgIdV3Api
/* 根据msgId查询消息详情 */
type Mq2QueryByMsgIdV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2QueryByMsgIdV3Api(client *core.CtyunClient) *Mq2QueryByMsgIdV3Api {
	return &Mq2QueryByMsgIdV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/message/queryByMsgId",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2QueryByMsgIdV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2QueryByMsgIdV3Request) (*Mq2QueryByMsgIdV3Response, error) {
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
	var resp Mq2QueryByMsgIdV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2QueryByMsgIdV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	MsgId      string `json:"msgId"`      /*  消息ID  */
}

type Mq2QueryByMsgIdV3Response struct {
	StatusCode *string                             `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败："900"  */
	Message    *string                             `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2QueryByMsgIdV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                             `json:"error"`      /*  错误码，描述错误信息  */
}

type Mq2QueryByMsgIdV3ReturnObjResponse struct {
	Data *Mq2QueryByMsgIdV3ReturnObjDataResponse `json:"data"` /*  消息详情数据对象  */
}

type Mq2QueryByMsgIdV3ReturnObjDataResponse struct {
	BrokerName                *string                                           `json:"brokerName"`                /*  Broker名称  */
	QueueId                   *int32                                            `json:"queueId"`                   /*  队列ID  */
	StoreSize                 *int32                                            `json:"storeSize"`                 /*  存储大小  */
	QueueOffset               *int32                                            `json:"queueOffset"`               /*  队列偏移量  */
	SysFlag                   *int32                                            `json:"sysFlag"`                   /*  系统标识  */
	BornTimestamp             *int64                                            `json:"bornTimestamp"`             /*  消息生成时间戳  */
	BornHost                  *string                                           `json:"bornHost"`                  /*  消息生成主机地址及端口  */
	StoreTimestamp            *int64                                            `json:"storeTimestamp"`            /*  消息存储时间戳  */
	StoreHost                 *string                                           `json:"storeHost"`                 /*  消息存储主机地址及端口  */
	MsgId                     *string                                           `json:"msgId"`                     /*  消息唯一标识  */
	CommitLogOffset           *int64                                            `json:"commitLogOffset"`           /*  提交日志偏移量  */
	BodyCRC                   *int32                                            `json:"bodyCRC"`                   /*  消息体CRC校验值  */
	ReconsumeTimes            *int32                                            `json:"reconsumeTimes"`            /*  重新消费次数  */
	PreparedTransactionOffset *int32                                            `json:"preparedTransactionOffset"` /*  预处理事务偏移量  */
	Topic                     *string                                           `json:"topic"`                     /*  消息主题  */
	Properties                *Mq2QueryByMsgIdV3ReturnObjDataPropertiesResponse `json:"properties"`                /*  消息属性集合  */
	MessageBody               *string                                           `json:"messageBody"`               /*  消息体内容  */
	Status                    *string                                           `json:"status"`                    /*  消息状态  */
	MessageBodyPath           *string                                           `json:"messageBodyPath"`           /*  消息体路径  */
	BodySize                  *int32                                            `json:"bodySize"`                  /*  消息体大小  */
}

type Mq2QueryByMsgIdV3ReturnObjDataPropertiesResponse struct {
	KEYS     *string `json:"KEYS"`     /*  keys  */
	UNIQ_KEY *string `json:"UNIQ_KEY"` /*  消息唯一键  */
	CLUSTER  *string `json:"CLUSTER"`  /*  集群ID  */
	TAGS     *string `json:"TAGS"`     /*  消息标签  */
	WAIT     *string `json:"WAIT"`     /*  是否等待标识  */
	DELAY    *string `json:"DELAY"`    /*  延迟时间  */
	BUYER_ID *string `json:"BUYER_ID"` /*  买家ID  */
}
