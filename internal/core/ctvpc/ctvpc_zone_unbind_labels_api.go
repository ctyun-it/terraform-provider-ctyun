package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcZoneUnbindLabelsApi
/* 内网 DNS 移除标签
 */type CtvpcZoneUnbindLabelsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcZoneUnbindLabelsApi(client *core.CtyunClient) *CtvpcZoneUnbindLabelsApi {
	return &CtvpcZoneUnbindLabelsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/private-zone/unbind-label",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcZoneUnbindLabelsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcZoneUnbindLabelsRequest) (*CtvpcZoneUnbindLabelsResponse, error) {
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
	var resp CtvpcZoneUnbindLabelsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcZoneUnbindLabelsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  区域ID  */
	ZoneID   string `json:"zoneID,omitempty"`   /*  内网 DNS ID  */
	LabelID  string `json:"labelID,omitempty"`  /*  标签ID  */
}

type CtvpcZoneUnbindLabelsResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcZoneUnbindLabelsReturnObjResponse `json:"returnObj"`             /*  检查结果  */
	Error       *string                                 `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcZoneUnbindLabelsReturnObjResponse struct {
	ZoneID *string `json:"zoneID,omitempty"` /*  内网 DNS ID  */
}
