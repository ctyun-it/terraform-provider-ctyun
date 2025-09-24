package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVpcCreateSubnetApi
/* 创建子网。
 */type CtvpcVpcCreateSubnetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVpcCreateSubnetApi(client *core.CtyunClient) *CtvpcVpcCreateSubnetApi {
	return &CtvpcVpcCreateSubnetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/create-subnet",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVpcCreateSubnetApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVpcCreateSubnetRequest) (*CtvpcVpcCreateSubnetResponse, error) {
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
	var resp CtvpcVpcCreateSubnetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVpcCreateSubnetRequest struct {
	ClientToken     string    `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID        string    `json:"regionID,omitempty"`        /*  资源池 ID  */
	VpcID           string    `json:"vpcID,omitempty"`           /*  虚拟私有云 ID  */
	Name            string    `json:"name,omitempty"`            /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Description     *string   `json:"description,omitempty"`     /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[\]·~！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
	CIDR            string    `json:"CIDR,omitempty"`            /*  子网网段  */
	EnableIpv6      *bool     `json:"enableIpv6"`                /*  是否开启 IPv6 网段。取值：false（默认值）:不开启，true: 开启  */
	DnsList         []*string `json:"dnsList"`                   /*  子网 dns 列表, 最多同时支持 4 个 dns 地址  */
	SubnetGatewayIP *string   `json:"subnetGatewayIP,omitempty"` /*  子网网关 IP  */
	SubnetType      *string   `json:"subnetType,omitempty"`      /*  子网类型：common（普通子网）/ cbm（裸金属子网），默认为普通子网  */
	DhcpIP          *string   `json:"dhcpIP,omitempty"`          /*  dhcpIP,和网关IP不能相同  */
}

type CtvpcVpcCreateSubnetResponse struct {
	StatusCode  int32                                  `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcVpcCreateSubnetReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVpcCreateSubnetReturnObjResponse struct {
	SubnetID *string `json:"subnetID,omitempty"` /*  subnet 示例 ID  */
}
