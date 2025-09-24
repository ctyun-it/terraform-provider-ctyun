package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosCreateRoleApi
/* 创建角色，通过STS创建的角色可以被授予特定的权限，这些权限可以用于访问不同的资源。
 */type ZosCreateRoleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosCreateRoleApi(client *core.CtyunClient) *ZosCreateRoleApi {
	return &ZosCreateRoleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/create-role",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosCreateRoleApi) Do(ctx context.Context, credential core.Credential, req *ZosCreateRoleRequest) (*ZosCreateRoleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosCreateRoleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosCreateRoleRequest struct {
	RegionID           string `json:"regionID,omitempty"`           /*  区域 ID  */
	RoleName           string `json:"roleName,omitempty"`           /*  角色名称，在资源池区域内唯一  */
	Path               string `json:"path,omitempty"`               /*  角色地址，必须以/开始和结尾  */
	AssumeName         string `json:"assumeName,omitempty"`         /*  被授权客户端，目前仅支持sts  */
	Note               string `json:"note,omitempty"`               /*  备注  */
	MaxSessionDuration int64  `json:"maxSessionDuration,omitempty"` /*  最大会话时间，单位秒，取值范围：3600-43200秒，即1-12小时，默认 3600  */
}

type ZosCreateRoleResponse struct {
	StatusCode  int64                           `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                          `json:"message,omitempty"`     /*  状态描述  */
	ReturnObj   *ZosCreateRoleReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	Description string                          `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string                          `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                          `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosCreateRoleReturnObjResponse struct {
	Arn string `json:"arn,omitempty"` /*  角色arn  */
}
