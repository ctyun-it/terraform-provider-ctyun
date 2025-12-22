package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaLinkProbeDeleteApi
/* 删除指定源目的ip的健康检查数据，支持批量删除 */
type CdaCdaLinkProbeDeleteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaLinkProbeDeleteApi(client *core.CtyunClient) *CdaCdaLinkProbeDeleteApi {
	return &CdaCdaLinkProbeDeleteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cda/link-probe/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaLinkProbeDeleteApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaLinkProbeDeleteRequest) (*CdaCdaLinkProbeDeleteResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaLinkProbeDeleteRequest
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
	var resp CdaCdaLinkProbeDeleteResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaLinkProbeDeleteRequest struct {
	IpInfoList []*CdaCdaLinkProbeDeleteIpInfoListRequest `json:"ipInfoList,omitempty"` /*  源ip、目的ip  */
}

type CdaCdaLinkProbeDeleteIpInfoListRequest struct {
	SrcIp string `json:"srcIp"` /*  源ip，无掩码。不填不带源ping  */
	DstIp string `json:"dstIp"` /*  目的ip，无掩码。  */
}

type CdaCdaLinkProbeDeleteResponse struct {
	StatusCode  *int32                                    `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CdaCdaLinkProbeDeleteReturnObjResponse   `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CdaCdaLinkProbeDeleteErrorDetailResponse `json:"errorDetail"` /*  错误明细  */
}

type CdaCdaLinkProbeDeleteReturnObjResponse struct {
	Result    *string `json:"result"`    /*  1成功， 0失败  */
	Data      *string `json:"data"`      /*  成功为null  */
	ErrorCode *string `json:"errorCode"` /*  错误代码，成功为null  */
	ErrorMsg  *string `json:"errorMsg"`  /*  成功为null  */
	TraceId   *string `json:"traceId"`   /*  日志跟踪ID  */
}

type CdaCdaLinkProbeDeleteErrorDetailResponse struct{}
