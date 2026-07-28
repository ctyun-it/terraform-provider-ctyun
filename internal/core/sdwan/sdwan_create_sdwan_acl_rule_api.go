package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanCreateSdwanAclRuleApi
/* 创建访问控制规则 */
type SdwanCreateSdwanAclRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanCreateSdwanAclRuleApi(client *core.CtyunClient) *SdwanCreateSdwanAclRuleApi {
	return &SdwanCreateSdwanAclRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/acl-rule/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanCreateSdwanAclRuleApi) Do(ctx context.Context, credential core.Credential, req *SdwanCreateSdwanAclRuleRequest) (*SdwanCreateSdwanAclRuleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanCreateSdwanAclRuleRequest
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
	var resp SdwanCreateSdwanAclRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanCreateSdwanAclRuleRequest struct {
	AclID    string                                    `json:"aclID"`              /*  ACL ID  */
	AddRules []*SdwanCreateSdwanAclRuleAddRulesRequest `json:"addRules,omitempty"` /*  新增规则  */
}

type SdwanCreateSdwanAclRuleAddRulesRequest struct {
	Direction    *string `json:"direction,omitempty"`    /*  本参数表示控制方向<br/><br/>取值范围:<br/>in:入方向<br/>out:出方向  */
	Protocol     *string `json:"protocol,omitempty"`     /*  本参数表示协议类型<br/><br/>取值范围:<br/>udp:UDP<br/>icmp:ICMP</br>all:ALL</br>tcp:TCP  */
	IpVersion    *string `json:"ipVersion,omitempty"`    /*  本参数表示IP协议版本<br/><br/>取值范围:<br/>IPv4:IPv4<br/>IPv6:IPv6  */
	DstCidr      *string `json:"dstCidr,omitempty"`      /*  目的网段  */
	DstPortRange *string `json:"dstPortRange,omitempty"` /*  目的端口范围（例如1-200， -1/-1为默认值，表示1-65535）  */
	Priority     int32   `json:"priority,omitempty"`     /*  priority  */
	Action       *string `json:"action,omitempty"`       /*  本参数表示策略类型<br/><br/>取值范围:<br/>allow:允许<br/>deny:拒绝  */
	SrcCidr      *string `json:"srcCidr,omitempty"`      /*  源网段  */
	SrcPortRange *string `json:"srcPortRange,omitempty"` /*  源端口范围（例如1-200， -1/-1为默认值，表示1-65535  */
}

type SdwanCreateSdwanAclRuleResponse struct {
	StatusCode  int32                                       `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                     `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                     `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                     `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *[]SdwanCreateSdwanAclRuleReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                     `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanCreateSdwanAclRuleReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作id  */
}
