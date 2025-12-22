package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanUpdateEdgeApi
/* 修改智能网关设备信息 */
type SdwanUpdateEdgeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanUpdateEdgeApi(client *core.CtyunClient) *SdwanUpdateEdgeApi {
	return &SdwanUpdateEdgeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge/update",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanUpdateEdgeApi) Do(ctx context.Context, credential core.Credential, req *SdwanUpdateEdgeRequest) (*SdwanUpdateEdgeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanUpdateEdgeRequest
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
	var resp SdwanUpdateEdgeResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanUpdateEdgeRequest struct {
	EdgeName    string  `json:"edgeName"`              /*  设备名称  */
	Description *string `json:"description,omitempty"` /*  edge描述  */
	EdgeID      string  `json:"edgeID"`                /*  智能网关ID  */
}

type SdwanUpdateEdgeResponse struct {
	StatusCode  int32   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	EdgeID      *string `json:"edgeID"`      /*  edge id  */
	Error       *string `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
