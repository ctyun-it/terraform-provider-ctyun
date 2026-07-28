package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcQueryRemainQuotaApi
/* 查询当前用户可创建的云间高速数量 */
type EcEcQueryRemainQuotaApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcQueryRemainQuotaApi(client *core.CtyunClient) *EcEcQueryRemainQuotaApi {
	return &EcEcQueryRemainQuotaApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/query-remain-quota",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcQueryRemainQuotaApi) Do(ctx context.Context, credential core.Credential, req *EcEcQueryRemainQuotaRequest) (*EcEcQueryRemainQuotaResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcQueryRemainQuotaResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcQueryRemainQuotaRequest struct{}

type EcEcQueryRemainQuotaResponse struct {
	StatusCode  *int32                                 `json:"statusCode"`  /*  返回状态码,<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                                `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcQueryRemainQuotaReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcQueryRemainQuotaReturnObjResponse struct {
	CurrentCount *int32                                          `json:"currentCount"` /*  当前页记录数  */
	TotalPage    *int32                                          `json:"totalPage"`    /*  总页数  */
	TotalCount   *int32                                          `json:"totalCount"`   /*  查询的总记录数  */
	Results      []*EcEcQueryRemainQuotaReturnObjResultsResponse `json:"results"`      /*  返回查询结果，Json数组  */
}

type EcEcQueryRemainQuotaReturnObjResultsResponse struct {
	RemainEcQuota *int32 `json:"remainEcQuota"` /*  当前用户可创建的云间高速数量  */
}
