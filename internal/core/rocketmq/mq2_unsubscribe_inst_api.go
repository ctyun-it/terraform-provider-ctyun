package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2UnsubscribeInstApi
/* 功能介绍：退订实例 */
type Mq2UnsubscribeInstApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2UnsubscribeInstApi(client *core.CtyunClient) *Mq2UnsubscribeInstApi {
	return &Mq2UnsubscribeInstApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/spuInst/unsubscribeInst",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2UnsubscribeInstApi) Do(ctx context.Context, credential core.Credential, req *Mq2UnsubscribeInstRequest) (*Mq2UnsubscribeInstResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*Mq2UnsubscribeInstRequest
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
	var resp Mq2UnsubscribeInstResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2UnsubscribeInstRequest struct {
	ResourceDetailJson *Mq2UnsubscribeInstResourceDetailJsonRequest `json:"resourceDetailJson,omitempty"` /*  下退订订单详细信息  */
	RawType            string                                       `json:"type"`                         /*  固定值：2  */
}

type Mq2UnsubscribeInstResourceDetailJsonRequest struct {
	Resources    []*Mq2UnsubscribeInstResourceDetailJsonResourcesRequest `json:"resources,omitempty"` /*  下退订订单详细信息  */
	AutoApproval string                                                  `json:"autoApproval"`        /*  是否自动审批  */
}

type Mq2UnsubscribeInstResourceDetailJsonResourcesRequest struct {
	Id            *string  `json:"id,omitempty"`            /*  指定订单id 只支持单独退订，批量退订不支持指定订单id  */
	RefundingCash *string  `json:"refundingCash,omitempty"` /*  出账金额  */
	ResourceIds   []string `json:"resourceIds"`             /*  主虚拟资源Id  */
}

type Mq2UnsubscribeInstResponse struct {
	StatusCode *int32                               `json:"statusCode"` /*  响应状态码  */
	Message    *string                              `json:"message"`    /*  响应信息  */
	ReturnObj  *Mq2UnsubscribeInstReturnObjResponse `json:"returnObj"`  /*  响应对象  */
}

type Mq2UnsubscribeInstReturnObjResponse struct{}
