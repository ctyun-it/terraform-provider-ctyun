package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVniaDeleteNetworkPathApi
/* 删除网络路径
 */type CtvpcVniaDeleteNetworkPathApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVniaDeleteNetworkPathApi(client *core.CtyunClient) *CtvpcVniaDeleteNetworkPathApi {
	return &CtvpcVniaDeleteNetworkPathApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vnia/delete-network-path",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVniaDeleteNetworkPathApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVniaDeleteNetworkPathRequest) (*CtvpcVniaDeleteNetworkPathResponse, error) {
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
	var resp CtvpcVniaDeleteNetworkPathResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVniaDeleteNetworkPathRequest struct {
	RegionID      string `json:"regionID,omitempty"`      /*  资源池 ID  */
	NetworkPathID string `json:"networkPathID,omitempty"` /*  网络路径 ID  */
}

type CtvpcVniaDeleteNetworkPathResponse struct {
	StatusCode  int32                                        `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                      `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                      `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                      `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcVniaDeleteNetworkPathReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVniaDeleteNetworkPathReturnObjResponse struct {
	NetworkPathID *string `json:"networkPathID,omitempty"` /*  网络路径 ID  */
}
