package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcCreateRegionPeerApi
/* 创建实例限速带宽 */
type EcCreateRegionPeerApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcCreateRegionPeerApi(client *core.CtyunClient) *EcCreateRegionPeerApi {
	return &EcCreateRegionPeerApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/region-peer/create",
			ContentType:  "application/json",
		},
	}
}

func (a *EcCreateRegionPeerApi) Do(ctx context.Context, credential core.Credential, req *EcCreateRegionPeerRequest) (*EcCreateRegionPeerResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcCreateRegionPeerRequest
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
	var resp EcCreateRegionPeerResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcCreateRegionPeerRequest struct {
	PeerName   string `json:"peerName"`
	EcID       string `json:"ecID"`
	SrcCgwID   string `json:"srcCgwID"`
	DstCgwID   string `json:"dstCgwID"`
	PacketID   string `json:"packetID"`
	Rate       int32  `json:"rate"`
	RouteLearn *int32 `json:"routeLearn"`
}

type EcCreateRegionPeerResponse struct {
	StatusCode  *int32                               `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                              `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                              `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                              `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                              `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcCreateRegionPeerReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                              `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type EcCreateRegionPeerReturnObjResponse struct {
	PeerID string `json:"peerId"`
}
