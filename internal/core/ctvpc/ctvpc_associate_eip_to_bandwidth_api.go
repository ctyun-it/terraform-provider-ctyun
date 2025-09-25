package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcAssociateEipToBandwidthApi
/* 调用此接口可添加EIPs至共享带宽。
 */type CtvpcAssociateEipToBandwidthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcAssociateEipToBandwidthApi(client *core.CtyunClient) *CtvpcAssociateEipToBandwidthApi {
	return &CtvpcAssociateEipToBandwidthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/bandwidth/associate-eip",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcAssociateEipToBandwidthApi) Do(ctx context.Context, credential core.Credential, req *CtvpcAssociateEipToBandwidthRequest) (*CtvpcAssociateEipToBandwidthResponse, error) {
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
	var resp CtvpcAssociateEipToBandwidthResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcAssociateEipToBandwidthRequest struct {
	ClientToken string   `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string   `json:"regionID,omitempty"`    /*  共享带宽的区域id。  */
	BandwidthID string   `json:"bandwidthID,omitempty"` /*  共享带宽id。  */
	EipIDs      []string `json:"eipIDs"`                /*  EIP数组，限额1-50  */
}

type CtvpcAssociateEipToBandwidthResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
