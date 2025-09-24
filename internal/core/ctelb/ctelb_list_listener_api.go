package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbListListenerApi
/* 查看监听器列表
 */type CtelbListListenerApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbListListenerApi(client *core.CtyunClient) *CtelbListListenerApi {
	return &CtelbListListenerApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/list-listener",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbListListenerApi) Do(ctx context.Context, credential core.Credential, req *CtelbListListenerRequest) (*CtelbListListenerResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.ClientToken != "" {
		ctReq.AddParam("clientToken", req.ClientToken)
	}
	ctReq.AddParam("regionID", req.RegionID)
	if req.ProjectID != "" {
		ctReq.AddParam("projectID", req.ProjectID)
	}
	if req.IDs != "" {
		ctReq.AddParam("IDs", req.IDs)
	}
	if req.Name != "" {
		ctReq.AddParam("name", req.Name)
	}
	if req.LoadBalancerID != "" {
		ctReq.AddParam("loadBalancerID", req.LoadBalancerID)
	}
	if req.AccessControlID != "" {
		ctReq.AddParam("accessControlID", req.AccessControlID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbListListenerResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbListListenerRequest struct {
	ClientToken     string /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一, 公共参数不支持修改, 长度 1 - 64  */
	RegionID        string /*  区域ID  */
	ProjectID       string /*  企业项目ID，默认为'0'  */
	IDs             string /*  监听器ID列表，以','分隔  */
	Name            string /*  监听器名称  */
	LoadBalancerID  string /*  负载均衡实例ID  */
	AccessControlID string /*  访问控制ID  */
}

type CtelbListListenerResponse struct {
	StatusCode  int32                                 `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbListListenerReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       string                                `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbListListenerReturnObjResponse struct {
	RegionID            string                                           `json:"regionID,omitempty"`            /*  区域ID  */
	AzName              string                                           `json:"azName,omitempty"`              /*  可用区名称  */
	ProjectID           string                                           `json:"projectID,omitempty"`           /*  项目ID  */
	ID                  string                                           `json:"ID,omitempty"`                  /*  监听器ID  */
	Name                string                                           `json:"name,omitempty"`                /*  监听器名称  */
	Description         string                                           `json:"description,omitempty"`         /*  描述  */
	LoadBalancerID      string                                           `json:"loadBalancerID,omitempty"`      /*  负载均衡实例ID  */
	Protocol            string                                           `json:"protocol,omitempty"`            /*  监听协议: TCP / UDP / HTTP / HTTPS  */
	ProtocolPort        int32                                            `json:"protocolPort,omitempty"`        /*  监听端口  */
	CertificateID       string                                           `json:"certificateID,omitempty"`       /*  证书ID  */
	CaEnabled           *bool                                            `json:"caEnabled"`                     /*  是否开启双向认证  */
	ClientCertificateID string                                           `json:"clientCertificateID,omitempty"` /*  双向认证的证书ID  */
	DefaultAction       *CtelbListListenerReturnObjDefaultActionResponse `json:"defaultAction"`                 /*  默认规则动作  */
	AccessControlID     string                                           `json:"accessControlID,omitempty"`     /*  访问控制ID  */
	AccessControlType   string                                           `json:"accessControlType,omitempty"`   /*  访问控制类型: Close / White / Black  */
	ForwardedForEnabled *bool                                            `json:"forwardedForEnabled"`           /*  是否开启x forward for功能  */
	Status              string                                           `json:"status,omitempty"`              /*  监听器状态: DOWN / ACTIVE  */
	CreatedTime         string                                           `json:"createdTime,omitempty"`         /*  创建时间，为UTC格式  */
	UpdatedTime         string                                           `json:"updatedTime,omitempty"`         /*  更新时间，为UTC格式  */
}

type CtelbListListenerReturnObjDefaultActionResponse struct {
	RawType            string                                                        `json:"type,omitempty"`               /*  默认规则动作类型: forward / redirect  */
	ForwardConfig      *CtelbListListenerReturnObjDefaultActionForwardConfigResponse `json:"forwardConfig"`                /*  转发配置  */
	RedirectListenerID string                                                        `json:"redirectListenerID,omitempty"` /*  重定向监听器ID  */
}

type CtelbListListenerReturnObjDefaultActionForwardConfigResponse struct {
	TargetGroups []*CtelbListListenerReturnObjDefaultActionForwardConfigTargetGroupsResponse `json:"targetGroups"` /*  后端服务组  */
}

type CtelbListListenerReturnObjDefaultActionForwardConfigTargetGroupsResponse struct {
	TargetGroupID string `json:"targetGroupID,omitempty"` /*  后端服务组ID  */
	Weight        int32  `json:"weight,omitempty"`        /*  权重  */
}
