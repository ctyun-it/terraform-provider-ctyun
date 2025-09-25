package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDeleteAclRuleApi
/* 删除 Acl 规则列表
 */type CtvpcDeleteAclRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDeleteAclRuleApi(client *core.CtyunClient) *CtvpcDeleteAclRuleApi {
	return &CtvpcDeleteAclRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/acl-rule/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDeleteAclRuleApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDeleteAclRuleRequest) (*CtvpcDeleteAclRuleResponse, error) {
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
	var resp CtvpcDeleteAclRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDeleteAclRuleRequest struct {
	ClientToken   string   `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID      string   `json:"regionID,omitempty"`    /*  资源池ID  */
	AclID         string   `json:"aclID,omitempty"`       /*  aclID  */
	AclRuleIDList []string `json:"aclRuleIDList"`         /*  aclID list  */
}

type CtvpcDeleteAclRuleResponse struct {
	StatusCode  int32                                `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                              `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                              `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                              `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcDeleteAclRuleReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                              `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcDeleteAclRuleReturnObjResponse struct {
	AclID *string `json:"aclID,omitempty"` /*  aclid+++++wxd是对的吗  */
}
