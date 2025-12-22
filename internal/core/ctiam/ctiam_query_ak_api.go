package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryAkApi
/* 查询密钥 */
type CtiamQueryAkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryAkApi(client *core.CtyunClient) *CtiamQueryAkApi {
	return &CtiamQueryAkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/credential/queryAk",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryAkApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryAkRequest) (*CtiamQueryAkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamQueryAkRequest
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
	var resp CtiamQueryAkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryAkRequest struct {
	UserIdList []string `json:"userIdList"` /*  用户id列表  */
}

type CtiamQueryAkResponse struct {
	StatusCode *string                        `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *CtiamQueryAkReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                        `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                        `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamQueryAkReturnObjResponse struct {
	AccessKeyUserList []*CtiamQueryAkReturnObjAccessKeyUserListResponse `json:"accessKeyUserList"` /*  用户密钥列表  */
}

type CtiamQueryAkReturnObjAccessKeyUserListResponse struct {
	UserId        *string                                                        `json:"userId"`        /*  用户id  */
	AccessKeyList []*CtiamQueryAkReturnObjAccessKeyUserListAccessKeyListResponse `json:"accessKeyList"` /*  密钥列表  */
}

type CtiamQueryAkReturnObjAccessKeyUserListAccessKeyListResponse struct {
	AccessKey   *string `json:"accessKey"`   /*  AK  */
	SecretKey   *string `json:"secretKey"`   /*  SM4算法加密后SK  */
	UserId      *string `json:"userId"`      /*  用户ID  */
	AccountId   *string `json:"accountId"`   /*  账号ID  */
	CreatedTime int64   `json:"createdTime"` /*  创建时间  */
	Status      *string `json:"status"`      /*  状态 1000  启用 1001  禁用  */
}
