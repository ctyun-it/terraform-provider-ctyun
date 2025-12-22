package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanQueryPriceRenewApi
/* 支持SDWAN智能网关包周期计费的续订询价。 */
type SdwanSdwanQueryPriceRenewApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanQueryPriceRenewApi(client *core.CtyunClient) *SdwanSdwanQueryPriceRenewApi {
	return &SdwanSdwanQueryPriceRenewApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/query-price-renew",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanQueryPriceRenewApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanQueryPriceRenewRequest) (*SdwanSdwanQueryPriceRenewResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanQueryPriceRenewRequest
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
	var resp SdwanSdwanQueryPriceRenewResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanQueryPriceRenewRequest struct {
	ResourceID string `json:"resourceID"` /*  SDWAN智能网关资源ID  */
	CycleType  string `json:"cycleType"`  /*  包周期类型，YEAR/MONTH。  */
	CycleCount int32  `json:"cycleCount"` /*  包周期数，周期最大长度不能超过36个月。  */
}

type SdwanSdwanQueryPriceRenewResponse struct {
	StatusCode  int32                                       `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                     `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                     `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanQueryPriceRenewReturnObjResponse `json:"returnObj"`   /*  返回参数列表  */
	ErrorCode   *string                                     `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码.  */
	Error       *string                                     `json:"error"`       /*  业务细分码，为product.module.code三段式码.  */
}

type SdwanSdwanQueryPriceRenewReturnObjResponse struct {
	TotalPrice     int32                                                       `json:"totalPrice"`     /*  总价，单位人民币/元  */
	IsSucceed      *bool                                                       `json:"isSucceed"`      /*  是否成功  */
	SubOrderPrices []*SdwanSdwanQueryPriceRenewReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  订单价格列表。  */
	FinalPrice     int32                                                       `json:"finalPrice"`     /*  最终价格，单位人民币/元  */
}

type SdwanSdwanQueryPriceRenewReturnObjSubOrderPricesResponse struct {
	ServiceTag      *string                                                                    `json:"serviceTag"`      /*  服务标识  */
	TotalPrice      int32                                                                      `json:"totalPrice"`      /*  总价，单位人民币/元  */
	OrderItemPrices []*SdwanSdwanQueryPriceRenewReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"` /*  单项价格列表。  */
	FinalPrice      int32                                                                      `json:"finalPrice"`      /*  最终价格，单位人民币/元  */
}

type SdwanSdwanQueryPriceRenewReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ItemId       *string `json:"itemId"`       /*  单项ID  */
	TotalPrice   int32   `json:"totalPrice"`   /*  总价，单位人民币/元  */
	ResourceType *string `json:"resourceType"` /*  资源类型。  */
	FinalPrice   int32   `json:"finalPrice"`   /*  最终价格，单位人民币/元  */
}
