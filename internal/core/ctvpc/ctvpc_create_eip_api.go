package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateEipApi
/* 调用此接口可创建弹性公网IP（Elastic IP Address，简称EIP）。
 */type CtvpcCreateEipApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateEipApi(client *core.CtyunClient) *CtvpcCreateEipApi {
	return &CtvpcCreateEipApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/eip/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateEipApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateEipRequest) (*CtvpcCreateEipResponse, error) {
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
	var resp CtvpcCreateEipResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateEipRequest struct {
	ClientToken       string  `json:"clientToken,omitempty"`       /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID          string  `json:"regionID,omitempty"`          /*  资源池 ID  */
	ProjectID         *string `json:"projectID,omitempty"`         /*  不填默认为默认企业项目，如果需要指定企业项目，则需要填写  */
	CycleType         string  `json:"cycleType,omitempty"`         /*  订购类型：month（包月） / year（包年） / on_demand（按需）  */
	CycleCount        int32   `json:"cycleCount,omitempty"`        /*  订购时长, 当 cycleType = month, 支持续订 1 - 11 个月; 当 cycleType = year, 支持续订 1 - 3 年, 当 cycleType = on_demand 时，可以不传  */
	Name              string  `json:"name,omitempty"`              /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Bandwidth         int32   `json:"bandwidth"`                   /*  弹性 IP 的带宽峰值，默认为 1 Mbps  */
	BandwidthID       *string `json:"bandwidthID,omitempty"`       /*  当 cycleType 为 on_demand 时，可以使用 bandwidthID，将弹性 IP 加入到共享带宽中  */
	DemandBillingType *string `json:"demandBillingType,omitempty"` /*  按需计费类型，当 cycleType 为 on_demand 时生效，支持 bandwidth（按带宽）/ upflowc（按流量）  */
	PayVoucherPrice   *string `json:"payVoucherPrice,omitempty"`   /*  代金券金额，支持到小数点后两位，仅包周期支持代金券  */
	LineType          *string `json:"lineType,omitempty"`          /*  线路类型，默认为163，支持163 / bgp / chinamobile / chinaunicom  */
	SegmentID         *string `json:"segmentID,omitempty"`         /*  专属地址池 segment id，先通过接口 /v4/eip/own-segments 获取  */
	ExclusiveName     *string `json:"exclusiveName,omitempty"`     /*  专属地址池 exclusiveName, 先通过接口 /v4/eip/own-segments 获取  */
}

type CtvpcCreateEipResponse struct {
	StatusCode  int32                            `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                          `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                          `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                          `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCreateEipReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       *string                          `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcCreateEipReturnObjResponse struct {
	MasterOrderID        *string `json:"masterOrderID,omitempty"`        /*  订单id。  */
	MasterOrderNO        *string `json:"masterOrderNO,omitempty"`        /*  订单编号, 可以为 null。  */
	MasterResourceStatus *string `json:"masterResourceStatus,omitempty"` /*  资源状态: started（启用） / renewed（续订） / refunded（退订） / destroyed（销毁） / failed（失败） / starting（正在启用） / changed（变配）/ expired（过期）/ unknown（未知）  */
	MasterResourceID     *string `json:"masterResourceID,omitempty"`     /*  可以为 null。  */
	RegionID             *string `json:"regionID,omitempty"`             /*  可用区id。  */
	EipID                *string `json:"eipID,omitempty"`                /*  弹性 IP id  */
}
