package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbCheckServerCertApi
/* 检查Server证书合法性
 */type CtelbCheckServerCertApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbCheckServerCertApi(client *core.CtyunClient) *CtelbCheckServerCertApi {
	return &CtelbCheckServerCertApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/check-certificate",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbCheckServerCertApi) Do(ctx context.Context, credential core.Credential, req *CtelbCheckServerCertRequest) (*CtelbCheckServerCertResponse, error) {
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
	var resp CtelbCheckServerCertResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbCheckServerCertRequest struct {
	Certificate string `json:"certificate,omitempty"` /*  证书内容  */
	PrivateKey  string `json:"privateKey,omitempty"`  /*  服务器证书私钥  */
}

type CtelbCheckServerCertResponse struct {
	StatusCode  int32                                  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       string                                 `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbCheckServerCertReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbCheckServerCertReturnObjResponse struct {
	IsValid *bool `json:"isValid"` /*  是否合法  */
}
