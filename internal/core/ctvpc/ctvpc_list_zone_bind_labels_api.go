package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListZoneBindLabelsApi
/* 获取内网 DNS 绑定的标签
 */type CtvpcListZoneBindLabelsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListZoneBindLabelsApi(client *core.CtyunClient) *CtvpcListZoneBindLabelsApi {
	return &CtvpcListZoneBindLabelsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/private-zone/list-labels",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListZoneBindLabelsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListZoneBindLabelsRequest) (*CtvpcListZoneBindLabelsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("zoneID", req.ZoneID)
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
	var resp CtvpcListZoneBindLabelsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListZoneBindLabelsRequest struct {
	RegionID string /*  区域ID  */
	ZoneID   string /*  内网 DNS ID  */
	PageNo   int32  /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize int32  /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcListZoneBindLabelsResponse struct {
	StatusCode   int32                                       `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	TotalCount   int32                                       `json:"totalCount"`            /*  列表条目数  */
	CurrentCount int32                                       `json:"currentCount"`          /*  分页查询时每页的行数。  */
	TotalPage    int32                                       `json:"totalPage"`             /*  总页数  */
	ReturnObj    []*CtvpcListZoneBindLabelsReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error        *string                                     `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListZoneBindLabelsReturnObjResponse struct {
	Results []*CtvpcListZoneBindLabelsReturnObjResultsResponse `json:"results"` /*  绑定的标签列表  */
}

type CtvpcListZoneBindLabelsReturnObjResultsResponse struct {
	LabelID    *string `json:"labelID,omitempty"`    /*  标签 id  */
	LabelKey   *string `json:"labelKey,omitempty"`   /*  标签名  */
	LabelValue *string `json:"labelValue,omitempty"` /*  标签值  */
}
