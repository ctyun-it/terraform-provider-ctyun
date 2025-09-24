package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosCreateZMSAssessmentTaskApi
/* 创建对象存储评估任务
 */type ZosCreateZMSAssessmentTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosCreateZMSAssessmentTaskApi(client *core.CtyunClient) *ZosCreateZMSAssessmentTaskApi {
	return &ZosCreateZMSAssessmentTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/zms/create-evaluation",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosCreateZMSAssessmentTaskApi) Do(ctx context.Context, credential core.Credential, req *ZosCreateZMSAssessmentTaskRequest) (*ZosCreateZMSAssessmentTaskResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosCreateZMSAssessmentTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosCreateZMSAssessmentTaskRequest struct {
	RegionID       string `json:"regionID,omitempty"`       /*  资源池 ID  */
	EvaluationName string `json:"evaluationName,omitempty"` /*  任务名称，必须为大小写字母、数字、横线或下划线，长度在4-32个字符之间，且名称不能重复  */
}

type ZosCreateZMSAssessmentTaskResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中/失败)  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为product.module.code三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
