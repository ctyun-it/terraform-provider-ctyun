package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamSetEpGroupPloyIncreaseApi
type CtiamSetEpGroupPloyIncreaseApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamSetEpGroupPloyIncreaseApi(client *core.CtyunClient) *CtiamSetEpGroupPloyIncreaseApi {
	return &CtiamSetEpGroupPloyIncreaseApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/project/setEpGroupPloyIncrease",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamSetEpGroupPloyIncreaseApi) Do(ctx context.Context, credential core.Credential, req *CtiamSetEpGroupPloyIncreaseRequest) (*CtiamSetEpGroupPloyIncreaseResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamSetEpGroupPloyIncreaseRequest
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
	var resp CtiamSetEpGroupPloyIncreaseResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamSetEpGroupPloyIncreaseRequest struct {
	Id        string `json:"id"`        /*  用户组id或用户id，根据value值不同区分  */
	ProjectId string `json:"projectId"` /*  企业项目id  */
	PloyIds   string `json:"ployIds"`   /*  策略id（策略ID,可传多个，以逗号分割）  */
	Value     string `json:"value"`     /*      GROUP("CUS_153_01_0001","用户组"),
	USER("CUS_153_01_0002","用户");  */
}

type CtiamSetEpGroupPloyIncreaseResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息
	 */
	Error     *string                                       `json:"error"`     /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj *CtiamSetEpGroupPloyIncreaseReturnObjResponse `json:"returnObj"` /*  返回体  */
}

type CtiamSetEpGroupPloyIncreaseReturnObjResponse struct {
	ProjectPrivilegeMessageBoList []*CtiamSetEpGroupPloyIncreaseReturnObjProjectPrivilegeMessageBoListResponse `json:"projectPrivilegeMessageBoList"` /*  授权信息列表  */
	EntProjectObjRelId            *string                                                                      `json:"entProjectObjRelId"`            /*  绑定关系id  */
}

type CtiamSetEpGroupPloyIncreaseReturnObjProjectPrivilegeMessageBoListResponse struct {
	PolyId      *string `json:"polyId"`      /*  策略id  */
	PrivilegeId *string `json:"privilegeId"` /*  授权id  */
}
