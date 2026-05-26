package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// Mq2GetByTimeV3Api
/* 查询指定时间段内订阅组内存在的死信消息 */
type Mq2GetByTimeV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2GetByTimeV3Api(client *core.CtyunClient) *Mq2GetByTimeV3Api {
	return &Mq2GetByTimeV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/dlq/getByTime",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2GetByTimeV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2GetByTimeV3Request) (*Mq2GetByTimeV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("groupName", req.GroupName)
	ctReq.AddParam("beginTime", strconv.FormatInt(int64(req.BeginTime), 10))
	ctReq.AddParam("endTime", strconv.FormatInt(int64(req.EndTime), 10))
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2GetByTimeV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2GetByTimeV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	GroupName  string `json:"groupName"`  /*  订阅组名称  */
	BeginTime  int64  `json:"beginTime"`  /*  开始时间的毫秒时间戳  */
	EndTime    int64  `json:"endTime"`    /*  结束时间的毫秒时间戳  */
}

type Mq2GetByTimeV3Response struct {
	StatusCode *string                          `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                          `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2GetByTimeV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                          `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2GetByTimeV3ReturnObjResponse struct {
	Total *int32                                 `json:"total"` /*  总记录数。  */
	Rows  []*Mq2GetByTimeV3ReturnObjRowsResponse `json:"rows"`  /*  消息详情列表。  */
}

type Mq2GetByTimeV3ReturnObjRowsResponse struct {
	QueueId                   *int32                                         `json:"queueId"`                   /*  队列ID。  */
	StoreSize                 *int32                                         `json:"storeSize"`                 /*  存储大小，单位B  */
	QueueOffset               *int32                                         `json:"queueOffset"`               /*  队列偏移量。  */
	SysFlag                   *int32                                         `json:"sysFlag"`                   /*  系统标记。  */
	BornTimestamp             *int64                                         `json:"bornTimestamp"`             /*  消息生成时间戳。  */
	BornHost                  *string                                        `json:"bornHost"`                  /*  消息生成主机地址。  */
	StoreTimestamp            *int64                                         `json:"storeTimestamp"`            /*  消息存储时间戳。  */
	StoreHost                 *string                                        `json:"storeHost"`                 /*  消息存储主机地址。  */
	MsgId                     *string                                        `json:"msgId"`                     /*  消息ID。  */
	CommitLogOffset           *int64                                         `json:"commitLogOffset"`           /*  提交日志偏移量。  */
	BodyCRC                   *int32                                         `json:"bodyCRC"`                   /*  消息体CRC校验值。  */
	ReconsumeTimes            *int32                                         `json:"reconsumeTimes"`            /*  重试次数。  */
	PreparedTransactionOffset *int64                                         `json:"preparedTransactionOffset"` /*  预处理事务偏移量。  */
	Topic                     *string                                        `json:"topic"`                     /*  主题名称。  */
	Properties                *Mq2GetByTimeV3ReturnObjRowsPropertiesResponse `json:"properties"`                /*  消息属性。  */
	MessageBody               *string                                        `json:"messageBody"`               /*  消息体内容。  */
}

type Mq2GetByTimeV3ReturnObjRowsPropertiesResponse struct {
	MIN_OFFSET        *string `json:"MIN_OFFSET"`        /*  最小偏移量。  */
	REAL_TOPIC        *string `json:"REAL_TOPIC"`        /*  真实主题名称。  */
	ORIGIN_MESSAGE_ID *string `json:"ORIGIN_MESSAGE_ID"` /*  原始消息ID。  */
	RETRY_TOPIC       *string `json:"RETRY_TOPIC"`       /*  重试主题名称。  */
	MAX_OFFSET        *string `json:"MAX_OFFSET"`        /*  最大偏移量。  */
	UNIQ_KEY          *string `json:"UNIQ_KEY"`          /*  唯一标识键。  */
	WAIT              *string `json:"WAIT"`              /*  等待标记。  */
	DELAY             *string `json:"DELAY"`             /*  延迟级别。  */
	REAL_QID          *string `json:"REAL_QID"`          /*  真实队列ID。  */
}
