package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosQueryNewOrderPriceApi
/* 查询ZOS资源包的价格。
 */type ZosQueryNewOrderPriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosQueryNewOrderPriceApi(client *core.CtyunClient) *ZosQueryNewOrderPriceApi {
	return &ZosQueryNewOrderPriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/new-order/query-price",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosQueryNewOrderPriceApi) Do(ctx context.Context, credential core.Credential, req *ZosQueryNewOrderPriceRequest) (*ZosQueryNewOrderPriceResponse, error) {
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
	var resp ZosQueryNewOrderPriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosQueryNewOrderPriceRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  区域 ID  */
	PkgType      string `json:"pkgType,omitempty"`      /*  资源包类型，可选参数如下：①zosSize（ZOS存储空间包）、②zosMzSize（ZOS多AZ存储空间包）、③zosBytesSend（ZOS流出流量资源包)、④zosRequest（ZOS请求次数包，仅支持自定义规格：即pkgSpecType为defined）、⑤zosRetrievalFlow（ZOS数据取回流量包，仅支持自定义规格：即pkgSpecType为defined，仅支持storageClass为STANDARD_IA（低频存储）和 GLACIER（归档存储）、⑥zosRetrievalFrequency（ZOS数据取回次数包，仅支持自定义规格：即pkgSpecType为defined，仅支持storageClass为STANDARD_IA（低频存储）和 GLACIER（归档存储）  */
	PkgSpecType  string `json:"pkgSpecType,omitempty"`  /*  资源包规格类型，可选参数如下：①fixed（固定规格）②defined（自定义规格）  */
	PkgSpec      int64  `json:"pkgSpec,omitempty"`      /*  资源包规格大小，单位：GB。当pkgType选择为请求次数包zosRequest和数据取回次数包zosRetrievalFrequency时，单位为：万次。说明：①当资源包规格为固定包且资源包类型为：ZOS存储空间包 或 ZOS多AZ存储空间包时，可选参数列表如下：[40, 100, 500, 1024, 2048, 5120, 10240, 20480, 51200, 102400, 204800, 307200, 512000, 1048576, 2097152]; ②当资源包规格为固定包且资源包类型为ZOS流出流量资源包时，可选参数列表如下：[50, 100, 300, 500, 1024, 2048, 10240, 30720, 51200, 102400, 307200, 512000, 1048576, 2097152]  */
	CycleCnt     int64  `json:"cycleCnt,omitempty"`     /*  订购周期（最大订购月数：36，最大订购年数：3）  */
	CycleType    string `json:"cycleType,omitempty"`    /*  订购周期类型，可选参数如下：①month（按月订购）、②year（按年订购）  */
	OrderNum     int64  `json:"orderNum,omitempty"`     /*  订购数量（最大订购数量：50）  */
	StorageClass string `json:"storageClass,omitempty"` /*  存储类型，可选参数如下：①STANDARD（标准存储）、②STANDARD_IA（低频存储）、③GLACIER（归档存储）  */
}

type ZosQueryNewOrderPriceResponse struct {
	StatusCode  int64                                   `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                  `json:"message,omitempty"`     /*  状态描述  */
	ReturnObj   *ZosQueryNewOrderPriceReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	Description string                                  `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string                                  `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                  `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosQueryNewOrderPriceReturnObjResponse struct {
	DiscountPrice  float32                                                 `json:"discountPrice"`  /*  折后价格，单位：元   */
	TotalPrice     float32                                                 `json:"totalPrice"`     /*  总价格，单位：元  */
	SubOrderPrices []*ZosQueryNewOrderPriceReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格信息  */
	IsSucceed      *bool                                                   `json:"isSucceed"`      /*  是否成功  */
	FinalPrice     float32                                                 `json:"finalPrice"`     /*  最终价格，单位：元  */
}

type ZosQueryNewOrderPriceReturnObjSubOrderPricesResponse struct {
	ServiceTag      string                                                                 `json:"serviceTag,omitempty"` /*  服务类型  */
	TotalPrice      float32                                                                `json:"totalPrice"`           /*  总价格，单位：元  */
	FinalPrice      float32                                                                `json:"finalPrice"`           /*  最终价格，单位：元  */
	OrderItemPrices []*ZosQueryNewOrderPriceReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"`      /*  item价格信息  */
}

type ZosQueryNewOrderPriceReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ItemId       string  `json:"itemId,omitempty"`       /*  itemId  */
	ResourceType string  `json:"resourceType,omitempty"` /*  对象存储资源包类型，总共9种类型：①"ZOS_SIZE_DIY": "ZOS自定义存储空间包"②"ZOS_MZ_SIZE_DIY": "ZOS多AZ自定义存储空间包"③"ZOS_BYTES_SEND_DIY": "ZOS自定义流出流量包"④"ZOS_REQUEST_DIY": "ZOS自定义请求次数包"⑤"ZOS_RETRIEVAL_FLOW_DIY": "ZOS自定义数据取回流量包"⑥"ZOS_RETRIEVAL_FREQUENCY_DIY": "ZOS自定义数据取回次数包"⑦"ZOS_MZ_SIZE_P": "ZOS多AZ存储空间资源包"⑧"ZOS_SIZE_P": "ZOS存储空间资源包"⑨"ZOS_BYTES_SEND_P": "ZOS流出流量资源包"  */
	TotalPrice   float32 `json:"totalPrice"`             /*  总价格，单位：元  */
	FinalPrice   float32 `json:"finalPrice"`             /*  最终价格，单位：元  */
	CtyunName    string  `json:"ctyunName,omitempty"`    /*  天翼云服务名称  */
	InstanceCnt  string  `json:"instanceCnt,omitempty"`  /*  订购套数  */
}
