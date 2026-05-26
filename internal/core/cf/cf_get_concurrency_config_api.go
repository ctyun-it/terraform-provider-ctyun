package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfGetConcurrencyConfigApi
/* 查询函数的并发配额 */
type CfGetConcurrencyConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetConcurrencyConfigApi(client *core.CtyunClient) *CfGetConcurrencyConfigApi {
	return &CfGetConcurrencyConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/resources/functions/{functionName}/reservedConcurrency",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetConcurrencyConfigApi) Do(ctx context.Context, credential core.Credential, req *CfGetConcurrencyConfigRequest) (*CfGetConcurrencyConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfGetConcurrencyConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetConcurrencyConfigRequest struct {
	FunctionName string `json:"functionName"` /*  函数名称  */
	RegionId     string `json:"regionId"`     /*  资源池ID，请参考<a target="_blank" href="https://www.ctyun.cn/document/10006234/10985593">资源池列表</a>  */
}

type CfGetConcurrencyConfigResponse struct {
	StatusCode *int32                                   `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Error      *string                                  `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                                  `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfGetConcurrencyConfigReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfGetConcurrencyConfigReturnObjResponse struct {
	ReservedConcurrency *int32  `json:"reservedConcurrency"` /*  最大并发实例数，默认配额时取值范围 [0, 100]  */
	FunctionName        *string `json:"functionName"`        /*  函数名称  */
}
