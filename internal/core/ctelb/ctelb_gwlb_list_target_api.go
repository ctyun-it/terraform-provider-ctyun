package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtelbGwlbListTargetApi
/* 查看target列表
 */type CtelbGwlbListTargetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbGwlbListTargetApi(client *core.CtyunClient) *CtelbGwlbListTargetApi {
	return &CtelbGwlbListTargetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/gwlb/list-target",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbGwlbListTargetApi) Do(ctx context.Context, credential core.Credential, req *CtelbGwlbListTargetRequest) (*CtelbGwlbListTargetResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("targetGroupID", req.TargetGroupID)
	if req.TargetID != "" {
		ctReq.AddParam("targetID", req.TargetID)
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
	var resp CtelbGwlbListTargetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbGwlbListTargetRequest struct {
	RegionID      string /*  资源池 ID  */
	TargetGroupID string /*  后端服务组 ID  */
	TargetID      string /*  后端服务 ID  */
	PageNumber    int32  /*  列表的页码，默认值为 1。  */
	PageSize      int32  /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtelbGwlbListTargetResponse struct {
	StatusCode  int32                                 `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbGwlbListTargetReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbGwlbListTargetReturnObjResponse struct {
	Results      []*CtelbGwlbListTargetReturnObjResultsResponse `json:"results"`                /*  接口业务数据  */
	TotalCount   int32                                          `json:"totalCount,omitempty"`   /*  列表条目数  */
	CurrentCount int32                                          `json:"currentCount,omitempty"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                          `json:"totalPage,omitempty"`    /*  总页数  */
}

type CtelbGwlbListTargetReturnObjResultsResponse struct {
	TargetID              string `json:"targetID,omitempty"`              /*  后端服务ID  */
	TargetGroupID         string `json:"targetGroupID,omitempty"`         /*  后端服务组ID  */
	InstanceType          string `json:"instanceType,omitempty"`          /*  实例类型，取值有: VM / BMS/ CBM  */
	InstanceID            string `json:"instanceID,omitempty"`            /*  实例 ID  */
	InstanceVpc           string `json:"instanceVpc,omitempty"`           /*  实例所在的 vpc  */
	Weight                int32  `json:"weight,omitempty"`                /*  权重  */
	HealthCheckStatus     string `json:"healthCheckStatus,omitempty"`     /*  ipv4 健康检查状态，取值: unknown / online / offline  */
	HealthCheckStatusIpv6 string `json:"healthCheckStatusIpv6,omitempty"` /*  ipv6 健康检查状态，取值: unknown / online / offline  */
	CreatedAt             string `json:"createdAt,omitempty"`             /*  创建时间  */
	UpdatedAt             string `json:"updatedAt,omitempty"`             /*  更新时间  */
}
