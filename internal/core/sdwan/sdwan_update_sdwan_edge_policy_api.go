package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanUpdateSdwanEdgePolicyApi
/* 智能选路策略修改 */
type SdwanUpdateSdwanEdgePolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanUpdateSdwanEdgePolicyApi(client *core.CtyunClient) *SdwanUpdateSdwanEdgePolicyApi {
	return &SdwanUpdateSdwanEdgePolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-policy/update",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanUpdateSdwanEdgePolicyApi) Do(ctx context.Context, credential core.Credential, req *SdwanUpdateSdwanEdgePolicyRequest) (*SdwanUpdateSdwanEdgePolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanUpdateSdwanEdgePolicyRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanUpdateSdwanEdgePolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanUpdateSdwanEdgePolicyRequest struct {
	EdgeID          string                                       `json:"edgeID"`                    /*  智能网关ID(该参数不可修改)  */
	PolicyID        string                                       `json:"policyID"`                  /*  策略id(该参数不可修改)  */
	PolicyName      string                                       `json:"policyName"`                /*  策略名称  */
	Priority        int32                                        `json:"priority"`                  /*  1~500  */
	SelectType      string                                       `json:"selectType"`                /*  本参数表示方式<br/><br/>取值范围:<br/>priority:优先级<br/>quality:质量<br/>weight:权重  */
	Protocol        string                                       `json:"protocol"`                  /*  本参数表示协议<br/><br/>取值范围:<br/>tcp:tcp<br/>udp:udp<br/>icmp:icmp  */
	SrcCidr         string                                       `json:"srcCidr"`                   /*  源网段  */
	SrcPort         string                                       `json:"srcPort"`                   /*  源端口范围  */
	DstCidr         string                                       `json:"dstCidr"`                   /*  目的网段  */
	DstPort         string                                       `json:"dstPort"`                   /*  目的端口范围  */
	LinkList        []*SdwanUpdateSdwanEdgePolicyLinkListRequest `json:"linkList,omitempty"`        /*  链路列表  */
	LinkQualityID   *string                                      `json:"linkQualityID,omitempty"`   /*  链路探测业务id  */
	LinkQualityName *string                                      `json:"linkQualityName,omitempty"` /*  链路探测业务名称  */
	AppList         []*SdwanUpdateSdwanEdgePolicyAppListRequest  `json:"appList,omitempty"`         /*  应用列表  */
	Quality         *SdwanUpdateSdwanEdgePolicyQualityRequest    `json:"quality,omitempty"`         /*  链路质量参数  */
}

type SdwanUpdateSdwanEdgePolicyLinkListRequest struct {
	LinkID    int32   `json:"linkID"`             /*  link_id  */
	LinkValue int32   `json:"linkValue"`          /*  当selectType为priority和quality时，取值为1，2，3，当selectType为weight时，取值为1~10  */
	NextHop   *string `json:"nextHop,omitempty"`  /*  next_hop  */
	DetectIP  *string `json:"detectIP,omitempty"` /*  detect_ip  */
}

type SdwanUpdateSdwanEdgePolicyAppListRequest struct {
	AppID int32 `json:"appID"` /*  应用ID  */
}

type SdwanUpdateSdwanEdgePolicyQualityRequest struct {
	Delay  int32 `json:"delay"`  /*  时延  */
	Jitter int32 `json:"jitter"` /*  抖动  */
	Lost   int32 `json:"lost"`   /*  丢包率  */
}

type SdwanUpdateSdwanEdgePolicyResponse struct {
	StatusCode  int32   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	OperationID *string `json:"operationID"` /*  操作id  */
	Error       *string `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
