package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanListElasticBandwidthPackageApi
/* 查询智能网关弹性带宽包 */
type SdwanSdwanListElasticBandwidthPackageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanListElasticBandwidthPackageApi(client *core.CtyunClient) *SdwanSdwanListElasticBandwidthPackageApi {
	return &SdwanSdwanListElasticBandwidthPackageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/elastic-bandwidth-package/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanListElasticBandwidthPackageApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanListElasticBandwidthPackageRequest) (*SdwanSdwanListElasticBandwidthPackageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanSdwanListElasticBandwidthPackageResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanListElasticBandwidthPackageRequest struct {
	EdgeID string `json:"edgeID"` /*  智能网关ID  */
}

type SdwanSdwanListElasticBandwidthPackageResponse struct {
	StatusCode  int32                                                   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                                 `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                                 `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                                 `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanListElasticBandwidthPackageReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                                 `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanSdwanListElasticBandwidthPackageReturnObjResponse struct {
	BandwidthPackageID *string `json:"bandwidthPackageID"` /*  弹性带宽实例Id  */
	EdgeID             *string `json:"edgeID"`             /*  智能网关Id  */
	CustomerID         *string `json:"customerID"`         /*  用户ID  */
	Status             *string `json:"status"`             /*  带宽状态  */
	Bandwidth          int32   `json:"bandwidth"`          /*  带宽大小  */
	ResourceID         *string `json:"resourceID"`         /*  资源ID  */
	EffectiveTime      *string `json:"effectiveTime"`      /*  生效时间  */
	ExpireTime         *string `json:"expireTime"`         /*  到期时间  */
}
