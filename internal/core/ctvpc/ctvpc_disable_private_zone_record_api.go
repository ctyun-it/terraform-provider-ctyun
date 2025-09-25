package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDisablePrivateZoneRecordApi
/* 停用解析记录状态
 */type CtvpcDisablePrivateZoneRecordApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDisablePrivateZoneRecordApi(client *core.CtyunClient) *CtvpcDisablePrivateZoneRecordApi {
	return &CtvpcDisablePrivateZoneRecordApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/private-zone-record/disable",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDisablePrivateZoneRecordApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDisablePrivateZoneRecordRequest) (*CtvpcDisablePrivateZoneRecordResponse, error) {
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
	var resp CtvpcDisablePrivateZoneRecordResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDisablePrivateZoneRecordRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID  */
	ZoneRecordID string `json:"zoneRecordID,omitempty"` /*  zoneRecordID  */
}

type CtvpcDisablePrivateZoneRecordResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
