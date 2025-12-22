package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamAttachGroupToUserApi
type CtiamAttachGroupToUserApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamAttachGroupToUserApi(client *core.CtyunClient) *CtiamAttachGroupToUserApi {
	return &CtiamAttachGroupToUserApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/attachGroupToUser",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamAttachGroupToUserApi) Do(ctx context.Context, credential core.Credential, req *CtiamAttachGroupToUserRequest) (*CtiamAttachGroupToUserResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamAttachGroupToUserRequest
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
	var resp CtiamAttachGroupToUserResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamAttachGroupToUserRequest struct {
	UserId   string   `json:"userId"`   /*  用户id  */
	GroupIds []string `json:"groupIds"` /*  用户组id列表  */
}

type CtiamAttachGroupToUserResponse struct {
	StatusCode *string `json:"statusCode"` /*  状态码  */
	Message    *string `json:"message"`    /*  补充信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
