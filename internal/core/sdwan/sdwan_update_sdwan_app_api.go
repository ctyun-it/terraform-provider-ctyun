package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanUpdateSdwanAppApi
/* 应用修改 */
type SdwanUpdateSdwanAppApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanUpdateSdwanAppApi(client *core.CtyunClient) *SdwanUpdateSdwanAppApi {
	return &SdwanUpdateSdwanAppApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/app/update",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanUpdateSdwanAppApi) Do(ctx context.Context, credential core.Credential, req *SdwanUpdateSdwanAppRequest) (*SdwanUpdateSdwanAppResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanUpdateSdwanAppRequest
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
	var resp SdwanUpdateSdwanAppResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanUpdateSdwanAppRequest struct {
	AppID          string                                   `json:"appID"`                 /*  应用ID  */
	AppName        string                                   `json:"appName"`               /*  应用组名称  */
	AddRuleList    []*SdwanUpdateSdwanAppAddRuleListRequest `json:"addRuleList,omitempty"` /*  增加规则列表(只支持五元组)  */
	DeleteRuleList []int32                                  `json:"deleteRuleList"`        /*  删除规则列表(只支持五元组，当还剩一条规则不能删)  */
	ModifyRule     *SdwanUpdateSdwanAppModifyRuleRequest    `json:"modifyRule,omitempty"`  /*  修改规则(只剩一条规则时能用)  */
}

type SdwanUpdateSdwanAppAddRuleListRequest struct {
	SrcCidr  string `json:"srcCidr"`  /*  源网段  */
	SrcPort  string `json:"srcPort"`  /*  源端口范围  */
	DstCidr  string `json:"dstCidr"`  /*  目的网段  */
	DstPort  string `json:"dstPort"`  /*  目的端口范围  */
	Protocol string `json:"protocol"` /*  本参数表示协议类型
	取值范围：
	tcp：tcp协议
	udp：udp协议
	icmp：icmp协议  */
}

type SdwanUpdateSdwanAppModifyRuleRequest struct {
	RuleID   string `json:"ruleID"`   /*  规则ID  */
	SrcCidr  string `json:"srcCidr"`  /*  源网段  */
	SrcPort  string `json:"srcPort"`  /*  源端口范围  */
	DstCidr  string `json:"dstCidr"`  /*  目的网段  */
	DstPort  string `json:"dstPort"`  /*  目的端口范围  */
	Protocol string `json:"protocol"` /*  本参数表示协议类型
	取值范围：
	tcp：tcp协议
	udp：udp协议
	icmp：icmp协议  */
	Hostname string `json:"hostname"` /*  访问地址  */
}

type SdwanUpdateSdwanAppResponse struct {
	StatusCode  *string                                 `json:"statusCode"`  /*  返回状态码('800为成功，900为失败) ，默认值:800  */
	ErrorCode   *string                                 `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                 `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                 `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*SdwanUpdateSdwanAppReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                 `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanUpdateSdwanAppReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作id  */
}
