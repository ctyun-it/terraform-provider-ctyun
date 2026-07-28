package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetEdgeCanConfigCountApi
/* 查找智能网关可配置数量 */
type SdwanGetEdgeCanConfigCountApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetEdgeCanConfigCountApi(client *core.CtyunClient) *SdwanGetEdgeCanConfigCountApi {
	return &SdwanGetEdgeCanConfigCountApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge/count",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetEdgeCanConfigCountApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetEdgeCanConfigCountRequest) (*SdwanGetEdgeCanConfigCountResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetEdgeCanConfigCountResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetEdgeCanConfigCountRequest struct{}

type SdwanGetEdgeCanConfigCountResponse struct {
	StatusCode     int32   `json:"statusCode"`     /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode      *string `json:"errorCode"`      /*  业务细分码，为product.module.code三段式码  */
	Message        *string `json:"message"`        /*  失败时的错误描述，一般为英文描述  */
	Description    *string `json:"description"`    /*  失败时的错误描述，一般为中文描述  */
	CanConfigCount int32   `json:"canConfigCount"` /*  可配置数量  */
	Error          *string `json:"error"`          /*  业务细分码，为product.module.code三段式码  */
}
