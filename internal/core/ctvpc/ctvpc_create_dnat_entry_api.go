package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateDnatEntryApi
/* 创建 dnat 规则
 */type CtvpcCreateDnatEntryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateDnatEntryApi(client *core.CtyunClient) *CtvpcCreateDnatEntryApi {
	return &CtvpcCreateDnatEntryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/create-dnat-entry",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateDnatEntryApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateDnatEntryRequest) (*CtvpcCreateDnatEntryResponse, error) {
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
	var resp CtvpcCreateDnatEntryResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateDnatEntryRequest struct {
	RegionID           string  `json:"regionID,omitempty"`         /*  区域id  */
	NatGatewayID       string  `json:"natGatewayID,omitempty"`     /*  natgw-1o5sdqb7i2  */
	ExternalID         string  `json:"externalID,omitempty"`       /*  弹性公网id  */
	ExternalPort       int32   `json:"externalPort"`               /*  弹性IP公网端口, 1 - 1024  */
	VirtualMachineID   *string `json:"virtualMachineID,omitempty"` /*  云主机  */
	VirtualMachineType int32   `json:"virtualMachineType"`         /*  云主机类型1-选择云主机，serverType字段必传
	2-自定义，internalIp必传  */
	InternalIp   *string `json:"internalIp,omitempty"`  /*  内部 IP  */
	ServerType   *string `json:"serverType,omitempty"`  /*  当 virtualMachineType 为 1 时，serverType 必传，支持: VM / BM （仅支持大写）  */
	InternalPort int32   `json:"internalPort"`          /*  主机内网端口，1-65535  */
	Protocol     string  `json:"protocol,omitempty"`    /*  支持协议：tcp/udp  */
	ClientToken  string  `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	Description  *string `json:"description,omitempty"` /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
}

type CtvpcCreateDnatEntryResponse struct {
	StatusCode  int32                                  `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCreateDnatEntryReturnObjResponse `json:"returnObj"`             /*  业务数据  */
	Error       *string                                `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcCreateDnatEntryReturnObjResponse struct {
	Dnat *CtvpcCreateDnatEntryReturnObjDnatResponse `json:"dnat"` /*  业务数据  */
}

type CtvpcCreateDnatEntryReturnObjDnatResponse struct {
	Status  *string `json:"status,omitempty"`  /*  绑定状态，取值 in_progress / done  */
	Message *string `json:"message,omitempty"` /*  绑定状态提示信息  */
	DnatID  *string `json:"dnatID,omitempty"`  /*  dnat id, 当 status != done 时，dnatID 为 null  */
}
