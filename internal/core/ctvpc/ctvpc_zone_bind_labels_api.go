package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcZoneBindLabelsApi
/* 内网 DNS 添加标签
 */type CtvpcZoneBindLabelsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcZoneBindLabelsApi(client *core.CtyunClient) *CtvpcZoneBindLabelsApi {
	return &CtvpcZoneBindLabelsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/private-zone/bind-label",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcZoneBindLabelsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcZoneBindLabelsRequest) (*CtvpcZoneBindLabelsResponse, error) {
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
	var resp CtvpcZoneBindLabelsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcZoneBindLabelsRequest struct {
	RegionID   string `json:"regionID,omitempty"`   /*  区域ID  */
	ZoneID     string `json:"zoneID,omitempty"`     /*  内网 DNS ID  */
	LabelKey   string `json:"labelKey,omitempty"`   /*  标签 key  */
	LabelValue string `json:"labelValue,omitempty"` /*  标签 取值  */
}

type CtvpcZoneBindLabelsResponse struct {
	StatusCode  int32                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcZoneBindLabelsReturnObjResponse `json:"returnObj"`             /*  检查结果  */
	Error       *string                               `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcZoneBindLabelsReturnObjResponse struct {
	ZoneID *string `json:"zoneID,omitempty"` /*  内网 DNS ID  */
}
