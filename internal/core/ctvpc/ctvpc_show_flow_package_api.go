package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowFlowPackageApi
/* 查询共享流量包详情。 */
type CtvpcShowFlowPackageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowFlowPackageApi(client *core.CtyunClient) *CtvpcShowFlowPackageApi {
	return &CtvpcShowFlowPackageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/flow_package/show",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowFlowPackageApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowFlowPackageRequest) (*CtvpcShowFlowPackageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sdpID", req.SdpID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowFlowPackageResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtvpcShowFlowPackageRequest struct {
	RegionID string `json:"regionID"` /*  资源池 ID  */
	SdpID    string `json:"sdpID"`    /*  记录标识  */
}

type CtvpcShowFlowPackageResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                `json:"message"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                `json:"description"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                `json:"errorCode"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowFlowPackageReturnObjResponse `json:"returnObj"`   /*  返回购买的共享流量包列表  */
}

type CtvpcShowFlowPackageReturnObjResponse struct {
	Id                     *string `json:"id"`                     /*  唯一标识  */
	Status                 *string `json:"status"`                 /*  购买共享流量包的状态，可能的取值：初始、有效、退订、过期、销毁  */
	CycleType              *string `json:"cycleType"`              /*  支持的取值：包小时、包天、包周、包月、包年  */
	EffectiveTime          *string `json:"effectiveTime"`          /*  生效时间  */
	ExpireTime             *string `json:"expireTime"`             /*  过期时间  */
	PackageName            *string `json:"packageName"`            /*  套餐名  */
	TotalVolumn            *string `json:"totalVolumn"`            /*  总流量  */
	LeftVolumn             *string `json:"leftVolumn"`             /*  剩余流量  */
	RegionID               *string `json:"regionID"`               /*  regionID  */
	MasterResourceBundleId *string `json:"masterResourceBundleId"` /*  资源包 ID  */
}
