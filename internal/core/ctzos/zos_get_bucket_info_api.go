package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetBucketInfoApi
/* 查询桶的基础信息和用量数据。
 */type ZosGetBucketInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetBucketInfoApi(client *core.CtyunClient) *ZosGetBucketInfoApi {
	return &ZosGetBucketInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-bucket-info",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetBucketInfoApi) Do(ctx context.Context, credential core.Credential, req *ZosGetBucketInfoRequest) (*ZosGetBucketInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetBucketInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetBucketInfoRequest struct {
	Bucket   string /*  存储桶名  */
	RegionID string /*  区域ID  */
}

type ZosGetBucketInfoResponse struct {
	StatusCode  int64                              `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                             `json:"message,omitempty"`     /*  状态描述  */
	Description string                             `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosGetBucketInfoReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                             `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                             `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetBucketInfoReturnObjResponse struct {
	ProjectID         string                                              `json:"projectID,omitempty"`         /*  企业项目 ID  */
	CmkUUID           *string                                             `json:"cmkUUID,omitempty"`           /*  kms加密ID，若此值为null，则表示未开启加密  */
	StorageType       string                                              `json:"storageType,omitempty"`       /*  存储类型，可选的值为 STANDARD, STANDARD_IA, GLACIER  */
	AZPolicy          string                                              `json:"AZPolicy,omitempty"`          /*  AZ策略，single-az 或 multi-az  */
	BucketQuota       *ZosGetBucketInfoReturnObjBucketQuotaResponse       `json:"bucketQuota"`                 /*  桶配额  */
	Tenant            string                                              `json:"tenant,omitempty"`            /*  租户信息  */
	Ctime             string                                              `json:"ctime,omitempty"`             /*  桶的创建时间  */
	BucketPreviewFlag int64                                               `json:"bucketPreviewFlag,omitempty"` /*  桶是否支持文件预览  */
	PlacementRule     string                                              `json:"placementRule,omitempty"`     /*  placement规则名  */
	Mtime             string                                              `json:"mtime,omitempty"`             /*  桶的最后修改时间  */
	Bucket            string                                              `json:"bucket,omitempty"`            /*  桶名  */
	Owner             string                                              `json:"owner,omitempty"`             /*  所属用户uid  */
	Usage             *ZosGetBucketInfoReturnObjUsageResponse             `json:"usage"`                       /*  使用信息  */
	NumShards         int64                                               `json:"numShards,omitempty"`         /*  分片数量  */
	TagMap            *ZosGetBucketInfoReturnObjTagMapResponse            `json:"tagMap"`                      /*  键值对形式的桶标签集  */
	IndexType         string                                              `json:"indexType,omitempty"`         /*  索引类型  */
	ExplicitPlacement *ZosGetBucketInfoReturnObjExplicitPlacementResponse `json:"explicitPlacement"`           /*  显示设置placement  */
	Zonegroup         string                                              `json:"zonegroup,omitempty"`         /*  zone组  */
}

type ZosGetBucketInfoReturnObjBucketQuotaResponse struct {
	Enabled    *bool `json:"enabled"`              /*  是否开启  */
	MaxSize    int64 `json:"maxSize,omitempty"`    /*  最大容量，单位byte，未设置时默认为-1  */
	MaxObjects int64 `json:"maxObjects,omitempty"` /*  最大对象数，未设置时默认为-1  */
	CheckOnRaw *bool `json:"checkOnRaw"`           /*  是否使用原始对象大小进行配额检查  */
	MaxSizeKb  int64 `json:"maxSizeKb,omitempty"`  /*  最大容量，单位kb  */
}

type ZosGetBucketInfoReturnObjUsageResponse struct {
	SizeKbUtilized      int64                                                      `json:"sizeKbUtilized,omitempty"`      /*  已使用容量，单位kb  */
	SizeActual          int64                                                      `json:"sizeActual,omitempty"`          /*  实际使用容量，单位byte  */
	SizeKbActual        int64                                                      `json:"sizeKbActual,omitempty"`        /*  实际使用容量，单位kb  */
	SizeKb              int64                                                      `json:"sizeKb,omitempty"`              /*  容量，单位kb  */
	StorageTypeIa       *ZosGetBucketInfoReturnObjUsageStorageTypeIaResponse       `json:"storageTypeIa"`                 /*  低频用量  */
	NumObjects          int64                                                      `json:"numObjects,omitempty"`          /*  对象数量  */
	EarlydelGlacierSize int64                                                      `json:"earlydelGlacierSize,omitempty"` /*  提前删除归档类型数据量  */
	StorageTypeGlacier  *ZosGetBucketInfoReturnObjUsageStorageTypeGlacierResponse  `json:"storageTypeGlacier"`            /*  归档用量  */
	StorageTypeStandard *ZosGetBucketInfoReturnObjUsageStorageTypeStandardResponse `json:"storageTypeStandard"`           /*  标准用量  */
	NumMultiparts       int64                                                      `json:"numMultiparts,omitempty"`       /*  碎片数量  */
	EarlydelIaSize      int64                                                      `json:"earlydelIaSize,omitempty"`      /*  提前删除低频类型数据量  */
	SizeUtilized        int64                                                      `json:"sizeUtilized,omitempty"`        /*  已使用容量，单位byte  */
	Size                int64                                                      `json:"size,omitempty"`                /*  容量，单位byte  */
}

type ZosGetBucketInfoReturnObjTagMapResponse struct{}

type ZosGetBucketInfoReturnObjExplicitPlacementResponse struct {
	DataExtraPool string `json:"dataExtraPool,omitempty"` /*  数据冗余池  */
	DataPool      string `json:"dataPool,omitempty"`      /*  数据池  */
	IndexPool     string `json:"indexPool,omitempty"`     /*  索引池  */
}

type ZosGetBucketInfoReturnObjUsageStorageTypeIaResponse struct {
	SizeKbUtilized int64 `json:"sizeKbUtilized,omitempty"` /*  已使用容量，单位kb  */
	SizeActual     int64 `json:"sizeActual,omitempty"`     /*  实际使用容量，单位byte  */
	SizeKbActual   int64 `json:"sizeKbActual,omitempty"`   /*  实际使用容量，单位kb  */
	NumObjects     int64 `json:"numObjects,omitempty"`     /*  对象数量  */
	SizeUtillized  int64 `json:"sizeUtillized,omitempty"`  /*  已使用容量，单位byte  */
	NumMultiparts  int64 `json:"numMultiparts,omitempty"`  /*  碎片数量  */
	SizeKb         int64 `json:"sizeKb,omitempty"`         /*  容量，单位kb  */
	Size           int64 `json:"size,omitempty"`           /*  容量，单位byte  */
}

type ZosGetBucketInfoReturnObjUsageStorageTypeGlacierResponse struct {
	SizeKbUtilized int64 `json:"sizeKbUtilized,omitempty"` /*  已使用容量，单位kb  */
	SizeActual     int64 `json:"sizeActual,omitempty"`     /*  实际使用容量，单位byte  */
	SizeKbActual   int64 `json:"sizeKbActual,omitempty"`   /*  实际使用容量，单位kb  */
	NumObjects     int64 `json:"numObjects,omitempty"`     /*  对象数量  */
	SizeUtillized  int64 `json:"sizeUtillized,omitempty"`  /*  已使用容量，单位byte  */
	NumMultiparts  int64 `json:"numMultiparts,omitempty"`  /*  碎片数量  */
	SizeKb         int64 `json:"sizeKb,omitempty"`         /*  容量，单位kb  */
	Size           int64 `json:"size,omitempty"`           /*  容量，单位byte  */
}

type ZosGetBucketInfoReturnObjUsageStorageTypeStandardResponse struct {
	SizeKbUtilized int64 `json:"sizeKbUtilized,omitempty"` /*  已使用容量，单位kb  */
	SizeActual     int64 `json:"sizeActual,omitempty"`     /*  实际使用容量，单位byte  */
	SizeKbActual   int64 `json:"sizeKbActual,omitempty"`   /*  实际使用容量，单位kb  */
	NumObjects     int64 `json:"numObjects,omitempty"`     /*  对象数量  */
	SizeUtillized  int64 `json:"sizeUtillized,omitempty"`  /*  已使用容量，单位byte  */
	NumMultiparts  int64 `json:"numMultiparts,omitempty"`  /*  碎片数量  */
	SizeKb         int64 `json:"sizeKb,omitempty"`         /*  容量，单位kb  */
	Size           int64 `json:"size,omitempty"`           /*  容量，单位byte  */
}
