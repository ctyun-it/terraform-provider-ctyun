package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtnatDeletePrivatenatIPApi
/* 删除中转IP
 */type CtnatDeletePrivatenatIPApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatDeletePrivatenatIPApi(client *core.CtyunClient) *CtnatDeletePrivatenatIPApi {
	return &CtnatDeletePrivatenatIPApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/privatenat/delete-ip",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatDeletePrivatenatIPApi) Do(ctx context.Context, credential core.Credential, req *CtnatDeletePrivatenatIPRequest) (*CtnatDeletePrivatenatIPResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtnatDeletePrivatenatIPRequest
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
	var resp CtnatDeletePrivatenatIPResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatDeletePrivatenatIPRequest struct {
	RegionID     string   `json:"regionID,omitempty"`     /*  区域id  */
	NatGatewayID string   `json:"natGatewayID,omitempty"` /*  私网NATID  */
	Addresses    []string `json:"addresses"`              /*  [192.168.8.56]  */
}

type CtnatDeletePrivatenatIPResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string `json:"message"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string `json:"description"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string `json:"errorCode"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
