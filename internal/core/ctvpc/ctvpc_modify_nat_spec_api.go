package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcModifyNatSpecApi
/* 变配 NAT 网关
 */type CtvpcModifyNatSpecApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcModifyNatSpecApi(client *core.CtyunClient) *CtvpcModifyNatSpecApi {
	return &CtvpcModifyNatSpecApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/modify-nat-gateway-spec",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcModifyNatSpecApi) Do(ctx context.Context, credential core.Credential, req *CtvpcModifyNatSpecRequest) (*CtvpcModifyNatSpecResponse, error) {
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
	var resp CtvpcModifyNatSpecResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcModifyNatSpecRequest struct {
	RegionID        string  `json:"regionID,omitempty"`        /*  NAT 网关所在的地域 ID。  */
	NatGatewayID    string  `json:"natGatewayID,omitempty"`    /*  NAT 网关 ID。  */
	Spec            int32   `json:"spec"`                      /*  规格(可传值：1-SMALL,2-MEDIUM,3-LARGE,4-XLARGE)  */
	ClientToken     string  `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	PayVoucherPrice *string `json:"payVoucherPrice,omitempty"` /*  代金券金额，支持到小数点后两位，仅包周期支持代金券  */
}

type CtvpcModifyNatSpecResponse struct {
	StatusCode  int32                                `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                              `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                              `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                              `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcModifyNatSpecReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       *string                              `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcModifyNatSpecReturnObjResponse struct {
	MasterOrderID *string `json:"masterOrderID,omitempty"` /*  订单 id  */
	MasterOrderNO *string `json:"masterOrderNO,omitempty"` /*  订单编号，可以为null  */
	RegionID      *string `json:"regionID,omitempty"`      /*  资源池ID  */
}
