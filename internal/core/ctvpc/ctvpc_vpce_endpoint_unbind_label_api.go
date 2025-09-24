package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVpceEndpointUnbindLabelApi
/* 终端节点解绑标签
 */type CtvpcVpceEndpointUnbindLabelApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVpceEndpointUnbindLabelApi(client *core.CtyunClient) *CtvpcVpceEndpointUnbindLabelApi {
	return &CtvpcVpceEndpointUnbindLabelApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/endpoint-unbind-label",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVpceEndpointUnbindLabelApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVpceEndpointUnbindLabelRequest) (*CtvpcVpceEndpointUnbindLabelResponse, error) {
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
	var resp CtvpcVpceEndpointUnbindLabelResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVpceEndpointUnbindLabelRequest struct {
	RegionID   string `json:"regionID,omitempty"`   /*  区域ID  */
	EndpointID string `json:"endpointID,omitempty"` /*  终端节点 ID  */
	LabelID    string `json:"labelID,omitempty"`    /*  标签 id  */
}

type CtvpcVpceEndpointUnbindLabelResponse struct {
	StatusCode  int32                                          `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                        `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                        `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                        `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcVpceEndpointUnbindLabelReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVpceEndpointUnbindLabelReturnObjResponse struct {
	EndpointID *string `json:"endpointID,omitempty"` /*  终端节点 ID  */
}
