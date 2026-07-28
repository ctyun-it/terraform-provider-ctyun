package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanOrderEdgeNewApi
/* 支持按需/包年包月订购SDWAN智能网关。 */
type SdwanSdwanOrderEdgeNewApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanOrderEdgeNewApi(client *core.CtyunClient) *SdwanSdwanOrderEdgeNewApi {
	return &SdwanSdwanOrderEdgeNewApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/new",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanOrderEdgeNewApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanOrderEdgeNewRequest) (*SdwanSdwanOrderEdgeNewResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanOrderEdgeNewRequest
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
	var resp SdwanSdwanOrderEdgeNewResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanOrderEdgeNewRequest struct {
	ClientToken  *string                                    `json:"clientToken,omitempty"`  /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一。  */
	DeviceType   string                                     `json:"deviceType"`             /*  本参数表示网关形态<br/>取值范围：<br/>hardware:硬件<br/>software:软件  */
	BaseInfo     *SdwanSdwanOrderEdgeNewBaseInfoRequest     `json:"baseInfo,omitempty"`     /*  网关开通基础信息  */
	GatewayInfo  *SdwanSdwanOrderEdgeNewGatewayInfoRequest  `json:"gatewayInfo,omitempty"`  /*  智能网关资源，deviceType=hareware必填  */
	SoftwareInfo *SdwanSdwanOrderEdgeNewSoftwareInfoRequest `json:"softwareInfo,omitempty"` /*  智能网关资源，deviceType=software必填  */
	AddrInfo     *SdwanSdwanOrderEdgeNewAddrInfoRequest     `json:"addrInfo,omitempty"`     /*  联系人地址信息  */
	AddedService *SdwanSdwanOrderEdgeNewAddedServiceRequest `json:"addedService,omitempty"` /*  增值业务  */
	OnDemand     *bool                                      `json:"onDemand,omitempty"`     /*  是否按需下单。默认为False  */
	CycleType    *string                                    `json:"cycleType,omitempty"`    /*  包周期类型，YEAR/MONTH。onDemand为false时，必须指定。  */
	CycleCount   int32                                      `json:"cycleCount"`             /*  包周期数。onDemand为false时必须指定。周期最大长度不能超过36个月  */
}

type SdwanSdwanOrderEdgeNewBaseInfoRequest struct {
	EdgeName  string  `json:"edgeName"`       /*  智能网关实例名称  */
	Bandwidth int32   `json:"bandwidth"`      /*  网关带宽  */
	Desc      *string `json:"desc,omitempty"` /*  描述  */
}

type SdwanSdwanOrderEdgeNewGatewayInfoRequest struct {
	IsInstall int32  `json:"isInstall"` /*  智能网关是否需要装维1-是， 0-否  */
	EdgeType  string `json:"edgeType"`  /*  本参数表示设备类型<br/>取值范围：<br/>economic:经济<br/>standard:标准<br/>enterprise:企业<br/>enhance:增强  */
	UseType   string `json:"useType"`   /*  本参数表示使用方式<br/>取值范围：<br/>singleNode:单机<br/>activeStandby:旁挂  */
}

type SdwanSdwanOrderEdgeNewSoftwareInfoRequest struct {
	EdgeType string `json:"edgeType"` /*  本参数表示设备类型<br/>取值范围：<br/>gateway:网关  */
	UseType  string `json:"useType"`  /*  本参数表示使用方式<br/>取值范围：<br/>singleNode:单机  */
}

type SdwanSdwanOrderEdgeNewAddrInfoRequest struct {
	Name        string  `json:"name"`                  /*  客户联系人姓名  */
	Mobile      string  `json:"mobile"`                /*  客户联系人电话  */
	SpareMobile *string `json:"spareMobile,omitempty"` /*  客户联系人备用电话  */
	Email       string  `json:"email"`                 /*  客户联系人邮箱  */
	Province    string  `json:"province"`              /*  客户联系人省份  */
	City        string  `json:"city"`                  /*  客户联系人城市  */
	Address     string  `json:"address"`               /*  客户联系人详细地址  */
}

type SdwanSdwanOrderEdgeNewAddedServiceRequest struct {
	Is5G   *bool `json:"is5G,omitempty"`   /*  是否支持5G,默认为否  */
	IsWIFI *bool `json:"isWIFI,omitempty"` /*  是否支持WIFI,默认为否  */
	IsLTE  *bool `json:"isLTE,omitempty"`  /*  是否支持LTE,默认为否  */
}

type SdwanSdwanOrderEdgeNewResponse struct {
	StatusCode  int32                                    `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                  `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                  `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanOrderEdgeNewReturnObjResponse `json:"returnObj"`   /*  返回参数列表  */
	ErrorCode   *string                                  `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码.  */
	Details     *string                                  `json:"details"`     /*  错误明细。一般情况下，会对订单侧(bss)的SDWAN智能网关订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode。  */
	Error       *string                                  `json:"error"`       /*  业务细分码，为product.module.code三段式码.  */
}

type SdwanSdwanOrderEdgeNewReturnObjResponse struct {
	MasterOrderID        *string                                           `json:"masterOrderID"`        /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态。  */
	MasterOrderNO        *string                                           `json:"masterOrderNO"`        /*  订单号  */
	MasterResourceID     *string                                           `json:"masterResourceID"`     /*  主资源ID。  */
	MasterResourceStatus *string                                           `json:"masterResourceStatus"` /*  主资源状态。只有主订单资源会返回  */
	RegionID             *string                                           `json:"regionID"`             /*  资源所属资源池ID  */
	Resources            *SdwanSdwanOrderEdgeNewReturnObjResourcesResponse `json:"resources"`            /*  资源明细列表,参考表resources  */
}

type SdwanSdwanOrderEdgeNewReturnObjResourcesResponse struct {
	ResourceID *string `json:"resourceID"` /*  单项资源的变配、续订、退订等需要该资源项的ID。<br/>比如某个云主机资源作为主资源，对其挂载  */
	OrderID    *string `json:"orderID"`    /*  无需关心  */
	StartTime  int32   `json:"startTime"`  /*  启动时刻，epoch时戳，毫秒精度  */
	ExpireTime int32   `json:"expireTime"` /*  过期时刻，epoch时戳，毫秒精度  */
	CreateTime int32   `json:"createTime"` /*  创建时刻，epoch时戳，毫秒精度  */
	UpdateTime int32   `json:"updateTime"` /*  更新时刻，epoch时戳，毫秒精度  */
	Status     int32   `json:"status"`     /*  资源状态，无需关心。参考masterResourceStatus  */
	IsMaster   *bool   `json:"isMaster"`   /*  是否是主资源项  */
	ItemValue  int32   `json:"itemValue"`  /*  资源规格，网关带宽大小  */
}
