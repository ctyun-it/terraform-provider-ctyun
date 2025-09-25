package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateOutFilterRuleApi
/* 创建出方向过滤规则
 */type CtvpcCreateOutFilterRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateOutFilterRuleApi(client *core.CtyunClient) *CtvpcCreateOutFilterRuleApi {
	return &CtvpcCreateOutFilterRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/mirrorflow/create-filter-out-rule",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateOutFilterRuleApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateOutFilterRuleRequest) (*CtvpcCreateOutFilterRuleResponse, error) {
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
	var resp CtvpcCreateOutFilterRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateOutFilterRuleRequest struct {
	RegionID         string `json:"regionID,omitempty"`       /*  区域ID  */
	MirrorFilterID   string `json:"mirrorFilterID,omitempty"` /*  过滤条件 ID  */
	DestCidr         string `json:"destCidr,omitempty"`       /*  目标 cidr  */
	SrcCidr          string `json:"srcCidr,omitempty"`        /*  源 cidr  */
	DestPort         string `json:"destPort,omitempty"`       /*  目的端口，格式为 1/2， 协议为 all / icmp，传值 -  */
	SrcPort          string `json:"srcPort,omitempty"`        /*  源端口，格式为 1/2，协议为 all / icmp，传值 -  */
	Protocol         string `json:"protocol,omitempty"`       /*  协议：all / tcp / udp / icmp  */
	EnableCollection bool   `json:"enableCollection"`         /*  是否开启采集，true 表示采集，false 表示不采集  */
}

type CtvpcCreateOutFilterRuleResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
