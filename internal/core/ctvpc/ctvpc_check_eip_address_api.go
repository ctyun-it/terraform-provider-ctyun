package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCheckEipAddressApi
/* 调用此接口可检查弹性公网IP地址是否已经被使用。
 */type CtvpcCheckEipAddressApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCheckEipAddressApi(client *core.CtyunClient) *CtvpcCheckEipAddressApi {
	return &CtvpcCheckEipAddressApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/eip/check-address",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCheckEipAddressApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCheckEipAddressRequest) (*CtvpcCheckEipAddressResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.ClientToken != nil {
		ctReq.AddParam("clientToken", *req.ClientToken)
	}
	ctReq.AddParam("regionID", req.RegionID)
	if req.ProjectID != nil {
		ctReq.AddParam("projectID", *req.ProjectID)
	}
	ctReq.AddParam("eipAddress", req.EipAddress)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcCheckEipAddressResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCheckEipAddressRequest struct {
	ClientToken *string /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string  /*  资源池 ID  */
	ProjectID   *string /*  企业项目 ID，默认为'0'  */
	EipAddress  string  /*  弹性公网IP地址  */
}

type CtvpcCheckEipAddressResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	EipAddress  *string `json:"eipAddress,omitempty"`  /*  弹性公网IP地址  */
	Used        *bool   `json:"used"`                  /*  是否被使用  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
