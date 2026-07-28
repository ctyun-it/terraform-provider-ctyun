package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanSetLinkQualityWeightApi
/* 设置链路质量评分 */
type SdwanSdwanSetLinkQualityWeightApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanSetLinkQualityWeightApi(client *core.CtyunClient) *SdwanSdwanSetLinkQualityWeightApi {
	return &SdwanSdwanSetLinkQualityWeightApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/set-link-quality-weight",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanSetLinkQualityWeightApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanSetLinkQualityWeightRequest) (*SdwanSdwanSetLinkQualityWeightResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanSetLinkQualityWeightRequest
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
	var resp SdwanSdwanSetLinkQualityWeightResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanSetLinkQualityWeightRequest struct {
	Delay  string `json:"delay"`  /*  时延权重  */
	Jitter string `json:"jitter"` /*  抖动权重  */
	Loss   string `json:"loss"`   /*  丢包率权重  */
}

type SdwanSdwanSetLinkQualityWeightResponse struct {
	StatusCode  int32                                            `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                          `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                          `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                          `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanSetLinkQualityWeightReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type SdwanSdwanSetLinkQualityWeightReturnObjResponse struct {
	Result *string `json:"result"` /*  配置结果  */
	Error  *string `json:"error"`  /*  业务细分码，为product.module.code三段式码  */
}
