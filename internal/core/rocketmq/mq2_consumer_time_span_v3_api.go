package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2ConsumerTimeSpanV3Api
/* 查询主题可重置的时间范围 */
type Mq2ConsumerTimeSpanV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2ConsumerTimeSpanV3Api(client *core.CtyunClient) *Mq2ConsumerTimeSpanV3Api {
	return &Mq2ConsumerTimeSpanV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/consumer/timeSpan",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2ConsumerTimeSpanV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2ConsumerTimeSpanV3Request) (*Mq2ConsumerTimeSpanV3Response, error) {
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
	var resp Mq2ConsumerTimeSpanV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2ConsumerTimeSpanV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	TopicName  string `json:"topicName"`  /*  Topic名字  */
}

type Mq2ConsumerTimeSpanV3Response struct {
	StatusCode *string                                 `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                                 `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2ConsumerTimeSpanV3ReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例'里面的注释  */
	Error      *string                                 `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2ConsumerTimeSpanV3ReturnObjResponse struct {
	Data *Mq2ConsumerTimeSpanV3ReturnObjDataResponse `json:"data"` /*  响应数据详情  */
}

type Mq2ConsumerTimeSpanV3ReturnObjDataResponse struct {
	TenantId     *string `json:"tenantId"`     /*  租户ID  */
	ClusterName  *string `json:"clusterName"`  /*  集群名称  */
	TopicName    *string `json:"topicName"`    /*  主题名称  */
	MinTimeStamp *int64  `json:"minTimeStamp"` /*  最小时间戳 单位:毫秒  */
	MaxTimeStamp *int64  `json:"maxTimeStamp"` /*  最大时间戳 单位:毫秒  */
}
