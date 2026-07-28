package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2TopicListV3Api
/* 获取Topic列表信息 */
type Mq2TopicListV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2TopicListV3Api(client *core.CtyunClient) *Mq2TopicListV3Api {
	return &Mq2TopicListV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/topic/list",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2TopicListV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2TopicListV3Request) (*Mq2TopicListV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	if req.TopicName != "" {
		ctReq.AddParam("topicName", req.TopicName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2TopicListV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2TopicListV3Request struct {
	RegionId   string `json:"regionId"`            /*  资源池编码  */
	ProdInstId string `json:"prodInstId"`          /*  实例ID  */
	TopicName  string `json:"topicName,omitempty"` /*  topic名字  */
}

type Mq2TopicListV3Response struct {
	StatusCode int32                            `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                          `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2TopicListV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                          `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2TopicListV3ReturnObjResponse struct {
	Rows  []*Mq2TopicListV3ReturnObjRowsResponse `json:"rows"`  /*  主题信息列表  */
	Total *int32                                 `json:"total"` /*  主题总数量  */
}

type Mq2TopicListV3ReturnObjRowsResponse struct {
	TopicName       *string `json:"topicName"`       /*  主题名称  */
	ReadQueueNums   *int32  `json:"readQueueNums"`   /*  读队列数量  */
	WriteQueueNums  *int32  `json:"writeQueueNums"`  /*  写队列数量  */
	Perm            *int32  `json:"perm"`            /*  读写权限   2   */
	TopicFilterType *string `json:"topicFilterType"` /*  主题过滤标志  */
	TopicSysFlag    *int32  `json:"topicSysFlag"`    /*  主题系统标识
	0 - 表示非系统主题
	1 - 系统内部主题  */
	Order       *bool   `json:"order"`       /*  是否为顺序消息  */
	Remark      string  `json:"remark"`      /*  备注  */
	ClusterName *string `json:"clusterName"` /*  集群名  */
	BrokerName  *string `json:"brokerName"`  /*  集群名称  */
	BrokerId    *int32  `json:"brokerId"`    /*  代理 id（数字类型）  */
	MessageType *string `json:"messageType"` /*  消息类型
	NORMAL-普通消息  */
}
