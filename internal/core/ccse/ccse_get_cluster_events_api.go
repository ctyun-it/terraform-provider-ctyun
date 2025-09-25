package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CcseGetClusterEventsApi
/* 查询指定集群事件列表
 */type CcseGetClusterEventsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetClusterEventsApi(client *core.CtyunClient) *CcseGetClusterEventsApi {
	return &CcseGetClusterEventsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/events/{clusterId}",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetClusterEventsApi) Do(ctx context.Context, credential core.Credential, req *CcseGetClusterEventsRequest) (*CcseGetClusterEventsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.EventType != "" {
		ctReq.AddParam("eventType", req.EventType)
	}
	if req.TaskId != "" {
		ctReq.AddParam("taskId", req.TaskId)
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
	var resp CcseGetClusterEventsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetClusterEventsRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	EventType  string /*  事件类型  */
	TaskId     string /*  任务ID  */
	PageNumber int32  /*  每页显示数量  */
	PageSize   int32  /*  分页查询页数  */
}

type CcseGetClusterEventsResponse struct {
	StatusCode int32                                  `json:"statusCode,omitempty"` /*  响应状态码  */
	RequestId  string                                 `json:"requestId,omitempty"`  /*  请求ID  */
	Message    string                                 `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *CcseGetClusterEventsReturnObjResponse `json:"returnObj"`            /*  请求结果  */
	Error      string                                 `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type CcseGetClusterEventsReturnObjResponse struct {
	Records []*CcseGetClusterEventsReturnObjRecordsResponse `json:"records"`           /*  诊断任务记录  */
	Total   int32                                           `json:"total,omitempty"`   /*  总记录数  */
	Pages   int32                                           `json:"pages,omitempty"`   /*  总页数  */
	Current int32                                           `json:"current,omitempty"` /*  当前页数  */
	Size    int32                                           `json:"size,omitempty"`    /*  每页大小  */
}

type CcseGetClusterEventsReturnObjRecordsResponse struct {
	TaskId       string `json:"taskId,omitempty"`       /*  任务ID  */
	ClusterId    string `json:"clusterId,omitempty"`    /*  集群ID  */
	EventId      string `json:"eventId,omitempty"`      /*  事件Id  */
	Source       string `json:"source,omitempty"`       /*  事件源  */
	Subject      string `json:"subject,omitempty"`      /*  事件关联操作对象ID  */
	EventType    string `json:"eventType,omitempty"`    /*  事件类型  */
	EventMessage string `json:"eventMessage,omitempty"` /*  事件内容  */
	CreatedTime  string `json:"createdTime,omitempty"`  /*  创建时间  */
}
