package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbShowRuleApi
/* 获取转发规则详情
 */type CtelbShowRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbShowRuleApi(client *core.CtyunClient) *CtelbShowRuleApi {
	return &CtelbShowRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/show-rule",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbShowRuleApi) Do(ctx context.Context, credential core.Credential, req *CtelbShowRuleRequest) (*CtelbShowRuleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ID != "" {
		ctReq.AddParam("ID", req.ID)
	}
	if req.PolicyID != "" {
		ctReq.AddParam("policyID", req.PolicyID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbShowRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbShowRuleRequest struct {
	RegionID string /*  区域ID  */
	ID       string /*  转发规则ID, 该字段后续废弃  */
	PolicyID string /*  转发规则ID, 推荐使用该字段, 当同时使用 ID 和 policyID 时，优先使用 policyID  */
}

type CtelbShowRuleResponse struct {
	StatusCode  int32                           `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                          `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                          `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                          `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbShowRuleReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       string                          `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbShowRuleReturnObjResponse struct {
	RegionID       string                                      `json:"regionID,omitempty"`       /*  区域ID  */
	AzName         string                                      `json:"azName,omitempty"`         /*  可用区名称  */
	ProjectID      string                                      `json:"projectID,omitempty"`      /*  项目ID  */
	ID             string                                      `json:"ID,omitempty"`             /*  转发规则ID  */
	LoadBalancerID string                                      `json:"loadBalancerID,omitempty"` /*  负载均衡ID  */
	ListenerID     string                                      `json:"listenerID,omitempty"`     /*  监听器ID  */
	Description    string                                      `json:"description,omitempty"`    /*  描述  */
	Conditions     []*CtelbShowRuleReturnObjConditionsResponse `json:"conditions"`               /*  匹配规则数据  */
	Action         *CtelbShowRuleReturnObjActionResponse       `json:"action"`                   /*  规则目标  */
	Status         string                                      `json:"status,omitempty"`         /*  状态: ACTIVE / DOWN  */
	CreatedTime    string                                      `json:"createdTime,omitempty"`    /*  创建时间，为UTC格式  */
	UpdatedTime    string                                      `json:"updatedTime,omitempty"`    /*  更新时间，为UTC格式  */
}

type CtelbShowRuleReturnObjConditionsResponse struct {
	RawType          string                                                    `json:"type,omitempty"`   /*  类型。取值范围：server_name（服务名称）、url_path（匹配路径）  */
	ServerNameConfig *CtelbShowRuleReturnObjConditionsServerNameConfigResponse `json:"serverNameConfig"` /*  服务名称  */
	UrlPathConfig    *CtelbShowRuleReturnObjConditionsUrlPathConfigResponse    `json:"urlPathConfig"`    /*  匹配路径  */
}

type CtelbShowRuleReturnObjActionResponse struct {
	RawType            string                                             `json:"type,omitempty"`               /*  默认规则动作类型: forward / redirect  */
	ForwardConfig      *CtelbShowRuleReturnObjActionForwardConfigResponse `json:"forwardConfig"`                /*  转发配置  */
	RedirectListenerID string                                             `json:"redirectListenerID,omitempty"` /*  重定向监听器ID  */
}

type CtelbShowRuleReturnObjConditionsServerNameConfigResponse struct {
	ServerName string `json:"serverName,omitempty"` /*  服务名称  */
}

type CtelbShowRuleReturnObjConditionsUrlPathConfigResponse struct {
	UrlPaths  string `json:"urlPaths,omitempty"`  /*  匹配路径  */
	MatchType string `json:"matchType,omitempty"` /*  匹配类型。取值范围：ABSOLUTE，PREFIX，REG  */
}

type CtelbShowRuleReturnObjActionForwardConfigResponse struct {
	TargetGroups []*CtelbCreateRuleActionForwardConfigTargetGroupsResponse `json:"targetGroups"` /*  后端服务组  */
}
type CtelbCreateRuleActionForwardConfigTargetGroupsResponse struct {
	TargetGroupID string `json:"targetGroupID,omitempty"` /*  后端服务组ID  */
	Weight        int32  `json:"weight,omitempty"`        /*  权重，取值范围：1-256。默认为100  */
}
