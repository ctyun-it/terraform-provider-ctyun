package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcUpdateMirrorFlowFilterRuleApi
/* 更新过滤规则
 */type CtvpcUpdateMirrorFlowFilterRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcUpdateMirrorFlowFilterRuleApi(client *core.CtyunClient) *CtvpcUpdateMirrorFlowFilterRuleApi {
	return &CtvpcUpdateMirrorFlowFilterRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/mirrorflow/update-filter-rule",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcUpdateMirrorFlowFilterRuleApi) Do(ctx context.Context, credential core.Credential, req *CtvpcUpdateMirrorFlowFilterRuleRequest) (*CtvpcUpdateMirrorFlowFilterRuleResponse, error) {
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
	var resp CtvpcUpdateMirrorFlowFilterRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcUpdateMirrorFlowFilterRuleRequest struct {
	RegionID           string `json:"regionID,omitempty"`           /*  区域ID  */
	MirrorFilterRuleID string `json:"mirrorFilterRuleID,omitempty"` /*  名称  */
	DestCidr           string `json:"destCidr,omitempty"`           /*  目标 cidr  */
	SrcCidr            string `json:"srcCidr,omitempty"`            /*  源 cidr  */
	DestPort           string `json:"destPort,omitempty"`           /*  目的端口，格式为 1/2， 协议为 all时，传值 -  */
	SrcPort            string `json:"srcPort,omitempty"`            /*  源端口，格式为 1/2，协议为 all时，传值 -  */
	Protocol           string `json:"protocol,omitempty"`           /*  协议：all / tcp / udp / icmp  */
	EnableCollection   bool   `json:"enableCollection"`             /*  是否开启采集，true 表示采集，false 表示不采集  */
}

type CtvpcUpdateMirrorFlowFilterRuleResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
