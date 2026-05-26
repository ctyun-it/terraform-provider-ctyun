package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfInvokeFunctionApi
/* 调用函数 */
type CfInvokeFunctionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfInvokeFunctionApi(client *core.CtyunClient) *CfInvokeFunctionApi {
	return &CfInvokeFunctionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/openapi/v1/functions/{functionName}/invocations",
			ContentType:  "application/json",
		},
	}
}

func (a *CfInvokeFunctionApi) Do(ctx context.Context, credential core.Credential, req *CfInvokeFunctionRequest) (*CfInvokeFunctionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfInvokeFunctionRequest
		RegionId     interface{} `json:"regionId,omitempty"`
		FunctionName interface{} `json:"functionName,omitempty"`
	}{
		req, nil, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfInvokeFunctionResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfInvokeFunctionRequest struct {
	FunctionName      string  `json:"functionName"`      /*  函数名称  */
	RegionId          string  `json:"regionId"`          /*  资源池id  */
	XFcInvocationType string  `json:"xFcInvocationType"` /*  调用类型，异步：Async  */
	XFcAsyncDelay     int32   `json:"xFcAsyncDelay"`     /*  延迟调用的时间，单位秒，范围0-3600  */
	Body              *string `json:"body,omitempty"`    /*  函数调用参数  */
	XFcAsyncTaskID    string  `json:"xFcAsyncTaskID"`    /*  异步任务 ID。长度<=128个字符，仅支持英文字符、数字、连字符、下划线  */
	Method            *string `json:"method,omitempty"`  /*  方法，GET/POST等HTTP方法  */
	Qualifier         string  `json:"qualifier"`         /*  版本或别名  */
}

type CfInvokeFunctionResponse struct {
	StatusCode *int32  `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Code       *string `json:"code"`       /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string `json:"message"`    /*  错误描述信息  */
}
