package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowSgRuleApi
/* 安全组规则详情。
 */type CtvpcShowSgRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowSgRuleApi(client *core.CtyunClient) *CtvpcShowSgRuleApi {
	return &CtvpcShowSgRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/describe-security-group-rule",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowSgRuleApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowSgRuleRequest) (*CtvpcShowSgRuleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("securityGroupID", req.SecurityGroupID)
	ctReq.AddParam("securityGroupRuleID", req.SecurityGroupRuleID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowSgRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowSgRuleRequest struct {
	RegionID            string /*  区域id  */
	SecurityGroupID     string /*  安全组 ID  */
	SecurityGroupRuleID string /*  安全组规则 ID  */
}

type CtvpcShowSgRuleResponse struct {
	StatusCode  int32                             `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                           `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                           `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                           `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowSgRuleReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtvpcShowSgRuleReturnObjResponse struct {
	Direction             *string `json:"direction,omitempty"`             /*  出方向-egress、入方向-ingress  */
	Priority              int32   `json:"priority"`                        /*  优先级:0~100  */
	Ethertype             *string `json:"ethertype,omitempty"`             /*  IP类型:IPv4、IPv6  */
	Protocol              *string `json:"protocol,omitempty"`              /*  协议: ANY、TCP、UDP、ICMP、ICMP6  */
	RawRange              *string `json:"range,omitempty"`                 /*  接口范围/ICMP类型:1-65535  */
	DestCidrIp            *string `json:"destCidrIp,omitempty"`            /*  远端地址:0.0.0.0/0  */
	Description           *string `json:"description,omitempty"`           /*  安全组规则描述信息。  */
	CreateTime            *string `json:"createTime,omitempty"`            /*  创建时间，UTC时间。  */
	Id                    *string `json:"id,omitempty"`                    /*  唯一标识ID  */
	SecurityGroupID       *string `json:"securityGroupID,omitempty"`       /*  安全组ID  */
	Action                *string `json:"action,omitempty"`                /*  拒绝策略:允许-accept 拒绝-drop  */
	Origin                *string `json:"origin,omitempty"`                /*  类型  */
	RemoteSecurityGroupID *string `json:"remoteSecurityGroupID,omitempty"` /*  远端安全组ID  */
	PrefixListID          *string `json:"prefixListID,omitempty"`          /*  前缀列表ID  */
}
