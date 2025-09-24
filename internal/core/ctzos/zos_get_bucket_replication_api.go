package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// ZosGetBucketReplicationApi
/* 获取某个源桶已创建的所有跨域复制规则。
 */type ZosGetBucketReplicationApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetBucketReplicationApi(client *core.CtyunClient) *ZosGetBucketReplicationApi {
	return &ZosGetBucketReplicationApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-bucket-replication",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetBucketReplicationApi) Do(ctx context.Context, credential core.Credential, req *ZosGetBucketReplicationRequest) (*ZosGetBucketReplicationResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	if req.Page != 0 {
		ctReq.AddParam("page", strconv.FormatInt(int64(req.Page), 10))
	}
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
	var resp ZosGetBucketReplicationResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetBucketReplicationRequest struct {
	Bucket   string /*  桶名  */
	RegionID string /*  区域 ID  */
	Page     int64  /*  页码，默认值为 1  */
	PageNo   int64  /*  页码，若与参数 page 同时存在，以 pageNo 为准。默认值为 1  */
	PageSize int64  /*  页大小，取值范围 1~50，默认为10  */
}

type ZosGetBucketReplicationResponse struct {
	StatusCode  int64                                     `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中/失败)  */
	Message     string                                    `json:"message,omitempty"`     /*  状态描述  */
	Description string                                    `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosGetBucketReplicationReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                                    `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                    `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetBucketReplicationReturnObjResponse struct {
	Result       []*ZosGetBucketReplicationReturnObjResultResponse `json:"result"`                 /*  规则详情的数组  */
	CurrentCount int64                                             `json:"currentCount,omitempty"` /*  当前页数量  */
	TotalCount   int64                                             `json:"totalCount,omitempty"`   /*  总数  */
}

type ZosGetBucketReplicationReturnObjResultResponse struct {
	Fuid             string   `json:"fuid,omitempty"`             /*  同步规则ID  */
	TargetRegionID   string   `json:"targetRegionID,omitempty"`   /*  目标区域ID  */
	TargetRegionIDv2 string   `json:"targetRegionIDv2,omitempty"` /*  新版目标区域ID  */
	TargetRegionName string   `json:"targetRegionName,omitempty"` /*  目标区域名称  */
	TargetBucket     string   `json:"targetBucket,omitempty"`     /*  目标桶  */
	Prefixes         []string `json:"prefixes"`                   /*  桶前缀  */
	Plot             *bool    `json:"plot"`                       /*  同步策略, 同步时是否允许删除  */
	History          *bool    `json:"history"`                    /*  是否同步历史数据  */
	Progress         float32  `json:"progress"`                   /*  同步进度  */
	LastUpdate       string   `json:"lastUpdate,omitempty"`       /*  ISO-8601 格式的日期字符串。若为空字符串 ""，表示未同步过。  */
}
