package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanCreateSdwanApi
/* 创建SD-WAN */
type SdwanCreateSdwanApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanCreateSdwanApi(client *core.CtyunClient) *SdwanCreateSdwanApi {
	return &SdwanCreateSdwanApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanCreateSdwanApi) Do(ctx context.Context, credential core.Credential, req *SdwanCreateSdwanRequest) (*SdwanCreateSdwanResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanCreateSdwanRequest
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
	var resp SdwanCreateSdwanResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanCreateSdwanRequest struct {
	SdwanName   string  `json:"sdwanName"`             /*  名称  */
	ProjectID   string  `json:"projectID"`             /*  企业项目  */
	Description *string `json:"description,omitempty"` /*  描述  */
}

type SdwanCreateSdwanResponse struct {
	StatusCode  int32   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	OperationID *string `json:"operationID"` /*  操作id  */
	Error       *string `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
