package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcNewPrivateZoneListApi
/* 内网 DNS 列表
 */type CtvpcNewPrivateZoneListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcNewPrivateZoneListApi(client *core.CtyunClient) *CtvpcNewPrivateZoneListApi {
	return &CtvpcNewPrivateZoneListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/private-zone/new-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcNewPrivateZoneListApi) Do(ctx context.Context, credential core.Credential, req *CtvpcNewPrivateZoneListRequest) (*CtvpcNewPrivateZoneListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.ClientToken != nil {
		ctReq.AddParam("clientToken", *req.ClientToken)
	}
	ctReq.AddParam("regionID", req.RegionID)
	if req.ZoneID != nil {
		ctReq.AddParam("zoneID", *req.ZoneID)
	}
	if req.ZoneName != nil {
		ctReq.AddParam("zoneName", *req.ZoneName)
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
	var resp CtvpcNewPrivateZoneListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcNewPrivateZoneListRequest struct {
	ClientToken *string /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string  /*  资源池ID  */
	ZoneID      *string /*  zoneID  */
	ZoneName    *string /*  zoneName  */
	PageNumber  int32   /*  列表的页码，默认值为1。  */
	PageNo      int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize    int32   /*  分页查询时每页的行数，最大值为50，默认值为10。  */
}

type CtvpcNewPrivateZoneListResponse struct {
	StatusCode  int32                                     `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                   `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                   `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                   `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcNewPrivateZoneListReturnObjResponse `json:"returnObj"`             /*  object  */
	Error       *string                                   `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcNewPrivateZoneListReturnObjResponse struct {
	Zones        []*CtvpcNewPrivateZoneListReturnObjZonesResponse `json:"zones"`        /*  dns 记录  */
	TotalCount   int32                                            `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                            `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                            `json:"totalPage"`    /*  总页数  */
}

type CtvpcNewPrivateZoneListReturnObjZonesResponse struct {
	ZoneID *string `json:"zoneID,omitempty"` /*  名称  */
	Name   *string `json:"name,omitempty"`   /*  zone record名称  */
}
