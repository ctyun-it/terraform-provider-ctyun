package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// ZosListBucketsApi
/* 查询某个地域下所有桶的信息。
 */type ZosListBucketsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosListBucketsApi(client *core.CtyunClient) *ZosListBucketsApi {
	return &ZosListBucketsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/list-buckets",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosListBucketsApi) Do(ctx context.Context, credential core.Credential, req *ZosListBucketsRequest) (*ZosListBucketsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ProjectID != "" {
		ctReq.AddParam("projectID", req.ProjectID)
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
	var resp ZosListBucketsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosListBucketsRequest struct {
	RegionID  string /*  区域 ID，将返回对应区域下的桶。若需要返回所有公共资源池的桶，请传入 public。  */
	ProjectID string /*  桶所属的企业项目。如有多个，使用逗号分隔；若不传此参数，将使用当前账号所有有权限的企业项目  */
	PageSize  int64  /*  页大小，默认值 10，取值范围 1~50  */
	PageNo    int64  /*  页码，默认值 1  */
}

type ZosListBucketsResponse struct {
	StatusCode  int64                            `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                           `json:"message,omitempty"`     /*  状态描述  */
	Description string                           `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosListBucketsReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                           `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                           `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosListBucketsReturnObjResponse struct {
	BucketList   []*ZosListBucketsReturnObjBucketListResponse `json:"bucketList"`             /*  桶数组  */
	BucketTotal  int64                                        `json:"bucketTotal,omitempty"`  /*  总记录数  */
	PageSize     int64                                        `json:"pageSize,omitempty"`     /*  页大小  */
	PageNo       int64                                        `json:"pageNo,omitempty"`       /*  页码  */
	TotalCount   int64                                        `json:"totalCount,omitempty"`   /*  总记录数  */
	CurrentCount int64                                        `json:"currentCount,omitempty"` /*  当前页记录数  */
}

type ZosListBucketsReturnObjBucketListResponse struct {
	CreationDate string `json:"creationDate,omitempty"` /*  创建日期，为 ISO 8601 格式  */
	CmkUUID      string `json:"cmkUUID,omitempty"`      /*  加密ID，若 isEncrypted 为 false，此值为空字符串 ""  */
	StorageType  string `json:"storageType,omitempty"`  /*  存储类型，可选的值为 STANDARD, STANDARD_IA, GLACIER  */
	ProjectID    string `json:"projectID,omitempty"`    /*  企业项目 ID  */
	Bucket       string `json:"bucket,omitempty"`       /*  桶名  */
	IsEncrypted  *bool  `json:"isEncrypted"`            /*  是否加密  */
	AZPolicy     string `json:"AZPolicy,omitempty"`     /*  AZ策略，single-az 或 multi-az  */
	RegionName   string `json:"regionName,omitempty"`   /*  区域名称  */
	RegionID     string `json:"regionID,omitempty"`     /*  区域ID  */
}
