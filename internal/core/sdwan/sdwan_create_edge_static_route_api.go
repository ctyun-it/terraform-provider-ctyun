package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanCreateEdgeStaticRouteApi
/* 增加静态路由 */
type SdwanCreateEdgeStaticRouteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanCreateEdgeStaticRouteApi(client *core.CtyunClient) *SdwanCreateEdgeStaticRouteApi {
	return &SdwanCreateEdgeStaticRouteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-static-route/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanCreateEdgeStaticRouteApi) Do(ctx context.Context, credential core.Credential, req *SdwanCreateEdgeStaticRouteRequest) (*SdwanCreateEdgeStaticRouteResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanCreateEdgeStaticRouteRequest
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
	var resp SdwanCreateEdgeStaticRouteResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanCreateEdgeStaticRouteRequest struct {
	EdgeID      string                                          `json:"edgeID"`                /*  智能网关Id  */
	Protocol    string                                          `json:"protocol"`              /*  本参数表示ip版本<br/><br/>取值范围:<br/>IPv4:ipv4<br/>IPv6:ipv6  */
	DstCIDR     string                                          `json:"dstCIDR"`               /*  目的子网  */
	NextHopList []*SdwanCreateEdgeStaticRouteNextHopListRequest `json:"nextHopList,omitempty"` /*  下一跳信息  */
}

type SdwanCreateEdgeStaticRouteNextHopListRequest struct {
	NextHop       string `json:"nextHop"`       /*  下一跳地址  */
	Priority      int32  `json:"priority"`      /*  优先级  */
	InterfaceName string `json:"interfaceName"` /*  出接口  */
}

type SdwanCreateEdgeStaticRouteResponse struct {
	StatusCode  int32     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*string `json:"returnObj"`   /*  返回参数  */
	Error       *string   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
