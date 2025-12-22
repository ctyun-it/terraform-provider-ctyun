package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetStatusByIDApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetStatusByIDApi(client *ctyunsdk.CtyunClient) *TeledbGetStatusByIDApi {
	return &TeledbGetStatusByIDApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/teledb-acceptor/v2/openapi/dcp-order-info/getInstState",
		},
	}
}

func (this *TeledbGetStatusByIDApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetStatusByIDRequest, header *TeledbGetStatusByIDRequestHeader) (destroyResp *TeledbGetStatusByIDResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
	}
	if req.InstID == "" {
		err = errors.New("instID为空")
		return
	}
	builder.AddParam("instID", req.InstID)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	destroyResp = &TeledbGetStatusByIDResponse{}
	err = resp.Parse(destroyResp)
	if err != nil {
		return
	}
	return destroyResp, nil
}

type TeledbGetStatusByIDRequest struct {
	InstID string `json:"instID"` // 实例ID，必填
}
type TeledbGetStatusByIDRequestHeader struct {
	ProjectID string `json:"projectID"`
}
type TeledbGetStatusByIDResponse struct {
	StatusCode int32                                 `json:"statusCode"` // 接口状态码
	Error      *string                               `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string                                `json:"message"`    // 描述信息
	ReturnObj  *TeledbGetStatusByIDResponseReturnObj `json:"returnObj"`  // 返回对象，类型为 DataObject
}

type TeledbGetStatusByIDResponseReturnObj struct {
	Data int `json:"data"` // 0-施工中 1-施工失败 2-运行中 3-已暂停 4-扩容中 5-已注销
}
