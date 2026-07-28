package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfReportTaskFailureApi
/* 上报云工作流异步任务失败 */
type CfReportTaskFailureApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfReportTaskFailureApi(client *core.CtyunClient) *CfReportTaskFailureApi {
	return &CfReportTaskFailureApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/openapi/v1/taskevent/reportTaskFailure",
			ContentType:  "application/json",
		},
	}
}

func (a *CfReportTaskFailureApi) Do(ctx context.Context, credential core.Credential, req *CfReportTaskFailureRequest) (*CfReportTaskFailureResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfReportTaskFailureRequest
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
	var resp CfReportTaskFailureResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfReportTaskFailureRequest struct {
	RegionId     string `json:"regionId"`     /*  资源池ID，您可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>接口获取所需资源池ID  */
	TaskToken    string `json:"taskToken"`    /*  任务令牌，通过<a href="https://www.ctyun.cn/document/10006234/11056633#section-6ed3df54ef975272">等待任务令牌</a>执行获取  */
	ErrorType    string `json:"errorType"`    /*  错误内容  */
	ErrorMessage string `json:"errorMessage"` /*  错误消息  */
}

type CfReportTaskFailureResponse struct {
	StatusCode *int32                                `json:"statusCode"` /*  状态码，0表示成功，非0表示不成功  */
	Error      *string                               `json:"error"`      /*  错误码  */
	Message    *string                               `json:"message"`    /*  结果描述  */
	ReturnObj  *CfReportTaskFailureReturnObjResponse `json:"returnObj"`  /*  结果数据  */
}

type CfReportTaskFailureReturnObjResponse struct {
	Ok      *bool   `json:"ok"`      /*  是否成功  */
	Code    *string `json:"code"`    /*  上报结果  */
	Message *string `json:"message"` /*  上报结果额外消息  */
}
