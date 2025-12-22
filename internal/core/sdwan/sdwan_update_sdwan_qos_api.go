package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanUpdateSdwanQosApi
/* 修改qos */
type SdwanUpdateSdwanQosApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanUpdateSdwanQosApi(client *core.CtyunClient) *SdwanUpdateSdwanQosApi {
	return &SdwanUpdateSdwanQosApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/qos/update",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanUpdateSdwanQosApi) Do(ctx context.Context, credential core.Credential, req *SdwanUpdateSdwanQosRequest) (*SdwanUpdateSdwanQosResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanUpdateSdwanQosRequest
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
	var resp SdwanUpdateSdwanQosResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanUpdateSdwanQosRequest struct {
	QosID          string                                      `json:"qosID"`                    /*  qos策略ID  */
	QosName        string                                      `json:"qosName"`                  /*  qos策略名称  */
	Description    string                                      `json:"description"`              /*  描述  */
	Bandwidth      string                                      `json:"bandwidth"`                /*  带宽峰值  */
	BandwidthType  string                                      `json:"bandwidthType"`            /*  本参数表示带宽类型<br/><br/>取值范围:<br/>internet:互联网带宽<br/>sdwan:SD-WAN带宽  */
	UpdateRuleList []*SdwanUpdateSdwanQosUpdateRuleListRequest `json:"updateRuleList,omitempty"` /*  修改规则列表  */
}

type SdwanUpdateSdwanQosUpdateRuleListRequest struct {
	RuleID     *string                                               `json:"ruleID,omitempty"`     /*  规则id  */
	Priority   string                                                `json:"priority"`             /*  优先级  */
	Rate       string                                                `json:"rate"`                 /*  速率  */
	RuleType   string                                                `json:"ruleType"`             /*  本参数表示规则类型<br/><br/>取值范围:<br/>ace:五元组<br/>app:应用组  */
	AddAceRule []*SdwanUpdateSdwanQosUpdateRuleListAddAceRuleRequest `json:"addAceRule,omitempty"` /*  增加QOS规则列表  */
	AddAppRule []*SdwanUpdateSdwanQosUpdateRuleListAddAppRuleRequest `json:"addAppRule,omitempty"` /*  增加QOS规则列表  */
	DelRule    []*string                                             `json:"delRule,omitempty"`    /*  删除QOS规则列表  ，值类型为string  */
}

type SdwanUpdateSdwanQosUpdateRuleListAddAceRuleRequest struct {
	AceRuleName  string  `json:"aceRuleName"`            /*  ace规则策略名称  */
	Protocol     string  `json:"protocol"`               /*  本参数表示协议<br/><br/>取值范围:<br/>tcp:tcp<br/>udp:udp<br/>icmp:icmp  */
	IpVersion    string  `json:"ipVersion"`              /*  本参数表示ip版本<br/><br/>取值范围:<br/>ipv4:ipv4<br/>ipv6:ipv6  */
	DstCidr      *string `json:"dstCidr,omitempty"`      /*  目的子网  */
	SrcCidr      *string `json:"srcCidr,omitempty"`      /*  源子网  */
	DstPortRange *string `json:"dstPortRange,omitempty"` /*  目的端口范围  */
	SrcPortRange *string `json:"srcPortRange,omitempty"` /*  源端口范围  */
}

type SdwanUpdateSdwanQosUpdateRuleListAddAppRuleRequest struct {
	AppGroupName string                                                    `json:"appGroupName"`   /*  应用组名称  */
	GroupType    string                                                    `json:"groupType"`      /*  本参数表示应用组类型<br/><br/>取值范围:<br/>system:系统<br/>custom:自定义  */
	Apps         []*SdwanUpdateSdwanQosUpdateRuleListAddAppRuleAppsRequest `json:"apps,omitempty"` /*  app列表  */
}

type SdwanUpdateSdwanQosUpdateRuleListAddAppRuleAppsRequest struct {
	AppName string `json:"appName"` /*  app名称  */
	AppID   string `json:"appID"`   /*  app id  */
}

type SdwanUpdateSdwanQosResponse struct {
	StatusCode  int32   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	Error       *string `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
