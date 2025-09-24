package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbShowHealthCheckApi
/* 查看健康检查详情
 */type CtelbShowHealthCheckApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbShowHealthCheckApi(client *core.CtyunClient) *CtelbShowHealthCheckApi {
	return &CtelbShowHealthCheckApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/show-health-check",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbShowHealthCheckApi) Do(ctx context.Context, credential core.Credential, req *CtelbShowHealthCheckRequest) (*CtelbShowHealthCheckResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.Id != "" {
		ctReq.AddParam("id", req.Id)
	}
	ctReq.AddParam("healthCheckID", req.HealthCheckID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbShowHealthCheckResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbShowHealthCheckRequest struct {
	RegionID      string /*  区域ID  */
	Id            string /*  健康检查ID, 后续废弃该字段  */
	HealthCheckID string /*  健康检查ID, 推荐使用该字段, 当同时使用 id 和 healthCheckID 时，优先使用 healthCheckID  */
}

type CtelbShowHealthCheckResponse struct {
	StatusCode  int32                                  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbShowHealthCheckReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       string                                 `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbShowHealthCheckReturnObjResponse struct {
	RegionID          string `json:"regionID,omitempty"`          /*  区域ID  */
	AzName            string `json:"azName,omitempty"`            /*  可用区名称  */
	ProjectID         string `json:"projectID,omitempty"`         /*  项目ID  */
	ID                string `json:"ID,omitempty"`                /*  健康检查ID  */
	Name              string `json:"name,omitempty"`              /*  健康检查名称  */
	Description       string `json:"description,omitempty"`       /*  描述  */
	Protocol          string `json:"protocol,omitempty"`          /*  健康检查协议: TCP / UDP / HTTP  */
	ProtocolPort      int32  `json:"protocolPort,omitempty"`      /*  健康检查端口  */
	Timeout           int32  `json:"timeout,omitempty"`           /*  健康检查响应的最大超时时间  */
	Interval          int32  `json:"interval,omitempty"`          /*  负载均衡进行健康检查的时间间隔  */
	MaxRetry          int32  `json:"maxRetry,omitempty"`          /*  最大重试次数  */
	HttpMethod        string `json:"httpMethod,omitempty"`        /*  HTTP请求的方法  */
	HttpUrlPath       string `json:"httpUrlPath,omitempty"`       /*  HTTP请求url路径  */
	HttpExpectedCodes string `json:"httpExpectedCodes,omitempty"` /*  HTTP预期码  */
	Status            int32  `json:"status,omitempty"`            /*  状态 1 表示 UP, 0 表示 DOWN  */
	CreateTime        string `json:"createTime,omitempty"`        /*  创建时间，为UTC格式  */
}
