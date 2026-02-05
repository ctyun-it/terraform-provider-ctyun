package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetParameterTemplateGroupDetailApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetParameterTemplateGroupDetailApi(client *ctyunsdk.CtyunClient) *TeledbGetParameterTemplateGroupDetailApi {
	return &TeledbGetParameterTemplateGroupDetailApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v1/open-api/parameter/describe-parameter-group",
		},
	}
}

func (this *TeledbGetParameterTemplateGroupDetailApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetParameterTemplateGroupDetailRequest, header *TeledbGetParameterTemplateGroupDetailRequestHeader) (GetParameterTemplateDetailResp *TeledbGetParameterTemplateGroupDetailResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}

	builder.AddHeader("regionId", header.RegionID)

	if req.ID == 0 {
		err = errors.New("id不能为空")
		return
	}
	builder.AddParam("id", fmt.Sprintf("%d", req.ID))

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	GetParameterTemplateDetailResp = &TeledbGetParameterTemplateGroupDetailResponse{}
	err = resp.Parse(GetParameterTemplateDetailResp)
	if err != nil {
		return
	}
	return GetParameterTemplateDetailResp, nil
}

type TeledbGetParameterTemplateGroupDetailRequest struct {
	ID int64 `json:"id"` //
}

type TeledbGetParameterTemplateGroupDetailRequestHeader struct {
	ProjectID *string `json:"Project-Id"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}

type TeledbGetParameterTemplateGroupDetailResponse struct {
	StatusCode int32                                                   `json:"statusCode"`      // 接口状态码
	Error      *string                                                 `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                                  `json:"message"`         // 描述信息
	ReturnObj  *TeledbGetParameterTemplateGroupDetailResponseReturnObj `json:"returnObj"`
}
type TeledbGetParameterTemplateGroupDetailResponseReturnObj struct {
	ParameterGroupName string                                                         `json:"parameterGroupName"`
	MysqlEngine        string                                                         `json:"mysqlEngine"` // M
	CreateTime         int64                                                          `json:"createTime"`
	Restart            bool                                                           `json:"restart"`
	ParameterCount     int64                                                          `json:"parameterCount"`
	Description        string                                                         `json:"description"`
	Isdefault          int32                                                          `json:"isdefault"`
	ID                 int64                                                          `json:"id"`
	ParamDetail        []TeledbGetParameterTemplateGroupDetailResponseReturnObjDeatil `json:"list"`
}

type TeledbGetParameterTemplateGroupDetailResponseReturnObjDeatil struct {
	ParamName  string `json:"ParamName"`
	ParamValue string `json:"ParamValue"`
}
