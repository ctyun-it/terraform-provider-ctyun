package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcGwlbListTargetGroupApi
/* 查看target_group列表
 */type CtvpcGwlbListTargetGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGwlbListTargetGroupApi(client *core.CtyunClient) *CtvpcGwlbListTargetGroupApi {
	return &CtvpcGwlbListTargetGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/gwlb/list-target-group",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGwlbListTargetGroupApi) Do(ctx context.Context, credential core.Credential, req *CtvpcGwlbListTargetGroupRequest) (*CtvpcGwlbListTargetGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.TargetGroupID != nil {
		ctReq.AddParam("targetGroupID", *req.TargetGroupID)
	}
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcGwlbListTargetGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGwlbListTargetGroupRequest struct {
	RegionID      string  /*  资源池 ID  */
	TargetGroupID *string /*  后端服务组 ID  */
	PageNumber    int32   /*  列表的页码，默认值为 1。  */
	PageSize      int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcGwlbListTargetGroupResponse struct {
	StatusCode  int32                                      `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcGwlbListTargetGroupReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcGwlbListTargetGroupReturnObjResponse struct {
	Results      []*CtvpcGwlbListTargetGroupReturnObjResultsResponse `json:"results"`      /*  接口业务数据  */
	TotalCount   int32                                               `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                               `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                               `json:"totalPage"`    /*  总页数  */
}

type CtvpcGwlbListTargetGroupReturnObjResultsResponse struct {
	TargetGroupID     *string `json:"targetGroupID,omitempty"` /*  后端服务组ID  */
	Name              *string `json:"name,omitempty"`          /*  名称  */
	Description       *string `json:"description,omitempty"`   /*  描述  */
	VpcID             *string `json:"vpcID,omitempty"`         /*  vpc id  */
	HealthCheckID     *string `json:"healthCheckID,omitempty"` /*  健康检查 ID  */
	SessionStickyMode int32   `json:"sessionStickyMode"`       /*  0:五元组, 4:二元组, 5:三元组, -1 未开启  */
	BypassType        int32   `json:"bypassType"`              /*  旁路类型 1 表示关闭，2 表示自动  */
	FailoverType      int32   `json:"failoverType"`            /*  故障转移类型 1 表示关闭，2 表示再平衡  */
}
