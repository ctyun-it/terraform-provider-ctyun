package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbListVmPoolApi
/* 查看后端服务组列表
 */type CtelbListVmPoolApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbListVmPoolApi(client *core.CtyunClient) *CtelbListVmPoolApi {
	return &CtelbListVmPoolApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/list-vm-pool",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbListVmPoolApi) Do(ctx context.Context, credential core.Credential, req *CtelbListVmPoolRequest) (*CtelbListVmPoolResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.TargetGroupID != "" {
		ctReq.AddParam("targetGroupID", req.TargetGroupID)
	}
	if req.Name != "" {
		ctReq.AddParam("name", req.Name)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbListVmPoolResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbListVmPoolRequest struct {
	RegionID      string /*  区域ID  */
	TargetGroupID string /*  后端服务组ID  */
	Name          string /*  后端服务组名称  */
}

type CtelbListVmPoolResponse struct {
	StatusCode  int32                               `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                              `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                              `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                              `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbListVmPoolReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbListVmPoolReturnObjResponse struct {
	RegionID      string                                           `json:"regionID,omitempty"`      /*  区域ID  */
	AzName        string                                           `json:"azName,omitempty"`        /*  可用区名称, 默认为 null  */
	ProjectID     string                                           `json:"projectID,omitempty"`     /*  项目ID  */
	ID            string                                           `json:"ID,omitempty"`            /*  后端服务组ID  */
	Name          string                                           `json:"name,omitempty"`          /*  后端服务组名称  */
	Description   string                                           `json:"description,omitempty"`   /*  描述  */
	VpcID         string                                           `json:"vpcID,omitempty"`         /*  vpc ID, 默认为 null  */
	HealthCheckID string                                           `json:"healthCheckID,omitempty"` /*  健康检查ID  */
	Algorithm     string                                           `json:"algorithm,omitempty"`     /*  调度算法  */
	SessionSticky []*CtelbListVmPoolReturnObjSessionStickyResponse `json:"sessionSticky"`           /*  会话保持配置  */
	Status        string                                           `json:"status,omitempty"`        /*  状态  */
	CreatedTime   string                                           `json:"createdTime,omitempty"`   /*  创建时间，为UTC格式  */
	UpdatedTime   string                                           `json:"updatedTime,omitempty"`   /*  更新时间，为UTC格式, 默认为 null  */
}

type CtelbListVmPoolReturnObjSessionStickyResponse struct {
	SessionStickyMode string `json:"sessionStickyMode,omitempty"` /*  会话保持模式  */
	CookieExpire      int32  `json:"cookieExpire,omitempty"`      /*  cookie过期时间  */
	RewriteCookieName string `json:"rewriteCookieName,omitempty"` /*  cookie重写名称  */
	SourceIpTimeout   int32  `json:"sourceIpTimeout,omitempty"`   /*  源IP会话保持超时时间。  */
}
