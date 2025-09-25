package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDisassociateIPv6FromBandwidthApi
/* 调用此接口可从共享带宽中移出IPv6s。
 */type CtvpcDisassociateIPv6FromBandwidthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDisassociateIPv6FromBandwidthApi(client *core.CtyunClient) *CtvpcDisassociateIPv6FromBandwidthApi {
	return &CtvpcDisassociateIPv6FromBandwidthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/bandwidth/disassociate-ipv6",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDisassociateIPv6FromBandwidthApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDisassociateIPv6FromBandwidthRequest) (*CtvpcDisassociateIPv6FromBandwidthResponse, error) {
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
	var resp CtvpcDisassociateIPv6FromBandwidthResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDisassociateIPv6FromBandwidthRequest struct {
	ClientToken string   `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string   `json:"regionID,omitempty"`    /*  共享带宽的区域id。  */
	BandwidthID string   `json:"bandwidthID,omitempty"` /*  共享带宽id。  */
	EipIDs      []string `json:"eipIDs"`                /*  portid数组。限额1-50  */
}

type CtvpcDisassociateIPv6FromBandwidthResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
