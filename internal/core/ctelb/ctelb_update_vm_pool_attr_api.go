package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbUpdateVmPoolAttrApi
/* 更新后端服务组
 */type CtelbUpdateVmPoolAttrApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbUpdateVmPoolAttrApi(client *core.CtyunClient) *CtelbUpdateVmPoolAttrApi {
	return &CtelbUpdateVmPoolAttrApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/update-vm-pool",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbUpdateVmPoolAttrApi) Do(ctx context.Context, credential core.Credential, req *CtelbUpdateVmPoolAttrRequest) (*CtelbUpdateVmPoolAttrResponse, error) {
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
	var resp CtelbUpdateVmPoolAttrResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbUpdateVmPoolAttrRequest struct {
	RegionID      string                                       `json:"regionID,omitempty"`      /*  区域ID  */
	TargetGroupID string                                       `json:"targetGroupID,omitempty"` /*  后端服务组ID  */
	Name          string                                       `json:"name,omitempty"`          /*  唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	HealthCheck   []*CtelbUpdateVmPoolAttrHealthCheckRequest   `json:"healthCheck"`             /*  当后端组已经有健康配置时，如果更新不传健康配置信息，表示移除当前后端组的健康检查配置  */
	SessionSticky []*CtelbUpdateVmPoolAttrSessionStickyRequest `json:"sessionSticky"`           /*  当后端组已经有会话配置时，如果更新不传会话配置信息，表示移除当前后端组的会话配置  */
}

type CtelbUpdateVmPoolAttrHealthCheckRequest struct {
	Protocol          string `json:"protocol,omitempty"`          /*  健康检查协议。取值范围：TCP、UDP、HTTP  */
	Timeout           int32  `json:"timeout,omitempty"`           /*  健康检查响应的最大超时时间，取值范围：2-60秒,默认2秒  */
	Interval          int32  `json:"interval,omitempty"`          /*  负载均衡进行健康检查的时间间隔，取值范围：1-20940秒，默认5秒  */
	MaxRetry          int32  `json:"maxRetry,omitempty"`          /*  最大重试次数，取值范围：1-10次，默认2次  */
	HttpMethod        string `json:"httpMethod,omitempty"`        /*  仅当protocol为HTTP时必填且生效,HTTP请求的方法默认GET，{GET/HEAD}  */
	HttpUrlPath       string `json:"httpUrlPath,omitempty"`       /*  仅当protocol为HTTP时必填且生效,支持的最大字符长度：80  */
	HttpExpectedCodes string `json:"httpExpectedCodes,omitempty"` /*  仅当protocol为HTTP时必填且生效，最长支持64个字符，只能是三位数，可以以,分隔表示多个，或者以-分割表示范围，默认200  */
}

type CtelbUpdateVmPoolAttrSessionStickyRequest struct {
	CookieName         string `json:"cookieName,omitempty"`         /*  cookie名称  */
	PersistenceTimeout int32  `json:"persistenceTimeout,omitempty"` /*  会话过期时间，1-86400  */
	SessionType        string `json:"sessionType,omitempty"`        /*  会话保持类型。取值范围：APP_COOKIE、HTTP_COOKIE、SOURCE_IP  */
}

type CtelbUpdateVmPoolAttrResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
