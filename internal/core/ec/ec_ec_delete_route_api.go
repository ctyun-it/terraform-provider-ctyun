package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcDeleteRouteApi
/* 删除路由 */
type EcEcDeleteRouteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcDeleteRouteApi(client *core.CtyunClient) *EcEcDeleteRouteApi {
	return &EcEcDeleteRouteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/route/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcDeleteRouteApi) Do(ctx context.Context, credential core.Credential, req *EcEcDeleteRouteRequest) (*EcEcDeleteRouteResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcDeleteRouteRequest
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
	var resp EcEcDeleteRouteResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcDeleteRouteRequest struct {
	RouteID string `json:"routeID"` /*  路由ID  */
	RtbID   string `json:"rtbID"`   /*  路由表ID  */
}

type EcEcDeleteRouteResponse struct {
	StatusCode  *int32                            `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                           `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                           `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                           `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                           `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcDeleteRouteReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcDeleteRouteReturnObjResponse struct{}
