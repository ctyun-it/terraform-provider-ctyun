package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfDeleteAsyncInvokeConfigApi
/* 删除的异步配置 */
type CfDeleteAsyncInvokeConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfDeleteAsyncInvokeConfigApi(client *core.CtyunClient) *CfDeleteAsyncInvokeConfigApi {
	return &CfDeleteAsyncInvokeConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodDelete,
			UrlPath:      "/openapi/v1/functions/{functionName}/async",
			ContentType:  "application/json",
		},
	}
}

func (a *CfDeleteAsyncInvokeConfigApi) Do(ctx context.Context, credential core.Credential, req *CfDeleteAsyncInvokeConfigRequest) (*CfDeleteAsyncInvokeConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.Qualifier != nil && *req.Qualifier != "" {
		ctReq.AddParam("qualifier", *req.Qualifier)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfDeleteAsyncInvokeConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfDeleteAsyncInvokeConfigRequest struct {
	FunctionName string  `json:"functionName"`        /*  函数名称  */
	RegionId     string  `json:"regionId"`            /*  资源池ID，请参考<a target="_blank" href="https://www.ctyun.cn/document/10006234/10985593">资源池列表</a>  */
	Qualifier    *string `json:"qualifier,omitempty"` /*  函数的版本或别名，默认为 LATEST  */
}

type CfDeleteAsyncInvokeConfigResponse struct {
	StatusCode *int32  `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Error      *string `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string `json:"message"`    /*  错误描述信息  */
}
