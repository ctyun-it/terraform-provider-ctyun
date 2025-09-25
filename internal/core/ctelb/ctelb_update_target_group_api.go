package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbUpdateTargetGroupApi
/* 更新后端服务组
 */type CtelbUpdateTargetGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbUpdateTargetGroupApi(client *core.CtyunClient) *CtelbUpdateTargetGroupApi {
	return &CtelbUpdateTargetGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/update-target-group",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbUpdateTargetGroupApi) Do(ctx context.Context, credential core.Credential, req *CtelbUpdateTargetGroupRequest) (*CtelbUpdateTargetGroupResponse, error) {
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
	var resp CtelbUpdateTargetGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbUpdateTargetGroupRequest struct {
	ClientToken   string                                      `json:"clientToken,omitempty"`   /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID      string                                      `json:"regionID,omitempty"`      /*  区域ID  */
	ProjectID     string                                      `json:"projectID,omitempty"`     /*  企业项目ID，默认为'0'  */
	ID            string                                      `json:"ID,omitempty"`            /*  后端服务组ID, 该字段后续废弃  */
	TargetGroupID string                                      `json:"targetGroupID,omitempty"` /*  后端服务组ID, 推荐使用该字段, 当同时使用 ID 和 targetGroupID 时，优先使用 targetGroupID  */
	Name          string                                      `json:"name,omitempty"`          /*  唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	HealthCheckID string                                      `json:"healthCheckID,omitempty"` /*  健康检查ID  */
	Algorithm     string                                      `json:"algorithm,omitempty"`     /*  调度算法。取值范围：rr（轮询）、wrr（带权重轮询）、lc（最少连接）、sh（源IP哈希）  */
	ProxyProtocol int32                                       `json:"proxyProtocol,omitempty"` /*  1 开启，0 关闭  */
	SessionSticky *CtelbUpdateTargetGroupSessionStickyRequest `json:"sessionSticky"`           /*  会话保持配置  */
	Protocol      string                                      `json:"protocol,omitempty"`      /*  协议，支持 TCP / UDP / HTTP / HTTPS / GENEVE  */
}

type CtelbUpdateTargetGroupSessionStickyRequest struct {
	SessionStickyMode string `json:"sessionStickyMode,omitempty"` /*  会话保持模式，支持取值：CLOSE（关闭）、INSERT（插入）、REWRITE（重写），当 algorithm 为 lc / sh 时，sessionStickyMode 必须为 CLOSE  */
	CookieExpire      int32  `json:"cookieExpire,omitempty"`      /*  cookie过期时间。INSERT模式必填  */
	RewriteCookieName string `json:"rewriteCookieName,omitempty"` /*  cookie重写名称，REWRITE模式必填  */
	SourceIpTimeout   int32  `json:"sourceIpTimeout,omitempty"`   /*  源IP会话保持超时时间。SOURCE_IP模式必填  */
}

type CtelbUpdateTargetGroupResponse struct {
	StatusCode  int32                                      `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbUpdateTargetGroupReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbUpdateTargetGroupReturnObjResponse struct {
	ID string `json:"ID,omitempty"` /*  后端服务组ID  */
}
