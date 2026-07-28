package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfGetAsyncInvokeConfigApi
/* 查询函数的异步配置 */
type CfGetAsyncInvokeConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetAsyncInvokeConfigApi(client *core.CtyunClient) *CfGetAsyncInvokeConfigApi {
	return &CfGetAsyncInvokeConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/functions/{functionName}/async",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetAsyncInvokeConfigApi) Do(ctx context.Context, credential core.Credential, req *CfGetAsyncInvokeConfigRequest) (*CfGetAsyncInvokeConfigResponse, error) {
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
	var resp CfGetAsyncInvokeConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetAsyncInvokeConfigRequest struct {
	FunctionName string  `json:"functionName"`        /*  函数名称  */
	RegionId     string  `json:"regionId"`            /*  资源池ID，请参考<a target="_blank" href="https://www.ctyun.cn/document/10006234/10985593">资源池列表</a>  */
	Qualifier    *string `json:"qualifier,omitempty"` /*  函数的版本或别名，默认为 LATEST  */
}

type CfGetAsyncInvokeConfigResponse struct {
	StatusCode *int32                                   `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Error      *string                                  `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                                  `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfGetAsyncInvokeConfigReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfGetAsyncInvokeConfigReturnObjResponse struct {
	AsyncTask                 *bool                                                     `json:"asyncTask"`                 /*  是否开启异步任务  */
	DestinationConfig         *CfGetAsyncInvokeConfigReturnObjDestinationConfigResponse `json:"destinationConfig"`         /*  目标配置  */
	MaxAsyncRetryAttempts     *int32                                                    `json:"maxAsyncRetryAttempts"`     /*  异步调用重试次数，取值范围 [0, 8]  */
	MaxAsyncEventAgeInSeconds *int32                                                    `json:"maxAsyncEventAgeInSeconds"` /*  事件最大存活时间，取值范围 [1, 604800]，单位为 秒  */
}

type CfGetAsyncInvokeConfigReturnObjDestinationConfigResponse struct {
	OnSuccess *CfGetAsyncInvokeConfigReturnObjDestinationConfigOnSuccessResponse `json:"onSuccess"` /*  成功的回调目标  */
	OnFailure *CfGetAsyncInvokeConfigReturnObjDestinationConfigOnFailureResponse `json:"onFailure"` /*  失败的回调目标  */
}

type CfGetAsyncInvokeConfigReturnObjDestinationConfigOnSuccessResponse struct {
	Byname *string `json:"byname"` /*  版本号或别名  */
	Flag   *string `json:"flag"`   /*  版本或别名标志。version：版本，alias：别名  */
	Ksvc   *string `json:"ksvc"`   /*  ksvc名称  */
	Name   *string `json:"name"`   /*  函数名称  */
	Svc    *string `json:"svc"`    /*  目标服务。fc：函数计算  */
}

type CfGetAsyncInvokeConfigReturnObjDestinationConfigOnFailureResponse struct {
	Byname *string `json:"byname"` /*  版本号或别名  */
	Flag   *string `json:"flag"`   /*  版本或别名标志。version：版本，alias：别名  */
	Ksvc   *string `json:"ksvc"`   /*  ksvc名称  */
	Name   *string `json:"name"`   /*  函数名称  */
	Svc    *string `json:"svc"`    /*  目标服务。fc：函数计算  */
}
