package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetBucketStatisticsApi
/* 查询桶统计信息<br>1. 传bucket参数，表示查询该桶的统计信息<br>2. 不传bucket参数，表示查询所有桶的统计信息，即该用户所有桶的统计信息<br>统计项：<br>&emsp;1). 每个小时累计的请求次数<br>&emsp;2). 每个小时累计的公网流出流量<br>&emsp;3). 每个小时累计的数据取回流量，仅在归档、低频存储类型时有此指标，归档的取回指归档对象解冻，低频的取回指低频对象的下载<br>&emsp;4). 每个小时累计的数据取回次数，仅在归档、低频存储类型时有此指标，归档的取回指归档对象解冻，低频的取回指低频对象的下载<br>&emsp;5). 存储容量在每个整点的值<br>&emsp;6). 每个小时累计的跨域复制公网流出流量<br>计算方法：<br>&emsp;1). 请求次数、公网流出流量为读API加写API的总值：<br>&emsp;&emsp;读API：'list_buckets', 'list_bucket', 'stat_bucket', 'get_obj', 'get_acls', 'get_cors', 'list_bucket_multiparts', 'list_multipart'<br>&emsp;&emsp;写API：'create_bucket', 'delete_bucket', 'put_obj', 'put_acls', 'put_cors', 'delete_cors', 'delete_obj', 'init_multipart', 'put_obj', 'complete_multipart', 'abort_multipart'<br>&emsp;&emsp;其中，请求次数为公网请求数加内网请求数的总值。
 */type ZosGetBucketStatisticsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetBucketStatisticsApi(client *core.CtyunClient) *ZosGetBucketStatisticsApi {
	return &ZosGetBucketStatisticsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-bucket-statistics",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetBucketStatisticsApi) Do(ctx context.Context, credential core.Credential, req *ZosGetBucketStatisticsRequest) (*ZosGetBucketStatisticsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.Bucket != "" {
		ctReq.AddParam("bucket", req.Bucket)
	}
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("startTime", req.StartTime)
	ctReq.AddParam("endTime", req.EndTime)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetBucketStatisticsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetBucketStatisticsRequest struct {
	Bucket    string /*  存储桶名  */
	RegionID  string /*  区域ID  */
	StartTime string /*  日期-小时 格式的时间字符串，时区为 UTC 时区  */
	EndTime   string /*  日期-小时 格式的时间字符串，时区为 UTC 时区  */
}

type ZosGetBucketStatisticsResponse struct {
	StatusCode  int64                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                   `json:"message,omitempty"`     /*  状态描述  */
	Description string                                   `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosGetBucketStatisticsReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                   `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetBucketStatisticsReturnObjResponse struct{}
