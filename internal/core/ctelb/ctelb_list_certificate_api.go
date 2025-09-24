package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbListCertificateApi
/* 获取证书列表
 */type CtelbListCertificateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbListCertificateApi(client *core.CtyunClient) *CtelbListCertificateApi {
	return &CtelbListCertificateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/list-certificate",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbListCertificateApi) Do(ctx context.Context, credential core.Credential, req *CtelbListCertificateRequest) (*CtelbListCertificateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.IDs != "" {
		ctReq.AddParam("IDs", req.IDs)
	}
	if req.Name != "" {
		ctReq.AddParam("name", req.Name)
	}
	if req.RawType != "" {
		ctReq.AddParam("type", req.RawType)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbListCertificateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbListCertificateRequest struct {
	RegionID string /*  资源池ID  */
	IDs      string /*  证书ID列表，以,分隔  */
	Name     string /*  证书名称，以,分隔，必须与ID顺序严格对应  */
	RawType  string /*  证书类型。Ca或Server，以,分隔，必须与ID和name的顺序严格对应  */
}

type CtelbListCertificateResponse struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                   `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                   `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbListCertificateReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       string                                   `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbListCertificateReturnObjResponse struct {
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
