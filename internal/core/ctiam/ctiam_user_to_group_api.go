package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamUserToGroupApi
/* 将用户移入用户组 */
type CtiamUserToGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamUserToGroupApi(client *core.CtyunClient) *CtiamUserToGroupApi {
	return &CtiamUserToGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/userGroup/userToGroup",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamUserToGroupApi) Do(ctx context.Context, credential core.Credential, req *CtiamUserToGroupRequest) (*CtiamUserToGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamUserToGroupRequest
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
	var resp CtiamUserToGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamUserToGroupRequest struct {
	GroupId string                            `json:"groupId"`           /*  用户组id  */
	UserIds []*CtiamUserToGroupUserIdsRequest `json:"userIds,omitempty"` /*  用户信息  */
}

type CtiamUserToGroupUserIdsRequest struct {
	Id string `json:"id"` /*  用户id  */
}

type CtiamUserToGroupResponse struct {
	StatusCode *string `json:"statusCode"` /*  状态码  */
	Message    *string `json:"message"`    /*  对于调用失败的补充信息说明  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
