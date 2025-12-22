package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaSwitchListApi
/* 查询已创建的云专线交换机。 */
type CdaCdaSwitchListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaSwitchListApi(client *core.CtyunClient) *CdaCdaSwitchListApi {
	return &CdaCdaSwitchListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/cda/switch/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaSwitchListApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaSwitchListRequest) (*CdaCdaSwitchListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaSwitchListRequest
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
	var resp CdaCdaSwitchListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaSwitchListRequest struct {
	SwitchId     *string `json:"switchId,omitempty"`     /*  交换机ID  */
	ResourcePool *string `json:"resourcePool,omitempty"` /*  资源池ID  */
	Hostname     *string `json:"hostname,omitempty"`     /*  交换机hostname  */
	Name         *string `json:"name,omitempty"`         /*  交换机name  */
	Ip           *string `json:"ip,omitempty"`           /*  交换机IP  */
}

type CdaCdaSwitchListResponse struct {
	StatusCode    *int32                               `json:"statusCode"`    /*  返回状态码(800为成功，900为失败)  */
	Message       *string                              `json:"message"`       /*  失败时的错误描述，一般为英文描述  */
	Description   *string                              `json:"description"`   /*  失败时的错误描述，一般为中文描述  */
	ErrorCode     *string                              `json:"errorCode"`     /*  业务细分码，为product.module.code三段式码.参考[结果码](#通用结果码)  */
	ErrorDetail   *CdaCdaSwitchListErrorDetailResponse `json:"errorDetail"`   /*  错误明细  */
	ErrorMsg      *string                              `json:"errorMsg"`      /*  错误信息  */
	SwitchId      *string                              `json:"switchId"`      /*  交换机ID  */
	Name          *string                              `json:"name"`          /*  交换机名字  */
	Factory       *string                              `json:"factory"`       /*  厂商（RUIJIE、华三）  */
	ResourcePool  *string                              `json:"resourcePool"`  /*  资源池ID  */
	ResourceName  *string                              `json:"resourceName"`  /*  资源池名字  */
	Hostname      *string                              `json:"hostname"`      /*  交换机hostname  */
	Ip            *string                              `json:"ip"`            /*  交换机IP  */
	LoginPort     *string                              `json:"loginPort"`     /*  登录port  */
	VtepIp        *string                              `json:"vtepIp"`        /*  VTEP IP  */
	VtepVlan      *string                              `json:"vtepVlan"`      /*  VTEP VLAN  */
	DeviceModel   *string                              `json:"deviceModel"`   /*  设备型号  */
	AccessPoint   *string                              `json:"accessPoint"`   /*  接入点  */
	As            *int32                               `json:"as"`            /*  as号  */
	HasBleafRoute *bool                                `json:"hasBleafRoute"` /*  标记交换机是否要配置BLEAF路由，默认为false（只有部分锐捷交换机需要配置）  */
	SysMac        *string                              `json:"sysMac"`        /*  交换机mac（多az并且是锐捷交换机则必填）（mac是查交换机配置查出来）  */
	ResourceType  *string                              `json:"resourceType"`  /*  资源池类型  */
	Error         *string                              `json:"error"`         /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaSwitchListErrorDetailResponse struct{}
