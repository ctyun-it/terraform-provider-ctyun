package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// Dcs2FindHistorySlowLogApi
/* 查询分布式缓存Redis实例历史慢日志
 */type Dcs2FindHistorySlowLogApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2FindHistorySlowLogApi(client *core.CtyunClient) *Dcs2FindHistorySlowLogApi {
	return &Dcs2FindHistorySlowLogApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/resourceMonitor/findHistorySlowLog",
			ContentType:  "",
		},
	}
}

func (a *Dcs2FindHistorySlowLogApi) Do(ctx context.Context, credential core.Credential, req *Dcs2FindHistorySlowLogRequest) (*Dcs2FindHistorySlowLogResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("nodeName", req.NodeName)
	ctReq.AddParam("redisSetName", req.RedisSetName)
	ctReq.AddParam("startTime", req.StartTime)
	ctReq.AddParam("endTime", req.EndTime)
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
	var resp Dcs2FindHistorySlowLogResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2FindHistorySlowLogRequest struct {
	RegionId     string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId   string /*  实例ID  */
	NodeName     string /*  节点名称<li>实例类型为Proxy集群、经典集群版时传入proxyName<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7741&isNormal=1&vid=270">查询实例的逻辑拓扑</a> AccessNode表proxyName字段<li>其他实例类型传入Redis节点名称<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7752&isNormal=1&vid=270">获取redis节点名列表</a> node表nodeName字段  */
	RedisSetName string /*  redis集群名  */
	StartTime    string /*  开始时间  */
	EndTime      string /*  结束时间  */
	Page         int32  /*  页码，默认1  */
	Rows         int32  /*  每页行数，默认10，允许值范围1~100  */
}

type Dcs2FindHistorySlowLogResponse struct {
	StatusCode int32                                    `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                   `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2FindHistorySlowLogReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                   `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                   `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                   `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2FindHistorySlowLogReturnObjResponse struct{}
