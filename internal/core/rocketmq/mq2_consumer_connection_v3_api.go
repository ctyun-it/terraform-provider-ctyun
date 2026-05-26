package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2ConsumerConnectionV3Api
/* 查询订阅组当前客户端的连接情况 */
type Mq2ConsumerConnectionV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2ConsumerConnectionV3Api(client *core.CtyunClient) *Mq2ConsumerConnectionV3Api {
	return &Mq2ConsumerConnectionV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/consumer/connection",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2ConsumerConnectionV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2ConsumerConnectionV3Request) (*Mq2ConsumerConnectionV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("groupName", req.GroupName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2ConsumerConnectionV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2ConsumerConnectionV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	GroupName  string `json:"groupName"`  /*  订阅组名字  */
}

type Mq2ConsumerConnectionV3Response struct {
	StatusCode *string                                   `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                                   `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2ConsumerConnectionV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                                   `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2ConsumerConnectionV3ReturnObjResponse struct {
	Data *Mq2ConsumerConnectionV3ReturnObjDataResponse `json:"data"` /*  响应数据详情  */
}

type Mq2ConsumerConnectionV3ReturnObjDataResponse struct {
	Group *Mq2ConsumerConnectionV3ReturnObjDataGroupResponse `json:"group"` /*  消费组详情  */
}

type Mq2ConsumerConnectionV3ReturnObjDataGroupResponse struct {
	ConnectionSet     []*Mq2ConsumerConnectionV3ReturnObjDataGroupConnectionSetResponse   `json:"connectionSet"`     /*  客户端连接集合  */
	SubscriptionTable *Mq2ConsumerConnectionV3ReturnObjDataGroupSubscriptionTableResponse `json:"subscriptionTable"` /*  订阅关系表  */
	ConsumeType       *string                                                             `json:"consumeType"`       /*  消费类型  */
	MessageModel      *string                                                             `json:"messageModel"`      /*  消费方式 CLUSTERING-集群消费;BROADCASTING-广播消费  */
	ConsumeFromWhere  *string                                                             `json:"consumeFromWhere"`  /*  消费起始位置  */
}

type Mq2ConsumerConnectionV3ReturnObjDataGroupConnectionSetResponse struct {
	ClientId   *string `json:"clientId"`   /*  客户端ID  */
	ClientAddr *string `json:"clientAddr"` /*  客户端地址  */
	Language   *string `json:"language"`   /*  客户端开发语言  */
	Version    *int32  `json:"version"`    /*  客户端版本  */
}

type Mq2ConsumerConnectionV3ReturnObjDataGroupSubscriptionTableResponse struct {
	ClassFilterMode   *bool     `json:"classFilterMode"`   /*  类过滤模式开关  */
	Topic             *string   `json:"topic"`             /*  订阅主题  */
	SubString         *string   `json:"subString"`         /*  订阅表达式  */
	TagsSet           []*string `json:"tagsSet"`           /*  标签集合  */
	CodeSet           []*string `json:"codeSet"`           /*  代码集合  */
	SubVersion        *int64    `json:"subVersion"`        /*  订阅版本号  */
	ExpressionType    *string   `json:"expressionType"`    /*  表达式类型  */
	FilterClassSource *string   `json:"filterClassSource"` /*  过滤类源码  */
}
