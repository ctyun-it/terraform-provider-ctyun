package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// ZosListPartsApi
/* 列出上传对象的全部分段。
 */type ZosListPartsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosListPartsApi(client *core.CtyunClient) *ZosListPartsApi {
	return &ZosListPartsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/list-parts",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosListPartsApi) Do(ctx context.Context, credential core.Credential, req *ZosListPartsRequest) (*ZosListPartsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("key", req.Key)
	if req.MaxParts != 0 {
		ctReq.AddParam("maxParts", strconv.FormatInt(int64(req.MaxParts), 10))
	}
	if req.PartNumberMarker != 0 {
		ctReq.AddParam("partNumberMarker", strconv.FormatInt(int64(req.PartNumberMarker), 10))
	}
	ctReq.AddParam("uploadID", req.UploadID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosListPartsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosListPartsRequest struct {
	RegionID         string /*  区域 ID  */
	Bucket           string /*  桶名  */
	Key              string /*  对象名  */
	MaxParts         int64  /*  返回的最大分块数，默认值与最大值均为1000  */
	PartNumberMarker int64  /*  指定列表应该开始的位置。只会列出比这个编号更高的分块。  */
	UploadID         string /*  uploadID  */
}

type ZosListPartsResponse struct {
	StatusCode  int64                          `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                         `json:"message,omitempty"`     /*  状态描述  */
	Description string                         `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosListPartsReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                         `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                         `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosListPartsReturnObjResponse struct {
	Bucket               string                                `json:"bucket,omitempty"`               /*  文件上传的Bcuket名称  */
	NextPartNumberMarker int64                                 `json:"nextPartNumberMarker,omitempty"` /*  下一次list的时候的分段起始编号，主要用于截断返回时(也就是已上传的分段数目大于当前返回的分段数目)，作为下一次list的分段起始编号  */
	Parts                []*ZosListPartsReturnObjPartsResponse `json:"parts"`                          /*  已经上传的分段信息  */
	UploadID             string                                `json:"uploadID,omitempty"`             /*  分段上传的ID  */
	StorageClass         string                                `json:"storageClass,omitempty"`         /*  分段上传的文件对应的存储级别  */
	Key                  string                                `json:"key,omitempty"`                  /*  分段上传的文件在集群中保存的名字  */
	Owner                *ZosListPartsReturnObjOwnerResponse   `json:"owner"`                          /*  分段上传的文件所属用户  */
	MaxParts             int64                                 `json:"maxParts,omitempty"`             /*  当前list最多返回的分段数目  */
	IsTruncated          *bool                                 `json:"isTruncated"`                    /*  表示返回的分段列表是否被截断。值为 true 表示列表已被截断。  */
	PartNumberMarker     int64                                 `json:"partNumberMarker,omitempty"`     /*  当前list的分段起始编号  */
}

type ZosListPartsReturnObjPartsResponse struct {
	PartNumber   int64  `json:"partNumber,omitempty"`   /*  分段编号  */
	Size         int64  `json:"size,omitempty"`         /*  分段大小（单位 Bytes）  */
	Etag         string `json:"Etag,omitempty"`         /*  该分段数据对应Etag  */
	LastModified string `json:"lastModified,omitempty"` /*  该分段上次被修改的时间，ISO-8601 格式的日期字符串  */
}

type ZosListPartsReturnObjOwnerResponse struct {
	DisplayName string `json:"displayName,omitempty"` /*  展示名称  */
	ID          string `json:"ID,omitempty"`          /*  用户名  */
}
