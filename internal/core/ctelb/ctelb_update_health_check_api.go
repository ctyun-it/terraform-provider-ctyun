package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbUpdateHealthCheckApi
/* 更新健康检查
 */type CtelbUpdateHealthCheckApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbUpdateHealthCheckApi(client *core.CtyunClient) *CtelbUpdateHealthCheckApi {
	return &CtelbUpdateHealthCheckApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/update-health-check",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbUpdateHealthCheckApi) Do(ctx context.Context, credential core.Credential, req *CtelbUpdateHealthCheckRequest) (*CtelbUpdateHealthCheckResponse, error) {
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
	var resp CtelbUpdateHealthCheckResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbUpdateHealthCheckRequest struct {
	ClientToken       string   `json:"clientToken,omitempty"`   /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID          string   `json:"regionID,omitempty"`      /*  区域ID,公共参数不支持修改  */
	ID                string   `json:"ID,omitempty"`            /*  健康检查ID, 后续废弃该字段  */
	HealthCheckID     string   `json:"healthCheckID,omitempty"` /*  健康检查ID, 推荐使用该字段, 当同时使用 ID 和 healthCheckID 时，优先使用 healthCheckID  */
	Name              string   `json:"name,omitempty"`          /*  唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Description       string   `json:"description,omitempty"`   /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
	Timeout           int32    `json:"timeout,omitempty"`       /*  健康检查响应的最大超时时间，取值范围：2-60秒，默认为2秒  */
	MaxRetry          int32    `json:"maxRetry,omitempty"`      /*  最大重试次数，取值范围：1-10次，默认为2次  */
	Interval          int32    `json:"interval,omitempty"`      /*  负载均衡进行健康检查的时间间隔，取值范围：1-20940秒，默认为5秒  */
	HttpMethod        string   `json:"httpMethod,omitempty"`    /*  HTTP请求的方法默认GET，{GET/HEAD/POST/PUT/DELETE/TRACE/OPTIONS/CONNECT/PATCH}（创建时仅当protocol为HTTP时必填且生效）  */
	HttpUrlPath       string   `json:"httpUrlPath,omitempty"`   /*  创建时仅当protocol为HTTP时必填且生效,支持的最大字符长度：80  */
	HttpExpectedCodes []string `json:"httpExpectedCodes"`       /*  仅当protocol为HTTP时必填且生效,支持{http_2xx/http_3xx/http_4xx/http_5xx}  */
}

type CtelbUpdateHealthCheckResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
