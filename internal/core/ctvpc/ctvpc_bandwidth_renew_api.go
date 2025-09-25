package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcBandwidthRenewApi
/* 续订共享带宽
 */type CtvpcBandwidthRenewApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcBandwidthRenewApi(client *core.CtyunClient) *CtvpcBandwidthRenewApi {
	return &CtvpcBandwidthRenewApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/bandwidth/renew",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcBandwidthRenewApi) Do(ctx context.Context, credential core.Credential, req *CtvpcBandwidthRenewRequest) (*CtvpcBandwidthRenewResponse, error) {
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
	var resp CtvpcBandwidthRenewResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcBandwidthRenewRequest struct {
	RegionID        string  `json:"regionID,omitempty"`        /*  共享带宽的区域id。  */
	ProjectID       *string `json:"projectID,omitempty"`       /*  企业项目 ID，默认为'0'  */
	BandwidthID     string  `json:"bandwidthID,omitempty"`     /*  共享带宽 ID  */
	ClientToken     string  `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	CycleType       string  `json:"cycleType,omitempty"`       /*  订购类型：month / year  */
	CycleCount      int32   `json:"cycleCount"`                /*  订购时长, 当 cycleType = month, 支持续订 1 - 11 个月; 当 cycleType = year, 支持续订 1 - 3 年  */
	PayVoucherPrice *string `json:"payVoucherPrice,omitempty"` /*  代金券金额，支持到小数点后两位  */
}

type CtvpcBandwidthRenewResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcBandwidthRenewReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                 `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcBandwidthRenewReturnObjResponse struct {
	MasterOrderID        *string `json:"masterOrderID,omitempty"`        /*  订单id。  */
	MasterOrderNO        *string `json:"masterOrderNO,omitempty"`        /*  订单编号, 可以为 null。  */
	MasterResourceStatus *string `json:"masterResourceStatus,omitempty"` /*  资源状态: started（启用） / renewed（续订） / refunded（退订） / destroyed（销毁） / failed（失败） / starting（正在启用） / changed（变配）/ expired（过期）/ unknown（未知）  */
	MasterResourceID     *string `json:"masterResourceID,omitempty"`     /*  可以为 null。  */
	RegionID             *string `json:"regionID,omitempty"`             /*  可用区id。  */
}
