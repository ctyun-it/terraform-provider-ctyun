package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVpceServiceBindLabelApi
/* 终端节点服务绑定标签
 */type CtvpcVpceServiceBindLabelApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVpceServiceBindLabelApi(client *core.CtyunClient) *CtvpcVpceServiceBindLabelApi {
	return &CtvpcVpceServiceBindLabelApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/service-bind-label",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVpceServiceBindLabelApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVpceServiceBindLabelRequest) (*CtvpcVpceServiceBindLabelResponse, error) {
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
	var resp CtvpcVpceServiceBindLabelResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVpceServiceBindLabelRequest struct {
	RegionID   string `json:"regionID,omitempty"`   /*  区域ID  */
	LabelKey   string `json:"labelKey,omitempty"`   /*  标签 key  */
	LabelValue string `json:"labelValue,omitempty"` /*  标签 取值  */
}

type CtvpcVpceServiceBindLabelResponse struct {
	StatusCode  int32                                       `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcVpceServiceBindLabelReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVpceServiceBindLabelReturnObjResponse struct {
	EndpointServiceID *string `json:"endpointServiceID,omitempty"` /*  终端节点服务 ID  */
}
