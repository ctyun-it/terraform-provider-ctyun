package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbListTargetGroupApi
/* 查看后端服务组列表
 */type CtelbListTargetGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbListTargetGroupApi(client *core.CtyunClient) *CtelbListTargetGroupApi {
	return &CtelbListTargetGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/list-target-group",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbListTargetGroupApi) Do(ctx context.Context, credential core.Credential, req *CtelbListTargetGroupRequest) (*CtelbListTargetGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.ClientToken != "" {
		ctReq.AddParam("clientToken", req.ClientToken)
	}
	ctReq.AddParam("regionID", req.RegionID)
	if req.IDs != "" {
		ctReq.AddParam("IDs", req.IDs)
	}
	if req.VpcID != "" {
		ctReq.AddParam("vpcID", req.VpcID)
	}
	if req.HealthCheckID != "" {
		ctReq.AddParam("healthCheckID", req.HealthCheckID)
	}
	if req.Name != "" {
		ctReq.AddParam("name", req.Name)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbListTargetGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbListTargetGroupRequest struct {
	ClientToken   string /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一  */
	RegionID      string /*  区域ID  */
	IDs           string /*  后端服务组ID列表，以,分隔  */
	VpcID         string /*  vpc ID  */
	HealthCheckID string /*  健康检查ID  */
	Name          string /*  后端服务组名称  */
}

type CtelbListTargetGroupResponse struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                   `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                   `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbListTargetGroupReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbListTargetGroupReturnObjResponse struct {
	RegionID      string                                              `json:"regionID,omitempty"`      /*  区域ID  */
	AzName        string                                              `json:"azName,omitempty"`        /*  可用区名称  */
	ProjectID     string                                              `json:"projectID,omitempty"`     /*  项目ID  */
	ID            string                                              `json:"ID,omitempty"`            /*  后端服务组ID  */
	Name          string                                              `json:"name,omitempty"`          /*  后端服务组名称  */
	Description   string                                              `json:"description,omitempty"`   /*  描述  */
	VpcID         string                                              `json:"vpcID,omitempty"`         /*  vpc ID  */
	HealthCheckID string                                              `json:"healthCheckID,omitempty"` /*  健康检查ID  */
	Algorithm     string                                              `json:"algorithm,omitempty"`     /*  调度算法  */
	SessionSticky *CtelbListTargetGroupReturnObjSessionStickyResponse `json:"sessionSticky"`           /*  会话保持配置  */
	Status        string                                              `json:"status,omitempty"`        /*  状态: DOWN / ACTIVE  */
	CreatedTime   string                                              `json:"createdTime,omitempty"`   /*  创建时间，为UTC格式  */
	UpdatedTime   string                                              `json:"updatedTime,omitempty"`   /*  更新时间，为UTC格式  */
}

type CtelbListTargetGroupReturnObjSessionStickyResponse struct {
	SessionStickyMode string `json:"sessionStickyMode,omitempty"` /*  会话保持模式，支持取值：CLOSE（关闭）、INSERT（插入）、REWRITE（重写  */
	CookieExpire      int32  `json:"cookieExpire,omitempty"`      /*  cookie过期时间  */
	RewriteCookieName string `json:"rewriteCookieName,omitempty"` /*  cookie重写名称  */
	SourceIpTimeout   int32  `json:"sourceIpTimeout,omitempty"`   /*  源IP会话保持超时时间。  */
}
