package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2QueryByKeyV3Api
/* 根据MessageKey查询消息 */
type Mq2QueryByKeyV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2QueryByKeyV3Api(client *core.CtyunClient) *Mq2QueryByKeyV3Api {
	return &Mq2QueryByKeyV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/message/queryByKey",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2QueryByKeyV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2QueryByKeyV3Request) (*Mq2QueryByKeyV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("topicName", req.TopicName)
	ctReq.AddParam("key", req.Key)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2QueryByKeyV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2QueryByKeyV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	TopicName  string `json:"topicName"`  /*  主题名字  */
	Key        string `json:"key"`        /*  用于查询的消息KEY值  */
}

type Mq2QueryByKeyV3Response struct {
	StatusCode *string                           `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                           `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2QueryByKeyV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                           `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2QueryByKeyV3ReturnObjResponse struct {
	Total *int32                                  `json:"total"` /*  总记录数。  */
	Rows  []*Mq2QueryByKeyV3ReturnObjRowsResponse `json:"rows"`  /*  消息详情列表。  */
}

type Mq2QueryByKeyV3ReturnObjRowsResponse struct {
	QueueId                   *int32      `json:"queueId"`                   /*  队列ID。  */
	StoreSize                 *int32      `json:"storeSize"`                 /*  存储大小。  */
	QueueOffset               *int32      `json:"queueOffset"`               /*  队列偏移量。  */
	SysFlag                   *int32      `json:"sysFlag"`                   /*  系统标记。  */
	BornTimestamp             *int64      `json:"bornTimestamp"`             /*  消息生成时间戳。  */
	BornHost                  *string     `json:"bornHost"`                  /*  消息生成主机地址。  */
	StoreTimestamp            *int64      `json:"storeTimestamp"`            /*  消息存储时间戳。  */
	StoreHost                 *string     `json:"storeHost"`                 /*  消息存储主机地址。  */
	MsgId                     *string     `json:"msgId"`                     /*  消息ID。  */
	CommitLogOffset           *int64      `json:"commitLogOffset"`           /*  提交日志偏移量。  */
	BodyCRC                   *int32      `json:"bodyCRC"`                   /*  消息体CRC校验值。  */
	ReconsumeTimes            *int32      `json:"reconsumeTimes"`            /*  重试次数。  */
	PreparedTransactionOffset *int64      `json:"preparedTransactionOffset"` /*  预处理事务偏移量。  */
	Topic                     *string     `json:"topic"`                     /*  主题名称。  */
	Properties                interface{} `json:"properties"`                /*  消息属性。  */
	MessageBody               *string     `json:"messageBody"`               /*  消息体内容。  */
	Status                    *string     `json:"status"`                    /*  消息状态。  */
}
