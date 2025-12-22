package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcUpdateWeightsApi
/* 更新sdwan网络实例权重 */
type EcEcUpdateWeightsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcUpdateWeightsApi(client *core.CtyunClient) *EcEcUpdateWeightsApi {
	return &EcEcUpdateWeightsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/sdwan-instance/update-weights",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcUpdateWeightsApi) Do(ctx context.Context, credential core.Credential, req *EcEcUpdateWeightsRequest) (*EcEcUpdateWeightsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcUpdateWeightsRequest
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
	var resp EcEcUpdateWeightsResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcUpdateWeightsRequest struct {
	InstanceID string `json:"instanceID"` /* sdwan ID */
	Weights    int32  `json:"weights"`    /* 权重值 */
}

type EcEcUpdateWeightsResponse struct {
	TraceID     *string                             `json:"traceID"`     /* 链路追踪ID */
	StatusCode  *int32                              `json:"statusCode"`  /* 返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败 */
	ErrorCode   *string                             `json:"errorCode"`   /* 业务细分码，为product.module.code三段式码 */
	Message     *string                             `json:"message"`     /* 失败时的错误描述，一般为英文描述 */
	Description *string                             `json:"description"` /* 失败时的错误描述，一般为中文描述 */
	ReturnObj   *EcEcUpdateWeightsReturnObjResponse `json:"returnObj"`   /* 返回参数 */
}

type EcEcUpdateWeightsReturnObjResponse struct {
	OplogID *string `json:"oplogID"` /* 操作日志id */
}
