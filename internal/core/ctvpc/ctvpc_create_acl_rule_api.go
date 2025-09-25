package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateAclRuleApi
/* 创建 Acl 规则
 */type CtvpcCreateAclRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateAclRuleApi(client *core.CtyunClient) *CtvpcCreateAclRuleApi {
	return &CtvpcCreateAclRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/acl-rule/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateAclRuleApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateAclRuleRequest) (*CtvpcCreateAclRuleResponse, error) {
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
	var resp CtvpcCreateAclRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateAclRuleRequest struct {
	ClientToken string                            `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string                            `json:"regionID,omitempty"`    /*  资源池ID  */
	AclID       string                            `json:"aclID,omitempty"`       /*  aclID  */
	Rules       []*CtvpcCreateAclRuleRulesRequest `json:"rules"`                 /*  rule 规则数组  */
}

type CtvpcCreateAclRuleRulesRequest struct {
	Direction            string  `json:"direction,omitempty"`            /*  类型,ingress, egress  */
	Priority             int32   `json:"priority"`                       /*  优先级 1 - 32766，不填默认100  */
	Protocol             string  `json:"protocol,omitempty"`             /*  all, icmp, tcp, udp, gre,  icmp6  */
	IpVersion            string  `json:"ipVersion,omitempty"`            /*  ipv4,  ipv6  */
	DestinationPort      *string `json:"destinationPort,omitempty"`      /*  开始和结束port以:隔开  */
	SourcePort           *string `json:"sourcePort,omitempty"`           /*  开始和结束port以:隔开  */
	SourceIpAddress      string  `json:"sourceIpAddress,omitempty"`      /*  类型,ingress, egress  */
	DestinationIpAddress string  `json:"destinationIpAddress,omitempty"` /*  类型,ingress, egress  */
	Action               string  `json:"action,omitempty"`               /*  accept, drop  */
	Enabled              string  `json:"enabled,omitempty"`              /*  disable, enable  */
}

type CtvpcCreateAclRuleResponse struct {
	StatusCode  int32                                  `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcCreateAclRuleReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcCreateAclRuleReturnObjResponse struct {
	AclID *string `json:"aclID,omitempty"` /*  名称  */
}
