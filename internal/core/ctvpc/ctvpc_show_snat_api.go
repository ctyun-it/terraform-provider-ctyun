package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowSnatApi
/* 获取SNAT详情
 */type CtvpcShowSnatApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowSnatApi(client *core.CtyunClient) *CtvpcShowSnatApi {
	return &CtvpcShowSnatApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/show-snat",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowSnatApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowSnatRequest) (*CtvpcShowSnatResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sNatID", req.SNatID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowSnatResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowSnatRequest struct {
	RegionID string /*  区域id  */
	SNatID   string /*  snat id  */
}

type CtvpcShowSnatResponse struct {
	StatusCode  int32                           `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                         `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                         `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                         `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowSnatReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtvpcShowSnatReturnObjResponse struct {
	SNatID       *string                               `json:"sNatID,omitempty"`       /*  snat id  */
	Description  *string                               `json:"description,omitempty"`  /*  描述信息  */
	SubnetCidr   *string                               `json:"subnetCidr,omitempty"`   /*  要查询的NAT网关所属VPC子网的cidr  */
	SubnetType   int32                                 `json:"subnetType"`             /*  子网类型：1-有vpcID的子网，0-自定义  */
	CreationTime *string                               `json:"creationTime,omitempty"` /*  创建时间  */
	Eips         []*CtvpcShowSnatReturnObjEipsResponse `json:"eips"`                   /*  绑定的 eip 信息  */
	SubnetID     *string                               `json:"subnetID,omitempty"`     /*  子网 ID  */
	NatGatewayID *string                               `json:"natGatewayID,omitempty"` /*  nat 网关 ID  */
}

type CtvpcShowSnatReturnObjEipsResponse struct {
	EipID     *string `json:"eipID,omitempty"`     /*  弹性 IP id  */
	IpAddress *string `json:"ipAddress,omitempty"` /*  弹性 IP 地址  */
}
