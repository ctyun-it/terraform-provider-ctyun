package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcsePauseClusterTaskApi
/* 暂停任务
 */type CcsePauseClusterTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcsePauseClusterTaskApi(client *core.CtyunClient) *CcsePauseClusterTaskApi {
	return &CcsePauseClusterTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/tasks/{taskId}/pause",
			ContentType:  "application/json",
		},
	}
}

func (a *CcsePauseClusterTaskApi) Do(ctx context.Context, credential core.Credential, req *CcsePauseClusterTaskRequest) (*CcsePauseClusterTaskResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("taskId", req.TaskId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcsePauseClusterTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcsePauseClusterTaskRequest struct {
	TaskId   string /*  任务ID  */
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
}

type CcsePauseClusterTaskResponse struct {
	StatusCode int32                                  `json:"statusCode,omitempty"` /*  响应状态码  */
	RequestId  string                                 `json:"requestId,omitempty"`  /*  请求ID  */
	Message    string                                 `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *CcsePauseClusterTaskReturnObjResponse `json:"returnObj"`            /*  请求结果  */
	Error      string                                 `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type CcsePauseClusterTaskReturnObjResponse struct{}
