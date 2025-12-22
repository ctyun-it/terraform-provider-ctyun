package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanCreateEdgeBindSdwanApi
/* edge绑定云连接网 */
type SdwanCreateEdgeBindSdwanApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanCreateEdgeBindSdwanApi(client *core.CtyunClient) *SdwanCreateEdgeBindSdwanApi {
	return &SdwanCreateEdgeBindSdwanApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-bind-sdwan/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanCreateEdgeBindSdwanApi) Do(ctx context.Context, credential core.Credential, req *SdwanCreateEdgeBindSdwanRequest) (*SdwanCreateEdgeBindSdwanResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanCreateEdgeBindSdwanRequest
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
	var resp SdwanCreateEdgeBindSdwanResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanCreateEdgeBindSdwanRequest struct {
	EdgeID      string  `json:"edgeID"`                /*  智能网关Id  */
	SdwanID     string  `json:"sdwanID"`               /*  sdwan的id  */
	CloudHighID *string `json:"cloudHighID,omitempty"` /*  云间高速Id  */
}

type SdwanCreateEdgeBindSdwanResponse struct {
	StatusCode  int32                                        `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                      `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                      `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                      `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*SdwanCreateEdgeBindSdwanReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                      `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanCreateEdgeBindSdwanReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作日志Id  */
}
