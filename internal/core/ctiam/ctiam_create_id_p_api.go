package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamCreateIdPApi
/* 创建身份提供商 */
type CtiamCreateIdPApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamCreateIdPApi(client *core.CtyunClient) *CtiamCreateIdPApi {
	return &CtiamCreateIdPApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/identityProvider/createIdP",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamCreateIdPApi) Do(ctx context.Context, credential core.Credential, req *CtiamCreateIdPRequest) (*CtiamCreateIdPResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamCreateIdPRequest
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
	var resp CtiamCreateIdPResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamCreateIdPRequest struct {
	FileName string  `json:"fileName"`         /*  文件名称（需携带后缀）  */
	File     []int8  `json:"file"`             /*  文件元数据（byte数组）  */
	RawType  int32   `json:"type"`             /*  类型，0 虚拟用户SSO，1 IAM用户SSO  */
	Name     string  `json:"name"`             /*  身份提供商名称  */
	Protocol int32   `json:"protocol"`         /*  协议类型，0 SAML协议，1 OIDC协议  */
	Remark   *string `json:"remark,omitempty"` /*  描述  */
}

type CtiamCreateIdPResponse struct {
	StatusCode *string                          `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                          `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                          `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamCreateIdPReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamCreateIdPReturnObjResponse struct {
	Id         *string `json:"id"`         /*  ID  */
	Name       *string `json:"name"`       /*  身份提供商名称  */
	RawType    *int32  `json:"type"`       /*  类型，0 虚拟用户SSO，1 IAM用户SSO  */
	Protocol   *int32  `json:"protocol"`   /*  协议类型，0 SAML协议，1 OIDC协议  */
	AccountId  *string `json:"accountId"`  /*  账号ID  */
	Remark     *string `json:"remark"`     /*  描述  */
	FileName   *string `json:"fileName"`   /*  文件名称  */
	CreateTime *string `json:"createTime"` /*  创建时间  */
	Uuid       *string `json:"uuid"`       /*  uuid  */
}
