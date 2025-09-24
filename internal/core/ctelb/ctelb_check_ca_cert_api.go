package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbCheckCaCertApi
/* 检查CA证书合法性
 */type CtelbCheckCaCertApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbCheckCaCertApi(client *core.CtyunClient) *CtelbCheckCaCertApi {
	return &CtelbCheckCaCertApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/check-ca-cert",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbCheckCaCertApi) Do(ctx context.Context, credential core.Credential, req *CtelbCheckCaCertRequest) (*CtelbCheckCaCertResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbCheckCaCertResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbCheckCaCertRequest struct {
	Certificate string `json:"certificate,omitempty"` /*  Ca证书Pem内容  */
}

type CtelbCheckCaCertResponse struct {
	StatusCode  int32                              `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                             `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                             `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                             `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       string                             `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbCheckCaCertReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbCheckCaCertReturnObjResponse struct {
	IsValid *bool `json:"isValid"` /*  是否合法  */
}
