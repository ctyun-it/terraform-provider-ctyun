package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListMirrorFlowFilterRuleApi
/* 查看过滤规则列表
 */type CtvpcListMirrorFlowFilterRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListMirrorFlowFilterRuleApi(client *core.CtyunClient) *CtvpcListMirrorFlowFilterRuleApi {
	return &CtvpcListMirrorFlowFilterRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/mirrorflow/list-filter-rule",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListMirrorFlowFilterRuleApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListMirrorFlowFilterRuleRequest) (*CtvpcListMirrorFlowFilterRuleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("mirrorFilterID", req.MirrorFilterID)
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
	ctReq.AddParam("direction", req.Direction)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcListMirrorFlowFilterRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListMirrorFlowFilterRuleRequest struct {
	RegionID       string  /*  区域ID  */
	MirrorFilterID string  /*  名称  */
	QueryContent   *string /*  按名字进行模糊过滤  */
	PageNumber     int32   /*  列表的页码，默认值为 1  */
	PageNo         int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize       int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
	Direction      string  /*  规则的出入方向: in / out  */
}

type CtvpcListMirrorFlowFilterRuleResponse struct {
	StatusCode  int32                                           `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                         `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                         `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                         `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcListMirrorFlowFilterRuleReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                         `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListMirrorFlowFilterRuleReturnObjResponse struct {
	MirrorFilterRules []*CtvpcListMirrorFlowFilterRuleReturnObjMirrorFilterRulesResponse `json:"mirrorFilterRules"` /*  流量镜像过滤规则列表  */
	TotalCount        int32                                                              `json:"totalCount"`        /*  列表条目数  */
	CurrentCount      int32                                                              `json:"currentCount"`      /*  分页查询时每页的行数。  */
	TotalPage         int32                                                              `json:"totalPage"`         /*  总页数  */
}

type CtvpcListMirrorFlowFilterRuleReturnObjMirrorFilterRulesResponse struct {
	DestCidr           *string `json:"destCidr,omitempty"`           /*  目标 cidr  */
	SrcCidr            *string `json:"srcCidr,omitempty"`            /*  源 cidr  */
	DestPort           int32   `json:"destPort"`                     /*  目的端口，格式为 1-2， 协议为 all，传值 -  */
	SrcPort            int32   `json:"srcPort"`                      /*  源端口，格式为 1-2，协议为 all，传值 -  */
	Protocol           *string `json:"protocol,omitempty"`           /*  协议：all / tcp / udp / icmp  */
	EnableCollection   *bool   `json:"enableCollection"`             /*  是否开启采集，true 表示采集，false 表示不采集  */
	MirrorFilterRuleID *string `json:"mirrorFilterRuleID,omitempty"` /*  流量镜像过滤规则 id  */
}
