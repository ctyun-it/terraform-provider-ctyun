package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetEdgeBranchEdgeApi
/* 互联查询EDGE(可以创建互联的edge) */
type SdwanGetEdgeBranchEdgeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetEdgeBranchEdgeApi(client *core.CtyunClient) *SdwanGetEdgeBranchEdgeApi {
	return &SdwanGetEdgeBranchEdgeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-branch/list-edge",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetEdgeBranchEdgeApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetEdgeBranchEdgeRequest) (*SdwanGetEdgeBranchEdgeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("sdwanID", req.SdwanID)
	if req.Zone != nil && *req.Zone != "" {
		ctReq.AddParam("zone", *req.Zone)
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.Search != nil && *req.Search != "" {
		ctReq.AddParam("search", *req.Search)
	}
	if req.EdgeID != nil && *req.EdgeID != "" {
		ctReq.AddParam("edgeID", *req.EdgeID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetEdgeBranchEdgeResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetEdgeBranchEdgeRequest struct {
	SdwanID  string  `json:"sdwanID"`          /*  sdwan id  */
	Zone     *string `json:"zone,omitempty"`   /*  区域  */
	PageNo   int32   `json:"pageNo"`           /*  页码  */
	PageSize int32   `json:"pageSize"`         /*  每页记录数目  */
	Search   *string `json:"search,omitempty"` /*  模糊查询  */
	EdgeID   *string `json:"edgeID,omitempty"` /*  edge id, 传这个盒子ID筛选出不包含该盒子以及与该盒子互联的盒子的其他盒子的信息  */
}

type SdwanGetEdgeBranchEdgeResponse struct {
	StatusCode  int32                                    `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                  `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                  `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                  `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetEdgeBranchEdgeReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                  `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetEdgeBranchEdgeReturnObjResponse struct {
	Message      *string                                          `json:"message"`      /*  message  */
	TotalCount   int32                                            `json:"totalCount"`   /*  总数  */
	CurrentCount int32                                            `json:"currentCount"` /*  当前页数量  */
	Code         *string                                          `json:"code"`         /*  状态码  */
	Result       []*SdwanGetEdgeBranchEdgeReturnObjResultResponse `json:"result"`       /*  edge列表  */
}

type SdwanGetEdgeBranchEdgeReturnObjResultResponse struct {
	EdgeBranchID *string `json:"edgeBranchID"` /*  互联信息id  */
	EdgeName     *string `json:"edgeName"`     /*  智能网关名称  */
}
