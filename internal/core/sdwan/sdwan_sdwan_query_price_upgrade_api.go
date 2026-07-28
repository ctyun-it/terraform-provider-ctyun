package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanQueryPriceUpgradeApi
/* 支持SDWAN智能网关变配询价，目前只支持升级带宽。 */
type SdwanSdwanQueryPriceUpgradeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanQueryPriceUpgradeApi(client *core.CtyunClient) *SdwanSdwanQueryPriceUpgradeApi {
	return &SdwanSdwanQueryPriceUpgradeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/query-price-upgrade",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanQueryPriceUpgradeApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanQueryPriceUpgradeRequest) (*SdwanSdwanQueryPriceUpgradeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanQueryPriceUpgradeRequest
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
	var resp SdwanSdwanQueryPriceUpgradeResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanQueryPriceUpgradeRequest struct {
	Bandwidth  int32  `json:"bandwidth"`  /*  带宽,范围为1-1000，单位mbps  */
	ResourceID string `json:"resourceID"` /*  SDWAN智能网关资源ID  */
}

type SdwanSdwanQueryPriceUpgradeResponse struct {
	StatusCode  int32                                         `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                       `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                       `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanQueryPriceUpgradeReturnObjResponse `json:"returnObj"`   /*  返回参数列表  */
	ErrorCode   *string                                       `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码.  */
	Error       *string                                       `json:"error"`       /*  业务细分码，为product.module.code三段式码.  */
}

type SdwanSdwanQueryPriceUpgradeReturnObjResponse struct {
	TotalPrice     int32                                                         `json:"totalPrice"`     /*  总价，单位人民币/元  */
	DiscountPrice  int32                                                         `json:"discountPrice"`  /*  折扣，单位人民币/元  */
	SubOrderPrices []*SdwanSdwanQueryPriceUpgradeReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  订单价格列表。  */
	FinalPrice     int32                                                         `json:"finalPrice"`     /*  最终价格，单位人民币/元  */
	IsSucceed      *bool                                                         `json:"isSucceed"`      /*  是否成功  */
}

type SdwanSdwanQueryPriceUpgradeReturnObjSubOrderPricesResponse struct {
	ServiceTag      *string                                                                      `json:"serviceTag"`      /*  服务标识  */
	TotalPrice      int32                                                                        `json:"totalPrice"`      /*  总价，单位人民币/元  */
	OrderItemPrices []*SdwanSdwanQueryPriceUpgradeReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"` /*  单项价格列表。  */
	FinalPrice      int32                                                                        `json:"finalPrice"`      /*  最终价格，单位人民币/元  */
}

type SdwanSdwanQueryPriceUpgradeReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ItemId       *string `json:"itemId"`       /*  单项ID  */
	TotalPrice   int32   `json:"totalPrice"`   /*  总价，单位人民币/元  */
	ResourceType *string `json:"resourceType"` /*  资源类型。  */
	FinalPrice   int32   `json:"finalPrice"`   /*  最终价格，单位人民币/元  */
}
