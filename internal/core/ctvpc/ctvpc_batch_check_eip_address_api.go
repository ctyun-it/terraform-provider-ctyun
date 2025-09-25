package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcBatchCheckEipAddressApi
/* 调用此接口可批量检查弹性公网IP地址是否已经被使用。
 */type CtvpcBatchCheckEipAddressApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcBatchCheckEipAddressApi(client *core.CtyunClient) *CtvpcBatchCheckEipAddressApi {
	return &CtvpcBatchCheckEipAddressApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/eip/check-addresses",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcBatchCheckEipAddressApi) Do(ctx context.Context, credential core.Credential, req *CtvpcBatchCheckEipAddressRequest) (*CtvpcBatchCheckEipAddressResponse, error) {
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
	var resp CtvpcBatchCheckEipAddressResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcBatchCheckEipAddressRequest struct {
	ClientToken  *string  `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID     string   `json:"regionID,omitempty"`    /*  资源池 ID  */
	ProjectID    *string  `json:"projectID,omitempty"`   /*  企业项目 ID，默认为'0'  */
	EipAddresses []string `json:"eipAddresses"`          /*  弹性公网IP地址列表  */
}

type CtvpcBatchCheckEipAddressResponse struct {
	StatusCode  int32                                         `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                       `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                       `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                       `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcBatchCheckEipAddressReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                       `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcBatchCheckEipAddressReturnObjResponse struct {
	EipAddress *string `json:"eipAddress,omitempty"` /*  弹性公网IP地址  */
	Used       *bool   `json:"used"`                 /*  是否被使用  */
}
