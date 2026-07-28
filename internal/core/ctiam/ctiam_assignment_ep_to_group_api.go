package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamAssignmentEpToGroupApi
/* 用户组与企业项目关联 */
type CtiamAssignmentEpToGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamAssignmentEpToGroupApi(client *core.CtyunClient) *CtiamAssignmentEpToGroupApi {
	return &CtiamAssignmentEpToGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/project/assignmentEpToGroup",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamAssignmentEpToGroupApi) Do(ctx context.Context, credential core.Credential, req *CtiamAssignmentEpToGroupRequest) (*CtiamAssignmentEpToGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamAssignmentEpToGroupRequest
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
	var resp CtiamAssignmentEpToGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamAssignmentEpToGroupRequest struct {
	ProjectId string   `json:"projectId"` /*  企业项目id  */
	GroupIds  []string `json:"groupIds"`  /*  用户组id  */
}

type CtiamAssignmentEpToGroupResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
