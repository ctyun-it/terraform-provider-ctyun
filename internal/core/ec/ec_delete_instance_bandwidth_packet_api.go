package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcDeleteInstanceBandwidthPacketApi
/* 删除实例带宽包 */
type EcDeleteInstanceBandwidthPacketApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcDeleteInstanceBandwidthPacketApi(client *core.CtyunClient) *EcDeleteInstanceBandwidthPacketApi {
	return &EcDeleteInstanceBandwidthPacketApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/instance-bandwidth-packet/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *EcDeleteInstanceBandwidthPacketApi) Do(ctx context.Context, credential core.Credential, req *EcDeleteInstanceBandwidthPacketRequest) (*EcDeleteInstanceBandwidthPacketResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcDeleteInstanceBandwidthPacketRequest
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
	var resp EcDeleteInstanceBandwidthPacketResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcDeleteInstanceBandwidthPacketRequest struct {
	IbpID string `json:"ibpID"` /*  实例带宽包ID  */
}

type EcDeleteInstanceBandwidthPacketResponse struct {
	StatusCode  *int32                                            `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                           `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                           `json:"message"`     /*  失败时的错误描述，一般为英文描述   */
	Description *string                                           `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                                           `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcDeleteInstanceBandwidthPacketReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                           `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type EcDeleteInstanceBandwidthPacketReturnObjResponse struct {
	Message *string `json:"message"` /*  更新结果  */
}
