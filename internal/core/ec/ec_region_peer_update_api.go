package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcRegionPeerUpdateApi
/* 修改云间高速 */
type EcRegionPeerUpdateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcRegionPeerUpdateApi(client *core.CtyunClient) *EcRegionPeerUpdateApi {
	return &EcRegionPeerUpdateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/region-peer/update",
			ContentType:  "application/json",
		},
	}
}

func (a *EcRegionPeerUpdateApi) Do(ctx context.Context, credential core.Credential, req *EcRegionPeerUpdateRequest) (*EcRegionPeerUpdateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcRegionPeerUpdateRequest
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
	var resp EcRegionPeerUpdateResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcRegionPeerUpdateRequest struct {
	PeerID string `json:"peerID"`
	Rate   int32  `json:"rate"`
}

type EcRegionPeerUpdateResponse struct {
	StatusCode  *int32                               `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                              `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                              `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                              `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                              `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcRegionPeerUpdateReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                              `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type EcRegionPeerUpdateReturnObjResponse struct {
}
