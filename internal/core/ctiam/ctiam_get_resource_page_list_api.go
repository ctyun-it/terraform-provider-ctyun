package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamGetResourcePageListApi
/* 分页查询资源信息 */
type CtiamGetResourcePageListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamGetResourcePageListApi(client *core.CtyunClient) *CtiamGetResourcePageListApi {
	return &CtiamGetResourcePageListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/resource/getResourcePageList",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamGetResourcePageListApi) Do(ctx context.Context, credential core.Credential, req *CtiamGetResourcePageListRequest) (*CtiamGetResourcePageListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamGetResourcePageListRequest
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
	var resp CtiamGetResourcePageListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamGetResourcePageListRequest struct {
	ProjectSetId string `json:"projectSetId"` /*  企业项目id  */
	PageNum      int32  `json:"pageNum"`      /*  页码  */
	PageSize     int32  `json:"pageSize"`     /*  每页条数  */
}

type CtiamGetResourcePageListResponse struct {
	StatusCode *string                                    `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                    `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	ReturnObj  *CtiamGetResourcePageListReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Error      *string                                    `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamGetResourcePageListReturnObjResponse struct {
	PageNum    *int32  `json:"pageNum"`    /*  页码  */
	PageSize   *int32  `json:"pageSize"`   /*  每页条数  */
	Total      *int32  `json:"total"`      /*  总条数  */
	Pages      *int32  `json:"pages"`      /*  总页数  */
	RecordList *string `json:"recordList"` /*  资源列表  */
}
