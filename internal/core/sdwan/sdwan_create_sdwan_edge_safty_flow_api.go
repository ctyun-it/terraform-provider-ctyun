package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanCreateSdwanEdgeSaftyFlowApi
/* 安全引流添加智能网关 */
type SdwanCreateSdwanEdgeSaftyFlowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanCreateSdwanEdgeSaftyFlowApi(client *core.CtyunClient) *SdwanCreateSdwanEdgeSaftyFlowApi {
	return &SdwanCreateSdwanEdgeSaftyFlowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-safty-flow/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanCreateSdwanEdgeSaftyFlowApi) Do(ctx context.Context, credential core.Credential, req *SdwanCreateSdwanEdgeSaftyFlowRequest) (*SdwanCreateSdwanEdgeSaftyFlowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanCreateSdwanEdgeSaftyFlowRequest
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
	var resp SdwanCreateSdwanEdgeSaftyFlowResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanCreateSdwanEdgeSaftyFlowRequest struct {
	EdgeID string `json:"edgeID"` /*  智能网关ID  */
}

type SdwanCreateSdwanEdgeSaftyFlowResponse struct {
	StatusCode  int32                                             `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                           `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                           `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                           `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*SdwanCreateSdwanEdgeSaftyFlowReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type SdwanCreateSdwanEdgeSaftyFlowReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作id  */
}
