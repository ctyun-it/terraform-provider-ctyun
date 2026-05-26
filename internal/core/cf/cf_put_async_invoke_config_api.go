package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfPutAsyncInvokeConfigApi
/* 创建或更新函数的异步配置 */
type CfPutAsyncInvokeConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfPutAsyncInvokeConfigApi(client *core.CtyunClient) *CfPutAsyncInvokeConfigApi {
	return &CfPutAsyncInvokeConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/openapi/v1/functions/{functionName}/async",
			ContentType:  "application/json",
		},
	}
}

func (a *CfPutAsyncInvokeConfigApi) Do(ctx context.Context, credential core.Credential, req *CfPutAsyncInvokeConfigRequest) (*CfPutAsyncInvokeConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfPutAsyncInvokeConfigRequest
		RegionId     interface{} `json:"regionId,omitempty"`
		FunctionName interface{} `json:"functionName,omitempty"`
	}{
		req, nil, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfPutAsyncInvokeConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfPutAsyncInvokeConfigRequest struct {
	FunctionName              string                                          `json:"functionName"`                /*  函数名称  */
	RegionId                  string                                          `json:"regionId"`                    /*  资源池ID，请参考<a target="_blank" href="https://www.ctyun.cn/document/10006234/10985593">资源池列表</a>  */
	AsyncTask                 bool                                            `json:"asyncTask"`                   /*  是否开启异步任务模式  */
	DestinationConfig         *CfPutAsyncInvokeConfigDestinationConfigRequest `json:"destinationConfig,omitempty"` /*  目标配置  */
	MaxAsyncRetryAttempts     int32                                           `json:"maxAsyncRetryAttempts"`       /*  异步调用重试次数，取值范围 [0, 8]  */
	MaxAsyncEventAgeInSeconds int32                                           `json:"maxAsyncEventAgeInSeconds"`   /*  事件最大存活时间，取值范围 [1, 604800]，单位为 秒  */
	Qualifier                 *string                                         `json:"qualifier,omitempty"`         /*  函数的版本或别名，默认为 LATEST  */
}

type CfPutAsyncInvokeConfigDestinationConfigRequest struct {
	OnSuccess *CfPutAsyncInvokeConfigDestinationConfigOnSuccessRequest `json:"onSuccess,omitempty"` /*  成功的回调目标  */
	OnFailure *CfPutAsyncInvokeConfigDestinationConfigOnFailureRequest `json:"onFailure,omitempty"` /*  失败的回调目标  */
}

type CfPutAsyncInvokeConfigDestinationConfigOnSuccessRequest struct {
	Byname string `json:"byname"` /*  版本号或别名  */
	Flag   string `json:"flag"`   /*  版本或别名标志。version：版本，alias：别名  */
	Ksvc   string `json:"ksvc"`   /*  ksvc名称  */
	Name   string `json:"name"`   /*  函数名称  */
	Svc    string `json:"svc"`    /*  目标服务。fc：函数计算  */
}

type CfPutAsyncInvokeConfigDestinationConfigOnFailureRequest struct {
	Byname string `json:"byname"` /*  版本号或别名  */
	Flag   string `json:"flag"`   /*  版本或别名标志。version：版本，alias：别名  */
	Ksvc   string `json:"ksvc"`   /*  ksvc名称  */
	Name   string `json:"name"`   /*  函数名称  */
	Svc    string `json:"svc"`    /*  目标服务。fc：函数计算  */
}

type CfPutAsyncInvokeConfigResponse struct {
	StatusCode *int32                                   `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Code       *string                                  `json:"code"`       /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                                  `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfPutAsyncInvokeConfigReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfPutAsyncInvokeConfigReturnObjResponse struct{}
