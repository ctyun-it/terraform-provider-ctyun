package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListSharedZoneRecordsApi
/* 获取共享 zone record
 */type CtvpcListSharedZoneRecordsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListSharedZoneRecordsApi(client *core.CtyunClient) *CtvpcListSharedZoneRecordsApi {
	return &CtvpcListSharedZoneRecordsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/private-zone-record/list-shared-records",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListSharedZoneRecordsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListSharedZoneRecordsRequest) (*CtvpcListSharedZoneRecordsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcListSharedZoneRecordsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListSharedZoneRecordsRequest struct {
	RegionID string /*  资源池ID  */
	PageNo   int32  /*  列表的页码，默认值为 1  */
	PageSize int32  /*  分页查询时每页的行数，最大值为50，默认值为10。  */
}

type CtvpcListSharedZoneRecordsResponse struct {
	StatusCode  int32                                        `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                      `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                      `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                      `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcListSharedZoneRecordsReturnObjResponse `json:"returnObj"`             /*  object  */
	Error       *string                                      `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListSharedZoneRecordsReturnObjResponse struct {
	Results      []*CtvpcListSharedZoneRecordsReturnObjResultsResponse `json:"results"`      /*  dns 记录  */
	TotalCount   int32                                                 `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                                 `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                                 `json:"totalPage"`    /*  总页数  */
}

type CtvpcListSharedZoneRecordsReturnObjResultsResponse struct {
	ZoneRecordDescription *string `json:"zoneRecordDescription,omitempty"` /*  描述  */
	ZoneRecordType        *string `json:"zoneRecordType,omitempty"`        /*  记录类型: SRV / CNAME / TXT / AAAA / A / TXT / PTR  */
	ZoneRecordName        *string `json:"zoneRecordName,omitempty"`        /*  记录名字  */
	TTL                   int32   `json:"TTL"`                             /*  zone ttl, default is 300  */
}
