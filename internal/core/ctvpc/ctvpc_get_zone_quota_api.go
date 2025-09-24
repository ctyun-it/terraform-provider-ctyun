package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcGetZoneQuotaApi
/* 获取用户 zone 配额
 */type CtvpcGetZoneQuotaApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGetZoneQuotaApi(client *core.CtyunClient) *CtvpcGetZoneQuotaApi {
	return &CtvpcGetZoneQuotaApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/private-zone/quota",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGetZoneQuotaApi) Do(ctx context.Context, credential core.Credential, req *CtvpcGetZoneQuotaRequest) (*CtvpcGetZoneQuotaResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcGetZoneQuotaResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGetZoneQuotaRequest struct {
	RegionID string /*  资源池ID  */
}

type CtvpcGetZoneQuotaResponse struct {
	StatusCode  int32                               `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                             `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                             `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                             `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcGetZoneQuotaReturnObjResponse `json:"returnObj"`             /*  object  */
	Error       *string                             `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcGetZoneQuotaReturnObjResponse struct {
	PrivateZoneQuota int32 `json:"privateZoneQuota"` /*  配额  */
}
