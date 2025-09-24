package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsCreateQueryPricesApi
/* 通过资源池ID和创建文件系统相关参数、查询文件系统创建订单价格
 */type SfsSfsCreateQueryPricesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsCreateQueryPricesApi(client *core.CtyunClient) *SfsSfsCreateQueryPricesApi {
	return &SfsSfsCreateQueryPricesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/new-order/query-prices",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsCreateQueryPricesApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsCreateQueryPricesRequest) (*SfsSfsCreateQueryPricesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsCreateQueryPricesRequest
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
	var resp SfsSfsCreateQueryPricesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsCreateQueryPricesRequest struct {
	RegionID   string `json:"regionID,omitempty"`   /*  资源池ID  */
	OrderNum   int32  `json:"orderNum,omitempty"`   /*  订购数量(最大订购数量:50)  */
	CycleType  string `json:"cycleType,omitempty"`  /*  订购时间单位(month/year)  */
	SfsSize    int32  `json:"sfsSize,omitempty"`    /*  订购规格(500-32768 单位GB)  */
	VolumeType string `json:"volumeType,omitempty"` /*  磁盘类型(标准型hdd, 性能型nvme, 标准专属型hdd_e)  */
	CycleCnt   int32  `json:"cycleCnt,omitempty"`   /*  订购时间长度(最大订购月数:36, 最大订购年数:3,选择包年可享受优惠，参考<a href="https://www.ctyun.cn/document/10027350/10237145" target="_blank">产品价格</a>)  */
}

type SfsSfsCreateQueryPricesResponse struct {
	StatusCode  int32                                     `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                    `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                    `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsCreateQueryPricesReturnObjResponse `json:"returnObj"`   /*  returnObj  */
	ErrorCode   string                                    `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                    `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsCreateQueryPricesReturnObjResponse struct {
	TotalPrice     float32                                                 `json:"totalPrice"`     /*  查询到的订单总价，单位：元  */
	FinalPrice     float32                                                 `json:"finalPrice"`     /*  最终价格，单位：元  */
	IsSucceed      *bool                                                   `json:"isSucceed"`      /*  查询成功状态  */
	SubOrderPrices *SfsSfsCreateQueryPricesReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  查询到的订单价格详情  */
}

type SfsSfsCreateQueryPricesReturnObjSubOrderPricesResponse struct {
	ServiceTag      string                                                                 `json:"serviceTag"`      /*  服务类型  */
	TotalPrice      float32                                                                `json:"totalPrice"`      /*  查询到的订单总价，单位：元  */
	OrderItemPrices *SfsSfsCreateQueryPricesReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"` /*  item价格信息  */
	FinalPrice      float32                                                                `json:"finalPrice"`      /*  最终价格，单位：元  */
}

type SfsSfsCreateQueryPricesReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ItemId       string  `json:"itemId"`       /*  itemId  */
	ResourceType string  `json:"resourceType"` /*  文件存储资源包类型(包周期为SFS_TURBOC)  */
	TotalPrice   float32 `json:"totalPrice"`   /*  总价格，单位：元  */
	FinalPrice   float32 `json:"finalPrice"`   /*  最终价格，单位：元  */
}
