package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfDeleteConcurrencyConfigApi
/* 删除函数的并发配额 */
type CfDeleteConcurrencyConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfDeleteConcurrencyConfigApi(client *core.CtyunClient) *CfDeleteConcurrencyConfigApi {
	return &CfDeleteConcurrencyConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodDelete,
			UrlPath:      "/openapi/v1/resources/functions/{functionName}/concurrency",
			ContentType:  "application/json",
		},
	}
}

func (a *CfDeleteConcurrencyConfigApi) Do(ctx context.Context, credential core.Credential, req *CfDeleteConcurrencyConfigRequest) (*CfDeleteConcurrencyConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfDeleteConcurrencyConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfDeleteConcurrencyConfigRequest struct {
	FunctionName string `json:"functionName"` /*  函数名称  */
	RegionId     string `json:"regionId"`     /*  资源池ID，请参考<a target="_blank" href="https://www.ctyun.cn/document/10006234/10985593">资源池列表</a>  */
}

type CfDeleteConcurrencyConfigResponse struct {
	StatusCode *int32  `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Error      *string `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string `json:"message"`    /*  错误描述信息  */
}
