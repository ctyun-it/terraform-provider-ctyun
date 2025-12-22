package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanEffectElasticBandwidthPackage1Api
/* 修改智能网关弹性带宽包生效 */
type SdwanSdwanEffectElasticBandwidthPackage1Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanEffectElasticBandwidthPackage1Api(client *core.CtyunClient) *SdwanSdwanEffectElasticBandwidthPackage1Api {
	return &SdwanSdwanEffectElasticBandwidthPackage1Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/effect-elastic-bandwidth-package",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanEffectElasticBandwidthPackage1Api) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanEffectElasticBandwidthPackage1Request) (*SdwanSdwanEffectElasticBandwidthPackage1Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanEffectElasticBandwidthPackage1Request
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
	var resp SdwanSdwanEffectElasticBandwidthPackage1Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanEffectElasticBandwidthPackage1Request struct {
	BandwidthPackageID string  `json:"bandwidthPackageID"`      /*  弹性带宽实例Id  */
	IsEffect           bool    `json:"isEffect"`                /*  是否立即生效,true/false, true表示立即生效，不用带effective,false表示延期生效，必带effectiveTime  */
	EffectiveTime      *string `json:"effectiveTime,omitempty"` /*  生效时间，如果设置时间只能设置当前时间点往后的时间  */
}

type SdwanSdwanEffectElasticBandwidthPackage1Response struct {
	StatusCode  int32                                                        `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                                      `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                                      `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                                      `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*SdwanSdwanEffectElasticBandwidthPackage1ReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                                      `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanSdwanEffectElasticBandwidthPackage1ReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作日志Id,当立即生效时返回会带，否则不会  */
}
