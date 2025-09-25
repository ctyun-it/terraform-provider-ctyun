package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcL2gwRenewApi
/* 续订。
 */type CtvpcL2gwRenewApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcL2gwRenewApi(client *core.CtyunClient) *CtvpcL2gwRenewApi {
	return &CtvpcL2gwRenewApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/l2gw/renew",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcL2gwRenewApi) Do(ctx context.Context, credential core.Credential, req *CtvpcL2gwRenewRequest) (*CtvpcL2gwRenewResponse, error) {
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
	var resp CtvpcL2gwRenewResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcL2gwRenewRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID  */
	CycleType   string `json:"cycleType,omitempty"`   /*  订购类型：month（包月） / year（包年）  */
	CycleCount  int32  `json:"cycleCount"`            /*  订购时长, ，包月1~11，包年1~3  */
	L2gwID      string `json:"l2gwID,omitempty"`      /*  l2gwID  */
}

type CtvpcL2gwRenewResponse struct {
	StatusCode  int32                              `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                            `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                            `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                            `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcL2gwRenewReturnObjResponse `json:"returnObj"`             /*  业务数据  */
}

type CtvpcL2gwRenewReturnObjResponse struct {
	MasterOrderID *string `json:"masterOrderID,omitempty"` /*  订单id。  */
	MasterOrderNO *string `json:"masterOrderNO,omitempty"` /*  订单编号, 可以为 null。  */
	RegionID      *string `json:"regionID,omitempty"`      /*  可用区id。  */
}
