package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryNewOrderPriceV41Api
/* 购买云产品时询价接口，支持云主机、云硬盘、弹性公网IP、NAT网关、共享带宽、物理机、性能保障型负载均衡、云主机备份存储库和云硬盘备份存储库产品的包年/包月或按量订单的询价功能
 */type CtecsQueryNewOrderPriceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryNewOrderPriceV41Api(client *core.CtyunClient) *CtecsQueryNewOrderPriceV41Api {
	return &CtecsQueryNewOrderPriceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/new-order/query-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryNewOrderPriceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryNewOrderPriceV41Request) (*CtecsQueryNewOrderPriceV41Response, error) {
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
	var resp CtecsQueryNewOrderPriceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryNewOrderPriceV41Request struct {
	RegionID        string                                         `json:"regionID,omitempty"`        /*  资源池ID  */
	ResourceType    string                                         `json:"resourceType,omitempty"`    /*  资源类型  */
	Count           int32                                          `json:"count,omitempty"`           /*  订购数量  */
	OnDemand        bool                                           `json:"onDemand"`                  /*  是否按需资源，true 按需 / false 包周期  */
	CycleType       string                                         `json:"cycleType,omitempty"`       /*  订购周期类型，当onDemand为false时为必填，可选值：MONTH 月/YEAR 年  */
	CycleCount      int32                                          `json:"cycleCount,omitempty"`      /*  订购周期大小，订购周期类型为MONTH时范围[1,60]，订购周期类型为YEAR时范围[1,5]，当onDemand为false时为必填  */
	FlavorName      string                                         `json:"flavorName,omitempty"`      /*  云主机规格，当resourceType为VM时必填  */
	ImageUUID       string                                         `json:"imageUUID,omitempty"`       /*  云主机镜像UUID，当resourceType为VM时必填  */
	SysDiskType     string                                         `json:"sysDiskType,omitempty"`     /*  云主机系统盘类型，当resourceType为VM时必填  */
	SysDiskSize     int32                                          `json:"sysDiskSize,omitempty"`     /*  云主机系统盘大小，范围[40,2048]，当resourceType为VM时必填  */
	Disks           []*CtecsQueryNewOrderPriceV41DisksRequest      `json:"disks"`                     /*  数据盘信息，当resourceType为VM选填，订购云主机时如果成套订购数据盘时需要该字段  */
	Bandwidth       int32                                          `json:"bandwidth,omitempty"`       /*  带宽大小，范围[1,2000]，当resourceType为IP时必填；当resourceType为VM时，如果成套订购弹性公网IP时需要该字段  */
	DiskType        string                                         `json:"diskType,omitempty"`        /*  磁盘类型，当resourceType为EBS时必填  */
	DiskSize        int32                                          `json:"diskSize,omitempty"`        /*  磁盘大小，范围[5,2000]，当resourceType为EBS时必填  */
	DiskMode        string                                         `json:"diskMode,omitempty"`        /*  磁盘模式(VBD/ISCSI/FCSAN)，当resourceType为EBS时必填  */
	NatType         string                                         `json:"natType,omitempty"`         /*  nat规格，当resourceType为NAT时必填  */
	IpPoolBandwidth int32                                          `json:"ipPoolBandwidth,omitempty"` /*  共享带宽大小，范围[5,2000]，当resourceType为IP_POOL时必填  */
	DeviceType      string                                         `json:"deviceType,omitempty"`      /*  物理机规格，当resourceType为BMS时必填  */
	AzName          string                                         `json:"azName,omitempty"`          /*  物理机规格可用区，当resourceType为BMS时必填  */
	OrderDisks      []*CtecsQueryNewOrderPriceV41OrderDisksRequest `json:"orderDisks"`                /*  物理机云硬盘信息，当resourceType为BMS选填  */
	ElbType         string                                         `json:"elbType,omitempty"`         /*  性能保障型负载均衡类型(支持standardI/standardII/enhancedI/enhancedII/higherI)，当resourceType为PGELB时必填  */
	CbrValue        int32                                          `json:"cbrValue,omitempty"`        /*  存储库大小，100-1024000GB，当resourceType为CBR_VM或CBR_VBS时必填  */
}

type CtecsQueryNewOrderPriceV41DisksRequest struct {
	DiskType string `json:"diskType,omitempty"` /*  磁盘类型  */
	DiskSize int32  `json:"diskSize,omitempty"` /*  磁盘大小  */
}

type CtecsQueryNewOrderPriceV41OrderDisksRequest struct {
	DiskType string `json:"diskType,omitempty"` /*  磁盘类型  */
	DiskSize int32  `json:"diskSize,omitempty"` /*  磁盘大小  */
}

type CtecsQueryNewOrderPriceV41Response struct {
	StatusCode  int32                                        `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)  */
	ErrorCode   string                                       `json:"errorCode,omitempty"`   /*  具体错误码标志  */
	Message     string                                       `json:"message,omitempty"`     /*  失败时的错误信息  */
	Description string                                       `json:"description,omitempty"` /*  失败时的错误描述  */
	ReturnObj   *CtecsQueryNewOrderPriceV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据，参见returnObj对象结构  */
	Error       string                                       `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码。请求成功时不返回该字段  */
}

type CtecsQueryNewOrderPriceV41ReturnObjResponse struct {
	TotalPrice     float32                                                      `json:"totalPrice"`     /*  总价格，单位CNY  */
	DiscountPrice  float32                                                      `json:"discountPrice"`  /*  折后价格，云主机相关产品有，单位CNY  */
	FinalPrice     float32                                                      `json:"finalPrice"`     /*  最终价格，单位CNY  */
	SubOrderPrices []*CtecsQueryNewOrderPriceV41ReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格信息  */
}

type CtecsQueryNewOrderPriceV41ReturnObjSubOrderPricesResponse struct {
	ServiceTag      string                                                                      `json:"serviceTag,omitempty"` /*  服务类型  */
	TotalPrice      float32                                                                     `json:"totalPrice"`           /*  子订单总价格，单位CNY  */
	OrderItemPrices []*CtecsQueryNewOrderPriceV41ReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"`      /*  item价格信息  */
}

type CtecsQueryNewOrderPriceV41ReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ResourceType string  `json:"resourceType,omitempty"` /*  资源类型  */
	TotalPrice   float32 `json:"totalPrice"`             /*  总价格，单位CNY  */
	FinalPrice   float32 `json:"finalPrice"`             /*  最终价格，单位CNY  */
}
