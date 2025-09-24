package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcTerminateCycleApi
/* 资源终止包周期
 */type CtvpcTerminateCycleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcTerminateCycleApi(client *core.CtyunClient) *CtvpcTerminateCycleApi {
	return &CtvpcTerminateCycleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/order/terminate-cycle",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcTerminateCycleApi) Do(ctx context.Context, credential core.Credential, req *CtvpcTerminateCycleRequest) (*CtvpcTerminateCycleResponse, error) {
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
	var resp CtvpcTerminateCycleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcTerminateCycleRequest struct {
	ResourceID   string `json:"resourceID,omitempty"`   /*  资源ID  */
	ResourceType string `json:"resourceType,omitempty"` /*  资源类型  */
	RegionID     string `json:"regionID,omitempty"`     /*  区域ID  */
}

type CtvpcTerminateCycleResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	Message     *string `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
}
