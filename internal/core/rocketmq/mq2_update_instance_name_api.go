package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2UpdateInstanceNameApi
/* 更新实例名称 */
type Mq2UpdateInstanceNameApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2UpdateInstanceNameApi(client *core.CtyunClient) *Mq2UpdateInstanceNameApi {
	return &Mq2UpdateInstanceNameApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/instance/updateName",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2UpdateInstanceNameApi) Do(ctx context.Context, credential core.Credential, req *Mq2UpdateInstanceNameRequest) (*Mq2UpdateInstanceNameResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*Mq2UpdateInstanceNameRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2UpdateInstanceNameResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2UpdateInstanceNameRequest struct {
	ProdInstId   string `json:"prodInstId"`   /*  实例ID  */
	InstanceName string `json:"instanceName"` /*  新的实例名称  */
}

type Mq2UpdateInstanceNameResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
