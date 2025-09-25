package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbAsyncCreateListenerApi
/* 创建监听器
 */type CtelbAsyncCreateListenerApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbAsyncCreateListenerApi(client *core.CtyunClient) *CtelbAsyncCreateListenerApi {
	return &CtelbAsyncCreateListenerApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/async-create-listener",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbAsyncCreateListenerApi) Do(ctx context.Context, credential core.Credential, req *CtelbAsyncCreateListenerRequest) (*CtelbAsyncCreateListenerResponse, error) {
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
	var resp CtelbAsyncCreateListenerResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbAsyncCreateListenerRequest struct {
	ClientToken         string                                      `json:"clientToken,omitempty"`         /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一,1-64  */
	RegionID            string                                      `json:"regionID,omitempty"`            /*  区域ID  */
	LoadBalanceID       string                                      `json:"loadBalanceID,omitempty"`       /*  负载均衡实例ID  */
	Name                string                                      `json:"name,omitempty"`                /*  唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Description         string                                      `json:"description,omitempty"`         /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
	Protocol            string                                      `json:"protocol,omitempty"`            /*  监听协议。取值范围：TCP、UDP、HTTP、HTTPS  */
	ProtocolPort        int32                                       `json:"protocolPort,omitempty"`        /*  负载均衡实例监听端口。取值：1-65535  */
	CertificateID       string                                      `json:"certificateID,omitempty"`       /*  证书ID。当protocol为HTTPS时,此参数必选  */
	CaEnabled           *bool                                       `json:"caEnabled"`                     /*  是否开启双向认证。false（不开启）、true（开启）  */
	ClientCertificateID string                                      `json:"clientCertificateID,omitempty"` /*  双向认证的证书ID  */
	TargetGroup         *CtelbAsyncCreateListenerTargetGroupRequest `json:"targetGroup"`                   /*  后端服务组  */
	AccessControlID     string                                      `json:"accessControlID,omitempty"`     /*  访问控制ID  */
	AccessControlType   string                                      `json:"accessControlType,omitempty"`   /*  访问控制类型。取值范围：Close（未启用）、White（白名单）、Black（黑名单）  */
	ForwardedForEnabled *bool                                       `json:"forwardedForEnabled"`           /*  x forward for功能。false（未开启）、true（开启）  */
}

type CtelbAsyncCreateListenerTargetGroupRequest struct {
	Name          string                                                   `json:"name,omitempty"`      /*  后端服务组名字  */
	Algorithm     string                                                   `json:"algorithm,omitempty"` /*  负载均衡算法，支持: rr (轮询), lc (最少链接)  */
	Targets       []*CtelbAsyncCreateListenerTargetGroupTargetsRequest     `json:"targets"`             /*  后端服务  */
	HealthCheck   *CtelbAsyncCreateListenerTargetGroupHealthCheckRequest   `json:"healthCheck"`         /*  健康检查配置  */
	SessionSticky *CtelbAsyncCreateListenerTargetGroupSessionStickyRequest `json:"sessionSticky"`       /*  会话保持  */
}

type CtelbAsyncCreateListenerTargetGroupTargetsRequest struct {
	InstanceID   string `json:"instanceID,omitempty"`   /*  后端服务主机 id  */
	ProtocolPort int32  `json:"protocolPort,omitempty"` /*  后端服务监听端口1-65535  */
	InstanceType string `json:"instanceType,omitempty"` /*  后端服务主机类型,仅支持vm类型  */
	Weight       int32  `json:"weight,omitempty"`       /*  后端服务主机权重: 1 - 256  */
	Address      string `json:"address,omitempty"`      /*  后端服务主机主网卡所在的 IP  */
}

type CtelbAsyncCreateListenerTargetGroupHealthCheckRequest struct {
	Protocol          string `json:"protocol,omitempty"`          /*  健康检查协议。取值范围：TCP、UDP、HTTP  */
	Timeout           int32  `json:"timeout,omitempty"`           /*  健康检查响应的最大超时时间，取值范围：2-60秒，默认为2秒  */
	Interval          int32  `json:"interval,omitempty"`          /*  负载均衡进行健康检查的时间间隔，取值范围：1-20940秒，默认为5秒  */
	MaxRetry          int32  `json:"maxRetry,omitempty"`          /*  最大重试次数，取值范围：1-10次，默认为2次  */
	HttpMethod        string `json:"httpMethod,omitempty"`        /*  仅当protocol为HTTP时必填且生效,HTTP请求的方法默认GET，{GET/HEAD}  */
	HttpUrlPath       string `json:"httpUrlPath,omitempty"`       /*  仅当protocol为HTTP时必填且生效,默认为'/',支持的最大字符长度：80  */
	HttpExpectedCodes string `json:"httpExpectedCodes,omitempty"` /*  仅当protocol为HTTP时必填且生效，最长支持64个字符，只能是三位数，可以以,分隔表示多个，或者以-分割表示范围，默认200  */
}

type CtelbAsyncCreateListenerTargetGroupSessionStickyRequest struct {
	SessionType        string `json:"sessionType,omitempty"`        /*  会话保持类型。取值范围：APP_COOKIE、HTTP_COOKIE、SOURCE_IP  */
	CookieName         string `json:"cookieName,omitempty"`         /*  cookie名称，当 sessionType 为 APP_COOKIE 时，为必填参数  */
	PersistenceTimeout int32  `json:"persistenceTimeout,omitempty"` /*  会话过期时间，当 sessionType 为 APP_COOKIE 或 SOURCE_IP 时，为必填参数  */
}

type CtelbAsyncCreateListenerResponse struct {
	StatusCode  int32                                      `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbAsyncCreateListenerReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                                     `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbAsyncCreateListenerReturnObjResponse struct {
	Status     string `json:"status,omitempty"`     /*  创建进度: in_progress / done  */
	Message    string `json:"message,omitempty"`    /*  进度说明  */
	ListenerID string `json:"listenerID,omitempty"` /*  监听器，可能为 null  */
}
