package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetAvailabilityZone struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetAvailabilityZone(client *ctyunsdk.CtyunClient) *TeledbGetAvailabilityZone {
	return &TeledbGetAvailabilityZone{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/teledb-acceptor/v2/openapi/accept-order-info/availabilityZone",
		},
	}
}

type TeledbGetAvailabilityZoneRequest struct {
	RegionId string `json:"regionId"`
}

type TeledbGetAvailabilityZoneRequestHeader struct {
	ProjectID *string `json:"project-id"`
}

type TeledbGetAvailabilityZoneResponseReturnObjData struct {
	AvailabilityZoneId   string `json:"availabilityZoneId"`
	AvailabilityZoneName string `json:"availabilityZoneName"`
	DisplayName          string `json:"displayName"`
}

type TeledbGetAvailabilityZoneResponseReturnObj struct {
	Data []TeledbGetAvailabilityZoneResponseReturnObjData `json:"data"`
}

type TeledbGetAvailabilityZoneResponse struct {
	StatusCode int32                                      `json:"statusCode"` // 接口状态码
	Error      string                                     `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string                                     `json:"message"`    // 描述信息
	ReturnObj  TeledbGetAvailabilityZoneResponseReturnObj `json:"returnObj"`
}

func (this *TeledbGetAvailabilityZone) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetAvailabilityZoneRequest, header *TeledbGetAvailabilityZoneRequestHeader) (bindResponse *TeledbGetAvailabilityZoneResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if req.RegionId == "" {
		err = errors.New("region id 不能为空")
	}
	builder.AddParam("regionId", req.RegionId)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	bindResponse = &TeledbGetAvailabilityZoneResponse{}
	err = resp.Parse(bindResponse)
	if err != nil {
		return
	}
	return bindResponse, nil
}
