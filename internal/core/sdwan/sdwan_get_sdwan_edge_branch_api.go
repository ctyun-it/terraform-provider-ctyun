package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetSdwanEdgeBranchApi
/* 查询edge互联信息列表(已经互联的edge) */
type SdwanGetSdwanEdgeBranchApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanEdgeBranchApi(client *core.CtyunClient) *SdwanGetSdwanEdgeBranchApi {
	return &SdwanGetSdwanEdgeBranchApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-branch/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanEdgeBranchApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanEdgeBranchRequest) (*SdwanGetSdwanEdgeBranchResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("sdwanID", req.SdwanID)
	if req.EdgeBranchID != nil && *req.EdgeBranchID != "" {
		ctReq.AddParam("edgeBranchID", *req.EdgeBranchID)
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
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanEdgeBranchResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanEdgeBranchRequest struct {
	SdwanID      string  `json:"sdwanID"`                /*  sdwan id  */
	EdgeBranchID *string `json:"edgeBranchID,omitempty"` /*  盒子互联信息ID  */
	PageNo       int32   `json:"pageNo"`                 /*  页码  */
	PageSize     int32   `json:"pageSize"`               /*  每页记录数目  */
	Search       *string `json:"search,omitempty"`       /*  模糊查询  */
}

type SdwanGetSdwanEdgeBranchResponse struct {
	StatusCode  int32                                     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetSdwanEdgeBranchReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetSdwanEdgeBranchReturnObjResponse struct {
	ReturnMessage *string                                           `json:"returnMessage"` /*  返回信息  */
	TotalCount    int32                                             `json:"totalCount"`    /*  总数  */
	CurrentCount  int32                                             `json:"currentCount"`  /*  当前页数量  */
	Code          *string                                           `json:"code"`          /*  code  */
	Result        []*SdwanGetSdwanEdgeBranchReturnObjResultResponse `json:"result"`        /*  列表  */
}

type SdwanGetSdwanEdgeBranchReturnObjResultResponse struct {
	EdgeBranchID     *string `json:"edgeBranchID"`     /*  互联信息id  */
	FuserLastUpdated *string `json:"fuserLastUpdated"` /*  用户最近更新时间  */
	EdgeName         *string `json:"edgeName"`         /*  edge名称  */
	CustomerID       *string `json:"customerID"`       /*  客户id  */
	Rate             *string `json:"rate"`             /*  速率  */
	DstDcName        *string `json:"dstDcName"`        /*  目的region名称  */
	DstDcID          *string `json:"dstDcID"`          /*  目的region id  */
	DstZone          *string `json:"dstZone"`          /*  目的区域  */
	DstEdgeID        *string `json:"dstEdgeID"`        /*  目的 edge id  */
	SrcDcName        *string `json:"srcDcName"`        /*  源 region名称  */
	SrcDcID          *string `json:"srcDcID"`          /*  源 region id  */
	SrcZone          *string `json:"srcZone"`          /*  源区域  */
	SrcEdgeID        *string `json:"srcEdgeID"`        /*  源 edge id  */
	SdwanID          *string `json:"sdwanID"`          /*  sdwan id  */
}
