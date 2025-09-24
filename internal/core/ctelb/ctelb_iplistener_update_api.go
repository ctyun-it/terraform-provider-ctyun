package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbIplistenerUpdateApi
/* 更新ip_listener
 */type CtelbIplistenerUpdateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbIplistenerUpdateApi(client *core.CtyunClient) *CtelbIplistenerUpdateApi {
	return &CtelbIplistenerUpdateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/iplistener/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbIplistenerUpdateApi) Do(ctx context.Context, credential core.Credential, req *CtelbIplistenerUpdateRequest) (*CtelbIplistenerUpdateResponse, error) {
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
	var resp CtelbIplistenerUpdateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbIplistenerUpdateRequest struct {
	RegionID      string `json:"regionID,omitempty"`      /*  资源池 ID  */
	IpListenerID  string `json:"ipListenerID,omitempty"`  /*  监听器 ID  */
	Name          string `json:"name,omitempty"`          /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	TargetGroupID string `json:"targetGroupID,omitempty"` /*  后端组 ID  */
	Weight        int32  `json:"weight,omitempty"`        /*  转发权重, 取值范围 0 - 256  */
}

type CtelbIplistenerUpdateResponse struct {
	StatusCode  int32                                   `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbIplistenerUpdateReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbIplistenerUpdateReturnObjResponse struct {
	IpListenerID string `json:"ipListenerID,omitempty"` /*  监听器 id  */
}
