package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbIplistenerDeleteApi
/* 删除ip_listener
 */type CtelbIplistenerDeleteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbIplistenerDeleteApi(client *core.CtyunClient) *CtelbIplistenerDeleteApi {
	return &CtelbIplistenerDeleteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/iplistener/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbIplistenerDeleteApi) Do(ctx context.Context, credential core.Credential, req *CtelbIplistenerDeleteRequest) (*CtelbIplistenerDeleteResponse, error) {
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
	var resp CtelbIplistenerDeleteResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbIplistenerDeleteRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池 ID  */
	IpListenerID string `json:"ipListenerID,omitempty"` /*  监听器 ID  */
}

type CtelbIplistenerDeleteResponse struct {
	StatusCode  int32                                   `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbIplistenerDeleteReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbIplistenerDeleteReturnObjResponse struct {
	IpListenerID string `json:"ipListenerID,omitempty"` /*  监听器 id  */
}
