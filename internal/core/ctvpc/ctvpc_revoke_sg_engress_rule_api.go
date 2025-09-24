package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcRevokeSgEngressRuleApi
/* 删除一条入方向安全组规则，撤销安全组出方向的权限设置。
 */type CtvpcRevokeSgEngressRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcRevokeSgEngressRuleApi(client *core.CtyunClient) *CtvpcRevokeSgEngressRuleApi {
	return &CtvpcRevokeSgEngressRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/revoke-security-group-ingress",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcRevokeSgEngressRuleApi) Do(ctx context.Context, credential core.Credential, req *CtvpcRevokeSgEngressRuleRequest) (*CtvpcRevokeSgEngressRuleResponse, error) {
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
	var resp CtvpcRevokeSgEngressRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcRevokeSgEngressRuleRequest struct {
	RegionID            string `json:"regionID,omitempty"`            /*  区域id  */
	SecurityGroupID     string `json:"securityGroupID,omitempty"`     /*  安全组ID  */
	SecurityGroupRuleID string `json:"securityGroupRuleID,omitempty"` /*  安全组规则ID  */
	ClientToken         string `json:"clientToken,omitempty"`         /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
}

type CtvpcRevokeSgEngressRuleResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
