package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaLinkProbeAddApi
/* 用户为专线网关创建健康检查 */
type CdaCdaLinkProbeAddApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaLinkProbeAddApi(client *core.CtyunClient) *CdaCdaLinkProbeAddApi {
	return &CdaCdaLinkProbeAddApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cda/link-probe/add",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaLinkProbeAddApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaLinkProbeAddRequest) (*CdaCdaLinkProbeAddResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaLinkProbeAddRequest
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
	var resp CdaCdaLinkProbeAddResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaLinkProbeAddRequest struct {
	IpInfoList                []*CdaCdaLinkProbeAddIpInfoListRequest `json:"ipInfoList,omitempty"`                /*  源ip、目的ip  */
	GatewayName               string                                 `json:"gatewayName"`                         /*  专线网关名字  */
	ServiceID                 string                                 `json:"serviceID"`                           /*  健康检查关联的业务id：物理lineId或者vpcID, 关于serviceID:如果是Ping客户侧那么传入物理lineId，如果是ping云侧那传入vpcID  */
	Interval                  *string                                `json:"interval,omitempty"`                  /*  发包时间间隔，单位毫秒，不带使用厂商默认时间间隔商<br/>   （h3c 200ms rg 100ms）  */
	NumberPackagesPerShipment *string                                `json:"numberPackagesPerShipment,omitempty"` /*  每次发包的数量，默认5个  */
}

type CdaCdaLinkProbeAddIpInfoListRequest struct {
	SrcIp string `json:"srcIp"` /*  源ip，无掩码。不填不带源ping  */
	DstIp string `json:"dstIp"` /*  目的ip，无掩码。  */
}

type CdaCdaLinkProbeAddResponse struct {
	StatusCode  *int32                                 `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaLinkProbeAddReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaLinkProbeAddErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
}

type CdaCdaLinkProbeAddReturnObjResponse struct {
	Result    *string `json:"result"`    /*  1成功， 0失败  */
	Data      *string `json:"data"`      /*  成功为空  */
	ErrorCode *string `json:"errorCode"` /*  错误代码，成功为空  */
	ErrorMsg  *string `json:"errorMsg"`  /*  成功为空  */
	TraceId   *string `json:"traceId"`   /*  日志跟踪ID  */
}

type CdaCdaLinkProbeAddErrorDetailResponse struct{}
