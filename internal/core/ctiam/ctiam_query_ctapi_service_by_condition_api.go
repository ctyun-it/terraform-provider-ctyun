package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryCtapiServiceByConditionApi
type CtiamQueryCtapiServiceByConditionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryCtapiServiceByConditionApi(client *core.CtyunClient) *CtiamQueryCtapiServiceByConditionApi {
	return &CtiamQueryCtapiServiceByConditionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/service/queryCtapiServiceByCondition",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryCtapiServiceByConditionApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryCtapiServiceByConditionRequest) (*CtiamQueryCtapiServiceByConditionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamQueryCtapiServiceByConditionRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamQueryCtapiServiceByConditionResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryCtapiServiceByConditionRequest struct {
	ServiceName *string `json:"serviceName,omitempty"` /*  服务名称（中文）  */
	ServiceType *int8   `json:"serviceType,omitempty"` /*  服务/产品类型
	1：资源池级云服务
	2：全局级云服务  */
}

type CtiamQueryCtapiServiceByConditionResponse struct {
	StatusCode *string                                             `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                             `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                             `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamQueryCtapiServiceByConditionReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamQueryCtapiServiceByConditionReturnObjResponse struct {
	ServiceList []*CtiamQueryCtapiServiceByConditionReturnObjServiceListResponse `json:"serviceList"` /*  服务列表  */
}

type CtiamQueryCtapiServiceByConditionReturnObjServiceListResponse struct {
	ServiceCode *string `json:"serviceCode"` /*  服务编码  */
	ServiceType *int8   `json:"serviceType"` /*  服务/产品类型
	1：资源池级云服务
	2：全局级云服务  */
	MainServiceName *string `json:"mainServiceName"` /*  服务名称  */
	ServiceDesc     *string `json:"serviceDesc"`     /*  服务描述  */
	Id              *int32  `json:"id"`              /*  云服务ID  */
}
