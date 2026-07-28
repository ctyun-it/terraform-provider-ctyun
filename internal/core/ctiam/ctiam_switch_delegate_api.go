package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamSwitchDelegateApi
/* 切换委托 */
type CtiamSwitchDelegateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamSwitchDelegateApi(client *core.CtyunClient) *CtiamSwitchDelegateApi {
	return &CtiamSwitchDelegateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/delegate/switchDelegate",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamSwitchDelegateApi) Do(ctx context.Context, credential core.Credential, req *CtiamSwitchDelegateRequest) (*CtiamSwitchDelegateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamSwitchDelegateRequest
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
	var resp CtiamSwitchDelegateResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamSwitchDelegateRequest struct {
	Name         *string `json:"name,omitempty"`         /*  委托名称（委托名称和委托用户id必填其中一个）  */
	AssumeUserId *string `json:"assumeUserId,omitempty"` /*  委托用户id（委托名称和委托用户id必填其中一个）  */
	AccountId    string  `json:"accountId"`              /*  账号id  */
	ValidTime    *int32  `json:"validTime,omitempty"`    /*  委托AK有效时间（单位：分钟），有效时间不得少于30分钟或大于36小时，不填默认30分钟  */
}

type CtiamSwitchDelegateResponse struct {
	StatusCode *string                               `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *CtiamSwitchDelegateReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                               `json:"message"`    /*  返回信息，请求失败时会回传错误信息。  */
	Error      *string                               `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamSwitchDelegateReturnObjResponse struct {
	AppId      *string `json:"appId"`      /*  ak  */
	AppKey     *string `json:"appKey"`     /*  sk（加密后）  */
	ExpireTime *string `json:"expireTime"` /*  过期时间  */
	Token      *string `json:"token"`      /*  token  */
}
