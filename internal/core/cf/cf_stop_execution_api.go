package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfStopExecutionApi
/* 停止执行云工作流（仅适用于标准工作流） */
type CfStopExecutionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfStopExecutionApi(client *core.CtyunClient) *CfStopExecutionApi {
	return &CfStopExecutionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/openapi/v1/execution/stop",
			ContentType:  "application/json",
		},
	}
}

func (a *CfStopExecutionApi) Do(ctx context.Context, credential core.Credential, req *CfStopExecutionRequest) (*CfStopExecutionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfStopExecutionRequest
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfStopExecutionResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfStopExecutionRequest struct {
	RegionId    string  `json:"regionId"`            /*  资源池ID，您可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>接口获取所需资源池ID  */
	ExecutionId string  `json:"executionId"`         /*  执行id  */
	ErrorType   *string `json:"errorType,omitempty"` /*  停止原因的错误类型  */
	Cause       *string `json:"cause,omitempty"`     /*  停止原因  */
}

type CfStopExecutionResponse struct {
	StatusCode *int32                            `json:"statusCode"` /*  状态码，0表示成功，非0表示不成功  */
	Error      *string                           `json:"error"`      /*  错误码  */
	Message    *string                           `json:"message"`    /*  结果描述  */
	ReturnObj  *CfStopExecutionReturnObjResponse `json:"returnObj"`  /*  结果数据  */
}

type CfStopExecutionReturnObjResponse struct {
	StartedTime *int32 `json:"startedTime"` /*  执行开始事件  */
	StoppedTime *int32 `json:"stoppedTime"` /*  执行结束时间  */
}
