package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// ZosListAllPartsApi
/* 列出桶内所有分段上传的碎片。
 */type ZosListAllPartsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosListAllPartsApi(client *core.CtyunClient) *ZosListAllPartsApi {
	return &ZosListAllPartsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/list-all-parts",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosListAllPartsApi) Do(ctx context.Context, credential core.Credential, req *ZosListAllPartsRequest) (*ZosListAllPartsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	if req.Page != 0 {
		ctReq.AddParam("page", strconv.FormatInt(int64(req.Page), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosListAllPartsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosListAllPartsRequest struct {
	Bucket   string /*  桶名  */
	RegionID string /*  区域 ID  */
	Page     int64  /*  页码。默认1  */
	PageSize int64  /*  每页展示的最大分段数量。取值范围 1~50，默认值为 10  */
	PageNo   int64  /*  页码，若与参数 page 同时存在，以 pageNo 为准。默认值 1  */
}

type ZosListAllPartsResponse struct {
	StatusCode  int64                             `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                            `json:"message,omitempty"`     /*  状态描述  */
	Description string                            `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosListAllPartsReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                            `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                            `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosListAllPartsReturnObjResponse struct {
	Result       []*ZosListAllPartsReturnObjResultResponse `json:"result"`                 /*  规则详情的数组  */
	TotalCount   int64                                     `json:"totalCount,omitempty"`   /*  总数  */
	CurrentCount int64                                     `json:"currentCount,omitempty"` /*  当前页记录数  */
}

type ZosListAllPartsReturnObjResultResponse struct {
	Bucket       string                                            `json:"bucket,omitempty"`       /*  桶名  */
	UploadID     string                                            `json:"uploadID,omitempty"`     /*  分段上传ID  */
	Key          string                                            `json:"key,omitempty"`          /*  对象名  */
	FragmentSize int64                                             `json:"fragmentSize,omitempty"` /*  当前段上传内碎片大小  */
	PartInfo     []*ZosListAllPartsReturnObjResultPartInfoResponse `json:"partInfo"`               /*  当前分段上传内的碎片列表  */
}

type ZosListAllPartsReturnObjResultPartInfoResponse struct {
	LastModified string `json:"lastModified,omitempty"` /*  最后修改时间，为 ISO 8601 格式  */
	PartNumber   int64  `json:"partNumber,omitempty"`   /*  分段编号  */
	ETag         string `json:"ETag,omitempty"`         /*  ETag  */
	Size         int64  `json:"size,omitempty"`         /*  分段大小（单位 Bytes）  */
}
