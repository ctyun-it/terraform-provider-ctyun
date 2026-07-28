package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamChangeResourceEpApi
/* 资源迁入/迁出企业项目 */
type CtiamChangeResourceEpApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamChangeResourceEpApi(client *core.CtyunClient) *CtiamChangeResourceEpApi {
	return &CtiamChangeResourceEpApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/resource/changeResourceEp",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamChangeResourceEpApi) Do(ctx context.Context, credential core.Credential, req *CtiamChangeResourceEpRequest) (*CtiamChangeResourceEpResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamChangeResourceEpRequest
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
	var resp CtiamChangeResourceEpResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamChangeResourceEpRequest struct {
	ProjectSetId string `json:"projectSetId"` /*  企业项目id  */
	ResourceId   string `json:"resourceId"`   /*  资源id（资源id,多个以逗号分割）  */
	IsEcs        string `json:"isEcs"`        /*  是否ecs关联迁入 ，1:是, 0:否  */
}

type CtiamChangeResourceEpResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
