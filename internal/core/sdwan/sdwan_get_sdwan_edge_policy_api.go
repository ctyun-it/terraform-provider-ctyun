package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetSdwanEdgePolicyApi
/* 智能选路策略列表查询 */
type SdwanGetSdwanEdgePolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanEdgePolicyApi(client *core.CtyunClient) *SdwanGetSdwanEdgePolicyApi {
	return &SdwanGetSdwanEdgePolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-policy/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanEdgePolicyApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanEdgePolicyRequest) (*SdwanGetSdwanEdgePolicyResponse, error) {
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
	ctReq.AddParam("edgeID", req.EdgeID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanEdgePolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanEdgePolicyRequest struct {
	PageNo   int32   `json:"pageNo"`           /*  页码  */
	PageSize int32   `json:"pageSize"`         /*  每页计算数目  */
	Search   *string `json:"search,omitempty"` /*  模糊查询  */
	EdgeID   string  `json:"edgeID"`           /*  智能网关ID  */
}

type SdwanGetSdwanEdgePolicyResponse struct {
	StatusCode  int32                                     `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetSdwanEdgePolicyReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetSdwanEdgePolicyReturnObjResponse struct {
	TotalCount   int32                                             `json:"totalCount"`   /*  总数  */
	CurrentCount int32                                             `json:"currentCount"` /*  当前页的数量  */
	Code         *string                                           `json:"code"`         /*  状态码  */
	Result       []*SdwanGetSdwanEdgePolicyReturnObjResultResponse `json:"result"`       /*  列表  */
}

type SdwanGetSdwanEdgePolicyReturnObjResultResponse struct {
	PolicyID        *string                                                     `json:"policyID"`        /*  策略id  */
	PolicyName      *string                                                     `json:"policyName"`      /*  策略名称  */
	SelectType      *string                                                     `json:"selectType"`      /*  本参数表示方式<br/><br/>取值范围:<br/>priority:优先级<br/>quality:质量<br/>weight:权重  */
	Protocol        *string                                                     `json:"protocol"`        /*  本参数表示协议<br/><br/>取值范围:<br/>tcp:tcp<br/>udp:udp<br/>icmp:icmp  */
	Priority        int32                                                       `json:"priority"`        /*  优先级,范围:1~500  */
	SrcCidr         *string                                                     `json:"srcCidr"`         /*  源网段  */
	SrcPort         *string                                                     `json:"srcPort"`         /*  源端口范围  */
	DstCidr         *string                                                     `json:"dstCidr"`         /*  目的网段  */
	DstPort         *string                                                     `json:"dstPort"`         /*  目的端口范围  */
	LinkList        []*SdwanGetSdwanEdgePolicyReturnObjResultLinkListResponse   `json:"linkList"`        /*  链路列表  */
	LinkConfig      []*SdwanGetSdwanEdgePolicyReturnObjResultLinkConfigResponse `json:"linkConfig"`      /*  链路配置列表  */
	LinkQualityID   *string                                                     `json:"linkQualityID"`   /*  链路探测业务id  */
	LinkQualityName *string                                                     `json:"linkQualityName"` /*  链路探测业务名称  */
	AppList         []*SdwanGetSdwanEdgePolicyReturnObjResultAppListResponse    `json:"appList"`         /*  应用列表  */
	Quality         *SdwanGetSdwanEdgePolicyReturnObjResultQualityResponse      `json:"quality"`         /*  链路质量参数  */
}

type SdwanGetSdwanEdgePolicyReturnObjResultLinkListResponse struct {
	LinkID    int32   `json:"linkID"`    /*  link_id  */
	LinkValue int32   `json:"linkValue"` /*  当selectType为priority和quality时，取值为1，2，3，当selectType为weight时，取值为1~10  */
	NextHop   *string `json:"nextHop"`   /*  next_hop  */
	DetectIP  *string `json:"detectIP"`  /*  detect_ip  */
}

type SdwanGetSdwanEdgePolicyReturnObjResultLinkConfigResponse struct {
	LinkID       int32   `json:"linkID"`       /*  链路ID  */
	LinkPort     *string `json:"linkPort"`     /*  本参数表示端口名称<br/>取值范围:<br/>WAN1:WAN1<br/>WAN2:WAN2<br/>LAN:LAN<br/>LTE:LTE  */
	TransportNet *string `json:"transportNet"` /*  本参数表示传输网络类型<br/><br/>取值范围:<br/>internet:internet<br/>mpls:mpls<br/>lte:lte  */
	Level        *string `json:"level"`        /*  本参数表示主备链路配置<br/>取值范围：<br/>master:主<br/>slave:备  */
}

type SdwanGetSdwanEdgePolicyReturnObjResultAppListResponse struct {
	AppID int32 `json:"appID"` /*  应用ID  */
}

type SdwanGetSdwanEdgePolicyReturnObjResultQualityResponse struct {
	Delay  int32 `json:"delay"`  /*  时延  */
	Jitter int32 `json:"jitter"` /*  抖动  */
	Lost   int32 `json:"lost"`   /*  丢包率  */
}
