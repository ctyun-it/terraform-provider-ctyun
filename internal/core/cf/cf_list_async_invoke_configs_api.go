package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CfListAsyncInvokeConfigsApi
/* 列出函数的异步配置 */
type CfListAsyncInvokeConfigsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfListAsyncInvokeConfigsApi(client *core.CtyunClient) *CfListAsyncInvokeConfigsApi {
	return &CfListAsyncInvokeConfigsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/async-configs",
			ContentType:  "application/json",
		},
	}
}

func (a *CfListAsyncInvokeConfigsApi) Do(ctx context.Context, credential core.Credential, req *CfListAsyncInvokeConfigsRequest) (*CfListAsyncInvokeConfigsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.FunctionName != nil && *req.FunctionName != "" {
		ctReq.AddParam("functionName", *req.FunctionName)
	}
	if req.PageIndex != nil && *req.PageIndex != 0 {
		ctReq.AddParam("pageIndex", strconv.FormatInt(int64(*req.PageIndex), 10))
	}
	if req.PageSize != nil && *req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(*req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfListAsyncInvokeConfigsResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfListAsyncInvokeConfigsRequest struct {
	RegionId     string  `json:"regionId"`               /*  资源池ID，请参考<a target="_blank" href="https://www.ctyun.cn/document/10006234/10985593">资源池列表</a>  */
	FunctionName *string `json:"functionName,omitempty"` /*  函数名称，若不指定则列出所有函数的异步调用配置  */
	PageIndex    *int32  `json:"pageIndex,omitempty"`    /*  页码，取值 >= 1，默认值为 1  */
	PageSize     *int32  `json:"pageSize,omitempty"`     /*  每页大小，取值范围[1, 100]，默认值为 50  */
}

type CfListAsyncInvokeConfigsResponse struct {
	StatusCode *int32                                     `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Error      *string                                    `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                                    `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfListAsyncInvokeConfigsReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfListAsyncInvokeConfigsReturnObjResponse struct {
	Data       []*CfListAsyncInvokeConfigsReturnObjDataResponse     `json:"data"`       /*  分页数据  */
	Pagination *CfListAsyncInvokeConfigsReturnObjPaginationResponse `json:"pagination"` /*  分页信息  */
}

type CfListAsyncInvokeConfigsReturnObjDataResponse struct {
	AsyncTask                 *bool                                                           `json:"asyncTask"`                 /*  是否开启异步任务  */
	DestinationConfig         *CfListAsyncInvokeConfigsReturnObjDataDestinationConfigResponse `json:"destinationConfig"`         /*  目标配置  */
	FunctionName              *string                                                         `json:"functionName"`              /*  函数名称  */
	MaxAsyncRetryAttempts     *int32                                                          `json:"maxAsyncRetryAttempts"`     /*  异步调用重试次数，取值范围 [0, 8]  */
	MaxAsyncEventAgeInSeconds *int32                                                          `json:"maxAsyncEventAgeInSeconds"` /*  事件最大存活时间，取值范围 [1, 604800]，单位为 秒  */
	Qualifier                 *string                                                         `json:"qualifier"`                 /*  函数的版本或别名  */
}

type CfListAsyncInvokeConfigsReturnObjPaginationResponse struct {
	PageIndex *int32 `json:"pageIndex"` /*  页码  */
	PageSize  *int32 `json:"pageSize"`  /*  每页大小  */
	Total     *int32 `json:"total"`     /*  总记录数  */
}

type CfListAsyncInvokeConfigsReturnObjDataDestinationConfigResponse struct {
	OnSuccess *CfListAsyncInvokeConfigsReturnObjDataDestinationConfigOnSuccessResponse `json:"onSuccess"` /*  成功的回调目标  */
	OnFailure *CfListAsyncInvokeConfigsReturnObjDataDestinationConfigOnFailureResponse `json:"onFailure"` /*  失败的回调目标  */
}

type CfListAsyncInvokeConfigsReturnObjDataDestinationConfigOnSuccessResponse struct {
	Byname *string `json:"byname"` /*  版本号或别名  */
	Flag   *string `json:"flag"`   /*  版本或别名标志。version：版本，alias：别名  */
	Ksvc   *string `json:"ksvc"`   /*  ksvc名称  */
	Name   *string `json:"name"`   /*  函数名称  */
	Svc    *string `json:"svc"`    /*  目标服务。fc：函数计算  */
}

type CfListAsyncInvokeConfigsReturnObjDataDestinationConfigOnFailureResponse struct {
	Byname *string `json:"byname"` /*  版本号或别名  */
	Flag   *string `json:"flag"`   /*  版本或别名标志。version：版本，alias：别名  */
	Ksvc   *string `json:"ksvc"`   /*  ksvc名称  */
	Name   *string `json:"name"`   /*  函数名称  */
	Svc    *string `json:"svc"`    /*  目标服务。fc：函数计算  */
}
