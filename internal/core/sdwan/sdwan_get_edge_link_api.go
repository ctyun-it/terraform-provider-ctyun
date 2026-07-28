package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetEdgeLinkApi
/* 查找智能网关链路信息 */
type SdwanGetEdgeLinkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetEdgeLinkApi(client *core.CtyunClient) *SdwanGetEdgeLinkApi {
	return &SdwanGetEdgeLinkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-link/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetEdgeLinkApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetEdgeLinkRequest) (*SdwanGetEdgeLinkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetEdgeLinkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetEdgeLinkRequest struct {
	EdgeID string `json:"edgeID"` /*  智能网关ID  */
}

type SdwanGetEdgeLinkResponse struct {
	StatusCode  int32                              `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                            `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                            `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                            `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetEdgeLinkReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                            `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetEdgeLinkReturnObjResponse struct {
	Result       []*SdwanGetEdgeLinkReturnObjResultResponse `json:"result"`       /*  查询edge 链路信息  */
	TotalCount   int32                                      `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                      `json:"currentCount"` /*  当前页数量  */
}

type SdwanGetEdgeLinkReturnObjResultResponse struct {
	LinkID       *string `json:"linkID"`       /*  链路Id  */
	LinkPort     *string `json:"linkPort"`     /*  本参数表示端口名称<br/>取值范围:<br/>WAN1:WAN1<br/>WAN2:WAN2<br/>LAN:LAN<br/>LTE:LTE  */
	TransportNet *string `json:"transportNet"` /*  本参数表示传输网络类型<br/><br/>取值范围:<br/>internet:internet<br/>mpls:mpls<br/>lte:lte  */
	Level        *string `json:"level"`        /*  本参数表示主备链路配置<br/>取值范围：<br/>master:主<br/>slave:备  */
	Status       *string `json:"status"`       /*  本参数表示端口状态<br/><br/>取值范围：<br/>up:开启<br/>down:关闭  */
	Active       *string `json:"active"`       /*  是否活跃  */
	LTEType      *string `json:"LTEType"`      /*  本参数表示LTE类型<br/>取值范围：<br/>4G:4G<br/>5G:5G  */
}
