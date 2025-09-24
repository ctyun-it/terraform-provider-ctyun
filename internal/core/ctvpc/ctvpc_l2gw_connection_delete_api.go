package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcL2gwConnectionDeleteApi
/* 删除l2gw_connection
 */type CtvpcL2gwConnectionDeleteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcL2gwConnectionDeleteApi(client *core.CtyunClient) *CtvpcL2gwConnectionDeleteApi {
	return &CtvpcL2gwConnectionDeleteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/l2gw_connection/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcL2gwConnectionDeleteApi) Do(ctx context.Context, credential core.Credential, req *CtvpcL2gwConnectionDeleteRequest) (*CtvpcL2gwConnectionDeleteResponse, error) {
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
	var resp CtvpcL2gwConnectionDeleteResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcL2gwConnectionDeleteRequest struct {
	RegionID       string `json:"regionID,omitempty"`       /*  资源池 ID  */
	L2ConnectionID string `json:"l2ConnectionID,omitempty"` /*  l2gw_connection 的 ID  */
}

type CtvpcL2gwConnectionDeleteResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
