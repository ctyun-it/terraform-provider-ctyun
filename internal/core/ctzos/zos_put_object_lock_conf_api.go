package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPutObjectLockConfApi
/* 设置桶的合规保留策略，已继承该桶的合规保留期限的对象在保留期内不可被彻底删除和篡改。
 */type ZosPutObjectLockConfApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPutObjectLockConfApi(client *core.CtyunClient) *ZosPutObjectLockConfApi {
	return &ZosPutObjectLockConfApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/put-object-lock-conf",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPutObjectLockConfApi) Do(ctx context.Context, credential core.Credential, req *ZosPutObjectLockConfRequest) (*ZosPutObjectLockConfResponse, error) {
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
	var resp ZosPutObjectLockConfResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPutObjectLockConfRequest struct {
	Bucket        string `json:"bucket,omitempty"`        /*  桶名  */
	RegionID      string `json:"regionID,omitempty"`      /*  区域 ID  */
	RetentionMode string `json:"retentionMode,omitempty"` /*  保留模式，必须为 COMPLIANCE 或 GOVERNANCE  */
	Days          int64  `json:"days,omitempty"`          /*  天数（days 与 years 参数必须存在其一，但不能同时存在）  */
	Years         int64  `json:"years,omitempty"`         /*  年数（days 与 years 参数必须存在其一，但不能同时存在）  */
}

type ZosPutObjectLockConfResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
