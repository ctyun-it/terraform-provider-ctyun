package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbUpdateCertificateApi
/* 更新证书
 */type CtelbUpdateCertificateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbUpdateCertificateApi(client *core.CtyunClient) *CtelbUpdateCertificateApi {
	return &CtelbUpdateCertificateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/update-certificate",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbUpdateCertificateApi) Do(ctx context.Context, credential core.Credential, req *CtelbUpdateCertificateRequest) (*CtelbUpdateCertificateResponse, error) {
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
	var resp CtelbUpdateCertificateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbUpdateCertificateRequest struct {
	ClientToken   string `json:"clientToken,omitempty"`   /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID      string `json:"regionID,omitempty"`      /*  资源池ID  */
	ProjectID     string `json:"projectID,omitempty"`     /*  企业项目ID，默认为0  */
	ID            string `json:"ID,omitempty"`            /*  证书ID, 该字段后续废弃  */
	CertificateID string `json:"certificateID,omitempty"` /*  证书ID, 推荐使用该字段, 当同时使用 ID 和 certificateID 时，优先使用 certificateID（至少填一个）  */
	Name          string `json:"name,omitempty"`          /*  唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Description   string `json:"description,omitempty"`   /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:"{},./;'[\]·~！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
}

type CtelbUpdateCertificateResponse struct {
	StatusCode  int32                                      `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbUpdateCertificateReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                                     `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbUpdateCertificateReturnObjResponse struct {
	ID string `json:"ID,omitempty"` /*  证书ID  */
}
