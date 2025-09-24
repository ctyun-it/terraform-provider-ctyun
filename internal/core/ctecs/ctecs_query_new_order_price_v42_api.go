package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryNewOrderPriceV42Api
/* 购买云产品时询价接口，支持云主机、云硬盘、弹性公网IP、NAT网关、共享带宽、物理机、性能保障型负载均衡、云主机备份存储库和云硬盘备份存储库产品的包年/包月或按量订单的询价功能
 */type CtecsQueryNewOrderPriceV42Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryNewOrderPriceV42Api(client *core.CtyunClient) *CtecsQueryNewOrderPriceV42Api {
	return &CtecsQueryNewOrderPriceV42Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/order/new-query-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryNewOrderPriceV42Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryNewOrderPriceV42Request) (*CtecsQueryNewOrderPriceV42Response, error) {
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
	var resp CtecsQueryNewOrderPriceV42Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryNewOrderPriceV42Request struct {
	RegionID        string                                         `json:"regionID,omitempty"`        /*  资源池ID  */
	ResourceType    string                                         `json:"resourceType,omitempty"`    /*  资源类型(VM:云主机,EBS:云硬盘,IP:弹性公网IP,IP_POOL:共享带宽,NAT:NAT网关,BMS:物理机,PGELB:性能保障型负载均衡,CBR_VM:云主机备份存储库,CBR_VBS:云硬盘备份存储库)  */
	Count           int32                                          `json:"count,omitempty"`           /*  订购数量  */
	OnDemand        bool                                           `json:"onDemand"`                  /*  是否按需资源，true 按需 / false 包周期  */
	CycleType       string                                         `json:"cycleType,omitempty"`       /*  订购周期类型，当onDemand为false时为必填，可选值：MONTH 月,YEAR 年  */
	CycleCount      int32                                          `json:"cycleCount,omitempty"`      /*  订购周期大小，订购周期类型为MONTH时范围[1,60]，订购周期类型为YEAR时范围[1,5]，当onDemand为false时为必填  */
	FlavorName      string                                         `json:"flavorName,omitempty"`      /*  云主机规格，当resourceType为VM时必填  */
	ImageUUID       string                                         `json:"imageUUID,omitempty"`       /*  云主机镜像UUID，当resourceType为VM时必填  */
	SysDiskType     string                                         `json:"sysDiskType,omitempty"`     /*  云主机系统盘类型(SAS:高IO,SATA:普通IO,SSD:超高IO,FAST-SSD:极速型SSD)，当resourceType为VM时必填  */
	SysDiskSize     int32                                          `json:"sysDiskSize,omitempty"`     /*  云主机系统盘大小，范围[40,2048]，当resourceType为VM时必填  */
	Disks           []*CtecsQueryNewOrderPriceV42DisksRequest      `json:"disks"`                     /*  数据盘信息，当resourceType为VM选填，订购云主机时如果成套订购数据盘时需要该字段  */
	Bandwidth       int32                                          `json:"bandwidth,omitempty"`       /*  带宽大小，范围[1,2000]，当resourceType为IP时必填；当resourceType为VM时，如果成套订购弹性公网IP时需要该字段  */
	DiskType        string                                         `json:"diskType,omitempty"`        /*  磁盘类型(SAS:高IO,SATA:普通IO,SSD:超高IO,FAST-SSD:极速型SSD)，当resourceType为EBS时必填  */
	DiskSize        int32                                          `json:"diskSize,omitempty"`        /*  磁盘大小，范围[5,2000]，当resourceType为EBS时必填  */
	DiskMode        string                                         `json:"diskMode,omitempty"`        /*  磁盘模式(VBD/ISCSI/FCSAN)，当resourceType为EBS时必填  */
	NatType         string                                         `json:"natType,omitempty"`         /*  nat规格(small:小型,medium:中型,large:大型,xlarge:超大型)，当resourceType为NAT时必填  */
	IpPoolBandwidth int32                                          `json:"ipPoolBandwidth,omitempty"` /*  共享带宽大小，范围[5,2000]，当resourceType为IP_POOL时必填  */
	DeviceType      string                                         `json:"deviceType,omitempty"`      /*  物理机规格，当resourceType为BMS时必填  */
	AzName          string                                         `json:"azName,omitempty"`          /*  物理机规格可用区，当resourceType为BMS时必填  */
	OrderDisks      []*CtecsQueryNewOrderPriceV42OrderDisksRequest `json:"orderDisks"`                /*  物理机云硬盘信息，当resourceType为BMS选填  */
	ElbType         string                                         `json:"elbType,omitempty"`         /*  性能保障型负载均衡类型(支持standardI/standardII/enhancedI/enhancedII/higherI)，当resourceType为PGELB时必填  */
	CbrValue        int32                                          `json:"cbrValue,omitempty"`        /*  存储库大小，100-1024000GB，当resourceType为CBR_VM或CBR_VBS时必填  */
}

type CtecsQueryNewOrderPriceV42DisksRequest struct {
	DiskType string `json:"diskType,omitempty"` /*  磁盘类型(SAS:高IO,SATA:普通IO,SSD:超高IO,FAST-SSD:极速型SSD)  */
	DiskSize int32  `json:"diskSize,omitempty"` /*  磁盘大小  */
}

type CtecsQueryNewOrderPriceV42OrderDisksRequest struct {
	DiskType string `json:"diskType,omitempty"` /*  磁盘类型(SAS:高IO,SATA:普通IO,SSD:超高IO,FAST-SSD:极速型SSD)  */
	DiskSize int32  `json:"diskSize,omitempty"` /*  磁盘大小  */
}

type CtecsQueryNewOrderPriceV42Response struct {
	StatusCode  int32                                        `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)  */
	ErrorCode   string                                       `json:"errorCode,omitempty"`   /*  具体错误码标志  */
	Message     string                                       `json:"message,omitempty"`     /*  失败时的错误信息  */
	Description string                                       `json:"description,omitempty"` /*  失败时的错误描述  */
	ReturnObj   *CtecsQueryNewOrderPriceV42ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据，参见returnObj对象结构  */
	Error       string                                       `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码。请求成功时不返回该字段  */
}

type CtecsQueryNewOrderPriceV42ReturnObjResponse struct {
	TotalPrice     float32                                                      `json:"totalPrice"`     /*  总价格，单位CNY  */
	DiscountPrice  float32                                                      `json:"discountPrice"`  /*  折后价格，云主机相关产品有，单位CNY  */
	FinalPrice     float32                                                      `json:"finalPrice"`     /*  最终价格，单位CNY  */
	SubOrderPrices []*CtecsQueryNewOrderPriceV42ReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格信息  */
}

type CtecsQueryNewOrderPriceV42ReturnObjSubOrderPricesResponse struct {
	ServiceTag      string                                                                      `json:"serviceTag,omitempty"` /*  服务类型  */
	TotalPrice      float32                                                                     `json:"totalPrice"`           /*  子订单总价格，单位CNY  */
	OrderItemPrices []*CtecsQueryNewOrderPriceV42ReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"`      /*  item价格信息  */
}

type CtecsQueryNewOrderPriceV42ReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ResourceType string  `json:"resourceType,omitempty"` /*  资源类型  */
	TotalPrice   float32 `json:"totalPrice"`             /*  总价格，单位CNY  */
	FinalPrice   float32 `json:"finalPrice"`             /*  最终价格，单位CNY  */
}
