package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbCreateListenerApi
/* 创建监听器
 */type CtelbCreateListenerApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbCreateListenerApi(client *core.CtyunClient) *CtelbCreateListenerApi {
	return &CtelbCreateListenerApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/create-listener",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbCreateListenerApi) Do(ctx context.Context, credential core.Credential, req *CtelbCreateListenerRequest) (*CtelbCreateListenerResponse, error) {
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
	var resp CtelbCreateListenerResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbCreateListenerRequest struct {
	ClientToken         string                                   `json:"clientToken,omitempty"`         /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID            string                                   `json:"regionID,omitempty"`            /*  区域ID  */
	LoadBalancerID      string                                   `json:"loadBalancerID,omitempty"`      /*  负载均衡实例ID  */
	Name                string                                   `json:"name,omitempty"`                /*  唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Description         string                                   `json:"description,omitempty"`         /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
	Protocol            string                                   `json:"protocol,omitempty"`            /*  监听协议。取值范围：TCP、UDP、HTTP、HTTPS  */
	ProtocolPort        int32                                    `json:"protocolPort,omitempty"`        /*  负载均衡实例监听端口。取值：1-65535  */
	CertificateID       string                                   `json:"certificateID,omitempty"`       /*  证书ID。当protocol为HTTPS时,此参数必选  */
	CaEnabled           *bool                                    `json:"caEnabled"`                     /*  是否开启双向认证。false（不开启）、true（开启）  */
	ClientCertificateID string                                   `json:"clientCertificateID,omitempty"` /*  双向认证的证书ID  */
	DefaultAction       *CtelbCreateListenerDefaultActionRequest `json:"defaultAction"`                 /*  默认规则动作  */
	AccessControlID     string                                   `json:"accessControlID,omitempty"`     /*  访问控制ID  */
	AccessControlType   string                                   `json:"accessControlType,omitempty"`   /*  访问控制类型。取值范围：Close（未启用）、White（白名单）、Black（黑名单）  */
	ForwardedForEnabled *bool                                    `json:"forwardedForEnabled"`           /*  x forward for功能。false（未开启）、true（开启）  */
}

type CtelbCreateListenerDefaultActionRequest struct {
	RawType            string                                                `json:"type,omitempty"`               /*  默认规则动作类型。取值范围：forward、redirect  */
	ForwardConfig      *CtelbCreateListenerDefaultActionForwardConfigRequest `json:"forwardConfig"`                /*  转发配置，当type为forward时，此字段必填  */
	RedirectListenerID string                                                `json:"redirectListenerID,omitempty"` /*  重定向监听器ID，当type为redirect时，此字段必填  */
}

type CtelbCreateListenerDefaultActionForwardConfigRequest struct {
	TargetGroups []*CtelbCreateListenerDefaultActionForwardConfigTargetGroupsRequest `json:"targetGroups"` /*  后端服务组  */
}

type CtelbCreateListenerDefaultActionForwardConfigTargetGroupsRequest struct {
	TargetGroupID string `json:"targetGroupID,omitempty"` /*  后端服务组ID  */
	Weight        int32  `json:"weight,omitempty"`        /*  后端主机权重，取值范围：1-256。默认为100  */
}

type CtelbCreateListenerResponse struct {
	StatusCode  int32                                   `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbCreateListenerReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                                  `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbCreateListenerReturnObjResponse struct {
	ID string `json:"ID,omitempty"` /*  监听器 ID  */
}
