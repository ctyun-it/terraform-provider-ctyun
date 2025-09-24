package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcSgRulePreCheckApi
/* 安全组规则检查
 */type CtvpcSgRulePreCheckApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcSgRulePreCheckApi(client *core.CtyunClient) *CtvpcSgRulePreCheckApi {
	return &CtvpcSgRulePreCheckApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/pre-check-sg-rule",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcSgRulePreCheckApi) Do(ctx context.Context, credential core.Credential, req *CtvpcSgRulePreCheckRequest) (*CtvpcSgRulePreCheckResponse, error) {
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
	var resp CtvpcSgRulePreCheckResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcSgRulePreCheckRequest struct {
	RegionID          string                                       `json:"regionID,omitempty"`        /*  区域 id  */
	SecurityGroupID   string                                       `json:"securityGroupID,omitempty"` /*  安全组 ID。  */
	SecurityGroupRule *CtvpcSgRulePreCheckSecurityGroupRuleRequest `json:"securityGroupRule"`         /*  规则信息  */
}

type CtvpcSgRulePreCheckSecurityGroupRuleRequest struct {
	Direction  string  `json:"direction,omitempty"`  /*  入方向  */
	Action     string  `json:"action,omitempty"`     /*  拒绝策略:允许-accept 拒绝-drop  */
	Priority   int32   `json:"priority"`             /*  优先级:1~100，取值越小优先级越大  */
	Protocol   string  `json:"protocol,omitempty"`   /*  协议: ANY、TCP、UDP、ICMP(v4)  */
	Ethertype  string  `json:"ethertype,omitempty"`  /*  IP 类型:IPv4、IPv6  */
	DestCidrIp string  `json:"destCidrIp,omitempty"` /*  远端地址:0.0.0.0/0  */
	RawRange   *string `json:"range,omitempty"`      /*  安全组开放的传输层协议相关的源端端口范围  */
}

type CtvpcSgRulePreCheckResponse struct {
	StatusCode  int32                                 `json:"statusCode"`            /*  返回状态码（800 为成功，900 为失败）  */
	Message     *string                               `json:"message,omitempty"`     /*  statusCode 为 900 时的错误信息; statusCode 为 800 时为 success, 英文  */
	Description *string                               `json:"description,omitempty"` /*  statusCode 为 900 时的错误信息; statusCode 为 800 时为成功, 中文  */
	ErrorCode   *string                               `json:"errorCode,omitempty"`   /*  statusCode 为 900 时为业务细分错误码，三段式：product.module.code; statusCode 为 800 时为 SUCCESS  */
	ReturnObj   *CtvpcSgRulePreCheckReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcSgRulePreCheckReturnObjResponse struct {
	SgRuleID *string `json:"sgRuleID,omitempty"` /*  和哪个规则重复  */
}
