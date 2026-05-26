package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2TopicStatusV3Api
/* 查询Topic状态 */
type Mq2TopicStatusV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2TopicStatusV3Api(client *core.CtyunClient) *Mq2TopicStatusV3Api {
	return &Mq2TopicStatusV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/topic/status",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2TopicStatusV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2TopicStatusV3Request) (*Mq2TopicStatusV3Response, error) {
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
	var resp Mq2TopicStatusV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2TopicStatusV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	TopicName  string `json:"topicName"`  /*  topic名字  */
}

type Mq2TopicStatusV3Response struct {
	StatusCode *string                            `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败："900"  */
	Message    *string                            `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2TopicStatusV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                            `json:"error"`      /*  错误码，描述错误信息  */
}

type Mq2TopicStatusV3ReturnObjResponse struct {
	Data *Mq2TopicStatusV3ReturnObjDataResponse `json:"data"` /*  响应数据对象  */
}

type Mq2TopicStatusV3ReturnObjDataResponse struct {
	Perm          *int32 `json:"perm"`          /*  权限标识  */
	TotalCount    *int32 `json:"totalCount"`    /*  总记录数  */
	LastTimeStamp *int64 `json:"lastTimeStamp"` /*  最新时间戳  */
}
