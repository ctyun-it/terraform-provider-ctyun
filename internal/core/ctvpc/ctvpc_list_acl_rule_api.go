package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcListAclRuleApi
/* 查看 Acl 规则列表
 */type CtvpcListAclRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListAclRuleApi(client *core.CtyunClient) *CtvpcListAclRuleApi {
	return &CtvpcListAclRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/acl-rule/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListAclRuleApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListAclRuleRequest) (*CtvpcListAclRuleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ProjectID != nil {
		ctReq.AddParam("projectID", *req.ProjectID)
	}
	ctReq.AddParam("aclID", req.AclID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcListAclRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListAclRuleRequest struct {
	RegionID  string  /*  资源池ID  */
	ProjectID *string /*  企业项目 ID，默认为'0'  */
	AclID     string  /*  aclID  */
}

type CtvpcListAclRuleResponse struct {
	StatusCode  int32                                `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                              `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                              `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                              `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcListAclRuleReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                              `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListAclRuleReturnObjResponse struct {
	AclID       *string                                      `json:"aclID,omitempty"`       /*  id  */
	Name        *string                                      `json:"name,omitempty"`        /*  名称  */
	Description *string                                      `json:"description,omitempty"` /*  描述  */
	VpcID       *string                                      `json:"vpcID,omitempty"`       /*  VPC  */
	Enabled     *string                                      `json:"enabled,omitempty"`     /*  disable,enable  */
	InPolicyID  *string                                      `json:"inPolicyID,omitempty"`  /*  入规则id数组  */
	OutPolicyID *string                                      `json:"outPolicyID,omitempty"` /*  出规则id数组  */
	InRules     []*CtvpcListAclRuleReturnObjInRulesResponse  `json:"inRules"`               /*  出规则id数组  */
	OutRules    []*CtvpcListAclRuleReturnObjOutRulesResponse `json:"outRules"`              /*  出规则id数组  */
	CreatedAt   *string                                      `json:"createdAt,omitempty"`   /*  创建时间  */
	UpdatedAt   *string                                      `json:"updatedAt,omitempty"`   /*  更新时间  */
	SubnetIDs   []*string                                    `json:"subnetIDs"`             /*  acl 绑定的子网 id  */
}

type CtvpcListAclRuleReturnObjInRulesResponse struct {
	AclRuleID            *string `json:"aclRuleID,omitempty"`            /*  aclRuleID  */
	Direction            *string `json:"direction,omitempty"`            /*  类型,ingress, egress  */
	Priority             int32   `json:"priority"`                       /*  优先级  */
	Protocol             *string `json:"protocol,omitempty"`             /*  all, icmp, tcp, udp, gre,  icmp6  */
	IpVersion            *string `json:"ipVersion,omitempty"`            /*  ipv4,  ipv6  */
	DestinationPort      *string `json:"destinationPort,omitempty"`      /*  开始和结束port以:隔开  */
	SourcePort           *string `json:"sourcePort,omitempty"`           /*  开始和结束port以:隔开  */
	SourceIpAddress      *string `json:"sourceIpAddress,omitempty"`      /*  类型,ingress, egress  */
	DestinationIpAddress *string `json:"destinationIpAddress,omitempty"` /*  类型,ingress, egress  */
	Action               *string `json:"action,omitempty"`               /*  accept, drop  */
	Enabled              *string `json:"enabled,omitempty"`              /*  disable, enable  */
	Description          *string `json:"description,omitempty"`          /*  描述  */
}

type CtvpcListAclRuleReturnObjOutRulesResponse struct {
	AclRuleID            *string `json:"aclRuleID,omitempty"`            /*  aclRuleID  */
	Direction            *string `json:"direction,omitempty"`            /*  类型,ingress, egress  */
	Priority             int32   `json:"priority"`                       /*  优先级  */
	Protocol             *string `json:"protocol,omitempty"`             /*  all, icmp, tcp, udp, gre,  icmp6  */
	IpVersion            *string `json:"ipVersion,omitempty"`            /*  ipv4,  ipv6  */
	DestinationPort      *string `json:"destinationPort,omitempty"`      /*  开始和结束port以:隔开  */
	SourcePort           *string `json:"sourcePort,omitempty"`           /*  开始和结束port以:隔开  */
	SourceIpAddress      *string `json:"sourceIpAddress,omitempty"`      /*  类型,ingress, egress  */
	DestinationIpAddress *string `json:"destinationIpAddress,omitempty"` /*  类型,ingress, egress  */
	Action               *string `json:"action,omitempty"`               /*  accept, drop  */
	Enabled              *string `json:"enabled,omitempty"`              /*  disable, enable  */
	Description          *string `json:"description,omitempty"`          /*  描述  */
}
