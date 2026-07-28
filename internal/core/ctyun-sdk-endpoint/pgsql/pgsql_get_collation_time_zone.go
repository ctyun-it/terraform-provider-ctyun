package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlGetCollationTimeZoneApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlGetCollationTimeZoneApi(client *ctyunsdk.CtyunClient) *PgsqlGetCollationTimeZoneApi {
	return &PgsqlGetCollationTimeZoneApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/PG/v1/database/collation-time-zone",
		},
	}
}

func (this *PgsqlGetCollationTimeZoneApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlGetCollationTimeZoneRequest, header *PgsqlGetCollationTimeZoneRequestHeader) (GetCollationTimeZoneResp *PgsqlGetCollationTimeZoneResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}

	builder.AddHeader("regionId", header.RegionID)

	if req.ProdInstId == "" {
		err = errors.New("missing required parameter: prodInstId")
		return
	}
	builder.AddParam("prodInstId", req.ProdInstId)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	GetCollationTimeZoneResp = &PgsqlGetCollationTimeZoneResponse{}
	err = resp.Parse(GetCollationTimeZoneResp)
	if err != nil {
		return
	}
	return GetCollationTimeZoneResp, nil
}

type PgsqlGetCollationTimeZoneRequest struct {
	ProdInstId string `json:"prodInstId"` // 外部实例ID，必填
}

type PgsqlGetCollationTimeZoneRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}

type PgsqlGetCollationTimeZoneResponse struct {
	StatusCode int32                                       `json:"statusCode"`      // 接口状态码
	Error      *string                                     `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                      `json:"message"`         // 描述信息
	ReturnObj  *PgsqlGetCollationTimeZoneResponseReturnObj `json:"returnObj"`
}

type PgsqlGetCollationTimeZoneResponseReturnObjCollation struct {
	CollName     string `json:"collname"`
	CollenCoding string `json:"collencoding"`
	CollCollate  string `json:"collcollate"`
	CollCtype    string `json:"collctype"`
}

type PgsqlGetCollationTimeZoneResponseReturnObj struct {
	StandardTimeOffset string                                                `json:"standardTimeOffset"`
	TimeZone           string                                                `json:"timeZone"`
	Collations         []PgsqlGetCollationTimeZoneResponseReturnObjCollation `json:"collations"`
	Description        string                                                `json:"description"`
}
