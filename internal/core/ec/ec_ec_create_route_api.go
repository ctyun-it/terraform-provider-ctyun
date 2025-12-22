package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcCreateRouteApi
/* 创建路由，针对已经加载到云间高速的网络实例，创建一条下一跳为该网络实例的路由。 */
type EcEcCreateRouteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcCreateRouteApi(client *core.CtyunClient) *EcEcCreateRouteApi {
	return &EcEcCreateRouteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/route/create",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcCreateRouteApi) Do(ctx context.Context, credential core.Credential, req *EcEcCreateRouteRequest) (*EcEcCreateRouteResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcCreateRouteRequest
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
	var resp EcEcCreateRouteResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcCreateRouteRequest struct {
	EcID             string  `json:"ecID"`                       /*  云间高速实例ID  */
	CgwID            string  `json:"cgwID"`                      /*  云网关ID  */
	RtbID            string  `json:"rtbID"`                      /*  路由表ID  */
	RouteType        string  `json:"routeType"`                  /*  路由类型<br/>取值范围:<br/>"1":自动学习<br/>"2":自定义（默认）  */
	RouteCIDR        string  `json:"routeCIDR"`                  /*  子网信息  */
	NexthopType      *string `json:"nexthopType,omitempty"`      /*  下一跳实例的类型，如不是黑洞路由则必填<br/>取值范围:<br/>"1":vpc<br/>"2":云专线<br/>"3":SDWAN<br/>"4":VPN<br/>"5":EDS<br/>"20":黑洞路由<br/>"30":跨域连接  */
	NexthopID        *string `json:"nexthopID,omitempty"`        /*  目的实例ID/跨域连接ID，如不是黑洞路由则必填  */
	RouteDescription *string `json:"routeDescription,omitempty"` /*  路由描述信息  */
	IPVersion        string  `json:"IPVersion"`                  /*  子网类型<br/>取值范围:<br/>"1":IPv4类型<br/>"2":IPv6类型  */
	IsBlackholeRoute *bool   `json:"isBlackholeRoute,omitempty"` /*  是否是黑洞路由, 如果选择true，则nexthopType、nexthopID字段可不传；<br/>取值范围:<br/>false: 否<br/>true：是<br/> 默认false(否)  */
}

type EcEcCreateRouteResponse struct {
	StatusCode  *int32                            `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                           `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                           `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                           `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                           `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcCreateRouteReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcCreateRouteReturnObjResponse struct {
	RouteID *string `json:"routeID"` /*  路由ID  */
}
