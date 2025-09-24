package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtelbIplistenerListApi
/* 查看ip_listener列表
 */type CtelbIplistenerListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbIplistenerListApi(client *core.CtyunClient) *CtelbIplistenerListApi {
	return &CtelbIplistenerListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/iplistener/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbIplistenerListApi) Do(ctx context.Context, credential core.Credential, req *CtelbIplistenerListRequest) (*CtelbIplistenerListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.IpListenerID != "" {
		ctReq.AddParam("ipListenerID", req.IpListenerID)
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
	var resp CtelbIplistenerListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbIplistenerListRequest struct {
	RegionID     string /*  资源池 ID  */
	IpListenerID string /*  监听器 ID  */
	PageNumber   int32  /*  列表的页码，默认值为 1。  */
	PageSize     int32  /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtelbIplistenerListResponse struct {
	StatusCode  int32                                 `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbIplistenerListReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbIplistenerListReturnObjResponse struct {
	Results      []*CtelbIplistenerListReturnObjResultsResponse `json:"results"`                /*  接口业务数据  */
	TotalCount   int32                                          `json:"totalCount,omitempty"`   /*  列表条目数  */
	CurrentCount int32                                          `json:"currentCount,omitempty"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                          `json:"totalPage,omitempty"`    /*  总页数  */
}

type CtelbIplistenerListReturnObjResultsResponse struct {
	GwElbID      string                                             `json:"gwElbID,omitempty"`      /*  网关负载均衡 ID  */
	Name         string                                             `json:"name,omitempty"`         /*  名字  */
	Description  string                                             `json:"description,omitempty"`  /*  描述  */
	IpListenerID string                                             `json:"ipListenerID,omitempty"` /*  监听器 id  */
	Action       *CtelbIplistenerListReturnObjResultsActionResponse `json:"action"`                 /*  转发配置  */
}

type CtelbIplistenerListReturnObjResultsActionResponse struct {
	RawType       string                                                          `json:"type,omitempty"` /*  默认规则动作类型: forward / redirect  */
	ForwardConfig *CtelbIplistenerListReturnObjResultsActionForwardConfigResponse `json:"forwardConfig"`  /*  转发配置  */
}

type CtelbIplistenerListReturnObjResultsActionForwardConfigResponse struct{}
