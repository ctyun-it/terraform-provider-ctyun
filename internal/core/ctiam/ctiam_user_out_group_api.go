package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamUserOutGroupApi
/* 将用户移出用户组 */
type CtiamUserOutGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamUserOutGroupApi(client *core.CtyunClient) *CtiamUserOutGroupApi {
	return &CtiamUserOutGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/userGroup/userOutGroup",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamUserOutGroupApi) Do(ctx context.Context, credential core.Credential, req *CtiamUserOutGroupRequest) (*CtiamUserOutGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamUserOutGroupRequest
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
	var resp CtiamUserOutGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamUserOutGroupRequest struct {
	GroupId string                             `json:"groupId"`           /*  用户组id  */
	UserIds []*CtiamUserOutGroupUserIdsRequest `json:"userIds,omitempty"` /*  用户列表  */
}

type CtiamUserOutGroupUserIdsRequest struct {
	Id string `json:"id"` /*  用户id  */
}

type CtiamUserOutGroupResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息。  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
