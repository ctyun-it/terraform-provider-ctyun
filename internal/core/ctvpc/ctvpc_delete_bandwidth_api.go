package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDeleteBandwidthApi
/* 调用此接口可删除共享带宽。
 */type CtvpcDeleteBandwidthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDeleteBandwidthApi(client *core.CtyunClient) *CtvpcDeleteBandwidthApi {
	return &CtvpcDeleteBandwidthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/bandwidth/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDeleteBandwidthApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDeleteBandwidthRequest) (*CtvpcDeleteBandwidthResponse, error) {
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
	var resp CtvpcDeleteBandwidthResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDeleteBandwidthRequest struct {
	ClientToken string  `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string  `json:"regionID,omitempty"`    /*  共享带宽的区域id。  */
	ProjectID   *string `json:"projectID,omitempty"`   /*  企业项目 ID，默认为"0"  */
	BandwidthID string  `json:"bandwidthID,omitempty"` /*  共享带宽id。  */
}

type CtvpcDeleteBandwidthResponse struct {
	StatusCode  int32                                  `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcDeleteBandwidthReturnObjResponse `json:"returnObj"`             /*  返回删除的结果  */
}

type CtvpcDeleteBandwidthReturnObjResponse struct {
	MasterOrderID        *string `json:"masterOrderID,omitempty"`        /*  55d531d7bf2d47658897c42ffb918423  */
	MasterOrderNO        *string `json:"masterOrderNO,omitempty"`        /*  20221021191602644224  */
	RegionID             *string `json:"regionID,omitempty"`             /*  81f7728662dd11ec810800155d307d5b  */
	MasterResourceStatus *string `json:"masterResourceStatus,omitempty"` /*  refuned  */
}
