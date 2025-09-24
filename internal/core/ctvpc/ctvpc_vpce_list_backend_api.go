package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVpceListBackendApi
/* 查看终端节点服务后端列表
 */type CtvpcVpceListBackendApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVpceListBackendApi(client *core.CtyunClient) *CtvpcVpceListBackendApi {
	return &CtvpcVpceListBackendApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpce/list-backends",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVpceListBackendApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVpceListBackendRequest) (*CtvpcVpceListBackendResponse, error) {
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
	var resp CtvpcVpceListBackendResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVpceListBackendRequest struct {
	RegionID          string `json:"regionID,omitempty"`          /*  区域ID  */
	EndpointServiceID string `json:"endpointServiceID,omitempty"` /*  终端节点服务 ID  */
}

type CtvpcVpceListBackendResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcVpceListBackendReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVpceListBackendReturnObjResponse struct {
	InstanceID   *string `json:"instanceID,omitempty"`   /*  实例 ID  */
	InstanceType *string `json:"instanceType,omitempty"` /*  实例类型支持：vip / lb / vm / bm / cbm / underlay  */
	InstanceName *string `json:"instanceName,omitempty"` /*  实例名  */
}
