package oceanfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// OceanfsOpendListSfsApi
/* 根据资源池 ID ，查询用户已开通文件系统
 */type OceanfsOpendListSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOceanfsOpendListSfsApi(client *core.CtyunClient) *OceanfsOpendListSfsApi {
	return &OceanfsOpendListSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oceanfs/opend-list-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *OceanfsOpendListSfsApi) Do(ctx context.Context, credential core.Credential, req *OceanfsOpendListSfsRequest) (*OceanfsOpendListSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
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
	var resp OceanfsOpendListSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type OceanfsOpendListSfsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  查询的页码。默认为1  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页的元素个数。默认为10  */
}

type OceanfsOpendListSfsResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}
