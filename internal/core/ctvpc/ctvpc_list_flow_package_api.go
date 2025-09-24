package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcListFlowPackageApi
/* 查询购买的共享流量包列表。
 */type CtvpcListFlowPackageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListFlowPackageApi(client *core.CtyunClient) *CtvpcListFlowPackageApi {
	return &CtvpcListFlowPackageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/flow_package/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListFlowPackageApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListFlowPackageRequest) (*CtvpcListFlowPackageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcListFlowPackageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListFlowPackageRequest struct {
	RegionID string /*  资源池 ID  */
}

type CtvpcListFlowPackageResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcListFlowPackageReturnObjResponse `json:"returnObj"`             /*  返回购买的共享流量包列表  */
	Error       *string                                  `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListFlowPackageReturnObjResponse struct {
	Id                     *string `json:"id,omitempty"`                     /*  唯一标识  */
	StatusText             *string `json:"statusText,omitempty"`             /*  购买共享流量包的状态，可能的取值：初始、有效、退订、过期、销毁  */
	Status                 int32   `json:"status"`                           /*  对应 statusText 的取值为：0、1、5、6、7  */
	CycleType              *string `json:"cycleType,omitempty"`              /*  支持的取值：包小时、包天、包周、包月、包年  */
	EffectiveTime          *string `json:"effectiveTime,omitempty"`          /*  生效时间  */
	ExpireTime             *string `json:"expireTime,omitempty"`             /*  过期时间  */
	PackageName            *string `json:"packageName,omitempty"`            /*  套餐名  */
	TotalVolumn            *string `json:"totalVolumn,omitempty"`            /*  总流量  */
	LeftVolumn             *string `json:"leftVolumn,omitempty"`             /*  剩余流量  */
	MasterResourceBundleId *string `json:"masterResourceBundleId,omitempty"` /*  资源包 ID  */
}
