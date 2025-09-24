package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcBindHavipApi
/* 将HaVip绑定到ECS实例上，由于绑定是异步操作，在第一次请求后，并不会立即返回绑定结果，调用者在获取到绑定状态为 in_progress 时，继续使用相同参数进行请求，获取最新的绑定结果，直到最后的绑定状态为 done 即可停止请求。
 */type CtvpcBindHavipApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcBindHavipApi(client *core.CtyunClient) *CtvpcBindHavipApi {
	return &CtvpcBindHavipApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/havip/bind",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcBindHavipApi) Do(ctx context.Context, credential core.Credential, req *CtvpcBindHavipRequest) (*CtvpcBindHavipResponse, error) {
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
	var resp CtvpcBindHavipResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcBindHavipRequest struct {
	ClientToken        string  `json:"clientToken,omitempty"`        /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID           string  `json:"regionID,omitempty"`           /*  资源池ID  */
	ResourceType       string  `json:"resourceType,omitempty"`       /*  绑定的实例类型，VM 表示虚拟机ECS, PM 表示裸金属, NETWORK 表示弹性 IP  */
	HaVipID            string  `json:"haVipID,omitempty"`            /*  高可用虚IP的ID  */
	NetworkInterfaceID *string `json:"networkInterfaceID,omitempty"` /*  虚拟网卡ID, 该网卡属于instanceID, 当 resourceType 为 VM / PM 时，必填  */
	InstanceID         *string `json:"instanceID,omitempty"`         /*  ECS示例ID，当 resourceType 为 VM / PM 时，必填  */
	FloatingID         *string `json:"floatingID,omitempty"`         /*  弹性IP ID，当 resourceType 为 NETWORK 时，必填  */
}

type CtvpcBindHavipResponse struct {
	StatusCode  int32                              `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                            `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                            `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                            `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcBindHavipReturnObjResponse `json:"returnObj"`             /*  绑定状态  */
	Error       *string                            `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcBindHavipReturnObjResponse struct {
	Status  *string `json:"status,omitempty"`  /*  绑定状态，取值 in_progress / done  */
	Message *string `json:"message,omitempty"` /*  绑定状态提示信息  */
}
