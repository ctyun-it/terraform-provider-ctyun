package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbListHealthCheckApi
/* 获取健康检查列表
 */type CtelbListHealthCheckApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbListHealthCheckApi(client *core.CtyunClient) *CtelbListHealthCheckApi {
	return &CtelbListHealthCheckApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/list-health-check",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbListHealthCheckApi) Do(ctx context.Context, credential core.Credential, req *CtelbListHealthCheckRequest) (*CtelbListHealthCheckResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.ClientToken != "" {
		ctReq.AddParam("clientToken", req.ClientToken)
	}
	ctReq.AddParam("regionID", req.RegionID)
	for _, ele0 := range req.IDs {
		if ele0 != "" {
			ctReq.AddParam("IDs", ele0)
		}
	}
	if req.Name != "" {
		ctReq.AddParam("name", req.Name)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbListHealthCheckResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbListHealthCheckRequest struct {
	ClientToken string   /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一  */
	RegionID    string   /*  区域ID  */
	IDs         []string /*  健康检查ID列表  */
	Name        string   /*  健康检查名称, 只能由数字，字母，-组成不能以数字和-开头，最大长度32  */
}

type CtelbListHealthCheckResponse struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                   `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                   `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbListHealthCheckReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       string                                   `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbListHealthCheckReturnObjResponse struct {
	RegionID          string `json:"regionID,omitempty"`          /*  区域ID  */
	AzName            string `json:"azName,omitempty"`            /*  可用区名称  */
	ProjectID         string `json:"projectID,omitempty"`         /*  项目ID  */
	ID                string `json:"ID,omitempty"`                /*  健康检查ID  */
	Name              string `json:"name,omitempty"`              /*  健康检查名称  */
	Description       string `json:"description,omitempty"`       /*  描述  */
	Protocol          string `json:"protocol,omitempty"`          /*  健康检查协议: TCP / UDP / HTTP  */
	ProtocolPort      int32  `json:"protocolPort,omitempty"`      /*  健康检查端口  */
	Timeout           int32  `json:"timeout,omitempty"`           /*  健康检查响应的最大超时时间  */
	Integererval      int32  `json:"Integererval,omitempty"`      /*  负载均衡进行健康检查的时间间隔  */
	MaxRetry          int32  `json:"maxRetry,omitempty"`          /*  最大重试次数  */
	HttpMethod        string `json:"httpMethod,omitempty"`        /*  HTTP请求的方法  */
	HttpUrlPath       string `json:"httpUrlPath,omitempty"`       /*  HTTP请求url路径  */
	HttpExpectedCodes string `json:"httpExpectedCodes,omitempty"` /*  HTTP预期码  */
	Status            int32  `json:"status,omitempty"`            /*  状态 1 表示 UP, 0 表示 DOWN  */
	CreateTime        string `json:"createTime,omitempty"`        /*  创建时间，为UTC格式  */
}
