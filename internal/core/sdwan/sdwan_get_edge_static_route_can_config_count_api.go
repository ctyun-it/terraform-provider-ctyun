package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetEdgeStaticRouteCanConfigCountApi
/* 查询智能网关可配置静态路由数量 */
type SdwanGetEdgeStaticRouteCanConfigCountApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetEdgeStaticRouteCanConfigCountApi(client *core.CtyunClient) *SdwanGetEdgeStaticRouteCanConfigCountApi {
	return &SdwanGetEdgeStaticRouteCanConfigCountApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-static-route/count",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetEdgeStaticRouteCanConfigCountApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetEdgeStaticRouteCanConfigCountRequest) (*SdwanGetEdgeStaticRouteCanConfigCountResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetEdgeStaticRouteCanConfigCountResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetEdgeStaticRouteCanConfigCountRequest struct {
	EdgeID string `json:"edgeID"` /*  智能网关Id  */
}

type SdwanGetEdgeStaticRouteCanConfigCountResponse struct {
	StatusCode     int32   `json:"statusCode"`     /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode      *string `json:"errorCode"`      /*  业务细分码，为product.module.code三段式码  */
	Message        *string `json:"message"`        /*  失败时的错误描述，一般为英文描述  */
	Description    *string `json:"description"`    /*  失败时的错误描述，一般为中文描述  */
	CanConfigCount int32   `json:"canConfigCount"` /*  可配置数量  */
	Error          *string `json:"error"`          /*  业务细分码，为product.module.code三段式码  */
}
