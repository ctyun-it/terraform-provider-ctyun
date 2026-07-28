package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamRemoveEpFromGroupApi
/* 移除企业项目关联用户组 */
type CtiamRemoveEpFromGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamRemoveEpFromGroupApi(client *core.CtyunClient) *CtiamRemoveEpFromGroupApi {
	return &CtiamRemoveEpFromGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/project/removeEpFromGroup",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamRemoveEpFromGroupApi) Do(ctx context.Context, credential core.Credential, req *CtiamRemoveEpFromGroupRequest) (*CtiamRemoveEpFromGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamRemoveEpFromGroupRequest
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
	var resp CtiamRemoveEpFromGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamRemoveEpFromGroupRequest struct {
	ProjectId string   `json:"projectId"` /*  企业项目id  */
	GroupIds  []string `json:"groupIds"`  /*  用户组id  */
}

type CtiamRemoveEpFromGroupResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
