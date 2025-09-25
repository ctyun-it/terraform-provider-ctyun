package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseGetClusterTaskApi
/* 查询任务详情
 */type CcseGetClusterTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetClusterTaskApi(client *core.CtyunClient) *CcseGetClusterTaskApi {
	return &CcseGetClusterTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/tasks/{taskId}",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetClusterTaskApi) Do(ctx context.Context, credential core.Credential, req *CcseGetClusterTaskRequest) (*CcseGetClusterTaskResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("taskId", req.TaskId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetClusterTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetClusterTaskRequest struct {
	TaskId   string /*  任务ID  */
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
}

type CcseGetClusterTaskResponse struct {
	StatusCode int32                                `json:"statusCode,omitempty"` /*  响应状态码  */
	RequestId  string                               `json:"requestId,omitempty"`  /*  请求ID  */
	Message    string                               `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *CcseGetClusterTaskReturnObjResponse `json:"returnObj"`            /*  请求结果  */
	Error      string                               `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type CcseGetClusterTaskReturnObjResponse struct {
	TaskId         string `json:"taskId,omitempty"`         /*  任务ID  */
	ClusterId      string `json:"clusterId,omitempty"`      /*  集群ID  */
	RegionId       string `json:"regionId,omitempty"`       /*  资源池Id  */
	TaskType       string `json:"taskType,omitempty"`       /*  任务类型  */
	TaskStatus     string `json:"taskStatus,omitempty"`     /*  任务状态  */
	ParallelNumber int32  `json:"parallelNumber,omitempty"` /*  并行数  */
	TaskContent    string `json:"taskContent,omitempty"`    /*  任务内容  */
	TaskResult     string `json:"taskResult,omitempty"`     /*  任务执行结果  */
	RetryTime      int32  `json:"retryTime,omitempty"`      /*  重试次数  */
	CreatedTime    string `json:"createdTime,omitempty"`    /*  创建时间  */
	ModifyTime     string `json:"modifyTime,omitempty"`     /*  修改时间  */
}
