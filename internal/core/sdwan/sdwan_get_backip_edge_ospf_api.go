package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetBackipEdgeOspfApi
/* 智能网关ospf模糊查询 */
type SdwanGetBackipEdgeOspfApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetBackipEdgeOspfApi(client *core.CtyunClient) *SdwanGetBackipEdgeOspfApi {
	return &SdwanGetBackipEdgeOspfApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/backup-edge-ospf/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetBackipEdgeOspfApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetBackipEdgeOspfRequest) (*SdwanGetBackipEdgeOspfResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.RegionID != nil && *req.RegionID != "" {
		ctReq.AddParam("regionID", *req.RegionID)
	}
	if req.EdgeID != nil && *req.EdgeID != "" {
		ctReq.AddParam("edgeID", *req.EdgeID)
	}
	if req.Search != nil && *req.Search != "" {
		ctReq.AddParam("search", *req.Search)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetBackipEdgeOspfResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetBackipEdgeOspfRequest struct {
	RegionID *string `json:"regionID,omitempty"` /*  资源池ID  */
	EdgeID   *string `json:"edgeID,omitempty"`   /*  智能网关ID  */
	Search   *string `json:"search,omitempty"`   /*  模糊查询  */
}

type SdwanGetBackipEdgeOspfResponse struct {
	StatusCode  int32                                    `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                  `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                  `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                  `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetBackipEdgeOspfReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                  `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetBackipEdgeOspfReturnObjResponse struct {
	Result       []*SdwanGetBackipEdgeOspfReturnObjResultResponse `json:"result"`       /*  模糊查询智能网关ospf  */
	TotalCount   int32                                            `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                            `json:"currentCount"` /*  当前页数量  */
}

type SdwanGetBackipEdgeOspfReturnObjResultResponse struct {
	EdgeName *string `json:"edgeName"` /*  智能网关名称  */
	EdgeID   *string `json:"edgeID"`   /*  智能网关ID  */
}
