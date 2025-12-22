package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanOrderElasticBandwidthPackageNewApi
/* 支持按需天订购智能网关弹性带宽包。 */
type SdwanSdwanOrderElasticBandwidthPackageNewApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanOrderElasticBandwidthPackageNewApi(client *core.CtyunClient) *SdwanSdwanOrderElasticBandwidthPackageNewApi {
	return &SdwanSdwanOrderElasticBandwidthPackageNewApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/elastic-bandwidth-package/new",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanOrderElasticBandwidthPackageNewApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanOrderElasticBandwidthPackageNewRequest) (*SdwanSdwanOrderElasticBandwidthPackageNewResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanOrderElasticBandwidthPackageNewRequest
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
	var resp SdwanSdwanOrderElasticBandwidthPackageNewResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanOrderElasticBandwidthPackageNewRequest struct {
	ClientToken   *string `json:"clientToken,omitempty"`   /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一。  */
	EdgeName      string  `json:"edgeName"`                /*  edge的名称  */
	EdgeID        string  `json:"edgeID"`                  /*  edge的ID  */
	Bandwidth     int32   `json:"bandwidth"`               /*  智能网关APP 带宽  */
	OnDemand      bool    `json:"onDemand"`                /*  是否按需下单。默认为false  */
	CycleType     *string `json:"cycleType,omitempty"`     /*  本参数表示包周期类型<br/><br/>取值范围:<br/>DAY:按天<br/>onDemand为False时，必须指定。  */
	CycleCount    int64   `json:"cycleCount"`              /*  包周期数。onDemand为False时必须指定。周期最大长度不能超过30天  */
	EffectiveTime *string `json:"effectiveTime,omitempty"` /*  生效时间，如果没有设置实时生效，如果设置时间只能设置当前时间之后的时间点  */
}

type SdwanSdwanOrderElasticBandwidthPackageNewResponse struct {
	StatusCode  int32                                                       `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                                     `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                                     `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanOrderElasticBandwidthPackageNewReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                                     `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Details     *string                                                     `json:"details"`     /*  错误明细。一般情况下，会对订单侧(bss)的SDWAN智能网关APP订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode。  */
	Error       *string                                                     `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanSdwanOrderElasticBandwidthPackageNewReturnObjResponse struct {
	MasterOrderID        *string                                                              `json:"masterOrderID"`        /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态。  */
	MasterOrderNO        *string                                                              `json:"masterOrderNO"`        /*  订单号  */
	MasterResourceID     *string                                                              `json:"masterResourceID"`     /*  主资源ID  */
	MasterResourceStatus *string                                                              `json:"masterResourceStatus"` /*  主资源状态。只有主订单资源会返回  */
	RegionID             *string                                                              `json:"regionID"`             /*  资源所属资源池ID  */
	Resources            *SdwanSdwanOrderElasticBandwidthPackageNewReturnObjResourcesResponse `json:"resources"`            /*  资源明细列表,参考表Resource  */
}

type SdwanSdwanOrderElasticBandwidthPackageNewReturnObjResourcesResponse struct {
	ResourceID *string `json:"resourceID"` /*  单项资源的变配、续订、退订等需要该资源项的ID。  */
	OrderID    *string `json:"orderID"`    /*  无需关心  */
	StartTime  int32   `json:"startTime"`  /*  启动时刻，epoch时戳，毫秒精度  */
	ExpireTime int32   `json:"expireTime"` /*  过期时刻，epoch时戳，毫秒精度  */
	CreateTime int32   `json:"createTime"` /*  创建时刻，epoch时戳，毫秒精度  */
	UpdateTime int32   `json:"updateTime"` /*  更新时刻，epoch时戳，毫秒精度  */
	Status     int32   `json:"status"`     /*  资源状态，无需关心。参考masterResourceStatus  */
	IsMaster   *bool   `json:"isMaster"`   /*  是否是主资源项  */
	ItemValue  int32   `json:"itemValue"`  /*  资源规格，带宽大小  */
}
