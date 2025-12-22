package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanEffectElasticBandwidthPackageApi
/* 删除智能网关弹性带宽包 */
type SdwanSdwanEffectElasticBandwidthPackageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanEffectElasticBandwidthPackageApi(client *core.CtyunClient) *SdwanSdwanEffectElasticBandwidthPackageApi {
	return &SdwanSdwanEffectElasticBandwidthPackageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/elastic-bandwidth-package/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanEffectElasticBandwidthPackageApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanEffectElasticBandwidthPackageRequest) (*SdwanSdwanEffectElasticBandwidthPackageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanEffectElasticBandwidthPackageRequest
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
	var resp SdwanSdwanEffectElasticBandwidthPackageResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanEffectElasticBandwidthPackageRequest struct {
	BandwidthPackageID string `json:"bandwidthPackageID"` /*  弹性带宽实例Id  */
}

type SdwanSdwanEffectElasticBandwidthPackageResponse struct {
	StatusCode  int32     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*string `json:"returnObj"`   /*  返回参数  */
	Error       *string   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
