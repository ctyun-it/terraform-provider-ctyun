package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowBandwidthApi
/* 调用此接口可查询共享带宽实例详情。
 */type CtvpcShowBandwidthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowBandwidthApi(client *core.CtyunClient) *CtvpcShowBandwidthApi {
	return &CtvpcShowBandwidthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/bandwidth/describe",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowBandwidthApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowBandwidthRequest) (*CtvpcShowBandwidthResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ProjectID != nil {
		ctReq.AddParam("projectID", *req.ProjectID)
	}
	ctReq.AddParam("bandwidthID", req.BandwidthID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowBandwidthResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowBandwidthRequest struct {
	RegionID    string  /*  共享带宽所在的区域id  */
	ProjectID   *string /*  企业项目 ID，默认为'0'  */
	BandwidthID string  /*  查询的共享带宽id。  */
}

type CtvpcShowBandwidthResponse struct {
	StatusCode  int32                                `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                              `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                              `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                              `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowBandwidthReturnObjResponse `json:"returnObj"`             /*  返回查询的共享带宽详细信息。  */
	Error       *string                              `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowBandwidthReturnObjResponse struct {
	Id        *string                                    `json:"id,omitempty"`        /*  共享带宽id。  */
	Status    *string                                    `json:"status,omitempty"`    /*  ACTIVE  */
	Bandwidth int32                                      `json:"bandwidth"`           /*  共享带宽的带宽峰值， 单位：Mbps。  */
	Name      *string                                    `json:"name,omitempty"`      /*  共享带宽名称。  */
	ExpireAt  *string                                    `json:"expireAt,omitempty"`  /*  过期时间  */
	CreatedAt *string                                    `json:"createdAt,omitempty"` /*  创建时间  */
	Eips      []*CtvpcShowBandwidthReturnObjEipsResponse `json:"eips"`                /*  绑定的弹性 IP 列表，见下表  */
}

type CtvpcShowBandwidthReturnObjEipsResponse struct {
	Ip    *string `json:"ip,omitempty"`    /*  弹性 IP 的 IP  */
	EipID *string `json:"eipID,omitempty"` /*  弹性 IP 的 ID  */
}
