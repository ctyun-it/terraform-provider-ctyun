package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// Dcs2DescribeParameterModificationHistoryApi
/* 查看分布式缓存Redis实例配置参数修改历史记录。
 */type Dcs2DescribeParameterModificationHistoryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeParameterModificationHistoryApi(client *core.CtyunClient) *Dcs2DescribeParameterModificationHistoryApi {
	return &Dcs2DescribeParameterModificationHistoryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instanceParam/describeParameterModificationHistory",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeParameterModificationHistoryApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeParameterModificationHistoryRequest) (*Dcs2DescribeParameterModificationHistoryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	if req.StartTime != "" {
		ctReq.AddParam("startTime", req.StartTime)
	}
	if req.EndTime != "" {
		ctReq.AddParam("endTime", req.EndTime)
	}
	if req.HistoryId != "" {
		ctReq.AddParam("historyId", req.HistoryId)
	}
	if req.Page != 0 {
		ctReq.AddParam("page", strconv.FormatInt(int64(req.Page), 10))
	}
	if req.Rows != 0 {
		ctReq.AddParam("rows", strconv.FormatInt(int64(req.Rows), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeParameterModificationHistoryResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeParameterModificationHistoryRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	StartTime  string /*  开始时间  */
	EndTime    string /*  结束时间  */
	HistoryId  string /*  记录ID<br>说明：若historyId为空，则需要传入startTime、endTime, 按时间查询所有记录  */
	Page       int32  /*  页码  */
	Rows       int32  /*  行数  */
}

type Dcs2DescribeParameterModificationHistoryResponse struct {
	StatusCode int32                                                      `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                                     `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeParameterModificationHistoryReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                                     `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                                     `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                                     `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeParameterModificationHistoryReturnObjResponse struct {
	Total int32                                                            `json:"total,omitempty"` /*  数量  */
	Rows  []*Dcs2DescribeParameterModificationHistoryReturnObjRowsResponse `json:"rows"`            /*  参数对象  */
}

type Dcs2DescribeParameterModificationHistoryReturnObjRowsResponse struct {
	ParamName     string `json:"paramName,omitempty"`     /*  实例ID  */
	OriginalValue string `json:"originalValue,omitempty"` /*  旧值  */
	CurrentValue  string `json:"currentValue,omitempty"`  /*  当前值  */
}
