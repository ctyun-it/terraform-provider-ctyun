package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcEnablePrivateZoneRecordApi
/* 开启解析记录状态
 */type CtvpcEnablePrivateZoneRecordApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcEnablePrivateZoneRecordApi(client *core.CtyunClient) *CtvpcEnablePrivateZoneRecordApi {
	return &CtvpcEnablePrivateZoneRecordApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/private-zone-record/enable",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcEnablePrivateZoneRecordApi) Do(ctx context.Context, credential core.Credential, req *CtvpcEnablePrivateZoneRecordRequest) (*CtvpcEnablePrivateZoneRecordResponse, error) {
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
	var resp CtvpcEnablePrivateZoneRecordResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcEnablePrivateZoneRecordRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID  */
	ZoneRecordID string `json:"zoneRecordID,omitempty"` /*  zoneRecordID  */
}

type CtvpcEnablePrivateZoneRecordResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
