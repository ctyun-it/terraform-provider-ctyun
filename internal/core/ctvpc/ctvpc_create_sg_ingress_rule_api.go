package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateSgIngressRuleApi
/* 创建安全组入向规则。
 */type CtvpcCreateSgIngressRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateSgIngressRuleApi(client *core.CtyunClient) *CtvpcCreateSgIngressRuleApi {
	return &CtvpcCreateSgIngressRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/create-security-group-ingress",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateSgIngressRuleApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateSgIngressRuleRequest) (*CtvpcCreateSgIngressRuleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcCreateSgIngressRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateSgIngressRuleRequest struct {
	RegionID           string                                               `json:"regionID,omitempty"`        /*  区域id  */
	SecurityGroupID    string                                               `json:"securityGroupID,omitempty"` /*  安全组ID。  */
	SecurityGroupRules []*CtvpcCreateSgIngressRuleSecurityGroupRulesRequest `json:"securityGroupRules"`        /*  规则信息  */
	ClientToken        string                                               `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
}

type CtvpcCreateSgIngressRuleSecurityGroupRulesRequest struct {
	Direction             string  `json:"direction,omitempty"`             /*  入方向  */
	RemoteType            int32   `json:"remoteType"`                      /*  remote 类型，0 表示使用 cidr，1 表示使用远端安全组，默认为 0   */
	RemoteSecurityGroupID *string `json:"remoteSecurityGroupID,omitempty"` /*  远端安全组 id  */
	Action                string  `json:"action,omitempty"`                /*  拒绝策略:允许-accept 拒绝-drop  */
	Priority              int32   `json:"priority"`                        /*  优先级:1~100，取值越小优先级越大  */
	Protocol              string  `json:"protocol,omitempty"`              /*  协议: ANY、TCP、UDP、ICMP(v4)  */
	Ethertype             string  `json:"ethertype,omitempty"`             /*  IP类型:IPv4、IPv6  */
	DestCidrIp            string  `json:"destCidrIp,omitempty"`            /*  远端地址:0.0.0.0/0  */
	RawRange              *string `json:"range,omitempty"`                 /*  安全组开放的传输层协议相关的源端端口范围  */
}

type CtvpcCreateSgIngressRuleResponse struct {
	StatusCode  int32                                      `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCreateSgIngressRuleReturnObjResponse `json:"returnObj"`             /*  业务数据  */
}

type CtvpcCreateSgIngressRuleReturnObjResponse struct {
	SgRuleIDs []*string `json:"sgRuleIDs"` /*  安全组规则 id 列表  */
}
