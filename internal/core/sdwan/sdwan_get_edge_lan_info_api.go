package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetEdgeLanInfoApi
/* 查找智能网关lan侧信息 */
type SdwanGetEdgeLanInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetEdgeLanInfoApi(client *core.CtyunClient) *SdwanGetEdgeLanInfoApi {
	return &SdwanGetEdgeLanInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-lan-info/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetEdgeLanInfoApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetEdgeLanInfoRequest) (*SdwanGetEdgeLanInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetEdgeLanInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetEdgeLanInfoRequest struct {
	EdgeID   string `json:"edgeID"`   /*  智能网关ID  */
	PageNo   int32  `json:"pageNo"`   /*  页数  */
	PageSize int32  `json:"pageSize"` /*  页大小  */
}

type SdwanGetEdgeLanInfoResponse struct {
	StatusCode  int32                                 `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                               `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                               `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                               `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetEdgeLanInfoReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                               `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetEdgeLanInfoReturnObjResponse struct {
	Result       []*SdwanGetEdgeLanInfoReturnObjResultResponse `json:"result"`       /*  查询edge lan侧信息  */
	TotalCount   int32                                         `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                         `json:"currentCount"` /*  当前页的数量  */
}

type SdwanGetEdgeLanInfoReturnObjResultResponse struct {
	LanIPDict *SdwanGetEdgeLanInfoReturnObjResultLanIPDictResponse `json:"lanIPDict"` /*  lan侧IP信息  */
}

type SdwanGetEdgeLanInfoReturnObjResultLanIPDictResponse struct {
	LanBusinessIP      *string `json:"lanBusinessIP"`      /*  lan侧业务IP  */
	SlaveLanBusinessIP *string `json:"slaveLanBusinessIP"` /*  备edge lan侧业务IP  */
	LanVrrpIP          *string `json:"lanVrrpIP"`          /*  lan侧vrrp IP  */
}
