package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanDeleteSdwanAppApi
/* 应用删除 */
type SdwanDeleteSdwanAppApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanDeleteSdwanAppApi(client *core.CtyunClient) *SdwanDeleteSdwanAppApi {
	return &SdwanDeleteSdwanAppApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/app/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanDeleteSdwanAppApi) Do(ctx context.Context, credential core.Credential, req *SdwanDeleteSdwanAppRequest) (*SdwanDeleteSdwanAppResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanDeleteSdwanAppRequest
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
	var resp SdwanDeleteSdwanAppResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanDeleteSdwanAppRequest struct {
	AppIds []string `json:"appIds"` /*  应用id列表  ，值类型为string  */
}

type SdwanDeleteSdwanAppResponse struct {
	StatusCode  int32                                   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                 `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                 `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                 `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*SdwanDeleteSdwanAppReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                 `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanDeleteSdwanAppReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作id  */
}
