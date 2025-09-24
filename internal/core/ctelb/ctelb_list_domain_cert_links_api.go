package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbListDomainCertLinksApi
/* 获取多证书
 */type CtelbListDomainCertLinksApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbListDomainCertLinksApi(client *core.CtyunClient) *CtelbListDomainCertLinksApi {
	return &CtelbListDomainCertLinksApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/list-domain-cert-links",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbListDomainCertLinksApi) Do(ctx context.Context, credential core.Credential, req *CtelbListDomainCertLinksRequest) (*CtelbListDomainCertLinksResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ListenerID != "" {
		ctReq.AddParam("listenerID", req.ListenerID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbListDomainCertLinksResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbListDomainCertLinksRequest struct {
	RegionID   string /*  资源池ID  */
	ListenerID string /*  监听器 ID  */
}

type CtelbListDomainCertLinksResponse struct {
	StatusCode  int32                                        `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                       `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                       `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                       `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbListDomainCertLinksReturnObjResponse `json:"returnObj"`             /*  检查结果  */
}

type CtelbListDomainCertLinksReturnObjResponse struct {
	CertificateName string `json:"certificateName,omitempty"` /*  多证书 id  */
	CertificateType string `json:"certificateType,omitempty"` /*  类型类型: ca / certificate  */
	ExtDomainName   string `json:"extDomainName,omitempty"`   /*  扩展域名  */
	CreatedTime     string `json:"createdTime,omitempty"`     /*  创建时间  */
	DomainCertID    string `json:"domainCertID,omitempty"`    /*  多证书 id  */
}
