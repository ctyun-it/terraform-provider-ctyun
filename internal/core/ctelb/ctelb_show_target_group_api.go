package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbShowTargetGroupApi
/* 查看后端服务组信息
 */type CtelbShowTargetGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbShowTargetGroupApi(client *core.CtyunClient) *CtelbShowTargetGroupApi {
	return &CtelbShowTargetGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/show-target-group",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbShowTargetGroupApi) Do(ctx context.Context, credential core.Credential, req *CtelbShowTargetGroupRequest) (*CtelbShowTargetGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ID != "" {
		ctReq.AddParam("ID", req.ID)
	}
	ctReq.AddParam("targetGroupID", req.TargetGroupID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbShowTargetGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbShowTargetGroupRequest struct {
	RegionID      string /*  区域ID  */
	ID            string /*  后端服务组ID, 该字段后续废弃  */
	TargetGroupID string /*  后端服务组ID, 推荐使用该字段, 当同时使用 ID 和 targetGroupID 时，优先使用 targetGroupID  */
}

type CtelbShowTargetGroupResponse struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                   `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                   `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbShowTargetGroupReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbShowTargetGroupReturnObjResponse struct {
	RegionID      string                                              `json:"regionID,omitempty"`      /*  区域ID  */
	AzName        string                                              `json:"azName,omitempty"`        /*  可用区名称  */
	ProjectID     string                                              `json:"projectID,omitempty"`     /*  项目ID  */
	ID            string                                              `json:"ID,omitempty"`            /*  后端服务组ID  */
	Name          string                                              `json:"name,omitempty"`          /*  后端服务组名称  */
	Description   string                                              `json:"description,omitempty"`   /*  描述  */
	VpcID         string                                              `json:"vpcID,omitempty"`         /*  vpc ID  */
	HealthCheckID string                                              `json:"healthCheckID,omitempty"` /*  健康检查ID  */
	Algorithm     string                                              `json:"algorithm,omitempty"`     /*  调度算法  */
	SessionSticky *CtelbShowTargetGroupReturnObjSessionStickyResponse `json:"sessionSticky"`           /*  会话保持配置  */
	Status        string                                              `json:"status,omitempty"`        /*  状态: DOWN / ACTIVE  */
	CreatedTime   string                                              `json:"createdTime,omitempty"`   /*  创建时间，为UTC格式  */
	UpdatedTime   string                                              `json:"updatedTime,omitempty"`   /*  更新时间，为UTC格式  */
	ProxyProtocol int32                                               `json:"proxyProtocol,omitempty"`
}

type CtelbShowTargetGroupReturnObjSessionStickyResponse struct {
	SessionStickyMode string `json:"sessionStickyMode,omitempty"` /*  会话保持模式，支持取值：CLOSE（关闭）、INSERT（插入）、REWRITE（重写  */
	CookieExpire      int32  `json:"cookieExpire,omitempty"`      /*  cookie过期时间  */
	RewriteCookieName string `json:"rewriteCookieName,omitempty"` /*  cookie重写名称  */
	SourceIpTimeout   int32  `json:"sourceIpTimeout,omitempty"`   /*  源IP会话保持超时时间。  */
}
