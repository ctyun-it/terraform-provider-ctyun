package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanUpdateSdwanAclRuleApi
/* 修改访问控制规则 */
type SdwanUpdateSdwanAclRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanUpdateSdwanAclRuleApi(client *core.CtyunClient) *SdwanUpdateSdwanAclRuleApi {
	return &SdwanUpdateSdwanAclRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/acl-rule/update",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanUpdateSdwanAclRuleApi) Do(ctx context.Context, credential core.Credential, req *SdwanUpdateSdwanAclRuleRequest) (*SdwanUpdateSdwanAclRuleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanUpdateSdwanAclRuleRequest
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
	var resp SdwanUpdateSdwanAclRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanUpdateSdwanAclRuleRequest struct {
	AclID       string                                       `json:"aclID"`                 /*  ACL ID  */
	UpdateRules []*SdwanUpdateSdwanAclRuleUpdateRulesRequest `json:"updateRules,omitempty"` /*  修改规则  */
}

type SdwanUpdateSdwanAclRuleUpdateRulesRequest struct {
	AclRuleID    *string `json:"aclRuleID,omitempty"`    /*  rule ID  */
	Direction    *string `json:"direction,omitempty"`    /*  本参数表示控制方向<br/><br/>取值范围:<br/>in:入方向<br/>out:出方向  */
	Protocol     *string `json:"protocol,omitempty"`     /*  本参数表示协议类型<br/><br/>取值范围:<br/>udp:UDP<br/>icmp:ICMP</br>all:ALL</br>tcp:TCP  */
	IpVersion    *string `json:"ipVersion,omitempty"`    /*  本参数表示IP协议版本<br/><br/>取值范围:<br/>IPv4:IPv4<br/>IPv6:IPv6  */
	DstCidr      *string `json:"dstCidr,omitempty"`      /*  目的网段  */
	DstPortRange *string `json:"dstPortRange,omitempty"` /*  目的端口范围  */
	Priority     int32   `json:"priority,omitempty"`     /*  priority  */
	Action       *string `json:"action,omitempty"`       /*  本参数表示策略类型<br/><br/>取值范围:<br/>allow:允许<br/>deny:拒绝  */
	SrcCidr      *string `json:"srcCidr,omitempty"`      /*  源网段  */
	SrcPortRange *string `json:"srcPortRange,omitempty"` /*  源端口范围  */
}

type SdwanUpdateSdwanAclRuleResponse struct {
	StatusCode  int32   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	Error       *string `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
