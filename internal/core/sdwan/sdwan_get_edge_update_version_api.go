package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetEdgeUpdateVersionApi
/* 查找智能网关可升级版本信息 */
type SdwanGetEdgeUpdateVersionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetEdgeUpdateVersionApi(client *core.CtyunClient) *SdwanGetEdgeUpdateVersionApi {
	return &SdwanGetEdgeUpdateVersionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-update-version/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetEdgeUpdateVersionApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetEdgeUpdateVersionRequest) (*SdwanGetEdgeUpdateVersionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetEdgeUpdateVersionResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetEdgeUpdateVersionRequest struct {
	EdgeID string `json:"edgeID"` /*  智能网关ID  */
}

type SdwanGetEdgeUpdateVersionResponse struct {
	StatusCode  int32                                       `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                     `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                     `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                     `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetEdgeUpdateVersionReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                     `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetEdgeUpdateVersionReturnObjResponse struct {
	Result       []*SdwanGetEdgeUpdateVersionReturnObjResultResponse `json:"result"`       /*  查询edge 可升级版本信息  */
	TotalCount   int32                                               `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                               `json:"currentCount"` /*  当前页数量  */
}

type SdwanGetEdgeUpdateVersionReturnObjResultResponse struct {
	Version     *string `json:"version"`     /*  edge版本  */
	VersionDesc *string `json:"versionDesc"` /*  edge版本  */
}
