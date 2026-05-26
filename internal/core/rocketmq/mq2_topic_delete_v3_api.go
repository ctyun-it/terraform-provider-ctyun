package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2TopicDeleteV3Api
/* 删除主题 */
type Mq2TopicDeleteV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2TopicDeleteV3Api(client *core.CtyunClient) *Mq2TopicDeleteV3Api {
	return &Mq2TopicDeleteV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/topic/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2TopicDeleteV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2TopicDeleteV3Request) (*Mq2TopicDeleteV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2TopicDeleteV3Request
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2TopicDeleteV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2TopicDeleteV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例id  */
	TopicName  string `json:"topicName"`  /*  主题名称  */
}

type Mq2TopicDeleteV3Response struct {
	StatusCode int32                              `json:"statusCode"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    *string                            `json:"message"`    /*  描述状态。  */
	ReturnObj  *Mq2TopicDeleteV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                            `json:"error"`      /*  错误码，描述错误信息。  */
}

type Mq2TopicDeleteV3ReturnObjResponse struct{}
