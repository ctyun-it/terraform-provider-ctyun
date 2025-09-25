package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcModifyBandwidthSpecApi
/* 调用此接口可修改共享带宽的带宽峰值。
 */type CtvpcModifyBandwidthSpecApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcModifyBandwidthSpecApi(client *core.CtyunClient) *CtvpcModifyBandwidthSpecApi {
	return &CtvpcModifyBandwidthSpecApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/bandwidth/modify-spec",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcModifyBandwidthSpecApi) Do(ctx context.Context, credential core.Credential, req *CtvpcModifyBandwidthSpecRequest) (*CtvpcModifyBandwidthSpecResponse, error) {
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
	var resp CtvpcModifyBandwidthSpecResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcModifyBandwidthSpecRequest struct {
	ClientToken     string  `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID        string  `json:"regionID,omitempty"`        /*  共享带宽的区域 id。  */
	BandwidthID     string  `json:"bandwidthID,omitempty"`     /*  共享带宽 id。  */
	Bandwidth       int32   `json:"bandwidth"`                 /*  共享带宽的带宽峰值。5-1000  */
	PayVoucherPrice *string `json:"payVoucherPrice,omitempty"` /*  代金券金额，支持到小数点后两位  */
}

type CtvpcModifyBandwidthSpecResponse struct {
	StatusCode  int32                                        `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                      `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                      `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                      `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcModifyBandwidthSpecReturnObjResponse `json:"returnObj"`             /*  业务数据  */
	Error       *string                                      `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcModifyBandwidthSpecReturnObjResponse struct {
	MasterOrderID *string `json:"masterOrderID,omitempty"` /*  55d531d7bf2d47658897c42ffb918423  */
	MasterOrderNO *string `json:"masterOrderNO,omitempty"` /*  20221021191602644224  */
	RegionID      *string `json:"regionID,omitempty"`      /*  81f7728662dd11ec810800155d307d5b  */
}
