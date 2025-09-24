package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetObjectLockConfApi
/* 获取指定桶的合规保留策略信息。
 */type ZosGetObjectLockConfApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetObjectLockConfApi(client *core.CtyunClient) *ZosGetObjectLockConfApi {
	return &ZosGetObjectLockConfApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-object-lock-conf",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetObjectLockConfApi) Do(ctx context.Context, credential core.Credential, req *ZosGetObjectLockConfRequest) (*ZosGetObjectLockConfResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetObjectLockConfResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetObjectLockConfRequest struct {
	Bucket   string /*  桶名  */
	RegionID string /*  区域 ID  */
}

type ZosGetObjectLockConfResponse struct {
	StatusCode  int64                                  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                 `json:"message,omitempty"`     /*  状态描述  */
	Description string                                 `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosGetObjectLockConfReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                                 `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                 `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetObjectLockConfReturnObjResponse struct {
	ObjectLockConfiguration *ZosGetObjectLockConfReturnObjObjectLockConfigurationResponse `json:"objectLockConfiguration"` /*  对象锁定配置  */
}

type ZosGetObjectLockConfReturnObjObjectLockConfigurationResponse struct {
	ObjectLockEnabled string                                                            `json:"objectLockEnabled,omitempty"` /*  应始终为 Enabled  */
	Rule              *ZosGetObjectLockConfReturnObjObjectLockConfigurationRuleResponse `json:"rule"`                        /*  规则对象  */
}

type ZosGetObjectLockConfReturnObjObjectLockConfigurationRuleResponse struct {
	DefaultRetention *ZosGetObjectLockConfReturnObjObjectLockConfigurationRuleDefaultRetentionResponse `json:"defaultRetention"` /*  默认保留配置  */
}

type ZosGetObjectLockConfReturnObjObjectLockConfigurationRuleDefaultRetentionResponse struct {
	Mode  *ZosGetObjectLockConfReturnObjObjectLockConfigurationRuleDefaultRetentionModeResponse `json:"mode"`            /*  保留模式  */
	Days  int64                                                                                 `json:"days,omitempty"`  /*  保留天数  */
	Years int64                                                                                 `json:"years,omitempty"` /*  保留年数  */
}

type ZosGetObjectLockConfReturnObjObjectLockConfigurationRuleDefaultRetentionModeResponse struct{}
