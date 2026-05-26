package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2GroupDeleteV3Api
/* 删除消费组 */
type Mq2GroupDeleteV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2GroupDeleteV3Api(client *core.CtyunClient) *Mq2GroupDeleteV3Api {
	return &Mq2GroupDeleteV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/group/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2GroupDeleteV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2GroupDeleteV3Request) (*Mq2GroupDeleteV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2GroupDeleteV3Request
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
	var resp Mq2GroupDeleteV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2GroupDeleteV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例id  */
	GroupName  string `json:"groupName"`  /*  消费组名称  */
}

type Mq2GroupDeleteV3Response struct {
	StatusCode int32                              `json:"statusCode"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    *string                            `json:"message"`    /*  描述状态。  */
	ReturnObj  *Mq2GroupDeleteV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                            `json:"error"`      /*  错误码，描述错误信息。  */
}

type Mq2GroupDeleteV3ReturnObjResponse struct{}
