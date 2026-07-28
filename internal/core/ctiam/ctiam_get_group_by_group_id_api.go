package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamGetGroupByGroupIdApi
/* 根据用户组ID查询用户组信息 */
type CtiamGetGroupByGroupIdApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamGetGroupByGroupIdApi(client *core.CtyunClient) *CtiamGetGroupByGroupIdApi {
	return &CtiamGetGroupByGroupIdApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/userGroup/getGroupByGroupId",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamGetGroupByGroupIdApi) Do(ctx context.Context, credential core.Credential, req *CtiamGetGroupByGroupIdRequest) (*CtiamGetGroupByGroupIdResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("groupId", req.GroupId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamGetGroupByGroupIdResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamGetGroupByGroupIdRequest struct {
	GroupId string `json:"groupId"` /*  用户组id  */
}

type CtiamGetGroupByGroupIdResponse struct {
	StatusCode *string                                  `json:"statusCode"` /*  状态码  */
	ReturnObj  *CtiamGetGroupByGroupIdReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                                  `json:"message"`    /*  对于调用失败的补充信息说明  */
	Error      *string                                  `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamGetGroupByGroupIdReturnObjResponse struct {
	Id         *string `json:"id"`         /*  用户组id  */
	HwGroupId  *string `json:"hwGroupId"`  /*  合营用户组id  */
	GroupName  *string `json:"groupName"`  /*  用户组名称  */
	AccountId  *string `json:"accountId"`  /*  账号id  */
	GroupIntro *string `json:"groupIntro"` /*  用户组描述  */
	CreateTime *string `json:"createTime"` /*  创建时间  */
}
