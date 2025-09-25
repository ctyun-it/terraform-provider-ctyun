package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcUpdatePortApi
/* 修改网卡属性
 */type CtvpcUpdatePortApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcUpdatePortApi(client *core.CtyunClient) *CtvpcUpdatePortApi {
	return &CtvpcUpdatePortApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ports/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcUpdatePortApi) Do(ctx context.Context, credential core.Credential, req *CtvpcUpdatePortRequest) (*CtvpcUpdatePortResponse, error) {
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
	var resp CtvpcUpdatePortResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcUpdatePortRequest struct {
	ClientToken        string    `json:"clientToken,omitempty"`        /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID           string    `json:"regionID,omitempty"`           /*  资源池ID  */
	NetworkInterfaceID string    `json:"networkInterfaceID,omitempty"` /*  网卡ID  */
	Name               *string   `json:"name,omitempty"`               /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Description        *string   `json:"description,omitempty"`        /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
	SecurityGroupIDs   []*string `json:"securityGroupIDs"`             /*  安全组列表  */
}

type CtvpcUpdatePortResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
