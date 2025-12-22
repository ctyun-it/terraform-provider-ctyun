package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcUpdateSDWANInstanceApi
/* 更新sdwan网络实例，将sdwan下所有网段同步到云间高速 */
type EcEcUpdateSDWANInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcUpdateSDWANInstanceApi(client *core.CtyunClient) *EcEcUpdateSDWANInstanceApi {
	return &EcEcUpdateSDWANInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/sdwan-instance/update",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcUpdateSDWANInstanceApi) Do(ctx context.Context, credential core.Credential, req *EcEcUpdateSDWANInstanceRequest) (*EcEcUpdateSDWANInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcUpdateSDWANInstanceRequest
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
	var resp EcEcUpdateSDWANInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcUpdateSDWANInstanceRequest struct {
	SdwanID string `json:"sdwanID"` /* sdwan ID */
}

type EcEcUpdateSDWANInstanceResponse struct {
	TraceID     *string                                   `json:"traceID"`     /* 链路追踪ID */
	StatusCode  *int32                                    `json:"statusCode"`  /* 返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败 */
	ErrorCode   *string                                   `json:"errorCode"`   /* 业务细分码，为product.module.code三段式码 */
	Message     *string                                   `json:"message"`     /* 失败时的错误描述，一般为英文描述 */
	Description *string                                   `json:"description"` /* 失败时的错误描述，一般为中文描述 */
	ReturnObj   *EcEcUpdateSDWANInstanceReturnObjResponse `json:"returnObj"`   /* 返回参数 */
}

type EcEcUpdateSDWANInstanceReturnObjResponse struct {
	OplogID *string `json:"oplogID"` /* 操作日志id */
}
