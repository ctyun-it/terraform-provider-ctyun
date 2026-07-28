package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanSdwanListSNATApi
/* 查询SNAT列表信息 */
type SdwanSdwanListSNATApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanListSNATApi(client *core.CtyunClient) *SdwanSdwanListSNATApi {
	return &SdwanSdwanListSNATApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/snat/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanListSNATApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanListSNATRequest) (*SdwanSdwanListSNATResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("siteID", req.SiteID)
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanSdwanListSNATResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanListSNATRequest struct {
	SiteID   string `json:"siteID"`   /*  站点 id  */
	PageNo   int32  `json:"pageNo"`   /*  页码  */
	PageSize int32  `json:"pageSize"` /*  每页记录数目  */
}

type SdwanSdwanListSNATResponse struct {
	StatusCode  int32                                `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                              `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                              `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                              `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanListSNATReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                              `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanSdwanListSNATReturnObjResponse struct {
	Result       []*SdwanSdwanListSNATReturnObjResultResponse `json:"result"`       /*  查询snat 信息  */
	TotalCount   int32                                        `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                        `json:"currentCount"` /*  当前页总数  */
}

type SdwanSdwanListSNATReturnObjResultResponse struct {
	SnatID  *string   `json:"snatID"`  /*  业务ID  */
	SiteID  *string   `json:"siteID"`  /*  站点ID  */
	SnatIP  *string   `json:"snatIP"`  /*  对外服务IP  */
	Subnets []*string `json:"subnets"` /*  私网子网，值类型为string  */
}
