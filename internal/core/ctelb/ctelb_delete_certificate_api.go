package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbDeleteCertificateApi
/* 删除证书
 */type CtelbDeleteCertificateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbDeleteCertificateApi(client *core.CtyunClient) *CtelbDeleteCertificateApi {
	return &CtelbDeleteCertificateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/delete-certificate",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbDeleteCertificateApi) Do(ctx context.Context, credential core.Credential, req *CtelbDeleteCertificateRequest) (*CtelbDeleteCertificateResponse, error) {
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
	var resp CtelbDeleteCertificateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbDeleteCertificateRequest struct {
	ClientToken   string `json:"clientToken,omitempty"`   /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID      string `json:"regionID,omitempty"`      /*  资源池ID  */
	ID            string `json:"ID,omitempty"`            /*  证书ID, 该字段后续废弃  */
	CertificateID string `json:"certificateID,omitempty"` /*  证书ID, 推荐使用该字段, 当同时使用 ID 和 certificateID 时，优先使用 certificateID  */
}

type CtelbDeleteCertificateResponse struct {
	StatusCode  int32                                      `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbDeleteCertificateReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                                     `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbDeleteCertificateReturnObjResponse struct {
	ID string `json:"ID,omitempty"` /*  证书ID  */
}
