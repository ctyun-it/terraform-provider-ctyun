package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbShowCertificateApi
/* 查看证书详情
 */type CtelbShowCertificateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbShowCertificateApi(client *core.CtyunClient) *CtelbShowCertificateApi {
	return &CtelbShowCertificateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/show-certificate",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbShowCertificateApi) Do(ctx context.Context, credential core.Credential, req *CtelbShowCertificateRequest) (*CtelbShowCertificateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ID != "" {
		ctReq.AddParam("ID", req.ID)
	}
	if req.CertificateID != "" {
		ctReq.AddParam("certificateID", req.CertificateID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbShowCertificateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbShowCertificateRequest struct {
	RegionID      string /*  资源池ID  */
	ID            string /*  证书ID, 该字段后续废弃  */
	CertificateID string /*  证书ID, 推荐使用该字段, 当同时使用 ID 和 certificateID 时，优先使用 certificateID  */
}

type CtelbShowCertificateResponse struct {
	StatusCode  int32                                  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbShowCertificateReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       string                                 `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbShowCertificateReturnObjResponse struct {
	RegionID    string `json:"regionID,omitempty"`    /*  资源池ID  */
	AzName      string `json:"azName,omitempty"`      /*  可用区名称  */
	ProjectID   string `json:"projectID,omitempty"`   /*  项目ID  */
	ID          string `json:"ID,omitempty"`          /*  证书ID  */
	Name        string `json:"name,omitempty"`        /*  名称  */
	Description string `json:"description,omitempty"` /*  描述  */
	RawType     string `json:"type,omitempty"`        /*  证书类型: certificate / ca  */
	PrivateKey  string `json:"privateKey,omitempty"`  /*  服务器证书私钥  */
	Certificate string `json:"certificate,omitempty"` /*  type为Server该字段表示服务器证书公钥Pem内容;type为Ca该字段表示Ca证书Pem内容  */
	Status      string `json:"status,omitempty"`      /*  状态: ACTIVE / INACTIVE  */
	CreatedTime string `json:"createdTime,omitempty"` /*  创建时间，为UTC格式  */
	UpdatedTime string `json:"updatedTime,omitempty"` /*  更新时间，为UTC格式  */
}
