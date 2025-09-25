package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcBuyFlowPackageApi
/* 购买共享流量包。
 */type CtvpcBuyFlowPackageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcBuyFlowPackageApi(client *core.CtyunClient) *CtvpcBuyFlowPackageApi {
	return &CtvpcBuyFlowPackageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/flow_package/buy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcBuyFlowPackageApi) Do(ctx context.Context, credential core.Credential, req *CtvpcBuyFlowPackageRequest) (*CtvpcBuyFlowPackageResponse, error) {
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
	var resp CtvpcBuyFlowPackageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcBuyFlowPackageRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID  */
	CycleType   string `json:"cycleType,omitempty"`   /*  订阅周期，仅支持 MONTH 和 YEAR  */
	CycleCount  int32  `json:"cycleCount"`            /*  订阅周期时长，仅支持购买 1 个月 / 购买 1 年  */
	Count       int32  `json:"count"`                 /*  购买数量，最大支持 20 个，最少 1 个  */
	Spec        int32  `json:"spec"`                  /*  规格说明：当 cycleType = MONTH 时，10-10GB,50-50GB,100-100GB,500-500GB,1024-1TB,5120-5TB,10240-10TB,51200-50TB;**当 cycleType = YEAR 时，120-120GB,512-512GB,8192-8TB,36864-36TB,122880-120TB,614400-600TB,1048576-1PB,2097152-2PB  */
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一，长度 1 - 64  */
}

type CtvpcBuyFlowPackageResponse struct {
	StatusCode  int32                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcBuyFlowPackageReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                               `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcBuyFlowPackageReturnObjResponse struct{}
