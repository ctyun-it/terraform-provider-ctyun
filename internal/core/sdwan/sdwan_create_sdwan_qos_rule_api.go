package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanCreateSdwanQosRuleApi
/* 增加qos规则五元组/应用组 */
type SdwanCreateSdwanQosRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanCreateSdwanQosRuleApi(client *core.CtyunClient) *SdwanCreateSdwanQosRuleApi {
	return &SdwanCreateSdwanQosRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/qos-rule/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanCreateSdwanQosRuleApi) Do(ctx context.Context, credential core.Credential, req *SdwanCreateSdwanQosRuleRequest) (*SdwanCreateSdwanQosRuleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanCreateSdwanQosRuleRequest
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
	var resp SdwanCreateSdwanQosRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanCreateSdwanQosRuleRequest struct {
	QosID    string                                   `json:"qosID"`             /*  qos策略ID  */
	RuleID   string                                   `json:"ruleID"`            /*  规则id  */
	RuleType string                                   `json:"ruleType"`          /*  本参数表示规则类型<br/><br/>取值范围:<br/>ace:五元组<br/>app:应用组  */
	AceList  []*SdwanCreateSdwanQosRuleAceListRequest `json:"aceList,omitempty"` /*  ace规则列表  */
	AppList  []*SdwanCreateSdwanQosRuleAppListRequest `json:"appList,omitempty"` /*  app应用组列表  */
}

type SdwanCreateSdwanQosRuleAceListRequest struct {
	AceRuleName  string  `json:"aceRuleName"`            /*  ace规则策略名称  */
	Protocol     string  `json:"protocol"`               /*  本参数表示协议<br/><br/>取值范围:<br/>tcp:tcp<br/>udp:udp<br/>icmp:icmp  */
	IpVersion    string  `json:"ipVersion"`              /*  本参数表示ip版本<br/><br/>取值范围:<br/>ipv4:ipv4<br/>ipv6:ipv6  */
	DstCidr      *string `json:"dstCidr,omitempty"`      /*  目的子网  */
	SrcCidr      *string `json:"srcCidr,omitempty"`      /*  源子网  */
	DstPortRange *string `json:"dstPortRange,omitempty"` /*  目的端口范围（例如1-200， -1/-1为默认值，表示1-65535)  */
	SrcPortRange *string `json:"srcPortRange,omitempty"` /*  源端口范围（例如1-200， -1/-1为默认值，表示1-65535）  */
}

type SdwanCreateSdwanQosRuleAppListRequest struct {
	AppGroupName string                                       `json:"appGroupName"`   /*  应用组名称  */
	GroupType    string                                       `json:"groupType"`      /*  本参数表示应用组类型<br/><br/>取值范围:<br/>system:系统<br/>custom:自定义  */
	Apps         []*SdwanCreateSdwanQosRuleAppListAppsRequest `json:"apps,omitempty"` /*  app列表  */
}

type SdwanCreateSdwanQosRuleAppListAppsRequest struct {
	AppName string `json:"appName"` /*  app名称  */
	AppID   string `json:"appID"`   /*  app id  */
}

type SdwanCreateSdwanQosRuleResponse struct {
	StatusCode  int32     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*string `json:"returnObj"`   /*  返回参数  */
	Error       *string   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
