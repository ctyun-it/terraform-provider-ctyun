package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanCreateSdwanEdgePolicyApi
/* 智能选路策略创建 */
type SdwanCreateSdwanEdgePolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanCreateSdwanEdgePolicyApi(client *core.CtyunClient) *SdwanCreateSdwanEdgePolicyApi {
	return &SdwanCreateSdwanEdgePolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-policy/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanCreateSdwanEdgePolicyApi) Do(ctx context.Context, credential core.Credential, req *SdwanCreateSdwanEdgePolicyRequest) (*SdwanCreateSdwanEdgePolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanCreateSdwanEdgePolicyRequest
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
	var resp SdwanCreateSdwanEdgePolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanCreateSdwanEdgePolicyRequest struct {
	EdgeID          string                                       `json:"edgeID"`                    /*  智能网关ID  */
	PolicyName      string                                       `json:"policyName"`                /*  策略名称  */
	SelectType      string                                       `json:"selectType"`                /*  本参数表示方式<br/><br/>取值范围:<br/>priority:优先级<br/>quality:质量<br/>weight:权重  */
	Protocol        string                                       `json:"protocol"`                  /*  本参数表示协议<br/><br/>取值范围:<br/>tcp:tcp<br/>udp:udp<br/>icmp:icmp  */
	Priority        int32                                        `json:"priority"`                  /*  1~500  */
	SrcCidr         string                                       `json:"srcCidr"`                   /*  源网段  */
	SrcPort         string                                       `json:"srcPort"`                   /*  源端口范围（例如1-200， -1/-1为默认值，表示1-65535）  */
	DstCidr         string                                       `json:"dstCidr"`                   /*  目的网段  */
	DstPort         string                                       `json:"dstPort"`                   /*  目的端口范围（例如1-200， -1/-1为默认值，表示1-65535）  */
	LinkList        []*SdwanCreateSdwanEdgePolicyLinkListRequest `json:"linkList,omitempty"`        /*  链路列表  */
	LinkQualityID   *string                                      `json:"linkQualityID,omitempty"`   /*  链路探测业务id,selectType选quality必填  */
	LinkQualityName *string                                      `json:"linkQualityName,omitempty"` /*  链路探测业务名称,selectType选quality必填  */
	AppList         []*SdwanCreateSdwanEdgePolicyAppListRequest  `json:"appList,omitempty"`         /*  应用列表  */
	Quality         *SdwanCreateSdwanEdgePolicyQualityRequest    `json:"quality,omitempty"`         /*  链路质量参数  */
}

type SdwanCreateSdwanEdgePolicyLinkListRequest struct {
	LinkID    int32   `json:"linkID"`             /*  link_id  */
	LinkValue int32   `json:"linkValue"`          /*  当selectType为priority和quality时，取值为1，2，3，当selectType为weight时，取值为1~10  */
	NextHop   *string `json:"nextHop,omitempty"`  /*  next_hop  */
	DetectIP  *string `json:"detectIP,omitempty"` /*  detect_ip  */
}

type SdwanCreateSdwanEdgePolicyAppListRequest struct {
	AppID int32 `json:"appID"` /*  应用ID  */
}

type SdwanCreateSdwanEdgePolicyQualityRequest struct {
	Delay  int32 `json:"delay"`  /*  时延  */
	Jitter int32 `json:"jitter"` /*  抖动  */
	Lost   int32 `json:"lost"`   /*  丢包率  */
}

type SdwanCreateSdwanEdgePolicyResponse struct {
	StatusCode  int32                                          `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                        `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                        `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                        `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*SdwanCreateSdwanEdgePolicyReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                        `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanCreateSdwanEdgePolicyReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作id  */
}
