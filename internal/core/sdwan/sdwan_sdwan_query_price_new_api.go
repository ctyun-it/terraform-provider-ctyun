package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanQueryPriceNewApi
/* 包年包月订购SDWAN智能网关询价。 */
type SdwanSdwanQueryPriceNewApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanQueryPriceNewApi(client *core.CtyunClient) *SdwanSdwanQueryPriceNewApi {
	return &SdwanSdwanQueryPriceNewApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/query-price-new",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanQueryPriceNewApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanQueryPriceNewRequest) (*SdwanSdwanQueryPriceNewResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanQueryPriceNewRequest
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
	var resp SdwanSdwanQueryPriceNewResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanQueryPriceNewRequest struct {
	DeviceType   string                                      `json:"deviceType"`             /*  本参数表示网关形态<br/>取值范围：<br/>hardware:硬件<br/>software:软件  */
	BaseInfo     *SdwanSdwanQueryPriceNewBaseInfoRequest     `json:"baseInfo,omitempty"`     /*  网关开通基础信息  */
	GatewayInfo  *SdwanSdwanQueryPriceNewGatewayInfoRequest  `json:"gatewayInfo,omitempty"`  /*  智能网关资源，当上面的参数deviceType=hareware，此参数必填  */
	SoftwareInfo *SdwanSdwanQueryPriceNewSoftwareInfoRequest `json:"softwareInfo,omitempty"` /*  智能网关资源，当上面的参数deviceType=software，此参数必填  */
	AddrInfo     *SdwanSdwanQueryPriceNewAddrInfoRequest     `json:"addrInfo,omitempty"`     /*  联系人地址信息  */
	AddedService *SdwanSdwanQueryPriceNewAddedServiceRequest `json:"addedService,omitempty"` /*  增值业务  */
	CycleType    string                                      `json:"cycleType"`              /*  包周期类型，YEAR/MONTH。必须指定。  */
	CycleCount   int32                                       `json:"cycleCount"`             /*  包周期数。onDemand为false时必须指定。周期最大长度不能超过36个月  */
}

type SdwanSdwanQueryPriceNewBaseInfoRequest struct {
	EdgeName  string  `json:"edgeName"`       /*  智能网关实例名称，长度不超过64位  */
	Bandwidth int32   `json:"bandwidth"`      /*  网关带宽,单位Mbps，最大值1000  */
	Desc      *string `json:"desc,omitempty"` /*  描述，长度不超过128位  */
}

type SdwanSdwanQueryPriceNewGatewayInfoRequest struct {
	IsInstall int32  `json:"isInstall"` /*  智能网关是否需要装维1-是， 0-否  */
	EdgeType  string `json:"edgeType"`  /*  本参数表示设备类型<br/>取值范围：<br/>economic:经济<br/>standard:标准<br/>enterprise:企业<br/>enhance:增强  */
	UseType   string `json:"useType"`   /*  本参数表示使用方式<br/>取值范围：<br/>singleNode:单机<br/>activeStandby:旁挂  */
}

type SdwanSdwanQueryPriceNewSoftwareInfoRequest struct {
	EdgeType string `json:"edgeType"` /*  本参数表示设备类型<br/>取值范围：<br/>gateway:网关  */
	UseType  string `json:"useType"`  /*  本参数表示使用方式<br/>取值范围：<br/>singleNode:单机  */
}

type SdwanSdwanQueryPriceNewAddrInfoRequest struct {
	Name        string  `json:"name"`                  /*  客户联系人姓名,长度不超过64位  */
	Mobile      string  `json:"mobile"`                /*  客户联系人电话,长度不超过64位  */
	SpareMobile *string `json:"spareMobile,omitempty"` /*  客户联系人备用电话,长度不超过64位  */
	Email       string  `json:"email"`                 /*  客户联系人邮箱,长度不超过64位  */
	Province    string  `json:"province"`              /*  客户联系人省份,长度不超过64位  */
	City        string  `json:"city"`                  /*  客户联系人城市,长度不超过64位  */
	Address     string  `json:"address"`               /*  客户联系人详细地址,长度不超过64位  */
}

type SdwanSdwanQueryPriceNewAddedServiceRequest struct {
	Is5G   *bool `json:"is5G,omitempty"`   /*  是否支持5G,默认为否  */
	IsWIFI *bool `json:"isWIFI,omitempty"` /*  是否支持WIFI,默认为否  */
	IsLTE  *bool `json:"isLTE,omitempty"`  /*  是否支持LTE,默认为否  */
}

type SdwanSdwanQueryPriceNewResponse struct {
	StatusCode  int32                                     `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanQueryPriceNewReturnObjResponse `json:"returnObj"`   /*  返回参数列表  */
	ErrorCode   *string                                   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码.  */
	Error       *string                                   `json:"error"`       /*  业务细分码，为product.module.code三段式码.  */
}

type SdwanSdwanQueryPriceNewReturnObjResponse struct {
	TotalPrice     int32                                                     `json:"totalPrice"`     /*  总价，单位人民币/元  */
	IsSucceed      *bool                                                     `json:"isSucceed"`      /*  是否成功  */
	SubOrderPrices []*SdwanSdwanQueryPriceNewReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  订单价格列表。  */
	FinalPrice     int32                                                     `json:"finalPrice"`     /*  最终价格，单位人民币/元  */
}

type SdwanSdwanQueryPriceNewReturnObjSubOrderPricesResponse struct {
	ServiceTag      *string                                                                  `json:"serviceTag"`      /*  服务标识  */
	TotalPrice      int32                                                                    `json:"totalPrice"`      /*  总价，单位人民币/元  */
	OrderItemPrices []*SdwanSdwanQueryPriceNewReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"` /*  单项价格列表。  */
	FinalPrice      int32                                                                    `json:"finalPrice"`      /*  最终价格，单位人民币/元  */
}

type SdwanSdwanQueryPriceNewReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ItemId      *string `json:"itemId"`      /*  单项ID  */
	TotalPrice  int32   `json:"totalPrice"`  /*  总价，单位人民币/元  */
	ResorceType *string `json:"resorceType"` /*  资源类型。  */
	FinalPrice  int32   `json:"finalPrice"`  /*  最终价格，单位人民币/元  */
}
