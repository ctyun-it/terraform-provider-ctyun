package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbShowListenerApi
/* 查看监听器详情
 */type CtelbShowListenerApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbShowListenerApi(client *core.CtyunClient) *CtelbShowListenerApi {
	return &CtelbShowListenerApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/show-listener",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbShowListenerApi) Do(ctx context.Context, credential core.Credential, req *CtelbShowListenerRequest) (*CtelbShowListenerResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ID != "" {
		ctReq.AddParam("ID", req.ID)
	}
	ctReq.AddParam("listenerID", req.ListenerID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbShowListenerResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbShowListenerRequest struct {
	RegionID   string /*  区域ID  */
	ID         string /*  监听器ID, 该字段后续废弃  */
	ListenerID string /*  监听器ID, 推荐使用该字段, 当同时使用 ID 和 listenerID 时，优先使用 listenerID  */
}

type CtelbShowListenerResponse struct {
	StatusCode  int32                                 `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbShowListenerReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       string                                `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbShowListenerReturnObjResponse struct {
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
	DefaultAction       *CtelbShowListenerReturnObjDefaultActionResponse `json:"defaultAction"`                 /*  默认规则动作  */
	AccessControlID     string                                           `json:"accessControlID,omitempty"`     /*  访问控制ID  */
	AccessControlType   string                                           `json:"accessControlType,omitempty"`   /*  访问控制类型: Close / White / Black  */
	ForwardedForEnabled *bool                                            `json:"forwardedForEnabled"`           /*  是否开启x forward for功能  */
	Status              string                                           `json:"status,omitempty"`              /*  监听器状态: DOWN / ACTIVE  */
	CreatedTime         string                                           `json:"createdTime,omitempty"`         /*  创建时间，为UTC格式  */
	UpdatedTime         string                                           `json:"updatedTime,omitempty"`         /*  更新时间，为UTC格式  */
	Cps                 int32                                            `json:"cps,omitempty"`
	Qps                 int32                                            `json:"qps,omitempty"`
	ResponseTimeout     int32                                            `json:"responseTimeout,omitempty"`
	EstablishTimeout    int32                                            `json:"establishTimeout,omitempty"`
	IdleTimeout         int32                                            `json:"idleTimeout,omitempty"`
	Nat64               int32                                            `json:"nat64,omitempty"`
}

type CtelbShowListenerReturnObjDefaultActionResponse struct {
	RawType            string                                                        `json:"type,omitempty"`               /*  默认规则动作类型: forward / redirect  */
	ForwardConfig      *CtelbShowListenerReturnObjDefaultActionForwardConfigResponse `json:"forwardConfig"`                /*  转发配置  */
	RedirectListenerID string                                                        `json:"redirectListenerID,omitempty"` /*  重定向监听器ID  */
}

type CtelbShowListenerReturnObjDefaultActionForwardConfigResponse struct {
	TargetGroups []CtelbShowListenerReturnObjTargetGroupResponse `json:"targetGroups"`
}
type CtelbShowListenerReturnObjTargetGroupResponse struct {
	TargetGroupID string `json:"targetGroupID,omitempty"`
	Weight        int32  `json:"weight,omitempty"`
}
