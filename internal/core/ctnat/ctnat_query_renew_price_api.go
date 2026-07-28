package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtnatQueryRenewPriceApi
/* 续订私网NAT询价。
 */type CtnatQueryRenewPriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatQueryRenewPriceApi(client *core.CtyunClient) *CtnatQueryRenewPriceApi {
	return &CtnatQueryRenewPriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/privatenat/query-renew-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatQueryRenewPriceApi) Do(ctx context.Context, credential core.Credential, req *CtnatQueryRenewPriceRequest) (*CtnatQueryRenewPriceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtnatQueryRenewPriceRequest
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
	var resp CtnatQueryRenewPriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatQueryRenewPriceRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  私网NAT所在区域的区域id。  */
	NatGatewayID string `json:"natGatewayID,omitempty"` /*  私网NAT的ID  */
	ClientToken  string `json:"clientToken,omitempty"`  /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	CycleType    string `json:"cycleType,omitempty"`    /*  订购类型：month / year  */
}

type CtnatQueryRenewPriceResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                 `json:"message"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                 `json:"description"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                 `json:"errorCode"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtnatQueryRenewPriceReturnObjResponse `json:"returnObj"`   /*  业务数据  */
}

type CtnatQueryRenewPriceReturnObjResponse struct{}
