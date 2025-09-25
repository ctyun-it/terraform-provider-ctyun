package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateEipWithIpAddressApi
/* 调用此接口可创建指定 IP 地址的弹性公网IP（Elastic IP Address，简称EIP）。
 */type CtvpcCreateEipWithIpAddressApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateEipWithIpAddressApi(client *core.CtyunClient) *CtvpcCreateEipWithIpAddressApi {
	return &CtvpcCreateEipWithIpAddressApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/eip/create-with-ipaddress",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateEipWithIpAddressApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateEipWithIpAddressRequest) (*CtvpcCreateEipWithIpAddressResponse, error) {
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
	var resp CtvpcCreateEipWithIpAddressResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateEipWithIpAddressRequest struct {
	ClientToken       string  `json:"clientToken,omitempty"`       /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID          string  `json:"regionID,omitempty"`          /*  资源池 ID  */
	ProjectID         *string `json:"projectID,omitempty"`         /*  企业项目 ID，默认为'0'    */
	CycleType         string  `json:"cycleType,omitempty"`         /*  订购类型：month（包月） / year（包年） / on_demand（按需）  */
	CycleCount        int32   `json:"cycleCount"`                  /*  订购时长, cycleType 是 on_demand 为可选，当 cycleType = month, 支持续订 1 - 11 个月; 当 cycleType = year, 支持续订 1 - 3 年  */
	Name              string  `json:"name,omitempty"`              /*  弹性 IP 名称  */
	Bandwidth         int32   `json:"bandwidth"`                   /*  弹性 IP 的带宽峰值，默认为 1 Mbps  */
	IpAddress         string  `json:"ipAddress,omitempty"`         /*  合法的公网 IP  */
	BandwidthID       *string `json:"bandwidthID,omitempty"`       /*  当 cycleType 为 on_demand 时，可以使用 bandwidthID，将弹性 IP 加入到共享带宽中  */
	DemandBillingType *string `json:"demandBillingType,omitempty"` /*  按需计费类型，当 cycleType 为 on_demand 时生效，支持 bandwidth（按带宽）/ upflowc（按流量）  */
	LineType          *string `json:"lineType,omitempty"`          /*  线路类型，默认为163，支持163 / bgp / chinamobile / chinaunicom  */
	PayVoucherPrice   *string `json:"payVoucherPrice,omitempty"`   /*  代金券金额，支持到小数点后两位，仅包周期支持代金券  */
	SegmentID         *string `json:"segmentID,omitempty"`         /*  专属地址池 segment id，先通过接口 /v4/eip/own-segments 获取  */
	ExclusiveName     *string `json:"exclusiveName,omitempty"`     /*  专属地址池 exclusiveName，先通过接口 /v4/eip/own-segments 获取  */
}

type CtvpcCreateEipWithIpAddressResponse struct {
	StatusCode  int32                                           `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                         `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                         `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                         `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcCreateEipWithIpAddressReturnObjResponse `json:"returnObj"`             /*  object  */
	Error       *string                                         `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcCreateEipWithIpAddressReturnObjResponse struct {
	MasterOrderID        *string `json:"masterOrderID,omitempty"`        /*  订单id。  */
	MasterOrderNO        *string `json:"masterOrderNO,omitempty"`        /*  订单编号, 可以为 null。  */
	MasterResourceStatus *string `json:"masterResourceStatus,omitempty"` /*  资源状态。  */
	MasterResourceID     *string `json:"masterResourceID,omitempty"`     /*  可以为 null。  */
	RegionID             *string `json:"regionID,omitempty"`             /*  可用区id。  */
	EipID                *string `json:"eipID,omitempty"`                /*  当 masterResourceStatus 不为 started 时，该值可能为  */
}
