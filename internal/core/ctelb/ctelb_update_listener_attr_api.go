package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbUpdateListenerAttrApi
/* 更新监听器
 */type CtelbUpdateListenerAttrApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbUpdateListenerAttrApi(client *core.CtyunClient) *CtelbUpdateListenerAttrApi {
	return &CtelbUpdateListenerAttrApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/update-listener-attr",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbUpdateListenerAttrApi) Do(ctx context.Context, credential core.Credential, req *CtelbUpdateListenerAttrRequest) (*CtelbUpdateListenerAttrResponse, error) {
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
	var resp CtelbUpdateListenerAttrResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbUpdateListenerAttrRequest struct {
	RegionID          string `json:"regionID,omitempty"`          /*  区域ID  */
	ListenerID        string `json:"listenerID,omitempty"`        /*  监听器 ID  */
	Name              string `json:"name,omitempty"`              /*  唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Description       string `json:"description,omitempty"`       /*  描述,支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
	AccessControlID   string `json:"accessControlID,omitempty"`   /*  访问控制ID  */
	AccessControlType string `json:"accessControlType,omitempty"` /*  访问控制类型。取值范围：Close（未启用）、White（白名单）、Black（黑名单）  */
}

type CtelbUpdateListenerAttrResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
