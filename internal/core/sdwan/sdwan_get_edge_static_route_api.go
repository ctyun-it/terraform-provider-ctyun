package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetEdgeStaticRouteApi
/* 查询静态路由 */
type SdwanGetEdgeStaticRouteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetEdgeStaticRouteApi(client *core.CtyunClient) *SdwanGetEdgeStaticRouteApi {
	return &SdwanGetEdgeStaticRouteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-static-route/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetEdgeStaticRouteApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetEdgeStaticRouteRequest) (*SdwanGetEdgeStaticRouteResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.Search != nil && *req.Search != "" {
		ctReq.AddParam("search", *req.Search)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetEdgeStaticRouteResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetEdgeStaticRouteRequest struct {
	EdgeID   string  `json:"edgeID"`           /*  智能网关Id  */
	PageNo   int32   `json:"pageNo"`           /*  页数  */
	PageSize int32   `json:"pageSize"`         /*  页大小  */
	Search   *string `json:"search,omitempty"` /*  模糊查询  */
}

type SdwanGetEdgeStaticRouteResponse struct {
	StatusCode    int32   `json:"statusCode"`    /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode     *string `json:"errorCode"`     /*  业务细分码，为product.module.code三段式码  */
	Message       *string `json:"message"`       /*  失败时的错误描述，一般为英文描述  */
	Description   *string `json:"description"`   /*  失败时的错误描述，一般为中文描述  */
	TotalCount    int32   `json:"totalCount"`    /*  总数量  */
	CurrentCount  int32   `json:"currentCount"`  /*  当前页数量  */
	RouteID       *string `json:"routeID"`       /*  静态路由Id  */
	Protocol      *string `json:"protocol"`      /*  本参数表示ip协议版本<br/>取值范围：<br/>IPv4：IPv4协议<br/>IPv6:  IPv6协议  */
	EdgeCIDR      *string `json:"edgeCIDR"`      /*  子网  */
	Status        *string `json:"status"`        /*  route的status  */
	NextHop       *string `json:"nextHop"`       /*  下一跳地址  */
	Priority      int32   `json:"priority"`      /*  优先级  */
	InterfaceName *string `json:"interfaceName"` /*  出接口  */
	Error         *string `json:"error"`         /*  业务细分码，为product.module.code三段式码  */
}
