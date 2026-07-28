package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaVpcUpdateApi
/* 给已创建的云专线网关修改VPC */
type CdaCdaVpcUpdateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaVpcUpdateApi(client *core.CtyunClient) *CdaCdaVpcUpdateApi {
	return &CdaCdaVpcUpdateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cda/vpc/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaVpcUpdateApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaVpcUpdateRequest) (*CdaCdaVpcUpdateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaVpcUpdateRequest
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
	var resp CdaCdaVpcUpdateResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaVpcUpdateRequest struct {
	RegionID      string    `json:"regionID"`                /*  资源池ID  */
	GatewayName   string    `json:"gatewayName"`             /*  专线网关名称（唯一）（只能是字母和数字）  */
	VpcID         string    `json:"vpcID"`                   /*  VPC ID  */
	VpcSubnet     []*string `json:"vpcSubnet,omitempty"`     /*  vpc ipv4子网列表(全量输入，ipVersion为IPV4和DUALSTACK时必填)  */
	VpcSubnetIPv6 []*string `json:"vpcSubnetIPv6,omitempty"` /*  vpc ipv6子网列表(全量输入，ipVersion为IPV6和DUALSTACK时必填)  */
	IpVersion     string    `json:"ipVersion"`               /*  本参数表示包周期类型。<br>取值范围：<br>IPV4<br>IPV6<br>DUALSTACK  */
	Bandwidth     int32     `json:"bandwidth"`               /*  带宽(M)  */
}

type CdaCdaVpcUpdateResponse struct {
	StatusCode  *int32                              `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                             `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                             `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaVpcUpdateReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                             `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaVpcUpdateErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
	Error       *string                             `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type CdaCdaVpcUpdateReturnObjResponse struct {
	Result   *string `json:"result"`   /*  1成功， 0失败  */
	Data     *string `json:"data"`     /*  成功为空  */
	ErrorMsg *string `json:"errorMsg"` /*  成功为空  */
}

type CdaCdaVpcUpdateErrorDetailResponse struct{}
