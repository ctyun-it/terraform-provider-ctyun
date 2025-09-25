package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDetachPortApi
/* 网卡解绑实例
 */type CtvpcDetachPortApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDetachPortApi(client *core.CtyunClient) *CtvpcDetachPortApi {
	return &CtvpcDetachPortApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ports/detach",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDetachPortApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDetachPortRequest) (*CtvpcDetachPortResponse, error) {
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
	var resp CtvpcDetachPortResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDetachPortRequest struct {
	ClientToken        string `json:"clientToken,omitempty"`        /*  客户端存根，用于保证订单幂等性。要求当个云平台账户内唯一  */
	RegionID           string `json:"regionID,omitempty"`           /*  资源池ID  */
	NetworkInterfaceID string `json:"networkInterfaceID,omitempty"` /*  网卡ID  */
}

type CtvpcDetachPortResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
