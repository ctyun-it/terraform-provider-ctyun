package pgsql

import (
	"context"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlGetParameterTemplateDetailApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlGetParameterTemplateDetailApi(client *ctyunsdk.CtyunClient) *PgsqlGetParameterTemplateDetailApi {
	return &PgsqlGetParameterTemplateDetailApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/PG/v1/template/template-detail",
		},
	}
}

func (this *PgsqlGetParameterTemplateDetailApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlGetParameterTemplateDetailRequest, header *PgsqlGetParameterTemplateDetailRequestHeader) (GetParameterTemplateDetailResp *PgsqlGetParameterTemplateDetailResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("Project-Id", *header.ProjectID)
	}

	builder.AddHeader("regionId", header.RegionID)
	builder.AddParam("templateId", fmt.Sprintf("%d", req.TemplateId))

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	GetParameterTemplateDetailResp = &PgsqlGetParameterTemplateDetailResponse{}
	err = resp.Parse(GetParameterTemplateDetailResp)
	if err != nil {
		return
	}
	return GetParameterTemplateDetailResp, nil
}

type PgsqlGetParameterTemplateDetailRequest struct {
	TemplateId int64 `json:"templateId"` //
}

type PgsqlGetParameterTemplateDetailRequestHeader struct {
	ProjectID *string `json:"Project-Id"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}

type PgsqlGetParameterTemplateDetailResponse struct {
	StatusCode int32                                              `json:"statusCode"`      // 接口状态码
	Error      *string                                            `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                             `json:"message"`         // 描述信息
	ReturnObj  []PgsqlGetParameterTemplateDetailResponseReturnObj `json:"returnObj"`
}
type PgsqlGetParameterTemplateDetailResponseReturnObj struct {
	ParameterName  string   `json:"paramName"`
	ParameterValue string   `json:"paramValue"`
	Description    string   `json:"description"`
	ValueType      string   `json:"valueType"`
	Restart        int32    `json:"restart"`
	Unit           string   `json:"unit"`
	MinVal         string   `json:"minVal"`
	MaxVal         string   `json:"maxVal"`
	EnumValues     []string `json:"enumValues"`
}

type PgsqlGetParameterTemplateDetailResponseReturnObjDeatil struct {
	ParameterGroupName string `json:"parameterGroupName"`
	ValueType          string `json:"valuetype"`
	DescriptionEn      string `json:"descriptionEn"`
	Restart            string `json:"restart"`
	Description        string `json:"description"`
	ID                 int64  `json:"id"`
	ParameterName      string `json:"parameterName"`
	ParameterValue     string `json:"parameterValue"`
	PermitValue        string `json:"permitValue"`
	UserID             int64  `json:"userId"`
}
