package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListFlowSessionsApi
/* 查看流量会话列表
 */type CtvpcListFlowSessionsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListFlowSessionsApi(client *core.CtyunClient) *CtvpcListFlowSessionsApi {
	return &CtvpcListFlowSessionsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/flowsession/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListFlowSessionsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListFlowSessionsRequest) (*CtvpcListFlowSessionsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.MirrorFilterID != nil {
		ctReq.AddParam("mirrorFilterID", *req.MirrorFilterID)
	}
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
	var resp CtvpcListFlowSessionsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListFlowSessionsRequest struct {
	RegionID       string  /*  区域ID  */
	MirrorFilterID *string /*  名称  */
	QueryContent   *string /*  按名字进行模糊过滤  */
	PageNumber     int32   /*  列表的页码，默认值为 1。  */
	PageNo         int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize       int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcListFlowSessionsResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcListFlowSessionsReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                 `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListFlowSessionsReturnObjResponse struct {
	FlowSessions []*CtvpcListFlowSessionsReturnObjFlowSessionsResponse `json:"flowSessions"` /*  流量镜像过滤规则列表  */
	TotalCount   int32                                                 `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                                 `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                                 `json:"totalPage"`    /*  总页数  */
}

type CtvpcListFlowSessionsReturnObjFlowSessionsResponse struct {
	MirrorFilterID *string `json:"mirrorFilterID,omitempty"` /*  过滤规则 ID  */
	SrcPort        *string `json:"srcPort,omitempty"`        /*  源弹性网卡 ID  */
	DstPort        *string `json:"dstPort,omitempty"`        /*  源弹性网卡 ID  */
	CreatedTime    *string `json:"createdTime,omitempty"`    /*  创建时间  */
	FlowSessionID  *string `json:"flowSessionID,omitempty"`  /*  会话 ID  */
	Name           *string `json:"name,omitempty"`           /*  会话名称  */
	Description    *string `json:"description,omitempty"`    /*  会话描述  */
	Vni            int32   `json:"vni"`                      /*  VXLAN 网络标识符  */
	DstPortType    *string `json:"dstPortType,omitempty"`    /*  目标网卡类型: VM  */
	Status         *string `json:"status,omitempty"`         /*  会话状态：on / off  */
}
