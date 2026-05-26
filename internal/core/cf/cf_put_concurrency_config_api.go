package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfPutConcurrencyConfigApi
/* 设置函数的最大并发实例数 */
type CfPutConcurrencyConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfPutConcurrencyConfigApi(client *core.CtyunClient) *CfPutConcurrencyConfigApi {
	return &CfPutConcurrencyConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPut,
			UrlPath:      "/openapi/v1/resources/reservedConcurrency",
			ContentType:  "application/json",
		},
	}
}

func (a *CfPutConcurrencyConfigApi) Do(ctx context.Context, credential core.Credential, req *CfPutConcurrencyConfigRequest) (*CfPutConcurrencyConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfPutConcurrencyConfigRequest
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfPutConcurrencyConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfPutConcurrencyConfigRequest struct {
	RegionId            string `json:"regionId"`            /*  资源池ID，请参考<a target="_blank" href="https://www.ctyun.cn/document/10006234/10985593">资源池列表</a>  */
	FunctionName        string `json:"functionName"`        /*  函数名称  */
	ReservedConcurrency int32  `json:"reservedConcurrency"` /*  最大并发实例数，默认配额时取值范围 [0, 100]  */
}

type CfPutConcurrencyConfigResponse struct {
	StatusCode *int32                                   `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Error      *string                                  `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                                  `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfPutConcurrencyConfigReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfPutConcurrencyConfigReturnObjResponse struct {
	ReservedConcurrency *int32 `json:"reservedConcurrency"` /*  最大并发实例数，默认配额时取值范围 [0, 100]  */
}
