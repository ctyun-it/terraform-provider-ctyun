package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsQueryPricesApi
/* 通过资源池ID和续订文件系统相关参数、查询文件系统续订订单价格
 */type SfsSfsQueryPricesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsQueryPricesApi(client *core.CtyunClient) *SfsSfsQueryPricesApi {
	return &SfsSfsQueryPricesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/renew-order/query-prices",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsQueryPricesApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsQueryPricesRequest) (*SfsSfsQueryPricesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsQueryPricesRequest
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
	var resp SfsSfsQueryPricesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsQueryPricesRequest struct {
	RegionID  string `json:"regionID,omitempty"`  /*  资源池ID  */
	SfsUID    string `json:"sfsUID,omitempty"`    /*  文件系统ID  */
	CycleType string `json:"cycleType,omitempty"` /*  续订时间单位(month/year)  */
	CycleCnt  int32  `json:"cycleCnt,omitempty"`  /*  续订周期数量  */
}

type SfsSfsQueryPricesResponse struct {
	StatusCode  int32                               `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                              `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                              `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsQueryPricesReturnObjResponse `json:"returnObj"`   /*  returnObj  */
	ErrorCode   string                              `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                              `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsQueryPricesReturnObjResponse struct {
	TotalPrice     float32                                             `json:"totalPrice"`     /*  查询到的订单总价  */
	FinalPrice     float32                                             `json:"finalPrice"`     /*  最终价格  */
	IsSucceed      *bool                                               `json:"isSucceed"`      /*  查询成功状态  */
	SubOrderPrices []*SfsSfsQueryPricesReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  查询到的订单价格详情  */
}

type SfsSfsQueryPricesReturnObjSubOrderPricesResponse struct {
	ServiceTag      string                                                             `json:"serviceTag"`      /*  服务类型  */
	TotalPrice      float32                                                            `json:"totalPrice"`      /*  查询到的订单总价  */
	OrderItemPrices []*SfsSfsQueryPricesReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"` /*  item价格信息  */
	FinalPrice      float32                                                            `json:"finalPrice"`      /*  最终价格  */
}

type SfsSfsQueryPricesReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ItemId       string  `json:"itemId"`       /*  itemId  */
	ResourceType string  `json:"resourceType"` /*  文件存储资源包类型(包周期为SFS_TURBOC)  */
	TotalPrice   float32 `json:"totalPrice"`   /*  总价格  */
	FinalPrice   float32 `json:"finalPrice"`   /*  最终价格  */
}
