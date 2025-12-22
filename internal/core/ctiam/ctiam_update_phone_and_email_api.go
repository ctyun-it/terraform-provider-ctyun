package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamUpdatePhoneAndEmailApi
type CtiamUpdatePhoneAndEmailApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamUpdatePhoneAndEmailApi(client *core.CtyunClient) *CtiamUpdatePhoneAndEmailApi {
	return &CtiamUpdatePhoneAndEmailApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/updatePhoneAndEmail",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamUpdatePhoneAndEmailApi) Do(ctx context.Context, credential core.Credential, req *CtiamUpdatePhoneAndEmailRequest) (*CtiamUpdatePhoneAndEmailResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamUpdatePhoneAndEmailRequest
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
	var resp CtiamUpdatePhoneAndEmailResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamUpdatePhoneAndEmailRequest struct {
	UserId      string  `json:"userId"`                /*  用户ID  */
	LoginEmail  *string `json:"loginEmail,omitempty"`  /*  邮箱（邮箱与手机号必填其中之一）  */
	MobilePhone *string `json:"mobilePhone,omitempty"` /*  手机号（邮箱与手机号必填其中之一）  */
}

type CtiamUpdatePhoneAndEmailResponse struct {
	StatusCode *string                                    `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                    `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                    `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamUpdatePhoneAndEmailReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamUpdatePhoneAndEmailReturnObjResponse struct {
	LoginEmail  *string `json:"loginEmail"`  /*  邮箱  */
	AccountId   *string `json:"accountId"`   /*  账号ID  */
	MobilePhone *string `json:"mobilePhone"` /*  手机号  */
}
