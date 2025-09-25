package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcNewPrivateZoneRecordListApi
/* 内网 DNS 记录列表
 */type CtvpcNewPrivateZoneRecordListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcNewPrivateZoneRecordListApi(client *core.CtyunClient) *CtvpcNewPrivateZoneRecordListApi {
	return &CtvpcNewPrivateZoneRecordListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/private-zone-record/new-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcNewPrivateZoneRecordListApi) Do(ctx context.Context, credential core.Credential, req *CtvpcNewPrivateZoneRecordListRequest) (*CtvpcNewPrivateZoneRecordListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ZoneRecordName != nil {
		ctReq.AddParam("zoneRecordName", *req.ZoneRecordName)
	}
	if req.ZoneID != nil {
		ctReq.AddParam("zoneID", *req.ZoneID)
	}
	if req.ZoneRecordID != nil {
		ctReq.AddParam("zoneRecordID", *req.ZoneRecordID)
	}
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
	}
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
	var resp CtvpcNewPrivateZoneRecordListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcNewPrivateZoneRecordListRequest struct {
	RegionID       string  /*  资源池ID  */
	ZoneRecordName *string /*  dns 记录集名字  */
	ZoneID         *string /*  zoneID  */
	ZoneRecordID   *string /*  zoneRecordID  */
	PageNumber     int32   /*  列表的页码，默认值为1。  */
	PageNo         int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize       int32   /*  分页查询时每页的行数，最大值为50，默认值为10。  */
}

type CtvpcNewPrivateZoneRecordListResponse struct {
	StatusCode  int32                                           `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                         `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                         `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                         `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcNewPrivateZoneRecordListReturnObjResponse `json:"returnObj"`             /*  object  */
	Error       *string                                         `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcNewPrivateZoneRecordListReturnObjResponse struct {
	ZoneRecords  []*CtvpcNewPrivateZoneRecordListReturnObjZoneRecordsResponse `json:"zoneRecords"`  /*  dns记录集  */
	TotalCount   int32                                                        `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                                        `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                                        `json:"totalPage"`    /*  总页数  */
}

type CtvpcNewPrivateZoneRecordListReturnObjZoneRecordsResponse struct {
	ZoneRecordID *string `json:"zoneRecordID,omitempty"` /*  名称  */
	Name         *string `json:"name,omitempty"`         /*  zone record名称  */
}
