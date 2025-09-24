package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVpceDeleteBackendApi
/* 删除终端节点服务后端
 */type CtvpcVpceDeleteBackendApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVpceDeleteBackendApi(client *core.CtyunClient) *CtvpcVpceDeleteBackendApi {
	return &CtvpcVpceDeleteBackendApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/delete-backend",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVpceDeleteBackendApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVpceDeleteBackendRequest) (*CtvpcVpceDeleteBackendResponse, error) {
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
	var resp CtvpcVpceDeleteBackendResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVpceDeleteBackendRequest struct {
	RegionID          string `json:"regionID,omitempty"`          /*  区域ID  */
	EndpointServiceID string `json:"endpointServiceID,omitempty"` /*  终端节点服务 ID  */
	InstanceID        string `json:"instanceID,omitempty"`        /*  实例 id  */
}

type CtvpcVpceDeleteBackendResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcVpceDeleteBackendReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVpceDeleteBackendReturnObjResponse struct {
	EndpointServiceID *string `json:"endpointServiceID,omitempty"` /*  终端节点服务 ID  */
}
