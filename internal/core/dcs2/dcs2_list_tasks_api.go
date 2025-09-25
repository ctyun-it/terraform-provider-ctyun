package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ListTasksApi
/* 查询任务列表。
 */type Dcs2ListTasksApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ListTasksApi(client *core.CtyunClient) *Dcs2ListTasksApi {
	return &Dcs2ListTasksApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/taskCenter/listTasks",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ListTasksApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ListTasksRequest) (*Dcs2ListTasksResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2ListTasksResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ListTasksRequest struct {
	RegionId  string                         /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	StartTime string                         `json:"startTime,omitempty"` /*  开始时间（yyyyMMdd）  */
	EndTime   string                         `json:"endTime,omitempty"`   /*  结束时间（yyyyMMdd）（最长时间跨度：31天）  */
	PageIndex int32                          `json:"pageIndex,omitempty"` /*  页码（默认：1）  */
	PageSize  int32                          `json:"pageSize,omitempty"`  /*  每页条数（默认：10，范围：1~100）  */
	Condition *Dcs2ListTasksConditionRequest `json:"condition"`           /*  查询条件  */
}

type Dcs2ListTasksConditionRequest struct {
	TaskTypeStr   string `json:"taskTypeStr,omitempty"`   /*  查询所有任务类型:空值，查询指定任务类型：输入指定taskType，多个时英文逗号分割；taskType可选值见taskType表  */
	Status        int32  `json:"status,omitempty"`        /*  查询指定任务状态<li>1：所有状态的任务<li>0：初始态任务<li>1：运行中的任务<li>2：成功的任务<li>3：失败的任务  */
	ProdInstId    string `json:"prodInstId,omitempty"`    /*  查询指定实例id的任务列表  */
	StartTimeDesc int32  `json:"startTimeDesc,omitempty"` /*  返回数据的排序<li>0：按创建时间自然排序<li>1：按创建时间降序排序  */
}

type Dcs2ListTasksResponse struct {
	StatusCode int32                           `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                          `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2ListTasksReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                          `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                          `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                          `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ListTasksReturnObjResponse struct {
	Total int32 `json:"total,omitempty"` /*  总数  */
}
