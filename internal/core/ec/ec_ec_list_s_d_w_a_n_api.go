package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EcEcListSDWANApi
/* 查询sdwan，云间高速console调用，查询当前租户所有的sdwan列表 */
type EcEcListSDWANApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcListSDWANApi(client *core.CtyunClient) *EcEcListSDWANApi {
	return &EcEcListSDWANApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/sdwan/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcListSDWANApi) Do(ctx context.Context, credential core.Credential, req *EcEcListSDWANRequest) (*EcEcListSDWANResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.SdwanID != nil && *req.SdwanID != "" {
		ctReq.AddParam("sdwanID", *req.SdwanID)
	}
	if *req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(*req.PageNo), 10))
	}
	if *req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(*req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcListSDWANResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcListSDWANRequest struct {
	SdwanID  *string `json:"sdwanID,omitempty"`  /*  sdwan ID  */
	PageNo   *int32  `json:"pageNo,omitempty"`   /*  查询的当前页，用于分页查询，从1开始  */
	PageSize *int32  `json:"pageSize,omitempty"` /*  分页查询页数  */
}

type EcEcListSDWANResponse struct {
	StatusCode  *int32                          `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                         `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                         `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                         `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                         `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcListSDWANReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcListSDWANReturnObjResponse struct {
	CurrentCount *int32                                   `json:"currentCount"` /*  当前页记录数  */
	TotalPage    *int32                                   `json:"totalPage"`    /*  总页数  */
	TotalCount   *int32                                   `json:"totalCount"`   /*  查询的总记录数  */
	Results      []*EcEcListSDWANReturnObjResultsResponse `json:"results"`      /*  返回查询结果，Json数组  */
}

type EcEcListSDWANReturnObjResultsResponse struct {
	SdwanID   *string `json:"sdwanID"`   /*  sdwan ID  */
	SdwanName *string `json:"sdwanName"` /*  sdwan名称  */
	EcID      *string `json:"ecID"`      /*  sdwan已经加入的云间高速ID  */
}
