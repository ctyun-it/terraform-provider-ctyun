package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcFlowsessionRefundApi
/* 退订付费购买的流量会话
 */type CtvpcFlowsessionRefundApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcFlowsessionRefundApi(client *core.CtyunClient) *CtvpcFlowsessionRefundApi {
	return &CtvpcFlowsessionRefundApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/flowsession/refund",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcFlowsessionRefundApi) Do(ctx context.Context, credential core.Credential, req *CtvpcFlowsessionRefundRequest) (*CtvpcFlowsessionRefundResponse, error) {
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
	var resp CtvpcFlowsessionRefundResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcFlowsessionRefundRequest struct {
	RegionID      string `json:"regionID,omitempty"`      /*  区域ID  */
	FlowSessionID string `json:"flowSessionID,omitempty"` /*  流量镜像会话 ID  */
}

type CtvpcFlowsessionRefundResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcFlowsessionRefundReturnObjResponse `json:"returnObj"`             /*  object  */
}

type CtvpcFlowsessionRefundReturnObjResponse struct {
	MasterOrderID        *string `json:"masterOrderID,omitempty"`        /*  订单id。  */
	MasterOrderNO        *string `json:"masterOrderNO,omitempty"`        /*  订单编号, 可以为 null。  */
	MasterResourceStatus *string `json:"masterResourceStatus,omitempty"` /*  资源状态: started（启用） / renewed（续订） / refunded（退订） / destroyed（销毁） / failed（失败） / starting（正在启用） / changed（变配）/ expired（过期）/ unknown（未知）  */
	MasterResourceID     *string `json:"masterResourceID,omitempty"`     /*  可以为 null。  */
}
