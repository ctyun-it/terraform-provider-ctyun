package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamRemoveGroupFromUserApi
type CtiamRemoveGroupFromUserApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamRemoveGroupFromUserApi(client *core.CtyunClient) *CtiamRemoveGroupFromUserApi {
	return &CtiamRemoveGroupFromUserApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/removeGroupFromUser",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamRemoveGroupFromUserApi) Do(ctx context.Context, credential core.Credential, req *CtiamRemoveGroupFromUserRequest) (*CtiamRemoveGroupFromUserResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamRemoveGroupFromUserRequest
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
	var resp CtiamRemoveGroupFromUserResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamRemoveGroupFromUserRequest struct {
	UserId   string   `json:"userId"`   /*  用户id  */
	GroupIds []string `json:"groupIds"` /*  用户组id列表  */
}

type CtiamRemoveGroupFromUserResponse struct {
	StatusCode *string `json:"statusCode"` /*  状态码  */
	Message    *string `json:"message"`    /*  对于调用失败的补充信息说明  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
