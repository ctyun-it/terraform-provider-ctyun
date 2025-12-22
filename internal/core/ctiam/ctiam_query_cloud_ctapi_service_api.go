package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryCloudCtapiServiceApi
type CtiamQueryCloudCtapiServiceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryCloudCtapiServiceApi(client *core.CtyunClient) *CtiamQueryCloudCtapiServiceApi {
	return &CtiamQueryCloudCtapiServiceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/delegate/queryCloudCtapiService",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryCloudCtapiServiceApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryCloudCtapiServiceRequest) (*CtiamQueryCloudCtapiServiceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamQueryCloudCtapiServiceResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryCloudCtapiServiceRequest struct{}

type CtiamQueryCloudCtapiServiceResponse struct {
	StatusCode *string                                       `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                       `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                       `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamQueryCloudCtapiServiceReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamQueryCloudCtapiServiceReturnObjResponse struct {
	ServiceList []*CtiamQueryCloudCtapiServiceReturnObjServiceListResponse `json:"serviceList"` /*  服务列表  */
}

type CtiamQueryCloudCtapiServiceReturnObjServiceListResponse struct {
	ServiceName *string `json:"serviceName"` /*  服务名称  */
	ServiceCode *string `json:"serviceCode"` /*  产品编码  */
	AccountId   *string `json:"accountId"`   /*  账号ID  */
}
