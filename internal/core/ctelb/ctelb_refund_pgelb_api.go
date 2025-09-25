package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbRefundPgelbApi
/* 保障型负载均衡退订
 */type CtelbRefundPgelbApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbRefundPgelbApi(client *core.CtyunClient) *CtelbRefundPgelbApi {
	return &CtelbRefundPgelbApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/refund-pgelb",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbRefundPgelbApi) Do(ctx context.Context, credential core.Credential, req *CtelbRefundPgelbRequest) (*CtelbRefundPgelbResponse, error) {
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
	var resp CtelbRefundPgelbResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbRefundPgelbRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string `json:"regionID,omitempty"`    /*  区域ID  */
	ElbID       string `json:"elbID,omitempty"`       /*  负载均衡 ID  */
}

type CtelbRefundPgelbResponse struct {
	StatusCode  int32                              `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                             `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                             `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                             `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbRefundPgelbReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       string                             `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbRefundPgelbReturnObjResponse struct {
	MasterOrderID        string `json:"masterOrderID,omitempty"`        /*  订单id。  */
	MasterOrderNO        string `json:"masterOrderNO,omitempty"`        /*  订单编号, 可以为 null。  */
	RegionID             string `json:"regionID,omitempty"`             /*  可用区id。  */
	MasterResourceStatus string `json:"masterResourceStatus,omitempty"` /*  资源状态: started（启用） / renewed（续订） / refunded（退订） / destroyed（销毁） / failed（失败） / starting（正在启用） / changed（变配）/ expired（过期）/ unknown（未知）  */
}
