package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbCreateRuleApi
/* 创建转发规则
 */type CtelbCreateRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbCreateRuleApi(client *core.CtyunClient) *CtelbCreateRuleApi {
	return &CtelbCreateRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/create-rule",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbCreateRuleApi) Do(ctx context.Context, credential core.Credential, req *CtelbCreateRuleRequest) (*CtelbCreateRuleResponse, error) {
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
	var resp CtelbCreateRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbCreateRuleRequest struct {
	ClientToken string                              `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string                              `json:"regionID,omitempty"`    /*  区域ID  */
	ListenerID  string                              `json:"listenerID,omitempty"`  /*  监听器ID  */
	Priority    int32                               `json:"priority,omitempty"`    /*  优先级，数字越小优先级越高，取值范围为：1-100(目前不支持配置此参数,只取默认值100)  */
	Conditions  []*CtelbCreateRuleConditionsRequest `json:"conditions"`            /*  匹配规则数据  */
	Action      *CtelbCreateRuleActionRequest       `json:"action"`                /*  规则目标  */
}

type CtelbCreateRuleConditionsRequest struct {
	RawType          string                                            `json:"type,omitempty"`   /*  类型。取值范围：server_name（服务名称）、url_path（匹配路径）  */
	ServerNameConfig *CtelbCreateRuleConditionsServerNameConfigRequest `json:"serverNameConfig"` /*  服务名称  */
	UrlPathConfig    *CtelbCreateRuleConditionsUrlPathConfigRequest    `json:"urlPathConfig"`    /*  匹配路径  */
}

type CtelbCreateRuleActionRequest struct {
	RawType            string                                     `json:"type,omitempty"`               /*  默认规则动作类型。取值范围：forward、redirect、deny(目前暂不支持配置为deny)  */
	ForwardConfig      *CtelbCreateRuleActionForwardConfigRequest `json:"forwardConfig"`                /*  转发配置  */
	RedirectListenerID string                                     `json:"redirectListenerID,omitempty"` /*  重定向监听器ID，当type为redirect时，此字段必填  */
}

type CtelbCreateRuleConditionsServerNameConfigRequest struct {
	ServerName string `json:"serverName,omitempty"` /*  服务名称  */
}

type CtelbCreateRuleConditionsUrlPathConfigRequest struct {
	UrlPaths  string `json:"urlPaths,omitempty"`  /*  匹配路径  */
	MatchType string `json:"matchType,omitempty"` /*  匹配类型。取值范围：ABSOLUTE，PREFIX，REG  */
}

type CtelbCreateRuleActionForwardConfigRequest struct {
	TargetGroups []*CtelbCreateRuleActionForwardConfigTargetGroupsRequest `json:"targetGroups"` /*  后端服务组  */
}

type CtelbCreateRuleActionForwardConfigTargetGroupsRequest struct {
	TargetGroupID string `json:"targetGroupID,omitempty"` /*  后端服务组ID  */
	Weight        int32  `json:"weight,omitempty"`        /*  权重，取值范围：1-256。默认为100  */
}

type CtelbCreateRuleResponse struct {
	StatusCode  int32                               `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                              `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                              `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                              `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbCreateRuleReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                              `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbCreateRuleReturnObjResponse struct {
	ID string `json:"ID,omitempty"` /*  转发规则 ID  */
}
