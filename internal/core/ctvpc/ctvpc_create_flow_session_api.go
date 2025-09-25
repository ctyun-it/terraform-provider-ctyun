package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateFlowSessionApi
/* 创建流量会话
 */type CtvpcCreateFlowSessionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateFlowSessionApi(client *core.CtyunClient) *CtvpcCreateFlowSessionApi {
	return &CtvpcCreateFlowSessionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/flowsession/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateFlowSessionApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateFlowSessionRequest) (*CtvpcCreateFlowSessionResponse, error) {
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
	var resp CtvpcCreateFlowSessionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateFlowSessionRequest struct {
	RegionID       string `json:"regionID,omitempty"`       /*  区域ID  */
	MirrorFilterID string `json:"mirrorFilterID,omitempty"` /*  过滤条件 ID  */
	SrcPort        string `json:"srcPort,omitempty"`        /*  源弹性网卡 ID，绑定类型只能是vm/bm  */
	DstPort        string `json:"dstPort,omitempty"`        /*  目的弹性网卡 ID，绑定类型只能是vm  */
	SubnetID       string `json:"subnetID,omitempty"`       /*  子网 ID  */
	Vni            int32  `json:"vni"`                      /*  VXLAN 网络标识符, 0 - 1677215  */
	Name           string `json:"name,omitempty"`           /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
}

type CtvpcCreateFlowSessionResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
