package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbCreateCertificateApi
/* 创建证书
 */type CtelbCreateCertificateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbCreateCertificateApi(client *core.CtyunClient) *CtelbCreateCertificateApi {
	return &CtelbCreateCertificateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/create-certificate",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbCreateCertificateApi) Do(ctx context.Context, credential core.Credential, req *CtelbCreateCertificateRequest) (*CtelbCreateCertificateResponse, error) {
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
	var resp CtelbCreateCertificateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbCreateCertificateRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池ID  */
	Name        string `json:"name,omitempty"`        /*  唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Description string `json:"description,omitempty"` /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
	RawType     string `json:"type,omitempty"`        /*  证书类型。取值范围：Server（服务器证书）、Ca（Ca证书）  */
	PrivateKey  string `json:"privateKey,omitempty"`  /*  服务器证书私钥，服务器证书此字段必填  */
	Certificate string `json:"certificate,omitempty"` /*  type为Server该字段表示服务器证书公钥Pem内容;type为Ca该字段表示Ca证书Pem内容  */
}

type CtelbCreateCertificateResponse struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                   `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                   `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbCreateCertificateReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                                   `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbCreateCertificateReturnObjResponse struct {
	ID string `json:"ID,omitempty"` /*  证书ID  */
}
