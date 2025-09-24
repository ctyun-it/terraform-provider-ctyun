package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsQueryPricesApi
/* 通过资源池ID和扩容升级文件系统相关参数、查询文件系统扩容升级订单价格
 */type SfsQueryPricesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsQueryPricesApi(client *core.CtyunClient) *SfsQueryPricesApi {
	return &SfsQueryPricesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/upgrade-order/query-prices",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsQueryPricesApi) Do(ctx context.Context, credential core.Credential, req *SfsQueryPricesRequest) (*SfsQueryPricesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsQueryPricesRequest
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
	var resp SfsQueryPricesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsQueryPricesRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  文件系统ID  */
	SfsSize  int32  `json:"sfsSize,omitempty"`  /*  须填写扩容后的文件系统容量大小。单位 GB，500GB起步，默认容量上限为32TB，即32768GB。如无法满足需求，可提交工单申请扩大配额，详细参见<a href="https://www.ctyun.cn/document/10027350/10192640" target="_blank">服务配额</a>  */
}

type SfsQueryPricesResponse struct {
	StatusCode  int32                            `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                           `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                           `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsQueryPricesReturnObjResponse `json:"returnObj"`   /*  returnObj  */
	ErrorCode   string                           `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                           `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsQueryPricesReturnObjResponse struct {
	TotalPrice     float32                                          `json:"totalPrice"`     /*  查询到的订单总价  */
	FinalPrice     float32                                          `json:"finalPrice"`     /*  最终价格  */
	IsSucceed      *bool                                            `json:"isSucceed"`      /*  查询成功状态  */
	SubOrderPrices []*SfsQueryPricesReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  查询到的订单价格详情  */
	DiscountPrice  float32                                          `json:"discountPrice"`  /*  折后价格  */
}

type SfsQueryPricesReturnObjSubOrderPricesResponse struct {
	ServiceTag      string                                                          `json:"serviceTag"`      /*  服务类型  */
	TotalPrice      float32                                                         `json:"totalPrice"`      /*  查询到的订单总价  */
	OrderItemPrices []*SfsQueryPricesReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"` /*  item价格信息  */
	FinalPrice      float32                                                         `json:"finalPrice"`      /*  最终价格  */
	CycleType       int32                                                           `json:"cycleType"`       /*  订购周期(年/月)  */
}

type SfsQueryPricesReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ItemId       string  `json:"itemId"`       /*  itemId  */
	ResourceType string  `json:"resourceType"` /*  文件存储资源包类型(包周期为SFS_TURBOC)  */
	TotalPrice   float32 `json:"totalPrice"`   /*  总价格  */
	FinalPrice   float32 `json:"finalPrice"`   /*  最终价格  */
}
