package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamSetEpGroupPloyApi
/* 设置企业项目所属用户组及策略 */
type CtiamSetEpGroupPloyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamSetEpGroupPloyApi(client *core.CtyunClient) *CtiamSetEpGroupPloyApi {
	return &CtiamSetEpGroupPloyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/project/setEpGroupPloy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamSetEpGroupPloyApi) Do(ctx context.Context, credential core.Credential, req *CtiamSetEpGroupPloyRequest) (*CtiamSetEpGroupPloyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamSetEpGroupPloyRequest
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
	var resp CtiamSetEpGroupPloyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamSetEpGroupPloyRequest struct {
	GroupId   string `json:"groupId"`   /*  用户组id  */
	ProjectId string `json:"projectId"` /*  企业项目id  */
	PloyIds   string `json:"ployIds"`   /*  策略id（策略ID,可传多个，以逗号分割）  */
}

type CtiamSetEpGroupPloyResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息
	 */
	Error *string `json:"error"` /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
