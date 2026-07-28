package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2DeleteGroupApi
/* 删除订阅组 */
type Mq2DeleteGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2DeleteGroupApi(client *core.CtyunClient) *Mq2DeleteGroupApi {
	return &Mq2DeleteGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/group/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2DeleteGroupApi) Do(ctx context.Context, credential core.Credential, req *Mq2DeleteGroupRequest) (*Mq2DeleteGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*Mq2DeleteGroupRequest
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
	var resp Mq2DeleteGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2DeleteGroupRequest struct {
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	GroupName  string `json:"groupName"`  /*  订阅组名字  */
}

type Mq2DeleteGroupResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
