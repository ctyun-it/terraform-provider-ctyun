package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanOfflineEdgeUnlockApi
/* 智能网关离线解锁定 */
type SdwanSdwanOfflineEdgeUnlockApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanOfflineEdgeUnlockApi(client *core.CtyunClient) *SdwanSdwanOfflineEdgeUnlockApi {
	return &SdwanSdwanOfflineEdgeUnlockApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/offline-edge/unlock",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanOfflineEdgeUnlockApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanOfflineEdgeUnlockRequest) (*SdwanSdwanOfflineEdgeUnlockResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanOfflineEdgeUnlockRequest
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
	var resp SdwanSdwanOfflineEdgeUnlockResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanOfflineEdgeUnlockRequest struct {
	SiteID string `json:"siteID"` /*  站点ID  */
}

type SdwanSdwanOfflineEdgeUnlockResponse struct {
	StatusCode  int32   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	OperationID *string `json:"operationID"` /*  操作日志Id  */
	Error       *string `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
