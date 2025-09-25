package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlBoundEipListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlBoundEipListApi(client *ctyunsdk.CtyunClient) *PgsqlBoundEipListApi {
	return &PgsqlBoundEipListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/Pgsql-dcp/v2/openapi/dcp-order-info/eips",
		},
	}
}

func (this *PgsqlBoundEipListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlBoundEipListRequest, header *PgsqlBoundEipListRequestHeader) (bindResponse *PgsqlBoundEipListResponse, err error) {
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
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	bindResponse = &PgsqlBoundEipListResponse{}
	err = resp.Parse(bindResponse)
	if err != nil {
		return
	}
	return bindResponse, nil
}

type PgsqlBoundEipListRequest struct {
	RegionID string  `json:"regionId"`
	EipID    *string `json:"eipId"`
	InstID   *string `json:"instId"`
}

type PgsqlBoundEipListRequestHeader struct {
	ProjectID *string `json:"project_id"`
}

type PgsqlBoundEipListResponseReturnObj struct {
	Data []PgsqlBoundEipListResponseReturnObjData `json:"data"`
}

type PgsqlBoundEipListResponse struct {
	StatusCode int32                               `json:"statusCode"` // 接口状态码
	Error      string                              `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    *string                             `json:"message"`    // 描述信息
	ReturnObj  *PgsqlBoundEipListResponseReturnObj `json:"returnObj"`
}

type PgsqlBoundEipListResponseReturnObjData struct {
	EipID      string `json:"eip_id"`
	Eip        string `json:"eip"`
	BindStatus int32  `json:"bindStatus"`
	Status     string `json:"status"`
	BandWidth  int32  `json:"bandWidth"`
}
