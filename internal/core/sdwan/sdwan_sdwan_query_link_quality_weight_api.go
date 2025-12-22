package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanQueryLinkQualityWeightApi
/* 查询链路质量评分 */
type SdwanSdwanQueryLinkQualityWeightApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanQueryLinkQualityWeightApi(client *core.CtyunClient) *SdwanSdwanQueryLinkQualityWeightApi {
	return &SdwanSdwanQueryLinkQualityWeightApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/query-link-quality-weight",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanQueryLinkQualityWeightApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanQueryLinkQualityWeightRequest) (*SdwanSdwanQueryLinkQualityWeightResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanSdwanQueryLinkQualityWeightResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanQueryLinkQualityWeightRequest struct{}

type SdwanSdwanQueryLinkQualityWeightResponse struct {
	StatusCode  int32                                              `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                            `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                            `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                            `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanQueryLinkQualityWeightReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                            `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanSdwanQueryLinkQualityWeightReturnObjResponse struct {
	Result       *SdwanSdwanQueryLinkQualityWeightReturnObjResultResponse `json:"result"`       /*  链路质量参数  */
	TotalCount   int32                                                    `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                                    `json:"currentCount"` /*  当前页数量  */
}

type SdwanSdwanQueryLinkQualityWeightReturnObjResultResponse struct {
	Delay          *string `json:"delay"`          /*  时延权重(delay,jitter,loss三者相加为100)  */
	Jitter         *string `json:"jitter"`         /*  抖动权重  */
	Loss           *string `json:"loss"`           /*  丢包率权重  */
	TenantID       *string `json:"tenantID"`       /*  租户ID  */
	DetectWindow   *string `json:"detectWindow"`   /*  探测窗口  */
	DetectInterval *string `json:"detectInterval"` /*  探测间隔  */
}
