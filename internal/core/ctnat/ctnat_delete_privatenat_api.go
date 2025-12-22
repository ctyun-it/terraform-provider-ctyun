package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtnatDeletePrivatenatApi
/* 删除私网NAT
 */type CtnatDeletePrivatenatApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatDeletePrivatenatApi(client *core.CtyunClient) *CtnatDeletePrivatenatApi {
	return &CtnatDeletePrivatenatApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/privatenat/delete-privatenat",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatDeletePrivatenatApi) Do(ctx context.Context, credential core.Credential, req *CtnatDeletePrivatenatRequest) (*CtnatDeletePrivatenatResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtnatDeletePrivatenatRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtnatDeletePrivatenatResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatDeletePrivatenatRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  区域id  */
	NatGatewayID string `json:"natGatewayID,omitempty"` /*  要删除的私网NAT的ID。  */
	ClientToken  string `json:"clientToken,omitempty"`  /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
}

type CtnatDeletePrivatenatResponse struct {
	StatusCode  int32                                   `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                  `json:"message"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                  `json:"description"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                  `json:"errorCode"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtnatDeletePrivatenatReturnObjResponse `json:"returnObj"`   /*  接口业务数据  */
}

type CtnatDeletePrivatenatReturnObjResponse struct {
	MasterOrderID        *string `json:"masterOrderID,omitempty"`        /*  订单id  */
	MasterOrderNO        *string `json:"masterOrderNO,omitempty"`        /*  订单编号  */
	RegionID             *string `json:"regionID,omitempty"`             /*  资源池ID  */
	MasterResourceStatus *string `json:"masterResourceStatus,omitempty"` /*  refuned  */
}
