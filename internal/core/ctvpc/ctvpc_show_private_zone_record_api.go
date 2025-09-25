package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowPrivateZoneRecordApi
/* 内网 DNS 记录详情
 */type CtvpcShowPrivateZoneRecordApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowPrivateZoneRecordApi(client *core.CtyunClient) *CtvpcShowPrivateZoneRecordApi {
	return &CtvpcShowPrivateZoneRecordApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/private-zone-record/query",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowPrivateZoneRecordApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowPrivateZoneRecordRequest) (*CtvpcShowPrivateZoneRecordResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("zoneRecordID", req.ZoneRecordID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowPrivateZoneRecordResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowPrivateZoneRecordRequest struct {
	RegionID     string /*  资源池ID  */
	ZoneRecordID string /*  zoneRecordID  */
}

type CtvpcShowPrivateZoneRecordResponse struct {
	StatusCode  int32                                        `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                      `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                      `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                      `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowPrivateZoneRecordReturnObjResponse `json:"returnObj"`             /*  object  */
	Error       *string                                      `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowPrivateZoneRecordReturnObjResponse struct {
	ZoneRecordID *string   `json:"zoneRecordID,omitempty"` /*  名称  */
	Name         *string   `json:"name,omitempty"`         /*  zoneRecord名称  */
	Description  *string   `json:"description,omitempty"`  /*  描述  */
	ZoneID       *string   `json:"zoneID,omitempty"`       /*  名称  */
	ZoneName     *string   `json:"zoneName,omitempty"`     /*  zone名称  */
	RawType      *string   `json:"type,omitempty"`         /*  zone record type: A,CNAME,MX,AAAA,TXT  */
	Value        []*string `json:"value"`                  /*  记录值  */
	TTL          int32     `json:"TTL"`                    /*  zone ttl, default is 300  */
	CreatedAt    *string   `json:"createdAt,omitempty"`    /*  创建时间  */
	UpdatedAt    *string   `json:"updatedAt,omitempty"`    /*  更新时间  */
}
