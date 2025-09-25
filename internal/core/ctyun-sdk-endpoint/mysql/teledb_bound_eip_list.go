package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbBoundEipListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbBoundEipListApi(client *ctyunsdk.CtyunClient) *TeledbBoundEipListApi {
	return &TeledbBoundEipListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/teledb-dcp/v2/openapi/dcp-order-info/eips",
		},
	}
}

func (this *TeledbBoundEipListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbBoundEipListRequest, header *TeledbBoundEipListRequestHeader) (bindResponse *TeledbBoundEipListResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if req.RegionID == "" {
		err = errors.New("region id is required")
		return
	}
	builder.AddParam("regionId", req.RegionID)
	if req.InstID != nil {
		builder.AddParam("instId", *req.InstID)
	}
	if req.EipID != nil {
		builder.AddParam("eipId", *req.EipID)
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	bindResponse = &TeledbBoundEipListResponse{}
	err = resp.Parse(bindResponse)
	if err != nil {
		return
	}
	return bindResponse, nil
}

type TeledbBoundEipListRequest struct {
	RegionID string  `json:"regionId"`
	EipID    *string `json:"eipId"`
	InstID   *string `json:"instId"`
}

type TeledbBoundEipListRequestHeader struct {
	ProjectID *string `json:"project_id"`
}

type TeledbBoundEipListResponseReturnObj struct {
	Data []TeledbBoundEipListResponseReturnObjData `json:"data"`
}

type TeledbBoundEipListResponse struct {
	StatusCode int32                                `json:"statusCode"` // 接口状态码
	Error      string                               `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string                               `json:"message"`    // 描述信息
	ReturnObj  *TeledbBoundEipListResponseReturnObj `json:"returnObj"`
}

type TeledbBoundEipListResponseReturnObjData struct {
	EipID      string `json:"eip_id"`
	Eip        string `json:"eip"`
	BindStatus int32  `json:"bindStatus"`
	Status     string `json:"status"`
	BandWidth  int32  `json:"bandWidth"`
}
