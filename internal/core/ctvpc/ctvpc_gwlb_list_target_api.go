package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcGwlbListTargetApi
/* 查看target列表
 */type CtvpcGwlbListTargetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGwlbListTargetApi(client *core.CtyunClient) *CtvpcGwlbListTargetApi {
	return &CtvpcGwlbListTargetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/gwlb/list-target",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGwlbListTargetApi) Do(ctx context.Context, credential core.Credential, req *CtvpcGwlbListTargetRequest) (*CtvpcGwlbListTargetResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("targetGroupID", req.TargetGroupID)
	if req.TargetID != nil {
		ctReq.AddParam("targetID", *req.TargetID)
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
	var resp CtvpcGwlbListTargetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGwlbListTargetRequest struct {
	RegionID      string  /*  资源池 ID  */
	TargetGroupID string  /*  后端服务组 ID  */
	TargetID      *string /*  后端服务 ID  */
	PageNumber    int32   /*  列表的页码，默认值为 1。  */
	PageSize      int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcGwlbListTargetResponse struct {
	StatusCode  int32                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcGwlbListTargetReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcGwlbListTargetReturnObjResponse struct {
	Results      []*CtvpcGwlbListTargetReturnObjResultsResponse `json:"results"`      /*  接口业务数据  */
	TotalCount   int32                                          `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                          `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                          `json:"totalPage"`    /*  总页数  */
}

type CtvpcGwlbListTargetReturnObjResultsResponse struct {
	TargetID              *string `json:"targetID,omitempty"`              /*  后端服务ID  */
	TargetGroupID         *string `json:"targetGroupID,omitempty"`         /*  后端服务组ID  */
	InstanceType          *string `json:"instanceType,omitempty"`          /*  实例类型，取值有: VM / BMS/ CBM  */
	InstanceID            *string `json:"instanceID,omitempty"`            /*  实例 ID  */
	InstanceVpc           *string `json:"instanceVpc,omitempty"`           /*  实例所在的 vpc  */
	Weight                int32   `json:"weight"`                          /*  权重  */
	HealthCheckStatus     *string `json:"healthCheckStatus,omitempty"`     /*  ipv4 健康检查状态，取值: unknown / online / offline  */
	HealthCheckStatusIpv6 *string `json:"healthCheckStatusIpv6,omitempty"` /*  ipv6 健康检查状态，取值: unknown / online / offline  */
	CreatedAt             *string `json:"createdAt,omitempty"`             /*  创建时间  */
	UpdatedAt             *string `json:"updatedAt,omitempty"`             /*  更新时间  */
}
