package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlGetCharacterSetApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlGetCharacterSetApi(client *ctyunsdk.CtyunClient) *PgsqlGetCharacterSetApi {
	return &PgsqlGetCharacterSetApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/PG/v1/database/character-set-name",
		},
	}
}

func (this *PgsqlGetCharacterSetApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlGetCharacterSetRequest, header *PgsqlGetCharacterSetRequestHeader) (GetCharacterSetResp *PgsqlGetCharacterSetResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}

	builder.AddHeader("regionId", header.RegionID)

	if req.Engine == "" {
		err = errors.New("engine is required")
		return
	}
	builder.AddParam("engine", req.Engine)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	GetCharacterSetResp = &PgsqlGetCharacterSetResponse{}
	err = resp.Parse(GetCharacterSetResp)
	if err != nil {
		return
	}
	return GetCharacterSetResp, nil
}

type PgsqlGetCharacterSetRequest struct {
	Engine string `json:"engine"` // 外部实例ID，必填 postgreSQL
}

type PgsqlGetCharacterSetRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}

type PgsqlGetCharacterSetResponse struct {
	StatusCode int32                                  `json:"statusCode"`      // 接口状态码
	Error      *string                                `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                 `json:"message"`         // 描述信息
	ReturnObj  *PgsqlGetCharacterSetResponseReturnObj `json:"returnObj"`
}

type PgsqlGetCharacterSetResponseReturnObj struct {
	Engine                string   `json:"engine"`
	CharacterSetNameItems []string `json:"characterSetNameItems"`
}
