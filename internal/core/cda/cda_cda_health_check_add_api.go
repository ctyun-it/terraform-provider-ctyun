package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaHealthCheckAddApi
/* 用户为专线网关创建健康检查 */
type CdaCdaHealthCheckAddApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaHealthCheckAddApi(client *core.CtyunClient) *CdaCdaHealthCheckAddApi {
	return &CdaCdaHealthCheckAddApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cda/health-check/add",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaHealthCheckAddApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaHealthCheckAddRequest) (*CdaCdaHealthCheckAddResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaHealthCheckAddRequest
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
	var resp CdaCdaHealthCheckAddResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaHealthCheckAddRequest struct {
	RegionID           string `json:"regionID"`           /*  资源池ID  */
	GatewayName        string `json:"gatewayName"`        /*  专线网关名字  */
	VpcID              string `json:"vpcID"`              /*  VPC ID  */
	VpcName            string `json:"vpcName"`            /*  VPC名字  */
	VpcSubnet          string `json:"vpcSubnet"`          /*  VPC子网  */
	VpcSubnetID        string `json:"vpcSubnetID"`        /*  VPC 子网ID  */
	SrcIP              string `json:"srcIP"`              /*  源IP，VPC子网范围内的任意空闲IP地址  */
	DstIP              string `json:"dstIP"`              /*  目的IP，远端互联IP地址或客户侧子网范围内的业务IP地址  */
	Interval           int32  `json:"interval"`           /*  发包间隔（s）：2、3、4、5  */
	Ntimest            int32  `json:"ntimest"`            /*  发包数量（个）：5、6、7、8、9、10  */
	AutoRouteSwitching bool   `json:"autoRouteSwitching"` /*  是否自动切换路由  */
}

type CdaCdaHealthCheckAddResponse struct {
	StatusCode  *int32                                   `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                  `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                  `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaHealthCheckAddReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                  `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaHealthCheckAddErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
}

type CdaCdaHealthCheckAddReturnObjResponse struct {
	Result    *string `json:"result"`    /*  1成功， 0失败  */
	Data      *string `json:"data"`      /*  成功为空  */
	ErrorCode *string `json:"errorCode"` /*  错误代码，成功为空  */
	ErrorMsg  *string `json:"errorMsg"`  /*  成功为空  */
	TraceId   *string `json:"traceId"`   /*  日志跟踪ID  */
}

type CdaCdaHealthCheckAddErrorDetailResponse struct{}
