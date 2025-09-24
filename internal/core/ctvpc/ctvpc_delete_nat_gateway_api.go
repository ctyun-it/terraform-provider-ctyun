package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDeleteNatGatewayApi
/* 删除 NAT 网关
 */type CtvpcDeleteNatGatewayApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDeleteNatGatewayApi(client *core.CtyunClient) *CtvpcDeleteNatGatewayApi {
	return &CtvpcDeleteNatGatewayApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/delete-nat-gateway",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDeleteNatGatewayApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDeleteNatGatewayRequest) (*CtvpcDeleteNatGatewayResponse, error) {
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
	var resp CtvpcDeleteNatGatewayResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDeleteNatGatewayRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  区域id  */
	NatGatewayID string `json:"natGatewayID,omitempty"` /*  要查询的NAT网关的ID。  */
	ClientToken  string `json:"clientToken,omitempty"`  /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
}

type CtvpcDeleteNatGatewayResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcDeleteNatGatewayReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcDeleteNatGatewayReturnObjResponse struct {
	MasterOrderID        *string `json:"masterOrderID,omitempty"`        /*  订单id  */
	MasterOrderNO        *string `json:"masterOrderNO,omitempty"`        /*  订单编号  */
	RegionID             *string `json:"regionID,omitempty"`             /*  资源池ID  */
	MasterResourceStatus *string `json:"masterResourceStatus,omitempty"` /*  refuned  */
}
