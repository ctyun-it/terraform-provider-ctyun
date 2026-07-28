package pgsql

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlUpdateParameterTemplateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlUpdateParameterTemplateApi(client *ctyunsdk.CtyunClient) *PgsqlUpdateParameterTemplateApi {
	return &PgsqlUpdateParameterTemplateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/PG/v1/template/modifyParams",
		},
	}
}

func (this *PgsqlUpdateParameterTemplateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlUpdateParameterTemplateRequest, header *PgsqlUpdateParameterTemplateRequestHeader) (UpdateParameterTemplateResp *PgsqlUpdateParameterTemplateResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}

	builder.AddHeader("regionId", header.RegionID)

	//builder.AddParam("id", fmt.Sprintf("%d", req.ID))
	//valueStr, err := valueToString(req.Value)
	//if err != nil {
	//	return nil, err
	//}
	//builder.AddParam("value", valueStr)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	UpdateParameterTemplateResp = &PgsqlUpdateParameterTemplateResponse{}
	err = resp.Parse(UpdateParameterTemplateResp)
	if err != nil {
		return
	}
	return UpdateParameterTemplateResp, nil
}

//func valueToString(value []ParameterObj) (string, error) {
//	jsonData, err := json.Marshal(value)
//	if err != nil {
//		return "", err
//	}
//	return string(jsonData), nil
//}

type ParameterObj struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PgsqlUpdateParameterTemplateRequest struct {
	TemplateId int64          `json:"templateId"`
	Params     []ParameterObj `json:"params"`
}

type PgsqlUpdateParameterTemplateRequestHeader struct {
	ProjectID *string `json:"Project-Id,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlUpdateParameterTemplateResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
