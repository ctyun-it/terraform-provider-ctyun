package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsDeleteCommandApi
/* 调用此接口可删除用户自己创建的云助手命令
 */type CtecsDeleteCommandApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsDeleteCommandApi(client *core.CtyunClient) *CtecsDeleteCommandApi {
	return &CtecsDeleteCommandApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cloud-assistant/delete-command",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsDeleteCommandApi) Do(ctx context.Context, credential core.Credential, req *CtecsDeleteCommandRequest) (*CtecsDeleteCommandResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsDeleteCommandResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsDeleteCommandRequest struct {
	RegionID  string `json:"regionID,omitempty"`  /*  资源池ID  */
	CommandID string `json:"commandID,omitempty"` /*  命令ID  */
}

type CtecsDeleteCommandResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码（800 为成功，900 为失败）  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，详见错误码说明  */
	Message     string `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
}
