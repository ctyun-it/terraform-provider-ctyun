package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamGetGroupUserApi
/* 分页查询用户组下的用户 */
type CtiamGetGroupUserApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamGetGroupUserApi(client *core.CtyunClient) *CtiamGetGroupUserApi {
	return &CtiamGetGroupUserApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/userGroup/getGroupUser",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamGetGroupUserApi) Do(ctx context.Context, credential core.Credential, req *CtiamGetGroupUserRequest) (*CtiamGetGroupUserResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamGetGroupUserRequest
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
	var resp CtiamGetGroupUserResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamGetGroupUserRequest struct {
	PageNum  int32                             `json:"pageNum"`          /*  页码  */
	PageSize int32                             `json:"pageSize"`         /*  每页条数  */
	Groups   []*CtiamGetGroupUserGroupsRequest `json:"groups,omitempty"` /*  用户组ID列表  */
}

type CtiamGetGroupUserGroupsRequest struct {
	Id string `json:"id"` /*  用户组ID  */
}

type CtiamGetGroupUserResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *string `json:"returnObj"`  /*  返回参数  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息。  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
