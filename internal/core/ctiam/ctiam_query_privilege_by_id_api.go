package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryPrivilegeByIdApi
type CtiamQueryPrivilegeByIdApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryPrivilegeByIdApi(client *core.CtyunClient) *CtiamQueryPrivilegeByIdApi {
	return &CtiamQueryPrivilegeByIdApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/perm/queryPrivilegeById",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryPrivilegeByIdApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryPrivilegeByIdRequest) (*CtiamQueryPrivilegeByIdResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("privilegeId", req.PrivilegeId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamQueryPrivilegeByIdResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryPrivilegeByIdRequest struct {
	PrivilegeId string `json:"privilegeId"` /*  授权策略id  */
}

type CtiamQueryPrivilegeByIdResponse struct {
	StatusCode *string                                   `json:"statusCode"` /*  响应码  */
	ReturnObj  *CtiamQueryPrivilegeByIdReturnObjResponse `json:"returnObj"`  /*  返回信息  */
	Message    *string                                   `json:"message"`    /*  对于调用失败的补充信息说明  */
	Error      *string                                   `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamQueryPrivilegeByIdReturnObjResponse struct {
	PrivilegeId   *string `json:"privilegeId"`   /*  授权id  */
	RegionId      *string `json:"regionId"`      /*  资源池id  */
	Id            *string `json:"id"`            /*  "principalType"为1时是id是用户id，为2是委托用户id，为0是用户组  */
	AccountId     *string `json:"accountId"`     /*  账号id  */
	PolicyId      *string `json:"policyId"`      /*  策略id  */
	PrincipalType *string `json:"principalType"` /*  授权策略类型  */
}
