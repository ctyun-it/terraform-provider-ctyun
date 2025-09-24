package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowFlowSessionApi
/* 查看流量会话详情
 */type CtvpcShowFlowSessionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowFlowSessionApi(client *core.CtyunClient) *CtvpcShowFlowSessionApi {
	return &CtvpcShowFlowSessionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/flowsession/show",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowFlowSessionApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowFlowSessionRequest) (*CtvpcShowFlowSessionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("flowSessionID", req.FlowSessionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowFlowSessionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowFlowSessionRequest struct {
	RegionID      string /*  区域ID  */
	FlowSessionID string /*  名称  */
}

type CtvpcShowFlowSessionResponse struct {
	StatusCode  int32                                  `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowFlowSessionReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowFlowSessionReturnObjResponse struct {
	MirrorFilterID *string `json:"mirrorFilterID,omitempty"` /*  过滤规则 ID  */
	SrcPort        *string `json:"srcPort,omitempty"`        /*  源弹性网卡 ID  */
	DstPort        *string `json:"dstPort,omitempty"`        /*  源弹性网卡 ID  */
	CreatedTime    *string `json:"createdTime,omitempty"`    /*  创建时间  */
	FlowSessionID  *string `json:"flowSessionID,omitempty"`  /*  会话 ID  */
	Name           *string `json:"name,omitempty"`           /*  会话名称  */
	Description    *string `json:"description,omitempty"`    /*  会话描述  */
	Vni            int32   `json:"vni"`                      /*  VXLAN 网络标识符  */
	DstPortType    *string `json:"dstPortType,omitempty"`    /*  目标网卡类型: VM  */
	Status         *string `json:"status,omitempty"`         /*  会话状态：on / off  */
}
