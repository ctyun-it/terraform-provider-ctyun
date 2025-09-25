package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcRenewEipApi
/* 续订 EIP
 */type CtvpcRenewEipApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcRenewEipApi(client *core.CtyunClient) *CtvpcRenewEipApi {
	return &CtvpcRenewEipApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/eip/renew",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcRenewEipApi) Do(ctx context.Context, credential core.Credential, req *CtvpcRenewEipRequest) (*CtvpcRenewEipResponse, error) {
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
	var resp CtvpcRenewEipResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcRenewEipRequest struct {
	ClientToken     string  `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID        string  `json:"regionID,omitempty"`        /*  资源池 ID  */
	ProjectID       *string `json:"projectID,omitempty"`       /*  企业项目 ID，默认为'0'    */
	CycleType       string  `json:"cycleType,omitempty"`       /*  订购类型：month / year  */
	CycleCount      int32   `json:"cycleCount"`                /*  订购时长, 当 cycleType = month, 支持续订 1 - 11 个月; 当 cycleType = year, 支持续订 1 - 3 年  */
	EipID           string  `json:"eipID,omitempty"`           /*  eip id  */
	PayVoucherPrice *string `json:"payVoucherPrice,omitempty"` /*  代金券金额，支持到小数点后两位  */
}

type CtvpcRenewEipResponse struct {
	StatusCode  int32                             `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                           `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                           `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                           `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcRenewEipReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       *string                           `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcRenewEipReturnObjResponse struct {
	MasterOrderID *string `json:"masterOrderID,omitempty"` /*  订单 id  */
	MasterOrderNO *string `json:"masterOrderNO,omitempty"` /*  订单编号  */
	RegionID      *string `json:"regionID,omitempty"`      /*  可用区 id  */
}
