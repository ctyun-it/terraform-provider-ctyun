package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryRegionByAccountIdApi
/* 查询账户资源池 */
type CtiamQueryRegionByAccountIdApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryRegionByAccountIdApi(client *core.CtyunClient) *CtiamQueryRegionByAccountIdApi {
	return &CtiamQueryRegionByAccountIdApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/region/queryRegionByAccountId",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryRegionByAccountIdApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryRegionByAccountIdRequest) (*CtiamQueryRegionByAccountIdResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamQueryRegionByAccountIdRequest
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
	var resp CtiamQueryRegionByAccountIdResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryRegionByAccountIdRequest struct {
	ZoneName *string `json:"zoneName,omitempty"` /*  资源池名称  */
	ZoneId   *string `json:"zoneId,omitempty"`   /*  资源池ID  */
}

type CtiamQueryRegionByAccountIdResponse struct {
	StatusCode *string                                       `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                       `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	ReturnObj  *CtiamQueryRegionByAccountIdReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Error      *string                                       `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamQueryRegionByAccountIdReturnObjResponse struct {
	RegionList []*CtiamQueryRegionByAccountIdReturnObjRegionListResponse `json:"regionList"` /*  资源池列表  */
}

type CtiamQueryRegionByAccountIdReturnObjRegionListResponse struct {
	ZoneId   *string `json:"zoneId"`   /*  资源池id  */
	ZoneName *string `json:"zoneName"` /*  资源池名称  */
}
