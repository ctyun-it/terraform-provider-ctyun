package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2UpdateGroupPermApi
/* 配置订阅组读取权限 */
type Mq2UpdateGroupPermApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2UpdateGroupPermApi(client *core.CtyunClient) *Mq2UpdateGroupPermApi {
	return &Mq2UpdateGroupPermApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/group/updatePerm",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2UpdateGroupPermApi) Do(ctx context.Context, credential core.Credential, req *Mq2UpdateGroupPermRequest) (*Mq2UpdateGroupPermResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*Mq2UpdateGroupPermRequest
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
	var resp Mq2UpdateGroupPermResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2UpdateGroupPermRequest struct {
	ProdInstId     string   `json:"prodInstId"`     /*  实例ID  */
	GroupName      string   `json:"groupName"`      /*  订阅组名字  */
	BrokerNameList []string `json:"brokerNameList"` /*  Broker名字列表  */
	Enable         bool     `json:"enable"`         /*  是否开启  */
}

type Mq2UpdateGroupPermResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
