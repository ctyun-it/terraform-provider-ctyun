package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcBindCloudGatewayApi
/* SDWAN绑定云网关 */
type EcEcBindCloudGatewayApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcBindCloudGatewayApi(client *core.CtyunClient) *EcEcBindCloudGatewayApi {
	return &EcEcBindCloudGatewayApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/sdwan/bind-cloud-gateway",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcBindCloudGatewayApi) Do(ctx context.Context, credential core.Credential, req *EcEcBindCloudGatewayRequest) (*EcEcBindCloudGatewayResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcBindCloudGatewayRequest
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
	var resp EcEcBindCloudGatewayResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcBindCloudGatewayRequest struct {
	EcID    string                                `json:"ecID"`    /* 云间高速ID */
	SdwanID string                                `json:"sdwanID"` /* sdwan ID */
	CgwList []*EcEcBindCloudGatewayCgwListRequest `json:"cgwList"` /* 当前需要绑定云网关信息，如果全部解绑则传空list */
}

type EcEcBindCloudGatewayCgwListRequest struct {
	RtbID string `json:"rtbID"` /* 云网关默认路由表id */
	CgwID string `json:"cgwID"` /* 云网关ID */
}

type EcEcBindCloudGatewayResponse struct {
	StatusCode  *int32                                 `json:"statusCode"`  /* 返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败 */
	ErrorCode   *string                                `json:"errorCode"`   /* 业务细分码，为product.module.code三段式码 */
	Message     *string                                `json:"message"`     /* 失败时的错误描述，一般为英文描述 */
	Description *string                                `json:"description"` /* 失败时的错误描述，一般为中文描述 */
	TraceID     *string                                `json:"traceID"`     /* 链路追踪ID */
	ReturnObj   *EcEcBindCloudGatewayReturnObjResponse `json:"returnObj"`   /* 返回参数 */
}

type EcEcBindCloudGatewayReturnObjResponse struct {
	OplogID *string `json:"oplogID"` /* 操作日志ID */
}
