package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamAttachEpPermissionApi
/* 授权企业项目 */
type CtiamAttachEpPermissionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamAttachEpPermissionApi(client *core.CtyunClient) *CtiamAttachEpPermissionApi {
	return &CtiamAttachEpPermissionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/project/attachEpPermission",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamAttachEpPermissionApi) Do(ctx context.Context, credential core.Credential, req *CtiamAttachEpPermissionRequest) (*CtiamAttachEpPermissionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamAttachEpPermissionRequest
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
	var resp CtiamAttachEpPermissionResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamAttachEpPermissionRequest struct {
	GroupName     *string  `json:"groupName,omitempty"`   /*  用户组名称（groupId与groupName必填其中之一）  */
	GroupIntro    *string  `json:"groupIntro,omitempty"`  /*  用户组描述（groupId为空时，可填, groupId不为空时,不用填）  */
	GroupId       *string  `json:"groupId,omitempty"`     /*  用户组ID（groupId与groupName必填其中之一））  */
	AssumeUserIds []string `json:"assumeUserIds"`         /*  用户ID列表  */
	ProjectName   *string  `json:"projectName,omitempty"` /*  企业项目名称（projectName与epId必填其中之一）  */
	Description   *string  `json:"description,omitempty"` /*  企业项目名称（epId为空时,可填，epId不为空时, 不用填）  */
	EpId          *string  `json:"epId,omitempty"`        /*  企业项目Id（projectName与epId必填其中之一）  */
	PloyIds       string   `json:"ployIds"`               /*  策略id（策略ID,可传多个，以逗号分割）  */
}

type CtiamAttachEpPermissionResponse struct {
	StatusCode *string                                   `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *CtiamAttachEpPermissionReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                                   `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                   `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamAttachEpPermissionReturnObjResponse struct {
	GroupName     *string   `json:"groupName"`     /*  用户组名称  */
	GroupIntro    *string   `json:"groupIntro"`    /*  用户组描述  */
	GroupId       *string   `json:"groupId"`       /*  用户组Id  */
	AssumeUserIds []*string `json:"assumeUserIds"` /*  用户Id列表  */
	EpId          *string   `json:"epId"`          /*  企业项目Id  */
	PloyIds       *string   `json:"ployIds"`       /*  策略id  */
}
