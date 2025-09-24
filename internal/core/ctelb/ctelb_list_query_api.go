package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbListQueryApi
/* 获取转发规则列表
 */type CtelbListQueryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbListQueryApi(client *core.CtyunClient) *CtelbListQueryApi {
	return &CtelbListQueryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/list-rule",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbListQueryApi) Do(ctx context.Context, credential core.Credential, req *CtelbListQueryRequest) (*CtelbListQueryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.IDs != "" {
		ctReq.AddParam("IDs", req.IDs)
	}
	if req.LoadBalancerID != "" {
		ctReq.AddParam("loadBalancerID", req.LoadBalancerID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbListQueryResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbListQueryRequest struct {
	RegionID       string /*  区域ID  */
	IDs            string /*  转发规则ID列表，以,分隔  */
	LoadBalancerID string /*  负载均衡实例ID  */
}

type CtelbListQueryResponse struct {
	StatusCode  int32                              `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                             `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                             `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                             `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbListQueryReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       string                             `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbListQueryReturnObjResponse struct {
	RegionID       string                                       `json:"regionID,omitempty"`       /*  区域ID  */
	AzName         string                                       `json:"azName,omitempty"`         /*  可用区名称  */
	ProjectID      string                                       `json:"projectID,omitempty"`      /*  项目ID  */
	ID             string                                       `json:"ID,omitempty"`             /*  转发规则ID  */
	LoadBalancerID string                                       `json:"loadBalancerID,omitempty"` /*  负载均衡ID  */
	ListenerID     string                                       `json:"listenerID,omitempty"`     /*  监听器ID  */
	Description    string                                       `json:"description,omitempty"`    /*  描述  */
	Conditions     []*CtelbListQueryReturnObjConditionsResponse `json:"conditions"`               /*  匹配规则数据  */
	Action         *CtelbListQueryReturnObjActionResponse       `json:"action"`                   /*  规则目标  */
	Status         string                                       `json:"status,omitempty"`         /*  状态: ACTIVE / DOWN  */
	CreatedTime    string                                       `json:"createdTime,omitempty"`    /*  创建时间，为UTC格式  */
	UpdatedTime    string                                       `json:"updatedTime,omitempty"`    /*  更新时间，为UTC格式  */
}

type CtelbListQueryReturnObjConditionsResponse struct {
	RawType          string                                                     `json:"type,omitempty"`   /*  类型。取值范围：server_name（服务名称）、url_path（匹配路径）  */
	ServerNameConfig *CtelbListQueryReturnObjConditionsServerNameConfigResponse `json:"serverNameConfig"` /*  服务名称  */
	UrlPathConfig    *CtelbListQueryReturnObjConditionsUrlPathConfigResponse    `json:"urlPathConfig"`    /*  匹配路径  */
}

type CtelbListQueryReturnObjActionResponse struct {
	RawType            string                                              `json:"type,omitempty"`               /*  默认规则动作类型  */
	ForwardConfig      *CtelbListQueryReturnObjActionForwardConfigResponse `json:"forwardConfig"`                /*  转发配置  */
	RedirectListenerID string                                              `json:"redirectListenerID,omitempty"` /*  重定向监听器ID  */
}

type CtelbListQueryReturnObjConditionsServerNameConfigResponse struct {
	ServerName string `json:"serverName,omitempty"` /*  服务名称  */
}

type CtelbListQueryReturnObjConditionsUrlPathConfigResponse struct {
	UrlPaths  string `json:"urlPaths,omitempty"`  /*  匹配路径  */
	MatchType string `json:"matchType,omitempty"` /*  匹配类型。取值范围：ABSOLUTE，PREFIX，REG  */
}

type CtelbListQueryReturnObjActionForwardConfigResponse struct {
	TargetGroups []*CtelbListQueryReturnObjActionForwardConfigTargetGroupsResponse `json:"targetGroups"` /*  后端服务组  */
}

type CtelbListQueryReturnObjActionForwardConfigTargetGroupsResponse struct {
	TargetGroupID string `json:"targetGroupID,omitempty"` /*  后端服务组ID  */
	Weight        int32  `json:"weight,omitempty"`        /*  权重  */
}
