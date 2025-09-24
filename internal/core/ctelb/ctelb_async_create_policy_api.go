package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbAsyncCreatePolicyApi
/* 创建转发规则
 */type CtelbAsyncCreatePolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbAsyncCreatePolicyApi(client *core.CtyunClient) *CtelbAsyncCreatePolicyApi {
	return &CtelbAsyncCreatePolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/async-create-policy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbAsyncCreatePolicyApi) Do(ctx context.Context, credential core.Credential, req *CtelbAsyncCreatePolicyRequest) (*CtelbAsyncCreatePolicyResponse, error) {
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
	var resp CtelbAsyncCreatePolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbAsyncCreatePolicyRequest struct {
	ClientToken string                                     `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string                                     `json:"regionID,omitempty"`    /*  区域ID  */
	ListenerID  string                                     `json:"listenerID,omitempty"`  /*  监听器ID  */
	Name        string                                     `json:"name,omitempty"`        /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Description string                                     `json:"description,omitempty"` /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
	Conditions  []*CtelbAsyncCreatePolicyConditionsRequest `json:"conditions"`            /*  匹配规则数据  */
	TargetGroup *CtelbAsyncCreatePolicyTargetGroupRequest  `json:"targetGroup"`           /*  后端服务组  */
}

type CtelbAsyncCreatePolicyConditionsRequest struct {
	RuleType   string `json:"ruleType,omitempty"`   /*  规则类型，支持 HOST（按照域名）、PATH（请求路径）  */
	MatchType  string `json:"matchType,omitempty"`  /*  匹配类型，支持 STARTS_WITH（前缀匹配）、EQUAL_TO（精确匹配）、REGEX（正则匹配）  */
	MatchValue string `json:"matchValue,omitempty"` /*  被匹配的值，如果 ruleType 为 PATH，不能用 / 进行匹配  */
}

type CtelbAsyncCreatePolicyTargetGroupRequest struct {
	Name          string                                                 `json:"name,omitempty"`      /*  后端服务组名字  */
	Algorithm     string                                                 `json:"algorithm,omitempty"` /*  负载均衡算法，支持: rr (轮询), lc (最少链接)  */
	Targets       []*CtelbAsyncCreatePolicyTargetGroupTargetsRequest     `json:"targets"`             /*  后端服务  */
	HealthCheck   *CtelbAsyncCreatePolicyTargetGroupHealthCheckRequest   `json:"healthCheck"`         /*  健康检查配置  */
	SessionSticky *CtelbAsyncCreatePolicyTargetGroupSessionStickyRequest `json:"sessionSticky"`       /*  会话保持  */
}

type CtelbAsyncCreatePolicyTargetGroupTargetsRequest struct {
	InstanceID   string `json:"instanceID,omitempty"`   /*  后端服务主机 id  */
	ProtocolPort int32  `json:"protocolPort,omitempty"` /*  后端服务监听端口  */
	InstanceType string `json:"instanceType,omitempty"` /*  后端服务主机类型，目前支持 vm  */
	Weight       int32  `json:"weight,omitempty"`       /*  后端服务主机权重: 1 - 256  */
	Address      string `json:"address,omitempty"`      /*  后端服务主机主网卡所在的 IP  */
}

type CtelbAsyncCreatePolicyTargetGroupHealthCheckRequest struct {
	Protocol          string `json:"protocol,omitempty"`          /*  健康检查协议。取值范围：TCP、UDP、HTTP  */
	Timeout           int32  `json:"timeout,omitempty"`           /*  健康检查响应的最大超时时间，取值范围：2-60秒,默认为2秒  */
	Interval          int32  `json:"interval,omitempty"`          /*  负载均衡进行健康检查的时间间隔，取值范围：1-20940秒，默认5秒  */
	MaxRetry          int32  `json:"maxRetry,omitempty"`          /*  最大重试次数，取值范围：1-10次，默认2次  */
	HttpMethod        string `json:"httpMethod,omitempty"`        /*  仅当protocol为HTTP时必填且生效,HTTP请求的方法默认GET，{GET/HEAD}  */
	HttpUrlPath       string `json:"httpUrlPath,omitempty"`       /*  仅当protocol为HTTP时必填且生效,支持的最大字符长度：80  */
	HttpExpectedCodes string `json:"httpExpectedCodes,omitempty"` /*  仅当protocol为HTTP时必填且生效，最长支持64个字符，只能是三位数，可以以,分隔表示多个，或者以-分割表示范围，默认200  */
}

type CtelbAsyncCreatePolicyTargetGroupSessionStickyRequest struct {
	CookieName         string `json:"cookieName,omitempty"`         /*  cookie名称，当 sessionType 为 APP_COOKIE 时，为必填参数  */
	PersistenceTimeout int32  `json:"persistenceTimeout,omitempty"` /*  会话过期时间，当 sessionType 为 APP_COOKIE 或 SOURCE_IP 时，为必填参数  */
	SessionType        string `json:"sessionType,omitempty"`        /*  会话保持类型。取值范围：APP_COOKIE、HTTP_COOKIE、SOURCE_IP  */
}

type CtelbAsyncCreatePolicyResponse struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                   `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                   `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbAsyncCreatePolicyReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                                   `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbAsyncCreatePolicyReturnObjResponse struct {
	Status   string `json:"status,omitempty"`   /*  创建进度: in_progress / done  */
	Message  string `json:"message,omitempty"`  /*  进度说明  */
	PolicyID string `json:"policyID,omitempty"` /*  转发策略 ID，可能为 null  */
}
