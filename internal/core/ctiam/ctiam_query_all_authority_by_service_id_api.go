package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtiamQueryAllAuthorityByServiceIdApi
type CtiamQueryAllAuthorityByServiceIdApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryAllAuthorityByServiceIdApi(client *core.CtyunClient) *CtiamQueryAllAuthorityByServiceIdApi {
	return &CtiamQueryAllAuthorityByServiceIdApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/service/queryAllAuthorityByServiceId",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryAllAuthorityByServiceIdApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryAllAuthorityByServiceIdRequest) (*CtiamQueryAllAuthorityByServiceIdResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("serviceId", strconv.FormatInt(int64(req.ServiceId), 10))
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamQueryAllAuthorityByServiceIdResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryAllAuthorityByServiceIdRequest struct {
	ServiceId int32 `json:"serviceId"` /*  云服务ID  */
}

type CtiamQueryAllAuthorityByServiceIdResponse struct {
	StatusCode *string                                             `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                             `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                             `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamQueryAllAuthorityByServiceIdReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamQueryAllAuthorityByServiceIdReturnObjResponse struct {
	AuthorityList []*CtiamQueryAllAuthorityByServiceIdReturnObjAuthorityListResponse `json:"authorityList"` /*  权限点列表  */
}

type CtiamQueryAllAuthorityByServiceIdReturnObjAuthorityListResponse struct {
	ServiceId      *int32  `json:"serviceId"`      /*  服务ID  */
	Name           *string `json:"name"`           /*  权限点名称  */
	Code           *string `json:"code"`           /*  权限点编码  */
	Description    *string `json:"description"`    /*  描述  */
	CtrntemplateId *string `json:"ctrntemplateId"` /*  资源路径模板ID  */
}
