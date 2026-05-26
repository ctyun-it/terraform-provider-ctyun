package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2DlqResendApi
/* 重发死信队列 */
type Mq2DlqResendApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2DlqResendApi(client *core.CtyunClient) *Mq2DlqResendApi {
	return &Mq2DlqResendApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/dlq/resend",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2DlqResendApi) Do(ctx context.Context, credential core.Credential, req *Mq2DlqResendRequest) (*Mq2DlqResendResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2DlqResendRequest
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
	var resp Mq2DlqResendResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2DlqResendRequest struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	GroupId    string `json:"groupId"`    /*  订阅组名称  */
	MsgId      string `json:"msgId"`      /*  消息ID  */
}

type Mq2DlqResendResponse struct {
	StatusCode *string                        `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                        `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2DlqResendReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例'里面的注释  */
	Error      *string                        `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2DlqResendReturnObjResponse struct{}
