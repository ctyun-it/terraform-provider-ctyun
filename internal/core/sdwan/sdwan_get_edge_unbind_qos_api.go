package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetEdgeUnbindQosApi
/* 查询qos未绑定智能网关信息列表 */
type SdwanGetEdgeUnbindQosApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetEdgeUnbindQosApi(client *core.CtyunClient) *SdwanGetEdgeUnbindQosApi {
	return &SdwanGetEdgeUnbindQosApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-unbind-qos/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetEdgeUnbindQosApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetEdgeUnbindQosRequest) (*SdwanGetEdgeUnbindQosResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.Search != nil && *req.Search != "" {
		ctReq.AddParam("search", *req.Search)
	}
	if req.SdwanName != nil && *req.SdwanName != "" {
		ctReq.AddParam("sdwanName", *req.SdwanName)
	}
	ctReq.AddParam("bandwidthType", req.BandwidthType)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetEdgeUnbindQosResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetEdgeUnbindQosRequest struct {
	PageNo        int32   `json:"pageNo"`              /*  页数  */
	PageSize      int32   `json:"pageSize"`            /*  页大小  */
	Search        *string `json:"search,omitempty"`    /*  模糊查询  */
	SdwanName     *string `json:"sdwanName,omitempty"` /*  sdwan名称  */
	BandwidthType string  `json:"bandwidthType"`       /*  本参数表示带宽类型<br/><br/>取值范围:<br/>internet:互联网带宽<br/>sdwan:SD-WAN带宽  */
}

type SdwanGetEdgeUnbindQosResponse struct {
	StatusCode  int32                                   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                 `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                 `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                 `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetEdgeUnbindQosReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                 `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetEdgeUnbindQosReturnObjResponse struct {
	Result       []*SdwanGetEdgeUnbindQosReturnObjResultResponse `json:"result"`       /*  查询qos未绑定的盒子信息  */
	TotalCount   int32                                           `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                           `json:"currentCount"` /*  当前页数量  */
}

type SdwanGetEdgeUnbindQosReturnObjResultResponse struct {
	EdgeID    *string `json:"edgeID"`    /*  智能网关ID  */
	EdgeName  *string `json:"edgeName"`  /*  智能网关名称  */
	SdwanName *string `json:"sdwanName"` /*  sdwan名称  */
	Status    *string `json:"status"`    /*  本参数表示设备状态<br/><br/>取值范围：<br/>online:在线<br/>offline:下线  */
	Bandwidth int32   `json:"bandwidth"` /*  带宽  */
}
