package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcModifySgEngressRuleApi
/* 修改安全组出方向规则的描述信息。该接口只能修改出方向描述信息。如果您需要修改安全组规则的策略、端口范围等信息，请在管理控制台修改。
 */type CtvpcModifySgEngressRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcModifySgEngressRuleApi(client *core.CtyunClient) *CtvpcModifySgEngressRuleApi {
	return &CtvpcModifySgEngressRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/modify-security-group-egress",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcModifySgEngressRuleApi) Do(ctx context.Context, credential core.Credential, req *CtvpcModifySgEngressRuleRequest) (*CtvpcModifySgEngressRuleResponse, error) {
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
	var resp CtvpcModifySgEngressRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcModifySgEngressRuleRequest struct {
	RegionID              string  `json:"regionID,omitempty"`              /*  区域id  */
	SecurityGroupID       string  `json:"securityGroupID,omitempty"`       /*  安全组ID。  */
	SecurityGroupRuleID   string  `json:"securityGroupRuleID,omitempty"`   /*  安全组规则ID。  */
	ClientToken           string  `json:"clientToken,omitempty"`           /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	Action                *string `json:"action,omitempty"`                /*  拒绝策略:允许-accept 拒绝-drop  */
	Priority              int32   `json:"priority"`                        /*  优先级:1~100，取值越小优先级越大  */
	Protocol              *string `json:"protocol,omitempty"`              /*  协议: ANY、TCP、UDP、ICMP(v4)  */
	RemoteSecurityGroupID *string `json:"remoteSecurityGroupID,omitempty"` /*  远端安全组id  */
	DestCidrIp            *string `json:"destCidrIp,omitempty"`            /*  cidr  */
	RemoteType            int32   `json:"remoteType"`                      /*  远端类型，0 表示 destCidrIp，1 表示 remoteSecurityGroupID, 2 表示 prefixlistID，默认为 0  */
	PrefixListID          *string `json:"prefixListID,omitempty"`          /*  前缀列表  */
}

type CtvpcModifySgEngressRuleResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
