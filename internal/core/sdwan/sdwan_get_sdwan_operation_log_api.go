package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetSdwanOperationLogApi
/* 查询任务状态 */
type SdwanGetSdwanOperationLogApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanOperationLogApi(client *core.CtyunClient) *SdwanGetSdwanOperationLogApi {
	return &SdwanGetSdwanOperationLogApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/operation/list-logs",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanOperationLogApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanOperationLogRequest) (*SdwanGetSdwanOperationLogResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("operationID", req.OperationID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanOperationLogResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanOperationLogRequest struct {
	OperationID string `json:"operationID"` /*  任务ID  */
}

type SdwanGetSdwanOperationLogResponse struct {
	StatusCode   int32                                      `json:"statusCode"`   /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode    *string                                    `json:"errorCode"`    /*  业务细分码，为product.module.code三段式码  */
	Message      *string                                    `json:"message"`      /*  失败时的错误描述，一般为英文描述  */
	Description  *string                                    `json:"description"`  /*  失败时的错误描述，一般为中文描述  */
	TotalCount   int32                                      `json:"totalCount"`   /*  总数  */
	CurrentCount int32                                      `json:"currentCount"` /*  当前页总数  */
	Code         *string                                    `json:"code"`         /*  状态码  */
	Result       []*SdwanGetSdwanOperationLogResultResponse `json:"result"`       /*  列表  */
	Error        *string                                    `json:"error"`        /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetSdwanOperationLogResultResponse struct {
	OperationID *string `json:"operationID"` /*  任务ID  */
	JobType     *string `json:"jobType"`     /*  本参数表示任务类型<br/><br/>取值范围:<br/>sdwan:sdwan  */
	Status      *string `json:"status"`      /*  本参数表示任务状态<br/><br/>取值范围:<br/>done:完成<br/>undone:未完成<br/>failed:失败<br/>unknown:未知  */
}
