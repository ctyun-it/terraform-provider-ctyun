package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbDeleteDomainCertLinksApi
/* 删除多证书
 */type CtelbDeleteDomainCertLinksApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbDeleteDomainCertLinksApi(client *core.CtyunClient) *CtelbDeleteDomainCertLinksApi {
	return &CtelbDeleteDomainCertLinksApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/delete-domain-cert-links",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbDeleteDomainCertLinksApi) Do(ctx context.Context, credential core.Credential, req *CtelbDeleteDomainCertLinksRequest) (*CtelbDeleteDomainCertLinksResponse, error) {
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
	var resp CtelbDeleteDomainCertLinksResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbDeleteDomainCertLinksRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID  */
	DomainCertID string `json:"domainCertID,omitempty"` /*  多证书 ID  */
}

type CtelbDeleteDomainCertLinksResponse struct {
	StatusCode  int32                                        `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                       `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                       `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                       `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbDeleteDomainCertLinksReturnObjResponse `json:"returnObj"`             /*  检查结果  */
	Error       string                                       `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbDeleteDomainCertLinksReturnObjResponse struct{}
