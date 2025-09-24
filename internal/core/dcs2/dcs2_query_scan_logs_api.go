package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// Dcs2QueryScanLogsApi
/* 查询分布式缓存Redis实例中清除过期Key的记录。
 */type Dcs2QueryScanLogsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2QueryScanLogsApi(client *core.CtyunClient) *Dcs2QueryScanLogsApi {
	return &Dcs2QueryScanLogsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/redisDataMgr/queryScanLogs",
			ContentType:  "",
		},
	}
}

func (a *Dcs2QueryScanLogsApi) Do(ctx context.Context, credential core.Credential, req *Dcs2QueryScanLogsRequest) (*Dcs2QueryScanLogsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("pageIndex", strconv.FormatInt(int64(req.PageIndex), 10))
	ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2QueryScanLogsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2QueryScanLogsRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	PageIndex  int32  /*  页码  */
	PageSize   int32  /*  每页数量, 范围[1,100]  */
}

type Dcs2QueryScanLogsResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                              `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2QueryScanLogsReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                              `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                              `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                              `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2QueryScanLogsReturnObjResponse struct {
	Total int32                                     `json:"total,omitempty"` /*  对象总数量  */
	Rows  []*Dcs2QueryScanLogsReturnObjRowsResponse `json:"rows"`            /*  对象列表  */
}

type Dcs2QueryScanLogsReturnObjRowsResponse struct {
	TaskId    string `json:"taskId,omitempty"`    /*  任务ID  */
	Status    int32  `json:"status,omitempty"`    /*  扫描记录状态<li>0：待完成<li>1：处理完成<li>2：处理异常  */
	RawType   int32  `json:"type,omitempty"`      /*  任务类型<li>0：手动分析<li>1：自动分析  */
	StartTime string `json:"startTime,omitempty"` /*  扫描开始时间  */
	EndTime   string `json:"endTime,omitempty"`   /*  扫描结束时间  */
	ScanKeys  int64  `json:"scanKeys,omitempty"`  /*  扫描key总数  */
	DelKeys   int64  `json:"delKeys,omitempty"`   /*  清理key数量  */
	Msg       string `json:"msg,omitempty"`       /*  错误信息  */
}
