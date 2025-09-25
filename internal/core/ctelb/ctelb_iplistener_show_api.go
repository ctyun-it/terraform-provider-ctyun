package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbIplistenerShowApi
/* 查看ip_listener详情
 */type CtelbIplistenerShowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbIplistenerShowApi(client *core.CtyunClient) *CtelbIplistenerShowApi {
	return &CtelbIplistenerShowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/iplistener/show",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbIplistenerShowApi) Do(ctx context.Context, credential core.Credential, req *CtelbIplistenerShowRequest) (*CtelbIplistenerShowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("ipListenerID", req.IpListenerID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbIplistenerShowResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbIplistenerShowRequest struct {
	RegionID     string /*  资源池 ID  */
	IpListenerID string /*  监听器 ID  */
}

type CtelbIplistenerShowResponse struct {
	StatusCode  int32                                 `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbIplistenerShowReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbIplistenerShowReturnObjResponse struct {
	GwElbID      string                                      `json:"gwElbID,omitempty"`      /*  网关负载均衡 ID  */
	Name         string                                      `json:"name,omitempty"`         /*  名字  */
	Description  string                                      `json:"description,omitempty"`  /*  描述  */
	IpListenerID string                                      `json:"ipListenerID,omitempty"` /*  监听器 id  */
	Action       *CtelbIplistenerShowReturnObjActionResponse `json:"action"`                 /*  转发配置  */
}

type CtelbIplistenerShowReturnObjActionResponse struct {
	RawType       string                                                   `json:"type,omitempty"` /*  默认规则动作类型: forward / redirect  */
	ForwardConfig *CtelbIplistenerShowReturnObjActionForwardConfigResponse `json:"forwardConfig"`  /*  转发配置  */
}

type CtelbIplistenerShowReturnObjActionForwardConfigResponse struct{}
