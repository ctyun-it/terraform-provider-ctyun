package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanListEdgeDatacentersApi
/* 查找智能网关可接入资源池信息 */
type SdwanSdwanListEdgeDatacentersApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanListEdgeDatacentersApi(client *core.CtyunClient) *SdwanSdwanListEdgeDatacentersApi {
	return &SdwanSdwanListEdgeDatacentersApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-datacenters/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanListEdgeDatacentersApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanListEdgeDatacentersRequest) (*SdwanSdwanListEdgeDatacentersResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
	ctReq.AddParam("managerType", req.ManagerType)
	if req.SearchName != nil && *req.SearchName != "" {
		ctReq.AddParam("searchName", *req.SearchName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanSdwanListEdgeDatacentersResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanListEdgeDatacentersRequest struct {
	EdgeID      string  `json:"edgeID"`               /*  edge id  */
	ManagerType string  `json:"managerType"`          /*  本参数表示资源池类型<br/>取值范围:<br/>DPOP:PON类型的智能网关接入的资源池<br/>DAGW:A8C类型的智能网关接入的资源池  */
	SearchName  *string `json:"searchName,omitempty"` /*  edge名称模糊查找  */
}

type SdwanSdwanListEdgeDatacentersResponse struct {
	StatusCode  int32                                           `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                         `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                         `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                         `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanListEdgeDatacentersReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                         `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanSdwanListEdgeDatacentersReturnObjResponse struct {
	Result       []*SdwanSdwanListEdgeDatacentersReturnObjResultResponse `json:"result"`       /*  查询edge 可接入资源池信息  */
	TotalCount   int32                                                   `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                                   `json:"currentCount"` /*  当前页数量  */
}

type SdwanSdwanListEdgeDatacentersReturnObjResultResponse struct {
	DcID      *string `json:"dcID"`      /*  资源池ID  */
	DcName    *string `json:"dcName"`    /*  资源池名称  */
	EdgeCount *string `json:"edgeCount"` /*  资源池接入智能网关计数  */
}
