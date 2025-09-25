package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeHotKeyTasksApi
/* 查询分布式缓存Redis实例的热key分析任务列表，需先创建热key分析任务。
 */type Dcs2DescribeHotKeyTasksApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeHotKeyTasksApi(client *core.CtyunClient) *Dcs2DescribeHotKeyTasksApi {
	return &Dcs2DescribeHotKeyTasksApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/keyAnalysisMgrServant/describeHotKeyTasks",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeHotKeyTasksApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeHotKeyTasksRequest) (*Dcs2DescribeHotKeyTasksResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeHotKeyTasksResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeHotKeyTasksRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
}

type Dcs2DescribeHotKeyTasksResponse struct {
	StatusCode int32                                     `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                    `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeHotKeyTasksReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                    `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                    `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeHotKeyTasksReturnObjResponse struct {
	Total int32                                           `json:"total,omitempty"` /*  对象总数量  */
	Rows  []*Dcs2DescribeHotKeyTasksReturnObjRowsResponse `json:"rows"`            /*  对象列表  */
}

type Dcs2DescribeHotKeyTasksReturnObjRowsResponse struct {
	TaskId       string `json:"taskId,omitempty"`       /*  任务ID  */
	Status       string `json:"status,omitempty"`       /*  状态<li>success：成功<li>processing：进行中<li>fail：失败  */
	Time         string `json:"time,omitempty"`         /*  开始时间,时间格式为“年月日时分秒”（如YYYYMMDDHHMMSS）  */
	ScanType     string `json:"scanType,omitempty"`     /*  执行方式  */
	FinishedTime string `json:"finishedTime,omitempty"` /*  结束时间  */
	TaskType     int32  `json:"taskType,omitempty"`     /*  任务类型<li>0：大key<li>1：热key  */
}
