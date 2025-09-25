package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcChangeToNeedApi
/* 资源转按需
 */type CtvpcChangeToNeedApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcChangeToNeedApi(client *core.CtyunClient) *CtvpcChangeToNeedApi {
	return &CtvpcChangeToNeedApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/order/change-to-need",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcChangeToNeedApi) Do(ctx context.Context, credential core.Credential, req *CtvpcChangeToNeedRequest) (*CtvpcChangeToNeedResponse, error) {
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
	var resp CtvpcChangeToNeedResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcChangeToNeedRequest struct {
	ResourceID   string `json:"resourceID,omitempty"`   /*  资源ID  */
	ResourceType string `json:"resourceType,omitempty"` /*  资源类型  */
	RegionID     string `json:"regionID,omitempty"`     /*  区域ID  */
	AutoToNeed   bool   `json:"autoToNeed"`             /*  到期后自动转按需，true:到期后自动转按需,false:取消到期后自动转按需  */
}

type CtvpcChangeToNeedResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	Message     *string `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
}
