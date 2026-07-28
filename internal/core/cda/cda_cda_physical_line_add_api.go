package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaPhysicalLineAddApi
/* 创建物理专线 */
type CdaCdaPhysicalLineAddApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaPhysicalLineAddApi(client *core.CtyunClient) *CdaCdaPhysicalLineAddApi {
	return &CdaCdaPhysicalLineAddApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cda/physical-line/add",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaPhysicalLineAddApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaPhysicalLineAddRequest) (*CdaCdaPhysicalLineAddResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaPhysicalLineAddRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CdaCdaPhysicalLineAddResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaPhysicalLineAddRequest struct {
	LineName          string  `json:"lineName"`                    /*  物理专线名字  */
	Account           string  `json:"account"`                     /*  天翼云客户邮箱  */
	ResourcePool      string  `json:"resourcePool"`                /*  资源池ID  */
	ResourcePoolName  string  `json:"resourcePoolName"`            /*  资源池名字  */
	IpVersion         string  `json:"ipVersion"`                   /*  IPV4/DUALSTACK/IPV6(三选一)  */
	LocalConnectIP    *string `json:"localConnectIP,omitempty"`    /*  本端互联IP(IPV4或DUALSTACK必填)  */
	RemoteConnectIP   *string `json:"remoteConnectIP,omitempty"`   /*  远端互联IP(IPV4或DUALSTACK必填)  */
	LocalConnectIPv6  *string `json:"localConnectIPv6,omitempty"`  /*  本端互联IPv6(IPV6或DUALSTACK必填)  */
	RemoteConnectIPv6 *string `json:"remoteConnectIPv6,omitempty"` /*  本端互联IPv6(IPV6或DUALSTACK必填)  */
	Hostname          string  `json:"hostname"`                    /*  交换机hostname  */
	PortType          string  `json:"portType"`                    /*  端口类型(10G、1G)  */
	PortName          string  `json:"portName"`                    /*  端口名字  */
	Bandwidth         int32   `json:"bandwidth"`                   /*  带宽(M)  */
	LineType          string  `json:"lineType"`                    /*  物理专线类型(PON/IPRAN等)  */
	Tag               int32   `json:"tag"`                         /*  是否带vlan tag, 1 为带， 0为不带  */
	Vlan              *int32  `json:"vlan,omitempty"`              /*  接入端口放行vlan, tag为1 ，此项必填  */
	IsShared          int32   `json:"isShared"`                    /*  端口类别: 独享(0), 共享(1)默认  */
	AccessPoint       *string `json:"accessPoint,omitempty"`       /*  接入点，多AZ必填: AP1, AP2  */
	Location          *string `json:"location,omitempty"`          /*  接入位置  */
	LineCode          *string `json:"lineCode,omitempty"`          /*  电路代号  */
	Description       *string `json:"description,omitempty"`       /*  描述  */
}

type CdaCdaPhysicalLineAddResponse struct {
	StatusCode  *int32                                    `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaPhysicalLineAddReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaPhysicalLineAddErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	Error       *string                                   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaPhysicalLineAddReturnObjResponse struct {
	Result   *string `json:"result"`   /*  1成功， 0失败  */
	ErrorMsg *string `json:"errorMsg"` /*  成功为空  */
}

type CdaCdaPhysicalLineAddErrorDetailResponse struct{}
