package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPauseZMSAssessmentTaskApi
/* 对象存储评估任务暂停
 */type ZosPauseZMSAssessmentTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPauseZMSAssessmentTaskApi(client *core.CtyunClient) *ZosPauseZMSAssessmentTaskApi {
	return &ZosPauseZMSAssessmentTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/zms/stop-evaluation",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPauseZMSAssessmentTaskApi) Do(ctx context.Context, credential core.Credential, req *ZosPauseZMSAssessmentTaskRequest) (*ZosPauseZMSAssessmentTaskResponse, error) {
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
	var resp ZosPauseZMSAssessmentTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPauseZMSAssessmentTaskRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池 ID  */
	EvaluationID string `json:"evaluationID,omitempty"` /*  评估任务ID  */
}

type ZosPauseZMSAssessmentTaskResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中/失败)  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为product.module.code三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
