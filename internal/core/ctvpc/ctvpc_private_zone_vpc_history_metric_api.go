package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcPrivateZoneVpcHistoryMetricApi
/* vpc维度历史监控
 */type CtvpcPrivateZoneVpcHistoryMetricApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcPrivateZoneVpcHistoryMetricApi(client *core.CtyunClient) *CtvpcPrivateZoneVpcHistoryMetricApi {
	return &CtvpcPrivateZoneVpcHistoryMetricApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/private-zone/vpc-history-metric",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcPrivateZoneVpcHistoryMetricApi) Do(ctx context.Context, credential core.Credential, req *CtvpcPrivateZoneVpcHistoryMetricRequest) (*CtvpcPrivateZoneVpcHistoryMetricResponse, error) {
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
	var resp CtvpcPrivateZoneVpcHistoryMetricResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcPrivateZoneVpcHistoryMetricRequest struct {
	RegionID  string `json:"regionID,omitempty"`  /*  资源池 ID  */
	VpcID     string `json:"vpcID,omitempty"`     /*  vpc id  */
	StartTime string `json:"startTime,omitempty"` /*  开始时间，YYYY-mmm-dd HH:MM:SS  */
	EndTime   string `json:"endTime,omitempty"`   /*  开始时间，YYYY-mmm-dd HH:MM:SS  */
}

type CtvpcPrivateZoneVpcHistoryMetricResponse struct {
	StatusCode  int32                                                `json:"statusCode"`            /*  返回状态码（800 为成功，900 为失败）  */
	Message     *string                                              `json:"message,omitempty"`     /*  statusCode 为 900 时的错误信息; statusCode 为 800 时为 success, 英文  */
	Description *string                                              `json:"description,omitempty"` /*  statusCode 为 900 时的错误信息; statusCode 为 800 时为成功, 中文  */
	ErrorCode   *string                                              `json:"errorCode,omitempty"`   /*  statusCode 为 900 时为业务细分错误码，三段式：product.module.code; statusCode 为 800 时为 SUCCESS  */
	ReturnObj   []*CtvpcPrivateZoneVpcHistoryMetricReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtvpcPrivateZoneVpcHistoryMetricReturnObjResponse struct {
	ItemName   *string                                                        `json:"itemName,omitempty"` /*  监控项名称  */
	ItemDesc   *string                                                        `json:"itemDesc,omitempty"` /*  监控项中文介绍  */
	ItemUnit   *string                                                        `json:"itemUnit,omitempty"` /*  监控项单位  */
	ItemData   []*CtvpcPrivateZoneVpcHistoryMetricReturnObjItemDataResponse   `json:"itemData"`           /*  监控项内容  */
	Dimensions []*CtvpcPrivateZoneVpcHistoryMetricReturnObjDimensionsResponse `json:"dimensions"`         /*  监控项标签  */
}

type CtvpcPrivateZoneVpcHistoryMetricReturnObjItemDataResponse struct {
	Value     float32 `json:"value"`     /*  监控项值，具体请参考对应监控项文档  */
	Timestamp int32   `json:"timestamp"` /*  监控数据 Unix 时间戳  */
}

type CtvpcPrivateZoneVpcHistoryMetricReturnObjDimensionsResponse struct {
	Name  *string `json:"name,omitempty"`  /*  监控项标签键  */
	Value *string `json:"value,omitempty"` /*  监控项标签键对应的值  */
}
