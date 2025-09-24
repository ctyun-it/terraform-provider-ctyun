package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListMirrorFlowFilterApi
/* 查看过滤条件列表
 */type CtvpcListMirrorFlowFilterApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListMirrorFlowFilterApi(client *core.CtyunClient) *CtvpcListMirrorFlowFilterApi {
	return &CtvpcListMirrorFlowFilterApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/mirrorflow/list-filter",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListMirrorFlowFilterApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListMirrorFlowFilterRequest) (*CtvpcListMirrorFlowFilterResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.QueryContent != nil {
		ctReq.AddParam("queryContent", *req.QueryContent)
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
	var resp CtvpcListMirrorFlowFilterResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListMirrorFlowFilterRequest struct {
	RegionID     string  /*  区域ID  */
	QueryContent *string /*  按名字进行模糊过滤  */
	PageNumber   int32   /*  列表的页码，默认值为 1  */
	PageNo       int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize     int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcListMirrorFlowFilterResponse struct {
	StatusCode  int32                                       `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcListMirrorFlowFilterReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                     `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListMirrorFlowFilterReturnObjResponse struct {
	MirrorFilters []*CtvpcListMirrorFlowFilterReturnObjMirrorFiltersResponse `json:"mirrorFilters"` /*  流量镜像过滤条件列表  */
	TotalCount    int32                                                      `json:"totalCount"`    /*  列表条目数  */
	CurrentCount  int32                                                      `json:"currentCount"`  /*  分页查询时每页的行数。  */
	TotalPage     int32                                                      `json:"totalPage"`     /*  总页数  */
}

type CtvpcListMirrorFlowFilterReturnObjMirrorFiltersResponse struct {
	Name            *string `json:"name,omitempty"`           /*  流量镜像名称  */
	Description     *string `json:"description,omitempty"`    /*  流量镜像描述  */
	MirrorFilterID  *string `json:"mirrorFilterID,omitempty"` /*  流量镜像ID  */
	CreatedTime     *string `json:"createdTime,omitempty"`    /*  创建时间  */
	InRuleCount     int32   `json:"inRuleCount"`              /*  出方向规则数  */
	OutRuleCount    int32   `json:"outRuleCount"`             /*  入方向规则数  */
	MirrorFlowCount int32   `json:"mirrorFlowCount"`          /*  关联流量镜像会话数  */
}
