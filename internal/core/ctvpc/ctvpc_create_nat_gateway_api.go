package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateNatGatewayApi
/* 创建 NAT 网关
 */type CtvpcCreateNatGatewayApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateNatGatewayApi(client *core.CtyunClient) *CtvpcCreateNatGatewayApi {
	return &CtvpcCreateNatGatewayApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/create-nat-gateway",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateNatGatewayApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateNatGatewayRequest) (*CtvpcCreateNatGatewayResponse, error) {
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
	var resp CtvpcCreateNatGatewayResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateNatGatewayRequest struct {
	RegionID        string  `json:"regionID,omitempty"`        /*  区域id  */
	VpcID           string  `json:"vpcID,omitempty"`           /*  需要创建 NAT 网关的 VPC 的 ID  */
	Spec            int32   `json:"spec"`                      /*  规格 1~4, 1表示小型, 2表示中型, 3表示大型, 4表示超大型  */
	Name            string  `json:"name,omitempty"`            /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Description     *string `json:"description,omitempty"`     /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:,'{},.,/;'[]·~！@#￥%……&*（） ——-+={}  */
	ClientToken     string  `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	CycleType       string  `json:"cycleType,omitempty"`       /*  订购类型：month（包月） / year（包年）/ on_demand（按需）  */
	CycleCount      *int32  `json:"cycleCount"`                /*  订购时长, 当 cycleType = month, 支持续订 1 - 11 个月; 当 cycleType = year, 支持续订 1 - 3 年  */
	AzName          string  `json:"azName,omitempty"`          /*  可用区名称  */
	PayVoucherPrice *string `json:"payVoucherPrice,omitempty"` /*  代金券金额，支持到小数点后两位  */
	ProjectID       *string `json:"projectID,omitempty"`       /*  企业项目，不传默认为 0  */
}

type CtvpcCreateNatGatewayResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCreateNatGatewayReturnObjResponse `json:"returnObj"`             /*  object  */
	Error       *string                                 `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcCreateNatGatewayReturnObjResponse struct {
	MasterOrderID        *string `json:"masterOrderID,omitempty"`        /*  订单id。  */
	MasterOrderNO        *string `json:"masterOrderNO,omitempty"`        /*  订单编号, 可以为 null。  */
	MasterResourceStatus *string `json:"masterResourceStatus,omitempty"` /*  资源状态: started（启用） / renewed（续订） / refunded（退订） / destroyed（销毁） / failed（失败） / starting（正在启用） / changed（变配）/ expired（过期）/ unknown（未知）  */
	MasterResourceID     *string `json:"masterResourceID,omitempty"`     /*  可以为 null。  */
	RegionID             *string `json:"regionID,omitempty"`             /*  可用区id。  */
	NatGatewayID         *string `json:"natGatewayID,omitempty"`         /*  nat 网关 ID，当 masterResourceStatus 不为 started，该字段为空字符串  */
}
