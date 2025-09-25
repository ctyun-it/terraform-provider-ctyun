package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbCreateDomainCertLinkApi
/* 创建多证书
 */type CtelbCreateDomainCertLinkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbCreateDomainCertLinkApi(client *core.CtyunClient) *CtelbCreateDomainCertLinkApi {
	return &CtelbCreateDomainCertLinkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/create-domain-cert-links",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbCreateDomainCertLinkApi) Do(ctx context.Context, credential core.Credential, req *CtelbCreateDomainCertLinkRequest) (*CtelbCreateDomainCertLinkResponse, error) {
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
	var resp CtelbCreateDomainCertLinkResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbCreateDomainCertLinkRequest struct {
	RegionID      string `json:"regionID,omitempty"`      /*  资源池ID  */
	DomainName    string `json:"domainName,omitempty"`    /*  整个域名的总长度不能超过 255 个字符，每个子域名（包括顶级域名）的长度不能超过 63 个字符，域名中的字符集包括大写字母、小写字母、数字和连字符（减号），连字符不能位于域名的开头  */
	CertificateID string `json:"certificateID,omitempty"` /*  证书 ID  */
	ListenerID    string `json:"listenerID,omitempty"`    /*  监听器 ID  */
	Description   string `json:"description,omitempty"`   /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
}

type CtelbCreateDomainCertLinkResponse struct {
	StatusCode  int32                                       `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                      `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                      `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                      `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbCreateDomainCertLinkReturnObjResponse `json:"returnObj"`             /*  检查结果  */
	Error       string                                      `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbCreateDomainCertLinkReturnObjResponse struct {
	DomainCertID string `json:"domainCertID,omitempty"` /*  多证书 id  */
}
