package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetCharacterSetApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetCharacterSetApi(client *ctyunsdk.CtyunClient) *TeledbGetCharacterSetApi {
	return &TeledbGetCharacterSetApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v1/open-api/database/show-character-set",
		},
	}
}

func (this *TeledbGetCharacterSetApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetCharacterSetRequest, header *TeledbGetCharacterSetRequestHeader) (GetCharacterSetResp *TeledbGetCharacterSetResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
	}
	if req.OuterProdInstId == "" || header.InstID == "" {
		err = errors.New("instId 为空")
		return
	}
	builder.AddHeader("inst-id", header.InstID)
	builder.AddHeader("regionId", header.RegionID)

	if req.OuterProdInstId == "" {
		err = errors.New("instId 为空")
		return
	}
	builder.AddHeader("inst-id", header.InstID)
	builder.AddParam("outerProdInstId", req.OuterProdInstId)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	GetCharacterSetResp = &TeledbGetCharacterSetResponse{}
	err = resp.Parse(GetCharacterSetResp)
	if err != nil {
		return
	}
	return GetCharacterSetResp, nil
}

type TeledbGetCharacterSetRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` // 外部实例ID，必填
}

type TeledbGetCharacterSetRequestHeader struct {
	ProjectID string `json:"project-id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}

type TeledbGetCharacterSetResponse struct {
	StatusCode int32                                   `json:"statusCode"`      // 接口状态码
	Error      *string                                 `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                  `json:"message"`         // 描述信息
	ReturnObj  *TeledbGetCharacterSetResponseReturnObj `json:"returnObj"`
}
type TeledbGetCharacterSetResponseReturnObjData struct {
	Charset          string `json:"charset"`
	Description      string `json:"description"`
	MaxLen           int32  `json:"maxlen"`
	DefaultCollation string `json:"defaultCollation"`
}

type TeledbGetCharacterSetResponseReturnObj struct {
	Data []TeledbGetCharacterSetResponseReturnObjData `json:"data"`
}
