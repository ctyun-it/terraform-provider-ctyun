package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcCreateRouteTableAutoLearnInstanceApi
/* 查询云专线信息列表 */
type EcEcCreateRouteTableAutoLearnInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcCreateRouteTableAutoLearnInstanceApi(client *core.CtyunClient) *EcEcCreateRouteTableAutoLearnInstanceApi {
	return &EcEcCreateRouteTableAutoLearnInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/cda/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcCreateRouteTableAutoLearnInstanceApi) Do(ctx context.Context, credential core.Credential, req *EcEcCreateRouteTableAutoLearnInstanceRequest) (*EcEcCreateRouteTableAutoLearnInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("account", req.Account)
	ctReq.AddParam("resourcePoolID", req.ResourcePoolID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcCreateRouteTableAutoLearnInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcCreateRouteTableAutoLearnInstanceRequest struct {
	Account        string `json:"account"`        /*  用户账号  */
	ResourcePoolID string `json:"resourcePoolID"` /*  资源池ID  */
}

type EcEcCreateRouteTableAutoLearnInstanceResponse struct {
	StatusCode  *int32  `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *string `json:"returnObj"`   /*  返回参数  */
}
