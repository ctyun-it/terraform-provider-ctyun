package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// Dcs2DescribeTopSlowLogRecoredsApi
/* 查询前N条慢日志
 */type Dcs2DescribeTopSlowLogRecoredsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeTopSlowLogRecoredsApi(client *core.CtyunClient) *Dcs2DescribeTopSlowLogRecoredsApi {
	return &Dcs2DescribeTopSlowLogRecoredsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/logMgr/describeTopSlowLogRecords",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeTopSlowLogRecoredsApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeTopSlowLogRecoredsRequest) (*Dcs2DescribeTopSlowLogRecoredsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("nodeName", req.NodeName)
	if req.Size != 0 {
		ctReq.AddParam("size", strconv.FormatInt(int64(req.Size), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeTopSlowLogRecoredsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeTopSlowLogRecoredsRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	NodeName   string /*  节点名称<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7752&isNormal=1&vid=270">获取redis节点名列表</a> node表nodeName字段  */
	Size       int32  /*  查询前N条  */
}

type Dcs2DescribeTopSlowLogRecoredsResponse struct {
	StatusCode int32                                            `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                           `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeTopSlowLogRecoredsReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                           `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                           `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                           `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeTopSlowLogRecoredsReturnObjResponse struct {
	Total int32                                                  `json:"total,omitempty"` /*  数量  */
	Rows  []*Dcs2DescribeTopSlowLogRecoredsReturnObjRowsResponse `json:"rows"`            /*  慢日志列表  */
}

type Dcs2DescribeTopSlowLogRecoredsReturnObjRowsResponse struct {
	Id        string `json:"id,omitempty"`        /*  日志ID  */
	BeginTime string `json:"beginTime,omitempty"` /*  开始时间,Unix时间戳格式  */
	Cost      string `json:"cost,omitempty"`      /*  消耗时间ms  */
	Command   string `json:"command,omitempty"`   /*  命令、参数  */
}
