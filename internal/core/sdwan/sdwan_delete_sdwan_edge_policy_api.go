package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanDeleteSdwanEdgePolicyApi
/* 智能选路策略删除 */
type SdwanDeleteSdwanEdgePolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanDeleteSdwanEdgePolicyApi(client *core.CtyunClient) *SdwanDeleteSdwanEdgePolicyApi {
	return &SdwanDeleteSdwanEdgePolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-policy/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanDeleteSdwanEdgePolicyApi) Do(ctx context.Context, credential core.Credential, req *SdwanDeleteSdwanEdgePolicyRequest) (*SdwanDeleteSdwanEdgePolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanDeleteSdwanEdgePolicyRequest
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
	var resp SdwanDeleteSdwanEdgePolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanDeleteSdwanEdgePolicyRequest struct {
	PolicyID string `json:"policyID"` /*  策略id  */
	EdgeID   string `json:"edgeID"`   /*  智能网关ID  */
}

type SdwanDeleteSdwanEdgePolicyResponse struct {
	StatusCode  int32                                        `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                      `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                      `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                      `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanDeleteSdwanEdgePolicyReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                      `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanDeleteSdwanEdgePolicyReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作id  */
}
