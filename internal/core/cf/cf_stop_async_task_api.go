package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfStopAsyncTaskApi
/* 停止异步任务 */
type CfStopAsyncTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfStopAsyncTaskApi(client *core.CtyunClient) *CfStopAsyncTaskApi {
	return &CfStopAsyncTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPut,
			UrlPath:      "/openapi/v1/functions/{functionName}/async-tasks/{taskId}/stop",
			ContentType:  "application/json",
		},
	}
}

func (a *CfStopAsyncTaskApi) Do(ctx context.Context, credential core.Credential, req *CfStopAsyncTaskRequest) (*CfStopAsyncTaskResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder = builder.ReplaceUrl("taskId", req.TaskId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfStopAsyncTaskRequest
		RegionId     interface{} `json:"regionId,omitempty"`
		FunctionName interface{} `json:"functionName,omitempty"`
		TaskId       interface{} `json:"taskId,omitempty"`
	}{
		req, nil, nil, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfStopAsyncTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfStopAsyncTaskRequest struct {
	FunctionName string `json:"functionName"` /*  函数名称，函数必须存在  */
	TaskId       string `json:"taskId"`       /*  异步任务 ID，任务ID必须存在  */
	RegionId     string `json:"regionId"`     /*  资源池id，标识不同的地区，如：华东1、西南1  */
	Qualifier    string `json:"qualifier"`    /*  版本或别名，版本包括1,2,...这样的普通版本，和特殊版本LATEST。  */
}

type CfStopAsyncTaskResponse struct {
	StatusCode *int32  `json:"statusCode"` /*  状态码。0表示成功，非0表示不成功  */
	Error      *string `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败。  */
	Message    *string `json:"message"`    /*  错误描述信息  */
}
