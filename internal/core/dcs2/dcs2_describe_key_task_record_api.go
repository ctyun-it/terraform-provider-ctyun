package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeKeyTaskRecordApi
/* 查询分布式缓存Redis实例大key、热key相关的任务结果
 */type Dcs2DescribeKeyTaskRecordApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeKeyTaskRecordApi(client *core.CtyunClient) *Dcs2DescribeKeyTaskRecordApi {
	return &Dcs2DescribeKeyTaskRecordApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/keyAnalysisMgrServant/describeKeyTaskRecord",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeKeyTaskRecordApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeKeyTaskRecordRequest) (*Dcs2DescribeKeyTaskRecordResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("taskId", req.TaskId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeKeyTaskRecordResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeKeyTaskRecordRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	TaskId     string /*  任务ID  */
}

type Dcs2DescribeKeyTaskRecordResponse struct {
	StatusCode int32                                       `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                      `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeKeyTaskRecordReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                      `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                      `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                      `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeKeyTaskRecordReturnObjResponse struct {
	Status   string                                                `json:"status,omitempty"` /*  任务状态  */
	Time     string                                                `json:"time,omitempty"`   /*  任务开始查询时间  */
	KeyNodes []*Dcs2DescribeKeyTaskRecordReturnObjKeyNodesResponse `json:"keyNodes"`         /*  热key(HotKeyNode)或大key(BigKeyNode)对象的集合  */
}

type Dcs2DescribeKeyTaskRecordReturnObjKeyNodesResponse struct {
	Key        string `json:"key,omitempty"`        /*  要查询的key  */
	QueryCount string `json:"queryCount,omitempty"` /*  请求次数  */
	Counter    int32  `json:"counter,omitempty"`    /*  频次  */
	Db         int32  `json:"db,omitempty"`         /*  热key所在分组db  */
	NodeUrl    string `json:"nodeUrl,omitempty"`    /*  热key所在节点的url  */
}
