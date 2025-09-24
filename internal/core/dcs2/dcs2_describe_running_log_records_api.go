package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// Dcs2DescribeRunningLogRecordsApi
/* 分页查询分布式缓存Redis实例运行日志
 */type Dcs2DescribeRunningLogRecordsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeRunningLogRecordsApi(client *core.CtyunClient) *Dcs2DescribeRunningLogRecordsApi {
	return &Dcs2DescribeRunningLogRecordsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/logMgr/describeRunningLogRecords",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeRunningLogRecordsApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeRunningLogRecordsRequest) (*Dcs2DescribeRunningLogRecordsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("nodeName", req.NodeName)
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PageIndex != 0 {
		ctReq.AddParam("pageIndex", strconv.FormatInt(int64(req.PageIndex), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeRunningLogRecordsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeRunningLogRecordsRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	NodeName   string /*  节点名称<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7752&isNormal=1&vid=270">获取redis节点名列表</a> node表nodeName字段  */
	PageSize   int32  /*  每页大小，默认值为10  */
	PageIndex  int32  /*  页码，默认值为1  */
}

type Dcs2DescribeRunningLogRecordsResponse struct {
	StatusCode int32                                           `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                          `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeRunningLogRecordsReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                          `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                          `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                          `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeRunningLogRecordsReturnObjResponse struct {
	Rows        []*Dcs2DescribeRunningLogRecordsReturnObjRowsResponse `json:"rows"`                  /*  rows  */
	Total       int32                                                 `json:"total,omitempty"`       /*  数量  */
	CurrentPage int32                                                 `json:"currentPage,omitempty"` /*  当前页  */
}

type Dcs2DescribeRunningLogRecordsReturnObjRowsResponse struct {
	GenTime string `json:"genTime,omitempty"` /*  产生时间  */
	LogInfo string `json:"logInfo,omitempty"` /*  日志信息  */
}
