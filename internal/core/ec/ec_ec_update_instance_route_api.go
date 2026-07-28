package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcUpdateInstanceRouteApi
/* 刷新VPC实例路由 */
type EcEcUpdateInstanceRouteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcUpdateInstanceRouteApi(client *core.CtyunClient) *EcEcUpdateInstanceRouteApi {
	return &EcEcUpdateInstanceRouteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/vpc-route/flush",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcUpdateInstanceRouteApi) Do(ctx context.Context, credential core.Credential, req *EcEcUpdateInstanceRouteRequest) (*EcEcUpdateInstanceRouteResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcUpdateInstanceRouteRequest
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
	var resp EcEcUpdateInstanceRouteResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcUpdateInstanceRouteRequest struct {
	EcID  string `json:"ecID"`  /*  云间高速实例ID  */
	VpcID string `json:"vpcID"` /*  vpc ID  */
}

type EcEcUpdateInstanceRouteResponse struct {
	StatusCode  *int32                                    `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                                   `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcUpdateInstanceRouteReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcUpdateInstanceRouteReturnObjResponse struct{}
