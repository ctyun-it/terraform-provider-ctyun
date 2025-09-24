package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcAttachPortBmApi
/* 调用此接口可将弹性公网IP（Elastic IP Address，简称EIP）与物理机网卡进行绑定。
 */type CtvpcAttachPortBmApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcAttachPortBmApi(client *core.CtyunClient) *CtvpcAttachPortBmApi {
	return &CtvpcAttachPortBmApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/eip/attach-port-bm",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcAttachPortBmApi) Do(ctx context.Context, credential core.Credential, req *CtvpcAttachPortBmRequest) (*CtvpcAttachPortBmResponse, error) {
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
	var resp CtvpcAttachPortBmResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcAttachPortBmRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID  */
	EipID       string `json:"eipID,omitempty"`       /*  绑定云产品实例的 EIP 的 ID  */
	SourceID    string `json:"sourceID,omitempty"`    /*  绑定云产品实例的 物理机 的 ID  */
	PortID      string `json:"portID,omitempty"`      /*  绑定的实例的 ID  */
}

type CtvpcAttachPortBmResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
